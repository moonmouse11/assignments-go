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
