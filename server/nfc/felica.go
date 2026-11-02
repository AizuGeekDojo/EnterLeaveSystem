package nfc

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

const (
	pollTimeout = 20 * time.Millisecond
	readTimeout = 100 * time.Millisecond
)

// StatusError is a FeliCa status flag error returned by the card itself,
// e.g. when the requested service does not exist on the card.
type StatusError struct {
	Flag1, Flag2 byte
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("nfc: felica status %02x %02x", e.Flag1, e.Flag2)
}

// poll runs a FeliCa Polling (system code FFFF) and returns the card IDm.
func (d *Device) poll() ([]byte, error) {
	res, err := d.inCommRF([]byte{0x06, 0x00, 0xff, 0xff, 0x01, 0x00}, pollTimeout)
	if err != nil {
		return nil, err
	}
	// LEN 01 IDm(8) PMm(8) [SystemCode(2)]
	if len(res) < 18 || res[1] != 0x01 {
		return nil, &RFError{Detail: fmt.Sprintf("invalid polling response % x", res)}
	}
	return append([]byte(nil), res[2:10]...), nil
}

// readBlock runs Read Without Encryption for one block of one service.
func (d *Device) readBlock(idm []byte, service uint16, block byte) ([]byte, error) {
	cmd := []byte{0, 0x06}
	cmd = append(cmd, idm...)
	cmd = append(cmd, 0x01, byte(service), byte(service>>8), 0x01, 0x80, block)
	cmd[0] = byte(len(cmd))

	res, err := d.inCommRF(cmd, readTimeout)
	if err != nil {
		return nil, err
	}
	// LEN 07 IDm(8) Status1 Status2 [NumBlocks Data(16)...]
	if len(res) < 12 || res[1] != 0x07 || !bytes.Equal(res[2:10], idm) {
		return nil, &RFError{Detail: fmt.Sprintf("invalid read response % x", res)}
	}
	if res[10] != 0 {
		return nil, &StatusError{Flag1: res[10], Flag2: res[11]}
	}
	if len(res) < 13+16 || res[12] < 1 {
		return nil, &RFError{Detail: fmt.Sprintf("short read response % x", res)}
	}
	return res[13 : 13+16], nil
}

func isTransient(err error) bool {
	var rf *RFError
	return errors.Is(err, errNoTarget) || errors.As(err, &rf)
}
