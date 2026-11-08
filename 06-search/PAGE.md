# Глава 6. Search — поиск по файлам

> Projectlet про grep-подобный инструмент:
> текст и регулярные выражения, файлы и директории, рекурсия.

## Цель главы

Написать ядро утилиты поиска: найти кандидата в файлах — по подстроке
или по регулярному выражению, без учёта регистра по умолчанию, с номерами
строк. Движок отделён от CLI — завтра к нему можно прикрутить GUI.

## Теория

- `regexp.Compile`; префикс `(?i)` — поиск без учёта регистра;
  `re.FindAllStringIndex(s, -1)` — все совпадения, их длина — их число.
- `regexp.QuoteMeta` экранирует метасимволы: подстрока превращается
  в литеральный паттерн — один движок на оба режима.
- `filepath.WalkDir` — обход дерева; `fs.DirEntry.IsDir()` отличает
  файлы от директорий; `fs.SkipDir` не пускает обход в поддиректорию.
- `filepath.Glob(pattern)` — раскрытие wildcard-имён вроде `*.txt`.

## Условие

1. `FindLine(line, needle, o)` — число вхождений needle в line и есть ли
   хоть одно. `o.Pattern` — needle это регулярное выражение, иначе
   подстрока. По умолчанию регистр не учитывается.
2. `SearchFile(path, needle, o)` — совпадения в файле: `[]Hit` (номер
   строки с 1, текст строки, число вхождений). Нет совпадений — пустой
   результат без ошибки. Нечитаемый файл — ошибка; некорректный
   Pattern — совпадений нет.
3. `ListFiles(paths, recurse)` — разворот имён в список файлов:
   директория → её файлы (с `recurse` — всё дерево через
   `filepath.WalkDir`); имя-файл → само имя; wildcard-имя → совпадения
   `filepath.Glob`. Несуществующее имя без совпадений — ошибка.

Wants (необязательно, как в `search.pdf`): подсветить совпавший текст.

## API для реализации

```go
type Options struct{ Pattern, CaseSensitive bool }
type Hit struct{ Line, Count int; Text string }
func FindLine(line, needle string, o Options) (int, bool)
func SearchFile(path, needle string, o Options) ([]Hit, error)
func ListFiles(paths []string, recurse bool) ([]string, error)
```

Сигнатуры и договоры — в `search.go`. Тесты — в `search_test.go`.

## Подсказки

- Пустой needle — нет вхождений: `(0, false)`, иначе regexp найдёт
  бесконечную серию пустых совпадений.
- `bufio.Scanner` читает файл построчно; счётчик строк веди сам.
- Порядок результата `ListFiles` — порядок paths, внутри — лексикографический.

## Проверка

```
go test ./06-search/
go test -tags solution ./06-search/
```

## Эталон

`solution.go` собирается только с тегом `solution`:
`go test -tags solution ./06-search/`.

## Углубление

- CLI-обвязка спеки: переключатели `-p/--pattern` (кандидат — регулярное
  выражение, не текст), `-r/--recurse`, `-c/--case` (по умолчанию регистр
  не учитывается); кандидат — arg0, имена файлов — всегда wildcard.
  Собирается на пакете `flag`.
- Формат вывода примера: перед совпадениями файла — заголовок
  `------ File ./Impl.cs ------`, строки — `NNNN: текст`; имя файла
  печатается только если кандидат найден (need 6 спеки).
- Подсветка совпавшего текста (wants 1): `re.FindAllStringIndex` даёт
  позиции всех вхождений в строке — обрами их маркерами или ANSI.
- Спека мотивирует отделение движка от CLI ради будущего GUI; ядро главы
  уже такое — GUI потребует внешней библиотеки (например, Fyne).

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./06-search/
```

## Заготовки (реализуйте)

### `search.go`

{% raw %}
```go
//go:build !solution

// Пакет search — глава 6. Реализуй функции так, чтобы прошли тесты
// в search_test.go. Условие — в README.md этой папки.
package search

// Options задаёт режим поиска: Pattern — needle это регулярное
// выражение (иначе подстрока); CaseSensitive — учитывать регистр
// (по умолчанию поиск без учёта регистра).
type Options struct {
	Pattern       bool
	CaseSensitive bool
}

// Hit — одно совпавшее место файла: Line — номер строки (с 1),
// Text — текст строки, Count — число вхождений кандидата в строку.
type Hit struct {
	Line  int
	Text  string
	Count int
}

// FindLine возвращает число вхождений needle в line и найдено ли
// хоть одно. Pattern → regexp, иначе подстрока; регистронезависимо,
// если не CaseSensitive. Пустой needle — (0, false).
func FindLine(line, needle string, o Options) (int, bool) {
	panic("TODO: FindLine")
}

// SearchFile ищет needle в файле path и возвращает совпадения
// по строкам в порядке возрастания. Нет совпадений — пустой
// результат и nil. Отсутствующий или нечитаемый файл — ошибка;
// некорректный Pattern — совпадений нет.
func SearchFile(path, needle string, o Options) ([]Hit, error) {
	panic("TODO: SearchFile")
}

// ListFiles разворачивает paths в список файлов. Директория —
// её файлы (с recurse — всё дерево через filepath.WalkDir, без —
// только верхний уровень); обычное имя файла — само имя; имя
// с wildcard-метасимволами — совпадения filepath.Glob. Путь,
// которого нет и который ни с чем не совпал, — ошибка.
func ListFiles(paths []string, recurse bool) ([]string, error) {
	panic("TODO: ListFiles")
}
```
{% endraw %}

## Тесты

### `search_test.go`

{% raw %}
```go
package search

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFindLine(t *testing.T) {
	cases := []struct {
		line, needle string
		o            Options
		count        int
		found        bool
	}{
		{"Go go GO", "go", Options{}, 3, true},
		{"Go go GO", "go", Options{CaseSensitive: true}, 1, true},
		{"hello world", "xyz", Options{}, 0, false},
		{"abc", "", Options{}, 0, false},
		{`a1b22c333`, `\d+`, Options{Pattern: true}, 3, true},
		{"ABC abc", `abc`, Options{Pattern: true}, 2, true},
		{"ABC abc", `abc`, Options{Pattern: true, CaseSensitive: true}, 1, true},
	}
	for _, tc := range cases {
		count, found := FindLine(tc.line, tc.needle, tc.o)
		if count != tc.count || found != tc.found {
			t.Errorf("FindLine(%q, %q, %+v) = (%d, %v), хотел (%d, %v)",
				tc.line, tc.needle, tc.o, count, found, tc.count, tc.found)
		}
	}
}

func TestSearchFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "note.txt")
	os.WriteFile(path, []byte("Alpha beta\nnothing here\nbeta Beta beta\n"), 0o644)

	hits, err := SearchFile(path, "beta", Options{})
	if err != nil {
		t.Fatalf("SearchFile: %v", err)
	}
	want := []Hit{
		{Line: 1, Text: "Alpha beta", Count: 1},
		{Line: 3, Text: "beta Beta beta", Count: 3},
	}
	if !reflect.DeepEqual(hits, want) {
		t.Errorf("SearchFile без регистра:\nполучил %+v\nхотел  %+v", hits, want)
	}

	hits, err = SearchFile(path, "BETA", Options{CaseSensitive: true})
	if err != nil {
		t.Fatalf("SearchFile: %v", err)
	}
	if len(hits) != 0 {
		t.Errorf("SearchFile с регистром: нашёл %+v, хотел пусто", hits)
	}
}

func TestSearchFilePattern(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	os.WriteFile(path, []byte("id=42 ok\nid=bad\nid=43 ok\n"), 0o644)

	hits, err := SearchFile(path, `id=\d+`, Options{Pattern: true})
	if err != nil {
		t.Fatalf("SearchFile: %v", err)
	}
	want := []Hit{
		{Line: 1, Text: "id=42 ok", Count: 1},
		{Line: 3, Text: "id=43 ok", Count: 1},
	}
	if !reflect.DeepEqual(hits, want) {
		t.Errorf("SearchFile паттерн:\nполучил %+v\nхотел  %+v", hits, want)
	}
}

func TestSearchFileMissing(t *testing.T) {
	if _, err := SearchFile(filepath.Join(t.TempDir(), "nope.txt"), "x", Options{}); err == nil {
		t.Fatal("SearchFile: хотел ошибку для отсутствующего файла, получил nil")
	}
}

func TestListFiles(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "sub"), 0o755)
	os.Mkdir(filepath.Join(dir, "sub", "deep"), 0o755)
	for _, name := range []string{
		"a.txt", "b.md",
		filepath.Join("sub", "c.txt"),
		filepath.Join("sub", "deep", "d.txt"),
	} {
		os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644)
	}

	// Директория без рекурсии — только файлы верхнего уровня.
	files, err := ListFiles([]string{dir}, false)
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("ListFiles без рекурсии: получил %v, хотел 2 файла верхнего уровня", files)
	}

	// С рекурсией — всё дерево, включая файл в поддиректории.
	files, err = ListFiles([]string{dir}, true)
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}
	if len(files) != 4 {
		t.Errorf("ListFiles рекурсия: получил %v, хотел 4 файла", files)
	}
	foundDeep := false
	for _, f := range files {
		if filepath.Base(f) == "d.txt" {
			foundDeep = true
		}
	}
	if !foundDeep {
		t.Error("ListFiles рекурсия: не нашёл файл в поддиректории")
	}

	// Одиночное имя файла — само себя.
	files, err = ListFiles([]string{filepath.Join(dir, "a.txt")}, false)
	if err != nil || len(files) != 1 || filepath.Base(files[0]) != "a.txt" {
		t.Errorf("ListFiles файл: получил %v, err %v", files, err)
	}

	// Wildcard-имя — раскрытие через glob.
	files, err = ListFiles([]string{filepath.Join(dir, "*.txt")}, false)
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}
	if len(files) != 1 || filepath.Base(files[0]) != "a.txt" {
		t.Errorf("ListFiles *.txt: получил %v, хотел a.txt", files)
	}

	// Несуществующее имя без совпадений — ошибка.
	if _, err := ListFiles([]string{filepath.Join(dir, "no-such-*")}, false); err == nil {
		t.Error("ListFiles: хотел ошибку для несуществующего имени, получил nil")
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

// Эталонное решение главы 6. Сравни со своей реализацией в search.go.
package search

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type Options struct{ Pattern, CaseSensitive bool }

type Hit struct {
	Line, Count int
	Text        string
}

func FindLine(line, needle string, o Options) (int, bool) {
	if needle == "" {
		return 0, false
	}
	if !o.Pattern {
		needle = regexp.QuoteMeta(needle) // подстрока — литеральный паттерн
	}
	if !o.CaseSensitive {
		needle = "(?i)" + needle
	}
	if re, err := regexp.Compile(needle); err == nil {
		n := len(re.FindAllStringIndex(line, -1))
		return n, n > 0
	}
	return 0, false
}

func SearchFile(path, needle string, o Options) ([]Hit, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	var hits []Hit
	for n := 0; sc.Scan(); n++ {
		if c, ok := FindLine(sc.Text(), needle, o); ok {
			hits = append(hits, Hit{Line: n + 1, Text: sc.Text(), Count: c})
		}
	}
	return hits, sc.Err()
}

func ListFiles(paths []string, recurse bool) ([]string, error) {
	var files []string
	for _, p := range paths {
		info, err := os.Stat(p)
		if err == nil && info.IsDir() {
			err = filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
				if err == nil && d.IsDir() && path != p && !recurse {
					return fs.SkipDir
				}
				if err == nil && !d.IsDir() {
					files = append(files, path)
				}
				return err
			})
			if err != nil {
				return nil, err
			}
		} else if err == nil {
			files = append(files, p)
		} else if m, _ := filepath.Glob(p); len(m) == 0 {
			return nil, err // os.Stat: пути нет и glob пуст
		} else {
			files = append(files, m...)
		}
	}
	return files, nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/06-search`, упражнение **search**. Локально: `go test -tags solution ./06-search/`._
