# Глава 12. NetGeo — география по IP, haversine

> Сервисы геолокации (ipstack) отвечают JSON;
> координаты двух точек дают оценку расстояния по формуле haversine.
> Сеть в курсе запрещена — работаем с фиксированными ответами.

## Цель главы

Разобрать JSON-ответ геолокации в `Location` и посчитать расстояние
между локациями по сфере: Москва — Питер ≈ 630 км. В оригинале —
api.ipstack.com и traceroute; здесь — те же вычисления без сети.

## Теория

- `encoding/json`: теги полей `json:"latitude"` маппят ответ сервиса
  на структуру; `json.Unmarshal(data, &loc)` — разбор, ошибка формата
  наружу.
- Формула haversine (в радианах!): для точек (φ, λ) и (φ₂, λ₂)
  a = sin²(Δφ/2) + cosφ·cosφ₂·sin²(Δλ/2);
  d = R · 2·atan2(√a, √(1−a)), R = 6371 км.
- `math.Sin/Cos/Pow/Sqrt/Atan2` — всё из stdlib; градусы в радианы
  через `x * math.Pi / 180`.
- Точность — оценка: высоты у нас нет, Земля не сфера — потому
  630 км и допуск ±30 км в тестах.

## Условие

1. `ParseLocation(data)` — JSON ipstack → `Location`: поля ip,
   latitude, longitude, city. Кривой JSON — ошибка.
2. `HaversineKm(a, b)` — расстояние в километмах между локациями
   по формуле haversine, R = 6371 км.
3. Одна и та же точка — расстояние 0.

Wants (необязательно, из `netgeo.pdf`): список hops из traceroute
и расстояния между соседними узлами, шифрование API-ключа.

## API для реализации

```go
type Location struct {
	IP        string
	Latitude  float64
	Longitude float64
	City      string
}
func ParseLocation(data []byte) (Location, error)
func HaversineKm(a, b Location) float64
```

Сигнатуры и договоры — в `netgeo.go`. Тесты — в `netgeo_test.go`.

## Подсказки

- Поле `City` в ответе может отсутствовать (null) — разбирайся без
  паники; тег `json:"city"` оставит пустую строку.
- Проверяй «одинаковая точка → 0» через `math.Abs(d) < 1e-9` —
  строгое равенство float — плохой тон.
- Москва (55.75, 37.62) — Питер (59.94, 30.31): получится ~635 км,
  внутри допуска 630±30.

## Проверка

```
go test ./12-netgeo/
go test -tags solution ./12-netgeo/
```

## Эталон

`go test -tags solution ./12-netgeo/`.

## Углубление

- `-t/--trace`: traceroute к аргументу, расстояние каждого hop;
  запуск утилиты программно — `os/exec`, разбор вывода — `regexp`;
  нужен установленный traceroute (внешняя зависимость).
- Живой api.ipstack.com: `net/http` + `encoding/json`, API-ключ по
  регистрации — требует сети; фикстуры тестов → реальные ответы.
- Скрытие API-ключа: у автора xor-шифрование в бинарнике (bindata —
  внешняя зависимость); в stdlib — `crypto/aes`, но ключ в бинарнике
  публичного репозитория не защитит (спека честно называет это
  иллюстрацией идеи симметричного шифрования).
- Enhancement спеки: другие облачные сервисы координат городов и
  расстояния между городами — сеть (`net/http`).
- CLI: `flag`, несколько `ipspec` — расстояния всех пар аргументов.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./12-netgeo/
```

## Заготовки (реализуйте)

### `netgeo.go`

{% raw %}
```go
//go:build !solution

// Пакет netgeo — глава 12. Реализуй функции так, чтобы прошли тесты
// в netgeo_test.go. Условие — в README.md этой папки.
package netgeo

// Location — точка на карте: IP-адрес узла, широта, долгота, город.
// Разбирается из JSON-ответа сервиса геолокации в духе ipstack.
type Location struct {
	IP        string
	Latitude  float64
	Longitude float64
	City      string
}

// ParseLocation разбирает JSON вида ipstack:
//
//	{"ip":"134.201.250.9","latitude":34.06,"longitude":-118.08,"city":"Los Angeles"}
//
// Поля city может не быть — тогда пустая строка. Кривой JSON — ошибка.
func ParseLocation(data []byte) (Location, error) {
	panic("TODO: ParseLocation")
}

// HaversineKm возвращает расстояние между a и b в километрах
// по формуле haversine (R = 6371 км, координаты в градусах).
// Одна и та же точка — 0.
func HaversineKm(a, b Location) float64 {
	panic("TODO: HaversineKm")
}
```
{% endraw %}

## Тесты

### `netgeo_test.go`

{% raw %}
```go
package netgeo

import (
	"math"
	"testing"
)

const ipstack = `{"ip":"134.201.250.9","type":"ipv4","continent_code":"NA",
"latitude":34.06,"longitude":-118.08,"city":"Los Angeles",
"region_code":"CA","country_code":"US"}`

func TestParseLocation(t *testing.T) {
	loc, err := ParseLocation([]byte(ipstack))
	if err != nil {
		t.Fatalf("ParseLocation: %v", err)
	}
	want := Location{IP: "134.201.250.9", Latitude: 34.06, Longitude: -118.08, City: "Los Angeles"}
	if loc != want {
		t.Errorf("ParseLocation: получил %+v — хотел %+v", loc, want)
	}
}

func TestParseLocationNoCity(t *testing.T) {
	// city может отсутствовать — тогда пустая строка, без ошибки.
	loc, err := ParseLocation([]byte(`{"ip":"8.8.8.8","latitude":37.4,"longitude":-122.1,"city":null}`))
	if err != nil {
		t.Fatalf("ParseLocation: %v", err)
	}
	if loc.City != "" || loc.IP != "8.8.8.8" {
		t.Errorf("ParseLocation без city: получил %+v", loc)
	}
}

func TestParseLocationErrors(t *testing.T) {
	for _, data := range [][]byte{
		[]byte(`{"ip":`),
		[]byte(`не JSON`),
		[]byte(``),
		[]byte(`{"latitude":"юг"}`), // широта обязана быть числом
	} {
		if _, err := ParseLocation(data); err == nil {
			t.Errorf("ParseLocation(%q): хотел ошибку — получил nil", data)
		}
	}
}

func TestHaversineMoscowSpb(t *testing.T) {
	msk := Location{IP: "1.1.1.1", Latitude: 55.75, Longitude: 37.62, City: "Moscow"}
	spb := Location{IP: "2.2.2.2", Latitude: 59.94, Longitude: 30.31, City: "Saint Petersburg"}
	d := HaversineKm(msk, spb)
	if math.Abs(d-630) > 30 {
		t.Errorf("Москва—Питер: получил %.1f км — хотел 630±30", d)
	}
	// Расстояние симметрично.
	if back := HaversineKm(spb, msk); math.Abs(d-back) > 1e-9 {
		t.Errorf("HaversineKm не симметричен: %.3f против %.3f", d, back)
	}
}

func TestHaversineSamePoint(t *testing.T) {
	p := Location{IP: "3.3.3.3", Latitude: 55.75, Longitude: 37.62, City: "Moscow"}
	if d := HaversineKm(p, p); math.Abs(d) > 1e-9 {
		t.Errorf("одна и та же точка: получил %v — хотел 0", d)
	}
}

func TestHaversineKnownDistance(t *testing.T) {
	// Лос-Анджелес — Нью-Йорк: около 3940 км по сфере.
	la := Location{Latitude: 34.06, Longitude: -118.08}
	ny := Location{Latitude: 40.71, Longitude: -74.01}
	d := HaversineKm(la, ny)
	if math.Abs(d-3940) > 100 {
		t.Errorf("Лос-Анджелес—Нью-Йорк: получил %.1f км — хотел 3940±100", d)
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

// Эталонное решение главы 12. Сравни со своей реализацией в netgeo.go.
package netgeo

import (
	"encoding/json"
	"fmt"
	"math"
)

// Location — точка на карте.
type Location struct {
	IP        string
	Latitude  float64
	Longitude float64
	City      string
}

func ParseLocation(data []byte) (Location, error) {
	var raw struct {
		IP        string  `json:"ip"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		City      string  `json:"city"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return Location{}, fmt.Errorf("netgeo: JSON: %v", err)
	}
	return Location{IP: raw.IP, Latitude: raw.Latitude, Longitude: raw.Longitude, City: raw.City}, nil
}

// HaversineKm — расстояние по сфере, R = 6371 км.
func HaversineKm(a, b Location) float64 {
	const r = 6371 // км
	lat1 := a.Latitude * math.Pi / 180
	lat2 := b.Latitude * math.Pi / 180
	dLat := (b.Latitude - a.Latitude) * math.Pi / 180
	dLon := (b.Longitude - a.Longitude) * math.Pi / 180
	x := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return r * 2 * math.Atan2(math.Sqrt(x), math.Sqrt(1-x))
}
```
{% endraw %}


</details>

_Пакет `assignments-go/12-netgeo`, упражнение **netgeo**. Локально: `go test -tags solution ./12-netgeo/`._
