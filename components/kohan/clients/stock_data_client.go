package clients

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	tax "github.com/amanhigh/go-fun/models/tax"
)

// StockDataClient is the interface for fetching stock price data from any source
type StockDataClient interface {
	FetchDailyPrices(ctx context.Context, ticker string) (tax.StockData, common.HttpError)
}
