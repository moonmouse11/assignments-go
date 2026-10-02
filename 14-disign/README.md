# Глава 14. DiSign — подпись документов

> По мотивам `disign.pdf`. Подпись файла приватным ключом и проверка
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
