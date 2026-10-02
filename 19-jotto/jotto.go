//go:build !solution

// Пакет jotto — глава 19. Реализуй Common и решателя так,
// чтобы прошли тесты в jotto_test.go. Условие — в README.md.
package jotto

// Common возвращает число общих букв строк a и b — пересечение
// мультимножеств букв: для каждой буквы минимум числа вхождений.
// Правильная реализация симметрична: Common(a, b) == Common(b, a)
// (наивная «буква есть в слове» — нет, см. README). Иди по рунам,
// а не по байтам: кириллическая буква занимает 2 байта в UTF-8.
//
//	Common("марс", "срам") == 4
//	Common("ааа", "аа")    == 2
func Common(a, b string) int {
	panic("TODO: Common")
}

// Solver — компьютер угадывает задуманное 4-буквенное слово,
// вычёркивая кандидатов по единственному счёту.
type Solver struct {
	// TODO: твои поля (например, слайс живых кандидатов).
}

// NewSolver создаёт решателя: кандидаты — 4-буквенные слова
// словаря Words из words.go.
func NewSolver() *Solver {
	panic("TODO: NewSolver")
}

// Guess возвращает текущую гипотезу — слово из живых кандидатов
// (детерминированную, например первого живого).
func (s *Solver) Guess() string {
	panic("TODO: Guess")
}

// Learn принимает ответ n = Common(secret, guess) на попытку guess.
// Попытка была неверна, поэтому вычёркивается и сама: остаются
// только кандидаты cand с Common(cand, guess) == n.
func (s *Solver) Learn(guess string, n int) {
	panic("TODO: Learn")
}
