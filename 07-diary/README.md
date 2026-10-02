# Глава 7. Diary — спринт-дневник

> По мотивам `diary.pdf`. Инструмент участника проекта: задачи спринта,
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
