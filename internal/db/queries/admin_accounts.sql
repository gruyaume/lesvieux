-- name: GetAdminAccount :one
SELECT * FROM admin_accounts
WHERE id = ? LIMIT 1;

-- name: GetAdminAccountByEmail :one
SELECT * FROM admin_accounts
WHERE email = ? LIMIT 1;

-- name: ListAdminAccounts :many
SELECT * FROM admin_accounts
ORDER BY email;

-- name: CreateAdminAccount :one
INSERT INTO admin_accounts (
  email, password_hash
) VALUES (
  ?, ?
)
RETURNING *;

-- name: UpdateAdminAccount :exec
UPDATE admin_accounts
set password_hash = ?
WHERE id = ?;

-- name: DeleteAdminAccount :exec
DELETE FROM admin_accounts
WHERE id = ?;

-- name: NumAdminAccounts :one
SELECT COUNT(*) FROM admin_accounts;