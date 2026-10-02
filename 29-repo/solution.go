//go:build solution

// Эталонное решение главы 29. Сравни со своей реализацией в repo.go.
package repo

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Project — один репозиторий группы.
type Project struct {
	RepoURL   string
	Reference string
	Path      string
}

// Config — разбор мини-YAML: workarea и список проектов.
type Config struct {
	Workarea string
	Projects []Project
}

// ParseConfig разбирает мини-YAML по строкам: до projects: — top-level
// ключи, после — список проектов по отступам.
func ParseConfig(r io.Reader) (Config, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Config{}, err
	}
	var (
		cfg        Config
		inProjects bool
		cur        *Project
	)
	for _, raw := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if !inProjects {
			key, val, err := splitKV(trimmed)
			if err != nil {
				return Config{}, err
			}
			switch key {
			case "workarea":
				if cfg.Workarea != "" {
					return Config{}, fmt.Errorf("повторный ключ workarea")
				}
				cfg.Workarea = val
			case "projects":
				if val != "" {
					return Config{}, fmt.Errorf("projects: не ждёт значения: %q", val)
				}
				inProjects = true
			default:
				return Config{}, fmt.Errorf("неизвестный ключ %q", key)
			}
			continue
		}
		if strings.HasPrefix(trimmed, "- ") {
			key, val, err := splitKV(strings.TrimSpace(trimmed[2:]))
			if err != nil {
				return Config{}, err
			}
			cfg.Projects = append(cfg.Projects, Project{})
			cur = &cfg.Projects[len(cfg.Projects)-1]
			if err := setField(cur, key, val); err != nil {
				return Config{}, err
			}
			continue
		}
		if cur == nil {
			return Config{}, fmt.Errorf("строка %q до первого элемента списка", trimmed)
		}
		key, val, err := splitKV(trimmed)
		if err != nil {
			return Config{}, err
		}
		if err := setField(cur, key, val); err != nil {
			return Config{}, err
		}
	}
	if !inProjects {
		return Config{}, fmt.Errorf("нет ключа projects:")
	}
	if cfg.Workarea == "" {
		return Config{}, fmt.Errorf("нет ключа workarea:")
	}
	return cfg, nil
}

// splitKV режет "ключ: значение" по первому двоеточию.
func splitKV(s string) (key, val string, err error) {
	i := strings.Index(s, ":")
	if i < 0 {
		return "", "", fmt.Errorf("нет двоеточия: %q", s)
	}
	return strings.TrimSpace(s[:i]), unquote(strings.TrimSpace(s[i+1:])), nil
}

// unquote срезает окружающие кавычки (двойные или одинарные).
func unquote(s string) string {
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}
	return s
}

func setField(p *Project, key, val string) error {
	switch key {
	case "repourl":
		p.RepoURL = val
	case "reference":
		p.Reference = val
	case "path":
		p.Path = val
	default:
		return fmt.Errorf("неизвестное поле проекта %q", key)
	}
	return nil
}

// Prepare проверяет пути на дубликаты, затем создаёт baseDir и каталоги
// всех проектов.
func Prepare(cfg Config, baseDir string) ([]string, error) {
	seen := map[string]string{}
	paths := make([]string, 0, len(cfg.Projects))
	for _, p := range cfg.Projects {
		if p.Path == "" {
			return nil, fmt.Errorf("пустой path у проекта %q", p.RepoURL)
		}
		if prev, dup := seen[p.Path]; dup {
			return nil, fmt.Errorf("дубликат path %q (у %q и %q)", p.Path, prev, p.RepoURL)
		}
		seen[p.Path] = p.RepoURL
		paths = append(paths, filepath.Join(baseDir, p.Path))
	}
	for _, p := range paths {
		if err := os.MkdirAll(p, 0o755); err != nil {
			return nil, err
		}
	}
	return paths, nil
}
