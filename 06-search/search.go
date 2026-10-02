//go:build !solution

// Пакет search — глава 6. Реализуй функции так, чтобы прошли тесты
// в search_test.go. Условие — в README.md этой папки.
package search

// Options задаёт режим поиска: Pattern — needle это регулярное
// выражение (иначе подстрока); CaseSensitive — учитывать регистр
// (по умолчанию поиск без учёта регистра).
type Options struct {
	Pattern       bool
	CaseSensitive bool
}

// Hit — одно совпавшее место файла: Line — номер строки (с 1),
// Text — текст строки, Count — число вхождений кандидата в строку.
type Hit struct {
	Line  int
	Text  string
	Count int
}

// FindLine возвращает число вхождений needle в line и найдено ли
// хоть одно. Pattern → regexp, иначе подстрока; регистронезависимо,
// если не CaseSensitive. Пустой needle — (0, false).
func FindLine(line, needle string, o Options) (int, bool) {
	panic("TODO: FindLine")
}

// SearchFile ищет needle в файле path и возвращает совпадения
// по строкам в порядке возрастания. Нет совпадений — пустой
// результат и nil. Отсутствующий или нечитаемый файл — ошибка;
// некорректный Pattern — совпадений нет.
func SearchFile(path, needle string, o Options) ([]Hit, error) {
	panic("TODO: SearchFile")
}

// ListFiles разворачивает paths в список файлов. Директория —
// её файлы (с recurse — всё дерево через filepath.WalkDir, без —
// только верхний уровень); обычное имя файла — само имя; имя
// с wildcard-метасимволами — совпадения filepath.Glob. Путь,
// которого нет и который ни с чем не совпал, — ошибка.
func ListFiles(paths []string, recurse bool) ([]string, error) {
	panic("TODO: ListFiles")
}
