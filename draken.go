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

	"github.com/joomcode/errorx"
)

type Draken struct {
	Config    Config
	Storage   Storage
	StartedAt time.Time
	R2        *R2
	Router    *Router
}

func New() (*Draken, error) {
	d := &Draken{}
	if err := d.setup(); err != nil {
		return nil, errorx.Decorate(err, "setup failed")
	}
	d.initStorage()
	d.initR2()

	log.Info().Msg("Created Draken app.")
	return d, nil
}

func (d *Draken) Serve() error {
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

		if err := d.Router.Echo.Shutdown(ctx); err != nil {
			log.Error().Err(errorx.IllegalState.New("shutdown failed")).Send()
		}
		close(idleConnsClosed)
	}()

	log.Info().Msgf("Listening on port %d", d.Config.Server.Port)
	if err := d.Router.Echo.Start(fmt.Sprintf(":%d", d.Config.Server.Port)); err != http.ErrServerClosed {
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

		if err := d.Router.Echo.Shutdown(ctx); err != nil {
			log.Error().Err(errorx.IllegalState.New("shutdown failed")).Send()
		}
		close(idleConnsClosed)
	}()

	log.Info().Msgf("Listening on port %d", d.Config.Server.Port)
	if err := d.Router.Echo.StartTLS(fmt.Sprintf(":%d", d.Config.Server.Port), tlsConfig.CertFile, tlsConfig.KeyFile); err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	log.Info().Msg("Graceful shutdown finished.")
	return nil
}

func DrakenHandler(w http.ResponseWriter, r *http.Request) (*Response, *Request) {
	return GetResponse(w), GetRequest(r)
}
