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

const (
	shutdownTimeout = 3 * time.Second
	readTimeout     = 5 * time.Second
	writeTimeout    = 5 * time.Second
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

// Start starts the server with graceful shutdown support using util.Shutdown
func (s *MonitorServer) Start(port int, shutdown util.Shutdown) error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           s.mux,
		ReadHeaderTimeout: readTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
	}

	errChan := make(chan error, 1)
	serverStopped := make(chan struct{})

	go s.runServer(srv, errChan, serverStopped, port)
	go s.handleShutdown(srv, shutdown, errChan, serverStopped)

	err := <-errChan
	close(serverStopped)
	if err != nil {
		log.Error().Err(err).Msg("Server error occurred")
		return err
	}
	return nil
}

func (s *MonitorServer) runServer(srv *http.Server, errChan chan<- error, serverStopped <-chan struct{}, port int) {
	log.Info().Int("port", port).Msg("Starting Monitor Server with graceful shutdown")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		select {
		case errChan <- err:
		case <-serverStopped:
		}
	}
}

func (s *MonitorServer) handleShutdown(srv *http.Server, shutdown util.Shutdown, errChan chan<- error, serverStopped <-chan struct{}) {
	shutdownCtx := shutdown.Wait()

	ctxTimed, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Info().Msg("Shutting down monitor server gracefully")
	if err := srv.Shutdown(ctxTimed); err != nil {
		log.Error().Err(err).Msg("Server forced shutdown")
		select {
		case errChan <- fmt.Errorf("graceful shutdown failed: %w", err):
		case <-serverStopped:
		}
		return
	}

	log.Info().Ctx(shutdownCtx).Msg("Monitor server shutdown complete")
	select {
	case errChan <- nil:
	case <-serverStopped:
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
