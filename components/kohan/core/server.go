package core

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// KohanServer serves all Kohan HTTP APIs (monitor + journal).
type KohanServer struct {
	*util.BaseHTTPServer
	monitorHandler handler.MonitorHandler
	journalHandler handler.JournalHandler
}

// NewKohanServer creates a KohanServer with interface-injected handlers.
func NewKohanServer(port int, monitorHandler handler.MonitorHandler, journalHandler handler.JournalHandler, shutdown util.Shutdown) *KohanServer {
	server := &KohanServer{
		// FIXME: Create BaseHTTPServer with DI Framework and inject here, shutdown also should be created by DI Framework.
		BaseHTTPServer: util.NewBaseHTTPServer("kohan", port, shutdown),
		monitorHandler: monitorHandler,
		journalHandler: journalHandler,
	}
	server.BaseHTTPServer.RegisterRoutes = server.registerRoutes
	return server
}

func (s *KohanServer) registerRoutes(engine *gin.Engine) {
	s.registerMonitorRoutes(engine)
	s.registerJournalRoutes(engine)
}

func (s *KohanServer) registerMonitorRoutes(engine *gin.Engine) {
	// HACK: Nil Handlers should not be checked rely on DI Framework to inject them.
	if s.monitorHandler == nil {
		return
	}
	engine.GET("/v1/ticker/:ticker/record", s.monitorHandler.HandleRecordTicker)
	engine.GET("/v1/clip/", s.monitorHandler.HandleReadClip)
	engine.POST("/v1/submap/:action", s.monitorHandler.HandleSubmapControl)
}

func (s *KohanServer) registerJournalRoutes(engine *gin.Engine) {
	// FIXME: Remove Handler Nil Checks.
	if s.journalHandler != nil {
		// BUG: Match Monitor Routes Pattern use /v1 only.
		v1 := engine.Group("/api/v1")
		{
			entries := v1.Group("/journal-entries")
			{
				entries.GET("", s.journalHandler.HandleListEntries)
				entries.GET("/:id", s.journalHandler.HandleGetEntry)
				entries.POST("", s.journalHandler.HandleCreateEntry)
			}
		}
	} else {
		log.Warn().Msg("JournalHandler is nil, journal API routes will not be registered")
	}
}
