//go:build !solution

// Пакет templates — глава 32. Реализуй функцию так, чтобы прошли тесты
// в templates_test.go. Условие — в README.md этой папки.
package templates

// Scaffold возвращает файлы заготовки проекта name:
//
//	kind "cli" — "main.go" (парсинг флагов -h/--help, -v/--version,
//	             --verbose по конвенциям главы 02) и "Makefile";
//	kind "web" — "main.go" (http.HandleFunc("/", ...) +
//	             ListenAndServe) и "go.mod" (module name).
//
// Неизвестный kind или пустое name — ошибка.
func Scaffold(kind, name string) (map[string]string, error) {
	panic("TODO: Scaffold")
}
