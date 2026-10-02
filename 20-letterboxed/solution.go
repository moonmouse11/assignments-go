//go:build solution

// Эталонное решение главы 20. Сравни со своей реализацией в letterboxed.go.
package letterboxed

import (
	"sort"
	"strings"
)

// Puzzle — 12 букв по 4 сторонам, по 3 буквы на стороне (Sides[i]).
type Puzzle struct{ Sides [4]string }

// sideOf — номер стороны буквы r или -1, если буквы нет на пазле.
func sideOf(p Puzzle, r rune) int {
	for i, s := range p.Sides {
		if strings.ContainsRune(s, r) {
			return i
		}
	}
	return -1
}

// ValidWord: в словаре, буквы на пазле, соседние буквы на разных сторонах.
func ValidWord(p Puzzle, w string, dict []string) bool {
	rs := []rune(w)
	if len(rs) < 2 {
		return false
	}
	found := false
	for _, d := range dict {
		if d == w {
			found = true
			break
		}
	}
	if !found {
		return false
	}
	prev := -1
	for _, r := range rs {
		s := sideOf(p, r)
		if s < 0 || s == prev {
			return false
		}
		prev = s
	}
	return true
}

// AllWords — валидные слова словаря, по возрастанию.
func AllWords(p Puzzle, dict []string) []string {
	var out []string
	for _, w := range dict {
		if ValidWord(p, w, dict) {
			out = append(out, w)
		}
	}
	sort.Strings(out)
	return out
}

// IsSolution: цепочка валидных слов, покрывающая все 12 букв.
func IsSolution(p Puzzle, seq []string, dict []string) bool {
	if len(seq) == 0 {
		return false
	}
	covered := map[rune]bool{}
	for i, w := range seq {
		if !ValidWord(p, w, dict) {
			return false
		}
		rs := []rune(w)
		if i > 0 { // следующее слово начинается с последней буквы предыдущего
			prev := []rune(seq[i-1])
			if rs[0] != prev[len(prev)-1] {
				return false
			}
		}
		for _, r := range rs {
			covered[r] = true
		}
	}
	for _, s := range p.Sides {
		for _, r := range s {
			if !covered[r] {
				return false
			}
		}
	}
	return true
}
