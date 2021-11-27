package util_test

import (
	"github.com/amanhigh/go-fun/apps/common/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"strconv"

	"bufio"

	"fmt"

	"strings"
)

var _ = Describe("Io", func() {
	Context("Int Scanner", func() {
		var (
			n          = 7
			line       = "100 100 50 40 40 20 10"
			line_array = []int{100, 100, 50, 40, 40, 20, 10}
		)

		It("should build", func() {
			Expect(util.NewStringScanner(line)).To(Not(BeNil()))
		})

		It("should read int", func() {
			scanner := util.NewStringScanner(strconv.Itoa(n))
			Expect(util.ReadInt(scanner)).To(Equal(n))
		})

		Context("Read Ints", func() {
			var (
				scanner *bufio.Scanner
			)
			BeforeEach(func() {
				scanner = util.NewStringScanner(line)
			})

			It("should read ints", func() {
				Expect(util.ReadInts(scanner, n)).To(Equal(line_array))
			})
		})

		It("should read mixed input", func() {
			text := fmt.Sprintf("%v\n%v", n, line)
			scanner := util.NewStringScanner(text)
			m := util.ReadInt(scanner)
			Expect(m).To(Equal(n))
			ints := util.ReadInts(scanner, m)
			Expect(ints).To(Equal(line_array))
		})
	})

	Context("String Scanner", func() {
		var (
			stringList = []string{"Test1", "Test2"}
		)
		It("should read strings", func() {
			scanner := util.NewStringScanner(strings.Join(stringList, "\n"))
			Expect(util.ReadStrings(scanner, len(stringList))).To(Equal(stringList))
		})
	})
})
