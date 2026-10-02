# Глава 17. Secchan — защищённый канал (sealed box)

> По мотивам `secchan.pdf`. Аналог libsodium sealed boxes на stdlib:
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
