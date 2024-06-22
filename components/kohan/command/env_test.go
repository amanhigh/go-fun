package command

import (
	"bytes"

	"github.com/amanhigh/go-fun/common/telemetry"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment Command", Ordered, func() {
	var (
		actual    = new(bytes.Buffer)
		logActual = new(bytes.Buffer)
	)

	BeforeAll(func() {
		telemetry.InitTestLogger(logActual)
	})

	BeforeEach(func() {
		actual.Reset()
		logActual.Reset()
	})

	Context("Debug", func() {
		BeforeEach(func() {
			// https://nayaktapan37.medium.com/testing-cobra-commands-in-golang-ca1fe4ad6657
			debugCmd.SetOut(actual)
			debugCmd.SetErr(actual)
		})

		It("should enable", func() {
			RootCmd.SetArgs([]string{"env", "debug", "true"})
			Expect(debugCmd.Execute()).Should(Succeed())
			Expect(logActual.String()).To(ContainSubstring("Enabling Debug Mode"))
		})

		It("should disable", func() {
			RootCmd.SetArgs([]string{"env", "debug", "false"})
			Expect(debugCmd.Execute()).Should(Succeed())
			Expect(logActual.String()).To(ContainSubstring("Disabling Debug Mode"))
		})

		It("should error on no args", func() {
			RootCmd.SetArgs([]string{"env", "debug"})
			err := debugCmd.Execute()
			Expect(err).ShouldNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("accepts 1 arg"))
		})

		It("should error on invalid args", func() {
			RootCmd.SetArgs([]string{"env", "debug", "invalid"})
			err := debugCmd.Execute()
			Expect(err).ShouldNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("invalid"))
		})
	})

})
