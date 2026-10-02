//go:build solution

// Эталонное решение главы 26. Сравни со своей реализацией в dsp.go.
package dsp

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func Synth(freq float64, secs int, rate int) []float64 {
	n := secs * rate
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = 0.5 * math.Sin(2*math.Pi*freq*float64(i)/float64(rate))
	}
	return out
}

func quantize(s float64) int16 {
	v := math.Round(s * 32767)
	if v > 32767 {
		v = 32767
	}
	if v < -32768 {
		v = -32768
	}
	return int16(v)
}

func WriteWav(w io.Writer, samples []float64, rate int) error {
	data := make([]byte, len(samples)*2)
	for i, s := range samples {
		binary.LittleEndian.PutUint16(data[2*i:], uint16(quantize(s)))
	}
	var h [44]byte
	copy(h[0:4], "RIFF")
	binary.LittleEndian.PutUint32(h[4:8], uint32(36+len(data)))
	copy(h[8:12], "WAVE")
	copy(h[12:16], "fmt ")
	binary.LittleEndian.PutUint32(h[16:20], 16) // размер fmt-чанка
	binary.LittleEndian.PutUint16(h[20:22], 1)  // PCM
	binary.LittleEndian.PutUint16(h[22:24], 1)  // моно
	binary.LittleEndian.PutUint32(h[24:28], uint32(rate))
	binary.LittleEndian.PutUint32(h[28:32], uint32(rate*2)) // байт/сек
	binary.LittleEndian.PutUint16(h[32:34], 2)              // блок = 2 байта
	binary.LittleEndian.PutUint16(h[34:36], 16)             // бит на сэмпл
	copy(h[36:40], "data")
	binary.LittleEndian.PutUint32(h[40:44], uint32(len(data)))
	if _, err := w.Write(h[:]); err != nil {
		return err
	}
	_, err := w.Write(data)
	return err
}

func ReadWav(r io.Reader) ([]float64, int, error) {
	var h [44]byte
	if _, err := io.ReadFull(r, h[:]); err != nil {
		return nil, 0, fmt.Errorf("ReadWav: короткий заголовок: %w", err)
	}
	if string(h[0:4]) != "RIFF" || string(h[8:12]) != "WAVE" {
		return nil, 0, fmt.Errorf("ReadWav: не RIFF/WAVE файл")
	}
	if string(h[12:16]) != "fmt " || string(h[36:40]) != "data" {
		return nil, 0, fmt.Errorf("ReadWav: нет чанков fmt/data")
	}
	rate := int(binary.LittleEndian.Uint32(h[24:28]))
	size := int(binary.LittleEndian.Uint32(h[40:44]))
	data := make([]byte, size)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, 0, fmt.Errorf("ReadWav: короткие данные: %w", err)
	}
	samples := make([]float64, size/2)
	for i := range samples {
		samples[i] = float64(int16(binary.LittleEndian.Uint16(data[2*i:]))) / 32767
	}
	return samples, rate, nil
}

func DominantFreq(samples []float64, rate int) float64 {
	best, bestMag := 0.0, -1.0
	for f := 20; f <= 4000; f++ {
		w := 2 * math.Pi * float64(f) / float64(rate)
		cw, sw := math.Cos(w), math.Sin(w)
		pr, pi := 1.0, 0.0 // e^{-iωt} при t=0
		var re, im float64
		for _, x := range samples {
			re += x * pr
			im += x * pi
			pr, pi = pr*cw-pi*sw, pr*sw+pi*cw // e^{-iω(t+1)}
		}
		if mag := math.Hypot(re, im); mag > bestMag {
			bestMag, best = mag, float64(f)
		}
	}
	return best
}
