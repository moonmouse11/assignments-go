//go:build !solution

// Капстоун 2: мини-CI. Реализуй клей поверх 28-srctrace и 30-buildsystem
// так, чтобы прошли тесты в minici_test.go. Условие — в README.md.
package minici

import (
	"assignments-go/28-srctrace"
	jobs "assignments-go/30-buildsystem"
)

// var _ — держит импорт srctrace: понадобится в Build (см. решение).
var _ = srctrace.Snippet

// CI — реестр сборок (jobs) + трассировка каждой сборки (srctrace),
// встраиваемая в приложение. Все файлы лежат в logDir.
type CI struct {
	// TODO: твои поля.
}

// New создаёт CI с логами и трассировкой в logDir.
func New(logDir string) *CI {
	panic("TODO: New")
}

// Build регистрирует сборку (ID с 1) и пишет её трассировку
// в logDir/trace-<ID>.go (Short — первые 7 символов commit).
func (c *CI) Build(branch, commit string) (jobs.Build, error) {
	panic("TODO: Build")
}

// All возвращает все сборки по возрастанию ID.
func (c *CI) All() []jobs.Build {
	panic("TODO: All")
}

// Trace возвращает traceability-текст сборки b: lang "go" или "ini";
// другой lang — ошибка.
func (c *CI) Trace(b jobs.Build, lang string) (string, error) {
	panic("TODO: Trace")
}

// Retain оставляет последние n сборок: логи и trace-файлы выпавших
// сборок удаляются. n < 1 — ошибка.
func (c *CI) Retain(n int) error {
	panic("TODO: Retain")
}
