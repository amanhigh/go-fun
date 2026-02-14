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
			tag = barkat.Tag{EntryID: entry.ID, Tag: "oe", Type: "reason"}
			Expect(repo.CreateTag(testCtx, &tag)).To(Succeed())
		})

		It("should create tag with generated ID", func() {
			Expect(tag.ID).ToNot(BeEmpty())
			Expect(tag.Tag).To(Equal("oe"))
			Expect(tag.Type).To(Equal("reason"))
		})

		Context("ListTags", func() {
			BeforeEach(func() {
				mgmt := barkat.Tag{EntryID: entry.ID, Tag: "enl", Type: "management"}
				Expect(repo.CreateTag(testCtx, &mgmt)).To(Succeed())
			})

			It("should list all tags for entry", func() {
				tags, err := repo.ListTags(testCtx, entry.ID, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(tags).To(HaveLen(2))
			})

			It("should filter by type=reason", func() {
				tags, err := repo.ListTags(testCtx, entry.ID, "reason")
				Expect(err).ToNot(HaveOccurred())
				Expect(tags).To(HaveLen(1))
				Expect(tags[0].Tag).To(Equal("oe"))
			})

			It("should filter by type=management", func() {
				tags, err := repo.ListTags(testCtx, entry.ID, "management")
				Expect(err).ToNot(HaveOccurred())
				Expect(tags).To(HaveLen(1))
				Expect(tags[0].Tag).To(Equal("enl"))
			})

			It("should return empty for unknown entry", func() {
				tags, err := repo.ListTags(testCtx, "unknown-id", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(tags).To(BeEmpty())
			})
		})

		Context("DeleteTag", func() {
			BeforeEach(func() {
				Expect(repo.DeleteTag(testCtx, entry.ID, tag.ID)).To(Succeed())
			})

			It("should remove tag", func() {
				tags, err := repo.ListTags(testCtx, entry.ID, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(tags).To(BeEmpty())
			})
		})

		Context("DeleteTag Not Found", func() {
			It("should return error for missing tag", func() {
				err := repo.DeleteTag(testCtx, entry.ID, "nonexistent")
				Expect(err).To(HaveOccurred())
			})

			It("should return error for wrong entry scope", func() {
				err := repo.DeleteTag(testCtx, "wrong-entry", tag.ID)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
