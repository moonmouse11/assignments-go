# Глава 26. Dsp — цифровой звук на stdlib

> Урезан до ядра: синтез, WAV-файлы, спектральный анализ —
> без единой внешней библиотеки.

## Цель главы

Сгенерировать синусоиду, записать и прочитать WAV-файл PCM16, найти
доминирующую частоту наивным DFT. Только `math` и `encoding/binary`.

## Теория

- Синус: `x(t) = 0.5·sin(2π·f·t/rate)`; сэмплов — `secs·rate`.
- WAV (RIFF): 44-байтовый заголовок + PCM16-моно данные; все числа —
  little-endian (`encoding/binary`).
- Квантование: float → int16 (умножить на 32767, округлить, зажать в
  [-32768, 32767]) и обратно.
- DFT: `X(f) = Σ x(t)·e^{-iωt}`; перебираем f от 20 до 4000 Гц с шагом 1,
  берём максимум модуля. Экспонента — cos и sin.
- 3981 частота × 8000 сэмплов: `math.Cos` на каждый шаг дорог; вращай
  фазу рекуррентно (умножение комплексных чисел) — на порядок быстрее.

## Условие

Мини-DSP без зависимостей.

1. `Synth(freq, secs, rate)` — синусоида амплитуды 0.5.
2. `WriteWav`/`ReadWav` — честный RIFF/WAVE PCM16 моно round-trip:
   порог квантования [-32768, 32767], сверху клиппинг.
3. `DominantFreq` — наивный DFT: частоты 20..4000 Гц с шагом 1,
   максимум |Σ x(t)·e^{-iωt}|.

## API для реализации

```go
func Synth(freq float64, secs int, rate int) []float64
func WriteWav(w io.Writer, samples []float64, rate int) error
func ReadWav(r io.Reader) (samples []float64, rate int, err error)
func DominantFreq(samples []float64, rate int) float64
```

Сигнатуры и договоры — в `dsp.go`. Тесты — в `dsp_test.go`.

## Подсказки

- Заголовок собери в `[44]byte`: `copy(h[0:4], "RIFF")` и
  `binary.LittleEndian.PutUint32(h[4:8], 36+len(data))`.
- `io.ReadFull` при чтении: короткий или кривой файл — ошибка.
- Проверь байты `"RIFF"`/`"WAVE"`/`"fmt "`/`"data"` до разбора полей.

## Проверка

```
go test ./26-dsp/
go test -tags solution ./26-dsp/
```

## Эталон

`go test -tags solution ./26-dsp/`

## Углубление

Полная спека — C++-компаньон ThinkDSP на GSL + libsndfile + gnuplot/octave;
в Go почти всё ядро доступно на stdlib.

- Зоопарк сигналов: периодические с амплитудой/частотой/фазой (сдвиг фазы
  в `Synth`), апериодические `GUUNoise`, `GBrownian` (накопление шума),
  `GGaussian` — `math`, `math/rand/v2`.
- `GTestSignals` (`GStep` и др.): сигналы-инструменты для вырезания
  сегментов/окон над другой волной — срезы и умножение, чистый Go.
- `GSpectrum` как полноценный тип: двусторонний обмен волна↔спектр
  (синтез из частотной области), а не только поиск максимума; сюда же
  цифровые фильтры из спеки.
- Внешние зависимости спеки: GSL (заменимо на `math`), gnuplot/octave для
  графиков, libsndfile для аудиоформатов кроме WAV — в главе обошлись
  `encoding/binary` и своим RIFF.

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./26-dsp/
```

## Заготовки (реализуйте)

### `dsp.go`

{% raw %}
```go
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
```
{% endraw %}

## Тесты

### `dsp_test.go`

{% raw %}
```go
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
```
{% endraw %}

## Эталонное решение

Сначала попробуйте сами. Решение скрыто — раскрывайте только для сверки.

<details>
<summary><strong>Эталон (`solution.go`)</strong> — нажмите, чтобы показать</summary>

{% raw %}
```go
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
```
{% endraw %}


</details>

_Пакет `assignments-go/26-dsp`, упражнение **dsp**. Локально: `go test -tags solution ./26-dsp/`._
