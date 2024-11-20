
-- name: ListPublicBlogPosts :many
SELECT * FROM blog_posts
WHERE status = 'published'
ORDER BY created_at DESC;

-- name: GetPublicBlogPost :one
SELECT * FROM blog_posts
WHERE id = ? AND status = 'published'
LIMIT 1;

-- name: ListBlogPosts :many
SELECT * FROM blog_posts
ORDER BY created_at DESC;

-- name: GetBlogPost :one
SELECT * FROM blog_posts
WHERE id = ? LIMIT 1;

-- name: ListPublicBlogPostsByAccount :many
SELECT * FROM blog_posts
WHERE status = 'published' AND account_id = ?
ORDER BY created_at DESC;

-- name: ListBlogPostsByAccount :many
SELECT * FROM blog_posts
WHERE account_id = ?
ORDER BY created_at DESC;

-- name: CreateBlogPost :one
INSERT INTO blog_posts (
  title, content, created_at, status, account_id
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateBlogPost :exec
UPDATE blog_posts
set title = ?, content = ?, status = ?
WHERE id = ?;

-- name: DeleteBlogPost :exec
DELETE FROM blog_posts
WHERE id = ?;