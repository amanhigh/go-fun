package util

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

// 3.2 FIXME: Introduce an HTTPServer interface plus base struct so fun-app and kohan can share lifecycle wiring while letting DI register handlers.
// NewHTTPServer creates an http.Server with standard timeouts for the given port and handler.
func NewHTTPServer(port int, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           handler,
		ReadHeaderTimeout: defaultReadTimeout,
		ReadTimeout:       defaultReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
	}
}

// RunHTTPServer starts the given http.Server and blocks until a graceful shutdown
// signal is received via the Shutdown interface. Returns nil on clean shutdown.
func RunHTTPServer(name string, srv *http.Server, shutdown Shutdown) error {
	errChan := make(chan error, 1)
	serverStopped := make(chan struct{})

	go listenAndServe(name, srv, errChan, serverStopped)
	go waitForShutdown(name, srv, shutdown, errChan, serverStopped)

	err := <-errChan
	close(serverStopped)
	if err != nil {
		log.Error().Err(err).Str("server", name).Msg("Server error occurred")
	}
	return err
}

func listenAndServe(name string, srv *http.Server, errChan chan<- error, stopped <-chan struct{}) {
	log.Info().Str("addr", srv.Addr).Str("server", name).Msg("Starting server")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		select {
		case errChan <- err:
		case <-stopped:
		}
	}
}

func waitForShutdown(name string, srv *http.Server, shutdown Shutdown, errChan chan<- error, stopped <-chan struct{}) {
	shutdownCtx := shutdown.Wait()

	ctxTimed, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	log.Info().Str("server", name).Msg("Shutting down server gracefully")
	if err := srv.Shutdown(ctxTimed); err != nil {
		log.Error().Err(err).Str("server", name).Msg("Server forced shutdown")
		select {
		case errChan <- fmt.Errorf("graceful shutdown failed: %w", err):
		case <-stopped:
		}
		return
	}

	log.Info().Ctx(shutdownCtx).Str("server", name).Msg("Server shutdown complete")
	select {
	case errChan <- nil:
	case <-stopped:
	}
}
