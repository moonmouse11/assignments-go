package hanoi

import (
	"reflect"
	"testing"
)

func TestSolutionsAgree(t *testing.T) {
	for n := 1; n <= 10; n++ {
		r := SolveRecursive(n)
		i := SolveIterative(n)
		if !reflect.DeepEqual(r, i) {
			t.Fatalf("n=%d: итеративное решение отличается от рекурсивного", n)
		}
		if len(r) != (1<<n)-1 {
			t.Errorf("n=%d: ходов %d, хотел %d", n, len(r), (1<<n)-1)
		}
		if err := Verify(r, n); err != nil {
			t.Errorf("n=%d: Verify: %v", n, err)
		}
	}
}

func TestVerifyRejectsBadMoves(t *testing.T) {
	cases := []struct {
		name  string
		moves []Move
		n     int
	}{
		{"ход с пустого колышка", []Move{{1, 2}}, 1},
		{"колышек вне диапазона", []Move{{0, 5}}, 1},
		{"задача не завершена", []Move{{0, 1}}, 2},
		{"больше на меньший", []Move{{0, 1}, {0, 1}}, 2},
	}
	for _, tc := range cases {
		if err := Verify(tc.moves, tc.n); err == nil {
			t.Errorf("%s: Verify должен вернуть ошибку", tc.name)
		}
	}
}

func TestVerifyAcceptsKnownGood(t *testing.T) {
	// n=3, каноническое решение.
	moves := []Move{
		{0, 2}, {0, 1}, {2, 1}, {0, 2}, {1, 0}, {1, 2}, {0, 2},
	}
	if err := Verify(moves, 3); err != nil {
		t.Errorf("Verify(каноническое решение n=3): %v", err)
	}
}

// Сравнение производительности — как требует исходная спецификация Hanoi.pdf:
//
//	go test -bench . ./04-hanoi
func BenchmarkRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SolveRecursive(20)
	}
}

func BenchmarkIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SolveIterative(20)
	}
}
