# Глава 19. Jotto — словесная угадайка с одним счётом

> Родственник главы 3 (cows and bulls), но секрет — слово,
> а счёт — один: сколько букв угадано.

## Цель главы

Строки как мультимножества букв, map-счётчики, фильтрация кандидатов — те же
приёмы, что в `03-cowsnbulls`, но с новым «одномерным» счётом и словарём
(`words.go`) вместо 5040 чисел.

## Теория

- В jotto счёт **один**: сколько букв попытки есть в секрете — без учёта позиций.
- Буквы повторяются: «less» против «sell» — 3 общие буквы. Правильный счёт —
  **пересечение мультимножеств**: для каждой буквы минимум числа вхождений.
- Наивная реализация «+1, если буква попытки есть в секрете» **некоммутативна**:
  `f(a, b) != f(b, a)` при повторах — в `jotto.pdf` это отмечено отдельно.
  Пересечение мультимножеств симметрично — используй его.
- Иди по **рунам** (`for _, r := range s`), не по байтам: кириллическая буква —
  2 байта UTF-8, побайтовый счёт для «марс»/«срам» даст 8, а не 4.

## Условие

Игра ведётся 4-буквенными словами из словаря `words.go` (строчные буквы).

1. `Common(a, b)` — сколько букв `a` есть в `b`, с учётом повторов
   (по минимуму вхождений каждой буквы).
2. `Solver` угадывает задуманное слово: `NewSolver()` — все кандидаты из
   `words.go`; `Guess()` — текущая гипотеза; `Learn(guess, n)` — человек называет
   `n = Common(secret, guess)`, решатель вычёркивает кандидатов с другим счётом.
   Сама попытка `guess` неверна — вычёркивай и её, иначе решатель «залипнет»
   на анаграммах (типа `kale`/`lake`).
3. Решатель обязан угадывать любое слово словаря не более чем за 12 попыток.

## API для реализации

```go
func Common(a, b string) int
type Solver struct{ /* ... */ }
func NewSolver() *Solver
func (s *Solver) Guess() string
func (s *Solver) Learn(guess string, n int)
```

## Подсказки

- `map[rune]int` — счётчик вхождений; второй проход уменьшает и считает.
- Фильтрация на месте: `kept := s.candidates[:0]` — как в главе 3.
- `Common(a, b) == Common(b, a)` — хороший самоконтроль симметричности.

## Проверка

```
go test ./19-jotto/
go test -tags solution ./19-jotto/
```

## Эталон

`go test -tags solution ./19-jotto/`

## Углубление

- Режим дилера: компьютер задумывает слово из словаря (ГСЧ от времени —
  `rand.New(rand.NewChaCha8)`), человек угадывает в цикле `bufio.Scanner`
  по stdin; программа отвечает только счётом `Common(secret, guess)`.
  Режим player (решатель) в главе есть — спека требует оба режима.
- Спека обещает угадывание за **5 попыток или меньше** (глава — за 12):
  в качестве попытки выбирай кандидата, максимизирующего отсечение, а не
  первого из списка; проверь максимум и среднее по всему словарю.
- Словарь спеки — `github.com/dwyl/english-words`: только 4-буквенных слов
  больше 7000. Это внешний источник (сеть); локально — урезанный словарь
  через `embed`, фильтр по длине слова тот же.
- Learning objective спеки — динамическая память: сравни фильтрацию на
  месте (`candidates[:0]`, ноль аллокаций) с пересборкой через `append` —
  например, бенчмарком `testing.B`.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./19-jotto/
```

## Заготовки (реализуйте)

### `jotto.go`

{% raw %}
```go
//go:build !solution

// Пакет jotto — глава 19. Реализуй Common и решателя так,
// чтобы прошли тесты в jotto_test.go. Условие — в README.md.
package jotto

// Common возвращает число общих букв строк a и b — пересечение
// мультимножеств букв: для каждой буквы минимум числа вхождений.
// Правильная реализация симметрична: Common(a, b) == Common(b, a)
// (наивная «буква есть в слове» — нет, см. README). Иди по рунам,
// а не по байтам: кириллическая буква занимает 2 байта в UTF-8.
//
//	Common("марс", "срам") == 4
//	Common("ааа", "аа")    == 2
func Common(a, b string) int {
	panic("TODO: Common")
}

// Solver — компьютер угадывает задуманное 4-буквенное слово,
// вычёркивая кандидатов по единственному счёту.
type Solver struct {
	// TODO: твои поля (например, слайс живых кандидатов).
}

// NewSolver создаёт решателя: кандидаты — 4-буквенные слова
// словаря Words из words.go.
func NewSolver() *Solver {
	panic("TODO: NewSolver")
}

// Guess возвращает текущую гипотезу — слово из живых кандидатов
// (детерминированную, например первого живого).
func (s *Solver) Guess() string {
	panic("TODO: Guess")
}

// Learn принимает ответ n = Common(secret, guess) на попытку guess.
// Попытка была неверна, поэтому вычёркивается и сама: остаются
// только кандидаты cand с Common(cand, guess) == n.
func (s *Solver) Learn(guess string, n int) {
	panic("TODO: Learn")
}
```
{% endraw %}

### `words.go`

{% raw %}
```go
// Пакет jotto — глава 19. Словарь для решателя: только данные, без кода.
package jotto

// Words — словарь кандидатов: 4-буквенные слова, строчные буквы.
var Words = []string{
	"acid", "away", "baby", "back", "ball", "band", "bank", "bath", "bean", "bear",
	"bell", "best", "bike", "bird", "blue", "boat", "body", "bone", "book", "burn",
	"busy", "cake", "calm", "camp", "card", "care", "cart", "case", "cash", "cast",
	"cell", "chef", "chin", "chip", "city", "club", "coal", "coat", "code", "coin",
	"cold", "cook", "cool", "copy", "cord", "core", "corn", "cost", "crew", "crop",
	"cube", "cure", "curl", "dark", "data", "date", "dawn", "dead", "deal", "dear",
	"debt", "deep", "desk", "dial", "dice", "diet", "dirt", "dish", "dive", "down",
	"each", "earn", "ease", "east", "edge", "evil", "exit", "face", "fact", "fail",
	"find", "fine", "fire", "firm", "fish", "five", "flag", "flat", "flow", "food",
	"grow", "hair", "half", "hall", "hand", "hard", "hate", "head", "hear", "heat",
	"help", "here", "hero", "hide", "high", "hold", "home", "hope", "hour", "huge",
	"hurt", "idea", "inch", "iron", "item", "join", "jump", "jury", "keep", "kind",
	"lamp", "land", "lane", "late", "lead", "leaf", "left", "less", "lift", "like",
	"lose", "luck", "mail", "main", "make", "mark", "mask", "math", "meal", "mean",
	"navy", "neck", "need", "news", "nice", "nine", "none", "noon", "nose", "note",
	"park", "part", "pass", "past", "path", "peak", "plan", "play", "plot", "plus",
	"rice", "rich", "ride", "ring", "rise", "risk", "road", "rock", "role", "roll",
	"seal", "seat", "seed", "ship", "shop", "show", "side", "silk", "skin", "soap",
	"tail", "take", "tall", "tank", "team", "tear", "tech", "test", "time", "tiny",
	"town", "tree", "trip", "twin", "view", "wall", "want", "west", "wind", "wine",
	"wing", "word", "work", "yard", "year", "zero", "zone",
}
```
{% endraw %}

## Тесты

### `jotto_test.go`

{% raw %}
```go
package jotto

import "testing"

func TestCommon(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"марс", "срам", 4},
		{"срам", "марс", 4}, // счёт симметричен
		{"ааа", "аа", 2},
		{"аа", "ааа", 2}, // и в обратную сторону
		{"кот", "слон", 1},
		{"less", "sell", 3},
		{"zero", "zone", 3},
		{"book", "read", 0}, // непересекающиеся слова
	}
	for _, tc := range cases {
		if got := Common(tc.a, tc.b); got != tc.want {
			t.Errorf("Common(%q, %q) = %d, хотел %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestNewSolver(t *testing.T) {
	s := NewSolver()
	g := s.Guess()
	if len(g) != 4 {
		t.Fatalf("длина гипотезы %d, хотел 4", len(g))
	}
	found := false
	for _, w := range Words {
		if w == g {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("гипотеза %q — не слово словаря", g)
	}
}

func TestLearnEliminates(t *testing.T) {
	s := NewSolver()
	g := s.Guess()
	s.Learn(g, 0) // ответ «ни одной общей буквы»
	next := s.Guess()
	if Common(next, g) != 0 {
		t.Errorf("после Learn(%q, 0) гипотеза %q имеет общие буквы с вычеркнутой", g, next)
	}
}

func TestSolverConverges(t *testing.T) {
	for _, secret := range Words {
		s := NewSolver()
		found := false
		for i := 0; i < 12; i++ {
			g := s.Guess()
			if g == secret {
				found = true
				break
			}
			s.Learn(g, Common(secret, g))
		}
		if !found {
			t.Errorf("решатель не угадал %q за 12 попыток", secret)
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

// Эталонное решение главы 19. Сравни со своей реализацией в jotto.go.
package jotto

// Common — пересечение мультимножеств букв (минимум вхождений каждой буквы).
func Common(a, b string) int {
	counts := make(map[rune]int)
	for _, r := range a {
		counts[r]++
	}
	n := 0
	for _, r := range b {
		if counts[r] > 0 {
			counts[r]--
			n++
		}
	}
	return n
}

// Solver — угадыватель jotto: живые кандидаты из words.go.
type Solver struct {
	candidates []string
}

// NewSolver — все 4-буквенные слова словаря Words.
func NewSolver() *Solver {
	c := make([]string, 0, len(Words))
	for _, w := range Words {
		if len(w) == 4 {
			c = append(c, w)
		}
	}
	return &Solver{candidates: c}
}

// Guess — первый живой кандидат.
func (s *Solver) Guess() string {
	return s.candidates[0]
}

// Learn оставляет кандидатов с тем же счётом против guess;
// сама попытка guess была неверна — вычёркиваем и её.
func (s *Solver) Learn(guess string, n int) {
	kept := s.candidates[:0]
	for _, cand := range s.candidates {
		if cand != guess && Common(cand, guess) == n {
			kept = append(kept, cand)
		}
	}
	s.candidates = kept
}
```
{% endraw %}


</details>

_Пакет `assignments-go/19-jotto`, упражнение **jotto**. Локально: `go test -tags solution ./19-jotto/`._
