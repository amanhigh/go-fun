CREATE TABLE IF NOT EXISTS entries (
    id TEXT PRIMARY KEY,
    ticker TEXT NOT NULL,
    sequence TEXT NOT NULL,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_entry_ticker_created ON entries (ticker, created_at DESC);
CREATE TABLE IF NOT EXISTS images (
    id TEXT PRIMARY KEY,
    entry_id TEXT NOT NULL REFERENCES entries(id),
    timeframe TEXT NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_image_entry_timeframe ON images (entry_id, timeframe);

CREATE TABLE IF NOT EXISTS tags (
    id TEXT PRIMARY KEY,
    entry_id TEXT NOT NULL REFERENCES entries(id),
    tag TEXT NOT NULL,
    type TEXT NOT NULL,
    override TEXT,
    created_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tag_entry_type ON tags (entry_id, type);
CREATE INDEX IF NOT EXISTS idx_tag_type_value ON tags (type, tag);

CREATE TABLE IF NOT EXISTS notes (
    id TEXT PRIMARY KEY,
    entry_id TEXT NOT NULL REFERENCES entries(id),
    status TEXT NOT NULL,
    content TEXT NOT NULL,
    format TEXT NOT NULL DEFAULT 'markdown',
    created_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_note_entry_status ON notes (entry_id, status);
