# assignments-go — учебник по Go задачами

Полный курс по мотивам projectlets из репозитория
[RajaSrinivasan/assignments](https://github.com/RajaSrinivasan/assignments):
33 главы, по одной на каждый projectlet (книги автора — `Ada-Skills-Sharpened.pdf`,
`Ada-Embedded.pdf`, `techada.pdf`, `codemd.pdf` — в курс не входят),
плюс интеграционные задачи — [капстоуны](capstones/README.md) в отдельной директории.
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
| 1 | [01-hello](01-hello/README.md) | CLI, файловый ввод/вывод, `fmt`, `bufio` | `hello.pdf` |
| 2 | [02-dump](02-dump/README.md) | байты, hex-дамп, обработка ошибок, конвенции CLI | `dump.pdf` |
| 3 | [03-cowsnbulls](03-cowsnbulls/README.md) | строки, слайсы, map, фильтрация 5040 кандидатов | `cowsnbulls.pdf` |
| 4 | [04-hanoi](04-hanoi/README.md) | рекурсия vs итерация, верификация, `testing.B` | `Hanoi.pdf` |
| 5 | [05-exec](05-exec/README.md) | горутины, `context`, `time.Ticker`, `os/exec` | `exec.pdf` (упрощён) |

### II. Текст, время, числа

| # | Глава | Тема | Исходный projectlet |
|---|-------|------|---------------------|
| 6 | [06-search](06-search/README.md) | поиск текста/регэкспов по файлам, обход каталогов | `search.pdf` |
| 7 | [07-diary](07-diary/README.md) | спринты, `time`, сериализация JSON | `diary.pdf` |
| 8 | [08-numbers](08-numbers/README.md) | библиотека теории чисел: делители, НОД, Фибоначчи | `numbers.pdf` |
| 9 | [09-random](09-random/README.md) | ГСЧ, детерминизм по seed, статистические проверки | `random.pdf` |

### III. Криптография и безопасность

| # | Глава | Тема | Исходный projectlet |
|---|-------|------|---------------------|
| 10 | [10-pwdgen](10-pwdgen/README.md) | генератор паролей из слов, seeded ChaCha8 | `pwdgen.pdf` |
| 11 | [11-ipadr](11-ipadr/README.md) | IPv4/IPv6: классы, приватные, подсети, broadcast | `ipadr.pdf` |
| 12 | [12-netgeo](12-netgeo/README.md) | JSON API, haversine, расстояния | `netgeo.pdf` |
| 13 | [13-crc](13-crc/README.md) | своя реализация CRC32 против stdlib | `crc.pdf` |
| 14 | [14-disign](14-disign/README.md) | Ed25519: подпись и проверка | `disign.pdf` |
| 15 | [15-ppm](15-ppm/README.md) | защищённый формат файла: PBKDF2 + AES-GCM | `ppm.pdf` |
| 16 | [16-password](16-password/README.md) | менеджер паролей поверх защищённого файла | `password.pdf` |
| 17 | [17-secchan](17-secchan/README.md) | защищённый канал: ECDH X25519 + AES-GCM | `secchan.pdf` |
| 18 | [18-spm](18-spm/README.md) | защищённые пакеты: tar.gz + манифест + подпись | `spm.pdf` |

### IV. Игры и алгоритмы

| # | Глава | Тема | Исходный projectlet |
|---|-------|------|---------------------|
| 19 | [19-jotto](19-jotto/README.md) | словесная угадайка: один счёт, решатель по словарю | `jotto.pdf` |
| 20 | [20-letterboxed](20-letterboxed/README.md) | NYT Letter Boxed: слова по сторонам, цепочки | `letterboxed.pdf` |
| 21 | [21-wordgames](21-wordgames/README.md) | модель Scrabble: очки, банк плиток, рука | `Wordgames.pdf` |
| 22 | [22-puzzle](22-puzzle/README.md) | пятнашки: ходы, перемешивание, чётность | `Puzzle.pdf` |
| 23 | [23-snakes](23-snakes/README.md) | змейка: игровое ядро без GUI | `Snakes.pdf` |
| 24 | [24-flash](24-flash/README.md) | флешкарты: seeded наборы, тайминг, `context` | `flash.pdf` |

### V. Сети и DSP

| # | Глава | Тема | Исходный projectlet |
|---|-------|------|---------------------|
| 25 | [25-sockgw](25-sockgw/README.md) | шлюз TCP → UDP: `net`, горутины, дедлайны | `sockgw.pdf` |
| 26 | [26-dsp](26-dsp/README.md) | WAV-файлы, синтез, DFT-спектр (thinkdsp) | `dsp.pdf` |

### VI. Инструменты разработки

| # | Глава | Тема | Исходный projectlet |
|---|-------|------|---------------------|
| 27 | [27-codex](27-codex/README.md) | извлечение фрагментов кода `from:count` | `codex.pdf` |
| 28 | [28-srctrace](28-srctrace/README.md) | сниппет трассировки к исходникам | `srctrace.pdf` |
| 29 | [29-repo](29-repo/README.md) | grepo: мини-YAML, группы репозиториев | `repo.pdf` |
| 30 | [30-buildsystem](30-buildsystem/README.md) | jobs: ID сборок, retention, traceability | `BuildSystem.pdf` |
| 31 | [31-project](31-project/README.md) | proj: клонирование проектов из шаблона | `project.pdf` |
| 32 | [32-templates](32-templates/README.md) | генератор заготовок cli/web | `templates.pdf` |
| 33 | [33-cshare](33-cshare/README.md) | облачное хранилище через интерфейс Storage | `cshare.pdf` |

## Просмотр как книги (Honkit / GitHub Pages)

Оглавление — [`SUMMARY.md`](SUMMARY.md), настройки — [`book.json`](book.json).
Страницы глав в книге — сгенерированные [`PAGE.md`](scripts/book-expand.mjs) (README + заготовки `.go` + тесты; `solution.go` — под спойлером «нажмите, чтобы показать»).

```bash
npm ci
npm run book:serve    # http://localhost:4000
npm run book:build    # book:expand + статика в _book/
```

**В интернете:** после мержа в `main` книга публикуется через [GitHub Actions](.github/workflows/pages.yml) на
**https://moonmouse11.github.io/assignments-go/** — один раз в репозитории: **Settings → Pages → Build and deployment → Source: GitHub Actions**, затем включить workflow **Deploy book to GitHub Pages** (или дождаться пуша в `main`).

## Команды

```
go test ./...                     # свои решения, глава за главой (все красные)
go test -tags solution ./...      # эталонные решения — все зелёные
go test -bench . ./04-hanoi       # сравнение производительности алгоритмов
go vet ./... && go vet -tags solution ./...
```

