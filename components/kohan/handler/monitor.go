package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/tools"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// MonitorHandler provides HTTP handlers for system monitoring operations.
type MonitorHandler interface {
	// HandleReadClip handles GET /v1/clip/
	HandleReadClip(ctx *gin.Context)
	// HandleRecordTicker handles GET /v1/ticker/:ticker/record
	HandleRecordTicker(ctx *gin.Context)
	// HandleSubmapControl handles POST /v1/submap/:action
	HandleSubmapControl(ctx *gin.Context)
}

type MonitorHandlerImpl struct {
	capturePath string
	autoManager manager.AutoManagerInterface
}

var _ MonitorHandler = (*MonitorHandlerImpl)(nil)

// NewMonitorHandler creates a new MonitorHandler.
func NewMonitorHandler(capturePath string, autoManager manager.AutoManagerInterface) *MonitorHandlerImpl {
	return &MonitorHandlerImpl{capturePath: capturePath, autoManager: autoManager}
}

// HandleReadClip handles GET /v1/clip/
func (h *MonitorHandlerImpl) HandleReadClip(ctx *gin.Context) {
	text, err := tools.ClipPaste()
	if err == nil {
		ctx.JSON(http.StatusOK, text)
	} else {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
}

// HandleRecordTicker handles GET /v1/ticker/:ticker/record
func (h *MonitorHandlerImpl) HandleRecordTicker(ctx *gin.Context) {
	ticker := ctx.Param("ticker")
	if err := h.autoManager.RecordTicker(ctx, ticker, h.capturePath); err == nil {
		ctx.JSON(http.StatusOK, "Success")
	} else {
		log.Error().Str("Ticker", ticker).Err(err).Msg("Record Ticker Failed")
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
}

// HandleSubmapControl handles POST /v1/submap/:action
func (h *MonitorHandlerImpl) HandleSubmapControl(ctx *gin.Context) {
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
