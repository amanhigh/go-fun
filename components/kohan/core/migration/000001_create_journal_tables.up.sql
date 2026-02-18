CREATE TABLE IF NOT EXISTS journal_entries (
    id TEXT PRIMARY KEY,
    ticker TEXT NOT NULL,
    sequence TEXT NOT NULL,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_entry_ticker_created ON journal_entries (ticker, created_at DESC);
CREATE TABLE IF NOT EXISTS journal_images (
    id TEXT PRIMARY KEY,
    entry_id TEXT NOT NULL REFERENCES journal_entries(id),
    timeframe TEXT NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_image_entry_timeframe ON journal_images (entry_id, timeframe);

CREATE TABLE IF NOT EXISTS journal_tags (
    id TEXT PRIMARY KEY,
    entry_id TEXT NOT NULL REFERENCES journal_entries(id),
    tag TEXT NOT NULL,
    type TEXT NOT NULL,
    override TEXT,
    created_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tag_entry_type ON journal_tags (entry_id, type);
CREATE INDEX IF NOT EXISTS idx_tag_type_value ON journal_tags (type, tag);

CREATE TABLE IF NOT EXISTS journal_notes (
    id TEXT PRIMARY KEY,
    entry_id TEXT NOT NULL REFERENCES journal_entries(id),
    status TEXT NOT NULL,
    content TEXT NOT NULL,
    format TEXT NOT NULL DEFAULT 'markdown',
    created_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_note_entry_status ON journal_notes (entry_id, status);
