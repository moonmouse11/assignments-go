# Глава 24. Flash — флешкарты для памяти

> Тренажёр запоминания, как у гериатров:
> показать набор слов и чисел, «вращать» показ, отбрасывая по одной карточке,
> затем опросить пациента. Ключевые темы — **горутины, время, отмена**
> (как в главе 5) плюс детерминированная случайность.

## Цель главы

Собрать детерминированный по seed набор «карточек», прокрутить его с
отбрасыванием по таймеру и оценить уровень припоминания.

## Теория

- `math/rand/v2`: `rand.New(rand.NewChaCha8(seed))` — свой ГСЧ с фиксированным
  seed; тесты воспроизводимы от запуска к запуску.
- `rand.Shuffle` перемешивает слайс на месте — «случайный порядок показа».
- Отбрасывание по таймеру — тот же паттерн, что `Periodic` в главе 5:
  `time.NewTicker` + `select` на `ctx.Done()`.
- Сравнение без учёта порядка: счётчик совпадений по паре (вид, значение);
  одна карточка засчитывается один раз.

## Условие

Тренажёр запоминания.

1. Набор: все m слов из words плюс nNums случайных чисел, порядок перемешан.
2. Показ «вращается»: каждые `interval` одна карточка **слева** отбрасывается
   (уходит в `onDrop`).
3. Отбрасывание прекращается, когда карточки кончились или ctx отменён;
   `Rotate` возвращает число отброшенных.
4. Опрос: припомненное верно, если совпадают и значение, и вид
   (слово ≠ число); порядок неважен, одна карточка — один балл.

## API для реализации

```go
type Item struct {
	Kind  byte // 'w' — слово | 'n' — число
	Value string
}
func Items(words []string, nNums int, seed uint64) []Item
func RecallScore(shown, recalled []Item) int
func Rotate(ctx context.Context, interval time.Duration, items []Item, onDrop func(Item)) int
```

Сигнатуры и договоры — в `flash.go`. Тесты — в `flash_test.go`.

## Подсказки

- `Rotate` — почти копия `Periodic` из 05-exec: первый отброс после первого
  тика, не сразу при старте.
- В `RecallScore` отмечай уже использованные элементы shown, иначе одна
  карточка засчитается дважды.
- Числа генерируй в узком диапазоне (например 0..99) — проще проверять.

## Проверка

```
go test ./24-flash/
go test -tags solution ./24-flash/
```

## Эталон

`go test -tags solution ./24-flash/`

## Углубление

- Интерактивный цикл целиком: показ карточек в терминале, очистка дисплея
  по таймеру и опрос пользователя через `os.Stdin`/`bufio.Scanner` —
  в главе только ядро (`Items`/`Rotate`/`RecallScore`), без UI-цикла.
- Второй поток из спеки: горутина-«диспетчер», которая через фиксированное
  время сама просит показ остановиться и очистить дисплей —
  `context.WithTimeout` поверх `Rotate`, а не только исчерпание карточек.
- «Вращение по часовой»: визуальная ротация дисплея со сдвигом карточек
  (escape-коды терминала через `fmt`) — stdlib, в главе только отбрасывание
  слева.
- Эталон автора — репозиторий codeberg.org/RajaSrinivasan/flash (сеть):
  сравни разбиение на потоки и управление таймингом.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./24-flash/
```

## Заготовки (реализуйте)

### `flash.go`

{% raw %}
```go
//go:build !solution

// Пакет flash — глава 24. Реализуй Items, RecallScore и Rotate так,
// чтобы прошли тесты в flash_test.go. Условие — в README.md этой папки.
package flash

import (
	"context"
	"time"
)

// Item — одна карточка показа: слово (Kind 'w') или число (Kind 'n').
// Значение — всегда строка: "бета" или "7".
type Item struct {
	Kind  byte   // 'w' — слово | 'n' — число
	Value string // текст значения
}

// Items собирает набор для показа: все слова words плюс nNums случайных
// чисел, порядок перемешан. Одинаковый seed — одинаковый набор.
func Items(words []string, nNums int, seed uint64) []Item {
	panic("TODO: Items")
}

// RecallScore считает, сколько карточек из recalled припомнено верно
// относительно shown: совпадают и значение, и вид; порядок неважен;
// каждая карточка shown засчитывается не больше одного раза.
func RecallScore(shown, recalled []Item) int {
	panic("TODO: RecallScore")
}

// Rotate каждые interval «отбрасывает» (вызывает onDrop) очередную карточку
// items слева направо. Первый отброс — через interval после старта.
// Возвращает число отброшенных до конца набора или отмены ctx.
func Rotate(ctx context.Context, interval time.Duration, items []Item, onDrop func(Item)) int {
	panic("TODO: Rotate")
}
```
{% endraw %}

## Тесты

### `flash_test.go`

{% raw %}
```go
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
```
{% endraw %}

## Эталонное решение

Сначала попробуйте сами. Решение скрыто — раскрывайте только для сверки.

<details>
<summary><strong>Эталон (`solution.go`)</strong> — нажмите, чтобы показать</summary>

{% raw %}
```go
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
```
{% endraw %}


</details>

_Пакет `assignments-go/24-flash`, упражнение **flash**. Локально: `go test -tags solution ./24-flash/`._
