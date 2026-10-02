//go:build solution

// Эталонное решение главы 24. Сравни со своей реализацией в flash.go.
package flash

import (
	"context"
	"encoding/binary"
	"math/rand/v2"
	"strconv"
	"time"
)

// Item — одна карточка показа: слово (Kind 'w') или число (Kind 'n').
type Item struct {
	Kind  byte
	Value string
}

func Items(words []string, nNums int, seed uint64) []Item {
	var sb [32]byte // uint64-seed разворачиваем в [32]byte
	binary.LittleEndian.PutUint64(sb[:], seed)
	r := rand.New(rand.NewChaCha8(sb))
	items := make([]Item, 0, len(words)+nNums)
	for _, w := range words {
		items = append(items, Item{Kind: 'w', Value: w})
	}
	for i := 0; i < nNums; i++ {
		items = append(items, Item{Kind: 'n', Value: strconv.Itoa(r.IntN(100))})
	}
	r.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })
	return items
}

func RecallScore(shown, recalled []Item) int {
	used := make([]bool, len(shown))
	score := 0
	for _, rec := range recalled {
		for i, sh := range shown {
			if !used[i] && sh.Kind == rec.Kind && sh.Value == rec.Value {
				used[i] = true
				score++
				break
			}
		}
	}
	return score
}

func Rotate(ctx context.Context, interval time.Duration, items []Item, onDrop func(Item)) int {
	t := time.NewTicker(interval)
	defer t.Stop()
	dropped := 0
	for _, it := range items {
		select {
		case <-ctx.Done():
			return dropped
		case <-t.C:
			if onDrop != nil {
				onDrop(it)
			}
			dropped++
		}
	}
	return dropped
}
