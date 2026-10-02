package ppm

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	payload := []byte("секретные данные паролей")
	if err := Save(payload, "hunter2", &buf); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load(&buf, "hunter2")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !bytes.Equal(got, payload) {
		t.Errorf("round-trip: получил %q, хотел %q", got, payload)
	}
}

func TestWrongPassword(t *testing.T) {
	var buf bytes.Buffer
	if err := Save([]byte("payload"), "правильный", &buf); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := Load(&buf, "неправильный"); err == nil {
		t.Fatal("Load: хотел ошибку для неверного пароля, получил nil")
	}
}

func TestCorruptPayload(t *testing.T) {
	var buf bytes.Buffer
	if err := Save([]byte("0123456789abcdef"), "pw", &buf); err != nil {
		t.Fatalf("Save: %v", err)
	}
	data := buf.Bytes()
	data[len(data)-1] ^= 0xFF // переворот последнего байта payload
	if _, err := Load(bytes.NewReader(data), "pw"); err == nil {
		t.Fatal("Load: хотел ошибку для испорченного payload, получил nil")
	}
}

func TestCorruptMagic(t *testing.T) {
	var buf bytes.Buffer
	if err := Save([]byte("data"), "pw", &buf); err != nil {
		t.Fatalf("Save: %v", err)
	}
	data := buf.Bytes()
	data[0] = 'X' // порча magic
	if _, err := Load(bytes.NewReader(data), "pw"); err == nil {
		t.Fatal("Load: хотел ошибку для испорченного magic, получил nil")
	}
}

func TestHeaderLayout(t *testing.T) {
	var buf bytes.Buffer
	if err := Save([]byte("hi"), "pw", &buf); err != nil {
		t.Fatalf("Save: %v", err)
	}
	data := buf.Bytes()
	if len(data) < 65 {
		t.Fatalf("файл короче 65 байт: %d", len(data))
	}
	if want := hex.EncodeToString([]byte("PPMF")); hex.EncodeToString(data[:4]) != want {
		t.Errorf("magic: получил % x, хотел % x", data[:4], []byte("PPMF"))
	}
	if data[4] != 1 {
		t.Errorf("версия: получил %d, хотел 1", data[4])
	}
}
