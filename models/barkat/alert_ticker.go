package barkat

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

// AlertTicker represents an Alert/Investing-side ticker attached to a TradingView ticker.
type AlertTicker struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	TickerID  uint64    `gorm:"column:ticker_id;not null;index:idx_alert_ticker_parent" json:"-"`
	Symbol    string    `gorm:"column:external_id;uniqueIndex;not null" json:"symbol" binding:"required,min=1,max=25,alert_symbol"`
	PairID    string    `gorm:"column:pair_id;not null;index:idx_alert_ticker_pair_id" json:"pair_id" binding:"required,min=1,max=64,numeric"`
	Name      string    `gorm:"column:name;not null" json:"name" binding:"required,min=1,max=100,alert_name"`
	Exchange  *string   `gorm:"column:exchange;index:idx_alert_ticker_exchange" json:"exchange" binding:"omitempty,min=1,max=10,alert_exchange"`
	Ticker    string    `gorm:"-" json:"ticker,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

// TableName maps AlertTicker to the PRD-defined alert_tickers table.
func (AlertTicker) TableName() string { return "alert_tickers" }

// PriceAlert represents a local Barkat price alert resolved to an Alert ticker.
type PriceAlert struct {
	ID            uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	AlertTickerID uint64    `gorm:"column:alert_ticker_id;not null;index:idx_price_alert_alert_ticker" json:"-"`
	AlertID       string    `gorm:"column:alert_id;uniqueIndex;not null" json:"alert_id"`
	TriggerPrice  float64   `gorm:"column:trigger_price;not null" json:"trigger_price"`
	Name          *string   `gorm:"column:name" json:"name"`
	CreatedAt     time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

// TableName maps PriceAlert to the PRD-defined price_alerts table.
func (PriceAlert) TableName() string { return "price_alerts" }

// AlertTickerPath binds the :symbol path parameter for Alert ticker APIs.
type AlertTickerPath struct {
	Symbol string `uri:"symbol" binding:"required,alert_symbol"`
}

// AlertTickerQuery holds query parameters for listing/filtering Alert tickers.
type AlertTickerQuery struct {
	common.Pagination
	Symbol   string `form:"symbol" binding:"omitempty,min=1,max=25,alert_symbol"`
	Ticker   string `form:"ticker" binding:"omitempty,tv_ticker_path"`
	PairID   string `form:"pair-id" binding:"omitempty,min=1,max=64,numeric"`
	Exchange string `form:"exchange" binding:"omitempty,min=1,max=10,alert_exchange"`
}

// NewAlertTickerQuery creates an AlertTickerQuery struct with default pagination values.
func NewAlertTickerQuery() AlertTickerQuery {
	return AlertTickerQuery{Pagination: common.Pagination{}}
}

// AlertTickerList is the paginated response for Alert tickers.
type AlertTickerList struct {
	AlertTickers []AlertTicker            `json:"alert_tickers"`
	Metadata     common.PaginatedResponse `json:"metadata"`
}
