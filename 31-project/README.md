# Глава 31. Project — клонирование из шаблона

> По мотивам `project.pdf`. Новый проект из шаблона: копируем дерево
> каталогов и подставляем значения вместо плейсхолдеров `{{NAME}}`.

## Цель главы

Освоить обход дерева каталогов (`filepath.WalkDir`), относительные пути
(`filepath.Rel`) и замену плейсхолдеров при копировании файлов.

## Теория

- `filepath.WalkDir(root, fn)` обходит дерево; fn получает путь и
  `fs.DirEntry` — каталоги и файлы различаются по `d.IsDir()`.
- Относительный путь: `filepath.Rel(tmplDir, p)` → тот же путь в dstDir
  через `filepath.Join(dstDir, rel)`.
- Замена: `strings.ReplaceAll(content, "{{"+name+"}}", val)` по всем
  ключам vars — плейсхолдеры разные, порядок обхода map не важен.
- `{{X}}` — простой маркер, а не выражения `text/template`: замена
  буквальная, никаких циклов и условий.

## Условие

1. `Instantiate` копирует дерево tmplDir (включая поддиректории) в dstDir.
2. В каждом файле вхождения `{{NAME}}` заменяются значением
   `vars["NAME"]` для всех ключей vars.
3. Файлы без плейсхолдеров копируются без изменений.
4. Отсутствующий tmplDir — ошибка.

Wants (необязательно): переименование файлов, meta-макросы
`$(Created)`/`$(Username)`, шаблон из git-репозитория.

## API для реализации

```go
func Instantiate(tmplDir, dstDir string, vars map[string]string) error
```

Сигнатуры и договоры — в `project.go`. Тесты — в `project_test.go`.

## Подсказки

- Сначала `os.Stat(tmplDir)` — понятная ошибка вместо пустого обхода.
- `os.ReadFile`/`os.WriteFile` достаточно: шаблоны — текст.
- Каталоги создавай `os.MkdirAll` — dstDir и вложенные пути сразу.

## Развитие

Настоящий `proj` берёт шаблон из git-репозитория, читает `Clone.cfg`
(INI-макросы с default/help), умеет переименовывать файлы по директивам
`proj: rename` и подставлять регулярки (`proj: subs`). Git и INI —
вне этой главы: только stdlib, без сети.

## Проверка

```
go vet ./31-project/ && go test ./31-project/            # красный: TODO
go vet -tags solution ./31-project/ && go test -tags solution ./31-project/
```

## Эталон

`solution.go` собирается только с тегом `solution`.
