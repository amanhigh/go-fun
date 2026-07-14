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
//
//nolint:funlen // Multi-step pipeline (validate, fetch splits, build timeline, peak, year-end)
func (v *ValuationManagerImpl) AnalyzeValuation(ctx context.Context, tickerSymbol string, trades []tax.Trade, year int) (tax.Valuation, common.HttpError) {
	// Step 1: Validate trade symbols first. This is crucial for the "Multiple Ticker Trades" test
	// to fail before attempting to get an opening position if symbols are inconsistent.
	if err := v.validateTradeSymbols(trades, tickerSymbol); err != nil {
		return tax.Valuation{}, err
	}

	// Step 2: Get opening positions (FirstPosition from Origin metadata, holding for Peak/Closing from Quantity)
	firstPosition, holdingPosition, err := v.getOpeningPositions(ctx, tickerSymbol, year)
	if err != nil {
		return tax.Valuation{}, common.NewServerError(fmt.Errorf("failed to get opening position for %s: %w", tickerSymbol, err))
	}

	// Step 3: Validate if trades exist or if there's a carry-over
	// (checked before GetSplits so invalid/no-position requests do not
	// query ticker data and a split-fetch error cannot mask this error).
	if err := v.validateTradesExistOrCarryOver(trades, holdingPosition, tickerSymbol); err != nil {
		return tax.Valuation{}, err
	}

	// Step 4: Validate first trade isn't sell on fresh start (before any market data calls)
	if err := v.validateFirstTradeNotSellOnFreshStart(trades, holdingPosition, tickerSymbol); err != nil {
		return tax.Valuation{}, err
	}

	// Step 5: Validate negative net quantity on fresh start (before any market data calls)
	if err := v.validateFirstDateNetNotNegative(trades, holdingPosition, tickerSymbol); err != nil {
		return tax.Valuation{}, err
	}

	// Step 6: Get split events for the calendar year (Jan 1 - Dec 31 inclusive)
	splitStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	splitEnd := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
	splits, splitErr := v.tickerManager.GetSplits(ctx, tickerSymbol, splitStart, splitEnd)
	if splitErr != nil {
		return tax.Valuation{}, splitErr
	}

	// setupFirstPosition derives FirstPosition for a fresh-start by aggregating
	// all trades on the first date (BUY added, SELL subtracted). Single BUY uses
	// the trade price directly; multiple trades with positive net use GetPrice;
	// zero net preserves the date without a price lookup. Carry-over positions
	// (non-zero Date) are returned unchanged.
	firstPos, fpErr := v.setupFirstPosition(ctx, trades, firstPosition, tickerSymbol)
	if fpErr != nil {
		return tax.Valuation{}, fpErr
	}

	analysis := tax.Valuation{Ticker: tickerSymbol, FirstPosition: firstPos}

	// Step 7: Build daily quantity timeline with event-date split awareness.
	// Split events are applied before trades on the same date,
	// then end-of-day quantity is recorded.
	// TODO: Genuine intraday support (e.g., split after trade on same day) is deliberately out of scope.
	quantityByDate := v.buildDailyQuantityTimeline(year, holdingPosition, trades, splits)

	// Step 8: Calculate daily peak value from the timeline (Tax.md daily peak calculation rule — mandatory)
	peakPosition, peakErr := v.calculateDailyPeak(ctx, tickerSymbol, year, holdingPosition, quantityByDate)
	if peakErr != nil {
		return tax.Valuation{}, common.NewServerError(
			fmt.Errorf("failed to calculate daily peak for %s: %w", tickerSymbol, peakErr))
	}
	analysis.PeakPosition = peakPosition

	// Step 9: Determine year-end position from timeline Dec 31 end-of-day quantity
	yearEndQuantity := v.getClosestValue(quantityByDate, splitEnd.Format(time.DateOnly))
	if detErr := v.determineYearEndPosition(ctx, &analysis, year, yearEndQuantity); detErr != nil {
		return tax.Valuation{}, detErr
	}

	return analysis, nil
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

// validateFirstDateNetNotNegative checks that the net quantity from aggregating
// all trades on the first trade date is not negative on a fresh start.
// Single-trade and first-trade-SELL scenarios are handled separately by
// validateFirstTradeNotSellOnFreshStart.
func (v *ValuationManagerImpl) validateFirstDateNetNotNegative(trades []tax.Trade, holdingPosition tax.Position, tickerSymbol string) common.HttpError {
	// Only applicable on fresh start with trades
	if holdingPosition.Quantity != 0 || len(trades) == 0 || trades[0].GetType() == tax.TRADE_TYPE_SELL {
		return nil
	}

	firstDate, dateErr := trades[0].GetDate()
	if dateErr != nil {
		return dateErr
	}

	_, netQty, aggErr := aggregateFirstDateTrades(trades, firstDate)
	if aggErr != nil {
		return aggErr
	}

	if netQty < 0 {
		return common.NewHttpError(fmt.Sprintf("negative net quantity on first date for %s", tickerSymbol), http.StatusBadRequest)
	}

	return nil
}

// aggregateFirstDateTrades collects all trades on the same calendar date as the
// first trade and calculates the net quantity (BUY added, SELL subtracted).
func aggregateFirstDateTrades(trades []tax.Trade, firstDate time.Time) (firstDateTrades []tax.Trade, netQty float64, httpErr common.HttpError) {
	firstDateStr := firstDate.Format(time.DateOnly)
	for _, trade := range trades {
		tradeDate, tErr := trade.GetDate()
		if tErr != nil {
			return nil, 0, tErr
		}
		if tradeDate.Format(time.DateOnly) != firstDateStr {
			break
		}
		firstDateTrades = append(firstDateTrades, trade)
		if trade.GetType() == tax.TRADE_TYPE_BUY {
			netQty += trade.Quantity
		} else {
			netQty -= trade.Quantity
		}
	}
	return firstDateTrades, netQty, nil
}

// setupFirstPosition determines the FirstPosition for a fresh-start scenario
// (no carry-over) with the following rules:
//   - Multiple trades on the first date: all trades are aggregated (BUY added,
//     SELL subtracted). If the net quantity is positive, GetPrice is called
//     for the historical closing price. If zero, the date is preserved with
//     Quantity=0 and USDPrice=0 without a price lookup.
//   - Single BUY trade on the first date: the trade's price is used directly
//     (preserving existing multi-day single-first-BUY behavior).
//   - Single SELL on fresh start: returns the zero firstPosition (the separate
//     validateFirstTradeNotSellOnFreshStart check produces the error).
//   - Carry-over (non-zero Date): always returned unchanged, even when
//     OriginQty=0 but OriginDate is valid.
func (v *ValuationManagerImpl) setupFirstPosition(ctx context.Context, trades []tax.Trade, firstPosition tax.Position, ticker string) (tax.Position, common.HttpError) {
	// Preserve carry-over positions (even when OriginQty=0 but OriginDate is valid)
	if !firstPosition.Date.IsZero() {
		return firstPosition, nil
	}

	// Fresh start with no trades — nothing to derive
	if len(trades) == 0 {
		return firstPosition, nil
	}

	firstDate, dateErr := trades[0].GetDate()
	if dateErr != nil {
		return tax.Position{}, dateErr
	}

	firstDateTrades, netQty, aggErr := aggregateFirstDateTrades(trades, firstDate)
	if aggErr != nil {
		return tax.Position{}, aggErr
	}

	// Multiple trades on the first date: aggregate and potentially call GetPrice
	if len(firstDateTrades) > 1 {
		if netQty > 0 {
			price, priceErr := v.tickerManager.GetPrice(ctx, ticker, firstDate)
			if priceErr != nil {
				return tax.Position{}, priceErr
			}
			return tax.Position{Date: firstDate, Quantity: netQty, USDPrice: price}, nil
		}
		// Zero net quantity: preserve the date without a price lookup
		return tax.Position{Date: firstDate}, nil
	}

	// Single BUY trade on the first date: preserve existing behavior
	if firstDateTrades[0].GetType() == tax.TRADE_TYPE_BUY {
		return tax.Position{
			Date:     firstDate,
			Quantity: firstDateTrades[0].Quantity,
			USDPrice: firstDateTrades[0].USDPrice,
		}, nil
	}

	return firstPosition, nil
}

// validateFirstTradeNotSellOnFreshStart returns an error when there is no
// carry-over position and the first trade is a SELL, which is invalid.
func (v *ValuationManagerImpl) validateFirstTradeNotSellOnFreshStart(trades []tax.Trade, holdingPosition tax.Position, tickerSymbol string) common.HttpError {
	if holdingPosition.Quantity == 0 && len(trades) > 0 && trades[0].GetType() == tax.TRADE_TYPE_SELL {
		return common.NewHttpError(fmt.Sprintf("first trade can't be sell on fresh start for %s", tickerSymbol), http.StatusBadRequest)
	}
	return nil
}

// getOpeningPositions returns two positions from a single GetRecord call:
//   - firstPosition: OriginQty-based (for First/Initial metadata)
//   - holdingPosition: Quantity-based (for Peak/Closing calculations)
//
// Raw account values are used directly — no split normalization is applied.
// In a fresh-start scenario, both are returned as zero positions.
func (v *ValuationManagerImpl) getOpeningPositions(ctx context.Context, ticker string, year int) (firstPosition, holdingPosition tax.Position, err common.HttpError) {
	account, accErr := v.accountManager.GetRecord(ctx, ticker, year-1)
	if accErr != nil {
		if errors.Is(accErr, common.ErrNotFound) {
			return tax.Position{}, tax.Position{}, nil
		}
		return tax.Position{}, tax.Position{}, accErr
	}

	// Quantity == 0 means the position was fully liquidated at year-end.
	// Treat as a fresh start even when stale origin metadata remains.
	if account.Quantity == 0 {
		return tax.Position{}, tax.Position{}, nil
	}

	// Account record found (carry-over scenario)
	// Reconstruct FirstPosition from Account metadata (original acquisition date/price)
	// OriginDate MUST be present for carry-over accounts - it's required for tax reporting
	originDate, parseErr := time.Parse(time.DateOnly, account.OriginDate)
	if parseErr != nil {
		return tax.Position{}, tax.Position{}, tax.NewInvalidDateError(
			fmt.Sprintf("failed to parse OriginDate '%s' for carry-over account %s: %v", account.OriginDate, ticker, parseErr))
	}

	firstPosition = tax.Position{
		Date:     originDate,
		Quantity: account.OriginQty,
		USDPrice: account.OriginPrice,
	}
	holdingPosition = tax.Position{
		Date:     originDate,
		Quantity: account.Quantity,
		USDPrice: account.OriginPrice,
	}
	return firstPosition, holdingPosition, nil
}

// calculateDailyPeak evaluates (Quantity × Market_Price × SBI_Rate) for every day in the year
// to find the true INR peak value during the calendar year.
// This ensures compliance with the Tax.md daily peak calculation rule.
// The quantityByDate timeline is pre-built by AnalyzeValuation with split events applied.
func (v *ValuationManagerImpl) calculateDailyPeak(
	ctx context.Context,
	ticker string,
	year int,
	openingPosition tax.Position,
	quantityByDate map[string]float64,
) (peakPosition tax.Position, err common.HttpError) {
	// Step 1: Get daily market prices for the ticker
	dailyPrices, priceErr := v.tickerManager.GetDailyPrices(ctx, ticker, year)
	if priceErr != nil {
		return tax.Position{}, priceErr
	}

	// Step 2: Get daily SBI TT Buy rates for the year
	dailyRates, rateErr := v.sbiManager.GetDailyRates(ctx, year)
	if rateErr != nil {
		return tax.Position{}, rateErr
	}

	// Step 3: Find the date with maximum INR value by iterating through each day
	return v.findPeakByIteratingYear(year, openingPosition, quantityByDate, dailyPrices, dailyRates), nil
}

// findPeakByIteratingYear finds maximum INR value (Qty × Price × Rate) across the year (Tax.md daily peak calculation rule).
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

// buildDailyQuantityTimeline creates a map of date → end-of-day quantity held.
// Split events are applied before trades on the same date (split-before-same-day-trades ordering).
// The timeline is used for both daily peak and Dec 31 YearEnd quantity.
//
// TODO: Genuine intraday support (e.g., split after trade on same day) is deliberately out of scope.
func (v *ValuationManagerImpl) buildDailyQuantityTimeline(
	year int,
	openingPosition tax.Position,
	trades []tax.Trade,
	splits []tax.YahooSplit,
) map[string]float64 {
	timeline := make(map[string]float64)
	currentQuantity := openingPosition.Quantity

	// Initialize the year with opening quantity
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	timeline[startDate.Format(time.DateOnly)] = currentQuantity

	// Group trades and splits by calendar date, then collect sorted event dates
	tradesByDate := groupTradesByDate(trades, year)
	splitsByDate := groupSplitsByDate(splits, year)
	dates := collectEventDates(tradesByDate, splitsByDate)

	// Process events chronologically: splits before trades on same date, then record end-of-day
	for _, dateStr := range dates {
		// Step 1: Apply all split events on this date first (before trades)
		for _, split := range splitsByDate[dateStr] {
			ratio := split.Ratio()
			currentQuantity *= ratio
		}

		// Step 2: Apply all trades on this date
		for _, trade := range tradesByDate[dateStr] {
			if trade.GetType() == tax.TRADE_TYPE_BUY {
				currentQuantity += trade.Quantity
			} else { // SELL
				currentQuantity -= trade.Quantity
			}
		}

		// Step 3: Record end-of-day quantity
		timeline[dateStr] = currentQuantity
	}

	return timeline
}

// groupTradesByDate partitions trades by calendar date, filtering out entries
// whose date cannot be parsed or falls outside the given year.
func groupTradesByDate(trades []tax.Trade, year int) map[string][]tax.Trade {
	tradesByDate := make(map[string][]tax.Trade)
	for _, trade := range trades {
		tradeDate, err := trade.GetDate()
		if err != nil || tradeDate.Year() != year {
			continue
		}
		dateStr := tradeDate.Format(time.DateOnly)
		tradesByDate[dateStr] = append(tradesByDate[dateStr], trade)
	}
	return tradesByDate
}

// groupSplitsByDate partitions Yahoo split events by calendar date, filtering
// out entries whose timestamp falls outside the given year.
func groupSplitsByDate(splits []tax.YahooSplit, year int) map[string][]tax.YahooSplit {
	splitsByDate := make(map[string][]tax.YahooSplit)
	for _, split := range splits {
		splitDate := split.EffectiveDate()
		if splitDate.Year() != year {
			continue
		}
		dateStr := splitDate.Format(time.DateOnly)
		splitsByDate[dateStr] = append(splitsByDate[dateStr], split)
	}
	return splitsByDate
}

// collectEventDates returns a chronologically sorted slice of unique dates
// present in either the trade-by-date or split-by-date maps.
func collectEventDates(tradesByDate map[string][]tax.Trade, splitsByDate map[string][]tax.YahooSplit) []string {
	dates := make([]string, 0, len(tradesByDate)+len(splitsByDate))
	for dateStr := range tradesByDate {
		dates = append(dates, dateStr)
	}
	for dateStr := range splitsByDate {
		dates = append(dates, dateStr)
	}
	slices.Sort(dates)
	return lo.Uniq(dates)
}

// getClosestValue finds the nearest previous value for a given date using backfill logic.
// If the exact date exists in the map, it returns that value immediately.
// Otherwise, it searches for the closest previous date with available data.
// Returns 0 if no previous data is found.
//
// Invariant: dateStr and all dateKey elements are canonical time.DateOnly
// strings ("YYYY-MM-DD"), so lexicographic ordering === chronological ordering.
func (v *ValuationManagerImpl) getClosestValue(dataMap map[string]float64, dateStr string) float64 {
	if value, exists := dataMap[dateStr]; exists {
		return value
	}

	// Backfill: find the closest previous date's value
	var closestValue float64
	var closestDate string

	for dateKey, value := range dataMap {
		if dateKey < dateStr && (closestDate == "" || dateKey > closestDate) {
			closestValue = value
			closestDate = dateKey
		}
	}

	return closestValue
}
