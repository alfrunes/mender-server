package stream

import (
	"context"
)

type Conn interface {
	Send(ctx context.Context, data []byte) error
	Recv(ctx context.Context) ([]byte, error)
	Close(ctx context.Context) error

	LocalAddr() string
	RemoteAddr() string
}

type Listener interface {
	Accept(ctx context.Context) (Conn, error)
	Close(ctx context.Context) error
	LocalAddr() string
}
