package util_test

import (
	util2 "github.com/amanhigh/go-fun/common/util"
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
			Expect(util2.NewStringScanner(line)).To(Not(BeNil()))
		})

		It("should read int", func() {
			scanner := util2.NewStringScanner(strconv.Itoa(n))
			Expect(util2.ReadInt(scanner)).To(Equal(n))
		})

		Context("Read Ints", func() {
			var (
				scanner *bufio.Scanner
			)
			BeforeEach(func() {
				scanner = util2.NewStringScanner(line)
			})

			It("should read ints", func() {
				Expect(util2.ReadInts(scanner, n)).To(Equal(line_array))
			})
		})

		It("should read mixed input", func() {
			text := fmt.Sprintf("%v\n%v", n, line)
			scanner := util2.NewStringScanner(text)
			m := util2.ReadInt(scanner)
			Expect(m).To(Equal(n))
			ints := util2.ReadInts(scanner, m)
			Expect(ints).To(Equal(line_array))
		})
	})

	Context("String Scanner", func() {
		var (
			stringList = []string{"Test1", "Test2"}
		)
		It("should read strings", func() {
			scanner := util2.NewStringScanner(strings.Join(stringList, "\n"))
			Expect(util2.ReadStrings(scanner, len(stringList))).To(Equal(stringList))
		})
	})
})
