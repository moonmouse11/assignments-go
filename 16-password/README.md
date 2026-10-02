# Глава 16. Password — менеджер паролей

> По мотивам `password.pdf`. Поверх защищённого файла из главы 15 строим
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
