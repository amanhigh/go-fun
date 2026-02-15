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
		Context("Happy Path", func() {
			var entry barkat.Entry

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

			It("should create tags with generated IDs", func() {
				Expect(entry.Tags).To(HaveLen(1))
				for _, tag := range entry.Tags {
					Expect(tag.ID).ToNot(BeEmpty())
					Expect(tag.EntryID).To(Equal(entry.ID))
				}
			})

			It("should create notes with generated IDs", func() {
				Expect(entry.Notes).To(HaveLen(1))
				for _, note := range entry.Notes {
					Expect(note.ID).ToNot(BeEmpty())
					Expect(note.EntryID).To(Equal(entry.ID))
				}
			})
		})

		Context("Edge Cases", func() {
			It("should create entry with minimal data", func() {
				minimalEntry := barkat.Entry{
					Ticker:   "MIN",
					Sequence: "d",
					Type:     "set",
					Status:   "taken",
				}
				err := repo.CreateEntry(testCtx, &minimalEntry)
				Expect(err).ToNot(HaveOccurred())
				Expect(minimalEntry.ID).ToNot(BeEmpty())
			})

			It("should create entry with empty associations", func() {
				emptyEntry := barkat.Entry{
					Ticker:   "EMPTY",
					Sequence: "d",
					Type:     "set",
					Status:   "taken",
					Images:   []barkat.Image{},
					Tags:     []barkat.Tag{},
					Notes:    []barkat.Note{},
				}
				err := repo.CreateEntry(testCtx, &emptyEntry)
				Expect(err).ToNot(HaveOccurred())
				Expect(emptyEntry.ID).ToNot(BeEmpty())
			})
		})
	})

	Context("Get", func() {
		Context("Happy Path", func() {
			var fetchedEntry barkat.Entry
			var testEntry barkat.Entry

			BeforeEach(func() {
				// Create an entry first for testing
				testEntry = barkat.Entry{
					Ticker:    "RELIANCE",
					Sequence:  "mwd",
					Type:      "rejected",
					Status:    "fail",
					CreatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					Images: []barkat.Image{
						{Timeframe: "DL"},
						{Timeframe: "WK"},
					},
					Tags: []barkat.Tag{
						{Tag: "oe", Type: "reason"},
					},
					Notes: []barkat.Note{
						{Status: "rejected", Content: "Strong OE at weekly level."},
					},
				}
				createErr := repo.CreateEntry(testCtx, &testEntry)
				Expect(createErr).ToNot(HaveOccurred())

				var getErr error
				fetchedEntry, getErr = repo.GetEntry(testCtx, testEntry.ID)
				Expect(getErr).ToNot(HaveOccurred())
			})

			It("should retrieve entry with associations", func() {
				Expect(fetchedEntry.ID).To(Equal(testEntry.ID))
				Expect(fetchedEntry.Ticker).To(Equal("RELIANCE"))
				Expect(fetchedEntry.Sequence).To(Equal("mwd"))
				Expect(fetchedEntry.Type).To(Equal("rejected"))
				Expect(fetchedEntry.Status).To(Equal("fail"))
				Expect(fetchedEntry.Images).To(HaveLen(2))
				Expect(fetchedEntry.Tags).To(HaveLen(1))
				Expect(fetchedEntry.Tags[0].Tag).To(Equal("oe"))
				Expect(fetchedEntry.Notes).To(HaveLen(1))
				Expect(fetchedEntry.Notes[0].Content).To(Equal("Strong OE at weekly level."))
			})

			It("should preserve image timeframes", func() {
				timeframes := []string{"DL", "WK"}
				for _, img := range fetchedEntry.Images {
					Expect(timeframes).To(ContainElement(img.Timeframe))
				}
			})

			It("should have valid timestamps", func() {
				Expect(fetchedEntry.CreatedAt).ToNot(BeZero())
			})
		})

		Context("Not Found", func() {
			It("should return error for missing ID", func() {
				_, err := repo.GetEntry(testCtx, "nonexistent-id")
				Expect(err).To(HaveOccurred())
			})

			It("should return error for empty ID", func() {
				_, err := repo.GetEntry(testCtx, "")
				Expect(err).To(HaveOccurred())
			})

			It("should return error for invalid UUID format", func() {
				_, err := repo.GetEntry(testCtx, "invalid-uuid")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("List", func() {
		Context("Happy Path", func() {
			var secondEntry barkat.Entry
			var firstEntry barkat.Entry

			BeforeEach(func() {
				// Create entries first for testing
				firstEntry = barkat.Entry{
					Ticker:    "RELIANCE",
					Sequence:  "mwd",
					Type:      "rejected",
					Status:    "fail",
					CreatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					Images: []barkat.Image{
						{Timeframe: "DL"},
						{Timeframe: "WK"},
					},
					Tags: []barkat.Tag{
						{Tag: "oe", Type: "reason"},
					},
					Notes: []barkat.Note{
						{Status: "rejected", Content: "Strong OE at weekly level."},
					},
				}
				createErr := repo.CreateEntry(testCtx, &firstEntry)
				Expect(createErr).ToNot(HaveOccurred())

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

			It("should sort by created_at descending by default", func() {
				query := barkat.EntryQuery{}
				query.Limit = 10
				entries, _, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(entries)).To(BeNumerically(">=", 2))
				// INFY (Feb 10) should come before RELIANCE (Jan 15) when sorted descending
				Expect(entries[0].Ticker).To(Equal("INFY"))
				Expect(entries[1].Ticker).To(Equal("RELIANCE"))
			})

			It("should sort by created_at ascending", func() {
				query := barkat.EntryQuery{}
				query.Limit = 10
				query.SortBy = "created_at"
				query.SortOrder = "asc"
				entries, _, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(entries)).To(BeNumerically(">=", 2))
				// RELIANCE (Jan 15) should come before INFY (Feb 10) when sorted ascending
				Expect(entries[0].Ticker).To(Equal("RELIANCE"))
				Expect(entries[1].Ticker).To(Equal("INFY"))
			})

			It("should sort by ticker ascending", func() {
				query := barkat.EntryQuery{}
				query.Limit = 10
				query.SortBy = "ticker"
				query.SortOrder = "asc"
				entries, _, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(entries)).To(BeNumerically(">=", 2))
				// INFY should come before RELIANCE alphabetically
				Expect(entries[0].Ticker).To(Equal("INFY"))
				Expect(entries[1].Ticker).To(Equal("RELIANCE"))
			})

			It("should sort by ticker descending", func() {
				query := barkat.EntryQuery{}
				query.Limit = 10
				query.SortBy = "ticker"
				query.SortOrder = "desc"
				entries, _, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(entries)).To(BeNumerically(">=", 2))
				// RELIANCE should come before INFY alphabetically when descending
				Expect(entries[0].Ticker).To(Equal("RELIANCE"))
				Expect(entries[1].Ticker).To(Equal("INFY"))
			})

			It("should sort by sequence ascending", func() {
				query := barkat.EntryQuery{}
				query.Limit = 10
				query.SortBy = "sequence"
				query.SortOrder = "asc"
				entries, _, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(entries)).To(BeNumerically(">=", 2))
				// mwd should come before yr alphabetically
				Expect(entries[0].Sequence).To(Equal("mwd"))
				Expect(entries[1].Sequence).To(Equal("yr"))
			})

			It("should sort by sequence descending", func() {
				query := barkat.EntryQuery{}
				query.Limit = 10
				query.SortBy = "sequence"
				query.SortOrder = "desc"
				entries, _, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(entries)).To(BeNumerically(">=", 2))
				// yr should come before mwd alphabetically when descending
				Expect(entries[0].Sequence).To(Equal("yr"))
				Expect(entries[1].Sequence).To(Equal("mwd"))
			})
		})

		Context("Edge Cases", func() {
			It("should handle empty query with no data", func() {
				query := barkat.EntryQuery{}
				entries, total, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(Equal(int64(0)))
				Expect(entries).To(BeEmpty())
			})

			It("should handle very large limit", func() {
				query := barkat.EntryQuery{}
				query.Limit = 1000
				_, total, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(BeNumerically(">=", 0))
			})
		})

		Context("Not Found", func() {
			It("should return empty for no matches", func() {
				query := barkat.EntryQuery{Ticker: "NONEXISTENT"}
				query.Limit = 10
				entries, total, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(Equal(int64(0)))
				Expect(entries).To(BeEmpty())
			})

			It("should return empty for filters with no matches", func() {
				query := barkat.EntryQuery{Ticker: "NONEXISTENT", Type: "set"}
				query.Limit = 10
				entries, total, err := repo.ListEntries(testCtx, query)
				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(Equal(int64(0)))
				Expect(entries).To(BeEmpty())
			})
		})
	})
})
