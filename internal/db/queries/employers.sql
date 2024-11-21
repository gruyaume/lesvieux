-- name: GetEmployer :one
SELECT * FROM employers
WHERE id = ? LIMIT 1;

-- name: ListEmployers :many
SELECT * FROM employers
ORDER BY name;

-- name: CreateEmployer :one
INSERT INTO employers (
  name
) VALUES (
  ?
)
RETURNING *;

-- name: UpdateEmployer :exec
UPDATE employers
SET name = ?
WHERE id = ?;

-- name: DeleteEmployer :exec
DELETE FROM employers
WHERE id = ?;

-- name: NumEmployers :one
SELECT COUNT(*) FROM employers;
