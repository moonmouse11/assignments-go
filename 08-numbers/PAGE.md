# Глава 8. Numbers — библиотека теории чисел

> Развлекательная теория чисел как повод
> построить слоистую библиотеку чистых функций и табличные тесты.

## Цель главы

Спроектировать «ядро» библиотеки: маленькие чистые функции без
состояния — цифры, делители, простые множители, НОД, Фибоначчи. Слои
выше (IsPrime, Лун, ISBN — см. `numbers.pdf`) строятся на примитивах.

## Теория

- `n%10` и `n/10` — последняя цифра и число без неё.
- Делители ищутся до `i*i <= n` парами `i` и `n/i`; дубли середины
  у квадратов убирает `slices.Compact` после `slices.Sort`.
- Разложение на множители: пробное деление на растущие `p`, пока
  `p*p <= n`; остаток больше 1 — сам простой.
- Алгоритм Евклида: `for b != 0 { a, b = b, a%b }` — НОД.
- Табличные тесты: `[]struct{ in, want … }` и цикл с `t.Errorf`,
  где видны оба значения — получил и хотел.

## Условие

1. `Digits(n)` — десятичные цифры n по порядку (0 → `[0]`; отрицательное — по модулю).
2. `FromDigits(d)` — число из цифр; ошибка, если есть не-цифры (вне 0..9); пустой список — 0 без ошибки.
3. `Divisors(n)` — все делители по возрастанию (n < 1 → nil).
4. `PrimeFactors(n)` — простые множители с повторами, по возрастанию (n < 2 → nil).
5. `Sum(xs)` / `Product(xs)` — сумма и произведение (пустой слайс: 0 и 1).
6. `GCD(a, b)` — наибольший общий делитель (результат ≥ 0; `GCD(0,0)=0`).
7. `Fibonacci(n)` — первые n чисел: `0, 1, 1, 2, 3, 5, …` (n < 1 — пустой слайс).
8. `Powers(xs, e)` — e-я степень каждого (e < 0 трактуется как 0).

Wants (необязательно, из `numbers.pdf`): IsPrime/IsPerfect/IsHarshad, попарные НОД Фибоначчи, taxicab-числа.

## API для реализации

```go
func Digits(n int) []int
func FromDigits(d []int) (int, error)
func Divisors(n int) []int
func PrimeFactors(n int) []int
func Sum(xs []int) int
func Product(xs []int) int
func GCD(a, b int) int
func Fibonacci(n int) []int
func Powers(xs []int, e int) []int
```

Сигнатуры и договоры — в `numbers.go`. Тесты — в `numbers_test.go`.

## Подсказки

- Раунд-трип: `FromDigits(Digits(n)) == n` для n ≥ 0.
- `Product(PrimeFactors(360)) == 360` — проверка разложения.
- 28 — совершенное число: `Sum(Divisors(28)) - 28 == 28`.
- `for range e` — e повторов внутреннего цикла (Go 1.22+).

## Проверка

```
go test ./08-numbers/
go test -tags solution ./08-numbers/
```

## Эталон

`go test -tags solution ./08-numbers/`.

## Углубление

- Applications спеки: Kaprekar-числа и Happy-числа (`IsHappy True` для
  4096 в примере) — их ещё нет даже в wants главы.
- Needs 13–14: квадраты, представимые суммой двух квадратов (найти
  максимум представлений), то же для кубов; пары чисел, чьё произведение —
  тоже сумма двух квадратов.
- Детали Лун и ISBN: тип карты (Visa/Mastercard/Amex) по префиксу и длине;
  перевод ISBN-10 ↔ ISBN-13, контрольная цифра, буква X в ISBN-10.
- Производительность — стержень спеки: timing парных НОД на 1024 и 4096
  элементах повтори бенчмарками (`go test -bench`).
- Wants: графики плотности простых и Фибоначчи — ASCII-гистограмма на
  stdlib или gonum/plot (внешняя зависимость).

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./08-numbers/
```

## Заготовки (реализуйте)

### `numbers.go`

{% raw %}
```go
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
```
{% endraw %}

## Тесты

### `numbers_test.go`

{% raw %}
```go
package numbers

import (
	"reflect"
	"testing"
)

func TestDigits(t *testing.T) {
	cases := []struct {
		in   int
		want []int
	}{
		{40585, []int{4, 0, 5, 8, 5}},
		{4096, []int{4, 0, 9, 6}},
		{7, []int{7}},
		{0, []int{0}},
		{-12, []int{1, 2}},
	}
	for _, tc := range cases {
		if got := Digits(tc.in); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Digits(%d) = %v, хотел %v", tc.in, got, tc.want)
		}
	}
}

func TestFromDigits(t *testing.T) {
	if got, err := FromDigits([]int{4, 0, 5, 8, 5}); err != nil || got != 40585 {
		t.Errorf("FromDigits([4 0 5 8 5]) = (%d, %v), хотел (40585, nil)", got, err)
	}
	for _, bad := range [][]int{{4, 12}, {10, -1}} {
		if _, err := FromDigits(bad); err == nil {
			t.Errorf("FromDigits(%v): хотел ошибку, получил nil", bad)
		}
	}
	if got, err := FromDigits(nil); err != nil || got != 0 {
		t.Errorf("FromDigits(nil) = (%d, %v), хотел (0, nil)", got, err)
	}
	// Раунд-трип: цифры → число.
	for _, n := range []int{0, 5, 42, 907, 40585} {
		if got, err := FromDigits(Digits(n)); err != nil || got != n {
			t.Errorf("FromDigits(Digits(%d)) = (%d, %v), хотел (%d, nil)", n, got, err, n)
		}
	}
}

func TestDivisors(t *testing.T) {
	cases := []struct {
		in   int
		want []int
	}{
		{12, []int{1, 2, 3, 4, 6, 12}},
		{1, []int{1}},
		{7, []int{1, 7}},
		{36, []int{1, 2, 3, 4, 6, 9, 12, 18, 36}},
		{28, []int{1, 2, 4, 7, 14, 28}},
		{0, nil},
		{-4, nil},
	}
	for _, tc := range cases {
		if got := Divisors(tc.in); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Divisors(%d) = %v, хотел %v", tc.in, got, tc.want)
		}
	}
}

func TestPrimeFactors(t *testing.T) {
	cases := []struct {
		in   int
		want []int
	}{
		{360, []int{2, 2, 2, 3, 3, 5}},
		{4096, []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}},
		{13, []int{13}},
		{1, nil},
		{0, nil},
	}
	for _, tc := range cases {
		if got := PrimeFactors(tc.in); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("PrimeFactors(%d) = %v, хотел %v", tc.in, got, tc.want)
		}
	}
}

func TestSumProduct(t *testing.T) {
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Errorf("Sum([1 2 3]) = %d, хотел 6", got)
	}
	if got := Sum(nil); got != 0 {
		t.Errorf("Sum(nil) = %d, хотел 0", got)
	}
	if got := Product([]int{2, 2, 5}); got != 20 {
		t.Errorf("Product([2 2 5]) = %d, хотел 20", got)
	}
	if got := Product(nil); got != 1 {
		t.Errorf("Product(nil) = %d, хотел 1", got)
	}
	if got := Product([]int{0, 3}); got != 0 {
		t.Errorf("Product([0 3]) = %d, хотел 0", got)
	}
}

func TestGCD(t *testing.T) {
	cases := []struct {
		a, b, want int
	}{
		{48, 18, 6},
		{18, 48, 6},
		{7, 13, 1},
		{270, 192, 6},
		{0, 5, 5},
		{5, 0, 5},
		{0, 0, 0},
		{-48, 18, 6},
	}
	for _, tc := range cases {
		if got := GCD(tc.a, tc.b); got != tc.want {
			t.Errorf("GCD(%d, %d) = %d, хотел %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestFibonacci(t *testing.T) {
	if got := Fibonacci(8); !reflect.DeepEqual(got, []int{0, 1, 1, 2, 3, 5, 8, 13}) {
		t.Errorf("Fibonacci(8) = %v, хотел [0 1 1 2 3 5 8 13]", got)
	}
	if got := Fibonacci(1); !reflect.DeepEqual(got, []int{0}) {
		t.Errorf("Fibonacci(1) = %v, хотел [0]", got)
	}
	if got := Fibonacci(0); len(got) != 0 {
		t.Errorf("Fibonacci(0) = %v, хотел пустой слайс", got)
	}
}

func TestPowers(t *testing.T) {
	if got := Powers([]int{1, 2, 3, 4}, 2); !reflect.DeepEqual(got, []int{1, 4, 9, 16}) {
		t.Errorf("Powers([1 2 3 4], 2) = %v, хотел [1 4 9 16]", got)
	}
	if got := Powers([]int{2, 3}, 3); !reflect.DeepEqual(got, []int{8, 27}) {
		t.Errorf("Powers([2 3], 3) = %v, хотел [8 27]", got)
	}
	if got := Powers([]int{5, 9}, 0); !reflect.DeepEqual(got, []int{1, 1}) {
		t.Errorf("Powers([5 9], 0) = %v, хотел [1 1]", got)
	}
}

func TestComposites(t *testing.T) {
	// Слои ядра складываются в проверки из numbers.pdf.
	if s := Sum(Divisors(28)) - 28; s != 28 {
		t.Errorf("28: сумма собственных делителей = %d, хотел 28 (совершенное число)", s)
	}
	if p := Product(PrimeFactors(360)); p != 360 {
		t.Errorf("произведение простых множителей 360 = %d, хотел 360", p)
	}
	if g := GCD(Sum([]int{48}), Product([]int{6, 3})); g != 6 {
		t.Errorf("GCD(48, 18) через Sum/Product = %d, хотел 6", g)
	}
}
```
{% endraw %}

## Эталонное решение

Сначала попробуйте сами. Решение скрыто — раскрывайте только для сверки.

<details>
<summary><strong>Эталон (`solution.go`)</strong> — нажмите, чтобы показать</summary>

{% raw %}
```go
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
```
{% endraw %}


</details>

_Пакет `assignments-go/08-numbers`, упражнение **numbers**. Локально: `go test -tags solution ./08-numbers/`._
