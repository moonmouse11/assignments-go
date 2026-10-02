//go:build solution

// Эталонное решение главы 19. Сравни со своей реализацией в jotto.go.
package jotto

// Common — пересечение мультимножеств букв (минимум вхождений каждой буквы).
func Common(a, b string) int {
	counts := make(map[rune]int)
	for _, r := range a {
		counts[r]++
	}
	n := 0
	for _, r := range b {
		if counts[r] > 0 {
			counts[r]--
			n++
		}
	}
	return n
}

// Solver — угадыватель jotto: живые кандидаты из words.go.
type Solver struct {
	candidates []string
}

// NewSolver — все 4-буквенные слова словаря Words.
func NewSolver() *Solver {
	c := make([]string, 0, len(Words))
	for _, w := range Words {
		if len(w) == 4 {
			c = append(c, w)
		}
	}
	return &Solver{candidates: c}
}

// Guess — первый живой кандидат.
func (s *Solver) Guess() string {
	return s.candidates[0]
}

// Learn оставляет кандидатов с тем же счётом против guess;
// сама попытка guess была неверна — вычёркиваем и её.
func (s *Solver) Learn(guess string, n int) {
	kept := s.candidates[:0]
	for _, cand := range s.candidates {
		if cand != guess && Common(cand, guess) == n {
			kept = append(kept, cand)
		}
	}
	s.candidates = kept
}
