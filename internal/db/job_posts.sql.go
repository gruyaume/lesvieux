// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: job_posts.sql

package db

import (
	"context"
)

const createJobPost = `-- name: CreateJobPost :one
INSERT INTO job_posts (
  title, content, created_at, status, employer_id
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING id, title, content, created_at, status, employer_id
`

type CreateJobPostParams struct {
	Title      string
	Content    string
	CreatedAt  string
	Status     string
	EmployerID int64
}

func (q *Queries) CreateJobPost(ctx context.Context, arg CreateJobPostParams) (JobPost, error) {
	row := q.db.QueryRowContext(ctx, createJobPost,
		arg.Title,
		arg.Content,
		arg.CreatedAt,
		arg.Status,
		arg.EmployerID,
	)
	var i JobPost
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Content,
		&i.CreatedAt,
		&i.Status,
		&i.EmployerID,
	)
	return i, err
}

const deleteJobPost = `-- name: DeleteJobPost :exec
DELETE FROM job_posts
WHERE id = ?
`

func (q *Queries) DeleteJobPost(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteJobPost, id)
	return err
}

const getJobPost = `-- name: GetJobPost :one
SELECT id, title, content, created_at, status, employer_id FROM job_posts
WHERE id = ? LIMIT 1
`

func (q *Queries) GetJobPost(ctx context.Context, id int64) (JobPost, error) {
	row := q.db.QueryRowContext(ctx, getJobPost, id)
	var i JobPost
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Content,
		&i.CreatedAt,
		&i.Status,
		&i.EmployerID,
	)
	return i, err
}

const listJobPosts = `-- name: ListJobPosts :many
SELECT id, title, content, created_at, status, employer_id FROM job_posts
ORDER BY created_at DESC
`

func (q *Queries) ListJobPosts(ctx context.Context) ([]JobPost, error) {
	rows, err := q.db.QueryContext(ctx, listJobPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []JobPost
	for rows.Next() {
		var i JobPost
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
			&i.Status,
			&i.EmployerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listJobPostsByAccount = `-- name: ListJobPostsByAccount :many
SELECT id, title, content, created_at, status, employer_id FROM job_posts
WHERE employer_id = ?
ORDER BY created_at DESC
`

func (q *Queries) ListJobPostsByAccount(ctx context.Context, employerID int64) ([]JobPost, error) {
	rows, err := q.db.QueryContext(ctx, listJobPostsByAccount, employerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []JobPost
	for rows.Next() {
		var i JobPost
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
			&i.Status,
			&i.EmployerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateJobPost = `-- name: UpdateJobPost :exec
UPDATE job_posts
set title = ?, content = ?, status = ?
WHERE id = ?
`

type UpdateJobPostParams struct {
	Title   string
	Content string
	Status  string
	ID      int64
}

func (q *Queries) UpdateJobPost(ctx context.Context, arg UpdateJobPostParams) error {
	_, err := q.db.ExecContext(ctx, updateJobPost,
		arg.Title,
		arg.Content,
		arg.Status,
		arg.ID,
	)
	return err
}
