# Глава 33. Cshare — облачное хранилище через интерфейс

> Облаков много (Drive, iCloud, OneDrive), API у
> всех разный — код не должен зависеть от конкретного. Тренируемся на
> `Storage` в памяти.

## Цель главы

Спроектировать интерфейс `Storage` и написать алгоритм синхронизации
каталога с хранилищем, работающий с любой реализацией — полиморфизм
вместо привязки к провайдеру.

## Теория

- Интерфейс — контракт: `List`/`Search`/`Upload`/`Download`; реализаций
  может быть сколько угодно, `MemStorage` — первая, для тестов.
- Синхронизация — сравнение множеств имён: локальное минус облачное →
  `Upload`, облачное минус локальное → `Download`.
- Полиморфизм: `SyncDir` принимает `Storage`, а не конкретный тип —
  тест пройдёт с любой реализацией.
- Имена сортируем (`sort.Strings`): отчёт детерминирован, тесты стабильны.

## Условие

1. `NewMemStorage` — `Storage` в памяти: `Upload` пишет всё из `r` под
   полным именем; `List` и `Search` возвращают отсортированные имена
   (`Search` — содержащие `namePart`); `Download` отсутствующего — ошибка.
2. `SyncDir(dir, s)` сравнивает локальные файлы dir (без подкаталогов)
   и `List()`: локальные без пары в облаке — `Upload` (содержимое файла);
   облачные без локальной пары — `Download` в dir.
3. Совпадающие имена не передаются никуда.
4. `SyncReport` содержит списки `Uploaded` и `Downloaded` по алфавиту.

Wants (необязательно): сравнение размеров по `--sync`, провайдеры
google/apple/microsoft/aws.

## API для реализации

```go
type Storage interface {
	List() ([]string, error)
	Search(namePart string) ([]string, error)
	Upload(name string, r io.Reader) error
	Download(name string, w io.Writer) error
}
func NewMemStorage() Storage
type SyncReport struct{ Uploaded, Downloaded []string }
func SyncDir(dir string, s Storage) (SyncReport, error)
```

Сигнатуры и договоры — в `cshare.go`. Тесты — в `cshare_test.go`.

## Подсказки

- Множества имён — `map[string]bool`; сортировка — `sort.Strings`.
- Для Upload открывай файл `os.Open` и закрывай после; для Download —
  `os.Create` в dir.
- `io.ReadAll` в Upload — всё содержимое из ридера разом.

## Проверка

```
go vet ./33-cshare/ && go test ./33-cshare/            # красный: TODO
go vet -tags solution ./33-cshare/ && go test -tags solution ./33-cshare/
```

## Эталон

`solution.go` собирается только с тегом `solution`.

## Углубление

- Подкоманды: `initialize` (креды в `~/.cshare/`, каталог создать при
  отсутствии — `os.MkdirAll`), `list <substr>` (дефолт `mp3`),
  `download <substr>` с `--local-store` (дефолт `.`), `upload <файлы...>`
  — каркас CLI на stdlib `flag`.
- `initialize --provider google|apple|microsoft|aws` (дефолт google),
  `--refresh` — забыть сохранённые креды; после первой авторизации
  токены переиспользуются, повторный вход не нужен.
- `--sync` при download: сравнить размеры локального и облачного файла,
  передавать только при расхождении (`os.Stat` + размер из Storage).
- Реальные провайдеры Google Drive/iCloud/OneDrive/AWS — OAuth и сеть:
  внешние зависимости (API-клиенты); в главе — `MemStorage`, stdlib.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./33-cshare/
```

## Заготовки (реализуйте)

### `cshare.go`

{% raw %}
```go
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
```
{% endraw %}

## Тесты

### `cshare_test.go`

{% raw %}
```go
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
```
{% endraw %}

## Эталонное решение

Сначала попробуйте сами. Решение скрыто — раскрывайте только для сверки.

<details>
<summary><strong>Эталон (`solution.go`)</strong> — нажмите, чтобы показать</summary>

{% raw %}
```go
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
```
{% endraw %}


</details>

_Пакет `assignments-go/33-cshare`, упражнение **cshare**. Локально: `go test -tags solution ./33-cshare/`._
