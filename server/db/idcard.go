package db

import (
	"database/sql"
	"errors"
)

// RegisterCard regist cardid with UID
func RegisterCard(CardID string, AINSID string, db *sql.DB) error {
	// Check user is exist
	gotAINSID, err := GetUserInfoByAinsID(AINSID, db)
	if err != nil {
		return err
	}
	if gotAINSID == "" {
		return errors.New("Your AINS ID (" + AINSID + ") is not registered in the system\nPlease contact to administrator.")
	}

	// Check cardID is not registered
	row := db.QueryRow(`SELECT ainsID FROM idcard WHERE idm=? AND ainsID=?`, CardID, AINSID)
	var sid string
	err = row.Scan(&sid)
	if err != sql.ErrNoRows {
		// No error (already registered) don't continue
		// Other error can't continue
		return err
	}

	// Register cardID into database
	_, err = db.Exec(`INSERT into idcard values(?,?)`, CardID, AINSID)
	return err
}