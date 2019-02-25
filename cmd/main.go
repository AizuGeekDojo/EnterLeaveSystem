package main

import (
	"fmt"
	"net/http"

	"github.com/AizuGeekDojo/EnterLeaveSystem/cmd/handler"
	"golang.org/x/net/websocket"
)

func main() {
	fmt.Println("Starting server...")

	panic("TODO: Caution: DB structure is modifyed. Please reset Database.")
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("dist"))))

	//API handler
	http.Handle("/socket/readCard", websocket.Handler(handler.ReadCardHandler))
	http.HandleFunc("/api/user", handler.UserAPIHandler)
	http.HandleFunc("/api/log", handler.LogAPIHandler)

	//Standby NFC card reader
	go handler.ReadCard()

	//Start web server
	fmt.Println("Start server")
	http.ListenAndServe(":3000", nil)
}
