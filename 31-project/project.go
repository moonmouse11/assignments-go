//go:build !solution

// Пакет proj — глава 31. Реализуй функцию так, чтобы прошли тесты
// в project_test.go. Условие — в README.md этой папки.
package proj

// Instantiate копирует дерево шаблона tmplDir (включая поддиректории)
// в dstDir, заменяя в каждом файле вхождения плейсхолдеров:
// {{NAME}} → vars["NAME"] для каждого ключа vars (буквальная замена,
// не text/template). Файлы без плейсхолдеров копируются как есть.
// Отсутствующий tmplDir — ошибка.
func Instantiate(tmplDir, dstDir string, vars map[string]string) error {
	panic("TODO: Instantiate")
}
