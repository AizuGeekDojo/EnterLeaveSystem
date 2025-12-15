package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AizuGeekDojo/EnterLeaveSystem/server/db"
)

// RegistUserInfo is user register data structue for request
type RegistUserInfo struct {
	UID    string `json:"SID"`
	CardID string `json:"CardID"`
}

// RegistUserResInfo is user register data structue for response
type RegistUserResInfo struct {
	Success bool   `json:"Success"`
	Reason  string `json:"Reason"`
}

// // UserInfoReq is user data structue for request
// type UserInfoReq struct {
// 	UID string `json:"SID"`
// }

// UserInfo is user data structue for response
type UserInfo struct {
	UID      string `json:"SID"`
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

	if len(r.Form["sid"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing 'sid' parameter")
		log.Printf("%v %v: Missing 'sid' parameter", r.Method, r.URL.Path)
		return
	}

	var uid = r.Form["sid"][0]

	if err := validateSID(uid); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid SID: %v", err)
		log.Printf("%v %v: Invalid SID: %v", r.Method, r.URL.Path, err)
		return
	}

	userresdat.UID = uid

	username, isenter, err := db.GetUserInfo(uid, d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v %v: db.GetUserInfo error: %v", r.Method, r.URL.Path, err)
		fmt.Fprintf(w, "{}")
		return
	} else if username == "" {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintf(w, "{}")
		return
	}
	userresdat.UserName = username
	userresdat.IsEnter = isenter

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
	var userdat RegistUserInfo
	var userresdat RegistUserResInfo

	if err := parseRequestBody(r, &userdat); err != nil {
		handleRequestError(w, r, http.StatusBadRequest, "Bad request", err)
		return
	}

	if err := validateSID(userdat.UID); err != nil {
		handleRequestError(w, r, http.StatusBadRequest, "Invalid SID", err)
		return
	}

	if err := validateCardID(userdat.CardID); err != nil {
		handleRequestError(w, r, http.StatusBadRequest, "Invalid CardID", err)
		return
	}

	err := db.RegisterCard(userdat.CardID, userdat.UID, d)
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
		return
	}
	w.Write(retbyte)
}
