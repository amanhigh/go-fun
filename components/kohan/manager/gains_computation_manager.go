package manager

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
)

//go:generate mockery --name GainsComputationManager
type GainsComputationManager interface {
	ComputeGainsFromTrades(ctx context.Context, trades []tax.Trade) ([]tax.Gains, common.HttpError)
}

// GainsComputationManagerImpl handles FIFO-based gains computation from trades
type GainsComputationManagerImpl struct{}

// NewGainsComputationManager creates a new GainsComputationManager
func NewGainsComputationManager() GainsComputationManager {
	return &GainsComputationManagerImpl{}
}

// PositionLot represents a buy position lot for FIFO matching
type PositionLot struct {
	BuyDate    time.Time
	Quantity   float64
	USDPrice   float64
	Commission float64
	Symbol     string
}

func (g *GainsComputationManagerImpl) ComputeGainsFromTrades(_ context.Context, trades []tax.Trade) ([]tax.Gains, common.HttpError) {
	sortedTrades := g.sortTradesByDate(trades)

	// Track buy positions by symbol using FIFO
	buyPositions := make(map[string][]PositionLot)
	var gains []tax.Gains

	for _, trade := range sortedTrades {
		tradeDate, err := time.Parse(time.DateOnly, trade.Date)
		if err != nil {
			return nil, common.NewHttpError(fmt.Sprintf("invalid trade date '%s': %v", trade.Date, err), http.StatusBadRequest)
		}

		switch trade.Type {
		case tax.TRADE_TYPE_BUY:
			// Add to buy positions for this symbol
			lot := PositionLot{
				BuyDate:    tradeDate,
				Quantity:   trade.Quantity,
				USDPrice:   trade.USDPrice,
				Commission: trade.Commission,
				Symbol:     trade.Symbol,
			}
			buyPositions[trade.Symbol] = append(buyPositions[trade.Symbol], lot)

		case tax.TRADE_TYPE_SELL:
			// Match against existing buy positions using FIFO
			sellGains, httpErr := g.matchSellWithBuys(trade, tradeDate, buyPositions)
			if httpErr != nil {
				return nil, httpErr
			}
			gains = append(gains, sellGains...)
		}
	}

	return gains, nil
}

// sortTradesByDate sorts trades chronologically by date
func (g *GainsComputationManagerImpl) sortTradesByDate(trades []tax.Trade) []tax.Trade {
	sortedTrades := make([]tax.Trade, len(trades))
	copy(sortedTrades, trades)

	sort.SliceStable(sortedTrades, func(i, j int) bool {
		dateI, errI := time.Parse(time.DateOnly, sortedTrades[i].Date)
		dateJ, errJ := time.Parse(time.DateOnly, sortedTrades[j].Date)
		if errI != nil || errJ != nil {
			return false // Keep original order if date parsing fails
		}
		return dateI.Before(dateJ)
	})

	return sortedTrades
}

// matchSellWithBuys matches a sell transaction with buy positions using FIFO
func (g *GainsComputationManagerImpl) matchSellWithBuys(sellTrade tax.Trade, sellDate time.Time, buyPositions map[string][]PositionLot) ([]tax.Gains, common.HttpError) {
	lots, exists := buyPositions[sellTrade.Symbol]
	if !exists || len(lots) == 0 {
		return nil, common.NewHttpError(fmt.Sprintf("no buy positions found for sell trade: %s", sellTrade.Symbol), http.StatusBadRequest)
	}

	gains, updatedLots, remainingQuantity := g.processFIFOMatching(sellTrade, sellDate, lots)

	// Update buy positions
	buyPositions[sellTrade.Symbol] = updatedLots

	// Check if we have unmatched sell quantity
	if remainingQuantity > 0 {
		return nil, common.NewHttpError(fmt.Sprintf("insufficient buy quantity for sell trade: %s, remaining: %.2f", sellTrade.Symbol, remainingQuantity), http.StatusBadRequest)
	}

	return gains, nil
}

// processFIFOMatching processes sell trade against buy lots using FIFO order
func (g *GainsComputationManagerImpl) processFIFOMatching(sellTrade tax.Trade, sellDate time.Time, lots []PositionLot) ([]tax.Gains, []PositionLot, float64) {
	var gains []tax.Gains
	var updatedLots []PositionLot
	remainingQuantity := sellTrade.Quantity

	for _, lot := range lots {
		if remainingQuantity <= 0 {
			updatedLots = append(updatedLots, lot)
			continue
		}

		quantityToMatch := g.calculateQuantityToMatch(remainingQuantity, lot.Quantity)
		gain := g.createGainFromLot(sellTrade, sellDate, lot, quantityToMatch)
		gains = append(gains, gain)

		// Update remaining quantities
		remainingQuantity -= quantityToMatch
		lot.Quantity -= quantityToMatch

		// Keep lot if it has remaining quantity
		if lot.Quantity > 0 {
			updatedLots = append(updatedLots, lot)
		}
	}

	return gains, updatedLots, remainingQuantity
}

// calculateQuantityToMatch determines how much quantity to match from a lot
func (g *GainsComputationManagerImpl) calculateQuantityToMatch(remainingQuantity, lotQuantity float64) float64 {
	if remainingQuantity > lotQuantity {
		return lotQuantity
	}
	return remainingQuantity
}

// createGainFromLot creates a gain record from matching a sell trade with a buy lot
func (g *GainsComputationManagerImpl) createGainFromLot(sellTrade tax.Trade, sellDate time.Time, lot PositionLot, quantityToMatch float64) tax.Gains {
	// Calculate PNL for this matched portion
	pnl := (sellTrade.USDPrice - lot.USDPrice) * quantityToMatch

	// Commission allocation - proportional to quantity
	buyCommissionPortion := lot.Commission * (quantityToMatch / lot.Quantity)
	sellCommissionPortion := sellTrade.Commission * (quantityToMatch / sellTrade.Quantity)
	totalCommission := buyCommissionPortion + sellCommissionPortion

	// Adjust PNL for commissions
	pnl -= totalCommission

	// Classify as STCG or LTCG (2-year rule for foreign stocks)
	gainType := g.classifyGainType(lot.BuyDate, sellDate)

	return tax.Gains{
		Symbol:     sellTrade.Symbol,
		BuyDate:    lot.BuyDate.Format(time.DateOnly),
		SellDate:   sellDate.Format(time.DateOnly),
		Quantity:   quantityToMatch,
		PNL:        pnl,
		Commission: totalCommission,
		Type:       gainType,
	}
}

// classifyGainType determines if gain is STCG or LTCG based on holding period
// For foreign stocks: LTCG >= 24 months, STCG < 24 months
func (g *GainsComputationManagerImpl) classifyGainType(buyDate, sellDate time.Time) string {
	holdingPeriod := sellDate.Sub(buyDate)

	// Foreign stocks: 2 years = 730 days (accounting for potential leap years)
	const (
		foreignStockLTCGDays = 730
		hoursPerDay          = 24
	)
	twoYears := foreignStockLTCGDays * hoursPerDay * time.Hour

	if holdingPeriod >= twoYears {
		return tax.GAIN_TYPE_LTCG
	}
	return tax.GAIN_TYPE_STCG
}
