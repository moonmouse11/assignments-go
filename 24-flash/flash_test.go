package flash

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func sameItems(a, b []Item) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestItemsDeterministic(t *testing.T) {
	words := []string{"альфа", "бета", "гамма", "дельта"}
	a := Items(words, 5, 42)
	b := Items(words, 5, 42)
	if len(a) != 9 {
		t.Fatalf("Items: длина %d, хотел 9 (4 слова + 5 чисел)", len(a))
	}
	if !sameItems(a, b) {
		t.Fatal("Items: одинаковый seed дал разные наборы — не детерминирован")
	}
	seen := make(map[string]bool)
	nums := 0
	for _, it := range a {
		switch it.Kind {
		case 'w':
			seen[it.Value] = true
		case 'n':
			nums++
			if _, err := strconv.Atoi(it.Value); err != nil {
				t.Errorf("Items: значение числа %q не парсится: %v", it.Value, err)
			}
		default:
			t.Errorf("Items: неожидаемый вид %q", it.Kind)
		}
	}
	if nums != 5 {
		t.Errorf("Items: чисел %d, хотел 5", nums)
	}
	for _, w := range words {
		if !seen[w] {
			t.Errorf("Items: слово %q потерялось", w)
		}
	}
	if sameItems(Items(words, 5, 1), Items(words, 5, 2)) {
		t.Error("Items: seed 1 и 2 дали одинаковые наборы — случайности нет")
	}
}

func TestRecallScore(t *testing.T) {
	shown := []Item{
		{Kind: 'w', Value: "альфа"},
		{Kind: 'n', Value: "7"},
		{Kind: 'w', Value: "бета"},
		{Kind: 'n', Value: "42"},
	}
	recalled := []Item{
		{Kind: 'n', Value: "42"},    // верно
		{Kind: 'w', Value: "бета"},  // верно, порядок другой — не важно
		{Kind: 'w', Value: "бета"},  // повтор — второй раз не считается
		{Kind: 'w', Value: "7"},     // значение есть, но вид не тот
		{Kind: 'w', Value: "гамма"}, // такого не показывали
	}
	if got := RecallScore(shown, recalled); got != 2 {
		t.Errorf("RecallScore = %d, хотел 2", got)
	}
	if got := RecallScore(shown, shown); got != 4 {
		t.Errorf("RecallScore(полный ответ) = %d, хотел 4", got)
	}
	if got := RecallScore(shown, nil); got != 0 {
		t.Errorf("RecallScore(пустой ответ) = %d, хотел 0", got)
	}
	if got := RecallScore(nil, shown); got != 0 {
		t.Errorf("RecallScore(пустой показ) = %d, хотел 0", got)
	}
}

func TestRotateCancel(t *testing.T) {
	items := Items([]string{"a", "b", "c", "d", "e"}, 2, 7)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(70 * time.Millisecond)
		cancel()
	}()
	var dropped []Item
	n := Rotate(ctx, 20*time.Millisecond, items, func(it Item) { dropped = append(dropped, it) })
	if n < 2 || n > 3 {
		t.Errorf("Rotate: отброшено %d, хотел 2..3 (интервал 20 мс, отмена 70 мс)", n)
	}
	if len(dropped) != n {
		t.Fatalf("Rotate: onDrop вызван %d раз, вернул %d", len(dropped), n)
	}
	for i, it := range dropped {
		if it != items[i] {
			t.Errorf("Rotate: отбрасывание не по порядку: dropped[%d] = %v, хотел %v", i, it, items[i])
		}
	}
}

func TestRotateAll(t *testing.T) {
	items := []Item{
		{Kind: 'n', Value: "1"},
		{Kind: 'n', Value: "2"},
		{Kind: 'n', Value: "3"},
		{Kind: 'n', Value: "4"},
	}
	var dropped []Item
	n := Rotate(context.Background(), 5*time.Millisecond, items, func(it Item) { dropped = append(dropped, it) })
	if n != len(items) {
		t.Errorf("Rotate без отмены = %d, хотел %d", n, len(items))
	}
	if !sameItems(dropped, items) {
		t.Errorf("Rotate: отброшено %v, хотел весь набор %v", dropped, items)
	}
}
