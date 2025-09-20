PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS weight (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  weight_gr INTEGER NOT NULL,
  logged_at TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_weight_user_time
  ON weight(user_id, logged_at DESC);