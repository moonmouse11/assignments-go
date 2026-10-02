# assignments-go — учебник по Go задачами

Полный курс по мотивам projectlets из репозитория
[RajaSrinivasan/assignments](https://github.com/RajaSrinivasan/assignments):
33 главы, по одной на каждый projectlet (книги автора — `Ada-Skills-Sharpened.pdf`,
`Ada-Embedded.pdf`, `techada.pdf`, `codemd.pdf` — в курс не входят),
плюс интеграционные задачи — [капстоуны](capstones/) в отдельной директории.
Спецификации взяты из PDF автора и адаптированы; решения — на Go, только stdlib.

## Как работать с курсом

1. Открой главу `NN-<name>/README.md`, прочитай цель, теорию и условие.
2. Запусти тесты — они **красные** (заготовка паникует `TODO`):
   ```
   go test ./01-hello/
   ```
3. Реализуй функции в файле `<name>.go` (заготовка с `TODO`).
4. Добейся зелёного: `go test ./NN-name/`
5. Сравни с эталоном `solution.go` — запусти его:
   ```
   go test -tags solution ./NN-name/
   ```

Файл `solution.go` собирается **только** с тегом `solution`
(`//go:build solution`), поэтому не мешает твоей реализации:
по умолчанию Go видит заготовку, с `-tags solution` — эталон.

## Разделы и главы (по нарастающей)

### I. Основы

| # | Глава | Тема | Исходный projectlet |
|---|-------|------|---------------------|
| 1 | [01-hello](01-hello/) | CLI, файловый ввод/вывод, `fmt`, `bufio` | `hello.pdf` |
| 2 | [02-dump](02-dump/) | байты, hex-дамп, обработка ошибок, конвенции CLI | `dump.pdf` |
| 3 | [03-cowsnbulls](03-cowsnbulls/) | строки, слайсы, map, фильтрация 5040 кандидатов | `cowsnbulls.pdf` |
| 4 | [04-hanoi](04-hanoi/) | рекурсия vs итерация, верификация, `testing.B` | `Hanoi.pdf` |
| 5 | [05-exec](05-exec/) | горутины, `context`, `time.Ticker`, `os/exec` | `exec.pdf` (упрощён) |

### II. Текст, время, числа

| # | Глава | Тема | Исходный projectlet |
|---|-------|------|---------------------|
| 6 | [06-search](06-search/) | поиск текста/регэкспов по файлам, обход каталогов | `search.pdf` |
| 7 | [07-diary](07-diary/) | спринты, `time`, сериализация JSON | `diary.pdf` |
| 8 | [08-numbers](08-numbers/) | библиотека теории чисел: делители, НОД, Фибоначчи | `numbers.pdf` |
| 9 | [09-random](09-random/) | ГСЧ, детерминизм по seed, статистические проверки | `random.pdf` |

### III. Криптография и безопасность

| # | Глава | Тема | Исходный projectlet |
|---|-------|------|---------------------|
| 10 | [10-pwdgen](10-pwdgen/) | генератор паролей из слов, seeded ChaCha8 | `pwdgen.pdf` |
| 11 | [11-ipadr](11-ipadr/) | IPv4/IPv6: классы, приватные, подсети, broadcast | `ipadr.pdf` |
| 12 | [12-netgeo](12-netgeo/) | JSON API, haversine, расстояния | `netgeo.pdf` |
| 13 | [13-crc](13-crc/) | своя реализация CRC32 против stdlib | `crc.pdf` |
| 14 | [14-disign](14-disign/) | Ed25519: подпись и проверка | `disign.pdf` |
| 15 | [15-ppm](15-ppm/) | защищённый формат файла: PBKDF2 + AES-GCM | `ppm.pdf` |
| 16 | [16-password](16-password/) | менеджер паролей поверх защищённого файла | `password.pdf` |
| 17 | [17-secchan](17-secchan/) | защищённый канал: ECDH X25519 + AES-GCM | `secchan.pdf` |
| 18 | [18-spm](18-spm/) | защищённые пакеты: tar.gz + манифест + подпись | `spm.pdf` |

### IV. Игры и алгоритмы

| # | Глава | Тема | Исходный projectlet |
|---|-------|------|---------------------|
| 19 | [19-jotto](19-jotto/) | словесная угадайка: один счёт, решатель по словарю | `jotto.pdf` |
| 20 | [20-letterboxed](20-letterboxed/) | NYT Letter Boxed: слова по сторонам, цепочки | `letterboxed.pdf` |
| 21 | [21-wordgames](21-wordgames/) | модель Scrabble: очки, банк плиток, рука | `Wordgames.pdf` |
| 22 | [22-puzzle](22-puzzle/) | пятнашки: ходы, перемешивание, чётность | `Puzzle.pdf` |
| 23 | [23-snakes](23-snakes/) | змейка: игровое ядро без GUI | `Snakes.pdf` |
| 24 | [24-flash](24-flash/) | флешкарты: seeded наборы, тайминг, `context` | `flash.pdf` |

### V. Сети и DSP

| # | Глава | Тема | Исходный projectlet |
|---|-------|------|---------------------|
| 25 | [25-sockgw](25-sockgw/) | шлюз TCP → UDP: `net`, горутины, дедлайны | `sockgw.pdf` |
| 26 | [26-dsp](26-dsp/) | WAV-файлы, синтез, DFT-спектр (thinkdsp) | `dsp.pdf` |

### VI. Инструменты разработки

| # | Глава | Тема | Исходный projectlet |
|---|-------|------|---------------------|
| 27 | [27-codex](27-codex/) | извлечение фрагментов кода `from:count` | `codex.pdf` |
| 28 | [28-srctrace](28-srctrace/) | сниппет трассировки к исходникам | `srctrace.pdf` |
| 29 | [29-repo](29-repo/) | grepo: мини-YAML, группы репозиториев | `repo.pdf` |
| 30 | [30-buildsystem](30-buildsystem/) | jobs: ID сборок, retention, traceability | `BuildSystem.pdf` |
| 31 | [31-project](31-project/) | proj: клонирование проектов из шаблона | `project.pdf` |
| 32 | [32-templates](32-templates/) | генератор заготовок cli/web | `templates.pdf` |
| 33 | [33-cshare](33-cshare/) | облачное хранилище через интерфейс Storage | `cshare.pdf` |

## Команды

```
go test ./...                     # свои решения, глава за главой (все красные)
go test -tags solution ./...      # эталонные решения — все зелёные
go test -bench . ./04-hanoi       # сравнение производительности алгоритмов
go vet ./... && go vet -tags solution ./...
```

## Что дальше

- **капстоуны** — [`capstones/`](capstones/): интеграционные задачи, собирающие
  главы в инструменты (два готовы: менеджер паролей, мини-CI); добавляй свои
  по тому же паттерну;
- **«Развитие»** в конце каждой главы — что добавить из полного projectlet
  (реальный git, cloud, GUI на fyne, `--clock` в exec и т.д.);
