package dump

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatBlock(t *testing.T) {
	got := FormatBlock(0, []byte{'A', 'B', 0x00, 'C'})
	want := "0000: 41 42 00 43 |AB.C|\n"
	if got != want {
		t.Errorf("FormatBlock:\nполучил %q\nхотел  %q", got, want)
	}
	if got := FormatBlock(8, []byte{}); got != "0008: ||\n" {
		t.Errorf("FormatBlock пустой блок: получил %q", got)
	}
}

func TestDump(t *testing.T) {
	// 10 байт: два блока по 8, второй неполный.
	data := "ABCDEFGHIJ"
	var buf bytes.Buffer
	if err := Dump(strings.NewReader(data), 8, &buf); err != nil {
		t.Fatalf("Dump: %v", err)
	}
	want := "0000: 41 42 43 44 45 46 47 48 |ABCDEFGH|\n" +
		"0008: 49 4a |IJ|\n"
	if buf.String() != want {
		t.Errorf("Dump:\nполучил %q\nхотел  %q", buf.String(), want)
	}
}

func TestDumpNonPrintable(t *testing.T) {
	var buf bytes.Buffer
	if err := Dump(strings.NewReader("\x01\x7f\n"), 8, &buf); err != nil {
		t.Fatalf("Dump: %v", err)
	}
	if !strings.Contains(buf.String(), "|...|") {
		t.Errorf("непечатаемые байты должны быть '.': %q", buf.String())
	}
}

func TestDumpInvalidBlockSize(t *testing.T) {
	for _, size := range []int{0, 4, 12, 64, -8} {
		if err := Dump(strings.NewReader("data"), size, &bytes.Buffer{}); err == nil {
			t.Errorf("blockSize=%d: хотел ошибку, получил nil", size)
		}
	}
}
