//go:build !solution

// Пакет random — глава 9. Реализуй функции так, чтобы прошли тесты
// в random_test.go. Условие — в README.md этой папки.
package random

// Uniform возвращает n чисел, равномерно распределённых на [lo, hi),
// из детерминированного источника ChaCha8 с семенем seed:
// одинаковое seed — в точности тот же слайс. n < 0 трактуется как 0.
func Uniform(seed uint64, lo, hi float64, n int) []float64 {
	panic("TODO: Uniform")
}

// Normal возвращает n чисел с нормальным распределением N(mean, std)
// из детерминированного источника ChaCha8 с семенем seed:
// одинаковое seed — в точности тот же слайс. n < 0 трактуется как 0.
func Normal(seed uint64, mean, std float64, n int) []float64 {
	panic("TODO: Normal")
}

// CryptoUniform возвращает n чисел, равномерно распределённых на
// [lo, hi), из криптографического источника crypto/rand (без семени).
// Ошибка источника возвращается наружу. n < 0 трактуется как 0.
func CryptoUniform(lo, hi float64, n int) ([]float64, error) {
	panic("TODO: CryptoUniform")
}
