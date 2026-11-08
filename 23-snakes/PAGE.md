# Глава 23. Snakes — ядро игры без GUI

> Змейка: спека целиком про **отделение алгоритма
> от представления** — ядро игры должно жить и тестироваться без интерфейса.

## Цель главы

Спроектировать тестируемое игровое ядро: состояние, направления, очередь
роста, столкновения — ни строки UI. Тесты управляют игрой ровно так, как
это делал бы интерфейс: `Turn`, `Grow`, `Step`.

## Теория

- Поле `w×h`, клетки — `Point{X, Y}`; барьеры — множество клеток.
- Направления North/East/South/West — цикл по часовой стрелке:
  противоположное — `(d+2)%4`; поэтому запрет разворота на 180° —
  одна проверка.
- Рост — «долг»: `Grow(n)` кладёт `n` в очередь, каждый `Step` съедает одну
  единицу — хвост в этот ход не двигается, змейка длиннее на клетку.
- Смерть — стена, барьер или собственное тело; после смерти `Step` всегда
  `false`. В клетку освобождающегося хвоста зайти можно (классика жанра).
- Спека отдельно подчёркивает тестирование: без UI ядро покрыть тестами
  тривиально — что мы и делаем.

## Условие

1. `NewGame(w, h, barriers, seed)` — поле `w×h` (не меньше 5×5) с барьерами;
   змейка из трёх клеток, голова в центре, направление East. `seed`
   зарезервирован для будущей еды — ядро детерминировано и без него.
2. `Turn(d)` — смена направления; разворот на 180° игнорируется.
3. `Grow(n)` — очередь роста.
4. `Step()` — ход на одну клетку; `false` — смерть; рост съедает очередь.

## API для реализации

```go
type Dir int
const ( North Dir = iota; East; South; West )
type Point struct{ X, Y int }
type Game struct{ /* ... */ }
func NewGame(w, h int, barriers []Point, seed uint64) *Game
func (g *Game) Turn(d Dir)
func (g *Game) Grow(n int)
func (g *Game) Step() bool
```

## Подсказки

- Тело удобно хранить слайсом: голова — последний элемент, `append` —
  новая голова, хвост отпадает пере-срезом `body[1:]`.
- Барьеры — `map[Point]bool`: проверка клетки — `m[p]`.
- Змейка в центре поля 9×9: голова `(4,4)`, хвост к западу. До стены
  на восток — 4 шага: это основа детерминированных тестов.

## Проверка

```
go test ./23-snakes/
go test -tags solution ./23-snakes/
```

## Эталон

`go test -tags solution ./23-snakes/`

## Углубление

- CLI спеки: `snakes [vertical|horizontal] [barriers]` — начальная
  ориентация змейки (глава жёстко задаёт East) и 2 барьера, поставленные
  поперёк движения; параметризуй `NewGame` через `os.Args`/`flag`.
- Пауза: состояние paused/resume (на скриншотах спеки — «SNAKE PAUSED»);
  ядро без UI хранит это тривиально.
- Игровой темп: спека гоняет ядро из цикла интерфейса — сделай обёртку
  на `time.Ticker` и горутине, вызывающую `Step`; само ядро от темпа не
  зависит (в этом суть «algorithm hiding» спеки).
- GUI-реализации спеки — внешние зависимости (fyne и gotk3/gtk, отдельные
  исполняемые `bin/fyne/snakes` и `bin/gtk3/snakes`); спека сравнивает
  обработку нажатий клавиш в двух библиотеках — дизайн ядра вдали от
  конкретных keycodes делает миграцию дешёвой.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./23-snakes/
```

## Заготовки (реализуйте)

### `snakes.go`

{% raw %}
```go
//go:build !solution

// Пакет snakes — глава 23. Ядро игры «змейка», полностью отделённое
// от UI (в этом и смысл спеки). Реализуй так, чтобы прошли тесты
// в snakes_test.go. Условие — в README.md.
package snakes

// Dir — направление движения.
type Dir int

// Направления образуют цикл по часовой стрелке:
// противоположное направление — (d+2)%4.
const (
	North Dir = iota
	East
	South
	West
)

// Point — клетка поля: X — столбец, Y — строка.
type Point struct{ X, Y int }

// Game — состояние игры: поле, барьеры, змейка, направление, рост.
type Game struct {
	// TODO: твои поля (тело, направление, очередь роста, барьеры, dead).
}

// NewGame создаёт игру на поле w×h (не меньше 5×5) с барьерами barriers.
// Змейка — три клетки с головой в центре поля, хвост — к западу;
// направление — East. seed зарезервирован для будущей еды: ядро
// детерминировано и без него.
func NewGame(w, h int, barriers []Point, seed uint64) *Game {
	panic("TODO: NewGame")
}

// Turn задаёт направление d. Разворот на 180° игнорируется.
func (g *Game) Turn(d Dir) {
	panic("TODO: Turn")
}

// Grow добавляет n клеток роста: каждый следующий Step удлиняет змейку
// на одну клетку (хвост в этот ход не двигается).
func (g *Game) Grow(n int) {
	panic("TODO: Grow")
}

// Step двигает голову на одну клетку в текущем направлении и возвращает
// false при смерти: стена, барьер или собственное тело. В клетку
// освобождающегося хвоста зайти можно. Рост съедает очередь роста.
// После смерти Step всегда false.
func (g *Game) Step() bool {
	panic("TODO: Step")
}
```
{% endraw %}

## Тесты

### `snakes_test.go`

{% raw %}
```go
package snakes

import "testing"

// Поле 9×9: голова стартует в центре (4,4), хвост — к западу, направление East.
const (
	w, h = 9, 9
)

func TestStepMovesHead(t *testing.T) {
	g := NewGame(w, h, nil, 1)
	for i := 1; i <= 4; i++ {
		if !g.Step() {
			t.Fatalf("шаг %d на восток: змейка должна жить (поле свободно)", i)
		}
	}
	// 5-й шаг: голова в (9,4) — стена
	if g.Step() {
		t.Fatal("шаг в стену должен возвращать false")
	}
	// после смерти Step всегда false
	if g.Step() {
		t.Fatal("после смерти Step должен возвращать false")
	}
}

func TestTurn(t *testing.T) {
	// разворот на 180° игнорируется: East -> West нельзя
	g := NewGame(w, h, nil, 1)
	g.Turn(West)
	for i := 1; i <= 4; i++ {
		if !g.Step() {
			t.Fatalf("шаг %d: змейка всё ещё движется на восток", i)
		}
	}
	if g.Step() {
		t.Fatal("после игнорирования Turn(West) змейка упирается в восточную стену")
	}

	// поворот работает: с центра на север — 4 шага до стены
	g = NewGame(w, h, nil, 1)
	g.Turn(North)
	for i := 1; i <= 4; i++ {
		if !g.Step() {
			t.Fatalf("шаг %d на север: 4 шага до стены", i)
		}
	}
	if g.Step() {
		t.Fatal("5-й шаг на север — в стену")
	}
}

func TestBarrier(t *testing.T) {
	g := NewGame(w, h, []Point{{X: 6, Y: 4}}, 1)
	if !g.Step() {
		t.Fatal("шаг на (5,4): барьера там нет")
	}
	if g.Step() {
		t.Fatal("шаг на (6,4) — барьер: хотел false")
	}
}

func TestGrowAndSelfCollision(t *testing.T) {
	// без роста: короткая змейка проходит петлю свободно
	g := NewGame(w, h, nil, 1)
	for i := 0; i < 3; i++ {
		if !g.Step() {
			t.Fatalf("шаг %d: поле свободно", i+1)
		}
	}
	g.Turn(South)
	g.Step()
	g.Turn(West)
	g.Step()
	g.Turn(North)
	if !g.Step() {
		t.Fatal("короткая змейка не задевает себя в петле")
	}

	// с ростом Grow(2): змейка длиннее на 2 клетки за 2 хода —
	// та же петля ведёт голову в собственное тело
	g = NewGame(w, h, nil, 1)
	g.Grow(2)
	for i := 0; i < 3; i++ {
		if !g.Step() {
			t.Fatalf("шаг %d с ростом: поле свободно", i+1)
		}
	}
	g.Turn(South)
	if !g.Step() {
		t.Fatal("поворот на юг: клетка свободна")
	}
	g.Turn(West)
	if !g.Step() {
		t.Fatal("поворот на запад: клетка свободна")
	}
	g.Turn(North)
	if g.Step() {
		t.Fatal("выросшая змейка кусает себя за тело: хотел false")
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

// Эталонное решение главы 23. Сравни со своей реализацией в snakes.go.
package snakes

// Dir — направление движения.
type Dir int

// Направления образуют цикл по часовой стрелке:
// противоположное направление — (d+2)%4.
const (
	North Dir = iota
	East
	South
	West
)

// Point — клетка поля: X — столбец, Y — строка.
type Point struct{ X, Y int }

// Game: тело — слайс, хвост первый, голова последняя.
type Game struct {
	w, h     int
	barriers map[Point]bool
	body     []Point
	dir      Dir
	grow     int
	dead     bool
}

// deltas[Dir] — сдвиг головы: North — вверх (Y-1) и т. д. по часовой стрелке.
var deltas = [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

// NewGame — змейка из трёх клеток с головой в центре, направление East.
func NewGame(w, h int, barriers []Point, seed uint64) *Game {
	bs := make(map[Point]bool, len(barriers))
	for _, p := range barriers {
		bs[p] = true
	}
	cx, cy := w/2, h/2
	return &Game{
		w: w, h: h, barriers: bs, dir: East,
		body: []Point{{X: cx - 2, Y: cy}, {X: cx - 1, Y: cy}, {X: cx, Y: cy}},
	}
}

// Turn — смена направления; разворот на 180° игнорируется.
func (g *Game) Turn(d Dir) {
	if (d+2)%4 == g.dir {
		return
	}
	g.dir = d
}

// Grow — очередь роста: один Step съедает одну единицу.
func (g *Game) Grow(n int) {
	if n > 0 {
		g.grow += n
	}
}

// Step — ход на одну клетку; false — смерть (стена, барьер, своё тело).
func (g *Game) Step() bool {
	if g.dead {
		return false
	}
	d := deltas[g.dir]
	head := g.body[len(g.body)-1]
	next := Point{X: head.X + d[0], Y: head.Y + d[1]}
	// стены
	if next.X < 0 || next.X >= g.w || next.Y < 0 || next.Y >= g.h {
		g.dead = true
		return false
	}
	// барьеры
	if g.barriers[next] {
		g.dead = true
		return false
	}
	// своё тело; без роста хвост уйдёт — его клетка свободна
	body := g.body
	if g.grow == 0 {
		body = body[1:]
	}
	for _, p := range body {
		if p == next {
			g.dead = true
			return false
		}
	}
	g.body = append(g.body, next)
	if g.grow == 0 {
		g.body = g.body[1:] // хвост отпадает
	} else {
		g.grow--
	}
	return true
}
```
{% endraw %}


</details>

_Пакет `assignments-go/23-snakes`, упражнение **snakes**. Локально: `go test -tags solution ./23-snakes/`._
