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

// HttpServerConfig holds the minimal configuration needed to create a BaseHTTPServer.
type HttpServerConfig struct {
	// FIXME: Move to Models Package.
	Name string
	Port int
	// BUG: Shutdown should be injected as a dependency, not part of config.
	Shutdown Shutdown
}

// ServerLifecycle defines hooks that customise HTTP server behaviour.
// Implement this interface and attach via SetLifecycle to override the default no-ops.
type ServerLifecycle interface {
	// RegisterRoutes is called once during Start to register application routes.
	RegisterRoutes(engine *gin.Engine)
	// BeforeStart runs after routes are registered but before the HTTP server begins listening.
	BeforeStart(ctx context.Context)
	// BeforeShutdown runs after the shutdown signal but before the HTTP server stops.
	BeforeShutdown(ctx context.Context)
	// AfterShutdown runs after the HTTP server has fully stopped.
	AfterShutdown(ctx context.Context)
}

// noopLifecycle is the default no-op ServerLifecycle.
type noopLifecycle struct{}

func (noopLifecycle) RegisterRoutes(_ *gin.Engine)     {}
func (noopLifecycle) BeforeStart(_ context.Context)    {}
func (noopLifecycle) BeforeShutdown(_ context.Context) {}
func (noopLifecycle) AfterShutdown(_ context.Context)  {}

// BaseHTTPServer provides shared HTTP server lifecycle: creation with standard
// timeouts, a default /health route, graceful start/stop, and overridable hooks.
//
// Call SetLifecycle to attach a ServerLifecycle implementation.
type BaseHTTPServer struct {
	Name      string
	Engine    *gin.Engine
	Server    *http.Server
	lifecycle ServerLifecycle
	shutdown  Shutdown
}

// NewBaseHTTPServer creates a BaseHTTPServer with a default gin.Engine, standard timeouts,
// and a /health route.
func NewBaseHTTPServer(cfg HttpServerConfig) *BaseHTTPServer {
	return NewBaseHTTPServerWithEngine(cfg, gin.Default())
}

// NewBaseHTTPServerWithEngine creates a BaseHTTPServer using a pre-configured gin.Engine
// (e.g. with custom middleware, rate limiting, prometheus) and a /health route.
func NewBaseHTTPServerWithEngine(cfg HttpServerConfig, engine *gin.Engine) *BaseHTTPServer {
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Statsviz: http://localhost:<port>/debug/statsviz/
	engine.GET("/debug/statsviz/*filepath", telemetry.StatvizMetrics)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           engine,
		ReadHeaderTimeout: defaultReadTimeout,
		ReadTimeout:       defaultReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
	}

	return &BaseHTTPServer{
		Name:      cfg.Name,
		Engine:    engine,
		Server:    srv,
		lifecycle: noopLifecycle{},
		shutdown:  cfg.Shutdown,
	}
}

// SetLifecycle attaches a ServerLifecycle to this server.
func (b *BaseHTTPServer) SetLifecycle(lc ServerLifecycle) {
	b.lifecycle = lc
}

// Start registers routes, runs BeforeStart hook, begins listening, and blocks until graceful shutdown completes.
func (b *BaseHTTPServer) Start() error {
	b.lifecycle.RegisterRoutes(b.Engine)
	b.lifecycle.BeforeStart(context.Background())

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

	b.lifecycle.BeforeShutdown(shutdownCtx)

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

	b.lifecycle.AfterShutdown(shutdownCtx)

	log.Info().Ctx(shutdownCtx).Str("server", b.Name).Msg("Server shutdown complete")
	select {
	case errChan <- nil:
	case <-stopped:
	}
}
