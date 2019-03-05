package handler

import (
	"encoding/json"
	"fmt"
	"log"
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
		return
	}

	if r.Method == "POST" {
		addLogHandler(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Unexpected method")
		log.Printf("%v %v: Unexpected method", r.Method, r.URL.Path)
	}
}

func addLogHandler(w http.ResponseWriter, r *http.Request) {
	var logdat LogInfo
	reqlen, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Cannot get Content-Length: %v", err)
		log.Printf("%v %v: Bad request: %v", r.Method, r.URL.Path, err)
		return
	}
	body := make([]byte, reqlen)
	_, err = r.Body.Read(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to read: %v", err)
		log.Printf("%v %v: Bad request: %v", r.Method, r.URL.Path, err)
		return
	}
	err = json.Unmarshal(body, &logdat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to parse JSON: %v", err)
		log.Printf("%v %v: Bad request: %v", r.Method, r.URL.Path, err)
		return
	}

	ts := time.Now()
	name, _, err := db.GetUserInfo(logdat.UID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error: %v", err)
		log.Printf("%v %v: db.GetUserInfo error: %v", r.Method, r.URL.Path, err)
		return
	}
	if name == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "The ID is not found.")
		return
	}

	err = db.AddLog(logdat.UID, (logdat.IsEnter == 1), ts, logdat.Ext)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error: %v", err)
		log.Printf("%v %v: db.AddLog error: %v", r.Method, r.URL.Path, err)
		return
	}

	err = slack.Notify(name, logdat.UID, (logdat.IsEnter == 1), ts, logdat.Ext)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error: %v", err)
		log.Printf("%v %v: slack.Notify error: %v", r.Method, r.URL.Path, err)
		return
	}
}
