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
