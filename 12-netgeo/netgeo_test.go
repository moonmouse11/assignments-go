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
