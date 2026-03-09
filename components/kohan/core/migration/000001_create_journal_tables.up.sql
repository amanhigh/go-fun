CREATE TABLE IF NOT EXISTS journals (
    id TEXT PRIMARY KEY,
    ticker TEXT NOT NULL,
    sequence TEXT NOT NULL,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_journal_ticker_created ON journals (ticker, created_at DESC);
CREATE TABLE IF NOT EXISTS images (
    id TEXT PRIMARY KEY,
    journal_id TEXT NOT NULL REFERENCES journals(id),
    timeframe TEXT NOT NULL,
    created_at DATETIME NOT NULL
);


CREATE TABLE IF NOT EXISTS tags (
    id TEXT PRIMARY KEY,
    journal_id TEXT NOT NULL REFERENCES journals(id),
    tag TEXT NOT NULL,
    type TEXT NOT NULL,
    override TEXT,
    created_at DATETIME NOT NULL
);


CREATE TABLE IF NOT EXISTS notes (
    id TEXT PRIMARY KEY,
    journal_id TEXT NOT NULL REFERENCES journals(id),
    status TEXT NOT NULL,
    content TEXT NOT NULL,
    format TEXT NOT NULL DEFAULT 'markdown',
    created_at DATETIME NOT NULL
);

