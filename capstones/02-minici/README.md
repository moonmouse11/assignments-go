# Капстоун 2. MiniCI — сборки с трассировкой

> Собирает [`28-srctrace`](../../28-srctrace/) и [`30-buildsystem`](../../30-buildsystem/)
> (по мотивам `srctrace.pdf`, `BuildSystem.pdf`).

## Цель капстоуна

Интеграция: реестр сборок из главы 30 регистрирует билды, а глава 28 превращает
сведения о билде в Go-сниппет, встраиваемый в приложение. Клей `CI` делает
то, что делает настоящий CI: сборка → лог → трассировка, с retention'ом
всех хвостов сразу.

## Теория

- Пакет главы 30 называется `jobs` (импорт `"assignments-go/30-buildsystem"`),
  пакет главы 28 — `srctrace`.
- `jobs.NewBuild` пишет лог `logDir/<id>.log`; трассировку кладём рядом:
  `logDir/trace-<id>.go`.
- Retention — общая забота: `CI.Retain` должен убирать **и** логи реестра,
  **и** файлы трассировки выпавших сборок.

## Условие

1. `New(logDir)` — CI с реестром сборок.
2. `Build(branch, commit)` — зарегистрировать сборку (ID с 1, дата UTC),
   сгенерировать `srctrace.RepoInfo{Branch, Short (первые 7 символов),
   Commit, Date}` и записать сниппет в `logDir/trace-<ID>.go`.
   Возвращает сборку.
3. `All()` — все сборки по возрастанию ID.
4. `Trace(b, lang)` — traceability-текст сборки (`"go"`, `"ini"`;
   другой язык — ошибка).
5. `Retain(n)` — оставить последние n сборок: логи и файлы трассировки
   выпавших сборок удалить. n < 1 — ошибка.

## API для реализации

```go
type CI struct{ /* TODO */ }
func New(logDir string) *CI
func (c *CI) Build(branch, commit string) (jobs.Build, error)
func (c *CI) All() []jobs.Build
func (c *CI) Trace(b jobs.Build, lang string) (string, error)
func (c *CI) Retain(n int) error
```

## Подсказки

- Проверь сниппет в тесте по-настоящему: `srctrace.Parse` прочитанного
  `trace-<id>.go` должен вернуть те же Branch/Commit/Date.
- `jobs.Retain` не знает про `trace-*.go` — удаляй их сам, до вызова.

## Проверка

```
go test ./capstones/02-minici/
go test -tags solution ./capstones/02-minici/
```

## Эталон

`go test -tags solution ./capstones/02-minici/`
