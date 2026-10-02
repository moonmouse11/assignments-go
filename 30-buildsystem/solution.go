//go:build solution

// Эталонное решение главы 30. Сравни со своей реализацией в buildsystem.go.
package jobs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Build — одна сборка: идентификатор, ветка, коммит, дата и лог.
type Build struct {
	ID     int
	Branch string
	Commit string
	Date   string
	Log    []byte
}

type Registry struct {
	logDir string
	builds []Build
	lastID int
}

func NewRegistry(logDir string) *Registry {
	return &Registry{logDir: logDir}
}

// NewBuild регистрирует сборку и создаёт её лог-файл.
func (r *Registry) NewBuild(branch, commit string) Build {
	r.lastID++
	b := Build{
		ID:     r.lastID,
		Branch: branch,
		Commit: commit,
		Date:   time.Now().UTC().Format(time.RFC3339),
	}
	b.Log = []byte(fmt.Sprintf("****jobs**** Build %d %s branch=%s commit=%s\n",
		b.ID, b.Date, b.Branch, b.Commit))
	if r.logDir != "" {
		os.MkdirAll(r.logDir, 0o755)
		os.WriteFile(r.logFile(b.ID), b.Log, 0o644)
	}
	r.builds = append(r.builds, b)
	return b
}

func (r *Registry) logFile(id int) string {
	return filepath.Join(r.logDir, fmt.Sprintf("%d.log", id))
}

// All возвращает копию реестра — по возрастанию ID.
func (r *Registry) All() []Build {
	out := make([]Build, len(r.builds))
	copy(out, r.builds)
	return out
}

// Retain оставляет последние n сборок, удаляя старые записи и их логи.
func (r *Registry) Retain(n int) error {
	if n < 1 {
		return fmt.Errorf("retention должен быть >= 1, получили %d", n)
	}
	for len(r.builds) > n {
		old := r.builds[0]
		if r.logDir != "" {
			os.Remove(r.logFile(old.ID))
		}
		r.builds = r.builds[1:]
	}
	return nil
}

// Traceability формирует сведения о сборке для встраивания в приложение.
func (r *Registry) Traceability(b Build, lang string) (string, error) {
	switch lang {
	case "go":
		return fmt.Sprintf(`// Сгенерировано jobs — не редактируйте вручную.
package revisions

const (
	Branch    = %q
	Commit    = %q
	BuildDate = %q
)
`, b.Branch, b.Commit, b.Date), nil
	case "ini":
		return fmt.Sprintf(`[jobs]
branch = %s
commit = %s
builddate = %s
`, b.Branch, b.Commit, b.Date), nil
	default:
		return "", fmt.Errorf("неизвестный язык трассировки %q (go или ini)", lang)
	}
}
