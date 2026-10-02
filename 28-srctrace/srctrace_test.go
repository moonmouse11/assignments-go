package srctrace

import (
	"strings"
	"testing"
)

var sample = RepoInfo{
	Branch: "master",
	Short:  "cefa722",
	Commit: "cefa72267cc4d3a07fbf5e672b0053116d582aa7",
	Date:   "2026-10-01T20:18:00Z",
}

func TestSnippetContainsValues(t *testing.T) {
	src := Snippet(sample)
	for _, v := range []string{sample.Branch, sample.Short, sample.Commit, sample.Date} {
		if !strings.Contains(src, `"`+v+`"`) {
			t.Errorf("Snippet не содержит значение %q:\n%s", v, src)
		}
	}
	if !strings.Contains(src, "package ") {
		t.Error("Snippet без строки package — не Go-файл:\n", src)
	}
	for _, name := range []string{"Branch", "ShortCommit", "Commit", "BuildDate"} {
		if !strings.Contains(src, name) {
			t.Errorf("Snippet без константы %s:\n%s", name, src)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	src := Snippet(sample)
	got, err := Parse(src)
	if err != nil {
		t.Fatalf("Parse(Snippet(...)): %v", err)
	}
	if got != sample {
		t.Errorf("round-trip:\nполучил %+v\nхотел  %+v", got, sample)
	}
}

func TestParseErrors(t *testing.T) {
	bad := []struct {
		name, src string
	}{
		{"пустой источник", ""},
		{"только пробелы", "   \n\t\n"},
		{"без package", "const Branch = \"main\"\n"},
		{"без констант", "package revisions\n"},
		{"не хватает Commit", "package revisions\n\nconst (\n\tBranch = \"main\"\n\tShortCommit = \"a\"\n\tBuildDate = \"d\"\n)\n"},
		{"значение без кавычек", "package revisions\nBranch = без кавычек\n"},
	}
	for _, tc := range bad {
		if _, err := Parse(tc.src); err == nil {
			t.Errorf("Parse(%s): хотел ошибку, получил nil", tc.name)
		}
	}
}
