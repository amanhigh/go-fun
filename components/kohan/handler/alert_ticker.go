package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
)

// AlertTickerHandler provides HTTP handlers for Barkat Alert ticker operations.
type AlertTickerHandler interface {
	HandleCreateAlertTicker(c *gin.Context)
	HandleGetAlertTicker(c *gin.Context)
	HandleDeleteAlertTicker(c *gin.Context)
	HandleListAlertTickers(c *gin.Context)
}

// AlertTickerHandlerImpl implements the AlertTickerHandler interface.
type AlertTickerHandlerImpl struct {
	alertTickerManager manager.AlertTickerManager
}

var _ AlertTickerHandler = (*AlertTickerHandlerImpl)(nil)

// NewAlertTickerHandler creates a new AlertTickerHandlerImpl.
func NewAlertTickerHandler(alertTickerManager manager.AlertTickerManager) *AlertTickerHandlerImpl {
	return &AlertTickerHandlerImpl{alertTickerManager: alertTickerManager}
}

// ---- Handlers ----

func (h *AlertTickerHandlerImpl) HandleCreateAlertTicker(c *gin.Context) {
	var path barkat.TickerPath
	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	var alert barkat.AlertTicker
	if bindErr := c.ShouldBindJSON(&alert); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	response, httpErr := h.alertTickerManager.CreateAlertTicker(c.Request.Context(), path.Ticker, &alert)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, common.NewEnvelope(response))
}

func (h *AlertTickerHandlerImpl) HandleGetAlertTicker(c *gin.Context) {
	var path barkat.AlertTickerPath
	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	response, httpErr := h.alertTickerManager.GetAlertTicker(c.Request.Context(), path.Symbol)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(response))
}

func (h *AlertTickerHandlerImpl) HandleDeleteAlertTicker(c *gin.Context) {
	var path barkat.AlertTickerPath
	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	if httpErr := h.alertTickerManager.DeleteAlertTicker(c.Request.Context(), path.Symbol); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *AlertTickerHandlerImpl) HandleListAlertTickers(c *gin.Context) {
	query := barkat.NewAlertTickerQuery()

	if bindErr := c.ShouldBindQuery(&query); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	response, httpErr := h.alertTickerManager.ListAlertTickers(c.Request.Context(), query)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(response))
}
