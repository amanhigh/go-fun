package barkat

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

// PriceAlert represents a local alert record imported from alertRepo (PRD Section 2.1.4).
// Alerts are linked through resolved Alert ticker rows. Import is deferred from Phase I.
type PriceAlert struct {
	ID            uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	AlertTickerID uint64    `gorm:"column:alert_ticker_id;not null;index:idx_price_alert_owner_price" json:"-"`
	AlertID       string    `gorm:"column:alert_id;uniqueIndex;not null" json:"alert_id" binding:"required,min=1,max=128"`
	TriggerPrice  float64   `gorm:"column:trigger_price;not null;type:decimal(18,6);index:idx_price_alert_owner_price,sort:asc" json:"trigger_price"`
	CreatedAt     time.Time `gorm:"column:created_at;not null" json:"created_at"`

	// Parent relation (not loaded by default in Phase I)
	AlertTicker *AlertTicker `gorm:"foreignKey:AlertTickerID;references:ID" json:"-"`
}

// TableName maps PriceAlert to the PRD-defined price_alerts table.
func (PriceAlert) TableName() string { return "price_alerts" }

// PriceAlertList is the paginated response for Price alerts.
type PriceAlertList struct {
	PriceAlerts []PriceAlert             `json:"price_alerts"`
	Metadata    common.PaginatedResponse `json:"metadata"`
}
