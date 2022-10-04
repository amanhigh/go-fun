package metrics

import (
	"fmt"
	"time"

	models2 "github.com/amanhigh/go-fun/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

/*
*
RequestId Generator for Gin
*/
func RequestId(c *gin.Context) {
	c.Set(models2.XRequestID, uuid.New())
	c.Next()
}

/* Gin Custom foramtter */
// This Formatter logs requestId over default gin formatter
var GinRequestIdFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %d | %15s |%s %-7s %s %#v | %s\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.BodySize,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.Keys[models2.XRequestID],
		param.ErrorMessage,
	)
}

/*
*
Processes Context Passed to Logger else ignores.
*/
type ContextLogHook struct {
}

func (h *ContextLogHook) Levels() []log.Level {
	return log.AllLevels
}

/*
*
Add RequestId from Context if Contexts is Present else ignore.
*/
func (h *ContextLogHook) Fire(e *log.Entry) error {
	if e.Context != nil {
		e.Data["RequestId"] = e.Context.Value(models2.XRequestID)
	}
	return nil
}
