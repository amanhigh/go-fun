package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// BarkatTickerManager provides business logic for barkat ticker CRUD operations.
type BarkatTickerManager interface {
	// CreateTicker creates a new ticker.
	CreateTicker(ctx context.Context, ticker *barkat.Ticker) common.HttpError
	// GetTicker retrieves a single ticker by ticker identity.
	GetTicker(ctx context.Context, ticker string) (barkat.Ticker, common.HttpError)
	// ListTickers returns a filtered, paginated list of tickers.
	ListTickers(ctx context.Context, query barkat.TickerQuery) (barkat.TickerList, common.HttpError)
	// UpdateTicker replaces mutable fields on an existing ticker.
	UpdateTicker(ctx context.Context, ticker string, req barkat.TickerUpdateRequest) (barkat.Ticker, common.HttpError)
	// PatchTickerLastOpened updates only the last_opened_at timestamp.
	PatchTickerLastOpened(ctx context.Context, ticker string, update barkat.TickerLastOpenedUpdate) (barkat.Ticker, common.HttpError)
	// DeleteTicker deletes a ticker and cascades to linked AlertTickers and PriceAlerts.
	DeleteTicker(ctx context.Context, ticker string) common.HttpError
}

type BarkatTickerManagerImpl struct {
	repo repository.TickerRepository
}

var _ BarkatTickerManager = (*BarkatTickerManagerImpl)(nil)

// NewBarkatTickerManager creates a new BarkatTickerManager.
func NewBarkatTickerManager(repo repository.TickerRepository) *BarkatTickerManagerImpl {
	return &BarkatTickerManagerImpl{repo: repo}
}

// ---- Ticker ----

func (m *BarkatTickerManagerImpl) CreateTicker(ctx context.Context, ticker *barkat.Ticker) common.HttpError {
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		return m.repo.Create(c, ticker)
	})
}

func (m *BarkatTickerManagerImpl) GetTicker(ctx context.Context, ticker string) (result barkat.Ticker, httpErr common.HttpError) {
	httpErr = m.repo.GetByExternalId(ctx, ticker, &result)
	return
}

func (m *BarkatTickerManagerImpl) ListTickers(ctx context.Context, query barkat.TickerQuery) (barkat.TickerList, common.HttpError) {
	tickers, total, httpErr := m.repo.ListTickers(ctx, query)
	if httpErr != nil {
		return barkat.TickerList{}, httpErr
	}
	return barkat.TickerList{
		Tickers: tickers,
		Metadata: common.PaginatedResponse{
			Total:  total,
			Offset: query.Offset,
			Limit:  query.Limit,
		},
	}, nil
}

func (m *BarkatTickerManagerImpl) UpdateTicker(ctx context.Context, ticker string, req barkat.TickerUpdateRequest) (updatedTicker barkat.Ticker, httpErr common.HttpError) {
	httpErr = m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Fetch existing ticker
		existing, httpErr := m.GetTicker(c, ticker)
		if httpErr != nil {
			return httpErr
		}

		// Copy mutable fields from update request
		if req.Exchange != nil {
			existing.Exchange = req.Exchange
		}
		existing.Timeframes = req.Timeframes
		existing.Type = req.Type
		existing.State = req.State
		existing.Trend = req.Trend
		existing.IsFNO = req.IsFNO

		// Save updated ticker
		if httpErr := m.repo.Update(c, &existing); httpErr != nil {
			return httpErr
		}

		updatedTicker = existing
		return nil
	})
	return
}

func (m *BarkatTickerManagerImpl) PatchTickerLastOpened(ctx context.Context, ticker string, update barkat.TickerLastOpenedUpdate) (updatedTicker barkat.Ticker, httpErr common.HttpError) {
	httpErr = m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Fetch existing ticker
		existing, httpErr := m.GetTicker(c, ticker)
		if httpErr != nil {
			return httpErr
		}

		// Update only last_opened_at
		existing.LastOpenedAt = update.LastOpenedAt

		// Save updated ticker
		if httpErr := m.repo.Update(c, &existing); httpErr != nil {
			return httpErr
		}

		updatedTicker = existing
		return nil
	})
	return
}

func (m *BarkatTickerManagerImpl) DeleteTicker(ctx context.Context, ticker string) common.HttpError {
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Fetch existing ticker to get internal ID
		existing, httpErr := m.GetTicker(c, ticker)
		if httpErr != nil {
			return httpErr
		}
		// Delete by internal ID — cascade to AlertTickers/PriceAlerts is handled
		// by DB foreign key constraints (ON DELETE CASCADE in production DBs).
		return m.repo.DeleteById(c, existing.ID, &barkat.Ticker{})
	})
}
