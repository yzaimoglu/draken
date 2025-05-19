package draken

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	"github.com/joomcode/errorx"
)

type Draken struct {
	Config    Config
	Storage   Storage
	StartedAt time.Time
	Chi       *chi.Mux
}

func New() (*Draken, error) {
	d := &Draken{}
	if err := d.setup(); err != nil {
		return nil, errorx.Decorate(err, "setup failed")
	}
	d.InitStorage()

	log.Info().Msg("Created Draken app.")
	return d, nil
}

func (d *Draken) CreateRouter() {
	d.Chi = chi.NewRouter()
	log.Info().Msg("Created Chi router.")
}

func (d *Draken) Get(route string, handler http.HandlerFunc) {
	d.Chi.Get(route, handler)
	log.Debug().Str("method", "GET").Str("route", route).Msg("Registered a handler")
}

func (d *Draken) Post(route string, handler http.HandlerFunc) {
	d.Chi.Post(route, handler)
	log.Debug().Str("method", "POST").Str("route", route).Msg("Registered a handler")
}

func (d *Draken) Put(route string, handler http.HandlerFunc) {
	d.Chi.Put(route, handler)
	log.Debug().Str("method", "PUT").Str("route", route).Msg("Registered a handler")
}

func (d *Draken) Patch(route string, handler http.HandlerFunc) {
	d.Chi.Patch(route, handler)
	log.Debug().Str("method", "PATCH").Str("route", route).Msg("Registered a handler")
}

func (d *Draken) Delete(route string, handler http.HandlerFunc) {
	d.Chi.Delete(route, handler)
	log.Debug().Str("method", "DELETE").Str("route", route).Msg("Registered a handler")
}

func (d *Draken) Serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", d.Config.Server.Port),
		Handler: d.Chi,
	}

	// Channel to listen for termination signals
	idleConnsClosed := make(chan struct{})

	go func() {
		// Listen for OS interrupt or termination signals
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		log.Debug().Msg("Graceful shutdown initiated.")
		// Initiate graceful shutdown with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(errorx.IllegalState.New("shutdown failed")).Send()
		}
		close(idleConnsClosed)
	}()

	log.Info().Msgf("Listening on port %d", d.Config.Server.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	log.Info().Msg("Graceful shutdown finished.")
	return nil
}

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

func (d *Draken) ServeTLS(tlsConfig TLSConfig) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", d.Config.Server.Port),
		Handler: d.Chi,
	}

	// Channel to listen for termination signals
	idleConnsClosed := make(chan struct{})

	go func() {
		// Listen for OS interrupt or termination signals
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		log.Debug().Msg("Graceful shutdown initiated.")
		// Initiate graceful shutdown with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(errorx.IllegalState.New("shutdown failed")).Send()
		}
		close(idleConnsClosed)
	}()

	log.Info().Msgf("Listening on port %d", d.Config.Server.Port)
	if err := srv.ListenAndServeTLS(tlsConfig.CertFile, tlsConfig.KeyFile); err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	log.Info().Msg("Graceful shutdown finished.")
	return nil
}

func DrakenHandler(w http.ResponseWriter, r *http.Request) (*Response, *Request) {
	return GetResponse(w), GetRequest(r)
}
