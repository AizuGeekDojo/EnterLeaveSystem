package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// GetUserEnterStatusByAinsID returns username the person is in the room or not.
// If the person in the room, return true
func GetUserEnterStatusByAinsID(ainsID string, db *sql.DB) (string, bool, error) {
	// Check cardID is not registered
	row := db.QueryRow(`
		SELECT
		  	u.name,
			(
				SELECT l.isEnter
				FROM log AS l
				WHERE l.ainsID = u.ainsID
				ORDER BY l.time DESC
				LIMIT 1
			) AS isEnter
		FROM users AS u
		WHERE u.ainsID = ?;

	`, ainsID)
	var isEnter *int
	var name *string
	err := row.Scan(&name, &isEnter)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	if name == nil {
		return "", (isEnter != nil && *isEnter == 1), nil
	}
	return *name, (isEnter != nil && *isEnter == 1), nil
}

// GetUserInfoByAinsID returns username by ainsID
func GetUserInfoByAinsID(ainsID string, db *sql.DB) (string, error) {
	// Check cardID is not registered
	row := db.QueryRow(`SELECT name FROM users WHERE ainsID=?`, ainsID)
	var name string
	err := row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return name, nil
}

// GetAinsIDByCardID is return ainsID by felica's IDm or ID code
// This is prepared for not student person
// If ainsID is not found, return nil
func GetAinsIDByCardID(cardID string, db *sql.DB) (string, error) {
	// Check cardID is not registered
	row := db.QueryRow(`SELECT ainsID FROM idcard WHERE idm=?`, cardID)
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

	timeMillis := time.Now().UnixNano() / int64(time.Millisecond)
	yesterdayStartTime := timeMillis - 24*60*60*1000
	count := 0

	// Select all users who are currently in the room
	rows, err := tx.Query(`
       SELECT ainsID
       FROM log l
       WHERE time = (
               SELECT MAX(time)
               FROM log
               WHERE ainsID = l.ainsID
                 AND time >= ?
       )
       AND isEnter = 1
    `, yesterdayStartTime)
	if err != nil {
		return fmt.Errorf("failed to query users in room: %w", err)
	}
	defer rows.Close()

	// Process each user
	for rows.Next() {
		var ainsID string
		if err := rows.Scan(&ainsID); err != nil {
			return fmt.Errorf("failed to scan user row: %w", err)
		}

		// Insert leave log
		_, err := tx.Exec(
			`INSERT INTO log (ainsID, isEnter, time, ext) VALUES (?, ?, ?, ?)`,
			ainsID, 0, timeMillis, "",
		)
		if err != nil {
			return fmt.Errorf("failed to insert log for user %s: %w", ainsID, err)
		}

		log.Printf("Force left: %s", ainsID)
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
