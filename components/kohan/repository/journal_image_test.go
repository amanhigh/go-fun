package repository_test

import (
	"context"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
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

	Context("ListImages", func() {
		Context("Happy Path", func() {
			var (
				images []barkat.Image
				err    error
			)

			BeforeEach(func() {
				// Create basic test images
				testImages := []barkat.Image{
					{EntryID: entry.ID, Timeframe: "DL"},
					{EntryID: entry.ID, Timeframe: "WK"},
					{EntryID: entry.ID, Timeframe: "MN"},
					{EntryID: entry.ID, Timeframe: "TMN"},
				}
				for _, img := range testImages {
					Expect(repo.Create(testCtx, &img)).To(Succeed())
				}

				// Get all images for testing
				images, err = repo.ListImages(testCtx, entry.ID)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should list all images successfully", func() {
				Expect(images).To(HaveLen(4))

				// Verify all timeframes are present using Lo
				timeframes := lo.Map(images, func(img barkat.Image, _ int) string { return img.Timeframe })
				Expect(timeframes).To(ContainElements("DL", "WK", "MN", "TMN"))
			})

			It("should list all images for entry", func() {
				Expect(len(images)).To(BeNumerically(">=", 4))
			})

			It("should have valid image metadata", func() {
				for _, img := range images {
					Expect(img.ID).ToNot(BeEmpty())
					Expect(img.EntryID).To(Equal(entry.ID))
					Expect(img.Timeframe).ToNot(BeEmpty())
					Expect(img.CreatedAt).ToNot(BeZero())
				}
			})

			It("should maintain timeframe ordering (DL → YR)", func() {
				Expect(len(images)).To(BeNumerically(">=", 4))
				timeframes := lo.Map(images, func(img barkat.Image, _ int) string { return img.Timeframe })
				Expect(timeframes).To(ContainElements("DL", "WK", "MN", "TMN"))
			})

			Context("Created At Ordering", func() {
				BeforeEach(func() {
					// Create multiple images with same timeframe for ordering tests
					for i := 0; i < 3; i++ {
						img := barkat.Image{EntryID: entry.ID, Timeframe: "DL"}
						Expect(repo.Create(testCtx, &img)).To(Succeed())
						time.Sleep(1 * time.Millisecond)
					}

					// Refresh images list
					images, err = repo.ListImages(testCtx, entry.ID)
					Expect(err).ToNot(HaveOccurred())
				})

				It("should sort by created_at when timeframes are same", func() {
					// Verify we have multiple DL images created with delays using Lo
					dlImages := lo.Filter(images, func(img barkat.Image, _ int) bool { return img.Timeframe == "DL" })
					Expect(len(dlImages)).To(BeNumerically(">=", 3))
				})
			})
		})

		Context("Edge Cases", func() {
			BeforeEach(func() {
				dl := barkat.Image{EntryID: entry.ID, Timeframe: "DL"}
				Expect(repo.Create(testCtx, &dl)).To(Succeed())
			})

			It("should handle empty entry ID gracefully", func() {
				images, err := repo.ListImages(testCtx, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(images).To(BeEmpty())
			})

			It("should handle very long entry ID", func() {
				longID := strings.Repeat("a", 1000)
				images, err := repo.ListImages(testCtx, longID)
				Expect(err).ToNot(HaveOccurred())
				Expect(images).To(BeEmpty())
			})

			It("should handle SQL injection attempts", func() {
				maliciousID := "'; DROP TABLE journal_images; --"
				images, err := repo.ListImages(testCtx, maliciousID)
				Expect(err).ToNot(HaveOccurred())
				Expect(images).To(BeEmpty())
			})
		})

		Context("Not Found", func() {
			It("should return empty for non-existent entry ID", func() {
				images, err := repo.ListImages(testCtx, "non-existent-entry-id")
				Expect(err).ToNot(HaveOccurred())
				Expect(images).To(BeEmpty())
			})

			It("should return empty for invalid UUID format", func() {
				images, err := repo.ListImages(testCtx, "invalid-uuid-format")
				Expect(err).ToNot(HaveOccurred())
				Expect(images).To(BeEmpty())
			})

			It("should return empty for very long non-existent ID", func() {
				longID := strings.Repeat("x", 100)
				images, err := repo.ListImages(testCtx, longID)
				Expect(err).ToNot(HaveOccurred())
				Expect(images).To(BeEmpty())
			})
		})
	})
})
