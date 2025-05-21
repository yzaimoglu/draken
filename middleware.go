package draken

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ContextKey string

const ContextKeyRequestId ContextKey = "draken-request-id"

func (r *Router) Middleware(middlewares ...echo.MiddlewareFunc) {
	r.Group.Use(middlewares...)
}

func (r *Router) EssentialMiddlewares() {
	if !r.Draken.Config.Server.Hidden {
		r.Middleware(WebserverMiddleware())
	}

	r.Middleware(RequestIdMiddleware())
	r.Middleware(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	r.Middleware(LoggerMiddleware(log.Logger))
	r.Middleware(middleware.Recover())
	if r.Draken.Config.Server.Security {
		r.Middleware(middleware.Secure())
	}

	if r.Draken.Config.Server.Heartbeat.Enabled {
		r.Get(r.Draken.Config.Server.Heartbeat.Endpoint, HeartbeatRoute)
	}
}

func RequestIdMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := xid.New().String()
			c.Set(string(ContextKeyRequestId), id)
			c.Response().Header().Set("X-Draken-Request-Id", id)
			return next(c)
		}
	}
}

func WebserverMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Response().Header()
			h.Set("Server", "draken")
			h.Set("X-Draken-Version", "v1")
			return next(c)
		}
	}
}

func LoggerMiddleware(l zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// run downstream
			err := next(c)

			res := c.Response()
			req := c.Request()

			l.Info().
				Str("method", req.Method).
				Str("url", req.URL.String()).
				Str("proto", req.Proto).
				Str("remote", c.RealIP()).
				Int("status", res.Status).
				Int64("bytes", res.Size).
				Dur("duration", time.Since(start)).
				Any("request_id", c.Get(string(ContextKeyRequestId))).
				Msgf(`%s %s %s from %s - %d %dB in %s`,
					req.Method,
					req.URL.String(),
					req.Proto,
					c.RealIP(),
					res.Status,
					res.Size,
					time.Since(start))

			return err
		}
	}
}

func HeartbeatRoute(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "im alive")
}
