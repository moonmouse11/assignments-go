package disign

import (
	"bytes"
	"crypto/ed25519"
	"testing"
)

var doc = []byte("Список файлов проекта\nREADME.md\ngo.mod\n")

func TestRoundTrip(t *testing.T) {
	kp, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	if len(kp.Private) != ed25519.PrivateKeySize || len(kp.Public) != ed25519.PublicKeySize {
		t.Fatalf("GenerateKey: длины ключей %d/%d — хотел %d/%d",
			len(kp.Private), len(kp.Public), ed25519.PrivateKeySize, ed25519.PublicKeySize)
	}
	sig := SignFile(kp, doc)
	if len(sig) != ed25519.SignatureSize {
		t.Fatalf("SignFile: длина подписи %d — хотел %d", len(sig), ed25519.SignatureSize)
	}
	if !VerifyFile(kp.Public, doc, sig) {
		t.Error("round-trip: подпись родного документа не прошла проверку")
	}
}

func TestTamperedData(t *testing.T) {
	kp, _ := GenerateKey()
	sig := SignFile(kp, doc)
	broken := bytes.Clone(doc)
	broken[len(broken)-1] ^= 0x01
	if VerifyFile(kp.Public, broken, sig) {
		t.Error("испорченный документ прошёл проверку — подпись не защищает")
	}
}

func TestTamperedSignature(t *testing.T) {
	kp, _ := GenerateKey()
	sig := bytes.Clone(SignFile(kp, doc))
	sig[0] ^= 0xFF
	if VerifyFile(kp.Public, doc, sig) {
		t.Error("испорченная подпись прошла проверку")
	}
}

func TestWrongKey(t *testing.T) {
	kp, _ := GenerateKey()
	other, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey (чужой): %v", err)
	}
	sig := SignFile(kp, doc)
	if VerifyFile(other.Public, doc, sig) {
		t.Error("чужой публичный ключ принял подпись")
	}
}

func TestBadInputs(t *testing.T) {
	kp, _ := GenerateKey()
	// Некорректный приватный ключ — nil, без паники.
	if got := SignFile(KeyPair{}, doc); got != nil {
		t.Errorf("SignFile с пустым ключом: получил %d байт — хотел nil", len(got))
	}
	// Кривые публичный ключ / подпись — false, без паники.
	if VerifyFile(nil, doc, SignFile(kp, doc)) {
		t.Error("пустой публичный ключ принял подпись")
	}
	if VerifyFile(kp.Public, doc, nil) || VerifyFile(kp.Public, doc, []byte{1, 2, 3}) {
		t.Error("кривая длина подписи принята")
	}
	if VerifyFile(kp.Public, doc, SignFile(kp, []byte("другой документ"))) {
		t.Error("подпись другого документа принята")
	}
}
