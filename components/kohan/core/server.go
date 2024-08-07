package core

import (
	"fmt"
	"net/http"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type MonitorServer struct {
	mux         *gin.Engine
	capturePath string
}

func NewMonitorServer(capturePath string) (server *MonitorServer) {
	// Build Server
	server = &MonitorServer{
		mux:         gin.Default(),
		capturePath: capturePath,
	}

	// Register Routes
	server.mux.GET("/v1/ticker/:ticker/record", server.HandleRecordTicker)
	server.mux.GET("/v1/clip/", server.HandleReadClip)

	return
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
	return
}

func (s *MonitorServer) HandleRecordTicker(ctx *gin.Context) {
	ticker := ctx.Param("ticker")
	if err := RecordTicker(ticker, s.capturePath); err == nil {
		ctx.JSON(http.StatusOK, "Success")

	} else {
		log.Error().Str("Ticker", ticker).Err(err).Msg("Record Ticker Failed")
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
}
