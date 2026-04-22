//nolint:dupl // Intentional repetition for field validation contexts
package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	managerMocks "github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/kohan"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("OS Handler Integration Tests", func() {
	var (
		osHandler   handler.OSHandler
		autoManager *managerMocks.AutoManagerInterface
		router      *gin.Engine
		req         *http.Request
		w           *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		core.RegisterJournalValidators()
		autoManager = managerMocks.NewAutoManagerInterface(GinkgoT())
		osHandler = handler.NewOSHandler("/tmp/kohan-capture", autoManager)

		router = gin.New()
		v1 := router.Group("/v1/api")
		os := v1.Group("/os")
		handler.SetupOSRoutes(os, osHandler)
	})

	Describe("POST /v1/api/os/screenshot", func() {
		Context("Happy Path", func() {
			Context("with valid FULL screenshot request", func() {
				BeforeEach(func() {
					autoManager.EXPECT().Screenshot(mock.Anything, kohan.ScreenshotTypeFull, "", "/tmp/kohan-capture/2025/08/test.png").Return(nil)

					payload := kohan.ScreenshotRequest{
						FileName: "test.png",
						SavePath: "2025/08",
						Type:     kohan.ScreenshotTypeFull,
						Notify:   false,
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK with screenshot metadata", func() {
					var envelope common.Envelope[kohan.ScreenshotResponse]
					util.AssertSuccess(w, http.StatusOK, &envelope)
					Expect(envelope.Data.FileName).To(Equal("test.png"))
					Expect(envelope.Data.RelativePath).To(Equal("2025/08/test.png"))
					Expect(envelope.Data.FullPath).To(Equal("/tmp/kohan-capture/2025/08/test.png"))
				})
			})

			Context("with valid REGION screenshot request and window", func() {
				BeforeEach(func() {
					autoManager.EXPECT().Screenshot(mock.Anything, kohan.ScreenshotTypeRegion, "TradingView", "/tmp/kohan-capture/test.png").Return(nil)

					payload := kohan.ScreenshotRequest{
						FileName: "test.png",
						SavePath: ".",
						Type:     kohan.ScreenshotTypeRegion,
						Window:   "TradingView",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK", func() {
					var envelope common.Envelope[kohan.ScreenshotResponse]
					util.AssertSuccess(w, http.StatusOK, &envelope)
					Expect(envelope.Data.FileName).To(Equal("test.png"))
				})
			})
		})

		Context("Field Validations - Binding Errors", func() {
			Context("FileName field", func() {
				It("should return 400 for missing FileName", func() {
					payload := map[string]string{
						"save_path": ".",
						"type":      "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "FileName", "required")
				})

				It("should return 400 for FileName exceeding max length", func() {
					payload := map[string]string{
						"file_name": "this_is_a_very_long_filename_that_exceeds_fifty_chars_limit.png",
						"save_path": ".",
						"type":      "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "FileName", "max")
				})

				It("should return 400 for non-png FileName", func() {
					payload := map[string]string{
						"file_name": "test.jpg",
						"save_path": ".",
						"type":      "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "FileName", "image_file")
				})

				It("should return 400 for FileName with path separators", func() {
					payload := map[string]string{
						"file_name": "../etc/passwd.png",
						"save_path": ".",
						"type":      "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "FileName", "image_file")
				})
			})

			Context("SavePath field", func() {
				It("should return 400 for missing SavePath", func() {
					payload := map[string]string{
						"file_name": "test.png",
						"type":      "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "SavePath", "required")
				})

				It("should return 400 for SavePath exceeding max length", func() {
					payload := map[string]string{
						"file_name": "test.png",
						"save_path": "this/is/a/very/long/path/that/exceeds/one/hundred/characters/limit/and/should/fail/validation/in/the/handler",
						"type":      "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "SavePath", "max")
				})

				It("should return 400 for SavePath with path traversal", func() {
					payload := map[string]string{
						"file_name": "test.png",
						"save_path": "../../../etc",
						"type":      "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "SavePath", "save_path")
				})

				It("should return 400 for absolute SavePath", func() {
					payload := map[string]string{
						"file_name": "test.png",
						"save_path": "/root",
						"type":      "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "SavePath", "save_path")
				})
			})

			Context("Type field", func() {
				It("should return 400 for missing Type", func() {
					payload := map[string]string{
						"file_name": "test.png",
						"save_path": ".",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "Type", "required")
				})

				It("should return 400 for invalid Type", func() {
					payload := map[string]string{
						"file_name": "test.png",
						"save_path": ".",
						"type":      "INVALID",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "Type", "oneof")
				})
			})

			Context("Window field", func() {
				It("should return 400 for Window exceeding max length", func() {
					payload := map[string]string{
						"file_name": "test.png",
						"save_path": ".",
						"type":      "FULL",
						"window":    "this_is_a_very_long_window_name_exceeding_thirty",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "Window", "max")
				})
			})
		})

		Context("Runtime Validations", func() {
			It("should return 400 for save_path outside allowed directory", func() {
				payload := map[string]string{
					"file_name": "test.png",
					"save_path": "../../../../etc",
					"type":      "FULL",
				}
				req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("Errors", func() {
			It("should return 400 for invalid JSON", func() {
				req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", "invalid json")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 500 when auto manager fails", func() {
				autoManager.EXPECT().Screenshot(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("screenshot failed"))

				payload := kohan.ScreenshotRequest{
					FileName: "test.png",
					SavePath: ".",
					Type:     kohan.ScreenshotTypeFull,
				}
				req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("Legacy OS Endpoints", func() {
		Context("GET /v1/api/os/ticker/:ticker/record", func() {
			It("should still handle legacy ticker recording", func() {
				autoManager.EXPECT().RecordTicker(mock.Anything, "AAPL", "/tmp/kohan-capture").Return(nil)

				req, w = util.CreateTestRequest("GET", "/v1/api/os/ticker/AAPL/record", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(ContainSubstring("Success"))
			})
		})
	})
})
