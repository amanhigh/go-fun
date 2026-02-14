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

var _ = Describe("JournalRepository Note", func() {
	var (
		repo      repository.NoteRepository
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
		repo = repository.NewNoteRepository(db)

		entry = barkat.Entry{
			Ticker: "DIXON", Sequence: "mwd", Type: "set", Status: "taken",
		}
		Expect(entryRepo.CreateEntry(testCtx, &entry)).To(Succeed())
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Context("CreateNote", func() {
		var note barkat.Note

		BeforeEach(func() {
			note = barkat.Note{EntryID: entry.ID, Status: "set", Content: "Trends\nHTF - Up\nMTF - Up"}
			Expect(repo.CreateNote(testCtx, &note)).To(Succeed())
		})

		It("should create note with generated ID and default format", func() {
			Expect(note.ID).ToNot(BeEmpty())
			Expect(note.Format).To(Equal("markdown"))
			Expect(note.Content).To(ContainSubstring("HTF - Up"))
		})

		Context("ListNotes", func() {
			BeforeEach(func() {
				taken := barkat.Note{EntryID: entry.ID, Status: "taken", Content: "Entered at 2450."}
				Expect(repo.CreateNote(testCtx, &taken)).To(Succeed())
			})

			It("should list all notes for entry", func() {
				notes, err := repo.ListNotes(testCtx, entry.ID, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(notes).To(HaveLen(2))
			})

			It("should filter by status", func() {
				notes, err := repo.ListNotes(testCtx, entry.ID, "taken")
				Expect(err).ToNot(HaveOccurred())
				Expect(notes).To(HaveLen(1))
				Expect(notes[0].Status).To(Equal("taken"))
			})

			It("should return empty for unknown entry", func() {
				notes, err := repo.ListNotes(testCtx, "unknown-id", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(notes).To(BeEmpty())
			})
		})

		Context("DeleteNote", func() {
			BeforeEach(func() {
				Expect(repo.DeleteNote(testCtx, entry.ID, note.ID)).To(Succeed())
			})

			It("should remove note", func() {
				notes, err := repo.ListNotes(testCtx, entry.ID, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(notes).To(BeEmpty())
			})
		})

		Context("DeleteNote Not Found", func() {
			It("should return error for missing note", func() {
				err := repo.DeleteNote(testCtx, entry.ID, "nonexistent")
				Expect(err).To(HaveOccurred())
			})

			It("should return error for wrong entry scope", func() {
				err := repo.DeleteNote(testCtx, "wrong-entry", note.ID)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
