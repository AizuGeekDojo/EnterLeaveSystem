package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/AizuGeekDojo/EnterLeaveSystem/server/db"
	"golang.org/x/net/websocket"
)

// IDCardInfo is structure for IDCard info
type IDCardInfo struct {
	IsCard bool   `json:"IsCard"`
	CardID string `json:"CardID"`
	SID    string `json:"SID"`
	IsNew  bool   `json:"IsNew"`

	ReaderErr string `json:"ReaderErr"`
}

// clients is websocket connections
var clients = []*websocket.Conn{}

// ReadCard runs card reader program, wait card data and send to clients.
func ReadCard(d *sql.DB) {
	for {
		var resdat IDCardInfo
		dat, err := exec.Command("python3", "nfc_reader.py").Output()
		if err != nil {
			log.Printf("socket: nfc reader error : %v\n", err)

			resdat.ReaderErr = err.Error()
			resdat.IsCard = false
			retbyte, err := json.Marshal(resdat)
			if err != nil {
				log.Printf("socket: json.Marshal error: %v", err)
				continue
			}

			for _, c := range clients {
				c.Write(retbyte)
			}
			time.Sleep(60 * time.Second)
			continue
		}

		datstrspl := strings.Split(string(dat), " ")
		if len(datstrspl) < 2 {
			continue
		}

		cardtype := datstrspl[0]
		cardid := strings.Split(datstrspl[1], "\n")[0]

		resdat.IsCard = true
		resdat.CardID = cardid

		if cardtype == "student" {
			resdat.SID = cardid
			resdat.IsNew = false
		} else if cardtype == "univ" || cardtype == "general" {
			resdat.SID, err = db.GetUIDByCardID(cardid, d)
			if err != nil {
				log.Printf("socket: db.GetUserInfo error: %v", err)
				continue
			}
			resdat.IsNew = (resdat.SID == "")
		} else {
			log.Printf("socket: unknown output: %v", cardtype)
			continue
		}

		retbyte, err := json.Marshal(resdat)

		if err != nil {
			log.Printf("socket: json.Marshal error: %v", err)
			continue
		}

		for _, c := range clients {
			c.Write(retbyte)
			c.Close()
		}
		clients = nil
	}
}

// ReadCardHandler handles Felica card reader.
func (h *Handler) ReadCardHandler(ws *websocket.Conn) {
	clients = append(clients, ws)
	dat := []byte{}
	var err error
	for err == nil {
		_, err = ws.Read(dat)
	}
}
