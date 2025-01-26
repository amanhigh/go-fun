package manager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	gocsv "github.com/gocarina/gocsv"
)

type CapitalGainsManager interface {
	AnalysePositions(ctx context.Context, year int) (map[string]tax.PositionAnalysis, error)
}

type CapitalGainsManagerImpl struct {
	tickerManager TickerManager
	downloadDir   string
	statementFile string
}

func NewCapitalGainsManager(
	tickerManager TickerManager,
	downloadDir string,
	statementFile string,
) *CapitalGainsManagerImpl {
	return &CapitalGainsManagerImpl{
		tickerManager: tickerManager,
		downloadDir:   downloadDir,
		statementFile: statementFile,
	}
}

func (c *CapitalGainsManagerImpl) AnalysePositions(ctx context.Context, year int) (map[string]tax.PositionAnalysis, error) {
	// TODO: Document Sample Formats for all CSV Files
	// Open and read CSV file
	filePath := filepath.Join(c.downloadDir, c.statementFile)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open broker statement: %w", err)
	}
	defer file.Close()

	// Parse CSV
	var transactions []tax.Transaction
	if err := gocsv.Unmarshal(file, &transactions); err != nil {
		return nil, fmt.Errorf("failed to parse broker statement: %w", err)
	}

	// Group by ticker
	tickerTransactions := make(map[string][]tax.Transaction)
	for _, t := range transactions {
		tickerTransactions[t.Security] = append(tickerTransactions[t.Security], t)
	}

	// Process each ticker
	result := make(map[string]tax.PositionAnalysis)
	for ticker, trades := range tickerTransactions {
		analysis, err := c.analyzeTickerPositions(ctx, ticker, trades, year)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze ticker %s: %w", ticker, err)
		}
		result[ticker] = analysis
	}

	return result, nil
}

func (c *CapitalGainsManagerImpl) analyzeTickerPositions(ctx context.Context, ticker string, transactions []tax.Transaction, year int) (tax.PositionAnalysis, error) {
	// Initialize analysis
	analysis := tax.PositionAnalysis{
		Ticker: ticker,
	}

	// Track running position
	var currentPosition float64
	var maxPosition float64
	var firstBuyDate time.Time
	// BUG: Track max Position Date
	// var maxPositionDate time.Time

	// Process transactions chronologically
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].DateAcquired < transactions[j].DateAcquired
	})

	for _, t := range transactions {
		// Parse dates
		acquiredDate, err := time.Parse(common.DateOnly, t.DateAcquired)
		if err != nil {
			return analysis, fmt.Errorf("invalid date format: %w", err)
		}

		// Update position
		currentPosition += t.QuantitySold // Note: Buy is positive, Sell is negative

		// Track first buy
		if firstBuyDate.IsZero() && t.QuantitySold > 0 {
			firstBuyDate = acquiredDate
			price, err := c.tickerManager.GetPrice(ctx, ticker, acquiredDate)
			if err != nil {
				return analysis, fmt.Errorf("failed to get price for first buy: %w", err)
			}
			analysis.FirstPosition = tax.Position{
				Date:     acquiredDate,
				Quantity: t.QuantitySold,
				USDPrice: price,
				USDValue: price * t.QuantitySold,
			}
		}

		// Track peak position with date
		if currentPosition > maxPosition {
			maxPosition = currentPosition
			price, err := c.tickerManager.GetPrice(ctx, ticker, acquiredDate)
			if err != nil {
				return analysis, fmt.Errorf("failed to get price for peak: %w", err)
			}
			analysis.PeakPosition = tax.Position{
				Date:     acquiredDate,
				Quantity: currentPosition,
				USDPrice: price,
				USDValue: price * currentPosition,
			}
		} else if currentPosition == maxPosition && acquiredDate.After(analysis.PeakPosition.Date) {
			// If we have the same position value but later date, update peak position
			price, err := c.tickerManager.GetPrice(ctx, ticker, acquiredDate)
			if err != nil {
				return analysis, fmt.Errorf("failed to get price for peak: %w", err)
			}
			analysis.PeakPosition = tax.Position{
				Date:     acquiredDate,
				Quantity: currentPosition,
				USDPrice: price,
				USDValue: price * currentPosition,
			}
		}
	}

	// Get year end position if any holdings exist
	if currentPosition > 0 {
		yearEndDate := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
		price, err := c.tickerManager.GetPrice(ctx, ticker, yearEndDate)
		if err != nil {
			return analysis, fmt.Errorf("failed to get year end price: %w", err)
		}
		analysis.YearEndPosition = tax.Position{
			Date:     yearEndDate,
			Quantity: currentPosition,
			USDPrice: price,
			USDValue: price * currentPosition,
		}
	}

	return analysis, nil
}
