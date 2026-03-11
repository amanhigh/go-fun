-- Create journals table with BIGINT primary keys and external_id
CREATE TABLE IF NOT EXISTS journals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    external_id VARCHAR(26) UNIQUE NOT NULL,
    ticker VARCHAR(10) NOT NULL,
    sequence VARCHAR(3) NOT NULL CHECK (sequence IN ('MWD', 'YR')),
    type VARCHAR(8) NOT NULL CHECK (type IN ('REJECTED', 'RESULT', 'SET')),
    status VARCHAR(10) NOT NULL CHECK (status IN ('SET', 'RUNNING', 'DROPPED', 'TAKEN', 'REJECTED', 'SUCCESS', 'FAIL', 'MISSED', 'JUST_LOSS', 'BROKEN')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- Create indexes for journals
CREATE UNIQUE INDEX IF NOT EXISTS idx_journal_external_id ON journals (external_id);
CREATE INDEX IF NOT EXISTS idx_journal_ticker ON journals (ticker);
CREATE INDEX IF NOT EXISTS idx_journal_created_at ON journals (created_at);

-- Create images table
CREATE TABLE IF NOT EXISTS images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    external_id VARCHAR(26) UNIQUE NOT NULL,
    journal_id INTEGER NOT NULL,
    timeframe VARCHAR(3) NOT NULL CHECK (timeframe IN ('DL', 'WK', 'MN', 'TMN', 'SMN', 'YR')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (journal_id) REFERENCES journals(id) ON DELETE CASCADE
);

-- Create indexes for images
CREATE UNIQUE INDEX IF NOT EXISTS idx_image_external_id ON images (external_id);
CREATE INDEX IF NOT EXISTS idx_image_journal_id ON images (journal_id);

-- Create tags table
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    external_id VARCHAR(26) UNIQUE NOT NULL,
    journal_id INTEGER NOT NULL,
    tag VARCHAR(10) NOT NULL,
    type VARCHAR(9) NOT NULL CHECK (type IN ('REASON', 'MANAGEMENT')),
    override VARCHAR(5),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (journal_id) REFERENCES journals(id) ON DELETE CASCADE
);

-- Create indexes for tags
CREATE UNIQUE INDEX IF NOT EXISTS idx_tag_external_id ON tags (external_id);
CREATE INDEX IF NOT EXISTS idx_tag_journal_id ON tags (journal_id);

-- Create notes table
CREATE TABLE IF NOT EXISTS notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    external_id VARCHAR(26) UNIQUE NOT NULL,
    journal_id INTEGER NOT NULL,
    status VARCHAR(10) NOT NULL CHECK (status IN ('SET', 'RUNNING', 'DROPPED', 'TAKEN', 'REJECTED', 'SUCCESS', 'FAIL', 'MISSED', 'JUST_LOSS', 'BROKEN')),
    content TEXT NOT NULL,
    format VARCHAR(10) NOT NULL CHECK (format IN ('MARKDOWN', 'PLAINTEXT')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (journal_id) REFERENCES journals(id) ON DELETE CASCADE
);

-- Create indexes for notes
CREATE UNIQUE INDEX IF NOT EXISTS idx_note_external_id ON notes (external_id);
CREATE INDEX IF NOT EXISTS idx_note_journal_id ON notes (journal_id);
CREATE INDEX IF NOT EXISTS idx_note_status ON notes (status);

