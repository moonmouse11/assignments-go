# Глава 31. Project — клонирование из шаблона

> Новый проект из шаблона: копируем дерево
> каталогов и подставляем значения вместо плейсхолдеров `{{NAME}}`.

## Цель главы

Освоить обход дерева каталогов (`filepath.WalkDir`), относительные пути
(`filepath.Rel`) и замену плейсхолдеров при копировании файлов.

## Теория

- `filepath.WalkDir(root, fn)` обходит дерево; fn получает путь и
  `fs.DirEntry` — каталоги и файлы различаются по `d.IsDir()`.
- Относительный путь: `filepath.Rel(tmplDir, p)` → тот же путь в dstDir
  через `filepath.Join(dstDir, rel)`.
- Замена: `strings.ReplaceAll(content, "{{"+name+"}}", val)` по всем
  ключам vars — плейсхолдеры разные, порядок обхода map не важен.
- `{{X}}` — простой маркер, а не выражения `text/template`: замена
  буквальная, никаких циклов и условий.

## Условие

1. `Instantiate` копирует дерево tmplDir (включая поддиректории) в dstDir.
2. В каждом файле вхождения `{{NAME}}` заменяются значением
   `vars["NAME"]` для всех ключей vars.
3. Файлы без плейсхолдеров копируются без изменений.
4. Отсутствующий tmplDir — ошибка.

Wants (необязательно): переименование файлов, meta-макросы
`$(Created)`/`$(Username)`, шаблон из git-репозитория.

## API для реализации

```go
func Instantiate(tmplDir, dstDir string, vars map[string]string) error
```

Сигнатуры и договоры — в `project.go`. Тесты — в `project_test.go`.

## Подсказки

- Сначала `os.Stat(tmplDir)` — понятная ошибка вместо пустого обхода.
- `os.ReadFile`/`os.WriteFile` достаточно: шаблоны — текст.
- Каталоги создавай `os.MkdirAll` — dstDir и вложенные пути сразу.

## Проверка

```
go vet ./31-project/ && go test ./31-project/            # красный: TODO
go vet -tags solution ./31-project/ && go test -tags solution ./31-project/
```

## Эталон

`solution.go` собирается только с тегом `solution`.

## Углубление

- `Clone.cfg` в корне шаблона — INI-секции-макросы с полями `default` и
  `help`; значения переопределяются из командной строки. Разбор —
  построчный, stdlib (`strings`, `bufio`).
- Директива ищется в строке по ключевому слову `proj:` (регистр важен),
  часть строки после вхождения не обрабатывается; директивы живут в
  комментариях языка — инструмент language-agnostic.
- `proj: rename $(Proj)` / `$(Proj)_option` — подстановка в basename
  (`filepath.Base`); `proj: subs word $(Macro)` — замена слова значением
  макроса, ненастроенный макрос уничтожает вхождение.
- Meta-макросы `$(Created)`/`$(Username)` нельзя переопределить из CLI.
- Усиления из спеки: директива `if` (как C-препроцессор), клонирование
  частичных структур, автогенерация каркаса — interface/implementation,
  unit-тесты, Makefile. У автора — C++ (glib/boost) и клон из git:
  внешние зависимости, в главе только stdlib.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./31-project/
```

## Заготовки (реализуйте)

### `project.go`

{% raw %}
```go
//go:build !solution

// Пакет proj — глава 31. Реализуй функцию так, чтобы прошли тесты
// в project_test.go. Условие — в README.md этой папки.
package proj

// Instantiate копирует дерево шаблона tmplDir (включая поддиректории)
// в dstDir, заменяя в каждом файле вхождения плейсхолдеров:
// {{NAME}} → vars["NAME"] для каждого ключа vars (буквальная замена,
// не text/template). Файлы без плейсхолдеров копируются как есть.
// Отсутствующий tmplDir — ошибка.
func Instantiate(tmplDir, dstDir string, vars map[string]string) error {
	panic("TODO: Instantiate")
}
```
{% endraw %}

## Тесты

### `project_test.go`

{% raw %}
```go
package proj

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	tmplMain  = "// Проект {{PROJECT}}.\n// Автор: {{USER}}.\npackage main\n\nfunc main() {}\n"
	tmplSub   = "{{PROJECT}} / {{USER}}\n"
	tmplPlain = "никаких плейсхолдеров тут нет\n"
)

var vars = map[string]string{"PROJECT": "goutil", "USER": "pots"}

func setupTemplate(t *testing.T) string {
	t.Helper()
	tmpl := filepath.Join(t.TempDir(), "tmpl")
	if err := os.MkdirAll(filepath.Join(tmpl, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	for path, content := range map[string]string{
		filepath.Join(tmpl, "main.go"):         tmplMain,
		filepath.Join(tmpl, "sub", "note.txt"): tmplSub,
		filepath.Join(tmpl, "README.txt"):      tmplPlain,
	} {
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return tmpl
}

func TestInstantiate(t *testing.T) {
	tmpl := setupTemplate(t)
	dst := filepath.Join(t.TempDir(), "dst")
	if err := Instantiate(tmpl, dst, vars); err != nil {
		t.Fatalf("Instantiate: %v", err)
	}
	checks := []struct{ rel, want string }{
		{"main.go", "// Проект goutil.\n// Автор: pots.\npackage main\n\nfunc main() {}\n"},
		{filepath.Join("sub", "note.txt"), "goutil / pots\n"},
		{"README.txt", tmplPlain}, // без плейсхолдеров — не тронут
	}
	for _, c := range checks {
		got, err := os.ReadFile(filepath.Join(dst, c.rel))
		if err != nil {
			t.Errorf("файл %s не создан: %v", c.rel, err)
			continue
		}
		if string(got) != c.want {
			t.Errorf("файл %s:\nполучил %q\nхотел  %q", c.rel, got, c.want)
		}
	}
}

func TestInstantiateMissingTemplate(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	dst := filepath.Join(t.TempDir(), "dst")
	if err := Instantiate(missing, dst, vars); err == nil {
		t.Fatal("Instantiate: хотел ошибку для отсутствующего шаблона, получил nil")
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

// Эталонное решение главы 31. Сравни со своей реализацией в project.go.
package proj

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Instantiate(tmplDir, dstDir string, vars map[string]string) error {
	if _, err := os.Stat(tmplDir); err != nil {
		return fmt.Errorf("шаблон %s недоступен: %w", tmplDir, err)
	}
	return filepath.WalkDir(tmplDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(tmplDir, p)
		if err != nil {
			return err
		}
		dst := filepath.Join(dstDir, rel)
		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		content := string(data)
		for name, val := range vars {
			content = strings.ReplaceAll(content, "{{"+name+"}}", val)
		}
		return os.WriteFile(dst, []byte(content), 0o644)
	})
}
```
{% endraw %}


</details>

_Пакет `assignments-go/31-project`, упражнение **project**. Локально: `go test -tags solution ./31-project/`._
