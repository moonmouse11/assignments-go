package passwordmanager

import (
	"path/filepath"
	"testing"
)

var words = []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot"}

func TestGenerateSaveOpen(t *testing.T) {
	m := New(42, words, "master-key")
	pw1, err := m.Generate("example.com", "alice", 3)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	pw2, err := m.Generate("gitlab.com", "bob", 3)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	path := filepath.Join(t.TempDir(), "vault.ppm")
	if err := m.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	opened, err := Open(path, "master-key")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	for ctx, want := range map[string]string{"example.com": pw1, "gitlab.com": pw2} {
		if got, ok := opened.Get(ctx); !ok || got != want {
			t.Errorf("Get(%q) после Open: получил %q, ok=%v, хотел %q", ctx, got, ok, want)
		}
	}
}

func TestOpenWrongMaster(t *testing.T) {
	m := New(42, words, "master-key")
	if _, err := m.Generate("example.com", "alice", 2); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	path := filepath.Join(t.TempDir(), "vault.ppm")
	if err := m.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := Open(path, "не-тот-master"); err == nil {
		t.Error("Open с неверным master: хотел ошибку, получил nil")
	}
}

func TestGenerateEmptyContext(t *testing.T) {
	m := New(42, words, "master-key")
	if _, err := m.Generate("", "alice", 2); err == nil {
		t.Error("Generate с пустым context: хотел ошибку")
	}
}

func TestOpenedManagerHasNoGenerator(t *testing.T) {
	m := New(42, words, "master-key")
	if _, err := m.Generate("example.com", "alice", 2); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	path := filepath.Join(t.TempDir(), "vault.ppm")
	if err := m.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	opened, err := Open(path, "master-key")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if _, err := opened.Generate("new.com", "carol", 2); err == nil {
		t.Error("Generate у открытого менеджера: хотел ошибку (нет генератора)")
	}
}

func TestListSorted(t *testing.T) {
	m := New(42, words, "master-key")
	m.Generate("zzz.example", "a", 2)
	m.Generate("aaa.example", "b", 2)
	list := m.List()
	if len(list) != 2 {
		t.Fatalf("List: записей %d, хотел 2", len(list))
	}
	if list[0].Context != "aaa.example" || list[1].Context != "zzz.example" {
		t.Errorf("List не отсортирован: %v, %v", list[0].Context, list[1].Context)
	}
}
