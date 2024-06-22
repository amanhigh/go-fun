package command

import (
	"bytes"

	"github.com/amanhigh/go-fun/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auto", Label(models.GINKGO_SETUP), func() {
	var (
		actual = new(bytes.Buffer)
	)

	BeforeEach(func() {
		actual.Reset()
		RootCmd.SetOut(actual)
		RootCmd.SetErr(actual)

		logActual.Reset()
	})

	Context("Run Or Focus", func() {
		It("should enable debug", func() {
			RootCmd.SetArgs([]string{"auto", "run-or-focus", "ls", "-d"})

			Expect(RootCmd.Execute()).Should(Succeed())
			Expect(logActual.String()).To(ContainSubstring("Debug Mode Enabled"))
		})

		It("should error with no args", func() {
			RootCmd.SetArgs([]string{"auto", "run-or-focus"})
			Expect(RootCmd.Execute()).ShouldNot(Succeed())
			Expect(actual.String()).To(ContainSubstring("accepts 1 arg"))
		})
	})

})
