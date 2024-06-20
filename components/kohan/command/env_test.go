package command

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment Command", func() {
	var (
		actual *bytes.Buffer
	)

	BeforeEach(func() {
		actual = new(bytes.Buffer)
	})

	Context("Debug", func() {
		BeforeEach(func() {
			// https://nayaktapan37.medium.com/testing-cobra-commands-in-golang-ca1fe4ad6657
			debugCmd.SetOut(actual)
			debugCmd.SetErr(actual)

			RootCmd.SetArgs([]string{"env", "debug"})
		})

		It("should enable debug mode", func() {
			// FIXME: #B Flags Should Work
			Expect(debugCmd.Execute()).Should(Succeed())
			// Expect(actual.String()).To(ContainSubstring("Enabling Debug Mode"))
		})
	})

})
