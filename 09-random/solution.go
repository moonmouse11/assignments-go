//go:build solution

// Эталонное решение главы 9. Сравни со своей реализацией в random.go.
package random

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand/v2"
)

// chacha — детерминированный источник: 8 байт семени, остальное нули.
func chacha(seed uint64) *rand.Rand {
	var b [32]byte
	binary.LittleEndian.PutUint64(b[:8], seed)
	return rand.New(rand.NewChaCha8(b))
}

func Uniform(seed uint64, lo, hi float64, n int) []float64 {
	r := chacha(seed)
	out := make([]float64, max(n, 0))
	for i := range out {
		out[i] = lo + r.Float64()*(hi-lo)
	}
	return out
}

func Normal(seed uint64, mean, std float64, n int) []float64 {
	r := chacha(seed)
	out := make([]float64, max(n, 0))
	for i := range out {
		out[i] = mean + std*r.NormFloat64()
	}
	return out
}

func CryptoUniform(lo, hi float64, n int) ([]float64, error) {
	out := make([]float64, max(n, 0))
	var b [8]byte
	for i := range out {
		if _, err := crand.Read(b[:]); err != nil {
			return nil, err
		}
		u := float64(binary.LittleEndian.Uint64(b[:])>>11) / (1 << 53)
		out[i] = lo + u*(hi-lo)
	}
	return out, nil
}
