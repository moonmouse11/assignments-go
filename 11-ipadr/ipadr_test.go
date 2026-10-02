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
