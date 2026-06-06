//nolint:dupl
package handler_test

import (
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

func decodeAlertTickerResponse(w *httptest.ResponseRecorder, expectedStatus int) barkat.AlertTicker {
	var envelope common.Envelope[barkat.AlertTicker]
	util.AssertSuccess(w, expectedStatus, &envelope)
	return envelope.Data
}

func newAlertTickerTestRouter(alertTickerHandler handler.AlertTickerHandler) *gin.Engine {
	router := util.CreateTestGinRouter()
	tickers := router.Group(barkat.TickerBase)
	handler.SetupTickerAlertRoutes(tickers, alertTickerHandler)
	alertTickers := router.Group(barkat.AlertTickerBase)
	handler.SetupAlertTickerRoutes(alertTickers, alertTickerHandler)
	return router
}

func createAlertTickerRequest(router *gin.Engine, ticker string, payload any) (*httptest.ResponseRecorder, barkat.AlertTicker) {
	req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+ticker+"/alert-tickers", payload)
	router.ServeHTTP(w, req)
	return w, decodeAlertTickerResponse(w, http.StatusCreated)
}

// AlertTickerHandler Integration CUD Tests - Comprehensive Master Specification.
// Tests complete HTTP → Handler → Manager → Repository → Database flow.
// Covers PRD Section 2.2.2 Alert Ticker APIs and Section 2.3.2 Alert Ticker DTO validations.
var _ = Describe("AlertTickerHandler Integration - CUD Tests - Section 2.2.2 Alert Ticker APIs", func() {
	var (
		alertTickerHandler      handler.AlertTickerHandler
		router                  *gin.Engine
		db                      *gorm.DB
		createdTicker           barkat.Ticker
		validAlertTickerPayload barkat.AlertTicker
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
		validAlertTickerPayload = barkat.AlertTicker{
			Symbol:   "MCIX",
			PairID:   "941982",
			Name:     "Multi Commodity Exchange of India",
			Exchange: new("NSE"),
		}
		alertTickerRepo := repository.NewAlertTickerRepository(db)
		alertTickerMgr := manager.NewAlertTickerManager(alertTickerRepo)
		alertTickerHandler = handler.NewAlertTickerHandler(alertTickerMgr)
		router = newAlertTickerTestRouter(alertTickerHandler)
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Describe("POST /v1/api/tickers/{ticker}/alert-tickers - Create Alert Ticker (2.2.2.1)", func() {
		Context("Happy Path", func() {
			Context("with valid alert ticker data", func() {
				var w *httptest.ResponseRecorder
				var response barkat.AlertTicker

				BeforeEach(func() {
					w, response = createAlertTickerRequest(router, createdTicker.Ticker, validAlertTickerPayload)
				})

				It("should return 201 Created", func() { Expect(w.Code).To(Equal(http.StatusCreated)) })
				It("should return Envelope success", func() {
					var envelope common.Envelope[barkat.AlertTicker]
					util.AssertSuccess(w, http.StatusCreated, &envelope)
					Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
				})
				It("should return created alert ticker inside data.alert_ticker", func() { Expect(response.Symbol).To(Equal("MCIX")) })
				It("should preserve symbol", func() { Expect(response.Symbol).To(Equal("MCIX")) })
				It("should preserve pair_id", func() { Expect(response.PairID).To(Equal("941982")) })
				It("should preserve name", func() { Expect(response.Name).To(Equal("Multi Commodity Exchange of India")) })
				It("should preserve exchange", func() { Expect(response.Exchange).To(Equal(new("NSE"))) })
				It("should include parent ticker", func() { Expect(response.TickerSymbol).To(Equal(createdTicker.Ticker)) })
				It("should set created_at timestamp", func() { Expect(response.CreatedAt).ToNot(BeZero()) })
				It("should set updated_at timestamp", func() { Expect(response.UpdatedAt).ToNot(BeZero()) })
				It("should persist alert ticker linked to parent ticker", func() {
					var persisted barkat.AlertTicker
					Expect(db.First(&persisted, "external_id = ?", "MCIX").Error).ToNot(HaveOccurred())
					Expect(persisted.TickerID).To(Equal(createdTicker.ID))
				})
			})
		})

		Context("Field Validations", func() {
			Context("Ticker Path Parameter", func() {
				Context("Allowed Values", func() {
					It("should accept existing valid parent ticker path", func() {
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, validAlertTickerPayload)
						Expect(response.TickerSymbol).To(Equal(createdTicker.Ticker))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for lowercase ticker path", func() {
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/mcx/alert-tickers", validAlertTickerPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for ticker path with whitespace", func() {
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/MC%20X/alert-tickers", validAlertTickerPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for ticker path with unsupported special character", func() {
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/MCX@/alert-tickers", validAlertTickerPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 404 for valid missing parent ticker", func() {
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/NOTFOUND/alert-tickers", validAlertTickerPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})

			Context("Symbol Field", func() {
				Context("Allowed Values", func() {
					It("should accept minimum symbol length 1", func() {
						payload := validAlertTickerPayload
						payload.Symbol = "A"
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Symbol).To(Equal("A"))
					})
					It("should accept maximum symbol length 25", func() {
						payload := validAlertTickerPayload
						payload.Symbol = strings.Repeat("A", 25)
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Symbol).To(HaveLen(25))
					})
					It("should accept letters and digits", func() {
						payload := validAlertTickerPayload
						payload.Symbol = "MCIX50"
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Symbol).To(Equal("MCIX50"))
					})
					It("should accept dot", func() {
						payload := validAlertTickerPayload
						payload.Symbol = "MC.IX"
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Symbol).To(Equal("MC.IX"))
					})
					It("should accept slash", func() {
						payload := validAlertTickerPayload
						payload.Symbol = "MC/IX"
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Symbol).To(Equal("MC/IX"))
					})
					It("should accept equals", func() {
						payload := validAlertTickerPayload
						payload.Symbol = "MC=IX"
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Symbol).To(Equal("MC=IX"))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for missing or empty symbol", func() {
						payload := validAlertTickerPayload
						payload.Symbol = ""
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Symbol", "required")
					})
					It("should return 400 for symbol exceeding 25 characters", func() {
						payload := validAlertTickerPayload
						payload.Symbol = strings.Repeat("A", 26)
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Symbol", "max")
					})
					It("should return 400 for symbol starting with unsupported character", func() {
						payload := validAlertTickerPayload
						payload.Symbol = ".MCIX"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Symbol", "alert_symbol")
					})
					It("should return 400 for symbol with unsupported special character", func() {
						payload := validAlertTickerPayload
						payload.Symbol = "MCIX@"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Symbol", "alert_symbol")
					})
					It("should return 400 for symbol with tab or newline", func() {
						payload := validAlertTickerPayload
						payload.Symbol = "MC\tIX"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Symbol", "alert_symbol")
					})
					It("should return 400 for symbol with underscore", func() {
						payload := validAlertTickerPayload
						payload.Symbol = "MC_IX"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Symbol", "alert_symbol")
					})
					It("should return 400 for symbol with exclamation", func() {
						payload := validAlertTickerPayload
						payload.Symbol = "MCIX!"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Symbol", "alert_symbol")
					})
					It("should return 400 for symbol with space", func() {
						payload := validAlertTickerPayload
						payload.Symbol = "MC IX"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Symbol", "alert_symbol")
					})
					It("should return 400 for symbol with hyphen", func() {
						payload := validAlertTickerPayload
						payload.Symbol = "MC-IX"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Symbol", "alert_symbol")
					})
					It("should return 400 for non-string symbol", func() {
						jsonPayload := `{"symbol":123,"pair_id":"941982","name":"Multi Commodity Exchange of India","exchange":"NSE"}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Pair ID Field", func() {
				Context("Allowed Values", func() {
					It("should accept minimum pair_id length 1", func() {
						payload := validAlertTickerPayload
						payload.PairID = "1"
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.PairID).To(Equal("1"))
					})
					It("should accept maximum pair_id length 64", func() {
						payload := validAlertTickerPayload
						payload.PairID = strings.Repeat("1", 64)
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.PairID).To(HaveLen(64))
					})
					It("should accept digits only", func() {
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, validAlertTickerPayload)
						Expect(response.PairID).To(Equal("941982"))
					})
					It("should preserve leading zeroes", func() {
						payload := validAlertTickerPayload
						payload.PairID = "00123"
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.PairID).To(Equal("00123"))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for missing or empty pair_id", func() {
						payload := validAlertTickerPayload
						payload.PairID = ""
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "PairID", "required")
					})
					It("should return 400 for pair_id exceeding 64 characters", func() {
						payload := validAlertTickerPayload
						payload.PairID = strings.Repeat("1", 65)
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "PairID", "max")
					})
					It("should return 400 for letters in pair_id", func() {
						payload := validAlertTickerPayload
						payload.PairID = "94A982"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "PairID", "number")
					})
					It("should return 400 for decimal pair_id", func() {
						payload := validAlertTickerPayload
						payload.PairID = "941.982"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "PairID", "number")
					})
					It("should return 400 for negative pair_id", func() {
						payload := validAlertTickerPayload
						payload.PairID = "-941982"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "PairID", "number")
					})
					It("should return 400 for whitespace in pair_id", func() {
						payload := validAlertTickerPayload
						payload.PairID = "941 982"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "PairID", "number")
					})
					It("should return 400 for non-string pair_id", func() {
						jsonPayload := `{"symbol":"MCIX","pair_id":941982,"name":"Multi Commodity Exchange of India","exchange":"NSE"}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Name Field", func() {
				Context("Allowed Values", func() {
					It("should accept minimum name length 1", func() {
						payload := validAlertTickerPayload
						payload.Name = "A"
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Name).To(Equal("A"))
					})
					It("should accept maximum name length 100", func() {
						payload := validAlertTickerPayload
						payload.Name = strings.Repeat("A", 100)
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Name).To(HaveLen(100))
					})
					It("should accept spaces", func() {
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, validAlertTickerPayload)
						Expect(response.Name).To(Equal("Multi Commodity Exchange of India"))
					})
					It("should accept dot", func() {
						payload := validAlertTickerPayload
						payload.Name = "M.C.X"
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Name).To(Equal("M.C.X"))
					})
					It("should accept ampersand", func() {
						payload := validAlertTickerPayload
						payload.Name = "M&M"
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Name).To(Equal("M&M"))
					})

					It("should accept apostrophe", func() {
						payload := validAlertTickerPayload
						payload.Name = "Trader's Index"
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Name).To(Equal("Trader's Index"))
					})
					It("should accept parentheses", func() {
						payload := validAlertTickerPayload
						payload.Name = "MCX (India)"
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Name).To(Equal("MCX (India)"))
					})
					It("should accept hyphen", func() {
						payload := validAlertTickerPayload
						payload.Name = "MCX-India"
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Name).To(Equal("MCX-India"))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for missing or empty name", func() {
						payload := validAlertTickerPayload
						payload.Name = ""
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Name", "required")
					})
					It("should return 400 for name exceeding 100 characters", func() {
						payload := validAlertTickerPayload
						payload.Name = strings.Repeat("A", 101)
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Name", "max")
					})
					It("should return 400 for name starting with unsupported character", func() {
						payload := validAlertTickerPayload
						payload.Name = ".MCX"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Name", "alert_name")
					})
					It("should return 400 for comma in name", func() {
						payload := validAlertTickerPayload
						payload.Name = "MCX, India"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Name", "alert_name")
					})
					It("should return 400 for underscore in name", func() {
						payload := validAlertTickerPayload
						payload.Name = "MCX_India"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Name", "alert_name")
					})
					It("should return 400 for slash in name", func() {
						payload := validAlertTickerPayload
						payload.Name = "MCX/India"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Name", "alert_name")
					})
					It("should return 400 for at-sign in name", func() {
						payload := validAlertTickerPayload
						payload.Name = "MCX@India"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Name", "alert_name")
					})
					It("should return 400 for non-string name", func() {
						jsonPayload := `{"symbol":"MCIX","pair_id":"941982","name":123,"exchange":"NSE"}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Exchange Field", func() {
				Context("Allowed Values", func() {
					It("should accept omitted exchange", func() {
						payload := validAlertTickerPayload
						payload.Exchange = nil
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Exchange).To(BeNil())
					})
					It("should accept null exchange", func() {
						payload := validAlertTickerPayload
						payload.Exchange = nil
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Exchange).To(BeNil())
					})
					It("should accept minimum exchange length 1", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new("N")
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(*response.Exchange).To(Equal("N"))
					})
					It("should accept maximum exchange length 10", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new(strings.Repeat("A", 10))
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(*response.Exchange).To(HaveLen(10))
					})
					It("should accept uppercase letters NSE", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new("NSE")
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(*response.Exchange).To(Equal("NSE"))
					})
					It("should accept lowercase letters nse", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new("nse")
						_, response := createAlertTickerRequest(router, createdTicker.Ticker, payload)
						Expect(*response.Exchange).To(Equal("nse"))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for empty exchange string", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new("")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "min")
					})
					It("should return 400 for exchange exceeding 10 characters", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new(strings.Repeat("A", 11))
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "max")
					})
					It("should return 400 for exchange with colon", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new("NSE:MCX")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "alert_exchange")
					})
					It("should return 400 for exchange with hyphen", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new("N-SE")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "alert_exchange")
					})
					It("should return 400 for exchange with whitespace", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new("NS E")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "alert_exchange")
					})
					It("should return 400 for exchange with unsupported special character", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new("NSE@")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "alert_exchange")
					})
					It("should return 400 for exchange with digit", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new("NSE1")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "alert_exchange")
					})
					It("should return 400 for exchange with dot", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new("N.SE")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "alert_exchange")
					})
					It("should return 400 for exchange with underscore", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new("N_SE")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "alert_exchange")
					})
					It("should return 400 for exchange starting with digit", func() {
						payload := validAlertTickerPayload
						payload.Exchange = new("1abc")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "alert_exchange")
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 400 for malformed JSON", func() {
				req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", `{"symbol":`)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 400 for empty request body", func() {
				req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", "")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 400 for null request body", func() {
				req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", "null")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 409 for duplicate symbol", func() {
				payload := validAlertTickerPayload
				_, _ = createAlertTickerRequest(router, createdTicker.Ticker, payload)
				req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", payload)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusConflict))
			})
			It("should return 409 for duplicate pair_id", func() {
				first := validAlertTickerPayload
				first.Symbol = "MCIXD1"
				_, _ = createAlertTickerRequest(router, createdTicker.Ticker, first)
				second := validAlertTickerPayload
				second.Symbol = "MCIXD2"
				req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", second)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusConflict))
			})
			It("should return 500 for persistence failure", func() {
				sqlDB, err := db.DB()
				Expect(err).ToNot(HaveOccurred())
				Expect(sqlDB.Close()).To(Succeed())
				req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase+"/"+createdTicker.Ticker+"/alert-tickers", validAlertTickerPayload)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("DELETE /v1/api/alert-tickers/{symbol} - Delete Alert Ticker (2.2.2.3)", func() {
		var createdAlertTicker barkat.AlertTicker

		BeforeEach(func() {
			createdAlertTicker = validAlertTickerPayload
			createdAlertTicker.TickerID = createdTicker.ID
			Expect(db.Create(&createdAlertTicker).Error).ToNot(HaveOccurred())
		})

		Context("Happy Path", func() {
			Context("with existing alert ticker", func() {
				It("should return 204 No Content", func() {
					req, w := util.CreateTestRequest(http.MethodDelete, barkat.AlertTickerBase+"/"+createdAlertTicker.Symbol, nil)
					router.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusNoContent))
				})
				It("should return empty body", func() {
					req, w := util.CreateTestRequest(http.MethodDelete, barkat.AlertTickerBase+"/"+createdAlertTicker.Symbol, nil)
					router.ServeHTTP(w, req)
					Expect(w.Body.String()).To(BeEmpty())
				})
				It("should delete alert ticker from database", func() {
					req, w := util.CreateTestRequest(http.MethodDelete, barkat.AlertTickerBase+"/"+createdAlertTicker.Symbol, nil)
					router.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusNoContent))
					var persisted barkat.AlertTicker
					Expect(db.First(&persisted, "symbol = ?", createdAlertTicker.Symbol).Error).To(HaveOccurred())
				})
			})
		})

		Context("Field Validations", func() {
			Context("Symbol Path Parameter", func() {
				Context("Allowed Values", func() {
					It("should accept valid existing symbol path", func() {
						req, w := util.CreateTestRequest(http.MethodDelete, barkat.AlertTickerBase+"/"+createdAlertTicker.Symbol, nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNoContent))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for invalid symbol path", func() {
						req, w := util.CreateTestRequest(http.MethodDelete, barkat.AlertTickerBase+"/.MCIX", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 404 for valid missing symbol path", func() {
						req, w := util.CreateTestRequest(http.MethodDelete, barkat.AlertTickerBase+"/NOTFOUND", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 404 on second delete", func() {
				req, w := util.CreateTestRequest(http.MethodDelete, barkat.AlertTickerBase+"/"+createdAlertTicker.Symbol, nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNoContent))
				req, w = util.CreateTestRequest(http.MethodDelete, barkat.AlertTickerBase+"/"+createdAlertTicker.Symbol, nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
			It("should return 500 for delete failure", func() {
				sqlDB, err := db.DB()
				Expect(err).ToNot(HaveOccurred())
				Expect(sqlDB.Close()).To(Succeed())
				req, w := util.CreateTestRequest(http.MethodDelete, barkat.AlertTickerBase+"/"+createdAlertTicker.Symbol, nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
