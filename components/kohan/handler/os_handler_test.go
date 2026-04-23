//nolint:dupl // Intentional repetition for field validation contexts
package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

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
		osHandler        handler.OSHandler
		autoManager      *managerMocks.AutoManagerInterface
		router           *gin.Engine
		req              *http.Request
		w                *httptest.ResponseRecorder
		screenShotTmpDir string
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		core.RegisterJournalValidators()
		autoManager = managerMocks.NewAutoManagerInterface(GinkgoT())
		dir, err := os.MkdirTemp("", "kohan-screenshots-*")
		Expect(err).ToNot(HaveOccurred())
		screenShotTmpDir = dir
		osHandler = handler.NewOSHandler(autoManager)

		router = gin.New()
		v1 := router.Group("/v1/api")
		os := v1.Group("/os")
		handler.SetupOSRoutes(os, osHandler)
	})

	AfterEach(func() {
		os.RemoveAll(screenShotTmpDir)
	})

	Describe("POST /v1/api/os/screenshot", func() {
		Context("Happy Path", func() {
			Context("with valid FULL screenshot request", func() {
				BeforeEach(func() {
					autoManager.EXPECT().Screenshot(mock.Anything, kohan.ScreenshotDirectoryTypeJournal, "test.png", kohan.ScreenshotTypeFull, "").Return(screenShotTmpDir+"/test.png", nil)

					payload := kohan.ScreenshotRequest{
						FileName:      "test.png",
						DirectoryType: kohan.ScreenshotDirectoryTypeJournal,
						Type:          kohan.ScreenshotTypeFull,
						Notify:        false,
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK with screenshot metadata", func() {
					var envelope common.Envelope[kohan.ScreenshotResponse]
					util.AssertSuccess(w, http.StatusOK, &envelope)
					Expect(envelope.Data.FileName).To(Equal("test.png"))
					Expect(envelope.Data.FullPath).To(Equal(screenShotTmpDir + "/test.png"))
				})
			})

			Context("with valid REGION screenshot request and window", func() {
				BeforeEach(func() {
					autoManager.EXPECT().Screenshot(mock.Anything, kohan.ScreenshotDirectoryTypeDownload, "test.png", kohan.ScreenshotTypeRegion, "TradingView").Return("/home/Downloads/test.png", nil)

					payload := kohan.ScreenshotRequest{
						FileName:      "test.png",
						DirectoryType: kohan.ScreenshotDirectoryTypeDownload,
						Type:          kohan.ScreenshotTypeRegion,
						Window:        "TradingView",
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
						"directory_type": "JOURNAL",
						"type":           "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "FileName", "required")
				})

				It("should return 400 for FileName exceeding max length", func() {
					payload := map[string]string{
						"file_name":      "this_is_a_very_long_filename_that_exceeds_fifty_chars_limit.png",
						"directory_type": "JOURNAL",
						"type":           "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "FileName", "max")
				})

				It("should return 400 for invalid FileName extension", func() {
					payload := map[string]string{
						"file_name":      "test.gif",
						"directory_type": "JOURNAL",
						"type":           "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "FileName", "image_file")
				})

				It("should return 400 for FileName with path separators", func() {
					payload := map[string]string{
						"file_name":      "../etc/passwd.png",
						"directory_type": "JOURNAL",
						"type":           "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "FileName", "image_file")
				})
			})

			Context("DirectoryType field", func() {
				It("should return 400 for missing DirectoryType", func() {
					payload := map[string]string{
						"file_name": "test.png",
						"type":      "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "DirectoryType", "required")
				})

				It("should return 400 for invalid DirectoryType", func() {
					payload := map[string]string{
						"file_name":      "test.png",
						"directory_type": "OTHER",
						"type":           "FULL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "DirectoryType", "oneof")
				})
			})

			Context("Type field", func() {
				It("should return 400 for missing Type", func() {
					payload := map[string]string{
						"file_name":      "test.png",
						"directory_type": "JOURNAL",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "Type", "required")
				})

				It("should return 400 for invalid Type", func() {
					payload := map[string]string{
						"file_name":      "test.png",
						"directory_type": "JOURNAL",
						"type":           "INVALID",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "Type", "oneof")
				})
			})

			Context("Window field", func() {
				It("should return 400 for Window exceeding max length", func() {
					payload := map[string]string{
						"file_name":      "test.png",
						"directory_type": "JOURNAL",
						"type":           "FULL",
						"window":         "this_is_a_very_long_window_name_exceeding_thirty",
					}
					req, w = util.CreateTestRequest("POST", "/v1/api/os/screenshot", payload)
					router.ServeHTTP(w, req)
					util.AssertError(w, "Window", "max")
				})
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
					FileName:      "test.png",
					DirectoryType: kohan.ScreenshotDirectoryTypeJournal,
					Type:          kohan.ScreenshotTypeFull,
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
				autoManager.EXPECT().RecordTicker(mock.Anything, "AAPL", screenShotTmpDir).Return(nil)

				req, w = util.CreateTestRequest("GET", "/v1/api/os/ticker/AAPL/record", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(ContainSubstring("Success"))
			})
		})
	})
})
