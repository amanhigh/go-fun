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
	AnalyzeValuation(ctx context.Context, trades []tax.Trade, year int) (tax.Valuation, common.HttpError)
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
func (v *ValuationManagerImpl) processTickerTrades(ctx context.Context, _ string, sortedTrades []tax.Trade, year int) (tax.Valuation, common.HttpError) {
	// Call the *existing* AnalyzeValuation method for this ticker's sorted trades
	valuation, analyzeErr := v.AnalyzeValuation(ctx, sortedTrades, year)
	if analyzeErr != nil {
		return tax.Valuation{}, analyzeErr
	}
	return valuation, nil
}

func (v *ValuationManagerImpl) AnalyzeValuation(ctx context.Context, trades []tax.Trade, year int) (tax.Valuation, common.HttpError) {
	if err := v.validateTrades(trades); err != nil {
		return tax.Valuation{}, err
	}

	analysis := tax.Valuation{Ticker: trades[0].Symbol}
	openingPosition, err := v.getOpeningPositionForPeriod(ctx, analysis.Ticker, year)
	if err != nil {
		return tax.Valuation{}, common.NewServerError(fmt.Errorf("failed to get opening position for %s: %w", analysis.Ticker, err))
	}

	analysis.FirstPosition = openingPosition
	analysis.PeakPosition = openingPosition

	currentQuantity, _, err := v.processTrades(&analysis, trades, openingPosition)
	if err != nil {
		return tax.Valuation{}, err
	}

	if err := v.determineYearEndPosition(ctx, &analysis, year, currentQuantity, openingPosition, trades); err != nil {
		return tax.Valuation{}, err
	}

	return analysis, nil
}

// processTrades updates analysis based on trades and opening position.
// It returns the final currentQuantity, maxQuantityDuringPeriod, and any error.
func (v *ValuationManagerImpl) processTrades(
	analysis *tax.Valuation,
	trades []tax.Trade,
	openingPeriodPosition tax.Position,
) (currentQuantity, maxQuantityDuringPeriod float64, err common.HttpError) { // Combined return types
	currentQuantity = openingPeriodPosition.Quantity
	maxQuantityDuringPeriod = currentQuantity
	firstPositionAlreadySetByTrade := false

	for _, trade := range trades {
		tradeDate, dateErr := trade.GetDate()
		if dateErr != nil {
			return 0, 0, dateErr
		}

		isFreshStartScenario := openingPeriodPosition.Date.IsZero()

		if isFreshStartScenario && !firstPositionAlreadySetByTrade && trade.Type == "BUY" {
			maxQuantityDuringPeriod, firstPositionAlreadySetByTrade = v.handleFreshStartTrade(analysis, trade, tradeDate)
		}

		if trade.Type == "BUY" {
			currentQuantity += trade.Quantity
		} else {
			currentQuantity -= trade.Quantity
		}

		peakCtx := peakPositionContext{
			analysis:                       analysis,
			trade:                          trade,
			tradeDate:                      tradeDate,
			currentQuantity:                currentQuantity,
			maxQuantityDuringPeriod:        maxQuantityDuringPeriod,
			firstPositionAlreadySetByTrade: firstPositionAlreadySetByTrade,
			isFreshStartScenario:           isFreshStartScenario,
		}
		maxQuantityDuringPeriod = v.updatePeakPositionIfNeeded(peakCtx)
	}
	return currentQuantity, maxQuantityDuringPeriod, nil
}

// handleFreshStartTrade updates analysis for the first BUY trade in a fresh start scenario.
// It returns the new maxQuantityDuringPeriod and sets firstPositionAlreadySetByTrade to true.
func (v *ValuationManagerImpl) handleFreshStartTrade(analysis *tax.Valuation, trade tax.Trade, tradeDate time.Time) (float64, bool) {
	analysis.FirstPosition = tax.Position{
		Date:     tradeDate,
		Quantity: trade.Quantity,
		USDPrice: trade.USDPrice,
	}
	analysis.PeakPosition = analysis.FirstPosition
	return trade.Quantity, true
}

type peakPositionContext struct {
	analysis                       *tax.Valuation
	trade                          tax.Trade
	tradeDate                      time.Time
	currentQuantity                float64
	maxQuantityDuringPeriod        float64
	firstPositionAlreadySetByTrade bool
	isFreshStartScenario           bool
}

// updatePeakPositionIfNeeded updates the peak position if the current quantity is higher.
// It returns the updated maxQuantityDuringPeriod.
func (v *ValuationManagerImpl) updatePeakPositionIfNeeded(ctx peakPositionContext) float64 {
	if (ctx.firstPositionAlreadySetByTrade || !ctx.isFreshStartScenario) && ctx.currentQuantity > ctx.maxQuantityDuringPeriod {
		newMaxQuantity := ctx.currentQuantity
		ctx.analysis.PeakPosition = tax.Position{
			Date:     ctx.tradeDate,
			Quantity: newMaxQuantity,
			USDPrice: ctx.trade.USDPrice,
		}
		return newMaxQuantity
	}
	return ctx.maxQuantityDuringPeriod
}

// determineYearEndPosition sets the YearEndPosition in the analysis.
func (v *ValuationManagerImpl) determineYearEndPosition(
	ctx context.Context,
	analysis *tax.Valuation,
	year int,
	currentQuantity float64,
	openingPeriodPosition tax.Position,
	trades []tax.Trade,
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
	case openingPeriodPosition.Quantity > 0 && len(trades) == 0:
		price, priceErr := v.tickerManager.GetPrice(ctx, analysis.Ticker, yearEndDate)
		if priceErr != nil {
			return common.NewServerError(fmt.Errorf("failed to get year end price for carry-over only asset %s: %w", analysis.Ticker, priceErr))
		}
		analysis.YearEndPosition = tax.Position{
			Date:     yearEndDate,
			Quantity: openingPeriodPosition.Quantity,
			USDPrice: price,
		}
	default:
		analysis.YearEndPosition = tax.Position{Date: yearEndDate}
	}
	return nil
}

func (v *ValuationManagerImpl) validateTrades(trades []tax.Trade) common.HttpError {
	if len(trades) == 0 {
		return common.NewHttpError("no trades provided", http.StatusBadRequest)
	}
	for _, t := range trades {
		if t.Symbol != trades[0].Symbol {
			return common.NewHttpError("multiple tickers found in trades", http.StatusBadRequest)
		}
	}
	return nil
}

func (v *ValuationManagerImpl) getOpeningPositionForPeriod(ctx context.Context, ticker string, year int) (position tax.Position, err common.HttpError) {
	// Last year's account record
	account, accErr := v.accountManager.GetRecord(ctx, ticker)
	if accErr != nil {
		if errors.Is(accErr, common.ErrNotFound) {
			// No account record found -> fresh start for this period.
			// Return a zero position. Its Date field will be the zero value for time.Time.
			return tax.Position{}, nil
		}
		return tax.Position{}, accErr // Other errors from accountManager
	}

	// Account record found (carry-over scenario)
	// Account record found (carry-over scenario)
	openingDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
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
