package handler

import "github.com/gin-gonic/gin"

// AlertTickerHandler provides HTTP handlers for Barkat Alert ticker operations.
type AlertTickerHandler interface {
	HandleCreateAlertTicker(c *gin.Context)
	HandleGetAlertTicker(c *gin.Context)
	HandleDeleteAlertTicker(c *gin.Context)
	HandleListAlertTickers(c *gin.Context)
}
