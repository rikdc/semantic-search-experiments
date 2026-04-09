package db

import (
	"database/sql"
	"fmt"
)

// Execute runs a SQL statement that does not return rows, such as INSERT, UPDATE, or DELETE.
// It returns the result metadata including rows affected and last insert ID.
func Execute(db *sql.DB, query string, args ...any) (sql.Result, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing statement: %w", err)
	}
	return result, nil
}

// QueryRow runs a query expected to return at most one row.
// The caller must call Scan on the returned Row to retrieve the value.
func QueryRow(db *sql.DB, query string, args ...any) *sql.Row {
	return db.QueryRow(query, args...)
}

// FetchAll runs a query and returns all matching rows. The caller is responsible
// for closing the returned Rows when done.
func FetchAll(db *sql.DB, query string, args ...any) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("fetching rows: %w", err)
	}
	return rows, nil
}
