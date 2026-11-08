# Глава 11. IPAdr — анализ IPv4/IPv6

> Разбор адреса: класс, приватность,
> зарезервированность, сеть/широковещание/маска — чистый интерфейс.

## Цель главы

По строке адреса выдать полную справку: класс, приватность,
зарезервированность (loopback, link-local), сеть, широковещание,
маску. Заодно отличить IPv6 от IPv4.

## Теория

- Классы по первому октету: A — 0–127, B — 128–191, C — 192–223,
  D (multicast) — 224–239, E (экспериментальный) — 240–255.
- Приватные (RFC 1918): `10/8`, `172.16/12`, `192.168/16`.
- Зарезервированные: `127/8` loopback, `169.254/16` link-local,
  `0.0.0.0/8` this network.
- Сеть = ip & mask, broadcast = сеть | ^mask; октеты → `uint32`
  сдвигами и обратно.
- `net.ParseIP` отличит IPv6; `To4() == nil` значит «это не IPv4».

## Условие

1. `ParseV4(s)` — строгий dotted-quad: 4 октета 0–255, только цифры;
   кривой адрес — ошибка.
2. `AnalyzeV4(ip, mask)` — справка; пустая маска → классовая
   (D/E → /32); формы `"255.255.255.192"`, `"/26"`, `"26"`;
   Network/Broadcast/Mask заполняются всегда.
3. `IsV6(s)` — true только для IPv6: `"2001:db8::1"` → true,
   `"1.2.3.4"` → false.
4. Кривая маска (несмежные биты, октет > 255, префикс > 32) — ошибка.

Wants (необязательно, из `ipadr.pdf`): подкоманды v4/v6, node id.

## API для реализации

```go
type V4Info struct {
	Class     string // "A".."E"
	Private   bool
	Reserved  string // "", "loopback", "link-local", "this network"
	Network   string
	Broadcast string
	Mask      string
}
func ParseV4(s string) ([4]byte, error)
func AnalyzeV4(ip string, mask string) (V4Info, error)
func IsV6(s string) bool
```

Сигнатуры и договоры — в `ipadr.go`. Тесты — в `ipadr_test.go`.

## Подсказки

- Собери `uint32` из октетов: `uint32(o)<<24 | ...` и обратно — сдвигами.
- Смежность маски: `inv := ^b; inv&(inv+1) == 0` — единицы слева
  без «дыр».
- Ошибки — `fmt.Errorf("ipadr: ...")`, в тестах смотри только err != nil.

## Проверка

```
go test ./11-ipadr/
go test -tags solution ./11-ipadr/
```

## Эталон

`go test -tags solution ./11-ipadr/`.

## Углубление

- Подкоманды `v4`/`v6` из спеки: `flag.NewFlagSet` — остальная часть
  командной строки обрабатывается в контексте v4/v6; спека
  подчёркивает выгоду чистого дизайна интерфейса (полиморфизм).
- Node id: если подсеть выводима — не только network/broadcast, но и
  номер узла: `ip & ^mask` (uint32-арифметика).
- `-m/--mask` без значения — аргумент сам трактуется как маска
  (syntax + semantic check); со значением — пара анализируется вместе.
- Заметки спеки: `.0` — network id, `.255` — broadcast — «плохие»
  адреса, заслуживают предупреждения.
- v6-контекст: `net.ParseIP` / `net.ParseCIDR`.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./11-ipadr/
```

## Заготовки (реализуйте)

### `ipadr.go`

{% raw %}
```go
//go:build !solution

// Пакет ipadr — глава 11. Реализуй функции так, чтобы прошли тесты
// в ipadr_test.go. Условие — в README.md этой папки.
package ipadr

// V4Info — справка по IPv4-адресу. Network, Broadcast и Mask
// заполняются всегда: по явной маске или классовой (A → /8,
// B → /16, C → /24, D/E → /32).
type V4Info struct {
	Class     string // "A".."E" по первому октету
	Private   bool   // RFC 1918: 10/8, 172.16/12, 192.168/16
	Reserved  string // "", "loopback", "link-local", "this network"
	Network   string // dotted-quad
	Broadcast string // dotted-quad
	Mask      string // dotted-quad
}

// ParseV4 разбирает dotted-quad "a.b.c.d" в 4 октета. Каждый октет —
// 1-3 цифры без знака, 0-255, частей ровно 4. Иначе — ошибка.
func ParseV4(s string) ([4]byte, error) {
	panic("TODO: ParseV4")
}

// AnalyzeV4 возвращает справку по адресу ip. Пустая mask — классовая
// по первому октету (D/E → /32). Маска бывает dotted-quad
// ("255.255.255.192"), префиксом ("/26" или "26"); несмежные биты
// или префикс > 32 — ошибка, как и кривой адрес.
func AnalyzeV4(ip string, mask string) (V4Info, error) {
	panic("TODO: AnalyzeV4")
}

// IsV6 сообщает, что s — IPv6-адрес ("2001:db8::1" — true).
// IPv4 и мусор — false.
func IsV6(s string) bool {
	panic("TODO: IsV6")
}
```
{% endraw %}

## Тесты

### `ipadr_test.go`

{% raw %}
```go
package ipadr

import "testing"

func TestParseV4(t *testing.T) {
	cases := map[string][4]byte{
		"8.8.8.8":         {8, 8, 8, 8},
		"192.168.1.130":   {192, 168, 1, 130},
		"0.0.0.0":         {0, 0, 0, 0},
		"255.255.255.192": {255, 255, 255, 192},
	}
	for s, want := range cases {
		got, err := ParseV4(s)
		if err != nil {
			t.Errorf("ParseV4(%q): неожиданная ошибка %v", s, err)
			continue
		}
		if got != want {
			t.Errorf("ParseV4(%q): получил %v — хотел %v", s, got, want)
		}
	}
}

func TestParseV4Errors(t *testing.T) {
	for _, s := range []string{"256.1.1.1", "1.2.3", "1.2.3.4.5", "a.b.c.d", "", "1.2.3.", ".1.2.3", " 1.2.3.4", "1.2.3.4 ", "01.2.3.4x", "+1.2.3.4", "1.2.3.-4"} {
		if _, err := ParseV4(s); err == nil {
			t.Errorf("ParseV4(%q): хотел ошибку — получил nil", s)
		}
	}
}

func TestAnalyzeV4(t *testing.T) {
	// 8.8.8.8 — класс A, публичный, ничего не зарезервировано.
	info, err := AnalyzeV4("8.8.8.8", "")
	if err != nil {
		t.Fatalf("AnalyzeV4(8.8.8.8): %v", err)
	}
	if info.Class != "A" || info.Private || info.Reserved != "" {
		t.Errorf("8.8.8.8: получил %+v — хотел класс A, публичный, без Reserved", info)
	}

	// 10.1.2.3 — приватный (10/8).
	info, err = AnalyzeV4("10.1.2.3", "")
	if err != nil || !info.Private {
		t.Errorf("10.1.2.3: получил %+v, err %v — хотел Private=true", info, err)
	}

	// 192.168.1.130/26 — сеть 192.168.1.128, broadcast .191, маска .192.
	for _, mask := range []string{"/26", "26", "255.255.255.192"} {
		info, err = AnalyzeV4("192.168.1.130", mask)
		if err != nil {
			t.Fatalf("AnalyzeV4(192.168.1.130, %q): %v", mask, err)
		}
		want := V4Info{Class: "C", Private: true, Network: "192.168.1.128", Broadcast: "192.168.1.191", Mask: "255.255.255.192"}
		if info != want {
			t.Errorf("192.168.1.130 %q: получил %+v — хотел %+v", mask, info, want)
		}
	}

	// 172.20.0.1/12 — приватный, сеть 172.16.0.0, broadcast 172.31.255.255.
	info, err = AnalyzeV4("172.20.0.1", "/12")
	if err != nil {
		t.Fatalf("AnalyzeV4(172.20.0.1, /12): %v", err)
	}
	want := V4Info{Class: "B", Private: true, Network: "172.16.0.0", Broadcast: "172.31.255.255", Mask: "255.240.0.0"}
	if info != want {
		t.Errorf("172.20.0.1/12: получил %+v — хотел %+v", info, want)
	}

	// 127.0.0.1 — loopback.
	info, err = AnalyzeV4("127.0.0.1", "")
	if err != nil || info.Reserved != "loopback" {
		t.Errorf("127.0.0.1: получил %+v, err %v — хотел Reserved=loopback", info, err)
	}

	// Классы D и E, link-local.
	if info, _ := AnalyzeV4("224.0.0.9", ""); info.Class != "D" {
		t.Errorf("224.0.0.9: класс %q — хотел D", info.Class)
	}
	if info, _ := AnalyzeV4("250.1.1.1", ""); info.Class != "E" {
		t.Errorf("250.1.1.1: класс %q — хотел E", info.Class)
	}
	if info, _ := AnalyzeV4("169.254.1.1", ""); info.Reserved != "link-local" {
		t.Errorf("169.254.1.1: Reserved %q — хотел link-local", info.Reserved)
	}
	if info, _ := AnalyzeV4("0.0.0.0", ""); info.Reserved != "this network" {
		t.Errorf("0.0.0.0: Reserved %q — хотел this network", info.Reserved)
	}
}

func TestAnalyzeV4Errors(t *testing.T) {
	for _, c := range [][2]string{
		{"256.1.1.1", ""}, {"1.2.3", ""}, {"abc", "/8"},
		{"1.2.3.4", "255.0.0.256"}, {"1.2.3.4", "/33"}, {"1.2.3.4", "abc"},
		{"1.2.3.4", "255.0.255.0"}, {"1.2.3.4", "-1"}, {"1.2.3.4", "/"},
	} {
		if _, err := AnalyzeV4(c[0], c[1]); err == nil {
			t.Errorf("AnalyzeV4(%q, %q): хотел ошибку — получил nil", c[0], c[1])
		}
	}
}

func TestIsV6(t *testing.T) {
	cases := map[string]bool{
		"2001:db8::1": true,
		"::1":         true,
		"fe80::1":     true,
		"1.2.3.4":     false,
		"301.2.3.4":   false,
		"hello":       false,
		"":            false,
	}
	for s, want := range cases {
		if got := IsV6(s); got != want {
			t.Errorf("IsV6(%q): получил %v — хотел %v", s, got, want)
		}
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

// Эталонное решение главы 11. Сравни со своей реализацией в ipadr.go.
package ipadr

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type V4Info struct {
	Class     string // "A".."E" по первому октету
	Private   bool
	Reserved  string // "", "loopback", "link-local", "this network"
	Network   string
	Broadcast string
	Mask      string
}

func ParseV4(s string) ([4]byte, error) {
	var ip [4]byte
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		return ip, fmt.Errorf("ipadr: %q", s)
	}
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n > 255 || len(p) > 3 || strings.ContainsAny(p, "+-") {
			return ip, fmt.Errorf("ipadr: %q", s)
		}
		ip[i] = byte(n)
	}
	return ip, nil
}

func maskBits(mask string, class byte) (uint32, error) {
	if mask == "" { // пустая → классовая: A /8, B /16, C /24, D/E /32
		return ^(^uint32(0) >> [8]int{8, 8, 8, 8, 16, 16, 24, 32}[class>>5]), nil
	}
	m := strings.TrimPrefix(mask, "/")
	if !strings.Contains(m, ".") { // префикс "/26" или "26"
		n, err := strconv.Atoi(m)
		if err != nil || n < 0 || n > 32 {
			return 0, fmt.Errorf("ipadr: префикс %q", mask)
		}
		return ^(^uint32(0) >> n), nil
	}
	o, err := ParseV4(m)
	if err != nil {
		return 0, fmt.Errorf("ipadr: маска %q", mask)
	}
	b := uint32(o[0])<<24 | uint32(o[1])<<16 | uint32(o[2])<<8 | uint32(o[3])
	if inv := ^b; inv&(inv+1) != 0 { // единицы — слева без дыр
		return 0, fmt.Errorf("ipadr: несмежная маска %q", mask)
	}
	return b, nil
}

var classes = [16]string{"A", "A", "A", "A", "A", "A", "A", "A", "B", "B", "B", "B", "C", "C", "D", "E"}

func AnalyzeV4(ip string, mask string) (V4Info, error) {
	o, err := ParseV4(ip)
	if err != nil {
		return V4Info{}, err
	}
	m, err := maskBits(mask, o[0])
	if err != nil {
		return V4Info{}, err
	}
	nw := (uint32(o[0])<<24 | uint32(o[1])<<16 | uint32(o[2])<<8 | uint32(o[3])) & m
	res := map[byte]string{127: "loopback", 0: "this network"}[o[0]]
	if o[0] == 169 && o[1] == 254 {
		res = "link-local"
	}
	return V4Info{
		Class:     classes[o[0]>>4],
		Private:   o[0] == 10 || o[0] == 172 && o[1]&0xF0 == 16 || o[0] == 192 && o[1] == 168,
		Reserved:  res,
		Network:   fmt.Sprintf("%d.%d.%d.%d", nw>>24, nw>>16&0xFF, nw>>8&0xFF, nw&0xFF),
		Broadcast: fmt.Sprintf("%d.%d.%d.%d", (nw|^m)>>24, (nw|^m)>>16&0xFF, (nw|^m)>>8&0xFF, (nw|^m)&0xFF),
		Mask:      fmt.Sprintf("%d.%d.%d.%d", m>>24, m>>16&0xFF, m>>8&0xFF, m&0xFF),
	}, nil
}

func IsV6(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && ip.To4() == nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/11-ipadr`, упражнение **ipadr**. Локально: `go test -tags solution ./11-ipadr/`._
