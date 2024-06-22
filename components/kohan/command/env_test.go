package command

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment Command", func() {
	var (
		actual = new(bytes.Buffer)
	)

	BeforeEach(func() {
		// https://nayaktapan37.medium.com/testing-cobra-commands-in-golang-ca1fe4ad6657
		actual.Reset()
		RootCmd.SetOut(actual)
		RootCmd.SetErr(actual)

		logActual.Reset()
	})

	Context("Debug", func() {
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
			Expect(debugCmd.Execute()).ShouldNot(Succeed())
			Expect(actual).To(ContainSubstring("accepts 1 arg"))
		})

		It("should error on invalid args", func() {
			RootCmd.SetArgs([]string{"env", "debug", "invalid"})
			err := debugCmd.Execute()
			Expect(err).ShouldNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("invalid"))
		})
	})

})
