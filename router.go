package draken

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Router struct {
	*chi.Mux
	Draken       *Draken
	ParentRouter *Router
	Subrouters   map[string]*Router
}

func (d *Draken) CreateRouter() {
	log.Debug().Msg("Creating router...")
	d.Router = &Router{
		Mux:          chi.NewRouter(),
		Draken:       d,
		ParentRouter: nil,
		Subrouters:   make(map[string]*Router),
	}
	log.Info().Msg("Created router.")
}

func (r *Router) CreateSubrouter(route string) *Router {
	log.Debug().Str("route", route).Msgf("Creating subrouter at %s...", route)
	sr := &Router{
		Mux:          chi.NewRouter(),
		Draken:       r.Draken,
		ParentRouter: r,
		Subrouters:   make(map[string]*Router),
	}

	r.Subrouters[route] = sr
	r.Mount(route, sr)
	log.Info().Str("route", route).Msgf("Created router at %s.", route)
	return sr
}

func (r *Router) Get(route string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	log.Debug().Str("method", "GET").Str("route", route).Msg("Registing handler...")
	if len(middlewares) != 0 {
		r.Mux.With(middlewares...).Get(route, handler)
	} else {
		r.Mux.Get(route, handler)
	}
	log.Debug().Str("method", "GET").Str("route", route).Msg("Registered a handler")
}

func (r *Router) Post(route string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	log.Debug().Str("method", "GET").Str("route", route).Msg("Registing handler...")
	if len(middlewares) != 0 {
		r.Mux.With(middlewares...).Post(route, handler)
	} else {
		r.Mux.Post(route, handler)
	}
	log.Debug().Str("method", "POST").Str("route", route).Msg("Registered a handler")
}

func (r *Router) Put(route string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	log.Debug().Str("method", "GET").Str("route", route).Msg("Registing handler...")
	if len(middlewares) != 0 {
		r.Mux.With(middlewares...).Put(route, handler)
	} else {
		r.Mux.Put(route, handler)
	}
	log.Debug().Str("method", "PUT").Str("route", route).Msg("Registered a handler")
}

func (r *Router) Patch(route string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	log.Debug().Str("method", "GET").Str("route", route).Msg("Registing handler...")
	if len(middlewares) != 0 {
		r.Mux.With(middlewares...).Patch(route, handler)
	} else {
		r.Mux.Patch(route, handler)
	}
	log.Debug().Str("method", "PATCH").Str("route", route).Msg("Registered a handler")
}

func (r *Router) Delete(route string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	log.Debug().Str("method", "GET").Str("route", route).Msg("Registing handler...")
	if len(middlewares) != 0 {
		r.Mux.With(middlewares...).Delete(route, handler)
	} else {
		r.Mux.Delete(route, handler)
	}
	log.Debug().Str("method", "DELETE").Str("route", route).Msg("Registered a handler")
}
