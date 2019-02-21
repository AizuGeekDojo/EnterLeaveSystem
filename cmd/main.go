package main

import (
	"fmt"
	"net/http"

	"github.com/AizuGeekDojo/EnterLeaveSystem/cmd/handler"
)

func main() {
	fmt.Println("Starting server...")

	//Static file handler
	// http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	//API handler
	http.HandleFunc("/socket/readcard", handler.SocketHandler)
	http.HandleFunc("/api/user", handler.UserAPIHandler)
	http.HandleFunc("/api/log", handler.LogAPIHandler)

	//Start web server
	fmt.Println("Start server")
	http.ListenAndServe(":3000", nil)
}
