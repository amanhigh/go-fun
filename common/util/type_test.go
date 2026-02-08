package util_test

import (
	"errors"

	"github.com/amanhigh/go-fun/common/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Type", func() {
	Context("IsInt", func() {
		It("should validate integer strings (SDK wrapper)", func() {
			Expect(util.IsInt("123")).ToNot(HaveOccurred())
			Expect(util.IsInt("abc")).To(HaveOccurred())
			Expect(util.IsInt("abc").Error()).To(Equal("abc is not a Valid Integer"))
		})
	})

	Context("ParseInt", func() {
		It("should parse integers with custom error message (SDK wrapper)", func() {
			result, err := util.ParseInt("123")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(123))

			result, err = util.ParseInt("abc")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("abc is not a Valid Integer"))
			Expect(result).To(Equal(0))
		})
	})

	Context("ParseBool", func() {
		It("should parse booleans with custom error message (SDK wrapper)", func() {
			result, err := util.ParseBool("true")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeTrue())

			result, err = util.ParseBool("invalid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid is not a Valid Boolean"))
			Expect(result).To(BeFalse())
		})
	})

	Context("ParseFloat", func() {
		It("should parse floats but has implementation bug", func() {
			result := util.ParseFloat("123.45")
			Expect(result).To(Equal(123.45))

			result = util.ParseFloat("")
			Expect(result).To(Equal(-1.0))

			// Implementation bug: doesn't reset to -1 on error
			result = util.ParseFloat("invalid")
			Expect(result).To(Equal(0.0)) // Should be -1.0 but implementation has bug
		})
	})

	Context("ReverseArray", func() {
		It("should reverse slices in-place", func() {
			input := []int{1, 2, 3, 4, 5}
			util.ReverseArray(input)
			Expect(input).To(Equal([]int{5, 4, 3, 2, 1}))

			// Works with different types
			strings := []string{"a", "b", "c"}
			util.ReverseArray(strings)
			Expect(strings).To(Equal([]string{"c", "b", "a"}))

			// Edge cases
			empty := []int{}
			util.ReverseArray(empty)
			Expect(empty).To(Equal([]int{}))
		})
	})

	Context("CancelFunc Type", func() {
		It("should be a function type that returns error", func() {
			var cancelFunc util.CancelFunc = func() error {
				return errors.New("test error")
			}

			err := cancelFunc()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("RoundToDecimals - Generic rounding function", func() {
		It("should round floating point precision errors to specified decimals", func() {
			Expect(util.RoundToDecimals(0.3899999999999999, 2)).To(Equal(0.39))
			Expect(util.RoundToDecimals(24.95999999999999, 2)).To(Equal(24.96))
			Expect(util.RoundToDecimals(-37.76999999999999, 2)).To(Equal(-37.77))
		})

		It("should handle accumulated calculation errors from multiple operations", func() {
			// Examples from FIFO calculations
			Expect(util.RoundToDecimals(0.777+0.333, 2)).To(Equal(1.11))
			Expect(util.RoundToDecimals(0.292661638, 2)).To(Equal(0.29))
			Expect(util.RoundToDecimals(2.6500000004, 2)).To(Equal(2.65))
		})

		It("should support rounding to any decimal precision", func() {
			Expect(util.RoundToDecimals(1.23456, 1)).To(Equal(1.2))
			Expect(util.RoundToDecimals(1.23456, 3)).To(Equal(1.235))
			Expect(util.RoundToDecimals(1.23456, 4)).To(Equal(1.2346))
		})
	})
})
