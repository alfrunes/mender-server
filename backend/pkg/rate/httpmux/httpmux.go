package httpmux

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mendersoftware/mender-server/pkg/rate"
	"github.com/mendersoftware/mender-server/pkg/requestid"
	"github.com/mendersoftware/mender-server/pkg/rest.utils"
)

var ErrTooManyRequests = errors.New("too many requests")

type RateMux struct {
	mux *http.ServeMux
}

type handle struct {
	*RateMux
	rate.Limiter
}

// ServeHTTP implements a basic http.Handler so that handler can be used
// as a handler for the mux. It will only write on errors and is expected
// to continue to the actual handler on success.
func (h *RateMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lim := h.Limiter(r)
	if lim == nil {
		return
	}
	res, err := lim.Reserve(ctx)
	if err != nil {
		hdr := w.Header()
		hdr.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(rest.Error{
			Err:       "internal server error",
			RequestID: requestid.FromContext(ctx),
		})
		return
	}
	if res.OK() {
		return
	} else {
		hdr := w.Header()
		hdr.Set("Content-Type", "application/json")
		retryAfter := math.Ceil(res.Delay().Abs().Seconds())
		hdr.Set("Retry-After", fmt.Sprint(retryAfter))
		w.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(w).Encode(rest.Error{
			Err:       ErrTooManyRequests.Error(),
			RequestID: requestid.FromContext(ctx),
		})
	}
}

func (h *RateMux) MiddlewareGin(c *gin.Context) {
	ctx := c.Request.Context()
	lim := h.Limiter(c.Request)
	if lim == nil {
		c.Next()
	}
	res, err := lim.Reserve(ctx)
	if err != nil {
		rest.RenderInternalError(c, err)
		c.Abort()
	}
	if res.OK() {
		c.Next()
	} else {
		retryAfter := math.Ceil(res.Delay().Abs().Seconds())
		c.Header("Retry-After", fmt.Sprint(retryAfter))
		rest.RenderError(c, http.StatusBadRequest, ErrTooManyRequests)
	}
}

func NewRateMux() *RateMux {
	return &RateMux{
		mux: http.NewServeMux(),
	}
}

func (m *RateMux) AddPattern(pattern string, limiter rate.Limiter) {
	m.mux.Handle(pattern, handle{Limiter: limiter, RateMux: m})
}

func (m *RateMux) Limiter(r *http.Request) rate.Limiter {
	h, _ := m.mux.Handler(r)
	hh, ok := h.(handle)
	if ok {
		return hh.Limiter
	}
	return nil
}
