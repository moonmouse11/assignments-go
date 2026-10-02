//go:build !solution

// Пакет dump — глава 2. Реализуй функции так, чтобы прошли тесты.
package dump

import "io"

// Version — версия утилиты в форме name-major.minor (конвенция из спецификации).
const Version = "dump-0.1"

// FormatBlock формирует одну строку дампа блока data, начинающегося
// с байта offset:
//
//	"0000: 41 42 00 43 |AB.C|"
//
// Смещение — 4 hex-цифры и двоеточие, байты — по два hex-символа через
// пробел, затем " |", текстовая колонка (непечатаемые байты — '.'),
// "|" и "\n".
func FormatBlock(offset int, data []byte) string {
	panic("TODO: FormatBlock")
}

// Dump пишет в out дамп всех данных из r блоками по blockSize байт
// (строки в формате FormatBlock, смещение — номер байта начала блока).
// Допустимые blockSize: 8, 16, 32; иначе — ошибка.
func Dump(r io.Reader, blockSize int, out io.Writer) error {
	panic("TODO: Dump")
}
