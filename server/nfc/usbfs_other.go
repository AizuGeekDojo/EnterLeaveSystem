//go:build !linux

package nfc

import (
	"errors"
	"time"
)

// usbDevice is only implemented on Linux (usbfs).
type usbDevice struct{}

func openUSB(vid, pid uint16) (*usbDevice, error) {
	return nil, errors.New("nfc: USB access is only supported on Linux")
}

func (d *usbDevice) write(data []byte, timeout time.Duration) error { return nil }

func (d *usbDevice) read(buf []byte, timeout time.Duration) (int, error) { return 0, nil }

func (d *usbDevice) close() error { return nil }
