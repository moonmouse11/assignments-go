//go:build !solution

// Пакет numbers — глава 8. Реализуй функции так, чтобы прошли тесты
// в numbers_test.go. Условие — в README.md этой папки.
package numbers

// Digits возвращает десятичные цифры n по возрастанию разрядов:
// Digits(40585) = [4 0 5 8 5]. 0 → [0]; отрицательное — по модулю.
func Digits(n int) []int {
	panic("TODO: Digits")
}

// FromDigits возвращает число, составленное из цифр d (пустой список —
// 0 без ошибки). Ошибка, если в d есть не-цифры (вне 0..9).
func FromDigits(d []int) (int, error) {
	panic("TODO: FromDigits")
}

// Divisors возвращает все делители n по возрастанию: Divisors(12) =
// [1 2 3 4 6 12]. Для n < 1 — nil.
func Divisors(n int) []int {
	panic("TODO: Divisors")
}

// PrimeFactors возвращает простые множители n с повторами,
// по возрастанию: PrimeFactors(360) = [2 2 2 3 3 5]. Для n < 2 — nil.
func PrimeFactors(n int) []int {
	panic("TODO: PrimeFactors")
}

// Sum возвращает сумму чисел xs (пустой слайс — 0).
func Sum(xs []int) int {
	panic("TODO: Sum")
}

// Product возвращает произведение чисел xs (пустой слайс — 1).
func Product(xs []int) int {
	panic("TODO: Product")
}

// GCD возвращает наибольший общий делитель a и b (результат ≥ 0;
// GCD(0, 0) = 0).
func GCD(a, b int) int {
	panic("TODO: GCD")
}

// Fibonacci возвращает первые n чисел Фибоначчи:
// Fibonacci(8) = [0 1 1 2 3 5 8 13]. Для n < 1 — пустой слайс.
func Fibonacci(n int) []int {
	panic("TODO: Fibonacci")
}

// Powers возвращает e-ю степень каждого элемента xs:
// Powers([1 2 3], 2) = [1 4 9]. e < 0 трактуется как 0 (все — 1).
func Powers(xs []int, e int) []int {
	panic("TODO: Powers")
}
