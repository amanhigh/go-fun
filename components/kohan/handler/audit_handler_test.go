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
	Expect(registry.RegisterPlugin(audit.NewStaleReviewPlugin(auditRepo))).ToNot(HaveOccurred())
	auditMgr := manager.NewAuditManager(registry)
	return handler.NewAuditHandler(auditMgr)
}

func seedAuditTicker(db *gorm.DB, ticker, state string) barkat.Ticker {
	result := barkat.Ticker{
		Ticker:       ticker,
		Exchange:     "NSE",
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
		req          *http.Request
		w            *httptest.ResponseRecorder
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

			BeforeEach(func() {
				req, w = util.CreateTestRequest(http.MethodGet, barkat.AuditBase, nil)
				router.ServeHTTP(w, req)
				response = decodeAuditCatalogResponse(w)
			})

			It("should return 200 OK", func() { Expect(w.Code).To(Equal(http.StatusOK)) })
			It("should return Envelope success", func() {
				var envelope common.Envelope[barkat.AuditCatalog]
				util.AssertSuccess(w, http.StatusOK, &envelope)
				Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			})
			It("should return the implemented audits in order", func() {
				Expect(response.Audits).To(HaveLen(2))
				Expect(response.Audits[0].ID).To(Equal("alert-coverage"))
				Expect(response.Audits[0].Title).To(Equal("Alert Coverage"))
				Expect(response.Audits[0].Description).ToNot(BeEmpty())
				Expect(response.Audits[0].Order).To(Equal(1))
				Expect(response.Audits[1].ID).To(Equal("stale-review"))
				Expect(response.Audits[1].Title).To(Equal("Stale Review"))
				Expect(response.Audits[1].Description).ToNot(BeEmpty())
				Expect(response.Audits[1].Order).To(Equal(2))
			})
		})
	})

	Describe("GET /v1/api/audits/{audit-id}/results - Execute Single Audit (2.2.2)", func() {
		Context("Happy Path", func() {
			Context("When tickers have mixed alert coverage", func() {
				var response barkat.AuditResult

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

					req, w = util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/alert-coverage/results", nil)
					router.ServeHTTP(w, req)
					response = decodeAuditResultResponse(w)
				})

				It("should return 200 OK", func() { Expect(w.Code).To(Equal(http.StatusOK)) })
				It("should echo the kebab-case audit ID", func() { Expect(response.AuditID).To(Equal("alert-coverage")) })
				It("should set generated_at", func() { Expect(response.GeneratedAt).ToNot(BeZero()) })
				It("should include full-result counts by finding code", func() {
					Expect(response.Counts).To(Equal(map[string]int{
						"NO_ALERT_TICKER": 2,
						"NO_ALERTS":       1,
						"SINGLE_ALERT":    1,
					}))
				})
				It("should return one finding for each coverage gap", func() { Expect(response.Findings).To(HaveLen(4)) })
				It("should mark total as the full unpaginated finding count", func() {
					Expect(response.Metadata.Total).To(Equal(int64(4)))
					Expect(response.Metadata.Offset).To(Equal(0))
					Expect(response.Metadata.Limit).To(Equal(20))
				})
				It("should distinguish missing mapping, no-alert, and single-alert gaps", func() {
					Expect(response.Findings).To(ContainElements(
						barkat.AuditFinding{Code: "NO_ALERT_TICKER", Target: "MCX", Severity: "HIGH", Data: map[string]string{"alert_ticker_count": "0", "price_alert_count": "0"}},
						barkat.AuditFinding{Code: "NO_ALERT_TICKER", Target: "NIFTY", Severity: "HIGH", Data: map[string]string{"alert_ticker_count": "0", "price_alert_count": "0"}},
						barkat.AuditFinding{Code: "NO_ALERTS", Target: "INFY", Severity: "MEDIUM", Data: map[string]string{"alert_ticker_count": "1", "price_alert_count": "0"}},
						barkat.AuditFinding{Code: "SINGLE_ALERT", Target: "TCS", Severity: "HIGH", Data: map[string]string{"alert_ticker_count": "1", "price_alert_count": "1"}},
					))
				})
				It("should include watched tickers and skip only blacklisted instruments", func() {
					targets := make([]string, 0, len(response.Findings))
					for _, finding := range response.Findings {
						targets = append(targets, finding.Target)
					}
					Expect(targets).To(ContainElement("NIFTY"))
					Expect(targets).ToNot(ContainElement("BANNED"))
				})
			})

			Context("When all READY tickers have sufficient alert coverage", func() {
				var response barkat.AuditResult

				BeforeEach(func() {
					validTicker := seedAuditTicker(db, "RELIANCE", "READY")
					validOwner := seedAuditAlertTicker(db, validTicker, "RELIANCE", "1003")
					seedAuditPriceAlert(db, validOwner, "2002")
					seedAuditPriceAlert(db, validOwner, "2003")

					req, w = util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/alert-coverage/results", nil)
					router.ServeHTTP(w, req)
					response = decodeAuditResultResponse(w)
				})

				It("should use metadata.total == 0 to signal pass", func() {
					Expect(response.Findings).To(BeEmpty())
					Expect(response.Counts).To(BeEmpty())
					Expect(response.Metadata.Total).To(Equal(int64(0)))
				})
			})
		})

		Context("When tracked tickers have stale review dates", func() {
			var response barkat.AuditResult

			BeforeEach(func() {
				// Ticker with old last_opened_at (stale - more than 180 days ago)
				staleTicker := barkat.Ticker{
					Ticker:       "STALE1",
					Exchange:     "NSE",
					Timeframes:   []string{"MN", "WK", "DL"},
					Type:         "EQUITY",
					State:        "WATCHED",
					Trend:        "UPTREND",
					LastOpenedAt: time.Date(2025, time.November, 1, 10, 30, 0, 0, time.UTC),
				}
				Expect(db.Create(&staleTicker).Error).ToNot(HaveOccurred())

				// Ticker opened recently (not stale)
				recentTicker := barkat.Ticker{
					Ticker:       "RECENT1",
					Exchange:     "NSE",
					Timeframes:   []string{"MN", "WK", "DL"},
					Type:         "EQUITY",
					State:        "WATCHED",
					Trend:        "UPTREND",
					LastOpenedAt: time.Now().UTC(),
				}
				Expect(db.Create(&recentTicker).Error).ToNot(HaveOccurred())

				// BLACKLIST ticker with old last_opened_at (stale, but excluded by BLACKLIST filter)
				blackTicker := barkat.Ticker{
					Ticker:       "BLACK1",
					Exchange:     "NSE",
					Timeframes:   []string{"MN", "WK", "DL"},
					Type:         "EQUITY",
					State:        "BLACKLIST",
					Trend:        "UPTREND",
					LastOpenedAt: time.Date(2025, time.November, 1, 10, 30, 0, 0, time.UTC),
				}
				Expect(db.Create(&blackTicker).Error).ToNot(HaveOccurred())

				req, w = util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/stale-review/results", nil)
				router.ServeHTTP(w, req)
				response = decodeAuditResultResponse(w)
			})

			It("should return 200 OK", func() { Expect(w.Code).To(Equal(http.StatusOK)) })
			It("should echo the kebab-case audit ID", func() { Expect(response.AuditID).To(Equal("stale-review")) })
			It("should set generated_at", func() { Expect(response.GeneratedAt).ToNot(BeZero()) })
			It("should find stale tickers older than the threshold", func() {
				targets := make([]string, 0, len(response.Findings))
				for _, finding := range response.Findings {
					targets = append(targets, finding.Target)
				}
				Expect(targets).To(ContainElement("STALE1"))
				Expect(targets).ToNot(ContainElement("BLACK1"))
			})
			It("should exclude recently opened tickers", func() {
				targets := make([]string, 0, len(response.Findings))
				for _, finding := range response.Findings {
					targets = append(targets, finding.Target)
				}
				Expect(targets).ToNot(ContainElement("RECENT1"))
			})
			It("should include last_opened_at in finding data", func() {
				for _, finding := range response.Findings {
					if finding.Target == "STALE1" {
						Expect(finding.Data).To(HaveKeyWithValue("last_opened_at", Not(BeEmpty())))
						break
					}
				}
			})
			It("should exclude blacklisted tickers from stale review", func() {
				targets := make([]string, 0, len(response.Findings))
				for _, finding := range response.Findings {
					targets = append(targets, finding.Target)
				}
				Expect(targets).ToNot(ContainElement("BLACK1"))
			})
			It("should include STALE_TICKER code in counts", func() {
				Expect(response.Counts).To(HaveKey("STALE_TICKER"))
				Expect(response.Counts["STALE_TICKER"]).To(Equal(1))
			})
		})

		Context("When no tickers are stale", func() {
			var response barkat.AuditResult

			BeforeEach(func() {
				recentTicker := barkat.Ticker{
					Ticker:       "FRESH1",
					Exchange:     "NSE",
					Timeframes:   []string{"MN", "WK", "DL"},
					Type:         "EQUITY",
					State:        "WATCHED",
					Trend:        "UPTREND",
					LastOpenedAt: time.Now().UTC(),
				}
				Expect(db.Create(&recentTicker).Error).ToNot(HaveOccurred())

				req, w = util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/stale-review/results", nil)
				router.ServeHTTP(w, req)
				response = decodeAuditResultResponse(w)
			})

			It("should use metadata.total == 0 to signal pass", func() {
				Expect(response.Findings).To(BeEmpty())
				Expect(response.Counts).To(BeEmpty())
				Expect(response.Metadata.Total).To(Equal(int64(0)))
			})
		})

		Context("Field Validations", func() {
			Context("Pagination Field", func() {
				Context("Allowed Values", func() {
					It("should paginate findings while keeping counts and total scoped to the full result set", func() {
						for i, ticker := range []string{"AAA", "BBB", "CCC"} {
							seededTicker := seedAuditTicker(db, ticker, "READY")
							if i > 0 {
								seedAuditAlertTicker(db, seededTicker, ticker, "30"+ticker)
							}
						}

						req, w = util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/alert-coverage/results?offset=1&limit=1", nil)
						router.ServeHTTP(w, req)
						response := decodeAuditResultResponse(w)

						Expect(response.Findings).To(HaveLen(1))
						Expect(response.Metadata.Total).To(Equal(int64(3)))
						Expect(response.Metadata.Offset).To(Equal(1))
						Expect(response.Metadata.Limit).To(Equal(1))
						Expect(response.Counts).To(Equal(map[string]int{"NO_ALERT_TICKER": 1, "NO_ALERTS": 2}))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid pagination query", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/alert-coverage/results?limit=0", nil)
						router.ServeHTTP(w, req)

						Expect(w.Code).To(Equal(http.StatusBadRequest))
						util.AssertError(w, "Limit", "min")
					})
				})
			})

			Context("Audit ID Field", func() {
				Context("Bad Values", func() {
					It("should return 404 for unknown audit id", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/unknown-audit/results", nil)
						router.ServeHTTP(w, req)

						Expect(w.Code).To(Equal(http.StatusNotFound))
						var response map[string]any
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).ToNot(HaveOccurred())
						Expect(response["status"]).To(Equal("fail"))
						Expect(response["data"]).To(HaveKeyWithValue("message", "Audit not found"))
					})

					It("should return 404 for snake_case audit id", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.AuditBase+"/alert_coverage/results", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})
	})

})
