package metrics

import (
	"github.com/amanhigh/go-fun/apps/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

/**
RequestId Generator for Gin
*/
func RequestId(c *gin.Context) {
	c.Set(models.XRequestID, uuid.New())
}

/**
RequestId Hook to Auto logs RequestId
from Context
*/
type RequestIdHook struct {
}

func (h *RequestIdHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *RequestIdHook) Fire(e *log.Entry) error {
	var value interface{}
	if e.Context == nil {
		value = "ContextMissing"
	} else {
		value = e.Context.Value(models.XRequestID)
	}
	e.Data["RequestId"] = value
	return nil
}
