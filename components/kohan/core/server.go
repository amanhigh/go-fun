package core

import (
	"fmt"
	"net/http"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type MonitorServer struct {
	mux         *gin.Engine
	capturePath string
	autoManager manager.AutoManagerInterface
}

func NewMonitorServer(capturePath string, autoManager manager.AutoManagerInterface) *MonitorServer {
	server := &MonitorServer{
		mux:         gin.Default(),
		capturePath: capturePath,
		autoManager: autoManager,
	}

	// Register Routes
	server.mux.GET("/v1/ticker/:ticker/record", server.HandleRecordTicker)
	server.mux.GET("/v1/clip/", server.HandleReadClip)
	server.mux.POST("/v1/submap/:action", server.HandleSubmapControl)

	return server
}

func (s *MonitorServer) Start(port int) (err error) {
	log.Info().Int("port", port).Msg("Starting Monitor Server")
	err = s.mux.Run(fmt.Sprintf(":%d", port))
	return
}

func (s *MonitorServer) HandleReadClip(ctx *gin.Context) {
	text, err := tools.ClipPaste()
	if err == nil {
		ctx.JSON(http.StatusOK, text)
	} else {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
}

func (s *MonitorServer) HandleRecordTicker(ctx *gin.Context) {
	ticker := ctx.Param("ticker")
	if err := s.autoManager.RecordTicker(ctx, ticker, s.capturePath); err == nil {
		ctx.JSON(http.StatusOK, "Success")
	} else {
		log.Error().Str("Ticker", ticker).Err(err).Msg("Record Ticker Failed")
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
}

func (s *MonitorServer) HandleSubmapControl(ctx *gin.Context) {
	action := ctx.Param("action")

	var request struct {
		Submap string `json:"submap"`
		Ticker string `json:"ticker,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&request); err == nil {
		switch action {
		case "enable":
			err = tools.HyperDispatch("submap " + request.Submap)
		case "disable":
			err = tools.HyperDispatch("submap reset")
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action. Use 'enable' or 'disable'"})
			return
		}

		if err == nil {
			log.Info().Str("Action", action).Str("Submap", request.Submap).Msg("Submap Control")
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "action": action, "submap": request.Submap})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
