package manager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/amanhigh/go-fun/models/fa"
	gocsv "github.com/gocarina/gocsv"
)

type CapitalGainsManager interface {
	ProcessBrokerStatement(ctx context.Context, year int) (map[string]fa.PositionAnalysis, error)
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

func (c *CapitalGainsManagerImpl) ProcessBrokerStatement(ctx context.Context, year int) (map[string]fa.PositionAnalysis, error) {
	// Open and read CSV file
	filePath := filepath.Join(c.downloadDir, c.statementFile)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open broker statement: %w", err)
	}
	defer file.Close()

	// Parse CSV
	var transactions []fa.Transaction
	if err := gocsv.Unmarshal(file, &transactions); err != nil {
		return nil, fmt.Errorf("failed to parse broker statement: %w", err)
	}

	// Group by ticker
	tickerTransactions := make(map[string][]fa.Transaction)
	for _, t := range transactions {
		tickerTransactions[t.Security] = append(tickerTransactions[t.Security], t)
	}

	// Process each ticker
	result := make(map[string]fa.PositionAnalysis)
	for ticker, trades := range tickerTransactions {
		analysis, err := c.analyzeTickerPositions(ctx, ticker, trades, year)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze ticker %s: %w", ticker, err)
		}
		result[ticker] = analysis
	}

	return result, nil
}

func (c *CapitalGainsManagerImpl) analyzeTickerPositions(ctx context.Context, ticker string, transactions []fa.Transaction, year int) (fa.PositionAnalysis, error) {
	// Initialize analysis
	analysis := fa.PositionAnalysis{
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
		acquiredDate, err := time.Parse("2006-01-02", t.DateAcquired)
		if err != nil {
			return analysis, fmt.Errorf("invalid date format: %w", err)
		}

		// Update position
		currentPosition += t.QuantitySold // Note: Buy is positive, Sell is negative

		// Track first buy
		if firstBuyDate.IsZero() && t.QuantitySold > 0 {
			firstBuyDate = acquiredDate
			price, err := c.tickerManager.GetPriceOnDate(ctx, ticker, acquiredDate)
			if err != nil {
				return analysis, fmt.Errorf("failed to get price for first buy: %w", err)
			}
			analysis.FirstPosition = fa.Position{
				Date:     acquiredDate,
				Quantity: t.QuantitySold,
				USDPrice: price,
				USDValue: price * t.QuantitySold,
			}
		}

		// Track peak position
		if currentPosition > maxPosition {
			maxPosition = currentPosition
			price, err := c.tickerManager.GetPriceOnDate(ctx, ticker, acquiredDate)
			if err != nil {
				return analysis, fmt.Errorf("failed to get price for peak: %w", err)
			}
			analysis.PeakPosition = fa.Position{
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
		price, err := c.tickerManager.GetPriceOnDate(ctx, ticker, yearEndDate)
		if err != nil {
			return analysis, fmt.Errorf("failed to get year end price: %w", err)
		}
		analysis.YearEndPosition = fa.Position{
			Date:     yearEndDate,
			Quantity: currentPosition,
			USDPrice: price,
			USDValue: price * currentPosition,
		}
	}

	return analysis, nil
}

func (c *CapitalGainsManagerImpl) GetPriceOnDate(ctx context.Context, ticker string, date time.Time) (float64, error) {
	// Use ticker manager to get price for the given date
	analysis, err := c.tickerManager.AnalyzeTicker(ctx, ticker, date.Year())
	if err != nil {
		return 0, err
	}

	// Use appropriate price based on date
	dateStr := date.Format("2006-01-02")
	if dateStr == analysis.YearEndDate {
		return analysis.YearEndPrice, nil
	}
	return analysis.PeakPrice, nil // Simplified for now
}
