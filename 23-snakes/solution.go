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
