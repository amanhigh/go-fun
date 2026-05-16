package handler

import (
	"net/http"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
)

// TickerHandler provides HTTP handlers for Barkat ticker operations.
type TickerHandler interface {
	HandleListTickers(c *gin.Context)
	HandleGetTicker(c *gin.Context)
	HandleCreateTicker(c *gin.Context)
	HandleUpdateTicker(c *gin.Context)
	HandlePatchTickerLastOpened(c *gin.Context)
	HandleDeleteTicker(c *gin.Context)
}

// TickerHandlerImpl implements the TickerHandler interface.
type TickerHandlerImpl struct {
	tickerManager manager.BarkatTickerManager
}

var _ TickerHandler = (*TickerHandlerImpl)(nil)

// NewTickerHandler creates a new TickerHandlerImpl.
func NewTickerHandler(tickerManager manager.BarkatTickerManager) *TickerHandlerImpl {
	return &TickerHandlerImpl{tickerManager: tickerManager}
}

// ---- Handlers ----

func (h *TickerHandlerImpl) HandleListTickers(c *gin.Context) {
	query := barkat.NewTickerQuery()

	if bindErr := c.ShouldBindQuery(&query); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	tickerList, httpErr := h.tickerManager.ListTickers(c.Request.Context(), query)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(tickerList))
}

func (h *TickerHandlerImpl) HandleGetTicker(c *gin.Context) {
	var path barkat.TickerPath

	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	ticker, httpErr := h.tickerManager.GetTicker(c.Request.Context(), path.Ticker)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(ticker))
}

func (h *TickerHandlerImpl) HandleCreateTicker(c *gin.Context) {
	var ticker barkat.Ticker
	if bindErr := c.ShouldBindJSON(&ticker); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	httpErr := h.tickerManager.CreateTicker(c.Request.Context(), &ticker)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusCreated, common.NewEnvelope(ticker))
}

func (h *TickerHandlerImpl) HandleUpdateTicker(c *gin.Context) {
	var path barkat.TickerPath

	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	var updateReq = barkat.TickerUpdateRequest{Ticker: path.Ticker}
	if bindErr := c.ShouldBindJSON(&updateReq); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	ticker, httpErr := h.tickerManager.UpdateTicker(c.Request.Context(), path.Ticker, updateReq)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(ticker))
}

func (h *TickerHandlerImpl) HandlePatchTickerLastOpened(c *gin.Context) {
	var path barkat.TickerPath
	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	var update barkat.TickerLastOpenedUpdate
	if bindErr := c.ShouldBindJSON(&update); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	ticker, httpErr := h.tickerManager.PatchTickerLastOpened(c.Request.Context(), path.Ticker, update)
	if httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusOK, common.NewEnvelope(ticker))
}

func (h *TickerHandlerImpl) HandleDeleteTicker(c *gin.Context) {
	var path barkat.TickerPath

	if bindErr := c.ShouldBindUri(&path); bindErr != nil {
		httpErr := util.ProcessValidationError(bindErr)
		c.JSON(httpErr.Code(), httpErr)
		return
	}

	if httpErr := h.tickerManager.DeleteTicker(c.Request.Context(), path.Ticker); httpErr != nil {
		c.JSON(httpErr.Code(), httpErr)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
