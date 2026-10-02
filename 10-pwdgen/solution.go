//go:build solution

// Эталонное решение главы 10. Сравни со своей реализацией в pwdgen.go.
package pwdgen

import (
	"encoding/binary"
	"math/rand/v2"
	"strings"
)

// Gen — генератор паролей: словарь слов, разделитель и ChaCha8.
type Gen struct {
	Words []string
	Sep   string
	rng   *rand.Rand
}

func NewGen(seed uint64, words []string, sep string) *Gen {
	if len(words) == 0 {
		panic("pwdgen: пустой словарь слов")
	}
	var b [32]byte
	binary.LittleEndian.PutUint64(b[:8], seed)
	return &Gen{Words: words, Sep: sep, rng: rand.New(rand.NewChaCha8(b))}
}

// Segment — слово + число из 2–4 цифр, склеенные через Sep.
func (g *Gen) Segment() string {
	word := g.Words[g.rng.IntN(len(g.Words))]
	digits := 2 + g.rng.IntN(3) // 2, 3 или 4 цифры
	lo := 1
	for range digits - 1 {
		lo *= 10
	}
	num := lo + g.rng.IntN(9*lo) // [lo, 10*lo) — без ведущих нулей
	return word + g.Sep + itoa(num)
}

// Password — nSeg сегментов через Sep.
func (g *Gen) Password(nSeg int) string {
	var b strings.Builder
	for i := range max(nSeg, 0) {
		if i > 0 {
			b.WriteString(g.Sep)
		}
		b.WriteString(g.Segment())
	}
	return b.String()
}

// itoa — десятичная запись n без ведущих нулей (strconv не нужен).
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var d [12]byte
	i := len(d)
	for n > 0 {
		i--
		d[i] = byte('0' + n%10)
		n /= 10
	}
	return string(d[i:])
}
