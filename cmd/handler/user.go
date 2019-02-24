package handler

import (
	"fmt"
	"net/http"
)

//AccountAPIHandler handles http request for account
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
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	//GetUserByCardIDで登録されてないか確認、重複チェック
	//AddCard
}
