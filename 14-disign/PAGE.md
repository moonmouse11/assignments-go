# Глава 14. DiSign — подпись документов

> Подпись файла приватным ключом и проверка
> публичным — модель openssl dgst -sign/-verify. В оригинале RSA;
> здесь — ed25519: быстрее, короче, современнее.

## Цель главы

Сгенерировать пару ключей (crypto/rand — только настоящий источник),
подписать документ и проверить подпись: round-trip true, порча данных
или подписи — false, чужой ключ — false. Подпись — это аутентичность:
кто подписал и что документ не меняли.

## Теория

- `crypto/ed25519`: `ed25519.GenerateKey(rand.Reader)` — пара ключей;
  приватный — 64 байта (семя + публичный), публичный — 32 байта.
- `ed25519.Sign(priv, msg)` — подпись 64 байта; `ed25519.Verify(pub,
  msg, sig)` — bool. Подпись зависит и от ключа, и от сообщения:
  изменился любой — проверка false.
- ed25519 сам хеширует сообщение (внутри SHA-512) — отдельный шаг
  «sha256 от файла» из `disign.pdf` встроен; ed25519 не требует
  внешнего хеша.
- `crypto/rand.Reader` — ключи каждый раз разные; ошибка чтения — наружу.

## Условие

1. `GenerateKey()` — пара ключей из crypto/rand; ошибка источника —
  наружу.
2. `SignFile(kp, data)` — подпись байтов data приватным ключом kp;
  некорректный приватный ключ — nil.
3. `VerifyFile(pub, data, sig)` — true только для родных ключа и данных.

Wants (необязательно, из `disign.pdf`): RSA-вариант, пароль на
приватный ключ, запись подписи в `<файл>.sig` на диске.

## API для реализации

```go
type KeyPair struct {
	Public  ed25519.PublicKey
	Private ed25519.PrivateKey
}
func GenerateKey() (KeyPair, error)
func SignFile(kp KeyPair, data []byte) []byte
func VerifyFile(pub ed25519.PublicKey, data, sig []byte) bool
```

Сигнатуры и договоры — в `disign.go`. Тесты — в `disign_test.go`.

## Подсказки

- `ed25519.PrivateKeySize` и `ed25519.SignatureSize` — константы
  stdlib для проверок длин; `len(kp.Private) != ed25519.PrivateKeySize`
  → nil из SignFile.
- Порча в тестах: `copy(sig, ...); sig[0] ^= 0xFF` — один перевёрнутый
  байт хватает.
- Тесты не зависят от того, какие ключи сгенерированы: проверяются
  отношения (подпись/проверка), а не конкретные байты.

## Проверка

```
go test ./14-disign/
go test -tags solution ./14-disign/
```

## Эталон

`go test -tags solution ./14-disign/`.

## Углубление

- Подкоманды `sign`/`authenticate` из спеки: `flag.NewFlagSet`;
  аргумент — список файлов (arg1 arg2 …, `flag.Args()`), не один.
- RSA-вариант спеки (модель — `openssl dgst -sign/-verify`):
  `crypto/rsa` + `crypto/sha256` (`rsa.SignPKCS1v15`); сверка с
  эталонным openssl — внешняя зависимость.
- Дефолтные ключи: `--private` → `~/.ssh/id_rsa`, `--public` →
  `~/.ssh/id_rsa.pub` — `os.UserHomeDir` + `os.ReadFile`.
- Passphrase на приватный ключ (wants): openssh/PEM-форматы —
  `golang.org/x/crypto/ssh` (внешняя зависимость).
- Подпись в `[basename].sig` рядом с файлом (`strings.TrimSuffix` +
  `filepath.Ext`), лог с таймстампами — `log`.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./14-disign/
```

## Заготовки (реализуйте)

### `disign.go`

{% raw %}
```go
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
```
{% endraw %}

## Тесты

### `disign_test.go`

{% raw %}
```go
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
```
{% endraw %}

## Эталонное решение

Сначала попробуйте сами. Решение скрыто — раскрывайте только для сверки.

<details>
<summary><strong>Эталон (`solution.go`)</strong> — нажмите, чтобы показать</summary>

{% raw %}
```go
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
```
{% endraw %}


</details>

_Пакет `assignments-go/14-disign`, упражнение **disign**. Локально: `go test -tags solution ./14-disign/`._
