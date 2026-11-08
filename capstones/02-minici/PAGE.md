# Капстоун 2. MiniCI — сборки с трассировкой

> Собирает [`28-srctrace`](../../28-srctrace/README.md) и [`30-buildsystem`](../../30-buildsystem/README.md)

## Цель капстоуна

Интеграция: реестр сборок из главы 30 регистрирует билды, а глава 28 превращает
сведения о билде в Go-сниппет, встраиваемый в приложение. Клей `CI` делает
то, что делает настоящий CI: сборка → лог → трассировка, с retention'ом
всех хвостов сразу.

## Теория

- Пакет главы 30 называется `jobs` (импорт `"assignments-go/30-buildsystem"`),
  пакет главы 28 — `srctrace`.
- `jobs.NewBuild` пишет лог `logDir/<id>.log`; трассировку кладём рядом:
  `logDir/trace-<id>.go`.
- Retention — общая забота: `CI.Retain` должен убирать **и** логи реестра,
  **и** файлы трассировки выпавших сборок.

## Условие

1. `New(logDir)` — CI с реестром сборок.
2. `Build(branch, commit)` — зарегистрировать сборку (ID с 1, дата UTC),
   сгенерировать `srctrace.RepoInfo{Branch, Short (первые 7 символов),
   Commit, Date}` и записать сниппет в `logDir/trace-<ID>.go`.
   Возвращает сборку.
3. `All()` — все сборки по возрастанию ID.
4. `Trace(b, lang)` — traceability-текст сборки (`"go"`, `"ini"`;
   другой язык — ошибка).
5. `Retain(n)` — оставить последние n сборок: логи и файлы трассировки
   выпавших сборок удалить. n < 1 — ошибка.

## API для реализации

```go
type CI struct{ /* TODO */ }
func New(logDir string) *CI
func (c *CI) Build(branch, commit string) (jobs.Build, error)
func (c *CI) All() []jobs.Build
func (c *CI) Trace(b jobs.Build, lang string) (string, error)
func (c *CI) Retain(n int) error
```

## Подсказки

- Проверь сниппет в тесте по-настоящему: `srctrace.Parse` прочитанного
  `trace-<id>.go` должен вернуть те же Branch/Commit/Date.
- `jobs.Retain` не знает про `trace-*.go` — удаляй их сам, до вызова.

## Проверка

```
go test ./capstones/02-minici/
go test -tags solution ./capstones/02-minici/
```

## Эталон

`go test -tags solution ./capstones/02-minici/`

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./capstones/02-minici/
```

## Заготовки (реализуйте)

### `minici.go`

{% raw %}
```go
//go:build !solution

// Капстоун 2: мини-CI. Реализуй клей поверх 28-srctrace и 30-buildsystem
// так, чтобы прошли тесты в minici_test.go. Условие — в README.md.
package minici

import (
	"assignments-go/28-srctrace"
	jobs "assignments-go/30-buildsystem"
)

// var _ — держит импорт srctrace: понадобится в Build (см. решение).
var _ = srctrace.Snippet

// CI — реестр сборок (jobs) + трассировка каждой сборки (srctrace),
// встраиваемая в приложение. Все файлы лежат в logDir.
type CI struct {
	// TODO: твои поля.
}

// New создаёт CI с логами и трассировкой в logDir.
func New(logDir string) *CI {
	panic("TODO: New")
}

// Build регистрирует сборку (ID с 1) и пишет её трассировку
// в logDir/trace-<ID>.go (Short — первые 7 символов commit).
func (c *CI) Build(branch, commit string) (jobs.Build, error) {
	panic("TODO: Build")
}

// All возвращает все сборки по возрастанию ID.
func (c *CI) All() []jobs.Build {
	panic("TODO: All")
}

// Trace возвращает traceability-текст сборки b: lang "go" или "ini";
// другой lang — ошибка.
func (c *CI) Trace(b jobs.Build, lang string) (string, error) {
	panic("TODO: Trace")
}

// Retain оставляет последние n сборок: логи и trace-файлы выпавших
// сборок удаляются. n < 1 — ошибка.
func (c *CI) Retain(n int) error {
	panic("TODO: Retain")
}
```
{% endraw %}

## Тесты

### `minici_test.go`

{% raw %}
```go
package minici

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"assignments-go/28-srctrace"
)

const commit = "cefa7227e2c04b2b4c1a4a1b4c1a4a1b4c1a4a1"

func TestBuildWritesTrace(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	b, err := c.Build("master", commit)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if b.ID != 1 {
		t.Errorf("ID первой сборки %d, хотел 1", b.ID)
	}
	src, err := os.ReadFile(filepath.Join(dir, "trace-1.go"))
	if err != nil {
		t.Fatalf("trace-1.go: %v", err)
	}
	// Настоящая интеграция: сниппет главы 28 должен читаться обратно.
	info, err := srctrace.Parse(string(src))
	if err != nil {
		t.Fatalf("srctrace.Parse(trace-1.go): %v", err)
	}
	if info.Branch != "master" || info.Commit != commit || info.Short != commit[:7] || info.Date != b.Date {
		t.Errorf("Parse(trace-1.go) = %+v, хотел branch=master commit=%s short=%s date=%s",
			info, commit, commit[:7], b.Date)
	}
}

func TestIDsAndAll(t *testing.T) {
	c := New(t.TempDir())
	c.Build("master", commit)
	c.Build("feature/x", "0000000000000000000000000000000000000001")
	all := c.All()
	if len(all) != 2 || all[0].ID != 1 || all[1].ID != 2 {
		t.Fatalf("All: %+v, хотел две сборки с ID 1,2", all)
	}
}

func TestTrace(t *testing.T) {
	c := New(t.TempDir())
	b, _ := c.Build("master", commit)
	got, err := c.Trace(b, "go")
	if err != nil {
		t.Fatalf("Trace(go): %v", err)
	}
	for _, want := range []string{"master", commit} {
		if !strings.Contains(got, want) {
			t.Errorf("Trace(go): нет %q в:\n%s", want, got)
		}
	}
	if got, err := c.Trace(b, "rust"); err == nil {
		t.Errorf("Trace(rust): хотел ошибку, получил:\n%s", got)
	}
}

func TestRetainRemovesEverything(t *testing.T) {
	dir := t.TempDir()
	c := New(dir)
	c.Build("master", commit)
	c.Build("dev", commit)
	c.Build("feature/x", commit)
	if err := c.Retain(1); err != nil {
		t.Fatalf("Retain: %v", err)
	}
	if all := c.All(); len(all) != 1 || all[0].ID != 3 {
		t.Errorf("All после Retain(1): %+v, хотел только сборку 3", all)
	}
	for _, id := range []int{1, 2} {
		if _, err := os.Stat(filepath.Join(dir, fmt.Sprintf("trace-%d.go", id))); !os.IsNotExist(err) {
			t.Errorf("trace-%d.go не удалён (err=%v)", id, err)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "trace-3.go")); err != nil {
		t.Errorf("trace-3.go должен остаться: %v", err)
	}
	if err := c.Retain(0); err == nil {
		t.Error("Retain(0): хотел ошибку")
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

// Эталон капстоуна 2.
package minici

import (
	"fmt"
	"os"
	"path/filepath"

	"assignments-go/28-srctrace"
	jobs "assignments-go/30-buildsystem"
)

type CI struct {
	reg    *jobs.Registry
	logDir string
}

func New(logDir string) *CI {
	return &CI{reg: jobs.NewRegistry(logDir), logDir: logDir}
}

func (c *CI) Build(branch, commit string) (jobs.Build, error) {
	b := c.reg.NewBuild(branch, commit)
	short := commit
	if len(short) > 7 {
		short = short[:7]
	}
	info := srctrace.RepoInfo{Branch: branch, Short: short, Commit: commit, Date: b.Date}
	src := srctrace.Snippet(info)
	path := filepath.Join(c.logDir, fmt.Sprintf("trace-%d.go", b.ID))
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		return b, err
	}
	return b, nil
}

func (c *CI) All() []jobs.Build {
	return c.reg.All()
}

func (c *CI) Trace(b jobs.Build, lang string) (string, error) {
	return c.reg.Traceability(b, lang)
}

func (c *CI) Retain(n int) error {
	// Сначала узнаём, кто выпадет, затем убираем trace-файлы и логи реестра.
	var old []int
	all := c.reg.All()
	if len(all) > n {
		for _, b := range all[:len(all)-n] {
			old = append(old, b.ID)
		}
	}
	if err := c.reg.Retain(n); err != nil {
		return err
	}
	for _, id := range old {
		if err := os.Remove(filepath.Join(c.logDir, fmt.Sprintf("trace-%d.go", id))); err != nil {
			return err
		}
	}
	return nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/capstones/02-minici`, упражнение **minici**. Локально: `go test -tags solution ./capstones/02-minici/`._
