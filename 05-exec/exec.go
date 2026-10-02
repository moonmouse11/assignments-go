//go:build !solution

// Пакет exec — глава 5 (упрощённый projectlet из exec.pdf).
// Реализуй RunCommand и Periodic так, чтобы прошли тесты.
package exec

import (
	"context"
	"time"
)

// RunCommand выполняет команду name с аргументами args и возвращает
// её вывод (stdout/stderr) или ошибку запуска.
func RunCommand(name string, args ...string) ([]byte, error) {
	panic("TODO: RunCommand")
}

// Periodic запускает fn каждые interval, начиная через interval после
// старта, пока ctx не отменят. Возвращает число исполнений fn.
func Periodic(ctx context.Context, interval time.Duration, fn func()) int {
	panic("TODO: Periodic")
}
