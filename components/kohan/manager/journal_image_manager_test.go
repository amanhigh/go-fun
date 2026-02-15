package manager_test

import (
	"context"
	"net/http"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	managerMocks "github.com/amanhigh/go-fun/components/kohan/manager/mocks"
	repoMocks "github.com/amanhigh/go-fun/components/kohan/repository/mocks"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("ImageManager", func() {
	var (
		entryMgr *managerMocks.JournalManager
		imgMgr   manager.ImageManager
		repo     *repoMocks.ImageRepository
		testCtx  = context.Background()
		image    barkat.Image
	)

	BeforeEach(func() {
		entryMgr = managerMocks.NewJournalManager(GinkgoT())
		repo = repoMocks.NewImageRepository(GinkgoT())
		imgMgr = manager.NewImageManager(entryMgr, repo)

		image = barkat.Image{
			Timeframe: "DL", // Valid timeframe from PRD spec
		}
	})

	Context("CreateImage", func() {
		Context("Happy Path", func() {
			BeforeEach(func() {
				// Mock entry exists
				entryMgr.EXPECT().EntryExists(testCtx, "test-entry-id").Return(nil)
				// Mock repository create within transaction
				repo.EXPECT().UseOrCreateTx(mock.Anything, mock.AnythingOfType("util.DbRun")).Return(nil)
			})

			It("should create image successfully with proper contract", func() {
				createdImage, err := imgMgr.CreateImage(testCtx, "test-entry-id", image)

				// Verify successful creation
				Expect(err).ToNot(HaveOccurred())
				Expect(createdImage).ToNot(BeNil())

				// Verify API contract compliance (FR-001, FR-004)
				Expect(createdImage.EntryID).To(Equal("test-entry-id"))
				Expect(createdImage.Timeframe).To(Equal("DL"))
				// Note: In mock environment, ID and CreatedAt are not actually set by BeforeCreate hook
				// In real implementation, these would be set by GORM hooks
			})
		})

		Context("Edge Cases", func() {
			Context("with repository error", func() {
				BeforeEach(func() {
					// Mock entry exists
					entryMgr.EXPECT().EntryExists(testCtx, "test-entry-id").Return(nil)
					// Mock repository transaction error
					repo.EXPECT().UseOrCreateTx(mock.Anything, mock.AnythingOfType("util.DbRun")).Return(&common.HttpErrorImpl{
						Msg:     "Database constraint violation",
						ErrCode: http.StatusInternalServerError,
					})
				})

				It("should return repository error and nil image", func() {
					createdImage, err := imgMgr.CreateImage(testCtx, "test-entry-id", image)
					Expect(err).To(HaveOccurred())
					Expect(err.Code()).To(Equal(http.StatusInternalServerError))
					Expect(err.Error()).To(Equal("Database constraint violation"))
					Expect(createdImage).To(BeNil()) // Critical: should not return partial data on error
				})
			})

			Context("with invalid timeframe", func() {
				BeforeEach(func() {
					// Mock entry exists
					entryMgr.EXPECT().EntryExists(testCtx, "test-entry-id").Return(nil)
					// Mock repository validation error
					repo.EXPECT().UseOrCreateTx(mock.Anything, mock.AnythingOfType("util.DbRun")).Return(&common.HttpErrorImpl{
						Msg:     "Invalid timeframe: must be one of DL, YR, MN, WK, D",
						ErrCode: http.StatusBadRequest,
					})
				})

				It("should validate timeframe and return 400", func() {
					invalidImage := barkat.Image{Timeframe: "INVALID"}
					createdImage, err := imgMgr.CreateImage(testCtx, "test-entry-id", invalidImage)
					Expect(err).To(HaveOccurred())
					Expect(err.Code()).To(Equal(http.StatusBadRequest))
					Expect(err.Error()).To(ContainSubstring("Invalid timeframe"))
					Expect(createdImage).To(BeNil())
				})
			})

			Context("with empty entry ID", func() {
				BeforeEach(func() {
					// Mock entry not found for empty ID
					entryMgr.EXPECT().EntryExists(testCtx, "").Return(common.ErrNotFound)
				})
				It("should handle empty entry ID gracefully", func() {
					createdImage, err := imgMgr.CreateImage(testCtx, "", image)
					Expect(err).To(HaveOccurred())
					Expect(err.Code()).To(Equal(http.StatusNotFound))
					Expect(createdImage).To(BeNil())
				})
			})
		})

		Context("Not Found", func() {
			BeforeEach(func() {
				// Mock entry not found
				entryMgr.EXPECT().EntryExists(testCtx, "nonexistent-entry").Return(common.ErrNotFound)
			})

			It("should return 404 for non-existent entry and nil image", func() {
				createdImage, err := imgMgr.CreateImage(testCtx, "nonexistent-entry", image)
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusNotFound))
				Expect(err.Error()).To(Equal("NotFound"))
				Expect(createdImage).To(BeNil()) // Critical: should not return partial data
			})
		})
	})

	Context("ListImages", func() {
		var (
			images []barkat.Image
		)

		Context("Happy Path", func() {
			BeforeEach(func() {
				// Mock entry exists
				entryMgr.EXPECT().EntryExists(testCtx, "test-entry-id").Return(nil)
				// Mock repository list within transaction
				repo.EXPECT().UseOrCreateTx(mock.Anything, mock.AnythingOfType("util.DbRun")).Return(nil)
			})

			It("should list images successfully with proper contract", func() {
				var err common.HttpError
				images, err = imgMgr.ListImages(testCtx, "test-entry-id")

				// Verify successful listing
				Expect(err).ToNot(HaveOccurred())
				Expect(images).To(BeEmpty())

				// Note: In real implementation, images would have proper timeframe metadata (FR-002)
			})
		})

		Context("Edge Cases", func() {
			Context("with repository error", func() {
				BeforeEach(func() {
					// Mock entry exists
					entryMgr.EXPECT().EntryExists(testCtx, "test-entry-id").Return(nil)
					// Mock repository transaction error
					repo.EXPECT().UseOrCreateTx(mock.Anything, mock.AnythingOfType("util.DbRun")).Return(&common.HttpErrorImpl{
						Msg:     "Database connection lost",
						ErrCode: http.StatusInternalServerError,
					})
				})

				It("should return repository error and nil images", func() {
					images, err := imgMgr.ListImages(testCtx, "test-entry-id")
					Expect(err).To(HaveOccurred())
					Expect(err.Code()).To(Equal(http.StatusInternalServerError))
					Expect(err.Error()).To(Equal("Database connection lost"))
					Expect(images).To(BeNil()) // Critical: should not return partial data on error
				})
			})

			Context("with empty entry ID", func() {
				BeforeEach(func() {
					// Mock entry not found for empty ID
					entryMgr.EXPECT().EntryExists(testCtx, "").Return(common.ErrNotFound)
				})
				It("should handle empty entry ID gracefully", func() {
					images, err := imgMgr.ListImages(testCtx, "")
					Expect(err).To(HaveOccurred())
					Expect(err.Code()).To(Equal(http.StatusNotFound))
					Expect(images).To(BeNil())
				})
			})
		})

		Context("Not Found", func() {
			BeforeEach(func() {
				// Mock entry not found
				entryMgr.EXPECT().EntryExists(testCtx, "nonexistent-entry").Return(common.ErrNotFound)
			})

			It("should return 404 for non-existent entry and nil images", func() {
				images, err := imgMgr.ListImages(testCtx, "nonexistent-entry")
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusNotFound))
				Expect(err.Error()).To(Equal("NotFound"))
				Expect(images).To(BeNil()) // Critical: should not return partial data
			})
		})
	})

	Context("DeleteImage", func() {
		Context("Happy Path", func() {
			BeforeEach(func() {
				// Mock entry exists
				entryMgr.EXPECT().EntryExists(testCtx, "test-entry-id").Return(nil)
				// Mock repository delete within transaction
				repo.EXPECT().UseOrCreateTx(mock.Anything, mock.AnythingOfType("util.DbRun")).Return(nil)
			})

			It("should delete image successfully with proper contract", func() {
				err := imgMgr.DeleteImage(testCtx, "test-entry-id", "test-image-id")

				// Verify successful deletion (FR-004.5)
				Expect(err).ToNot(HaveOccurred())
				// Delete operations should return nil on success, not 204 (that's HTTP layer)
			})
		})

		Context("Edge Cases", func() {
			Context("with repository error", func() {
				BeforeEach(func() {
					// Mock entry exists
					entryMgr.EXPECT().EntryExists(testCtx, "test-entry-id").Return(nil)
					// Mock repository transaction error (image not found)
					repo.EXPECT().UseOrCreateTx(mock.Anything, mock.AnythingOfType("util.DbRun")).Return(&common.HttpErrorImpl{
						Msg:     "Image not found",
						ErrCode: http.StatusNotFound,
					})
				})

				It("should return 404 when image does not exist", func() {
					err := imgMgr.DeleteImage(testCtx, "test-entry-id", "nonexistent-image")
					Expect(err).To(HaveOccurred())
					Expect(err.Code()).To(Equal(http.StatusNotFound))
					Expect(err.Error()).To(Equal("Image not found"))
				})
			})

			Context("with empty entry ID", func() {
				BeforeEach(func() {
					// Mock entry not found for empty ID
					entryMgr.EXPECT().EntryExists(testCtx, "").Return(common.ErrNotFound)
				})
				It("should handle empty entry ID gracefully", func() {
					err := imgMgr.DeleteImage(testCtx, "", "test-image-id")
					Expect(err).To(HaveOccurred())
					Expect(err.Code()).To(Equal(http.StatusNotFound))
				})
			})
		})

		Context("Not Found", func() {
			BeforeEach(func() {
				// Mock entry not found
				entryMgr.EXPECT().EntryExists(testCtx, "nonexistent-entry").Return(common.ErrNotFound)
			})

			It("should return 404 for non-existent entry", func() {
				err := imgMgr.DeleteImage(testCtx, "nonexistent-entry", "test-image-id")
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusNotFound))
				Expect(err.Error()).To(Equal("NotFound"))
			})
		})
	})
})
