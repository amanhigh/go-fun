package play_fast

import (
	"os"

	"github.com/amanhigh/go-fun/common/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// https://www.sentinelone.com/blog/log-formatting-best-practices-readable/

var _ = FDescribe("Logging", func() {
	var (
		log_file  = "/tmp/log_test"
		file      *os.File
		msgFile   = "I am Testing Logging on File"
		msgStdout = "I am Testing Logging on Stdout"
		field1    = "I am Param1"
		field2    = "I am Param2"
	)
	Context("Logrus", func() {
		var (
			logger *logrus.Logger
			err    error
			rus    = "logrus"
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
				logger.WithField("Logger", rus).Info(msgStdout)
			})

			It("should print fields", func() {
				fields := logrus.Fields{
					"Logger": rus,
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
			err    error
			zapStr = "zap"
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
				logger.Info(msgStdout, zap.String("Logger", zapStr))
			})

			It("should print fields", func() {
				fields := zap.Fields(
					zap.String("Logger", zapStr),
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

})
