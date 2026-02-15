package core

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/gin-gonic/gin"
)

// KohanServer serves all Kohan HTTP APIs (monitor + journal).
type KohanServer struct {
	*util.BaseHTTPServer
	MonitorHandler handler.MonitorHandler
	JournalHandler handler.JournalHandler
	ImageHandler   handler.ImageHandler
	NoteHandler    handler.NoteHandler
	TagHandler     handler.TagHandler
}

// NewKohanServer creates a KohanServer for testing with explicit handler injection.
func NewKohanServer(base *util.BaseHTTPServer, monitorHandler handler.MonitorHandler,
	journalHandler handler.JournalHandler, imageHandler handler.ImageHandler,
	noteHandler handler.NoteHandler, tagHandler handler.TagHandler) *KohanServer {
	server := &KohanServer{
		BaseHTTPServer: base,
		MonitorHandler: monitorHandler,
		JournalHandler: journalHandler,
		ImageHandler:   imageHandler,
		NoteHandler:    noteHandler,
		TagHandler:     tagHandler,
	}
	server.RegisterRoutes = server.registerRoutes
	return server
}

func (s *KohanServer) registerRoutes(engine *gin.Engine) {
	s.registerMonitorRoutes(engine)
	s.registerJournalRoutes(engine)
}

func (s *KohanServer) registerMonitorRoutes(engine *gin.Engine) {
	if s.MonitorHandler == nil {
		return
	}
	engine.GET("/v1/ticker/:ticker/record", s.MonitorHandler.HandleRecordTicker)
	engine.GET("/v1/clip/", s.MonitorHandler.HandleReadClip)
	engine.POST("/v1/submap/:action", s.MonitorHandler.HandleSubmapControl)
}

func (s *KohanServer) registerJournalRoutes(engine *gin.Engine) {
	entries := engine.Group("/v1/journal-entries")
	{
		entries.GET("", s.JournalHandler.HandleListEntries)
		entries.GET("/:id", s.JournalHandler.HandleGetEntry)
		entries.POST("", s.JournalHandler.HandleCreateEntry)

		entries.POST("/:id/images", s.ImageHandler.HandleCreateImage)
		entries.GET("/:id/images", s.ImageHandler.HandleListImages)
		entries.DELETE("/:id/images/:imageId", s.ImageHandler.HandleDeleteImage)

		entries.POST("/:id/notes", s.NoteHandler.HandleCreateNote)
		entries.GET("/:id/notes", s.NoteHandler.HandleListNotes)
		entries.DELETE("/:id/notes/:noteId", s.NoteHandler.HandleDeleteNote)

		entries.POST("/:id/tags", s.TagHandler.HandleCreateTag)
		entries.GET("/:id/tags", s.TagHandler.HandleListTags)
		entries.DELETE("/:id/tags/:tagId", s.TagHandler.HandleDeleteTag)
	}
}
