package play_fast

import (
	"bytes"
	"os"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// https://www.sentinelone.com/blog/log-formatting-best-practices-readable/

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
				// HACK: How to Print Request Id from Context.
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
})
