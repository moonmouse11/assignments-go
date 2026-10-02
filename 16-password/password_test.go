package password

import (
	"path/filepath"
	"reflect"
	"testing"
)

func sampleVault(t *testing.T) *Vault {
	t.Helper()
	v := NewVault()
	entries := []Entry{
		{Context: "z-site", Username: "zed", Password: "p3"},
		{Context: "mail", Username: "me", Password: "p1"},
		{Context: "git", Username: "dev", Password: "p2"},
	}
	for _, e := range entries {
		if err := v.Add(e); err != nil {
			t.Fatalf("Add(%q): %v", e.Context, err)
		}
	}
	return v
}

func TestAddGetList(t *testing.T) {
	v := sampleVault(t)

	got, ok := v.Get("mail")
	if !ok || got != (Entry{Context: "mail", Username: "me", Password: "p1"}) {
		t.Errorf("Get(mail): получил %+v, ok=%v", got, ok)
	}
	if _, ok := v.Get("нет-такого"); ok {
		t.Error("Get: хотел ok=false для отсутствующего контекста")
	}

	// Add обновляет запись того же контекста (UN-5).
	if err := v.Add(Entry{Context: "mail", Username: "me", Password: "новый"}); err != nil {
		t.Fatalf("Add повторный: %v", err)
	}
	if got, _ := v.Get("mail"); got.Password != "новый" {
		t.Errorf("Add не обновил запись: %+v", got)
	}

	// List — по алфавиту контекстов.
	want := []string{"git", "mail", "z-site"}
	list := v.List()
	if len(list) != len(want) {
		t.Fatalf("List: получил %d записей, хотел %d", len(list), len(want))
	}
	for i, e := range list {
		if e.Context != want[i] {
			t.Errorf("List[%d]: контекст %q, хотел %q", i, e.Context, want[i])
		}
	}
}

func TestAddEmptyContext(t *testing.T) {
	v := NewVault()
	if err := v.Add(Entry{Username: "u", Password: "p"}); err == nil {
		t.Fatal("Add: хотел ошибку для пустого контекста, получил nil")
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	v := sampleVault(t)
	path := filepath.Join(t.TempDir(), "vault.dat")

	if err := v.SaveTo(path, "master"); err != nil {
		t.Fatalf("SaveTo: %v", err)
	}
	got, err := LoadFrom(path, "master")
	if err != nil {
		t.Fatalf("LoadFrom: %v", err)
	}
	if !reflect.DeepEqual(got.List(), v.List()) {
		t.Errorf("round-trip: получил %+v, хотел %+v", got.List(), v.List())
	}
}

func TestLoadWrongMaster(t *testing.T) {
	v := sampleVault(t)
	path := filepath.Join(t.TempDir(), "vault.dat")
	if err := v.SaveTo(path, "master"); err != nil {
		t.Fatalf("SaveTo: %v", err)
	}
	if _, err := LoadFrom(path, "не-тот-пароль"); err == nil {
		t.Fatal("LoadFrom: хотел ошибку для неверного мастер-пароля, получил nil")
	}
}
