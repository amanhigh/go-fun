package manager_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xuri/excelize/v2"
)

var _ = Describe("DriveWealthManagerImpl", func() {
	var (
		tempTestDir        string
		basePath           string
		driveWealthManager manager.Broker
	)

	BeforeEach(func() {
		var err error
		tempTestDir, err = os.MkdirTemp("", "drivewealth_test_*")
		Expect(err).ToNot(HaveOccurred())

		basePath = filepath.Join(tempTestDir, "vested")
	})

	AfterEach(func() {
		err := os.RemoveAll(tempTestDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("with a valid Excel file", func() {
		BeforeEach(func() {
			f := excelize.NewFile()
			sheetName := "Income"
			_, err := f.NewSheet(sheetName)
			Expect(err).ToNot(HaveOccurred())

			headers := []string{"Date", "Time (in UTC)", "Activity", "Ticker", "Gross Cash Amount (in USD)"}
			err = f.SetSheetRow(sheetName, "A1", &headers)
			Expect(err).ToNot(HaveOccurred())

			rows := [][]interface{}{
				{"2024-06-06", "05:24:40 AM", "Dividend", "IEF", 158.17},
				{"2024-06-06", "05:24:39 AM", "Tax", "IEF", -39.54},
				{"2024-06-06", "05:24:37 AM", "Dividend", "TLT", 3.51},
				{"2024-06-06", "05:24:37 AM", "Tax", "BIL", -35.52},
				{"2024-06-06", "05:24:36 AM", "Dividend", "BIL", 142.07},
				{"2024-06-03", "05:34:52 AM", "Interest", "", 0.59},
				{"2024-05-02", "04:57:05 AM", "Interest", "", 1.18},
			}

			for i, rowData := range rows {
				err = f.SetSheetRow(sheetName, fmt.Sprintf("A%d", i+2), &rowData)
				Expect(err).ToNot(HaveOccurred())
			}

			Expect(f.DeleteSheet("Sheet1")).To(Succeed())

			tradeSheet := "Trades"
			_, err = f.NewSheet(tradeSheet)
			Expect(err).ToNot(HaveOccurred())

			tradeHeaders := []string{"Date", "Time (in UTC)", "Name", "Ticker", "Activity", "Order Type", "Quantity", "Price Per Share (in USD)", "Cash Amount (in USD)", "Commission Charges (in USD)"}
			err = f.SetSheetRow(tradeSheet, "A1", &tradeHeaders)
			Expect(err).ToNot(HaveOccurred())

			tradeRows := [][]interface{}{
				{"2024-04-03", "04:53:52 PM", "Vanguard Russell 2000 ETF", "VTWO", "Buy", "Market", 70, 77.41, 5418.7, 0},
				{"2024-04-03", "04:26:02 PM", "Europe ETF FTSE Vanguard", "VGK", "Buy", "Market", 60, 70.18, 4210.8, 0},
				{"2024-02-14", "02:30:00 PM", "Barclays 1-3 Month T-Bill ETF SPDR", "BIL", "Sell", "Market", 1, 91.58, 91.58, 0},
			}

			for i, rowData := range tradeRows {
				err = f.SetSheetRow(tradeSheet, fmt.Sprintf("A%d", i+2), &rowData)
				Expect(err).ToNot(HaveOccurred())
			}

			err = f.SaveAs(basePath + "_2024.xlsx")
			Expect(err).ToNot(HaveOccurred())

			driveWealthManager = manager.NewDriveWealthManagerImpl(basePath)
		})

		Context("when parsing interests", func() {
			It("should extract interest entries correctly", func() {
				info, err := driveWealthManager.Parse(testYear)
				Expect(err).ToNot(HaveOccurred())
				Expect(info.Interests).To(HaveLen(2))
				Expect(info.Interests[0].Amount).To(Equal(0.59))
				Expect(info.Interests[1].Amount).To(Equal(1.18))
			})
		})

		Context("when parsing dividends", func() {
			It("should extract dividend entries correctly", func() {
				info, err := driveWealthManager.Parse(testYear)
				Expect(err).ToNot(HaveOccurred())
				Expect(info.Dividends).To(HaveLen(3))

				Expect(info.Dividends[0].Symbol).To(Equal("IEF"))
				Expect(info.Dividends[0].Amount).To(Equal(158.17))
				Expect(info.Dividends[0].Tax).To(Equal(39.54))
				Expect(info.Dividends[0].Net).To(Equal(118.63))

				Expect(info.Dividends[1].Symbol).To(Equal("TLT"))
				Expect(info.Dividends[1].Amount).To(Equal(3.51))
				Expect(info.Dividends[1].Tax).To(Equal(0.0))
				Expect(info.Dividends[1].Net).To(Equal(3.51))

				Expect(info.Dividends[2].Symbol).To(Equal("BIL"))
				Expect(info.Dividends[2].Amount).To(Equal(142.07))
				Expect(info.Dividends[2].Tax).To(Equal(35.52))
				Expect(info.Dividends[2].Net).To(BeNumerically("~", 106.55, 0.01))
			})
		})

		Context("when parsing trades", func() {
			It("should extract trade entries correctly", func() {
				info, err := driveWealthManager.Parse(testYear)
				Expect(err).ToNot(HaveOccurred())
				Expect(info.Trades).To(HaveLen(3))

				Expect(info.Trades[0].Symbol).To(Equal("VTWO"))
				Expect(info.Trades[0].Quantity).To(Equal(70.0))
				Expect(info.Trades[0].USDPrice).To(Equal(77.41))
				Expect(info.Trades[0].USDValue).To(Equal(5418.7))
				Expect(info.Trades[0].Type).To(Equal("Buy"))

				Expect(info.Trades[1].Symbol).To(Equal("VGK"))
				Expect(info.Trades[1].Quantity).To(Equal(60.0))
				Expect(info.Trades[1].USDPrice).To(Equal(70.18))
				Expect(info.Trades[1].USDValue).To(Equal(4210.8))
				Expect(info.Trades[1].Type).To(Equal("Buy"))

				Expect(info.Trades[2].Symbol).To(Equal("BIL"))
				Expect(info.Trades[2].Quantity).To(Equal(1.0))
				Expect(info.Trades[2].USDPrice).To(Equal(91.58))
				Expect(info.Trades[2].USDValue).To(Equal(91.58))
				Expect(info.Trades[2].Type).To(Equal("Sell"))
			})
		})
	})

	Context("with an invalid or malformed Excel file", func() {
		Context("when the Excel file is missing", func() {
			It("should return an error", func() {
				nonExistentManager := manager.NewDriveWealthManagerImpl(basePath)
				_, err := nonExistentManager.Parse(testYear)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the 'Income' sheet is missing", func() {
			It("should return an error", func() {
				f := excelize.NewFile()
				_, err := f.NewSheet("OtherSheet")
				Expect(err).ToNot(HaveOccurred())
				Expect(f.DeleteSheet("Sheet1")).To(Succeed())
				err = f.SaveAs(basePath + "_2024.xlsx")
				Expect(err).ToNot(HaveOccurred())

				driveWealthManager = manager.NewDriveWealthManagerImpl(basePath)
				_, err = driveWealthManager.Parse(testYear)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the 'Trades' sheet is missing", func() {
			It("should return an error", func() {
				f := excelize.NewFile()
				_, err := f.NewSheet("Income")
				Expect(err).ToNot(HaveOccurred())
				Expect(f.DeleteSheet("Sheet1")).To(Succeed())
				err = f.SaveAs(basePath + "_2024.xlsx")
				Expect(err).ToNot(HaveOccurred())

				driveWealthManager = manager.NewDriveWealthManagerImpl(basePath)
				_, err = driveWealthManager.Parse(testYear)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("commission fallback from All Transactions sheet", func() {
		const testYear = 2024

		Context("when commission fallback logic is applied", func() {
			BeforeEach(func() {
				f := excelize.NewFile()

				// Setup Income sheet
				sheetName := "Income"
				_, err := f.NewSheet(sheetName)
				Expect(err).ToNot(HaveOccurred())

				headers := []string{"Date", "Time (in UTC)", "Activity", "Ticker", "Gross Cash Amount (in USD)"}
				err = f.SetSheetRow(sheetName, "A1", &headers)
				Expect(err).ToNot(HaveOccurred())

				Expect(f.DeleteSheet("Sheet1")).To(Succeed())

				// Create Trades sheet with MIXED commissions (some zero, some non-zero)
				tradeSheet := "Trades"
				_, err = f.NewSheet(tradeSheet)
				Expect(err).ToNot(HaveOccurred())

				tradeHeaders := []string{
					"Date", "Time (in UTC)", "Name", "Ticker", "Activity",
					"Order Type", "Quantity", "Price Per Share (in USD)",
					"Cash Amount (in USD)", "Commission Charges (in USD)",
				}
				err = f.SetSheetRow(tradeSheet, "A1", &tradeHeaders)
				Expect(err).ToNot(HaveOccurred())

				// Mixed: zero commissions (Buy/Sell MCD) and non-zero (Buy AAPL)
				tradeRows := [][]interface{}{
					{"2024-07-01", "02:26:51 PM", "McDonald's Corp.", "MCD", "Buy", "Market", 5, 251.54, 1257.7, 0},  // Zero -> fallback
					{"2024-07-09", "03:12:32 PM", "McDonald's Corp.", "MCD", "Sell", "Market", 5, 245.10, 1225.5, 0}, // Zero -> fallback
					{"2024-07-15", "10:30:00 AM", "Apple Inc.", "AAPL", "Buy", "Market", 10, 225.50, 2255.0, 5.00},   // Non-zero -> use as-is
				}

				for i, rowData := range tradeRows {
					err = f.SetSheetRow(tradeSheet, fmt.Sprintf("A%d", i+2), &rowData)
					Expect(err).ToNot(HaveOccurred())
				}

				// Create All Transactions sheet with commission data for MCD only
				commSheet := "All Transactions"
				_, err = f.NewSheet(commSheet)
				Expect(err).ToNot(HaveOccurred())

				commHeaders := []string{
					"Date", "Time (in UTC)", "Type", "Amount", "Account Balance", "Comment",
				}
				err = f.SetSheetRow(commSheet, "A1", &commHeaders)
				Expect(err).ToNot(HaveOccurred())

				// Commission entries for MCD trades only (AAPL not present to test non-matched)
				commRows := [][]interface{}{
					{"2024-07-01", "02:26:51 PM", "COMM", 2.51, 65358.96, "COMM Buy MCD base=2.51"},
					{"2024-07-09", "03:12:32 PM", "COMM", 3.06, 66844.73, "COMM Sell MCD base=3.06"},
				}

				for i, rowData := range commRows {
					err = f.SetSheetRow(commSheet, fmt.Sprintf("A%d", i+2), &rowData)
					Expect(err).ToNot(HaveOccurred())
				}

				err = f.SaveAs(basePath + "_2024.xlsx")
				Expect(err).ToNot(HaveOccurred())

				driveWealthManager = manager.NewDriveWealthManagerImpl(basePath)
			})

			It("should fallback to All Transactions when Trades sheet has zero commission", func() {
				info, err := driveWealthManager.Parse(testYear)
				Expect(err).ToNot(HaveOccurred())

				Expect(info.Trades).To(HaveLen(3))

				// Find MCD buy trade (zero -> fallback to 2.51)
				var mcdBuy *tax.Trade
				for i := range info.Trades {
					if info.Trades[i].Symbol == "MCD" && info.Trades[i].Type == "Buy" {
						mcdBuy = &info.Trades[i]
						break
					}
				}
				Expect(mcdBuy).ToNot(BeNil())
				Expect(mcdBuy.Commission).To(Equal(2.51))

				// Find MCD sell trade (zero -> fallback to 3.06)
				var mcdSell *tax.Trade
				for i := range info.Trades {
					if info.Trades[i].Symbol == "MCD" && info.Trades[i].Type == "Sell" {
						mcdSell = &info.Trades[i]
						break
					}
				}
				Expect(mcdSell).ToNot(BeNil())
				Expect(mcdSell.Commission).To(Equal(3.06))
			})

			It("should prefer Trades sheet commission over All Transactions when non-zero", func() {
				info, err := driveWealthManager.Parse(testYear)
				Expect(err).ToNot(HaveOccurred())

				Expect(info.Trades).To(HaveLen(3))

				// Find AAPL buy trade (has 5.00 in Trades sheet -> should use 5.00, NOT fallback)
				var aaplBuy *tax.Trade
				for i := range info.Trades {
					if info.Trades[i].Symbol == "AAPL" && info.Trades[i].Type == "Buy" {
						aaplBuy = &info.Trades[i]
						break
					}
				}

				Expect(aaplBuy).ToNot(BeNil())
				Expect(aaplBuy.Commission).To(Equal(5.00))
			})
		})
	})
})
