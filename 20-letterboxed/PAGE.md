# Глава 20. Letter Boxed — 12 букв на 4 сторонах

> Пазл NYT [Letter Boxed](https://www.nytimes.com/puzzles/letter-boxed).

## Цель главы

Перебор с фильтрацией: проверка слов по правилам пазла, сортировка, цепочки
слов и покрытие букв. Спека решает пазл «жадным» brute force — здесь та же
идея, но в тестируемых функциях.

## Теория

- Пазл — 12 букв, по 3 на каждой из 4 сторон (`Puzzle.Sides`).
- Слово валидно, если каждая его буква есть на пазле, а **соседние буквы
  лежат на разных сторонах**; начинать можно с любой буквы.
- Решение — цепочка слов, где следующее слово начинается с последней буквы
  предыдущего, а вместе они покрывают **все 12 букв**.
- Спека упоминает пазл «paf uot mie xlr» с десятками тысяч решений из 4 слов —
  решений обычно много, важен перебор и проверка покрытия.
- Словарь `words.go`: все слова валидны для учебного пазла `["abc","def","ghi","jkl"]`.

## Условие

1. `ValidWord(p, w, dict)` — слово `w`: длина не меньше 2, буквы на пазле,
   соседние буквы с разных сторон, и `w` есть в словаре `dict`.
2. `AllWords(p, dict)` — все валидные слова словаря `dict` на пазле `p`,
   отсортированные по возрастанию.
3. `IsSolution(p, seq, dict)` — `seq`: все слова валидны, следующее начинается
   с последней буквы предыдущего, покрыты все 12 букв. Пустая `seq` — не решение.

## API для реализации

```go
type Puzzle struct{ Sides [4]string }
func ValidWord(p Puzzle, w string, dict []string) bool
func AllWords(p Puzzle, dict []string) []string
func IsSolution(p Puzzle, seq []string, dict []string) bool
```

## Подсказки

- Номер стороны буквы: `strings.ContainsRune` по `p.Sides`.
- Словарь маленький — линейного поиска достаточно.
- Покрытие: `map[rune]bool` букв всех слов `seq`; затем проверь каждую букву
  всех сторон.
- Цепочка: сравни первую руну слова `i` с последней руной слова `i-1`.

## Проверка

```
go test ./20-letterboxed/
go test -tags solution ./20-letterboxed/
```

## Эталон

`go test -tags solution ./20-letterboxed/`

## Углубление

- Спека не проверяет готовое решение, а **находит все**: обход в глубину
  от каждой стартовой буквы, следующее слово — по последней букве
  предыдущего; для пазла «paf uot mie xlr» — 84964 решения из 4 слов.
  Ядро — те же `ValidWord`/`IsSolution`.
- Таблицы спеки `w`/`wi`: слова по каждой из 12 позиций пазла и индекс по
  первой букве — в Go `map[byte][]string`; без них цепочки упираются в
  линейный поиск по всему словарю.
- Предфильтрация словаря из «Potential Improvements»: выбрасывать слова,
  не начинающиеся с одной из 12 букв (около 60% словаря), затем слова с
  буквами вне пазла.
- Замерь полный перебор бенчмарком `testing.B` — в спеке это 11.10 с на
  MacBook Pro; оптимизации выше должны давать заметный выигрыш.
- Словарь примеров спеки — 38617 слов (внешний источник, сеть); локально —
  словарь через `embed`, начиная с учебного `words.go`.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./20-letterboxed/
```

## Заготовки (реализуйте)

### `letterboxed.go`

{% raw %}
```go
//go:build !solution

// Пакет letterboxed — глава 20. Реализуй ValidWord, AllWords и IsSolution
// так, чтобы прошли тесты в letterboxed_test.go. Условие — в README.md.
package letterboxed

// Puzzle — 12 букв по 4 сторонам, по 3 буквы на стороне (Sides[i]).
type Puzzle struct{ Sides [4]string }

// ValidWord сообщает, валидно ли слово w на пазле p: длина не меньше 2,
// каждая буква w лежит на пазле, соседние буквы w — на РАЗНЫХ сторонах,
// и w есть в словаре dict. Начинать можно с любой буквы пазла.
func ValidWord(p Puzzle, w string, dict []string) bool {
	panic("TODO: ValidWord")
}

// AllWords возвращает все валидные слова словаря dict на пазле p,
// отсортированные по возрастанию (лексикографически).
func AllWords(p Puzzle, dict []string) []string {
	panic("TODO: AllWords")
}

// IsSolution сообщает, является ли seq решением пазла p: каждое слово
// валидно, следующее слово начинается с последней буквы предыдущего,
// и вместе слова покрывают все 12 букв пазла. Пустая seq — не решение.
func IsSolution(p Puzzle, seq []string, dict []string) bool {
	panic("TODO: IsSolution")
}
```
{% endraw %}

### `words.go`

{% raw %}
```go
// Пакет letterboxed — глава 20. Словарь для тестов: только данные, без кода.
package letterboxed

// Words — учебный словарь: строчные слова из букв a–l. Каждое слово
// ВАЛИДНО для учебного пазла со сторонами ["abc","def","ghi","jkl"]:
// соседние буквы всегда лежат на разных сторонах.
// bdjceglfk — служебное «слово» для теста решения из двух слов.
var Words = []string{
	// 3 буквы
	"age", "aid", "ale", "beg", "bid", "did", "die", "eld", "gel", "hag",
	"ice", "jib", "keg", "kid", "lei", "leg", "lie", "lid",
	// 4 буквы
	"blah", "dale", "dice", "dike", "fail", "gale", "hadj", "haji", "hail",
	"jail", "jibe", "kale", "kick", "lake", "lick", "like", "glia", "glib",
	"flag", "flake", "bleak", "bilge",
	// 5 букв
	"ahead", "algae", "alibi", "alike", "galah", "gelid", "glebe", "iliad",
	"khadi", "khaki", "legal", "ladle", "algal", "field", "fjeld",
	// 6 букв
	"afield", "beadle", "behead", "beheld",
	// 9 букв — служебное слово для теста IsSolution
	"bdjceglfk",
}
```
{% endraw %}

## Тесты

### `letterboxed_test.go`

{% raw %}
```go
package letterboxed

import (
	"sort"
	"testing"
)

// Учебный пазл: side 0 = abc, side 1 = def, side 2 = ghi, side 3 = jkl.
var p = Puzzle{Sides: [4]string{"abc", "def", "ghi", "jkl"}}

func TestValidWord(t *testing.T) {
	dict := []string{"khaki", "dice", "back", "cab", "idea", "hello", "a"}
	cases := []struct {
		w    string
		want bool
	}{
		{"khaki", true},
		{"dice", true},
		{"back", false},   // b и a на одной стороне
		{"cab", false},    // c и a на одной стороне
		{"idea", false},   // d и e на одной стороне
		{"hello", false},  // буквы l нет на пазле
		{"beglad", false}, // по сторонам валидно, но слова нет в словаре
		{"a", false},      // слишком короткое
	}
	for _, tc := range cases {
		if got := ValidWord(p, tc.w, dict); got != tc.want {
			t.Errorf("ValidWord(%q) = %v, хотел %v", tc.w, got, tc.want)
		}
	}
}

func TestAllWords(t *testing.T) {
	dict := []string{"khaki", "back", "legal", "idea", "dice", "beg", "alibi", "cab", "field", "hello", "jail"}
	want := []string{"alibi", "beg", "dice", "field", "jail", "khaki", "legal"}
	got := AllWords(p, dict)
	if len(got) != len(want) {
		t.Fatalf("AllWords: получил %d слов %v, хотел %d", len(got), got, len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("AllWords[%d] = %q, хотел %q", i, got[i], want[i])
		}
	}
}

func TestAllWordsDict(t *testing.T) {
	// весь words.go валиден для учебного пазла: AllWords = Words по возрастанию
	got := AllWords(p, Words)
	want := append([]string(nil), Words...)
	sort.Strings(want)
	if len(got) != len(want) {
		t.Fatalf("AllWords: получил %d слов, хотел %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("AllWords[%d] = %q, хотел %q", i, got[i], want[i])
		}
	}
}

func TestIsSolution(t *testing.T) {
	// два слова покрывают все 12 букв: bdjceglfk -> khaki
	if !IsSolution(p, []string{"bdjceglfk", "khaki"}, Words) {
		t.Error("IsSolution: цепочка bdjceglfk + khaki покрывает все 12 букв")
	}
	// не все буквы покрыты
	if IsSolution(p, []string{"khaki", "jibe"}, Words) {
		t.Error("IsSolution: khaki + jibe не покрывают все 12 букв")
	}
	// цепочка разорвана: alibi кончается на i, а legal начинается с l
	if IsSolution(p, []string{"alibi", "legal"}, Words) {
		t.Error("IsSolution: alibi -> legal — разорванная цепочка")
	}
	// одно слово — мало
	if IsSolution(p, []string{"khaki"}, Words) {
		t.Error("IsSolution: одно слово не покрывает все 12 букв")
	}
	// пустая последовательность — не решение
	if IsSolution(p, nil, Words) {
		t.Error("IsSolution: пустая последовательность — не решение")
	}
}

func TestSpecPuzzle(t *testing.T) {
	// пример из letterboxed.pdf: «paf uot mie xlr»
	p2 := Puzzle{Sides: [4]string{"paf", "uot", "mie", "xlr"}}
	dict := []string{"pupil", "later", "rumor", "reflex", "romp", "prefix", "tote", "flap"}
	for _, w := range []string{"pupil", "later", "rumor", "reflex", "romp", "prefix"} {
		if !ValidWord(p2, w, dict) {
			t.Errorf("ValidWord(%q) на пазле из спеки: хотел true", w)
		}
	}
	for _, w := range []string{"tote", "flap"} { // t и o, a и p — одна сторона
		if ValidWord(p2, w, dict) {
			t.Errorf("ValidWord(%q) на пазле из спеки: хотел false", w)
		}
	}
	// решение из 4 слов из самой спеки
	if !IsSolution(p2, []string{"pupil", "later", "rumor", "reflex"}, dict) {
		t.Error("IsSolution: pupil + later + rumor + reflex — решение из спеки")
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

// Эталонное решение главы 20. Сравни со своей реализацией в letterboxed.go.
package letterboxed

import (
	"sort"
	"strings"
)

// Puzzle — 12 букв по 4 сторонам, по 3 буквы на стороне (Sides[i]).
type Puzzle struct{ Sides [4]string }

// sideOf — номер стороны буквы r или -1, если буквы нет на пазле.
func sideOf(p Puzzle, r rune) int {
	for i, s := range p.Sides {
		if strings.ContainsRune(s, r) {
			return i
		}
	}
	return -1
}

// ValidWord: в словаре, буквы на пазле, соседние буквы на разных сторонах.
func ValidWord(p Puzzle, w string, dict []string) bool {
	rs := []rune(w)
	if len(rs) < 2 {
		return false
	}
	found := false
	for _, d := range dict {
		if d == w {
			found = true
			break
		}
	}
	if !found {
		return false
	}
	prev := -1
	for _, r := range rs {
		s := sideOf(p, r)
		if s < 0 || s == prev {
			return false
		}
		prev = s
	}
	return true
}

// AllWords — валидные слова словаря, по возрастанию.
func AllWords(p Puzzle, dict []string) []string {
	var out []string
	for _, w := range dict {
		if ValidWord(p, w, dict) {
			out = append(out, w)
		}
	}
	sort.Strings(out)
	return out
}

// IsSolution: цепочка валидных слов, покрывающая все 12 букв.
func IsSolution(p Puzzle, seq []string, dict []string) bool {
	if len(seq) == 0 {
		return false
	}
	covered := map[rune]bool{}
	for i, w := range seq {
		if !ValidWord(p, w, dict) {
			return false
		}
		rs := []rune(w)
		if i > 0 { // следующее слово начинается с последней буквы предыдущего
			prev := []rune(seq[i-1])
			if rs[0] != prev[len(prev)-1] {
				return false
			}
		}
		for _, r := range rs {
			covered[r] = true
		}
	}
	for _, s := range p.Sides {
		for _, r := range s {
			if !covered[r] {
				return false
			}
		}
	}
	return true
}
```
{% endraw %}


</details>

_Пакет `assignments-go/20-letterboxed`, упражнение **letterboxed**. Локально: `go test -tags solution ./20-letterboxed/`._
