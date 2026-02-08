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
