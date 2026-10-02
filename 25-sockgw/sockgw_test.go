package sockgw

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGatewayErrors(t *testing.T) {
	var g Gateway
	if err := g.SetForward("bogus-addr"); err == nil {
		t.Error("SetForward: хотел ошибку для кривого адреса")
	}
	if err := g.Listen("127.0.0.1:99999"); err == nil {
		t.Error("Listen: хотел ошибку для кривого порта")
	}
	bad := filepath.Join(t.TempDir(), "no", "such", "dir", "l.log")
	if err := g.SetLog(bad); err == nil {
		t.Error("SetLog: хотел ошибку для отсутствующего каталога")
	}
	if g.Addr() != nil {
		t.Error("Addr: хотел nil до Listen")
	}
}

func TestGatewayForward(t *testing.T) {
	// UDP-приёмник.
	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("ListenPacket: %v", err)
	}
	defer pc.Close()

	// Шлюз.
	var g Gateway
	defer g.Close()
	if err := g.Listen("127.0.0.1:0"); err != nil {
		t.Fatalf("Listen: %v", err)
	}
	tcp := g.Addr()
	if tcp == nil || tcp.(*net.TCPAddr).Port == 0 {
		t.Fatalf("Addr: хотел выданный порт, получил %v", tcp)
	}
	if err := g.SetForward(pc.LocalAddr().String()); err != nil {
		t.Fatalf("SetForward: %v", err)
	}
	logPath := filepath.Join(t.TempDir(), "archive.bin")
	if err := g.SetLog(logPath); err != nil {
		t.Fatalf("SetLog: %v", err)
	}

	// Serve в горутине с ctx.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	served := make(chan struct{})
	go func() {
		g.Serve(ctx)
		close(served)
	}()

	// Клиент: три строки.
	conn, err := net.DialTimeout("tcp", tcp.String(), 2*time.Second)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	lines := []string{"первая запись", "second line", "No: 3"}
	for _, l := range lines {
		fmt.Fprintf(conn, "%s\n", l)
	}
	conn.Close()

	// Три UDP-датаграммы (с дедлайном, чтобы не виснуть).
	buf := make([]byte, 1024)
	var got []string
	if err := pc.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	for i := 0; i < len(lines); i++ {
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			t.Fatalf("ReadFrom (датаграмма %d): %v", i+1, err)
		}
		got = append(got, string(buf[:n]))
	}
	if strings.Join(got, "|") != strings.Join(lines, "|") {
		t.Errorf("датаграммы:\nполучил %q\nхотел  %q", got, lines)
	}

	// Лог: строка с \n, в том же порядке (лог пишется до отправки UDP).
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile(лог): %v", err)
	}
	if want := strings.Join(lines, "\n") + "\n"; string(data) != want {
		t.Errorf("лог:\nполучил %q\nхотел  %q", data, want)
	}

	// Отмена ctx завершает Serve.
	cancel()
	select {
	case <-served:
	case <-time.After(2 * time.Second):
		t.Fatal("Serve не завершился после отмены ctx")
	}
}
