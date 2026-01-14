package stream

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// [prefix.]streamv1.<localAddr>.<remoteAddr>.<messageType>[:sequenceNum]
//
// messageType = enum{
//   data
//   ack
//   hello
//   bye
// }
//
// Subscriptions:
//   listener: [prefix.]streamv1.l.<localAddr>.*.hello
//   session:  [prefix.]streamv1.s.<localAddr>.<remoteAddr>.*

type natsStream struct {
	nc         *nats.Conn
	subSession *nats.Subscription
	msgChan    chan *nats.Msg

	retryTimer *time.Timer

	closed   chan struct{}
	recvChan chan *nats.Msg
	ackChan  chan *nats.Msg
	sendOnce chan struct{}
	recvOnce chan struct{}

	localAddr  string
	remoteAddr string

	remoteSeq uint32
	localSeq  uint32

	closeFunc func()
	termOnce  sync.Once
}

func newStream(nc *nats.Conn, localAddr, remoteAddr string) (*natsStream, error) {
	msgChan := make(chan *nats.Msg, 8)
	sub, err := nc.ChanSubscribe(fmtSessionSubject(localAddr, remoteAddr, "*"), msgChan)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to session: %w", err)
	}
	ret := &natsStream{
		nc:         nc,
		subSession: sub,
		msgChan:    msgChan,

		retryTimer: time.NewTimer(0),

		localAddr:  localAddr,
		remoteAddr: remoteAddr,
		recvChan:   make(chan *nats.Msg, 4),
		ackChan:    make(chan *nats.Msg, 4),
		closed:     make(chan struct{}),
		sendOnce:   make(chan struct{}, 1),
		recvOnce:   make(chan struct{}, 1),

		remoteSeq: 1,
		localSeq:  0,
	}
	go ret.protocol()
	return ret, nil
}

func fmtHelloSubject(localAddr, remoteAddr string) string {
	return fmt.Sprintf("streamv1.l.%s.%s.hello", localAddr, remoteAddr)
}

func fmtSessionSubject(localAddr, remoteAddr, msgType string) string {
	return fmt.Sprintf("streamv1.s.%s.%s.%s", localAddr, remoteAddr, msgType)
}

func cutLast(s, substr string) (before, after string, found bool) {
	if i := strings.LastIndex(s, substr); i >= 0 {
		return s[:i], s[i+len(substr):], true
	}
	return s, "", false
}

var (
	ErrProtocol          = errors.New("protocol error")
	ErrClosed            = errors.New("stream closed")
	ErrConnectionRefused = errors.New("connection refused")
)

func ConnectNATS(ctx context.Context, nc *nats.Conn, srcAddr, dstAddr string) (Conn, error) {
	stream, err := newStream(nc, srcAddr, dstAddr)
	defer func() {
		if err != nil {
			stream.subSession.Unsubscribe()
		}
	}()

	err = nc.PublishRequest(
		fmtHelloSubject(dstAddr, srcAddr),
		fmtSessionSubject(srcAddr, dstAddr, "ack:0"),
		nil)
	if err != nil {
		err = fmt.Errorf("failed to submit handshake to peer: %w", err)
		return nil, err
	}
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return nil, err
	case msg := <-stream.ackChan:
		if msg.Header.Get("Status") == "503" {
			return nil, ErrConnectionRefused
		}
	}

	return stream, nil
}

func (stream *natsStream) protocol() {
	for {
		select {
		case <-stream.closed:
			return
		case msg := <-stream.msgChan:
			_, msgType, _ := cutLast(msg.Subject, ".")
			switch {
			case msgType == "bye":
				if msg.Reply != "" {
					msg.Respond(nil)
				}
				stream.term()
			case msgType == "byeack":
				stream.term()

			case strings.HasPrefix(msgType, "data"):
				stream.recvChan <- msg

			case strings.HasPrefix(msgType, "ack"):
				stream.ackChan <- msg
			}
		}
	}
}

func (stream *natsStream) term() error {
	err := ErrClosed
	stream.termOnce.Do(func() {
		close(stream.closed)
		if stream.closeFunc != nil {
			stream.closeFunc()
		}
		err = stream.subSession.Unsubscribe()
	})
	return err
}

func (stream *natsStream) Send(ctx context.Context, data []byte) error {
	stream.sendOnce <- struct{}{}
	defer func() { <-stream.sendOnce }()
	stream.localSeq++
	msgTypeData := fmt.Sprintf("data:%d", stream.localSeq)
	msgTypeAck := fmt.Sprintf("ack:%d", stream.localSeq)
	err := stream.nc.PublishRequest(
		fmtSessionSubject(stream.remoteAddr, stream.localAddr, msgTypeData),
		fmtSessionSubject(stream.localAddr, stream.remoteAddr, msgTypeAck),
		data,
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	if !stream.retryTimer.Stop() {
		<-stream.retryTimer.C
	}
	var retryDelay = time.Millisecond * 200
	stream.retryTimer.Reset(retryDelay)
	for {
		select {
		case <-stream.closed:
			return ErrClosed
		case msg := <-stream.ackChan:
			_, seqNumStr, _ := cutLast(msg.Subject, ":")
			seqNum, _ := strconv.ParseUint(seqNumStr, 10, 32)
			// NOTE: "modulo arithmetic"
			diff := stream.localSeq - uint32(seqNum)
			if diff == 0 {
				return nil
			} else if diff == 1 {
				// ignore ack for last package
				continue
			}
			errClose := stream.Close(ctx)
			if errClose != nil {
				return fmt.Errorf("error closing stream on protocol violation: %w", err)
			}
			return ErrProtocol

		case <-stream.retryTimer.C:
			err := stream.nc.PublishRequest(
				fmtSessionSubject(stream.remoteAddr, stream.localAddr, msgTypeData),
				fmtSessionSubject(stream.localAddr, stream.remoteAddr, msgTypeAck),
				data,
			)
			if err != nil {
				return fmt.Errorf("failed to publish message: %w", err)
			}
			retryDelay *= 2
			stream.retryTimer.Reset(retryDelay)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (stream *natsStream) Recv(ctx context.Context) ([]byte, error) {
	stream.recvOnce <- struct{}{}
	defer func() { <-stream.recvOnce }()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case msg := <-stream.recvChan:
			err := msg.Respond(nil)
			if err != nil {
				return nil, fmt.Errorf("failed to ack message: %w", err)
			}
			_, seqNumStr, _ := cutLast(msg.Subject, ":")
			seqNum, _ := strconv.ParseUint(seqNumStr, 10, 32)
			diff := stream.remoteSeq - uint32(seqNum)
			if diff == 0 {
				stream.remoteSeq++
				return msg.Data, nil
			} else if diff == 1 {
				// ignore ack for last package
				continue
			}
			errClose := stream.Close(ctx)
			if errClose != nil {
				return nil, fmt.Errorf("error closing stream on protocol violation: %w", errClose)
			}
			return nil, ErrProtocol
		}
	}
}

func (stream *natsStream) Close(ctx context.Context) error {
	err := stream.nc.PublishRequest(
		fmtSessionSubject(stream.remoteAddr, stream.localAddr, "bye"),
		fmtSessionSubject(stream.localAddr, stream.remoteAddr, "byeack"),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to send close to peer: %w", err)
	}
	select {
	case <-stream.closed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (stream *natsStream) RemoteAddr() string {
	return stream.remoteAddr
}

func (stream *natsStream) LocalAddr() string {
	return stream.localAddr
}

type natsListener struct {
	nc        *nats.Conn
	subListen *nats.Subscription

	msgChan chan *nats.Msg
	addr    string

	// openConns track open connections
	openConns map[string]Conn
	mu        sync.Mutex
}

func ListenNATS(nc *nats.Conn, addr string) (Listener, error) {
	msgChan := make(chan *nats.Msg, 3)
	sub, err := nc.ChanSubscribe(fmtHelloSubject(addr, "*"), msgChan)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to listener subject: %w", err)
	}

	return &natsListener{
		nc:        nc,
		subListen: sub,
		addr:      addr,
		msgChan:   msgChan,
		openConns: make(map[string]Conn),
	}, nil
}

func (l *natsListener) Close(ctx context.Context) error {
	err := l.subListen.Unsubscribe()
	if err != nil {
		return err
	}
	for _, stream := range l.openConns {
		err = stream.Close(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *natsListener) Accept(ctx context.Context) (Conn, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg := <-l.msgChan:
		rest, msgType, _ := cutLast(msg.Subject, ".")
		if msgType != "hello" {
			return nil, ErrProtocol
		}
		_, addr, _ := cutLast(rest, ".")

		_, ok := l.openConns[addr]
		if ok {
			return nil, ErrProtocol
		}
		stream, err := newStream(l.nc, l.addr, addr)
		if err != nil {
			return nil, err
		}
		l.mu.Lock()
		if _, ok := l.openConns[addr]; ok {
			l.mu.Unlock()
			_ = stream.term()
			return nil, ErrProtocol
		} else {
			l.openConns[addr] = stream
		}
		l.mu.Unlock()
		stream.closeFunc = func() {
			l.mu.Lock()
			defer l.mu.Unlock()
			delete(l.openConns, addr)
		}
		err = msg.Respond(nil)
		if err != nil {
			l.openConns[addr].Close(ctx)
			return nil, fmt.Errorf("failed to complete handshake: %w", err)
		}
		return stream, nil
	}
}

func (l *natsListener) LocalAddr() string {
	return l.addr
}
