package play_fast

import (
	"os"

	"github.com/amanhigh/go-fun/common/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
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
		)

		It("should build", func() {
			logger = logrus.New()
			logger.SetFormatter(&logrus.TextFormatter{
				FullTimestamp: true,
				ForceColors:   true,
			})
			Expect(logger).To(Not(BeNil()))
		})

		Context("StdOut", func() {
			BeforeEach(func() {
				logger.Out = os.Stdout
			})

			It("should write log", func() {
				logger.WithField("Logger", "logrus").Info(msgStdout)
			})

			It("should print fields", func() {
				fields := logrus.Fields{
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
})
