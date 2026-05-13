package barkat

import (
	"time"

	"github.com/amanhigh/go-fun/models/common"
)

// Ticker API route constants.
const (
	TickerBase      = common.APIV1 + "/tickers"
	AlertTickerBase = common.APIV1 + "/alert-tickers"
)

// Ticker represents a TradingView-side ticker imported from Barkat tickerRepo.
type Ticker struct {
	ID               uint64        `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	Ticker           string        `gorm:"column:external_id;uniqueIndex;not null" json:"ticker" binding:"required,min=1,max=50,ticker"`
	Exchange         *string       `gorm:"column:exchange;index:idx_ticker_exchange" json:"exchange" binding:"omitempty,min=1,max=10,ticker_exchange"`
	Timeframes       []string      `gorm:"column:timeframes;serializer:json;not null" json:"timeframes" binding:"required,min=1,max=6,dive,oneof=YR SMN TMN MN WK DL"`
	Type             string        `gorm:"column:type;not null;index:idx_ticker_type" json:"type" binding:"required,oneof=EQUITY INDEX CRYPTO COMMODITY FX BOND COMPOSITE"`
	State            string        `gorm:"column:state;not null;default:WATCHED;index:idx_ticker_state" json:"state" binding:"required,oneof=WATCHED READY BLACKLIST"`
	Trend            string        `gorm:"column:trend;not null;index:idx_ticker_trend" json:"trend" binding:"required,oneof=UPTREND SIDEWAYS DOWNTREND"`
	LastOpenedAt     time.Time     `gorm:"column:last_opened_at;not null;index:idx_ticker_last_opened_at" json:"last_opened_at" binding:"required"`
	IsFNO            bool          `gorm:"column:is_fno;not null;default:false;index:idx_ticker_is_fno" json:"is_fno"`
	CreatedAt        time.Time     `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt        time.Time     `gorm:"column:updated_at;not null" json:"updated_at"`
	AlertTickers     []AlertTicker `gorm:"foreignKey:TickerID;references:ID" json:"alert_tickers,omitempty"`
	AlertTickerCount int64         `gorm:"-" json:"alert_ticker_count,omitempty"`
}

// TableName maps Ticker to the PRD-defined tradingview_tickers table.
func (Ticker) TableName() string { return "tradingview_tickers" }

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

// TickerPath binds the :ticker path parameter.
type TickerPath struct {
	Ticker string `uri:"ticker" binding:"required,tv_ticker_path"`
}

// AlertTickerPath binds the :symbol path parameter for Alert ticker APIs.
type AlertTickerPath struct {
	Symbol string `uri:"symbol" binding:"required,alert_symbol"`
}

// TickerQuery holds query parameters for listing/filtering tickers.
type TickerQuery struct {
	common.Pagination
	Search      string `form:"search" binding:"omitempty,min=1,max=50"`
	Exchange    string `form:"exchange" binding:"omitempty,min=1,max=10,ticker_exchange"`
	Type        string `form:"type" binding:"omitempty,oneof=EQUITY INDEX CRYPTO COMMODITY FX BOND COMPOSITE"`
	State       string `form:"state" binding:"omitempty,oneof=WATCHED READY BLACKLIST"`
	Trend       string `form:"trend" binding:"omitempty,oneof=UPTREND SIDEWAYS DOWNTREND"`
	IsFNO       *bool  `form:"is-fno" binding:"omitempty"`
	OpenedAfter string `form:"opened-after" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	SortBy      string `form:"sort-by" binding:"omitempty,oneof=ticker exchange type state trend last_opened_at"`
	SortOrder   string `form:"sort-order" binding:"omitempty,oneof=asc desc"`
}

// NewTickerQuery creates a TickerQuery struct with default pagination values.
func NewTickerQuery() TickerQuery {
	return TickerQuery{Pagination: common.Pagination{}}
}

// TickerList is the paginated response for tickers.
type TickerList struct {
	Tickers  []Ticker                 `json:"tickers"`
	Metadata common.PaginatedResponse `json:"metadata"`
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

// TickerUpdateRequest represents PUT update ticker request body.
// Ticker field is excluded (comes from URL path).
type TickerUpdateRequest struct {
	Exchange   *string  `json:"exchange" binding:"omitempty,min=1,max=10,ticker_exchange"`
	Timeframes []string `json:"timeframes" binding:"required,min=1,max=6,dive,oneof=YR SMN TMN MN WK DL"`
	Type       string   `json:"type" binding:"required,oneof=EQUITY INDEX CRYPTO COMMODITY FX BOND COMPOSITE"`
	State      string   `json:"state" binding:"required,oneof=WATCHED READY BLACKLIST"`
	Trend      string   `json:"trend" binding:"required,oneof=UPTREND SIDEWAYS DOWNTREND"`
	IsFNO      *bool    `json:"is_fno" binding:"required"`
}

// TickerLastOpenedUpdate represents PATCH ticker last_opened_at request body.
type TickerLastOpenedUpdate struct {
	LastOpenedAt time.Time `json:"last_opened_at" binding:"required"`
}
