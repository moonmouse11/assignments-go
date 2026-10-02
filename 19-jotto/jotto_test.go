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
