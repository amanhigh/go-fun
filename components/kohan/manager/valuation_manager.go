package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/samber/lo"
)

//go:generate mockery --name=ValuationManager
type ValuationManager interface {
	AnalyzeValuation(ctx context.Context, tickerSymbol string, trades []tax.Trade, year int) (tax.Valuation, common.HttpError)
	// GetYearlyValuationsUSD calculates the base USD Valuation (First, Peak, YearEnd)
	// for all relevant tickers based on trade history up to the end of the specified calendar year.
	GetYearlyValuationsUSD(ctx context.Context, year int) ([]tax.Valuation, common.HttpError)
}

type ValuationManagerImpl struct {
	tickerManager   TickerManager
	accountManager  AccountManager
	tradeRepository repository.TradeRepository
	fyManager       FinancialYearManager[tax.Trade]
}

func NewValuationManager(
	tickerManager TickerManager,
	accountManager AccountManager,
	tradeRepository repository.TradeRepository,
	fyManager FinancialYearManager[tax.Trade],
) ValuationManager {
	return &ValuationManagerImpl{
		tickerManager:   tickerManager,
		accountManager:  accountManager,
		tradeRepository: tradeRepository,
		fyManager:       fyManager,
	}
}

func (v *ValuationManagerImpl) GetYearlyValuationsUSD(ctx context.Context, year int) (valuations []tax.Valuation, err common.HttpError) {
	allTrades, repoErr := v.tradeRepository.GetAllRecords(ctx)
	if repoErr != nil {
		return nil, repoErr // Return other errors
	}

	// Filter trades for the specified US financial year (calendar year) and sort them
	yearTrades, filterErr := v.fyManager.FilterUS(ctx, allTrades, year)
	if filterErr != nil {
		return nil, filterErr
	}
	if len(yearTrades) == 0 {
		return []tax.Valuation{}, common.NewHttpError(fmt.Sprintf("no trades found for year %d", year), http.StatusNotFound)
	}

	// Group filtered trades by Ticker Symbol
	tradesByTicker := lo.GroupBy(yearTrades, func(trade tax.Trade) string {
		return trade.Symbol
	})

	// Process trades for all tickers using the helper function
	valuations, err = v.processTradesByTicker(ctx, tradesByTicker, year)
	if err != nil {
		return nil, err
	}

	return valuations, nil // Return aggregated results if all analyses succeeded
}

// processTradesByTicker iterates through tickers, sorts their trades, and processes them.
func (v *ValuationManagerImpl) processTradesByTicker(ctx context.Context, tradesByTicker map[string][]tax.Trade, year int) ([]tax.Valuation, common.HttpError) {
	valuations := make([]tax.Valuation, 0, len(tradesByTicker))

	// Get tickers for processing (Sort Order Helps in Tests)
	tickers := lo.Keys(tradesByTicker)
	slices.Sort(tickers)

	// Iterate through tickers
	for _, ticker := range tickers {
		tickerTrades := tradesByTicker[ticker]

		// Process trades for the current ticker (trades are assumed sorted by FilterUS)
		valuation, processErr := v.processTickerTrades(ctx, ticker, tickerTrades, year)
		if processErr != nil {
			// Fail fast: return immediately upon the first analysis error
			return nil, processErr
		}
		valuations = append(valuations, valuation)
	}
	return valuations, nil
}

// processTickerTrades analyzes the valuation for a single ticker's sorted trades.
func (v *ValuationManagerImpl) processTickerTrades(ctx context.Context, tickerSymbol string, sortedTrades []tax.Trade, year int) (tax.Valuation, common.HttpError) {
	// Call AnalyzeValuation with the known tickerSymbol
	valuation, analyzeErr := v.AnalyzeValuation(ctx, tickerSymbol, sortedTrades, year)
	if analyzeErr != nil {
		return tax.Valuation{}, analyzeErr
	}
	return valuation, nil
}

// AnalyzeValuation calculates valuation based on trades and opening position for a given ticker.
func (v *ValuationManagerImpl) AnalyzeValuation(ctx context.Context, tickerSymbol string, trades []tax.Trade, year int) (tax.Valuation, common.HttpError) {
	// Step 1: Validate trade symbols first. This is crucial for the "Multiple Ticker Trades" test
	// to fail before attempting to get an opening position if symbols are inconsistent.
	if err := v.validateTradeSymbols(trades, tickerSymbol); err != nil {
		return tax.Valuation{}, err
	}

	// Step 2: Get opening position
	openingPosition, err := v.getOpeningPositionForPeriod(ctx, tickerSymbol, year)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			// Fresh start - create zero position with valid date
			openingPosition = tax.Position{
				Date:     time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
				Quantity: 0,
				USDPrice: 0,
			}
		} else {
			return tax.Valuation{}, common.NewServerError(fmt.Errorf("failed to get opening position for %s: %w", tickerSymbol, err))
		}
	}

	// Step 3: Validate if trades exist or if there's a carry-over
	if err := v.validateTradesExistOrCarryOver(trades, openingPosition, tickerSymbol); err != nil {
		return tax.Valuation{}, err
	}

	analysis := tax.Valuation{Ticker: tickerSymbol}
	analysis.FirstPosition = openingPosition
	analysis.PeakPosition = openingPosition

	currentQuantity, processErr := v.processTrades(&analysis, trades, openingPosition)
	if processErr != nil {
		return tax.Valuation{}, processErr
	}

	if detErr := v.determineYearEndPosition(ctx, &analysis, year, currentQuantity); detErr != nil {
		return tax.Valuation{}, detErr
	}

	return analysis, nil
}

// processTrades updates analysis based on trades and opening position.
// It returns the final currentQuantity, and any error.
func (v *ValuationManagerImpl) processTrades(
	analysis *tax.Valuation,
	trades []tax.Trade,
	openingPeriodPosition tax.Position,
) (currentQuantity float64, err common.HttpError) {
	if openingPeriodPosition.Quantity == 0 && len(trades) > 0 && trades[0].Type == "SELL" {
		return 0, common.NewHttpError(fmt.Sprintf("first trade can't be sell on fresh start for %s", analysis.Ticker), http.StatusBadRequest)
	}

	currentQuantity = openingPeriodPosition.Quantity
	if openingPeriodPosition.Quantity > 0 {
		analysis.PeakPosition = openingPeriodPosition
	}

	for i, trade := range trades {
		if i == 0 && openingPeriodPosition.Quantity == 0 {
			currentQuantity, err = v.handleFirstTrade(analysis, trade)
			if err != nil {
				return 0, err
			}
			continue
		}

		currentQuantity, err = v.applyTrade(analysis, trade, currentQuantity)
		if err != nil {
			return 0, err
		}
	}
	return currentQuantity, nil
}

// handleFirstTrade handles the very first BUY trade in a fresh start scenario.
func (v *ValuationManagerImpl) handleFirstTrade(analysis *tax.Valuation, trade tax.Trade) (currentQuantity float64, err common.HttpError) {
	tradeDate, dateErr := trade.GetDate()
	if dateErr != nil {
		return 0, dateErr
	}

	if trade.Type == "BUY" {
		analysis.FirstPosition = tax.Position{
			Date:     tradeDate,
			Quantity: trade.Quantity,
			USDPrice: trade.USDPrice,
		}
		analysis.PeakPosition = analysis.FirstPosition // Initial peak is the first buy
		return trade.Quantity, nil
	}
	return 0, nil // Should not happen due to initial check in processTrades
}

// applyTrade processes subsequent trades or trades in a carry-over scenario.
func (v *ValuationManagerImpl) applyTrade(analysis *tax.Valuation, trade tax.Trade, currentQuantity float64) (float64, common.HttpError) {
	tradeDate, dateErr := trade.GetDate()
	if dateErr != nil {
		return 0, dateErr
	}

	if trade.Type == "BUY" {
		currentQuantity += trade.Quantity
		if currentQuantity > analysis.PeakPosition.Quantity {
			analysis.PeakPosition = tax.Position{
				Date:     tradeDate,
				Quantity: currentQuantity,
				USDPrice: trade.USDPrice,
			}
		}
	} else { // SELL
		currentQuantity -= trade.Quantity
	}
	return currentQuantity, nil
}

// determineYearEndPosition sets the YearEndPosition in the analysis.
func (v *ValuationManagerImpl) determineYearEndPosition(
	ctx context.Context,
	analysis *tax.Valuation,
	year int,
	currentQuantity float64,
) common.HttpError {
	yearEndDate := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
	switch {
	case currentQuantity > 0:
		price, priceErr := v.tickerManager.GetPrice(ctx, analysis.Ticker, yearEndDate)
		if priceErr != nil {
			return common.NewServerError(fmt.Errorf("failed to get year end price for %s: %w", analysis.Ticker, priceErr))
		}
		analysis.YearEndPosition = tax.Position{
			Date:     yearEndDate,
			Quantity: currentQuantity,
			USDPrice: price,
		}
	default:
		analysis.YearEndPosition = tax.Position{Date: yearEndDate}
	}
	return nil
}

// validateTradeSymbols checks if the trades (if any) have symbols consistent with the expectedTicker.
func (v *ValuationManagerImpl) validateTradeSymbols(trades []tax.Trade, expectedTicker string) common.HttpError {
	if expectedTicker == "" {
		return common.NewHttpError("expected ticker symbol cannot be empty", http.StatusBadRequest)
	}
	for _, t := range trades {
		if t.Symbol != expectedTicker {
			// Assuming t.Date is a string or has a String() method.
			return common.NewHttpError(fmt.Sprintf("trade symbol mismatch: expected %s but found %s in trade dated %s", expectedTicker, t.Symbol, t.Date), http.StatusBadRequest)
		}
	}
	return nil
}

// validateTradesExistOrCarryOver checks if there are trades or a carry-over position.
func (v *ValuationManagerImpl) validateTradesExistOrCarryOver(trades []tax.Trade, openingPosition tax.Position, expectedTicker string) common.HttpError {
	if len(trades) == 0 && openingPosition.Quantity == 0 {
		return common.NewHttpError(fmt.Sprintf("no trades or carry-over position provided for ticker %s", expectedTicker), http.StatusBadRequest)
	}
	return nil
}

func (v *ValuationManagerImpl) getOpeningPositionForPeriod(ctx context.Context, ticker string, year int) (position tax.Position, err common.HttpError) {
	// Smart account detection for previous year
	account, accErr := v.accountManager.GetRecord(ctx, ticker, year)
	if accErr != nil {
		if errors.Is(accErr, common.ErrNotFound) {
			// Fresh start - no previous year account
			return tax.Position{}, common.ErrNotFound
		}
		return tax.Position{}, accErr // Other errors from accountManager
	}

	// Account record found (carry-over scenario)
	// Use December 31st of previous year for carry-over positions
	openingDate := time.Date(year-1, 12, 31, 23, 59, 59, 0, time.UTC)
	var openingPrice float64
	if account.Quantity > 0 { // Avoid division by zero
		openingPrice = account.MarketValue / account.Quantity
	} else {
		openingPrice = 0 // If prior year quantity was zero, opening price is zero
	}

	return tax.Position{
		Date:     openingDate,
		Quantity: account.Quantity,
		USDPrice: openingPrice,
	}, nil
}
