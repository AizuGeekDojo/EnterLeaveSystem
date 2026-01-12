package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/AizuGeekDojo/EnterLeaveSystem/server/db"
	"golang.org/x/net/websocket"
)

const (
	// ReaderErrorRetryDelay is the delay before retrying after NFC reader error
	ReaderErrorRetryDelay = 5 * time.Second
)

// IDCardInfo is structure for IDCard info
type IDCardInfo struct {
	IsCard bool   `json:"IsCard"`
	CardID string `json:"CardID"`
	AINSID string `json:"AINSID"`
	IsNew  bool   `json:"IsNew"`

	ReaderErr string `json:"ReaderErr"`
}

var (
	// clients is websocket connections
	clients = []*websocket.Conn{}
	// clientsMutex protects clients slice from concurrent access
	clientsMutex sync.RWMutex
)

// ReadCard runs card reader program, wait card data and send to clients.
func ReadCard(d *sql.DB) {
	for {
		var CardData IDCardInfo
		data, err := exec.Command("python3", "server/nfc/nfc_reader.py").Output()
		// data, err := exec.Command("python3", "server/test/nfc_reader.py").Output()
		if err != nil {
			log.Printf("socket: nfc reader error : %v\n", err)

			CardData.ReaderErr = err.Error()
			CardData.IsCard = false
			retbyte, err := json.Marshal(CardData)
			if err != nil {
				log.Printf("socket: json.Marshal error: %v", err)
				continue
			}

			clientsMutex.RLock()
			for _, c := range clients {
				if _, err := c.Write(retbyte); err != nil {
					log.Printf("socket: failed to write to client: %v", err)
				}
			}
			clientsMutex.RUnlock()
			time.Sleep(ReaderErrorRetryDelay)
			continue
		}

		datstrspl := strings.Split(string(data), " ")
		if len(datstrspl) < 2 {
			continue
		}

		cardtype := datstrspl[0]
		cardid := strings.Split(datstrspl[1], "\n")[0]

		log.Printf("socket: card read: type=%s id=%s\n", cardtype, cardid)

		CardData.IsCard = true
		CardData.CardID = cardid

		if cardtype == "student" {
			CardData.AINSID = cardid
			CardData.IsNew = false
		} else if cardtype == "faculty" || cardtype == "general" {
			CardData.AINSID, err = db.GetAinsIDByCardID(cardid, d)
			if err != nil {
				log.Printf("socket: db.GetUserInfo error: %v", err)
				continue
			}
			CardData.IsNew = (CardData.AINSID == "")
		} else {
			log.Printf("socket: unknown output: %v", cardtype)
			continue
		}

		retbyte, err := json.Marshal(CardData)

		if err != nil {
			log.Printf("socket: json.Marshal error: %v", err)
			continue
		}

		clientsMutex.Lock()
		for _, c := range clients {
			if _, err := c.Write(retbyte); err != nil {
				log.Printf("socket: failed to write to client: %v", err)
			}
			c.Close()
		}
		clients = nil
		clientsMutex.Unlock()
	}
}

// ReadCardHandler handles Felica card reader.
func (h *Handler) ReadCardHandler(ws *websocket.Conn) {
	clientsMutex.Lock()
	clients = append(clients, ws)
	clientsMutex.Unlock()

	// Keep connection alive until client disconnects
	dat := make([]byte, 1024)
	for {
		_, err := ws.Read(dat)
		if err != nil {
			// Client disconnected, remove from clients list
			clientsMutex.Lock()
			for i, c := range clients {
				if c == ws {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			clientsMutex.Unlock()
			break
		}
	}
}
