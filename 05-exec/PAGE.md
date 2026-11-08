# Глава 5. Exec — планировщик команд (упрощённый)

> Урезан до ядра: **горутины, время, отмена**.
> Полная версия (подкоманды `status`/`stop`, `--clock`, логи в `~/logs`, IPC)
> — в плане расширения курса.

## Цель главы

Первая задача с параллелизмом: `context`, `select`, `time.Ticker`, отмена
работы горутины. Плюс запуск внешних процессов — `os/exec`.

## Теория

- `context.Context` — канал отмены: `ctx.Done()` закрывается при `cancel()`.
- `time.NewTicker(d)` шлёт в `t.C` тик каждые `d`; не забудь `defer t.Stop()`.
- `select` ждёт **первый** готовый канал — отмена проверяется наравне
  с тиками.
- `os/exec`: `osexec.Command(name, args...).Output()` — запуск и вывод.
  (`osexec` — алиас: пакет нашей задачи тоже называется `exec`, stdlib
  нельзя импортировать под тем же именем.)
- Тесты с реальным временем — с **запасом**: тики через 20 мс, отмена через
  70 мс, проверяем диапазон, а не точное число.

## Условие

Планировщик, исполняющий команду по расписанию.

1. Без расписания команда исполняется один раз — это `RunCommand`.
2. `--timer <интервал>`: команда исполняется **первый раз через интервал**,
   затем каждый интервал, пока не отменили — это `Periodic`.
3. Отмена — через `context`: все горутины завершаются, лишних исполнений
   после отмены нет.

`RunCommand(name, args...)` — выполнить команду, вернуть её stdout/stderr.
`Periodic(ctx, interval, fn)` — запускать `fn` по тикам таймера, вернуть
число исполнений.

## API для реализации

```go
func RunCommand(name string, args ...string) ([]byte, error)
func Periodic(ctx context.Context, interval time.Duration, fn func()) int
```

## Подсказки

- `Periodic` — один `for` с `select` на два канала: `ctx.Done()` и `t.C`.
- Первое исполнение — после первого тика, не сразу при старте.
- Число исполнений считай сам (локальная переменная — гонок нет).

## Проверка

```
go test ./05-exec/
go test -tags solution ./05-exec/
```

## Эталон

`go test -tags solution ./05-exec/`

## Углубление

- Полный CLI автора: подкоманды `run`/`status`/`stop`/`version`, флаги
  `--config` (`$HOME/.cli/cli.yaml`), `--verbose`.
- `--clock HH:MM` — исполнение в заданное время сегодня (или завтра,
  если уже прошло): `time.Until` до ближайшего следующего момента.
- Лог `~/logs/<script>.log` с шапкой `Command/Started` и подвалом
  `Terminated/Duration` — как в примере спеки.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./05-exec/
```

## Заготовки (реализуйте)

### `exec.go`

{% raw %}
```go
//go:build !solution

// Пакет exec — глава 5 (упрощённый projectlet из exec.pdf).
// Реализуй RunCommand и Periodic так, чтобы прошли тесты.
package exec

import (
	"context"
	"time"
)

// RunCommand выполняет команду name с аргументами args и возвращает
// её вывод (stdout/stderr) или ошибку запуска.
func RunCommand(name string, args ...string) ([]byte, error) {
	panic("TODO: RunCommand")
}

// Periodic запускает fn каждые interval, начиная через interval после
// старта, пока ctx не отменят. Возвращает число исполнений fn.
func Periodic(ctx context.Context, interval time.Duration, fn func()) int {
	panic("TODO: Periodic")
}
```
{% endraw %}

## Тесты

### `exec_test.go`

{% raw %}
```go
package exec

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestRunCommand(t *testing.T) {
	out, err := RunCommand("echo", "hello", "gobook")
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}
	if !bytes.Equal(out, []byte("hello gobook\n")) {
		t.Errorf("RunCommand: получил %q, хотел %q", out, "hello gobook\n")
	}
}

func TestRunCommandError(t *testing.T) {
	if _, err := RunCommand("no-such-command-xyz"); err == nil {
		t.Error("RunCommand: хотел ошибку для несуществующей команды")
	}
}

func TestPeriodicCounts(t *testing.T) {
	// Тики каждые 20 мс, отмена через 70 мс → исполнений 2..3.
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(70 * time.Millisecond)
		cancel()
	}()
	runs := Periodic(ctx, 20*time.Millisecond, func() {})
	if runs < 2 || runs > 3 {
		t.Errorf("Periodic: исполнений %d, хотел 2..3", runs)
	}
}

func TestPeriodicCancelImmediately(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	runs := Periodic(ctx, time.Hour, func() { t.Error("fn не должен выполниться после отмены") })
	if runs != 0 {
		t.Errorf("Periodic: исполнений %d, хотел 0", runs)
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

// Эталонное решение главы 5.
package exec

import (
	"context"
	"os/exec"
	"time"
)

func RunCommand(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

func Periodic(ctx context.Context, interval time.Duration, fn func()) int {
	t := time.NewTicker(interval)
	defer t.Stop()
	runs := 0
	for {
		select {
		case <-ctx.Done():
			return runs
		case <-t.C:
			fn()
			runs++
		}
	}
}
```
{% endraw %}


</details>

_Пакет `assignments-go/05-exec`, упражнение **exec**. Локально: `go test -tags solution ./05-exec/`._
