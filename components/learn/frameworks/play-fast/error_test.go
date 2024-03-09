package play_fast

import (
	"errors"
	"fmt"

	perr "github.com/pkg/errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rotisserie/eris"
)

// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
// https://blog.dharnitski.com/2019/09/09/go-errors-are-not-pkg-errors/
var _ = Describe("Error", func() {
	var (
		err error
	)

	BeforeEach(func() {
		err = errors.New("error:0")
	})

	It("should be created", func() {
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("error:0"))
	})

	It("should join", func() {
		err2 := fmt.Errorf("error:%d", 1)
		err = errors.Join(err, nil, err2)

		Expect(err).ToNot(BeNil())
		Expect(err.Error()).Should(Equal("error:0\nerror:1"))
	})

	Context("Wrap", func() {
		BeforeEach(func() {
			err = fmt.Errorf("error:%d %w", 1, err)
		})

		It("should be created", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).Should(Equal("error:1 error:0"))
		})

		It("should have cause", func() {
			Expect(err).ToNot(BeNil())
			Expect(errors.Unwrap(err)).Should(MatchError("error:0"))
		})
	})

	Context("pkg/errors", func() {
		BeforeEach(func() {
			err = perr.Wrapf(err, "error:%d", 1)
		})

		It("should be created", func() {
			Expect(err).ToNot(BeNil())
			Expect(err).Should(MatchError("error:1: error:0"))
		})

		It("should print stacktrace", func() {
			trace := fmt.Sprintf("%+v", err)
			Expect(trace).Should(ContainSubstring("error_test.go"))
		})

		Context("Unwrap", func() {
			BeforeEach(func() {
				err = perr.Cause(err)
			})

			It("should work", func() {
				Expect(err).ToNot(BeNil())
				Expect(err).Should(MatchError("error:0"))
			})
		})
	})

	Context("Eris", func() {
		BeforeEach(func() {
			err = eris.New("error:0")
			err = eris.Wrapf(err, "error:%d", 1)
		})

		It("should be created", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).Should(Equal("error:1: error:0"))
		})

		It("should print stacktrace", func() {
			trace := fmt.Sprintf("%+v", err)
			Expect(eris.ToString(err, true)).Should(Equal(trace))
			Expect(trace).Should(ContainSubstring("error_test.go"))
		})

		It("should print custom stacktrace", func() {
			// format the error to a string with custom separators
			trace := eris.ToCustomString(err, eris.StringFormat{
				Options: eris.FormatOptions{
					WithTrace: true, // flag that enables stack trace output
				},
				MsgStackSep:  "\n",  // separator between error messages and stack frame data
				PreStackSep:  "\t",  // separator at the beginning of each stack frame
				StackElemSep: " | ", // separator between elements of each stack frame
				ErrorSep:     "\n",  // separator between each error in the chain
			})
			Expect(trace).Should(ContainSubstring("error_test.go"))
		})

		It("should generate json", func() {
			json := eris.ToJSON(err, true)
			Expect(json).ShouldNot(BeEmpty())
		})

		It("should have Cause", func() {
			Expect(eris.Cause(err)).Should(MatchError("error:0"))
		})
	})

})
