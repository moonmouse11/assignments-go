# Глава 32. Templates — генератор заготовок

> По мотивам `templates.pdf`. Каждый новый проект начинается одинаково:
> CLI с --help/--version/--verbose или веб-сервер. Генерируем заготовку
> одной функцией — map «имя файла → содержимое».

## Цель главы

Собрать файлы заготовки для двух типов проектов и закрепить конвенции
CLI из главы 02 в виде кода-шаблона.

## Теория

- Конвенции CLI (глава 02): `-h/--help` — usage, `-v/--version` — версия,
  `--verbose` — диагностика; пароли — не в командной строке.
- Флаги stdlib: `flag.Bool("help", ...)` принимает и `-help`, и `--help`;
  короткие формы — отдельные флаги `"h"`, `"v"`.
- Веб-минимум: `http.HandleFunc("/", ...)` + `http.ListenAndServe`.
- Шаблон — обычная строка, имя проекта вставляется `fmt.Sprintf`;
  `%q` корректно закавычит его в генерируемом коде.

## Условие

1. `Scaffold("cli", name)` возвращает main.go (парсинг флагов -h/--help,
   -v/--version, --verbose по конвенциям главы 02) и Makefile.
2. `Scaffold("web", name)` возвращает main.go (`http.HandleFunc("/", ...)`
   + `ListenAndServe`) и go.mod (`module <name>`).
3. Неизвестный kind или пустое name — ошибка.

Wants (необязательно): https/TLS, база пользователей, команда install,
запись файлов на диск.

## API для реализации

```go
func Scaffold(kind, name string) (map[string]string, error)
```

Сигнатуры и договоры — в `templates.go`. Тесты — в `templates_test.go`.

## Подсказки

- Аргументы проверяй до построения map — ошибки дешевле генерации.
- Шаблоны не интерпретируй: это просто строки, тесты проверяют вхождения
  (имя проекта, help/version/verbose, ListenAndServe, module).

## Проверка

```
go vet ./32-templates/ && go test ./32-templates/            # красный: TODO
go vet -tags solution ./32-templates/ && go test -tags solution ./32-templates/
```

## Эталон

`solution.go` собирается только с тегом `solution`.

## Углубление

- CLI-шаблон автора — с подкомандами (help/version) и `--config` (дефолт
  `$HOME/.cli/cli.yaml`); в веб-шаблоне `--verbose` — уровень 1..16
  (`flag.Int`), а не bool; пароли — не в командной строке (env-переменная).
- Веб-шаблон: подкоманды `install`, `show`, `user`; https/TLS и вход
  обязательны до любых сервисов; admin и обычный пользователь;
  статический html (`http.FileServer`).
- `install` — интерактив (каталог, порт, URL, пароли admin/user) и
  self-signed сертификат: keypair.pvt/pub.pem + certfile — `crypto/x509`,
  `crypto/rsa`, `encoding/pem`, stdlib.
- База пользователей: пароли не в открытом виде — `golang.org/x/crypto/
  bcrypt` (внешняя) или свой salt+hash на `crypto/sha256` (stdlib).
- Шаблоны — отдельные файлы, пишутся на диск (`os.WriteFile`, `embed`);
  у автора живут в git-репозитории.
