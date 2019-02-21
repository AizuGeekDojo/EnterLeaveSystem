package handler

import (
	"fmt"
	"net/http"
)

//AccountAPIHandler handles http request for account
func SocketHandler(w http.ResponseWriter, r *http.Request) {
	//Cors Header
	w.Header().Add("Access-Control-Allow-Origin", "*")

	//Cors Option check
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		w.WriteHeader(200)
		return
	}

	//Routing
	switch r.URL.Path {
	case "/socket/readcard":
		// loginHandler(w, r)
	default:
		w.WriteHeader(404)
		fmt.Fprintf(w, "Unknown API")
	}
}
