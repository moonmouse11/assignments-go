//go:build !solution

// Пакет dsp — глава 26. Реализуй Synth, WriteWav, ReadWav и DominantFreq
// так, чтобы прошли тесты в dsp_test.go. Условие — в README.md этой папки.
package dsp

import "io"

// Synth возвращает синусоиду частоты freq Гц амплитуды 0.5 длиной
// secs*rate сэмплов: x(t) = 0.5·sin(2π·freq·t/rate).
func Synth(freq float64, secs int, rate int) []float64 {
	panic("TODO: Synth")
}

// WriteWav пишет samples в w как RIFF/WAVE PCM16 моно: 44-байтовый
// заголовок (little-endian) + данные. Значения квантуются умножением
// на 32767 с округлением и клиппингом в [-32768, 32767].
func WriteWav(w io.Writer, samples []float64, rate int) error {
	panic("TODO: WriteWav")
}

// ReadWav читает RIFF/WAVE PCM16 моно из r: проверяет заголовок,
// возвращает сэмплы (float, /32767) и sample rate.
// Кривой или короткий файл — ошибка.
func ReadWav(r io.Reader) (samples []float64, rate int, err error) {
	panic("TODO: ReadWav")
}

// DominantFreq ищет доминирующую частоту наивным DFT: перебирает
// частоты 20..4000 Гц с шагом 1, максимум |Σ x(t)·e^{-iωt}|.
// Возвращает частоту с наибольшей амплитудой.
func DominantFreq(samples []float64, rate int) float64 {
	panic("TODO: DominantFreq")
}
