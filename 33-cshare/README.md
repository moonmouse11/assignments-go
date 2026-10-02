# Глава 33. Cshare — облачное хранилище через интерфейс

> По мотивам `cshare.pdf`. Облаков много (Drive, iCloud, OneDrive), API у
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

## Развитие

Реальные провайдеры: OAuth-инициализация (креды в `~/.cshare/`),
сохранённые токены (без повторного входа), выбор `--provider`, докачка
расходящихся по размеру файлов. Сеть и авторизация — вне этой главы.

## Проверка

```
go vet ./33-cshare/ && go test ./33-cshare/            # красный: TODO
go vet -tags solution ./33-cshare/ && go test -tags solution ./33-cshare/
```

## Эталон

`solution.go` собирается только с тегом `solution`.
