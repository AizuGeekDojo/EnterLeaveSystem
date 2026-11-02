package nfc

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	// serviceStudentID holds the student / faculty number in block 0.
	serviceStudentID = 0x300b

	pollInterval = 200 * time.Millisecond
	// removeMisses is how many empty polls in a row mean the card was taken away.
	removeMisses = 2
)

// Card is a card read from the reader.
type Card struct {
	// Type is "student", "faculty" or "general".
	Type string
	// ID is the AINS ID for students, the card number for faculty
	// and the IDm (hex) for general cards.
	ID string
}

// Run polls for cards and calls onCard once per card touch until ctx is
// cancelled or the reader fails. A card left on the reader is reported once.
// Transient RF errors are retried silently; only device failures are returned.
func (d *Device) Run(ctx context.Context, onCard func(Card)) error {
	var current []byte
	misses := 0

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}

		idm, err := d.poll()
		if errors.Is(err, errNoTarget) {
			if misses++; misses >= removeMisses {
				current = nil
			}
			continue
		}
		if isTransient(err) {
			log.Printf("nfc: polling: %v", err)
			continue
		}
		if err != nil {
			return err
		}

		misses = 0
		if bytes.Equal(idm, current) {
			continue
		}

		card, err := d.readCard(idm)
		if isTransient(err) {
			// Leave current unset so the next poll retries the read.
			log.Printf("nfc: reading card %x: %v", idm, err)
			continue
		}
		if err != nil {
			return err
		}
		current = idm
		if card.Type == "" {
			log.Printf("nfc: unsupported card data from %x", idm)
			continue
		}
		onCard(card)
	}
}

func (d *Device) readCard(idm []byte) (Card, error) {
	block, err := d.readBlock(idm, serviceStudentID, 0)
	var status *StatusError
	if errors.As(err, &status) {
		// The card answered but has no student ID service: a general card.
		return Card{Type: "general", ID: hex.EncodeToString(idm)}, nil
	}
	if err != nil {
		return Card{}, err
	}
	return classifyBlock(block, idm), nil
}

// classifyBlock interprets block 0 of the student ID service.
// A zero Card means the data is not a known format.
func classifyBlock(block, idm []byte) Card {
	if !utf8.Valid(block) {
		return Card{Type: "general", ID: hex.EncodeToString(idm)}
	}
	data := string(block)
	if strings.HasPrefix(data, "1000") {
		sid := data[4:11]
		switch sid[0] {
		case '1':
			return Card{Type: "student", ID: "s" + sid}
		case '5':
			return Card{Type: "student", ID: "m" + sid}
		case '8':
			return Card{Type: "student", ID: "d" + sid}
		}
		return Card{}
	}
	// Faculty card numbers are stored as the raw block text up to the first
	// space or newline, as the previous Python reader produced them.
	if i := strings.IndexAny(data, " \n"); i >= 0 {
		data = data[:i]
	}
	return Card{Type: "faculty", ID: data}
}

// RunMock reads "<type> <id>" lines (e.g. "student s1300071") from r as card
// touches, for development without a reader.
func RunMock(ctx context.Context, r io.Reader, onCard func(Card)) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		fields := strings.Fields(sc.Text())
		if len(fields) != 2 {
			log.Printf("nfc: mock: expected \"<type> <id>\", got %q", sc.Text())
			continue
		}
		onCard(Card{Type: fields[0], ID: fields[1]})
	}
	if err := sc.Err(); err != nil {
		return err
	}
	<-ctx.Done()
	return ctx.Err()
}
