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
