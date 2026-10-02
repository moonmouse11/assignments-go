//go:build !solution

// Пакет hanoi — глава 4. Реализуй так, чтобы прошли тесты.
package hanoi

import "fmt"

// Move — ход: перенос одного диска с колышка From на колышек To (0..2).
type Move struct {
	From, To int
}

// SolveRecursive возвращает последовательность ходов переноса n дисков
// с колышка 0 на колышек 2 (рекурсивное решение).
func SolveRecursive(n int) []Move {
	panic("TODO: SolveRecursive")
}

// SolveIterative решает ту же задачу итеративно; последовательность ходов
// совпадает с SolveRecursive(n) в точности.
func SolveIterative(n int) []Move {
	panic("TODO: SolveIterative")
}

// Verify симулирует ходы: каждый ход законен (с непустого колышка,
// меньший диск на больший, колышки в 0..2) и в конце все n дисков
// на колышке 2. Незаконный ход — error с номером хода.
func Verify(moves []Move, n int) error {
	_ = fmt.Sprintf // убрать, когда появится настоящий код
	panic("TODO: Verify")
}
