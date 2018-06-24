package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"strconv"

	"bufio"

	"fmt"

	. "github.com/amanhigh/go-fun/util"
)

var _ = Describe("Io", func() {
	Context("String Scanner", func() {
		var (
			n          = 7
			line       = "100 100 50 40 40 20 10"
			line_array = []int{100, 100, 50, 40, 40, 20, 10}
		)

		It("should build", func() {
			Expect(NewStringScanner(line)).To(Not(BeNil()))
		})

		It("should read int", func() {
			scanner := NewStringScanner(strconv.Itoa(n))
			Expect(ReadInt(scanner)).To(Equal(n))
		})

		Context("Read Ints", func() {
			var (
				scanner *bufio.Scanner
			)
			BeforeEach(func() {
				scanner = NewStringScanner(line)
			})

			It("should read ints", func() {
				Expect(ReadInts(scanner, n)).To(Equal(line_array))
			})
		})

		It("should read mixed input", func() {
			text := fmt.Sprintf("%v\n%v", n, line)
			scanner := NewStringScanner(text)
			m := ReadInt(scanner)
			Expect(m).To(Equal(n))
			ints := ReadInts(scanner, m)
			Expect(ints).To(Equal(line_array))
		})
	})
})
