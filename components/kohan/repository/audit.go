package repository

import (
	"context"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"gorm.io/gorm"
)

// AuditRepository provides read operations required by audit plugins.
type AuditRepository interface {
	util.BaseDbRepositoryInterface
	// ListAlertCoverageRows returns coverage counts for active non-watched tickers.
	ListAlertCoverageRows(ctx context.Context) ([]barkat.AlertCoverageRow, common.HttpError)
}

type AuditRepositoryImpl struct {
	util.BaseDbRepository
}

var _ AuditRepository = (*AuditRepositoryImpl)(nil)

// NewAuditRepository creates a new AuditRepository backed by GORM.
func NewAuditRepository(db *gorm.DB) *AuditRepositoryImpl {
	return &AuditRepositoryImpl{BaseDbRepository: util.NewBaseDbRepository(db)}
}

func (r *AuditRepositoryImpl) ListAlertCoverageRows(ctx context.Context) ([]barkat.AlertCoverageRow, common.HttpError) {
	var rows []barkat.AlertCoverageRow
	err := r.SafeTx(ctx).Model(&barkat.Ticker{}).
		Select(`tickers.external_id AS ticker,
			COUNT(DISTINCT alert_tickers.id) AS alert_ticker_count,
			COUNT(price_alerts.id) AS price_alert_count`).
		Joins("LEFT JOIN alert_tickers ON alert_tickers.ticker_id = tickers.id").
		Joins("LEFT JOIN price_alerts ON price_alerts.alert_ticker_id = alert_tickers.id").
		Where("tickers.state = ?", "READY").
		Group("tickers.id, tickers.external_id").
		Order("tickers.external_id ASC").
		Scan(&rows).Error
	return rows, util.GormErrorMapper(err)
}
