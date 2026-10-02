//go:build !solution

// Пакет repo — глава 29. Реализуй функции так, чтобы прошли тесты
// в repo_test.go. Условие — в README.md этой папки.
package repo

import "io"

// Project — один репозиторий группы: откуда клонировать (RepoURL),
// какую версию брать (Reference: ветка/тег/коммит) и куда положить
// (Path — путь относительно workarea, уникальный в группе).
type Project struct {
	RepoURL   string
	Reference string
	Path      string
}

// Config — разбор мини-YAML: workarea и список проектов.
type Config struct {
	Workarea string
	Projects []Project
}

// ParseConfig разбирает подмножество YAML (см. README): top-level ключ
// workarea, затем projects: с элементами "- repourl: ...",
// "reference: ...", "path: ..." по отступам. Значения в кавычках —
// строки (кавычки срезаются). Некорректная конфигурация — ошибка.
func ParseConfig(r io.Reader) (Config, error) {
	panic("TODO: ParseConfig")
}

// Prepare создаёт baseDir и каталоги Path всех проектов (под baseDir),
// возвращает список полных путей в порядке объявления. Дубликат Path
// или пустой Path — ошибка.
func Prepare(cfg Config, baseDir string) ([]string, error) {
	panic("TODO: Prepare")
}
