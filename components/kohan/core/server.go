package core

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type MonitorServer struct {
	mux         *gin.Engine
	capturePath string
	autoManager manager.AutoManagerInterface
}

func NewMonitorServer(capturePath string, autoManager manager.AutoManagerInterface) *MonitorServer {
	server := &MonitorServer{
		mux:         gin.Default(),
		capturePath: capturePath,
		autoManager: autoManager,
	}

	// Register Routes
	server.mux.GET("/v1/ticker/:ticker/record", server.HandleRecordTicker)
	server.mux.GET("/v1/clip/", server.HandleReadClip)
	server.mux.POST("/v1/submap/:action", server.HandleSubmapControl)

	return server
}

func (s *MonitorServer) Start(port int) (err error) {
	log.Info().Int("port", port).Msg("Starting Monitor Server")
	err = s.mux.Run(fmt.Sprintf(":%d", port))
	return
}

// StartWithShutdownHandler starts the server with graceful shutdown support using util.Shutdown
func (s *MonitorServer) StartWithShutdownHandler(port int, shutdown util.Shutdown) error {
	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.mux,
	}

	// Start server in goroutine so it won't block graceful shutdown handling
	errChan := make(chan error, 1)
	go func() {
		log.Info().Int("port", port).Msg("Starting Monitor Server with graceful shutdown")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Check for startup errors or wait for shutdown signal (following fun-app pattern)
	select {
	case err := <-errChan:
		log.Error().Err(err).Msg("Failed to start monitor server")
		return fmt.Errorf("server start failed: %w", err)
	case <-time.After(time.Second):
		// No error occurred, wait for graceful shutdown signal
		shutdownCtx := shutdown.Wait() // This blocks until SIGTERM/SIGINT or programmatic stop

		// Graceful shutdown with timeout
		ctxTimed, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		log.Info().Msg("Shutting down monitor server gracefully")
		if err := srv.Shutdown(ctxTimed); err != nil {
			log.Error().Err(err).Msg("Server forced shutdown")
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		log.Info().Ctx(shutdownCtx).Msg("Monitor server shutdown complete")
		return nil
	}
}

func (s *MonitorServer) HandleReadClip(ctx *gin.Context) {
	text, err := tools.ClipPaste()
	if err == nil {
		ctx.JSON(http.StatusOK, text)
	} else {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
}

func (s *MonitorServer) HandleRecordTicker(ctx *gin.Context) {
	ticker := ctx.Param("ticker")
	if err := s.autoManager.RecordTicker(ctx, ticker, s.capturePath); err == nil {
		ctx.JSON(http.StatusOK, "Success")
	} else {
		log.Error().Str("Ticker", ticker).Err(err).Msg("Record Ticker Failed")
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
}

func (s *MonitorServer) HandleSubmapControl(ctx *gin.Context) {
	action := ctx.Param("action")

	var request struct {
		Submap string `json:"submap"`
		Ticker string `json:"ticker,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&request); err == nil {
		switch action {
		case "enable":
			err = tools.HyperDispatch("submap " + request.Submap)
		case "disable":
			err = tools.HyperDispatch("submap reset")
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action. Use 'enable' or 'disable'"})
			return
		}

		if err == nil {
			log.Info().Str("Action", action).Str("Submap", request.Submap).Msg("Submap Control")
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "action": action, "submap": request.Submap})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
