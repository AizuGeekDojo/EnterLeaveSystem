package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// GetUserInfo returns username the person is in the room or not.
// If the person in the room, return true
func GetUserInfo(UID string, db *sql.DB) (string, bool, error) {
	// Check cardID is not registered
	row := db.QueryRow(`SELECT name,isenter FROM users WHERE sid=?`, UID)
	var isenter int
	var name string
	err := row.Scan(&name, &isenter)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return name, (isenter == 1), nil
}

// ForceLeave sets all users to leave status (called daily at midnight)
func ForceLeave(d *sql.DB) error {
	// Start transaction
	tx, err := d.Begin()
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

	// Query users who are currently in the room
	rows, err := tx.Query(`SELECT sid, name FROM users WHERE isenter=1`)
	if err != nil {
		return fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	ts := time.Now()
	tsMillis := ts.UnixNano() / int64(time.Millisecond)
	count := 0

	// Process each user
	for rows.Next() {
		var sid, name string
		if err := rows.Scan(&sid, &name); err != nil {
			return fmt.Errorf("failed to scan user row: %w", err)
		}

		// Insert leave log
		_, err := tx.Exec(
			`INSERT INTO log (sid, isenter, time, ext) VALUES (?, ?, ?, ?)`,
			sid, 0, tsMillis, "",
		)
		if err != nil {
			return fmt.Errorf("failed to insert log for user %s: %w", sid, err)
		}

		// Update user status to left
		_, err = tx.Exec(`UPDATE users SET isenter=? WHERE sid=?`, 0, sid)
		if err != nil {
			return fmt.Errorf("failed to update user %s: %w", sid, err)
		}

		log.Printf("Force left: %s (%s)", sid, name)
		count++
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Force leave completed: %d users processed", count)
	return nil
}
