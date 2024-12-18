// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: employers.sql

package db

import (
	"context"
)

const createEmployer = `-- name: CreateEmployer :one
INSERT INTO employers (
  name
) VALUES (
  ?
)
RETURNING id, name
`

func (q *Queries) CreateEmployer(ctx context.Context, name string) (Employer, error) {
	row := q.db.QueryRowContext(ctx, createEmployer, name)
	var i Employer
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteEmployer = `-- name: DeleteEmployer :exec
DELETE FROM employers
WHERE id = ?
`

func (q *Queries) DeleteEmployer(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteEmployer, id)
	return err
}

const getEmployer = `-- name: GetEmployer :one
SELECT id, name FROM employers
WHERE id = ? LIMIT 1
`

func (q *Queries) GetEmployer(ctx context.Context, id int64) (Employer, error) {
	row := q.db.QueryRowContext(ctx, getEmployer, id)
	var i Employer
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listEmployers = `-- name: ListEmployers :many
SELECT id, name FROM employers
ORDER BY name
`

func (q *Queries) ListEmployers(ctx context.Context) ([]Employer, error) {
	rows, err := q.db.QueryContext(ctx, listEmployers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Employer
	for rows.Next() {
		var i Employer
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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

const numEmployers = `-- name: NumEmployers :one
SELECT COUNT(*) FROM employers
`

func (q *Queries) NumEmployers(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, numEmployers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const updateEmployer = `-- name: UpdateEmployer :exec
UPDATE employers
SET name = ?
WHERE id = ?
`

type UpdateEmployerParams struct {
	Name string
	ID   int64
}

func (q *Queries) UpdateEmployer(ctx context.Context, arg UpdateEmployerParams) error {
	_, err := q.db.ExecContext(ctx, updateEmployer, arg.Name, arg.ID)
	return err
}
