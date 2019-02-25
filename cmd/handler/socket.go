package handler

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/AizuGeekDojo/EnterLeaveSystem/cmd/db"
	"golang.org/x/net/websocket"
)

var clients = []*websocket.Conn{}

// ReadCardHandler handles Felica card reader.
func ReadCardHandler(ws *websocket.Conn) {
	clients = append(clients, ws)
	dat := []byte{}
	var err error
	for err == nil {
		_, err = ws.Read(dat)
	}
}

type IDCardInfo struct {
	IsCard bool   `json:"IsCard"`
	CardID string `json:"CardID"`
	SID    string `json:"SID"`
	IsNew  bool   `json:"IsNew"`
}

func sendData(dat IDCardInfo) {
	retbyte, _ := json.Marshal(dat)
	for _, c := range clients {
		c.Write(retbyte)
		c.Close()
	}
	clients = nil
}

func ReadCard() {
	for {
		dat, err := exec.Command("./test").Output()
		// dat, err := exec.Command("python2", "nfc_reader.py").Output()
		if err != nil {
			panic(err)
		}
		datstrspl := strings.Split(string(dat), " ")
		if len(datstrspl) < 2 {
			break
		}
		cardtype := datstrspl[0]
		cardid := strings.Split(datstrspl[1], "\n")[0]
		var resdat IDCardInfo
		resdat.IsCard = true
		resdat.CardID = cardid
		if cardtype == "student" {
			resdat.SID = cardid
			resdat.IsNew = false
		} else if cardtype == "univ" || cardtype == "general" {
			resdat.SID, _ = db.GetUIDByCardID(cardid)
			resdat.IsNew = (resdat.SID == "")
		} else {
			continue
		}
		sendData(resdat)
	}
}
