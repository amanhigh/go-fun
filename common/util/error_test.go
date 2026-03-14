package util_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
)

var _ = Describe("Error", func() {

	// ============================================================================
	// SECTION 1: HTTP Response Processing (ResponseProcessor)
	// ============================================================================
	Context("ResponseProcessor", func() {
		var client *resty.Client

		BeforeEach(func() {
			client = resty.New()
		})

		Context("Client Errors", func() {
			It("should return server error when resty error exists", func() {
				response := &resty.Response{}
				response.Request = client.R()
				restyErr := errors.New("connection failed")

				result := util.ResponseProcessor(response, restyErr)

				Expect(result).To(HaveOccurred())
				Expect(result.Code()).To(Equal(http.StatusInternalServerError))
				Expect(result.Error()).To(Equal("connection failed"))
			})
		})

		Context("Status Code Mapping", func() {
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

			It("should return ErrPayloadTooLarge for 413 status", func() {
				response := &resty.Response{}
				response.Request = client.R()
				response.RawResponse = &http.Response{StatusCode: http.StatusRequestEntityTooLarge}

				result := util.ResponseProcessor(response, nil)
				Expect(result).To(Equal(common.ErrPayloadTooLarge))
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
	})

	// ============================================================================
	// SECTION 2: Validation Error Processing (ProcessValidationError)
	// ============================================================================
	Context("ProcessValidationError", func() {
		var validate *validator.Validate

		BeforeEach(func() {
			validate = validator.New()
		})

		// --------------------------------------------------------------------
		// 2.1: Validator Library Errors (binding tags)
		// --------------------------------------------------------------------
		Context("Validator Library Errors", func() {

			Context("Required Tag", func() {
				It("should handle required field violation", func() {
					type TestStruct struct {
						Name string `validate:"required"`
					}
					testData := TestStruct{Name: ""}
					validationErr := validate.Struct(testData)

					result := util.ProcessValidationError(validationErr)

					// First assert that it's a FieldHttpError
					fieldErr, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErr.Error()).To(ContainSubstring("Name"))
					Expect(fieldErr.Error()).To(ContainSubstring("required"))
					Expect(fieldErr.Field()).To(Equal("Name"))
				})
			})

			Context("Min/Max Tags", func() {
				It("should handle min violation for integers", func() {
					type TestStruct struct {
						Age int `validate:"min=18"`
					}
					testData := TestStruct{Age: 15}
					validationErr := validate.Struct(testData)

					result := util.ProcessValidationError(validationErr)

					// First assert that it's a FieldHttpError
					fieldErr, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErr.Error()).To(ContainSubstring("Age"))
					Expect(fieldErr.Error()).To(ContainSubstring("min"))
					Expect(fieldErr.Field()).To(Equal("Age"))
				})

				It("should handle max violation for strings", func() {
					type TestStruct struct {
						Ticker string `validate:"max=10"`
					}
					testData := TestStruct{Ticker: "VERYLONGTICKER"}
					validationErr := validate.Struct(testData)

					result := util.ProcessValidationError(validationErr)

					// First assert that it's a FieldHttpError
					fieldErr, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErr.Error()).To(ContainSubstring("Ticker"))
					Expect(fieldErr.Error()).To(ContainSubstring("max"))
					Expect(fieldErr.Field()).To(Equal("Ticker"))
				})

				It("should handle min violation for slices", func() {
					type TestStruct struct {
						Images []string `validate:"min=4"`
					}
					testData := TestStruct{Images: []string{"a", "b"}}
					validationErr := validate.Struct(testData)

					result := util.ProcessValidationError(validationErr)

					// First assert that it's a FieldHttpError
					fieldErr, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErr.Error()).To(ContainSubstring("Images"))
					Expect(fieldErr.Error()).To(ContainSubstring("min"))
					Expect(fieldErr.Field()).To(Equal("Images"))
				})

				It("should handle max violation for slices", func() {
					type TestStruct struct {
						Notes []string `validate:"max=1"`
					}
					testData := TestStruct{Notes: []string{"a", "b", "c"}}
					validationErr := validate.Struct(testData)

					result := util.ProcessValidationError(validationErr)

					// First assert that it's a FieldHttpError
					fieldErr, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErr.Error()).To(ContainSubstring("Notes"))
					Expect(fieldErr.Error()).To(ContainSubstring("max"))
					Expect(fieldErr.Field()).To(Equal("Notes"))
				})
			})

			Context("Oneof Tag", func() {
				It("should handle oneof violation", func() {
					type TestStruct struct {
						Sequence string `validate:"oneof=MWD YR"`
					}
					testData := TestStruct{Sequence: "INVALID"}
					validationErr := validate.Struct(testData)

					result := util.ProcessValidationError(validationErr)

					// First assert that it's a FieldHttpError
					fieldErr, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErr.Error()).To(ContainSubstring("Sequence"))
					Expect(fieldErr.Error()).To(ContainSubstring("oneof"))
					Expect(fieldErr.Field()).To(Equal("Sequence"))
				})

				It("should handle oneof with multiple values", func() {
					type TestStruct struct {
						Status string `validate:"oneof=SET RUNNING DROPPED TAKEN"`
					}
					testData := TestStruct{Status: "UNKNOWN"}
					validationErr := validate.Struct(testData)

					result := util.ProcessValidationError(validationErr)

					// First assert that it's a FieldHttpError
					fieldErr, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErr.Error()).To(ContainSubstring("Status"))
					Expect(fieldErr.Field()).To(Equal("Status"))
				})
			})

			Context("Email Tag", func() {
				It("should handle email validation failure", func() {
					type TestStruct struct {
						Email string `validate:"email"`
					}
					testData := TestStruct{Email: "invalid-email"}
					validationErr := validate.Struct(testData)

					result := util.ProcessValidationError(validationErr)

					// First assert that it's a FieldHttpError
					fieldErr, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErr.Error()).To(ContainSubstring("Email"))
					Expect(fieldErr.Error()).To(ContainSubstring("email"))
					Expect(fieldErr.Field()).To(Equal("Email"))
				})
			})

			Context("Dive Tag (Nested Validation)", func() {
				It("should handle dive validation for nested structs", func() {
					type Image struct {
						Timeframe string `validate:"required,oneof=DL WK MN"`
					}
					type Journal struct {
						Images []Image `validate:"dive"`
					}
					testData := Journal{
						Images: []Image{{Timeframe: "INVALID"}},
					}
					validationErr := validate.Struct(testData)

					result := util.ProcessValidationError(validationErr)

					// First assert that it's a FieldHttpError
					fieldErr, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErr.Error()).To(ContainSubstring("Timeframe"))
					Expect(fieldErr.Field()).To(Equal("Timeframe"))
				})
			})

			Context("Time Format Tag", func() {
				It("should handle time parsing error for created-before field", func() {
					timeErr := &time.ParseError{
						Layout: "2006-01-02",
						Value:  "invalid-date",
					}
					// Wrap with field context to simulate Gin binding error
					fieldErr := fmt.Errorf("created-before: %w", timeErr)

					result := util.ProcessValidationError(fieldErr)

					// First assert that it's a FieldHttpError
					fieldErrResult, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErrResult.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErrResult.Field()).To(Equal("CreatedBefore"))
					Expect(fieldErrResult.Error()).To(Equal("Must be YYYY-MM-DD format"))
				})

				It("should handle time parsing error for created-after field", func() {
					timeErr := &time.ParseError{
						Layout: "2006-01-02",
						Value:  "15-02-2024",
					}
					// Wrap with field context
					fieldErr := fmt.Errorf("created-after: %w", timeErr)

					result := util.ProcessValidationError(fieldErr)

					// First assert that it's a FieldHttpError
					fieldErrResult, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErrResult.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErrResult.Field()).To(Equal("CreatedAfter"))
					Expect(fieldErrResult.Error()).To(Equal("Must be YYYY-MM-DD format"))
				})
			})

			Context("Multiple Validation Errors", func() {
				It("should return first error when multiple fields fail", func() {
					type TestStruct struct {
						Name  string `validate:"required"`
						Age   int    `validate:"min=18"`
						Email string `validate:"email"`
					}
					testData := TestStruct{Name: "", Age: 15, Email: "invalid"}
					validationErr := validate.Struct(testData)

					result := util.ProcessValidationError(validationErr)

					// First assert that it's a FieldHttpError
					fieldErr, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErr.Error()).To(ContainSubstring("Name"))
					Expect(fieldErr.Field()).To(Equal("Name"))
				})
			})
		})

		// --------------------------------------------------------------------
		// 2.2: JSON Parsing Errors
		// --------------------------------------------------------------------
		Context("JSON Parsing Errors", func() {

			Context("Syntax Errors", func() {
				It("should handle invalid JSON syntax", func() {
					invalidJSON := `{"name": "test", "age": }`
					var test struct {
						Name string `json:"name"`
						Age  int    `json:"age"`
					}
					err := json.Unmarshal([]byte(invalidJSON), &test)

					result := util.ProcessValidationError(err)

					Expect(result).To(HaveOccurred())
					Expect(result.Code()).To(Equal(http.StatusBadRequest))
					Expect(result.Error()).To(ContainSubstring("Invalid JSON"))
				})

				It("should handle unclosed brackets", func() {
					invalidJSON := `{"name": "test"`
					var test struct {
						Name string `json:"name"`
					}
					err := json.Unmarshal([]byte(invalidJSON), &test)

					result := util.ProcessValidationError(err)

					Expect(result).To(HaveOccurred())
					Expect(result.Code()).To(Equal(http.StatusBadRequest))
				})
			})

			Context("Type Mismatch Errors", func() {
				It("should handle string where int expected", func() {
					invalidJSON := `{"age": "not-a-number"}`
					var test struct {
						Age int `json:"age"`
					}
					err := json.Unmarshal([]byte(invalidJSON), &test)

					result := util.ProcessValidationError(err)

					// First assert that it's a FieldHttpError
					fieldErr, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErr.Error()).To(ContainSubstring("expects"))
					Expect(fieldErr.Field()).To(Equal("age"))
				})

				It("should handle int where string expected", func() {
					invalidJSON := `{"name": 123}`
					var test struct {
						Name string `json:"name"`
					}
					err := json.Unmarshal([]byte(invalidJSON), &test)

					result := util.ProcessValidationError(err)

					// First assert that it's a FieldHttpError
					fieldErr, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErr.Error()).To(ContainSubstring("expects"))
					Expect(fieldErr.Field()).To(Equal("name"))
				})

				It("should handle array where object expected", func() {
					invalidJSON := `{"data": [1,2,3]}`
					var test struct {
						Data struct{ ID int } `json:"data"`
					}
					err := json.Unmarshal([]byte(invalidJSON), &test)

					result := util.ProcessValidationError(err)

					// First assert that it's a FieldHttpError
					fieldErr, ok := result.(common.FieldHttpError)
					Expect(ok).To(BeTrue())

					Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
					Expect(fieldErr.Field()).To(Equal("data"))
				})
			})

			Context("Empty Body Errors", func() {
				It("should handle EOF error", func() {
					result := util.ProcessValidationError(io.EOF)

					Expect(result).To(HaveOccurred())
					Expect(result.Code()).To(Equal(http.StatusBadRequest))
					Expect(result.Error()).To(Equal("Request body cannot be empty or malformed JSON"))
				})

				It("should handle Unexpected EOF error", func() {
					result := util.ProcessValidationError(io.ErrUnexpectedEOF)

					Expect(result).To(HaveOccurred())
					Expect(result.Code()).To(Equal(http.StatusBadRequest))
					Expect(result.Error()).To(Equal("Request body cannot be empty or malformed JSON"))
				})
			})
		})

		// --------------------------------------------------------------------
		// 2.3: Query Parameter Errors
		// --------------------------------------------------------------------
		Context("Query Parameter Errors", func() {
			It("should handle numeric parsing errors", func() {
				_, err := strconv.Atoi("not-a-number")

				result := util.ProcessValidationError(err)

				Expect(result).To(HaveOccurred())
				Expect(result.Code()).To(Equal(http.StatusBadRequest))
				Expect(result.Error()).To(ContainSubstring("must be numeric"))
			})
		})

		// --------------------------------------------------------------------
		// 2.4: HttpError Passthrough
		// --------------------------------------------------------------------
		Context("HttpError Passthrough", func() {
			It("should return existing HttpError unchanged", func() {
				httpErr := common.NewHttpError("custom validation failed", http.StatusBadRequest)

				result := util.ProcessValidationError(httpErr)

				Expect(result).To(Equal(httpErr))
			})

			It("should preserve status code from HttpError", func() {
				httpErr := common.NewHttpError("not found", http.StatusNotFound)

				result := util.ProcessValidationError(httpErr)

				Expect(result.Code()).To(Equal(http.StatusNotFound))
			})
		})

		// --------------------------------------------------------------------
		// 2.5: Edge Cases and Fallback
		// --------------------------------------------------------------------
		Context("Edge Cases", func() {
			It("should handle nil error", func() {
				result := util.ProcessValidationError(nil)

				Expect(result).To(HaveOccurred())
				Expect(result.Code()).To(Equal(http.StatusInternalServerError))
				Expect(result.Error()).To(Equal("Invalid validation error format"))
			})

			It("should handle unknown error types gracefully", func() {
				genericErr := errors.New("some generic error")

				result := util.ProcessValidationError(genericErr)

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

				result := util.ProcessValidationError(validationErr)

				Expect(result).To(HaveOccurred())
				Expect(result.Code()).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
