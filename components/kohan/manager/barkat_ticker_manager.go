package manager

import (
	"context"
	"fmt"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// BarkatTickerManager provides business logic for barkat ticker CRUD operations.
type BarkatTickerManager interface {
	// CreateTicker creates a new ticker.
	CreateTicker(ctx context.Context, ticker *barkat.Ticker) common.HttpError
	// GetTicker retrieves a single ticker by tv_symbol with AlertTickers preloaded.
	GetTicker(ctx context.Context, tvSymbol string) (barkat.Ticker, common.HttpError)
	// ListTickers returns a filtered, paginated list of tickers.
	ListTickers(ctx context.Context, query barkat.TickerQuery) (barkat.TickerList, common.HttpError)
	// UpdateTicker replaces mutable fields on an existing ticker.
	UpdateTicker(ctx context.Context, tvSymbol string, req barkat.TickerUpdateRequest) (barkat.Ticker, common.HttpError)
	// PatchTickerLastOpened updates only the last_opened_at timestamp.
	PatchTickerLastOpened(ctx context.Context, tvSymbol string, update barkat.TickerLastOpenedUpdate) (barkat.Ticker, common.HttpError)
	// DeleteTicker deletes a ticker and cascades to linked AlertTickers and PriceAlerts.
	DeleteTicker(ctx context.Context, tvSymbol string) common.HttpError
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

func (m *BarkatTickerManagerImpl) GetTicker(ctx context.Context, tvSymbol string) (barkat.Ticker, common.HttpError) {
	ticker, err := m.repo.GetByTvSymbol(ctx, tvSymbol)
	if err != nil {
		return barkat.Ticker{}, common.ErrNotFound
	}
	// Compute the non-persisted alert_ticker_count from the preloaded association
	ticker.AlertTickerCount = int64(len(ticker.AlertTickers))
	return ticker, nil
}

func (m *BarkatTickerManagerImpl) ListTickers(ctx context.Context, query barkat.TickerQuery) (barkat.TickerList, common.HttpError) {
	tickers, total, err := m.repo.ListTickers(ctx, query)
	if err != nil {
		return barkat.TickerList{}, common.NewServerError(fmt.Errorf("failed to list tickers: %w", err))
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

func (m *BarkatTickerManagerImpl) UpdateTicker(ctx context.Context, tvSymbol string, req barkat.TickerUpdateRequest) (barkat.Ticker, common.HttpError) {
	var updatedTicker barkat.Ticker
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Fetch existing ticker
		ticker, httpErr := m.GetTicker(c, tvSymbol)
		if httpErr != nil {
			return httpErr
		}

		// Copy mutable fields from update request
		if req.Exchange != nil {
			ticker.Exchange = req.Exchange
		}
		ticker.Timeframes = req.Timeframes
		ticker.Type = req.Type
		ticker.State = req.State
		ticker.Trend = req.Trend
		ticker.IsFNO = req.IsFNO

		// Save updated ticker
		if httpErr := m.repo.Update(c, &ticker); httpErr != nil {
			return httpErr
		}

		updatedTicker = ticker
		return nil
	})
	if err != nil {
		return barkat.Ticker{}, err
	}
	return updatedTicker, nil
}

func (m *BarkatTickerManagerImpl) PatchTickerLastOpened(ctx context.Context, tvSymbol string, update barkat.TickerLastOpenedUpdate) (barkat.Ticker, common.HttpError) {
	var updatedTicker barkat.Ticker
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Fetch existing ticker
		ticker, httpErr := m.GetTicker(c, tvSymbol)
		if httpErr != nil {
			return httpErr
		}

		// Update only last_opened_at
		ticker.LastOpenedAt = update.LastOpenedAt

		// Save updated ticker
		if httpErr := m.repo.Update(c, &ticker); httpErr != nil {
			return httpErr
		}

		updatedTicker = ticker
		return nil
	})
	if err != nil {
		return barkat.Ticker{}, err
	}
	return updatedTicker, nil
}

func (m *BarkatTickerManagerImpl) DeleteTicker(ctx context.Context, tvSymbol string) common.HttpError {
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Fetch existing ticker to get internal ID
		ticker, httpErr := m.GetTicker(c, tvSymbol)
		if httpErr != nil {
			return httpErr
		}
		// Delete by internal ID — cascade to AlertTickers/PriceAlerts is handled
		// by DB foreign key constraints (ON DELETE CASCADE in production DBs).
		return m.repo.DeleteById(c, ticker.ID, &barkat.Ticker{})
	})
}
