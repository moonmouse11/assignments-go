//go:build solution

// Эталонное решение главы 5.
package exec

import (
	"context"
	"os/exec"
	"time"
)

func RunCommand(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

func Periodic(ctx context.Context, interval time.Duration, fn func()) int {
	t := time.NewTicker(interval)
	defer t.Stop()
	runs := 0
	for {
		select {
		case <-ctx.Done():
			return runs
		case <-t.C:
			fn()
			runs++
		}
	}
}
