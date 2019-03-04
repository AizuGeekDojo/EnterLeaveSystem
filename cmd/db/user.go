package db

import (
	"database/sql"
	"errors"
)

// GetUserInfo returns username the person is in the room or not.
// If the person in the room, return true
func GetUserInfo(UID string) (string, bool, error) {
	db, err := openDB()
	if err != nil {
		return "", false, err
	}
	defer db.Close()

	// Check cardID is not registered
	row := db.QueryRow(`SELECT name,isenter FROM users WHERE sid=?`, UID)
	var isenter int
	var name string
	err = row.Scan(&name, &isenter)
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
func GetUIDByCardID(CardID string) (string, error) {
	db, err := openDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Check cardID is not registered
	row := db.QueryRow(`SELECT sid FROM idcard WHERE idm=?`, CardID)
	var sid string
	err = row.Scan(&sid)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return sid, nil
}

// RegisterCard regist cardid with UID
func RegisterCard(CardID string, UID string) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Check user is exist
	gotuid, _, err := GetUserInfo(UID)
	if err != nil {
		return err
	}
	if gotuid == "" {
		return errors.New("ID \"" + UID + "\" is not found.")
	}

	// Check cardID is not registered
	row := db.QueryRow(`SELECT sid FROM idcard WHERE idm=? AND sid=?`, CardID, UID)
	var sid string
	err = row.Scan(&sid)
	if err != sql.ErrNoRows {
		// No error (already registered) don't continue
		// Other error can't continue
		return err
	}

	// Register cardID into database
	_, err = db.Exec(`insert into idcard values(?,?)`, CardID, UID)
	return err
}
