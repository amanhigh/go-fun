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
Processes Context Passed to Logger else ignores.
*/
type ContextLogHook struct {
}

func (h *ContextLogHook) Levels() []log.Level {
	return log.AllLevels
}

/**
Add RequestId from Context if Contexts is Present else ignore.
*/
func (h *ContextLogHook) Fire(e *log.Entry) error {
	if e.Context != nil {
		e.Data["RequestId"] = e.Context.Value(models.XRequestID)
	}
	return nil
}
