# Глава 15. PPM — защищённый формат файла

> Пароли нельзя хранить в открытом виде: разработаем
> формат файла, доступ к которому даёт только правильный пароль.

## Цель главы

Освоить симметричную криптографию из stdlib: KDF (`crypto/pbkdf2`),
AES-GCM (`crypto/aes`, `crypto/cipher`) и контроль целостности через
`crypto/sha256`.

## Теория

- Пароль — не ключ: его надо «растянуть» в 32 байта через PBKDF2 с солью.
  Соль (16 байт) и число итераций мешают атакам по словарям.
- AES-GCM — шифрование + имитовставок: подмена payload не останется
  незамеченной, nonce (12 байт) хранится в открытом виде — это нормально.
- SHA-256 от заголовка ловит порчу полей до шифротекста (magic, соль, nonce).

## Условие

Файл PPMF устроен так:

1. Заголовок: magic `PPMF` (4 байта), версия `1` (1 байт), `salt` (16),
   `nonce` (12), затем SHA-256 от этих 33 байт (32 байта) — итого 65.
2. Далее — AES-GCM-шифротекст payload.
3. Ключ = PBKDF2-HMAC-SHA-256(password, salt, 1000 итераций, 32 байта).

`Load` обязан отклонять: неверный magic/версию, порчу заголовка
(несовпадение SHA-256), неверный пароль и порчу payload.

## API для реализации

```go
func Save(payload []byte, password string, w io.Writer) error
func Load(r io.Reader, password string) ([]byte, error)
```

Сигнатуры и договоры — в `ppm.go`. Тесты — в `ppm_test.go`.

## Подсказки

- `pbkdf2.Key(sha256.New, password, salt, 1000, 32)` — с Go 1.24 в stdlib.
- `io.ReadAll` прочитает весь файл, дальше — проверка длин и полей.
- Ошибка `cipher.AEAD.Open` уже покрывает и неверный пароль, и порчу payload.

## Проверка

```
go test ./15-ppm/            # красный: TODO
go test -tags solution ./15-ppm/   # зелёный: эталон
```

## Эталон

`solution.go` собирается только с тегом `solution`.

## Углубление

- Полная спека — не только формат файла, а CLI-приложение Personal Password
  Manager: подкоманды `create`, `set`, `reset`, `show` и флаги
  `-p/--password-file`, `-f/--force`, `-k/--keep`, `-v/--verbose` — наведи
  CLI поверх `Save`/`Load` (`flag`).
- `reset` — смена пароля файла (re-keying): прочитать базу старым паролем,
  записать новым; пример дважды спрашивает «Password» и «New Password».
- Payload в спеке — простой CSV `context,username,password`, а не JSON:
  `encoding/csv` как альтернатива сериализации.
- `set -g` генерирует пароль при записи: встроенный wordlist, число сегментов,
  разделитель (`-s`, `-t`, `-m`) — реюз генератора из главы 10.
- Оригинал построен на OpenSSL (внешняя зависимость); в главе всё заменено
  на stdlib: `crypto/pbkdf2`, `crypto/aes`, `crypto/cipher`, `crypto/sha256`.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./15-ppm/
```

## Заготовки (реализуйте)

### `ppm.go`

{% raw %}
```go
//go:build !solution

// Пакет ppm — глава 15. Реализуй функции так, чтобы прошли тесты
// в ppm_test.go. Условие — в README.md этой папки.
package ppm

import "io"

// Формат файла PPMF: заголовок 65 байт, затем AES-GCM-шифротекст.
//
//	[0:4]    magic "PPMF"
//	[4:5]    версия формата, 1
//	[5:21]   salt (16 байт, PBKDF2)
//	[21:33]  nonce GCM (12 байт)
//	[33:65]  sha256 от байтов [0:33]
//	[65:]    AES-GCM(payload), ключ — pbkdf2(sha256, пароль, salt, iter, 32)

// iterations — число итераций PBKDF2. В учебных целях маленькое:
// тесты обязаны выполняться быстрее секунды.
const iterations = 1000

// Save шифрует payload паролем password (pbkdf2 + AES-GCM, случайные
// salt и nonce) и записывает файл формата PPMF в w. Возвращает ошибку
// при сбое записи или ГСЧ.
func Save(payload []byte, password string, w io.Writer) error {
	panic("TODO: Save")
}

// Load читает из r файл формата PPMF и расшифровывает payload паролем
// password. Ошибки: короткий файл, неверный magic/версия, порча
// заголовка (sha256 не сходится), неверный пароль или порча payload
// (ошибка расшифровки GCM).
func Load(r io.Reader, password string) ([]byte, error) {
	panic("TODO: Load")
}
```
{% endraw %}

## Тесты

### `ppm_test.go`

{% raw %}
```go
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
```
{% endraw %}

## Эталонное решение

Сначала попробуйте сами. Решение скрыто — раскрывайте только для сверки.

<details>
<summary><strong>Эталон (`solution.go`)</strong> — нажмите, чтобы показать</summary>

{% raw %}
```go
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
```
{% endraw %}


</details>

_Пакет `assignments-go/15-ppm`, упражнение **ppm**. Локально: `go test -tags solution ./15-ppm/`._
