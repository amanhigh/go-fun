package util_test

import (
	"errors"

	"github.com/amanhigh/go-fun/common/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate", func() {
	Context("Verify", func() {
		Context("Single Error", func() {
			It("should return nil when no error provided", func() {
				err := util.Verify()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return nil when single nil error provided", func() {
				err := util.Verify(nil)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return the error when single error provided", func() {
				testErr := errors.New("test error")
				err := util.Verify(testErr)
				Expect(err).To(Equal(testErr))
			})
		})

		Context("Multiple Errors", func() {
			It("should return nil when all errors are nil", func() {
				err := util.Verify(nil, nil, nil)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return first non-nil error", func() {
				firstErr := errors.New("first error")
				secondErr := errors.New("second error")
				err := util.Verify(nil, firstErr, secondErr)
				Expect(err).To(Equal(firstErr))
			})

			It("should return first error even if others are nil", func() {
				testErr := errors.New("test error")
				err := util.Verify(testErr, nil, nil)
				Expect(err).To(Equal(testErr))
			})

			It("should skip nil errors and return first non-nil", func() {
				firstErr := errors.New("first error")
				secondErr := errors.New("second error")
				err := util.Verify(nil, nil, firstErr, secondErr)
				Expect(err).To(Equal(firstErr))
			})
		})

		Context("Edge Cases", func() {
			It("should handle mixed error types", func() {
				customErr := &customError{msg: "custom error"}
				standardErr := errors.New("standard error")
				err := util.Verify(nil, customErr, standardErr)
				Expect(err).To(Equal(customErr))
			})
		})
	})

	Context("ValidateEnumArg", func() {
		var validEnum []string

		BeforeEach(func() {
			validEnum = []string{"option1", "option2", "option3"}
		})

		Context("Valid Arguments", func() {
			It("should return nil for valid enum value", func() {
				err := util.ValidateEnumArg("option1", validEnum)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return nil for each valid enum value", func() {
				for _, option := range validEnum {
					err := util.ValidateEnumArg(option, validEnum)
					Expect(err).ToNot(HaveOccurred())
				}
			})

			It("should be case sensitive", func() {
				err := util.ValidateEnumArg("Option1", validEnum)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Option1 is not a Valid Argument"))
			})
		})

		Context("Invalid Arguments", func() {
			It("should return error for invalid enum value", func() {
				err := util.ValidateEnumArg("invalid", validEnum)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid is not a Valid Argument. Valid Values: [option1 option2 option3]"))
			})

			It("should return error for empty string", func() {
				err := util.ValidateEnumArg("", validEnum)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is not a Valid Argument"))
			})

			It("should return error for whitespace string", func() {
				err := util.ValidateEnumArg(" ", validEnum)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is not a Valid Argument"))
			})

			It("should return error for partial match", func() {
				err := util.ValidateEnumArg("option", validEnum)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("option is not a Valid Argument"))
			})
		})

		Context("Edge Cases", func() {
			It("should handle empty enum slice", func() {
				emptyEnum := []string{}
				err := util.ValidateEnumArg("anything", emptyEnum)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("anything is not a Valid Argument. Valid Values: []"))
			})

			It("should handle enum with special characters", func() {
				specialEnum := []string{"option@1", "ñáéíóú"}
				err := util.ValidateEnumArg("ñáéíóú", specialEnum)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Error Message Format", func() {
			It("should include argument and valid values in error message", func() {
				invalidArg := "invalidArg"
				err := util.ValidateEnumArg(invalidArg, validEnum)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(invalidArg))
				Expect(err.Error()).To(ContainSubstring("Valid Values:"))
				Expect(err.Error()).To(ContainSubstring("[option1 option2 option3]"))
			})

			It("should format single enum value correctly", func() {
				singleEnum := []string{"only_option"}
				err := util.ValidateEnumArg("wrong", singleEnum)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("wrong is not a Valid Argument. Valid Values: [only_option]"))
			})
		})
	})
})

// Custom error type for testing mixed error types
type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}
