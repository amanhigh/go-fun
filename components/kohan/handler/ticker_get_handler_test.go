//nolint:dupl
package handler_test

import (
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

func decodeTickerGetResponse(w *httptest.ResponseRecorder) barkat.Ticker {
	var envelope common.Envelope[barkat.Ticker]
	util.AssertSuccess(w, http.StatusOK, &envelope)
	return envelope.Data
}

func decodeTickerListResponse(w *httptest.ResponseRecorder) barkat.TickerList {
	var envelope common.Envelope[barkat.TickerList]
	util.AssertSuccess(w, http.StatusOK, &envelope)
	return envelope.Data
}

func seedTicker(ctx context.Context, db *gorm.DB, ticker barkat.Ticker) barkat.Ticker {
	Expect(db.WithContext(ctx).Create(&ticker).Error).ToNot(HaveOccurred())
	return ticker
}

// TickerHandler Integration GET/List Tests - Comprehensive Master Specification
// Tests complete HTTP → Handler → Manager → Repository → Database flow for PRD Section 2.2.1.2 and 2.2.1.6.
var _ = Describe("TickerHandler Integration - GET/List Tests - Section 2.2.1 Primary Ticker APIs", func() {
	const (
		tickerTypeComposite = "COMPOSITE"
		tickerNiftyUSDINR   = "NIFTY/USDINR"
		exchangeNSE         = "NSE"
		tickerUS10YUS02Y    = "US10Y-US02Y"
	)
	var (
		tickerHandler      handler.TickerHandler
		router             *gin.Engine
		testCtx            = context.Background()
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
			Exchange:     exchangeNSE,
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
	// 2.2.1.2 GET /v1/api/tickers/{ticker} - Retrieve Primary Ticker
	// ============================================================================
	Describe("GET /v1/api/tickers/{ticker} - Retrieve Primary Ticker (2.2.1.2)", func() {
		var createdTicker barkat.Ticker

		BeforeEach(func() {
			createdTicker = seedTicker(testCtx, db, validTickerPayload)
			alertExchange := exchangeNSE
			Expect(db.Create(&barkat.AlertTicker{
				TickerID: createdTicker.ID, Symbol: "MCIX", PairID: "941982",
				Name: "Multi Commodity Exchange of India", Exchange: &alertExchange, Type: barkat.AlertTickerTypePrimary,
			}).Error).ToNot(HaveOccurred())
		})

		Context("Happy Path", func() {
			Context("with existing ticker", func() {
				var response barkat.Ticker

				BeforeEach(func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"/"+createdTicker.Ticker, nil)
					router.ServeHTTP(w, req)
					response = decodeTickerGetResponse(w)
				})

				It("should return 200 OK", func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"/"+createdTicker.Ticker, nil)
					router.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusOK))
				})

				It("should return Envelope success", func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"/"+createdTicker.Ticker, nil)
					router.ServeHTTP(w, req)
					var envelope common.Envelope[barkat.Ticker]
					util.AssertSuccess(w, http.StatusOK, &envelope)
					Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
				})

				It("should return ticker with correct fields", func() {
					Expect(response.Ticker).To(Equal("MCX"))
					Expect(response.Exchange).To(Equal(exchangeNSE))
					Expect(response.Timeframes).To(Equal([]string{"MN", "WK", "DL"}))
					Expect(response.Type).To(Equal("EQUITY"))
					Expect(response.State).To(Equal("WATCHED"))
					Expect(response.Trend).To(Equal("UPTREND"))
					Expect(response.IsFNO).To(BeTrue())
				})

				It("should include created_at and updated_at", func() {
					Expect(response.CreatedAt).ToNot(BeZero())
					Expect(response.UpdatedAt).ToNot(BeZero())
				})

				It("should include mapped alert_tickers array", func() {
					Expect(response.AlertTickers).To(HaveLen(1))
					Expect(response.AlertTickers[0].Symbol).To(Equal("MCIX"))
					Expect(response.AlertTickers[0].PairID).To(Equal("941982"))
					Expect(response.AlertTickers[0].Name).To(Equal("Multi Commodity Exchange of India"))
					Expect(response.AlertTickers[0].Type).To(Equal(barkat.AlertTickerTypePrimary))
				})
			})
		})

		Context("Field Validations", func() {
			Context("Ticker Path Parameter", func() {
				Context("Allowed Values", func() {
					It("should accept simple ticker path MCX", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"/MCX", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})

					It("should accept futures ticker path GOLD1!", func() {
						payload := validTickerPayload
						payload.Ticker = "GOLD1!"
						seedTicker(testCtx, db, payload)
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"/GOLD1!", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})

					It("should accept ticker path with underscore", func() {
						payload := validTickerPayload
						payload.Ticker = "ABC_DEF"
						seedTicker(testCtx, db, payload)
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"/ABC_DEF", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})

					It("should accept encoded composite slash ticker NIFTY/USDINR", func() {
						payload := validTickerPayload
						payload.Ticker = tickerNiftyUSDINR
						payload.Type = tickerTypeComposite
						seedTicker(testCtx, db, payload)
						req, w := util.CreateTestRequest(http.MethodGet,
							barkat.TickerBase+"/"+url.PathEscape(tickerNiftyUSDINR), nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
						response := decodeTickerGetResponse(w)
						Expect(response.Ticker).To(Equal(tickerNiftyUSDINR))
						Expect(response.Type).To(Equal(tickerTypeComposite))
					})

					It("should accept composite minus ticker US10Y-US02Y", func() {
						payload := validTickerPayload
						payload.Ticker = tickerUS10YUS02Y
						payload.Type = tickerTypeComposite
						seedTicker(testCtx, db, payload)
						req, w := util.CreateTestRequest(http.MethodGet,
							barkat.TickerBase+"/"+url.PathEscape(tickerUS10YUS02Y), nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
						response := decodeTickerGetResponse(w)
						Expect(response.Ticker).To(Equal(tickerUS10YUS02Y))
						Expect(response.Type).To(Equal(tickerTypeComposite))
					})
				})

				Context("Bad Values", func() {
					It("should return 400 for lowercase ticker path", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"/mcx", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for ticker path with whitespace", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"/MC%20X", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for ticker path with unsupported special character", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"/MCX@", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for ticker path with ampersand", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"/M%26M", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 404 for valid ticker format but missing ticker", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"/NOTFOUND", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 500 when repository load fails", func() {
				sqlDB, err := db.DB()
				Expect(err).ToNot(HaveOccurred())
				Expect(sqlDB.Close()).To(Succeed())
				req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"/"+createdTicker.Ticker, nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	// ============================================================================
	// 2.2.1.6 GET /v1/api/tickers - List Primary Tickers
	// ============================================================================
	Describe("GET /v1/api/tickers - List Primary Tickers (2.2.1.6)", func() {
		var (
			mcxTicker   barkat.Ticker
			btcTicker   barkat.Ticker
			niftyTicker barkat.Ticker
		)

		BeforeEach(func() {
			mcxTicker = barkat.Ticker{
				Ticker:       "MCX",
				Exchange:     exchangeNSE,
				Timeframes:   []string{"MN", "WK", "DL"},
				Type:         "EQUITY",
				State:        "WATCHED",
				Trend:        "UPTREND",
				LastOpenedAt: time.Date(2026, time.May, 5, 10, 30, 0, 0, time.UTC),
				IsFNO:        true,
			}
			mcxTicker = seedTicker(testCtx, db, mcxTicker)
			btcTicker = barkat.Ticker{
				Ticker:       "BTCUSD",
				Exchange:     "BINANCE",
				Timeframes:   []string{"DL"},
				Type:         "CRYPTO",
				State:        "READY",
				Trend:        "SIDEWAYS",
				LastOpenedAt: time.Date(2026, time.May, 6, 10, 30, 0, 0, time.UTC),
				IsFNO:        false,
			}
			btcTicker = seedTicker(testCtx, db, btcTicker)
			niftyTicker = barkat.Ticker{
				Ticker:       tickerNiftyUSDINR,
				Exchange:     exchangeNSE,
				Timeframes:   []string{"YR", "MN"},
				Type:         tickerTypeComposite,
				State:        "BLACKLIST",
				Trend:        "DOWNTREND",
				LastOpenedAt: time.Date(2026, time.May, 7, 10, 30, 0, 0, time.UTC),
				IsFNO:        false,
			}
			niftyTicker = seedTicker(testCtx, db, niftyTicker)

			// Seed alert tickers with distinct counts: MCX=2, BTCUSD=1, NIFTY/USDINR=0
			alertExchange := exchangeNSE
			Expect(db.Create(&barkat.AlertTicker{TickerID: mcxTicker.ID, Symbol: "MCIX1", PairID: "941982", Name: "MC Alert 1", Exchange: &alertExchange}).Error).ToNot(HaveOccurred())
			Expect(db.Create(&barkat.AlertTicker{TickerID: mcxTicker.ID, Symbol: "MCIX2", PairID: "941983", Name: "MC Alert 2", Exchange: &alertExchange}).Error).ToNot(HaveOccurred())
			Expect(db.Create(&barkat.AlertTicker{TickerID: btcTicker.ID, Symbol: "BTCALERT", PairID: "941984", Name: "BTC Alert", Exchange: &alertExchange}).Error).ToNot(HaveOccurred())
		})

		Context("Happy Path", func() {
			Context("default pagination", func() {
				var response barkat.TickerList

				BeforeEach(func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase, nil)
					router.ServeHTTP(w, req)
					response = decodeTickerListResponse(w)
				})

				It("should return 200 OK", func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase, nil)
					router.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusOK))
				})
				It("should return Envelope success", func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase, nil)
					router.ServeHTTP(w, req)
					var envelope common.Envelope[barkat.TickerList]
					util.AssertSuccess(w, http.StatusOK, &envelope)
					Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
				})
				It("should return tickers array", func() { Expect(response.Tickers).To(HaveLen(3)) })
				It("should return metadata offset 0", func() { Expect(response.Metadata.Offset).To(Equal(0)) })
				It("should return metadata limit 20", func() { Expect(response.Metadata.Limit).To(Equal(20)) })
				It("should return metadata total", func() { Expect(response.Metadata.Total).To(Equal(int64(3))) })
			})

			Context("response shape", func() {
				var ticker barkat.Ticker
				var tickersBySymbol map[string]barkat.Ticker

				BeforeEach(func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase, nil)
					router.ServeHTTP(w, req)
					response := decodeTickerListResponse(w)
					ticker = response.Tickers[0]
					tickersBySymbol = make(map[string]barkat.Ticker, len(response.Tickers))
					for _, t := range response.Tickers {
						tickersBySymbol[t.Ticker] = t
					}
				})

				It("should include ticker", func() { Expect(ticker.Ticker).ToNot(BeEmpty()) })
				It("should include exchange", func() { Expect(ticker.Exchange).ToNot(BeEmpty()) })
				It("should include timeframes", func() { Expect(ticker.Timeframes).ToNot(BeEmpty()) })
				It("should include type", func() { Expect(ticker.Type).ToNot(BeEmpty()) })
				It("should include state", func() { Expect(ticker.State).ToNot(BeEmpty()) })
				It("should include trend", func() { Expect(ticker.Trend).ToNot(BeEmpty()) })
				It("should include last_opened_at", func() { Expect(ticker.LastOpenedAt).ToNot(BeZero()) })
				It("should include is_fno", func() { Expect(ticker.IsFNO).To(BeAssignableToTypeOf(false)) })
				It("should include alert_ticker_count", func() {
					Expect(tickersBySymbol["MCX"].AlertTickerCount).To(Equal(int64(2)))
					Expect(tickersBySymbol["BTCUSD"].AlertTickerCount).To(Equal(int64(1)))
					Expect(tickersBySymbol[tickerNiftyUSDINR].AlertTickerCount).To(Equal(int64(0)))
				})
			})
		})

		Context("Field Validations", func() {
			Context("Search Query Parameter", func() {
				Context("Allowed Values", func() {
					It("should filter by case-insensitive ticker substring", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?search=mc", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						Expect(response.Tickers).To(HaveLen(1))
						Expect(response.Tickers[0].Ticker).To(Equal("MCX"))
					})
					It("should return empty list for no match", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?search=ZZZ", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						Expect(response.Tickers).To(BeEmpty())
						Expect(response.Metadata.Total).To(Equal(int64(0)))
					})
					It("should return empty results for search with no matching chars", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?search=MCX@", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
				})
			})

			Context("Exchange Query Parameter", func() {
				Context("Allowed Values", func() {
					It("should filter by exact exchange", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?exchange=NSE", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						Expect(response.Tickers).To(HaveLen(2))
						Expect(response.Tickers[0].Exchange).To(Equal(exchangeNSE))
					})
					It("should accept uppercase letters NSE", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?exchange=NSE", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should filter by BINANCE exchange", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?exchange=BINANCE", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						Expect(response.Tickers).To(HaveLen(1))
						Expect(response.Tickers[0].Exchange).To(Equal("BINANCE"))
					})
					It("should accept NASDAQ exchange", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?exchange=NASDAQ", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should treat empty exchange query as no filter", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?exchange=", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for exchange exceeding 10 characters", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?exchange="+strings.Repeat("A", 11), nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for lowercase exchange", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?exchange=nse", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "oneof")
					})
					It("should return 400 for European exchange LSE", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?exchange=LSE", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Exchange", "oneof")
					})
				})
			})

			Context("Type Query Parameter", func() {
				Context("Allowed Values", func() {
					for _, typeValue := range []string{"EQUITY", "INDEX", "CRYPTO", "COMMODITY", "FX", "BOND", tickerTypeComposite} {
						value := typeValue
						It("should filter by "+value, func() {
							req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?type="+value, nil)
							router.ServeHTTP(w, req)
							response := decodeTickerListResponse(w)
							for _, ticker := range response.Tickers {
								Expect(ticker.Type).To(Equal(value))
							}
						})
					}
				})
				Context("Bad Values", func() {
					It("should return 400 for lowercase type query", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?type=equity", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})
					It("should return 400 for unsupported type query", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?type=METAL", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Type", "oneof")
					})
				})
			})

			Context("State Query Parameter", func() {
				Context("Allowed Values", func() {
					for _, stateValue := range []string{"WATCHED", "READY", "BLACKLIST"} {
						value := stateValue
						It("should filter by "+value, func() {
							req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?state="+value, nil)
							router.ServeHTTP(w, req)
							response := decodeTickerListResponse(w)
							for _, ticker := range response.Tickers {
								Expect(ticker.State).To(Equal(value))
							}
						})
					}
				})
				Context("Bad Values", func() {
					It("should return 400 for lowercase state query", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?state=watched", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "State", "oneof")
					})
					It("should return 400 for unsupported state query", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?state=ARCHIVED", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "State", "oneof")
					})
				})
			})

			Context("Trend Query Parameter", func() {
				Context("Allowed Values", func() {
					for _, trendValue := range []string{"UPTREND", "SIDEWAYS", "DOWNTREND"} {
						value := trendValue
						It("should filter by "+value, func() {
							req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?trend="+value, nil)
							router.ServeHTTP(w, req)
							response := decodeTickerListResponse(w)
							for _, ticker := range response.Tickers {
								Expect(ticker.Trend).To(Equal(value))
							}
						})
					}
				})
				Context("Bad Values", func() {
					It("should return 400 for lowercase trend query", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?trend=uptrend", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Trend", "oneof")
					})
					It("should return 400 for unsupported trend query", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?trend=NEUTRAL", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Trend", "oneof")
					})
				})
			})

			Context("Is FNO Query Parameter", func() {
				Context("Allowed Values", func() {
					It("should filter by is-fno=true", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?is-fno=true", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						for _, ticker := range response.Tickers {
							Expect(ticker.IsFNO).To(BeTrue())
						}
					})
					It("should filter by is-fno=false", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?is-fno=false", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						for _, ticker := range response.Tickers {
							Expect(ticker.IsFNO).To(BeFalse())
						}
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for non-boolean is-fno", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?is-fno=maybe", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Opened After Query Parameter", func() {
				Context("Allowed Values", func() {
					It("should filter inclusively by opened-after RFC3339 timestamp", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?opened-after=2026-05-06T10:30:00Z", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						for _, ticker := range response.Tickers {
							Expect(ticker.LastOpenedAt).To(BeTemporally(">=", time.Date(2026, time.May, 6, 10, 30, 0, 0, time.UTC)))
						}
					})
					It("should return empty list when no ticker opened after timestamp", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?opened-after=2030-01-01T00:00:00Z", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						Expect(response.Tickers).To(BeEmpty())
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for date-only opened-after", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?opened-after=2026-05-06", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "OpenedAfter", "datetime")
					})
					It("should return 400 for timestamp without timezone", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?opened-after=2026-05-06T10:30:00", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "OpenedAfter", "datetime")
					})
					It("should return 400 for invalid opened-after text", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?opened-after=invalid", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "OpenedAfter", "datetime")
					})
				})
			})

			Context("Sorting Query Parameters", func() {
				Context("Allowed Values", func() {
					It("should sort by ticker ascending by default", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase, nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						Expect(response.Tickers[0].Ticker).To(Equal("BTCUSD"))
					})
					It("should sort by ticker descending", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?sort-by=ticker&sort-order=desc", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						Expect(response.Tickers[0].Ticker).To(Equal(tickerNiftyUSDINR))
					})

					It("should sort by last_opened_at ascending", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?sort-by=last_opened_at&sort-order=asc", nil)
						router.ServeHTTP(w, req)
						tickers := decodeTickerListResponse(w).Tickers
						Expect(tickers).To(HaveLen(3))
						Expect(tickers[0].Ticker).To(Equal("MCX"))
						Expect(tickers[1].Ticker).To(Equal("BTCUSD"))
						Expect(tickers[2].Ticker).To(Equal(tickerNiftyUSDINR))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for unsupported sort-by", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?sort-by=unsupported", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "SortBy", "oneof")
					})
					It("should return 400 for unsupported sort-order", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?sort-order=up", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "SortOrder", "oneof")
					})
					It("should return 400 for uppercase sort-order", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?sort-order=DESC", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "SortOrder", "oneof")
					})
				})
			})

			Context("Pagination Query Parameters", func() {
				Context("Allowed Values", func() {
					It("should accept offset 0", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?offset=0", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						Expect(response.Metadata.Offset).To(Equal(0))
					})
					It("should accept positive offset", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?offset=1", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						Expect(response.Metadata.Offset).To(Equal(1))
					})
					It("should accept limit 1", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?limit=1", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						Expect(response.Tickers).To(HaveLen(1))
						Expect(response.Metadata.Limit).To(Equal(1))
					})
					It("should accept limit 100", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?limit=100", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						Expect(response.Metadata.Limit).To(Equal(100))
					})
					It("should return empty list for offset beyond total", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?offset=100", nil)
						router.ServeHTTP(w, req)
						response := decodeTickerListResponse(w)
						Expect(response.Tickers).To(BeEmpty())
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for negative offset", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?offset=-1", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Offset", "min")
					})
					It("should return 400 for non-numeric offset", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?offset=abc", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "message", "numeric")
					})
					It("should return 400 for limit 0", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?limit=0", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Limit", "min")
					})
					It("should return 400 for limit greater than 100", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?limit=101", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Limit", "max")
					})
					It("should return 400 for non-numeric limit", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase+"?limit=abc", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "message", "numeric")
					})
				})
			})
		})

		Context("Errors", func() {
			It("should return 500 when repository list fails", func() {
				sqlDB, err := db.DB()
				Expect(err).ToNot(HaveOccurred())
				Expect(sqlDB.Close()).To(Succeed())
				req, w := util.CreateTestRequest(http.MethodGet, barkat.TickerBase, nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
