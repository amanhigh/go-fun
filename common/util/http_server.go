package util

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/common/telemetry"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

// BaseHTTPServer provides shared HTTP server lifecycle: creation with standard
// timeouts, a default /health route, graceful start/stop, and overridable hooks.
//
// Embedders override func fields to customise route registration and shutdown behaviour.
type BaseHTTPServer struct {
	Name   string
	Engine *gin.Engine
	Server *http.Server

	// RegisterRoutes is called once during Start to register application routes.
	// Default is a no-op (/health is already registered by the constructor).
	// Embedders replace this to add their own routes on the gin.Engine.
	RegisterRoutes func(engine *gin.Engine)

	// BeforeStart runs after routes are registered but before the HTTP server begins listening.
	// Default is a no-op. Embedders replace this to start background services.
	BeforeStart func(ctx context.Context)

	// BeforeShutdown runs after the shutdown signal but before the HTTP server stops.
	// Default is a no-op. Embedders replace this to add custom pre-shutdown logic.
	BeforeShutdown func(ctx context.Context)

	// AfterShutdown runs after the HTTP server has fully stopped.
	// Default is a no-op. Embedders replace this to add custom post-shutdown logic.
	AfterShutdown func(ctx context.Context)

	shutdown Shutdown
}

// NewBaseHTTPServer creates a BaseHTTPServer with a default gin.Engine, standard timeouts,
// and a /health route.
func NewBaseHTTPServer(name string, port int, shutdown Shutdown) *BaseHTTPServer {
	return NewBaseHTTPServerWithEngine(name, port, shutdown, gin.Default())
}

// NewBaseHTTPServerWithEngine creates a BaseHTTPServer using a pre-configured gin.Engine
// (e.g. with custom middleware, rate limiting, prometheus) and a /health route.
func NewBaseHTTPServerWithEngine(name string, port int, shutdown Shutdown, engine *gin.Engine) *BaseHTTPServer {
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Statsviz: http://localhost:<port>/debug/statsviz/
	engine.GET("/debug/statsviz/*filepath", telemetry.StatvizMetrics)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           engine,
		ReadHeaderTimeout: defaultReadTimeout,
		ReadTimeout:       defaultReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
	}

	noop := func(_ context.Context) {}
	noopRoutes := func(_ *gin.Engine) {}
	return &BaseHTTPServer{
		Name:           name,
		Engine:         engine,
		Server:         srv,
		RegisterRoutes: noopRoutes,
		BeforeStart:    noop,
		BeforeShutdown: noop,
		AfterShutdown:  noop,
		shutdown:       shutdown,
	}
}

// Start registers routes, runs BeforeStart hook, begins listening, and blocks until graceful shutdown completes.
func (b *BaseHTTPServer) Start() error {
	b.RegisterRoutes(b.Engine)
	b.BeforeStart(context.Background())

	errChan := make(chan error, 1)
	serverStopped := make(chan struct{})

	go b.listenAndServe(errChan, serverStopped)
	go b.waitForShutdown(errChan, serverStopped)

	err := <-errChan
	close(serverStopped)
	if err != nil {
		log.Error().Err(err).Str("server", b.Name).Msg("Server error occurred")
	}
	return err
}

// Stop triggers graceful shutdown programmatically.
func (b *BaseHTTPServer) Stop(ctx context.Context) {
	b.shutdown.Stop(ctx)
}

func (b *BaseHTTPServer) listenAndServe(errChan chan<- error, stopped <-chan struct{}) {
	log.Info().Str("addr", b.Server.Addr).Str("server", b.Name).Msg("Starting server")
	if err := b.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		select {
		case errChan <- err:
		case <-stopped:
		}
	}
}

func (b *BaseHTTPServer) waitForShutdown(errChan chan<- error, stopped <-chan struct{}) {
	shutdownCtx := b.shutdown.Wait()

	b.BeforeShutdown(shutdownCtx)

	ctxTimed, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	log.Info().Str("server", b.Name).Msg("Shutting down server gracefully")
	if err := b.Server.Shutdown(ctxTimed); err != nil {
		log.Error().Err(err).Str("server", b.Name).Msg("Server forced shutdown")
		select {
		case errChan <- fmt.Errorf("graceful shutdown failed: %w", err):
		case <-stopped:
		}
		return
	}

	b.AfterShutdown(shutdownCtx)

	log.Info().Ctx(shutdownCtx).Str("server", b.Name).Msg("Server shutdown complete")
	select {
	case errChan <- nil:
	case <-stopped:
	}
}
