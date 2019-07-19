package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" //Sqlite database driver
)

const dbfilename string = "database.db"

func OpenDB() (*sql.DB, error) {
	//Open Database
	db, err := sql.Open("sqlite3", dbfilename)
	if err != nil {
		return nil, err
	}
	//Create tables if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "users" (sid TEXT,name TEXT,isenter INTEGER, FOREIGN KEY (isenter) REFERENCES io_status(id))`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "idcard" (idm TEXT,sid TEXT)`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "log" (sid TEXT,isenter INTEGER,time INTEGER,ext TEXT, FOREIGN KEY (isenter) REFERENCES io_status(id)))`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "io_status" (id TEXT, status TEXT)`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`INSERT INTO "io_status"(id, status) VALUES ('0', 'leave')`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`INSERT INTO "io_status"(id, status) VALUES ('1', 'enter')`)
	return db, err
}
