package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/AizuGeekDojo/EnterLeaveSystem/server/db"
	"io"
	"log"
	"net/http"
	"strconv"
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
	case "POST":
		borrowReturnItemHandler(w, r, h.DB)
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
	ID int `json:"ID"`
	BarCode string `json:"BarCode"`
	Name string `json:"Name"`
}

type UserBorrowingTransaction struct {
	UID string `json:"SID"`
	Operations []UserBorrowingOperation `json:"Operations"`
}

type UserBorrowingOperation struct {
	ItemID string `json:"ID"`
	Operation string `json:"Op"`
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

func borrowReturnItemHandler(w http.ResponseWriter, r *http.Request, d *sql.DB) {
	var usertrans UserBorrowingTransaction
	r.ParseForm()

	reqlen, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Cannot get Content-Length: %v", err)
		log.Printf("%v %v: Bad request: %v", r.Method, r.URL.Path, err)
		return
	}
	body := make([]byte, reqlen)
	n, err := r.Body.Read(body)
	if err != nil {
		if err != io.EOF || n == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Failed to read: %v", err)
			log.Printf("%v %v: Bad request: %v", r.Method, r.URL.Path, err)
			return
		}
	}
	err = json.Unmarshal(body, &usertrans)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to parse JSON: %v", err)
		log.Printf("%v %v: Bad request: %v", r.Method, r.URL.Path, err)
		return
	}

	for _, operation := range usertrans.Operations {
		// TODO: アイテムの貸し出し処理の実装
	}

	retbyte, err := json.Marshal(usertrans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v %v: json.Marshal error: %v", r.Method, r.URL.Path, err)
		fmt.Fprintf(w, "{}")
		return
	}
	w.Write(retbyte)
}