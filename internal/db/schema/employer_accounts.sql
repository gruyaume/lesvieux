CREATE TABLE IF NOT EXISTS employer_accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
    employer_id INTEGER NOT NULL,
    FOREIGN KEY (employer_id) REFERENCES employers(id) ON DELETE CASCADE
);