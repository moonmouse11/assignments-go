# Глава 22. Пятнашки — классическая Game of Fifteen

> Классическая [15 puzzle](https://en.wikipedia.org/wiki/15_puzzle).
> Изоляция алгоритма от GUI — здесь только модель доски.

## Цель главы

Модель доски с инвариантами: валидные ходы, стабильный `String`,
детерминированное перемешивание и теория решаемости перестановок.

## Теория

- Доска N×N, клетки удобно держать row-major слайсом; пустая клетка — `0`.
- Ход — обмен плитки с пустой клеткой; плитка обязана быть её соседом
  (общая сторона).
- Перемешивание **случайными валидными ходами** всегда даёт решаемую доску —
  это проще и честнее, чем случайная перестановка.
- Решаемость ([критерий](https://en.wikipedia.org/wiki/15_puzzle#Solvability)):
  для нечётного N — чётность числа инверсий; для чётного N — сумма инверсий
  и строки пустой клетки (считая снизу, нижняя = 1) должна быть нечётной.
- ГСЧ: `rand.New(rand.NewChaCha8(sb))`, `sb [32]byte` из uint64-seed —
  одинаковое засеивание, одинаковая доска.

## Условие

1. `New(n)` — решённая доска: `1..n²-1` по порядку, `0` в конце.
2. `Move(tile)` — обмен тайла с пустой клеткой; не сосед или нет такой
   плитки — ошибка.
3. `IsSolved()` — решена ли; `String()` — стабильный формат: строки через
   `\n`, клетки `%2d` через пробел, пустая — `« _»`, в конце `\n`.
4. `Shuffle(n, seed)` — решённая доска n×n, перемешанная случайными
   валидными ходами (всегда решаемая).
5. `Solvable(tiles, n)` — решаема ли перестановка (row-major, длина n²,
   пустая клетка — 0).

## API для реализации

```go
type Board struct{ /* ... */ }
func New(n int) *Board
func (b *Board) Move(tile int) error
func (b *Board) IsSolved() bool
func (b *Board) String() string
func Shuffle(n int, seed uint64) *Board
func Solvable(tiles []int, n int) bool
```

## Подсказки

- Соседство: `dr, dc` между плиткой и пустой; сосед — `dr == 0 && |dc| == 1`
  или наоборот.
- В `Shuffle` удобно на каждом шагу собрать соседей пустой клетки и выбрать
  случайного — это гарантированно валидный ход.
- Инверсии: пары `(i, j)`, `i < j`, `tiles[i] > tiles[j]`, пустая не считается.

## Проверка

```
go test ./22-puzzle/
go test -tags solution ./22-puzzle/
```

## Эталон

`go test -tags solution ./22-puzzle/`

## Углубление

- Лица-плитки из картинки: разрезание изображения на N² тайлов — `image`,
  `image/png`, `image/jpeg` и `SubImage` (stdlib); решённая доска
  восстанавливает исходную картинку, число плиток — параметр CLI.
- Undo из «Potential Improvements»: история ходов (слайс передвижений) и
  отмена последнего хода — ядро модели это уже позволяет.
- Паттерны вместо случайных перемешиваний: числа по периметру по часовой
  стрелке, пустая клетка ближе к центру — детерминированные «красивые»
  стартовые позиции.
- Решатель из «Potential Improvements»: A*/IDA* поиска оптимального
  решения от произвольного состояния (анимация ходов — уже GUI).
- GUI-обвязка спеки — внешние зависимости (fyne или gotk3/gtk, параметр
  `[fyne|gtk]` в CLI): клики мышью по плиткам; gtk требует всех
  изменений GUI из одного потока — таймер «времени на решение» считается
  отдельной горутиной.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./22-puzzle/
```

## Заготовки (реализуйте)

### `puzzle.go`

{% raw %}
```go
//go:build !solution

// Пакет puzzle — глава 22. Реализуй доску, ходы, перемешивание
// и Solvable так, чтобы прошли тесты в puzzle_test.go.
// Условие — в README.md.
package puzzle

// Board — доска N×N: клетки в row-major слайсе, пустая клетка — 0.
type Board struct {
	// TODO: твои поля (размер и слайс клеток).
}

// New возвращает решённую доску n×n: плитки 1..n²-1 по порядку,
// пустая клетка 0 — в конце.
func New(n int) *Board {
	panic("TODO: New")
}

// Move меняет местами плитку tile и пустую клетку. Плитка обязана
// соседствовать с пустой клеткой (общая сторона), иначе ошибка.
// Нет такой плитки — тоже ошибка.
func (b *Board) Move(tile int) error {
	panic("TODO: Move")
}

// IsSolved сообщает, решена ли доска (плитки по порядку, 0 в конце).
func (b *Board) IsSolved() bool {
	panic("TODO: IsSolved")
}

// String — стабильный формат: строки через \n, клетки через пробел,
// каждая клетка — %2d, пустая клетка — « _». В конце — \n.
//
//	New(3).String() == " 1  2  3\n 4  5  6\n 7  8  _\n"
func (b *Board) String() string {
	panic("TODO: String")
}

// Shuffle возвращает решённую доску n×n, перемешанную случайными
// ВАЛИДНЫМИ ходами — поэтому всегда решаемую. ГСЧ — math/rand/v2:
// rand.New(rand.NewChaCha8(sb)), где sb [32]byte — seed, развёрнутый
// в байты. Одинаковый seed — одинаковая доска.
func Shuffle(n int, seed uint64) *Board {
	panic("TODO: Shuffle")
}

// Solvable сообщает, решаема ли перестановка tiles (row-major, длина n²,
// пустая клетка — 0): для нечётного n — чётность инверсий; для чётного n —
// инверсии + строка пустой клетки, считая снизу (см. README).
func Solvable(tiles []int, n int) bool {
	panic("TODO: Solvable")
}
```
{% endraw %}

## Тесты

### `puzzle_test.go`

{% raw %}
```go
package puzzle

import (
	"strconv"
	"strings"
	"testing"
)

func TestNewSolved(t *testing.T) {
	b := New(4)
	if !b.IsSolved() {
		t.Fatal("New(4) должен быть решённым")
	}
	want := " 1  2  3  4\n 5  6  7  8\n 9 10 11 12\n13 14 15  _\n"
	if got := b.String(); got != want {
		t.Errorf("String():\nполучил:\n%s\nхотел:\n%s", got, want)
	}
	if got := New(3).String(); got != " 1  2  3\n 4  5  6\n 7  8  _\n" {
		t.Errorf("New(3).String() = %q", got)
	}
}

func TestMove(t *testing.T) {
	b := New(4)
	// 15 соседствует с пустой клеткой
	if err := b.Move(15); err != nil {
		t.Fatalf("Move(15): %v", err)
	}
	if b.IsSolved() {
		t.Fatal("после хода доска не решена")
	}
	// вернуть 15 на место — доска снова решена
	if err := b.Move(15); err != nil {
		t.Fatalf("Move(15) назад: %v", err)
	}
	if !b.IsSolved() {
		t.Fatal("после обратного хода доска решена")
	}
}

func TestMoveErrors(t *testing.T) {
	b := New(4)
	for _, tile := range []int{1, 5, 11, 14, 0, 16, -1} {
		if err := b.Move(tile); err == nil {
			t.Errorf("Move(%d) в решённой доске: хотел ошибку — тайл не сосед пустой клетки", tile)
		}
	}
	// после хода 15 соседями пустой клетки становятся 11 и 14
	if err := b.Move(15); err != nil {
		t.Fatalf("Move(15): %v", err)
	}
	if err := b.Move(11); err != nil {
		t.Errorf("Move(11) после Move(15): 11 — сосед пустой клетки: %v", err)
	}
}

// tilesOf восстанавливает плитки из стабильного формата String().
func tilesOf(t *testing.T, b *Board) []int {
	t.Helper()
	var tiles []int
	for _, line := range strings.Split(strings.TrimSpace(b.String()), "\n") {
		for _, f := range strings.Fields(line) {
			if f == "_" {
				tiles = append(tiles, 0)
				continue
			}
			n, err := strconv.Atoi(f)
			if err != nil {
				t.Fatalf("не разобрать поле %q: %v", f, err)
			}
			tiles = append(tiles, n)
		}
	}
	return tiles
}

func TestShuffle(t *testing.T) {
	b := Shuffle(4, 1)
	if b.IsSolved() {
		t.Fatal("Shuffle не должен возвращать решённую доску")
	}
	tiles := tilesOf(t, b)
	if !Solvable(tiles, 4) {
		t.Errorf("доска после Shuffle должна быть решаемой: %v", tiles)
	}
	if Shuffle(4, 1).String() != b.String() {
		t.Error("одинаковый seed — одинаковая доска")
	}
}

func TestSolvable(t *testing.T) {
	// решённая доска 4×4
	solved4 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0}
	if !Solvable(solved4, 4) {
		t.Error("Solvable(решённая 4×4): хотел true")
	}
	// классический нерешимый случай: 14 и 15 поменены местами
	unsolvable := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 15, 14, 0}
	if Solvable(unsolvable, 4) {
		t.Error("Solvable(1..13,15,14): хотел false — знаменитый случай Самуэля Ллойда")
	}
	// нечётная ширина: важны только инверсии
	solved3 := []int{1, 2, 3, 4, 5, 6, 7, 8, 0}
	if !Solvable(solved3, 3) {
		t.Error("Solvable(решённая 3×3): хотел true")
	}
	if Solvable([]int{1, 2, 3, 4, 5, 6, 8, 7, 0}, 3) {
		t.Error("Solvable(7 и 8 поменяны местами): хотел false")
	}
	// один валидный ход от решённой — решаемо
	moved := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 0, 15}
	if !Solvable(moved, 4) {
		t.Error("Solvable(ход от решённой): хотел true")
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

// Эталонное решение главы 22. Сравни со своей реализацией в puzzle.go.
package puzzle

import (
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
)

// Board — клетки в row-major слайсе, пустая клетка — 0.
type Board struct {
	n, blank int
	cells    []int
}

// New — решённая доска: 1..n²-1 по порядку, 0 в конце.
func New(n int) *Board {
	b := &Board{n: n, blank: n*n - 1, cells: make([]int, n*n)}
	for i := range b.cells {
		b.cells[i] = (i + 1) % (n * n) // последняя клетка — 0
	}
	return b
}

// Move — обмен плитки с пустой клеткой, если они соседи.
func (b *Board) Move(tile int) error {
	if tile < 1 || tile > b.n*b.n-1 {
		return fmt.Errorf("тайла %d нет на доске", tile)
	}
	ti := slices.Index(b.cells, tile)
	dr, dc := ti/b.n-b.blank/b.n, ti%b.n-b.blank%b.n
	if dr*dr+dc*dc != 1 { // соседи — манхэттеново расстояние 1
		return fmt.Errorf("тайл %d не соседствует с пустой клеткой", tile)
	}
	b.cells[ti], b.cells[b.blank] = 0, tile
	b.blank = ti
	return nil
}

// IsSolved — клетки равны решённой доске того же размера.
func (b *Board) IsSolved() bool {
	return slices.Equal(b.cells, New(b.n).cells)
}

// String — строки через \n, клетки %2d через пробел, пустая — « _».
func (b *Board) String() string {
	rows := make([]string, b.n)
	for r := range rows {
		cells := make([]string, b.n)
		for c := range cells {
			if v := b.cells[r*b.n+c]; v == 0 {
				cells[c] = " _"
			} else {
				cells[c] = fmt.Sprintf("%2d", v)
			}
		}
		rows[r] = strings.Join(cells, " ")
	}
	return strings.Join(rows, "\n") + "\n"
}

// Shuffle — решённая доска, перемешанная случайным блужданием пустой клетки.
func Shuffle(n int, seed uint64) *Board {
	b := New(n)
	var sb [32]byte
	binary.LittleEndian.PutUint64(sb[:8], seed)
	r := rand.New(rand.NewChaCha8(sb))
	for i := 0; i < 10*n*n; i++ {
		d := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}[r.IntN(4)]
		if r2, c2 := b.blank/n+d[0], b.blank%n+d[1]; r2 >= 0 && r2 < n && c2 >= 0 && c2 < n {
			b.cells[b.blank], b.cells[r2*n+c2] = b.cells[r2*n+c2], 0
			b.blank = r2*n + c2
		}
	}
	return b
}

// Solvable — нечётное n: инверсии чётны; чётное n: сумма инверсий и строки пустой клетки снизу нечётна.
func Solvable(tiles []int, n int) bool {
	inv, blank := 0, 0
	for i, a := range tiles {
		if a == 0 {
			blank = i
		}
		for _, b2 := range tiles[i+1:] {
			if a != 0 && b2 != 0 && b2 < a {
				inv++
			}
		}
	}
	if n%2 == 1 {
		return inv%2 == 0
	}
	return (inv+n-blank/n)%2 == 1 // n - blank/n — строка пустой клетки снизу
}
```
{% endraw %}


</details>

_Пакет `assignments-go/22-puzzle`, упражнение **puzzle**. Локально: `go test -tags solution ./22-puzzle/`._
