package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AizuGeekDojo/EnterLeaveSystem/cmd/db"
	slack "github.com/AizuGeekDojo/EnterLeaveSystem/cmd/el_slack"
)

// LogInfo is log data structue
type LogInfo struct {
	UID     string `json:"SID"`
	isEnter int    `json:"IsEnter"`
	Ext     string `json:"Ext"`
}

//LogAPIHandler handles http request for logging
func LogAPIHandler(w http.ResponseWriter, r *http.Request) {
	//Cors Header
	w.Header().Add("Access-Control-Allow-Origin", "*")

	//Cors Option check
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		w.WriteHeader(200)
		return
	}

	if r.Method == "POST" {
		addLogHandler(w, r)
	} else {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Unexpected method")
	}
}

func addLogHandler(w http.ResponseWriter, r *http.Request) {
	var logdat LogInfo
	reqlen, _ := strconv.Atoi(r.Header.Get("Content-Length"))
	body := make([]byte, reqlen)
	r.Body.Read(body)
	json.Unmarshal(body, &logdat)

	ts := time.Now()
	db.AddLog(logdat.UID, (logdat.isEnter == 1), ts, logdat.Ext)

	slack.SlackNotify(db.GetUserName(logdat.UID), logdat.UID, (logdat.isEnter == 1), ts, logdat.Ext)
}
