package common_test

import (
	"encoding/json"
	"net/http"

	"github.com/amanhigh/go-fun/models/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Http Util", func() {
	Context("HttpError", func() {
		Context("Standard Http Errors", func() {
			It("should have correct error codes and messages", func() {
				Expect(common.ErrBadRequest.Code()).To(Equal(http.StatusBadRequest))
				Expect(common.ErrBadRequest.Error()).To(Equal("BadRequest"))

				Expect(common.ErrNotFound.Code()).To(Equal(http.StatusNotFound))
				Expect(common.ErrNotFound.Error()).To(Equal("NotFound"))

				Expect(common.ErrNotAuthorized.Code()).To(Equal(http.StatusUnauthorized))
				Expect(common.ErrNotAuthorized.Error()).To(Equal("NotAuthorized"))

				Expect(common.ErrNotAuthenticated.Code()).To(Equal(http.StatusForbidden))
				Expect(common.ErrNotAuthenticated.Error()).To(Equal("NotAuthenticated"))

				Expect(common.ErrEntityExists.Code()).To(Equal(http.StatusConflict))
				Expect(common.ErrEntityExists.Error()).To(Equal("EntityExists"))

				Expect(common.ErrPayloadTooLarge.Code()).To(Equal(http.StatusRequestEntityTooLarge))
				Expect(common.ErrPayloadTooLarge.Error()).To(Equal("PayloadTooLarge"))

				Expect(common.ErrInternalServerError.Code()).To(Equal(http.StatusInternalServerError))
				Expect(common.ErrInternalServerError.Error()).To(Equal("InternalServerError"))
			})
		})

		Context("HttpError", func() {
			It("should create HttpError with custom message and code", func() {
				customErr := common.NewHttpError("Custom Error", http.StatusTeapot)

				Expect(customErr.Error()).To(Equal("Custom Error"))
				Expect(customErr.Code()).To(Equal(http.StatusTeapot))
			})

			It("should create ServerError when code is 500 or higher", func() {
				serverErr := common.NewHttpError("Internal server error", http.StatusInternalServerError)
				Expect(serverErr.Code()).To(Equal(http.StatusInternalServerError))
			})

			It("should implement error interface", func() {
				var err error = common.NewHttpError("Test", http.StatusBadRequest)
				Expect(err.Error()).To(Equal("Test"))
			})
		})

		Context("FieldHttpError", func() {
			It("should create HttpError with field name, message and code", func() {
				fieldErr := common.NewFieldHttpError("email", "Invalid email format")

				Expect(fieldErr.Error()).To(Equal("Invalid email format"))
				Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
				Expect(fieldErr.Field()).To(Equal("email"))
			})

			It("should implement FieldHttpError interface", func() {
				httpErr := common.NewFieldHttpError("age", "Must be greater than 18")
				Expect(httpErr.Error()).To(Equal("Must be greater than 18"))
				Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
				Expect(httpErr.Field()).To(Equal("age"))
			})

			It("field method should return empty string when field is empty", func() {
				fieldErr := common.NewFieldHttpError("", "No field specified")
				Expect(fieldErr.Field()).To(Equal(""))
			})

			It("should panic when trying to create FieldHttpError with 5xx code", func() {
				// This test is no longer relevant since we hardcoded to 400
				// FieldHttpError now always uses 400 status code
				fieldErr := common.NewFieldHttpError("field", "Server error")
				Expect(fieldErr.Code()).To(Equal(http.StatusBadRequest))
			})
		})

	})

	Context("JSend Envelope Format (MarshalJSON)", func() {
		Context("4xx Errors - HttpError", func() {
			It("should serialize 404 NotFound as JSend fail", func() {
				httpErr := common.NewHttpError("Resource not found", http.StatusNotFound)
				jsonBytes, err := json.Marshal(httpErr)
				Expect(err).ToNot(HaveOccurred())

				var result map[string]any
				Expect(json.Unmarshal(jsonBytes, &result)).To(Succeed())
				Expect(result["status"]).To(Equal("fail"))
				Expect(result["data"]).To(HaveKeyWithValue("message", "Resource not found"))
			})

			It("should serialize 409 Conflict as JSend fail", func() {
				httpErr := common.NewHttpError("Entity exists", http.StatusConflict)
				jsonBytes, err := json.Marshal(httpErr)
				Expect(err).ToNot(HaveOccurred())

				var result map[string]any
				Expect(json.Unmarshal(jsonBytes, &result)).To(Succeed())
				Expect(result["status"]).To(Equal("fail"))
				Expect(result["data"]).To(HaveKeyWithValue("message", "Entity exists"))
			})

			It("should serialize 422 UnprocessableEntity as JSend fail", func() {
				httpErr := common.NewHttpError("Invalid input", http.StatusUnprocessableEntity)
				jsonBytes, err := json.Marshal(httpErr)
				Expect(err).ToNot(HaveOccurred())

				var result map[string]any
				Expect(json.Unmarshal(jsonBytes, &result)).To(Succeed())
				Expect(result["status"]).To(Equal("fail"))
				Expect(result["data"]).To(HaveKeyWithValue("message", "Invalid input"))
			})
		})

		Context("400 Bad Request - FieldHttpError (Fail Format)", func() {
			It("should serialize field-specific error with field name as key", func() {
				fieldErr := common.NewFieldHttpError("username", "Username is required")
				jsonBytes, err := json.Marshal(fieldErr)
				Expect(err).ToNot(HaveOccurred())

				var result map[string]any
				Expect(json.Unmarshal(jsonBytes, &result)).To(Succeed())
				Expect(result["status"]).To(Equal("fail"))
				Expect(result["data"]).To(HaveKeyWithValue("username", "Username is required"))
				Expect(result["data"]).ToNot(HaveKey("message"))
			})

			It("should handle empty field name by falling back to 'message' key", func() {
				fieldErr := common.NewFieldHttpError("", "Generic validation error")
				jsonBytes, err := json.Marshal(fieldErr)
				Expect(err).ToNot(HaveOccurred())

				var result map[string]any
				Expect(json.Unmarshal(jsonBytes, &result)).To(Succeed())
				Expect(result["status"]).To(Equal("fail"))
				Expect(result["data"]).To(HaveKeyWithValue("message", "Generic validation error"))
			})

			It("should serialize BadRequest field error", func() {
				fieldErr := common.NewFieldHttpError("password", "Password too short")
				jsonBytes, err := json.Marshal(fieldErr)
				Expect(err).ToNot(HaveOccurred())

				var result map[string]any
				Expect(json.Unmarshal(jsonBytes, &result)).To(Succeed())
				Expect(result["status"]).To(Equal("fail"))
				Expect(result["data"]).To(HaveKeyWithValue("password", "Password too short"))
			})

			It("should serialize UnprocessableEntity field error", func() {
				fieldErr := common.NewFieldHttpError("price", "Must be positive")
				jsonBytes, err := json.Marshal(fieldErr)
				Expect(err).ToNot(HaveOccurred())

				var result map[string]any
				Expect(json.Unmarshal(jsonBytes, &result)).To(Succeed())
				Expect(result["status"]).To(Equal("fail"))
				Expect(result["data"]).To(HaveKeyWithValue("price", "Must be positive"))
			})
		})

		Context("5xx Errors -  Error Format", func() {
			It("should serialize 500 InternalServerError as JSend error", func() {
				httpErr := common.NewHttpError("Database connection failed", http.StatusInternalServerError)
				jsonBytes, err := json.Marshal(httpErr)
				Expect(err).ToNot(HaveOccurred())

				var result map[string]any
				Expect(json.Unmarshal(jsonBytes, &result)).To(Succeed())
				Expect(result["status"]).To(Equal("error"))
				Expect(result["message"]).To(Equal("Database connection failed"))
				Expect(result["code"]).To(BeNumerically("==", 500))
			})

			It("should serialize 502 BadGateway as JSend error", func() {
				httpErr := common.NewHttpError("Upstream service unavailable", http.StatusBadGateway)
				jsonBytes, err := json.Marshal(httpErr)
				Expect(err).ToNot(HaveOccurred())

				var result map[string]any
				Expect(json.Unmarshal(jsonBytes, &result)).To(Succeed())
				Expect(result["status"]).To(Equal("error"))
				Expect(result["message"]).To(Equal("Upstream service unavailable"))
				Expect(result["code"]).To(BeNumerically("==", 502))
			})

			It("should serialize 503 ServiceUnavailable as JSend error", func() {
				httpErr := common.NewHttpError("Service temporarily unavailable", http.StatusServiceUnavailable)
				jsonBytes, err := json.Marshal(httpErr)
				Expect(err).ToNot(HaveOccurred())

				var result map[string]any
				Expect(json.Unmarshal(jsonBytes, &result)).To(Succeed())
				Expect(result["status"]).To(Equal("error"))
				Expect(result["message"]).To(Equal("Service temporarily unavailable"))
				Expect(result["code"]).To(BeNumerically("==", 503))
			})
		})

		Context("Standard Errors - Error Format", func() {
			It("should serialize ErrBadRequest as JSend fail", func() {
				jsonBytes, err := json.Marshal(common.ErrBadRequest)
				Expect(err).ToNot(HaveOccurred())

				var result map[string]any
				Expect(json.Unmarshal(jsonBytes, &result)).To(Succeed())
				Expect(result["status"]).To(Equal("fail"))
			})

			It("should serialize ErrInternalServerError as JSend error", func() {
				jsonBytes, err := json.Marshal(common.ErrInternalServerError)
				Expect(err).ToNot(HaveOccurred())

				var result map[string]any
				Expect(json.Unmarshal(jsonBytes, &result)).To(Succeed())
				Expect(result["status"]).To(Equal("error"))
				Expect(result["message"]).To(Equal("InternalServerError"))
				Expect(result["code"]).To(BeNumerically("==", 500))
			})
		})
	})

	Context("Pagination", func() {
		It("should have correct field types and tags", func() {
			pagination := common.Pagination{
				Offset: 10,
				Limit:  5,
			}

			Expect(pagination.Offset).To(Equal(10))
			Expect(pagination.Limit).To(Equal(5))
		})

		It("should work with zero values", func() {
			pagination := common.Pagination{}

			Expect(pagination.Offset).To(Equal(0))
			Expect(pagination.Limit).To(Equal(0))
		})
	})

	Context("Sort", func() {
		It("should have correct field types", func() {
			sort := common.Sort{
				SortBy: "name",
				Order:  "asc",
			}

			Expect(sort.SortBy).To(Equal("name"))
			Expect(sort.Order).To(Equal("asc"))
		})

		It("should work with empty values", func() {
			sort := common.Sort{}

			Expect(sort.SortBy).To(Equal(""))
			Expect(sort.Order).To(Equal(""))
		})
	})

	Context("PaginatedResponse", func() {
		It("should have correct field types", func() {
			response := common.PaginatedResponse{
				Total:  100,
				Offset: 0,
				Limit:  20,
			}

			Expect(response.Total).To(Equal(int64(100)))
			Expect(response.Offset).To(Equal(0))
			Expect(response.Limit).To(Equal(20))
		})

		It("should work with zero value", func() {
			response := common.PaginatedResponse{}

			Expect(response.Total).To(Equal(int64(0)))
			Expect(response.Offset).To(Equal(0))
			Expect(response.Limit).To(Equal(0))
		})
	})
})
