//go:build !solution

// Пакет letterboxed — глава 20. Реализуй ValidWord, AllWords и IsSolution
// так, чтобы прошли тесты в letterboxed_test.go. Условие — в README.md.
package letterboxed

// Puzzle — 12 букв по 4 сторонам, по 3 буквы на стороне (Sides[i]).
type Puzzle struct{ Sides [4]string }

// ValidWord сообщает, валидно ли слово w на пазле p: длина не меньше 2,
// каждая буква w лежит на пазле, соседние буквы w — на РАЗНЫХ сторонах,
// и w есть в словаре dict. Начинать можно с любой буквы пазла.
func ValidWord(p Puzzle, w string, dict []string) bool {
	panic("TODO: ValidWord")
}

// AllWords возвращает все валидные слова словаря dict на пазле p,
// отсортированные по возрастанию (лексикографически).
func AllWords(p Puzzle, dict []string) []string {
	panic("TODO: AllWords")
}

// IsSolution сообщает, является ли seq решением пазла p: каждое слово
// валидно, следующее слово начинается с последней буквы предыдущего,
// и вместе слова покрывают все 12 букв пазла. Пустая seq — не решение.
func IsSolution(p Puzzle, seq []string, dict []string) bool {
	panic("TODO: IsSolution")
}
