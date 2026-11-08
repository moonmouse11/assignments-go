# Глава 1. Hello — первый CLI-инструмент

> Первый projectlet — «привет, мир» инженера:
> спецификация → код → тесты, без внешних зависимостей.

## Цель главы

Настроить окружение и пройти полный цикл разработки: прочитать требования,
написать утилиту, проверить тестами. С этого начинается каждый проект в этом
курсе.

## Теория

- `os.Open` открывает файл; не забудь `defer f.Close()`.
- `bufio.Scanner` читает файл построчно (`sc.Scan()` / `sc.Text()`).
- `fmt.Fprintf(w, "%04d %s\n", n, line)` — формат `NNNN` с паддингом нулями.
- Вывод всегда через `io.Writer` (stdout, `bytes.Buffer`, файл) — так утилиту
  легко тестировать без захвата настоящего stdout.

## Условие

Утилита `hello` — командная строка, каждый аргумент — имя файла.

1. Каждый аргумент трактуется как имя файла.
2. Утилитет печатает строки каждого файла с номерами.
3. Номер — 4 цифры, слева добитые нулями: `NNNN`, затем пробел и строка: `0007 текст`.
4. Нумерация **начинается заново для каждого файла** (с `0001`).
5. Вывод — в `stdout` (можно перенаправить в файл-листинг).
6. Если файл не читается — вернуть ошибку (не молчать).

Wants (необязательно, повысит оценку): печатать полный путь файла
(`filepath.Abs`) и дату создания перед содержимым.

## API для реализации

```go
func NumberLines(lines []string) string
func Run(filenames []string, out io.Writer) error
```

Сигнатуры и договоры — в `hello.go`. Тесты — в `hello_test.go`.

## Подсказки

- `os.WriteFile` и `t.TempDir()` в тестах — быстрые фикстуры.
- Пустая строка тоже строка: `0003 ` (номер, пробел, пусто).

## Проверка

```
go test ./01-hello/
go test -tags solution ./01-hello/
```

## Эталон

`solution.go` собирается только с тегом `solution`:
`go test -tags solution ./01-hello/`.

## Углубление

- Wants спеки: полный путь файла (`filepath.Abs`), дата создания, размер —
  заголовком перед листингом и только по явному запросу (флаг `-l`).
- У автора утилита называется `lister` — собери CLI: `flag.Parse()`,
  `flag.Arg(i)`; вывод примера имеет формат `NNNN : строка` (с пробелами
  вокруг двоеточия).

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./01-hello/
```

## Заготовки (реализуйте)

### `hello.go`

{% raw %}
```go
//go:build !solution

// Пакет hello — глава 1. Реализуй функции так, чтобы прошли тесты
// в hello_test.go. Условие — в README.md этой папки.
package hello

import "io"

// NumberLines возвращает строки lines, пронумерованные с 1
// в формате NNNN: "0007 текст". Строки разделяются \n,
// в конце — завершающий \n.
//
//	NumberLines([]string{"a", "b"}) == "0001 a\n0002 b\n"
func NumberLines(lines []string) string {
	panic("TODO: NumberLines")
}

// Run печатает в out каждый файл из filenames построчно с номерами
// (нумерация с 0001 заново для каждого файла, формат как в NumberLines).
// Возвращает ошибку, если файл не открывается или читается.
func Run(filenames []string, out io.Writer) error {
	panic("TODO: Run")
}
```
{% endraw %}

## Тесты

### `hello_test.go`

{% raw %}
```go
package hello

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestNumberLines(t *testing.T) {
	got := NumberLines([]string{"alpha", "beta", "", "gamma"})
	want := "0001 alpha\n0002 beta\n0003 \n0004 gamma\n"
	if got != want {
		t.Errorf("NumberLines:\nполучил %q\nхотел  %q", got, want)
	}
}

func TestRun(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.txt")
	f2 := filepath.Join(dir, "b.txt")
	os.WriteFile(f1, []byte("one\ntwo\n"), 0o644)
	os.WriteFile(f2, []byte("three"), 0o644)

	var buf bytes.Buffer
	if err := Run([]string{f1, f2}, &buf); err != nil {
		t.Fatalf("Run: %v", err)
	}
	// Нумерация начинается заново для каждого файла.
	want := "0001 one\n0002 two\n0001 three\n"
	if buf.String() != want {
		t.Errorf("Run:\nполучил %q\nхотел  %q", buf.String(), want)
	}
}

func TestRunMissingFile(t *testing.T) {
	if err := Run([]string{filepath.Join(t.TempDir(), "nope.txt")}, io.Discard); err == nil {
		t.Fatal("Run: хотел ошибку для отсутствующего файла, получил nil")
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

// Эталонное решение главы 1. Сравни со своей реализацией в hello.go.
package hello

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func NumberLines(lines []string) string {
	var b strings.Builder
	for i, line := range lines {
		fmt.Fprintf(&b, "%04d %s\n", i+1, line)
	}
	return b.String()
}

func Run(filenames []string, out io.Writer) error {
	for _, name := range filenames {
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		n := 0
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			n++
			fmt.Fprintf(out, "%04d %s\n", n, sc.Text())
		}
		if err := sc.Err(); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/01-hello`, упражнение **hello**. Локально: `go test -tags solution ./01-hello/`._
