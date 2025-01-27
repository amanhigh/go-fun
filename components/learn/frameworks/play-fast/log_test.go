package play_fast

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// https://www.sentinelone.com/blog/log-formatting-best-practices-readable/

// Custom ContextHandler
type ContextHandler struct {
	slog.Handler
}

const (
	requestIDKey  common.ContextKey = "RequestId"
	testRequestID string            = "test-id"
)

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		r.AddAttrs(slog.String(string(requestIDKey), requestID))
	}
	return h.Handler.Handle(ctx, r)
}

var _ = Describe("Logging", func() {
	var (
		log_file  = "/tmp/log_test"
		file      *os.File
		msgFile   = "I am Testing Logging on File"
		msgStdout = "I am Testing Logging on Stdout"
		field1    = "I am Param1"
		field2    = "I am Param2"
		err       error
	)
	Context("Logrus", func() {
		var (
			logger *logrus.Logger
			name   = "logrus"
		)

		It("should build", func() {
			logger = logrus.New()
			Expect(logger).To(Not(BeNil()))
		})

		Context("StdOut", func() {
			BeforeEach(func() {
				logger.Out = os.Stdout
				// Formatter
				logger.SetFormatter(&logrus.TextFormatter{
					FullTimestamp: true,
					ForceColors:   true,
				})
			})

			It("should write log", func() {
				logger.WithField("Logger", name).Info(msgStdout)
			})

			It("should print fields", func() {
				fields := logrus.Fields{
					"Logger": name,
					"Field1": field1,
					"Field2": field2,
				}

				logger.WithFields(fields).Info(msgStdout)
			})
		})

		Context("File", func() {
			BeforeEach(func() {
				file, err = util.OpenOrCreateFile(log_file)
				Expect(err).To(BeNil())
				logger.Out = file
			})

			AfterEach(func() {
				err = os.Remove(log_file)
				Expect(err).To(BeNil())
			})

			It("should write log", func() {
				logger.Info(msgFile)
				lines := util.ReadAllLines(log_file)
				Expect(len(lines)).To(Equal(1))
				Expect(lines[0]).To(ContainSubstring(msgFile))
			})

			It("should write json log", func() {
				logger.SetFormatter(&logrus.JSONFormatter{})
				logger.Info(msgFile)
				lines := util.ReadAllLines(log_file)
				Expect(len(lines)).To(Equal(1))
				Expect(lines[0]).To(ContainSubstring(msgFile))
				Expect(lines[0]).To(ContainSubstring(`"level":"info"`))
			})
		})

		It("should have test logger", func() {
			logger, hook := test.NewNullLogger()
			logger.Info(msgFile)

			Expect(len(hook.AllEntries())).To(Equal(1))
			Expect(hook.LastEntry().Message).To(Equal(msgFile))
			Expect(hook.LastEntry().Level).To(Equal(logrus.InfoLevel))
			hook.Reset()

			Expect(hook.LastEntry()).To(BeNil())
		})
	})

	Context("Zap", func() {
		var (
			logger *zap.Logger
			name   = "zap"
		)

		It("should build", func() {
			logger, err = zap.NewProduction()
			Expect(err).To(BeNil())
			Expect(logger).To(Not(BeNil()))
		})

		Context("StdOut", func() {
			BeforeEach(func() {
				// Formatter
				config := zap.NewDevelopmentConfig()
				config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
				config.EncoderConfig.ConsoleSeparator = " | "

				logger, err = config.Build()
				Expect(err).To(BeNil())
			})

			It("should write log", func() {
				logger.Info(msgStdout, zap.String("Logger", name))
			})

			It("should print fields", func() {
				fields := zap.Fields(
					zap.String("Logger", name),
					zap.String("Field1", field1),
					zap.String("Field2", field2),
				)

				logger.WithOptions(fields).Info(msgStdout)
			})
		})

		Context("File", func() {
			BeforeEach(func() {
				file, err = util.OpenOrCreateFile(log_file)
				Expect(err).To(BeNil())
				logger = zap.New(zapcore.NewCore(
					zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
					zapcore.AddSync(file),
					zap.InfoLevel,
				))
			})

			AfterEach(func() {
				err = os.Remove(log_file)
				Expect(err).To(BeNil())
			})

			It("should write log", func() {
				logger.Info(msgFile)
				lines := util.ReadAllLines(log_file)
				Expect(len(lines)).To(Equal(1))
				Expect(lines[0]).To(ContainSubstring(msgFile))
			})

			It("should write json log", func() {
				logger.Info(msgFile)
				lines := util.ReadAllLines(log_file)
				Expect(len(lines)).To(Equal(1))
				Expect(lines[0]).To(ContainSubstring(msgFile))
				Expect(lines[0]).To(ContainSubstring(`"level":"info"`))
			})
		})

		It("should have test logger", func() {
			core, recorded := observer.New(zap.InfoLevel)
			logger := zap.New(core)

			logger.Info(msgFile)

			Expect(len(recorded.All())).To(Equal(1))
			Expect(recorded.All()[0].Message).To(Equal(msgFile))
			Expect(recorded.All()[0].Level).To(Equal(zap.InfoLevel))
			recorded.TakeAll()

			Expect(len(recorded.All())).To(BeZero())
		})
	})

	// https://github.com/rs/zerolog?tab=readme-ov-file#multiple-log-output
	// https://github.com/rs/zerolog?tab=readme-ov-file#global-settings
	Context("ZeroLog", func() {
		var (
			logger zerolog.Logger
			name   = "zerolog"
		)

		It("should build", func() {
			logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
			Expect(logger).To(Not(BeNil()))
		})

		Context("StdOut", func() {
			BeforeEach(func() {
				// https://github.com/rs/zerolog?tab=readme-ov-file#create-logger-instance-to-manage-different-outputs
				output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
				output.FormatFieldName = func(i any) string {
					return color.YellowString("%s -> ", i)
				}
				output.FormatFieldValue = func(i any) string {
					return color.BlueString("%s", i)
				}
				output.FormatMessage = func(i any) string {
					white := color.New(color.FgWhite, color.Bold)
					return white.Sprintf(" | %s | ", i)
				}
				logger = zerolog.New(output).With().Timestamp().Logger()
			})

			It("should write log", func() {
				logger.Info().Str("Logger", name).Msg(msgStdout)
			})

			It("should print fields", func() {
				subLogger := logger.With().
					Str("Logger", name).
					Str("Field1", field1).
					Str("Field2", field2).
					Logger()
				subLogger.Info().Msg(msgStdout)
			})
		})

		Context("File", func() {
			BeforeEach(func() {
				file, err = util.OpenOrCreateFile(log_file)
				Expect(err).To(BeNil())
				logger = zerolog.New(file).With().Timestamp().Logger()
			})

			AfterEach(func() {
				err = os.Remove(log_file)
				Expect(err).To(BeNil())
			})

			It("should write log", func() {
				logger.Info().Str("Logger", name).Msg(msgFile)
				lines := util.ReadAllLines(log_file)
				Expect(len(lines)).To(Equal(1))
				Expect(lines[0]).To(ContainSubstring(msgFile))
			})

			It("should write json log", func() {
				logger.Info().Str("Logger", name).Msg(msgFile)
				lines := util.ReadAllLines(log_file)
				Expect(len(lines)).To(Equal(1))
				Expect(lines[0]).To(ContainSubstring(msgFile))
				Expect(lines[0]).To(ContainSubstring(`"level":"info"`))
			})

			It("should log context values", func() {
				// Create a child logger for concurrency safety
				child := logger.With().Logger()

				child.UpdateContext(func(c zerolog.Context) zerolog.Context {
					return c.Str("Child", name)
				})
				child.Info().Msg(msgFile)
				lines := util.ReadAllLines(log_file)
				Expect(len(lines)).To(Equal(1))
				Expect(lines[0]).To(ContainSubstring(msgFile))
				Expect(lines[0]).To(ContainSubstring(`"Child":"` + name))
			})

			It("should sample logs", func() {
				// https://github.com/rs/zerolog?tab=readme-ov-file#log-sampling
				sampled := logger.Sample(&zerolog.BasicSampler{N: 10})
				sampled.Info().Msg("will be logged every 10 messages")
			})

			It("should print stacktrace", func() {
				zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
				err := errors.New("test error")
				logger.Error().Stack().Err(err).Msg(msgFile)
				lines := util.ReadAllLines(log_file)
				Expect(len(lines)).To(Equal(1))
				Expect(lines[0]).To(ContainSubstring(msgFile))
				Expect(lines[0]).To(ContainSubstring(`"level":"error"`))
				Expect(lines[0]).To(ContainSubstring(`"stack":`))
			})
		})

		Context("Context Logger", func() {

			It("should have RequestId", func() {
				var buf bytes.Buffer
				logger := zerolog.New(&buf).With().Str(string(requestIDKey), testRequestID).Logger()

				// Creating a context with the logger
				ctx := logger.WithContext(context.Background())

				// Retrieve logger from context and log a message
				zerolog.Ctx(ctx).Info().Msg("Context Test")

				logOutput := buf.String()

				// Check that the log output contains the request ID and the message
				Expect(logOutput).To(ContainSubstring(string(requestIDKey)))
				Expect(logOutput).To(ContainSubstring(testRequestID))
			})

			It("should not be nil without creation", func() {
				logger := zerolog.Ctx(context.Background())
				Expect(logger).ShouldNot(BeNil())
			})
		})

		It("should have test logger", func() {
			var buf bytes.Buffer
			logger := zerolog.New(&buf)

			logger.Info().Msg(msgFile)

			logOutput := buf.String()
			Expect(logOutput).To(ContainSubstring(msgFile))
			Expect(logOutput).To(ContainSubstring(`"level":"info"`))

			buf.Reset()

			Expect(buf.String()).To(Equal(""))
		})
	})

	Context("SLog", func() {
		var (
			logger *slog.Logger
			name   = "slog"
		)

		BeforeEach(func() {
			// Setup logger with Custom Time Formatter
			opts := &slog.HandlerOptions{
				ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						if t, ok := a.Value.Any().(time.Time); ok {
							return slog.String(slog.TimeKey, t.Format("2006-01-02 15:04:05"))
						}
					}
					return a
				},
			}
			logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
		})

		It("should build", func() {
			Expect(logger).To(Not(BeNil()))
		})

		It("should write log", func() {
			logger.Info(msgStdout, "Logger", name)
		})

		It("should print fields", func() {
			logger.Info(msgStdout,
				"Logger", name,
				"Field1", field1,
				"Field2", field2,
			)
		})

		Context("File", func() {
			var file *os.File

			BeforeEach(func() {
				var err error
				file, err = util.OpenOrCreateFile(log_file)
				Expect(err).To(BeNil())

				logger = slog.New(slog.NewJSONHandler(file, nil))
			})

			AfterEach(func() {
				err := os.Remove(log_file)
				Expect(err).To(BeNil())
			})

			It("should write log", func() {
				logger.Info(msgFile, "Logger", name)
				lines := util.ReadAllLines(log_file)
				Expect(len(lines)).To(Equal(1))
				Expect(lines[0]).To(ContainSubstring(msgFile))
			})

			It("should write json log", func() {
				logger.Info(msgFile, "Logger", name)
				lines := util.ReadAllLines(log_file)
				Expect(len(lines)).To(Equal(1))
				Expect(lines[0]).To(ContainSubstring(msgFile))
				Expect(lines[0]).To(ContainSubstring(`"level":"INFO"`))
			})
		})

		Context("Context Logger", func() {

			It("should log context values automatically", func() {
				var buf bytes.Buffer
				baseHandler := slog.NewJSONHandler(&buf, nil)
				contextHandler := ContextHandler{Handler: baseHandler}
				logger := slog.New(contextHandler)

				ctx := context.WithValue(context.Background(), requestIDKey, testRequestID)

				// Now we don't need to explicitly add the RequestId to the log
				logger.InfoContext(ctx, "Context Test")

				logOutput := buf.String()
				Expect(logOutput).To(ContainSubstring(string(requestIDKey)))
				Expect(logOutput).To(ContainSubstring(testRequestID))
			})
		})

		Context("Group Logging", func() {
			var (
				logger *slog.Logger
				buf    bytes.Buffer
			)

			BeforeEach(func() {
				buf.Reset()
				logger = slog.New(slog.NewJSONHandler(&buf, nil))
			})

			It("should log basic group", func() {
				logger.Info("User action",
					slog.Group("user",
						slog.String("id", "123"),
						slog.String("name", "John Doe"),
					),
				)

				logOutput := buf.String()
				Expect(logOutput).To(ContainSubstring(`"user":{"id":"123","name":"John Doe"}`))
			})

			It("should log nested groups", func() {
				logger.Info("Complex data",
					slog.Group("outer",
						slog.String("key1", "value1"),
						slog.Group("inner",
							slog.Int("number", 42),
							slog.Bool("flag", true),
						),
					),
				)

				logOutput := buf.String()
				Expect(logOutput).To(ContainSubstring(`"outer":{"key1":"value1","inner":{"number":42,"flag":true}}`))
			})

			It("should omit empty groups", func() {
				logger.Info("Empty group test",
					slog.Group("empty"),
				)

				logOutput := buf.String()
				Expect(logOutput).To(ContainSubstring(`"msg":"Empty group test"`))
				Expect(logOutput).NotTo(ContainSubstring(`"empty"`))
			})
		})
	})
})
