package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AizuGeekDojo/EnterLeaveSystem/cmd/db"
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

// UserInfoReq is user data structue for request
type UserInfoReq struct {
	UID string `json:"SID"`
}

// UserInfo is user data structue for response
type UserInfo struct {
	UID      string `json:"SID"`
	UserName string `json:"UserName"`
	IsEnter  bool   `json:"IsEnter"`
}

//UserAPIHandler handles http request for user management
func UserAPIHandler(w http.ResponseWriter, r *http.Request) {
	//Cors Header
	w.Header().Add("Access-Control-Allow-Origin", "*")

	//Cors Option check
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		w.WriteHeader(200)
		return
	}

	switch r.Method {
	case "GET":
		getUserHandler(w, r)
	case "POST":
		createUserHandler(w, r)
	default:
		w.WriteHeader(405)
		fmt.Fprintf(w, "Unexpected method")
	}
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	//学番を確認して、ユーザーが存在するかチェック
	//結果を返す
	var userreqdat UserInfoReq
	var userresdat UserInfo

	reqlen, _ := strconv.Atoi(r.Header.Get("Content-Length"))
	body := make([]byte, reqlen)
	r.Body.Read(body)
	json.Unmarshal(body, &userreqdat)

	// ts := time.Now()

	userresdat.UID = userreqdat.UID
	username, isenter, err := db.GetUserInfo(userreqdat.UID)
	if err != nil {
		w.WriteHeader(400)
	}
	userresdat.UserName = username
	userresdat.IsEnter = isenter

	retbyte, _ := json.Marshal(userresdat)
	w.Write(retbyte)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userdat RegistUserInfo
	var userresdat RegistUserResInfo

	reqlen, _ := strconv.Atoi(r.Header.Get("Content-Length"))
	body := make([]byte, reqlen)
	r.Body.Read(body)
	json.Unmarshal(body, &userdat)

	// ts := time.Now()

	err := db.RegisterCard(userdat.CardID, userdat.UID)
	if err == nil {
		userresdat.Success = true
		userresdat.Reason = ""
	} else {
		userresdat.Success = false
		userresdat.Reason = fmt.Sprintf("%#v\n", err)
		w.WriteHeader(400)
	}

	retbyte, _ := json.Marshal(userresdat)
	w.Write(retbyte)
}
