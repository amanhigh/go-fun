package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/manager/audit"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func decodeAuditCatalogResponse(w *httptest.ResponseRecorder) barkat.AuditCatalog {
	var envelope common.Envelope[barkat.AuditCatalog]
	util.AssertSuccess(w, http.StatusOK, &envelope)
	return envelope.Data
}

func decodeAuditResultResponse(w *httptest.ResponseRecorder) barkat.AuditResult {
	var envelope common.Envelope[barkat.AuditResult]
	util.AssertSuccess(w, http.StatusOK, &envelope)
	return envelope.Data
}

func newAuditTestRouter(auditHandler handler.AuditHandler) *gin.Engine {
	router := util.CreateTestGinRouter()
	audits := router.Group(barkat.AuditBase)
	handler.SetupAuditRoutes(audits, auditHandler)
	return router
}

func newAuditTestHandler(db *gorm.DB) handler.AuditHandler {
	auditRepo := repository.NewAuditRepository(db)
	registry := audit.NewPluginRegistry()
	Expect(registry.RegisterPlugin(audit.NewAlertCoveragePlugin(auditRepo))).ToNot(HaveOccurred())
	auditMgr := manager.NewAuditManager(registry)
	return handler.NewAuditHandler(auditMgr)
}

func seedAuditTicker(db *gorm.DB, ticker, state string) barkat.Ticker {
	result := barkat.Ticker{
		Ticker:       ticker,
		Exchange:     new("NSE"),
		Timeframes:   []string{"MN", "WK", "DL"},
		Type:         "EQUITY",
		State:        state,
		Trend:        "UPTREND",
		LastOpenedAt: time.Date(2026, time.May, 5, 10, 30, 0, 0, time.UTC),
	}
	Expect(db.Create(&result).Error).ToNot(HaveOccurred())
	return result
}

func seedAuditAlertTicker(db *gorm.DB, ticker barkat.Ticker, symbol, pairID string) barkat.AlertTicker {
	result := barkat.AlertTicker{
		TickerID: ticker.ID,
		Symbol:   symbol,
		PairID:   pairID,
		Name:     symbol + " Ltd",
		Exchange: new("NSE"),
	}
	Expect(db.Create(&result).Error).ToNot(HaveOccurred())
	return result
}

func seedAuditPriceAlert(db *gorm.DB, alertTicker barkat.AlertTicker, alertID string) {
	alert := barkat.PriceAlert{
		AlertTickerID: alertTicker.ID,
		AlertID:       &alertID,
		TriggerPrice:  100.25,
	}
	Expect(db.Create(&alert).Error).ToNot(HaveOccurred())
}

// AuditHandler Integration Tests.
// Tests complete HTTP → Handler → Manager → Repository → Database flow.
// Covers PRD Section 2.2 Audit APIs and FR-002 Alert Coverage Audit Plugin.
var _ = Describe("AuditHandler Integration - Section 2.2 Audit APIs", func() {
	var (
		auditHandler handler.AuditHandler
		router       *gin.Engine
		db           *gorm.DB
	)

	BeforeEach(func() {
		var err error
		core.RegisterJournalValidators()
		db, err = core.CreateTestBarkatDB()
		Expect(err).ToNot(HaveOccurred())

		auditHandler = newAuditTestHandler(db)
		router = newAuditTestRouter(auditHandler)
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Describe("GET /v1/api/audits - List Audit Catalog (2.2.1)", func() {
		Context("Happy Path", func() {
			var response barkat.AuditCatalog
			var w *httptest.ResponseRecorder

			BeforeEach(func() {
				req, recorder := util.CreateTestRequest(http.MethodGet, barkat.AuditBase, nil)
				w = recorder
				router.ServeHTTP(w, req)
				response = decodeAuditCatalogResponse(w)
			})

			It("should return 200 OK", func() { Expect(w.Code).To(Equal(http.StatusOK)) })
			It("should return Envelope success", func() {
				var envelope common.Envelope[barkat.AuditCatalog]
				util.AssertSuccess(w, http.StatusOK, &envelope)
				Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			})
			It("should return the implemented alert coverage audit", func() {
				Expect(response.Audits).To(HaveLen(1))
				Expect(response.Audits[0].ID).To(Equal("alert-coverage"))
				Expect(response.Audits[0].Title).To(Equal("Alert Coverage"))
				Expect(response.Audits[0].Description).ToNot(BeEmpty())
				Expect(response.Audits[0].Order).To(Equal(1))
			})
		})
	})

	Describe("GET /v1/api/audits/{audit-id}/results - Execute Single Audit (2.2.2)", func() {
		Context("Alert Coverage Plugin", func() {
			var (
				response barkat.AuditResult
				w        *httptest.ResponseRecorder
			)

			BeforeEach(func() {
				missingMappingTicker := seedAuditTicker(db, "MCX", "READY")
				Expect(missingMappingTicker.ID).ToNot(BeZero())

				noAlertsTicker := seedAuditTicker(db, "INFY", "READY")
				seedAuditAlertTicker(db, noAlertsTicker, "INFY", "1001")

				singleAlertTicker := seedAuditTicker(db, "TCS", "READY")
				singleAlertOwner := seedAuditAlertTicker(db, singleAlertTicker, "TCS", "1002")
				seedAuditPriceAlert(db, singleAlertOwner, "2001")

				validTicker := seedAuditTicker(db, "RELIANCE", "READY")
				validAlertOwner := seedAuditAlertTicker(db, validTicker, "RELIANCE", "1003")
				seedAuditPriceAlert(db, validAlertOwner, "2002")
				seedAuditPriceAlert(db, validAlertOwner, "2003")

				watchedTicker := seedAuditTicker(db, "NIFTY", "WATCHED")
				Expect(watchedTicker.ID).ToNot(BeZero())
				blacklistedTicker := seedAuditTicker(db, "BANNED", "BLACKLIST")
				Expect(blacklistedTicker.ID).ToNot(BeZero())

				req, recorder := util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/alert-coverage/results", nil)
				w = recorder
				router.ServeHTTP(w, req)
				response = decodeAuditResultResponse(w)
			})

			It("should return 200 OK", func() { Expect(w.Code).To(Equal(http.StatusOK)) })
			It("should echo the kebab-case audit ID", func() { Expect(response.AuditID).To(Equal("alert-coverage")) })
			It("should set generated_at", func() { Expect(response.GeneratedAt).ToNot(BeZero()) })
			It("should include full-result counts by finding code", func() {
				Expect(response.Counts).To(Equal(map[string]int{
					"NO_ALERT_TICKER": 1,
					"NO_ALERTS":       1,
					"SINGLE_ALERT":    1,
				}))
			})
			It("should return one finding for each coverage gap", func() { Expect(response.Findings).To(HaveLen(3)) })
			It("should mark total as the full unpaginated finding count", func() {
				Expect(response.Metadata.Total).To(Equal(int64(3)))
				Expect(response.Metadata.Offset).To(Equal(0))
				Expect(response.Metadata.Limit).To(Equal(20))
			})
			It("should distinguish missing mapping, no-alert, and single-alert gaps", func() {
				Expect(response.Findings).To(ContainElements(
					barkat.AuditFinding{Code: "NO_ALERT_TICKER", Target: "MCX", Severity: "HIGH", Data: map[string]string{"alert_ticker_count": "0", "price_alert_count": "0"}},
					barkat.AuditFinding{Code: "NO_ALERTS", Target: "INFY", Severity: "MEDIUM", Data: map[string]string{"alert_ticker_count": "1", "price_alert_count": "0"}},
					barkat.AuditFinding{Code: "SINGLE_ALERT", Target: "TCS", Severity: "HIGH", Data: map[string]string{"alert_ticker_count": "1", "price_alert_count": "1"}},
				))
			})
			It("should skip actively watched and blacklisted instruments", func() {
				targets := make([]string, 0, len(response.Findings))
				for _, finding := range response.Findings {
					targets = append(targets, finding.Target)
				}
				Expect(targets).ToNot(ContainElement("NIFTY"))
				Expect(targets).ToNot(ContainElement("BANNED"))
			})
		})

		Context("Pagination", func() {
			It("should paginate findings while keeping counts and total scoped to the full result set", func() {
				for i, ticker := range []string{"AAA", "BBB", "CCC"} {
					seededTicker := seedAuditTicker(db, ticker, "READY")
					if i > 0 {
						seedAuditAlertTicker(db, seededTicker, ticker, "30"+ticker)
					}
				}

				req, w := util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/alert-coverage/results?offset=1&limit=1", nil)
				router.ServeHTTP(w, req)
				response := decodeAuditResultResponse(w)

				Expect(response.Findings).To(HaveLen(1))
				Expect(response.Metadata.Total).To(Equal(int64(3)))
				Expect(response.Metadata.Offset).To(Equal(1))
				Expect(response.Metadata.Limit).To(Equal(1))
				Expect(response.Counts).To(Equal(map[string]int{"NO_ALERT_TICKER": 1, "NO_ALERTS": 2}))
			})
		})

		Context("Clean Audit", func() {
			It("should use metadata.total == 0 to signal pass", func() {
				validTicker := seedAuditTicker(db, "RELIANCE", "READY")
				validOwner := seedAuditAlertTicker(db, validTicker, "RELIANCE", "1003")
				seedAuditPriceAlert(db, validOwner, "2002")
				seedAuditPriceAlert(db, validOwner, "2003")

				req, w := util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/alert-coverage/results", nil)
				router.ServeHTTP(w, req)
				response := decodeAuditResultResponse(w)

				Expect(response.Findings).To(BeEmpty())
				Expect(response.Counts).To(BeEmpty())
				Expect(response.Metadata.Total).To(Equal(int64(0)))
			})
		})

		Context("Field Validations", func() {
			It("should return 400 for invalid pagination query", func() {
				req, w := util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/alert-coverage/results?limit=0", nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				util.AssertError(w, "Limit", "min")
			})

			It("should return 404 for unknown audit id", func() {
				req, w := util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/unknown-audit/results", nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
				var response map[string]any
				Expect(json.Unmarshal(w.Body.Bytes(), &response)).ToNot(HaveOccurred())
				Expect(response["status"]).To(Equal("fail"))
				Expect(response["data"]).To(HaveKeyWithValue("audit-id", "Audit not found"))
			})

			It("should return 404 for snake_case audit id", func() {
				req, w := util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/alert_coverage/results", nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

})
