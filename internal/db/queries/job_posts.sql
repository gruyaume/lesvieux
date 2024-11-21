-- name: ListJobPosts :many
SELECT * FROM job_posts
ORDER BY created_at DESC;

-- name: GetJobPost :one
SELECT * FROM job_posts
WHERE id = ? LIMIT 1;

-- name: ListJobPostsByAccount :many
SELECT * FROM job_posts
WHERE employer_id = ?
ORDER BY created_at DESC;

-- name: CreateJobPost :one
INSERT INTO job_posts (
  title, content, created_at, status, employer_id
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateJobPost :exec
UPDATE job_posts
set title = ?, content = ?, status = ?
WHERE id = ?;

-- name: DeleteJobPost :exec
DELETE FROM job_posts
WHERE id = ?;