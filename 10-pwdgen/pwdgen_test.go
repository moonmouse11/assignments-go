package pwdgen

import (
	"regexp"
	"strings"
	"testing"
)

var words = []string{"vein", "Bone", "protein", "virus", "azalea", "Nettle"}

func TestDeterminism(t *testing.T) {
	g1 := NewGen(42, words, "-")
	g2 := NewGen(42, words, "-")
	if p1, p2 := g1.Password(3), g2.Password(3); p1 != p2 {
		t.Errorf("одинаковое семя:\nполучил %q и %q — хотел совпадение", p1, p2)
	}
	g3 := NewGen(43, words, "-")
	if p3, p4 := g3.Password(3), g1.Password(3); p3 == p4 {
		t.Errorf("разные семени:\nполучил одинаковые %q — хотел разные", p3)
	}
}

func TestFormat(t *testing.T) {
	g := NewGen(7, words, "-")
	p := g.Password(3)
	segs := strings.Split(p, "-")
	if len(segs) != 6 {
		t.Fatalf("Password(3): получил %q — хотел 6 частей (3 слова + 3 числа)", p)
	}
	wordSet := make(map[string]bool, len(words))
	for _, w := range words {
		wordSet[w] = true
	}
	numRe := regexp.MustCompile(`^\d{2,4}$`)
	for i, part := range segs {
		if i%2 == 0 {
			if !wordSet[part] {
				t.Errorf("часть %q не из словаря (пароль %q)", part, p)
			}
		} else if !numRe.MatchString(part) {
			t.Errorf("число %q не 2-4 цифры (пароль %q)", part, p)
		}
	}
}

func TestSeparator(t *testing.T) {
	g := NewGen(7, words, "~")
	p := g.Password(2)
	if len(strings.Split(p, "~")) != 4 {
		t.Errorf("Password(2) с sep ~: получил %q — хотел 4 части", p)
	}
}

func TestSmallCounts(t *testing.T) {
	g := NewGen(1, words, "-")
	if p := g.Password(0); p != "" {
		t.Errorf("Password(0): получил %q — хотел пустую строку", p)
	}
	if p := g.Password(-3); p != "" {
		t.Errorf("Password(-3): получил %q — хотел пустую строку", p)
	}
}

func TestEmptyWordsPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("NewGen с пустым words: хотел panic — не случился")
		}
	}()
	NewGen(1, nil, "-")
}
