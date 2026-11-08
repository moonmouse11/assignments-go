# Глава 10. Pwdgen — генератор паролей из слов

> Сегмент — слово плюс число, пароль —
> несколько сегментов через разделитель. Начало серии о паролях.

## Цель главы

Собрать детерминированный генератор паролей: по словарю слов,
разделителю и семени он выдаёт воспроизводимые «запоминаемые» пароли
вида `protein~57038~Nettle~33965`. Воспроизводимость (seed → тот же
пароль) — ключ к тестированию генераторов.

## Теория

- `math/rand/v2`: `rand.New(rand.NewChaCha8(seed))` — быстрый
  детерминированный источник; одинаковое семя → та же серия чисел.
- Семя `uint64` укладывается в 32 байта ChaCha8 через
  `binary.LittleEndian.PutUint64` (глава 9), остальное — нули.
- Выбор: `r.IntN(len(words))` — слово; число — `lo + r.IntN(hi-lo)`.
- Пароль «запоминаемый»: инженер возьмёт технический словарь,
  музыкант — музыкальный (`pwdgen.pdf`).

## Условие

1. `NewGen(seed, words, sep)` — генератор с семенем seed, словарём
   words и разделителем sep.
2. `Segment()` — слово из словаря + случайное число из 2–4 цифр,
   склеенные через sep: `star-29660`.
3. `Password(nSeg)` — nSeg сегментов, склеенные через sep:
   `vein-47519-Bone-52061`.
4. Одинаковое семя и словарь — в точности тот же пароль.
5. Пустой словарь — panic (нет слов — нет пароля; договор в `pwdgen.go`).

Wants (необязательно, из `pwdgen.pdf`): `--max-word-length` (фильтр
длинных слов), деривация ключа из пароля (`sha256`, итерации).

## API для реализации

```go
type Gen struct {
	Words []string
	Sep   string
	rng   *rand.Rand
}
func NewGen(seed uint64, words []string, sep string) *Gen
func (g *Gen) Segment() string
func (g *Gen) Password(nSeg int) string
```

Сигнатуры и договоры — в `pwdgen.go`. Тесты — в `pwdgen_test.go`.

## Подсказки

- Число из 2–4 цифр: сначала длина `2 + r.IntN(3)`, потом значение
  `pown(10, d-1) + r.IntN(9*pown(10, d-1))` — без ведущих нулей.
- Собирай пароль `strings.Builder`-ом, сегменты между собой — тоже sep.
- nSeg ≤ 0 — пустая строка, без паники.
- Тест формата: `strings.Split(p, sep)` → чётные части — слова
  словаря, нечётные — `^\d{2,4}$`.

## Проверка

```
go test ./10-pwdgen/
go test -tags solution ./10-pwdgen/
```

## Эталон

`go test -tags solution ./10-pwdgen/`.

## Углубление

- CLI из спеки: `flag` — `-s/--segments`, `-p/--separator`,
  `-m/--max-word-length`, `-v/--verbose`; usage-вывод как у `--help`.
- Словари: встроенный (`-b/--builtin-wordlist`) — `embed`; свой из
  файла (`-f/--word-list-file-name`) — `os.Open` + `bufio.Scanner`.
- `-k/--derive-key`: ключ фиксированной длины из пароля, сложность
  растёт с `-i/--iterations` — `crypto/sha256`, вывод `encoding/hex`;
  у автора деривация на OpenSSL (внешняя зависимость, в Go не нужна).
- `-d/--dump-spec` — дамп Ada-спеки: артефакт оригинальной
  реализации (Ada/GNAT — внешняя зависимость).
- Оригинал: gitlab.com/ada23/pwdgen (сеть).

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./10-pwdgen/
```

## Заготовки (реализуйте)

### `pwdgen.go`

{% raw %}
```go
//go:build !solution

// Пакет pwdgen — глава 10. Реализуй функции так, чтобы прошли тесты
// в pwdgen_test.go. Условие — в README.md этой папки.
package pwdgen

import "math/rand/v2"

// Gen — генератор паролей: словарь слов, разделитель сегментов и
// детерминированный источник ChaCha8. Создаётся через NewGen.
type Gen struct {
	Words []string
	Sep   string
	rng   *rand.Rand
}

// NewGen возвращает генератор с семенем seed, словарём words и
// разделителем sep. Одинаковое семя и словарь — одинаковые пароли.
// Паникует, если words пуст (нет слов — нет пароля).
func NewGen(seed uint64, words []string, sep string) *Gen {
	panic("TODO: NewGen")
}

// Segment возвращает один сегмент пароля: случайное слово из словаря,
// склеенное через Sep со случайным числом из 2–4 цифр (без ведущих
// нулей). Например, при Sep "-" — "star-29660".
func (g *Gen) Segment() string {
	panic("TODO: Segment")
}

// Password возвращает nSeg сегментов, склеенных через Sep:
// "vein-47519-Bone-52061". nSeg <= 0 — пустая строка.
func (g *Gen) Password(nSeg int) string {
	panic("TODO: Password")
}
```
{% endraw %}

## Тесты

### `pwdgen_test.go`

{% raw %}
```go
package pwdgen

import (
	"regexp"
	"strings"
	"testing"
)

var words = []string{"vein", "Bone", "protein", "virus", "azalea", "Nettle"}

func TestDeterminism(t *testing.T) {
	g1 := NewGen(42, words, "-")
	g2 := NewGen(42, words, "-")
	if p1, p2 := g1.Password(3), g2.Password(3); p1 != p2 {
		t.Errorf("одинаковое семя:\nполучил %q и %q — хотел совпадение", p1, p2)
	}
	g3 := NewGen(43, words, "-")
	if p3, p4 := g3.Password(3), g1.Password(3); p3 == p4 {
		t.Errorf("разные семени:\nполучил одинаковые %q — хотел разные", p3)
	}
}

func TestFormat(t *testing.T) {
	g := NewGen(7, words, "-")
	p := g.Password(3)
	segs := strings.Split(p, "-")
	if len(segs) != 6 {
		t.Fatalf("Password(3): получил %q — хотел 6 частей (3 слова + 3 числа)", p)
	}
	wordSet := make(map[string]bool, len(words))
	for _, w := range words {
		wordSet[w] = true
	}
	numRe := regexp.MustCompile(`^\d{2,4}$`)
	for i, part := range segs {
		if i%2 == 0 {
			if !wordSet[part] {
				t.Errorf("часть %q не из словаря (пароль %q)", part, p)
			}
		} else if !numRe.MatchString(part) {
			t.Errorf("число %q не 2-4 цифры (пароль %q)", part, p)
		}
	}
}

func TestSeparator(t *testing.T) {
	g := NewGen(7, words, "~")
	p := g.Password(2)
	if len(strings.Split(p, "~")) != 4 {
		t.Errorf("Password(2) с sep ~: получил %q — хотел 4 части", p)
	}
}

func TestSmallCounts(t *testing.T) {
	g := NewGen(1, words, "-")
	if p := g.Password(0); p != "" {
		t.Errorf("Password(0): получил %q — хотел пустую строку", p)
	}
	if p := g.Password(-3); p != "" {
		t.Errorf("Password(-3): получил %q — хотел пустую строку", p)
	}
}

func TestEmptyWordsPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("NewGen с пустым words: хотел panic — не случился")
		}
	}()
	NewGen(1, nil, "-")
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
```
{% endraw %}


</details>

_Пакет `assignments-go/10-pwdgen`, упражнение **pwdgen**. Локально: `go test -tags solution ./10-pwdgen/`._
