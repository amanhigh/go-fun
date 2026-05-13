package handler

import "github.com/gin-gonic/gin"

// TickerHandler provides HTTP handlers for Barkat ticker operations.
type TickerHandler interface {
	HandleListTickers(c *gin.Context)
	HandleGetTicker(c *gin.Context)
	HandleCreateTicker(c *gin.Context)
	HandleUpdateTicker(c *gin.Context)
	HandlePatchTickerLastOpened(c *gin.Context)
	HandleDeleteTicker(c *gin.Context)
}
