package draken

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Router struct {
	Echo         *echo.Echo
	Group        *echo.Group
	Draken       *Draken
	ParentRouter *Router
	Subrouters   map[string]*Router
}

func (d *Draken) CreateRouter() {
	log.Debug().Msg("Creating router...")
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	g := e.Group("")
	d.Router = &Router{
		Echo:         e,
		Group:        g,
		ParentRouter: nil,
		Draken:       d,
		Subrouters:   make(map[string]*Router),
	}
	log.Info().Msg("Created router.")
}

func (r *Router) CreateSubrouter(route string) *Router {
	log.Debug().Str("route", route).Msgf("Creating subrouter at %s...", route)
	sr := &Router{
		Echo:         r.Echo,
		Group:        r.Group.Group(route),
		Draken:       r.Draken,
		ParentRouter: r,
		Subrouters:   make(map[string]*Router),
	}

	r.Subrouters[route] = sr
	log.Info().Str("route", route).Msgf("Created router at %s.", route)
	return sr
}

func (r *Router) Get(route string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	log.Debug().Str("method", "GET").Str("route", route).Msg("Registing handler...")
	r.Group.GET(route, handler, middlewares...)
	log.Debug().Str("method", "GET").Str("route", route).Msg("Registered a handler")
}

func (r *Router) Post(route string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	log.Debug().Str("method", "POST").Str("route", route).Msg("Registing handler...")
	r.Group.POST(route, handler, middlewares...)
	log.Debug().Str("method", "POST").Str("route", route).Msg("Registered a handler")
}

func (r *Router) Put(route string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	log.Debug().Str("method", "PUT").Str("route", route).Msg("Registing handler...")
	r.Group.PUT(route, handler, middlewares...)
	log.Debug().Str("method", "PUT").Str("route", route).Msg("Registered a handler")
}

func (r *Router) Patch(route string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	log.Debug().Str("method", "PATCH").Str("route", route).Msg("Registing handler...")
	r.Group.PATCH(route, handler, middlewares...)
	log.Debug().Str("method", "PATCH").Str("route", route).Msg("Registered a handler")
}

func (r *Router) Delete(route string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	log.Debug().Str("method", "DELETE").Str("route", route).Msg("Registing handler...")
	r.Group.DELETE(route, handler, middlewares...)
	log.Debug().Str("method", "DELETE").Str("route", route).Msg("Registered a handler")
}
