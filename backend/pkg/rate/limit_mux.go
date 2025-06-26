package rate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"

	"github.com/mendersoftware/mender-server/pkg/identity"
	"github.com/mendersoftware/mender-server/pkg/requestid"
	"github.com/mendersoftware/mender-server/pkg/rest.utils"
)

var ErrTooManyRequests = errors.New("too many requests")

type LimiterGroup struct {
	EventLimiter
	*template.Template
}

type HTTPLimiter struct {
	templateData  func(r *http.Request) any
	template      *template.Template
	mux           *http.ServeMux
	DefaultGroup  *LimiterGroup
	limiterGroups map[string]*LimiterGroup
}

func defaultTemplateData(r *http.Request) any {
	id := identity.FromContext(r.Context())
	ctx := map[string]any{
		"Identity": id,
		"Request":  r,
	}
	return ctx
}

func NewHTTPLimiter(eventLimiter EventLimiter, eventTemplate string) (*HTTPLimiter, error) {
	template, err := template.New("").Parse(eventTemplate)
	if err != nil {
		return nil, fmt.Errorf("invalid eventTemplate: %w", err)
	}
	return &HTTPLimiter{
		template:     template.New("").Option("missingkey=zero"),
		mux:          http.NewServeMux(),
		templateData: defaultTemplateData,
		DefaultGroup: &LimiterGroup{
			EventLimiter: eventLimiter,
			Template:     template,
		},
		limiterGroups: make(map[string]*LimiterGroup),
	}, nil
}

func (h *HTTPLimiter) TemplateDataFunc(f func(*http.Request) any) *HTTPLimiter {
	h.templateData = f
	return h
}

func (h *HTTPLimiter) TemplateFuncs(funcs map[string]any) *HTTPLimiter {
	h.template.Funcs(funcs)
	return h
}

func (h *HTTPLimiter) AddRateLimitGroup(limiter EventLimiter, group, eventTemplate string) error {
	t, err := h.template.Clone()
	if err == nil {
		_, err = t.Parse(eventTemplate)
	}
	if err != nil {
		return fmt.Errorf("failed to compile event template: %w", err)
	}
	h.limiterGroups[group] = &LimiterGroup{
		EventLimiter: limiter,
		Template:     t,
	}
	return nil
}

func (h *HTTPLimiter) MatchHTTPPattern(
	pattern, groupTemplate string,
) error {
	var (
		t   *template.Template
		err error
	)
	if groupTemplate != "" {
		// Compile eventTemplate:
		t, err = h.template.Clone()
		if err == nil {
			_, err = t.Parse(groupTemplate)
		}
		if err != nil {
			return fmt.Errorf("error parsing group_template: %w", err)
		}
	}
	limiterHandle := handle{
		HTTPLimiter:   h,
		groupTemplate: t,
	}
	h.mux.Handle(pattern, limiterHandle)
	return nil
}

type handle struct {
	*HTTPLimiter
	groupTemplate *template.Template
}

// ServeHTTP implements a basic http.Handler so that handler can be used
// as a handler for the mux. It will only write on errors and is expected
// to continue to the actual handler on success.
func (h *HTTPLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := h.Reserve(r)
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
	if res == nil || res.OK() {
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

func (h *HTTPLimiter) MiddlewareGin(c *gin.Context) {
	res, err := h.Reserve(c.Request)
	if err != nil {
		rest.RenderInternalError(c, err)
		c.Abort()
	}
	if res == nil || res.OK() {
		c.Next()
	} else {
		retryAfter := math.Ceil(res.Delay().Abs().Seconds())
		c.Header("Retry-After", fmt.Sprint(retryAfter))
		rest.RenderError(c, http.StatusBadRequest, ErrTooManyRequests)
	}
}

func (m *HTTPLimiter) Reserve(r *http.Request) (Reservation, error) {
	h, _ := m.mux.Handler(r)
	hh, ok := h.(handle)
	if ok {
		ctx := r.Context()
		var b bytes.Buffer
		templateData := m.templateData(r)
		var limiter *LimiterGroup
		if hh.groupTemplate != nil {
			err := hh.groupTemplate.Execute(&b, templateData)
			if err != nil {
				return nil, fmt.Errorf("error executing ratelimit group template: %w", err)
			}
			limiter = m.limiterGroups[b.String()]
		}
		if limiter == nil {
			limiter = m.DefaultGroup
		}
		err := limiter.Execute(&b, templateData)
		if err != nil {
			return nil, fmt.Errorf("error executing template for event ID: %w", err)
		}
		if limiter == nil {
			return nil, nil
		}
		return limiter.ReserveEvent(ctx, b.String())
	}
	return nil, nil
}
