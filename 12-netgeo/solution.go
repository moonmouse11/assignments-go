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
