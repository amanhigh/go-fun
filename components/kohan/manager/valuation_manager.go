package manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	repository "github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/samber/lo"
)

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
	sbiManager      SBIManager
}

func NewValuationManager(
	tickerManager TickerManager,
	accountManager AccountManager,
	tradeRepository repository.TradeRepository,
	fyManager FinancialYearManager[tax.Trade],
	sbiManager SBIManager,
) ValuationManager {
	return &ValuationManagerImpl{
		tickerManager:   tickerManager,
		accountManager:  accountManager,
		tradeRepository: tradeRepository,
		fyManager:       fyManager,
		sbiManager:      sbiManager,
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

	// Group filtered trades by Ticker Symbol
	tradesByTicker := lo.GroupBy(yearTrades, func(trade tax.Trade) string {
		return trade.Symbol
	})

	// TDD FIX: Get all tickers with carry-over positions from previous year
	// This ensures tickers without trades in target year are still included
	prevYearAccounts, accErr := v.accountManager.GetAllRecords(ctx, year-1)
	if accErr != nil && !errors.Is(accErr, common.ErrNotFound) {
		return nil, accErr
	}

	// Add carry-over tickers that don't have trades in target year
	for _, account := range prevYearAccounts {
		if account.Quantity > 0 && tradesByTicker[account.Symbol] == nil {
			tradesByTicker[account.Symbol] = []tax.Trade{} // Empty trades, will use carry-over position
		}
	}

	if len(tradesByTicker) == 0 {
		return []tax.Valuation{}, common.NewHttpError(fmt.Sprintf("no trades or carry-over positions found for year %d", year), http.StatusNotFound)
	}

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
		valuation, analyzeErr := v.AnalyzeValuation(ctx, ticker, tickerTrades, year)
		if analyzeErr != nil {
			// Fail fast: return immediately upon the first analysis error
			return nil, analyzeErr
		}

		valuations = append(valuations, valuation)
	}
	return valuations, nil
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
		// getOpeningPositionForPeriod returns (tax.Position{}, nil) for common.ErrNotFound (fresh start)
		// So, any non-nil err here is an actual error.
		return tax.Valuation{}, common.NewServerError(fmt.Errorf("failed to get opening position for %s: %w", tickerSymbol, err))
	}

	// Step 3: Validate if trades exist or if there's a carry-over
	if err := v.validateTradesExistOrCarryOver(trades, openingPosition, tickerSymbol); err != nil {
		return tax.Valuation{}, err
	}

	analysis := tax.Valuation{Ticker: tickerSymbol}
	analysis.FirstPosition = openingPosition
	analysis.PeakPosition = openingPosition

	// Step 4: Calculate daily peak value (Tax.md Line 124 compliance - MANDATORY)
	// Daily peak calculation is the authoritative method for determining peak INR value
	peakPosition, peakErr := v.calculateDailyPeak(ctx, tickerSymbol, year, openingPosition, trades)
	if peakErr != nil {
		return tax.Valuation{}, common.NewServerError(
			fmt.Errorf("failed to calculate daily peak for %s: %w", tickerSymbol, peakErr))
	}
	analysis.PeakPosition = peakPosition

	// Step 5: Determine current quantity at year end
	currentQuantity, processErr := v.processTrades(&analysis, trades, openingPosition)
	if processErr != nil {
		return tax.Valuation{}, processErr
	}

	// Step 6: Determine year end position
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
	// Use GetType() for normalized type comparison
	if openingPeriodPosition.Quantity == 0 && len(trades) > 0 && trades[0].GetType() == tax.TRADE_TYPE_SELL {
		return 0, common.NewHttpError(fmt.Sprintf("first trade can't be sell on fresh start for %s", analysis.Ticker), http.StatusBadRequest)
	}

	currentQuantity = openingPeriodPosition.Quantity
	// Peak will be calculated by calculateDailyPeak() using INR values

	for i, trade := range trades {
		if i == 0 && openingPeriodPosition.Quantity == 0 {
			currentQuantity, err = v.handleFirstTrade(analysis, trade)
			if err != nil {
				return 0, err
			}
			continue
		}

		currentQuantity = v.applyTrade(trade, currentQuantity)
	}
	return currentQuantity, nil
}

// handleFirstTrade handles the very first BUY trade in a fresh start scenario.
func (v *ValuationManagerImpl) handleFirstTrade(analysis *tax.Valuation, trade tax.Trade) (currentQuantity float64, err common.HttpError) {
	tradeDate, dateErr := trade.GetDate()
	if dateErr != nil {
		return 0, dateErr
	}

	// Use GetType() for normalized type comparison with constants
	// Real data has "Buy"/"Sell" from DriveWealth, GetType() normalizes to uppercase
	if trade.GetType() == tax.TRADE_TYPE_BUY {
		analysis.FirstPosition = tax.Position{
			Date:     tradeDate,
			Quantity: trade.Quantity,
			USDPrice: trade.USDPrice,
		}
		// Peak will be calculated by calculateDailyPeak() using INR values
		return trade.Quantity, nil
	}
	return 0, nil // Should not happen due to initial check in processTrades
}

// applyTrade processes subsequent trades or trades in a carry-over scenario.
func (v *ValuationManagerImpl) applyTrade(trade tax.Trade, currentQuantity float64) float64 {
	if trade.GetType() == tax.TRADE_TYPE_BUY {
		return currentQuantity + trade.Quantity
	}
	return currentQuantity - trade.Quantity
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
	// Last year's account record
	account, accErr := v.accountManager.GetRecord(ctx, ticker, year-1)
	if accErr != nil {
		if errors.Is(accErr, common.ErrNotFound) {
			// No account record found -> fresh start for this period.
			// Return a zero position. Its Date field will be the zero value for time.Time.
			return tax.Position{}, nil
		}
		return tax.Position{}, accErr // Other errors from accountManager
	}

	// Account record found (carry-over scenario)
	// Reconstruct FirstPosition from Account metadata (original acquisition date/price)
	// OriginDate MUST be present for carry-over accounts - it's required for tax reporting
	originDate, parseErr := time.Parse(time.DateOnly, account.OriginDate)
	if parseErr != nil {
		return tax.Position{}, tax.NewInvalidDateError(fmt.Sprintf("failed to parse OriginDate '%s' for carry-over account %s: %v", account.OriginDate, ticker, parseErr))
	}

	return tax.Position{
		Date:     originDate,
		Quantity: account.OriginQty,
		USDPrice: account.OriginPrice,
	}, nil
}

// calculateDailyPeak evaluates (Quantity × Market_Price × SBI_Rate) for every day in the year
// to find the true INR peak value during the calendar year.
// This ensures compliance with Tax.md Line 124 requirement for daily evaluation.
func (v *ValuationManagerImpl) calculateDailyPeak(
	ctx context.Context,
	ticker string,
	year int,
	openingPosition tax.Position,
	trades []tax.Trade,
) (peakPosition tax.Position, err common.HttpError) {
	// Step 1: Build daily quantity timeline by processing trades chronologically
	quantityByDate := v.buildDailyQuantityTimeline(year, openingPosition, trades)
	if len(quantityByDate) == 0 {
		// No holdings during the year
		return tax.Position{}, nil
	}

	// Step 2: Get daily market prices for the ticker
	dailyPrices, priceErr := v.tickerManager.GetDailyPrices(ctx, ticker, year)
	if priceErr != nil {
		return tax.Position{}, priceErr
	}

	// Step 3: Get daily SBI TT Buy rates for the year
	dailyRates, rateErr := v.sbiManager.GetDailyRates(ctx, year)
	if rateErr != nil {
		return tax.Position{}, rateErr
	}

	// Step 4: Find the date with maximum INR value by iterating through each day
	return v.findPeakByIteratingYear(year, openingPosition, quantityByDate, dailyPrices, dailyRates), nil
}

// findPeakByIteratingYear finds maximum INR value (Qty × Price × Rate) across the year (Tax.md Line 124).
func (v *ValuationManagerImpl) findPeakByIteratingYear(
	year int,
	openingPosition tax.Position,
	quantityByDate map[string]float64,
	dailyPrices map[string]float64,
	dailyRates map[string]float64,
) tax.Position {
	peakPos := openingPosition
	maxINRValue := 0.0
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)

	for currDate := startDate; !currDate.After(endDate); currDate = currDate.AddDate(0, 0, 1) {
		dateStr := currDate.Format(time.DateOnly)
		quantity := v.getClosestValue(quantityByDate, dateStr)
		if quantity == 0 {
			continue
		}
		price := v.getClosestValue(dailyPrices, dateStr)
		if price == 0 {
			continue
		}
		rate := v.getClosestValue(dailyRates, dateStr)
		if rate == 0 {
			continue
		}

		inrValue := quantity * price * rate
		if inrValue > maxINRValue {
			maxINRValue = inrValue
			peakPos = tax.Position{Date: currDate, Quantity: quantity, USDPrice: price}
		}
	}
	return peakPos
}

// buildDailyQuantityTimeline creates a map of date → quantity held by processing trades
func (v *ValuationManagerImpl) buildDailyQuantityTimeline(
	year int,
	openingPosition tax.Position,
	trades []tax.Trade,
) map[string]float64 {
	timeline := make(map[string]float64)
	currentQuantity := openingPosition.Quantity

	// Initialize the year with opening quantity
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	timeline[startDate.Format(time.DateOnly)] = currentQuantity

	// Process each trade chronologically
	for _, trade := range trades {
		tradeDate, dateErr := trade.GetDate()
		if dateErr != nil {
			continue // Skip trades with invalid dates
		}

		if tradeDate.Year() != year {
			continue // Skip trades outside the target year
		}

		// Update quantity based on trade type
		if trade.GetType() == tax.TRADE_TYPE_BUY {
			currentQuantity += trade.Quantity
		} else { // SELL
			currentQuantity -= trade.Quantity
		}

		timeline[tradeDate.Format(time.DateOnly)] = currentQuantity
	}

	return timeline
}

// getClosestValue finds the nearest previous value for a given date using backfill logic.
// If the exact date exists in the map, it returns that value immediately.
// Otherwise, it searches for the closest previous date with available data.
// Returns 0 if no previous data is found.
func (v *ValuationManagerImpl) getClosestValue(dataMap map[string]float64, dateStr string) float64 {
	if value, exists := dataMap[dateStr]; exists {
		return value
	}

	// Backfill: find the closest previous date's value
	parsedDate, _ := time.Parse(time.DateOnly, dateStr)
	var closestValue float64
	var closestDate time.Time

	for dateKey, value := range dataMap {
		keyDate, _ := time.Parse(time.DateOnly, dateKey)
		if !keyDate.After(parsedDate) && (closestDate.IsZero() || keyDate.After(closestDate)) {
			closestValue = value
			closestDate = keyDate
		}
	}

	return closestValue
}
