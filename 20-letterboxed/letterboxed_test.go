package letterboxed

import (
	"sort"
	"testing"
)

// Учебный пазл: side 0 = abc, side 1 = def, side 2 = ghi, side 3 = jkl.
var p = Puzzle{Sides: [4]string{"abc", "def", "ghi", "jkl"}}

func TestValidWord(t *testing.T) {
	dict := []string{"khaki", "dice", "back", "cab", "idea", "hello", "a"}
	cases := []struct {
		w    string
		want bool
	}{
		{"khaki", true},
		{"dice", true},
		{"back", false},   // b и a на одной стороне
		{"cab", false},    // c и a на одной стороне
		{"idea", false},   // d и e на одной стороне
		{"hello", false},  // буквы l нет на пазле
		{"beglad", false}, // по сторонам валидно, но слова нет в словаре
		{"a", false},      // слишком короткое
	}
	for _, tc := range cases {
		if got := ValidWord(p, tc.w, dict); got != tc.want {
			t.Errorf("ValidWord(%q) = %v, хотел %v", tc.w, got, tc.want)
		}
	}
}

func TestAllWords(t *testing.T) {
	dict := []string{"khaki", "back", "legal", "idea", "dice", "beg", "alibi", "cab", "field", "hello", "jail"}
	want := []string{"alibi", "beg", "dice", "field", "jail", "khaki", "legal"}
	got := AllWords(p, dict)
	if len(got) != len(want) {
		t.Fatalf("AllWords: получил %d слов %v, хотел %d", len(got), got, len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("AllWords[%d] = %q, хотел %q", i, got[i], want[i])
		}
	}
}

func TestAllWordsDict(t *testing.T) {
	// весь words.go валиден для учебного пазла: AllWords = Words по возрастанию
	got := AllWords(p, Words)
	want := append([]string(nil), Words...)
	sort.Strings(want)
	if len(got) != len(want) {
		t.Fatalf("AllWords: получил %d слов, хотел %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("AllWords[%d] = %q, хотел %q", i, got[i], want[i])
		}
	}
}

func TestIsSolution(t *testing.T) {
	// два слова покрывают все 12 букв: bdjceglfk -> khaki
	if !IsSolution(p, []string{"bdjceglfk", "khaki"}, Words) {
		t.Error("IsSolution: цепочка bdjceglfk + khaki покрывает все 12 букв")
	}
	// не все буквы покрыты
	if IsSolution(p, []string{"khaki", "jibe"}, Words) {
		t.Error("IsSolution: khaki + jibe не покрывают все 12 букв")
	}
	// цепочка разорвана: alibi кончается на i, а legal начинается с l
	if IsSolution(p, []string{"alibi", "legal"}, Words) {
		t.Error("IsSolution: alibi -> legal — разорванная цепочка")
	}
	// одно слово — мало
	if IsSolution(p, []string{"khaki"}, Words) {
		t.Error("IsSolution: одно слово не покрывает все 12 букв")
	}
	// пустая последовательность — не решение
	if IsSolution(p, nil, Words) {
		t.Error("IsSolution: пустая последовательность — не решение")
	}
}

func TestSpecPuzzle(t *testing.T) {
	// пример из letterboxed.pdf: «paf uot mie xlr»
	p2 := Puzzle{Sides: [4]string{"paf", "uot", "mie", "xlr"}}
	dict := []string{"pupil", "later", "rumor", "reflex", "romp", "prefix", "tote", "flap"}
	for _, w := range []string{"pupil", "later", "rumor", "reflex", "romp", "prefix"} {
		if !ValidWord(p2, w, dict) {
			t.Errorf("ValidWord(%q) на пазле из спеки: хотел true", w)
		}
	}
	for _, w := range []string{"tote", "flap"} { // t и o, a и p — одна сторона
		if ValidWord(p2, w, dict) {
			t.Errorf("ValidWord(%q) на пазле из спеки: хотел false", w)
		}
	}
	// решение из 4 слов из самой спеки
	if !IsSolution(p2, []string{"pupil", "later", "rumor", "reflex"}, dict) {
		t.Error("IsSolution: pupil + later + rumor + reflex — решение из спеки")
	}
}
