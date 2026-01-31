package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AizuGeekDojo/EnterLeaveSystem/server/db"
)

// RegisterUserInfo is user register data structue for request
type RegisterUserInfo struct {
	AINSID string `json:"AINSID"`
	CardID string `json:"CardID"`
}

// RegisterUserResInfo is user register data structue for response
type RegisterUserResInfo struct {
	Success bool   `json:"Success"`
	Reason  string `json:"Reason"`
}

// // UserInfoReq is user data structue for request
// type UserInfoReq struct {
// 	UID string `json:"SID"`
// }

// UserInfo is user data structue for response
type UserInfo struct {
	AINSID   string `json:"AINSID"`
	UserName string `json:"UserName"`
	IsEnter  bool   `json:"IsEnter"`
}

//UserAPIHandler handles http request for user management
func (h *Handler) UserAPIHandler(w http.ResponseWriter, r *http.Request) {
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
		getUserHandler(w, r, h.DB)
	case "POST":
		createUserHandler(w, r, h.DB)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Unexpected method")
		log.Printf("%v %v: Unexpected method", r.Method, r.URL.Path)
	}
}

func getUserHandler(w http.ResponseWriter, r *http.Request, d *sql.DB) {
	var userresdat UserInfo
	r.ParseForm()

	if len(r.Form["AINSID"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing 'ainsID' parameter")
		log.Printf("%v %v: Missing 'ainsID' parameter", r.Method, r.URL.Path)
		return
	}

	var ainsID = r.Form["AINSID"][0]
	if err := validateAinsID(ainsID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid AINS ID: %v", err)
		log.Printf("%v %v: Invalid AINS ID: %v", r.Method, r.URL.Path, err)
		return
	}

	userresdat.AINSID = ainsID
	name, isEnter, err := db.GetUserEnterStatusByAinsID(ainsID, d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v %v: db.GetUserEnterStatusByAinsID error: %v", r.Method, r.URL.Path, err)
		fmt.Fprintf(w, "{}")
		return
	}
	userresdat.UserName = name
	userresdat.IsEnter = isEnter

	retbyte, err := json.Marshal(userresdat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v %v: json.Marshal error: %v", r.Method, r.URL.Path, err)
		fmt.Fprintf(w, "{}")
		return
	}
	w.Write(retbyte)
}

func createUserHandler(w http.ResponseWriter, r *http.Request, d *sql.DB) {
	var userdat RegisterUserInfo
	var userresdat RegisterUserResInfo

	if err := parseRequestBody(r, &userdat); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request: %v", err)
		log.Printf("%v %v: Bad request: %v", r.Method, r.URL.Path, err)
		return
	}

	if err := validateAinsID(userdat.AINSID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid AINS ID: %v", err)
		log.Printf("%v %v: Invalid AINS ID: %v", r.Method, r.URL.Path, err)
		return
	}

	err := db.RegisterCard(userdat.CardID, userdat.AINSID, d)
	if err == nil {
		userresdat.Success = true
		userresdat.Reason = ""
	} else {
		userresdat.Success = false
		userresdat.Reason = fmt.Sprintf("%v\n", err)
		log.Printf("%v %v: Bad request: %v", r.Method, r.URL.Path, err)
		w.WriteHeader(http.StatusBadRequest)
	}

	retbyte, err := json.Marshal(userresdat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v %v: json.Marshal error: %v", r.Method, r.URL.Path, err)
		fmt.Fprintf(w, "{}")
	}
	w.Write(retbyte)
}
