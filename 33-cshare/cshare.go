//go:build !solution

// Пакет cshare — глава 33. Реализуй типы и функции так, чтобы прошли
// тесты в cshare_test.go. Условие — в README.md этой папки.
package cshare

import "io"

// Storage — хранилище файлов: перечислить (List, отсортировано),
// найти по части имени (Search), загрузить (Upload читает всё из r
// под полным именем name) и скачать (Download пишет в w). Отсутствующий
// при Download файл — ошибка.
type Storage interface {
	List() ([]string, error)
	Search(namePart string) ([]string, error)
	Upload(name string, r io.Reader) error
	Download(name string, w io.Writer) error
}

// NewMemStorage возвращает Storage поверх map в памяти.
func NewMemStorage() Storage {
	panic("TODO: NewMemStorage")
}

// SyncReport — итог синхронизации: имена файлов, залитых в облако
// (Uploaded) и скачанных из облака (Downloaded), по алфавиту.
type SyncReport struct {
	Uploaded   []string
	Downloaded []string
}

// SyncDir синхронизирует локальный каталог dir (файлы верхнего уровня)
// с хранилищем s: локальные файлы, которых нет в облаке, — Upload;
// файлы облака, которых нет локально, — Download в dir. Совпадающие
// имена не передаются.
func SyncDir(dir string, s Storage) (SyncReport, error) {
	panic("TODO: SyncDir")
}
