package manager

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/gocarina/gocsv"
)

type DividendManager interface {
	GetDividends(ctx context.Context) ([]tax.Dividend, error)
}

type DividendManagerImpl struct {
	sbiManager   SBIManager
	downloadsDir string
	dividendFile string
}

func NewDividendManager(sbiManager SBIManager, downloadsDir, dividendFile string) DividendManager {
	return &DividendManagerImpl{
		sbiManager:   sbiManager,
		downloadsDir: downloadsDir,
		dividendFile: dividendFile,
	}
}

func (d *DividendManagerImpl) GetDividends(_ context.Context) ([]tax.Dividend, error) {
	// Open CSV file
	file, err := os.Open(filepath.Join(d.downloadsDir, d.dividendFile))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse CSV rows
	var rows []tax.DividendBase
	if err := gocsv.UnmarshalFile(file, &rows); err != nil {
		return nil, err
	}

	// Process each dividend row
	var transactions []tax.Dividend
	for _, row := range rows {
		date, err := time.Parse(common.DateOnly, row.DividendDate)
		if err != nil {
			return nil, err
		}

		// Get TT rate for dividend date
		ttRate, err := d.sbiManager.GetTTBuyRate(date)
		if err != nil {
			return nil, err
		}

		// Create transaction with INR conversions
		transaction := tax.Dividend{
			DividendBase:   row,
			USDINRRate:     ttRate,
			NetDividendINR: row.NetDividend * ttRate,
			DividendTaxINR: row.DividendTax * ttRate,
		}
		transactions = append(transactions, transaction)
	}

	// TODO: Cache Transactions in Memory ?

	return transactions, nil
}
