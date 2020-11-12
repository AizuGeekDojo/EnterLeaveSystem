package main

import (
	"fmt"
	"net/http"

	"github.com/AizuGeekDojo/EnterLeaveSystem/server/db"
	"github.com/AizuGeekDojo/EnterLeaveSystem/server/handler"
	"github.com/AizuGeekDojo/EnterLeaveSystem/server/utils"
	"golang.org/x/net/websocket"
)

func main() {
	d, err := db.OpenDB()
	if err != nil {
		panic(err)
	}
	defer d.Close()
	h := handler.NewHandler(d)

	fmt.Println("Starting server...")

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("dist"))))

	//API handler
	http.Handle("/socket/readCard", websocket.Handler(h.ReadCardHandler))
	http.HandleFunc("/api/user", h.UserAPIHandler)
	http.HandleFunc("/api/log", h.LogAPIHandler)
	http.HandleFunc("/api/borrow", h.BorrowAPIHandler)

	//Standby NFC card reader
	go handler.ReadCard(d)

	//Start cron
	err = utils.CronInit(d)
	if err != nil {
		panic(err)
	}

	//Start web server
	fmt.Println("Start server")
	http.ListenAndServe(":3000", nil)
}
