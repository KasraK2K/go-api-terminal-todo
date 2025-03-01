CREATE TABLE IF NOT EXISTS todos
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT NOT NULL,
    description TEXT,
    completed   BOOLEAN DEFAULT FALSE

);