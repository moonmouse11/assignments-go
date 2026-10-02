# Глава 29. Repo — группы репозиториев

> По мотивам `repo.pdf` (grepo). Большой проект — это много репозиториев:
> конфигурация группы в мини-YAML и подготовка рабочей области (workarea)
> по списку путей проектов.

## Цель главы

Написать парсер подмножества YAML руками — без внешних либ — и научиться
готовить структуру каталогов из конфигурации.

## Теория

- YAML держится на отступах: top-level ключи — без отступа, элементы
  списка начинаются с `-`, поля элемента — с отступом под ним.
- Значения в кавычках (`"$HOME/x"`) — просто строки: кавычки срезаем,
  переменные окружения не разворачиваем.
- Полный путь проекта — `filepath.Join(workarea, path)`; уникальность
  `path` проверяется до создания каталогов (map «path → проект»).
- `os.MkdirAll` создаёт цепочку каталогов целиком — baseDir отдельно
  создавать не нужно.

## Условие

Конфигурация в духе `AllProjects.yaml` (упрощённо — одна группа):

```yaml
workarea: "$HOME/grepotest/pub"
projects:
  - repourl: "git@gitlab.com:projtemplates/go.git"
    reference: "master"
    path: "go/lib"
```

1. `ParseConfig` разбирает `workarea` и список проектов (поля `repourl`,
   `reference`, `path`); пустые строки и комментарии `#` пропускаются.
2. Некорректная конфигурация — ошибка: строка без двоеточия, неизвестное
   поле, нет `workarea` или нет `projects:`.
3. `Prepare` создаёт baseDir и каталоги `Path` всех проектов, возвращает
   список полных путей в порядке объявления; дубликат `Path` — ошибка.

Wants (необязательно): разворачивание `$HOME`, группы public/private,
поля `repo`/`server`/`build`.

## API для реализации

```go
type Project struct{ RepoURL, Reference, Path string }
type Config struct { Workarea string; Projects []Project }
func ParseConfig(r io.Reader) (Config, error)
func Prepare(cfg Config, baseDir string) ([]string, error)
```

Сигнатуры и договоры — в `repo.go`. Тесты — в `repo_test.go`.

## Подсказки

- Разбор по строкам: `strings.TrimSpace` + проверка префикса `- `.
- Пара «ключ: значение» — `strings.Index(s, ":")`; кавычки `"` и `'` — срезать.
- Указатель `cur *Project` на последний элемент списка упрощает дополнение
  полей следующих строк.

## Проверка

```
go vet ./29-repo/ && go test ./29-repo/            # красный: TODO
go vet -tags solution ./29-repo/ && go test -tags solution ./29-repo/
```

## Эталон

`solution.go` собирается только с тегом `solution`.

## Углубление

- Групп две — public/private: public — только reference (клон с detached
  HEAD, без push и tag); `server` группы + `repo` проекта склеиваются в
  полный URL; `reference` — дефолт группы (ветка, тег или commit id),
  переопределяется в проекте; `build` — shell-команда после clone/pull.
- `workarea` переводится через переменные окружения (`os.Getenv`);
  воркспейсы не должны лежать внутри чужого git-дерева.
- Подкоманды: `init` (каталоги + `git clone` + `git checkout reference` +
  build; `init <branch>` создаёт ветку во всех private-репо), `pull`,
  `push` (commit с общим сообщением по группе), `status`, `diff`, `tag`;
  флаги `--config` (дефолт `Project.yaml`), `--verbose` — каркас на `flag`.
- Клонирование, pull/push и build-команды — `git` через `os/exec` и сеть:
  внешние зависимости; в главе — только конфиг и каталоги, stdlib.
