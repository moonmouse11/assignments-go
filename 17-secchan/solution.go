//go:build solution

// Эталонное решение главы 17. Сравни со своей реализацией в secchan.go.
package secchan

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

const (
	pubLen, nonceLen = 32, 12                 // ключ X25519 и nonce GCM
	minSealed        = pubLen + nonceLen + 16 // + имитовставок GCM
)

// Party — участник канала: пара ключей X25519.
type Party struct{ priv *ecdh.PrivateKey }

func NewParty() (Party, error) {
	priv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return Party{}, err
	}
	return Party{priv: priv}, nil
}

func (p Party) Public() []byte { return p.priv.PublicKey().Bytes() }

// deriveAEAD: общий секрет ECDH → HKDF-SHA-256 (salt — оба публичных
// ключа в стабильном порядке) → AES-GCM.
func deriveAEAD(shared, pub1, pub2 []byte) (cipher.AEAD, error) {
	key, err := hkdf.Key(sha256.New, shared, append(append([]byte{}, pub1...), pub2...), "secchan/v1", 32)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func Seal(receiverPublic []byte, msg []byte) ([]byte, error) {
	if len(receiverPublic) != pubLen {
		return nil, fmt.Errorf("публичный ключ получателя: %d байт, хотел %d", len(receiverPublic), pubLen)
	}
	rkey, err := ecdh.X25519().NewPublicKey(receiverPublic)
	if err != nil {
		return nil, err
	}
	eph, err := ecdh.X25519().GenerateKey(rand.Reader) // эфемерная пара отправителя
	if err != nil {
		return nil, err
	}
	shared, err := eph.ECDH(rkey)
	if err != nil {
		return nil, err
	}
	aead, err := deriveAEAD(shared, eph.PublicKey().Bytes(), receiverPublic)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, nonceLen)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return append(append(eph.PublicKey().Bytes(), nonce...), aead.Seal(nil, nonce, msg, nil)...), nil
}

func Open(receiver Party, sealed []byte) ([]byte, error) {
	if len(sealed) < minSealed {
		return nil, fmt.Errorf("сообщение короче минимума: %d байт", len(sealed))
	}
	ephPub, err := ecdh.X25519().NewPublicKey(sealed[:pubLen])
	if err != nil {
		return nil, err
	}
	shared, err := receiver.priv.ECDH(ephPub)
	if err != nil {
		return nil, err
	}
	aead, err := deriveAEAD(shared, sealed[:pubLen], receiver.Public())
	if err != nil {
		return nil, err
	}
	msg, err := aead.Open(nil, sealed[pubLen:pubLen+nonceLen], sealed[pubLen+nonceLen:], nil)
	if err != nil {
		return nil, fmt.Errorf("расшифровка не удалась: %w", err)
	}
	return msg, nil
}
