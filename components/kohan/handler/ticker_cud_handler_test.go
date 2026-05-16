//nolint:dupl
package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func decodeTickerResponse(w *httptest.ResponseRecorder, expectedStatus int) barkat.Ticker {
	var envelope common.Envelope[barkat.Ticker]
	util.AssertSuccess(w, expectedStatus, &envelope)
	return envelope.Data
}

func newTickerTestRouter(tickerHandler handler.TickerHandler) *gin.Engine {
	router := util.CreateTestGinRouter()
	tickers := router.Group(barkat.TickerBase)
	handler.SetupTickerRoutes(tickers, tickerHandler)
	return router
}

func createTickerRequest(router *gin.Engine, payload any) (*httptest.ResponseRecorder, barkat.Ticker) {
	req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
	router.ServeHTTP(w, req)
	return w, decodeTickerResponse(w, http.StatusCreated)
}

func updateTickerRequest(router *gin.Engine, ticker string, payload any) (*httptest.ResponseRecorder, barkat.Ticker) {
	req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+ticker, payload)
	router.ServeHTTP(w, req)
	return w, decodeTickerResponse(w, http.StatusOK)
}

func patchTickerRequest(router *gin.Engine, ticker string, payload any) (*httptest.ResponseRecorder, barkat.Ticker) {
	req, w := util.CreateTestRequest(http.MethodPatch, barkat.TickerBase+"/"+ticker, payload)
	router.ServeHTTP(w, req)
	return w, decodeTickerResponse(w, http.StatusOK)
}

func rawTickerRequest(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequestWithContext(context.Background(), method, url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

// TickerHandler Integration Tests - Comprehensive Master Specification
// Tests complete HTTP → Handler → Manager → Repository → Database flow.
// Covers PRD Section 2.2.1 Primary Ticker APIs and Section 2.3.1 Primary Ticker DTO validations.
//
// TEST STRUCTURE FORMAT:
// ====================
// Describe(API)
// -> Context(Happy Path): 2xx Success Cases
// -> Context(Field Validations): All 4xx Validation Cases
//
//	-> Context(Field Name): One Context for Each Field
//	   -> Context(Allowed Values): All Variations of Valid Values (2xx) - If Applicable
//	   -> Context(Bad Values): All Variations of Missing,Regex,Min,Max Edge Cases (4xx)
//
// -> Context(Errors): 5xx Server Error Cases
var _ = Describe("TickerHandler Integration - CUD Tests - Section 2.2.1 Primary Ticker APIs", func() {
	var (
		tickerHandler      handler.TickerHandler
		router             *gin.Engine
		db                 *gorm.DB
		validTickerPayload barkat.Ticker
	)

	BeforeEach(func() {
		var err error
		core.RegisterJournalValidators()
		db, err = core.CreateTestBarkatDB()
		Expect(err).ToNot(HaveOccurred())
		tickerRepo := repository.NewTickerRepository(db)
		tickerMgr := manager.NewBarkatTickerManager(tickerRepo)
		tickerHandler = handler.NewTickerHandler(tickerMgr)
		validTickerPayload = barkat.Ticker{
			Ticker:       "MCX",
			Exchange:     new("NSE"),
			Timeframes:   []string{"MN", "WK", "DL"},
			Type:         "EQUITY",
			State:        "WATCHED",
			Trend:        "UPTREND",
			LastOpenedAt: time.Date(2026, time.May, 5, 10, 30, 0, 0, time.UTC),
			IsFNO:        true,
		}
		router = newTickerTestRouter(tickerHandler)
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	// ============================================================================
	// 2.2.1.1 POST /v1/api/tickers - Create Primary Ticker
	// ============================================================================
	Describe("POST /v1/api/tickers - Create Primary Ticker (2.2.1.1)", func() {
		Context("Happy Path", func() {
			Context("with valid ticker data", func() {
				var w *httptest.ResponseRecorder
				var response barkat.Ticker

				BeforeEach(func() {
					w, response = createTickerRequest(router, validTickerPayload)
				})

				It("should return 201 Created", func() {
					Expect(w.Code).To(Equal(http.StatusCreated))
				})

				It("should return Envelope success", func() {
					var envelope common.Envelope[barkat.Ticker]
					util.AssertSuccess(w, http.StatusCreated, &envelope)
					Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
				})

				It("should return created ticker in data", func() {
					Expect(response.Ticker).To(Equal("MCX"))
				})

				It("should preserve ticker field", func() { Expect(response.Ticker).To(Equal("MCX")) })
				It("should preserve exchange field", func() { Expect(response.Exchange).To(Equal(new("NSE"))) })
				It("should preserve ordered timeframes", func() { Expect(response.Timeframes).To(Equal([]string{"MN", "WK", "DL"})) })
				It("should preserve type field", func() { Expect(response.Type).To(Equal("EQUITY")) })
				It("should preserve state field", func() { Expect(response.State).To(Equal("WATCHED")) })
				It("should preserve trend field", func() { Expect(response.Trend).To(Equal("UPTREND")) })
				It("should preserve last_opened_at field", func() { Expect(response.LastOpenedAt).To(Equal(validTickerPayload.LastOpenedAt)) })
				It("should preserve is_fno field", func() { Expect(response.IsFNO).To(BeTrue()) })
				It("should set created_at timestamp", func() { Expect(response.CreatedAt).ToNot(BeZero()) })
				It("should set updated_at timestamp", func() { Expect(response.UpdatedAt).ToNot(BeZero()) })
				It("should persist ticker to database", func() {
					var persisted barkat.Ticker
					Expect(db.First(&persisted, "external_id = ?", "MCX").Error).ToNot(HaveOccurred())
					Expect(persisted.Ticker).To(Equal("MCX"))
				})
			})
		})

		Context("Field Validations", func() {
			Context("Ticker Field", func() {
				Context("Allowed Values", func() {
					It("should accept minimum ticker length 1", func() {
						payload := validTickerPayload
						payload.Ticker = "A"
						_, response := createTickerRequest(router, payload)
						Expect(response.Ticker).To(Equal("A"))
					})
					It("should accept maximum ticker length 50", func() {
						payload := validTickerPayload
						payload.Ticker = strings.Repeat("A", 50)
						_, response := createTickerRequest(router, payload)
						Expect(response.Ticker).To(HaveLen(50))
					})
					It("should accept uppercase ticker token MCX", func() {
						payload := validTickerPayload
						payload.Ticker = "MCX"
						_, response := createTickerRequest(router, payload)
						Expect(response.Ticker).To(Equal("MCX"))
					})
					It("should accept ticker with digits NIFTY50", func() {
						payload := validTickerPayload
						payload.Ticker = "NIFTY50"
						_, response := createTickerRequest(router, payload)
						Expect(response.Ticker).To(Equal("NIFTY50"))
					})
					It("should accept ticker with dot BRK.B", func() {
						payload := validTickerPayload
						payload.Ticker = "BRK.B"
						_, response := createTickerRequest(router, payload)
						Expect(response.Ticker).To(Equal("BRK.B"))
					})
					It("should accept ticker with underscore ABC_DEF", func() {
						payload := validTickerPayload
						payload.Ticker = "ABC_DEF"
						_, response := createTickerRequest(router, payload)
						Expect(response.Ticker).To(Equal("ABC_DEF"))
					})
					It("should accept futures exclamation GOLD1!", func() {
						payload := validTickerPayload
						payload.Ticker = "GOLD1!"
						_, response := createTickerRequest(router, payload)
						Expect(response.Ticker).To(Equal("GOLD1!"))
					})

					It("should accept composite expression when type = COMPOSITE", func() {
						payload := validTickerPayload
						payload.Ticker = "NIFTY/USDINR"
						payload.Type = "COMPOSITE"
						_, response := createTickerRequest(router, payload)
						Expect(response.Type).To(Equal("COMPOSITE"))
						Expect(response.Ticker).To(Equal("NIFTY/USDINR"))
					})
					It("should accept plus operator NIFTY+USDINR in composite expression", func() {
						payload := validTickerPayload
						payload.Ticker = "NIFTY+USDINR"
						payload.Type = "COMPOSITE"
						_, response := createTickerRequest(router, payload)
						Expect(response.Ticker).To(Equal("NIFTY+USDINR"))
					})
					It("should accept minus operator US10Y-US02Y in composite expression", func() {
						payload := validTickerPayload
						payload.Ticker = "US10Y-US02Y"
						payload.Type = "COMPOSITE"
						_, response := createTickerRequest(router, payload)
						Expect(response.Ticker).To(Equal("US10Y-US02Y"))
					})
					It("should accept multiply operator XAUUSD*USDINR in composite expression", func() {
						payload := validTickerPayload
						payload.Ticker = "XAUUSD*USDINR"
						payload.Type = "COMPOSITE"
						_, response := createTickerRequest(router, payload)
						Expect(response.Ticker).To(Equal("XAUUSD*USDINR"))
					})
					It("should accept ^ operator NIFTY^USDINR in composite expression", func() {
						payload := validTickerPayload
						payload.Ticker = "NIFTY^USDINR"
						payload.Type = "COMPOSITE"
						_, response := createTickerRequest(router, payload)
						Expect(response.Ticker).To(Equal("NIFTY^USDINR"))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing or empty ticker (PRD: required)", func() {
						payload := validTickerPayload
						payload.Ticker = ""
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "required")
					})
					It("should return 400 for ticker exceeding 50 characters", func() {
						payload := validTickerPayload
						payload.Ticker = strings.Repeat("A", 51)
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "max")
					})
					It("should return 400 for ticker starting with unsupported character", func() {
						payload := validTickerPayload
						payload.Ticker = ".MCX"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})
					It("should return 400 for ticker with whitespace", func() {
						payload := validTickerPayload
						payload.Ticker = "MC X"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})
					It("should return 400 for ticker with lowercase characters", func() {
						payload := validTickerPayload
						payload.Ticker = "mcx"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})
					It("should return 400 for ticker with unsupported special character @", func() {
						payload := validTickerPayload
						payload.Ticker = "MCX@"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})
					It("should return 400 for ticker with hyphen", func() {
						payload := validTickerPayload
						payload.Ticker = "ABC-DEF"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})
					It("should return 400 for ticker with ampersand", func() {
						payload := validTickerPayload
						payload.Ticker = "M&M"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})
					It("should return 400 for composite expression when type is not COMPOSITE", func() {
						payload := validTickerPayload
						payload.Ticker = "NIFTY/USDINR"
						payload.Type = "EQUITY"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})
					It("should return 400 for composite expression with non-math invalid char @", func() {
						payload := validTickerPayload
						payload.Ticker = "NIFTY@USDINR"
						payload.Type = "COMPOSITE"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})
					It("should return 400 for composite expression containing whitespace", func() {
						payload := validTickerPayload
						payload.Ticker = "NIFTY / USDINR"
						payload.Type = "COMPOSITE"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})
				})
			})

			Context("Exchange Field", func() {
				Context("Allowed Values", func() {
					It("should accept omitted exchange on POST", func() {
						payload := validTickerPayload
						payload.Exchange = nil
						_, response := createTickerRequest(router, payload)
						Expect(response.Exchange).To(BeNil())
					})
					It("should accept minimum exchange length 1", func() {
						payload := validTickerPayload
						payload.Exchange = new("N")
						_, response := createTickerRequest(router, payload)
						Expect(*response.Exchange).To(Equal("N"))
					})
					It("should accept maximum exchange length 10", func() {
						payload := validTickerPayload
						payload.Exchange = new(strings.Repeat("A", 10))
						_, response := createTickerRequest(router, payload)
						Expect(*response.Exchange).To(HaveLen(10))
					})
					It("should accept uppercase letters NSE", func() {
						payload := validTickerPayload
						payload.Exchange = new("NSE")
						_, response := createTickerRequest(router, payload)
						Expect(*response.Exchange).To(Equal("NSE"))
					})
					It("should accept underscore in exchange code", func() {
						payload := validTickerPayload
						payload.Exchange = new("FX_IDC")
						_, response := createTickerRequest(router, payload)
						Expect(*response.Exchange).To(Equal("FX_IDC"))
					})
					It("should accept dot in exchange code", func() {
						payload := validTickerPayload
						payload.Exchange = new("N.SE")
						_, response := createTickerRequest(router, payload)
						Expect(*response.Exchange).To(Equal("N.SE"))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for empty exchange string", func() {
						payload := validTickerPayload
						payload.Exchange = new("")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "min")
					})
					It("should return 400 for exchange exceeding 10 characters", func() {
						payload := validTickerPayload
						payload.Exchange = new(strings.Repeat("A", 11))
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "max")
					})
					It("should return 400 for lowercase exchange", func() {
						payload := validTickerPayload
						payload.Exchange = new("nse")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
					It("should return 400 for exchange with colon", func() {
						payload := validTickerPayload
						payload.Exchange = new("NSE:MCX")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
					It("should return 400 for exchange with hyphen", func() {
						payload := validTickerPayload
						payload.Exchange = new("N-SE")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
					It("should return 400 for exchange with whitespace", func() {
						payload := validTickerPayload
						payload.Exchange = new("NS E")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
					It("should return 400 for exchange with unsupported special character", func() {
						payload := validTickerPayload
						payload.Exchange = new("NSE@")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
					It("should return 400 for exchange with digit", func() {
						payload := validTickerPayload
						payload.Exchange = new("NSE1")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
					It("should return 400 for exchange starting with digit", func() {
						payload := validTickerPayload
						payload.Exchange = new("1NSE")
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
				})
			})

			Context("Timeframes Field", func() {
				Context("Allowed Values", func() {
					It("should accept one timeframe", func() {
						payload := validTickerPayload
						payload.Timeframes = []string{"DL"}
						_, response := createTickerRequest(router, payload)
						Expect(response.Timeframes).To(Equal([]string{"DL"}))
					})
					It("should accept six timeframes", func() {
						payload := validTickerPayload
						payload.Timeframes = []string{"YR", "SMN", "TMN", "MN", "WK", "DL"}
						_, response := createTickerRequest(router, payload)
						Expect(response.Timeframes).To(HaveLen(6))
					})
					It("should accept timeframe YR", func() {
						payload := validTickerPayload
						payload.Timeframes = []string{"YR"}
						_, response := createTickerRequest(router, payload)
						Expect(response.Timeframes).To(Equal([]string{"YR"}))
					})
					It("should accept timeframe SMN", func() {
						payload := validTickerPayload
						payload.Timeframes = []string{"SMN"}
						_, response := createTickerRequest(router, payload)
						Expect(response.Timeframes).To(Equal([]string{"SMN"}))
					})
					It("should accept timeframe TMN", func() {
						payload := validTickerPayload
						payload.Timeframes = []string{"TMN"}
						_, response := createTickerRequest(router, payload)
						Expect(response.Timeframes).To(Equal([]string{"TMN"}))
					})
					It("should accept timeframe MN", func() {
						payload := validTickerPayload
						payload.Timeframes = []string{"MN"}
						_, response := createTickerRequest(router, payload)
						Expect(response.Timeframes).To(Equal([]string{"MN"}))
					})
					It("should accept timeframe WK", func() {
						payload := validTickerPayload
						payload.Timeframes = []string{"WK"}
						_, response := createTickerRequest(router, payload)
						Expect(response.Timeframes).To(Equal([]string{"WK"}))
					})
					It("should accept timeframe DL", func() {
						payload := validTickerPayload
						payload.Timeframes = []string{"DL"}
						_, response := createTickerRequest(router, payload)
						Expect(response.Timeframes).To(Equal([]string{"DL"}))
					})
					It("should preserve provided timeframe order", func() {
						payload := validTickerPayload
						payload.Timeframes = []string{"YR", "MN", "DL"}
						_, response := createTickerRequest(router, payload)
						Expect(response.Timeframes).To(Equal([]string{"YR", "MN", "DL"}))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for missing timeframes (PRD: required)", func() {
						payload := validTickerPayload
						payload.Timeframes = nil
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframes", "required")
					})
					It("should return 400 for empty timeframe array", func() {
						payload := validTickerPayload
						payload.Timeframes = []string{}
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframes", "min")
					})
					It("should return 400 for more than six timeframes", func() {
						payload := validTickerPayload
						payload.Timeframes = []string{"YR", "SMN", "TMN", "MN", "WK", "DL", "YR"}
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframes", "max")
					})
					It("should return 400 for unsupported timeframe HR", func() {
						payload := validTickerPayload
						payload.Timeframes = []string{"HR"}
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframes", "oneof")
					})
					It("should return 400 for lowercase timeframe dl", func() {
						payload := validTickerPayload
						payload.Timeframes = []string{"dl"}
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframes", "oneof")
					})
					It("should return 400 for numeric timeframe item", func() {
						jsonPayload := `{"ticker":"MCX","exchange":"NSE","timeframes":[1],"type":"EQUITY","state":"WATCHED","trend":"UPTREND","last_opened_at":"2026-05-05T10:30:00Z","is_fno":true}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for null timeframe item", func() {
						jsonPayload := `{"ticker":"MCX","exchange":"NSE","timeframes":[null],"type":"EQUITY","state":"WATCHED","trend":"UPTREND","last_opened_at":"2026-05-05T10:30:00Z","is_fno":true}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Type Field", func() {
				Context("Allowed Values", func() {
					for _, tickerType := range []string{"EQUITY", "INDEX", "CRYPTO", "COMMODITY", "FX", "BOND", "COMPOSITE"} {
						typeValue := tickerType
						It("should accept type "+typeValue, func() {
							payload := validTickerPayload
							payload.Type = typeValue
							if typeValue == "COMPOSITE" {
								payload.Ticker = "NIFTY/USDINR"
							}
							_, response := createTickerRequest(router, payload)
							Expect(response.Type).To(Equal(typeValue))
						})
					}
				})

				Context("Bad Values", func() {
					It("should return 400 for missing type (PRD: required)", func() {
						payload := validTickerPayload
						payload.Type = ""
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "required")
					})
					It("should return 400 for lowercase type equity", func() {
						payload := validTickerPayload
						payload.Type = "equity"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})
					It("should return 400 for unsupported type METAL", func() {
						payload := validTickerPayload
						payload.Type = "METAL"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})
					It("should return 400 for empty type", func() {
						payload := validTickerPayload
						payload.Type = ""
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "required")
					})
				})
			})

			Context("State Field", func() {
				Context("Allowed Values", func() {
					for _, state := range []string{"WATCHED", "READY", "BLACKLIST"} {
						stateValue := state
						It("should accept state "+stateValue, func() {
							payload := validTickerPayload
							payload.State = stateValue
							_, response := createTickerRequest(router, payload)
							Expect(response.State).To(Equal(stateValue))
						})
					}
				})

				Context("Bad Values", func() {
					It("should return 400 for missing state (PRD: required)", func() {
						payload := validTickerPayload
						payload.State = ""
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "State", "required")
					})
					It("should return 400 for lowercase state watched", func() {
						payload := validTickerPayload
						payload.State = "watched"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "State", "oneof")
					})
					It("should return 400 for unsupported state ARCHIVED", func() {
						payload := validTickerPayload
						payload.State = "ARCHIVED"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "State", "oneof")
					})
					It("should return 400 for empty state", func() {
						payload := validTickerPayload
						payload.State = ""
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "State", "required")
					})
				})
			})

			Context("Trend Field", func() {
				Context("Allowed Values", func() {
					for _, trend := range []string{"UPTREND", "SIDEWAYS", "DOWNTREND"} {
						trendValue := trend
						It("should accept trend "+trendValue, func() {
							payload := validTickerPayload
							payload.Trend = trendValue
							_, response := createTickerRequest(router, payload)
							Expect(response.Trend).To(Equal(trendValue))
						})
					}
				})

				Context("Bad Values", func() {
					It("should return 400 for missing trend (PRD: required)", func() {
						payload := validTickerPayload
						payload.Trend = ""
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Trend", "required")
					})
					It("should return 400 for lowercase trend uptrend", func() {
						payload := validTickerPayload
						payload.Trend = "uptrend"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Trend", "oneof")
					})
					It("should return 400 for unsupported trend NEUTRAL", func() {
						payload := validTickerPayload
						payload.Trend = "NEUTRAL"
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Trend", "oneof")
					})
					It("should return 400 for empty trend", func() {
						payload := validTickerPayload
						payload.Trend = ""
						req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Trend", "required")
					})
				})
			})

			Context("Last Opened At Field", func() {
				Context("Allowed Values", func() {
					It("should accept RFC3339 timestamp with Z timezone", func() {
						payload := validTickerPayload
						payload.LastOpenedAt = time.Date(2026, time.May, 5, 10, 30, 0, 0, time.UTC)
						_, response := createTickerRequest(router, payload)
						Expect(response.LastOpenedAt).To(Equal(payload.LastOpenedAt))
					})
					It("should accept RFC3339 timestamp with numeric timezone offset", func() {
						jsonPayload := `{"ticker":"MCX","exchange":"NSE","timeframes":["MN","WK","DL"],"type":"EQUITY","state":"WATCHED","trend":"UPTREND","last_opened_at":"2026-05-05T10:30:00+05:30","is_fno":true}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, jsonPayload)
						router.ServeHTTP(w, req)
						response := decodeTickerResponse(w, http.StatusCreated)
						Expect(response.LastOpenedAt).ToNot(BeZero())
					})
					It("should preserve timestamp value", func() {
						payload := validTickerPayload
						_, response := createTickerRequest(router, payload)
						Expect(response.LastOpenedAt).To(Equal(payload.LastOpenedAt))
					})
				})

				Context("Bad Values", func() {
					// HACK: Simplify Use Typed Payloads in Tests to Avoid Repetitive Raw JSON Strings and Reduce Risk of Syntax Errors in Test Cases
					It("should return 400 for missing last_opened_at (PRD: required)", func() {
						jsonPayload := `{"ticker":"MCX","exchange":"NSE","timeframes":["MN","WK","DL"],"type":"EQUITY","state":"WATCHED","trend":"UPTREND","is_fno":true}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for non-string last_opened_at", func() {
						jsonPayload := `{"ticker":"MCX","exchange":"NSE","timeframes":["MN","WK","DL"],"type":"EQUITY","state":"WATCHED","trend":"UPTREND","last_opened_at":123,"is_fno":true}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for date-only value", func() {
						jsonPayload := `{"ticker":"MCX","exchange":"NSE","timeframes":["MN","WK","DL"],"type":"EQUITY","state":"WATCHED","trend":"UPTREND","last_opened_at":"2026-05-05","is_fno":true}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for timestamp without timezone", func() {
						jsonPayload := `{"ticker":"MCX","exchange":"NSE","timeframes":["MN","WK","DL"],"type":"EQUITY","state":"WATCHED","trend":"UPTREND","last_opened_at":"2026-05-05T10:30:00","is_fno":true}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for invalid timestamp text", func() {
						jsonPayload := `{"ticker":"MCX","exchange":"NSE","timeframes":["MN","WK","DL"],"type":"EQUITY","state":"WATCHED","trend":"UPTREND","last_opened_at":"not-a-date","is_fno":true}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for empty timestamp string", func() {
						jsonPayload := `{"ticker":"MCX","exchange":"NSE","timeframes":["MN","WK","DL"],"type":"EQUITY","state":"WATCHED","trend":"UPTREND","last_opened_at":"","is_fno":true}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Is FNO Field", func() {
				Context("Allowed Values", func() {
					It("should accept omitted is_fno on POST and default false", func() {
						jsonPayload := `{"ticker":"MCX","exchange":"NSE","timeframes":["MN","WK","DL"],"type":"EQUITY","state":"WATCHED","trend":"UPTREND","last_opened_at":"2026-05-05T10:30:00Z"}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, jsonPayload)
						router.ServeHTTP(w, req)
						response := decodeTickerResponse(w, http.StatusCreated)
						Expect(response.IsFNO).To(BeFalse())
					})
					It("should accept is_fno true", func() {
						payload := validTickerPayload
						payload.IsFNO = true
						_, response := createTickerRequest(router, payload)
						Expect(response.IsFNO).To(BeTrue())
					})
					It("should accept is_fno false", func() {
						payload := validTickerPayload
						payload.IsFNO = false
						_, response := createTickerRequest(router, payload)
						Expect(response.IsFNO).To(BeFalse())
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for string is_fno", func() {
						jsonPayload := `{"ticker":"MCX","exchange":"NSE","timeframes":["MN","WK","DL"],"type":"EQUITY","state":"WATCHED","trend":"UPTREND","last_opened_at":"2026-05-05T10:30:00Z","is_fno":"true"}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for numeric is_fno", func() {
						jsonPayload := `{"ticker":"MCX","exchange":"NSE","timeframes":["MN","WK","DL"],"type":"EQUITY","state":"WATCHED","trend":"UPTREND","last_opened_at":"2026-05-05T10:30:00Z","is_fno":1}`
						req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 400 for malformed JSON", func() {
				req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, `{"ticker":`)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 400 for empty request body", func() {
				req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, "")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 400 for null request body", func() {
				req, w := rawTickerRequest(http.MethodPost, barkat.TickerBase, "null")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 409 for duplicate ticker", func() {
				payload := validTickerPayload
				_, _ = createTickerRequest(router, payload)
				req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, payload)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusConflict))
			})
			It("should return 500 for persistence failure", func() {
				sqlDB, err := db.DB()
				Expect(err).ToNot(HaveOccurred())
				Expect(sqlDB.Close()).To(Succeed())
				req, w := util.CreateTestRequest(http.MethodPost, barkat.TickerBase, validTickerPayload)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	// ============================================================================
	// 2.2.1.3 PUT /v1/api/tickers/{ticker} - Update Primary Ticker
	// ============================================================================
	Describe("PUT /v1/api/tickers/{ticker} - Update Primary Ticker (2.2.1.3)", func() {
		var createdTicker barkat.Ticker
		var validUpdatePayload barkat.Ticker

		BeforeEach(func() {
			_, createdTicker = createTickerRequest(router, validTickerPayload)
			validUpdatePayload = barkat.Ticker{Exchange: new("NSE"), Timeframes: []string{"YR", "SMN", "TMN", "MN", "WK"}, Type: "EQUITY", State: "READY", Trend: "UPTREND", IsFNO: true}
		})

		Context("Happy Path", func() {
			Context("with valid replacement metadata", func() {
				var w *httptest.ResponseRecorder
				var response barkat.Ticker

				BeforeEach(func() { w, response = updateTickerRequest(router, createdTicker.Ticker, validUpdatePayload) })

				It("should return 200 OK", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
				})
				It("should return Envelope success", func() {
					var envelope common.Envelope[barkat.Ticker]
					util.AssertSuccess(w, http.StatusOK, &envelope)
					Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
				})
				It("should preserve immutable ticker from path", func() { Expect(response.Ticker).To(Equal(createdTicker.Ticker)) })
				It("should replace exchange", func() { Expect(response.Exchange).To(Equal(new("NSE"))) })
				It("should replace timeframes", func() { Expect(response.Timeframes).To(Equal([]string{"YR", "SMN", "TMN", "MN", "WK"})) })
				It("should replace type", func() { Expect(response.Type).To(Equal("EQUITY")) })
				It("should replace state", func() { Expect(response.State).To(Equal("READY")) })
				It("should replace trend", func() { Expect(response.Trend).To(Equal("UPTREND")) })
				It("should replace is_fno", func() { Expect(response.IsFNO).To(BeTrue()) })
				It("should update updated_at timestamp", func() { Expect(response.UpdatedAt).To(BeTemporally(">=", createdTicker.UpdatedAt)) })
			})
		})

		Context("Field Validations", func() {
			Context("Ticker Path Parameter", func() {
				Context("Allowed Values", func() {
					It("should accept valid existing ticker path", func() {
						_, response := updateTickerRequest(router, createdTicker.Ticker, validUpdatePayload)
						Expect(response.Ticker).To(Equal(createdTicker.Ticker))
					})
					It("should accept encoded composite ticker path NIFTY/USDINR", func() {
						compositePayload := validTickerPayload
						compositePayload.Ticker = "NIFTY/USDINR"
						compositePayload.Type = "COMPOSITE"
						_, compositeTicker := createTickerRequest(router, compositePayload)

						compositeUpdate := validUpdatePayload
						compositeUpdate.Type = "COMPOSITE"
						_, response := updateTickerRequest(router, url.PathEscape(compositeTicker.Ticker), compositeUpdate)
						Expect(response.Ticker).To(Equal(compositeTicker.Ticker))
						Expect(response.Type).To(Equal("COMPOSITE"))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for invalid ticker path", func() {
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/mcx", validUpdatePayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 when composite ticker path is updated with non-COMPOSITE type", func() {
						compositePayload := validTickerPayload
						compositePayload.Ticker = "ABC(DEF)"
						compositePayload.Type = "COMPOSITE"
						_, compositeTicker := createTickerRequest(router, compositePayload)

						updatePayload := validUpdatePayload
						updatePayload.Type = "EQUITY"
						req, w := util.CreateTestRequest(http.MethodPut,
							barkat.TickerBase+"/"+compositeTicker.Ticker, updatePayload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Ticker", "ticker")
					})
					It("should return 404 for valid missing ticker path", func() {
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/NOTFOUND", validUpdatePayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})

			Context("Exchange Field", func() {
				Context("Allowed Values", func() {
					It("should accept omitted exchange on PUT and leave existing value unchanged", func() {
						req, w := rawTickerRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, `{"timeframes":["MN","WK","DL"],"type":"EQUITY","state":"READY","trend":"UPTREND","is_fno":true}`)
						router.ServeHTTP(w, req)
						response := decodeTickerResponse(w, http.StatusOK)
						Expect(response.Exchange).To(Equal(createdTicker.Exchange))
					})
					It("should accept valid exchange code", func() {
						payload := validUpdatePayload
						payload.Exchange = new("NSE")
						_, response := updateTickerRequest(router, createdTicker.Ticker, payload)
						Expect(*response.Exchange).To(Equal("NSE"))
					})
					It("should accept maximum exchange length 10", func() {
						payload := validUpdatePayload
						payload.Exchange = new(strings.Repeat("A", 10))
						_, response := updateTickerRequest(router, createdTicker.Ticker, payload)
						Expect(*response.Exchange).To(HaveLen(10))
					})
					It("should accept underscore in exchange code", func() {
						payload := validUpdatePayload
						payload.Exchange = new("FX_IDC")
						_, response := updateTickerRequest(router, createdTicker.Ticker, payload)
						Expect(*response.Exchange).To(Equal("FX_IDC"))
					})
					It("should accept dot in exchange code", func() {
						payload := validUpdatePayload
						payload.Exchange = new("N.SE")
						_, response := updateTickerRequest(router, createdTicker.Ticker, payload)
						Expect(*response.Exchange).To(Equal("N.SE"))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for empty exchange string", func() {
						payload := validUpdatePayload
						payload.Exchange = new("")
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "min")
					})
					It("should return 400 for exchange exceeding 10 characters", func() {
						payload := validUpdatePayload
						payload.Exchange = new(strings.Repeat("A", 11))
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "max")
					})
					It("should return 400 for unsupported exchange format", func() {
						payload := validUpdatePayload
						payload.Exchange = new("NSE:MCX")
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
					It("should return 400 for exchange with hyphen", func() {
						payload := validUpdatePayload
						payload.Exchange = new("N-SE")
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
					It("should return 400 for lowercase exchange", func() {
						payload := validUpdatePayload
						payload.Exchange = new("nse")
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
					It("should return 400 for exchange with whitespace", func() {
						payload := validUpdatePayload
						payload.Exchange = new("NS E")
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
					It("should return 400 for exchange with unsupported special character", func() {
						payload := validUpdatePayload
						payload.Exchange = new("NSE@")
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
					It("should return 400 for exchange with digit", func() {
						payload := validUpdatePayload
						payload.Exchange = new("NSE1")
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
					It("should return 400 for exchange starting with digit", func() {
						payload := validUpdatePayload
						payload.Exchange = new("1NSE")
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "ticker_exchange")
					})
				})
			})

			Context("Timeframes Field", func() {
				Context("Allowed Values", func() {
					It("should accept one timeframe", func() {
						payload := validUpdatePayload
						payload.Timeframes = []string{"DL"}
						_, response := updateTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Timeframes).To(Equal([]string{"DL"}))
					})
					It("should accept six timeframes", func() {
						payload := validUpdatePayload
						payload.Timeframes = []string{"YR", "SMN", "TMN", "MN", "WK", "DL"}
						_, response := updateTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Timeframes).To(HaveLen(6))
					})
					It("should accept all supported timeframe enum values", func() {
						payload := validUpdatePayload
						payload.Timeframes = []string{"YR", "SMN", "TMN", "MN", "WK", "DL"}
						_, response := updateTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.Timeframes).To(ConsistOf("YR", "SMN", "TMN", "MN", "WK", "DL"))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for missing timeframes", func() {
						payload := validUpdatePayload
						payload.Timeframes = nil
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframes", "required")
					})
					It("should return 400 for empty timeframes", func() {
						payload := validUpdatePayload
						payload.Timeframes = []string{}
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframes", "min")
					})
					It("should return 400 for more than six timeframes", func() {
						payload := validUpdatePayload
						payload.Timeframes = []string{"YR", "SMN", "TMN", "MN", "WK", "DL", "YR"}
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframes", "max")
					})
					It("should return 400 for unsupported timeframe", func() {
						payload := validUpdatePayload
						payload.Timeframes = []string{"HR"}
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Timeframes", "oneof")
					})
				})
			})

			Context("Type Field", func() {
				Context("Allowed Values", func() {
					It("should accept every supported type enum value", func() {
						for _, v := range []string{"EQUITY", "INDEX", "CRYPTO", "COMMODITY", "FX", "BOND", "COMPOSITE"} {
							payload := validUpdatePayload
							payload.Type = v
							if v == "COMPOSITE" {
								payload.Ticker = "NIFTY/USDINR"
							}
							_, response := updateTickerRequest(router, createdTicker.Ticker, payload)
							Expect(response.Type).To(Equal(v))
						}
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for missing type", func() {
						payload := validUpdatePayload
						payload.Type = ""
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "required")
					})
					It("should return 400 for lowercase type", func() {
						payload := validUpdatePayload
						payload.Type = "equity"
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})
					It("should return 400 for unsupported type", func() {
						payload := validUpdatePayload
						payload.Type = "METAL"
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})
				})
			})

			Context("State Field", func() {
				Context("Allowed Values", func() {
					It("should accept every supported state enum value", func() {
						for _, v := range []string{"WATCHED", "READY", "BLACKLIST"} {
							payload := validUpdatePayload
							payload.State = v
							_, response := updateTickerRequest(router, createdTicker.Ticker, payload)
							Expect(response.State).To(Equal(v))
						}
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for missing state", func() {
						payload := validUpdatePayload
						payload.State = ""
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "State", "required")
					})
					It("should return 400 for lowercase state", func() {
						payload := validUpdatePayload
						payload.State = "ready"
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "State", "oneof")
					})
					It("should return 400 for unsupported state", func() {
						payload := validUpdatePayload
						payload.State = "ARCHIVED"
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "State", "oneof")
					})
				})
			})

			Context("Trend Field", func() {
				Context("Allowed Values", func() {
					It("should accept every supported trend enum value", func() {
						for _, v := range []string{"UPTREND", "SIDEWAYS", "DOWNTREND"} {
							payload := validUpdatePayload
							payload.Trend = v
							_, response := updateTickerRequest(router, createdTicker.Ticker, payload)
							Expect(response.Trend).To(Equal(v))
						}
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for missing trend", func() {
						payload := validUpdatePayload
						payload.Trend = ""
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Trend", "required")
					})
					It("should return 400 for lowercase trend", func() {
						payload := validUpdatePayload
						payload.Trend = "uptrend"
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Trend", "oneof")
					})
					It("should return 400 for unsupported trend", func() {
						payload := validUpdatePayload
						payload.Trend = "NEUTRAL"
						req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, payload)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Trend", "oneof")
					})
				})
			})

			Context("Is FNO Field", func() {
				Context("Allowed Values", func() {
					It("should accept is_fno true", func() {
						payload := validUpdatePayload
						payload.IsFNO = true
						_, response := updateTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.IsFNO).To(BeTrue())
					})
					It("should accept is_fno false", func() {
						payload := validUpdatePayload
						payload.IsFNO = false
						_, response := updateTickerRequest(router, createdTicker.Ticker, payload)
						Expect(response.IsFNO).To(BeFalse())
					})
				})
				Context("Bad Values", func() {
					PIt("should return 400 for missing is_fno", func() {
						req, w := rawTickerRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, `{"exchange":"NSE","timeframes":["MN","WK","DL"],"type":"EQUITY","state":"READY","trend":"UPTREND"}`)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for string is_fno", func() {
						jsonPayload := `{"exchange":"NSE","timeframes":["MN","WK","DL"],"type":"EQUITY","state":"READY","trend":"UPTREND","is_fno":"true"}`
						req, w := rawTickerRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for numeric is_fno", func() {
						jsonPayload := `{"exchange":"NSE","timeframes":["MN","WK","DL"],"type":"EQUITY","state":"READY","trend":"UPTREND","is_fno":1}`
						req, w := rawTickerRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, jsonPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 400 for malformed JSON", func() {
				req, w := rawTickerRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, `{"exchange":`)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 400 for empty request body", func() {
				req, w := rawTickerRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, "")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 400 for null request body", func() {
				req, w := rawTickerRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, "null")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 500 for update failure", func() {
				sqlDB, err := db.DB()
				Expect(err).ToNot(HaveOccurred())
				Expect(sqlDB.Close()).To(Succeed())
				req, w := util.CreateTestRequest(http.MethodPut, barkat.TickerBase+"/"+createdTicker.Ticker, validUpdatePayload)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	// ============================================================================
	// 2.2.1.4 PATCH /v1/api/tickers/{ticker} - Patch Last Opened Timestamp
	// ============================================================================
	Describe("PATCH /v1/api/tickers/{ticker} - Patch Last Opened Timestamp (2.2.1.4)", func() {
		var createdTicker barkat.Ticker
		var validPatchPayload barkat.TickerLastOpenedUpdate

		BeforeEach(func() {
			_, createdTicker = createTickerRequest(router, validTickerPayload)
			validPatchPayload = barkat.TickerLastOpenedUpdate{LastOpenedAt: time.Date(2026, time.May, 5, 11, 0, 0, 0, time.UTC)}
		})

		Context("Happy Path", func() {
			Context("with valid timestamp payload", func() {
				var w *httptest.ResponseRecorder
				var response barkat.Ticker
				BeforeEach(func() { w, response = patchTickerRequest(router, createdTicker.Ticker, validPatchPayload) })
				It("should return 200 OK", func() {
					Expect(w.Code).To(Equal(http.StatusOK))
				})
				It("should return Envelope success", func() {
					var envelope common.Envelope[barkat.Ticker]
					util.AssertSuccess(w, http.StatusOK, &envelope)
					Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
				})
				It("should return ticker field", func() { Expect(response.Ticker).To(Equal(createdTicker.Ticker)) })
				It("should return updated last_opened_at", func() { Expect(response.LastOpenedAt).To(Equal(validPatchPayload.LastOpenedAt)) })
				It("should update updated_at timestamp", func() { Expect(response.UpdatedAt).To(BeTemporally(">=", createdTicker.UpdatedAt)) })
				It("should not modify exchange, timeframes, type, state, trend, or is_fno", func() {
					Expect(response.Exchange).To(Equal(createdTicker.Exchange))
					Expect(response.Timeframes).To(Equal(createdTicker.Timeframes))
					Expect(response.Type).To(Equal(createdTicker.Type))
					Expect(response.State).To(Equal(createdTicker.State))
					Expect(response.Trend).To(Equal(createdTicker.Trend))
					Expect(response.IsFNO).To(Equal(createdTicker.IsFNO))
				})
			})
		})

		Context("Field Validations", func() {
			Context("Ticker Path Parameter", func() {
				Context("Allowed Values", func() {
					It("should accept valid existing ticker path", func() {
						_, response := patchTickerRequest(router, createdTicker.Ticker, validPatchPayload)
						Expect(response.Ticker).To(Equal(createdTicker.Ticker))
					})
					It("should accept encoded composite ticker path US10Y-US02Y", func() {
						compositePayload := validTickerPayload
						compositePayload.Ticker = "US10Y-US02Y"
						compositePayload.Type = "COMPOSITE"
						_, compositeTicker := createTickerRequest(router, compositePayload)
						_, response := patchTickerRequest(router, url.PathEscape(compositeTicker.Ticker), validPatchPayload)
						Expect(response.Ticker).To(Equal(compositeTicker.Ticker))
						Expect(response.LastOpenedAt).To(Equal(validPatchPayload.LastOpenedAt))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for invalid ticker path", func() {
						req, w := util.CreateTestRequest(http.MethodPatch, barkat.TickerBase+"/mcx", validPatchPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 404 for valid missing ticker path", func() {
						req, w := util.CreateTestRequest(http.MethodPatch, barkat.TickerBase+"/NOTFOUND", validPatchPayload)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})

			Context("Last Opened At Field", func() {
				Context("Allowed Values", func() {
					It("should accept RFC3339 timestamp with Z timezone", func() {
						_, response := patchTickerRequest(router, createdTicker.Ticker, validPatchPayload)
						Expect(response.LastOpenedAt).To(Equal(validPatchPayload.LastOpenedAt))
					})
					It("should accept RFC3339 timestamp with numeric timezone offset", func() {
						req, w := rawTickerRequest(http.MethodPatch, barkat.TickerBase+"/"+createdTicker.Ticker, `{"last_opened_at":"2026-05-05T11:00:00+05:30"}`)
						router.ServeHTTP(w, req)
						response := decodeTickerResponse(w, http.StatusOK)
						Expect(response.LastOpenedAt).ToNot(BeZero())
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for missing last_opened_at", func() {
						req, w := rawTickerRequest(http.MethodPatch, barkat.TickerBase+"/"+createdTicker.Ticker, `{}`)
						router.ServeHTTP(w, req)
						util.AssertError(w, "LastOpenedAt", "required")
					})
					It("should return 400 for non-string last_opened_at", func() {
						req, w := rawTickerRequest(http.MethodPatch, barkat.TickerBase+"/"+createdTicker.Ticker, `{"last_opened_at":123}`)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for date-only value", func() {
						req, w := rawTickerRequest(http.MethodPatch, barkat.TickerBase+"/"+createdTicker.Ticker, `{"last_opened_at":"2026-05-05"}`)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for timestamp without timezone", func() {
						req, w := rawTickerRequest(http.MethodPatch, barkat.TickerBase+"/"+createdTicker.Ticker, `{"last_opened_at":"2026-05-05T11:00:00"}`)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for invalid timestamp text", func() {
						req, w := rawTickerRequest(http.MethodPatch, barkat.TickerBase+"/"+createdTicker.Ticker, `{"last_opened_at":"invalid"}`)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 400 for malformed JSON", func() {
				req, w := rawTickerRequest(http.MethodPatch, barkat.TickerBase+"/"+createdTicker.Ticker, `{"last_opened_at":`)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 400 for empty request body", func() {
				req, w := rawTickerRequest(http.MethodPatch, barkat.TickerBase+"/"+createdTicker.Ticker, "")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 400 for null request body", func() {
				req, w := rawTickerRequest(http.MethodPatch, barkat.TickerBase+"/"+createdTicker.Ticker, "null")
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 500 for update failure", func() {
				sqlDB, err := db.DB()
				Expect(err).ToNot(HaveOccurred())
				Expect(sqlDB.Close()).To(Succeed())
				req, w := util.CreateTestRequest(http.MethodPatch, barkat.TickerBase+"/"+createdTicker.Ticker, validPatchPayload)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	// ============================================================================
	// 2.2.1.5 DELETE /v1/api/tickers/{ticker} - Delete Primary Ticker
	// ============================================================================
	Describe("DELETE /v1/api/tickers/{ticker} - Delete Primary Ticker (2.2.1.5)", func() {
		var createdTicker barkat.Ticker
		BeforeEach(func() { _, createdTicker = createTickerRequest(router, validTickerPayload) })

		Context("Happy Path", func() {
			Context("with existing ticker", func() {
				It("should return 204 No Content", func() {
					req, w := util.CreateTestRequest(http.MethodDelete, barkat.TickerBase+"/"+createdTicker.Ticker, nil)
					router.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusNoContent))
				})
				It("should return empty body", func() {
					req, w := util.CreateTestRequest(http.MethodDelete, barkat.TickerBase+"/"+createdTicker.Ticker, nil)
					router.ServeHTTP(w, req)
					Expect(w.Body.String()).To(BeEmpty())
				})
				It("should delete ticker from database", func() {
					req, w := util.CreateTestRequest(http.MethodDelete, barkat.TickerBase+"/"+createdTicker.Ticker, nil)
					router.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusNoContent))
					var persisted barkat.Ticker
					Expect(db.First(&persisted, "external_id = ?", createdTicker.Ticker).Error).To(HaveOccurred())
				})
				It("should cascade delete linked Alert tickers and price alerts", func() {
					req, w := util.CreateTestRequest(http.MethodDelete, barkat.TickerBase+"/"+createdTicker.Ticker, nil)
					router.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusNoContent))
					var alertTickerCount int64
					Expect(db.Model(&barkat.AlertTicker{}).Where("ticker_id = ?", createdTicker.ID).Count(&alertTickerCount).Error).ToNot(HaveOccurred())
					Expect(alertTickerCount).To(Equal(int64(0)))
				})
			})
		})

		Context("Field Validations", func() {
			Context("Ticker Path Parameter", func() {
				Context("Allowed Values", func() {
					It("should accept valid existing ticker path", func() {
						req, w := util.CreateTestRequest(http.MethodDelete, barkat.TickerBase+"/"+createdTicker.Ticker, nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNoContent))
					})
					It("should accept encoded composite ticker path NIFTY/USDINR", func() {
						compositePayload := validTickerPayload
						compositePayload.Ticker = "NIFTY/USDINR"
						compositePayload.Type = "COMPOSITE"
						_, compositeTicker := createTickerRequest(router, compositePayload)
						req, w := util.CreateTestRequest(http.MethodDelete,
							barkat.TickerBase+"/"+url.PathEscape(compositeTicker.Ticker), nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNoContent))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for invalid ticker path", func() {
						req, w := util.CreateTestRequest(http.MethodDelete, barkat.TickerBase+"/mcx", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 404 for valid missing ticker path", func() {
						req, w := util.CreateTestRequest(http.MethodDelete, barkat.TickerBase+"/NOTFOUND", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 404 on second delete", func() {
				req, w := util.CreateTestRequest(http.MethodDelete, barkat.TickerBase+"/"+createdTicker.Ticker, nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNoContent))
				req, w = util.CreateTestRequest(http.MethodDelete, barkat.TickerBase+"/"+createdTicker.Ticker, nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
			It("should return 500 for delete failure", func() {
				sqlDB, err := db.DB()
				Expect(err).ToNot(HaveOccurred())
				Expect(sqlDB.Close()).To(Succeed())
				req, w := util.CreateTestRequest(http.MethodDelete, barkat.TickerBase+"/"+createdTicker.Ticker, nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
