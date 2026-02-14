package core

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// KohanServer serves all Kohan HTTP APIs (monitor + journal).
type KohanServer struct {
	mux *gin.Engine
}

// NewKohanServer creates a KohanServer with monitor and journal routes registered.
func NewKohanServer(capturePath string, autoManager manager.AutoManagerInterface, journalHandler *handler.JournalHandler) *KohanServer {
	mux := gin.Default()

	// Monitor routes
	monitorHandler := handler.NewMonitorHandler(capturePath, autoManager)
	mux.GET("/v1/ticker/:ticker/record", monitorHandler.HandleRecordTicker)
	mux.GET("/v1/clip/", monitorHandler.HandleReadClip)
	mux.POST("/v1/submap/:action", monitorHandler.HandleSubmapControl)

	// Journal API routes
	// 3.3 BUG: Don’t silently accept nil handlers; expose RegisterHandlers hook so DI supplies a complete handler set and fail fast otherwise.
	if journalHandler != nil {
		v1 := mux.Group("/api/v1")
		{
			entries := v1.Group("/journal-entries")
			{
				entries.GET("", journalHandler.HandleListEntries)
				entries.GET("/:id", journalHandler.HandleGetEntry)
				entries.POST("", journalHandler.HandleCreateEntry)
			}
		}
	} else {
		log.Warn().Msg("JournalHandler is nil, journal API routes will not be registered")
	}

	return &KohanServer{mux: mux}
}

// Start starts the server with graceful shutdown support using util.Shutdown.
func (s *KohanServer) Start(port int, shutdown util.Shutdown) error {
	return util.RunHTTPServer("kohan", util.NewHTTPServer(port, s.mux), shutdown)
}
