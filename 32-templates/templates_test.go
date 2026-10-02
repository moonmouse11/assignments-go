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
