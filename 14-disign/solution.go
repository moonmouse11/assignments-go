//go:build solution

// Эталонное решение главы 14. Сравни со своей реализацией в disign.go.
package disign

import (
	"crypto/ed25519"
	crand "crypto/rand"
)

// KeyPair — пара ключей ed25519.
type KeyPair struct {
	Public  ed25519.PublicKey
	Private ed25519.PrivateKey
}

func GenerateKey() (KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(crand.Reader)
	if err != nil {
		return KeyPair{}, err
	}
	return KeyPair{Public: pub, Private: priv}, nil
}

func SignFile(kp KeyPair, data []byte) []byte {
	if len(kp.Private) != ed25519.PrivateKeySize {
		return nil // некорректный приватный ключ
	}
	return ed25519.Sign(ed25519.PrivateKey(kp.Private), data)
}

func VerifyFile(pub ed25519.PublicKey, data, sig []byte) bool {
	// ed25519.Verify panics for incorrect key/signature lengths — we return false.
	if len(pub) != ed25519.PublicKeySize || len(sig) != ed25519.SignatureSize {
		return false
	}
	return ed25519.Verify(ed25519.PublicKey(pub), data, sig)
}
