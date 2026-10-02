package wordgames

import "testing"

func TestTileScore(t *testing.T) {
	cases := []struct {
		letter byte
		want   int
	}{
		{'a', 1}, {'e', 1}, {'i', 1}, {'o', 1}, {'u', 1}, {'l', 1},
		{'n', 1}, {'s', 1}, {'t', 1}, {'r', 1},
		{'d', 2}, {'g', 2},
		{'b', 3}, {'c', 3}, {'m', 3}, {'p', 3},
		{'f', 4}, {'h', 4}, {'v', 4}, {'w', 4}, {'y', 4},
		{'k', 5},
		{'j', 8}, {'x', 8},
		{'q', 10}, {'z', 10},
		{'A', 1}, {'Q', 10}, // заглавные — как строчные
		{'.', 0}, // не-буква — 0
	}
	for _, tc := range cases {
		if got := TileScore(tc.letter); got != tc.want {
			t.Errorf("TileScore(%q) = %d, хотел %d", tc.letter, got, tc.want)
		}
	}
}

func TestDrawDeterministic(t *testing.T) {
	b1, b2 := NewBag(42), NewBag(42)
	d1, err := b1.Draw(7)
	if err != nil {
		t.Fatalf("Draw: %v", err)
	}
	d2, err := b2.Draw(7)
	if err != nil {
		t.Fatalf("Draw: %v", err)
	}
	if string(d1) != string(d2) {
		t.Errorf("одинаковый seed — одинаковые плитки: %q != %q", d1, d2)
	}
	for _, c := range d1 {
		if c < 'a' || c > 'z' {
			t.Errorf("плитка %q — не строчная буква", c)
		}
	}
	if d, err := NewBag(42).Draw(0); err != nil || len(d) != 0 {
		t.Errorf("Draw(0): хотел пустой слайс и nil, получил %d, %v", len(d), err)
	}
}

func TestDrawAll(t *testing.T) {
	b := NewBag(7)
	tiles, err := b.Draw(98)
	if err != nil {
		t.Fatalf("Draw(98): %v", err)
	}
	if len(tiles) != 98 {
		t.Fatalf("получил %d плиток, хотел 98", len(tiles))
	}
	// частоты совпадают с банком
	counts := map[byte]int{}
	for _, c := range tiles {
		counts[c]++
	}
	want := map[byte]int{
		'a': 9, 'b': 2, 'c': 2, 'd': 4, 'e': 12, 'f': 2, 'g': 3, 'h': 2,
		'i': 9, 'j': 1, 'k': 1, 'l': 4, 'm': 2, 'n': 6, 'o': 8, 'p': 2,
		'q': 1, 'r': 6, 's': 4, 't': 6, 'u': 4, 'v': 2, 'w': 2, 'x': 1,
		'y': 2, 'z': 1,
	}
	for c, n := range want {
		if counts[c] != n {
			t.Errorf("плиток %q: %d, хотел %d", c, counts[c], n)
		}
	}
	// мешок пуст — ошибка
	if _, err := b.Draw(1); err == nil {
		t.Error("Draw из пустого мешка: хотел ошибку")
	}
	// больше, чем есть — ошибка
	if _, err := NewBag(7).Draw(99); err == nil {
		t.Error("Draw(99) из мешка в 98 плиток: хотел ошибку")
	}
}

func TestScoreWord(t *testing.T) {
	hand := []byte("scrabble")
	// валидное слово из руки — сумма очков плиток (пример из спеки: SCRABBLE = 14)
	score, err := ScoreWord(hand, "scrabble", Words)
	if err != nil {
		t.Fatalf("ScoreWord(scrabble): %v", err)
	}
	if score != 14 {
		t.Errorf("ScoreWord(scrabble) = %d, хотел 14", score)
	}
	// другое слово из той же руки
	score, err = ScoreWord(hand, "able", Words)
	if err != nil {
		t.Fatalf("ScoreWord(able): %v", err)
	}
	if score != 6 { // a1 + b3 + l1 + e1
		t.Errorf("ScoreWord(able) = %d, хотел 6", score)
	}
	// слова нет в словаре
	if _, err := ScoreWord(hand, "belcas", Words); err == nil {
		t.Error("ScoreWord(belcas): слова нет в словаре — хотел ошибку")
	}
	// буквы нет в руке (t)
	if _, err := ScoreWord(hand, "tables", Words); err == nil {
		t.Error("ScoreWord(tables): в руке нет t — хотел ошибку")
	}
	// плиток не хватает: нужно две s, в руке одна
	if _, err := ScoreWord(hand, "scrabbles", Words); err == nil {
		t.Error("ScoreWord(scrabbles): не хватает плиток — хотел ошибку")
	}
}
