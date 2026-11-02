package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/AizuGeekDojo/EnterLeaveSystem/server/db"
	"github.com/AizuGeekDojo/EnterLeaveSystem/server/nfc"
	"github.com/AizuGeekDojo/EnterLeaveSystem/server/utils"
	"golang.org/x/net/websocket"
)

const (
	// ReaderErrorRetryDelay is the delay before retrying after NFC reader error
	ReaderErrorRetryDelay = 5 * time.Second
	// ReaderReopenDelay is the delay before reopening a reader that stopped working
	ReaderReopenDelay = 1 * time.Second
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

// ReadCard drives the card reader and sends read cards to clients until ctx is cancelled.
func ReadCard(ctx context.Context, d *sql.DB) {
	onCard := func(card nfc.Card) { sendCard(card, d) }

	if utils.GetReaderInfo().Mock {
		log.Println("socket: using mock card reader on stdin")
		if err := nfc.RunMock(ctx, os.Stdin, onCard); err != nil && ctx.Err() == nil {
			log.Printf("socket: mock reader error: %v", err)
		}
		return
	}

	readerErrorActive := false
	for ctx.Err() == nil {
		dev, err := nfc.Open()
		if err != nil {
			log.Printf("socket: nfc reader error : %v\n", err)

			// Notify Slack only on the first occurrence of a reader error.
			if !readerErrorActive {
				readerErrorActive = true
				if notifyErr := utils.SlackNotifyReaderError(err.Error(), false); notifyErr != nil {
					log.Printf("socket: failed to send Slack reader error notification: %v", notifyErr)
				}
			}
			sendReaderError(err)

			select {
			case <-ctx.Done():
			case <-time.After(ReaderErrorRetryDelay):
			}
			continue
		}

		// Reader recovered from an error state – notify Slack.
		if readerErrorActive {
			readerErrorActive = false
			if notifyErr := utils.SlackNotifyReaderError("", true); notifyErr != nil {
				log.Printf("socket: failed to send Slack reader recovery notification: %v", notifyErr)
			}
		}

		log.Println("socket: nfc reader ready")
		err = dev.Run(ctx, onCard)
		dev.Close()
		if ctx.Err() != nil {
			return
		}
		// Reopen once before reporting: the reader may just have been reset.
		log.Printf("socket: nfc reader stopped, reopening: %v", err)
		time.Sleep(ReaderReopenDelay)
	}
}

// sendCard looks up the card owner and sends it to waiting clients.
func sendCard(card nfc.Card, d *sql.DB) {
	log.Printf("socket: card read: type=%s id=%s\n", card.Type, card.ID)

	var CardData IDCardInfo
	var err error
	CardData.IsCard = true
	CardData.CardID = card.ID

	if card.Type == "student" {
		CardData.AINSID = card.ID
		CardData.IsNew = false
	} else if card.Type == "faculty" || card.Type == "general" {
		CardData.AINSID, err = db.GetAinsIDByCardID(card.ID, d)
		if err != nil {
			log.Printf("socket: db.GetUserInfo error: %v", err)
			return
		}
		CardData.IsNew = (CardData.AINSID == "")
	} else {
		log.Printf("socket: unknown output: %v", card.Type)
		return
	}

	retbyte, err := json.Marshal(CardData)
	if err != nil {
		log.Printf("socket: json.Marshal error: %v", err)
		return
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

// sendReaderError tells waiting clients that the reader is unavailable.
func sendReaderError(readerErr error) {
	retbyte, err := json.Marshal(IDCardInfo{IsCard: false, ReaderErr: readerErr.Error()})
	if err != nil {
		log.Printf("socket: json.Marshal error: %v", err)
		return
	}

	clientsMutex.RLock()
	for _, c := range clients {
		if _, err := c.Write(retbyte); err != nil {
			log.Printf("socket: failed to write to client: %v", err)
		}
	}
	clientsMutex.RUnlock()
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
