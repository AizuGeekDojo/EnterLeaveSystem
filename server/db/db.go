package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" //Sqlite database driver
)

const (
	// DBFilename is the default SQLite database file name
	DBFilename = "database.db"
	// MaxOpenConns is the maximum number of open connections to the database
	MaxOpenConns = 25
	// MaxIdleConns is the maximum number of idle connections in the pool
	MaxIdleConns = 5
	// ConnMaxLifetime is the maximum amount of time a connection may be reused
	ConnMaxLifetime = 5 * time.Minute
)

// OpenDB opens and initializes the SQLite database
func OpenDB() (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open("sqlite3", DBFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetConnMaxLifetime(ConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables if not exists
	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// createTables creates the necessary database tables if they don't exist
func createTables(db *sql.DB) error {
	tables := []struct {
		name   string
		schema string
	}{
		{
			name:   "users",
			schema: `CREATE TABLE IF NOT EXISTS "users" (ainsID TEXT PRIMARY KEY, name TEXT NOT NULL)`,
		},
		{
			name:   "idcard",
			schema: `CREATE TABLE IF NOT EXISTS "idcard" (idm TEXT PRIMARY KEY, ainsID TEXT NOT NULL)`,
		},
		{
			name:   "log",
			schema: `CREATE TABLE IF NOT EXISTS "log" (id INTEGER PRIMARY KEY AUTOINCREMENT, ainsID TEXT NOT NULL, isenter INTEGER NOT NULL, time INTEGER NOT NULL, ext TEXT)`,
		},
	}

	for _, table := range tables {
		if _, err := db.Exec(table.schema); err != nil {
			return fmt.Errorf("failed to create table %s: %w", table.name, err)
		}
	}

	// Create indexes for better query performance
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_users_ainsID ON users(ainsID)`,
		`CREATE INDEX IF NOT EXISTS idx_idcard_idm ON idcard(idm)`,
		`CREATE INDEX IF NOT EXISTS idx_log_ainsID ON log(ainsID)`,
		`CREATE INDEX IF NOT EXISTS idx_log_time ON log(time)`,
	}

	for _, idx := range indexes {
		if _, err := db.Exec(idx); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}
