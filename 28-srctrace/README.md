# Глава 28. Srctrace — трассировка к исходникам

> По мотивам `srctrace.pdf`. Приложение обязано уметь рассказать, из какого
> коммита оно собрано: генерируем фрагмент исходника с данными репозитория,
> компилируем его в приложение и разбираем обратно.

## Цель главы

Освоить генерацию кода (`go:generate` в миниатюре): превратить данные
репозитория в валидный Go-файл и написать парсер для обратного разора —
round-trip без потерь.

## Теория

- Трассировка — не лог, а константы в бинарнике: `Branch`, `ShortCommit`,
  `Commit`, `BuildDate` компилируются в приложение и видны в runtime.
- Генератор — обычная функция над `strings.Builder` и `fmt.Fprintf` с
  глаголом `%q`: он корректно закавычит любую строку.
- Парсер — regexp по строкам: `(?m)^имя = "значение"`; doc-комментарии
  (`// Branch — ...`) не совпадают с шаблоном и не мешают.
- Договор двусторонний: `Parse(Snippet(info)) == info` для любых данных.

## Условие

1. `Snippet(info)` возвращает валидный Go-файл: строка `package revisions`
   и блок `const` с константами `Branch`, `ShortCommit`, `Commit`,
   `BuildDate` — значения из `info`.
2. `Parse(src)` разбирает такой файл обратно в `RepoInfo`.
3. Некорректный источник — ошибка: пустой, без `package`, без какой-либо
   из четырёх констант.

Wants (необязательно): реальные данные из git, другие языки вывода.

## API для реализации

```go
type RepoInfo struct{ Branch, Short, Commit, Date string }
func Snippet(info RepoInfo) string
func Parse(src string) (RepoInfo, error)
```

Сигнатуры и договоры — в `srctrace.go`. Тесты — в `srctrace_test.go`.

## Подсказки

- `regexp.MustCompile` — один раз на уровне пакета, а не в функции.
- `package` ищи построчно (`strings.HasPrefix`), константы — регуляркой.
- Собери совпадения в map «имя → значение», потом проверь наличие всех
  четырёх имён — это и есть проверка «некорректного».

## Проверка

```
go vet ./28-srctrace/ && go test ./28-srctrace/            # красный: TODO
go vet -tags solution ./28-srctrace/ && go test -tags solution ./28-srctrace/
```

## Эталон

`solution.go` собирается только с тегом `solution`.

## Углубление

- Полей в спеке больше наших четырёх: `VERSION_MAJOR`/`VERSION_MINOR`
  (дефолт 0), `BUILD_NUMBER` (дефолт 9999), `REPO_NAME`, короткий
  `COMMIT_ID` и полный `COMMIT_ID_FULL`, `BUILD_TIME`.
- CLI: `flag` — `-major`/`-minor`/`-build` со значениями по умолчанию,
  `-lang` (C, C++, Ada, Python, Go), `-output=name` (дефолт `revisions`
  с расширением по языку: `revisions.h`, `revisions.ads`).
- Формат по языкам: C — блок `#define VERSION_MAJOR (2)`, Ada —
  `package revisions is ... end revisions;` — та же генерация на
  `strings.Builder` и `fmt.Fprintf`, stdlib.
- Реальное наполнение `RepoInfo` — `git rev-parse --abbrev-ref HEAD`,
  `git log -1` через `os/exec`, режим `go:generate`; wants — отчёт
  по конкретному файлу (аргумент CLI). Git — внешняя зависимость,
  в этой главе не трогаем.
