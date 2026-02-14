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
	*util.BaseHTTPServer
}

// NewKohanServer creates a KohanServer with monitor and journal routes.
func NewKohanServer(port int, capturePath string, autoManager manager.AutoManagerInterface, journalHandler *handler.JournalHandler, shutdown util.Shutdown) *KohanServer {
	server := &KohanServer{
		BaseHTTPServer: util.NewBaseHTTPServer("kohan", port, shutdown),
	}

	// FIXME: All Handlers should be interface injected by DI
	// FIXME: Route should call handler not Manager Directly create Handler for Auto Manager
	// HACK: Each Handler sholud be interface injected by DI not Struct Rename it to JournalHandlerImpl
	server.BaseHTTPServer.RegisterRoutes = func(engine *gin.Engine) {
		// FIXME: Register Routes should be named function here not lampda.
		registerMonitorRoutes(engine, capturePath, autoManager)
		registerJournalRoutes(engine, journalHandler)
	}

	return server
}

func registerMonitorRoutes(engine *gin.Engine, capturePath string, autoManager manager.AutoManagerInterface) {
	monitorHandler := handler.NewMonitorHandler(capturePath, autoManager)
	engine.GET("/v1/ticker/:ticker/record", monitorHandler.HandleRecordTicker)
	engine.GET("/v1/clip/", monitorHandler.HandleReadClip)
	engine.POST("/v1/submap/:action", monitorHandler.HandleSubmapControl)
}

func registerJournalRoutes(engine *gin.Engine, journalHandler *handler.JournalHandler) {
	if journalHandler != nil {
		v1 := engine.Group("/api/v1")
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
}
