package clients_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/amanhigh/go-fun/components/kohan/clients"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("YahooClient", func() {
	var (
		server              *httptest.Server
		yahooClient         *clients.YahooClient
		ctx                 context.Context
		ticker              string
		tickerDataStartYear int
	)

	BeforeEach(func() {
		ctx = context.Background()
		ticker = "AAPL"
		tickerDataStartYear = 2020
	})

	AfterEach(func() {
		if server != nil {
			server.Close()
		}
	})

	Context("FetchDailyPrices", func() {
		Context("successful fetch with valid response", func() {
			var stockData tax.StockData
			var err error

			BeforeEach(func() {
				responseBody := `{
					"chart": {
						"result": [
							{
								"meta": {
									"currency": "USD",
									"symbol": "AAPL",
									"exchangeName": "NMS"
								},
								"timestamp": [1705276800, 1705363200, 1705449600],
								"indicators": {
									"quote": [
										{
											"open": [190.0, 191.0, 192.0],
											"high": [192.0, 193.0, 194.0],
											"low": [189.0, 190.0, 191.0],
											"close": [191.5, 192.5, 193.5],
											"volume": [50000000, 51000000, 52000000]
										}
									]
								}
							}
						],
						"error": null
					}
				}`

				server = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						Expect(r.URL.Path).To(ContainSubstring("/v8/finance/chart/AAPL"))
						Expect(r.URL.Query().Get("interval")).To(Equal("1d"))
						Expect(r.URL.Query().Get("period1")).To(Equal("1577836800")) // 2020-01-01
						Expect(r.URL.Query().Get("period2")).ToNot(BeEmpty())        // Current timestamp
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(responseBody))
					}),
				)

				client := resty.NewWithClient(&http.Client{})
				yahooClient = clients.NewYahooClient(client, server.URL, tickerDataStartYear)
				stockData, err = yahooClient.FetchDailyPrices(ctx, ticker)
			})

			It("should return no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return correct number of prices", func() {
				Expect(stockData.Prices).To(HaveLen(3))
			})

			It("should correctly parse closing prices for each date", func() {
				Expect(stockData.Prices["2024-01-15"]).To(Equal(191.5))
				Expect(stockData.Prices["2024-01-16"]).To(Equal(192.5))
				Expect(stockData.Prices["2024-01-17"]).To(Equal(193.5))
			})

			It("should correctly convert Unix timestamps to YYYY-MM-DD format", func() {
				Expect(stockData.Prices).To(HaveKey("2024-01-15"))
				Expect(stockData.Prices).To(HaveKey("2024-01-16"))
				Expect(stockData.Prices).To(HaveKey("2024-01-17"))
			})

			It("should return non-nil empty splits when no events present", func() {
				Expect(stockData.Splits).ToNot(BeNil())
				Expect(stockData.Splits).To(BeEmpty())
			})
		})

		Context("empty chart result", func() {
			var err error

			BeforeEach(func() {
				responseBody := `{
					"chart": {
						"result": [],
						"error": null
					}
				}`

				server = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(responseBody))
					}),
				)

				client := resty.NewWithClient(&http.Client{})
				yahooClient = clients.NewYahooClient(client, server.URL, tickerDataStartYear)
				_, err = yahooClient.FetchDailyPrices(ctx, "INVALID")
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should indicate no data found", func() {
				Expect(err.Error()).To(ContainSubstring("no data found for INVALID"))
			})
		})

		Context("missing quote data", func() {
			var err error

			BeforeEach(func() {
				responseBody := `{
					"chart": {
						"result": [
							{
								"meta": {
									"currency": "USD",
									"symbol": "AAPL"
								},
								"timestamp": [1705276800],
								"indicators": {
									"quote": []
								}
							}
						],
						"error": null
					}
				}`

				server = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(responseBody))
					}),
				)

				client := resty.NewWithClient(&http.Client{})
				yahooClient = clients.NewYahooClient(client, server.URL, tickerDataStartYear)
				_, err = yahooClient.FetchDailyPrices(ctx, ticker)
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should indicate no quote data found", func() {
				Expect(err.Error()).To(ContainSubstring("no quote data found for AAPL"))
			})
		})

		Context("empty close prices", func() {
			var err error

			BeforeEach(func() {
				responseBody := `{
					"chart": {
						"result": [
							{
								"meta": {
									"currency": "USD",
									"symbol": "AAPL"
								},
								"timestamp": [1705276800],
								"indicators": {
									"quote": [
										{
											"open": [190.0],
											"high": [192.0],
											"low": [189.0],
											"close": [],
											"volume": [50000000]
										}
									]
								}
							}
						],
						"error": null
					}
				}`

				server = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(responseBody))
					}),
				)

				client := resty.NewWithClient(&http.Client{})
				yahooClient = clients.NewYahooClient(client, server.URL, tickerDataStartYear)
				_, err = yahooClient.FetchDailyPrices(ctx, ticker)
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should indicate no price data found", func() {
				Expect(err.Error()).To(ContainSubstring("no price data found for AAPL"))
			})
		})

		Context("mismatched timestamp and close price counts", func() {
			var stockData tax.StockData
			var err error

			BeforeEach(func() {
				responseBody := `{
					"chart": {
						"result": [
							{
								"meta": {
									"currency": "USD",
									"symbol": "AAPL"
								},
								"timestamp": [1705276800, 1705363200, 1705449600],
								"indicators": {
									"quote": [
										{
											"open": [190.0, 191.0, 192.0],
											"high": [192.0, 193.0, 194.0],
											"low": [189.0, 190.0, 191.0],
											"close": [191.5, 192.5],
											"volume": [50000000, 51000000, 52000000]
										}
									]
								}
							}
						],
						"error": null
					}
				}`

				server = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(responseBody))
					}),
				)

				client := resty.NewWithClient(&http.Client{})
				yahooClient = clients.NewYahooClient(client, server.URL, tickerDataStartYear)
				stockData, err = yahooClient.FetchDailyPrices(ctx, ticker)
			})

			It("should return no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should only use available close prices", func() {
				Expect(stockData.Prices).To(HaveLen(2))
			})

			It("should have prices for first two dates", func() {
				Expect(stockData.Prices).To(HaveKey("2024-01-15"))
				Expect(stockData.Prices).To(HaveKey("2024-01-16"))
			})

			It("should not have price for third date", func() {
				Expect(stockData.Prices).ToNot(HaveKey("2024-01-17"))
			})
		})

		Context("request includes User-Agent header", func() {
			var userAgentVerified bool
			var err error

			BeforeEach(func() {
				responseBody := `{
					"chart": {
						"result": [
							{
								"meta": {
									"currency": "USD",
									"symbol": "AAPL"
								},
								"timestamp": [1705276800],
								"indicators": {
									"quote": [
										{
											"open": [190.0],
											"high": [192.0],
											"low": [189.0],
											"close": [191.5],
											"volume": [50000000]
										}
									]
								}
							}
						],
						"error": null
					}
				}`

				server = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						userAgent := r.Header.Get("User-Agent")
						Expect(userAgent).To(ContainSubstring("Mozilla/5.0"))
						userAgentVerified = true
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(responseBody))
					}),
				)

				client := resty.NewWithClient(&http.Client{})
				yahooClient = clients.NewYahooClient(client, server.URL, tickerDataStartYear)
				_, err = yahooClient.FetchDailyPrices(ctx, ticker)
			})

			It("should return no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should have verified User-Agent header", func() {
				Expect(userAgentVerified).To(BeTrue())
			})
		})

		Context("split events parsing", func() {
			var (
				stockData      tax.StockData
				err            error
				eventsVerified bool
			)

			BeforeEach(func() {
				// Two out-of-order events: 1705363200 listed first, 1705276800 listed second.
				// Includes a 4:1 split (numerator:4, denominator:1).
				responseBody := `{
					"chart": {
						"result": [
							{
								"meta": {"currency": "USD", "symbol": "AAPL"},
								"timestamp": [1705276800],
								"indicators": {
									"quote": [{"open": [190.0], "high": [192.0], "low": [189.0], "close": [191.5], "volume": [50000000]}]
								},
								"events": {
									"splits": {
										"1705363200": {
											"date": 1705363200,
											"numerator": 2,
											"denominator": 1
										},
										"1705276800": {
											"date": 1705276800,
											"numerator": 4,
											"denominator": 1
										}
									}
								}
							}
						],
						"error": null
					}
				}`

				server = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						Expect(r.URL.Query().Get("events")).To(Equal("splits"))
						eventsVerified = true
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(responseBody))
					}),
				)

				client := resty.NewWithClient(&http.Client{})
				yahooClient = clients.NewYahooClient(client, server.URL, tickerDataStartYear)
				stockData, err = yahooClient.FetchDailyPrices(ctx, ticker)
			})

			It("should return no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should have verified events=splits query parameter", func() {
				Expect(eventsVerified).To(BeTrue())
			})

			It("should have two splits in chronological order", func() {
				Expect(stockData.Splits).To(HaveLen(2))
				Expect(stockData.Splits[0].Date).To(Equal(int64(1705276800)))
				Expect(stockData.Splits[0].Numerator).To(Equal(4.0))
				Expect(stockData.Splits[0].Denominator).To(Equal(1.0))
				Expect(stockData.Splits[1].Date).To(Equal(int64(1705363200)))
				Expect(stockData.Splits[1].Numerator).To(Equal(2.0))
				Expect(stockData.Splits[1].Denominator).To(Equal(1.0))
			})
		})

		Context("GetSecurityInfo", func() {
			var (
				query   string
				results []tax.SecurityInfo
				httpErr error
			)

			Context("ISIN resolves to FISV equity candidate", func() {
				BeforeEach(func() {
					query = "US...FISV"
					responseBody := `{
						"quotes": [
							{
								"symbol": "FISV",
								"longname": "Fiserv Inc.",
								"shortname": "Fiserv",
								"exchange": "NMS",
								"quoteType": "EQUITY"
							}
						]
					}`

					server = httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							Expect(r.URL.Path).To(Equal("/v1/finance/search"))
							Expect(r.URL.Query().Get("q")).To(Equal(query))
							Expect(r.URL.Query().Get("quotesCount")).To(Equal("20"))
							Expect(r.URL.Query().Get("newsCount")).To(Equal("0"))
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusOK)
							_, _ = w.Write([]byte(responseBody))
						}),
					)

					client := resty.NewWithClient(&http.Client{})
					yahooClient = clients.NewYahooClient(client, server.URL, tickerDataStartYear)
					results, httpErr = yahooClient.GetSecurityInfo(ctx, query)
				})

				It("should return no error", func() {
					Expect(httpErr).ToNot(HaveOccurred())
				})

				It("should return exactly one candidate", func() {
					Expect(results).To(HaveLen(1))
				})

				It("should use long name", func() {
					Expect(results[0].Name).To(Equal("Fiserv Inc."))
				})

				It("should have correct symbol and type", func() {
					Expect(results[0].Symbol).To(Equal("FISV"))
					Expect(results[0].Type).To(Equal("EQUITY"))
				})
			})

			Context("ETF classification for VTI", func() {
				BeforeEach(func() {
					query = "VTI"
					responseBody := `{
						"quotes": [
							{
								"symbol": "VTI",
								"longname": "Vanguard Total Stock Market ETF",
								"shortname": "Vanguard Total Stock Market",
								"exchange": "NMS",
								"quoteType": "ETF"
							}
						]
					}`

					server = httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							Expect(r.URL.Path).To(Equal("/v1/finance/search"))
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusOK)
							_, _ = w.Write([]byte(responseBody))
						}),
					)

					client := resty.NewWithClient(&http.Client{})
					yahooClient = clients.NewYahooClient(client, server.URL, tickerDataStartYear)
					results, httpErr = yahooClient.GetSecurityInfo(ctx, query)
				})

				It("should return no error", func() {
					Expect(httpErr).ToNot(HaveOccurred())
				})

				It("should return exactly one candidate", func() {
					Expect(results).To(HaveLen(1))
				})

				It("should classify as ETF", func() {
					Expect(results[0].Type).To(Equal("ETF"))
				})

				It("should have correct symbol and name", func() {
					Expect(results[0].Symbol).To(Equal("VTI"))
					Expect(results[0].Name).To(Equal("Vanguard Total Stock Market ETF"))
				})
			})

			Context("multiple candidates without selection", func() {
				BeforeEach(func() {
					query = "BRK"
					responseBody := `{
						"quotes": [
							{
								"symbol": "BRK-A",
								"longname": "Berkshire Hathaway Inc.",
								"shortname": "Berkshire Hathaway",
								"exchange": "NYSE",
								"quoteType": "EQUITY"
							},
							{
								"symbol": "BRK-B",
								"longname": "",
								"shortname": "Berkshire Hathaway B",
								"exchange": "NYSE",
								"quoteType": "EQUITY"
							}
						]
					}`

					server = httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusOK)
							_, _ = w.Write([]byte(responseBody))
						}),
					)

					client := resty.NewWithClient(&http.Client{})
					yahooClient = clients.NewYahooClient(client, server.URL, tickerDataStartYear)
					results, httpErr = yahooClient.GetSecurityInfo(ctx, query)
				})

				It("should return no error", func() {
					Expect(httpErr).ToNot(HaveOccurred())
				})

				It("should return both candidates", func() {
					Expect(results).To(HaveLen(2))
				})

				It("should use long name when available", func() {
					Expect(results[0].Name).To(Equal("Berkshire Hathaway Inc."))
				})

				It("should fall back to short name when long name is empty", func() {
					Expect(results[1].Name).To(Equal("Berkshire Hathaway B"))
				})

				It("should not select one over the other", func() {
					Expect(results[0].Symbol).To(Equal("BRK-A"))
					Expect(results[1].Symbol).To(Equal("BRK-B"))
				})
			})

			Context("empty results return non-nil empty slice", func() {
				BeforeEach(func() {
					query = "ZZXXYYZZ"
					responseBody := `{
						"quotes": []
					}`

					server = httptest.NewServer(
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusOK)
							_, _ = w.Write([]byte(responseBody))
						}),
					)

					client := resty.NewWithClient(&http.Client{})
					yahooClient = clients.NewYahooClient(client, server.URL, tickerDataStartYear)
					results, httpErr = yahooClient.GetSecurityInfo(ctx, query)
				})

				It("should return no error", func() {
					Expect(httpErr).ToNot(HaveOccurred())
				})

				It("should return empty non-nil slice", func() {
					Expect(results).ToNot(BeNil())
					Expect(results).To(BeEmpty())
				})
			})
		})

		Context("multiple different tickers", func() {
			var aaplData tax.StockData
			var msftData tax.StockData
			var aaplErr error
			var msftErr error
			var requestCount int

			BeforeEach(func() {
				requestCount = 0
				server = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						tickerParam := r.URL.Path[len("/v8/finance/chart/"):]
						var response string

						switch tickerParam {
						case "AAPL":
							response = `{
								"chart": {
									"result": [
										{
											"meta": {"currency": "USD", "symbol": "AAPL"},
											"timestamp": [1705276800],
											"indicators": {
												"quote": [{"open": [190.0], "high": [192.0], "low": [189.0], "close": [191.5], "volume": [50000000]}]
											}
										}
									],
									"error": null
								}
							}`
						case "MSFT":
							response = `{
								"chart": {
									"result": [
										{
											"meta": {"currency": "USD", "symbol": "MSFT"},
											"timestamp": [1705276800],
											"indicators": {
												"quote": [{"open": [380.0], "high": [385.0], "low": [378.0], "close": [382.5], "volume": [20000000]}]
											}
										}
									],
									"error": null
								}
							}`
						}

						requestCount++
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(response))
					}),
				)

				client := resty.NewWithClient(&http.Client{})
				yahooClient = clients.NewYahooClient(client, server.URL, tickerDataStartYear)
				aaplData, aaplErr = yahooClient.FetchDailyPrices(ctx, "AAPL")
				msftData, msftErr = yahooClient.FetchDailyPrices(ctx, "MSFT")
			})

			It("should return no error for AAPL", func() {
				Expect(aaplErr).ToNot(HaveOccurred())
			})

			It("should return no error for MSFT", func() {
				Expect(msftErr).ToNot(HaveOccurred())
			})

			It("should return correct price for AAPL", func() {
				Expect(aaplData.Prices["2024-01-15"]).To(Equal(191.5))
			})

			It("should return correct price for MSFT", func() {
				Expect(msftData.Prices["2024-01-15"]).To(Equal(382.5))
			})

			It("should have made two requests", func() {
				Expect(requestCount).To(Equal(2))
			})
		})
	})
})
