package db

import (
	"context"
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema/admin_accounts.sql
var adminAccountsTableDdl string

//go:embed schema/employers.sql
var employersTableDdl string

//go:embed schema/employer_accounts.sql
var employerAccountsTableDdl string

//go:embed schema/job_posts.sql
var jobPostsTableDdl string

func Initialize(dbPath string) (*Queries, error) {
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err := database.ExecContext(context.Background(), adminAccountsTableDdl); err != nil {
		return nil, err
	}
	if _, err := database.ExecContext(context.Background(), employersTableDdl); err != nil {
		return nil, err
	}
	if _, err := database.ExecContext(context.Background(), employerAccountsTableDdl); err != nil {
		return nil, err
	}
	if _, err := database.ExecContext(context.Background(), jobPostsTableDdl); err != nil {
		return nil, err
	}
	queries := New(database)
	return queries, nil
}
