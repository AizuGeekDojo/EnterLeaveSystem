package db

import "database/sql"

type ProductDB struct {
	Id       string
	Name string
	Barcode  string
	Borrower string
}

// GetUserBorrowing is return requester's borrowing products
func GetUserBorrowing(UID string, db *sql.DB) ([]ProductDB, error) {
	var products []ProductDB
	rows, err := db.Query(`SELECT id,name,barcode,borrowersid FROM products WHERE borrowersid=?`, UID)
	if err != nil {
		return products, err
	}

	for rows.Next() {
		var id string
		var name string
		var barcode string
		var borrowersid string
		if err := rows.Scan(&id, &name, &barcode, &borrowersid); err != nil {
			return products, err
		}
		products = append(products, ProductDB{
			Id: id,
			Name: name,
			Barcode:  barcode,
			Borrower: borrowersid,
		})
	}

	defer rows.Close()
	if err := rows.Err(); err != nil {
		return products, err
	}

	return products, err
}

func CreateLendingEquipment(name string, barcode string, db *sql.DB) error {
	// バーコードの重複などを調べるエラー対処
	_, err := db.Exec(`INSERT INTO products(name,barcode) values(name,barcode)`, name, barcode)
	if err != nil {
		return err
	}
	return nil
}

func UpdateLendingEquipmentByID(id string, name string, barcode string, db *sql.DB) error {
	_, err := db.Exec(`UPDATE products SET name=?,barcode=? WHERE id=?`, name, barcode, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateLendingEquipmentByBarCode(name string, barcode string, db *sql.DB) error {
	_, err := db.Exec(`UPDATE products SET name=? WHERE barcode=?`, name, barcode)
	if err != nil {
		return err
	}
	return nil
}

func Borrow(id string, sid string, db *sql.DB) error {
	_, err := db.Exec(`UPDATE products SET borrowersid=? WHERE id=?`, sid, id)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Borrowが使えるようにバーコードとIDでFindできる関数の追加