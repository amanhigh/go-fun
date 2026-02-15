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

		image = barkat.Image{Timeframe: "DL"}
	})

	Context("CreateImage", func() {
		Context("with valid entry", func() {
			BeforeEach(func() {
				// Mock entry exists
				entryMgr.EXPECT().EntryExists(testCtx, "test-entry-id").Return(nil)
				// Mock repository create
				repo.EXPECT().UseOrCreateTx(mock.Anything, mock.AnythingOfType("util.DbRun")).Return(nil)
			})

			It("should create image successfully", func() {
				err := imgMgr.CreateImage(testCtx, "test-entry-id", &image)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("with non-existent entry", func() {
			BeforeEach(func() {
				// Mock entry not found
				entryMgr.EXPECT().EntryExists(testCtx, "nonexistent").Return(common.ErrNotFound)
			})

			It("should return 404 error", func() {
				err := imgMgr.CreateImage(testCtx, "nonexistent", &image)
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusNotFound))
			})
		})
	})

	Context("ListImages", func() {
		Context("with valid entry", func() {
			BeforeEach(func() {
				// Mock entry exists
				entryMgr.EXPECT().EntryExists(testCtx, "test-entry-id").Return(nil)
				// Mock repository list
				repo.EXPECT().UseOrCreateTx(mock.Anything, mock.AnythingOfType("util.DbRun")).Return(nil)
			})

			It("should list images successfully", func() {
				images, err := imgMgr.ListImages(testCtx, "test-entry-id")
				Expect(err).ToNot(HaveOccurred())
				// Note: With mocks, we can't easily test the actual return value
				// This would be better tested with integration tests
				_ = images
			})
		})

		Context("with unknown entry", func() {
			BeforeEach(func() {
				// Mock entry not found
				entryMgr.EXPECT().EntryExists(testCtx, "unknown-id").Return(common.ErrNotFound)
			})

			It("should return 404 error", func() {
				_, err := imgMgr.ListImages(testCtx, "unknown-id")
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusNotFound))
			})
		})
	})

	Context("DeleteImage", func() {
		Context("with valid entry and image", func() {
			BeforeEach(func() {
				// Mock entry exists
				entryMgr.EXPECT().EntryExists(testCtx, "test-entry-id").Return(nil)
				// Mock repository delete
				repo.EXPECT().DeleteById(testCtx, "test-image-id", mock.AnythingOfType("*barkat.Image")).Return(nil)
			})

			It("should delete image successfully", func() {
				err := imgMgr.DeleteImage(testCtx, "test-entry-id", "test-image-id")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("with non-existent entry", func() {
			BeforeEach(func() {
				// Mock entry not found
				entryMgr.EXPECT().EntryExists(testCtx, "nonexistent").Return(common.ErrNotFound)
			})

			It("should return 404 error", func() {
				err := imgMgr.DeleteImage(testCtx, "nonexistent", "test-image-id")
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusNotFound))
			})
		})

		Context("with repository error", func() {
			BeforeEach(func() {
				// Mock entry exists
				entryMgr.EXPECT().EntryExists(testCtx, "test-entry-id").Return(nil)
				// Mock repository error
				repo.EXPECT().DeleteById(testCtx, "nonexistent-image", mock.AnythingOfType("*barkat.Image")).Return(common.ErrNotFound)
			})

			It("should return repository error", func() {
				err := imgMgr.DeleteImage(testCtx, "test-entry-id", "nonexistent-image")
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusNotFound))
			})
		})
	})
})
