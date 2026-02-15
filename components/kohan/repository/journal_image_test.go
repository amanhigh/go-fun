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

		util.SetTxInContext(testCtx, db)
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
		BeforeEach(func() {
			dl := barkat.Image{EntryID: entry.ID, Timeframe: "DL"}
			Expect(repo.Create(testCtx, &dl)).To(Succeed())
			wk := barkat.Image{EntryID: entry.ID, Timeframe: "WK"}
			Expect(repo.Create(testCtx, &wk)).To(Succeed())
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
})
