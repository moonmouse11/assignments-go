//go:build solution

// Эталонное решение главы 13. Сравни со своей реализацией в crc.go.
package crc

// CRC32 — битовая реализация CRC-32/IEEE:
// полином 0xEDB88320 (reflected), init 0xFFFFFFFF, xorout 0xFFFFFFFF.
func CRC32(data []byte) uint32 {
	const poly = 0xEDB88320
	crc := uint32(0xFFFFFFFF)
	for _, b := range data {
		crc ^= uint32(b)
		for range 8 {
			if crc&1 == 1 {
				crc = crc>>1 ^ poly
			} else {
				crc >>= 1
			}
		}
	}
	return ^crc // xorout 0xFFFFFFFF
}
