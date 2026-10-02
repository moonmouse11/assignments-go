# Глава 26. Dsp — цифровой звук на stdlib

> По мотивам `dsp.pdf` (C++-компаньон ThinkDSP), урезан до ядра:
> синтез, WAV-файлы, спектральный анализ — без единой внешней библиотеки.

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
