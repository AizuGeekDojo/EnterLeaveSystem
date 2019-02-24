package handler

import (
	"fmt"
	"net/http"
)

//AccountAPIHandler handles http request for account
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
	// TODO:
	// AddLog
	// Slack notify
}
