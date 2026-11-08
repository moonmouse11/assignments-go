# Глава 16. Password — менеджер паролей

> Поверх защищённого файла из главы 15 строим
> менеджер паролей: хранилище контекст/логин/пароль под мастер-паролем.

## Цель главы

Собрать криптографиию в приложение: связка «структуры + JSON + KDF + AES-GCM
+ файл» — типичный каркас реального password manager.

## Теория

- Каждая запись — тройка: контекст (сайт/сервис), имя пользователя, пароль.
  Ключ поиска — контекст.
- Вся база сериализуется в JSON и шифруется целиком: по файлу снаружи
  не видно ни числа записей, ни контекстов.
- KDF и AES-GCM реализуем локально в пакете (pbkdf2-hmac-sha256 + GCM),
 salt/nonce случайные и лежат в заголовке файла рядом с шифротекстом.

## Условие

1. `Add` добавляет запись; пустой контекст — ошибка; повторный `Add`
   того же контекста обновляет запись (UN-4 и UN-5).
2. `Get` возвращает запись по контексту; `List` — все записи по алфавиту
   контекстов.
3. `SaveTo`/`LoadFrom` пишут/читают шифрованный JSON-файл под мастер-паролем.
4. Любая порча файла и неверный мастер-пароль должны детектироваться
   (UN-2, UN-8).

## API для реализации

```go
type Entry struct{ Context, Username, Password string }
type Vault struct{ /* entries map[string]Entry */ }
func NewVault() *Vault
func (v *Vault) Add(e Entry) error
func (v *Vault) Get(context string) (Entry, bool)
func (v *Vault) List() []Entry
func (v *Vault) SaveTo(path, master string) error
func LoadFrom(path, master string) (*Vault, error)
```

Договоры — в `password.go`. Тесты — в `password_test.go`.

## Подсказки

- Храни записи в `map[string]Entry`, сортируй ключи в `List`.
- Свой шифрофайл: заголовок `PSWD` + версия + salt + nonce, затем GCM.
  Итерации KDF — константа 1000, чтобы тесты бежали <1с.
- `encoding/json` сортирует ключи map сам — файл получается стабильным.

## Проверка

```
go test ./16-password/                  # красный: TODO
go test -tags solution ./16-password/  # зелёный: эталон
```

## Эталон

`solution.go` собирается только с тегом `solution`.

## Углубление

- UN-7: несколько identity на один контекст — ключ не `context`, а пара
  `(context, username)`: `map[string]map[string]Entry` и `Get(ctx, user)`.
- UW-2: мастер-пароль вводится по приглашению и никогда — в командной строке:
  промпт через stdin (`bufio`); скрытие ввода — `golang.org/x/term`
  (внешняя зависимость).
- DS-2: CLI `ppm` с флагом `-f` и файлом по умолчанию `~/.ppm/password.dat`
  (`flag`, `os.UserHomeDir`), подкоманды `create/add/update/show/list`;
  при промахе печатается `Unable to find Context ... Username ...`.
- DS-3: в спеке файл — зашифрованный CSV (`encoding/csv`), а не JSON.
- UW-1: пароль без повторного ввода в login-контексте — keyring/агент
  (внешние зависимости, в курсе не обязательны).

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./16-password/
```

## Заготовки (реализуйте)

### `password.go`

{% raw %}
```go
//go:build !solution

// Пакет password — глава 16. Реализуй типы и методы так, чтобы прошли
// тесты в password_test.go. Условие — в README.md этой папки.
package password

// Entry — одна запись менеджера паролей: контекст (сайт/сервис),
// имя пользователя и пароль для него. Ключ — Context.
type Entry struct {
	Context  string
	Username string
	Password string
}

// Vault — хранилище записей. Ключ — Context записи.
type Vault struct {
	// TODO: поле для хранения записей (map по Context).
}

// NewVault возвращает пустое хранилище.
func NewVault() *Vault {
	panic("TODO: NewVault")
}

// Add добавляет запись e. Пустой Context — ошибка. Повторный Add
// того же Context обновляет запись (UN-4, UN-5).
func (v *Vault) Add(e Entry) error {
	panic("TODO: Add")
}

// Get возвращает запись по контексту context; ok=false, если записи нет.
func (v *Vault) Get(context string) (Entry, bool) {
	panic("TODO: Get")
}

// List возвращает все записи, отсортированные по алфавиту контекстов.
func (v *Vault) List() []Entry {
	panic("TODO: List")
}

// SaveTo сериализует хранилище в JSON и записывает в path зашифрованным
// (pbkdf2-hmac-sha256 + AES-GCM, реализация локальная в пакете) файлом
// под мастер-паролем master. Заголовок файла: magic "PSWD", версия,
// salt[16], nonce[12], sha256 заголовка; затем шифротекст.
func (v *Vault) SaveTo(path, master string) error {
	panic("TODO: SaveTo")
}

// LoadFrom читает зашифрованный файл path, расшифровывает его мастер-паролем
// master и возвращает хранилище. Ошибки: неверный формат/портченый файл,
// неверный мастер-пароль.
func LoadFrom(path, master string) (*Vault, error) {
	panic("TODO: LoadFrom")
}
```
{% endraw %}

## Тесты

### `password_test.go`

{% raw %}
```go
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
```
{% endraw %}

## Эталонное решение

Сначала попробуйте сами. Решение скрыто — раскрывайте только для сверки.

<details>
<summary><strong>Эталон (`solution.go`)</strong> — нажмите, чтобы показать</summary>

{% raw %}
```go
//go:build solution

// Эталонное решение главы 16. Сравни со своей реализацией в password.go.
package password

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
)

const (
	magic   = "PSWD"
	headLen = len(magic) + 1 + 16 + 12 + sha256.Size // 65
)

// Entry — одна запись менеджера паролей: контекст, логин, пароль.
type Entry struct{ Context, Username, Password string }

// Vault — хранилище записей (ключ — Context). SaveTo/LoadFrom шифруют
// JSON-слепок хранилища файлом формата PSWD (заголовок + AES-GCM).
type Vault struct{ entries map[string]Entry }

func NewVault() *Vault { return &Vault{entries: map[string]Entry{}} }

func (v *Vault) Add(e Entry) error {
	if e.Context == "" {
		return errors.New("пустой контекст")
	}
	v.entries[e.Context] = e
	return nil
}

func (v *Vault) Get(context string) (Entry, bool) {
	e, ok := v.entries[context]
	return e, ok
}

func (v *Vault) List() []Entry {
	out := make([]Entry, 0, len(v.entries))
	for _, e := range v.entries {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Context < out[j].Context })
	return out
}

// gcm выводит ключ из мастер-пароля (pbkdf2, 1000 итераций) и собирает AES-GCM.
func gcm(master string, salt []byte) (cipher.AEAD, error) {
	key, err := pbkdf2.Key(sha256.New, master, salt, 1000, 32)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func (v *Vault) SaveTo(path, master string) error {
	blob, err := json.Marshal(v.entries)
	if err != nil {
		return err
	}
	head := make([]byte, headLen-sha256.Size) // magic + версия + salt + nonce
	copy(head, magic)
	head[len(magic)] = 1
	if _, err := rand.Read(head[len(magic)+1:]); err != nil {
		return err
	}
	salt, nonce := head[len(magic)+1:len(magic)+1+16], head[len(magic)+1+16:]
	aead, err := gcm(master, salt)
	if err != nil {
		return err
	}
	h := sha256.Sum256(head)
	file := append(append(head, h[:]...), aead.Seal(nil, nonce, blob, nil)...)
	return os.WriteFile(path, file, 0o600)
}

func LoadFrom(path, master string) (*Vault, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256(data[:headLen-sha256.Size])
	if len(data) < headLen || string(data[:len(magic)]) != magic || data[len(magic)] != 1 || string(h[:]) != string(data[headLen-sha256.Size:headLen]) {
		return nil, errors.New("плохой или порченый заголовок файла (длина, magic, версия, sha256)")
	}
	salt, nonce := data[len(magic)+1:len(magic)+1+16], data[len(magic)+1+16:headLen-sha256.Size]
	aead, err := gcm(master, salt)
	if err != nil {
		return nil, err
	}
	blob, err := aead.Open(nil, nonce, data[headLen:], nil)
	if err != nil {
		return nil, fmt.Errorf("расшифровка не удалась (мастер-пароль или файл): %w", err)
	}
	v := NewVault()
	if err := json.Unmarshal(blob, &v.entries); err != nil {
		return nil, fmt.Errorf("повреждённый JSON внутри: %w", err)
	}
	return v, nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/16-password`, упражнение **password**. Локально: `go test -tags solution ./16-password/`._
