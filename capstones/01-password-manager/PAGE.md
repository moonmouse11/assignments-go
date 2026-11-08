# Капстоун 1. Password Manager — свой менеджер паролей

> Собирает [`10-pwdgen`](../../10-pwdgen/README.md) и [`16-password`](../../16-password/README.md)

## Цель капстоуна

Интеграция: генератор из главы 10 выдаёт пароль, хранилище из главы 16
запоминает его за контекстом и прячет в защищённый файл. Твоя работа —
только клей `Manager`; генерация, шифрование и сериализация — библиотеки глав.

## Теория

- Импорт соседних пакетов модуля: `import "assignments-go/16-password"`
  (пакет называется `password`), `"assignments-go/10-pwdgen"` (пакет `pwdgen`).
- Заготовки глав паникуют `TODO`, пока не решены, — капстоун честно красный,
  пока красные главы. С `-tags solution` всё собирается из эталонов.
- Ошибка `Open` с неверным master-паролем приходит из главы 16 — не дублируй её.

## Условие

1. `New(seed, words, master)` — менеджер со свежим хранилищем и генератором
   (разделитель сегментов — `-`).
2. `Generate(context, username, nSeg)` — пароль из `nSeg` сегментов генератора;
   кладётся в хранилище и возвращается. Пустой `context` — ошибка (от хранилища).
3. `Get(context)` — пароль контекста.
4. `Save(path)` — хранилище в защищённый файл с master-паролем менеджера.
5. `Open(path, master)` — менеджер из файла; неверный master — ошибка.
   У открытого менеджера генератора нет: `Generate` возвращает ошибку.
6. `List()` — записи по алфавиту контекстов.

## API для реализации

```go
type Manager struct{ /* TODO */ }
func New(seed uint64, words []string, master string) *Manager
func (m *Manager) Generate(context, username string, nSeg int) (string, error)
func (m *Manager) Get(context string) (string, bool)
func (m *Manager) Save(path string) error
func Open(path, master string) (*Manager, error)
func (m *Manager) List() []password.Entry
```

## Подсказки

- Поля Manager: хранилище, генератор (может быть nil), master.
- `Generate`: `gen.Password(nSeg)` → `vault.Add(Entry{...})` → вернуть пароль.
- Тесты сравнивают пароли до/после `Save`+`Open` — детерминизм генератора
  по seed тут не нужен, сравнивай то, что вернул `Generate`.

## Проверка

```
go test ./capstones/01-password-manager/
go test -tags solution ./capstones/01-password-manager/
```

## Эталон

`go test -tags solution ./capstones/01-password-manager/`

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./capstones/01-password-manager/
```

## Заготовки (реализуйте)

### `manager.go`

{% raw %}
```go
//go:build !solution

// Капстоун 1: менеджер паролей. Реализуй клей поверх 10-pwdgen и 16-password
// так, чтобы прошли тесты в manager_test.go. Условие — в README.md.
package passwordmanager

import (
	"assignments-go/10-pwdgen"
	"assignments-go/16-password"
)

// var _ — держит импорт pwdgen: понадобится в New (см. решение).
var _ = pwdgen.NewGen

// Manager — менеджер паролей: хранилище (16-password) + генератор (10-pwdgen)
// + master-пароль. У открытого из файла генератора нет.
type Manager struct {
	// TODO: твои поля.
}

// New создаёт менеджер: свежее хранилище и генератор (разделитель "-").
func New(seed uint64, words []string, master string) *Manager {
	panic("TODO: New")
}

// Generate генерирует пароль из nSeg сегментов, кладёт его в хранилище
// под context/username и возвращает. Пустой context — ошибка.
// У менеджера без генератора (Open) — ошибка.
func (m *Manager) Generate(context, username string, nSeg int) (string, error) {
	panic("TODO: Generate")
}

// Get возвращает пароль контекста.
func (m *Manager) Get(context string) (string, bool) {
	panic("TODO: Get")
}

// Save пишет хранилище в защищённый файл path с master-паролем менеджера.
func (m *Manager) Save(path string) error {
	panic("TODO: Save")
}

// Open читает менеджер из защищённого файла. Неверный master — ошибка.
func Open(path, master string) (*Manager, error) {
	panic("TODO: Open")
}

// List возвращает записи по алфавиту контекстов.
func (m *Manager) List() []password.Entry {
	panic("TODO: List")
}
```
{% endraw %}

## Тесты

### `manager_test.go`

{% raw %}
```go
package passwordmanager

import (
	"path/filepath"
	"testing"
)

var words = []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot"}

func TestGenerateSaveOpen(t *testing.T) {
	m := New(42, words, "master-key")
	pw1, err := m.Generate("example.com", "alice", 3)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	pw2, err := m.Generate("gitlab.com", "bob", 3)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	path := filepath.Join(t.TempDir(), "vault.ppm")
	if err := m.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	opened, err := Open(path, "master-key")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	for ctx, want := range map[string]string{"example.com": pw1, "gitlab.com": pw2} {
		if got, ok := opened.Get(ctx); !ok || got != want {
			t.Errorf("Get(%q) после Open: получил %q, ok=%v, хотел %q", ctx, got, ok, want)
		}
	}
}

func TestOpenWrongMaster(t *testing.T) {
	m := New(42, words, "master-key")
	if _, err := m.Generate("example.com", "alice", 2); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	path := filepath.Join(t.TempDir(), "vault.ppm")
	if err := m.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := Open(path, "не-тот-master"); err == nil {
		t.Error("Open с неверным master: хотел ошибку, получил nil")
	}
}

func TestGenerateEmptyContext(t *testing.T) {
	m := New(42, words, "master-key")
	if _, err := m.Generate("", "alice", 2); err == nil {
		t.Error("Generate с пустым context: хотел ошибку")
	}
}

func TestOpenedManagerHasNoGenerator(t *testing.T) {
	m := New(42, words, "master-key")
	if _, err := m.Generate("example.com", "alice", 2); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	path := filepath.Join(t.TempDir(), "vault.ppm")
	if err := m.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	opened, err := Open(path, "master-key")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if _, err := opened.Generate("new.com", "carol", 2); err == nil {
		t.Error("Generate у открытого менеджера: хотел ошибку (нет генератора)")
	}
}

func TestListSorted(t *testing.T) {
	m := New(42, words, "master-key")
	m.Generate("zzz.example", "a", 2)
	m.Generate("aaa.example", "b", 2)
	list := m.List()
	if len(list) != 2 {
		t.Fatalf("List: записей %d, хотел 2", len(list))
	}
	if list[0].Context != "aaa.example" || list[1].Context != "zzz.example" {
		t.Errorf("List не отсортирован: %v, %v", list[0].Context, list[1].Context)
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

// Эталон капстоуна 1.
package passwordmanager

import (
	"fmt"

	"assignments-go/10-pwdgen"
	"assignments-go/16-password"
)

type Manager struct {
	vault  *password.Vault
	gen    *pwdgen.Gen
	master string
}

func New(seed uint64, words []string, master string) *Manager {
	return &Manager{
		vault:  password.NewVault(),
		gen:    pwdgen.NewGen(seed, words, "-"),
		master: master,
	}
}

func (m *Manager) Generate(context, username string, nSeg int) (string, error) {
	if m.gen == nil {
		return "", fmt.Errorf("генератор не настроен (менеджер открыт из файла)")
	}
	pw := m.gen.Password(nSeg)
	err := m.vault.Add(password.Entry{Context: context, Username: username, Password: pw})
	if err != nil {
		return "", err
	}
	return pw, nil
}

func (m *Manager) Get(context string) (string, bool) {
	e, ok := m.vault.Get(context)
	if !ok {
		return "", false
	}
	return e.Password, true
}

func (m *Manager) Save(path string) error {
	return m.vault.SaveTo(path, m.master)
}

func Open(path, master string) (*Manager, error) {
	v, err := password.LoadFrom(path, master)
	if err != nil {
		return nil, err
	}
	return &Manager{vault: v, master: master}, nil
}

func (m *Manager) List() []password.Entry {
	return m.vault.List()
}
```
{% endraw %}


</details>

_Пакет `assignments-go/capstones/01-password-manager`, упражнение **password-manager**. Локально: `go test -tags solution ./capstones/01-password-manager/`._
