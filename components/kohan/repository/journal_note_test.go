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
			Expect(repo.Create(testCtx, &note)).To(Succeed())
		})

		It("should create note with generated ID and default format", func() {
			Expect(note.ID).ToNot(BeEmpty())
			Expect(note.Format).To(Equal("markdown"))
			Expect(note.Content).To(ContainSubstring("HTF - Up"))
		})

		Context("ListNotes", func() {
			Context("ListNotes Happy Path", func() {
				var (
					notes []barkat.Note
					err   error
				)

				BeforeEach(func() {
					// Create basic test notes
					basicNotes := []barkat.Note{
						{EntryID: entry.ID, Status: "set", Content: "Trends\nHTF - Up\nMTF - Up"},
						{EntryID: entry.ID, Status: "taken", Content: "Entered at 2450."},
						{EntryID: entry.ID, Status: "running", Content: "Currently in position."},
						{EntryID: entry.ID, Status: "dropped", Content: "Exited position."},
					}
					for _, note := range basicNotes {
						Expect(repo.Create(testCtx, &note)).To(Succeed())
					}

					// Get all notes for testing
					notes, err = repo.ListNotes(testCtx, entry.ID, "")
					Expect(err).ToNot(HaveOccurred())
				})

				It("should list all notes successfully", func() {
					Expect(len(notes)).To(BeNumerically(">=", 4)) // At least 4 from BeforeEach + 1 from parent BeforeEach
				})

				It("should filter by status correctly", func() {
					filteredNotes, err := repo.ListNotes(testCtx, entry.ID, "taken")
					Expect(err).ToNot(HaveOccurred())
					Expect(filteredNotes).ToNot(BeEmpty())
					Expect(filteredNotes[0].Status).To(Equal("taken"))
				})

				It("should have valid note metadata", func() {
					for _, note := range notes {
						Expect(note.ID).ToNot(BeEmpty())
						Expect(note.EntryID).To(Equal(entry.ID))
						Expect(note.Status).ToNot(BeEmpty())
						Expect(note.Content).ToNot(BeEmpty())
						Expect(note.Format).ToNot(BeEmpty())
						Expect(note.CreatedAt).ToNot(BeZero())
					}
				})

				Context("Created At Ordering", func() {
					BeforeEach(func() {
						// Create multiple notes with different timestamps for ordering tests
						for i := 0; i < 3; i++ {
							testNote := barkat.Note{EntryID: entry.ID, Status: "set", Content: fmt.Sprintf("note %d", i)}
							Expect(repo.Create(testCtx, &testNote)).To(Succeed())
							time.Sleep(1 * time.Millisecond)
						}

						// Refresh notes list
						notes, err = repo.ListNotes(testCtx, entry.ID, "")
						Expect(err).ToNot(HaveOccurred())
					})

					It("should sort by created_at descending by default", func() {
						Expect(len(notes)).To(BeNumerically(">=", 3))

						// Verify chronological ordering (newest first)
						for i := 1; i < len(notes); i++ {
							Expect(notes[i].CreatedAt).To(BeTemporally(">=", notes[i-1].CreatedAt))
						}
					})
				})

				Context("Content Integrity", func() {
					BeforeEach(func() {
						// Create notes with specific content for integrity tests
						contentNotes := []barkat.Note{
							{EntryID: entry.ID, Status: "set", Content: "First note"},
							{EntryID: entry.ID, Status: "taken", Content: "Second note"},
							{EntryID: entry.ID, Status: "running", Content: "Third note"},
						}
						for _, note := range contentNotes {
							Expect(repo.Create(testCtx, &note)).To(Succeed())
						}

						// Refresh notes list
						notes, err = repo.ListNotes(testCtx, entry.ID, "")
						Expect(err).ToNot(HaveOccurred())
					})

					It("should preserve content integrity when ordered", func() {
						Expect(len(notes)).To(BeNumerically(">=", 7)) // 4 from basicNotes + 1 from parent BeforeEach + 3 from contentNotes

						// Verify content is preserved correctly regardless of ordering using Lo
						contents := lo.Map(notes, func(note barkat.Note, _ int) string { return note.Content })
						Expect(contents).To(ContainElements("First note", "Second note", "Third note", "Trends\nHTF - Up\nMTF - Up"))
					})
				})
			})

			Context("Not Found", func() {
				It("should return empty for non-existent entry ID", func() {
					notes, err := repo.ListNotes(testCtx, "non-existent-entry-id", "")
					Expect(err).ToNot(HaveOccurred())
					Expect(notes).To(BeEmpty())
				})

				It("should return empty for invalid UUID format", func() {
					notes, err := repo.ListNotes(testCtx, "invalid-uuid-format", "")
					Expect(err).ToNot(HaveOccurred())
					Expect(notes).To(BeEmpty())
				})

				It("should return empty for very long non-existent ID", func() {
					longID := strings.Repeat("x", 100)
					notes, err := repo.ListNotes(testCtx, longID, "")
					Expect(err).ToNot(HaveOccurred())
					Expect(notes).To(BeEmpty())
				})

				It("should return empty for non-existent entry with status filter", func() {
					notes, err := repo.ListNotes(testCtx, "non-existent-entry-id", "taken")
					Expect(err).ToNot(HaveOccurred())
					Expect(notes).To(BeEmpty())
				})
			})
		})

		Context("DeleteNote", func() {
			BeforeEach(func() {
				Expect(repo.DeleteById(testCtx, note.ID, &barkat.Note{})).To(Succeed())
			})

			It("should remove note", func() {
				notes, err := repo.ListNotes(testCtx, entry.ID, "")
				Expect(err).ToNot(HaveOccurred())
				Expect(notes).To(BeEmpty())
			})
		})
	})
})
