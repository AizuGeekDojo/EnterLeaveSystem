package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	cfg := GetSlackInfo()
	UPLOADURL := "https://slack.com/api/files.upload"
	TOKEN := cfg.CSVLogToken
	CHANNEL := cfg.CSVLogChannelID

	csv, err := csvExport(d)
	if err != nil {
		return err
	}
	resp, err := http.PostForm(
		UPLOADURL,
		url.Values{
			"token":           {TOKEN},
			"channels":        {CHANNEL},
			"filename":        {"log.csv"},
			"initial_comment": {"Enter leave log (csv format) by csvexport"},
			"title":           {"log.csv"},
			"content":         {csv},
		},
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var res = make(map[string]interface{})
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}
	ok, exists := res["ok"].(bool)
	if !exists {
		return errors.New("invalid Slack API response: 'ok' field is missing or not a boolean")
	}
	if !ok {
		errMsg, ok := res["error"].(string)
		if !ok {
			return errors.New("Slack API request failed with unknown error")
		}
		return errors.New(errMsg)
	}

	_, err = d.Exec(`delete from log`)
	if err != nil {
		return err
	}

	return nil
}
