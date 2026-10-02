//go:build solution

// Эталонное решение главы 1. Сравни со своей реализацией в hello.go.
package hello

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func NumberLines(lines []string) string {
	var b strings.Builder
	for i, line := range lines {
		fmt.Fprintf(&b, "%04d %s\n", i+1, line)
	}
	return b.String()
}

func Run(filenames []string, out io.Writer) error {
	for _, name := range filenames {
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		n := 0
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			n++
			fmt.Fprintf(out, "%04d %s\n", n, sc.Text())
		}
		if err := sc.Err(); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}
