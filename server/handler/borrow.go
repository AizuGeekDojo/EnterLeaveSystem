package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/AizuGeekDojo/EnterLeaveSystem/server/db"
	"log"
	"net/http"
)

//UserAPIHandler handles http request for user management
func (h *Handler) BorrowAPIHandler(w http.ResponseWriter, r *http.Request) {
	//Cors Header
	w.Header().Add("Access-Control-Allow-Origin", "*")

	//Cors Option check
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		return
	}

	switch r.Method {
	case "GET":
		getUserBorrowingHandler(w, r, h.DB)
	//case "POST":
	//	createNewBorrowing(w, r, h.DB)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Unexpected method")
		log.Printf("%v %v: Unexpected method", r.Method, r.URL.Path)
	}
}

// UserBorrowingInfo is user data structure for response
type UserBorrowingInfo struct {
	Products []UserBorrowingProduct `json:"Products"`
}

type UserBorrowingProduct struct {
	ID string `json:"ID"`
	BarCode string `json:"BarCode"`
	Name string `json:"Name"`
}

func getUserBorrowingHandler(w http.ResponseWriter, r *http.Request, d *sql.DB) {
	var userborrowingdat UserBorrowingInfo
	r.ParseForm()
	var uid = r.Form["sid"][0]

	username, _, err := db.GetUserInfo(uid, d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v %v: db.GetUserInfo error: %v", r.Method, r.URL.Path, err)
	} else if username == "" {
		w.WriteHeader(http.StatusNoContent)
	}

	borrowingdb, err := db.GetUserBorrowing(uid, d)
	for _, bdb := range borrowingdb {
		userborrowingdat.Products = append(userborrowingdat.Products, UserBorrowingProduct{
			ID:     bdb.Id,
			BarCode: bdb.Barcode,
			Name:    bdb.Name,
		})
	}

	retbyte, err := json.Marshal(borrowingdb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v %v: json.Marshal error: %v", r.Method, r.URL.Path, err)
		fmt.Fprintf(w, "{}")
		return
	}
	w.Write(retbyte)
}
