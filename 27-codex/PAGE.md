# Глава 27. Codex — извлечение фрагментов кода

> Автору книги нужны иллюстрации к тексту:
> вырезать из большого исходника фрагменты «строка N, длина M» и превратить
> в картинки. Мы делаем ядро утилиты: парсинг спецификаций и вырезание строк.

## Цель главы

Разбор аргументов командной строки формата `from:count`, аккуратные ошибки
с именем виновника, работа со срезами с 1-индексацией.

## Теория

- `strings.Cut(s, ":")` делит строку один раз: `from, count, ok`.
- `strconv.Atoi` + проверка диапазона: `from ≥ 1`, `count ≥ 1`.
- Ошибка обязана называть аргумент: `fmt.Errorf("codex: spec %q: ...", arg)`
  — по сообщению видно, какой аргумент кривой.
- Срез `lines[from-1 : from-1+count]` — не забудь проверить
  `from+count-1 ≤ len(lines)` до вырезания, иначе паника на границах.

## Условие

Извлекатель фрагментов кода.

1. Аргументы вида `from:count`: from — номер первой строки (нумерация с 1),
   count — число строк; невалидный аргумент — ошибка с именем этого
   аргумента.
2. `Extract` возвращает строки `From .. From+Count-1` исходного слайса;
   выход за границы файла — ошибка.
3. Пустой список аргументов — пустой результат, без ошибки.

## API для реализации

```go
type Fragment struct {
	From  int // номер первой строки, 1-индексация
	Count int // число строк
}
func ParseSpecs(args []string) ([]Fragment, error)
func Extract(lines []string, f Fragment) ([]string, error)
```

Сигнатуры и договоры — в `codex.go`. Тесты — в `codex_test.go`.

## Подсказки

- `Extract` возвращает срез исходного слайса — копия не нужна.
- `"1:2:3"` невалиден: после `Cut` останется `"2:3"`, `Atoi` его не съест.
- В ошибке `Extract` полезно назвать размер файла: сразу видно, где вышли.

## Проверка

```
go test ./27-codex/
go test -tags solution ./27-codex/
```

## Эталон

`go test -tags solution ./27-codex/`

## Углубление

- Полный CLI с подкомандой `convert`: `codex convert sourcefile graphicbase
  from:count ...`, глобальные `--verbose` и `-h/--help` —
  `flag.NewFlagSet` + кастомный `Usage`.
- Выходные файлы `graphicbase.<N>.png` — по одному на фрагмент, нумерация
  с 0 (`fmt.Sprintf("%s.%d.png", ...)`); лог хода «Will create out.0.png
  with lines 34 to 42» — `log`.
- Рендер фрагмента в PNG: `image`, `image/draw`, `image/png` из stdlib, но
  тексту нужен шрифт — внешний модуль `golang.org/x/image/font/basicfont`;
  вариант без зависимостей — выдавать listing-текст вместо картинки.
- Эталон автора — репозиторий gitlab.com/RajaSrinivasan/codex (сеть).

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./27-codex/
```

## Заготовки (реализуйте)

### `codex.go`

{% raw %}
```go
//go:build !solution

// Пакет codex — глава 27. Реализуй ParseSpecs и Extract так, чтобы прошли
// тесты в codex_test.go. Условие — в README.md этой папки.
package codex

// Fragment — спецификация фрагмента кода: Count строк начиная со строки
// From (нумерация с 1).
type Fragment struct {
	From  int // номер первой строки, 1-индексация
	Count int // число строк, ≥ 1
}

// ParseSpecs разбирает аргументы вида "from:count" в фрагменты.
// count обязана быть ≥ 1; невалидный аргумент — ошибка, называющая его.
func ParseSpecs(args []string) ([]Fragment, error) {
	panic("TODO: ParseSpecs")
}

// Extract возвращает строки lines с From по From+Count-1 (1-индексация).
// Выход за границы файла (From < 1, Count < 1 или From+Count-1 > len(lines))
// — ошибка.
func Extract(lines []string, f Fragment) ([]string, error) {
	panic("TODO: Extract")
}
```
{% endraw %}

## Тесты

### `codex_test.go`

{% raw %}
```go
package codex

import (
	"strconv"
	"strings"
	"testing"
)

func TestParseSpecsValid(t *testing.T) {
	fs, err := ParseSpecs([]string{"1:5", "34:9", "71:11"})
	if err != nil {
		t.Fatalf("ParseSpecs: %v", err)
	}
	want := []Fragment{{From: 1, Count: 5}, {From: 34, Count: 9}, {From: 71, Count: 11}}
	if len(fs) != len(want) {
		t.Fatalf("ParseSpecs: фрагментов %d, хотел %d", len(fs), len(want))
	}
	for i := range want {
		if fs[i] != want[i] {
			t.Errorf("ParseSpecs: [%d] = %+v, хотел %+v", i, fs[i], want[i])
		}
	}
	if fs, err = ParseSpecs(nil); err != nil || len(fs) != 0 {
		t.Errorf("ParseSpecs(без аргументов) = %+v, %v; хотел пусто и nil", fs, err)
	}
}

func TestParseSpecsErrors(t *testing.T) {
	bad := []string{"34", "a:9", "34:x", "0:9", "34:0", "34:-1", "1:2:3", ":"}
	for _, arg := range bad {
		fs, err := ParseSpecs([]string{"1:5", arg})
		if err == nil {
			t.Errorf("ParseSpecs(%q): хотел ошибку, получил %+v", arg, fs)
			continue
		}
		// Ошибка обязана называть виновника.
		if !strings.Contains(err.Error(), arg) {
			t.Errorf("ParseSpecs(%q): ошибка %q не называет аргумент", arg, err)
		}
	}
}

func testLines(t *testing.T) []string {
	t.Helper()
	var lines []string
	for i := 1; i <= 10; i++ {
		lines = append(lines, "L"+strconv.Itoa(i))
	}
	return lines
}

func TestExtract(t *testing.T) {
	lines := testLines(t)

	got, err := Extract(lines, Fragment{From: 3, Count: 4})
	if err != nil {
		t.Fatalf("Extract(3:4): %v", err)
	}
	if want := "L3|L4|L5|L6"; strings.Join(got, "|") != want {
		t.Errorf("Extract(3:4) = %q, хотел %q", strings.Join(got, "|"), want)
	}

	// Последний фрагмент: до конца файла.
	got, err = Extract(lines, Fragment{From: 7, Count: 4})
	if err != nil {
		t.Fatalf("Extract(7:4): %v", err)
	}
	if want := "L7|L8|L9|L10"; strings.Join(got, "|") != want {
		t.Errorf("Extract(7:4) = %q, хотел %q", strings.Join(got, "|"), want)
	}

	// Одна строка.
	got, err = Extract(lines, Fragment{From: 10, Count: 1})
	if err != nil || len(got) != 1 || got[0] != "L10" {
		t.Errorf("Extract(10:1) = %q, %v; хотел [L10], nil", got, err)
	}
}

func TestExtractOutOfBounds(t *testing.T) {
	lines := testLines(t)
	for _, f := range []Fragment{
		{From: 8, Count: 4},  // 8+4-1 = 11 > 10
		{From: 11, Count: 1}, // начинается за файлом
		{From: 0, Count: 1},  // 1-индексация: 0 — не строка
		{From: 5, Count: 0},  // пустой фрагмент запрещён
		{From: 5, Count: -1},
	} {
		if got, err := Extract(lines, f); err == nil {
			t.Errorf("Extract(%+v) = %q, хотел ошибку за границей файла из %d строк", f, got, len(lines))
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

// Эталонное решение главы 27. Сравни со своей реализацией в codex.go.
package codex

import (
	"fmt"
	"strconv"
	"strings"
)

// Fragment — спецификация фрагмента кода.
type Fragment struct {
	From  int
	Count int
}

func ParseSpecs(args []string) ([]Fragment, error) {
	out := make([]Fragment, 0, len(args))
	for _, arg := range args {
		fromS, countS, ok := strings.Cut(arg, ":")
		if !ok {
			return nil, fmt.Errorf("codex: spec %q: хочу from:count", arg)
		}
		from, err1 := strconv.Atoi(fromS)
		count, err2 := strconv.Atoi(countS)
		if err1 != nil || err2 != nil || from < 1 || count < 1 {
			return nil, fmt.Errorf("codex: spec %q: from и count — целые ≥ 1", arg)
		}
		out = append(out, Fragment{From: from, Count: count})
	}
	return out, nil
}

func Extract(lines []string, f Fragment) ([]string, error) {
	if f.From < 1 || f.Count < 1 || f.From+f.Count-1 > len(lines) {
		return nil, fmt.Errorf("codex: фрагмент %d:%d выходит за границы файла из %d строк", f.From, f.Count, len(lines))
	}
	return lines[f.From-1 : f.From-1+f.Count], nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/27-codex`, упражнение **codex**. Локально: `go test -tags solution ./27-codex/`._
