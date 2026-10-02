//go:build solution

// Эталонное решение главы 31. Сравни со своей реализацией в project.go.
package proj

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Instantiate(tmplDir, dstDir string, vars map[string]string) error {
	if _, err := os.Stat(tmplDir); err != nil {
		return fmt.Errorf("шаблон %s недоступен: %w", tmplDir, err)
	}
	return filepath.WalkDir(tmplDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(tmplDir, p)
		if err != nil {
			return err
		}
		dst := filepath.Join(dstDir, rel)
		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		content := string(data)
		for name, val := range vars {
			content = strings.ReplaceAll(content, "{{"+name+"}}", val)
		}
		return os.WriteFile(dst, []byte(content), 0o644)
	})
}
