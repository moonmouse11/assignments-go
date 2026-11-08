# Глава 25. Sockgw — шлюз stream-socket → UDP

> Реальная история: Qt-приложение тонет в потоках
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

## Углубление

- Полный CLI с режимами: `-g` шлюз, `-c` sample-клиент (генератор тестового
  трафика, шлёт записи периодически), `-u` UDP-приёмник, печатающий
  датаграммы, — подрежимы через `flag`/`flag.NewFlagSet`; в главе приёмник
  живёт только в тестах.
- Уровни протоколирования: `-v/--verbose`, `-L/--logging-level` — `log/slog`
  с уровнями и переключением по флагу.
- Имена логов у автора: `archive_YYYYMMDD_HHMMSS.bin` — таймштамп старта
  через `time.Now().Format("20060102_150405")`, по файлу на запуск.
- Статистика по соединениям: «Server No N» и «Forwarded N bytes» на каждую
  запись — счётчики байт на соединение под `sync.Mutex`/`sync/atomic`.
- Эталон автора — Ada-реализация gitlab.com/ada23/sockgw (сеть).

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./25-sockgw/
```

## Заготовки (реализуйте)

### `sockgw.go`

{% raw %}
```go
//go:build !solution

// Пакет sockgw — глава 25. Реализуй методы Gateway так, чтобы прошли
// тесты в sockgw_test.go. Условие — в README.md этой папки.
package sockgw

import (
	"context"
	"net"
)

// Gateway — шлюз stream-socket → UDP: принимает TCP-соединения,
// каждую строку-запись пишет в лог-файл и пересылает UDP-датаграммой.
// Поля — на твой выбор (listener, UDP-адрес, лог-файл, мьютекс для Close).
type Gateway struct{}

// Listen поднимает TCP-listener на tcpAddr (пустая строка → "127.0.0.1:0").
// Возвращает ошибку, если адрес не слушается.
func (g *Gateway) Listen(tcpAddr string) error {
	panic("TODO: Listen")
}

// Addr возвращает адрес TCP-listener'а (nil, если Listen не вызывался).
func (g *Gateway) Addr() net.Addr {
	panic("TODO: Addr")
}

// SetForward задаёт UDP-адрес пересылки (net.ResolveUDPAddr).
// Возвращает ошибку, если адрес не резолвится.
func (g *Gateway) SetForward(udpAddr string) error {
	panic("TODO: SetForward")
}

// SetLog открывает path для записи лога (создаёт при отсутствии).
// Возвращает ошибку, если файл не открывается.
func (g *Gateway) SetLog(path string) error {
	panic("TODO: SetLog")
}

// Serve блокируется до отмены ctx: accept — в горутине, каждое соединение —
// в своей горутине; по каждой строке — запись в лог (с "\n") и
// UDP-датаграмма. По отмене ctx закрывает всё (Close).
func (g *Gateway) Serve(ctx context.Context) {
	panic("TODO: Serve")
}

// Close закрывает listener, лог и UDP-соединение. Идемпотентно:
// повторный вызов не ошибка.
func (g *Gateway) Close() {
	panic("TODO: Close")
}
```
{% endraw %}

## Тесты

### `sockgw_test.go`

{% raw %}
```go
package sockgw

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGatewayErrors(t *testing.T) {
	var g Gateway
	if err := g.SetForward("bogus-addr"); err == nil {
		t.Error("SetForward: хотел ошибку для кривого адреса")
	}
	if err := g.Listen("127.0.0.1:99999"); err == nil {
		t.Error("Listen: хотел ошибку для кривого порта")
	}
	bad := filepath.Join(t.TempDir(), "no", "such", "dir", "l.log")
	if err := g.SetLog(bad); err == nil {
		t.Error("SetLog: хотел ошибку для отсутствующего каталога")
	}
	if g.Addr() != nil {
		t.Error("Addr: хотел nil до Listen")
	}
}

func TestGatewayForward(t *testing.T) {
	// UDP-приёмник.
	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("ListenPacket: %v", err)
	}
	defer pc.Close()

	// Шлюз.
	var g Gateway
	defer g.Close()
	if err := g.Listen("127.0.0.1:0"); err != nil {
		t.Fatalf("Listen: %v", err)
	}
	tcp := g.Addr()
	if tcp == nil || tcp.(*net.TCPAddr).Port == 0 {
		t.Fatalf("Addr: хотел выданный порт, получил %v", tcp)
	}
	if err := g.SetForward(pc.LocalAddr().String()); err != nil {
		t.Fatalf("SetForward: %v", err)
	}
	logPath := filepath.Join(t.TempDir(), "archive.bin")
	if err := g.SetLog(logPath); err != nil {
		t.Fatalf("SetLog: %v", err)
	}

	// Serve в горутине с ctx.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	served := make(chan struct{})
	go func() {
		g.Serve(ctx)
		close(served)
	}()

	// Клиент: три строки.
	conn, err := net.DialTimeout("tcp", tcp.String(), 2*time.Second)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	lines := []string{"первая запись", "second line", "No: 3"}
	for _, l := range lines {
		fmt.Fprintf(conn, "%s\n", l)
	}
	conn.Close()

	// Три UDP-датаграммы (с дедлайном, чтобы не виснуть).
	buf := make([]byte, 1024)
	var got []string
	if err := pc.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	for i := 0; i < len(lines); i++ {
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			t.Fatalf("ReadFrom (датаграмма %d): %v", i+1, err)
		}
		got = append(got, string(buf[:n]))
	}
	if strings.Join(got, "|") != strings.Join(lines, "|") {
		t.Errorf("датаграммы:\nполучил %q\nхотел  %q", got, lines)
	}

	// Лог: строка с \n, в том же порядке (лог пишется до отправки UDP).
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile(лог): %v", err)
	}
	if want := strings.Join(lines, "\n") + "\n"; string(data) != want {
		t.Errorf("лог:\nполучил %q\nхотел  %q", data, want)
	}

	// Отмена ctx завершает Serve.
	cancel()
	select {
	case <-served:
	case <-time.After(2 * time.Second):
		t.Fatal("Serve не завершился после отмены ctx")
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

// Эталонное решение главы 25.
package sockgw

import (
	"bufio"
	"context"
	"net"
	"os"
	"strings"
	"sync"
)

// Gateway — шлюз stream-socket → UDP.
type Gateway struct {
	mu  sync.Mutex
	ln  net.Listener
	fwd *net.UDPAddr
	log *os.File
	udp *net.UDPConn
}

func (g *Gateway) Listen(addr string) (err error) {
	if addr == "" {
		addr = "127.0.0.1:0"
	}
	g.ln, err = net.Listen("tcp", addr)
	return err
}

func (g *Gateway) Addr() net.Addr {
	if g.ln == nil {
		return nil
	}
	return g.ln.Addr()
}

func (g *Gateway) SetForward(addr string) (err error) {
	g.fwd, err = net.ResolveUDPAddr("udp", addr)
	return err
}

func (g *Gateway) SetLog(path string) (err error) {
	g.log, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	return err
}

// Serve блокируется до отмены ctx; accept и соединения — в своих горутинах.
func (g *Gateway) Serve(ctx context.Context) {
	g.mu.Lock()
	if g.fwd != nil && g.udp == nil {
		if udp, err := net.DialUDP("udp", nil, g.fwd); err == nil {
			g.udp = udp
		}
	}
	ln := g.ln
	g.mu.Unlock()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go g.handleConn(conn)
		}
	}()
	<-ctx.Done()
	g.Close()
}

func (g *Gateway) handleConn(conn net.Conn) {
	defer conn.Close()
	for r := bufio.NewReader(conn); ; {
		line, err := r.ReadString('\n')
		if len(line) > 0 {
			g.record(strings.TrimRight(line, "\r\n"))
		}
		if err != nil {
			return
		}
	}
}

// record: лог пишется до отправки UDP; под мьютексом — Close не мешает.
func (g *Gateway) record(rec string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.log != nil {
		g.log.WriteString(rec + "\n")
	}
	if g.udp != nil {
		g.udp.Write([]byte(rec))
	}
}

func (g *Gateway) Close() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.ln != nil {
		g.ln.Close()
	}
	if g.udp != nil {
		g.udp.Close()
	}
	if g.log != nil {
		g.log.Close()
	}
	g.ln, g.udp, g.log = nil, nil, nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/25-sockgw`, упражнение **sockgw**. Локально: `go test -tags solution ./25-sockgw/`._
