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
func (v *ValuationManagerImpl) applyTrade(_ *tax.Valuation, trade tax.Trade, currentQuantity float64) (float64, common.HttpError) {
	// Use GetType() for normalized type comparison with constants
	if trade.GetType() == tax.TRADE_TYPE_BUY {
		currentQuantity += trade.Quantity
		// Peak calculation happens in calculateDailyPeak() which uses INR values
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
	// Opening position date should be Dec 31st of previous year for carry-over positions
	openingDate := time.Date(year-1, 12, 31, 0, 0, 0, 0, time.UTC)
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

	// Step 4: Find the date with maximum INR value
	return v.findDailyPeakPosition(year, openingPosition, quantityByDate, dailyPrices, dailyRates), nil
}

// peakCalculationData holds data needed for daily peak calculation
type peakCalculationData struct {
	quantityByDate map[string]float64
	dailyPrices    map[string]float64
	dailyRates     map[string]float64
	maxINRValue    float64
}

// findDailyPeakPosition iterates through each day of the year to find the maximum INR value
func (v *ValuationManagerImpl) findDailyPeakPosition(
	year int,
	openingPosition tax.Position,
	quantityByDate map[string]float64,
	dailyPrices map[string]float64,
	dailyRates map[string]float64,
) tax.Position {
	data := peakCalculationData{
		quantityByDate: quantityByDate,
		dailyPrices:    dailyPrices,
		dailyRates:     dailyRates,
		maxINRValue:    0,
	}

	peakPosition := openingPosition
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)

	for currDate := startDate; !currDate.After(endDate); currDate = currDate.AddDate(0, 0, 1) {
		dateStr := currDate.Format(time.DateOnly)
		peakPosition = v.updatePeakIfHigher(peakPosition, currDate, dateStr, &data)
	}

	return peakPosition
}

// updatePeakIfHigher checks if the current date's INR value is higher than the current peak
func (v *ValuationManagerImpl) updatePeakIfHigher(
	peakPos tax.Position,
	currDate time.Time,
	dateStr string,
	data *peakCalculationData,
) tax.Position {
	quantity := v.getQuantityForDate(data.quantityByDate, dateStr)
	if quantity == 0 {
		return peakPos
	}

	price := v.getClosestPrice(data.dailyPrices, dateStr)
	if price == 0 {
		return peakPos
	}

	rate := v.getClosestRate(data.dailyRates, dateStr)
	if rate == 0 {
		return peakPos
	}

	inrValue := quantity * price * rate
	if inrValue > data.maxINRValue {
		data.maxINRValue = inrValue
		return tax.Position{
			Date:     currDate,
			Quantity: quantity,
			USDPrice: price,
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

// getQuantityForDate returns the quantity held on a given date by using the last known quantity
func (v *ValuationManagerImpl) getQuantityForDate(timeline map[string]float64, dateStr string) float64 {
	if qty, exists := timeline[dateStr]; exists {
		return qty
	}

	// Backfill: find the closest previous date's quantity
	parsedDate, _ := time.Parse(time.DateOnly, dateStr)
	var closestQty float64
	var closestDate time.Time

	for dateKey, qty := range timeline {
		keyDate, _ := time.Parse(time.DateOnly, dateKey)
		if !keyDate.After(parsedDate) && (closestDate.IsZero() || keyDate.After(closestDate)) {
			closestQty = qty
			closestDate = keyDate
		}
	}

	return closestQty
}

// getClosestPrice finds the nearest previous market price for a given date
func (v *ValuationManagerImpl) getClosestPrice(dailyPrices map[string]float64, dateStr string) float64 {
	if price, exists := dailyPrices[dateStr]; exists {
		return price
	}

	// Backfill: find closest previous date with available price
	parsedDate, _ := time.Parse(time.DateOnly, dateStr)
	var closestPrice float64
	var closestDate time.Time

	for priceDate, price := range dailyPrices {
		keyDate, _ := time.Parse(time.DateOnly, priceDate)
		if !keyDate.After(parsedDate) && (closestDate.IsZero() || keyDate.After(closestDate)) {
			closestPrice = price
			closestDate = keyDate
		}
	}

	return closestPrice
}

// getClosestRate finds the nearest previous SBI exchange rate for a given date
func (v *ValuationManagerImpl) getClosestRate(dailyRates map[string]float64, dateStr string) float64 {
	if rate, exists := dailyRates[dateStr]; exists {
		return rate
	}

	// Backfill: find closest previous date with available rate
	parsedDate, _ := time.Parse(time.DateOnly, dateStr)
	var closestRate float64
	var closestDate time.Time

	for rateDate, rate := range dailyRates {
		keyDate, _ := time.Parse(time.DateOnly, rateDate)
		if !keyDate.After(parsedDate) && (closestDate.IsZero() || keyDate.After(closestDate)) {
			closestRate = rate
			closestDate = keyDate
		}
	}

	return closestRate
}
