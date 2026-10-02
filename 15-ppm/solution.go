//go:build solution

// Эталонное решение главы 15. Сравни со своей реализацией в ppm.go.
package ppm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

const (
	magic   = "PPMF"
	ver     = byte(1)
	headLen = len(magic) + 1 + 16 + 12 + sha256.Size // 65
)

// gcm выводит ключ из пароля (pbkdf2, 1000 итераций) и собирает AES-GCM.
func gcm(password string, salt []byte) (cipher.AEAD, error) {
	key, err := pbkdf2.Key(sha256.New, password, salt, 1000, 32)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func Save(payload []byte, password string, w io.Writer) error {
	head := make([]byte, headLen-sha256.Size) // magic + версия + salt + nonce
	copy(head, magic)
	head[len(magic)] = ver
	if _, err := rand.Read(head[len(magic)+1:]); err != nil {
		return err
	}
	salt, nonce := head[len(magic)+1:len(magic)+1+16], head[len(magic)+1+16:]
	aead, err := gcm(password, salt)
	if err != nil {
		return err
	}
	h := sha256.Sum256(head)
	all := append(append(head, h[:]...), aead.Seal(nil, nonce, payload, nil)...)
	_, err = w.Write(all)
	return err
}

func Load(r io.Reader, password string) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if len(data) < headLen || string(data[:len(magic)]) != magic || data[len(magic)] != ver {
		return nil, errors.New("плохой заголовок файла (длина, magic или версия)")
	}
	if h := sha256.Sum256(data[:headLen-sha256.Size]); string(h[:]) != string(data[headLen-sha256.Size:headLen]) {
		return nil, errors.New("порчен заголовок: sha256 не сходится")
	}
	salt, nonce := data[len(magic)+1:len(magic)+1+16], data[len(magic)+1+16:headLen-sha256.Size]
	aead, err := gcm(password, salt)
	if err != nil {
		return nil, err
	}
	payload, err := aead.Open(nil, nonce, data[headLen:], nil)
	if err != nil {
		return nil, fmt.Errorf("расшифровка не удалась (пароль или payload): %w", err)
	}
	return payload, nil
}
