-- name: GetEmployerAccount :one
SELECT * FROM employer_accounts
where employer_id = ? and id = ? LIMIT 1;

-- name: GetEmployerAccountByEmail :one
SELECT * FROM employer_accounts
WHERE email = ? LIMIT 1;

-- name: ListEmployerAccounts :many
SELECT * FROM employer_accounts
where employer_id = ?
ORDER BY email;

-- name: CreateEmployerAccount :one
INSERT INTO employer_accounts (
  email, password_hash, employer_id
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: UpdateEmployerAccount :exec
UPDATE employer_accounts
set password_hash = ?
WHERE id = ?;

-- name: DeleteEmployerAccount :exec
DELETE FROM employer_accounts
where employer_id = ? and id = ?;

-- name: NumEmployerAccounts :one
SELECT COUNT(*) FROM employer_accounts;