# Глава 2. Dump — hex-дамп и конвенции CLI

> Файлы как бинарные потоки + шаблон командной утилиты.

## Цель главы

Научиться смотреть на файл как на поток байтов, освоить hex-форматирование и
обработку аргументов командной строки. Появляются **конвенции CLI** — им
будут следовать все следующие утилиты курса.

## Теория

- `io.Reader` — поток байтов; `io.ReadFull` читает ровно N байт.
- `fmt.Fprintf(w, "%02x", b)` — байт в hex; `"%04x"` — смещение.
- ASCII-печатными считаем байты `0x20..0x7e`, всё остальное в текстовой
  колонке — `.`.
- Ошибки возвращаем, а не печатаем: `fmt.Errorf(...)`.

## Условие

Утилита `dump` показывает содержимое файла в виде дампа и в тексте.

1. Данные показываются **блоками**: в каждой строке — один блок.
2. Размер блока по умолчанию — 32; допустимы только 8, 16, 32
   (перекрывается переключателем `-b`/`--blocksize`; иначе — ошибка).
3. Формат строки: смещение блока в hex (`0000:`), затем байты блока в hex
   через пробел, затем `|`, текстовое представление блока, `|`:

   ```
   0000: 41 42 00 43 44 |AB.CD|
   ```

4. В текстовой колонке непечатаемые байты показываются как `.`.
5. Последний (неполный) блок печатается с его реальной длиной.

### Конвенции CLI (для `main`, опционально в этой главе)

- без аргументов — печать help;
- `-h`/`--help` — help;
- версия в форме `name-major.minor` (см. `Version`).

## API для реализации

```go
const Version = "dump-0.1"
func FormatBlock(offset int, data []byte) string
func Dump(r io.Reader, blockSize int, out io.Writer) error
```

## Подсказки

- `FormatBlock` — чистая функция: удобно гонять таблицей в тестах.
- В `Dump` читай циклом по `blockSize`, пока не кончатся данные.
- Смещение — номер **байта** начала блока, не номер блока.

## Проверка

```
go test ./02-dump/
go test -tags solution ./02-dump/
```

## Эталон

`go test -tags solution ./02-dump/`

## Углубление

- `-o`/`--octal` — октальный вывод байтов вместо hex (`%04o` на байт).
- Заголовок перед дампом каждого файла: полный путь, размер, дата создания.
- Вывод автора шире: 8 байт в блоке, hex без пробелов, смещение `%08x` —
  сравни форматы, спека оставляет пространство для вариаций.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./02-dump/
```

## Заготовки (реализуйте)

### `dump.go`

{% raw %}
```go
//go:build !solution

// Пакет dump — глава 2. Реализуй функции так, чтобы прошли тесты.
package dump

import "io"

// Version — версия утилиты в форме name-major.minor (конвенция из спецификации).
const Version = "dump-0.1"

// FormatBlock формирует одну строку дампа блока data, начинающегося
// с байта offset:
//
//	"0000: 41 42 00 43 |AB.C|"
//
// Смещение — 4 hex-цифры и двоеточие, байты — по два hex-символа через
// пробел, затем " |", текстовая колонка (непечатаемые байты — '.'),
// "|" и "\n".
func FormatBlock(offset int, data []byte) string {
	panic("TODO: FormatBlock")
}

// Dump пишет в out дамп всех данных из r блоками по blockSize байт
// (строки в формате FormatBlock, смещение — номер байта начала блока).
// Допустимые blockSize: 8, 16, 32; иначе — ошибка.
func Dump(r io.Reader, blockSize int, out io.Writer) error {
	panic("TODO: Dump")
}
```
{% endraw %}

## Тесты

### `dump_test.go`

{% raw %}
```go
package dump

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatBlock(t *testing.T) {
	got := FormatBlock(0, []byte{'A', 'B', 0x00, 'C'})
	want := "0000: 41 42 00 43 |AB.C|\n"
	if got != want {
		t.Errorf("FormatBlock:\nполучил %q\nхотел  %q", got, want)
	}
	if got := FormatBlock(8, []byte{}); got != "0008: ||\n" {
		t.Errorf("FormatBlock пустой блок: получил %q", got)
	}
}

func TestDump(t *testing.T) {
	// 10 байт: два блока по 8, второй неполный.
	data := "ABCDEFGHIJ"
	var buf bytes.Buffer
	if err := Dump(strings.NewReader(data), 8, &buf); err != nil {
		t.Fatalf("Dump: %v", err)
	}
	want := "0000: 41 42 43 44 45 46 47 48 |ABCDEFGH|\n" +
		"0008: 49 4a |IJ|\n"
	if buf.String() != want {
		t.Errorf("Dump:\nполучил %q\nхотел  %q", buf.String(), want)
	}
}

func TestDumpNonPrintable(t *testing.T) {
	var buf bytes.Buffer
	if err := Dump(strings.NewReader("\x01\x7f\n"), 8, &buf); err != nil {
		t.Fatalf("Dump: %v", err)
	}
	if !strings.Contains(buf.String(), "|...|") {
		t.Errorf("непечатаемые байты должны быть '.': %q", buf.String())
	}
}

func TestDumpInvalidBlockSize(t *testing.T) {
	for _, size := range []int{0, 4, 12, 64, -8} {
		if err := Dump(strings.NewReader("data"), size, &bytes.Buffer{}); err == nil {
			t.Errorf("blockSize=%d: хотел ошибку, получил nil", size)
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

// Эталонное решение главы 2.
package dump

import (
	"fmt"
	"io"
	"strings"
)

const Version = "dump-0.1"

func FormatBlock(offset int, data []byte) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%04x:", offset)
	for _, c := range data {
		fmt.Fprintf(&b, " %02x", c)
	}
	b.WriteString(" |")
	for _, c := range data {
		if c >= 0x20 && c < 0x7f {
			b.WriteByte(c)
		} else {
			b.WriteByte('.')
		}
	}
	b.WriteString("|\n")
	return b.String()
}

func Dump(r io.Reader, blockSize int, out io.Writer) error {
	switch blockSize {
	case 8, 16, 32:
	default:
		return fmt.Errorf("недопустимый размер блока %d (допустимо 8, 16, 32)", blockSize)
	}
	buf := make([]byte, blockSize)
	for off := 0; ; off += blockSize {
		n, err := io.ReadFull(r, buf)
		if n > 0 {
			fmt.Fprint(out, FormatBlock(off, buf[:n]))
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
```
{% endraw %}


</details>

_Пакет `assignments-go/02-dump`, упражнение **dump**. Локально: `go test -tags solution ./02-dump/`._
