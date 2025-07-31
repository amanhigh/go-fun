package manager_test

import (
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

		// Create a dummy Excel file for testing
		f := excelize.NewFile()
		sheetName := "Income"
		_, err = f.NewSheet(sheetName)
		Expect(err).ToNot(HaveOccurred())

		// Add headers
		headers := []string{"Date", "Time (in UTC)", "Activity", "Ticker", "Gross Cash Amount (in USD)"}
		err = f.SetSheetRow(sheetName, "A1", &headers)
		Expect(err).ToNot(HaveOccurred())

		// Add sample data
		rows := [][]interface{}{
			{"2024-01-15", "10:00:00", "Dividend", "AAPL", "100.00"},
			{"2024-01-20", "11:00:00", "Interest", "CASH", "50.00"},
			{"2024-02-10", "12:00:00", "Tax", "AAPL", "-15.00"},
			{"2024-02-15", "13:00:00", "Interest", "CASH", "25.50"},
		}

		for i, rowData := range rows {
			err = f.SetSheetRow(sheetName, "A"+string(rune('2'+i)), &rowData)
			Expect(err).ToNot(HaveOccurred())
		}

		//Remove Default Sheet
		f.DeleteSheet("Sheet1")

		err = f.SaveAs(sampleExcelPath)
		Expect(err).ToNot(HaveOccurred())

		driveWealthManager = manager.NewDriveWealthManager(sampleExcelPath)
	})

	AfterEach(func() {
		err := os.RemoveAll(tempTestDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when parsing a valid Excel file", func() {
		It("should extract interest entries correctly", func() {
			interestEntries, err := driveWealthManager.Parse()
			Expect(err).ToNot(HaveOccurred())
			Expect(interestEntries).To(HaveLen(2))

			// Verify first interest entry
			Expect(interestEntries[0].Symbol).To(Equal("CASH"))
			Expect(interestEntries[0].Date).To(Equal("2024-01-20"))
			Expect(interestEntries[0].Amount).To(Equal("50.00"))
			Expect(interestEntries[0].Tax).To(Equal("0"))
			Expect(interestEntries[0].Net).To(Equal("50.00"))

			// Verify second interest entry
			Expect(interestEntries[1].Symbol).To(Equal("CASH"))
			Expect(interestEntries[1].Date).To(Equal("2024-02-15"))
			Expect(interestEntries[1].Amount).To(Equal("25.50"))
			Expect(interestEntries[1].Tax).To(Equal("0"))
			Expect(interestEntries[1].Net).To(Equal("25.50"))
		})
	})

	Context("when the Excel file is missing", func() {
		It("should return an error", func() {
			nonExistentManager := manager.NewDriveWealthManager("non_existent_file.xlsx")
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
			f.DeleteSheet("Sheet1")
			err = f.SaveAs(sampleExcelPath)
			Expect(err).ToNot(HaveOccurred())

			driveWealthManager = manager.NewDriveWealthManager(sampleExcelPath)
			_, err = driveWealthManager.Parse()
			Expect(err).To(HaveOccurred())
		})
	})
})
