//go:build !solution

// Пакет jobs — глава 30. Реализуй типы и методы так, чтобы прошли тесты
// в buildsystem_test.go. Условие — в README.md этой папки.
package jobs

// Build — одна сборка: идентификатор (с 1), ветка, коммит, дата
// (UTC, RFC3339) и начальное содержимое лога.
type Build struct {
	ID     int
	Branch string
	Commit string
	Date   string
	Log    []byte
}

// Registry — реестр сборок; лог-файлы лежат в logDir/<id>.log.
// Добавь свои поля.
type Registry struct{}

// NewRegistry создаёт реестр сборок с логами в logDir.
func NewRegistry(logDir string) *Registry {
	panic("TODO: NewRegistry")
}

// NewBuild регистрирует новую сборку: ID — следующий по счёту (с 1),
// Date — текущее время UTC в RFC3339, лог-файл logDir/<id>.log создаётся
// и наполняется заголовком (Build.Log — то же содержимое).
func (r *Registry) NewBuild(branch, commit string) Build {
	panic("TODO: NewBuild")
}

// All возвращает все сборки в порядке возрастания ID.
func (r *Registry) All() []Build {
	panic("TODO: All")
}

// Retain оставляет последние n сборок: старые записи удаляются из
// реестра, их лог-файлы стираются. n < 1 — ошибка.
func (r *Registry) Retain(n int) error {
	panic("TODO: Retain")
}

// Traceability возвращает сведения о сборке b для встраивания в
// приложение: lang "go" — Go-константы Branch/Commit/BuildDate,
// "ini" — ini-текст с секцией [jobs]. Другой lang — ошибка.
func (r *Registry) Traceability(b Build, lang string) (string, error) {
	panic("TODO: Traceability")
}
