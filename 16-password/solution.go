//go:build solution

// Эталонное решение главы 16. Сравни со своей реализацией в password.go.
package password

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
)

const (
	magic   = "PSWD"
	headLen = len(magic) + 1 + 16 + 12 + sha256.Size // 65
)

// Entry — одна запись менеджера паролей: контекст, логин, пароль.
type Entry struct{ Context, Username, Password string }

// Vault — хранилище записей (ключ — Context). SaveTo/LoadFrom шифруют
// JSON-слепок хранилища файлом формата PSWD (заголовок + AES-GCM).
type Vault struct{ entries map[string]Entry }

func NewVault() *Vault { return &Vault{entries: map[string]Entry{}} }

func (v *Vault) Add(e Entry) error {
	if e.Context == "" {
		return errors.New("пустой контекст")
	}
	v.entries[e.Context] = e
	return nil
}

func (v *Vault) Get(context string) (Entry, bool) {
	e, ok := v.entries[context]
	return e, ok
}

func (v *Vault) List() []Entry {
	out := make([]Entry, 0, len(v.entries))
	for _, e := range v.entries {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Context < out[j].Context })
	return out
}

// gcm выводит ключ из мастер-пароля (pbkdf2, 1000 итераций) и собирает AES-GCM.
func gcm(master string, salt []byte) (cipher.AEAD, error) {
	key, err := pbkdf2.Key(sha256.New, master, salt, 1000, 32)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func (v *Vault) SaveTo(path, master string) error {
	blob, err := json.Marshal(v.entries)
	if err != nil {
		return err
	}
	head := make([]byte, headLen-sha256.Size) // magic + версия + salt + nonce
	copy(head, magic)
	head[len(magic)] = 1
	if _, err := rand.Read(head[len(magic)+1:]); err != nil {
		return err
	}
	salt, nonce := head[len(magic)+1:len(magic)+1+16], head[len(magic)+1+16:]
	aead, err := gcm(master, salt)
	if err != nil {
		return err
	}
	h := sha256.Sum256(head)
	file := append(append(head, h[:]...), aead.Seal(nil, nonce, blob, nil)...)
	return os.WriteFile(path, file, 0o600)
}

func LoadFrom(path, master string) (*Vault, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256(data[:headLen-sha256.Size])
	if len(data) < headLen || string(data[:len(magic)]) != magic || data[len(magic)] != 1 || string(h[:]) != string(data[headLen-sha256.Size:headLen]) {
		return nil, errors.New("плохой или порченый заголовок файла (длина, magic, версия, sha256)")
	}
	salt, nonce := data[len(magic)+1:len(magic)+1+16], data[len(magic)+1+16:headLen-sha256.Size]
	aead, err := gcm(master, salt)
	if err != nil {
		return nil, err
	}
	blob, err := aead.Open(nil, nonce, data[headLen:], nil)
	if err != nil {
		return nil, fmt.Errorf("расшифровка не удалась (мастер-пароль или файл): %w", err)
	}
	v := NewVault()
	if err := json.Unmarshal(blob, &v.entries); err != nil {
		return nil, fmt.Errorf("повреждённый JSON внутри: %w", err)
	}
	return v, nil
}
