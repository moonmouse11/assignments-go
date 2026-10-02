//go:build solution

// Эталонное решение главы 28. Сравни со своей реализацией в srctrace.go.
package srctrace

import (
	"fmt"
	"regexp"
	"strings"
)

// RepoInfo — сведения о репозитории на момент сборки.
type RepoInfo struct {
	Branch string // имя ветки, например "master"
	Short  string // короткий id коммита, например "cefa722"
	Commit string // полный id коммита (40 hex-символов)
	Date   string // дата/время сборки, например RFC3339
}

// Snippet генерирует валидный Go-файл с константами трассировки.
func Snippet(info RepoInfo) string {
	var b strings.Builder
	b.WriteString("// Сгенерировано srctrace — не редактируйте вручную.\n")
	b.WriteString("package revisions\n\n")
	b.WriteString("const (\n")
	fmt.Fprintf(&b, "\tBranch      = %q\n", info.Branch)
	fmt.Fprintf(&b, "\tShortCommit = %q\n", info.Short)
	fmt.Fprintf(&b, "\tCommit      = %q\n", info.Commit)
	fmt.Fprintf(&b, "\tBuildDate   = %q\n", info.Date)
	b.WriteString(")\n")
	return b.String()
}

// constRe ловит строки вида `Branch = "master"` (начало строки, имя,
// знак =, закавыченное значение). Комментарии не совпадают.
var constRe = regexp.MustCompile(`(?m)^\s*(Branch|ShortCommit|Commit|BuildDate)\s*=\s*"([^"\n]*)"`)

// Parse разбирает файл, сгенерированный Snippet.
func Parse(src string) (RepoInfo, error) {
	if strings.TrimSpace(src) == "" {
		return RepoInfo{}, fmt.Errorf("пустой источник")
	}
	hasPkg := false
	for _, line := range strings.Split(src, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "package ") {
			hasPkg = true
			break
		}
	}
	if !hasPkg {
		return RepoInfo{}, fmt.Errorf("нет строки package")
	}
	var info RepoInfo
	found := map[string]bool{}
	for _, m := range constRe.FindAllStringSubmatch(src, -1) {
		found[m[1]] = true
		switch m[1] {
		case "Branch":
			info.Branch = m[2]
		case "ShortCommit":
			info.Short = m[2]
		case "Commit":
			info.Commit = m[2]
		case "BuildDate":
			info.Date = m[2]
		}
	}
	for _, name := range []string{"Branch", "ShortCommit", "Commit", "BuildDate"} {
		if !found[name] {
			return RepoInfo{}, fmt.Errorf("нет константы %s", name)
		}
	}
	return info, nil
}
