package manager

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/tax"
)

// MatchDividendWithTax applies withholding tax from taxMap to dividend and calculates net.
// It matches dividend to tax by symbol+date, then removes the tax from the pool.
// If no matching tax is found, tax is set to 0 and net equals amount.
// The Net value is rounded to 2 decimals to avoid floating-point precision errors.
func MatchDividendWithTax(dividend *tax.Dividend, taxMap map[string]map[string]float64) {
	if dateTaxes, ok := taxMap[dividend.Symbol]; ok {
		if taxAmount, ok := dateTaxes[dividend.Date]; ok {
			dividend.Tax = taxAmount
			delete(dateTaxes, dividend.Date)
		}
	}
	dividend.Net = util.RoundToDecimals(dividend.Amount-dividend.Tax, 2)
}
