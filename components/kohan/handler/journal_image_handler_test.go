package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

// ImageHandler Integration Tests - Tests behavior with real SQLite DB, managers, and repositories
// This tests the complete HTTP → Handler → Manager → Repository → Database flow
var _ = PDescribe("ImageHandler Integration", func() {
	var (
		imageHandler *handler.ImageHandlerImpl
		router       *gin.Engine
		testCtx      = context.Background()
		db           *gorm.DB
		entryMgr     manager.JournalManager
		imgMgr       manager.ImageManager
		entry        barkat.Journal
		req          *http.Request
		w            *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		var err error
		// Create real SQLite database for testing using proper migrations
		db, err = core.CreateTestBarkatDB()
		Expect(err).ToNot(HaveOccurred())

		// Create real managers and repositories (no mocks)
		entryRepo := repository.NewJournalRepository(db)
		entryMgr = manager.NewJournalManager(entryRepo)
		imgMgr = manager.NewImageManager(entryMgr, repository.NewImageRepository(db))
		imageHandler = handler.NewImageHandler(imgMgr)

		// Setup Gin router using helper
		router = util.CreateTestGinRouter()
		v1 := router.Group("/v1")
		journal := v1.Group("/journal")
		handler.SetupImageRoutes(journal, imageHandler)

		// Create test entry for image operations
		entry = barkat.Journal{
			Ticker:   "GRSE",
			Sequence: "MWD",
			Type:     "REJECTED",
			Status:   "FAIL",
			Images: []barkat.Image{
				{Timeframe: "DL"},
				{Timeframe: "WK"},
				{Timeframe: "MN"},
				{Timeframe: "TMN"},
			},
		}
		Expect(entryMgr.CreateJournal(testCtx, &entry)).To(Succeed())
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Context("HandleCreateImage", func() {
		Context("with valid request", func() {
			BeforeEach(func() {
				// Setup HTTP request with real image data
				image := barkat.Image{
					Timeframe: "DL",
				}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+entry.ID+"/images", image)
			})

			It("should create image and return 201", func() {
				router.ServeHTTP(w, req)

				var response barkat.Image
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)

				// Verify the created image has proper data
				Expect(response.Timeframe).To(Equal("DL"))
				Expect(response.JournalID).To(Equal(entry.ID))
				Expect(response.ID).ToNot(BeEmpty())
				Expect(response.CreatedAt).ToNot(BeZero())

				// Verify image is actually in database
				images, err := imgMgr.ListImages(testCtx, entry.ID)
				Expect(err).ToNot(HaveOccurred())
				Expect(images).To(HaveLen(1))
				Expect(images[0].ID).To(Equal(response.ID))
			})
		})

		Context("with valid timeframe variations", func() {
			It("should accept DL timeframe", func() {
				image := barkat.Image{Timeframe: "DL"}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+entry.ID+"/images", image)
				router.ServeHTTP(w, req)
				var response barkat.Image
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Timeframe).To(Equal("DL"))
			})

			It("should accept WK timeframe", func() {
				image := barkat.Image{Timeframe: "WK"}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+entry.ID+"/images", image)
				router.ServeHTTP(w, req)
				var response barkat.Image
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Timeframe).To(Equal("WK"))
			})

			It("should accept MN timeframe", func() {
				image := barkat.Image{Timeframe: "MN"}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+entry.ID+"/images", image)
				router.ServeHTTP(w, req)
				var response barkat.Image
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Timeframe).To(Equal("MN"))
			})

			It("should accept TMN timeframe", func() {
				image := barkat.Image{Timeframe: "TMN"}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+entry.ID+"/images", image)
				router.ServeHTTP(w, req)
				var response barkat.Image
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Timeframe).To(Equal("TMN"))
			})

			It("should accept SMN timeframe", func() {
				image := barkat.Image{Timeframe: "SMN"}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+entry.ID+"/images", image)
				router.ServeHTTP(w, req)
				var response barkat.Image
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Timeframe).To(Equal("SMN"))
			})

			It("should accept YR timeframe", func() {
				image := barkat.Image{Timeframe: "YR"}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+entry.ID+"/images", image)
				router.ServeHTTP(w, req)
				var response barkat.Image
				util.AssertJSONAndStatus(w, http.StatusCreated, &response)
				Expect(response.Timeframe).To(Equal("YR"))
			})
		})

		Context("field validation", func() {
			It("should reject empty timeframe", func() {
				image := barkat.Image{Timeframe: ""}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+entry.ID+"/images", image)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				// Should return validation error about required timeframe
				var errorResponse map[string]any
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("required"))
			})

			It("should reject invalid timeframe", func() {
				image := barkat.Image{Timeframe: "INVALID"}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+entry.ID+"/images", image)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				// Should return validation error about oneof constraint
				var errorResponse map[string]any
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("oneof"))
			})

			It("should reject lowercase timeframe", func() {
				image := barkat.Image{Timeframe: "dl"}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+entry.ID+"/images", image)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				// Should return validation error about oneof constraint
				var errorResponse map[string]any
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("oneof"))
			})

			It("should reject extra whitespace timeframe", func() {
				image := barkat.Image{Timeframe: " DL "}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+entry.ID+"/images", image)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				// Should return validation error about oneof constraint
				var errorResponse map[string]any
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("oneof"))
			})

			It("should reject missing timeframe field", func() {
				image := barkat.Image{}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+entry.ID+"/images", image)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
				// Should return validation error about required timeframe
				var errorResponse map[string]any
				util.AssertJSONAndStatus(w, http.StatusBadRequest, &errorResponse)
				Expect(errorResponse["error"]).To(ContainSubstring("required"))
			})
		})

		Context("with invalid JSON", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/"+entry.ID+"/images", []byte("invalid json"))
			})

			It("should return 400 error", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with non-existent entry", func() {
			BeforeEach(func() {
				image := barkat.Image{Timeframe: "DL"}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/nonexistent/images", image)
			})

			It("should return 404 for non-existent entry", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with malformed entry ID", func() {
			BeforeEach(func() {
				image := barkat.Image{Timeframe: "DL"}
				req, w = util.CreateTestRequest("POST", barkat.JournalBase+"/invalid-id/images", image)
			})

			It("should return 404 for malformed entry ID", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with empty entry ID", func() {
			BeforeEach(func() {
				image := barkat.Image{Timeframe: "DL"}
				req, w = util.CreateTestRequest("POST", "barkat.JournalBase//images", image)
			})

			It("should return 400 for empty entry ID (route not found)", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Context("HandleListImages", func() {
		var (
			createdImages []barkat.Image
		)

		BeforeEach(func() {
			// Create multiple images for testing
			timeframes := []string{"DL", "WK", "MN"}
			for _, tf := range timeframes {
				image := barkat.Image{Timeframe: tf}
				created, err := imgMgr.CreateImage(testCtx, entry.ID, image)
				Expect(err).ToNot(HaveOccurred())
				createdImages = append(createdImages, *created)
			}
		})

		Context("with valid entry", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", barkat.JournalBase+"/"+entry.ID+"/images", nil)
			})

			It("should list images and return 200", func() {
				router.ServeHTTP(w, req)

				var response map[string][]barkat.Image
				util.AssertJSONAndStatus(w, http.StatusOK, &response)

				// Verify all images are returned
				Expect(response["images"]).To(HaveLen(3))

				// Verify timeframes are preserved
				timeframes := []string{}
				for _, img := range response["images"] {
					timeframes = append(timeframes, img.Timeframe)
				}
				Expect(timeframes).To(ContainElements("DL", "WK", "MN"))

				// Verify each image has proper metadata
				for _, img := range response["images"] {
					Expect(img.JournalID).To(Equal(entry.ID))
					Expect(img.ID).ToNot(BeEmpty())
					Expect(img.CreatedAt).ToNot(BeZero())
				}
			})

			It("should return images in creation order", func() {
				router.ServeHTTP(w, req)

				var response map[string][]barkat.Image
				util.AssertJSONAndStatus(w, http.StatusOK, &response)

				// Verify images are returned in creation order (DL, WK, MN)
				Expect(response["images"][0].Timeframe).To(Equal("DL"))
				Expect(response["images"][1].Timeframe).To(Equal("WK"))
				Expect(response["images"][2].Timeframe).To(Equal("MN"))
			})
		})

		Context("with non-existent entry", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", barkat.JournalBase+"/nonexistent/images", nil)
			})

			It("should return 404 for non-existent entry", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with malformed entry ID", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", barkat.JournalBase+"/invalid-id/images", nil)
			})

			It("should return 404 for malformed entry ID", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with empty entry ID", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("GET", "barkat.JournalBase//images", nil)
			})

			It("should return 404 for empty entry ID (route not found)", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with no images for entry", func() {
			BeforeEach(func() {
				// Delete all images for this entry
				for _, img := range createdImages {
					err := imgMgr.DeleteImage(testCtx, entry.ID, img.ID)
					Expect(err).ToNot(HaveOccurred())
				}
				req, w = util.CreateTestRequest("GET", barkat.JournalBase+"/"+entry.ID+"/images", nil)
			})

			It("should return empty array for entry with no images", func() {
				router.ServeHTTP(w, req)

				var response map[string][]barkat.Image
				util.AssertJSONAndStatus(w, http.StatusOK, &response)
				Expect(response["images"]).To(BeEmpty())
			})
		})
	})

	Context("HandleDeleteImage", func() {
		var imageToDelete barkat.Image

		BeforeEach(func() {
			// Create an image to delete
			image := barkat.Image{Timeframe: "DL"}
			created, err := imgMgr.CreateImage(testCtx, entry.ID, image)
			Expect(err).ToNot(HaveOccurred())
			imageToDelete = *created
		})

		Context("with valid entry and image", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("DELETE", "barkat.JournalBase/"+entry.ID+"/images/"+imageToDelete.ID, nil)
			})

			It("should delete image and return 204", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNoContent))
				Expect(w.Body.String()).To(BeEmpty())

				// Verify image is actually deleted from database
				images, err := imgMgr.ListImages(testCtx, entry.ID)
				Expect(err).ToNot(HaveOccurred())
				Expect(images).To(BeEmpty())
			})

			It("should return 204 on first delete and 404 on second delete", func() {
				// Delete once
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNoContent))

				// Delete again (should return 404 since image is gone)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
				Expect(w.Body.String()).ToNot(BeEmpty())
			})
		})

		Context("with non-existent image", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("DELETE", "barkat.JournalBase/"+entry.ID+"/images/nonexistent-image", nil)
			})

			It("should return 404 for non-existent image", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with non-existent entry", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("DELETE", barkat.JournalBase+"/nonexistent/images/"+imageToDelete.ID, nil)
			})

			It("should return 404 for non-existent entry", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with malformed entry ID", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("DELETE", barkat.JournalBase+"/invalid-id/images/"+imageToDelete.ID, nil)
			})

			It("should return 404 for malformed entry ID", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with empty entry ID", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("DELETE", "barkat.JournalBase//images/"+imageToDelete.ID, nil)
			})

			It("should return 400 for empty entry ID (route not found)", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with empty image ID", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("DELETE", "barkat.JournalBase/"+entry.ID+"/images/", nil)
			})

			It("should return 400 for empty image ID (route not found)", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with malformed image ID", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest("DELETE", "barkat.JournalBase/"+entry.ID+"/images/invalid-id", nil)
			})

			It("should return 404 for malformed image ID", func() {
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with concurrent deletion safety", func() {
			It("should return 404 when delete request races after first succeeds", func() {
				// First delete
				req1, w1 := util.CreateTestRequest("DELETE", "barkat.JournalBase-entries/"+entry.ID+"/images/"+imageToDelete.ID, nil)
				router.ServeHTTP(w1, req1)
				Expect(w1.Code).To(Equal(http.StatusNoContent))

				// Second delete should report missing since image no longer exists
				req2, w2 := util.CreateTestRequest("DELETE", "barkat.JournalBase-entries/"+entry.ID+"/images/"+imageToDelete.ID, nil)
				router.ServeHTTP(w2, req2)
				Expect(w2.Code).To(Equal(http.StatusNotFound))

				// Verify image is deleted
				images, err := imgMgr.ListImages(testCtx, entry.ID)
				Expect(err).ToNot(HaveOccurred())
				Expect(images).To(BeEmpty())
			})
		})
	})
})
