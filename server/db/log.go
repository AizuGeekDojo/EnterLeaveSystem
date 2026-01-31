package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// AddLog adds enter/leave log and updates user status in a transaction
func AddLog(ainsID string, isEnter bool, Timestamp time.Time, Ext string, db *sql.DB) error {
	// Start transaction for atomic operation
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	isEnterInt := 0
	if isEnter {
		isEnterInt = 1
	}
	tsMillis := Timestamp.UnixNano() / int64(time.Millisecond)

	// Insert log entry
	_, err = tx.Exec(
		`INSERT INTO log (ainsID, isenter, time, ext) VALUES (?, ?, ?, ?)`,
		ainsID, isEnterInt, tsMillis, Ext,
	)
	if err != nil {
		return fmt.Errorf("failed to insert log: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
