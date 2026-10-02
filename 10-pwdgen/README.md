# Глава 10. Pwdgen — генератор паролей из слов

> По мотивам `pwdgen.pdf`. Сегмент — слово плюс число, пароль —
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
