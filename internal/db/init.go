package db

import (
	"context"
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema/accounts.sql
var accountsTableDdl string

//go:embed schema/job_posts.sql
var jobPostsTableDdl string

func Initialize(dbPath string) (*Queries, error) {
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err := database.ExecContext(context.Background(), accountsTableDdl); err != nil {
		return nil, err
	}
	if _, err := database.ExecContext(context.Background(), jobPostsTableDdl); err != nil {
		return nil, err
	}
	queries := New(database)
	return queries, nil
}
