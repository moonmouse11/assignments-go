//go:build !solution

// Пакет flash — глава 24. Реализуй Items, RecallScore и Rotate так,
// чтобы прошли тесты в flash_test.go. Условие — в README.md этой папки.
package flash

import (
	"context"
	"time"
)

// Item — одна карточка показа: слово (Kind 'w') или число (Kind 'n').
// Значение — всегда строка: "бета" или "7".
type Item struct {
	Kind  byte   // 'w' — слово | 'n' — число
	Value string // текст значения
}

// Items собирает набор для показа: все слова words плюс nNums случайных
// чисел, порядок перемешан. Одинаковый seed — одинаковый набор.
func Items(words []string, nNums int, seed uint64) []Item {
	panic("TODO: Items")
}

// RecallScore считает, сколько карточек из recalled припомнено верно
// относительно shown: совпадают и значение, и вид; порядок неважен;
// каждая карточка shown засчитывается не больше одного раза.
func RecallScore(shown, recalled []Item) int {
	panic("TODO: RecallScore")
}

// Rotate каждые interval «отбрасывает» (вызывает onDrop) очередную карточку
// items слева направо. Первый отброс — через interval после старта.
// Возвращает число отброшенных до конца набора или отмены ctx.
func Rotate(ctx context.Context, interval time.Duration, items []Item, onDrop func(Item)) int {
	panic("TODO: Rotate")
}
