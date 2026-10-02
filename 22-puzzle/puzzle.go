//go:build !solution

// Пакет puzzle — глава 22. Реализуй доску, ходы, перемешивание
// и Solvable так, чтобы прошли тесты в puzzle_test.go.
// Условие — в README.md.
package puzzle

// Board — доска N×N: клетки в row-major слайсе, пустая клетка — 0.
type Board struct {
	// TODO: твои поля (размер и слайс клеток).
}

// New возвращает решённую доску n×n: плитки 1..n²-1 по порядку,
// пустая клетка 0 — в конце.
func New(n int) *Board {
	panic("TODO: New")
}

// Move меняет местами плитку tile и пустую клетку. Плитка обязана
// соседствовать с пустой клеткой (общая сторона), иначе ошибка.
// Нет такой плитки — тоже ошибка.
func (b *Board) Move(tile int) error {
	panic("TODO: Move")
}

// IsSolved сообщает, решена ли доска (плитки по порядку, 0 в конце).
func (b *Board) IsSolved() bool {
	panic("TODO: IsSolved")
}

// String — стабильный формат: строки через \n, клетки через пробел,
// каждая клетка — %2d, пустая клетка — « _». В конце — \n.
//
//	New(3).String() == " 1  2  3\n 4  5  6\n 7  8  _\n"
func (b *Board) String() string {
	panic("TODO: String")
}

// Shuffle возвращает решённую доску n×n, перемешанную случайными
// ВАЛИДНЫМИ ходами — поэтому всегда решаемую. ГСЧ — math/rand/v2:
// rand.New(rand.NewChaCha8(sb)), где sb [32]byte — seed, развёрнутый
// в байты. Одинаковый seed — одинаковая доска.
func Shuffle(n int, seed uint64) *Board {
	panic("TODO: Shuffle")
}

// Solvable сообщает, решаема ли перестановка tiles (row-major, длина n²,
// пустая клетка — 0): для нечётного n — чётность инверсий; для чётного n —
// инверсии + строка пустой клетки, считая снизу (см. README).
func Solvable(tiles []int, n int) bool {
	panic("TODO: Solvable")
}
