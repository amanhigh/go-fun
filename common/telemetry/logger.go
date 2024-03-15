package telemetry

import (
	"os"

	"github.com/amanhigh/go-fun/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
)

/*
*
Processes Context Passed to Logger else ignores.
*/
type ContextLogHook struct {
}

func (h *ContextLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

/*
*
Add RequestId from Context if Contexts is Present else ignore.
*/
func (h *ContextLogHook) Fire(e *logrus.Entry) error {
	if e.Context != nil {
		e.Data["RequestId"] = e.Context.Value(models.XRequestID)
	}
	return nil
}

func InitLogrus(level logrus.Level) {
	//Auto Log RequestId
	//TODO: Move to Zap or Zerologger once they support Context and OTEL.
	logrus.AddHook(&ContextLogHook{})
	logrus.SetLevel(level)

	// Tracing
	logrus.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	)))
}

// InitLogger initializes the logger with the specified level.
// It takes a parameter level of type zerolog.Level.
//
// level - zerolog.DebugLevel (Verbose) to zerolog.ErrorLevel (Limited), or zerolog.FatalLevel (Critical)
func InitLogger(level zerolog.Level) {
	// Level
	zerolog.SetGlobalLevel(level)

	// Formatter
	// HACK: Add Environment Support to Switch Dev vs Prod
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}
