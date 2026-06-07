-- Create tickers table for TradingView-side ticker identity (PRD Section 2.1.2)
CREATE TABLE IF NOT EXISTS tickers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    external_id VARCHAR(50) UNIQUE NOT NULL,
    exchange VARCHAR(15) NOT NULL,
    timeframes TEXT NOT NULL DEFAULT '["MN","WK","DL"]',
    type VARCHAR(15) NOT NULL CHECK (type IN ('EQUITY', 'INDEX', 'CRYPTO', 'COMMODITY', 'FX', 'BOND', 'COMPOSITE')),
    state VARCHAR(10) NOT NULL DEFAULT 'WATCHED' CHECK (state IN ('WATCHED', 'READY', 'BLACKLIST')),
    trend VARCHAR(10) NOT NULL CHECK (trend IN ('UPTREND', 'SIDEWAYS', 'DOWNTREND')),
    last_opened_at DATETIME NOT NULL,
    is_fno INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for tickers
CREATE UNIQUE INDEX IF NOT EXISTS idx_ticker_external_id ON tickers (external_id);
CREATE INDEX IF NOT EXISTS idx_ticker_exchange ON tickers (exchange);
CREATE INDEX IF NOT EXISTS idx_ticker_type ON tickers (type);
CREATE INDEX IF NOT EXISTS idx_ticker_state ON tickers (state);
CREATE INDEX IF NOT EXISTS idx_ticker_trend ON tickers (trend);
CREATE INDEX IF NOT EXISTS idx_ticker_last_opened_at ON tickers (last_opened_at);
CREATE INDEX IF NOT EXISTS idx_ticker_is_fno ON tickers (is_fno);

-- Create alert_tickers table for Investing-side ticker identity (PRD Section 2.1.3)
CREATE TABLE IF NOT EXISTS alert_tickers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticker_id INTEGER NOT NULL,
    external_id VARCHAR(32) UNIQUE NOT NULL,
    pair_id VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    exchange VARCHAR(15),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ticker_id) REFERENCES tickers(id) ON DELETE CASCADE
);

-- Indexes for alert_tickers
CREATE UNIQUE INDEX IF NOT EXISTS idx_alert_ticker_external_id ON alert_tickers (external_id);
CREATE INDEX IF NOT EXISTS idx_alert_ticker_parent ON alert_tickers (ticker_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_alert_ticker_pair_id ON alert_tickers (pair_id);
CREATE INDEX IF NOT EXISTS idx_alert_ticker_exchange ON alert_tickers (exchange);

-- Create price_alerts table for local alert records (PRD Section 2.1.4)
CREATE TABLE IF NOT EXISTS price_alerts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alert_ticker_id INTEGER NOT NULL,
    alert_id VARCHAR(128) UNIQUE,
    trigger_price DECIMAL(18,6) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (alert_ticker_id) REFERENCES alert_tickers(id) ON DELETE CASCADE
);

-- Indexes for price_alerts
CREATE UNIQUE INDEX IF NOT EXISTS idx_price_alert_alert_id ON price_alerts (alert_id);
CREATE INDEX IF NOT EXISTS idx_price_alert_owner_price ON price_alerts (alert_ticker_id, trigger_price);
