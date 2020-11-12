package db

import (
	"database/sql"
	"errors"
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


type ProductDB struct {
	id string
	barcode string
	borrower string
}
// GetUserBorrowing is return requester's borrowing products
func GetUserBorrowing(UID string, db *sql.DB) ([]ProductDB, error) {
	var products []ProductDB
	rows, err := db.Query(`SELECT id,name,barcode FROM products WHERE borrowersid=?`, UID)
	if err != nil {
		return products, err
	}

	for rows.Next() {
		var id string
		var name string
		var barcode string
		if err := rows.Scan(&id, &name, &barcode); err != nil {
			return products, err
		}
		products = append(products, ProductDB{
			id:       id,
			barcode:  name,
			borrower: barcode,
		})
	}

	defer rows.Close()
	if err := rows.Err(); err != nil {
		return products, err
	}

	return products, err
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

// RegisterCard regist cardid with UID
func RegisterCard(CardID string, UID string, db *sql.DB) error {
	// Check user is exist
	gotuid, _, err := GetUserInfo(UID, db)
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

// ForceLeave sets all users leave status
func ForceLeave(d *sql.DB) error {
	rows, err := d.Query(`SELECT * FROM users where isenter=1`)
	if err != nil {
		return err
	}
	ts := time.Now()

	tx, err := d.Begin()
	for rows.Next() {
		var (
			sid     string
			name    string
			isenter int64
		)
		if err := rows.Scan(&sid, &name, &isenter); err != nil {
			return err
		}

		tsint64 := ts.UnixNano() / int64(time.Millisecond)

		_, err := tx.Exec(`insert into log values(?,?,?,?)`, sid, 0, tsint64, "")
		if err != nil {
			return err
		}
		_, err = tx.Exec(`update users set isenter=? where sid=?`, 0, sid)
		if err != nil {
			return err
		}
		log.Printf("Force left: %v(%v)", sid, name)
	}
	err = rows.Close()
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}