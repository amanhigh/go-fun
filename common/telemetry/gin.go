package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/amanhigh/go-fun/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

/*
*
RequestId Generator for Gin
*/
func RequestId(c *gin.Context) {
	// Generate UUID
	uuid := uuid.New()
	c.Set(models.XRequestID, uuid)

	// Logger with RequestId
	idLogger := log.With().Str(models.XRequestID, uuid.String()).Logger()

	//Add UUID & Logger to Request Context as well
	ctx := context.WithValue(c.Request.Context(), models.XRequestID, uuid)
	ctx = idLogger.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

/* Gin Custom foramtter */
// This Formatter logs requestId over default gin formatter
var GinRequestIdFormatter = func(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[GIN] %v | %3d | %5d | %d | %15s | %s | %s ",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		param.StatusCode,
		param.Latency.Microseconds(),
		param.BodySize,
		param.ClientIP,
		param.Keys[models.XRequestID],
		param.Method,
	)
}
