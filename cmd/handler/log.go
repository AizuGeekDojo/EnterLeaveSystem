package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AizuGeekDojo/EnterLeaveSystem/cmd/db"
	"github.com/AizuGeekDojo/EnterLeaveSystem/cmd/slack"
)

// LogInfo is log data structue
type LogInfo struct {
	UID     string `json:"SID"`
	IsEnter int    `json:"IsEnter"`
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
	err := json.Unmarshal(body, &logdat)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Failed to parse JSON: %v", err)
		return
	}

	ts := time.Now()
	name, _, err := db.GetUserInfo(logdat.UID)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error: %v", err)
		return
	}
	if name == "" {
		w.WriteHeader(404)
		fmt.Fprintf(w, "The ID is not found.")
		return
	}

	err = db.AddLog(logdat.UID, (logdat.IsEnter == 1), ts, logdat.Ext)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error: %v", err)
		return
	}

	err = slack.Notify(name, logdat.UID, (logdat.IsEnter == 1), ts, logdat.Ext)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal server error: %v", err)
		return
	}
}
