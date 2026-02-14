package repository_test

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/common/util"
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

		err = db.AutoMigrate(&barkat.Entry{}, &barkat.Image{})
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
				Outcome:   "fail",
				Trend:     "trend",
				CreatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				Images: []barkat.Image{
					{Position: 1, Path: "assets/trading/2024/01/RELIANCE.mwd.trend.rejected__20240115__100000.png"},
					{Position: 2, Path: "assets/trading/2024/01/RELIANCE.mwd.trend.rejected__20240115__100001.png"},
					{Position: 3, Path: "assets/trading/2024/01/RELIANCE.mwd.trend.rejected__20240115__100002.png"},
					{Position: 4, Path: "assets/trading/2024/01/RELIANCE.mwd.trend.rejected__20240115__100003.png"},
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

			It("should retrieve entry with images", func() {
				Expect(fetchedEntry.ID).To(Equal(entry.ID))
				Expect(fetchedEntry.Ticker).To(Equal("RELIANCE"))
				Expect(fetchedEntry.Sequence).To(Equal("mwd"))
				Expect(fetchedEntry.Type).To(Equal("rejected"))
				Expect(fetchedEntry.Outcome).To(Equal("fail"))
				Expect(fetchedEntry.Trend).To(Equal("trend"))
				Expect(fetchedEntry.Images).To(HaveLen(4))
			})

			It("should preserve image positions", func() {
				for i, img := range fetchedEntry.Images {
					Expect(img.Position).To(Equal(i + 1))
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
				notes := "Trends\nHTF - Up\nMTF - Up"
				secondEntry = barkat.Entry{
					Ticker:        "INFY",
					Sequence:      "yr",
					Type:          "set",
					Outcome:       "taken",
					Trend:         "trend",
					NotesMarkdown: &notes,
					CreatedAt:     time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
					Images: []barkat.Image{
						{Position: 1, Path: "assets/trading/2024/02/INFY.yr.set__20240210__100000.png"},
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

			It("should filter by outcome", func() {
				query := barkat.EntryQuery{Outcome: "taken"}
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

			It("should preload images on list", func() {
				query := barkat.EntryQuery{Ticker: "RELIANCE"}
				query.Limit = 10
				entries, _, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(entries[0].Images).To(HaveLen(4))
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
