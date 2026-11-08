# Глава 17. Secchan — защищённый канал (sealed box)

> Аналог libsodium sealed boxes на stdlib:
> отправитель шифрует под публичным ключом получателя, не имея своего
> долгосрочного секрета.

## Цель главы

Освоить гибридную криптографию: ECDH (`crypto/ecdh`, X25519) для общего
секрета, HKDF-SHA-256 для ключа, AES-GCM для сообщения — всё из stdlib.

## Теория

- Sealed box: отправитель создаёт **эфемерную** пару X25519, считает ECDH
  с публичным ключом получателя и кладёт свой эфемерный публичный ключ
  в сообщение. Получатель повторяет ECDH своим секретом.
- Общий секрет ECDH нельзя использовать как ключ напрямую: прогоняем
  его через HKDF (`crypto/hkdf`) с info-строкой — получаем 32 байта AES-256.
- AES-GCM даёт и конфиденциальность, и целостность; nonce (12 байт)
  кладём рядом с шифротекстом.

## Условие

1. `Party` — пара ключей X25519 (`crypto/rand`); `Public` — публичный ключ.
2. `Seal`: эфемерная пара отправителя → ECDH с публичным ключом получателя
   → HKDF-SHA-256 → AES-GCM. Выход: `ephemeralPub[32] || nonce[12] ||
   ciphertext`.
3. `Open` у получателя разбирает запись и расшифровывает.
4. Ошибки: неверная длина публичного ключа, короткое сообщение, чужой
   получатель, порча ciphertext.

## API для реализации

```go
type Party struct{ /* приватный ключ ecdh.X25519 */ }
func NewParty() (Party, error)
func (p Party) Public() []byte
func Seal(receiverPublic []byte, msg []byte) ([]byte, error)
func Open(receiver Party, sealed []byte) ([]byte, error)
```

Договоры — в `secchan.go`. Тесты — в `secchan_test.go`.

## Подсказки

- `ecdh.X25519().GenerateKey(rand.Reader)`, `NewPublicKey(b)`, `ECDH`.
- HKDF: `hkdf.Key(sha256.New, secret, nil, "secchan/v1", 32)`.
- Проверь длину `sealed` до нарезки: минимум 32+12+16 байт.

## Проверка

```
go test ./17-secchan/                 # красный: TODO
go test -tags solution ./17-secchan/ # зелёный: эталон
```

## Эталон

`solution.go` собирается только с тегом `solution`.

## Углубление

- Установка сессии поверх `net.Conn` (`net`): клиент первым сообщением
  присылает свой публичный ключ, сервер отвечает своим — обмен ключами
  вместо долговременных пар у обеих сторон.
- Фрейминг сообщений: спека кладёт поле размера payload перед шифротекстом —
  length-prefix (`encoding/binary`, `io.ReadFull`) для записей разной длины.
- Пример приложения `secpwdgen`: сервер-генератор паролей, клиент задаёт
  сегменты/разделитель/количество (`-s`, `-r`, `-n`); транзакции строго
  «запрос клиента → ответ сервера» (`flag`).
- Иные паттерны допустимы, пока каждый payload запечатан публичным ключом
  партнёра.
- Оригинал — libsodium sealed boxes и Ada-биндинги (внешняя зависимость);
  в главе всё заменено на stdlib: `crypto/ecdh`, `crypto/hkdf`, `crypto/aes`.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./17-secchan/
```

## Заготовки (реализуйте)

### `secchan.go`

{% raw %}
```go
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
```
{% endraw %}

## Тесты

### `secchan_test.go`

{% raw %}
```go
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
```
{% endraw %}

## Эталонное решение

Сначала попробуйте сами. Решение скрыто — раскрывайте только для сверки.

<details>
<summary><strong>Эталон (`solution.go`)</strong> — нажмите, чтобы показать</summary>

{% raw %}
```go
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
```
{% endraw %}


</details>

_Пакет `assignments-go/17-secchan`, упражнение **secchan**. Локально: `go test -tags solution ./17-secchan/`._
