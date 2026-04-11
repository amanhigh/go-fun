package core

import (
	"context"
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/templui/templui/utils"
)

// KohanServerLifecycle implements ServerLifecycle for the Kohan HTTP server.
type KohanServerLifecycle struct {
	OSHandler      handler.OSHandler      `container:"type"`
	JournalHandler handler.JournalHandler `container:"type"`
	ImageHandler   handler.ImageHandler   `container:"type"`
	NoteHandler    handler.NoteHandler    `container:"type"`
	TagHandler     handler.TagHandler     `container:"type"`
	IndexPortal    handler.IndexPortal    `container:"type"`
	JournalPortal  handler.JournalPortal  `container:"type"`
}

type PortalHandlers struct {
	IndexPortal   handler.IndexPortal
	JournalPortal handler.JournalPortal
}

var _ util.ServerLifecycle = (*KohanServerLifecycle)(nil)

// NewKohanServerLifecycle creates a KohanServerLifecycle for testing with explicit handler injection.
func NewKohanServerLifecycle(osHandler handler.OSHandler,
	journalHandler handler.JournalHandler, imageHandler handler.ImageHandler,
	noteHandler handler.NoteHandler, tagHandler handler.TagHandler,
	portalHandlers PortalHandlers) *KohanServerLifecycle {
	return &KohanServerLifecycle{
		OSHandler:      osHandler,
		JournalHandler: journalHandler,
		ImageHandler:   imageHandler,
		NoteHandler:    noteHandler,
		TagHandler:     tagHandler,
		IndexPortal:    portalHandlers.IndexPortal,
		JournalPortal:  portalHandlers.JournalPortal,
	}
}

func (s *KohanServerLifecycle) RegisterRoutes(engine *gin.Engine) {
	s.registerOSRoutes(engine)
	s.registerJournalRoutes(engine)
	s.registerPortalRoutes(engine)
}

func (s *KohanServerLifecycle) RegisterSwagger(engine *gin.Engine) {
	// TODO: Generate swagger docs for Kohan handlers and add annotations
	// make swag-kohan
	// Add Swagger - https://github.com/swaggo/gin-swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *KohanServerLifecycle) BeforeStart(_ context.Context)    {}
func (s *KohanServerLifecycle) BeforeShutdown(_ context.Context) {}
func (s *KohanServerLifecycle) AfterShutdown(_ context.Context)  {}

func (s *KohanServerLifecycle) registerOSRoutes(engine *gin.Engine) {
	os := engine.Group(common.OSBase)
	handler.SetupOSRoutes(os, s.OSHandler)
}

func (s *KohanServerLifecycle) registerJournalRoutes(engine *gin.Engine) {
	journal := engine.Group(barkat.JournalBase)
	handler.SetupJournalRoutes(journal, s.JournalHandler)
	handler.SetupImageRoutes(journal, s.ImageHandler)
	handler.SetupNoteRoutes(journal, s.NoteHandler)
	handler.SetupTagRoutes(journal, s.TagHandler)
}

func (s *KohanServerLifecycle) registerPortalRoutes(engine *gin.Engine) {
	engine.Static("/assets", "assets")
	mux := http.NewServeMux()
	utils.SetupScriptRoutes(mux, true)
	engine.Any("/templui/*filepath", gin.WrapH(mux))
	engine.GET("/", s.IndexPortal.HandleIndex)
	engine.GET("/journal", s.JournalPortal.HandleJournal)
}
