package handler

import (
	"fmt"
	"net/http"
	"strings"
)

//AccountAPIHandler handles http request for account
func LogAPIHandler(w http.ResponseWriter, r *http.Request) {
	//Cors Header
	w.Header().Add("Access-Control-Allow-Origin", strings.Join([]string{"http://", common.OriginDomain}, ""))

	//Cors Option check
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		w.WriteHeader(200)
		return
	}

	// //Routing
	// switch r.URL.Path {
	// case "/api/log":
	// 	// loginHandler(w, r)
	// default:
	// 	w.WriteHeader(404)
	// 	fmt.Fprintf(w, "Unknown API")
	// }

	if r.Method == "POST" {
		// downloadHandler(w, r)
	} else {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Unexpected method")
	}
}
