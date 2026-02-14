package core

import (
	"context"
	"fmt"
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// BarkatServer serves the Barkat Journal Explorer API.
// FIXME: Merge with KohanServer and combine command as well to start the server.
// FIXME: Introduce Handler package and move all handlers there.
type BarkatServer struct {
	mux           *gin.Engine
	barkatManager manager.BarkatManager
}

// NewBarkatServer creates a new BarkatServer with registered routes.
func NewBarkatServer(barkatManager manager.BarkatManager) *BarkatServer {
	server := &BarkatServer{
		mux:           gin.Default(),
		barkatManager: barkatManager,
	}

	v1 := server.mux.Group("/api/v1")
	{
		entries := v1.Group("/journal-entries")
		{
			entries.GET("", server.HandleListEntries)
			entries.GET("/:id", server.HandleGetEntry)
			entries.POST("", server.HandleCreateEntry)
		}
	}

	return server
}

// NewBarkatServerWithDB creates a BarkatServer from an existing *gorm.DB (useful for testing).
func NewBarkatServerWithDB(db *gorm.DB) *BarkatServer {
	repo := repository.NewBarkatRepository(db)
	mgr := manager.NewBarkatManager(repo)
	return NewBarkatServer(mgr)
}

// StartOnPort starts the barkat server on the given port (blocking, no graceful shutdown).
func (s *BarkatServer) StartOnPort(port int) error {
	if err := s.mux.Run(fmt.Sprintf(":%d", port)); err != nil {
		return fmt.Errorf("barkat server failed: %w", err)
	}
	return nil
}

// Start starts the barkat server with graceful shutdown support.
func (s *BarkatServer) Start(port int, shutdown util.Shutdown) error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           s.mux,
		ReadHeaderTimeout: readTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
	}

	errChan := make(chan error, 1)
	serverStopped := make(chan struct{})

	go s.listenAndServe(srv, errChan, serverStopped, port)
	go s.waitForShutdown(srv, shutdown, errChan, serverStopped)

	err := <-errChan
	close(serverStopped)
	if err != nil {
		log.Error().Err(err).Msg("Barkat server error")
	}
	return err
}

func (s *BarkatServer) listenAndServe(srv *http.Server, errChan chan<- error, stopped <-chan struct{}, port int) {
	log.Info().Int("port", port).Msg("Starting Barkat Server")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		select {
		case errChan <- err:
		case <-stopped:
		}
	}
}

func (s *BarkatServer) waitForShutdown(srv *http.Server, shutdown util.Shutdown, errChan chan<- error, stopped <-chan struct{}) {
	shutdownCtx := shutdown.Wait()
	ctxTimed, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Info().Msg("Shutting down barkat server")
	if err := srv.Shutdown(ctxTimed); err != nil {
		log.Error().Err(err).Msg("Barkat server forced shutdown")
		select {
		case errChan <- fmt.Errorf("graceful shutdown failed: %w", err):
		case <-stopped:
		}
		return
	}

	log.Info().Ctx(shutdownCtx).Msg("Barkat server shutdown complete")
	select {
	case errChan <- nil:
	case <-stopped:
	}
}

// HandleListEntries handles GET /api/v1/journal-entries
func (s *BarkatServer) HandleListEntries(c *gin.Context) {
	var query barkat.EntryQuery
	query.Limit = 10 // Default limit

	if err := c.ShouldBindQuery(&query); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	entryList, httpErr := s.barkatManager.ListEntries(c.Request.Context(), query)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, entryList)
}

// HandleGetEntry handles GET /api/v1/journal-entries/:id
func (s *BarkatServer) HandleGetEntry(c *gin.Context) {
	var path barkat.EntryPath
	if err := c.ShouldBindUri(&path); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	entry, httpErr := s.barkatManager.GetEntry(c.Request.Context(), path.ID)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, entry)
}

// HandleCreateEntry handles POST /api/v1/journal-entries
func (s *BarkatServer) HandleCreateEntry(c *gin.Context) {
	var entry barkat.Entry
	if err := c.ShouldBindJSON(&entry); err != nil {
		err = util.ProcessValidationError(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if httpErr := s.barkatManager.CreateEntry(c.Request.Context(), &entry); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, entry)
}
