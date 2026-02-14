package manager_test

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _ = Describe("JournalManager", func() {
	var (
		mgr     manager.JournalManager
		repo    repository.JournalRepository
		testCtx = context.Background()
		db      *gorm.DB
	)

	BeforeEach(func() {
		var err error
		db, err = util.CreateTestDb(logger.Warn)
		Expect(err).ToNot(HaveOccurred())
		Expect(core.SetupBarkatDB(db)).To(Succeed())
		repo = repository.NewJournalRepository(db)
		mgr = manager.NewJournalManager(repo)
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Context("Entry", func() {
		var entry barkat.Entry

		BeforeEach(func() {
			entry = barkat.Entry{
				Ticker: "GRSE", Sequence: "mwd", Type: "rejected", Status: "fail",
				Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}},
				Tags:   []barkat.Tag{{Tag: "tto", Type: "reason"}},
			}
			Expect(mgr.CreateEntry(testCtx, &entry)).To(BeNil())
		})

		It("should create entry with associations", func() {
			Expect(entry.ID).ToNot(BeEmpty())
			Expect(entry.Images).To(HaveLen(2))
			Expect(entry.Tags).To(HaveLen(1))
		})

		Context("GetEntry", func() {
			var fetched barkat.Entry

			BeforeEach(func() {
				var httpErr interface{ Code() int }
				fetched, httpErr = mgr.GetEntry(testCtx, entry.ID)
				Expect(httpErr).To(BeNil())
			})

			It("should return full entry with preloaded associations", func() {
				Expect(fetched.Ticker).To(Equal("GRSE"))
				Expect(fetched.Images).To(HaveLen(2))
				Expect(fetched.Tags).To(HaveLen(1))
			})
		})

		Context("GetEntry Not Found", func() {
			It("should return 404 for missing entry", func() {
				_, httpErr := mgr.GetEntry(testCtx, "nonexistent")
				Expect(httpErr).ToNot(BeNil())
				Expect(httpErr.Code()).To(Equal(404))
			})
		})

		Context("ListEntries", func() {
			BeforeEach(func() {
				second := barkat.Entry{
					Ticker: "DIXON", Sequence: "yr", Type: "set", Status: "taken",
				}
				Expect(mgr.CreateEntry(testCtx, &second)).To(BeNil())
			})

			It("should list with pagination metadata", func() {
				query := barkat.EntryQuery{}
				query.Limit = 10
				result, httpErr := mgr.ListEntries(testCtx, query)
				Expect(httpErr).To(BeNil())
				Expect(result.Records).To(HaveLen(2))
				Expect(result.Metadata.Total).To(Equal(int64(2)))
			})

			It("should filter by ticker", func() {
				query := barkat.EntryQuery{Ticker: "GRSE"}
				query.Limit = 10
				result, httpErr := mgr.ListEntries(testCtx, query)
				Expect(httpErr).To(BeNil())
				Expect(result.Records).To(HaveLen(1))
				Expect(result.Records[0].Ticker).To(Equal("GRSE"))
			})
		})

		Context("Image", func() {
			var image barkat.Image

			BeforeEach(func() {
				image = barkat.Image{Timeframe: "MN"}
				Expect(mgr.CreateImage(testCtx, entry.ID, &image)).To(BeNil())
			})

			It("should attach image to entry", func() {
				Expect(image.ID).ToNot(BeEmpty())
				Expect(image.EntryID).To(Equal(entry.ID))
			})

			Context("ListImages", func() {
				It("should list all images for entry", func() {
					images, httpErr := mgr.ListImages(testCtx, entry.ID)
					Expect(httpErr).To(BeNil())
					Expect(images).To(HaveLen(3))
				})
			})

			Context("DeleteImage", func() {
				BeforeEach(func() {
					Expect(mgr.DeleteImage(testCtx, entry.ID, image.ID)).To(BeNil())
				})

				It("should remove image", func() {
					images, httpErr := mgr.ListImages(testCtx, entry.ID)
					Expect(httpErr).To(BeNil())
					Expect(images).To(HaveLen(2))
				})
			})

			Context("CreateImage on missing entry", func() {
				It("should return 404", func() {
					img := barkat.Image{Timeframe: "DL"}
					httpErr := mgr.CreateImage(testCtx, "nonexistent", &img)
					Expect(httpErr).ToNot(BeNil())
					Expect(httpErr.Code()).To(Equal(404))
				})
			})

			Context("ListImages on missing entry", func() {
				It("should return 404", func() {
					_, httpErr := mgr.ListImages(testCtx, "nonexistent")
					Expect(httpErr).ToNot(BeNil())
					Expect(httpErr.Code()).To(Equal(404))
				})
			})
		})

		Context("Note", func() {
			var note barkat.Note

			BeforeEach(func() {
				note = barkat.Note{Status: "set", Content: "Trends\nHTF - Up"}
				Expect(mgr.CreateNote(testCtx, entry.ID, &note)).To(BeNil())
			})

			It("should attach note to entry", func() {
				Expect(note.ID).ToNot(BeEmpty())
				Expect(note.EntryID).To(Equal(entry.ID))
			})

			Context("ListNotes", func() {
				BeforeEach(func() {
					taken := barkat.Note{Status: "taken", Content: "Entered at 2450."}
					Expect(mgr.CreateNote(testCtx, entry.ID, &taken)).To(BeNil())
				})

				It("should list all notes", func() {
					notes, httpErr := mgr.ListNotes(testCtx, entry.ID, "")
					Expect(httpErr).To(BeNil())
					Expect(notes).To(HaveLen(2))
				})

				It("should filter by status", func() {
					notes, httpErr := mgr.ListNotes(testCtx, entry.ID, "taken")
					Expect(httpErr).To(BeNil())
					Expect(notes).To(HaveLen(1))
					Expect(notes[0].Status).To(Equal("taken"))
				})
			})

			Context("DeleteNote", func() {
				BeforeEach(func() {
					Expect(mgr.DeleteNote(testCtx, entry.ID, note.ID)).To(BeNil())
				})

				It("should remove note", func() {
					notes, httpErr := mgr.ListNotes(testCtx, entry.ID, "")
					Expect(httpErr).To(BeNil())
					Expect(notes).To(BeEmpty())
				})
			})

			Context("DeleteNote Not Found", func() {
				It("should return 404 for missing note", func() {
					httpErr := mgr.DeleteNote(testCtx, entry.ID, "nonexistent")
					Expect(httpErr).ToNot(BeNil())
					Expect(httpErr.Code()).To(Equal(404))
				})
			})

			Context("CreateNote on missing entry", func() {
				It("should return 404", func() {
					n := barkat.Note{Status: "set", Content: "test"}
					httpErr := mgr.CreateNote(testCtx, "nonexistent", &n)
					Expect(httpErr).ToNot(BeNil())
					Expect(httpErr.Code()).To(Equal(404))
				})
			})
		})

		Context("Tag", func() {
			var tag barkat.Tag

			BeforeEach(func() {
				tag = barkat.Tag{Tag: "enl", Type: "management"}
				Expect(mgr.CreateTag(testCtx, entry.ID, &tag)).To(BeNil())
			})

			It("should attach tag to entry", func() {
				Expect(tag.ID).ToNot(BeEmpty())
				Expect(tag.EntryID).To(Equal(entry.ID))
			})

			Context("ListTags", func() {
				BeforeEach(func() {
					reason := barkat.Tag{Tag: "dep", Type: "reason"}
					Expect(mgr.CreateTag(testCtx, entry.ID, &reason)).To(BeNil())
				})

				It("should list all tags", func() {
					tags, httpErr := mgr.ListTags(testCtx, entry.ID, "")
					Expect(httpErr).To(BeNil())
					Expect(tags).To(HaveLen(3))
				})

				It("should filter by type", func() {
					tags, httpErr := mgr.ListTags(testCtx, entry.ID, "management")
					Expect(httpErr).To(BeNil())
					Expect(tags).To(HaveLen(1))
					Expect(tags[0].Tag).To(Equal("enl"))
				})
			})

			Context("DeleteTag", func() {
				BeforeEach(func() {
					Expect(mgr.DeleteTag(testCtx, entry.ID, tag.ID)).To(BeNil())
				})

				It("should remove tag", func() {
					tags, httpErr := mgr.ListTags(testCtx, entry.ID, "")
					Expect(httpErr).To(BeNil())
					Expect(tags).To(HaveLen(1))
				})
			})

			Context("DeleteTag Not Found", func() {
				It("should return 404 for missing tag", func() {
					httpErr := mgr.DeleteTag(testCtx, entry.ID, "nonexistent")
					Expect(httpErr).ToNot(BeNil())
					Expect(httpErr.Code()).To(Equal(404))
				})
			})

			Context("CreateTag on missing entry", func() {
				It("should return 404", func() {
					t := barkat.Tag{Tag: "test", Type: "reason"}
					httpErr := mgr.CreateTag(testCtx, "nonexistent", &t)
					Expect(httpErr).ToNot(BeNil())
					Expect(httpErr.Code()).To(Equal(404))
				})
			})
		})
	})
})
