package util_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test Helpers", func() {
	Context("AssertJSendFail", func() {
		Context("Valid JSend Fail Responses", func() {
			It("should decode field-specific error correctly", func() {
				w := httptest.NewRecorder()
				fieldErr := common.NewFieldHttpError("username", "Username is required")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(mustMarshal(fieldErr))

				data := util.AssertJSendFail(w, http.StatusBadRequest)
				Expect(data).To(HaveKeyWithValue("username", "Username is required"))
			})

			It("should decode multiple field errors correctly", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"fail","data":{"email":"Invalid email","password":"Too short"}}`))

				data := util.AssertJSendFail(w, http.StatusBadRequest)
				Expect(data).To(HaveKeyWithValue("email", "Invalid email"))
				Expect(data).To(HaveKeyWithValue("password", "Too short"))
				Expect(data).To(HaveLen(2))
			})

			It("should decode message-based error correctly", func() {
				w := httptest.NewRecorder()
				httpErr := common.NewHttpError("Resource not found", http.StatusNotFound)
				w.WriteHeader(http.StatusNotFound)
				w.Write(mustMarshal(httpErr))

				data := util.AssertJSendFail(w, http.StatusNotFound)
				Expect(data).To(HaveKeyWithValue("message", "Resource not found"))
			})

			It("should handle 409 Conflict as JSend fail", func() {
				w := httptest.NewRecorder()
				httpErr := common.NewHttpError("Entity already exists", http.StatusConflict)
				w.WriteHeader(http.StatusConflict)
				w.Write(mustMarshal(httpErr))

				data := util.AssertJSendFail(w, http.StatusConflict)
				Expect(data).To(HaveKeyWithValue("message", "Entity already exists"))
			})
		})

		Context("Error Cases", func() {
			It("should panic when status code does not match", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"fail","data":{"field":"error"}}`))

				Expect(func() {
					util.AssertJSendFail(w, http.StatusNotFound)
				}).To(Panic())
			})

			It("should panic when response is not valid JSON", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`invalid json`))

				Expect(func() {
					util.AssertJSendFail(w, http.StatusBadRequest)
				}).To(Panic())
			})

			It("should panic when status is not 'fail'", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"error","message":"Server error","code":500}`))

				Expect(func() {
					util.AssertJSendFail(w, http.StatusBadRequest)
				}).To(Panic())
			})

			It("should panic when 'data' field is missing", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"fail"}`))

				Expect(func() {
					util.AssertJSendFail(w, http.StatusBadRequest)
				}).To(Panic())
			})

			It("should panic when 'data' field is not a map", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"fail","data":"not a map"}`))

				Expect(func() {
					util.AssertJSendFail(w, http.StatusBadRequest)
				}).To(Panic())
			})

			It("should panic when data value is not a string", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"fail","data":{"field":123}}`))

				Expect(func() {
					util.AssertJSendFail(w, http.StatusBadRequest)
				}).To(Panic())
			})
		})
	})

	Context("AssertJSendError", func() {
		Context("Valid JSend Error Responses", func() {
			It("should decode 500 Internal Server Error correctly", func() {
				w := httptest.NewRecorder()
				httpErr := common.NewHttpError("Database connection failed", http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(mustMarshal(httpErr))

				message, code := util.AssertJSendError(w, http.StatusInternalServerError)
				Expect(message).To(Equal("Database connection failed"))
				Expect(code).To(Equal(500))
			})

			It("should decode 502 Bad Gateway correctly", func() {
				w := httptest.NewRecorder()
				httpErr := common.NewHttpError("Upstream service unavailable", http.StatusBadGateway)
				w.WriteHeader(http.StatusBadGateway)
				w.Write(mustMarshal(httpErr))

				message, code := util.AssertJSendError(w, http.StatusBadGateway)
				Expect(message).To(Equal("Upstream service unavailable"))
				Expect(code).To(Equal(502))
			})

			It("should decode 503 Service Unavailable correctly", func() {
				w := httptest.NewRecorder()
				httpErr := common.NewHttpError("Service temporarily unavailable", http.StatusServiceUnavailable)
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write(mustMarshal(httpErr))

				message, code := util.AssertJSendError(w, http.StatusServiceUnavailable)
				Expect(message).To(Equal("Service temporarily unavailable"))
				Expect(code).To(Equal(503))
			})

			It("should handle code as float64 from JSON unmarshaling", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","message":"Server error","code":500.0}`))

				message, code := util.AssertJSendError(w, http.StatusInternalServerError)
				Expect(message).To(Equal("Server error"))
				Expect(code).To(Equal(500))
			})
		})

		Context("Error Cases", func() {
			It("should panic when status code does not match", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","message":"Server error","code":500}`))

				Expect(func() {
					util.AssertJSendError(w, http.StatusBadGateway)
				}).To(Panic())
			})

			It("should panic when response is not valid JSON", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`invalid json`))

				Expect(func() {
					util.AssertJSendError(w, http.StatusInternalServerError)
				}).To(Panic())
			})

			It("should panic when status is not 'error'", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"fail","data":{"field":"error"}}`))

				Expect(func() {
					util.AssertJSendError(w, http.StatusInternalServerError)
				}).To(Panic())
			})

			It("should panic when 'message' field is missing", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","code":500}`))

				Expect(func() {
					util.AssertJSendError(w, http.StatusInternalServerError)
				}).To(Panic())
			})

			It("should panic when 'code' field is missing", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","message":"Server error"}`))

				Expect(func() {
					util.AssertJSendError(w, http.StatusInternalServerError)
				}).To(Panic())
			})

			It("should panic when 'message' field is not a string", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","message":123,"code":500}`))

				Expect(func() {
					util.AssertJSendError(w, http.StatusInternalServerError)
				}).To(Panic())
			})

			It("should panic when 'code' field is not a number", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"status":"error","message":"Server error","code":"not a number"}`))

				Expect(func() {
					util.AssertJSendError(w, http.StatusInternalServerError)
				}).To(Panic())
			})
		})
	})

	Context("AssertFieldError", func() {
		Context("Valid Field Errors", func() {
			It("should extract single field error correctly", func() {
				w := httptest.NewRecorder()
				fieldErr := common.NewFieldHttpError("email", "Invalid email format")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(mustMarshal(fieldErr))

				errorMsg := util.AssertFieldError(w, http.StatusBadRequest, "email")
				Expect(errorMsg).To(Equal("Invalid email format"))
			})

			It("should extract specific field from multiple errors", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"fail","data":{"email":"Invalid email","password":"Too short","age":"Must be positive"}}`))

				errorMsg := util.AssertFieldError(w, http.StatusBadRequest, "password")
				Expect(errorMsg).To(Equal("Too short"))
			})
		})

		Context("Error Cases", func() {
			It("should panic when field is not found", func() {
				w := httptest.NewRecorder()
				fieldErr := common.NewFieldHttpError("email", "Invalid email")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(mustMarshal(fieldErr))

				Expect(func() {
					util.AssertFieldError(w, http.StatusBadRequest, "username")
				}).To(Panic())
			})

			It("should panic when status code does not match", func() {
				w := httptest.NewRecorder()
				fieldErr := common.NewFieldHttpError("email", "Invalid email")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(mustMarshal(fieldErr))

				Expect(func() {
					util.AssertFieldError(w, http.StatusNotFound, "email")
				}).To(Panic())
			})
		})
	})
})

// mustMarshal is a test helper that marshals or panics
func mustMarshal(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
