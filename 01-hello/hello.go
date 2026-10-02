//go:build !solution

// Пакет hello — глава 1. Реализуй функции так, чтобы прошли тесты
// в hello_test.go. Условие — в README.md этой папки.
package hello

import "io"

// NumberLines возвращает строки lines, пронумерованные с 1
// в формате NNNN: "0007 текст". Строки разделяются \n,
// в конце — завершающий \n.
//
//	NumberLines([]string{"a", "b"}) == "0001 a\n0002 b\n"
func NumberLines(lines []string) string {
	panic("TODO: NumberLines")
}

// Run печатает в out каждый файл из filenames построчно с номерами
// (нумерация с 0001 заново для каждого файла, формат как в NumberLines).
// Возвращает ошибку, если файл не открывается или читается.
func Run(filenames []string, out io.Writer) error {
	panic("TODO: Run")
}
