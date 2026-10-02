//go:build solution

// Эталон капстоуна 2.
package minici

import (
	"fmt"
	"os"
	"path/filepath"

	"assignments-go/28-srctrace"
	jobs "assignments-go/30-buildsystem"
)

type CI struct {
	reg    *jobs.Registry
	logDir string
}

func New(logDir string) *CI {
	return &CI{reg: jobs.NewRegistry(logDir), logDir: logDir}
}

func (c *CI) Build(branch, commit string) (jobs.Build, error) {
	b := c.reg.NewBuild(branch, commit)
	short := commit
	if len(short) > 7 {
		short = short[:7]
	}
	info := srctrace.RepoInfo{Branch: branch, Short: short, Commit: commit, Date: b.Date}
	src := srctrace.Snippet(info)
	path := filepath.Join(c.logDir, fmt.Sprintf("trace-%d.go", b.ID))
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		return b, err
	}
	return b, nil
}

func (c *CI) All() []jobs.Build {
	return c.reg.All()
}

func (c *CI) Trace(b jobs.Build, lang string) (string, error) {
	return c.reg.Traceability(b, lang)
}

func (c *CI) Retain(n int) error {
	// Сначала узнаём, кто выпадет, затем убираем trace-файлы и логи реестра.
	var old []int
	all := c.reg.All()
	if len(all) > n {
		for _, b := range all[:len(all)-n] {
			old = append(old, b.ID)
		}
	}
	if err := c.reg.Retain(n); err != nil {
		return err
	}
	for _, id := range old {
		if err := os.Remove(filepath.Join(c.logDir, fmt.Sprintf("trace-%d.go", id))); err != nil {
			return err
		}
	}
	return nil
}
