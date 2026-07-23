package util_test

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"

	"github.com/amanhigh/go-fun/common/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xuri/excelize/v2"
)

var _ = Describe("Excel", func() {
	var (
		f         *excelize.File
		sheetName string
	)

	BeforeEach(func() {
		f = excelize.NewFile()
		sheetName = f.GetSheetName(0) // Default "Sheet1"
	})

	AfterEach(func() {
		f.Close()
	})

	Context("WriteRow", func() {
		It("should write headers and data rows correctly", func() {
			headers := []any{"Symbol", "Price", "Quantity"}
			err := util.WriteRow(f, sheetName, 1, headers)
			Expect(err).ToNot(HaveOccurred())

			data := []any{"AAPL", 150.25, 10}
			err = util.WriteRow(f, sheetName, 2, data)
			Expect(err).ToNot(HaveOccurred())

			// Verify header cells
			symbol, err := f.GetCellValue(sheetName, "A1")
			Expect(err).ToNot(HaveOccurred())
			Expect(symbol).To(Equal("Symbol"))

			price, err := f.GetCellValue(sheetName, "B1")
			Expect(err).ToNot(HaveOccurred())
			Expect(price).To(Equal("Price"))

			// Verify data cells
			val, err := f.GetCellValue(sheetName, "A2")
			Expect(err).ToNot(HaveOccurred())
			Expect(val).To(Equal("AAPL"))

			price2, err := f.GetCellValue(sheetName, "B2")
			Expect(err).ToNot(HaveOccurred())
			Expect(price2).To(Equal("150.25"))
		})

		It("should handle empty data slice", func() {
			err := util.WriteRow(f, sheetName, 1, []any{})
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("WriteFormulaCell", func() {
		BeforeEach(func() {
			// Write values that formulas will reference
			err := util.WriteRow(f, sheetName, 1, []any{"ColumnA", "ColumnB"})
			Expect(err).ToNot(HaveOccurred())
			err = util.WriteRow(f, sheetName, 2, []any{10, 20})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should write a valid formula and calculate correctly", func() {
			err := util.WriteFormulaCell(f, sheetName, 2, 3, "=A2+B2")
			Expect(err).ToNot(HaveOccurred())

			// Verify formula string
			formula, err := f.GetCellFormula(sheetName, "C2")
			Expect(err).ToNot(HaveOccurred())
			Expect(formula).To(Equal("=A2+B2"))

			// Verify calculated value
			calculated, err := f.CalcCellValue(sheetName, "C2")
			Expect(err).ToNot(HaveOccurred())
			Expect(calculated).To(Equal("30"))
		})
	})

	Context("WriteFormulaRange", func() {
		BeforeEach(func() {
			err := util.WriteRow(f, sheetName, 2, []any{100, 2, 50})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should write multiple formulas and calculate correctly", func() {
			formulas := map[int]string{
				4: "=A2*B2", // D2 = 100 * 2 = 200
				5: "=C2*B2", // E2 = 50 * 2 = 100
			}
			err := util.WriteFormulaRange(f, sheetName, 2, formulas)
			Expect(err).ToNot(HaveOccurred())

			// Verify formula strings
			formula1, err := f.GetCellFormula(sheetName, "D2")
			Expect(err).ToNot(HaveOccurred())
			Expect(formula1).To(Equal("=A2*B2"))

			formula2, err := f.GetCellFormula(sheetName, "E2")
			Expect(err).ToNot(HaveOccurred())
			Expect(formula2).To(Equal("=C2*B2"))

			// Verify calculated values
			val1, err := f.CalcCellValue(sheetName, "D2")
			Expect(err).ToNot(HaveOccurred())
			Expect(val1).To(Equal("200"))

			val2, err := f.CalcCellValue(sheetName, "E2")
			Expect(err).ToNot(HaveOccurred())
			Expect(val2).To(Equal("100"))
		})
	})

	Context("ApplyAutoFilter", func() {
		BeforeEach(func() {
			// Write headers and a data row
			err := util.WriteRow(f, sheetName, 1, []any{"Symbol", "Price", "Quantity"})
			Expect(err).ToNot(HaveOccurred())
			err = util.WriteRow(f, sheetName, 2, []any{"AAPL", 150.25, 10})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should apply AutoFilter through the last data row", func() {
			err := util.ApplyAutoFilter(f, sheetName, "C", 2)
			Expect(err).ToNot(HaveOccurred())

			// Serialize the workbook and inspect the worksheet XML for the persisted autoFilter range
			var buf bytes.Buffer
			err = f.Write(&buf)
			Expect(err).ToNot(HaveOccurred())

			// Open as ZIP
			zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
			Expect(err).ToNot(HaveOccurred())

			// Find sheet1.xml (default sheet)
			var sheetXML []byte
			for _, zf := range zipReader.File {
				if zf.Name == "xl/worksheets/sheet1.xml" {
					rc, openErr := zf.Open()
					Expect(openErr).ToNot(HaveOccurred())
					sheetXML, err = io.ReadAll(rc)
					rc.Close()
					Expect(err).ToNot(HaveOccurred())
					break
				}
			}
			Expect(sheetXML).ToNot(BeNil(), "sheet1.xml not found in workbook")

			// Parse autoFilter ref from worksheet XML
			type autoFilterElement struct {
				Ref string `xml:"ref,attr"`
			}
			type worksheet struct {
				AutoFilter *autoFilterElement `xml:"autoFilter"`
			}
			var ws worksheet
			err = xml.Unmarshal(sheetXML, &ws)
			Expect(err).ToNot(HaveOccurred())
			Expect(ws.AutoFilter).ToNot(BeNil(), "autoFilter element should exist in worksheet XML")
			Expect(ws.AutoFilter.Ref).To(Equal("$A$1:$C$2"))
		})

		It("should return error when lastColumn is empty", func() {
			err := util.ApplyAutoFilter(f, sheetName, "", 2)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid AutoFilter range"))
		})

		It("should return error when lastDataRow is less than 1", func() {
			err := util.ApplyAutoFilter(f, sheetName, "C", 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid AutoFilter range"))
		})
	})
})
