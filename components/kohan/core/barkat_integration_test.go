package core_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testPort = 19020

var baseURL string

// httpDo is a helper for making HTTP requests with method, url, and optional JSON body.
func httpDo(method, url string, body any) (*http.Response, []byte) {
	var reqBody *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		Expect(err).ToNot(HaveOccurred())
		reqBody = bytes.NewReader(b)
	} else {
		reqBody = bytes.NewReader(nil)
	}
	req, err := http.NewRequest(method, url, reqBody)
	Expect(err).ToNot(HaveOccurred())
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	Expect(err).ToNot(HaveOccurred())
	respBody, err := io.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())
	resp.Body.Close()
	return resp, respBody
}

// httpDoEntry is a helper for making HTTP requests that return an Entry envelope response.
func httpDoEntry(method, url string, body any) barkat.Journal {
	resp, responseBody := httpDo(method, url, body)
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	var envelope common.Envelope[barkat.Journal]
	Expect(json.Unmarshal(responseBody, &envelope)).To(Succeed())
	Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))

	return envelope.Data
}

var _ = PDescribe("Barkat Integration Test", func() {
	// TODO: #B Decide Scenarios for Integration Test and Fix it.
	BeforeEach(func() {
		if baseURL == "" {
			db, err := core.CreateTestBarkatDB()
			Expect(err).ToNot(HaveOccurred())

			entryRepo := repository.NewJournalRepository(db)
			entryMgr := manager.NewJournalManager(entryRepo)
			journalHandler := handler.NewJournalHandler(entryMgr)
			imageHandler := handler.NewImageHandler(manager.NewImageManager(entryMgr, repository.NewImageRepository(db)))
			noteHandler := handler.NewNoteHandler(manager.NewNoteManager(entryMgr, repository.NewNoteRepository(db)))
			tagHandler := handler.NewTagHandler(manager.NewTagManager(entryMgr, repository.NewTagRepository(db)))

			shutdown := util.NewGracefulShutdown()
			base := util.NewHttpServer(config.HttpServerConfig{Name: "kohan", Port: testPort}, gin.Default(), shutdown)
			lifecycle := core.NewKohanServerLifecycle(nil, journalHandler, imageHandler, noteHandler, tagHandler)
			base.SetLifecycle(lifecycle)
			baseURL = fmt.Sprintf("http://localhost:%d", testPort)

			go func() {
				defer GinkgoRecover()
				_ = base.Start()
			}()

			Eventually(func() error {
				_, err := http.Get(baseURL + "/v1/journal-entries")
				return err
			}, 5*time.Second, 100*time.Millisecond).Should(Succeed())
		}
	})

	// ---- Real Production Data: GRSE rejected/fail with reason tag tto-loc (2023-06-15) ----
	Context("GRSE Rejected Entry", func() {
		var createdEntry barkat.Journal

		BeforeEach(func() {
			entry := barkat.Journal{
				Ticker:   "GRSE",
				Sequence: "MWD",
				Type:     "REJECTED",
				Status:   "FAIL",
				Images: []barkat.Image{
					{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"},
				},
				Tags: []barkat.Tag{
					{Tag: "tto", Type: "reason", Override: new("loc")},
				},
			}
			resp, body := httpDo("POST", baseURL+"/v1/journal-entries", entry)
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			// Handle the envelope response for creation
			var envelope common.Envelope[barkat.Journal]
			Expect(json.Unmarshal(body, &envelope)).To(Succeed())
			Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			createdEntry = envelope.Data
		})

		It("should create with all fields", func() {
			Expect(createdEntry.ID).ToNot(BeEmpty())
			Expect(createdEntry.Ticker).To(Equal("GRSE"))
			Expect(createdEntry.Status).To(Equal("FAIL"))
			Expect(createdEntry.Images).To(HaveLen(4))
			Expect(createdEntry.Tags).To(HaveLen(1))
			Expect(createdEntry.Tags[0].Tag).To(Equal("tto"))
			Expect(*createdEntry.Tags[0].Override).To(Equal("loc"))
		})

		Context("Get", func() {
			var fetchedEntry barkat.Journal

			BeforeEach(func() {
				fetchedEntry = httpDoEntry("GET", baseURL+"/v1/journal-entries/"+createdEntry.ID, nil)
			})

			It("should return full entry with associations", func() {
				Expect(fetchedEntry.ID).To(Equal(createdEntry.ID))
				Expect(fetchedEntry.Ticker).To(Equal("GRSE"))
				Expect(fetchedEntry.Images).To(HaveLen(4))
				Expect(fetchedEntry.Tags).To(HaveLen(1))
			})
		})
	})

	// ---- Real Production Data: DIXON set/taken with notes (2023-06-15) ----
	Context("DIXON Set Entry with Notes", func() {
		var createdEntry barkat.Journal

		BeforeEach(func() {
			entry := barkat.Journal{
				Ticker:   "DIXON",
				Sequence: "MWD",
				Type:     "SET",
				Status:   "TAKEN",
				Images: []barkat.Image{
					{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"},
				},
				Notes: []barkat.Note{
					{Status: "set", Content: "Trends\nMN - DN\nWK - Up\nD1 - Up\n\nPlan: Shorts @ WK SZ nested in MN SZ", Format: "markdown"},
				},
			}
			resp, body := httpDo("POST", baseURL+"/v1/journal-entries", entry)
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			// Handle the envelope response for creation
			var envelope common.Envelope[barkat.Journal]
			Expect(json.Unmarshal(body, &envelope)).To(Succeed())
			Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			createdEntry = envelope.Data
		})

		It("should create with inline note", func() {
			Expect(createdEntry.Notes).To(HaveLen(1))
			Expect(createdEntry.Notes[0].Content).To(ContainSubstring("Shorts @ WK SZ"))
			Expect(createdEntry.Notes[0].Format).To(Equal("markdown"))
		})

		Context("Add Note via API", func() {
			var addedNote barkat.Note

			BeforeEach(func() {
				note := barkat.Note{Status: "taken", Content: "Entered at 2450, SL at 2420."}
				resp, body := httpDo("POST", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/notes", note)
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				// Sub-resource handlers return direct responses, not enveloped
				Expect(json.Unmarshal(body, &addedNote)).To(Succeed())
			})

			It("should attach note to entry", func() {
				Expect(addedNote.ID).ToNot(BeEmpty())
				Expect(addedNote.JournalID).To(Equal(createdEntry.ID))
				Expect(addedNote.Status).To(Equal("taken"))
			})

			Context("List Notes", func() {
				It("should list all notes", func() {
					resp, body := httpDo("GET", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/notes", nil)
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					var result map[string][]barkat.Note
					Expect(json.Unmarshal(body, &result)).To(Succeed())
					Expect(result["notes"]).To(HaveLen(2))
				})

				It("should filter notes by status", func() {
					resp, body := httpDo("GET", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/notes?note_status=taken", nil)
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					var result map[string][]barkat.Note
					Expect(json.Unmarshal(body, &result)).To(Succeed())
					Expect(result["notes"]).To(HaveLen(1))
					Expect(result["notes"][0].Status).To(Equal("taken"))
				})
			})

			Context("Delete Note", func() {
				BeforeEach(func() {
					resp, _ := httpDo("DELETE", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/notes/"+addedNote.ID, nil)
					Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
				})

				It("should remove note from entry", func() {
					resp, body := httpDo("GET", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/notes", nil)
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					var result map[string][]barkat.Note
					Expect(json.Unmarshal(body, &result)).To(Succeed())
					Expect(result["notes"]).To(HaveLen(1))
				})
			})
		})
	})

	// ---- Real Production Data: CEATLTD set/success with management tags (2025-08-21) ----
	Context("CEATLTD Success with Management Tags", func() {
		var createdEntry barkat.Journal

		BeforeEach(func() {
			entry := barkat.Journal{
				Ticker:   "CEATLTD",
				Sequence: "MWD",
				Type:     "SET",
				Status:   "SUCCESS",
				Images: []barkat.Image{
					{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"},
				},
				Tags: []barkat.Tag{
					{Tag: "enl", Type: "management"},
					{Tag: "ntr", Type: "management"},
				},
				Notes: []barkat.Note{
					{Status: "set", Content: "Trends\nHTF - Up\nMTF - Up\nTTF - Up\n\nPlan: Longs @ WK DZ\n\nSupport:\n- MN EMA", Format: "markdown"},
				},
			}
			resp, body := httpDo("POST", baseURL+"/v1/journal-entries", entry)
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			// Handle the envelope response for creation
			var envelope common.Envelope[barkat.Journal]
			Expect(json.Unmarshal(body, &envelope)).To(Succeed())
			Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			createdEntry = envelope.Data
		})

		It("should create with management tags", func() {
			Expect(createdEntry.Tags).To(HaveLen(2))
			Expect(createdEntry.Tags[0].Type).To(Equal("management"))
		})

		Context("Tag Sub-resource APIs", func() {
			var addedTag barkat.Tag

			BeforeEach(func() {
				tag := barkat.Tag{Tag: "er", Type: "reason"}
				resp, body := httpDo("POST", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/tags", tag)
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				Expect(json.Unmarshal(body, &addedTag)).To(Succeed())
			})

			It("should add tag to entry", func() {
				Expect(addedTag.ID).ToNot(BeEmpty())
				Expect(addedTag.Tag).To(Equal("er"))
				Expect(addedTag.Type).To(Equal("reason"))
			})

			Context("List Tags", func() {
				It("should list all tags", func() {
					resp, body := httpDo("GET", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/tags", nil)
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					var result map[string][]barkat.Tag
					Expect(json.Unmarshal(body, &result)).To(Succeed())
					Expect(result["tags"]).To(HaveLen(3))
				})

				It("should filter by type=management", func() {
					resp, body := httpDo("GET", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/tags?type=management", nil)
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					var result map[string][]barkat.Tag
					Expect(json.Unmarshal(body, &result)).To(Succeed())
					Expect(result["tags"]).To(HaveLen(2))
				})

				It("should filter by type=reason", func() {
					resp, body := httpDo("GET", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/tags?type=reason", nil)
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					var result map[string][]barkat.Tag
					Expect(json.Unmarshal(body, &result)).To(Succeed())
					Expect(result["tags"]).To(HaveLen(1))
					Expect(result["tags"][0].Tag).To(Equal("er"))
				})
			})

			Context("Delete Tag", func() {
				BeforeEach(func() {
					resp, _ := httpDo("DELETE", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/tags/"+addedTag.ID, nil)
					Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
				})

				It("should remove tag from entry", func() {
					resp, body := httpDo("GET", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/tags", nil)
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					var result map[string][]barkat.Tag
					Expect(json.Unmarshal(body, &result)).To(Succeed())
					Expect(result["tags"]).To(HaveLen(2))
				})
			})
		})

		Context("Image Sub-resource APIs", func() {
			var addedImage barkat.Image

			BeforeEach(func() {
				image := barkat.Image{Timeframe: "SMN"}
				resp, body := httpDo("POST", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/images", image)
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				Expect(json.Unmarshal(body, &addedImage)).To(Succeed())
			})

			It("should add image to entry", func() {
				Expect(addedImage.ID).ToNot(BeEmpty())
				Expect(addedImage.Timeframe).To(Equal("SMN"))
			})

			Context("List Images", func() {
				It("should list all images", func() {
					resp, body := httpDo("GET", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/images", nil)
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					var result map[string][]barkat.Image
					Expect(json.Unmarshal(body, &result)).To(Succeed())
					Expect(result["images"]).To(HaveLen(5))
				})
			})

			Context("Delete Image", func() {
				BeforeEach(func() {
					resp, _ := httpDo("DELETE", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/images/"+addedImage.ID, nil)
					Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
				})

				It("should remove image from entry", func() {
					resp, body := httpDo("GET", baseURL+"/v1/journal-entries/"+createdEntry.ID+"/images", nil)
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					var result map[string][]barkat.Image
					Expect(json.Unmarshal(body, &result)).To(Succeed())
					Expect(result["images"]).To(HaveLen(4))
				})
			})
		})
	})

	// ---- List with Filters (multiple entries from production patterns) ----
	Context("List with Filters", func() {
		BeforeEach(func() {
			entries := []barkat.Journal{
				{
					Ticker: "KEI", Sequence: "MWD", Type: "REJECTED", Status: "FAIL",
					Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}},
					Tags:   []barkat.Tag{{Tag: "dep", Type: "reason"}},
				},
				{
					Ticker: "SJVN", Sequence: "MWD", Type: "REJECTED", Status: "FAIL",
					Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}},
					Tags:   []barkat.Tag{{Tag: "zn", Type: "reason", Override: new("big")}},
				},
				{
					Ticker: "PDSL", Sequence: "YR", Type: "SET", Status: "RUNNING",
					Images: []barkat.Image{{Timeframe: "DL"}, {Timeframe: "WK"}, {Timeframe: "MN"}, {Timeframe: "TMN"}},
					Notes:  []barkat.Note{{Status: "set", Content: "Trends\nHTF - Up\nMTF - Up\nTTF - Up\n\nPlan: Longs @ TTF DZ", Format: "markdown"}},
				},
			}
			for i := range entries {
				resp, _ := httpDo("POST", baseURL+"/v1/journal-entries", entries[i])
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			}
		})

		It("should filter by ticker", func() {
			resp, body := httpDo("GET", baseURL+"/v1/journal-entries?ticker=KEI&limit=10", nil)
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			var envelope common.Envelope[barkat.JournalList]
			Expect(json.Unmarshal(body, &envelope)).To(Succeed())
			Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			Expect(envelope.Data.Records).To(HaveLen(1))
			Expect(envelope.Data.Records[0].Ticker).To(Equal("KEI"))
		})

		It("should filter by sequence=yr", func() {
			resp, body := httpDo("GET", baseURL+"/v1/journal-entries?sequence=YR&limit=10", nil)
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			var envelope common.Envelope[barkat.JournalList]
			Expect(json.Unmarshal(body, &envelope)).To(Succeed())
			Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			for _, e := range envelope.Data.Records {
				Expect(e.Sequence).To(Equal("YR"))
			}
		})

		It("should filter by status=running", func() {
			resp, body := httpDo("GET", baseURL+"/v1/journal-entries?status=RUNNING&limit=10", nil)
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			var envelope common.Envelope[barkat.JournalList]
			Expect(json.Unmarshal(body, &envelope)).To(Succeed())
			Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			for _, e := range envelope.Data.Records {
				Expect(e.Status).To(Equal("running"))
			}
		})

		It("should return lightweight summaries without associations", func() {
			resp, body := httpDo("GET", baseURL+"/v1/journal-entries?limit=10", nil)
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			var envelope common.Envelope[barkat.JournalList]
			Expect(json.Unmarshal(body, &envelope)).To(Succeed())
			Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			Expect(envelope.Data.Records).ToNot(BeEmpty())
			for _, e := range envelope.Data.Records {
				Expect(e.Images).To(BeEmpty())
				Expect(e.Tags).To(BeEmpty())
				Expect(e.Notes).To(BeEmpty())
			}
		})
	})

	// ---- Error Cases ----
	Context("Error Cases", func() {
		It("should return 404 for missing entry", func() {
			resp, _ := httpDo("GET", baseURL+"/v1/journal-entries/nonexistent-id", nil)
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should return 404 for sub-resource on missing entry", func() {
			resp, _ := httpDo("GET", baseURL+"/v1/journal-entries/nonexistent-id/images", nil)
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should return 404 for deleting missing image", func() {
			resp, _ := httpDo("DELETE", baseURL+"/v1/journal-entries/nonexistent-id/images/missing-img", nil)
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should return 404 for deleting missing note", func() {
			resp, _ := httpDo("DELETE", baseURL+"/v1/journal-entries/nonexistent-id/notes/missing-note", nil)
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should return 404 for deleting missing tag", func() {
			resp, _ := httpDo("DELETE", baseURL+"/v1/journal-entries/nonexistent-id/tags/missing-tag", nil)
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
