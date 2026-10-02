package search

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFindLine(t *testing.T) {
	cases := []struct {
		line, needle string
		o            Options
		count        int
		found        bool
	}{
		{"Go go GO", "go", Options{}, 3, true},
		{"Go go GO", "go", Options{CaseSensitive: true}, 1, true},
		{"hello world", "xyz", Options{}, 0, false},
		{"abc", "", Options{}, 0, false},
		{`a1b22c333`, `\d+`, Options{Pattern: true}, 3, true},
		{"ABC abc", `abc`, Options{Pattern: true}, 2, true},
		{"ABC abc", `abc`, Options{Pattern: true, CaseSensitive: true}, 1, true},
	}
	for _, tc := range cases {
		count, found := FindLine(tc.line, tc.needle, tc.o)
		if count != tc.count || found != tc.found {
			t.Errorf("FindLine(%q, %q, %+v) = (%d, %v), хотел (%d, %v)",
				tc.line, tc.needle, tc.o, count, found, tc.count, tc.found)
		}
	}
}

func TestSearchFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "note.txt")
	os.WriteFile(path, []byte("Alpha beta\nnothing here\nbeta Beta beta\n"), 0o644)

	hits, err := SearchFile(path, "beta", Options{})
	if err != nil {
		t.Fatalf("SearchFile: %v", err)
	}
	want := []Hit{
		{Line: 1, Text: "Alpha beta", Count: 1},
		{Line: 3, Text: "beta Beta beta", Count: 3},
	}
	if !reflect.DeepEqual(hits, want) {
		t.Errorf("SearchFile без регистра:\nполучил %+v\nхотел  %+v", hits, want)
	}

	hits, err = SearchFile(path, "BETA", Options{CaseSensitive: true})
	if err != nil {
		t.Fatalf("SearchFile: %v", err)
	}
	if len(hits) != 0 {
		t.Errorf("SearchFile с регистром: нашёл %+v, хотел пусто", hits)
	}
}

func TestSearchFilePattern(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	os.WriteFile(path, []byte("id=42 ok\nid=bad\nid=43 ok\n"), 0o644)

	hits, err := SearchFile(path, `id=\d+`, Options{Pattern: true})
	if err != nil {
		t.Fatalf("SearchFile: %v", err)
	}
	want := []Hit{
		{Line: 1, Text: "id=42 ok", Count: 1},
		{Line: 3, Text: "id=43 ok", Count: 1},
	}
	if !reflect.DeepEqual(hits, want) {
		t.Errorf("SearchFile паттерн:\nполучил %+v\nхотел  %+v", hits, want)
	}
}

func TestSearchFileMissing(t *testing.T) {
	if _, err := SearchFile(filepath.Join(t.TempDir(), "nope.txt"), "x", Options{}); err == nil {
		t.Fatal("SearchFile: хотел ошибку для отсутствующего файла, получил nil")
	}
}

func TestListFiles(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "sub"), 0o755)
	os.Mkdir(filepath.Join(dir, "sub", "deep"), 0o755)
	for _, name := range []string{
		"a.txt", "b.md",
		filepath.Join("sub", "c.txt"),
		filepath.Join("sub", "deep", "d.txt"),
	} {
		os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644)
	}

	// Директория без рекурсии — только файлы верхнего уровня.
	files, err := ListFiles([]string{dir}, false)
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("ListFiles без рекурсии: получил %v, хотел 2 файла верхнего уровня", files)
	}

	// С рекурсией — всё дерево, включая файл в поддиректории.
	files, err = ListFiles([]string{dir}, true)
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}
	if len(files) != 4 {
		t.Errorf("ListFiles рекурсия: получил %v, хотел 4 файла", files)
	}
	foundDeep := false
	for _, f := range files {
		if filepath.Base(f) == "d.txt" {
			foundDeep = true
		}
	}
	if !foundDeep {
		t.Error("ListFiles рекурсия: не нашёл файл в поддиректории")
	}

	// Одиночное имя файла — само себя.
	files, err = ListFiles([]string{filepath.Join(dir, "a.txt")}, false)
	if err != nil || len(files) != 1 || filepath.Base(files[0]) != "a.txt" {
		t.Errorf("ListFiles файл: получил %v, err %v", files, err)
	}

	// Wildcard-имя — раскрытие через glob.
	files, err = ListFiles([]string{filepath.Join(dir, "*.txt")}, false)
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}
	if len(files) != 1 || filepath.Base(files[0]) != "a.txt" {
		t.Errorf("ListFiles *.txt: получил %v, хотел a.txt", files)
	}

	// Несуществующее имя без совпадений — ошибка.
	if _, err := ListFiles([]string{filepath.Join(dir, "no-such-*")}, false); err == nil {
		t.Error("ListFiles: хотел ошибку для несуществующего имени, получил nil")
	}
}
