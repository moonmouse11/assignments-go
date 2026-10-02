package cshare

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestMemStorageRoundTrip(t *testing.T) {
	s := NewMemStorage()
	data := []byte("содержимое файла a")
	if err := s.Upload("a.txt", bytes.NewReader(data)); err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if err := s.Upload("b.md", bytes.NewReader([]byte("байты b"))); err != nil {
		t.Fatalf("Upload: %v", err)
	}
	names, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if want := []string{"a.txt", "b.md"}; !reflect.DeepEqual(names, want) {
		t.Errorf("List = %v, хотел %v (сортировка)", names, want)
	}
	found, err := s.Search("b.")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if want := []string{"b.md"}; !reflect.DeepEqual(found, want) {
		t.Errorf("Search(\"b.\") = %v, хотел %v", found, want)
	}
	if none, err := s.Search("zzz"); err != nil || len(none) != 0 {
		t.Errorf("Search(\"zzz\") = %v (err=%v), хотел пусто", none, err)
	}
	var buf bytes.Buffer
	if err := s.Download("a.txt", &buf); err != nil {
		t.Fatalf("Download: %v", err)
	}
	if !bytes.Equal(buf.Bytes(), data) {
		t.Errorf("Download: получил %q, хотел %q", buf.Bytes(), data)
	}
	if err := s.Download("nope.txt", &buf); err == nil {
		t.Error("Download(nope.txt): хотел ошибку для отсутствующего файла, получил nil")
	}
}

func TestSyncDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "local-only.txt"), []byte("локальный"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "shared.txt"), []byte("совместный"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := NewMemStorage()
	if err := s.Upload("cloud-only.txt", bytes.NewReader([]byte("облачный"))); err != nil {
		t.Fatal(err)
	}
	if err := s.Upload("shared.txt", bytes.NewReader([]byte("совместный"))); err != nil {
		t.Fatal(err)
	}

	rep, err := SyncDir(dir, s)
	if err != nil {
		t.Fatalf("SyncDir: %v", err)
	}
	if want := []string{"local-only.txt"}; !reflect.DeepEqual(rep.Uploaded, want) {
		t.Errorf("Uploaded = %v, хотел %v", rep.Uploaded, want)
	}
	if want := []string{"cloud-only.txt"}; !reflect.DeepEqual(rep.Downloaded, want) {
		t.Errorf("Downloaded = %v, хотел %v", rep.Downloaded, want)
	}

	// Скачанный файл появился локально с содержимым облака.
	got, err := os.ReadFile(filepath.Join(dir, "cloud-only.txt"))
	if err != nil {
		t.Fatalf("файл cloud-only.txt не появился после SyncDir: %v", err)
	}
	if string(got) != "облачный" {
		t.Errorf("cloud-only.txt = %q, хотел %q", got, "облачный")
	}

	// Совпадающий файл не передан никуда: локально не тронут,
	// в облаке — прежнее содержимое.
	localShared, err := os.ReadFile(filepath.Join(dir, "shared.txt"))
	if err != nil || string(localShared) != "совместный" {
		t.Errorf("shared.txt = %q (err=%v), хотел %q", localShared, err, "совместный")
	}
	var cloudShared bytes.Buffer
	if err := s.Download("shared.txt", &cloudShared); err != nil || cloudShared.String() != "совместный" {
		t.Errorf("shared.txt в облаке = %q (err=%v)", cloudShared.String(), err)
	}

	// После синхронизации повтор не должен ничего передавать.
	rep2, err := SyncDir(dir, s)
	if err != nil {
		t.Fatalf("повторный SyncDir: %v", err)
	}
	if len(rep2.Uploaded) != 0 || len(rep2.Downloaded) != 0 {
		t.Errorf("повторный SyncDir передал что-то: %+v", rep2)
	}
}
