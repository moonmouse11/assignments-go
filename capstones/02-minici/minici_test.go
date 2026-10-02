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
