package repo

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sample = `
# Пример в духе AllProjects.yaml.
workarea: "$HOME/grepotest/pub"
projects:
  - repourl: "git@gitlab.com:projtemplates/go.git"
    reference: "master"
    path: "go/lib"
  - repourl: "https://github.com/RajaSrinivasan/disign.git"
    path: "disign"
`

func TestParseConfig(t *testing.T) {
	cfg, err := ParseConfig(strings.NewReader(sample))
	if err != nil {
		t.Fatalf("ParseConfig: %v", err)
	}
	if cfg.Workarea != "$HOME/grepotest/pub" {
		t.Errorf("Workarea = %q, хотел %q (строка как есть)", cfg.Workarea, "$HOME/grepotest/pub")
	}
	if len(cfg.Projects) != 2 {
		t.Fatalf("Projects: %d шт., хотел 2", len(cfg.Projects))
	}
	first := cfg.Projects[0]
	if first.RepoURL != "git@gitlab.com:projtemplates/go.git" ||
		first.Reference != "master" || first.Path != "go/lib" {
		t.Errorf("первый проект = %+v", first)
	}
	if cfg.Projects[1].RepoURL != "https://github.com/RajaSrinivasan/disign.git" ||
		cfg.Projects[1].Reference != "" || cfg.Projects[1].Path != "disign" {
		t.Errorf("второй проект = %+v", cfg.Projects[1])
	}
}

func TestParseConfigErrors(t *testing.T) {
	bad := []string{
		"workarea: /tmp\nprojects:\n  - repourl: x.git\n    oops\n", // строка без двоеточия
		"projects:\n  - repourl: \"x.git\"\n    path: p\n",          // нет workarea
		"workarea: /tmp\n", // нет projects
		"garbage\n",        // нет двоеточия
		"workarea: /tmp\nprojects:\n  - repourl: x.git\n    path: p\n    build: make\n", // неизвестное поле
	}
	for _, src := range bad {
		if _, err := ParseConfig(strings.NewReader(src)); err == nil {
			t.Errorf("ParseConfig(%q): хотел ошибку, получил nil", src)
		}
	}
}

func TestPrepare(t *testing.T) {
	cfg, err := ParseConfig(strings.NewReader(sample))
	if err != nil {
		t.Fatalf("ParseConfig: %v", err)
	}
	base := filepath.Join(t.TempDir(), "base")
	paths, err := Prepare(cfg, base)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("Prepare вернул %d путей, хотел 2: %v", len(paths), paths)
	}
	for i, p := range paths {
		if want := filepath.Join(base, cfg.Projects[i].Path); p != want {
			t.Errorf("путь #%d = %q, хотел %q", i, p, want)
		}
		if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
			t.Errorf("каталог %s не создан (err=%v)", p, err)
		}
	}
}

func TestPrepareDuplicatePath(t *testing.T) {
	cfg := Config{
		Workarea: "/tmp/wa",
		Projects: []Project{
			{RepoURL: "a.git", Path: "dup"},
			{RepoURL: "b.git", Path: "dup"},
		},
	}
	if _, err := Prepare(cfg, filepath.Join(t.TempDir(), "base")); err == nil {
		t.Fatal("Prepare: хотел ошибку для дубликата path, получил nil")
	}
}
