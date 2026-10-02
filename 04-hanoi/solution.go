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
