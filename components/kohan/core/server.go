package core

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/gin-gonic/gin"
)

// KohanServer serves all Kohan HTTP APIs (monitor + journal).
type KohanServer struct {
	*util.BaseHTTPServer
	// FIXME: Inject via Named Tags leave constructor only for test.
	monitorHandler handler.MonitorHandler
	journalHandler handler.JournalHandler
	imageHandler   handler.ImageHandler
	noteHandler    handler.NoteHandler
	tagHandler     handler.TagHandler
}

// NewKohanServer creates a KohanServer with interface-injected handlers.
func NewKohanServer(base *util.BaseHTTPServer, monitorHandler handler.MonitorHandler,
	journalHandler handler.JournalHandler, imageHandler handler.ImageHandler,
	noteHandler handler.NoteHandler, tagHandler handler.TagHandler) *KohanServer {
	server := &KohanServer{
		BaseHTTPServer: base,
		monitorHandler: monitorHandler,
		journalHandler: journalHandler,
		imageHandler:   imageHandler,
		noteHandler:    noteHandler,
		tagHandler:     tagHandler,
	}
	server.RegisterRoutes = server.registerRoutes
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

		entries.POST("/:id/images", s.imageHandler.HandleCreateImage)
		entries.GET("/:id/images", s.imageHandler.HandleListImages)
		entries.DELETE("/:id/images/:imageId", s.imageHandler.HandleDeleteImage)

		entries.POST("/:id/notes", s.noteHandler.HandleCreateNote)
		entries.GET("/:id/notes", s.noteHandler.HandleListNotes)
		entries.DELETE("/:id/notes/:noteId", s.noteHandler.HandleDeleteNote)

		entries.POST("/:id/tags", s.tagHandler.HandleCreateTag)
		entries.GET("/:id/tags", s.tagHandler.HandleListTags)
		entries.DELETE("/:id/tags/:tagId", s.tagHandler.HandleDeleteTag)
	}
}
