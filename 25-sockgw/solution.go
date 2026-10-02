//go:build solution

// Эталонное решение главы 25.
package sockgw

import (
	"bufio"
	"context"
	"net"
	"os"
	"strings"
	"sync"
)

// Gateway — шлюз stream-socket → UDP.
type Gateway struct {
	mu  sync.Mutex
	ln  net.Listener
	fwd *net.UDPAddr
	log *os.File
	udp *net.UDPConn
}

func (g *Gateway) Listen(addr string) (err error) {
	if addr == "" {
		addr = "127.0.0.1:0"
	}
	g.ln, err = net.Listen("tcp", addr)
	return err
}

func (g *Gateway) Addr() net.Addr {
	if g.ln == nil {
		return nil
	}
	return g.ln.Addr()
}

func (g *Gateway) SetForward(addr string) (err error) {
	g.fwd, err = net.ResolveUDPAddr("udp", addr)
	return err
}

func (g *Gateway) SetLog(path string) (err error) {
	g.log, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	return err
}

// Serve блокируется до отмены ctx; accept и соединения — в своих горутинах.
func (g *Gateway) Serve(ctx context.Context) {
	g.mu.Lock()
	if g.fwd != nil && g.udp == nil {
		if udp, err := net.DialUDP("udp", nil, g.fwd); err == nil {
			g.udp = udp
		}
	}
	ln := g.ln
	g.mu.Unlock()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go g.handleConn(conn)
		}
	}()
	<-ctx.Done()
	g.Close()
}

func (g *Gateway) handleConn(conn net.Conn) {
	defer conn.Close()
	for r := bufio.NewReader(conn); ; {
		line, err := r.ReadString('\n')
		if len(line) > 0 {
			g.record(strings.TrimRight(line, "\r\n"))
		}
		if err != nil {
			return
		}
	}
}

// record: лог пишется до отправки UDP; под мьютексом — Close не мешает.
func (g *Gateway) record(rec string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.log != nil {
		g.log.WriteString(rec + "\n")
	}
	if g.udp != nil {
		g.udp.Write([]byte(rec))
	}
}

func (g *Gateway) Close() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.ln != nil {
		g.ln.Close()
	}
	if g.udp != nil {
		g.udp.Close()
	}
	if g.log != nil {
		g.log.Close()
	}
	g.ln, g.udp, g.log = nil, nil, nil
}
