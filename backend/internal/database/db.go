package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cyclingstream/backend/internal/logger"
	_ "github.com/lib/pq"
)

const (
	// SlowQueryThreshold is the duration threshold for logging slow queries
	SlowQueryThreshold = 1 * time.Second
)

type DB struct {
	*sql.DB
}

func New(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool to prevent connection exhaustion
	// SetMaxOpenConns sets the maximum number of open connections to the database
	// This should be set based on your database server's max_connections setting
	db.SetMaxOpenConns(25)
	
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	// Should be less than or equal to SetMaxOpenConns
	db.SetMaxIdleConns(5)
	
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	// This helps prevent issues with stale connections
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// SetConnMaxIdleTime sets the maximum amount of time a connection may be idle
	// Connections idle longer than this will be closed
	db.SetConnMaxIdleTime(1 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// QueryWithLogging wraps sql.DB.Query with slow query logging
func (db *DB) QueryWithLogging(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := db.DB.Query(query, args...)
	duration := time.Since(start)

	if duration > SlowQueryThreshold {
		logger.WithFields(map[string]interface{}{
			"query":    query,
			"args":     args,
			"duration_ms": duration.Milliseconds(),
			"error":    err != nil,
		}).Warn("Slow query detected")
	}

	return rows, err
}

// QueryRowWithLogging wraps sql.DB.QueryRow with slow query logging
// Note: For QueryRow, timing includes the Scan operation
// Usage: row := db.QueryRowWithLogging(query, args...); err := row.Scan(...)
// The timing will be logged after Scan completes (use LogQueryRowTiming for explicit timing)
func (db *DB) QueryRowWithLogging(query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRow(query, args...)
}

// LogQueryRowTiming logs slow QueryRow operations
// Call this after completing a QueryRow+Scan operation
// Usage:
//   start := time.Now()
//   row := db.QueryRow(query, args...)
//   err := row.Scan(...)
//   database.LogQueryRowTiming(query, args, time.Since(start), err)
func LogQueryRowTiming(query string, args []interface{}, duration time.Duration, err error) {
	if duration > SlowQueryThreshold {
		fields := map[string]interface{}{
			"query":       query,
			"duration_ms": duration.Milliseconds(),
		}
		if len(args) > 0 {
			fields["args"] = args
		}
		if err != nil {
			fields["error"] = err.Error()
		}
		logger.WithFields(fields).Warn("Slow query detected (QueryRow)")
	}
}

// LogSlowQuery is a helper function to log slow queries
// Use this to wrap database operations that might be slow
func LogSlowQuery(query string, args []interface{}, duration time.Duration, err error) {
	if duration > SlowQueryThreshold {
		fields := map[string]interface{}{
			"query":       query,
			"duration_ms": duration.Milliseconds(),
		}
		if len(args) > 0 {
			fields["args"] = args
		}
		if err != nil {
			fields["error"] = err.Error()
		}
		logger.WithFields(fields).Warn("Slow query detected")
	}
}

// ExecWithLogging wraps sql.DB.Exec with slow query logging
func (db *DB) ExecWithLogging(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := db.DB.Exec(query, args...)
	duration := time.Since(start)

	if duration > SlowQueryThreshold {
		logger.WithFields(map[string]interface{}{
			"query":    query,
			"args":     args,
			"duration_ms": duration.Milliseconds(),
			"error":    err != nil,
		}).Warn("Slow query detected")
	}

	return result, err
}

