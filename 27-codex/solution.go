//go:build solution

// Эталонное решение главы 27. Сравни со своей реализацией в codex.go.
package codex

import (
	"fmt"
	"strconv"
	"strings"
)

// Fragment — спецификация фрагмента кода.
type Fragment struct {
	From  int
	Count int
}

func ParseSpecs(args []string) ([]Fragment, error) {
	out := make([]Fragment, 0, len(args))
	for _, arg := range args {
		fromS, countS, ok := strings.Cut(arg, ":")
		if !ok {
			return nil, fmt.Errorf("codex: spec %q: хочу from:count", arg)
		}
		from, err1 := strconv.Atoi(fromS)
		count, err2 := strconv.Atoi(countS)
		if err1 != nil || err2 != nil || from < 1 || count < 1 {
			return nil, fmt.Errorf("codex: spec %q: from и count — целые ≥ 1", arg)
		}
		out = append(out, Fragment{From: from, Count: count})
	}
	return out, nil
}

func Extract(lines []string, f Fragment) ([]string, error) {
	if f.From < 1 || f.Count < 1 || f.From+f.Count-1 > len(lines) {
		return nil, fmt.Errorf("codex: фрагмент %d:%d выходит за границы файла из %d строк", f.From, f.Count, len(lines))
	}
	return lines[f.From-1 : f.From-1+f.Count], nil
}
