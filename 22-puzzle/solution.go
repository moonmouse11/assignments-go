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
