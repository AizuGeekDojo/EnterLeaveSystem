// Package nfc drives a Sony RC-S380 (PaSoRi) reader directly over USB and
// reads FeliCa student / faculty ID cards.
//
// The chip protocol follows nfcpy's port100 driver.
package nfc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

const (
	vendorSony    = 0x054c
	productRCS380 = 0x06c3

	usbTimeout = 500 * time.Millisecond
)

var (
	// ErrDeviceNotFound is returned when no RC-S380 is connected.
	ErrDeviceNotFound = errors.New("nfc: RC-S380 not found")
	// errNoTarget means no card answered within the RF timeout.
	errNoTarget = errors.New("nfc: no card")
	// errUSBTimeout means the reader itself did not answer.
	errUSBTimeout = errors.New("nfc: usb timeout")

	ackFrame = []byte{0x00, 0x00, 0xff, 0x00, 0xff, 0x00}

	// InSetProtocol defaults from nfcpy port100.Chipset.in_set_protocol_defaults
	inSetProtocolDefaults = []byte{
		0x00, 0x18, 0x01, 0x01, 0x02, 0x01, 0x03, 0x00, 0x04, 0x00, 0x05, 0x00,
		0x06, 0x00, 0x07, 0x08, 0x08, 0x00, 0x09, 0x00, 0x0a, 0x00, 0x0b, 0x00,
		0x0c, 0x00, 0x0e, 0x04, 0x0f, 0x00, 0x10, 0x00, 0x11, 0x00, 0x12, 0x00,
		0x13, 0x06,
	}
)

// RFError is an RF communication failure reported by the reader chip
// (e.g. the card was moved away mid-transfer). It is transient.
type RFError struct {
	Status uint32
	Detail string
}

func (e *RFError) Error() string {
	if e.Detail != "" {
		return "nfc: rf communication error: " + e.Detail
	}
	return fmt.Sprintf("nfc: rf communication error (status 0x%08x)", e.Status)
}

// Device is an opened RC-S380 reader.
type Device struct {
	usb *usbDevice
	buf []byte
}

// Open opens the first connected RC-S380 and configures it for 212kbps FeliCa.
func Open() (*Device, error) {
	usb, err := openUSB(vendorSony, productRCS380)
	if err != nil {
		return nil, err
	}
	d := &Device{usb: usb, buf: make([]byte, 512)}
	if err := d.init(); err != nil {
		usb.close()
		return nil, fmt.Errorf("nfc: init reader: %w", err)
	}
	return d, nil
}

// Close turns the RF field off and releases the device.
func (d *Device) Close() error {
	d.command(0x06, []byte{0x00}, usbTimeout) // SwitchRF off
	return d.usb.close()
}

func (d *Device) init() error {
	// Cancel whatever command a previous process may have left behind.
	if err := d.usb.write(ackFrame, usbTimeout); err != nil {
		return err
	}
	d.usb.read(d.buf, 50*time.Millisecond)

	steps := []struct {
		code byte
		data []byte
	}{
		{0x2a, []byte{0x01}},                   // SetCommandType 1
		{0x06, []byte{0x00}},                   // SwitchRF off
		{0x00, []byte{0x01, 0x01, 0x0f, 0x01}}, // InSetRF 212F
		{0x02, inSetProtocolDefaults},          // InSetProtocol defaults
		{0x02, []byte{0x00, 0x18}},             // InSetProtocol initial guard time
	}
	for _, s := range steps {
		res, err := d.command(s.code, s.data, usbTimeout)
		if err != nil {
			return fmt.Errorf("command 0x%02x: %w", s.code, err)
		}
		if len(res) > 0 && res[0] != 0 {
			return fmt.Errorf("command 0x%02x: status 0x%02x", s.code, res[0])
		}
	}
	return nil
}

// inCommRF sends data to the card and returns its response.
func (d *Device) inCommRF(data []byte, timeout time.Duration) ([]byte, error) {
	units := timeout/(100*time.Microsecond) + 10
	if units > 0xffff {
		units = 0xffff
	}
	req := binary.LittleEndian.AppendUint16(nil, uint16(units))
	req = append(req, data...)

	res, err := d.command(0x04, req, timeout+usbTimeout)
	if err != nil {
		return nil, err
	}
	if len(res) < 4 {
		return nil, fmt.Errorf("nfc: short InCommRF response % x", res)
	}
	// Error responses carry only the 4-byte status.
	if status := binary.LittleEndian.Uint32(res[0:4]); status != 0 {
		if status == 0x80 {
			return nil, errNoTarget
		}
		return nil, &RFError{Status: status}
	}
	if len(res) < 5 {
		return nil, fmt.Errorf("nfc: short InCommRF response % x", res)
	}
	return res[5:], nil
}

// command sends a chip command and returns the response payload.
func (d *Device) command(code byte, data []byte, timeout time.Duration) ([]byte, error) {
	if err := d.usb.write(buildFrame(code, data), usbTimeout); err != nil {
		return nil, err
	}

	n, err := d.usb.read(d.buf, usbTimeout)
	if err != nil {
		return nil, fmt.Errorf("read ack: %w", err)
	}
	frame := d.buf[:n]
	if !bytes.HasPrefix(frame, ackFrame) {
		return nil, fmt.Errorf("nfc: expected ack, got % x", frame)
	}
	if frame = frame[len(ackFrame):]; len(frame) == 0 {
		if n, err = d.usb.read(d.buf, timeout); err != nil {
			if errors.Is(err, errUSBTimeout) {
				// Abort the pending command so the chip accepts the next one.
				d.usb.write(ackFrame, usbTimeout)
			}
			return nil, fmt.Errorf("read response: %w", err)
		}
		frame = d.buf[:n]
	}
	return parseFrame(frame, code)
}

// buildFrame wraps a command in an extended information frame:
// 00 00 FF FF FF LEN(2, LE) LCS D6 CMD DATA... DCS 00
func buildFrame(code byte, data []byte) []byte {
	body := append([]byte{0xd6, code}, data...)
	frame := []byte{0x00, 0x00, 0xff, 0xff, 0xff}
	frame = binary.LittleEndian.AppendUint16(frame, uint16(len(body)))
	frame = append(frame, checksum(frame[5:7]))
	frame = append(frame, body...)
	frame = append(frame, checksum(body), 0x00)
	return frame
}

// parseFrame validates a response frame for code and returns its payload.
func parseFrame(frame []byte, code byte) ([]byte, error) {
	if len(frame) < 10 || !bytes.Equal(frame[0:5], []byte{0x00, 0x00, 0xff, 0xff, 0xff}) {
		return nil, fmt.Errorf("nfc: invalid frame % x", frame)
	}
	n := int(binary.LittleEndian.Uint16(frame[5:7]))
	if checksum(frame[5:7]) != frame[7] {
		return nil, fmt.Errorf("nfc: bad length checksum % x", frame)
	}
	if len(frame) < 8+n+1 {
		return nil, fmt.Errorf("nfc: truncated frame % x", frame)
	}
	body := frame[8 : 8+n]
	if checksum(body) != frame[8+n] {
		return nil, fmt.Errorf("nfc: bad data checksum % x", frame)
	}
	if n < 2 || body[0] != 0xd7 || body[1] != code+1 {
		return nil, fmt.Errorf("nfc: unexpected response % x to command 0x%02x", body, code)
	}
	return body[2:], nil
}

func checksum(b []byte) byte {
	var sum byte
	for _, v := range b {
		sum += v
	}
	return -sum
}
