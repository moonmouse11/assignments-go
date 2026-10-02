package exec

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestRunCommand(t *testing.T) {
	out, err := RunCommand("echo", "hello", "gobook")
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}
	if !bytes.Equal(out, []byte("hello gobook\n")) {
		t.Errorf("RunCommand: получил %q, хотел %q", out, "hello gobook\n")
	}
}

func TestRunCommandError(t *testing.T) {
	if _, err := RunCommand("no-such-command-xyz"); err == nil {
		t.Error("RunCommand: хотел ошибку для несуществующей команды")
	}
}

func TestPeriodicCounts(t *testing.T) {
	// Тики каждые 20 мс, отмена через 70 мс → исполнений 2..3.
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(70 * time.Millisecond)
		cancel()
	}()
	runs := Periodic(ctx, 20*time.Millisecond, func() {})
	if runs < 2 || runs > 3 {
		t.Errorf("Periodic: исполнений %d, хотел 2..3", runs)
	}
}

func TestPeriodicCancelImmediately(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	runs := Periodic(ctx, time.Hour, func() { t.Error("fn не должен выполниться после отмены") })
	if runs != 0 {
		t.Errorf("Periodic: исполнений %d, хотел 0", runs)
	}
}
