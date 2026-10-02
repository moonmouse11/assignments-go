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
