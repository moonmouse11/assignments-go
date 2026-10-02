//go:build solution

// Эталонное решение главы 8. Сравни со своей реализацией в numbers.go.
package numbers

import "fmt"
import "slices"

func Digits(n int) []int {
	if n = max(n, -n); n < 10 {
		return []int{n}
	}
	return append(Digits(n/10), n%10)
}

func FromDigits(d []int) (int, error) {
	for _, x := range d {
		if x < 0 || x > 9 {
			return 0, fmt.Errorf("numbers: %d — не цифра", x)
		}
	}
	return fold(d, 0, func(acc, x int) int { return acc*10 + x }), nil
}

func Divisors(n int) (ds []int) {
	for i := 1; i*i <= n; i++ {
		if n%i == 0 {
			ds = append(ds, i, n/i)
		}
	}
	return slices.Compact(slices.Sorted(slices.Values(ds)))
}

func PrimeFactors(n int) (f []int) {
	for p := 2; p*p <= n; p++ {
		for n%p == 0 {
			f, n = append(f, p), n/p
		}
	}
	if n > 1 {
		f = append(f, n)
	}
	return
}

func fold(xs []int, acc int, op func(acc, x int) int) int {
	for _, x := range xs {
		acc = op(acc, x)
	}
	return acc
}

func Sum(xs []int) int     { return fold(xs, 0, func(acc, x int) int { return acc + x }) }
func Product(xs []int) int { return fold(xs, 1, func(acc, x int) int { return acc * x }) }

func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return max(a, -a)
}

func Fibonacci(n int) []int {
	f, a, b := make([]int, max(n, 0)), 0, 1
	for i := range f {
		f[i], a, b = a, b, a+b
	}
	return f
}

func Powers(xs []int, e int) []int {
	out := make([]int, len(xs))
	for i, x := range xs {
		out[i] = 1
		for range e {
			out[i] *= x
		}
	}
	return out
}
