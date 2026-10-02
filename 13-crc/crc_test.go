package crc

import (
	"hash/crc32"
	"testing"
)

func TestKnownVector(t *testing.T) {
	// Общепринятый тест-вектор CRC-32/IEEE (catalogue of CRCs).
	if got := CRC32([]byte("123456789")); got != 0xCBF43926 {
		t.Errorf("CRC32(\"123456789\"): получил %#08x — хотел 0xCBF43926", got)
	}
}

func TestEmpty(t *testing.T) {
	// CRC32 пустого сообщения == xorout-обнулённый init == 0x00000000.
	if got := CRC32(nil); got != 0 {
		t.Errorf("CRC32(nil): получил %#08x — хотел 0x00000000", got)
	}
	if got := CRC32([]byte{}); got != 0 {
		t.Errorf("CRC32(пусто): получил %#08x — хотел 0x00000000", got)
	}
}

func TestMatchesStdlib(t *testing.T) {
	// stdlib — измерительный прибор: совпадение доказывает схему.
	zero1024 := make([]byte, 1024)
	zero1024[512] = 1 // один-единственный бит в 1024 нулях
	messages := [][]byte{
		[]byte("a"),
		[]byte("abc"),
		[]byte("abcdefghijklmnopqrstuvwxyz"),
		[]byte("The quick brown fox jumps over the lazy dog"),
		{0x00, 0x01, 0x02, 0xFE, 0xFF},
		{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
		zero1024,
	}
	for _, m := range messages {
		if got, want := CRC32(m), crc32.ChecksumIEEE(m); got != want {
			t.Errorf("CRC32(%q): получил %#08x — хотел %#08x (stdlib)", m, got, want)
		}
	}
}

func TestSingleBitChange(t *testing.T) {
	// Один испорченный бит меняет контрольную сумму — потому CRC
	// и ловит ошибки линий (want из crc.pdf — шум).
	orig := []byte("perfect-29787-Dendrite-38203")
	c0 := CRC32(orig)
	for i := range orig {
		broken := append([]byte(nil), orig...)
		broken[i] ^= 0x01
		if c := CRC32(broken); c == c0 {
			t.Errorf("порча байта %d не меняет CRC — полином проморгал ошибку", i)
		}
	}
}
