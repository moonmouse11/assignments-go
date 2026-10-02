package dsp

import (
	"bytes"
	"math"
	"testing"
)

func TestSynth(t *testing.T) {
	s := Synth(440, 1, 8000)
	if len(s) != 8000 {
		t.Fatalf("Synth: длина %d, хотел 8000 (1 сек × 8000 Гц)", len(s))
	}
	for i, x := range s {
		if math.Abs(x) > 0.5 {
			t.Fatalf("Synth: [%d] = %v, амплитуда вне [-0.5, 0.5]", i, x)
		}
	}
	if s[0] != 0 {
		t.Errorf("Synth: x(0) = %v, хотел 0", s[0])
	}
}

func TestDominantFreq(t *testing.T) {
	s := Synth(440, 1, 8000)
	if got := DominantFreq(s, 8000); math.Abs(got-440) > 2 {
		t.Errorf("DominantFreq(440 Гц) = %v, хотел ≈440 (±2)", got)
	}
	s = Synth(220, 1, 8000)
	if got := DominantFreq(s, 8000); math.Abs(got-220) > 2 {
		t.Errorf("DominantFreq(220 Гц) = %v, хотел ≈220 (±2)", got)
	}
}

func TestWavRoundTrip(t *testing.T) {
	src := Synth(330, 1, 8000)
	var buf bytes.Buffer
	if err := WriteWav(&buf, src, 8000); err != nil {
		t.Fatalf("WriteWav: %v", err)
	}
	if buf.Len() != 44+2*len(src) {
		t.Fatalf("WriteWav: размер %d, хотел 44+2×%d", buf.Len(), len(src))
	}
	got, rate, err := ReadWav(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("ReadWav: %v", err)
	}
	if rate != 8000 {
		t.Errorf("ReadWav: rate = %d, хотел 8000", rate)
	}
	if len(got) != len(src) {
		t.Fatalf("ReadWav: сэмплов %d, хотел %d", len(got), len(src))
	}
	for i := range got {
		if math.Abs(got[i]-src[i]) > 1e-3 {
			t.Fatalf("round-trip: [%d] = %v, хотел ≈%v (порог 1e-3)", i, got[i], src[i])
		}
	}
}

func TestWavClipping(t *testing.T) {
	src := []float64{2.5, -2.5}
	var buf bytes.Buffer
	if err := WriteWav(&buf, src, 8000); err != nil {
		t.Fatalf("WriteWav: %v", err)
	}
	got, _, err := ReadWav(&buf)
	if err != nil {
		t.Fatalf("ReadWav: %v", err)
	}
	// Порог квантования [-32768, 32767]: сверху зажимается в 32767/32767 = 1,
	// снизу — в -32768/32767.
	if got[0] != 1.0 {
		t.Errorf("клиппинг сверху: %v, хотел 1.0", got[0])
	}
	if want := -32768.0 / 32767.0; got[1] != want {
		t.Errorf("клиппинг снизу: %v, хотел %v", got[1], want)
	}
}

func TestReadWavBadHeader(t *testing.T) {
	if _, _, err := ReadWav(bytes.NewReader([]byte("NOPE NOPE NOPE NOPE NOPE"))); err == nil {
		t.Error("ReadWav: хотел ошибку для не-RIFF файла")
	}
	if _, _, err := ReadWav(bytes.NewReader([]byte("RIFFxxxxWAVEjunk"))); err == nil {
		t.Error("ReadWav: хотел ошибку для кривого заголовка")
	}
	if _, _, err := ReadWav(bytes.NewReader(nil)); err == nil {
		t.Error("ReadWav: хотел ошибку для пустого файла")
	}
}
