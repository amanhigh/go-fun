package manager_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xuri/excelize/v2"
)

var _ = Describe("DriveWealthManager", func() {
	var (
		tempTestDir        string
		sampleExcelPath    string
		driveWealthManager *manager.DriveWealthManager
	)

	BeforeEach(func() {
		var err error
		tempTestDir, err = os.MkdirTemp("", "drivewealth_test_*")
		Expect(err).ToNot(HaveOccurred())

		sampleExcelPath = filepath.Join(tempTestDir, "vested_transactions.xlsx")
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

			//Remove Default Sheet
			f.DeleteSheet("Sheet1")

			err = f.SaveAs(sampleExcelPath)
			Expect(err).ToNot(HaveOccurred())

			driveWealthManager = manager.NewDriveWealthManager(sampleExcelPath)
		})

		Context("when parsing interests", func() {
			It("should extract interest entries correctly", func() {
				interests, _, err := driveWealthManager.Parse()
				Expect(err).ToNot(HaveOccurred())
				Expect(interests).To(HaveLen(2))
				Expect(interests[0].Amount).To(BeNumerically("~", 0.59))
				Expect(interests[1].Amount).To(BeNumerically("~", 1.18))
			})
		})

		Context("when parsing dividends", func() {
			It("should extract dividend entries correctly", func() {
				_, dividends, err := driveWealthManager.Parse()
				Expect(err).ToNot(HaveOccurred())
				Expect(dividends).To(HaveLen(3))

				Expect(dividends[0].Symbol).To(Equal("IEF"))
				Expect(dividends[0].Amount).To(BeNumerically("~", 158.17))
				Expect(dividends[0].Tax).To(BeNumerically("~", 39.54))
				Expect(dividends[0].Net).To(BeNumerically("~", 118.63))

				Expect(dividends[1].Symbol).To(Equal("TLT"))
				Expect(dividends[1].Amount).To(BeNumerically("~", 3.51))
				Expect(dividends[1].Tax).To(BeNumerically("~", 0.0))
				Expect(dividends[1].Net).To(BeNumerically("~", 3.51))

				Expect(dividends[2].Symbol).To(Equal("BIL"))
				Expect(dividends[2].Amount).To(BeNumerically("~", 142.07))
				Expect(dividends[2].Tax).To(BeNumerically("~", 35.52))
				Expect(dividends[2].Net).To(BeNumerically("~", 106.55))
			})
		})
	})

	Context("with an invalid or malformed Excel file", func() {
		Context("when the Excel file is missing", func() {
			It("should return an error", func() {
				nonExistentManager := manager.NewDriveWealthManager("non_existent_file.xlsx")
				_, _, err := nonExistentManager.Parse()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the 'Income' sheet is missing", func() {
			It("should return an error", func() {
				// Create a new Excel file without the "Income" sheet
				f := excelize.NewFile()
				_, err := f.NewSheet("OtherSheet")
				Expect(err).ToNot(HaveOccurred())
				f.DeleteSheet("Sheet1")
				err = f.SaveAs(sampleExcelPath)
				Expect(err).ToNot(HaveOccurred())

				driveWealthManager = manager.NewDriveWealthManager(sampleExcelPath)
				_, _, err = driveWealthManager.Parse()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
