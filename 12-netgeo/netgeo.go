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
