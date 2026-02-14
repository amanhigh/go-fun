package repository_test

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _ = Describe("JournalRepository Image", func() {
	var (
		repo      repository.ImageRepository
		entryRepo repository.JournalRepository
		testCtx   = context.Background()
		db        *gorm.DB
		entry     barkat.Entry
	)

	BeforeEach(func() {
		var err error
		db, err = util.CreateTestDb(logger.Warn)
		Expect(err).ToNot(HaveOccurred())
		Expect(core.SetupBarkatDB(db)).To(Succeed())
		entryRepo = repository.NewJournalRepository(db)
		repo = repository.NewImageRepository(db)

		entry = barkat.Entry{
			Ticker: "RELIANCE", Sequence: "mwd", Type: "rejected", Status: "fail",
		}
		Expect(entryRepo.CreateEntry(testCtx, &entry)).To(Succeed())
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Context("CreateImage", func() {
		var image barkat.Image

		BeforeEach(func() {
			image = barkat.Image{EntryID: entry.ID, Timeframe: "DL"}
			Expect(repo.CreateImage(testCtx, &image)).To(Succeed())
		})

		It("should create image with generated ID", func() {
			Expect(image.ID).ToNot(BeEmpty())
			Expect(image.EntryID).To(Equal(entry.ID))
			Expect(image.Timeframe).To(Equal("DL"))
		})

		Context("ListImages", func() {
			BeforeEach(func() {
				wk := barkat.Image{EntryID: entry.ID, Timeframe: "WK"}
				Expect(repo.CreateImage(testCtx, &wk)).To(Succeed())
			})

			It("should list all images for entry", func() {
				images, err := repo.ListImages(testCtx, entry.ID)
				Expect(err).ToNot(HaveOccurred())
				Expect(images).To(HaveLen(2))
			})

			It("should return empty for unknown entry", func() {
				images, err := repo.ListImages(testCtx, "unknown-id")
				Expect(err).ToNot(HaveOccurred())
				Expect(images).To(BeEmpty())
			})
		})

		Context("DeleteImage", func() {
			BeforeEach(func() {
				Expect(repo.DeleteImage(testCtx, entry.ID, image.ID)).To(Succeed())
			})

			It("should remove image", func() {
				images, err := repo.ListImages(testCtx, entry.ID)
				Expect(err).ToNot(HaveOccurred())
				Expect(images).To(BeEmpty())
			})
		})

		Context("DeleteImage Not Found", func() {
			It("should return error for missing image", func() {
				err := repo.DeleteImage(testCtx, entry.ID, "nonexistent")
				Expect(err).To(HaveOccurred())
			})

			It("should return error for wrong entry scope", func() {
				err := repo.DeleteImage(testCtx, "wrong-entry", image.ID)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
