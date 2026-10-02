# Капстоун 1. Password Manager — свой менеджер паролей

> Собирает [`10-pwdgen`](../../10-pwdgen/) и [`16-password`](../../16-password/)
> (по мотивам `pwdgen.pdf`, `ppm.pdf`, `password.pdf`).

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
