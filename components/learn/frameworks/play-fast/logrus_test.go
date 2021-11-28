package play_fast_test

import (
	"github.com/amanhigh/go-fun/common/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"os"
)

var _ = Describe("Logrus", func() {
	var (
		logger *logrus.Logger
		err    error
	)

	It("should build", func() {
		logger = logrus.New()

		Expect(logger).To(Not(BeNil()))
	})

	Context("File", func() {
		var (
			log_file = "/tmp/logrus_test"
			file     *os.File
		)
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
			msg := "Writing to File"
			logger.Info(msg)
			lines := util.ReadAllLines(log_file)
			Expect(len(lines)).To(Equal(1))
			Expect(lines[0]).To(ContainSubstring(msg))
		})
	})

	Context("StdOut", func() {
		BeforeEach(func() {
			logger.Out = os.Stdout
		})

		It("should write log", func() {
			logger.Info("Writing to Stdout")
		})
	})

})
