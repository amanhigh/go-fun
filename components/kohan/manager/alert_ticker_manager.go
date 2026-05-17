package manager

import (
	"context"
	"net/http"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// AlertTickerManager provides business logic for Alert ticker CRUD operations.
type AlertTickerManager interface {
	// CreateAlertTicker creates a new Alert ticker under a parent ticker.
	CreateAlertTicker(ctx context.Context, ticker string, alert *barkat.AlertTicker) (barkat.AlertTicker, common.HttpError)
	// GetAlertTicker retrieves a single Alert ticker by symbol.
	GetAlertTicker(ctx context.Context, symbol string) (barkat.AlertTicker, common.HttpError)
	// DeleteAlertTicker deletes an Alert ticker by symbol.
	DeleteAlertTicker(ctx context.Context, symbol string) common.HttpError
	// ListAlertTickers returns a filtered, paginated list of Alert tickers.
	ListAlertTickers(ctx context.Context, query barkat.AlertTickerQuery) (barkat.AlertTickerList, common.HttpError)
}

type AlertTickerManagerImpl struct {
	alertRepo repository.AlertTickerRepository
}

var _ AlertTickerManager = (*AlertTickerManagerImpl)(nil)

// NewAlertTickerManager creates a new AlertTickerManager.
func NewAlertTickerManager(alertRepo repository.AlertTickerRepository) *AlertTickerManagerImpl {
	return &AlertTickerManagerImpl{alertRepo: alertRepo}
}

// ---- Alert Ticker ----

func (m *AlertTickerManagerImpl) CreateAlertTicker(ctx context.Context, ticker string, alert *barkat.AlertTicker) (result barkat.AlertTicker, httpErr common.HttpError) {
	httpErr = m.alertRepo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Look up parent ticker
		var parent barkat.Ticker
		if httpErr := m.alertRepo.GetByExternalId(c, ticker, &parent); httpErr != nil {
			return httpErr
		}

		alert.TickerID = parent.ID

		if httpErr := m.alertRepo.Create(c, alert); httpErr != nil {
			return httpErr
		}

		// Populate parent ticker string for response
		result = *alert
		result.Ticker = parent.Ticker
		return nil
	})
	return
}

func (m *AlertTickerManagerImpl) GetAlertTicker(ctx context.Context, symbol string) (barkat.AlertTicker, common.HttpError) {
	return m.alertRepo.GetBySymbol(ctx, symbol)
}

func (m *AlertTickerManagerImpl) DeleteAlertTicker(ctx context.Context, symbol string) common.HttpError {
	return m.alertRepo.DeleteBySymbol(ctx, symbol)
}

func (m *AlertTickerManagerImpl) ListAlertTickers(ctx context.Context, query barkat.AlertTickerQuery) (barkat.AlertTickerList, common.HttpError) {
	// When ticker filter is provided, validate parent ticker exists
	if query.Ticker != "" {
		var parent barkat.Ticker
		if httpErr := m.alertRepo.GetByExternalId(ctx, query.Ticker, &parent); httpErr != nil {
			if httpErr.Code() == http.StatusNotFound {
				return barkat.AlertTickerList{}, common.NewHttpError("Ticker not found", http.StatusNotFound)
			}
			return barkat.AlertTickerList{}, httpErr
		}
	}

	alertTickers, total, httpErr := m.alertRepo.ListAlertTickers(ctx, query)
	if httpErr != nil {
		return barkat.AlertTickerList{}, httpErr
	}

	return barkat.AlertTickerList{
		AlertTickers: alertTickers,
		Metadata: common.PaginatedResponse{
			Total:  total,
			Offset: query.Offset,
			Limit:  query.Limit,
		},
	}, nil
}
