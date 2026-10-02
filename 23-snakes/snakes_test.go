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
