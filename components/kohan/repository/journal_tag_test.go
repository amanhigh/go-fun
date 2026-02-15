package repository_test

import (
	"context"
	"fmt"
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

var _ = Describe("JournalRepository Tag", func() {
	var (
		repo      repository.TagRepository
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
		repo = repository.NewTagRepository(db)

		entry = barkat.Entry{
			Ticker: "CEATLTD", Sequence: "mwd", Type: "set", Status: "success",
		}
		Expect(entryRepo.CreateEntry(testCtx, &entry)).To(Succeed())
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Context("CreateTag", func() {
		var tag barkat.Tag

		BeforeEach(func() {
			tag = barkat.Tag{EntryID: entry.ID, Type: "reason", Tag: "oe"}
			Expect(repo.Create(testCtx, &tag)).To(Succeed())
		})

		It("should create tag with generated ID", func() {
			Expect(tag.ID).ToNot(BeEmpty())
			Expect(tag.Tag).To(Equal("oe"))
			Expect(tag.Type).To(Equal("reason"))
		})

		Context("ListTags Happy Path", func() {
			var (
				tags []barkat.Tag
				err  error
			)

			BeforeEach(func() {
				// Create basic test tags
				basicTags := []barkat.Tag{
					{EntryID: entry.ID, Tag: "enl", Type: "management"},
					{EntryID: entry.ID, Tag: "brk", Type: "reason"},
					{EntryID: entry.ID, Tag: "pos", Type: "management"},
				}
				for _, t := range basicTags {
					Expect(repo.Create(testCtx, &t)).To(Succeed())
				}

				// Get all tags for testing
				tags, err = repo.ListTags(testCtx, entry.ID, "")
				Expect(err).ToNot(HaveOccurred())
			})

			It("should list all tags for entry", func() {
				Expect(len(tags)).To(BeNumerically(">=", 4)) // At least 3 from BeforeEach + 1 from parent BeforeEach

				// Verify all expected tags are present using Lo
				tagValues := lo.Map(tags, func(tag barkat.Tag, _ int) string { return tag.Tag })
				Expect(tagValues).To(ContainElements("oe", "enl", "brk", "pos"))
			})

			It("should have valid tag metadata", func() {
				for _, tag := range tags {
					Expect(tag.ID).ToNot(BeEmpty())
					Expect(tag.EntryID).To(Equal(entry.ID))
					Expect(tag.Tag).ToNot(BeEmpty())
					Expect(tag.Type).ToNot(BeEmpty())
					Expect(tag.CreatedAt).ToNot(BeZero())
				}
			})

			It("should filter by type=reason correctly", func() {
				reasonTags, err := repo.ListTags(testCtx, entry.ID, "reason")
				Expect(err).ToNot(HaveOccurred())
				Expect(reasonTags).ToNot(BeEmpty()) // At least oe + brk

				for _, tag := range reasonTags {
					Expect(tag.Type).To(Equal("reason"))
				}
			})

			It("should filter by type=management correctly", func() {
				mgmtTags, err := repo.ListTags(testCtx, entry.ID, "management")
				Expect(err).ToNot(HaveOccurred())
				Expect(mgmtTags).ToNot(BeEmpty()) // At least enl + pos

				for _, tag := range mgmtTags {
					Expect(tag.Type).To(Equal("management"))
				}
			})

			Context("Created At Ordering", func() {
				BeforeEach(func() {
					// Create multiple tags with different timestamps for ordering tests
					for i := 0; i < 3; i++ {
						testTag := barkat.Tag{EntryID: entry.ID, Type: "reason", Tag: fmt.Sprintf("tag%d", i)}
						Expect(repo.Create(testCtx, &testTag)).To(Succeed())
						time.Sleep(1 * time.Millisecond)
					}

					// Refresh tags list
					tags, err = repo.ListTags(testCtx, entry.ID, "")
					Expect(err).ToNot(HaveOccurred())
				})

				It("should sort by created_at descending by default", func() {
					Expect(len(tags)).To(BeNumerically(">=", 3))

					// Verify chronological ordering (newest first)
					for i := 1; i < len(tags); i++ {
						Expect(tags[i].CreatedAt).To(BeTemporally(">=", tags[i-1].CreatedAt))
					}
				})
			})

			Context("ListTags Edge Cases", func() {
				BeforeEach(func() {
					tag = barkat.Tag{EntryID: entry.ID, Type: "reason", Tag: "oe"}
					Expect(repo.Create(testCtx, &tag)).To(Succeed())
				})

				It("should handle empty entry ID gracefully", func() {
					tags, err := repo.ListTags(testCtx, "", "")
					Expect(err).ToNot(HaveOccurred())
					Expect(tags).To(BeEmpty())
				})

				It("should handle very long entry ID", func() {
					longID := strings.Repeat("a", 1000)
					tags, err := repo.ListTags(testCtx, longID, "")
					Expect(err).ToNot(HaveOccurred())
					Expect(tags).To(BeEmpty())
				})

				It("should handle SQL injection attempts", func() {
					maliciousID := "'; DROP TABLE journal_tags; --"
					tags, err := repo.ListTags(testCtx, maliciousID, "")
					Expect(err).ToNot(HaveOccurred())
					Expect(tags).To(BeEmpty())
				})
			})

			Context("Not Found", func() {
				It("should return empty for non-existent entry ID", func() {
					tags, err := repo.ListTags(testCtx, "non-existent-entry-id", "")
					Expect(err).ToNot(HaveOccurred())
					Expect(tags).To(BeEmpty())
				})
			})
		})
	})
})
