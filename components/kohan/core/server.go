package core

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/gin-gonic/gin"
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
	if s.monitorHandler == nil {
		return
	}
	engine.GET("/v1/ticker/:ticker/record", s.monitorHandler.HandleRecordTicker)
	engine.GET("/v1/clip/", s.monitorHandler.HandleReadClip)
	engine.POST("/v1/submap/:action", s.monitorHandler.HandleSubmapControl)
}

func (s *KohanServer) registerJournalRoutes(engine *gin.Engine) {
	entries := engine.Group("/v1/journal-entries")
	{
		entries.GET("", s.journalHandler.HandleListEntries)
		entries.GET("/:id", s.journalHandler.HandleGetEntry)
		entries.POST("", s.journalHandler.HandleCreateEntry)

		entries.POST("/:id/images", s.journalHandler.HandleCreateImage)
		entries.GET("/:id/images", s.journalHandler.HandleListImages)
		entries.DELETE("/:id/images/:imageId", s.journalHandler.HandleDeleteImage)

		entries.POST("/:id/notes", s.journalHandler.HandleCreateNote)
		entries.GET("/:id/notes", s.journalHandler.HandleListNotes)
		entries.DELETE("/:id/notes/:noteId", s.journalHandler.HandleDeleteNote)

		entries.POST("/:id/tags", s.journalHandler.HandleCreateTag)
		entries.GET("/:id/tags", s.journalHandler.HandleListTags)
		entries.DELETE("/:id/tags/:tagId", s.journalHandler.HandleDeleteTag)
	}
}
