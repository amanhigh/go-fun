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
	Context("AssertError", func() {
		Context("Valid Field Errors", func() {
			It("should extract single field error correctly", func() {
				w := httptest.NewRecorder()
				fieldErr := common.NewFieldHttpError("email", "Invalid email format")
				w.WriteHeader(http.StatusBadRequest)

				// Marshal with proper error handling instead of panic
				fieldErrBytes, err := json.Marshal(fieldErr)
				Expect(err).ToNot(HaveOccurred())
				_, err = w.Write(fieldErrBytes)
				Expect(err).To(Succeed())

				util.AssertError(w, "email", "Invalid email format")
			})

			It("should extract specific field from multiple errors", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write([]byte(`{"status":"fail","data":{"email":"Invalid email","password":"Too short","age":"Must be positive"}}`))
				Expect(err).To(Succeed())

				util.AssertError(w, "password", "Too short")
			})
		})
	})

	Context("AssertSuccess", func() {
		Context("Valid JSON Responses", func() {
			It("should unmarshal envelope responses correctly", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusCreated)
				_, err := w.Write([]byte(`{"status":"success","data":{"id":"123","ticker":"GRSE"}}`))
				Expect(err).To(Succeed())

				var envelope common.Envelope[map[string]any]
				util.AssertSuccess(w, http.StatusCreated, &envelope)
				Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
				Expect(envelope.Data["id"]).To(Equal("123"))
				Expect(envelope.Data["ticker"]).To(Equal("GRSE"))
			})

			It("should unmarshal direct JSON responses correctly", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"message":"success","count":42}`))
				Expect(err).To(Succeed())

				var response map[string]any
				util.AssertSuccess(w, http.StatusOK, &response)
				Expect(response["message"]).To(Equal("success"))
				Expect(response["count"]).To(Equal(float64(42)))
			})
		})

		Context("Error Cases", func() {
			It("should fail when status code does not match", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write([]byte(`{"error":"bad request"}`))
				Expect(err).To(Succeed())

				var response map[string]any
				failures := InterceptGomegaFailures(func() {
					util.AssertSuccess(w, http.StatusOK, &response)
				})
				Expect(failures).To(HaveLen(1))
				// Should fail because expected status 200 doesn't match actual 400
				Expect(failures[0]).To(ContainSubstring("Expected status 200, got 400"))
			})

			It("should fail when JSON is invalid", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`invalid json`))
				Expect(err).To(Succeed())

				var response map[string]any
				failures := InterceptGomegaFailures(func() {
					util.AssertSuccess(w, http.StatusOK, &response)
				})
				Expect(failures).To(HaveLen(1))
				// Should fail because JSON cannot be unmarshaled
				Expect(failures[0]).To(ContainSubstring("Failed to unmarshal Test JSON"))
				Expect(failures[0]).To(ContainSubstring("invalid character"))
			})

			It("should fail when non-2xx status code is passed", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write([]byte(`{"error":"bad request"}`))
				Expect(err).To(Succeed())

				var response map[string]any
				failures := InterceptGomegaFailures(func() {
					util.AssertSuccess(w, http.StatusBadRequest, &response)
				})
				Expect(failures).To(HaveLen(1))
				// Should fail because 400 is not a 2xx status code
				Expect(failures[0]).To(ContainSubstring("AssertSuccess only accepts 2xx status codes, got 400"))
			})

			It("should fail when 3xx status code is passed", func() {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusFound)
				_, err := w.Write([]byte(`{"redirect":"/new-location"}`))
				Expect(err).To(Succeed())

				var response map[string]any
				failures := InterceptGomegaFailures(func() {
					util.AssertSuccess(w, http.StatusFound, &response)
				})
				Expect(failures).To(HaveLen(1))
				// Should fail because 302 is not a 2xx status code
				Expect(failures[0]).To(ContainSubstring("AssertSuccess only accepts 2xx status codes, got 302"))
			})
		})
	})

	Context("CreateHTMLTestRequest", func() {
		Context("HTML Request Creation", func() {
			It("should create request with text/html content type", func() {
				req, w := util.CreateHTMLTestRequest("GET", "/test/path")
				Expect(req).ToNot(BeNil())
				Expect(w).ToNot(BeNil())
				Expect(req.Method).To(Equal("GET"))
				Expect(req.URL.Path).To(Equal("/test/path"))
				Expect(req.Header.Get("Content-Type")).To(Equal("text/html"))
			})

			It("should create request with different HTTP methods", func() {
				postReq, postW := util.CreateHTMLTestRequest("POST", "/submit")
				Expect(postReq.Method).To(Equal("POST"))
				Expect(postW).ToNot(BeNil())

				putReq, putW := util.CreateHTMLTestRequest("PUT", "/update")
				Expect(putReq.Method).To(Equal("PUT"))
				Expect(putW).ToNot(BeNil())
			})
		})
	})
})
