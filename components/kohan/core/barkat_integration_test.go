//nolint:dupl
package core_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/go-resty/resty/v2"
	"github.com/golang-sql/civil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// standardImages provides 4 required images for journal creation
var standardImages = []barkat.Image{
	{Timeframe: "DL", FileName: "daily.png"},
	{Timeframe: "WK", FileName: "weekly.png"},
	{Timeframe: "MN", FileName: "monthly.png"},
	{Timeframe: "TMN", FileName: "trend_monthly.png"},
}

// decodeJournalResponse unmarshals response body into Journal envelope
func decodeJournalResponse(resp *resty.Response) barkat.Journal {
	var envelope common.Envelope[barkat.Journal]
	ExpectWithOffset(1, json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
	ExpectWithOffset(1, envelope.Status).To(Equal(common.EnvelopeSuccess))
	return envelope.Data
}

// Barkat E2E Test Suite
//
// Tests critical paths through real HTTP server with in-memory SQLite DB.
// Focuses on scenarios that add value beyond unit/integration tests:
// - Full CRUD lifecycle with associations
// - Cascade delete (FK constraints)
// - Validation through real HTTP stack
// - Review status workflow
//
// Server is started/stopped in core_suite_test.go BeforeSuite/AfterSuite
var _ = Describe("Barkat E2E Test", func() {
	var client *resty.Client

	BeforeEach(func() {
		// Create fresh client for each test - server is managed by core_suite_test.go
		client = resty.New()
		client.SetTimeout(5 * time.Second)
		client.SetHeader("Content-Type", "application/json")
		client.SetBaseURL(fmt.Sprintf("http://localhost:%d", testPort))
	})

	Context("Portal", func() {
		It("should render the index page", func() {
			resp, err := client.R().Get("/")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			Expect(resp.Header().Get("Content-Type")).To(ContainSubstring("text/html"))
			Expect(resp.String()).To(ContainSubstring("Shadow Gate"))
			Expect(resp.String()).To(ContainSubstring("Welcome to the Kohan portal."))
		})

		It("should render the journal page", func() {
			resp, err := client.R().Get("/journal")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			Expect(resp.Header().Get("Content-Type")).To(ContainSubstring("text/html"))
			Expect(resp.String()).To(ContainSubstring("Journal"))
			Expect(resp.String()).To(ContainSubstring("Kohan Portal"))
		})

		It("should render journal detail page", func() {
			resp, err := client.R().Get("/journal/jrn_1234abcd")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			Expect(resp.Header().Get("Content-Type")).To(ContainSubstring("text/html"))
			Expect(resp.String()).To(ContainSubstring("Journal Detail"))
			Expect(resp.String()).To(ContainSubstring("jrn_1234abcd"))
		})

		It("should serve journal images from static route", func() {
			resp, err := client.R().Get("/journal/images/2024/01/sample.png")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			Expect(resp.String()).To(Equal("sample-image"))
		})
	})

	// Admin Endpoints - Tests server administration functionality
	Context("Admin Endpoints", func() {
		It("should handle health endpoint", func() {
			resp, err := client.R().Get("/health")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			Expect(resp.String()).To(ContainSubstring("ok"))
		})
	})

	// OS Endpoints - Tests system-level operations
	Context("OS Endpoints", func() {
		It("should handle OS endpoint", func() {
			resp, err := client.R().SetBody(map[string]string{"submap": "test"}).
				Post("/v1/api/os/submap/disable")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var result map[string]any
			Expect(json.Unmarshal(resp.Body(), &result)).To(Succeed())
			Expect(result["status"]).To(Equal("success"))
			Expect(result["action"]).To(Equal("disable"))
		})

		It("should handle ticker recording endpoint", func() {
			resp, err := client.R().Get("/v1/api/os/ticker/AAPL/record")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			Expect(resp.String()).To(ContainSubstring("Success"))
		})
	})

	// Full CRUD Lifecycle - Tests complete flow through real HTTP + DB
	Context("Journal CRUD Lifecycle", func() {
		var createdJournal barkat.Journal

		BeforeEach(func() {
			journal := barkat.Journal{
				Ticker:   "LIFECYCLE",
				Sequence: "MWD",
				Type:     "SET",
				Status:   "RUNNING",
				Images:   standardImages,
				Tags:     []barkat.Tag{{Tag: "oe", Type: "REASON"}},
				Notes:    []barkat.Note{{Status: "SET", Content: "Initial setup note", Format: "MARKDOWN"}},
			}
			resp, err := client.R().SetBody(journal).Post(barkat.JournalBase)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
			createdJournal = decodeJournalResponse(resp)
		})

		It("should create journal with all associations", func() {
			Expect(createdJournal.ExternalID).To(HavePrefix("jrn_"))
			Expect(createdJournal.Images).To(HaveLen(4))
			Expect(createdJournal.Tags).To(HaveLen(1))
			Expect(createdJournal.Notes).To(HaveLen(1))

			// Verify association IDs
			for _, img := range createdJournal.Images {
				Expect(img.ExternalID).To(HavePrefix("img_"))
			}
			Expect(createdJournal.Tags[0].ExternalID).To(HavePrefix("tag_"))
			Expect(createdJournal.Notes[0].ExternalID).To(HavePrefix("not_"))
		})

		It("should retrieve journal with associations", func() {
			resp, err := client.R().Get(barkat.JournalBase + "/" + createdJournal.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			fetched := decodeJournalResponse(resp)
			Expect(fetched.Ticker).To(Equal("LIFECYCLE"))
			Expect(fetched.Images).To(HaveLen(4))
			Expect(fetched.Tags).To(HaveLen(1))
			Expect(fetched.Notes).To(HaveLen(1))
		})

		It("should list entries with pagination", func() {
			resp, err := client.R().Get(barkat.JournalBase + "?limit=10")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var envelope common.Envelope[barkat.JournalList]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			Expect(envelope.Data.Journals).ToNot(BeEmpty())
		})

		It("should delete journal and cascade to associations", func() {
			resp, err := client.R().Delete(barkat.JournalBase + "/" + createdJournal.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusNoContent))

			// Verify journal is deleted
			resp, err = client.R().Get(barkat.JournalBase + "/" + createdJournal.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusNotFound))
		})
	})

	// Image Management - Tests FR-003: Journal Image Management
	Context("Image Management", func() {
		var createdJournal barkat.Journal
		var createdImage barkat.Image

		BeforeEach(func() {
			journal := barkat.Journal{
				Ticker:   "IMGTEST",
				Sequence: "MWD",
				Type:     "SET",
				Status:   "RUNNING",
				Images:   standardImages,
			}
			resp, err := client.R().SetBody(journal).Post(barkat.JournalBase)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
			createdJournal = decodeJournalResponse(resp)
			createdImage = createdJournal.Images[0]
		})

		It("should add additional timeframe image", func() {
			newImage := barkat.Image{Timeframe: "WK", FileName: "weekly.png"}
			resp, err := client.R().SetBody(newImage).Post(barkat.JournalBase + "/" + createdJournal.ExternalID + "/images")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))

			var envelope common.Envelope[barkat.Image]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			addedImage := envelope.Data
			Expect(addedImage.ExternalID).To(HavePrefix("img_"))
			Expect(addedImage.Timeframe).To(Equal("WK"))
		})

		It("should list all images for journal", func() {
			resp, err := client.R().Get(barkat.JournalBase + "/" + createdJournal.ExternalID + "/images")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var envelope common.Envelope[barkat.ImageList]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			images := envelope.Data.Images
			Expect(images).To(HaveLen(4)) // standardImages has 4 images
			// Check for standard timeframes
			timeframes := make(map[string]bool)
			for _, img := range images {
				timeframes[img.Timeframe] = true
			}
			Expect(timeframes["DL"]).To(BeTrue())
			Expect(timeframes["WK"]).To(BeTrue())
			Expect(timeframes["MN"]).To(BeTrue())
			Expect(timeframes["TMN"]).To(BeTrue())
		})

		It("should delete individual image", func() {
			resp, err := client.R().Delete(barkat.JournalBase + "/" + createdJournal.ExternalID + "/images/" + createdImage.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusNoContent))

			// Verify image is deleted (should have 3 remaining)
			resp, err = client.R().Get(barkat.JournalBase + "/" + createdJournal.ExternalID + "/images")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var envelope common.Envelope[barkat.ImageList]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			Expect(envelope.Data.Images).To(HaveLen(3)) // Started with 4, deleted 1
		})
	})

	// Note Management - Tests FR-004: Journal Note Management
	Context("Note Management", func() {
		var createdJournal barkat.Journal
		var createdNote barkat.Note

		BeforeEach(func() {
			journal := barkat.Journal{
				Ticker:   "NOTETEST",
				Sequence: "YR",
				Type:     "RESULT",
				Status:   "SUCCESS",
				Images:   standardImages,
			}
			resp, err := client.R().SetBody(journal).Post(barkat.JournalBase)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
			createdJournal = decodeJournalResponse(resp)

			// Create a note for use in tests
			note := barkat.Note{
				Status:  "SET",
				Content: "Trade execution plan: Long at support with 2:1 RR",
				Format:  "MARKDOWN",
			}
			resp, err = client.R().SetBody(note).Post(barkat.JournalBase + "/" + createdJournal.ExternalID + "/notes")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))

			var envelope common.Envelope[barkat.Note]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			createdNote = envelope.Data
			Expect(createdNote.ExternalID).To(HavePrefix("not_"))
			Expect(createdNote.Status).To(Equal("SET"))
		})

		It("should have created note available", func() {
			// Verify the note created in BeforeEach is available
			Expect(createdNote.ExternalID).To(HavePrefix("not_"))
			Expect(createdNote.Status).To(Equal("SET"))
			Expect(createdNote.Content).To(Equal("Trade execution plan: Long at support with 2:1 RR"))
			Expect(createdNote.Format).To(Equal("MARKDOWN"))
		})

		It("should list all notes for journal", func() {
			resp, err := client.R().Get(barkat.JournalBase + "/" + createdJournal.ExternalID + "/notes")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var envelope common.Envelope[barkat.NoteList]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			notes := envelope.Data.Notes
			Expect(notes).ToNot(BeEmpty())

			// Verify our created note is in the list
			found := false
			for _, n := range notes {
				if n.ExternalID == createdNote.ExternalID {
					found = true
					Expect(n.Status).To(Equal("SET"))
					break
				}
			}
			Expect(found).To(BeTrue(), "Created note should be in the list")
		})

		It("should delete individual note", func() {
			// Delete the note created in BeforeEach
			resp, err := client.R().Delete(barkat.JournalBase + "/" + createdJournal.ExternalID + "/notes/" + createdNote.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusNoContent))

			// Verify note deletion by checking the notes list
			resp, err = client.R().Get(barkat.JournalBase + "/" + createdJournal.ExternalID + "/notes")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var listEnvelope common.Envelope[barkat.NoteList]
			Expect(json.Unmarshal(resp.Body(), &listEnvelope)).To(Succeed())
			notes := listEnvelope.Data.Notes

			// Check that deleted note is not in the list
			for _, n := range notes {
				Expect(n.ExternalID).ToNot(Equal(createdNote.ExternalID))
			}
		})
	})

	// Tag Management - Tests FR-005: Journal Tag Management
	Context("Tag Management", func() {
		var createdJournal barkat.Journal
		var createdTag barkat.Tag

		BeforeEach(func() {
			journal := barkat.Journal{
				Ticker:   "TAGTEST",
				Sequence: "MWD",
				Type:     "REJECTED",
				Status:   "FAIL",
				Images:   standardImages,
			}
			resp, err := client.R().SetBody(journal).Post(barkat.JournalBase)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
			createdJournal = decodeJournalResponse(resp)

			// Create a tag for use in tests
			tag := barkat.Tag{
				Tag:  "oe",
				Type: "REASON",
			}
			resp, err = client.R().SetBody(tag).Post(barkat.JournalBase + "/" + createdJournal.ExternalID + "/tags")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))

			var envelope common.Envelope[barkat.Tag]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			createdTag = envelope.Data
			Expect(createdTag.ExternalID).To(HavePrefix("tag_"))
			Expect(createdTag.Tag).To(Equal("oe"))
			Expect(createdTag.Type).To(Equal("REASON"))
		})

		It("should have created tag available", func() {
			// Verify the tag created in BeforeEach is available
			Expect(createdTag.ExternalID).To(HavePrefix("tag_"))
			Expect(createdTag.Tag).To(Equal("oe"))
			Expect(createdTag.Type).To(Equal("REASON"))
		})

		It("should list all tags for journal", func() {
			resp, err := client.R().Get(barkat.JournalBase + "/" + createdJournal.ExternalID + "/tags")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var envelope common.Envelope[barkat.TagList]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			tags := envelope.Data.Tags
			Expect(tags).ToNot(BeEmpty())

			// Verify our created tag is in the list
			found := false
			for _, t := range tags {
				if t.ExternalID == createdTag.ExternalID {
					found = true
					Expect(t.Tag).To(Equal("oe"))
					Expect(t.Type).To(Equal("REASON"))
					break
				}
			}
			Expect(found).To(BeTrue(), "Created tag should be in the list")
		})

		It("should delete individual tag", func() {
			// Delete the tag created in BeforeEach
			resp, err := client.R().Delete(barkat.JournalBase + "/" + createdJournal.ExternalID + "/tags/" + createdTag.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusNoContent))

			// Verify tag deletion by checking the tags list
			resp, err = client.R().Get(barkat.JournalBase + "/" + createdJournal.ExternalID + "/tags")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var listEnvelope common.Envelope[barkat.TagList]
			Expect(json.Unmarshal(resp.Body(), &listEnvelope)).To(Succeed())
			tags := listEnvelope.Data.Tags

			// Check that deleted tag is not in the list
			for _, t := range tags {
				Expect(t.ExternalID).ToNot(Equal(createdTag.ExternalID))
			}
		})
	})

	// Review Status Workflow - Tests FR-009: Journal Review
	Context("Review Status Workflow", func() {
		var createdJournal barkat.Journal

		BeforeEach(func() {
			journal := barkat.Journal{
				Ticker:   "REVIEW",
				Sequence: "YR",
				Type:     "RESULT",
				Status:   "SUCCESS",
				Images:   standardImages,
			}
			resp, err := client.R().SetBody(journal).Post(barkat.JournalBase)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
			createdJournal = decodeJournalResponse(resp)
		})

		It("should mark journal as reviewed", func() {
			reviewDate := civil.Date{Year: 2024, Month: 1, Day: 15}
			payload := barkat.JournalReviewUpdate{ReviewedAt: &reviewDate}

			resp, err := client.R().SetBody(payload).Patch(barkat.JournalBase + "/" + createdJournal.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var envelope common.Envelope[barkat.UpdateJournalStatusResponse]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			Expect(envelope.Data.ReviewedAt).ToNot(BeNil())
		})

		It("should clear reviewed status", func() {
			// First mark as reviewed
			reviewDate := civil.Date{Year: 2024, Month: 1, Day: 15}
			resp, _ := client.R().SetBody(barkat.JournalReviewUpdate{ReviewedAt: &reviewDate}).Patch(barkat.JournalBase + "/" + createdJournal.ExternalID)
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			// Then clear
			resp, err := client.R().SetBody(barkat.JournalReviewUpdate{ReviewedAt: nil}).
				Patch(barkat.JournalBase + "/" + createdJournal.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))

			var envelope common.Envelope[barkat.UpdateJournalStatusResponse]
			Expect(json.Unmarshal(resp.Body(), &envelope)).To(Succeed())
			Expect(envelope.Data.ReviewedAt).To(BeNil())
		})
	})

	// Validation Through HTTP Stack - Ensures validators are registered
	Context("Validation Errors", func() {
		It("should reject invalid ticker format", func() {
			journal := barkat.Journal{
				Ticker:   "lowercase", // PRD: must be uppercase
				Sequence: "MWD",
				Type:     "SET",
				Status:   "RUNNING",
				Images:   standardImages,
			}
			resp, err := client.R().SetBody(journal).Post(barkat.JournalBase)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusBadRequest))
		})

		It("should reject insufficient images", func() {
			journal := barkat.Journal{
				Ticker:   "VALID",
				Sequence: "MWD",
				Type:     "SET",
				Status:   "RUNNING",
				Images:   []barkat.Image{{Timeframe: "DL", FileName: "only_one.png"}}, // PRD: min 4
			}
			resp, err := client.R().SetBody(journal).Post(barkat.JournalBase)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusBadRequest))
		})

		It("should reject future review date", func() {
			// Create journal first
			journal := barkat.Journal{
				Ticker:   "FUTURE",
				Sequence: "MWD",
				Type:     "SET",
				Status:   "RUNNING",
				Images:   standardImages,
			}
			resp, _ := client.R().SetBody(journal).Post(barkat.JournalBase)
			created := decodeJournalResponse(resp)

			// Try to set future date
			futureDate := civil.Date{Year: 2099, Month: 12, Day: 31}
			resp, err := client.R().SetBody(barkat.JournalReviewUpdate{ReviewedAt: &futureDate}).
				Patch(barkat.JournalBase + "/" + created.ExternalID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusBadRequest))
		})
	})

	// Error Handling - Tests 404 and invalid ID scenarios
	Context("Error Handling", func() {
		It("should return 404 for non-existent journal", func() {
			// Use valid format (8 hex chars) but non-existent ID
			resp, err := client.R().Get(barkat.JournalBase + "/jrn_12345678")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusNotFound))
		})

		It("should return 400 for invalid ID format", func() {
			resp, err := client.R().Get(barkat.JournalBase + "/invalid_format")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusBadRequest))
		})
	})
})
