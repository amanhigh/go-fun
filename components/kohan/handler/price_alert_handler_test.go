package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

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

func decodePriceAlertResponse(w *httptest.ResponseRecorder, expectedStatus int) barkat.PriceAlert {
	var envelope common.Envelope[barkat.PriceAlert]
	util.AssertSuccess(w, expectedStatus, &envelope)
	return envelope.Data
}

func decodePriceAlertReplaceResponse(w *httptest.ResponseRecorder) barkat.PriceAlertReplaceResult {
	var envelope common.Envelope[barkat.PriceAlertReplaceResult]
	util.AssertSuccess(w, http.StatusOK, &envelope)
	return envelope.Data
}

func decodePriceAlertListResponse(w *httptest.ResponseRecorder) barkat.PriceAlertList {
	var envelope common.Envelope[barkat.PriceAlertList]
	util.AssertSuccess(w, http.StatusOK, &envelope)
	return envelope.Data
}

func newPriceAlertTestRouter(alertTickerHandler handler.AlertTickerHandler, priceAlertHandler handler.PriceAlertHandler) *gin.Engine {
	router := util.CreateTestGinRouter()
	tickers := router.Group(barkat.TickerBase)
	handler.SetupTickerAlertRoutes(tickers, alertTickerHandler)
	handler.SetupTickerPriceAlertRoutes(tickers, priceAlertHandler)
	alertTickers := router.Group(barkat.AlertTickerBase)
	handler.SetupAlertTickerRoutes(alertTickers, alertTickerHandler)
	alerts := router.Group(barkat.PriceAlertBase)
	handler.SetupPriceAlertRoutes(alerts, priceAlertHandler)
	return router
}

func seedPriceAlert(ctx context.Context, db *gorm.DB, alert barkat.PriceAlert) barkat.PriceAlert {
	Expect(db.WithContext(ctx).Create(&alert).Error).ToNot(HaveOccurred())
	return alert
}

// PriceAlertHandler Integration Tests - Comprehensive Master Specification.
// Tests complete HTTP → Handler → Manager → Repository → Database flow.
// Covers PRD Section 2.2.3 Price Alert APIs and Section 2.3.3 Price Alert DTO validations.
var _ = Describe("PriceAlertHandler Integration - Section 2.2.3 Price Alert APIs", func() {
	var (
		alertTickerHandler handler.AlertTickerHandler
		priceAlertHandler  handler.PriceAlertHandler
		router             *gin.Engine
		testCtx            = context.Background()
		db                 *gorm.DB
		createdTicker      barkat.Ticker
		createdAlertTicker barkat.AlertTicker
		req                *http.Request
		w                  *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		var err error
		core.RegisterJournalValidators()
		db, err = core.CreateTestBarkatDB()
		Expect(err).ToNot(HaveOccurred())

		createdTicker = barkat.Ticker{
			Ticker:       "MCX",
			Exchange:     new("NSE"),
			Timeframes:   []string{"MN", "WK", "DL"},
			Type:         "EQUITY",
			State:        "WATCHED",
			Trend:        "UPTREND",
			LastOpenedAt: time.Date(2026, time.May, 5, 10, 30, 0, 0, time.UTC),
			IsFNO:        true,
		}
		Expect(db.Create(&createdTicker).Error).ToNot(HaveOccurred())

		createdAlertTicker = barkat.AlertTicker{
			TickerID: createdTicker.ID,
			Symbol:   "MCIX",
			PairID:   "941982",
			Name:     "Multi Commodity Exchange of India",
			Exchange: new("NSE"),
		}
		Expect(db.Create(&createdAlertTicker).Error).ToNot(HaveOccurred())

		alertTickerRepo := repository.NewAlertTickerRepository(db)
		alertTickerMgr := manager.NewAlertTickerManager(alertTickerRepo)
		alertTickerHandler = handler.NewAlertTickerHandler(alertTickerMgr)

		priceAlertRepo := repository.NewPriceAlertRepository(db)
		priceAlertMgr := manager.NewPriceAlertManager(priceAlertRepo)
		priceAlertHandler = handler.NewPriceAlertHandler(priceAlertMgr)
		router = newPriceAlertTestRouter(alertTickerHandler, priceAlertHandler)
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	// ============================================================================
	// 2.2.3.1 PUT /v1/api/alerts - Replace Price Alerts
	// ============================================================================
	Describe("PUT /v1/api/alerts - Replace Price Alerts (2.2.3.1)", func() {
		var replacePayload barkat.PriceAlertReplaceRequest

		BeforeEach(func() {
			replacePayload = barkat.PriceAlertReplaceRequest{Alerts: []barkat.PriceAlertInput{
				{PairID: "941982", AlertID: "158741518", TriggerPrice: 1.0632},
				{PairID: "941982", AlertID: "158741514", TriggerPrice: 1.2401},
			}}
		})

		Context("Happy Path", func() {
			Context("with complete refreshed alert rows", func() {
				var response barkat.PriceAlertReplaceResult

				BeforeEach(func() {
					req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, replacePayload)
					router.ServeHTTP(w, req)
					response = decodePriceAlertReplaceResponse(w)
				})

				It("should return 200 OK", func() { Expect(w.Code).To(Equal(http.StatusOK)) })
				It("should return Envelope success", func() {
					var envelope common.Envelope[barkat.PriceAlertReplaceResult]
					util.AssertSuccess(w, http.StatusOK, &envelope)
					Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
				})
				It("should report replaced pair count", func() { Expect(response.PairsReplaced).To(Equal(1)) })
				It("should report created alert count", func() { Expect(response.AlertsCreated).To(Equal(2)) })
				It("should persist price alerts under the resolved alert ticker", func() {
					var persisted []barkat.PriceAlert
					Expect(db.Where("alert_ticker_id = ?", createdAlertTicker.ID).Find(&persisted).Error).ToNot(HaveOccurred())
					Expect(persisted).To(HaveLen(2))
				})
			})

			It("should replace existing alerts for submitted pair ids only", func() {
				otherTicker := barkat.Ticker{Ticker: "NIFTY", Exchange: new("NSE"), Timeframes: []string{"MN"}, Type: "EQUITY", State: "WATCHED", Trend: "SIDEWAYS", LastOpenedAt: time.Now()}
				Expect(db.Create(&otherTicker).Error).ToNot(HaveOccurred())
				otherAlertTicker := barkat.AlertTicker{TickerID: otherTicker.ID, Symbol: "NIFTY50", PairID: "17940", Name: "Nifty 50"}
				Expect(db.Create(&otherAlertTicker).Error).ToNot(HaveOccurred())
				oldID := "111111"
				otherID := "222222"
				seedPriceAlert(testCtx, db, barkat.PriceAlert{AlertTickerID: createdAlertTicker.ID, AlertID: &oldID, TriggerPrice: 10})
				seedPriceAlert(testCtx, db, barkat.PriceAlert{AlertTickerID: otherAlertTicker.ID, AlertID: &otherID, TriggerPrice: 20})

				req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, replacePayload)
				router.ServeHTTP(w, req)
				response := decodePriceAlertReplaceResponse(w)
				Expect(response.AlertsCreated).To(Equal(2))
				var oldCount int64
				Expect(db.Model(&barkat.PriceAlert{}).Where("alert_id = ?", oldID).Count(&oldCount).Error).ToNot(HaveOccurred())
				Expect(oldCount).To(Equal(int64(0)))
				var otherCount int64
				Expect(db.Model(&barkat.PriceAlert{}).Where("alert_id = ?", otherID).Count(&otherCount).Error).ToNot(HaveOccurred())
				Expect(otherCount).To(Equal(int64(1)))
			})

			It("should replace alerts across multiple submitted pair ids", func() {
				otherTicker := barkat.Ticker{Ticker: "NIFTY", Exchange: new("NSE"), Timeframes: []string{"MN"}, Type: "EQUITY", State: "WATCHED", Trend: "SIDEWAYS", LastOpenedAt: time.Now()}
				Expect(db.Create(&otherTicker).Error).ToNot(HaveOccurred())
				otherAlertTicker := barkat.AlertTicker{TickerID: otherTicker.ID, Symbol: "NIFTY50", PairID: "17940", Name: "Nifty 50"}
				Expect(db.Create(&otherAlertTicker).Error).ToNot(HaveOccurred())

				multiPairPayload := barkat.PriceAlertReplaceRequest{Alerts: []barkat.PriceAlertInput{
					{PairID: createdAlertTicker.PairID, AlertID: "158741518", TriggerPrice: 1.0632},
					{PairID: otherAlertTicker.PairID, AlertID: "158741515", TriggerPrice: 200.50},
				}}
				req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, multiPairPayload)
				router.ServeHTTP(w, req)
				response := decodePriceAlertReplaceResponse(w)
				Expect(response.PairsReplaced).To(Equal(2))
				Expect(response.AlertsCreated).To(Equal(2))
			})
		})

		Context("Field Validations", func() {
			Context("Alerts Field", func() {
				Context("Allowed Values", func() {
					It("should accept empty alerts array", func() {
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, barkat.PriceAlertReplaceRequest{Alerts: []barkat.PriceAlertInput{}})
						router.ServeHTTP(w, req)
						response := decodePriceAlertReplaceResponse(w)
						Expect(response.PairsReplaced).To(Equal(0))
						Expect(response.AlertsCreated).To(Equal(0))
					})

					It("should accept exactly 100 alerts", func() {
						alerts := make([]barkat.PriceAlertInput, 100)
						for i := range alerts {
							alerts[i] = barkat.PriceAlertInput{PairID: "941982", AlertID: fmt.Sprintf("%d", 100000+i), TriggerPrice: 1}
						}
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, barkat.PriceAlertReplaceRequest{Alerts: alerts})
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing alerts", func() {
						req, w = rawTickerRequest(http.MethodPut, barkat.PriceAlertBase, `{}`)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Alerts", "required")
					})

					It("should return 413 for more than 100 alerts", func() {
						alerts := make([]barkat.PriceAlertInput, 101)
						for i := range alerts {
							alerts[i] = barkat.PriceAlertInput{PairID: "941982", AlertID: fmt.Sprintf("%d", 100000+i), TriggerPrice: 1}
						}
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, barkat.PriceAlertReplaceRequest{Alerts: alerts})
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusRequestEntityTooLarge))
					})
				})
			})

			Context("Pair ID Field", func() {
				Context("Allowed Values", func() {
					It("should accept minimum pair_id length 1", func() {
						minPairAlertTicker := barkat.AlertTicker{TickerID: createdTicker.ID, Symbol: "MINP", PairID: "1", Name: "Min Pair", Exchange: new("NSE")}
						Expect(db.Create(&minPairAlertTicker).Error).ToNot(HaveOccurred())

						payload := replacePayload
						payload.Alerts[0].PairID = "1"
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})

					It("should accept maximum pair_id length 64", func() {
						maxPair := strings.Repeat("1", 64)
						maxPairAlertTicker := barkat.AlertTicker{TickerID: createdTicker.ID, Symbol: "MAXP", PairID: maxPair, Name: "Max Pair", Exchange: new("NSE")}
						Expect(db.Create(&maxPairAlertTicker).Error).ToNot(HaveOccurred())

						payload := replacePayload
						payload.Alerts[0].PairID = maxPair
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing pair_id", func() {
						payload := replacePayload
						payload.Alerts[0].PairID = ""
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "PairID", "required")
					})

					It("should return 400 for pair_id exceeding 64 characters", func() {
						payload := replacePayload
						payload.Alerts[0].PairID = strings.Repeat("1", 65)
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "PairID", "max")
					})

					It("should return 400 for non-digit pair_id", func() {
						payload := replacePayload
						payload.Alerts[0].PairID = "94A982"
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "PairID", "number")
					})

					It("should return 404 for unresolved pair_id", func() {
						payload := replacePayload
						payload.Alerts[0].PairID = "999999"
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})

			Context("Alert ID Field", func() {
				Context("Allowed Values", func() {
					It("should accept minimum alert_id length 1", func() {
						payload := replacePayload
						payload.Alerts[0].AlertID = "1"
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})

					It("should accept maximum alert_id length 128", func() {
						payload := replacePayload
						payload.Alerts[0].AlertID = strings.Repeat("1", 128)
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing alert_id", func() {
						payload := replacePayload
						payload.Alerts[0].AlertID = ""
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "AlertID", "required")
					})

					It("should return 400 for alert_id exceeding 128 characters", func() {
						payload := replacePayload
						payload.Alerts[0].AlertID = strings.Repeat("1", 129)
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "AlertID", "max")
					})

					It("should return 400 for non-digit alert_id", func() {
						payload := replacePayload
						payload.Alerts[0].AlertID = "158A"
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "AlertID", "number")
					})

					It("should return 409 for duplicate alert_id within the request", func() {
						payload := replacePayload
						payload.Alerts[1].AlertID = payload.Alerts[0].AlertID
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusConflict))
					})
				})
			})

			Context("Trigger Price Field", func() {
				Context("Allowed Values", func() {
					It("should accept positive decimal trigger_price", func() {
						payload := replacePayload
						payload.Alerts[0].TriggerPrice = 0.0001
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						response := decodePriceAlertReplaceResponse(w)
						Expect(response.AlertsCreated).To(Equal(2))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing trigger_price", func() {
						req, w = rawTickerRequest(http.MethodPut, barkat.PriceAlertBase, `{"alerts":[{"pair_id":"941982","alert_id":"158741518"}]}`)
						router.ServeHTTP(w, req)
						util.AssertError(w, "TriggerPrice", "required")
					})

					It("should return 400 for zero trigger_price", func() {
						payload := replacePayload
						payload.Alerts[0].TriggerPrice = 0
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "TriggerPrice", "required")
					})

					It("should return 400 for negative trigger_price", func() {
						payload := replacePayload
						payload.Alerts[0].TriggerPrice = -1
						req, w = util.CreateTestRequest(http.MethodPut, barkat.PriceAlertBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "TriggerPrice", "gt")
					})

					It("should return 400 for non-numeric trigger_price", func() {
						req, w = rawTickerRequest(http.MethodPut, barkat.PriceAlertBase, `{"alerts":[{"pair_id":"941982","alert_id":"158741518","trigger_price":"1.2"}]}`)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 400 for malformed JSON", func() {
				req, w = rawTickerRequest(http.MethodPut, barkat.PriceAlertBase, `{"alerts":[`)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for null body", func() {
				req, w = rawTickerRequest(http.MethodPut, barkat.PriceAlertBase, "null")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	// ============================================================================
	// 2.2.3.2 POST /v1/api/tickers/{ticker}/alerts - Create Pending Price Alert
	// ============================================================================
	Describe("POST /v1/api/tickers/{ticker}/alerts - Create Pending Price Alert (2.2.3.2)", func() {
		Context("Happy Path", func() {
			var response barkat.PriceAlert

			BeforeEach(func() {
				req, w = util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alerts", barkat.PendingPriceAlertRequest{TriggerPrice: 42.25})
				router.ServeHTTP(w, req)
				response = decodePriceAlertResponse(w, http.StatusCreated)
			})

			It("should return 201 Created", func() { Expect(w.Code).To(Equal(http.StatusCreated)) })
			It("should return Envelope success", func() {
				var envelope common.Envelope[barkat.PriceAlert]
				util.AssertSuccess(w, http.StatusCreated, &envelope)
				Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			})
			It("should derive pair_id from the parent ticker mapping", func() { Expect(response.PairID).To(Equal(createdAlertTicker.PairID)) })
			It("should omit canonical alert_id until refresh", func() { Expect(response.AlertID).To(BeNil()) })
			It("should preserve trigger_price", func() { Expect(response.TriggerPrice).To(Equal(42.25)) })
			It("should persist pending alert under the first alert ticker", func() {
				var count int64
				Expect(db.Model(&barkat.PriceAlert{}).Where("alert_ticker_id = ? AND alert_id IS NULL", createdAlertTicker.ID).Count(&count).Error).ToNot(HaveOccurred())
				Expect(count).To(Equal(int64(1)))
			})
		})

		Context("Field Validations", func() {
			Context("Ticker Path Parameter", func() {
				Context("Allowed Values", func() {
					It("should accept existing valid parent ticker path", func() {
						req, w = util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alerts", barkat.PendingPriceAlertRequest{TriggerPrice: 1})
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusCreated))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid ticker path", func() {
						req, w = util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/mcx/alerts", barkat.PendingPriceAlertRequest{TriggerPrice: 1})
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 404 when primary ticker is missing", func() {
						req, w = util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/NOTFOUND/alerts", barkat.PendingPriceAlertRequest{TriggerPrice: 1})
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 404 when ticker has no alert ticker mapping", func() {
						other := barkat.Ticker{Ticker: "NIFTY", Exchange: new("NSE"), Timeframes: []string{"MN"}, Type: "EQUITY", State: "WATCHED", Trend: "SIDEWAYS", LastOpenedAt: time.Now()}
						Expect(db.Create(&other).Error).ToNot(HaveOccurred())
						req, w = util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/NIFTY/alerts", barkat.PendingPriceAlertRequest{TriggerPrice: 1})
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})

			Context("Trigger Price Field", func() {
				Context("Allowed Values", func() {
					It("should accept positive decimal trigger_price", func() {
						req, w = util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alerts", barkat.PendingPriceAlertRequest{TriggerPrice: 0.0001})
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusCreated))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing trigger_price", func() {
						req, w = rawTickerRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alerts", `{}`)
						router.ServeHTTP(w, req)
						util.AssertError(w, "TriggerPrice", "required")
					})

					It("should return 400 for zero trigger_price", func() {
						req, w = util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alerts", barkat.PendingPriceAlertRequest{TriggerPrice: 0})
						router.ServeHTTP(w, req)
						util.AssertError(w, "TriggerPrice", "required")
					})

					It("should return 400 for negative trigger_price", func() {
						req, w = util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alerts", barkat.PendingPriceAlertRequest{TriggerPrice: -1})
						router.ServeHTTP(w, req)
						util.AssertError(w, "TriggerPrice", "gt")
					})

					It("should return 400 for non-numeric trigger_price", func() {
						req, w = rawTickerRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alerts", `{"trigger_price":"abc"}`)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 400 for malformed JSON", func() {
				req, w = rawTickerRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alerts", `{"trigger_price":`)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 for null body", func() {
				req, w = rawTickerRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alerts", "null")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	// ============================================================================
	// 2.2.3.3 DELETE /v1/api/alerts/{alert-id} - Delete Price Alert
	// ============================================================================
	Describe("DELETE /v1/api/alerts/{alert-id} - Delete Price Alert (2.2.3.3)", func() {
		var alertID string

		BeforeEach(func() {
			alertID = "158741518"
			seedPriceAlert(testCtx, db, barkat.PriceAlert{AlertTickerID: createdAlertTicker.ID, AlertID: &alertID, TriggerPrice: 1.0632})
		})

		Context("Happy Path", func() {
			BeforeEach(func() {
				req, w = util.CreateTestRequest(http.MethodDelete, barkat.PriceAlertBase+"/"+alertID, nil)
				router.ServeHTTP(w, req)
			})

			It("should return 204 No Content", func() {
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})

			It("should delete alert from database", func() {
				var count int64
				Expect(db.Model(&barkat.PriceAlert{}).Where("alert_id = ?", alertID).Count(&count).Error).ToNot(HaveOccurred())
				Expect(count).To(Equal(int64(0)))
			})
		})

		Context("Field Validations", func() {
			Context("Alert ID Path Parameter", func() {
				Context("Allowed Values", func() {
					It("should accept valid numeric alert-id path", func() {
						req, w = util.CreateTestRequest(http.MethodDelete, barkat.PriceAlertBase+"/"+alertID, nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNoContent))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for invalid alert-id path", func() {
						req, w = util.CreateTestRequest(http.MethodDelete, barkat.PriceAlertBase+"/ABC", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 400 for alert-id exceeding 128 characters", func() {
						req, w = util.CreateTestRequest(http.MethodDelete, barkat.PriceAlertBase+"/"+strings.Repeat("1", 129), nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})

					It("should return 404 for missing alert-id", func() {
						req, w = util.CreateTestRequest(http.MethodDelete, barkat.PriceAlertBase+"/999999", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})
	})

	// ============================================================================
	// 2.2.3.4 GET /v1/api/alerts - List Price Alerts
	// ============================================================================
	Describe("GET /v1/api/alerts - List Price Alerts (2.2.3.4)", func() {
		BeforeEach(func() {
			firstID := "158741518"
			secondID := "158741514"
			seedPriceAlert(testCtx, db, barkat.PriceAlert{AlertTickerID: createdAlertTicker.ID, AlertID: &firstID, TriggerPrice: 1.0632})
			seedPriceAlert(testCtx, db, barkat.PriceAlert{AlertTickerID: createdAlertTicker.ID, AlertID: &secondID, TriggerPrice: 1.2401})
		})

		Context("Happy Path", func() {
			var response barkat.PriceAlertList

			BeforeEach(func() {
				req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase, nil)
				router.ServeHTTP(w, req)
				response = decodePriceAlertListResponse(w)
			})

			It("should return 200 OK", func() { Expect(w.Code).To(Equal(http.StatusOK)) })
			It("should return Envelope success", func() {
				var envelope common.Envelope[barkat.PriceAlertList]
				util.AssertSuccess(w, http.StatusOK, &envelope)
				Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
			})
			It("should return alerts array with default pagination", func() {
				Expect(response.PriceAlerts).To(HaveLen(2))
				Expect(response.Metadata.Offset).To(Equal(0))
				Expect(response.Metadata.Limit).To(Equal(10))
				Expect(response.Metadata.Total).To(Equal(int64(2)))
			})
			It("should include pair_id, alert_id, trigger_price, and created_at in each alert", func() {
				for _, alert := range response.PriceAlerts {
					Expect(alert.PairID).To(Equal(createdAlertTicker.PairID))
					Expect(alert.AlertID).ToNot(BeNil())
					Expect(alert.TriggerPrice).To(BeNumerically(">", 0))
					Expect(alert.CreatedAt).ToNot(BeZero())
				}
			})

			It("should filter by primary ticker without leaking other ticker alerts", func() {
				otherTicker := barkat.Ticker{Ticker: "NIFTY", Exchange: new("NSE"), Timeframes: []string{"MN"}, Type: "EQUITY", State: "WATCHED", Trend: "SIDEWAYS", LastOpenedAt: time.Now()}
				Expect(db.Create(&otherTicker).Error).ToNot(HaveOccurred())
				otherAlertTicker := barkat.AlertTicker{TickerID: otherTicker.ID, Symbol: "NIFTY50", PairID: "17940", Name: "Nifty 50"}
				Expect(db.Create(&otherAlertTicker).Error).ToNot(HaveOccurred())
				otherAlertID := "222222"
				seedPriceAlert(testCtx, db, barkat.PriceAlert{AlertTickerID: otherAlertTicker.ID, AlertID: &otherAlertID, TriggerPrice: 20})

				req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?ticker=MCX", nil)
				router.ServeHTTP(w, req)
				listResponse := decodePriceAlertListResponse(w)
				Expect(listResponse.PriceAlerts).To(HaveLen(2))
				for _, alert := range listResponse.PriceAlerts {
					Expect(alert.PairID).To(Equal(createdAlertTicker.PairID))
				}
			})

			It("should sort by trigger_price descending", func() {
				req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?sort-by=trigger_price&sort-order=desc", nil)
				router.ServeHTTP(w, req)
				listResponse := decodePriceAlertListResponse(w)
				Expect(listResponse.PriceAlerts[0].TriggerPrice).To(BeNumerically(">", listResponse.PriceAlerts[1].TriggerPrice))
			})

			It("should paginate using limit and offset", func() {
				req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?offset=1&limit=1", nil)
				router.ServeHTTP(w, req)
				listResponse := decodePriceAlertListResponse(w)
				Expect(listResponse.PriceAlerts).To(HaveLen(1))
				Expect(listResponse.Metadata.Offset).To(Equal(1))
				Expect(listResponse.Metadata.Limit).To(Equal(1))
			})
		})

		Context("Field Validations", func() {
			Context("Ticker Filter", func() {
				Context("Allowed Values", func() {
					It("should filter by existing primary ticker", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?ticker=MCX", nil)
						router.ServeHTTP(w, req)
						listResponse := decodePriceAlertListResponse(w)
						Expect(listResponse.PriceAlerts).To(HaveLen(2))
					})
				})

				Context("Bad Values", func() {
					It("should return 404 when ticker filter references missing parent", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?ticker=NOTFOUND", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})

					It("should return 400 for invalid ticker filter", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?ticker=mcx", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Sort By Field", func() {
				Context("Allowed Values", func() {
					It("should sort by trigger_price ascending (default)", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?sort-by=trigger_price&sort-order=asc", nil)
						router.ServeHTTP(w, req)
						listResponse := decodePriceAlertListResponse(w)
						Expect(listResponse.PriceAlerts[0].TriggerPrice).To(BeNumerically("<", listResponse.PriceAlerts[1].TriggerPrice))
					})

					It("should sort by created_at", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?sort-by=created_at&sort-order=asc", nil)
						router.ServeHTTP(w, req)
						listResponse := decodePriceAlertListResponse(w)
						Expect(listResponse.PriceAlerts).To(HaveLen(2))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for unsupported sort-by", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?sort-by=alert_id", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Sort Order Field", func() {
				Context("Allowed Values", func() {
					It("should accept asc", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?sort-order=asc", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})

					It("should accept desc", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?sort-order=desc", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for unsupported sort-order", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?sort-order=up", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Offset Field", func() {
				Context("Allowed Values", func() {
					It("should accept positive offset with pagination", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?offset=1&limit=1", nil)
						router.ServeHTTP(w, req)
						listResponse := decodePriceAlertListResponse(w)
						Expect(listResponse.PriceAlerts).To(HaveLen(1))
						Expect(listResponse.Metadata.Offset).To(Equal(1))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for negative offset", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?offset=-1", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Offset", "min")
					})
				})
			})

			Context("Limit Field", func() {
				Context("Allowed Values", func() {
					It("should accept max limit 10", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?limit=10", nil)
						router.ServeHTTP(w, req)
						listResponse := decodePriceAlertListResponse(w)
						Expect(listResponse.Metadata.Limit).To(Equal(10))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for zero limit", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?limit=0", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Limit", "min")
					})

					It("should return 400 for limit greater than 10", func() {
						req, w = util.CreateTestRequest(http.MethodGet, barkat.PriceAlertBase+"?limit=11", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Limit", "max")
					})
				})
			})
		})
	})
})
