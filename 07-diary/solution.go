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
