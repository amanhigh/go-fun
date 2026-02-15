package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/common/util"
	kohanhandler "github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("ImageHandler", func() {
	var (
		imageMgr *mocks.ImageManager
		handler  *kohanhandler.ImageHandlerImpl
		router   *gin.Engine
		testCtx  = context.Background()
		image    barkat.Image
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()
		imageMgr = mocks.NewImageManager(GinkgoT())
		handler = kohanhandler.NewImageHandler(imageMgr)

		image = barkat.Image{
			ID:        "test-image-id",
			EntryID:   "test-entry-id",
			Timeframe: "DL",
		}

		// Setup routes using reusable function
		v1 := router.Group("/v1")
		kohanhandler.SetupImageRoutes(v1, handler)
	})

	Context("HandleCreateImage", func() {
		Context("with valid request", func() {
			var (
				req *http.Request
				w   *httptest.ResponseRecorder
			)

			BeforeEach(func() {
				// Mock manager success
				imageMgr.EXPECT().CreateImage(testCtx, "test-entry-id", mock.AnythingOfType("*barkat.Image")).Return(nil)

				// Setup HTTP request using helper
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries/test-entry-id/images", image)
			})

			It("should create image and return 201", func() {
				router.ServeHTTP(w, req)

				var response barkat.Image
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Timeframe).To(Equal("DL"))
			})
		})

		Context("with invalid JSON", func() {
			var (
				req *http.Request
				w   *httptest.ResponseRecorder
			)

			BeforeEach(func() {
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries/test-entry-id/images", []byte("invalid json"))
			})

			It("should return 400 error", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with manager error", func() {
			var (
				req *http.Request
				w   *httptest.ResponseRecorder
			)

			BeforeEach(func() {
				// Mock manager error
				imageMgr.EXPECT().CreateImage(testCtx, "test-entry-id", mock.AnythingOfType("*barkat.Image")).Return(common.ErrNotFound)

				// Setup HTTP request using helper
				req, w = util.CreateTestRequest("POST", "/v1/journal-entries/test-entry-id/images", image)
			})

			It("should return manager error", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	Context("HandleListImages", func() {
		Context("with valid entry", func() {
			var (
				req *http.Request
				w   *httptest.ResponseRecorder
			)

			BeforeEach(func() {
				// Mock manager success
				expectedImages := []barkat.Image{image}
				imageMgr.EXPECT().ListImages(testCtx, "test-entry-id").Return(expectedImages, nil)

				// Setup HTTP request using helper
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries/test-entry-id/images", nil)
			})

			It("should list images and return 200", func() {
				router.ServeHTTP(w, req)

				var response map[string][]barkat.Image
				util.AssertJSONAndStatus(w, http.StatusOK, &response)
				Expect(response["images"]).To(HaveLen(1))
				Expect(response["images"][0].Timeframe).To(Equal("DL"))
			})
		})

		Context("with manager error", func() {
			var (
				req *http.Request
				w   *httptest.ResponseRecorder
			)

			BeforeEach(func() {
				// Mock manager error
				imageMgr.EXPECT().ListImages(testCtx, "unknown-id").Return(nil, common.ErrNotFound)

				// Setup HTTP request using helper
				req, w = util.CreateTestRequest("GET", "/v1/journal-entries/unknown-id/images", nil)
			})

			It("should return manager error", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	Context("HandleDeleteImage", func() {
		Context("with valid entry and image", func() {
			var (
				req *http.Request
				w   *httptest.ResponseRecorder
			)

			BeforeEach(func() {
				// Mock manager success
				imageMgr.EXPECT().DeleteImage(testCtx, "test-entry-id", "test-image-id").Return(nil)

				// Setup HTTP request using helper
				req, w = util.CreateTestRequest("DELETE", "/v1/journal-entries/test-entry-id/images/test-image-id", nil)
			})

			It("should delete image and return 204", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNoContent))
				Expect(w.Body.String()).To(BeEmpty())
			})
		})

		Context("with nonexistant image", func() {
			var (
				req *http.Request
				w   *httptest.ResponseRecorder
			)

			BeforeEach(func() {
				// Mock manager error
				imageMgr.EXPECT().DeleteImage(testCtx, "test-entry-id", "nonexistent-image").Return(common.ErrNotFound)

				// Setup HTTP request using helper
				req, w = util.CreateTestRequest("DELETE", "/v1/journal-entries/test-entry-id/images/nonexistent-image", nil)
			})

			It("should return not found", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
