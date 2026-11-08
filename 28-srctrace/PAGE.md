# Глава 28. Srctrace — трассировка к исходникам

> Приложение обязано уметь рассказать, из какого
> коммита оно собрано: генерируем фрагмент исходника с данными репозитория,
> компилируем его в приложение и разбираем обратно.

## Цель главы

Освоить генерацию кода (`go:generate` в миниатюре): превратить данные
репозитория в валидный Go-файл и написать парсер для обратного разора —
round-trip без потерь.

## Теория

- Трассировка — не лог, а константы в бинарнике: `Branch`, `ShortCommit`,
  `Commit`, `BuildDate` компилируются в приложение и видны в runtime.
- Генератор — обычная функция над `strings.Builder` и `fmt.Fprintf` с
  глаголом `%q`: он корректно закавычит любую строку.
- Парсер — regexp по строкам: `(?m)^имя = "значение"`; doc-комментарии
  (`// Branch — ...`) не совпадают с шаблоном и не мешают.
- Договор двусторонний: `Parse(Snippet(info)) == info` для любых данных.

## Условие

1. `Snippet(info)` возвращает валидный Go-файл: строка `package revisions`
   и блок `const` с константами `Branch`, `ShortCommit`, `Commit`,
   `BuildDate` — значения из `info`.
2. `Parse(src)` разбирает такой файл обратно в `RepoInfo`.
3. Некорректный источник — ошибка: пустой, без `package`, без какой-либо
   из четырёх констант.

Wants (необязательно): реальные данные из git, другие языки вывода.

## API для реализации

```go
type RepoInfo struct{ Branch, Short, Commit, Date string }
func Snippet(info RepoInfo) string
func Parse(src string) (RepoInfo, error)
```

Сигнатуры и договоры — в `srctrace.go`. Тесты — в `srctrace_test.go`.

## Подсказки

- `regexp.MustCompile` — один раз на уровне пакета, а не в функции.
- `package` ищи построчно (`strings.HasPrefix`), константы — регуляркой.
- Собери совпадения в map «имя → значение», потом проверь наличие всех
  четырёх имён — это и есть проверка «некорректного».

## Проверка

```
go vet ./28-srctrace/ && go test ./28-srctrace/            # красный: TODO
go vet -tags solution ./28-srctrace/ && go test -tags solution ./28-srctrace/
```

## Эталон

`solution.go` собирается только с тегом `solution`.

## Углубление

- Полей в спеке больше наших четырёх: `VERSION_MAJOR`/`VERSION_MINOR`
  (дефолт 0), `BUILD_NUMBER` (дефолт 9999), `REPO_NAME`, короткий
  `COMMIT_ID` и полный `COMMIT_ID_FULL`, `BUILD_TIME`.
- CLI: `flag` — `-major`/`-minor`/`-build` со значениями по умолчанию,
  `-lang` (C, C++, Ada, Python, Go), `-output=name` (дефолт `revisions`
  с расширением по языку: `revisions.h`, `revisions.ads`).
- Формат по языкам: C — блок `#define VERSION_MAJOR (2)`, Ada —
  `package revisions is ... end revisions;` — та же генерация на
  `strings.Builder` и `fmt.Fprintf`, stdlib.
- Реальное наполнение `RepoInfo` — `git rev-parse --abbrev-ref HEAD`,
  `git log -1` через `os/exec`, режим `go:generate`; wants — отчёт
  по конкретному файлу (аргумент CLI). Git — внешняя зависимость,
  в этой главе не трогаем.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./28-srctrace/
```

## Заготовки (реализуйте)

### `srctrace.go`

{% raw %}
```go
//go:build !solution

// Пакет srctrace — глава 28. Реализуй функции так, чтобы прошли тесты
// в srctrace_test.go. Условие — в README.md этой папки.
package srctrace

// RepoInfo — сведения о репозитории на момент сборки.
type RepoInfo struct {
	Branch string // имя ветки, например "master"
	Short  string // короткий id коммита, например "cefa722"
	Commit string // полный id коммита (40 hex-символов)
	Date   string // дата/время сборки, например RFC3339
}

// Snippet возвращает валидный Go-файл (строка package + блок const
// с константами Branch, ShortCommit, Commit, BuildDate), компилируемый
// в приложение и хранящий данные info. Обратный разбор — Parse.
func Snippet(info RepoInfo) string {
	panic("TODO: Snippet")
}

// Parse разбирает Go-файл, сгенерированный Snippet, и возвращает
// исходные данные. Некорректный источник — ошибка: пустой, без строки
// package, без какой-либо из четырёх констант.
func Parse(src string) (RepoInfo, error) {
	panic("TODO: Parse")
}
```
{% endraw %}

## Тесты

### `srctrace_test.go`

{% raw %}
```go
package srctrace

import (
	"strings"
	"testing"
)

var sample = RepoInfo{
	Branch: "master",
	Short:  "cefa722",
	Commit: "cefa72267cc4d3a07fbf5e672b0053116d582aa7",
	Date:   "2026-10-01T20:18:00Z",
}

func TestSnippetContainsValues(t *testing.T) {
	src := Snippet(sample)
	for _, v := range []string{sample.Branch, sample.Short, sample.Commit, sample.Date} {
		if !strings.Contains(src, `"`+v+`"`) {
			t.Errorf("Snippet не содержит значение %q:\n%s", v, src)
		}
	}
	if !strings.Contains(src, "package ") {
		t.Error("Snippet без строки package — не Go-файл:\n", src)
	}
	for _, name := range []string{"Branch", "ShortCommit", "Commit", "BuildDate"} {
		if !strings.Contains(src, name) {
			t.Errorf("Snippet без константы %s:\n%s", name, src)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	src := Snippet(sample)
	got, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse(Snippet(...)): %v", err)
	}
	if got != sample {
		t.Errorf("round-trip:\nполучил %+v\nхотел  %+v", got, sample)
	}
}

func TestParseErrors(t *testing.T) {
	bad := []struct {
		name, src string
	}{
		{"пустой источник", ""},
		{"только пробелы", "   \n\t\n"},
		{"без package", "const Branch = \"main\"\n"},
		{"без констант", "package revisions\n"},
		{"не хватает Commit", "package revisions\n\nconst (\n\tBranch = \"main\"\n\tShortCommit = \"a\"\n\tBuildDate = \"d\"\n)\n"},
		{"значение без кавычек", "package revisions\nBranch = без кавычек\n"},
	}
	for _, tc := range bad {
		if _, err := Parse(tc.src); err == nil {
			t.Errorf("Parse(%s): хотел ошибку, получил nil", tc.name)
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

// Эталонное решение главы 28. Сравни со своей реализацией в srctrace.go.
package srctrace

import (
	"fmt"
	"regexp"
	"strings"
)

// RepoInfo — сведения о репозитории на момент сборки.
type RepoInfo struct {
	Branch string // имя ветки, например "master"
	Short  string // короткий id коммита, например "cefa722"
	Commit string // полный id коммита (40 hex-символов)
	Date   string // дата/время сборки, например RFC3339
}

// Snippet генерирует валидный Go-файл с константами трассировки.
func Snippet(info RepoInfo) string {
	var b strings.Builder
	b.WriteString("// Сгенерировано srctrace — не редактируйте вручную.\n")
	b.WriteString("package revisions\n\n")
	b.WriteString("const (\n")
	fmt.Fprintf(&b, "\tBranch      = %q\n", info.Branch)
	fmt.Fprintf(&b, "\tShortCommit = %q\n", info.Short)
	fmt.Fprintf(&b, "\tCommit      = %q\n", info.Commit)
	fmt.Fprintf(&b, "\tBuildDate   = %q\n", info.Date)
	b.WriteString(")\n")
	return b.String()
}

// constRe ловит строки вида `Branch = "master"` (начало строки, имя,
// знак =, закавыченное значение). Комментарии не совпадают.
var constRe = regexp.MustCompile(`(?m)^\s*(Branch|ShortCommit|Commit|BuildDate)\s*=\s*"([^"\n]*)"`)

// Parse разбирает файл, сгенерированный Snippet.
func Parse(src string) (RepoInfo, error) {
	if strings.TrimSpace(src) == "" {
		return RepoInfo{}, fmt.Errorf("пустой источник")
	}
	hasPkg := false
	for _, line := range strings.Split(src, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "package ") {
			hasPkg = true
			break
		}
	}
	if !hasPkg {
		return RepoInfo{}, fmt.Errorf("нет строки package")
	}
	var info RepoInfo
	found := map[string]bool{}
	for _, m := range constRe.FindAllStringSubmatch(src, -1) {
		found[m[1]] = true
		switch m[1] {
		case "Branch":
			info.Branch = m[2]
		case "ShortCommit":
			info.Short = m[2]
		case "Commit":
			info.Commit = m[2]
		case "BuildDate":
			info.Date = m[2]
		}
	}
	for _, name := range []string{"Branch", "ShortCommit", "Commit", "BuildDate"} {
		if !found[name] {
			return RepoInfo{}, fmt.Errorf("нет константы %s", name)
		}
	}
	return info, nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/28-srctrace`, упражнение **srctrace**. Локально: `go test -tags solution ./28-srctrace/`._
