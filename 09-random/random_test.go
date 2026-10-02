package random

import (
	"math"
	"reflect"
	"testing"
)

func TestUniformDeterministic(t *testing.T) {
	a := Uniform(1729, 0, 1, 100)
	b := Uniform(1729, 0, 1, 100)
	if len(a) != 100 {
		t.Fatalf("Uniform: длина %d, хотел 100", len(a))
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("Uniform: одинаковое семя должно давать в точности ту же серию")
	}
	if c := Uniform(1728, 0, 1, 100); c[0] == a[0] && c[99] == a[99] {
		t.Error("Uniform: разные семена дали одинаковые серии")
	}
}

func TestUniformBounds(t *testing.T) {
	for _, x := range Uniform(42, -5, 5, 1000) {
		if x < -5 || x >= 5 {
			t.Fatalf("Uniform: значение %v вне [-5, 5)", x)
		}
	}
}

func TestUniformBuckets(t *testing.T) {
	// n=10000 в 10 бакетах: пустой бакет — провал качества генератора.
	const n, buckets = 10000, 10
	var counts [buckets]int
	for _, x := range Uniform(1729, 0, 1, n) {
		if x < 0 || x >= 1 {
			t.Fatalf("Uniform: значение %v вне [0, 1)", x)
		}
		counts[int(x*buckets)]++
	}
	for i, c := range counts {
		if c == 0 {
			t.Errorf("Uniform: бакет %d пуст (n=%d, качество плохое)", i, n)
		}
	}
}

func TestNormalDeterministic(t *testing.T) {
	a := Normal(1729, 10, 2, 100)
	b := Normal(1729, 10, 2, 100)
	if len(a) != 100 {
		t.Fatalf("Normal: длина %d, хотел 100", len(a))
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("Normal: одинаковое семя должно давать в точности ту же серию")
	}
}

func TestNormalMean(t *testing.T) {
	// Среднее N(10, 2) из 10000 выборок: допуск 3 сигмы = 3*std/√n.
	const mean, std, n = 10.0, 2.0, 10000
	sum := 0.0
	for _, x := range Normal(1729, mean, std, n) {
		sum += x
	}
	tol := 3 * std / math.Sqrt(n)
	if diff := sum/n - mean; diff > tol || diff < -tol {
		t.Errorf("Normal: среднее %.4f, хотел %.4f ± %.4f", sum/n, mean, tol)
	}
}

func TestCryptoUniform(t *testing.T) {
	xs, err := CryptoUniform(0, 100, 100)
	if err != nil {
		t.Fatalf("CryptoUniform: %v", err)
	}
	if len(xs) != 100 {
		t.Fatalf("CryptoUniform: длина %d, хотел 100", len(xs))
	}
	for i, x := range xs {
		if x < 0 || x >= 100 {
			t.Fatalf("CryptoUniform: элемент %d = %v вне [0, 100)", i, x)
		}
	}
}
