package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
)

// PriceAlertHandler provides HTTP handlers for Barkat price alert operations.
type PriceAlertHandler interface {
	HandleReplacePriceAlerts(c *gin.Context)
	HandleCreatePendingPriceAlert(c *gin.Context)
	HandleDeletePriceAlert(c *gin.Context)
	HandleListPriceAlerts(c *gin.Context)
}

type PriceAlertHandlerImpl struct {
	priceAlertManager manager.PriceAlertManager
}

var _ PriceAlertHandler = (*PriceAlertHandlerImpl)(nil)

// NewPriceAlertHandler creates a new PriceAlertHandlerImpl.
func NewPriceAlertHandler(priceAlertManager manager.PriceAlertManager) *PriceAlertHandlerImpl {
	return &PriceAlertHandlerImpl{priceAlertManager: priceAlertManager}
}

func (h *PriceAlertHandlerImpl) HandleReplacePriceAlerts(c *gin.Context) {
	var request barkat.PriceAlertReplaceRequest
	if bindErr := c.ShouldBindJSON(&request); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	response, httpErr := h.priceAlertManager.ReplacePriceAlerts(c.Request.Context(), request)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(response))
}

func (h *PriceAlertHandlerImpl) HandleCreatePendingPriceAlert(c *gin.Context) {
	var path barkat.TickerPath
	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	var request barkat.PendingPriceAlertRequest
	if bindErr := c.ShouldBindJSON(&request); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	response, httpErr := h.priceAlertManager.CreatePendingPriceAlert(c.Request.Context(), path.Ticker, request)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, common.NewEnvelope(response))
}

func (h *PriceAlertHandlerImpl) HandleDeletePriceAlert(c *gin.Context) {
	var path barkat.PriceAlertPath
	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	if httpErr := h.priceAlertManager.DeletePriceAlert(c.Request.Context(), path.AlertID); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *PriceAlertHandlerImpl) HandleListPriceAlerts(c *gin.Context) {
	var query barkat.PriceAlertQuery
	if bindErr := c.ShouldBindQuery(&query); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	response, httpErr := h.priceAlertManager.ListPriceAlerts(c.Request.Context(), query)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(response))
}
