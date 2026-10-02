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
