# Глава 30. Buildsystem — мини-CI

> Ядро build-системы: реестр сборок со
> сквозной нумерацией, лог-файлы, retention и трассировка «сборка →
> исходник».

## Цель главы

Смоделировать жизненный цикл CI: каждая сборка получает уникальный ID,
пишет лог-файл, устаревшие сборки (и их логи) удаляются, а связь
с исходником встраивается в приложение.

## Теория

- ID — монотонный счётчик с 1; после Retain он не «сбрасывается»,
  поэтому храним `lastID`, а не считаем `len(builds)`.
- Дата — `time.Now().UTC().Format(time.RFC3339)`: единый формат,
  сортируемый как строка.
- Retention: оставить последние n — хвост слайса `builds[len-n:]`,
  лог-файлы удалённых сборок — `os.Remove`.
- Traceability — те же данные в двух форматах: Go-константы (как в
  главе 28) и ini-текст с секцией `[jobs]`.

## Условие

1. `NewRegistry(logDir)` — реестр сборок; лог сборки с ID `<id>` —
   в файле `logDir/<id>.log`.
2. `NewBuild(branch, commit)` — следующая сборка: ID с 1, `Date` —
   текущее время UTC в RFC3339, лог-файл создаётся и наполняется
   заголовком; `Build.Log` — то же содержимое.
3. `All()` — все сборки в порядке возрастания ID.
4. `Retain(n)` — оставить последние n сборок: старые записи удаляются
   из реестра, их лог-файлы стираются; `n < 1` — ошибка.
5. `Traceability(b, lang)`: `"go"` — Go-константы Branch/Commit/
   BuildDate; `"ini"` — ini-текст; другой lang — ошибка.

Wants (необязательно): стадии init/build/publish/archive, скрипты
сборки, toml-конфиги, hostname build-машины.

## API для реализации

```go
type Build struct { ID int; Branch, Commit, Date string; Log []byte }
type Registry struct{ /* свои поля */ }
func NewRegistry(logDir string) *Registry
func (r *Registry) NewBuild(branch, commit string) Build
func (r *Registry) All() []Build
func (r *Registry) Retain(n int) error
func (r *Registry) Traceability(b Build, lang string) (string, error)
```

Сигнатуры и договоры — в `buildsystem.go`. Тесты — в
`buildsystem_test.go`.

## Подсказки

- `All` возвращает копию слайса: тесты не должны видеть внутренности.
- Порядок в Retain: сначала проверка `n`, потом удаления — `Retain(0)`
  ничего не трогает.
- Имена файлов — `filepath.Join(logDir, fmt.Sprintf("%d.log", id))`.

## Проверка

```
go vet ./30-buildsystem/ && go test ./30-buildsystem/            # красный: TODO
go vet -tags solution ./30-buildsystem/ && go test -tags solution ./30-buildsystem/
```

## Эталон

`solution.go` собирается только с тегом `solution`.

## Углубление

- Подкоманды: `build <job> <script> [branch]` (без branch — дефолт из
  init), `purge` (чистка по retention), `show` (полный лог),
  `enumerate --current` (ветки), `fetch --set-default` — каркас на `flag`.
- Флаги: `--jobspace` (дефолт `$JOBSPACE`, имя репо — имя джоба),
  `--init=scriptdir` (каталог скриптов в репо: shell/python3),
  `--trace=(C,C++,Ada,Python,Go,Text)` с дефолтом Ada, `--retention=3`,
  `--pull` (git pull перед сборкой).
- Сборки — каталоги `<jobspace>/<job>/<branch>/<id>` с `jobs.log`
  (таймстемп + `......> [I] сообщение`); скрипту через env передаются
  `BUILD_NUMBER`, `OUTPUT_DIR`, `SOURCE_DIR`, `REVISION_FILE`
  (`os/exec.Cmd.Env`).
- Ещё stdlib: hostname build-машины (`os.Hostname`) в трассировке, конфиги
  `jobs.toml`. Клонирование и запуск скриптов — `git` + сеть: внешние
  зависимости, в главе только реестр.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./30-buildsystem/
```

## Заготовки (реализуйте)

### `buildsystem.go`

{% raw %}
```go
//go:build !solution

// Пакет jobs — глава 30. Реализуй типы и методы так, чтобы прошли тесты
// в buildsystem_test.go. Условие — в README.md этой папки.
package jobs

// Build — одна сборка: идентификатор (с 1), ветка, коммит, дата
// (UTC, RFC3339) и начальное содержимое лога.
type Build struct {
	ID     int
	Branch string
	Commit string
	Date   string
	Log    []byte
}

// Registry — реестр сборок; лог-файлы лежат в logDir/<id>.log.
// Добавь свои поля.
type Registry struct{}

// NewRegistry создаёт реестр сборок с логами в logDir.
func NewRegistry(logDir string) *Registry {
	panic("TODO: NewRegistry")
}

// NewBuild регистрирует новую сборку: ID — следующий по счёту (с 1),
// Date — текущее время UTC в RFC3339, лог-файл logDir/<id>.log создаётся
// и наполняется заголовком (Build.Log — то же содержимое).
func (r *Registry) NewBuild(branch, commit string) Build {
	panic("TODO: NewBuild")
}

// All возвращает все сборки в порядке возрастания ID.
func (r *Registry) All() []Build {
	panic("TODO: All")
}

// Retain оставляет последние n сборок: старые записи удаляются из
// реестра, их лог-файлы стираются. n < 1 — ошибка.
func (r *Registry) Retain(n int) error {
	panic("TODO: Retain")
}

// Traceability возвращает сведения о сборке b для встраивания в
// приложение: lang "go" — Go-константы Branch/Commit/BuildDate,
// "ini" — ini-текст с секцией [jobs]. Другой lang — ошибка.
func (r *Registry) Traceability(b Build, lang string) (string, error) {
	panic("TODO: Traceability")
}
```
{% endraw %}

## Тесты

### `buildsystem_test.go`

{% raw %}
```go
package jobs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func newTestRegistry(t *testing.T) (*Registry, string) {
	t.Helper()
	logDir := filepath.Join(t.TempDir(), "logs")
	return NewRegistry(logDir), logDir
}

func TestNewBuild(t *testing.T) {
	r, logDir := newTestRegistry(t)
	cases := []struct{ branch, commit string }{
		{"master", "aaaa1111"},
		{"master", "bbbb2222"},
		{"feat/trace", "cccc3333"},
	}
	for i, tc := range cases {
		b := r.NewBuild(tc.branch, tc.commit)
		if b.ID != i+1 {
			t.Errorf("NewBuild #%d: ID = %d, хотел %d", i+1, b.ID, i+1)
		}
		if b.Branch != tc.branch || b.Commit != tc.commit {
			t.Errorf("NewBuild #%d: получил %+v", i+1, b)
		}
		if _, err := time.Parse(time.RFC3339, b.Date); err != nil {
			t.Errorf("NewBuild #%d: Date %q не RFC3339: %v", i+1, b.Date, err)
		}
		logFile := filepath.Join(logDir, fmt.Sprintf("%d.log", b.ID))
		if _, err := os.Stat(logFile); err != nil {
			t.Errorf("NewBuild #%d: лог-файл %s не создан: %v", b.ID, logFile, err)
		}
		if len(b.Log) == 0 {
			t.Errorf("NewBuild #%d: пустой Log", b.ID)
		}
	}
	if all := r.All(); len(all) != 3 {
		t.Errorf("All: %d сборок, хотел 3", len(all))
	}
}

func TestRetain(t *testing.T) {
	r, logDir := newTestRegistry(t)
	for i := 0; i < 4; i++ {
		r.NewBuild("master", fmt.Sprintf("commit-%d", i+1))
	}
	if err := r.Retain(2); err != nil {
		t.Fatalf("Retain(2): %v", err)
	}
	all := r.All()
	if len(all) != 2 || all[0].ID != 3 || all[1].ID != 4 {
		t.Fatalf("Retain(2): остались сборки с ID %v, хотел [3 4]", ids(all))
	}
	for _, id := range []int{1, 2} {
		if _, err := os.Stat(filepath.Join(logDir, fmt.Sprintf("%d.log", id))); !os.IsNotExist(err) {
			t.Errorf("Retain(2): лог-файл сборки %d не удалён", id)
		}
	}
	for _, id := range []int{3, 4} {
		if _, err := os.Stat(filepath.Join(logDir, fmt.Sprintf("%d.log", id))); err != nil {
			t.Errorf("Retain(2): лог-файл сборки %d пропал: %v", id, err)
		}
	}
	// Новый ID продолжает счётчик, а не начинается заново.
	if b := r.NewBuild("master", "dddd4444"); b.ID != 5 {
		t.Errorf("после Retain новый ID = %d, хотел 5", b.ID)
	}
}

func TestRetainInvalid(t *testing.T) {
	r, _ := newTestRegistry(t)
	if err := r.Retain(0); err == nil {
		t.Error("Retain(0): хотел ошибку, получил nil")
	}
}

func TestTraceability(t *testing.T) {
	r, _ := newTestRegistry(t)
	b := r.NewBuild("master", "cafe1234")
	for _, lang := range []string{"go", "ini"} {
		out, err := r.Traceability(b, lang)
		if err != nil {
			t.Fatalf("Traceability(%s): %v", lang, err)
		}
		for _, v := range []string{b.Branch, b.Commit, b.Date} {
			if !strings.Contains(out, v) {
				t.Errorf("Traceability(%s) не содержит %q:\n%s", lang, v, out)
			}
		}
	}
	if _, err := r.Traceability(b, "ada"); err == nil {
		t.Error("Traceability(ada): хотел ошибку, получил nil")
	}
}

func ids(bs []Build) []int {
	out := make([]int, len(bs))
	for i, b := range bs {
		out[i] = b.ID
	}
	return out
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

// Эталонное решение главы 30. Сравни со своей реализацией в buildsystem.go.
package jobs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Build — одна сборка: идентификатор, ветка, коммит, дата и лог.
type Build struct {
	ID     int
	Branch string
	Commit string
	Date   string
	Log    []byte
}

type Registry struct {
	logDir string
	builds []Build
	lastID int
}

func NewRegistry(logDir string) *Registry {
	return &Registry{logDir: logDir}
}

// NewBuild регистрирует сборку и создаёт её лог-файл.
func (r *Registry) NewBuild(branch, commit string) Build {
	r.lastID++
	b := Build{
		ID:     r.lastID,
		Branch: branch,
		Commit: commit,
		Date:   time.Now().UTC().Format(time.RFC3339),
	}
	b.Log = []byte(fmt.Sprintf("****jobs**** Build %d %s branch=%s commit=%s\n",
		b.ID, b.Date, b.Branch, b.Commit))
	if r.logDir != "" {
		os.MkdirAll(r.logDir, 0o755)
		os.WriteFile(r.logFile(b.ID), b.Log, 0o644)
	}
	r.builds = append(r.builds, b)
	return b
}

func (r *Registry) logFile(id int) string {
	return filepath.Join(r.logDir, fmt.Sprintf("%d.log", id))
}

// All возвращает копию реестра — по возрастанию ID.
func (r *Registry) All() []Build {
	out := make([]Build, len(r.builds))
	copy(out, r.builds)
	return out
}

// Retain оставляет последние n сборок, удаляя старые записи и их логи.
func (r *Registry) Retain(n int) error {
	if n < 1 {
		return fmt.Errorf("retention должен быть >= 1, получили %d", n)
	}
	for len(r.builds) > n {
		old := r.builds[0]
		if r.logDir != "" {
			os.Remove(r.logFile(old.ID))
		}
		r.builds = r.builds[1:]
	}
	return nil
}

// Traceability формирует сведения о сборке для встраивания в приложение.
func (r *Registry) Traceability(b Build, lang string) (string, error) {
	switch lang {
	case "go":
		return fmt.Sprintf(`// Сгенерировано jobs — не редактируйте вручную.
package revisions

const (
	Branch    = %q
	Commit    = %q
	BuildDate = %q
)
`, b.Branch, b.Commit, b.Date), nil
	case "ini":
		return fmt.Sprintf(`[jobs]
branch = %s
commit = %s
builddate = %s
`, b.Branch, b.Commit, b.Date), nil
	default:
		return "", fmt.Errorf("неизвестный язык трассировки %q (go или ini)", lang)
	}
}
```
{% endraw %}


</details>

_Пакет `assignments-go/30-buildsystem`, упражнение **buildsystem**. Локально: `go test -tags solution ./30-buildsystem/`._
