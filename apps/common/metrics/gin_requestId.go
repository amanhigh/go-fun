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
	c.Next()
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

/**
Add RequestId from Context if Contexts is Present else ignore.
*/
func (h *RequestIdHook) Fire(e *log.Entry) error {
	if e.Context != nil {
		e.Data["RequestId"] = e.Context.Value(models.XRequestID)
	}
	return nil
}
