# Глава 32. Templates — генератор заготовок

> Каждый новый проект начинается одинаково:
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

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./32-templates/
```

## Заготовки (реализуйте)

### `templates.go`

{% raw %}
```go
//go:build !solution

// Пакет templates — глава 32. Реализуй функцию так, чтобы прошли тесты
// в templates_test.go. Условие — в README.md этой папки.
package templates

// Scaffold возвращает файлы заготовки проекта name:
//
//	kind "cli" — "main.go" (парсинг флагов -h/--help, -v/--version,
//	             --verbose по конвенциям главы 02) и "Makefile";
//	kind "web" — "main.go" (http.HandleFunc("/", ...) +
//	             ListenAndServe) и "go.mod" (module name).
//
// Неизвестный kind или пустое name — ошибка.
func Scaffold(kind, name string) (map[string]string, error) {
	panic("TODO: Scaffold")
}
```
{% endraw %}

## Тесты

### `templates_test.go`

{% raw %}
```go
package templates

import (
	"strings"
	"testing"
)

func TestScaffoldCLI(t *testing.T) {
	files, err := Scaffold("cli", "mytool")
	if err != nil {
		t.Fatalf("Scaffold(cli): %v", err)
	}
	main, ok := files["main.go"]
	if !ok {
		t.Fatal("Scaffold(cli): нет main.go")
	}
	if _, ok := files["Makefile"]; !ok {
		t.Fatal("Scaffold(cli): нет Makefile")
	}
	for _, s := range []string{"mytool", "help", "version", "verbose", "flag"} {
		if !strings.Contains(main, s) {
			t.Errorf("cli main.go не содержит %q:\n%s", s, main)
		}
	}
	if !strings.Contains(files["Makefile"], "mytool") {
		t.Errorf("Makefile не содержит имя проекта:\n%s", files["Makefile"])
	}
}

func TestScaffoldWeb(t *testing.T) {
	files, err := Scaffold("web", "myweb")
	if err != nil {
		t.Fatalf("Scaffold(web): %v", err)
	}
	main, ok := files["main.go"]
	if !ok {
		t.Fatal("Scaffold(web): нет main.go")
	}
	gomod, ok := files["go.mod"]
	if !ok {
		t.Fatal("Scaffold(web): нет go.mod")
	}
	for _, s := range []string{"myweb", "http.HandleFunc", "ListenAndServe"} {
		if !strings.Contains(main, s) {
			t.Errorf("web main.go не содержит %q:\n%s", s, main)
		}
	}
	if !strings.Contains(gomod, "module myweb") {
		t.Errorf("go.mod без module myweb:\n%s", gomod)
	}
}

func TestScaffoldErrors(t *testing.T) {
	if _, err := Scaffold("lib", "x"); err == nil {
		t.Error("Scaffold(lib): хотел ошибку для неизвестного типа, получил nil")
	}
	if _, err := Scaffold("cli", ""); err == nil {
		t.Error("Scaffold(cli, \"\"): хотел ошибку для пустого имени, получил nil")
	}
	if _, err := Scaffold("web", ""); err == nil {
		t.Error("Scaffold(web, \"\"): хотел ошибку для пустого имени, получил nil")
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

// Эталонное решение главы 32. Сравни со своей реализацией в templates.go.
package templates

import "fmt"

func Scaffold(kind, name string) (map[string]string, error) {
	if name == "" {
		return nil, fmt.Errorf("пустое имя проекта")
	}
	switch kind {
	case "cli":
		return map[string]string{
			"main.go":  cliMain(name),
			"Makefile": makefile(name),
		}, nil
	case "web":
		return map[string]string{
			"main.go": webMain(name),
			"go.mod":  goMod(name),
		}, nil
	default:
		return nil, fmt.Errorf("неизвестный тип шаблона %q (cli или web)", kind)
	}
}

func cliMain(name string) string {
	return fmt.Sprintf(`// Утилита %[1]s — заготовка CLI по конвенциям главы 02.
package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	verbose := flag.Bool("verbose", false, "подробный вывод")
	showVersion := flag.Bool("version", false, "версия утилиты")
	showVersionShort := flag.Bool("v", false, "версия утилиты (короткая форма)")
	showHelp := flag.Bool("help", false, "помощь по утилите")
	showHelpShort := flag.Bool("h", false, "помощь (короткая форма)")
	flag.Parse()
	switch {
	case *showHelp || *showHelpShort:
		fmt.Printf("Usage: %[1]s [flags]\n", %[1]q)
		flag.PrintDefaults()
	case *showVersion || *showVersionShort:
		fmt.Printf("%[1]s %%s\n", version)
	}
	if *verbose {
		fmt.Fprintln(os.Stderr, "verbose: аргументы:", flag.Args())
	}
}
`, name)
}

func makefile(name string) string {
	return fmt.Sprintf(`# Makefile для %[1]s
build:
	go build -o %[1]s .

test:
	go test ./...

clean:
	rm -f %[1]s
`, name)
}

func webMain(name string) string {
	return fmt.Sprintf(`// Веб-сервер %[1]s — заготовка.
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Это %%s.\n", %[1]q)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
`, name)
}

func goMod(name string) string {
	return fmt.Sprintf("module %s\n\ngo 1.26\n", name)
}
```
{% endraw %}


</details>

_Пакет `assignments-go/32-templates`, упражнение **templates**. Локально: `go test -tags solution ./32-templates/`._
