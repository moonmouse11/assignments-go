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
