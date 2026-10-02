//go:build solution

// Эталонное решение главы 3.
package cowsnbulls

// Score: bulls — цифра на своём месте, cows — есть в секрете, но не там.
func Score(secret, guess string) (bulls, cows int) {
	inSecret := make(map[byte]bool, len(secret))
	for i := range secret {
		inSecret[secret[i]] = true
	}
	for i := range guess {
		switch {
		case guess[i] == secret[i]:
			bulls++
		case inSecret[guess[i]]:
			cows++
		}
	}
	return bulls, cows
}

type GuessingSolver struct {
	candidates []string
}

func NewGuessingSolver() *GuessingSolver {
	const digits = "0123456789"
	var c []string
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if j == i {
				continue
			}
			for k := 0; k < 10; k++ {
				if k == i || k == j {
					continue
				}
				for l := 0; l < 10; l++ {
					if l == i || l == j || l == k {
						continue
					}
					c = append(c, string([]byte{digits[i], digits[j], digits[k], digits[l]}))
				}
			}
		}
	}
	return &GuessingSolver{candidates: c}
}

func (g *GuessingSolver) Guess() string {
	return g.candidates[0]
}

func (g *GuessingSolver) Learn(guess string, bulls, cows int) {
	kept := g.candidates[:0]
	for _, cand := range g.candidates {
		if b, c := Score(cand, guess); b == bulls && c == cows {
			kept = append(kept, cand)
		}
	}
	g.candidates = kept
}
