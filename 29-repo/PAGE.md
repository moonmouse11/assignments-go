# Глава 29. Repo — группы репозиториев

> Большой проект — это много репозиториев: конфигурация группы в
> мини-YAML и подготовка рабочей области (workarea) по списку путей проектов.

## Цель главы

Написать парсер подмножества YAML руками — без внешних либ — и научиться
готовить структуру каталогов из конфигурации.

## Теория

- YAML держится на отступах: top-level ключи — без отступа, элементы
  списка начинаются с `-`, поля элемента — с отступом под ним.
- Значения в кавычках (`"$HOME/x"`) — просто строки: кавычки срезаем,
  переменные окружения не разворачиваем.
- Полный путь проекта — `filepath.Join(workarea, path)`; уникальность
  `path` проверяется до создания каталогов (map «path → проект»).
- `os.MkdirAll` создаёт цепочку каталогов целиком — baseDir отдельно
  создавать не нужно.

## Условие

Конфигурация в духе `AllProjects.yaml` (упрощённо — одна группа):

```yaml
workarea: "$HOME/grepotest/pub"
projects:
  - repourl: "git@gitlab.com:projtemplates/go.git"
    reference: "master"
    path: "go/lib"
```

1. `ParseConfig` разбирает `workarea` и список проектов (поля `repourl`,
   `reference`, `path`); пустые строки и комментарии `#` пропускаются.
2. Некорректная конфигурация — ошибка: строка без двоеточия, неизвестное
   поле, нет `workarea` или нет `projects:`.
3. `Prepare` создаёт baseDir и каталоги `Path` всех проектов, возвращает
   список полных путей в порядке объявления; дубликат `Path` — ошибка.

Wants (необязательно): разворачивание `$HOME`, группы public/private,
поля `repo`/`server`/`build`.

## API для реализации

```go
type Project struct{ RepoURL, Reference, Path string }
type Config struct { Workarea string; Projects []Project }
func ParseConfig(r io.Reader) (Config, error)
func Prepare(cfg Config, baseDir string) ([]string, error)
```

Сигнатуры и договоры — в `repo.go`. Тесты — в `repo_test.go`.

## Подсказки

- Разбор по строкам: `strings.TrimSpace` + проверка префикса `- `.
- Пара «ключ: значение» — `strings.Index(s, ":")`; кавычки `"` и `'` — срезать.
- Указатель `cur *Project` на последний элемент списка упрощает дополнение
  полей следующих строк.

## Проверка

```
go vet ./29-repo/ && go test ./29-repo/            # красный: TODO
go vet -tags solution ./29-repo/ && go test -tags solution ./29-repo/
```

## Эталон

`solution.go` собирается только с тегом `solution`.

## Углубление

- Групп две — public/private: public — только reference (клон с detached
  HEAD, без push и tag); `server` группы + `repo` проекта склеиваются в
  полный URL; `reference` — дефолт группы (ветка, тег или commit id),
  переопределяется в проекте; `build` — shell-команда после clone/pull.
- `workarea` переводится через переменные окружения (`os.Getenv`);
  воркспейсы не должны лежать внутри чужого git-дерева.
- Подкоманды: `init` (каталоги + `git clone` + `git checkout reference` +
  build; `init <branch>` создаёт ветку во всех private-репо), `pull`,
  `push` (commit с общим сообщением по группе), `status`, `diff`, `tag`;
  флаги `--config` (дефолт `Project.yaml`), `--verbose` — каркас на `flag`.
- Клонирование, pull/push и build-команды — `git` через `os/exec` и сеть:
  внешние зависимости; в главе — только конфиг и каталоги, stdlib.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./29-repo/
```

## Заготовки (реализуйте)

### `repo.go`

{% raw %}
```go
//go:build !solution

// Пакет repo — глава 29. Реализуй функции так, чтобы прошли тесты
// в repo_test.go. Условие — в README.md этой папки.
package repo

import "io"

// Project — один репозиторий группы: откуда клонировать (RepoURL),
// какую версию брать (Reference: ветка/тег/коммит) и куда положить
// (Path — путь относительно workarea, уникальный в группе).
type Project struct {
	RepoURL   string
	Reference string
	Path      string
}

// Config — разбор мини-YAML: workarea и список проектов.
type Config struct {
	Workarea string
	Projects []Project
}

// ParseConfig разбирает подмножество YAML (см. README): top-level ключ
// workarea, затем projects: с элементами "- repourl: ...",
// "reference: ...", "path: ..." по отступам. Значения в кавычках —
// строки (кавычки срезаются). Некорректная конфигурация — ошибка.
func ParseConfig(r io.Reader) (Config, error) {
	panic("TODO: ParseConfig")
}

// Prepare создаёт baseDir и каталоги Path всех проектов (под baseDir),
// возвращает список полных путей в порядке объявления. Дубликат Path
// или пустой Path — ошибка.
func Prepare(cfg Config, baseDir string) ([]string, error) {
	panic("TODO: Prepare")
}
```
{% endraw %}

## Тесты

### `repo_test.go`

{% raw %}
```go
package repo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sample = `
# Пример в духе AllProjects.yaml.
workarea: "$HOME/grepotest/pub"
projects:
  - repourl: "git@gitlab.com:projtemplates/go.git"
    reference: "master"
    path: "go/lib"
  - repourl: "https://github.com/RajaSrinivasan/disign.git"
    path: "disign"
`

func TestParseConfig(t *testing.T) {
	cfg, err := ParseConfig(strings.NewReader(sample))
	if err != nil {
		t.Fatalf("ParseConfig: %v", err)
	}
	if cfg.Workarea != "$HOME/grepotest/pub" {
		t.Errorf("Workarea = %q, хотел %q (строка как есть)", cfg.Workarea, "$HOME/grepotest/pub")
	}
	if len(cfg.Projects) != 2 {
		t.Fatalf("Projects: %d шт., хотел 2", len(cfg.Projects))
	}
	first := cfg.Projects[0]
	if first.RepoURL != "git@gitlab.com:projtemplates/go.git" ||
		first.Reference != "master" || first.Path != "go/lib" {
		t.Errorf("первый проект = %+v", first)
	}
	if cfg.Projects[1].RepoURL != "https://github.com/RajaSrinivasan/disign.git" ||
		cfg.Projects[1].Reference != "" || cfg.Projects[1].Path != "disign" {
		t.Errorf("второй проект = %+v", cfg.Projects[1])
	}
}

func TestParseConfigErrors(t *testing.T) {
	bad := []string{
		"workarea: /tmp\nprojects:\n  - repourl: x.git\n    oops\n", // строка без двоеточия
		"projects:\n  - repourl: \"x.git\"\n    path: p\n",          // нет workarea
		"workarea: /tmp\n", // нет projects
		"garbage\n",        // нет двоеточия
		"workarea: /tmp\nprojects:\n  - repourl: x.git\n    path: p\n    build: make\n", // неизвестное поле
	}
	for _, src := range bad {
		if _, err := ParseConfig(strings.NewReader(src)); err == nil {
			t.Errorf("ParseConfig(%q): хотел ошибку, получил nil", src)
		}
	}
}

func TestPrepare(t *testing.T) {
	cfg, err := ParseConfig(strings.NewReader(sample))
	if err != nil {
		t.Fatalf("ParseConfig: %v", err)
	}
	base := filepath.Join(t.TempDir(), "base")
	paths, err := Prepare(cfg, base)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("Prepare вернул %d путей, хотел 2: %v", len(paths), paths)
	}
	for i, p := range paths {
		if want := filepath.Join(base, cfg.Projects[i].Path); p != want {
			t.Errorf("путь #%d = %q, хотел %q", i, p, want)
		}
		if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
			t.Errorf("каталог %s не создан (err=%v)", p, err)
		}
	}
}

func TestPrepareDuplicatePath(t *testing.T) {
	cfg := Config{
		Workarea: "/tmp/wa",
		Projects: []Project{
			{RepoURL: "a.git", Path: "dup"},
			{RepoURL: "b.git", Path: "dup"},
		},
	}
	if _, err := Prepare(cfg, filepath.Join(t.TempDir(), "base")); err == nil {
		t.Fatal("Prepare: хотел ошибку для дубликата path, получил nil")
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

// Эталонное решение главы 29. Сравни со своей реализацией в repo.go.
package repo

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Project — один репозиторий группы.
type Project struct {
	RepoURL   string
	Reference string
	Path      string
}

// Config — разбор мини-YAML: workarea и список проектов.
type Config struct {
	Workarea string
	Projects []Project
}

// ParseConfig разбирает мини-YAML по строкам: до projects: — top-level
// ключи, после — список проектов по отступам.
func ParseConfig(r io.Reader) (Config, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Config{}, err
	}
	var (
		cfg        Config
		inProjects bool
		cur        *Project
	)
	for _, raw := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if !inProjects {
			key, val, err := splitKV(trimmed)
			if err != nil {
				return Config{}, err
			}
			switch key {
			case "workarea":
				if cfg.Workarea != "" {
					return Config{}, fmt.Errorf("повторный ключ workarea")
				}
				cfg.Workarea = val
			case "projects":
				if val != "" {
					return Config{}, fmt.Errorf("projects: не ждёт значения: %q", val)
				}
				inProjects = true
			default:
				return Config{}, fmt.Errorf("неизвестный ключ %q", key)
			}
			continue
		}
		if strings.HasPrefix(trimmed, "- ") {
			key, val, err := splitKV(strings.TrimSpace(trimmed[2:]))
			if err != nil {
				return Config{}, err
			}
			cfg.Projects = append(cfg.Projects, Project{})
			cur = &cfg.Projects[len(cfg.Projects)-1]
			if err := setField(cur, key, val); err != nil {
				return Config{}, err
			}
			continue
		}
		if cur == nil {
			return Config{}, fmt.Errorf("строка %q до первого элемента списка", trimmed)
		}
		key, val, err := splitKV(trimmed)
		if err != nil {
			return Config{}, err
		}
		if err := setField(cur, key, val); err != nil {
			return Config{}, err
		}
	}
	if !inProjects {
		return Config{}, fmt.Errorf("нет ключа projects:")
	}
	if cfg.Workarea == "" {
		return Config{}, fmt.Errorf("нет ключа workarea:")
	}
	return cfg, nil
}

// splitKV режет "ключ: значение" по первому двоеточию.
func splitKV(s string) (key, val string, err error) {
	i := strings.Index(s, ":")
	if i < 0 {
		return "", "", fmt.Errorf("нет двоеточия: %q", s)
	}
	return strings.TrimSpace(s[:i]), unquote(strings.TrimSpace(s[i+1:])), nil
}

// unquote срезает окружающие кавычки (двойные или одинарные).
func unquote(s string) string {
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}
	return s
}

func setField(p *Project, key, val string) error {
	switch key {
	case "repourl":
		p.RepoURL = val
	case "reference":
		p.Reference = val
	case "path":
		p.Path = val
	default:
		return fmt.Errorf("неизвестное поле проекта %q", key)
	}
	return nil
}

// Prepare проверяет пути на дубликаты, затем создаёт baseDir и каталоги
// всех проектов.
func Prepare(cfg Config, baseDir string) ([]string, error) {
	seen := map[string]string{}
	paths := make([]string, 0, len(cfg.Projects))
	for _, p := range cfg.Projects {
		if p.Path == "" {
			return nil, fmt.Errorf("пустой path у проекта %q", p.RepoURL)
		}
		if prev, dup := seen[p.Path]; dup {
			return nil, fmt.Errorf("дубликат path %q (у %q и %q)", p.Path, prev, p.RepoURL)
		}
		seen[p.Path] = p.RepoURL
		paths = append(paths, filepath.Join(baseDir, p.Path))
	}
	for _, p := range paths {
		if err := os.MkdirAll(p, 0o755); err != nil {
			return nil, err
		}
	}
	return paths, nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/29-repo`, упражнение **repo**. Локально: `go test -tags solution ./29-repo/`._
