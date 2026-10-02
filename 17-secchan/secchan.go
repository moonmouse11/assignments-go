//go:build !solution

// Пакет secchan — глава 17. Реализуй тип и функции так, чтобы прошли
// тесты в secchan_test.go. Условие — в README.md этой папки.
package secchan

import "crypto/ecdh"

// Формат запечатанного сообщения (sealed box):
//
//	[0:32]   эфемерный публичный ключ X25519 отправителя
//	[32:44]  nonce GCM (12 байт)
//	[44:]    AES-GCM(msg), ключ = HKDF-SHA-256(общий секрет ECDH)

// Party — участник канала: пара ключей X25519.
type Party struct {
	// TODO: приватный ключ *ecdh.PrivateKey.
	priv *ecdh.PrivateKey
}

// NewParty генерирует новую пару ключей X25519 из crypto/rand.
func NewParty() (Party, error) {
	panic("TODO: NewParty")
}

// Public возвращает публичный ключ X25519 стороны в байтах (32).
func (p Party) Public() []byte {
	panic("TODO: Public")
}

// Seal запечатывает msg для получателя с публичным ключом receiverPublic:
// эфемерная пара X25519, ECDH с ключом получателя, HKDF-SHA-256 → AES-GCM.
// Выход: ephemeralPub || nonce || ciphertext. Ошибки: неверная длина
// receiverPublic, сбой ГСЧ, недопустимый ключ.
func Seal(receiverPublic []byte, msg []byte) ([]byte, error) {
	panic("TODO: Seal")
}

// Open расшифровывает sealed у получателя receiver. Ошибки: короткое
// сообщение, чужой/недопустимый эфемерный ключ, порча ciphertext.
func Open(receiver Party, sealed []byte) ([]byte, error) {
	panic("TODO: Open")
}
