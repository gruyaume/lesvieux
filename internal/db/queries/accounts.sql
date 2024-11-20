-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1;

-- name: GetAccountByUsername :one
SELECT * FROM accounts
WHERE username = ? LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY username;

-- name: CreateAccount :one
INSERT INTO accounts (
  username, password_hash, role
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: UpdateAccount :exec
UPDATE accounts
set password_hash = ?
WHERE id = ?;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = ?;

-- name: NumAccounts :one
SELECT COUNT(*) FROM accounts;
