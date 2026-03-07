package util

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/common/telemetry"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

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

// HttpServer provides shared HTTP server lifecycle: creation with standard
// timeouts, a default /health route, graceful start/stop, and overridable hooks.
//
// Call SetLifecycle to attach a ServerLifecycle implementation.
// HACK: Extract Interface and rename struct to HttpServerImpl
type HttpServer struct {
	Name      string
	Engine    *gin.Engine
	Server    *http.Server
	lifecycle ServerLifecycle
	shutdown  Shutdown
}

// NewHttpServer creates a HttpServer with the provided gin.Engine, standard timeouts,
// and a /health route.
func NewHttpServer(cfg config.HttpServerConfig, engine *gin.Engine, shutdown Shutdown) *HttpServer {
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

	return &HttpServer{
		Name:      cfg.Name,
		Engine:    engine,
		Server:    srv,
		lifecycle: noopLifecycle{},
		shutdown:  shutdown,
	}
}

// SetLifecycle attaches a ServerLifecycle to this server.
func (h *HttpServer) SetLifecycle(lc ServerLifecycle) {
	h.lifecycle = lc
}

// Start registers routes, runs BeforeStart hook, begins listening, and blocks until graceful shutdown completes.
func (h *HttpServer) Start() error {
	h.lifecycle.RegisterRoutes(h.Engine)
	h.lifecycle.BeforeStart(context.Background())

	errChan := make(chan error, 1)
	serverStopped := make(chan struct{})

	go h.listenAndServe(errChan, serverStopped)
	go h.waitForShutdown(errChan, serverStopped)

	err := <-errChan
	close(serverStopped)
	if err != nil {
		log.Error().Err(err).Str("server", h.Name).Msg("Server error occurred")
	}
	return err
}

// Stop triggers graceful shutdown programmatically.
func (h *HttpServer) Stop(ctx context.Context) {
	h.shutdown.Stop(ctx)
}

func (h *HttpServer) listenAndServe(errChan chan<- error, stopped <-chan struct{}) {
	log.Info().Str("addr", h.Server.Addr).Str("server", h.Name).Msg("Starting server")
	if err := h.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		select {
		case errChan <- err:
		case <-stopped:
		}
	}
}

func (h *HttpServer) waitForShutdown(errChan chan<- error, stopped <-chan struct{}) {
	shutdownCtx := h.shutdown.Wait()

	h.lifecycle.BeforeShutdown(shutdownCtx)

	ctxTimed, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	log.Info().Str("server", h.Name).Msg("Shutting down server gracefully")
	if err := h.Server.Shutdown(ctxTimed); err != nil {
		log.Error().Err(err).Str("server", h.Name).Msg("Server forced shutdown")
		select {
		case errChan <- fmt.Errorf("graceful shutdown failed: %w", err):
		case <-stopped:
		}
		return
	}

	h.lifecycle.AfterShutdown(shutdownCtx)

	log.Info().Ctx(shutdownCtx).Str("server", h.Name).Msg("Server shutdown complete")
	select {
	case errChan <- nil:
	case <-stopped:
	}
}
