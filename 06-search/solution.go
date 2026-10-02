//go:build solution

// Эталонное решение главы 6. Сравни со своей реализацией в search.go.
package search

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type Options struct{ Pattern, CaseSensitive bool }

type Hit struct {
	Line, Count int
	Text        string
}

func FindLine(line, needle string, o Options) (int, bool) {
	if needle == "" {
		return 0, false
	}
	if !o.Pattern {
		needle = regexp.QuoteMeta(needle) // подстрока — литеральный паттерн
	}
	if !o.CaseSensitive {
		needle = "(?i)" + needle
	}
	if re, err := regexp.Compile(needle); err == nil {
		n := len(re.FindAllStringIndex(line, -1))
		return n, n > 0
	}
	return 0, false
}

func SearchFile(path, needle string, o Options) ([]Hit, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	var hits []Hit
	for n := 0; sc.Scan(); n++ {
		if c, ok := FindLine(sc.Text(), needle, o); ok {
			hits = append(hits, Hit{Line: n + 1, Text: sc.Text(), Count: c})
		}
	}
	return hits, sc.Err()
}

func ListFiles(paths []string, recurse bool) ([]string, error) {
	var files []string
	for _, p := range paths {
		info, err := os.Stat(p)
		if err == nil && info.IsDir() {
			err = filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
				if err == nil && d.IsDir() && path != p && !recurse {
					return fs.SkipDir
				}
				if err == nil && !d.IsDir() {
					files = append(files, path)
				}
				return err
			})
			if err != nil {
				return nil, err
			}
		} else if err == nil {
			files = append(files, p)
		} else if m, _ := filepath.Glob(p); len(m) == 0 {
			return nil, err // os.Stat: пути нет и glob пуст
		} else {
			files = append(files, m...)
		}
	}
	return files, nil
}
