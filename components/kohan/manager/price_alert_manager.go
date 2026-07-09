package manager

import (
	"context"
	"fmt"
	"net/http"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
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

// ---- Price Alerts ----

func (m *PriceAlertManagerImpl) ReplacePriceAlerts(ctx context.Context, request barkat.PriceAlertReplaceRequest) (result barkat.PriceAlertReplaceResult, httpErr common.HttpError) {
	httpErr = m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		alertTickerByPairID, err := m.resolveAlertTickersForInputs(c, request.Alerts)
		if err != nil {
			return err
		}

		pairIDs := m.collectPairIDs(alertTickerByPairID)
		alerts := m.buildReplacementAlerts(request.Alerts, alertTickerByPairID)

		// Delete existing alerts for submitted pair IDs, then create replacements.
		if err := m.repo.DeleteByPairIDs(c, pairIDs); err != nil {
			return err
		}
		if err := m.repo.CreateAlerts(c, alerts); err != nil {
			return err
		}

		result = barkat.PriceAlertReplaceResult{PairsReplaced: len(alertTickerByPairID), AlertsCreated: len(alerts)}
		return nil
	})
	return result, httpErr
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

// ---- Private Replacement Helpers ----

func (m *PriceAlertManagerImpl) collectPairIDs(alertTickerByPairID map[string]barkat.AlertTicker) []string {
	return lo.Keys(alertTickerByPairID)
}

func (m *PriceAlertManagerImpl) buildReplacementAlerts(inputs []barkat.PriceAlertInput, alertTickerByPairID map[string]barkat.AlertTicker) []barkat.PriceAlert {
	alerts := make([]barkat.PriceAlert, 0, len(inputs))
	for _, input := range inputs {
		alertID := input.AlertID
		alerts = append(alerts, barkat.PriceAlert{
			AlertTickerID: alertTickerByPairID[input.PairID].ID,
			AlertID:       &alertID,
			TriggerPrice:  input.TriggerPrice,
		})
	}
	return alerts
}

func (m *PriceAlertManagerImpl) resolveAlertTickersForInputs(ctx context.Context, inputs []barkat.PriceAlertInput) (map[string]barkat.AlertTicker, common.HttpError) {
	alertTickerByPairID := make(map[string]barkat.AlertTicker)

	for _, input := range inputs {
		if _, ok := alertTickerByPairID[input.PairID]; ok {
			continue
		}
		alertTicker, err := m.repo.GetByPairId(ctx, input.PairID, string(barkat.AlertTickerTypePrimary))
		if err != nil {
			// Try secondary lookup to provide a readable message
			msg := buildUnresolvedAlertMessage(ctx, m.repo, input)
			log.Error().
				Str("pair_id", input.PairID).
				Str("alert_id", input.AlertID).
				Float64("trigger_price", input.TriggerPrice).
				Msg("ReplacePriceAlerts — unresolved alert ticker")
			return nil, common.NewHttpError(msg, http.StatusNotFound)
		}
		alertTickerByPairID[input.PairID] = alertTicker
	}

	return alertTickerByPairID, nil
}

// buildUnresolvedAlertMessage attempts to find any alert ticker by pair_id to provide
// a readable error message. If an alert ticker exists (e.g. SECONDARY), its human-readable
// fields are included. If none exists, only the incoming request fields are shown.
func buildUnresolvedAlertMessage(ctx context.Context, repo repository.PriceAlertRepository, input barkat.PriceAlertInput) string {
	anyTicker, anyErr := repo.GetByPairId(ctx, input.PairID)
	if anyErr == nil {
		exchange := ""
		if anyTicker.Exchange != nil {
			exchange = *anyTicker.Exchange
		}
		return fmt.Sprintf(
			"PRIMARY alert ticker not found for pair_id=%s, alert_id=%s, trigger_price=%.4f. "+
				"Found %s alert ticker instead: symbol=%s, name=%q, exchange=%s.",
			input.PairID, input.AlertID, input.TriggerPrice,
			anyTicker.Type, anyTicker.Symbol, anyTicker.Name, exchange,
		)
	}
	return fmt.Sprintf(
		"PRIMARY alert ticker not found for pair_id=%s, alert_id=%s, trigger_price=%.4f. "+
			"No alert ticker mapping exists in Barkat.",
		input.PairID, input.AlertID, input.TriggerPrice,
	)
}
