package core

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
)

// KohanServerLifecycle implements ServerLifecycle for the Kohan HTTP server.
type KohanServerLifecycle struct {
	monitorHandler handler.MonitorHandler
	journalHandler handler.JournalHandler
	imageHandler   handler.ImageHandler
	noteHandler    handler.NoteHandler
	tagHandler     handler.TagHandler
}

var _ util.ServerLifecycle = (*KohanServerLifecycle)(nil)

// NewKohanServerLifecycle creates a KohanServerLifecycle for testing with explicit handler injection.
func NewKohanServerLifecycle(monitorHandler handler.MonitorHandler,
	journalHandler handler.JournalHandler, imageHandler handler.ImageHandler,
	noteHandler handler.NoteHandler, tagHandler handler.TagHandler) *KohanServerLifecycle {
	return &KohanServerLifecycle{
		monitorHandler: monitorHandler,
		journalHandler: journalHandler,
		imageHandler:   imageHandler,
		noteHandler:    noteHandler,
		tagHandler:     tagHandler,
	}
}

func (s *KohanServerLifecycle) RegisterRoutes(engine *gin.Engine) {
	s.registerMonitorRoutes(engine)
	s.registerJournalRoutes(engine)
}

func (s *KohanServerLifecycle) BeforeStart(_ context.Context)    {}
func (s *KohanServerLifecycle) BeforeShutdown(_ context.Context) {}
func (s *KohanServerLifecycle) AfterShutdown(_ context.Context)  {}

func (s *KohanServerLifecycle) registerMonitorRoutes(engine *gin.Engine) {
	monitor := engine.Group(common.MonitorBase)
	handler.SetupMonitorRoutes(monitor, s.monitorHandler)
}

func (s *KohanServerLifecycle) registerJournalRoutes(engine *gin.Engine) {
	journal := engine.Group(barkat.JournalBase)
	handler.SetupJournalEntryRoutes(journal, s.journalHandler)
	handler.SetupImageRoutes(journal, s.imageHandler)
	handler.SetupNoteRoutes(journal, s.noteHandler)
	handler.SetupTagRoutes(journal, s.tagHandler)
}
