package proj

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	tmplMain  = "// Проект {{PROJECT}}.\n// Автор: {{USER}}.\npackage main\n\nfunc main() {}\n"
	tmplSub   = "{{PROJECT}} / {{USER}}\n"
	tmplPlain = "никаких плейсхолдеров тут нет\n"
)

var vars = map[string]string{"PROJECT": "goutil", "USER": "pots"}

func setupTemplate(t *testing.T) string {
	t.Helper()
	tmpl := filepath.Join(t.TempDir(), "tmpl")
	if err := os.MkdirAll(filepath.Join(tmpl, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	for path, content := range map[string]string{
		filepath.Join(tmpl, "main.go"):         tmplMain,
		filepath.Join(tmpl, "sub", "note.txt"): tmplSub,
		filepath.Join(tmpl, "README.txt"):      tmplPlain,
	} {
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return tmpl
}

func TestInstantiate(t *testing.T) {
	tmpl := setupTemplate(t)
	dst := filepath.Join(t.TempDir(), "dst")
	if err := Instantiate(tmpl, dst, vars); err != nil {
		t.Fatalf("Instantiate: %v", err)
	}
	checks := []struct{ rel, want string }{
		{"main.go", "// Проект goutil.\n// Автор: pots.\npackage main\n\nfunc main() {}\n"},
		{filepath.Join("sub", "note.txt"), "goutil / pots\n"},
		{"README.txt", tmplPlain}, // без плейсхолдеров — не тронут
	}
	for _, c := range checks {
		got, err := os.ReadFile(filepath.Join(dst, c.rel))
		if err != nil {
			t.Errorf("файл %s не создан: %v", c.rel, err)
			continue
		}
		if string(got) != c.want {
			t.Errorf("файл %s:\nполучил %q\nхотел  %q", c.rel, got, c.want)
		}
	}
}

func TestInstantiateMissingTemplate(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	dst := filepath.Join(t.TempDir(), "dst")
	if err := Instantiate(missing, dst, vars); err == nil {
		t.Fatal("Instantiate: хотел ошибку для отсутствующего шаблона, получил nil")
	}
}
