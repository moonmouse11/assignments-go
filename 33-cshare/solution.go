//go:build solution

// Эталонное решение главы 33. Сравни со своей реализацией в cshare.go.
package cshare

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Storage — контракт хранилища файлов.
type Storage interface {
	List() ([]string, error)
	Search(namePart string) ([]string, error)
	Upload(name string, r io.Reader) error
	Download(name string, w io.Writer) error
}

// SyncReport — итог синхронизации.
type SyncReport struct {
	Uploaded   []string
	Downloaded []string
}

// memStorage — Storage поверх map «имя → содержимое».
type memStorage struct {
	files map[string][]byte
}

func NewMemStorage() Storage {
	return &memStorage{files: map[string][]byte{}}
}

func (m *memStorage) List() ([]string, error) {
	return m.matching("")
}

func (m *memStorage) Search(namePart string) ([]string, error) {
	return m.matching(namePart)
}

func (m *memStorage) matching(part string) ([]string, error) {
	var names []string
	for n := range m.files {
		if part == "" || strings.Contains(n, part) {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	return names, nil
}

func (m *memStorage) Upload(name string, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	m.files[name] = data
	return nil
}

func (m *memStorage) Download(name string, w io.Writer) error {
	data, ok := m.files[name]
	if !ok {
		return fmt.Errorf("файла %q нет в хранилище", name)
	}
	_, err := w.Write(data)
	return err
}

// SyncDir доливает локальные файлы в облако и докачивает облачные
// в каталог; совпадающие имена не трогает.
func SyncDir(dir string, s Storage) (SyncReport, error) {
	var rep SyncReport
	entries, err := os.ReadDir(dir)
	if err != nil {
		return rep, err
	}
	local := map[string]bool{}
	for _, e := range entries {
		if !e.IsDir() {
			local[e.Name()] = true
		}
	}
	remote, err := s.List()
	if err != nil {
		return rep, err
	}
	cloud := map[string]bool{}
	for _, n := range remote {
		cloud[n] = true
	}
	for _, name := range sorted(local) {
		if cloud[name] {
			continue
		}
		f, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			return rep, err
		}
		err = s.Upload(name, f)
		f.Close()
		if err != nil {
			return rep, err
		}
		rep.Uploaded = append(rep.Uploaded, name)
	}
	for _, name := range sorted(cloud) {
		if local[name] {
			continue
		}
		f, err := os.Create(filepath.Join(dir, name))
		if err != nil {
			return rep, err
		}
		err = s.Download(name, f)
		f.Close()
		if err != nil {
			return rep, err
		}
		rep.Downloaded = append(rep.Downloaded, name)
	}
	return rep, nil
}

func sorted(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for n := range set {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}
