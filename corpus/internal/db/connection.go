package db

import (
	"database/sql"
	"fmt"
)

// Connect opens a connection to the database using the provided DSN and verifies
// it is reachable before returning. The caller is responsible for calling Close.
func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("verifying database connection: %w", err)
	}
	return db, nil
}

// Close shuts down the database connection pool and releases all held resources.
func Close(db *sql.DB) error {
	return db.Close()
}

// Ping checks that the database connection is still alive. Use this for health checks.
func Ping(db *sql.DB) error {
	return db.Ping()
}
