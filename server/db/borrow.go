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
