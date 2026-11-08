# Глава 13. CRC — циклический избыточный код

> Контрольные суммы в линиях связи: считаем
> CRC-32 IEEE своими битами, сверяемся со stdlib. Эксперименты
> с шумом и полиномами — в wants.

## Цель главы

Написать CRC-32 «на бумажных битах»: полином IEEE (reflected
`0xEDB88320`), init `0xFFFFFFFF`, xorout `0xFFFFFFFF`. Правильность
докажем эталонным вектором `CRC32("123456789") == 0xCBF43926`
и совпадением с `hash/crc32.ChecksumIEEE` (импорт — только в тестах).

## Теория

- CRC — остаток от деления сообщения на полином в GF(2). Reflected-
  схема обрабатывает байты младшим битом вперёд: `crc ^= byte`,
  затем 8 раз `crc>>1 ^ 0xEDB88320`, если младший бит был 1.
- Init `0xFFFFFFFF` и финальный `^0xFFFFFFFF` — стандарт IEEE 802.3
  (тот же, что у `hash/crc32.ChecksumIEEE` и zip/PNG).
- Ошибка обнаруживается, пока полином «не делит» повреждение —
  потому в `crc.pdf` и играют шумом: сколько бит можно испортить,
  чтобы CRC изменился (любой один бит — гарантированно).
- Тест-вектор `0xCBF43926` для `"123456789"` — общепринятая проверка
  всех CRC-32 реализаций (catalogue of CRCs).

## Условие

1. `CRC32(data)` — контрольная сумма байтов data по схеме IEEE:
   полином `0xEDB88320` (reflected), init `0xFFFFFFFF`,
   xorout `0xFFFFFFFF`. Только своя битовая реализация —
   `hash/crc32` в решении запрещён (в тестах — можно).
2. Пустой ввод — допустим: `CRC32(nil)` == CRC32 пустого сообщения.

Wants (необязательно, из `crc.pdf`): CRC-16 с выбором полинома,
таблица `-t`, шум `-n` (сколько бит переживает проверку).

## API для реализации

```go
func CRC32(data []byte) uint32
```

Сигнатуры и договоры — в `crc.go`. Тесты — в `crc_test.go`.

## Подсказки

- Тело цикла: `for each byte { crc ^= uint32(b); for 8 { if crc&1==1
  { crc = crc>>1 ^ 0xEDB88320 } else { crc >>= 1 } } }`; в конце
  `return ^crc` — это и есть xorout.
- Прогони цикл на бумаге для одного байта, прежде чем писать код.
- Сравнение со stdlib — в тестах: `crc32.ChecksumIEEE(data)`
  из `hash/crc32` — это разрешённый «измерительный прибор».

## Проверка

```
go test ./13-crc/
go test -tags solution ./13-crc/
```

## Эталон

`go test -tags solution ./13-crc/`.

## Углубление

- Главный герой спеки — CRC-16: `-p/--polynomial` (16 бит; дефолт —
  reflected 0xA001, CCITT — 0x1021), своя реализация + таблица
  256 × `uint16`.
- `-t/--table`: печать начальной таблицы полинома — 32 строки по 16
  значений, `fmt.Printf("%04x")`.
- Источники данных: `-x/--hex-string` — `encoding/hex`; `-f/--file` —
  `os.ReadFile`; `-l/--lines` — построчно (null-terminated), hex+lines
  комбинируются.
- `-n/--noise N`: инъекция N случайных битовых ошибок (случайный бит
  в случайном байте) — `math/rand`; чувствительность полинома к
  ошибкам — суть экспериментов спеки.
- `-v/--verbose`: эхо ввода, без него печатается только CRC. Теория
  выбора полиномов — ross.net/crc «crc_v3.txt» (сеть).

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./13-crc/
```

## Заготовки (реализуйте)

### `crc.go`

{% raw %}
```go
//go:build !solution

// Пакет crc — глава 13. Реализуй функцию так, чтобы прошли тесты
// в crc_test.go. Условие — в README.md этой папки.
package crc

// CRC32 возвращает контрольную сумму data по схеме CRC-32/IEEE:
// полином 0xEDB88320 (reflected), init 0xFFFFFFFF, xorout 0xFFFFFFFF.
// Только битовая реализация — hash/crc32 здесь не помощник.
// Пустой ввод — корректный случай: контрольная сумма пустого сообщения.
func CRC32(data []byte) uint32 {
	panic("TODO: CRC32")
}
```
{% endraw %}

## Тесты

### `crc_test.go`

{% raw %}
```go
package crc

import (
	"hash/crc32"
	"testing"
)

func TestKnownVector(t *testing.T) {
	// Общепринятый тест-вектор CRC-32/IEEE (catalogue of CRCs).
	if got := CRC32([]byte("123456789")); got != 0xCBF43926 {
		t.Errorf("CRC32(\"123456789\"): получил %#08x — хотел 0xCBF43926", got)
	}
}

func TestEmpty(t *testing.T) {
	// CRC32 пустого сообщения == xorout-обнулённый init == 0x00000000.
	if got := CRC32(nil); got != 0 {
		t.Errorf("CRC32(nil): получил %#08x — хотел 0x00000000", got)
	}
	if got := CRC32([]byte{}); got != 0 {
		t.Errorf("CRC32(пусто): получил %#08x — хотел 0x00000000", got)
	}
}

func TestMatchesStdlib(t *testing.T) {
	// stdlib — измерительный прибор: совпадение доказывает схему.
	zero1024 := make([]byte, 1024)
	zero1024[512] = 1 // один-единственный бит в 1024 нулях
	messages := [][]byte{
		[]byte("a"),
		[]byte("abc"),
		[]byte("abcdefghijklmnopqrstuvwxyz"),
		[]byte("The quick brown fox jumps over the lazy dog"),
		{0x00, 0x01, 0x02, 0xFE, 0xFF},
		{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
		zero1024,
	}
	for _, m := range messages {
		if got, want := CRC32(m), crc32.ChecksumIEEE(m); got != want {
			t.Errorf("CRC32(%q): получил %#08x — хотел %#08x (stdlib)", m, got, want)
		}
	}
}

func TestSingleBitChange(t *testing.T) {
	// Один испорченный бит меняет контрольную сумму — потому CRC
	// и ловит ошибки линий (want из crc.pdf — шум).
	orig := []byte("perfect-29787-Dendrite-38203")
	c0 := CRC32(orig)
	for i := range orig {
		broken := append([]byte(nil), orig...)
		broken[i] ^= 0x01
		if c := CRC32(broken); c == c0 {
			t.Errorf("порча байта %d не меняет CRC — полином проморгал ошибку", i)
		}
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

// Эталонное решение главы 13. Сравни со своей реализацией в crc.go.
package crc

// CRC32 — битовая реализация CRC-32/IEEE:
// полином 0xEDB88320 (reflected), init 0xFFFFFFFF, xorout 0xFFFFFFFF.
func CRC32(data []byte) uint32 {
	const poly = 0xEDB88320
	crc := uint32(0xFFFFFFFF)
	for _, b := range data {
		crc ^= uint32(b)
		for range 8 {
			if crc&1 == 1 {
				crc = crc>>1 ^ poly
			} else {
				crc >>= 1
			}
		}
	}
	return ^crc // xorout 0xFFFFFFFF
}
```
{% endraw %}


</details>

_Пакет `assignments-go/13-crc`, упражнение **crc**. Локально: `go test -tags solution ./13-crc/`._
