package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      app.routes(), // TODO: handle this route
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)
	go func() {
		// Create a quit channel which carries os.Signal values. Use buffered
		quit := make(chan os.Signal, 1)

		// Use signal.Notify() to listen for incoming SIGINT and SIGTERM signals and relay
		// them to the quit channel. Any other signal will not be caught by signal.Notify()
		// and will retain their default behavior.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Read the signal from the quit channel. This code will block until a signal is
		// received.
		s := <-quit

		app.logger.Info().Str("signal", s.String()).Msg("Caught signal termination")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info().
			Dict("server", zerolog.Dict().Str("addr", srv.Addr).Str("env", app.config.Env)).
			Msg("Completing background tasks")

		app.wg.Wait()
		shutdownError <- nil

	}()

	app.logger.Info().
		Dict("server", zerolog.Dict().Str("addr", srv.Addr).Str("env", app.config.Env)).
		Msg("Server running!")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info().
		Dict("server", zerolog.Dict().Str("addr", srv.Addr).Str("env", app.config.Env)).
		Msg("Server stop")

	return nil
}
