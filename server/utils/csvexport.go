package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

const (
	// TimestampMillisecondDivisor converts Unix timestamp from milliseconds to seconds
	TimestampMillisecondDivisor = 1000
)

// csvExport exports log data as CSV text
func csvExport(d *sql.DB) (string, error) {
	rows, err := d.Query(`select time,log.sid,name,log.isenter,ext from log,users where log.sid=users.sid`)
	if err != nil {
		return "", err
	}

	csv := "Date,StudentID,Name,Enter/Leave,Purpose,Comment\n"
	for rows.Next() {
		var (
			ts      int64
			sid     string
			name    string
			isenter int64
			ext     string
		)
		if err := rows.Scan(&ts, &sid, &name, &isenter, &ext); err != nil {
			return "", err
		}

		datefmted := time.Unix(ts/TimestampMillisecondDivisor, 0).Format("2006-01-02 15:04:05")
		entstr := "Leave"
		if isenter == 1 {
			entstr = "Enter"
		}

		if ext != "" {
			var RawJSON = []byte(ext)
			var ExtList = make(map[string]interface{})

			err := json.Unmarshal(RawJSON, &ExtList)
			if err != nil {
				return "", err
			}

			useage, ok := ExtList["Use"].([]interface{})
			if !ok {
				return "", errors.New("invalid Ext format: 'Use' field is not an array")
			}
			mess, ok := ExtList["message"].(string)
			if !ok {
				return "", errors.New("invalid Ext format: 'message' field is not a string")
			}
			mess = strings.Replace(mess, "\"", "\"\"", -1)

			csv += fmt.Sprintf("%v,%v,%v,%v,%v,\"%v\"\n", datefmted, sid, name, entstr, useage, mess)
		} else {
			csv += fmt.Sprintf("%v,%v,%v,%v,,\n", datefmted, sid, name, entstr)
		}
	}
	err = rows.Close()
	if err != nil {
		return "", err
	}

	return csv, nil
}

// sendMonthlyLog sends csv log file via slack
func sendMonthlyLog(d *sql.DB) error {
	// Use initialized Slack client (socket mode capable) like slack.go
	client := GetSlackClient()
	if client == nil {
		return errors.New("Slack client is not initialized")
	}

	cfg := GetSlackInfo()
	csv, err := csvExport(d)
	if err != nil {
		return err
	}

	// Upload the CSV content using UploadFileV2 (files.upload v2)
	_, err = client.UploadFileV2(slack.UploadFileV2Parameters{
		Filename:       "log.csv",
		Title:          "log.csv",
		InitialComment: "Enter leave log (csv format) by csvexport",
		Content:        csv,
		FileSize:       len([]byte(csv)),
		Channel:        cfg.CSVLogChannelID,
	})
	if err != nil {
		return err
	}

	if _, err = d.Exec(`delete from log`); err != nil {
		return err
	}
	return nil
}
