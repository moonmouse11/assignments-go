# Глава 21. Wordgames — модель Scrabble: мешок, рука, счёт

> Модель из спеки: плитки с очками, мешок
> с частотами букв, рука из плиток и слова из руки.

## Цель главы

Байты как буквы, таблицы очков, детерминированная случайность
(`math/rand/v2`, ChaCha8) и мультимножественная проверка «составляется ли
слово из руки».

## Теория

- Шкала Scrabble (упрощённая): `AEIOULNSTR=1, DG=2, BCMP=3, FHVWY=4, K=5, JX=8, QZ=10`.
- Банк: у каждой буквы своя частота (A — 9 плиток, X — одна); всего 98 плиток,
  бланков нет. Мешок — слайс «по плитке на вхождение», тогда равномерный
  `Shuffle` честно отражает частоты (приём прямо из спеки).
- Рука — мультимножество плиток: слово «babble» требует **три** `b`,
  даже если в руке их две.
- ГСЧ только `rand.New(rand.NewChaCha8(seedBytes))`, где `seedBytes [32]byte` —
  uint64-seed, развёрнутый в байты (`encoding/binary`). Одинаковый seed —
  одинаковый мешок: тесты детерминированы.

## Условие

1. `TileScore(letter)` — очки буквы; заглавные как строчные, не-буква — 0.
2. `Bag` + `NewBag(seed)` — полный перемешанный мешок (98 плиток,
   частоты банка). `Draw(n)` вынимает `n` плиток; не хватает — ошибка
   (в том числе из пустого); `n == 0` — пустой слайс и `nil`.
3. `ScoreWord(hand, word, dict)` — очки слова, составленного из плиток руки
   (каждая плитка не более одного раза). Слово не в словаре — ошибка;
   в руке нет нужной плитки (или их мало) — ошибка.

## API для реализации

```go
func TileScore(letter byte) int
type Bag struct{ /* ... */ }
func NewBag(seed uint64) *Bag
func (b *Bag) Draw(n int) ([]byte, error)
func ScoreWord(hand []byte, word string, dict []string) (int, error)
```

## Подсказки

- Мешок удобно хранить слайсом: `Draw` отрезает префикс
  (`b.tiles = b.tiles[n:]`).
- Проверка руки: `map[byte]int` счётчик плиток; по букве слова —
  уменьшение, ноль — ошибка.
- Словарь `words.go` — для тестов `ScoreWord`.

## Проверка

```
go test ./21-wordgames/
go test -tags solution ./21-wordgames/
```

## Эталон

`go test -tags solution ./21-wordgames/`

## Углубление

- Бланки: в полном банке 2 пустые плитки с 0 очков (в главе — 98 плиток
  без бланков); бланк замещает любую букву при составлении слова.
- Подписи-анаграммы спеки: буква → простое число, слово → произведение
  **уникальных** букв (`uint64`), словарь — `map[uint64][]string`. Буква
  может входить в слово любое число раз (PRECIPICE из одной E) — так рука
  CRPEI мгновенно даёт CREEPIE, PIERCE, RECIPE, PRECIPICE…
- Перебор подмножеств руки битовой маской (`1<<i` — буква включена).
  Режимы спеки: player — перечислить все составляемые слова с очками и
  суммой; dealer — перегенерация руки, если потенциальных слов меньше
  порога или гласных больше 3. REPL-команды примеров: `score`, `show`,
  `hand`.
- Полный scrabble-словарь из спеки — 161524 слова (внешний источник,
  сеть); локально — словарь через `embed`.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./21-wordgames/
```

## Заготовки (реализуйте)

### `wordgames.go`

{% raw %}
```go
//go:build !solution

// Пакет wordgames — глава 21. Реализуй TileScore, Bag и ScoreWord так,
// чтобы прошли тесты в wordgames_test.go. Условие — в README.md.
package wordgames

// TileScore возвращает очки плитки с буквой letter по упрощённой шкале
// Scrabble: AEIOULNSTR=1, DG=2, BCMP=3, FHVWY=4, K=5, JX=8, QZ=10.
// Заглавные приводятся к строчным; не-буква — 0.
//
//	TileScore('a') == 1
//	TileScore('q') == 10
func TileScore(letter byte) int {
	panic("TODO: TileScore")
}

// Bag — мешок плиток: упрощённый банк Scrabble по частотам букв
// (без бланков, 98 плиток), перемешанный детерминированно по seed.
type Bag struct {
	// TODO: твои поля (например, слайс оставшихся плиток).
}

// NewBag создаёт полный перемешанный мешок. ГСЧ — math/rand/v2:
// rand.New(rand.NewChaCha8(sb)), где sb [32]byte — seed, развёрнутый
// в байты (например, encoding/binary). Одинаковый seed — одинаковый мешок.
func NewBag(seed uint64) *Bag {
	panic("TODO: NewBag")
}

// Draw вынимает n плиток из мешка. Если в мешке меньше n плиток — ошибка
// (в том числе из пустого). n == 0 — пустой слайс и nil.
func (b *Bag) Draw(n int) ([]byte, error) {
	panic("TODO: Draw")
}

// ScoreWord считает очки слова word, составленного из плиток руки hand:
// каждая плитка руки используется не более одного раза (мультимножество),
// очки — сумма TileScore по буквам. Слово обязано быть в словаре dict,
// иначе ошибка; если в руке не хватает плиток — тоже ошибка.
//
//	ScoreWord([]byte("scrabble"), "scrabble", dict) == 14, nil
func ScoreWord(hand []byte, word string, dict []string) (int, error) {
	panic("TODO: ScoreWord")
}
```
{% endraw %}

### `words.go`

{% raw %}
```go
// Пакет wordgames — глава 21. Словарь для ScoreWord: только данные, без кода.
package wordgames

// Words — учебный словарь: строчные слова разных длин.
var Words = []string{
	"able", "about", "above", "actor", "adapt", "after", "again", "agent",
	"agree", "ahead", "album", "alert", "alike", "alive", "allow", "alone",
	"along", "alter", "among", "anger", "angle", "angry", "apart", "apple",
	"apply", "arena", "argue", "arise", "array", "aside", "asset", "audio",
	"avoid", "awake", "award", "badly", "baker", "basic", "beach", "began",
	"begin", "being", "below", "bench", "birth", "black", "blame", "blind",
	"block", "blood", "board", "boost", "booth", "brain", "brand", "bread",
	"break", "breed", "brick", "brief", "bring", "broad", "brown", "build",
	"built", "bunch", "cable", "candle", "castle", "chance", "check", "chess",
	"civil", "claim", "class", "clean", "clear", "click", "climb", "clock",
	"close", "coach", "coast", "label", "tables", "scrabble", "scrabbles",
}
```
{% endraw %}

## Тесты

### `wordgames_test.go`

{% raw %}
```go
package wordgames

import "testing"

func TestTileScore(t *testing.T) {
	cases := []struct {
		letter byte
		want   int
	}{
		{'a', 1}, {'e', 1}, {'i', 1}, {'o', 1}, {'u', 1}, {'l', 1},
		{'n', 1}, {'s', 1}, {'t', 1}, {'r', 1},
		{'d', 2}, {'g', 2},
		{'b', 3}, {'c', 3}, {'m', 3}, {'p', 3},
		{'f', 4}, {'h', 4}, {'v', 4}, {'w', 4}, {'y', 4},
		{'k', 5},
		{'j', 8}, {'x', 8},
		{'q', 10}, {'z', 10},
		{'A', 1}, {'Q', 10}, // заглавные — как строчные
		{'.', 0}, // не-буква — 0
	}
	for _, tc := range cases {
		if got := TileScore(tc.letter); got != tc.want {
			t.Errorf("TileScore(%q) = %d, хотел %d", tc.letter, got, tc.want)
		}
	}
}

func TestDrawDeterministic(t *testing.T) {
	b1, b2 := NewBag(42), NewBag(42)
	d1, err := b1.Draw(7)
	if err != nil {
		t.Fatalf("Draw: %v", err)
	}
	d2, err := b2.Draw(7)
	if err != nil {
		t.Fatalf("Draw: %v", err)
	}
	if string(d1) != string(d2) {
		t.Errorf("одинаковый seed — одинаковые плитки: %q != %q", d1, d2)
	}
	for _, c := range d1 {
		if c < 'a' || c > 'z' {
			t.Errorf("плитка %q — не строчная буква", c)
		}
	}
	if d, err := NewBag(42).Draw(0); err != nil || len(d) != 0 {
		t.Errorf("Draw(0): хотел пустой слайс и nil, получил %d, %v", len(d), err)
	}
}

func TestDrawAll(t *testing.T) {
	b := NewBag(7)
	tiles, err := b.Draw(98)
	if err != nil {
		t.Fatalf("Draw(98): %v", err)
	}
	if len(tiles) != 98 {
		t.Fatalf("получил %d плиток, хотел 98", len(tiles))
	}
	// частоты совпадают с банком
	counts := map[byte]int{}
	for _, c := range tiles {
		counts[c]++
	}
	want := map[byte]int{
		'a': 9, 'b': 2, 'c': 2, 'd': 4, 'e': 12, 'f': 2, 'g': 3, 'h': 2,
		'i': 9, 'j': 1, 'k': 1, 'l': 4, 'm': 2, 'n': 6, 'o': 8, 'p': 2,
		'q': 1, 'r': 6, 's': 4, 't': 6, 'u': 4, 'v': 2, 'w': 2, 'x': 1,
		'y': 2, 'z': 1,
	}
	for c, n := range want {
		if counts[c] != n {
			t.Errorf("плиток %q: %d, хотел %d", c, counts[c], n)
		}
	}
	// мешок пуст — ошибка
	if _, err := b.Draw(1); err == nil {
		t.Error("Draw из пустого мешка: хотел ошибку")
	}
	// больше, чем есть — ошибка
	if _, err := NewBag(7).Draw(99); err == nil {
		t.Error("Draw(99) из мешка в 98 плиток: хотел ошибку")
	}
}

func TestScoreWord(t *testing.T) {
	hand := []byte("scrabble")
	// валидное слово из руки — сумма очков плиток (пример из спеки: SCRABBLE = 14)
	score, err := ScoreWord(hand, "scrabble", Words)
	if err != nil {
		t.Fatalf("ScoreWord(scrabble): %v", err)
	}
	if score != 14 {
		t.Errorf("ScoreWord(scrabble) = %d, хотел 14", score)
	}
	// другое слово из той же руки
	score, err = ScoreWord(hand, "able", Words)
	if err != nil {
		t.Fatalf("ScoreWord(able): %v", err)
	}
	if score != 6 { // a1 + b3 + l1 + e1
		t.Errorf("ScoreWord(able) = %d, хотел 6", score)
	}
	// слова нет в словаре
	if _, err := ScoreWord(hand, "belcas", Words); err == nil {
		t.Error("ScoreWord(belcas): слова нет в словаре — хотел ошибку")
	}
	// буквы нет в руке (t)
	if _, err := ScoreWord(hand, "tables", Words); err == nil {
		t.Error("ScoreWord(tables): в руке нет t — хотел ошибку")
	}
	// плиток не хватает: нужно две s, в руке одна
	if _, err := ScoreWord(hand, "scrabbles", Words); err == nil {
		t.Error("ScoreWord(scrabbles): не хватает плиток — хотел ошибку")
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

// Эталонное решение главы 21. Сравни со своей реализацией в wordgames.go.
package wordgames

import (
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"slices"
)

// bankSpec — упрощённый банк Scrabble: частота каждой буквы, без бланков.
var bankSpec = []struct {
	letter byte
	count  int
}{
	{'a', 9}, {'b', 2}, {'c', 2}, {'d', 4}, {'e', 12}, {'f', 2}, {'g', 3},
	{'h', 2}, {'i', 9}, {'j', 1}, {'k', 1}, {'l', 4}, {'m', 2}, {'n', 6},
	{'o', 8}, {'p', 2}, {'q', 1}, {'r', 6}, {'s', 4}, {'t', 6}, {'u', 4},
	{'v', 2}, {'w', 2}, {'x', 1}, {'y', 2}, {'z', 1},
}

// TileScore — очки буквы по упрощённой шкале Scrabble.
func TileScore(letter byte) int {
	if letter >= 'A' && letter <= 'Z' {
		letter += 'a' - 'A'
	}
	switch letter {
	case 'a', 'e', 'i', 'o', 'u', 'l', 'n', 's', 't', 'r':
		return 1
	case 'd', 'g':
		return 2
	case 'b', 'c', 'm', 'p':
		return 3
	case 'f', 'h', 'v', 'w', 'y':
		return 4
	case 'k':
		return 5
	case 'j', 'x':
		return 8
	case 'q', 'z':
		return 10
	}
	return 0
}

// Bag — перемешанный мешок; Draw вынимает с начала.
type Bag struct {
	tiles []byte
}

// NewBag — полный мешок (98 плиток), перемешанный ChaCha8 по seed.
func NewBag(seed uint64) *Bag {
	var tiles []byte
	for _, s := range bankSpec {
		for i := 0; i < s.count; i++ {
			tiles = append(tiles, s.letter)
		}
	}
	var sb [32]byte // uint64 -> [32]byte: детерминированное засеивание
	binary.LittleEndian.PutUint64(sb[:8], seed)
	r := rand.New(rand.NewChaCha8(sb))
	r.Shuffle(len(tiles), func(i, j int) { tiles[i], tiles[j] = tiles[j], tiles[i] })
	return &Bag{tiles: tiles}
}

// Draw вынимает n плиток; не хватает — ошибка.
func (b *Bag) Draw(n int) ([]byte, error) {
	if n < 0 || n > len(b.tiles) {
		return nil, fmt.Errorf("в мешке только %d плиток — нельзя вытянуть %d", len(b.tiles), n)
	}
	drawn := b.tiles[:n:n]
	b.tiles = b.tiles[n:]
	return drawn, nil
}

// ScoreWord — сумма очков плиток руки, составивших слово из словаря.
func ScoreWord(hand []byte, word string, dict []string) (int, error) {
	if !slices.Contains(dict, word) {
		return 0, fmt.Errorf("слово %q не в словаре", word)
	}
	avail := map[byte]int{}
	for _, t := range hand {
		avail[t]++
	}
	score := 0
	for i := 0; i < len(word); i++ {
		c := word[i]
		if avail[c] == 0 {
			return 0, fmt.Errorf("в руке нет плитки для буквы %q", string(rune(c)))
		}
		avail[c]--
		score += TileScore(c)
	}
	return score, nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/21-wordgames`, упражнение **wordgames**. Локально: `go test -tags solution ./21-wordgames/`._
