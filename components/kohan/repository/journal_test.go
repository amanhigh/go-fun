package repository_test

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _ = Describe("JournalRepository", func() {
	var (
		repo    repository.JournalRepository
		testCtx = context.Background()
		db      *gorm.DB
	)

	BeforeEach(func() {
		var err error
		db, err = util.CreateTestDb(logger.Warn)
		Expect(err).ToNot(HaveOccurred())

		err = core.SetupBarkatDB(db)
		Expect(err).ToNot(HaveOccurred())

		repo = repository.NewJournalRepository(db)
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Context("Create", func() {
		var (
			entry barkat.Entry
		)

		BeforeEach(func() {
			entry = barkat.Entry{
				Ticker:    "RELIANCE",
				Sequence:  "mwd",
				Type:      "rejected",
				Status:    "fail",
				CreatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				Images: []barkat.Image{
					{Timeframe: "DL"},
					{Timeframe: "WK"},
					{Timeframe: "MN"},
					{Timeframe: "TMN"},
				},
				Tags: []barkat.Tag{
					{Tag: "oe", Type: "reason"},
				},
				Notes: []barkat.Note{
					{Status: "rejected", Content: "Strong OE at weekly level."},
				},
			}

			err := repo.CreateEntry(testCtx, &entry)
			Expect(err).ToNot(HaveOccurred())
			Expect(entry.ID).ToNot(BeEmpty())
		})

		It("should create entry with generated ID", func() {
			Expect(entry.ID).ToNot(BeEmpty())
			Expect(entry.Ticker).To(Equal("RELIANCE"))
		})

		It("should create images with generated IDs", func() {
			Expect(entry.Images).To(HaveLen(4))
			for _, img := range entry.Images {
				Expect(img.ID).ToNot(BeEmpty())
				Expect(img.EntryID).To(Equal(entry.ID))
			}
		})

		Context("Get", func() {
			var fetchedEntry barkat.Entry

			BeforeEach(func() {
				var err error
				fetchedEntry, err = repo.GetEntry(testCtx, entry.ID)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should retrieve entry with associations", func() {
				Expect(fetchedEntry.ID).To(Equal(entry.ID))
				Expect(fetchedEntry.Ticker).To(Equal("RELIANCE"))
				Expect(fetchedEntry.Sequence).To(Equal("mwd"))
				Expect(fetchedEntry.Type).To(Equal("rejected"))
				Expect(fetchedEntry.Status).To(Equal("fail"))
				Expect(fetchedEntry.Images).To(HaveLen(4))
				Expect(fetchedEntry.Tags).To(HaveLen(1))
				Expect(fetchedEntry.Tags[0].Tag).To(Equal("oe"))
				Expect(fetchedEntry.Notes).To(HaveLen(1))
				Expect(fetchedEntry.Notes[0].Content).To(Equal("Strong OE at weekly level."))
			})

			It("should preserve image timeframes", func() {
				timeframes := []string{"DL", "WK", "MN", "TMN"}
				for _, img := range fetchedEntry.Images {
					Expect(timeframes).To(ContainElement(img.Timeframe))
				}
			})
		})

		Context("Get Not Found", func() {
			It("should return error for missing ID", func() {
				_, err := repo.GetEntry(testCtx, "nonexistent-id")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("List", func() {
			var secondEntry barkat.Entry

			BeforeEach(func() {
				secondEntry = barkat.Entry{
					Ticker:    "INFY",
					Sequence:  "yr",
					Type:      "set",
					Status:    "taken",
					CreatedAt: time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
					Images: []barkat.Image{
						{Timeframe: "DL"},
					},
					Notes: []barkat.Note{
						{Status: "set", Content: "Trends\nHTF - Up\nMTF - Up"},
					},
				}

				err := repo.CreateEntry(testCtx, &secondEntry)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should list all entries with pagination", func() {
				query := barkat.EntryQuery{}
				query.Limit = 10
				entries, total, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(Equal(int64(2)))
				Expect(entries).To(HaveLen(2))
				// Ordered by created_at DESC — INFY first
				Expect(entries[0].Ticker).To(Equal("INFY"))
				Expect(entries[1].Ticker).To(Equal("RELIANCE"))
			})

			It("should filter by ticker", func() {
				query := barkat.EntryQuery{Ticker: "INFY"}
				query.Limit = 10
				entries, total, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(Equal(int64(1)))
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].Ticker).To(Equal("INFY"))
			})

			It("should filter by type", func() {
				query := barkat.EntryQuery{Type: "rejected"}
				query.Limit = 10
				entries, total, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(Equal(int64(1)))
				Expect(entries[0].Ticker).To(Equal("RELIANCE"))
			})

			It("should filter by status", func() {
				query := barkat.EntryQuery{Status: "taken"}
				query.Limit = 10
				entries, total, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(Equal(int64(1)))
				Expect(entries[0].Ticker).To(Equal("INFY"))
			})

			It("should filter by sequence", func() {
				query := barkat.EntryQuery{Sequence: "yr"}
				query.Limit = 10
				entries, total, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(Equal(int64(1)))
				Expect(entries[0].Ticker).To(Equal("INFY"))
			})

			It("should paginate with offset", func() {
				query := barkat.EntryQuery{}
				query.Limit = 1
				query.Offset = 0
				entries, total, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(Equal(int64(2)))
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].Ticker).To(Equal("INFY"))

				// Second page
				query.Offset = 1
				entries, total, err = repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(Equal(int64(2)))
				Expect(entries).To(HaveLen(1))
				Expect(entries[0].Ticker).To(Equal("RELIANCE"))
			})

			It("should return lightweight summaries without associations", func() {
				query := barkat.EntryQuery{Ticker: "RELIANCE"}
				query.Limit = 10
				entries, _, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(entries[0].Ticker).To(Equal("RELIANCE"))
				Expect(entries[0].Images).To(BeEmpty())
				Expect(entries[0].Tags).To(BeEmpty())
				Expect(entries[0].Notes).To(BeEmpty())
			})

			It("should return empty for no matches", func() {
				query := barkat.EntryQuery{Ticker: "NONEXISTENT"}
				query.Limit = 10
				entries, total, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(Equal(int64(0)))
				Expect(entries).To(BeEmpty())
			})
		})
	})
})
