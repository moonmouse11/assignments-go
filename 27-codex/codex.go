//go:build !solution

// Пакет codex — глава 27. Реализуй ParseSpecs и Extract так, чтобы прошли
// тесты в codex_test.go. Условие — в README.md этой папки.
package codex

// Fragment — спецификация фрагмента кода: Count строк начиная со строки
// From (нумерация с 1).
type Fragment struct {
	From  int // номер первой строки, 1-индексация
	Count int // число строк, ≥ 1
}

// ParseSpecs разбирает аргументы вида "from:count" в фрагменты.
// count обязана быть ≥ 1; невалидный аргумент — ошибка, называющая его.
func ParseSpecs(args []string) ([]Fragment, error) {
	panic("TODO: ParseSpecs")
}

// Extract возвращает строки lines с From по From+Count-1 (1-индексация).
// Выход за границы файла (From < 1, Count < 1 или From+Count-1 > len(lines))
// — ошибка.
func Extract(lines []string, f Fragment) ([]string, error) {
	panic("TODO: Extract")
}
