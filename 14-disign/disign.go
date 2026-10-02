//go:build !solution

// Пакет disign — глава 14. Реализуй функции так, чтобы прошли тесты
// в disign_test.go. Условие — в README.md этой папки.
package disign

import "crypto/ed25519"

// KeyPair — пара ключей ed25519: приватный (подпись) и публичный
// (проверка). Создаётся через GenerateKey.
type KeyPair struct {
	Public  ed25519.PublicKey
	Private ed25519.PrivateKey
}

// GenerateKey возвращает свежую пару ключей из криптографического
// источника crypto/rand. Ошибка чтения источника — наружу.
func GenerateKey() (KeyPair, error) {
	panic("TODO: GenerateKey")
}

// SignFile подписывает data приватным ключом kp и возвращает подпись
// (len == ed25519.SignatureSize). Некорректный приватный ключ — nil.
func SignFile(kp KeyPair, data []byte) []byte {
	panic("TODO: SignFile")
}

// VerifyFile проверяет подпись sig байтов data публичным ключом pub:
// true, только если sig — подпись именно этих data именно этим ключом.
// Кривые длины ключа или подписи — false, без паники.
func VerifyFile(pub ed25519.PublicKey, data, sig []byte) bool {
	panic("TODO: VerifyFile")
}
