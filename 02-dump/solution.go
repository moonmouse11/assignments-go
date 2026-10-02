//go:build solution

// Эталонное решение главы 2.
package dump

import (
	"fmt"
	"io"
	"strings"
)

const Version = "dump-0.1"

func FormatBlock(offset int, data []byte) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%04x:", offset)
	for _, c := range data {
		fmt.Fprintf(&b, " %02x", c)
	}
	b.WriteString(" |")
	for _, c := range data {
		if c >= 0x20 && c < 0x7f {
			b.WriteByte(c)
		} else {
			b.WriteByte('.')
		}
	}
	b.WriteString("|\n")
	return b.String()
}

func Dump(r io.Reader, blockSize int, out io.Writer) error {
	switch blockSize {
	case 8, 16, 32:
	default:
		return fmt.Errorf("недопустимый размер блока %d (допустимо 8, 16, 32)", blockSize)
	}
	buf := make([]byte, blockSize)
	for off := 0; ; off += blockSize {
		n, err := io.ReadFull(r, buf)
		if n > 0 {
			fmt.Fprint(out, FormatBlock(off, buf[:n]))
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
