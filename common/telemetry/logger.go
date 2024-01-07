package telemetry

import (
	"github.com/amanhigh/go-fun/models"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
)

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
		e.Data["RequestId"] = e.Context.Value(models.HEADER_REQUESTID)
	}
	return nil
}

func InitLogger(level log.Level) {
	//Auto Log RequestId
	log.AddHook(&ContextLogHook{})
	log.SetLevel(level)

	// Tracing
	logrus.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	)))
}
