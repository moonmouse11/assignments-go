package secchan

import (
	"bytes"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	alice, err := NewParty()
	if err != nil {
		t.Fatalf("NewParty(alice): %v", err)
	}
	bob, err := NewParty()
	if err != nil {
		t.Fatalf("NewParty(bob): %v", err)
	}
	if len(alice.Public()) != 32 {
		t.Fatalf("Public: длина %d, хотел 32", len(alice.Public()))
	}
	if bytes.Equal(alice.Public(), bob.Public()) {
		t.Fatal("Public: ключи двух сторон совпали")
	}

	msg := []byte("запрос: сгенерируй пароль из 3 сегментов")
	sealed, err := Seal(bob.Public(), msg)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	got, err := Open(bob, sealed)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if !bytes.Equal(got, msg) {
		t.Errorf("round-trip: получил %q, хотел %q", got, msg)
	}
}

func TestWrongReceiver(t *testing.T) {
	bob, err := NewParty()
	if err != nil {
		t.Fatalf("NewParty(bob): %v", err)
	}
	mallory, err := NewParty()
	if err != nil {
		t.Fatalf("NewParty(mallory): %v", err)
	}
	sealed, err := Seal(bob.Public(), []byte("секрет"))
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	if _, err := Open(mallory, sealed); err == nil {
		t.Fatal("Open: хотел ошибку для чужого получателя, получил nil")
	}
}

func TestCorruptCiphertext(t *testing.T) {
	bob, err := NewParty()
	if err != nil {
		t.Fatalf("NewParty(bob): %v", err)
	}
	sealed, err := Seal(bob.Public(), []byte("данные"))
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	sealed[len(sealed)-1] ^= 0xFF // переворот последнего байта ciphertext
	if _, err := Open(bob, sealed); err == nil {
		t.Fatal("Open: хотел ошибку для испорченного ciphertext, получил nil")
	}
}

func TestSealBadPublicKey(t *testing.T) {
	if _, err := Seal([]byte("короткий"), []byte("msg")); err == nil {
		t.Fatal("Seal: хотел ошибку для ключа неверной длины, получил nil")
	}
}
