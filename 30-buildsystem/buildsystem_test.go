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
