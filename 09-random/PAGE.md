# Глава 9. Random — качество генераторов

> Изучаем генераторы случайных чисел:
> детерминированный (семя → та же серия) и криптографический,
> качество проверяем статистикой.

## Цель главы

Сгенерировать серии чисел и убедиться, что они «хорошие»: одинаковое
семя → в точности та же серия (воспроизводимость при разработке);
значения в границах; среднее Normal — около матожидания; равномерность —
все бакеты гистограммы непустые. В оригинале серии рисуют графиками —
здесь проверяем числами.

## Теория

- `math/rand/v2`: `rand.New(rand.NewChaCha8(seed [32]byte))` — источник
  ChaCha8; `r.Float64()` — [0, 1), `r.NormFloat64()` — N(0, 1).
  Семя `uint64` укладывается в 32 байта, например, через
  `binary.LittleEndian.PutUint64` — остальное нули.
- `crypto/rand.Read(buf)` — криптографический источник: без семени,
  каждая серия другая; ошибка чтения передаётся наружу.
- Из [0, 1) в [lo, hi): `lo + u*(hi-lo)`.
- Закон больших чисел: среднее n выборок из N(mean, std) лежит около
  mean с разбросом std/√n — допуск 3·std/√n почти неотказуем.
- Гистограмма: n=10000 в 10 бакетах; пустой бакет — провал качества,
  как «дыры» в гистограммах `random.pdf`.

## Условие

1. `Uniform(seed, lo, hi, n)` — n чисел, равномерных на [lo, hi),
   из ChaCha8 с семенем seed. Одно семя — в точности тот же слайс.
2. `Normal(seed, mean, std, n)` — n чисел N(mean, std) из того же
   детерминированного источника.
3. `CryptoUniform(lo, hi, n)` — n чисел, равномерных на [lo, hi),
   из криптографического источника; ошибка источника — наружу.

Wants (необязательно, из `random.pdf`): гистограмма и series-плот в PNG,
другие распределения (Вейбулл), качество при росте n.

## API для реализации

```go
func Uniform(seed uint64, lo, hi float64, n int) []float64
func Normal(seed uint64, mean, std float64, n int) []float64
func CryptoUniform(lo, hi float64, n int) ([]float64, error)
```

Сигнатуры и договоры — в `random.go`. Тесты — в `random_test.go`.

## Подсказки

- n < 0 трактуется как 0 — пустой слайс, без паники `make`.
- Для CryptoUniform: 8 байт из `rand.Read` → `Uint64` → сдвиг вправо
  на 11 и деление на 2⁵³ — равномерное [0, 1).
- Среднее считай циклом — внешних библиотек нет.
- Тесты детерминированы: seed зафиксирован, поэтому проверки бакетов
  и среднего стабильны от запуска к запуску.

## Проверка

```
go test ./09-random/
go test -tags solution ./09-random/
```

## Эталон

`go test -tags solution ./09-random/`.

## Углубление

- CLI спеки: подкоманды uniform/normal с позиционными аргументами
  (min max / mean std, умолчания 0.0 1.0) и переключателями `-C/--crypto`,
  `-s/--samples` (1000), `-c/--slices` (10), `-d/--seed` (1729),
  `-o/--output`, `--series`, `-t/--table`; плюс команды help и version.
- `--table <файл>`: табличный дамп значений в файл в дополнение к графику
  (в примере — `normal -t table --samples 1001`).
- Seed 1729 — таксикаб-число Рамануджана (перекличка с 08-numbers);
  линейные конгруэнтные генераторы повторяют серию при том же семени —
  потому гистограммы примера воспроизводимы.
- PNG-графики wants: столбики можно нарисовать stdlib (`image`,
  `image/png`); готовые серии и гистограммы — gonum/plot (внешняя
  зависимость). Эталон автора — репозиторий на GitLab (сеть).

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./09-random/
```

## Заготовки (реализуйте)

### `random.go`

{% raw %}
```go
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
```
{% endraw %}

## Тесты

### `random_test.go`

{% raw %}
```go
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
```
{% endraw %}

## Эталонное решение

Сначала попробуйте сами. Решение скрыто — раскрывайте только для сверки.

<details>
<summary><strong>Эталон (`solution.go`)</strong> — нажмите, чтобы показать</summary>

{% raw %}
```go
//go:build solution

// Эталонное решение главы 9. Сравни со своей реализацией в random.go.
package random

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand/v2"
)

// chacha — детерминированный источник: 8 байт семени, остальное нули.
func chacha(seed uint64) *rand.Rand {
	var b [32]byte
	binary.LittleEndian.PutUint64(b[:8], seed)
	return rand.New(rand.NewChaCha8(b))
}

func Uniform(seed uint64, lo, hi float64, n int) []float64 {
	r := chacha(seed)
	out := make([]float64, max(n, 0))
	for i := range out {
		out[i] = lo + r.Float64()*(hi-lo)
	}
	return out
}

func Normal(seed uint64, mean, std float64, n int) []float64 {
	r := chacha(seed)
	out := make([]float64, max(n, 0))
	for i := range out {
		out[i] = mean + std*r.NormFloat64()
	}
	return out
}

func CryptoUniform(lo, hi float64, n int) ([]float64, error) {
	out := make([]float64, max(n, 0))
	var b [8]byte
	for i := range out {
		if _, err := crand.Read(b[:]); err != nil {
			return nil, err
		}
		u := float64(binary.LittleEndian.Uint64(b[:])>>11) / (1 << 53)
		out[i] = lo + u*(hi-lo)
	}
	return out, nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/09-random`, упражнение **random**. Локально: `go test -tags solution ./09-random/`._
