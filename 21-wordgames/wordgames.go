//go:build !solution

// Пакет wordgames — глава 21. Реализуй TileScore, Bag и ScoreWord так,
// чтобы прошли тесты в wordgames_test.go. Условие — в README.md.
package wordgames

// TileScore возвращает очки плитки с буквой letter по упрощённой шкале
// Scrabble: AEIOULNSTR=1, DG=2, BCMP=3, FHVWY=4, K=5, JX=8, QZ=10.
// Заглавные приводятся к строчным; не-буква — 0.
//
//	TileScore('a') == 1
//	TileScore('q') == 10
func TileScore(letter byte) int {
	panic("TODO: TileScore")
}

// Bag — мешок плиток: упрощённый банк Scrabble по частотам букв
// (без бланков, 98 плиток), перемешанный детерминированно по seed.
type Bag struct {
	// TODO: твои поля (например, слайс оставшихся плиток).
}

// NewBag создаёт полный перемешанный мешок. ГСЧ — math/rand/v2:
// rand.New(rand.NewChaCha8(sb)), где sb [32]byte — seed, развёрнутый
// в байты (например, encoding/binary). Одинаковый seed — одинаковый мешок.
func NewBag(seed uint64) *Bag {
	panic("TODO: NewBag")
}

// Draw вынимает n плиток из мешка. Если в мешке меньше n плиток — ошибка
// (в том числе из пустого). n == 0 — пустой слайс и nil.
func (b *Bag) Draw(n int) ([]byte, error) {
	panic("TODO: Draw")
}

// ScoreWord считает очки слова word, составленного из плиток руки hand:
// каждая плитка руки используется не более одного раза (мультимножество),
// очки — сумма TileScore по буквам. Слово обязано быть в словаре dict,
// иначе ошибка; если в руке не хватает плиток — тоже ошибка.
//
//	ScoreWord([]byte("scrabble"), "scrabble", dict) == 14, nil
func ScoreWord(hand []byte, word string, dict []string) (int, error) {
	panic("TODO: ScoreWord")
}
