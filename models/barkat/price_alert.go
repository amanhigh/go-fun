package barkat

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"gorm.io/gorm"
)

// PriceAlert API route constants.
const (
	PriceAlertBase         = common.APIV1 + "/alerts"
	MaxPriceAlertBatchSize = 100
	DefaultPriceAlertLimit = 10
)

// PriceAlert represents a local alert record imported from alertRepo (PRD Section 2.1.4).
// Alerts are linked through resolved Alert ticker rows. Import is deferred from Phase I.
type PriceAlert struct {
	ID            uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	AlertTickerID uint64    `gorm:"column:alert_ticker_id;not null;index:idx_price_alert_owner_price" json:"-"`
	AlertID       *string   `gorm:"column:alert_id;uniqueIndex" json:"alert_id,omitempty"`
	PairID        string    `gorm:"-" json:"pair_id,omitempty"`
	TriggerPrice  float64   `gorm:"column:trigger_price;not null;type:decimal(18,6);index:idx_price_alert_owner_price,sort:asc" json:"trigger_price"`
	CreatedAt     time.Time `gorm:"column:created_at;not null" json:"created_at"`

	// Parent relation (not loaded by default in Phase I)
	AlertTicker *AlertTicker `gorm:"foreignKey:AlertTickerID;references:ID" json:"-"`
}

// TableName maps PriceAlert to the PRD-defined price_alerts table.
func (PriceAlert) TableName() string { return "price_alerts" }

// AfterFind populates denormalized response fields from the preloaded AlertTicker relation.
func (p *PriceAlert) AfterFind(_ *gorm.DB) error {
	if p.AlertTicker != nil {
		p.PairID = p.AlertTicker.PairID
	}
	return nil
}

// PriceAlertInput is one canonical refreshed Investing.com price-alert row.
type PriceAlertInput struct {
	PairID       string  `json:"pair_id" binding:"required,min=1,max=64,number"`
	AlertID      string  `json:"alert_id" binding:"required,min=1,max=128,number"`
	TriggerPrice float64 `json:"trigger_price" binding:"required,gt=0"`
}

// PriceAlertReplaceRequest replaces all alerts for pair IDs included in Alerts.
type PriceAlertReplaceRequest struct {
	Alerts []PriceAlertInput `json:"alerts" binding:"required,dive"`
}

// PriceAlertReplaceResult summarizes replacement counts.
type PriceAlertReplaceResult struct {
	PairsReplaced int `json:"pairs_replaced"`
	AlertsCreated int `json:"alerts_created"`
}

// PendingPriceAlertRequest creates a local pending alert without a canonical Investing alert id.
type PendingPriceAlertRequest struct {
	TriggerPrice float64 `json:"trigger_price" binding:"required,gt=0"`
}

// PriceAlertPath binds the :alert-id path parameter for price alert APIs.
type PriceAlertPath struct {
	AlertID string `uri:"alert-id" binding:"required,min=1,max=128,number"`
}

// PriceAlertQuery holds query parameters for listing price alerts.
type PriceAlertQuery struct {
	Offset    int              `form:"offset,default=0" binding:"min=0"`
	Limit     int              `form:"limit,default=10" binding:"min=1,max=10"`
	SortOrder common.SortOrder `form:"sort-order,default=asc" binding:"omitempty,oneof=asc desc"`
	Ticker    string           `form:"ticker" binding:"omitempty,ticker_path"`
	SortBy    string           `form:"sort-by,default=trigger_price" binding:"omitempty,oneof=trigger_price created_at"`
}

// PriceAlertList is the paginated response for Price alerts.
type PriceAlertList struct {
	PriceAlerts []PriceAlert             `json:"alerts"`
	Metadata    common.PaginatedResponse `json:"metadata"`
}
