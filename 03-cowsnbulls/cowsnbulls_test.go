package cowsnbulls

import "testing"

func TestScore(t *testing.T) {
	cases := []struct {
		secret, guess string
		bulls, cows   int
	}{
		{"1234", "1234", 4, 0},
		{"1234", "4321", 0, 4},
		{"1234", "5678", 0, 0},
		{"1234", "1203", 2, 1}, // 1 и 2 на месте; 3 не на месте
		{"0987", "0897", 2, 2},
	}
	for _, tc := range cases {
		b, c := Score(tc.secret, tc.guess)
		if b != tc.bulls || c != tc.cows {
			t.Errorf("Score(%q, %q) = (%d, %d), хотел (%d, %d)",
				tc.secret, tc.guess, b, c, tc.bulls, tc.cows)
		}
	}
}

func TestNewGuessingSolverSize(t *testing.T) {
	// 10*9*8*7 = 5040 кандидатов, все с неповторяющимися цифрами.
	s := NewGuessingSolver()
	g := s.Guess()
	if len(g) != 4 {
		t.Fatalf("длина гипотезы %d, хотел 4", len(g))
	}
	seen := map[byte]bool{}
	for i := 0; i < len(g); i++ {
		if seen[g[i]] {
			t.Fatalf("в гипотезе %s повторяется цифра %c", g, g[i])
		}
		seen[g[i]] = true
	}
}

func TestLearnEliminates(t *testing.T) {
	s := NewGuessingSolver()
	s.Learn("0123", 0, 0)
	g := s.Guess()
	for i := 0; i < len(g); i++ {
		if g[i] == '0' || g[i] == '1' || g[i] == '2' || g[i] == '3' {
			t.Fatalf("после Learn(0123, 0, 0) гипотеза %s содержит вычеркнутую цифру", g)
		}
	}
}

func TestSolverConverges(t *testing.T) {
	for _, secret := range []string{"0123", "4930", "8765", "1204"} {
		s := NewGuessingSolver()
		found := false
		for i := 0; i < 25; i++ {
			g := s.Guess()
			if g == secret {
				found = true
				break
			}
			b, c := Score(secret, g)
			s.Learn(g, b, c)
		}
		if !found {
			t.Fatalf("решатель не угадал секрет %q за 25 попыток", secret)
		}
	}
}
