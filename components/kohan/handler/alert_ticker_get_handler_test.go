//nolint:dupl
package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func decodeAlertTickerOKResponse(w *httptest.ResponseRecorder) barkat.AlertTicker {
	var envelope common.Envelope[map[string]barkat.AlertTicker]
	util.AssertSuccess(w, http.StatusOK, &envelope)
	return envelope.Data["alert_ticker"]
}

func decodeAlertTickerListResponse(w *httptest.ResponseRecorder) barkat.AlertTickerList {
	var envelope common.Envelope[barkat.AlertTickerList]
	util.AssertSuccess(w, http.StatusOK, &envelope)
	return envelope.Data
}

func seedAlertTicker(ctx context.Context, db *gorm.DB, alertTicker barkat.AlertTicker) barkat.AlertTicker {
	Expect(db.WithContext(ctx).Create(&alertTicker).Error).ToNot(HaveOccurred())
	return alertTicker
}

// AlertTickerHandler Integration GET/List Tests - Comprehensive Master Specification.
// Tests complete HTTP → Handler → Manager → Repository → Database flow for PRD Section 2.2.2.2 and 2.2.2.4.
var _ = PDescribe("AlertTickerHandler Integration - GET/List Tests - Section 2.2.2 Alert Ticker APIs", func() {
	var (
		alertTickerHandler      handler.AlertTickerHandler
		router                  *gin.Engine
		testCtx                 = context.Background()
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
		router = newAlertTickerTestRouter(alertTickerHandler)
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).ToNot(HaveOccurred())
		sqlDB.Close()
	})

	Describe("GET /v1/api/alert-tickers/{symbol} - Retrieve Alert Ticker (2.2.2.2)", func() {
		var createdAlertTicker barkat.AlertTicker

		BeforeEach(func() {
			createdAlertTicker = validAlertTickerPayload
			createdAlertTicker.TickerID = createdTicker.ID
			createdAlertTicker = seedAlertTicker(testCtx, db, createdAlertTicker)
		})

		Context("Happy Path", func() {
			Context("with existing alert ticker", func() {
				var response barkat.AlertTicker

				BeforeEach(func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"/"+createdAlertTicker.Symbol, nil)
					router.ServeHTTP(w, req)
					response = decodeAlertTickerOKResponse(w)
				})

				It("should return 200 OK", func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"/"+createdAlertTicker.Symbol, nil)
					router.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusOK))
				})
				It("should return Envelope success", func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"/"+createdAlertTicker.Symbol, nil)
					router.ServeHTTP(w, req)
					var envelope common.Envelope[map[string]barkat.AlertTicker]
					util.AssertSuccess(w, http.StatusOK, &envelope)
					Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
				})
				It("should return alert ticker with correct fields", func() {
					Expect(response.Symbol).To(Equal("MCIX"))
					Expect(response.PairID).To(Equal("941982"))
					Expect(response.Name).To(Equal("Multi Commodity Exchange of India"))
					Expect(response.Exchange).To(Equal(new("NSE")))
				})
				It("should include parent ticker reference", func() { Expect(response.Ticker).To(Equal(createdTicker.Ticker)) })
				It("should include created_at and updated_at", func() {
					Expect(response.CreatedAt).ToNot(BeZero())
					Expect(response.UpdatedAt).ToNot(BeZero())
				})
			})
		})

		Context("Field Validations", func() {
			Context("Symbol Path Parameter", func() {
				Context("Allowed Values", func() {
					It("should accept valid existing symbol path", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"/"+createdAlertTicker.Symbol, nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for invalid symbol path", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"/.MCIX", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 404 for valid missing symbol path", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"/NOTFOUND", nil)
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
				req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"/"+createdAlertTicker.Symbol, nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("GET /v1/api/alert-tickers - List Alert Tickers (2.2.2.4)", func() {
		BeforeEach(func() {
			seedAlertTicker(testCtx, db, barkat.AlertTicker{TickerID: createdTicker.ID, Symbol: "MCIX", PairID: "941982", Name: "Multi Commodity Exchange of India", Exchange: new("NSE")})
			seedAlertTicker(testCtx, db, barkat.AlertTicker{TickerID: createdTicker.ID, Symbol: "GOLD1!", PairID: "100200", Name: "Gold Index", Exchange: new("N.SE")})
			seedAlertTicker(testCtx, db, barkat.AlertTicker{TickerID: createdTicker.ID, Symbol: "BTCUSD", PairID: "000777", Name: "Bitcoin USD", Exchange: new("BINANCE")})
		})

		Context("Happy Path", func() {
			Context("default pagination", func() {
				var response barkat.AlertTickerList

				BeforeEach(func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase, nil)
					router.ServeHTTP(w, req)
					response = decodeAlertTickerListResponse(w)
				})

				It("should return 200 OK", func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase, nil)
					router.ServeHTTP(w, req)
					Expect(w.Code).To(Equal(http.StatusOK))
				})
				It("should return Envelope success", func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase, nil)
					router.ServeHTTP(w, req)
					var envelope common.Envelope[barkat.AlertTickerList]
					util.AssertSuccess(w, http.StatusOK, &envelope)
					Expect(envelope.Status).To(Equal(common.EnvelopeSuccess))
				})
				It("should return alert_tickers array", func() { Expect(response.AlertTickers).To(HaveLen(3)) })
				It("should return metadata offset 0", func() { Expect(response.Metadata.Offset).To(Equal(0)) })
				It("should return metadata limit 20", func() { Expect(response.Metadata.Limit).To(Equal(20)) })
				It("should return metadata total", func() { Expect(response.Metadata.Total).To(Equal(int64(3))) })
			})

			Context("response shape", func() {
				var alertTicker barkat.AlertTicker

				BeforeEach(func() {
					req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase, nil)
					router.ServeHTTP(w, req)
					response := decodeAlertTickerListResponse(w)
					alertTicker = response.AlertTickers[0]
				})

				It("should include symbol", func() { Expect(alertTicker.Symbol).ToNot(BeEmpty()) })
				It("should include pair_id", func() { Expect(alertTicker.PairID).ToNot(BeEmpty()) })
				It("should include name", func() { Expect(alertTicker.Name).ToNot(BeEmpty()) })
				It("should include exchange", func() { Expect(alertTicker.Exchange).ToNot(BeNil()) })
				It("should include ticker", func() { Expect(alertTicker.Ticker).ToNot(BeEmpty()) })
			})
		})

		Context("Field Validations", func() {
			Context("Symbol Query Parameter", func() {
				Context("Allowed Values", func() {
					It("should filter by exact symbol", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?symbol=MCIX", nil)
						router.ServeHTTP(w, req)
						response := decodeAlertTickerListResponse(w)
						Expect(response.AlertTickers).To(HaveLen(1))
						Expect(response.AlertTickers[0].Symbol).To(Equal("MCIX"))
					})
					It("should accept minimum symbol length 1", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?symbol=A", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should accept maximum symbol length 25", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?symbol="+strings.Repeat("A", 25), nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should accept letters and digits in symbol query", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?symbol=ABC123", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should accept dot in symbol query", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?symbol=ABC.D", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should accept underscore in symbol query", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?symbol=ABC_D", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should return empty list for no symbol match", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?symbol=ZZZ", nil)
						router.ServeHTTP(w, req)
						response := decodeAlertTickerListResponse(w)
						Expect(response.AlertTickers).To(BeEmpty())
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for invalid symbol query starting with dot", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?symbol=.MCIX", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for symbol query exceeding 25 characters", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?symbol="+strings.Repeat("A", 26), nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for symbol query with unsupported special character", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?symbol=MCIX@", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for symbol query with hyphen", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?symbol=ABC-D", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Ticker Query Parameter", func() {
				Context("Allowed Values", func() {
					It("should filter by parent ticker", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?ticker="+createdTicker.Ticker, nil)
						router.ServeHTTP(w, req)
						response := decodeAlertTickerListResponse(w)
						Expect(response.AlertTickers).To(HaveLen(3))
					})
					It("should return empty list for no child match under existing ticker", func() {
						otherTicker := barkat.Ticker{
							Ticker:       "NIFTY",
							Exchange:     new("NSE"),
							Timeframes:   []string{"MN", "WK", "DL"},
							Type:         "EQUITY",
							State:        "WATCHED",
							Trend:        "UPTREND",
							LastOpenedAt: time.Date(2026, time.May, 5, 10, 30, 0, 0, time.UTC),
							IsFNO:        true,
						}
						Expect(db.Create(&otherTicker).Error).ToNot(HaveOccurred())
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?ticker=NIFTY", nil)
						router.ServeHTTP(w, req)
						response := decodeAlertTickerListResponse(w)
						Expect(response.AlertTickers).To(BeEmpty())
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for invalid ticker query", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?ticker=mcx", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 404 for valid missing ticker query", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?ticker=NOTFOUND", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusNotFound))
					})
				})
			})

			Context("Pair ID Query Parameter", func() {
				Context("Allowed Values", func() {
					It("should filter by exact pair-id", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?pair-id=941982", nil)
						router.ServeHTTP(w, req)
						response := decodeAlertTickerListResponse(w)
						Expect(response.AlertTickers).To(HaveLen(1))
						Expect(response.AlertTickers[0].PairID).To(Equal("941982"))
					})
					It("should accept minimum pair-id length 1", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?pair-id=1", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should accept maximum pair-id length 64", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?pair-id="+strings.Repeat("1", 64), nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should preserve leading zeroes", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?pair-id=00123", nil)
						router.ServeHTTP(w, req)
						response := decodeAlertTickerListResponse(w)
						Expect(response.AlertTickers).To(HaveLen(1))
						Expect(response.AlertTickers[0].PairID).To(Equal("000777"))
					})
					It("should return empty list for no pair-id match", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?pair-id=999999", nil)
						router.ServeHTTP(w, req)
						response := decodeAlertTickerListResponse(w)
						Expect(response.AlertTickers).To(BeEmpty())
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for non-digit pair-id", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?pair-id=94A982", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for pair-id exceeding 64 characters", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?pair-id="+strings.Repeat("1", 65), nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for negative pair-id", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?pair-id=-941982", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for whitespace in pair-id", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?pair-id=941%20982", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Exchange Query Parameter", func() {
				Context("Allowed Values", func() {
					It("should filter by exact exchange", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange=NSE", nil)
						router.ServeHTTP(w, req)
						response := decodeAlertTickerListResponse(w)
						Expect(response.AlertTickers).To(HaveLen(1))
						Expect(*response.AlertTickers[0].Exchange).To(Equal("NSE"))
					})
					It("should accept minimum exchange length 1", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange=N", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should accept maximum exchange length 10", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange="+strings.Repeat("A", 10), nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should accept uppercase letters NSE", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange=NSE", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should accept lowercase letters nse", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange=nse", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should accept digits in exchange code", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange=NSE1", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should accept dot in exchange code", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange=N.SE", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
					It("should accept underscore in exchange code", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange=N_SE", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusOK))
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for empty exchange query", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange=", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for exchange exceeding 10 characters", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange="+strings.Repeat("A", 11), nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for exchange with colon", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange=NSE:MCX", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for exchange with hyphen", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange=N-SE", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for exchange with whitespace", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange=NS%20E", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
					It("should return 400 for exchange with unsupported special character", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?exchange=NSE@", nil)
						router.ServeHTTP(w, req)
						Expect(w.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})

			Context("Pagination Query Parameters", func() {
				Context("Allowed Values", func() {
					It("should accept offset 0", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?offset=0", nil)
						router.ServeHTTP(w, req)
						response := decodeAlertTickerListResponse(w)
						Expect(response.Metadata.Offset).To(Equal(0))
					})
					It("should accept positive offset", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?offset=1", nil)
						router.ServeHTTP(w, req)
						response := decodeAlertTickerListResponse(w)
						Expect(response.Metadata.Offset).To(Equal(1))
					})
					It("should accept limit 1", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?limit=1", nil)
						router.ServeHTTP(w, req)
						response := decodeAlertTickerListResponse(w)
						Expect(response.AlertTickers).To(HaveLen(1))
						Expect(response.Metadata.Limit).To(Equal(1))
					})
					It("should accept limit 100", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?limit=100", nil)
						router.ServeHTTP(w, req)
						response := decodeAlertTickerListResponse(w)
						Expect(response.Metadata.Limit).To(Equal(100))
					})
					It("should return empty list for offset beyond total", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?offset=100", nil)
						router.ServeHTTP(w, req)
						response := decodeAlertTickerListResponse(w)
						Expect(response.AlertTickers).To(BeEmpty())
					})
				})
				Context("Bad Values", func() {
					It("should return 400 for negative offset", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?offset=-1", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Offset", "min")
					})
					It("should return 400 for non-numeric offset", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?offset=abc", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "message", "numeric")
					})
					It("should return 400 for limit 0", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?limit=0", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Limit", "min")
					})
					It("should return 400 for limit greater than 100", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?limit=101", nil)
						router.ServeHTTP(w, req)
						util.AssertError(w, "Limit", "max")
					})
					It("should return 400 for non-numeric limit", func() {
						req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase+"?limit=abc", nil)
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
				req, w := util.CreateTestRequest(http.MethodGet, barkat.AlertTickerBase, nil)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
