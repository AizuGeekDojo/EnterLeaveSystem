//go:build linux

package nfc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// usbdevfs_bulktransfer from <linux/usbdevice_fs.h>
type bulkTransfer struct {
	ep      uint32
	len     uint32
	timeout uint32 // milliseconds
	data    unsafe.Pointer
}

var (
	ioctlClaimInterface   = ioc(2, 'U', 15, 4) // _IOR('U', 15, unsigned int)
	ioctlReleaseInterface = ioc(2, 'U', 16, 4) // _IOR('U', 16, unsigned int)
	ioctlBulk             = ioc(3, 'U', 2, unsafe.Sizeof(bulkTransfer{}))
)

func ioc(dir, typ, nr, size uintptr) uintptr {
	return dir<<30 | size<<16 | typ<<8 | nr
}

// usbDevice is a USB device opened through Linux usbfs.
type usbDevice struct {
	f     *os.File
	iface uint32
	epIn  uint32
	epOut uint32
}

// openUSB finds the first device with the given VID/PID and claims interface 0.
func openUSB(vid, pid uint16) (*usbDevice, error) {
	sysdir, err := findSysfsDevice(vid, pid)
	if err != nil {
		return nil, err
	}
	busnum, err := readSysfsInt(filepath.Join(sysdir, "busnum"), 10)
	if err != nil {
		return nil, err
	}
	devnum, err := readSysfsInt(filepath.Join(sysdir, "devnum"), 10)
	if err != nil {
		return nil, err
	}
	epIn, epOut, err := findBulkEndpoints(sysdir)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(fmt.Sprintf("/dev/bus/usb/%03d/%03d", busnum, devnum), os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	d := &usbDevice{f: f, iface: 0, epIn: epIn, epOut: epOut}
	if err := d.ioctl(ioctlClaimInterface, unsafe.Pointer(&d.iface)); err != nil {
		f.Close()
		return nil, fmt.Errorf("claim interface: %w", err)
	}
	return d, nil
}

func findSysfsDevice(vid, pid uint16) (string, error) {
	dirs, err := filepath.Glob("/sys/bus/usb/devices/*")
	if err != nil {
		return "", err
	}
	for _, dir := range dirs {
		v, err := readSysfsInt(filepath.Join(dir, "idVendor"), 16)
		if err != nil || v != int(vid) {
			continue
		}
		p, err := readSysfsInt(filepath.Join(dir, "idProduct"), 16)
		if err != nil || p != int(pid) {
			continue
		}
		return dir, nil
	}
	return "", ErrDeviceNotFound
}

func findBulkEndpoints(sysdir string) (in, out uint32, err error) {
	eps, err := filepath.Glob(filepath.Join(sysdir, filepath.Base(sysdir)+":1.0", "ep_*"))
	if err != nil {
		return 0, 0, err
	}
	for _, ep := range eps {
		typ, _ := os.ReadFile(filepath.Join(ep, "type"))
		if strings.TrimSpace(string(typ)) != "Bulk" {
			continue
		}
		addr, err := strconv.ParseUint(strings.TrimPrefix(filepath.Base(ep), "ep_"), 16, 8)
		if err != nil {
			continue
		}
		if addr&0x80 != 0 {
			in = uint32(addr)
		} else {
			out = uint32(addr)
		}
	}
	if in == 0 || out == 0 {
		return 0, 0, fmt.Errorf("bulk endpoints not found in %s", sysdir)
	}
	return in, out, nil
}

func readSysfsInt(path string, base int) (int, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseInt(strings.TrimSpace(string(b)), base, 32)
	return int(v), err
}

func (d *usbDevice) ioctl(req uintptr, arg unsafe.Pointer) error {
	_, err := d.ioctlN(req, arg)
	return err
}

func (d *usbDevice) ioctlN(req uintptr, arg unsafe.Pointer) (int, error) {
	r, _, errno := syscall.Syscall(syscall.SYS_IOCTL, d.f.Fd(), req, uintptr(arg))
	if errno != 0 {
		return 0, errno
	}
	return int(r), nil
}

func (d *usbDevice) bulk(ep uint32, buf []byte, timeout time.Duration) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	t := bulkTransfer{
		ep:      ep,
		len:     uint32(len(buf)),
		timeout: uint32(timeout / time.Millisecond),
		data:    unsafe.Pointer(&buf[0]),
	}
	n, err := d.ioctlN(ioctlBulk, unsafe.Pointer(&t))
	runtime.KeepAlive(buf)
	if errors.Is(err, syscall.ETIMEDOUT) {
		return 0, errUSBTimeout
	}
	return n, err
}

func (d *usbDevice) write(data []byte, timeout time.Duration) error {
	_, err := d.bulk(d.epOut, data, timeout)
	return err
}

func (d *usbDevice) read(buf []byte, timeout time.Duration) (int, error) {
	return d.bulk(d.epIn, buf, timeout)
}

func (d *usbDevice) close() error {
	d.ioctl(ioctlReleaseInterface, unsafe.Pointer(&d.iface))
	return d.f.Close()
}
