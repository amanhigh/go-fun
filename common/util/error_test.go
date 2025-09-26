package util_test

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
)

var _ = Describe("Error", func() {

	Context("ResponseProcessor", func() {
		var client *resty.Client

		BeforeEach(func() {
			client = resty.New()
		})

		It("should return server error when resty error exists", func() {
			response := &resty.Response{}
			response.Request = client.R()
			restyErr := errors.New("connection failed")
			result := util.ResponseProcessor(response, restyErr)

			Expect(result).To(HaveOccurred())
			Expect(result.Code()).To(Equal(http.StatusInternalServerError))
			Expect(result.Error()).To(Equal("connection failed"))
		})

		It("should return ErrBadRequest for 400 status", func() {
			response := &resty.Response{}
			response.Request = client.R()
			response.RawResponse = &http.Response{StatusCode: http.StatusBadRequest}

			result := util.ResponseProcessor(response, nil)
			Expect(result).To(Equal(common.ErrBadRequest))
		})

		It("should return ErrNotFound for 404 status", func() {
			response := &resty.Response{}
			response.Request = client.R()
			response.RawResponse = &http.Response{StatusCode: http.StatusNotFound}

			result := util.ResponseProcessor(response, nil)
			Expect(result).To(Equal(common.ErrNotFound))
		})

		It("should return ErrNotAuthorized for 401 status", func() {
			response := &resty.Response{}
			response.Request = client.R()
			response.RawResponse = &http.Response{StatusCode: http.StatusUnauthorized}

			result := util.ResponseProcessor(response, nil)
			Expect(result).To(Equal(common.ErrNotAuthorized))
		})

		It("should return ErrNotAuthenticated for 403 status", func() {
			response := &resty.Response{}
			response.Request = client.R()
			response.RawResponse = &http.Response{StatusCode: http.StatusForbidden}

			result := util.ResponseProcessor(response, nil)
			Expect(result).To(Equal(common.ErrNotAuthenticated))
		})

		It("should return ErrEntityExists for 409 status", func() {
			response := &resty.Response{}
			response.Request = client.R()
			response.RawResponse = &http.Response{StatusCode: http.StatusConflict}

			result := util.ResponseProcessor(response, nil)
			Expect(result).To(Equal(common.ErrEntityExists))
		})

		It("should return ErrInternalServerError for 500 status", func() {
			response := &resty.Response{}
			response.Request = client.R()
			response.RawResponse = &http.Response{StatusCode: http.StatusInternalServerError}

			result := util.ResponseProcessor(response, nil)
			Expect(result).To(Equal(common.ErrInternalServerError))
		})

		It("should return nil for unhandled status codes", func() {
			response := &resty.Response{}
			response.Request = client.R()
			response.RawResponse = &http.Response{StatusCode: http.StatusTeapot}

			result := util.ResponseProcessor(response, nil)
			Expect(result).ToNot(HaveOccurred())
		})

		It("should return nil for successful status codes", func() {
			response := &resty.Response{}
			response.Request = client.R()
			response.RawResponse = &http.Response{StatusCode: http.StatusOK}

			result := util.ResponseProcessor(response, nil)
			Expect(result).ToNot(HaveOccurred())
		})
	})

	Context("ProcessValidationError", func() {
		var validate *validator.Validate

		BeforeEach(func() {
			validate = validator.New()
		})

		It("should process validator.ValidationErrors correctly", func() {
			type TestStruct struct {
				Name string `validate:"required"`
				Age  int    `validate:"min=18"`
			}

			testData := TestStruct{Name: "", Age: 15}
			validationErr := validate.Struct(testData)

			result := util.ProcessValidationError(validationErr)

			Expect(result).To(HaveOccurred())
			Expect(result.Code()).To(Equal(http.StatusBadRequest))
			Expect(result.Error()).To(ContainSubstring("'Name'"))
			Expect(result.Error()).To(ContainSubstring("required"))
		})

		It("should return existing HttpError when passed", func() {
			httpErr := common.NewHttpError("validation failed", http.StatusBadRequest)

			result := util.ProcessValidationError(httpErr)
			Expect(result).To(Equal(httpErr))
		})

		It("should handle non-validation errors gracefully", func() {
			genericErr := errors.New("some generic error")

			result := util.ProcessValidationError(genericErr)

			Expect(result).To(HaveOccurred())
			Expect(result.Code()).To(Equal(http.StatusInternalServerError))
			Expect(result.Error()).To(Equal("Invalid validation error format"))
		})

		It("should handle nil validation error", func() {
			result := util.ProcessValidationError(nil)

			Expect(result).To(HaveOccurred())
			Expect(result.Code()).To(Equal(http.StatusInternalServerError))
			Expect(result.Error()).To(Equal("Invalid validation error format"))
		})

		It("should handle valid struct with no validation errors", func() {
			type TestStruct struct {
				Name string `validate:"required"`
			}

			testData := TestStruct{Name: "valid"}
			validationErr := validate.Struct(testData)

			// This should be nil for valid struct
			result := util.ProcessValidationError(validationErr)

			Expect(result).To(HaveOccurred())
			Expect(result.Code()).To(Equal(http.StatusInternalServerError))
		})

		It("should handle multiple validation errors", func() {
			type TestStruct struct {
				Name  string `validate:"required"`
				Age   int    `validate:"min=18"`
				Email string `validate:"email"`
			}

			testData := TestStruct{Name: "", Age: 15, Email: "invalid-email"}
			validationErr := validate.Struct(testData)

			result := util.ProcessValidationError(validationErr)

			Expect(result).To(HaveOccurred())
			Expect(result.Code()).To(Equal(http.StatusBadRequest))
			// Should only return first error due to break in loop
			Expect(result.Error()).To(ContainSubstring("'Name'"))
		})
	})
})
