# Глава 7. Diary — спринт-дневник

> Инструмент участника проекта: задачи спринта,
> статусы, отчёт и сериализация в JSON.

## Цель главы

Поработать с календарём и сериализацией: спринты с воскресенья
по субботу, задачи с оценкой усилия, JSON round-trip. В оригинале —
CLI с подкомандами `todo`/`done`/`blocked`/`report`; здесь ядро без обвязки.

## Теория

- `Weekday()` — день недели (`time.Sunday` — 0); откат к воскресенью
  недели: `at.AddDate(0, 0, -int(at.Weekday()))`.
- `time.Date(y, m, day, 0, 0, 0, 0, loc)` нормализует значения вне
  диапазона: `day = -1` уйдёт в конец прошлого месяца.
- `encoding/json`: `json.NewEncoder(w).Encode(v)` и
  `json.NewDecoder(r).Decode(&v)`; `time.Time` сериализуется как RFC 3339.
- `Add` меняет `Sprint` (`*Sprint`), `Report` — чистая функция от значения.

## Условие

1. `NewSprint(at, weeks)` — спринт на `weeks` недель, включающий `at`:
   `Start` — воскресенье недели `at` (полночь), `End` — суббота `weeks`-й недели.
2. `Add(desc, effort)` — добавить задачу в статусе `StatusTodo`.
3. `Report()` — заголовок спринта и по строке на задачу со статусом.
4. `Save`/`Load` — сериализация спринта в JSON и обратное чтение;
   испорченный JSON — ошибка.

Wants (необязательно): подкоманды `done`/`blocked` поверх ядра, отчёт всех спринтов.

## API для реализации

```go
type Status int
const (
	StatusTodo     Status = iota
	StatusDone
	StatusBlocked
)
type Item struct{ Desc string; Effort int; Status Status }
type Sprint struct{ Start, End time.Time; Items []Item }
func NewSprint(at time.Time, weeks int) Sprint
func (s *Sprint) Add(desc string, effort int)
func (s Sprint) Report() string
func Save(s Sprint, w io.Writer) error
func Load(r io.Reader) (Sprint, error)
```

Сигнатуры и договоры — в `diary.go`. Тесты — в `diary_test.go`.

## Подсказки

- `weeks < 1` трактуется как 1 — спринта короче недели не бывает.
- `End = Start.AddDate(0, 0, 7*weeks-1)`: при weeks=1 — суббота той же недели.
- Статусы меняет обвязка: `s.Items[i].Status = StatusDone` — поля открыты.
- Round-trip: `Save` → `bytes.Buffer` → `Load`; времена — `Equal`,
  задачи — `reflect.DeepEqual`.

## Проверка

```
go test ./07-diary/
go test -tags solution ./07-diary/
```

## Эталон

`go test -tags solution ./07-diary/`.

## Углубление

- Полная модель задачи: автоматические ID вида `S0:T0` (спринт:задача,
  выбор по `-t/--task`) и поле Blockage — причина блокировки (аргумент
  подкоманды blocked).
- CLI спеки: подкоманды create/modify/todo/done/blocked/report;
  `-f/--filename` с умолчанием `~/.diary/diary.dia` (`os.UserHomeDir`),
  `-d/--duration <n>w`, `-e/--estimate`, `-a/--all`; спринты нумеруются
  подряд, длительность при create — `1w…4w`.
- Отчёт спеки считает суммарный effort по задачам и статусам; close
  (из примера) закрывает текущий спринт досрочно, с предупреждением.
- Сериализация: спека допускает JSON, XML и базы данных. XML —
  `encoding/xml` из stdlib; БД потребует внешнего драйвера.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./07-diary/
```

## Заготовки (реализуйте)

### `diary.go`

{% raw %}
```go
//go:build !solution

// Пакет diary — глава 7. Реализуй функции так, чтобы прошли тесты
// в diary_test.go. Условие — в README.md этой папки.
package diary

import (
	"io"
	"time"
)

// Status — состояние задачи спринта.
type Status int

const (
	StatusTodo    Status = iota // задача поставлена
	StatusDone                  // задача завершена
	StatusBlocked               // задача заблокирована
)

// Item — задача: описание, оценка усилия (единицы произвольны),
// статус.
type Item struct {
	Desc   string
	Effort int
	Status Status
}

// Sprint — спринт: недели с воскресенья по субботу и задачи.
type Sprint struct {
	Start time.Time
	End   time.Time
	Items []Item
}

// NewSprint создаёт спринт длиной weeks недель, включающий дату at:
// Start — полночь воскресенья недели at, End — суббота weeks-й недели
// (weeks < 1 трактуется как 1).
func NewSprint(at time.Time, weeks int) Sprint {
	panic("TODO: NewSprint")
}

// Add добавляет задачу desc с оценкой усилия effort
// в статусе StatusTodo.
func (s *Sprint) Add(desc string, effort int) {
	panic("TODO: Add")
}

// Report возвращает отчёт: строка-заголовок спринта и по строке
// на задачу с её статусом (todo/done/blocked) и описанием.
func (s Sprint) Report() string {
	panic("TODO: Report")
}

// Save записывает спринт s в w в формате JSON.
func Save(s Sprint, w io.Writer) error {
	panic("TODO: Save")
}

// Load читает спринт из r (формат Save) и возвращает его.
// Испорченный JSON — ошибка.
func Load(r io.Reader) (Sprint, error) {
	panic("TODO: Load")
}
```
{% endraw %}

## Тесты

### `diary_test.go`

{% raw %}
```go
package diary

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Среда 2026-09-30: неделя начинается в воскресенье 2026-09-27.
var wednesday = time.Date(2026, time.September, 30, 15, 30, 0, 0, time.UTC)

func TestNewSprint(t *testing.T) {
	s := NewSprint(wednesday, 1)
	if s.Start.Weekday() != time.Sunday {
		t.Errorf("NewSprint: начало %v — %s, хотел воскресенье", s.Start, s.Start.Weekday())
	}
	if got := s.Start.Format("2006-01-02"); got != "2026-09-27" {
		t.Errorf("NewSprint: начало %s, хотел 2026-09-27 (воскресенье той же недели)", got)
	}
	if s.End.Weekday() != time.Saturday {
		t.Errorf("NewSprint: конец %v — %s, хотел субботу", s.End, s.End.Weekday())
	}
	if got := s.End.Format("2006-01-02"); got != "2026-10-03" {
		t.Errorf("NewSprint: конец %s, хотел 2026-10-03", got)
	}
	if !s.End.Equal(s.Start.AddDate(0, 0, 6)) {
		t.Error("NewSprint: спринт в 1 неделю должен кончаться через 6 дней после начала")
	}
}

func TestNewSprintSundayStart(t *testing.T) {
	// Воскресенье само — начало своей недели, даже поздним вечером.
	s := NewSprint(time.Date(2026, time.September, 27, 23, 0, 0, 0, time.UTC), 1)
	if got := s.Start.Format("2006-01-02"); got != "2026-09-27" {
		t.Errorf("NewSprint: для воскресенья начало %s, хотел 2026-09-27", got)
	}
}

func TestNewSprintWeeks(t *testing.T) {
	s := NewSprint(wednesday, 2)
	if got := s.Start.Format("2006-01-02"); got != "2026-09-27" {
		t.Errorf("NewSprint: начало %s, хотел 2026-09-27", got)
	}
	// 2 недели: 27.09 + 13 дней = суббота 2026-10-10.
	if got := s.End.Format("2006-01-02"); got != "2026-10-10" {
		t.Errorf("NewSprint: конец %s, хотел 2026-10-10", got)
	}
}

func TestAdd(t *testing.T) {
	s := NewSprint(wednesday, 1)
	s.Add("написать тесты", 3)
	s.Add("ревью PR", 1)
	if len(s.Items) != 2 {
		t.Fatalf("Add: задач %d, хотел 2", len(s.Items))
	}
	want := Item{Desc: "написать тесты", Effort: 3, Status: StatusTodo}
	if s.Items[0] != want {
		t.Errorf("Add: получил %+v, хотел %+v", s.Items[0], want)
	}
}

func TestReport(t *testing.T) {
	s := NewSprint(wednesday, 1)
	s.Add("написать тесты", 3)
	s.Add("ревью PR", 1)
	s.Add("ждать сборку CI", 2)
	s.Items[1].Status = StatusDone
	s.Items[2].Status = StatusBlocked

	rep := s.Report()
	if got := strings.Count(rep, "\n"); got != 4 { // заголовок + 3 задачи
		t.Errorf("Report: строк %d, хотел 4:\n%s", got, rep)
	}
	for _, want := range []string{
		"todo", "done", "blocked",
		"написать тесты", "ревью PR", "ждать сборку CI",
	} {
		if !strings.Contains(rep, want) {
			t.Errorf("Report: нет %q в отчёте:\n%s", want, rep)
		}
	}
}

func TestSaveLoad(t *testing.T) {
	s := NewSprint(wednesday, 2)
	s.Add("написать тесты", 3)
	s.Add("ревью PR", 1)
	s.Items[0].Status = StatusDone

	var buf bytes.Buffer
	if err := Save(s, &buf); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load(&buf)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !got.Start.Equal(s.Start) || !got.End.Equal(s.End) {
		t.Errorf("Load: промежуток %s..%s, хотел %s..%s",
			got.Start, got.End, s.Start, s.End)
	}
	if !reflect.DeepEqual(got.Items, s.Items) {
		t.Errorf("Load: задачи %+v, хотел %+v", got.Items, s.Items)
	}
}

func TestLoadBadJSON(t *testing.T) {
	if _, err := Load(strings.NewReader("{не JSON")); err == nil {
		t.Error("Load: хотел ошибку для испорченного JSON, получил nil")
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

// Эталонное решение главы 7. Сравни со своей реализацией в diary.go.
package diary

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type Status int

const (
	StatusTodo Status = iota
	StatusDone
	StatusBlocked
)

var statusName = [...]string{"todo", "done", "blocked"}

type Item struct {
	Desc   string
	Effort int
	Status Status
}

type Sprint struct {
	Start time.Time
	End   time.Time
	Items []Item
}

func NewSprint(at time.Time, weeks int) Sprint {
	if weeks < 1 {
		weeks = 1
	}
	// Полночь воскресенья недели, содержащей at.
	start := time.Date(at.Year(), at.Month(), at.Day()-int(at.Weekday()), 0, 0, 0, 0, at.Location())
	return Sprint{Start: start, End: start.AddDate(0, 0, 7*weeks-1)}
}

func (s *Sprint) Add(desc string, effort int) {
	s.Items = append(s.Items, Item{Desc: desc, Effort: effort, Status: StatusTodo})
}

func (s Sprint) Report() string {
	var b strings.Builder
	fmt.Fprintf(&b, "спринт %s .. %s, задач %d\n",
		s.Start.Format("2006-01-02"), s.End.Format("2006-01-02"), len(s.Items))
	for _, it := range s.Items {
		fmt.Fprintf(&b, "%s: %s (effort %d)\n", statusName[it.Status], it.Desc, it.Effort)
	}
	return b.String()
}

func Save(s Sprint, w io.Writer) error {
	return json.NewEncoder(w).Encode(s)
}

func Load(r io.Reader) (Sprint, error) {
	var s Sprint
	if err := json.NewDecoder(r).Decode(&s); err != nil {
		return Sprint{}, err
	}
	return s, nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/07-diary`, упражнение **diary**. Локально: `go test -tags solution ./07-diary/`._
