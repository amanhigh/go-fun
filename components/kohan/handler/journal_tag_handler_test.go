//nolint:dupl
package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func decodeTagResponse(w *httptest.ResponseRecorder) barkat.Tag {
	var envelope common.Envelope[barkat.Tag]
	util.AssertSuccess(w, http.StatusCreated, &envelope)
	return envelope.Data
}

func decodeTagListResponse(w *httptest.ResponseRecorder) []barkat.Tag {
	var envelope common.Envelope[map[string][]barkat.Tag]
	util.AssertSuccess(w, http.StatusOK, &envelope)
	return envelope.Data["tags"]
}

// TagHandler Integration Tests - Comprehensive Master Specification
// Tests complete HTTP → Handler → Manager → Repository → Database flow
// Covers all PRD validations for Section 2.4 JournalTag APIs
//
// TEST STRUCTURE FORMAT:
// ====================
// Describe(API)
// -> Context(Happy Path): 2xx Success Cases
// -> Context(Field Validations): All 4xx Validation Cases
//    -> Context(Field Name): One Context for Each Field
//       -> Context(Allowed Values): All Variations of Valid Values (2xx) - If Applicable
//       -> Context(Bad Values): All Variations of Missing,Regex,Min,Max Edge Cases (4xx)
// -> Context(Errors): 5xx Server Error Cases

var _ = Describe("TagHandler Integration - Section 2.4 JournalTag APIs", func() {
	var (
		tagHandler *handler.TagHandlerImpl
		router     *gin.Engine
		testCtx    = context.Background()
		db         *gorm.DB
		journalMgr manager.JournalManager
		tagMgr     manager.TagManager
		journal    barkat.Journal
		req        *http.Request
		w          *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		var err error
		db, err = core.CreateTestBarkatDB()
		Expect(err).ToNot(HaveOccurred())

		journalRepo := repository.NewJournalRepository(db)
		journalMgr = manager.NewJournalManager(journalRepo)
		tagMgr = manager.NewTagManager(journalMgr, repository.NewTagRepository(db))
		tagHandler = handler.NewTagHandler(tagMgr)

		router = util.CreateTestGinRouter()
		v1 := router.Group("/v1")
		journalGroup := v1.Group("/journals")
		handler.SetupTagRoutes(journalGroup, tagHandler)

		// Create base journal for tag operations
		journal = barkat.Journal{
			Ticker:   "GRSE",
			Sequence: "MWD",
			Type:     "REJECTED",
			Status:   "FAIL",
			Images: []barkat.Image{
				{Timeframe: "DL", FileName: "test-dl.png"},
				{Timeframe: "WK", FileName: "test-wk.png"},
				{Timeframe: "MN", FileName: "test-mn.png"},
				{Timeframe: "TMN", FileName: "test-tmn.png"},
			},
		}
		Expect(journalMgr.CreateJournal(testCtx, &journal)).To(Succeed())
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	// ============================================================================
	// 2.4.1 POST /v1/journals/{journal-id}/tags - Add Tag
	// ============================================================================
	Describe("POST /v1/journals/{journal-id}/tags - Add Tag (2.4.1)", func() {
		Context("Happy Path", func() {
			Context("with valid tag data", func() {
				var response barkat.Tag

				BeforeEach(func() {
					tag := barkat.Tag{
						Tag:  "dep",
						Type: "REASON",
					}
					req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
					router.ServeHTTP(w, req)
				})

				It("should return 201 Created", func() {
					Expect(w.Code).To(Equal(http.StatusCreated))
				})

				It("should return Envelope success", func() {
					var envelope common.Envelope[barkat.Tag]
					util.AssertSuccess(w, http.StatusCreated, &envelope)
					Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
				})

				It("should return created tag with external ID", func() {
					response = decodeTagResponse(w)
					Expect(response.ExternalID).To(HavePrefix("tag_"))
				})

				It("should preserve tag field", func() {
					response = decodeTagResponse(w)
					Expect(response.Tag).To(Equal("dep"))
				})

				It("should preserve type field", func() {
					response = decodeTagResponse(w)
					Expect(response.Type).To(Equal("REASON"))
				})

				It("should set created_at timestamp", func() {
					response = decodeTagResponse(w)
					Expect(response.CreatedAt).ToNot(BeZero())
				})

				It("should persist tag to database", func() {
					tags, err := tagMgr.ListTags(testCtx, journal.ExternalID, "")
					Expect(err).ToNot(HaveOccurred())
					Expect(tags).To(HaveLen(1))
				})
			})

			// HACK: Remove as it is Covered in Field Validation.
			Context("with optional override field", func() {
				It("should accept tag with override", func() {
					override := "loc"
					tag := barkat.Tag{
						Tag:      "dep",
						Type:     "REASON",
						Override: &override,
					}
					req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
					router.ServeHTTP(w, req)
					response := decodeTagResponse(w)
					Expect(response.Override).ToNot(BeNil())
					Expect(*response.Override).To(Equal("loc"))
				})

				It("should accept tag without override (nil)", func() {
					tag := barkat.Tag{
						Tag:  "oe",
						Type: "REASON",
					}
					req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
					router.ServeHTTP(w, req)
					response := decodeTagResponse(w)
					Expect(response.Override).To(BeNil())
				})
			})
		})

		Context("Field Validations", func() {
			Context("Tag Field", func() {
				Context("Allowed Values", func() {
					It("should accept minimum tag length (1 char)", func() {
						tag := barkat.Tag{Tag: "a", Type: "REASON"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						response := decodeTagResponse(w)
						Expect(response.Tag).To(Equal("a"))
					})

					It("should accept maximum tag length (10 chars)", func() {
						tag := barkat.Tag{Tag: "abcdefghij", Type: "REASON"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						response := decodeTagResponse(w)
						Expect(response.Tag).To(Equal("abcdefghij"))
					})

					It("should accept tag with alphanumeric characters", func() {
						tag := barkat.Tag{Tag: "dep123", Type: "REASON"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						response := decodeTagResponse(w)
						Expect(response.Tag).To(Equal("dep123"))
					})

					It("should accept tag with hyphens", func() {
						tag := barkat.Tag{Tag: "dep-loc", Type: "REASON"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						response := decodeTagResponse(w)
						Expect(response.Tag).To(Equal("dep-loc"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing tag (PRD: required)", func() {
						tag := barkat.Tag{Tag: "", Type: "REASON"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Tag", "required")
					})

					It("should return 400 for tag exceeding max length (PRD: max 10 chars)", func() {
						tag := barkat.Tag{Tag: "abcdefghijk", Type: "REASON"} // 11 chars
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Tag", "max")
					})

					It("should return 400 for tag with invalid characters (PRD: alphanumeric with hyphens)", func() {
						tag := barkat.Tag{Tag: "dep@loc", Type: "REASON"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Tag", "tag")
					})

					It("should return 400 for tag with spaces", func() {
						tag := barkat.Tag{Tag: "dep loc", Type: "REASON"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Tag", "tag")
					})
				})
			})

			Context("Type Field", func() {
				Context("Allowed Values", func() {
					It("should accept type = REASON", func() {
						tag := barkat.Tag{Tag: "dep", Type: "REASON"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						response := decodeTagResponse(w)
						Expect(response.Type).To(Equal("REASON"))
					})

					It("should accept type = MANAGEMENT", func() {
						tag := barkat.Tag{Tag: "ptr", Type: "MANAGEMENT"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						response := decodeTagResponse(w)
						Expect(response.Type).To(Equal("MANAGEMENT"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing type (PRD: required)", func() {
						tag := barkat.Tag{Tag: "dep", Type: ""}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "required")
					})

					It("should return 400 for invalid type enum (PRD: must be REASON or MANAGEMENT)", func() {
						tag := barkat.Tag{Tag: "dep", Type: "INVALID"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})

					It("should return 400 for lowercase type (PRD: case-sensitive)", func() {
						tag := barkat.Tag{Tag: "dep", Type: "reason"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})
				})
			})

			Context("Override Field", func() {
				Context("Allowed Values", func() {
					It("should accept override as nil (optional field)", func() {
						tag := barkat.Tag{
							Tag:  "dep",
							Type: "REASON",
							// Override is nil by default
						}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						response := decodeTagResponse(w)
						Expect(response.Override).To(BeNil())
					})

					It("should accept minimum override length (1 char)", func() {
						override := "a"
						tag := barkat.Tag{Tag: "dep", Type: "REASON", Override: &override}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						response := decodeTagResponse(w)
						Expect(*response.Override).To(Equal("a"))
					})

					It("should accept maximum override length (5 chars)", func() {
						override := "abcde"
						tag := barkat.Tag{Tag: "dep", Type: "REASON", Override: &override}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						response := decodeTagResponse(w)
						Expect(*response.Override).To(Equal("abcde"))
					})

					It("should accept override with letters only", func() {
						override := "loc"
						tag := barkat.Tag{Tag: "dep", Type: "REASON", Override: &override}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						response := decodeTagResponse(w)
						Expect(*response.Override).To(Equal("loc"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for override exceeding max length (PRD: max 5 chars)", func() {
						override := "abcdef" // 6 chars
						tag := barkat.Tag{Tag: "dep", Type: "REASON", Override: &override}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Override", "max")
					})

					It("should return 400 for override with numbers (PRD: letters only)", func() {
						override := "loc1"
						tag := barkat.Tag{Tag: "dep", Type: "REASON", Override: &override}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Override", "override")
					})

					It("should return 400 for override with special characters", func() {
						override := "loc@"
						tag := barkat.Tag{Tag: "dep", Type: "REASON", Override: &override}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Override", "override")
					})
				})
			})

			Context("Journal ID Path Parameter", func() {
				Context("Bad Values", func() {
					It("should return 404 for non-existent journal ID", func() {
						tag := barkat.Tag{Tag: "dep", Type: "REASON"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/nonexistent-id/tags", tag)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for malformed journal ID", func() {
						tag := barkat.Tag{Tag: "dep", Type: "REASON"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/invalid-uuid-format/tags", tag)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for valid UUID format but non-existent", func() {
						tag := barkat.Tag{Tag: "dep", Type: "REASON"}
						req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/550e8400-e29b-41d4-a716-446655440000/tags", tag)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 400 for invalid JSON", func() {
				req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", []byte("invalid json"))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for empty request body", func() {
				req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", []byte(""))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for null request body", func() {
				req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", []byte("null"))
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 409 for duplicate tag (PRD: unique tag per journal)", func() {
				// Create first tag
				tag1 := barkat.Tag{Tag: "dep", Type: "REASON"}
				_, err := tagMgr.CreateTag(testCtx, journal.ExternalID, tag1)
				Expect(err).ToNot(HaveOccurred())

				// Attempt to create duplicate
				tag2 := barkat.Tag{Tag: "dep", Type: "REASON"}
				req, w = util.CreateTestRequest("POST", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", tag2)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusConflict))
			})
		})
	})

	// ============================================================================
	// 2.4.2 GET /v1/journals/{journal-id}/tags - List Tags
	// ============================================================================
	Describe("GET /v1/journals/{journal-id}/tags - List Tags (2.4.2)", func() {
		Context("Happy Path", func() {
			Context("with journal having tags", func() {
				var tags []barkat.Tag

				BeforeEach(func() {
					// Create multiple tags for testing
					tag1 := barkat.Tag{Tag: "dep", Type: "REASON"}
					_, err := tagMgr.CreateTag(testCtx, journal.ExternalID, tag1)
					Expect(err).ToNot(HaveOccurred())

					tag2 := barkat.Tag{Tag: "ptr", Type: "MANAGEMENT"}
					_, err = tagMgr.CreateTag(testCtx, journal.ExternalID, tag2)
					Expect(err).ToNot(HaveOccurred())

					req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", nil)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
				})

				It("should return all tags for journal", func() {
					tags = decodeTagListResponse(w)
					Expect(tags).To(HaveLen(2))
				})

				It("should return tags with correct types", func() {
					tags = decodeTagListResponse(w)
					types := []string{}
					for _, tag := range tags {
						types = append(types, tag.Type)
					}
					Expect(types).To(ContainElements("REASON", "MANAGEMENT"))
				})

				It("should return tags with external IDs", func() {
					tags = decodeTagListResponse(w)
					for _, tag := range tags {
						Expect(tag.ExternalID).To(HavePrefix("tag_"))
					}
				})

				It("should return tags with created_at timestamps", func() {
					tags = decodeTagListResponse(w)
					for _, tag := range tags {
						Expect(tag.CreatedAt).ToNot(BeZero())
					}
				})
			})

			Context("with journal having no tags", func() {
				BeforeEach(func() {
					req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/"+journal.ExternalID+"/tags", nil)
					router.ServeHTTP(w, req)
				})

				It("should return 200 OK with empty array", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
					tags := decodeTagListResponse(w)
					Expect(tags).To(BeEmpty())
				})
			})
		})

		Context("Field Validations", func() {
			Context("Type Query Parameter", func() {
				BeforeEach(func() {
					// Create tags for filtering tests
					tag1 := barkat.Tag{Tag: "dep", Type: "REASON"}
					tag2 := barkat.Tag{Tag: "oe", Type: "MANAGEMENT"}
					_, err1 := tagMgr.CreateTag(testCtx, journal.ExternalID, tag1)
					_, err2 := tagMgr.CreateTag(testCtx, journal.ExternalID, tag2)
					Expect(err1).ToNot(HaveOccurred())
					Expect(err2).ToNot(HaveOccurred())
				})

				Context("Allowed Values", func() {
					It("should filter tags by type = REASON", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/"+journal.ExternalID+"/tags?type=REASON", nil)
						router.ServeHTTP(w, req)
						tags := decodeTagListResponse(w)
						Expect(tags).To(HaveLen(1))
						Expect(tags[0].Type).To(Equal("REASON"))
						Expect(tags[0].Tag).To(Equal("dep"))
					})

					It("should filter tags by type = MANAGEMENT", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/"+journal.ExternalID+"/tags?type=MANAGEMENT", nil)
						router.ServeHTTP(w, req)
						tags := decodeTagListResponse(w)
						Expect(tags).To(HaveLen(1))
						Expect(tags[0].Type).To(Equal("MANAGEMENT"))
						Expect(tags[0].Tag).To(Equal("oe"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid type query parameter", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/"+journal.ExternalID+"/tags?type=INVALID", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})
			Context("Journal ID Path Parameter", func() {
				Context("Bad Values", func() {
					It("should return 404 for non-existent journal ID", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/nonexistent-id/tags", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for malformed journal ID", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/invalid-uuid-format/tags", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for valid UUID format but non-existent", func() {
						req, w = util.CreateTestRequest("GET", barkat.JournalEntries+"/550e8400-e29b-41d4-a716-446655440000/tags", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			// No server error scenarios for GET currently
		})
	})

	// ============================================================================
	// 2.4.3 DELETE /v1/journals/{journal-id}/tags/{tag-id} - Remove Tag
	// ============================================================================
	Describe("DELETE /v1/journals/{journal-id}/tags/{tag-id} - Remove Tag (2.4.3)", func() {
		var tagToDelete *barkat.Tag

		BeforeEach(func() {
			// Create a tag to delete
			tag := barkat.Tag{Tag: "dep", Type: "REASON"}
			var err error
			tagToDelete, err = tagMgr.CreateTag(testCtx, journal.ExternalID, tag)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Happy Path", func() {
			Context("with valid journal and tag IDs", func() {
				BeforeEach(func() {
					req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+journal.ExternalID+"/tags/"+tagToDelete.ExternalID, nil)
					router.ServeHTTP(w, req)
				})

				It("should return 204 No Content", func() {
					Expect(w.Code).To(Equal(http.StatusNoContent))
				})

				It("should return empty body", func() {
					Expect(w.Body.String()).To(BeEmpty())
				})

				It("should actually delete the tag from database", func() {
					tags, err := tagMgr.ListTags(testCtx, journal.ExternalID, "")
					Expect(err).ToNot(HaveOccurred())
					Expect(tags).To(BeEmpty())
				})
			})
		})

		Context("Field Validations", func() {
			Context("Journal ID Path Parameter", func() {
				Context("Bad Values", func() {
					It("should return 404 for non-existent journal ID", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/nonexistent-id/tags/"+tagToDelete.ExternalID, nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for malformed journal ID", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/invalid-uuid-format/tags/"+tagToDelete.ExternalID, nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})

			Context("Tag ID Path Parameter", func() {
				Context("Bad Values", func() {
					It("should return 404 for non-existent tag ID", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+journal.ExternalID+"/tags/nonexistent-tag", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for malformed tag ID", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+journal.ExternalID+"/tags/invalid-uuid-format", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 for valid UUID format but non-existent", func() {
						req, w = util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+journal.ExternalID+"/tags/550e8400-e29b-41d4-a716-446655440000", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 404 on second delete (idempotency check)", func() {
				// First delete
				req1, w1 := util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+journal.ExternalID+"/tags/"+tagToDelete.ExternalID, nil)
				router.ServeHTTP(w1, req1)
				Expect(w1.Code).To(Equal(http.StatusNoContent))

				// Second delete should return 404
				req2, w2 := util.CreateTestRequest("DELETE", barkat.JournalEntries+"/"+journal.ExternalID+"/tags/"+tagToDelete.ExternalID, nil)
				router.ServeHTTP(w2, req2)
				Expect(w2.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
