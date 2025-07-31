package manager_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xuri/excelize/v2"
)

var _ = Describe("DriveWealthManager", func() {
	var (
		tempTestDir        string
		sampleExcelPath    string
		driveWealthManager manager.DriveWealthManager
		taxConfig          config.TaxConfig
	)

	BeforeEach(func() {
		var err error
		tempTestDir, err = os.MkdirTemp("", "drivewealth_test_*")
		Expect(err).ToNot(HaveOccurred())

		sampleExcelPath = filepath.Join(tempTestDir, "vested_transactions.xlsx")
		taxConfig = config.TaxConfig{
			InterestFilePath: filepath.Join(tempTestDir, "interest.csv"),
			TradesPath:       filepath.Join(tempTestDir, "trades.csv"),
			DividendFilePath: filepath.Join(tempTestDir, "dividends.csv"),
			DriveWealthPath:  sampleExcelPath,
		}
	})

	AfterEach(func() {
		err := os.RemoveAll(tempTestDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("with a valid Excel file", func() {
		BeforeEach(func() {
			// Create a dummy Excel file for testing
			f := excelize.NewFile()
			sheetName := "Income"
			_, err := f.NewSheet(sheetName)
			Expect(err).ToNot(HaveOccurred())

			// Add headers
			headers := []string{"Date", "Time (in UTC)", "Activity", "Ticker", "Gross Cash Amount (in USD)"}
			err = f.SetSheetRow(sheetName, "A1", &headers)
			Expect(err).ToNot(HaveOccurred())

			// Add sample data
			rows := [][]interface{}{
				{"2025-06-06", "05:24:40 AM", "Dividend", "IEF", 158.17},
				{"2025-06-06", "05:24:39 AM", "Tax", "IEF", -39.54},
				{"2025-06-06", "05:24:37 AM", "Dividend", "TLT", 3.51},
				{"2025-06-06", "05:24:37 AM", "Tax", "BIL", -35.52},
				{"2025-06-06", "05:24:36 AM", "Dividend", "BIL", 142.07},
				{"2025-06-03", "05:34:52 AM", "Interest", "", 0.59},
				{"2025-05-02", "04:57:05 AM", "Interest", "", 1.18},
			}

			for i, rowData := range rows {
				err = f.SetSheetRow(sheetName, fmt.Sprintf("A%d", i+2), &rowData)
				Expect(err).ToNot(HaveOccurred())
			}

			// Remove Default Sheet
			Expect(f.DeleteSheet("Sheet1")).To(Succeed())

			/* Create Trades Sheet */
			tradeSheet := "Trades"
			_, err = f.NewSheet(tradeSheet)
			Expect(err).ToNot(HaveOccurred())

			tradeHeaders := []string{"Date", "Time (in UTC)", "Name", "Ticker", "Activity", "Order Type", "Quantity", "Price Per Share (in USD)", "Cash Amount (in USD)", "Commission Charges (in USD)"}
			err = f.SetSheetRow(tradeSheet, "A1", &tradeHeaders)
			Expect(err).ToNot(HaveOccurred())

			tradeRows := [][]interface{}{
				{"2025-04-03", "04:53:52 PM", "Vanguard Russell 2000 ETF", "VTWO", "Buy", "Market", 70, 77.41, 5418.7, 0},
				{"2025-04-03", "04:26:02 PM", "Europe ETF FTSE Vanguard", "VGK", "Buy", "Market", 60, 70.18, 4210.8, 0},
				{"2025-02-14", "02:30:00 PM", "Barclays 1-3 Month T-Bill ETF SPDR", "BIL", "Sell", "Market", 1, 91.58, 91.58, 0},
			}

			for i, rowData := range tradeRows {
				err = f.SetSheetRow(tradeSheet, fmt.Sprintf("A%d", i+2), &rowData)
				Expect(err).ToNot(HaveOccurred())
			}

			err = f.SaveAs(sampleExcelPath)
			Expect(err).ToNot(HaveOccurred())

			driveWealthManager = manager.NewDriveWealthManager(taxConfig)
		})

		Context("when parsing interests", func() {
			It("should extract interest entries correctly", func() {
				info, err := driveWealthManager.Parse()
				Expect(err).ToNot(HaveOccurred())
				Expect(info.Interests).To(HaveLen(2))
				Expect(info.Interests[0].Amount).To(BeNumerically("~", 0.59))
				Expect(info.Interests[1].Amount).To(BeNumerically("~", 1.18))
			})
		})

		Context("when parsing dividends", func() {
			It("should extract dividend entries correctly", func() {
				info, err := driveWealthManager.Parse()
				Expect(err).ToNot(HaveOccurred())
				Expect(info.Dividends).To(HaveLen(3))

				Expect(info.Dividends[0].Symbol).To(Equal("IEF"))
				Expect(info.Dividends[0].Amount).To(BeNumerically("~", 158.17))
				Expect(info.Dividends[0].Tax).To(BeNumerically("~", 39.54))
				Expect(info.Dividends[0].Net).To(BeNumerically("~", 118.63))

				Expect(info.Dividends[1].Symbol).To(Equal("TLT"))
				Expect(info.Dividends[1].Amount).To(BeNumerically("~", 3.51))
				Expect(info.Dividends[1].Tax).To(BeNumerically("~", 0.0))
				Expect(info.Dividends[1].Net).To(BeNumerically("~", 3.51))

				Expect(info.Dividends[2].Symbol).To(Equal("BIL"))
				Expect(info.Dividends[2].Amount).To(BeNumerically("~", 142.07))
				Expect(info.Dividends[2].Tax).To(BeNumerically("~", 35.52))
				Expect(info.Dividends[2].Net).To(BeNumerically("~", 106.55))
			})
		})

		Context("when parsing trades", func() {
			It("should extract trade entries correctly", func() {
				info, err := driveWealthManager.Parse()
				Expect(err).ToNot(HaveOccurred())
				Expect(info.Trades).To(HaveLen(3))

				Expect(info.Trades[0].Symbol).To(Equal("VTWO"))
				Expect(info.Trades[0].Quantity).To(BeNumerically("~", 70))
				Expect(info.Trades[0].USDPrice).To(BeNumerically("~", 77.41))
				Expect(info.Trades[0].USDValue).To(BeNumerically("~", 5418.7))
				Expect(info.Trades[0].Type).To(Equal("Buy"))

				Expect(info.Trades[1].Symbol).To(Equal("VGK"))
				Expect(info.Trades[1].Quantity).To(BeNumerically("~", 60))
				Expect(info.Trades[1].USDPrice).To(BeNumerically("~", 70.18))
				Expect(info.Trades[1].USDValue).To(BeNumerically("~", 4210.8))
				Expect(info.Trades[1].Type).To(Equal("Buy"))

				Expect(info.Trades[2].Symbol).To(Equal("BIL"))
				Expect(info.Trades[2].Quantity).To(BeNumerically("~", 1))
				Expect(info.Trades[2].USDPrice).To(BeNumerically("~", 91.58))
				Expect(info.Trades[2].USDValue).To(BeNumerically("~", 91.58))
				Expect(info.Trades[2].Type).To(Equal("Sell"))
			})
		})

		Context("when generating CSV", func() {
			It("should create valid csv files", func() {
				info := tax.DriveWealthInfo{
					Interests: []tax.Interest{
						{Symbol: "CASH", Date: "2025-06-03", Amount: 0.59, Tax: 0, Net: 0.59},
						{Symbol: "CASH", Date: "2025-05-02", Amount: 1.18, Tax: 0, Net: 1.18},
					},
					Trades: []tax.Trade{
						{Symbol: "VTWO", Date: "2025-04-03", Type: "Buy", Quantity: 70, USDPrice: 77.41, USDValue: 5418.7, Commission: 0},
					},
					Dividends: []tax.Dividend{
						{Symbol: "IEF", Date: "2025-06-06", Amount: 158.17, Tax: 39.54, Net: 118.63},
					},
				}

				err := driveWealthManager.GenerateCsv(info)
				Expect(err).ToNot(HaveOccurred())

				// Verify Interest file content
				data, err := os.ReadFile(taxConfig.InterestFilePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("CASH,2025-06-03,0.59,0,0.59"))
				Expect(string(data)).To(ContainSubstring("CASH,2025-05-02,1.18,0,1.18"))

				// Verify Trade file content
				data, err = os.ReadFile(taxConfig.TradesPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("VTWO,2025-04-03,Buy,70,77.41,5418.7,0"))

				// Verify Dividend file content
				data, err = os.ReadFile(taxConfig.DividendFilePath)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(data)).To(ContainSubstring("IEF,2025-06-06,158.17,39.54,118.63"))
			})
		})
	})

	Context("with an invalid or malformed Excel file", func() {
		Context("when the Excel file is missing", func() {
			It("should return an error", func() {
				// This test works because the top-level BeforeEach sets up the path
				// but does not create the file. The manager is initialized with a
				// path that points to a non-existent file, so Parse should fail.
				nonExistentManager := manager.NewDriveWealthManager(taxConfig)
				_, err := nonExistentManager.Parse()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the 'Income' sheet is missing", func() {
			It("should return an error", func() {
				// Create a new Excel file without the "Income" sheet
				f := excelize.NewFile()
				_, err := f.NewSheet("OtherSheet")
				Expect(err).ToNot(HaveOccurred())
				Expect(f.DeleteSheet("Sheet1")).To(Succeed())
				err = f.SaveAs(sampleExcelPath)
				Expect(err).ToNot(HaveOccurred())

				driveWealthManager = manager.NewDriveWealthManager(taxConfig)
				_, err = driveWealthManager.Parse()
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when the 'Trades' sheet is missing", func() {
			It("should return an error", func() {
				// Create a new Excel file without the "Trades" sheet
				f := excelize.NewFile()
				_, err := f.NewSheet("Income")
				Expect(err).ToNot(HaveOccurred())
				Expect(f.DeleteSheet("Sheet1")).To(Succeed())
				err = f.SaveAs(sampleExcelPath)
				Expect(err).ToNot(HaveOccurred())

				driveWealthManager = manager.NewDriveWealthManager(taxConfig)
				_, err = driveWealthManager.Parse()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
