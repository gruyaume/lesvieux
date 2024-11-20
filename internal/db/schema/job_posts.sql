CREATE TABLE IF NOT EXISTS job_posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft',
    account_id INTEGER NOT NULL,
    FOREIGN KEY(account_id) REFERENCES accounts(account_id)
);