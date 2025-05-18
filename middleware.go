package draken

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/xid"
)

type ContextKey string

const ContextKeyRequestId ContextKey = "draken-request-id"

func (d *Draken) Middleware(middleware func(http.Handler) http.Handler) {
	d.Chi.Use(middleware)
}

func (d *Draken) EssentialMiddlewares() {
	d.Middleware(WebserverMiddleware())
	d.Middleware(RequestIdMiddleware())
	d.Middleware(middleware.RealIP)
	d.Middleware(middleware.Logger)
	d.Middleware(middleware.Recoverer)
	d.Middleware(middleware.Heartbeat("/health"))
	d.Middleware(SecurityMiddleware())
}

type SecurityMiddlewareConfig struct {
	XContentTypeOptions     string
	XFrameOptions           string
	XXSSProtection          string
	ReferrerPolicy          string
	ContentSecurityPolicy   string
	CacheControl            string
	StrictTransportSecurity string
}

// DefaultSecurityMiddlewareConfig returns a secure baseline config.
func DefaultSecurityMiddlewareConfig() SecurityMiddlewareConfig {
	return SecurityMiddlewareConfig{
		XContentTypeOptions:     "nosniff",
		XFrameOptions:           "DENY",
		XXSSProtection:          "1; mode=block",
		ReferrerPolicy:          "no-referrer",
		ContentSecurityPolicy:   "default-src 'self'",
		CacheControl:            "no-store",
		StrictTransportSecurity: "max-age=63072000; includeSubDomains; preload",
	}
}

func SecurityMiddleware(config ...SecurityMiddlewareConfig) func(http.Handler) http.Handler {
	c := DefaultSecurityMiddlewareConfig()
	if len(config) != 0 {
		c = config[0]
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c.XContentTypeOptions != "" {
				w.Header().Set("X-Content-Type-Options", c.XContentTypeOptions)
			}
			if c.XFrameOptions != "" {
				w.Header().Set("X-Frame-Options", c.XFrameOptions)
			}
			if c.XXSSProtection != "" {
				w.Header().Set("X-XSS-Protection", c.XXSSProtection)
			}
			if c.ReferrerPolicy != "" {
				w.Header().Set("Referrer-Policy", c.ReferrerPolicy)
			}
			if c.ContentSecurityPolicy != "" {
				w.Header().Set("Content-Security-Policy", c.ContentSecurityPolicy)
			}
			if c.CacheControl != "" {
				w.Header().Set("Cache-Control", c.CacheControl)
			}
			if c.StrictTransportSecurity != "" {
				w.Header().Set("Strict-Transport-Security", c.StrictTransportSecurity)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequestIdMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqId := xid.New().String()
			ctx := context.WithValue(r.Context(), ContextKeyRequestId, reqId)
			w.Header().Set("X-Draken-Request-Id", reqId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WebserverMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Server", "draken")
			w.Header().Set("X-Draken-Version", "v1")

			next.ServeHTTP(w, r)
		})
	}
}
