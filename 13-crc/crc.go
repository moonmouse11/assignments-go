//go:build !solution

// Пакет crc — глава 13. Реализуй функцию так, чтобы прошли тесты
// в crc_test.go. Условие — в README.md этой папки.
package crc

// CRC32 возвращает контрольную сумму data по схеме CRC-32/IEEE:
// полином 0xEDB88320 (reflected), init 0xFFFFFFFF, xorout 0xFFFFFFFF.
// Только битовая реализация — hash/crc32 здесь не помощник.
// Пустой ввод — корректный случай: контрольная сумма пустого сообщения.
func CRC32(data []byte) uint32 {
	panic("TODO: CRC32")
}
