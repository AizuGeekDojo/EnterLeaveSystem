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

// GetUIDByCardID is return UID by felica's IDm or ID code
// This is prepared for not student person
// If UID is not found, return nil
func GetUIDByCardID(CardID string, db *sql.DB) (string, error) {
	// Check cardID is not registered
	row := db.QueryRow(`SELECT sid FROM idcard WHERE idm=?`, CardID)
	var sid string
	err := row.Scan(&sid)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return sid, nil
}

// RegisterCard registers a card ID with a user ID
func RegisterCard(CardID string, UID string, db *sql.DB) error {
	// Check user exists
	gotuid, _, err := GetUserInfo(UID, db)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}
	if gotuid == "" {
		return fmt.Errorf("user ID %q not found", UID)
	}

	// Check cardID is not already registered
	row := db.QueryRow(`SELECT sid FROM idcard WHERE idm=? AND sid=?`, CardID, UID)
	var sid string
	err = row.Scan(&sid)
	if err != sql.ErrNoRows {
		if err == nil {
			return fmt.Errorf("card %q is already registered to user %q", CardID, UID)
		}
		return fmt.Errorf("failed to check card registration: %w", err)
	}

	// Register cardID into database
	_, err = db.Exec(`INSERT INTO idcard (idm, sid) VALUES (?, ?)`, CardID, UID)
	if err != nil {
		return fmt.Errorf("failed to register card: %w", err)
	}
	return nil
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
