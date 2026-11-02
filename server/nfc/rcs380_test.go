package nfc

import (
	"bytes"
	"testing"
)

func TestBuildFrame(t *testing.T) {
	// SetCommandType(1) as sent by nfcpy
	want := []byte{0x00, 0x00, 0xff, 0xff, 0xff, 0x03, 0x00, 0xfd, 0xd6, 0x2a, 0x01, 0xff, 0x00}
	if got := buildFrame(0x2a, []byte{0x01}); !bytes.Equal(got, want) {
		t.Errorf("buildFrame = % x, want % x", got, want)
	}
}

func TestParseFrame(t *testing.T) {
	frame := []byte{0x00, 0x00, 0xff, 0xff, 0xff, 0x03, 0x00, 0xfd, 0xd7, 0x2b, 0x00, 0xfe, 0x00}
	got, err := parseFrame(frame, 0x2a)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, []byte{0x00}) {
		t.Errorf("payload = % x", got)
	}

	if _, err := parseFrame(frame, 0x06); err == nil {
		t.Error("expected error for response to another command")
	}
	bad := append([]byte(nil), frame...)
	bad[11] ^= 0xff
	if _, err := parseFrame(bad, 0x2a); err == nil {
		t.Error("expected checksum error")
	}
}

func TestClassifyBlock(t *testing.T) {
	idm := []byte{0x01, 0x2e, 0x4c, 0xd3, 0x9a, 0x11, 0x22, 0x33}
	tests := []struct {
		block string
		want  Card
	}{
		{"10001320137\x00\x00\x00\x00\x00", Card{"student", "s1320137"}},
		{"10005220001\x00\x00\x00\x00\x00", Card{"student", "m5220001"}},
		{"10008220001\x00\x00\x00\x00\x00", Card{"student", "d8220001"}},
		{"10009220001\x00\x00\x00\x00\x00", Card{}},
		{"1100131201400000", Card{"faculty", "1100131201400000"}},
		{"\xff\xfe00000000000000", Card{"general", "012e4cd39a112233"}},
	}
	for _, tt := range tests {
		if got := classifyBlock([]byte(tt.block), idm); got != tt.want {
			t.Errorf("classifyBlock(%q) = %+v, want %+v", tt.block, got, tt.want)
		}
	}
}
