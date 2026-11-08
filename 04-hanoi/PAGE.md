# Глава 4. Hanoi — рекурсия против итерации

> Размер задачи определяет цену решения:
> сравниваем два алгоритма одной и той же задачи инструментами Go.

## Цель главы

Рекурсия и итерация дают **один и тот же** результат — докажи это тестом.
Затем измерь цену: `go test -bench` — сравнение производительности
алгоритмов, как того требует исходная спецификация. Плюс важная идея:
алгоритм отделён от визуализации — ходы это просто данные `[]Move`.

## Теория

- Ханойские башни: n дисков, 3 колышка (0, 1, 2), перенос 0 → 2,
  диск кладётся только на больший. Минимум ходов: 2ⁿ − 1.
- Рекурсия: перенести n−1 на «запасной», большой диск на цель, n−1 сверху.
- Итерация (правило чётности): нечётное n — самый маленький диск ходит
  по кругу 0→2→1; чётное — 0→1→2; между его ходами — единственный законный
  ход остальных двух колышков.
- `testing.B`: `func BenchmarkX(b *testing.B)` — `b.N` повторений.
- `reflect.DeepEqual` — сравнение слайсов в тестах.

## Условие

1. `SolveRecursive(n)` — решение рекурсией: слайс ходов `[]Move`
   переноса n дисков с колышка 0 на колышек 2.
2. `SolveIterative(n)` — та же задача итерацией. Результат — **в точности
   такой же слайс**, как у рекурсии.
3. `Verify(moves, n)` — симуляция: проверить ходы на законность
   (ход с пустого колышка нельзя; больше на меньший нельзя; колышки 0..2)
   и что в конце все диски на колышке 2. Ошибка — `error` с описанием хода.
4. Число ходов — ровно 2ⁿ − 1 (проверяется тестом, но и ты это докажи).

## API для реализации

```go
type Move struct{ From, To int }
func SolveRecursive(n int) []Move
func SolveIterative(n int) []Move
func Verify(moves []Move, n int) error
```

## Подсказки

- Колышки в симуляции — стеки: `[][]int` (срез-стек, диск 1 — самый маленький).
- В итерации удобно хранить, где сейчас диск 1, и шаг цикла:
  нечётное n → +2 mod 3, чётное → +1 mod 3.
- Ошибка `fmt.Errorf("ход %d: ...", i, ...)` — с номером хода.

## Проверка

```
go test ./04-hanoi/
go test -tags solution ./04-hanoi/
go test -bench . ./04-hanoi        # рекурсия vs итерация
```

## Эталон

`go test -tags solution ./04-hanoi/`

## Углубление

- Профили памяти и CPU на решение (`go test -memprofile`, `-cpuprofile`):
  глубина рекурсии ~ n — прокси стека; итерация — глубина 0, но растут
  внутренние счётчики.
- Таблица итераций из спеки: диски 16/20/24/28 → итерации
  21845/349525/5592405/89478485 — воспроизведи счётчик итераций
  в `SolveIterative` и проверь формулу.
- Визуализация ходов с паузой/возобновлением (fyne/gtk3 у автора) —
  `[]Move` уже отделяет алгоритм от представления.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./04-hanoi/
```

## Заготовки (реализуйте)

### `hanoi.go`

{% raw %}
```go
//go:build !solution

// Пакет hanoi — глава 4. Реализуй так, чтобы прошли тесты.
package hanoi

import "fmt"

// Move — ход: перенос одного диска с колышка From на колышек To (0..2).
type Move struct {
	From, To int
}

// SolveRecursive возвращает последовательность ходов переноса n дисков
// с колышка 0 на колышек 2 (рекурсивное решение).
func SolveRecursive(n int) []Move {
	panic("TODO: SolveRecursive")
}

// SolveIterative решает ту же задачу итеративно; последовательность ходов
// совпадает с SolveRecursive(n) в точности.
func SolveIterative(n int) []Move {
	panic("TODO: SolveIterative")
}

// Verify симулирует ходы: каждый ход законен (с непустого колышка,
// меньший диск на больший, колышки в 0..2) и в конце все n дисков
// на колышке 2. Незаконный ход — error с номером хода.
func Verify(moves []Move, n int) error {
	_ = fmt.Sprintf // убрать, когда появится настоящий код
	panic("TODO: Verify")
}
```
{% endraw %}

## Тесты

### `hanoi_test.go`

{% raw %}
```go
package hanoi

import (
	"reflect"
	"testing"
)

func TestSolutionsAgree(t *testing.T) {
	for n := 1; n <= 10; n++ {
		r := SolveRecursive(n)
		i := SolveIterative(n)
		if !reflect.DeepEqual(r, i) {
			t.Fatalf("n=%d: итеративное решение отличается от рекурсивного", n)
		}
		if len(r) != (1<<n)-1 {
			t.Errorf("n=%d: ходов %d, хотел %d", n, len(r), (1<<n)-1)
		}
		if err := Verify(r, n); err != nil {
			t.Errorf("n=%d: Verify: %v", n, err)
		}
	}
}

func TestVerifyRejectsBadMoves(t *testing.T) {
	cases := []struct {
		name  string
		moves []Move
		n     int
	}{
		{"ход с пустого колышка", []Move{{1, 2}}, 1},
		{"колышек вне диапазона", []Move{{0, 5}}, 1},
		{"задача не завершена", []Move{{0, 1}}, 2},
		{"больше на меньший", []Move{{0, 1}, {0, 1}}, 2},
	}
	for _, tc := range cases {
		if err := Verify(tc.moves, tc.n); err == nil {
			t.Errorf("%s: Verify должен вернуть ошибку", tc.name)
		}
	}
}

func TestVerifyAcceptsKnownGood(t *testing.T) {
	// n=3, каноническое решение.
	moves := []Move{
		{0, 2}, {0, 1}, {2, 1}, {0, 2}, {1, 0}, {1, 2}, {0, 2},
	}
	if err := Verify(moves, 3); err != nil {
		t.Errorf("Verify(каноническое решение n=3): %v", err)
	}
}

// Сравнение производительности — как требует исходная спецификация Hanoi.pdf:
//
//	go test -bench . ./04-hanoi
func BenchmarkRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SolveRecursive(20)
	}
}

func BenchmarkIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SolveIterative(20)
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

// Эталонное решение главы 4.
package hanoi

import "fmt"

type Move struct {
	From, To int
}

func SolveRecursive(n int) []Move {
	var moves []Move
	var rec func(n, from, to, via int)
	rec = func(n, from, to, via int) {
		if n == 0 {
			return
		}
		rec(n-1, from, via, to)
		moves = append(moves, Move{from, to})
		rec(n-1, via, to, from)
	}
	rec(n, 0, 2, 1)
	return moves
}

func SolveIterative(n int) []Move {
	pegs := make([][]int, 3)
	for d := n; d >= 1; d-- {
		pegs[0] = append(pegs[0], d)
	}
	top := func(p int) int {
		if len(pegs[p]) == 0 {
			return 1 << 30
		}
		return pegs[p][len(pegs[p])-1]
	}
	pop := func(p int) int {
		d := pegs[p][len(pegs[p])-1]
		pegs[p] = pegs[p][:len(pegs[p])-1]
		return d
	}

	// Диск 1 ходит по кругу: нечётное n — 0→2→1, чётное — 0→1→2.
	step := 2
	if n%2 == 0 {
		step = 1
	}
	p1 := 0 // где сейчас диск 1
	total := (1 << n) - 1
	moves := make([]Move, 0, total)
	for len(moves) < total {
		var mv Move
		if len(moves)%2 == 0 {
			mv = Move{p1, (p1 + step) % 3}
			p1 = mv.To
		} else {
			// Единственный законный ход между двумя колышками без диска 1:
			// меньший верхний — на больший.
			var others [2]int
			k := 0
			for p := 0; p < 3; p++ {
				if p != p1 {
					others[k] = p
					k++
				}
			}
			a, b := others[0], others[1]
			if top(a) < top(b) {
				mv = Move{a, b}
			} else {
				mv = Move{b, a}
			}
		}
		pegs[mv.To] = append(pegs[mv.To], pop(mv.From))
		moves = append(moves, mv)
	}
	return moves
}

func Verify(moves []Move, n int) error {
	pegs := make([][]int, 3)
	for d := n; d >= 1; d-- {
		pegs[0] = append(pegs[0], d)
	}
	for i, m := range moves {
		if m.From < 0 || m.From > 2 || m.To < 0 || m.To > 2 {
			return fmt.Errorf("ход %d: колышек %d→%d вне 0..2", i, m.From, m.To)
		}
		if len(pegs[m.From]) == 0 {
			return fmt.Errorf("ход %d: колышек %d пуст", i, m.From)
		}
		d := pegs[m.From][len(pegs[m.From])-1]
		if dst := pegs[m.To]; len(dst) > 0 && dst[len(dst)-1] < d {
			return fmt.Errorf("ход %d: диск %d на меньший %d", i, d, dst[len(dst)-1])
		}
		pegs[m.From] = pegs[m.From][:len(pegs[m.From])-1]
		pegs[m.To] = append(pegs[m.To], d)
	}
	if len(pegs[2]) != n {
		return fmt.Errorf("в конце на колышке 2 дисков %d, хотел %d", len(pegs[2]), n)
	}
	return nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/04-hanoi`, упражнение **hanoi**. Локально: `go test -tags solution ./04-hanoi/`._
