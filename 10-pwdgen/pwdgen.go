//go:build !solution

// Пакет pwdgen — глава 10. Реализуй функции так, чтобы прошли тесты
// в pwdgen_test.go. Условие — в README.md этой папки.
package pwdgen

import "math/rand/v2"

// Gen — генератор паролей: словарь слов, разделитель сегментов и
// детерминированный источник ChaCha8. Создаётся через NewGen.
type Gen struct {
	Words []string
	Sep   string
	rng   *rand.Rand
}

// NewGen возвращает генератор с семенем seed, словарём words и
// разделителем sep. Одинаковое семя и словарь — одинаковые пароли.
// Паникует, если words пуст (нет слов — нет пароля).
func NewGen(seed uint64, words []string, sep string) *Gen {
	panic("TODO: NewGen")
}

// Segment возвращает один сегмент пароля: случайное слово из словаря,
// склеенное через Sep со случайным числом из 2–4 цифр (без ведущих
// нулей). Например, при Sep "-" — "star-29660".
func (g *Gen) Segment() string {
	panic("TODO: Segment")
}

// Password возвращает nSeg сегментов, склеенных через Sep:
// "vein-47519-Bone-52061". nSeg <= 0 — пустая строка.
func (g *Gen) Password(nSeg int) string {
	panic("TODO: Password")
}
