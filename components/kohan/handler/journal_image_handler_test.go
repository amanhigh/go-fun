//nolint:dupl
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
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func decodeImageResponse(w *httptest.ResponseRecorder) barkat.Image {
	var envelope common.Envelope[barkat.Image]
	util.AssertSuccess(w, http.StatusCreated, &envelope)
	return envelope.Data
}

func decodeImageListResponse(w *httptest.ResponseRecorder) []barkat.Image {
	var envelope common.Envelope[map[string][]barkat.Image]
	util.AssertSuccess(w, http.StatusOK, &envelope)
	return envelope.Data["images"]
}

// ImageHandler Integration Tests - Comprehensive Master Specification
// Tests complete HTTP → Handler → Manager → Repository → Database flow
// Covers all PRD validations for Section 2.2 JournalImage APIs
//
// TEST STRUCTURE FORMAT:
// ====================
// Describe(API)
// -> Context(Happy Path): 2xx Success Cases
// -> Context(Field Validations): All 4xx Validation Cases
//    -> Context(Field Name): One Context for Each Field
//       -> Context(Allowed Values): All Variations of Valid Values (2xx) - If Applicable
//       -> Context(Bad Values): All Variations of Missing,Regex,Min,Max Edge Cases (4xx)
// -> Context(Errors): 5xx Server Error Cases

var _ = Describe("ImageHandler Integration - Section 2.2 JournalImage APIs", func() {
	var (
		imageHandler *handler.ImageHandlerImpl
		router       *gin.Engine
		testCtx      = context.Background()
		db           *gorm.DB
		journalMgr   manager.JournalManager
		imgMgr       manager.ImageManager
		journal      barkat.Journal
		req          *http.Request
		w            *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		var err error
		db, err = core.CreateTestBarkatDB()
		Expect(err).ToNot(HaveOccurred())

		journalRepo := repository.NewJournalRepository(db)
		journalMgr = manager.NewJournalManager(journalRepo)
		imgMgr = manager.NewImageManager(journalMgr, repository.NewImageRepository(db))
		imageHandler = handler.NewImageHandler(imgMgr)

		router = util.CreateTestGinRouter()
		v1 := router.Group("/v1")
		journalGroup := v1.Group("/journals")
		handler.SetupImageRoutes(journalGroup, imageHandler)

		// Create base journal for image operations (with minimal images)
		journal = barkat.Journal{
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
		Expect(journalMgr.CreateJournal(testCtx, &journal)).To(Succeed())
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	// ============================================================================
	// 2.2.1 POST /v1/journals/{journal-id}/images - Add Image
	// ============================================================================
	Describe("POST /v1/journals/{journal-id}/images - Add Image (2.2.1)", func() {
		Context("Happy Path", func() {
			Context("with valid timeframe image", func() {
				var response barkat.Image

				BeforeEach(func() {
					image := barkat.Image{
						Timeframe: "SMN",
						FileName:  "RELIANCE.mwd.rejected.oe__20240115_132138.png",
					}
					req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
					router.ServeHTTP(w, req)
				})

				It("should return 201 Created", func() {
					Expect(w.Code).To(Equal(http.StatusCreated))
				})

				It("should return Envelope success", func() {
					var envelope common.Envelope[barkat.Image]
					util.AssertSuccess(w, http.StatusCreated, &envelope)
					Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
				})

				It("should return created image with external ID", func() {
					response = decodeImageResponse(w)
					Expect(response.ExternalID).To(HavePrefix("img_"))
				})

				It("should preserve timeframe field", func() {
					response = decodeImageResponse(w)
					Expect(response.Timeframe).To(Equal("SMN"))
				})

				It("should preserve file_name field", func() {
					response = decodeImageResponse(w)
					Expect(response.FileName).To(Equal("RELIANCE.mwd.rejected.oe__20240115_132138.png"))
				})

				It("should set created_at timestamp", func() {
					response = decodeImageResponse(w)
					Expect(response.CreatedAt).ToNot(BeZero())
				})

				It("should persist image to database", func() {
					response = decodeImageResponse(w)
					images, err := imgMgr.ListImages(testCtx, journal.ExternalID)
					Expect(err).ToNot(HaveOccurred())
					// 4 from journal creation + 1 new = 5
					Expect(images).To(HaveLen(5))
				})
			})
		})

		Context("Field Validations", func() {
			Context("Timeframe Field", func() {
				Context("Allowed Values", func() {
					It("should accept timeframe = DL", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						response := decodeImageResponse(w)
						Expect(response.Timeframe).To(Equal("DL"))
					})

					It("should accept timeframe = WK", func() {
						image := barkat.Image{Timeframe: "WK", FileName: "RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						response := decodeImageResponse(w)
						Expect(response.Timeframe).To(Equal("WK"))
					})

					It("should accept timeframe = MN", func() {
						image := barkat.Image{Timeframe: "MN", FileName: "RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						response := decodeImageResponse(w)
						Expect(response.Timeframe).To(Equal("MN"))
					})

					It("should accept timeframe = TMN", func() {
						image := barkat.Image{Timeframe: "TMN", FileName: "RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						response := decodeImageResponse(w)
						Expect(response.Timeframe).To(Equal("TMN"))
					})

					It("should accept timeframe = SMN", func() {
						image := barkat.Image{Timeframe: "SMN", FileName: "RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						response := decodeImageResponse(w)
						Expect(response.Timeframe).To(Equal("SMN"))
					})

					It("should accept timeframe = YR", func() {
						image := barkat.Image{Timeframe: "YR", FileName: "RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						response := decodeImageResponse(w)
						Expect(response.Timeframe).To(Equal("YR"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing timeframe (PRD: required)", func() {
						image := barkat.Image{Timeframe: "", FileName: "RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframe", "required")
					})

					It("should return 400 for invalid timeframe enum (PRD: must be DL,WK,MN,TMN,SMN,YR)", func() {
						image := barkat.Image{Timeframe: "INVALID", FileName: "RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframe", "oneof")
					})

					It("should return 400 for lowercase timeframe (PRD: case-sensitive)", func() {
						image := barkat.Image{Timeframe: "dl", FileName: "RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframe", "oneof")
					})

					It("should return 400 for timeframe with whitespace", func() {
						image := barkat.Image{Timeframe: " DL ", FileName: "RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframe", "oneof")
					})
				})
			})

			Context("FileName Field", func() {
				Context("Allowed Values", func() {
					It("should accept minimum file name length (1 char)", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "a.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						response := decodeImageResponse(w)
						Expect(response.FileName).To(Equal("a.png"))
					})

					It("should accept maximum file name length (255 chars)", func() {
						longFileName := ""
						for i := 0; i < 251; i++ { // 251 + ".png" = 255
							longFileName += "a"
						}
						longFileName += ".png"
						image := barkat.Image{Timeframe: "DL", FileName: longFileName}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						response := decodeImageResponse(w)
						Expect(response.FileName).To(Equal(longFileName))
					})

					It("should accept PNG file extension", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						response := decodeImageResponse(w)
						Expect(response.FileName).To(Equal("RELIANCE.mwd.test.png"))
					})

					It("should accept JPG file extension", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "RELIANCE.mwd.test.jpg"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						response := decodeImageResponse(w)
						Expect(response.FileName).To(Equal("RELIANCE.mwd.test.jpg"))
					})

					It("should accept JPEG file extension", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "RELIANCE.mwd.test.jpeg"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						response := decodeImageResponse(w)
						Expect(response.FileName).To(Equal("RELIANCE.mwd.test.jpeg"))
					})

					It("should accept file name with special characters (dots, hyphens, underscores)", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "RELIANCE.mwd.rejected.oe__20240115_132138.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						response := decodeImageResponse(w)
						Expect(response.FileName).To(Equal("RELIANCE.mwd.rejected.oe__20240115_132138.png"))
					})

					It("should accept file name with numbers", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "RELIANCE123.mwd.456.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						response := decodeImageResponse(w)
						Expect(response.FileName).To(Equal("RELIANCE123.mwd.456.test.png"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing file_name (PRD: required)", func() {
						image := barkat.Image{Timeframe: "DL", FileName: ""}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						util.AssertError(w, "FileName", "required")
					})

					It("should return 400 for file_name exceeding max length (PRD: max 255 chars)", func() {
						longFileName := ""
						for i := 0; i < 256; i++ { // 256 chars exceeds limit
							longFileName += "a"
						}
						image := barkat.Image{Timeframe: "DL", FileName: longFileName}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						util.AssertError(w, "FileName", "max")
					})

					It("should return 400 for file_name without extension", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "RELIANCE.mwd.test"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						util.AssertError(w, "FileName", "file_name")
					})

					It("should return 400 for file_name with invalid extension", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "RELIANCE.mwd.test.txt"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						util.AssertError(w, "FileName", "file_name")
					})

					It("should return 400 for file_name with invalid characters (PRD: alphanumeric, dots, hyphens, underscores only)", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "RELIANCE@mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						util.AssertError(w, "FileName", "file_name")
					})

					It("should return 400 for file_name with spaces", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "RELIANCE mwd test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						util.AssertError(w, "FileName", "file_name")
					})

					It("should return 400 for file_name with path separators", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "path/RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
						router.ServeHTTP(w, req)
						util.AssertError(w, "FileName", "file_name")
					})
				})
			})

			Context("Journal ID Path Parameter", func() {
				Context("Bad Values", func() {
					It("should return 404 for non-existent journal ID", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/nonexistent-id/images", image)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for malformed journal ID", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "RELIANCE.mwd.test.png"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/invalid-uuid-format/images", image)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for valid UUID format but non-existent", func() {
						image := barkat.Image{Timeframe: "DL", FileName: "RELIANCE.mwd.test.png"}
						// BUG: uid format not matching journal format.
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/550e8400-e29b-41d4-a716-446655440000/images", image)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 400 for invalid JSON", func() {
				req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", []byte("invalid json"))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for empty request body", func() {
				req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", []byte(""))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for null request body", func() {
				req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", []byte("null"))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 409 for duplicate timeframe (PRD: unique timeframe per journal)", func() {
				// First image with DL timeframe already exists from journal creation
				image := barkat.Image{
					Timeframe: "DL",
					FileName:  "RELIANCE.mwd.duplicate.png",
				}
				req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/images", image)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusConflict))
			})
		})
	})

	// ============================================================================
	// 2.2.2 GET /v1/journals/{journal-id}/images - List Images
	// ============================================================================
	Describe("GET /v1/journals/{journal-id}/images - List Images (2.2.2)", func() {
		Context("Happy Path", func() {
			Context("with journal having images", func() {
				var images []barkat.Image

				BeforeEach(func() {
					req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/"+journal.ExternalID+"/images", nil)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
				})

				It("should return all images for journal", func() {
					images = decodeImageListResponse(w)
					Expect(images).To(HaveLen(4)) // 4 images from journal creation
				})

				It("should return images with correct timeframes", func() {
					images = decodeImageListResponse(w)
					timeframes := []string{}
					for _, img := range images {
						timeframes = append(timeframes, img.Timeframe)
					}
					Expect(timeframes).To(ContainElements("DL", "WK", "MN", "TMN"))
				})

				It("should return images with external IDs", func() {
					images = decodeImageListResponse(w)
					for _, img := range images {
						Expect(img.ExternalID).To(HavePrefix("img_"))
					}
				})

				It("should return images with created_at timestamps", func() {
					images = decodeImageListResponse(w)
					for _, img := range images {
						Expect(img.CreatedAt).ToNot(BeZero())
					}
				})
			})

			Context("with journal having no images", func() {
				var emptyJournal barkat.Journal

				BeforeEach(func() {
					// Create a new journal and manually delete its images
					emptyJournal = barkat.Journal{
						Ticker:   "EMPTY",
						Sequence: "YR",
						Type:     "SET",
						Status:   "TAKEN",
						Images: []barkat.Image{
							{Timeframe: "DL"},
							{Timeframe: "WK"},
							{Timeframe: "MN"},
							{Timeframe: "TMN"},
						},
					}
					Expect(journalMgr.CreateJournal(testCtx, &emptyJournal)).To(Succeed())

					// Delete all images
					for _, img := range emptyJournal.Images {
						err := imgMgr.DeleteImage(testCtx, emptyJournal.ExternalID, img.ExternalID)
						Expect(err).ToNot(HaveOccurred())
					}

					req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/"+emptyJournal.ExternalID+"/images", nil)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK with empty array", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
					images := decodeImageListResponse(w)
					Expect(images).To(BeEmpty())
				})
			})
		})

		Context("Field Validations", func() {
			Context("Journal ID Path Parameter", func() {
				Context("Bad Values", func() {
					It("should return 404 for non-existent journal ID", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/nonexistent-id/images", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for malformed journal ID", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/invalid-uuid-format/images", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for valid UUID format but non-existent", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/550e8400-e29b-41d4-a716-446655440000/images", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			// No server error scenarios for GET currently
		})
	})

	// ============================================================================
	// 2.2.3 DELETE /v1/journals/{journal-id}/images/{image-id} - Remove Image
	// ============================================================================
	Describe("DELETE /v1/journals/{journal-id}/images/{image-id} - Remove Image (2.2.3)", func() {
		var imageToDelete barkat.Image

		BeforeEach(func() {
			// Get first image from journal to delete
			imageToDelete = journal.Images[0]
		})

		Context("Happy Path", func() {
			Context("with valid journal and image IDs", func() {
				BeforeEach(func() {
					req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+journal.ExternalID+"/images/"+imageToDelete.ExternalID, nil)
					router.ServeHTTP(w, req)
				})

				It("should return 204 No Content", func() {
					Expect(w.Code).To(Equal(http.StatusNoContent))
				})

				It("should return empty body", func() {
					Expect(w.Body.String()).To(BeEmpty())
				})

				It("should actually delete the image from database", func() {
					images, err := imgMgr.ListImages(testCtx, journal.ExternalID)
					Expect(err).ToNot(HaveOccurred())
					Expect(images).To(HaveLen(3)) // 4 - 1 = 3
				})
			})
		})

		Context("Field Validations", func() {
			Context("Journal ID Path Parameter", func() {
				Context("Bad Values", func() {
					It("should return 404 for non-existent journal ID", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/nonexistent-id/images/"+imageToDelete.ExternalID, nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for malformed journal ID", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/invalid-uuid-format/images/"+imageToDelete.ExternalID, nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})

			Context("Image ID Path Parameter", func() {
				Context("Bad Values", func() {
					It("should return 404 for non-existent image ID", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+journal.ExternalID+"/images/nonexistent-image", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for malformed image ID", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+journal.ExternalID+"/images/invalid-uuid-format", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for valid UUID format but non-existent", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+journal.ExternalID+"/images/550e8400-e29b-41d4-a716-446655440000", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 404 on second delete (idempotency check)", func() {
				// First delete
				req1, w1 := util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+journal.ExternalID+"/images/"+imageToDelete.ExternalID, nil)
				router.ServeHTTP(w1, req1)
				Expect(w1.Code).To(Equal(http.StatusNoContent))

				// Second delete should return 404
				req2, w2 := util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+journal.ExternalID+"/images/"+imageToDelete.ExternalID, nil)
				router.ServeHTTP(w2, req2)
				Expect(w2.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
