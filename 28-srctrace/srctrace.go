//go:build !solution

// Пакет srctrace — глава 28. Реализуй функции так, чтобы прошли тесты
// в srctrace_test.go. Условие — в README.md этой папки.
package srctrace

// RepoInfo — сведения о репозитории на момент сборки.
type RepoInfo struct {
	Branch string // имя ветки, например "master"
	Short  string // короткий id коммита, например "cefa722"
	Commit string // полный id коммита (40 hex-символов)
	Date   string // дата/время сборки, например RFC3339
}

// Snippet возвращает валидный Go-файл (строка package + блок const
// с константами Branch, ShortCommit, Commit, BuildDate), компилируемый
// в приложение и хранящий данные info. Обратный разбор — Parse.
func Snippet(info RepoInfo) string {
	panic("TODO: Snippet")
}

// Parse разбирает Go-файл, сгенерированный Snippet, и возвращает
// исходные данные. Некорректный источник — ошибка: пустой, без строки
// package, без какой-либо из четырёх констант.
func Parse(src string) (RepoInfo, error) {
	panic("TODO: Parse")
}
