package db

import "time"

// AddLog adds enter/leave log
// and change isenter status
func AddLog(UID string, isEnter bool, Timestamp time.Time, Ext string) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	isEnterInt := 0
	if isEnter {
		isEnterInt = 1
	}
	tsint64 := Timestamp.UnixNano() / int64(time.Millisecond)

	_, err = db.Exec(`insert into log values(?,?,?,?)`, UID, isEnterInt, tsint64, Ext)
	if err != nil {
		return err
	}
	_, err = db.Exec(`update users set isenter=? where sid=?`, isEnterInt, UID)
	if err != nil {
		return err
	}
	return nil
}
