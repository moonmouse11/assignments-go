package codex

import (
	"strconv"
	"strings"
	"testing"
)

func TestParseSpecsValid(t *testing.T) {
	fs, err := ParseSpecs([]string{"1:5", "34:9", "71:11"})
	if err != nil {
		t.Fatalf("ParseSpecs: %v", err)
	}
	want := []Fragment{{From: 1, Count: 5}, {From: 34, Count: 9}, {From: 71, Count: 11}}
	if len(fs) != len(want) {
		t.Fatalf("ParseSpecs: фрагментов %d, хотел %d", len(fs), len(want))
	}
	for i := range want {
		if fs[i] != want[i] {
			t.Errorf("ParseSpecs: [%d] = %+v, хотел %+v", i, fs[i], want[i])
		}
	}
	if fs, err = ParseSpecs(nil); err != nil || len(fs) != 0 {
		t.Errorf("ParseSpecs(без аргументов) = %+v, %v; хотел пусто и nil", fs, err)
	}
}

func TestParseSpecsErrors(t *testing.T) {
	bad := []string{"34", "a:9", "34:x", "0:9", "34:0", "34:-1", "1:2:3", ":"}
	for _, arg := range bad {
		fs, err := ParseSpecs([]string{"1:5", arg})
		if err == nil {
			t.Errorf("ParseSpecs(%q): хотел ошибку, получил %+v", arg, fs)
			continue
		}
		// Ошибка обязана называть виновника.
		if !strings.Contains(err.Error(), arg) {
			t.Errorf("ParseSpecs(%q): ошибка %q не называет аргумент", arg, err)
		}
	}
}

func testLines(t *testing.T) []string {
	t.Helper()
	var lines []string
	for i := 1; i <= 10; i++ {
		lines = append(lines, "L"+strconv.Itoa(i))
	}
	return lines
}

func TestExtract(t *testing.T) {
	lines := testLines(t)

	got, err := Extract(lines, Fragment{From: 3, Count: 4})
	if err != nil {
		t.Fatalf("Extract(3:4): %v", err)
	}
	if want := "L3|L4|L5|L6"; strings.Join(got, "|") != want {
		t.Errorf("Extract(3:4) = %q, хотел %q", strings.Join(got, "|"), want)
	}

	// Последний фрагмент: до конца файла.
	got, err = Extract(lines, Fragment{From: 7, Count: 4})
	if err != nil {
		t.Fatalf("Extract(7:4): %v", err)
	}
	if want := "L7|L8|L9|L10"; strings.Join(got, "|") != want {
		t.Errorf("Extract(7:4) = %q, хотел %q", strings.Join(got, "|"), want)
	}

	// Одна строка.
	got, err = Extract(lines, Fragment{From: 10, Count: 1})
	if err != nil || len(got) != 1 || got[0] != "L10" {
		t.Errorf("Extract(10:1) = %q, %v; хотел [L10], nil", got, err)
	}
}

func TestExtractOutOfBounds(t *testing.T) {
	lines := testLines(t)
	for _, f := range []Fragment{
		{From: 8, Count: 4},  // 8+4-1 = 11 > 10
		{From: 11, Count: 1}, // начинается за файлом
		{From: 0, Count: 1},  // 1-индексация: 0 — не строка
		{From: 5, Count: 0},  // пустой фрагмент запрещён
		{From: 5, Count: -1},
	} {
		if got, err := Extract(lines, f); err == nil {
			t.Errorf("Extract(%+v) = %q, хотел ошибку за границей файла из %d строк", f, got, len(lines))
		}
	}
}
