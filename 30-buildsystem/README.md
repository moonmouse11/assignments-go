# Глава 30. Buildsystem — мини-CI

> По мотивам `BuildSystem.pdf`. Ядро build-системы: реестр сборок со
> сквозной нумерацией, лог-файлы, retention и трассировка «сборка →
> исходник».

## Цель главы

Смоделировать жизненный цикл CI: каждая сборка получает уникальный ID,
пишет лог-файл, устаревшие сборки (и их логи) удаляются, а связь
с исходником встраивается в приложение.

## Теория

- ID — монотонный счётчик с 1; после Retain он не «сбрасывается»,
  поэтому храним `lastID`, а не считаем `len(builds)`.
- Дата — `time.Now().UTC().Format(time.RFC3339)`: единый формат,
  сортируемый как строка.
- Retention: оставить последние n — хвост слайса `builds[len-n:]`,
  лог-файлы удалённых сборок — `os.Remove`.
- Traceability — те же данные в двух форматах: Go-константы (как в
  главе 28) и ini-текст с секцией `[jobs]`.

## Условие

1. `NewRegistry(logDir)` — реестр сборок; лог сборки с ID `<id>` —
   в файле `logDir/<id>.log`.
2. `NewBuild(branch, commit)` — следующая сборка: ID с 1, `Date` —
   текущее время UTC в RFC3339, лог-файл создаётся и наполняется
   заголовком; `Build.Log` — то же содержимое.
3. `All()` — все сборки в порядке возрастания ID.
4. `Retain(n)` — оставить последние n сборок: старые записи удаляются
   из реестра, их лог-файлы стираются; `n < 1` — ошибка.
5. `Traceability(b, lang)`: `"go"` — Go-константы Branch/Commit/
   BuildDate; `"ini"` — ini-текст; другой lang — ошибка.

Wants (необязательно): стадии init/build/publish/archive, скрипты
сборки, toml-конфиги, hostname build-машины.

## API для реализации

```go
type Build struct { ID int; Branch, Commit, Date string; Log []byte }
type Registry struct{ /* свои поля */ }
func NewRegistry(logDir string) *Registry
func (r *Registry) NewBuild(branch, commit string) Build
func (r *Registry) All() []Build
func (r *Registry) Retain(n int) error
func (r *Registry) Traceability(b Build, lang string) (string, error)
```

Сигнатуры и договоры — в `buildsystem.go`. Тесты — в
`buildsystem_test.go`.

## Подсказки

- `All` возвращает копию слайса: тесты не должны видеть внутренности.
- Порядок в Retain: сначала проверка `n`, потом удаления — `Retain(0)`
  ничего не трогает.
- Имена файлов — `filepath.Join(logDir, fmt.Sprintf("%d.log", id))`.

## Проверка

```
go vet ./30-buildsystem/ && go test ./30-buildsystem/            # красный: TODO
go vet -tags solution ./30-buildsystem/ && go test -tags solution ./30-buildsystem/
```

## Эталон

`solution.go` собирается только с тегом `solution`.

## Углубление

- Подкоманды: `build <job> <script> [branch]` (без branch — дефолт из
  init), `purge` (чистка по retention), `show` (полный лог),
  `enumerate --current` (ветки), `fetch --set-default` — каркас на `flag`.
- Флаги: `--jobspace` (дефолт `$JOBSPACE`, имя репо — имя джоба),
  `--init=scriptdir` (каталог скриптов в репо: shell/python3),
  `--trace=(C,C++,Ada,Python,Go,Text)` с дефолтом Ada, `--retention=3`,
  `--pull` (git pull перед сборкой).
- Сборки — каталоги `<jobspace>/<job>/<branch>/<id>` с `jobs.log`
  (таймстемп + `......> [I] сообщение`); скрипту через env передаются
  `BUILD_NUMBER`, `OUTPUT_DIR`, `SOURCE_DIR`, `REVISION_FILE`
  (`os/exec.Cmd.Env`).
- Ещё stdlib: hostname build-машины (`os.Hostname`) в трассировке, конфиги
  `jobs.toml`. Клонирование и запуск скриптов — `git` + сеть: внешние
  зависимости, в главе только реестр.
