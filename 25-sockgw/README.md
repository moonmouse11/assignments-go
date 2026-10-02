# Глава 25. Sockgw — шлюз stream-socket → UDP

> По мотивам `sockgw.pdf`. Реальная история: Qt-приложение тонет в потоках
> данных от нижних контроллеров; логирование выносят в отдельный
> процесс-шлюз. Мы строим ядро такого шлюза: TCP-сервер, лог-файл,
> пересылка записей по UDP.

## Цель главы

Первое сетевое приложение: слушать TCP, разбирать байтовый поток на записи
по `\n`, писать лог и пересылать датаграммы — всё останавливаемо через
`context`.

## Теория

- `net.Listen("tcp", "127.0.0.1:0")`: порт 0 — «дай свободный»; реальный
  адрес узнаётся через `ln.Addr()`.
- Stream-сокет — **байтовый поток**: границ записей в нём нет, записи
  восстанавливаются сами (`bufio.Reader.ReadString('\n')`).
- UDP (`net.DialUDP`, `net.ListenPacket`) — датаграммы: сообщение целиком,
  но без гарантий доставки; на localhost — надёжно.
- `Serve` блокируется до отмены ctx; accept — в своей горутине, каждое
  соединение — в своей.
- Таймауты обязательны: `SetReadDeadline`/`SetDeadline` на сокетах, иначе
  тест повиснет навечно при ошибке.

## Условие

Шлюз: принимает TCP-соединения, каждую принятую строку-запись пишет в
лог-файл (с `\n`) и пересылает по UDP на заданный адрес.

1. `Listen` поднимает TCP-listener; пустая строка адреса = `127.0.0.1:0`.
2. `SetForward` задаёт UDP-адрес (resolve), `SetLog` — путь лога (файл).
3. `Serve` блокируется: accept в горутине, по строке — лог + UDP-датаграмма.
4. `Close` закрывает listener, лог и UDP-соединение; повторный вызов —
   не ошибка (идемпотентно).

## API для реализации

```go
type Gateway struct{ /* поля на твой выбор */ }
func (g *Gateway) Listen(tcpAddr string) error
func (g *Gateway) Addr() net.Addr
func (g *Gateway) SetForward(udpAddr string) error
func (g *Gateway) SetLog(path string) error
func (g *Gateway) Serve(ctx context.Context)
func (g *Gateway) Close()
```

Сигнатуры и договоры — в `sockgw.go`. Тесты — в `sockgw_test.go`.

## Подсказки

- UDP-приёмник в тесте — `net.ListenPacket("udp", "127.0.0.1:0")`,
  адрес для `SetForward` — `pc.LocalAddr().String()`.
- Пиши лог **до** отправки UDP: получив датаграмму, тест уже знает,
  что строка попала в лог.
- Мьютекс защитит `Close` от гонки с записью лога.
- Все адреса — только localhost.

## Проверка

```
go test ./25-sockgw/
go test -tags solution ./25-sockgw/
```

## Эталон

`go test -tags solution ./25-sockgw/`
