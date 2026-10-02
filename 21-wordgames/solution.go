//go:build solution

// Эталонное решение главы 21. Сравни со своей реализацией в wordgames.go.
package wordgames

import (
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"slices"
)

// bankSpec — упрощённый банк Scrabble: частота каждой буквы, без бланков.
var bankSpec = []struct {
	letter byte
	count  int
}{
	{'a', 9}, {'b', 2}, {'c', 2}, {'d', 4}, {'e', 12}, {'f', 2}, {'g', 3},
	{'h', 2}, {'i', 9}, {'j', 1}, {'k', 1}, {'l', 4}, {'m', 2}, {'n', 6},
	{'o', 8}, {'p', 2}, {'q', 1}, {'r', 6}, {'s', 4}, {'t', 6}, {'u', 4},
	{'v', 2}, {'w', 2}, {'x', 1}, {'y', 2}, {'z', 1},
}

// TileScore — очки буквы по упрощённой шкале Scrabble.
func TileScore(letter byte) int {
	if letter >= 'A' && letter <= 'Z' {
		letter += 'a' - 'A'
	}
	switch letter {
	case 'a', 'e', 'i', 'o', 'u', 'l', 'n', 's', 't', 'r':
		return 1
	case 'd', 'g':
		return 2
	case 'b', 'c', 'm', 'p':
		return 3
	case 'f', 'h', 'v', 'w', 'y':
		return 4
	case 'k':
		return 5
	case 'j', 'x':
		return 8
	case 'q', 'z':
		return 10
	}
	return 0
}

// Bag — перемешанный мешок; Draw вынимает с начала.
type Bag struct {
	tiles []byte
}

// NewBag — полный мешок (98 плиток), перемешанный ChaCha8 по seed.
func NewBag(seed uint64) *Bag {
	var tiles []byte
	for _, s := range bankSpec {
		for i := 0; i < s.count; i++ {
			tiles = append(tiles, s.letter)
		}
	}
	var sb [32]byte // uint64 -> [32]byte: детерминированное засеивание
	binary.LittleEndian.PutUint64(sb[:8], seed)
	r := rand.New(rand.NewChaCha8(sb))
	r.Shuffle(len(tiles), func(i, j int) { tiles[i], tiles[j] = tiles[j], tiles[i] })
	return &Bag{tiles: tiles}
}

// Draw вынимает n плиток; не хватает — ошибка.
func (b *Bag) Draw(n int) ([]byte, error) {
	if n < 0 || n > len(b.tiles) {
		return nil, fmt.Errorf("в мешке только %d плиток — нельзя вытянуть %d", len(b.tiles), n)
	}
	drawn := b.tiles[:n:n]
	b.tiles = b.tiles[n:]
	return drawn, nil
}

// ScoreWord — сумма очков плиток руки, составивших слово из словаря.
func ScoreWord(hand []byte, word string, dict []string) (int, error) {
	if !slices.Contains(dict, word) {
		return 0, fmt.Errorf("слово %q не в словаре", word)
	}
	avail := map[byte]int{}
	for _, t := range hand {
		avail[t]++
	}
	score := 0
	for i := 0; i < len(word); i++ {
		c := word[i]
		if avail[c] == 0 {
			return 0, fmt.Errorf("в руке нет плитки для буквы %q", string(rune(c)))
		}
		avail[c]--
		score += TileScore(c)
	}
	return score, nil
}
