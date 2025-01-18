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
	idLogger := log.With().Str("RequestId", uuid.String()).Logger()

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

	// XXX: Implement CLF Field Authentication $remote_user: - if no authentication is used.
	return fmt.Sprintf("[GIN] %s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\" \"%s\" \"%d\"\n",
		param.ClientIP,
		param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.BodySize,
		param.Request.Referer(),
		param.Request.UserAgent(),
		param.Keys[models.XRequestID],
		param.Latency.Microseconds(),
	)
}
