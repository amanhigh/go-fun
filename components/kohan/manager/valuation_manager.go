package manager

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	"github.com/samber/lo"
)

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
}

func NewValuationManager(
	tickerManager TickerManager,
	accountManager AccountManager,
	tradeRepository repository.TradeRepository,
) ValuationManager {
	return &ValuationManagerImpl{
		tickerManager:   tickerManager,
		accountManager:  accountManager,
		tradeRepository: tradeRepository,
	}
}

func (v *ValuationManagerImpl) GetYearlyValuationsUSD(ctx context.Context, year int) (valuations []tax.Valuation, err common.HttpError) {
	allTrades, repoErr := v.tradeRepository.GetAllRecords(ctx)
	if repoErr != nil {
		return nil, repoErr // Return other errors
	}

	if len(allTrades) == 0 {
		return []tax.Valuation{}, common.NewHttpError("no trades found", http.StatusNotFound)
	}

	// Group trades directly by Ticker Symbol
	tradesByTicker := lo.GroupBy(allTrades, func(trade tax.Trade) string {
		return trade.Symbol
	})

	valuations = make([]tax.Valuation, 0, len(tradesByTicker))

	// Iterate through the grouped map (ticker -> list of trades)
	for _, tickerTrades := range tradesByTicker {
		// Sort the trades for this specific ticker chronologically
		slices.SortFunc(tickerTrades, func(a, b tax.Trade) int {
			if a.GetDate().Before(b.GetDate()) {
				return -1
			} else if a.GetDate().After(b.GetDate()) {
				return 1
			}
			return 0
		})

		// Call the *existing* AnalyzeValuation method for this ticker's sorted trades
		valuation, analyzeErr := v.AnalyzeValuation(ctx, tickerTrades, year)
		if analyzeErr != nil {
			// Fail fast: return immediately upon the first analysis error
			return nil, analyzeErr
		}
		valuations = append(valuations, valuation)
	}

	return valuations, nil // Return aggregated results if all analyses succeeded
}

func (v *ValuationManagerImpl) AnalyzeValuation(ctx context.Context, trades []tax.Trade, year int) (tax.Valuation, common.HttpError) {
	if len(trades) == 0 {
		return tax.Valuation{}, common.NewHttpError("no trades provided", http.StatusBadRequest)
	}

	if err := v.validateTrades(trades); err != nil {
		return tax.Valuation{}, err
	}

	// Set ticker and get starting position
	analysis := tax.Valuation{
		Ticker: trades[0].Symbol,
	}
	startPosition, err := v.getStartingPosition(ctx, analysis.Ticker, year)
	if err != nil {
		return analysis, err
	}

	// Process trades with starting position
	// HACK: Simplify this Class
	currentPosition := v.trackPositions(&analysis, startPosition, trades)

	// Update year-end position if there are remaining holdings
	if currentPosition > 0 {
		yearEndDate := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
		price, err := v.tickerManager.GetPrice(ctx, analysis.Ticker, yearEndDate)
		if err != nil {
			return analysis, common.NewHttpError("failed to get year end price", http.StatusInternalServerError)
		}

		analysis.YearEndPosition = tax.Position{
			Date:     yearEndDate,
			Quantity: currentPosition,
			USDPrice: price,
		}
	}

	return analysis, nil
}

func (v *ValuationManagerImpl) validateTrades(trades []tax.Trade) common.HttpError {
	for _, t := range trades {
		if t.Symbol != trades[0].Symbol {
			return common.NewHttpError("multiple tickers found in trades", http.StatusBadRequest)
		}
	}
	return nil
}

func (v *ValuationManagerImpl) getStartingPosition(ctx context.Context, ticker string, year int) (position tax.Position, err common.HttpError) {
	// Last year's account record
	account, err := v.accountManager.GetRecord(ctx, ticker)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			// Return zero position for fresh start
			return position, nil
		}
		return position, err
	}

	// Convert account to last year-end position
	lastYearEnd := time.Date(year-1, 12, 31, 0, 0, 0, 0, time.UTC)
	position = tax.Position{
		Date:     lastYearEnd,
		Quantity: account.Quantity,
		USDPrice: account.MarketValue / account.Quantity,
	}
	return
}

func (v *ValuationManagerImpl) trackPositions(analysis *tax.Valuation, startPosition tax.Position, trades []tax.Trade) (currentPosition float64) {
	currentPosition = startPosition.Quantity
	maxPosition := currentPosition

	// Set initial position
	analysis.FirstPosition = startPosition
	analysis.PeakPosition = startPosition

	// Process all trades
	for _, t := range trades {
		tradeDate := t.GetDate()
		if t.Type == "BUY" {
			currentPosition += t.Quantity
		} else {
			currentPosition -= t.Quantity
		}

		// Update first position if starting from zero
		if startPosition.Quantity == 0 && t.Type == "BUY" && analysis.FirstPosition.Quantity == 0 {
			analysis.FirstPosition = tax.Position{
				Date:     tradeDate,
				Quantity: t.Quantity,
				USDPrice: t.USDPrice,
			}
		}

		// Track peak position
		if currentPosition > maxPosition {
			maxPosition = currentPosition
			analysis.PeakPosition = tax.Position{
				Date:     tradeDate,
				Quantity: maxPosition,
				USDPrice: t.USDPrice,
			}
		}
	}
	return
}
