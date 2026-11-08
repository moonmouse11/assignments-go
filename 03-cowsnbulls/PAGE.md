# Глава 3. Cows and Bulls — игра и решатель

> Классическая игра
> ([Bulls and Cows](https://en.wikipedia.org/wiki/Bulls_and_Cows)):
> обе стороны монеты — и дилер, и игрок.

## Цель главы

Строки, слайсы, map, алгоритм фильтрации кандидатов. Первая задача, где
компьютер **угадывает** сам: перебор пространства кандидатов с отсечением
по ответу — ядро многих алгоритмов (и то же самое потом в `jotto.pdf`).

## Теория

- Bulls — быки: цифра на своём месте. Cows — коровы: цифра есть, но не там.
- `map[byte]bool` — множество для быстрого «есть такая цифра?».
- Фильтрация на месте: пишем результат в тот же слайс
  `out := s[:0]` — ноль новых аллокаций.
- Пространство: 4-значные числа с неповторяющимися цифрами — 10·9·8·7 = 5040
  кандидатов.

## Условие

Игра ведётся 4-значными числами с **неповторяющимися** цифрами.

### Часть A — счёт (режим дилера)

`Score(secret, guess)` возвращает `(bulls, cows)` — оценку попытки `guess`
против секрета `secret`. Гарантируется: длины равны 4, цифры не повторяются
(в `secret` — обязательно, в `guess` — тоже, договоримся для простоты).

### Часть B — решатель (режим игрока)

`GuessingSolver` — компьютер угадывает число, задуманное человеком:

1. `NewGuessingSolver()` — все 5040 кандидатов.
2. `Guess()` — текущая гипотеза (детерминированная, например первая).
3. `Learn(guess, bulls, cows)` — человек называет счёт, решатель вычёркивает
   всех кандидатов, дающих другой счёт.
4. Решатель обязан находить секрет не более чем за 25 попыток
   (гарантированно хватит меньше — проверь экспериментом).

## API для реализации

```go
func Score(secret, guess string) (bulls, cows int)
type GuessingSolver struct{ /* ... */ }
func NewGuessingSolver() *GuessingSolver
func (g *GuessingSolver) Guess() string
func (g *GuessingSolver) Learn(guess string, bulls, cows int)
```

## Подсказки

- Сначала `Score`, она нужна решателю для самопроверки кандидатов.
- Генерация кандидатов: 4 вложенных цикла с пропуском повторов.
- `Learn` можно писать так: оставить `cand`, если
  `Score(cand, guess) == (bulls, cows)`.

## Проверка

```
go test ./03-cowsnbulls/
go test -tags solution ./03-cowsnbulls/
```

## Эталон

`go test -tags solution ./03-cowsnbulls/`

## Углубление

- Dealer-режим в CLI: компьютер загадывает (`crypto/rand`), человек
  угадывает, счёт после каждой попытки.
- Спека требует инструментировать аллокации (автор — gnatmem/valgrind):
  в Go — `go test -memprofile=mem.out`, бенчмарк `Learn` на полном
  пространстве кандидатов.
- Теория из спеки: решатель обязан угадывать за ≤5 попыток — проверь
  худший случай статистикой по всем 5040 секретам.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./03-cowsnbulls/
```

## Заготовки (реализуйте)

### `cowsnbulls.go`

{% raw %}
```go
//go:build !solution

// Пакет cowsnbulls — глава 3. Реализуй Score и решатель так,
// чтобы прошли тесты. Условие — в README.md.
package cowsnbulls

// Score оценивает попытку guess против секрета secret
// (4 цифры, без повторов): bulls — цифра на своём месте,
// cows — правильная цифра не на своём месте.
func Score(secret, guess string) (bulls, cows int) {
	panic("TODO: Score")
}

// GuessingSolver — компьютер угадывает задуманное число,
// вычёркивая кандидатов по ответам человека.
type GuessingSolver struct {
	// TODO: твои поля (например, слайс живых кандидатов).
}

// NewGuessingSolver создаёт решателя со всеми 4-значными
// числами с неповторяющимися цифрами (5040 кандидатов).
func NewGuessingSolver() *GuessingSolver {
	panic("TODO: NewGuessingSolver")
}

// Guess возвращает текущую гипотезу решателя.
func (g *GuessingSolver) Guess() string {
	panic("TODO: Guess")
}

// Learn принимает ответ человека на попытку guess и вычёркивает
// всех кандидатов, дающих другой счёт.
func (g *GuessingSolver) Learn(guess string, bulls, cows int) {
	panic("TODO: Learn")
}
```
{% endraw %}

## Тесты

### `cowsnbulls_test.go`

{% raw %}
```go
package cowsnbulls

import "testing"

func TestScore(t *testing.T) {
	cases := []struct {
		secret, guess string
		bulls, cows   int
	}{
		{"1234", "1234", 4, 0},
		{"1234", "4321", 0, 4},
		{"1234", "5678", 0, 0},
		{"1234", "1203", 2, 1}, // 1 и 2 на месте; 3 не на месте
		{"0987", "0897", 2, 2},
	}
	for _, tc := range cases {
		b, c := Score(tc.secret, tc.guess)
		if b != tc.bulls || c != tc.cows {
			t.Errorf("Score(%q, %q) = (%d, %d), хотел (%d, %d)",
				tc.secret, tc.guess, b, c, tc.bulls, tc.cows)
		}
	}
}

func TestNewGuessingSolverSize(t *testing.T) {
	// 10*9*8*7 = 5040 кандидатов, все с неповторяющимися цифрами.
	s := NewGuessingSolver()
	g := s.Guess()
	if len(g) != 4 {
		t.Fatalf("длина гипотезы %d, хотел 4", len(g))
	}
	seen := map[byte]bool{}
	for i := 0; i < len(g); i++ {
		if seen[g[i]] {
			t.Fatalf("в гипотезе %s повторяется цифра %c", g, g[i])
		}
		seen[g[i]] = true
	}
}

func TestLearnEliminates(t *testing.T) {
	s := NewGuessingSolver()
	s.Learn("0123", 0, 0)
	g := s.Guess()
	for i := 0; i < len(g); i++ {
		if g[i] == '0' || g[i] == '1' || g[i] == '2' || g[i] == '3' {
			t.Fatalf("после Learn(0123, 0, 0) гипотеза %s содержит вычеркнутую цифру", g)
		}
	}
}

func TestSolverConverges(t *testing.T) {
	for _, secret := range []string{"0123", "4930", "8765", "1204"} {
		s := NewGuessingSolver()
		found := false
		for i := 0; i < 25; i++ {
			g := s.Guess()
			if g == secret {
				found = true
				break
			}
			b, c := Score(secret, g)
			s.Learn(g, b, c)
		}
		if !found {
			t.Fatalf("решатель не угадал секрет %q за 25 попыток", secret)
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

// Эталонное решение главы 3.
package cowsnbulls

// Score: bulls — цифра на своём месте, cows — есть в секрете, но не там.
func Score(secret, guess string) (bulls, cows int) {
	inSecret := make(map[byte]bool, len(secret))
	for i := range secret {
		inSecret[secret[i]] = true
	}
	for i := range guess {
		switch {
		case guess[i] == secret[i]:
			bulls++
		case inSecret[guess[i]]:
			cows++
		}
	}
	return bulls, cows
}

type GuessingSolver struct {
	candidates []string
}

func NewGuessingSolver() *GuessingSolver {
	const digits = "0123456789"
	var c []string
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if j == i {
				continue
			}
			for k := 0; k < 10; k++ {
				if k == i || k == j {
					continue
				}
				for l := 0; l < 10; l++ {
					if l == i || l == j || l == k {
						continue
					}
					c = append(c, string([]byte{digits[i], digits[j], digits[k], digits[l]}))
				}
			}
		}
	}
	return &GuessingSolver{candidates: c}
}

func (g *GuessingSolver) Guess() string {
	return g.candidates[0]
}

func (g *GuessingSolver) Learn(guess string, bulls, cows int) {
	kept := g.candidates[:0]
	for _, cand := range g.candidates {
		if b, c := Score(cand, guess); b == bulls && c == cows {
			kept = append(kept, cand)
		}
	}
	g.candidates = kept
}
```
{% endraw %}


</details>

_Пакет `assignments-go/03-cowsnbulls`, упражнение **cowsnbulls**. Локально: `go test -tags solution ./03-cowsnbulls/`._
