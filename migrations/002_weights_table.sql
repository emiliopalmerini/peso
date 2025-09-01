-- Weights table
CREATE TABLE IF NOT EXISTS weights (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    value REAL NOT NULL,
    unit TEXT NOT NULL CHECK (unit IN ('kg', 'lb')),
    measured_at DATETIME NOT NULL,
    notes TEXT DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
)