package core

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/gin-gonic/gin"
)

// KohanServer serves all Kohan HTTP APIs (monitor + journal).
type KohanServer struct {
	// HACK: Should we inject pointer to BaseHTTPServer or not?
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
	monitor := engine.Group("/v1/monitor")
	handler.SetupMonitorRoutes(monitor, s.MonitorHandler)
}

func (s *KohanServer) registerJournalRoutes(engine *gin.Engine) {
	journal := engine.Group("/v1/journal")
	handler.SetupJournalEntryRoutes(journal, s.JournalHandler)
	handler.SetupImageRoutes(journal, s.ImageHandler)
	handler.SetupNoteRoutes(journal, s.NoteHandler)
	handler.SetupTagRoutes(journal, s.TagHandler)
}
