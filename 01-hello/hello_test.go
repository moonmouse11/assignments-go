package hello

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestNumberLines(t *testing.T) {
	got := NumberLines([]string{"alpha", "beta", "", "gamma"})
	want := "0001 alpha\n0002 beta\n0003 \n0004 gamma\n"
	if got != want {
		t.Errorf("NumberLines:\nполучил %q\nхотел  %q", got, want)
	}
}

func TestRun(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.txt")
	f2 := filepath.Join(dir, "b.txt")
	os.WriteFile(f1, []byte("one\ntwo\n"), 0o644)
	os.WriteFile(f2, []byte("three"), 0o644)

	var buf bytes.Buffer
	if err := Run([]string{f1, f2}, &buf); err != nil {
		t.Fatalf("Run: %v", err)
	}
	// Нумерация начинается заново для каждого файла.
	want := "0001 one\n0002 two\n0001 three\n"
	if buf.String() != want {
		t.Errorf("Run:\nполучил %q\nхотел  %q", buf.String(), want)
	}
}

func TestRunMissingFile(t *testing.T) {
	if err := Run([]string{filepath.Join(t.TempDir(), "nope.txt")}, io.Discard); err == nil {
		t.Fatal("Run: хотел ошибку для отсутствующего файла, получил nil")
	}
}
