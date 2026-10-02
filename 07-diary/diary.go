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
