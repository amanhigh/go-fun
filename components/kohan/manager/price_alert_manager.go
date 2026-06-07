package manager

import (
	"context"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// PriceAlertManager provides business logic for price alert operations.
type PriceAlertManager interface {
	// ReplacePriceAlerts replaces all price alerts for pair IDs included in the request.
	ReplacePriceAlerts(ctx context.Context, request barkat.PriceAlertReplaceRequest) (barkat.PriceAlertReplaceResult, common.HttpError)
	// CreatePendingPriceAlert creates a local pending alert for a primary ticker.
	CreatePendingPriceAlert(ctx context.Context, ticker string, request barkat.PendingPriceAlertRequest) (barkat.PriceAlert, common.HttpError)
	// DeletePriceAlert deletes one canonical alert by alert id.
	DeletePriceAlert(ctx context.Context, alertID string) common.HttpError
	// ListPriceAlerts returns a filtered, sorted, paginated list of price alerts.
	ListPriceAlerts(ctx context.Context, query barkat.PriceAlertQuery) (barkat.PriceAlertList, common.HttpError)
}

type PriceAlertManagerImpl struct {
	repo repository.PriceAlertRepository
}

var _ PriceAlertManager = (*PriceAlertManagerImpl)(nil)

// NewPriceAlertManager creates a new PriceAlertManager.
func NewPriceAlertManager(repo repository.PriceAlertRepository) *PriceAlertManagerImpl {
	return &PriceAlertManagerImpl{repo: repo}
}

func (m *PriceAlertManagerImpl) ReplacePriceAlerts(ctx context.Context, request barkat.PriceAlertReplaceRequest) (result barkat.PriceAlertReplaceResult, httpErr common.HttpError) {
	httpErr = m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		alertTickerByPairID, err := m.resolveAlertTickersForInputs(c, request.Alerts)
		if err != nil {
			return err
		}

		alerts := make([]barkat.PriceAlert, 0, len(request.Alerts))
		for _, input := range request.Alerts {
			alertID := input.AlertID
			alerts = append(alerts, barkat.PriceAlert{
				AlertTickerID: alertTickerByPairID[input.PairID].ID,
				AlertID:       &alertID,
				TriggerPrice:  input.TriggerPrice,
			})
		}

		if err := m.repo.ReplaceAlerts(c, alerts); err != nil {
			return err
		}
		result = barkat.PriceAlertReplaceResult{PairsReplaced: len(alertTickerByPairID), AlertsCreated: len(alerts)}
		return nil
	})
	return result, httpErr
}

func (m *PriceAlertManagerImpl) resolveAlertTickersForInputs(ctx context.Context, inputs []barkat.PriceAlertInput) (map[string]barkat.AlertTicker, common.HttpError) {
	alertTickerByPairID := make(map[string]barkat.AlertTicker)

	for _, input := range inputs {
		if _, ok := alertTickerByPairID[input.PairID]; ok {
			continue
		}
		alertTicker, err := m.repo.GetByPairId(ctx, input.PairID)
		if err != nil {
			return nil, err
		}
		alertTickerByPairID[input.PairID] = alertTicker
	}

	return alertTickerByPairID, nil
}

func (m *PriceAlertManagerImpl) CreatePendingPriceAlert(ctx context.Context, ticker string, request barkat.PendingPriceAlertRequest) (result barkat.PriceAlert, httpErr common.HttpError) {
	httpErr = m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		alertTicker, err := m.repo.GetByTicker(c, ticker)
		if err != nil {
			return err
		}

		alert := barkat.PriceAlert{AlertTickerID: alertTicker.ID, TriggerPrice: request.TriggerPrice}
		if err := m.repo.Create(c, &alert); err != nil {
			return err
		}
		result = alert
		result.PairID = alertTicker.PairID
		return nil
	})
	return result, httpErr
}

func (m *PriceAlertManagerImpl) DeletePriceAlert(ctx context.Context, alertID string) common.HttpError {
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		return m.repo.DeleteByAlertID(c, alertID)
	})
}

func (m *PriceAlertManagerImpl) ListPriceAlerts(ctx context.Context, query barkat.PriceAlertQuery) (barkat.PriceAlertList, common.HttpError) {
	if query.Ticker != "" {
		if httpErr := m.repo.GetByExternalId(ctx, query.Ticker, &barkat.Ticker{}); httpErr != nil {
			return barkat.PriceAlertList{}, httpErr
		}
	}

	alerts, total, httpErr := m.repo.ListPriceAlerts(ctx, query)
	if httpErr != nil {
		return barkat.PriceAlertList{}, httpErr
	}
	return barkat.PriceAlertList{
		PriceAlerts: alerts,
		Metadata: common.PaginatedResponse{
			Total:  total,
			Offset: query.Offset,
			Limit:  query.Limit,
		},
	}, nil
}
