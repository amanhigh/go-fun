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
	repo repository.AlertTickerRepository
}

var _ AlertTickerManager = (*AlertTickerManagerImpl)(nil)

// NewAlertTickerManager creates a new AlertTickerManager.
func NewAlertTickerManager(repo repository.AlertTickerRepository) *AlertTickerManagerImpl {
	return &AlertTickerManagerImpl{repo: repo}
}

// ---- Helpers ----

// hydrateParentTicker populates the Ticker (parent external_id) field on an AlertTicker.
func (m *AlertTickerManagerImpl) hydrateParentTicker(ctx context.Context, alert *barkat.AlertTicker) common.HttpError {
	var parent barkat.Ticker
	if httpErr := m.repo.FindById(ctx, alert.TickerID, &parent); httpErr != nil {
		return httpErr
	}
	alert.Ticker = parent.Ticker
	return nil
}

// ---- Alert Ticker ----

func (m *AlertTickerManagerImpl) CreateAlertTicker(ctx context.Context, ticker string, alert *barkat.AlertTicker) (result barkat.AlertTicker, httpErr common.HttpError) {
	httpErr = m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Validate and look up parent ticker
		var parent barkat.Ticker
		if httpErr := m.repo.GetByExternalId(c, ticker, &parent); httpErr != nil {
			return httpErr
		}

		alert.TickerID = parent.ID

		if httpErr := m.repo.Create(c, alert); httpErr != nil {
			return httpErr
		}

		result = *alert
		result.Ticker = parent.Ticker
		return nil
	})
	return
}

func (m *AlertTickerManagerImpl) GetAlertTicker(ctx context.Context, symbol string) (result barkat.AlertTicker, httpErr common.HttpError) {
	httpErr = m.repo.GetByExternalId(ctx, symbol, &result)
	if httpErr != nil {
		return
	}
	httpErr = m.hydrateParentTicker(ctx, &result)
	return
}

func (m *AlertTickerManagerImpl) DeleteAlertTicker(ctx context.Context, symbol string) common.HttpError {
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		var alert barkat.AlertTicker
		if httpErr := m.repo.GetByExternalId(c, symbol, &alert); httpErr != nil {
			return httpErr
		}
		return m.repo.DeleteById(c, alert.ID, &barkat.AlertTicker{})
	})
}

func (m *AlertTickerManagerImpl) ListAlertTickers(ctx context.Context, query barkat.AlertTickerQuery) (barkat.AlertTickerList, common.HttpError) {
	// When ticker filter is provided, validate parent ticker exists
	if query.Ticker != "" {
		var parent barkat.Ticker
		if httpErr := m.repo.GetByExternalId(ctx, query.Ticker, &parent); httpErr != nil {
			if httpErr.Code() == http.StatusNotFound {
				return barkat.AlertTickerList{}, common.NewHttpError("Ticker not found", http.StatusNotFound)
			}
			return barkat.AlertTickerList{}, httpErr
		}
	}

	alertTickers, total, httpErr := m.repo.ListAlertTickers(ctx, query)
	if httpErr != nil {
		return barkat.AlertTickerList{}, httpErr
	}

	// Hydrate parent ticker for each result
	for i := range alertTickers {
		if httpErr := m.hydrateParentTicker(ctx, &alertTickers[i]); httpErr != nil {
			return barkat.AlertTickerList{}, httpErr
		}
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
