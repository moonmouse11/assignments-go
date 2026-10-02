//go:build !solution

// Пакет sockgw — глава 25. Реализуй методы Gateway так, чтобы прошли
// тесты в sockgw_test.go. Условие — в README.md этой папки.
package sockgw

import (
	"context"
	"net"
)

// Gateway — шлюз stream-socket → UDP: принимает TCP-соединения,
// каждую строку-запись пишет в лог-файл и пересылает UDP-датаграммой.
// Поля — на твой выбор (listener, UDP-адрес, лог-файл, мьютекс для Close).
type Gateway struct{}

// Listen поднимает TCP-listener на tcpAddr (пустая строка → "127.0.0.1:0").
// Возвращает ошибку, если адрес не слушается.
func (g *Gateway) Listen(tcpAddr string) error {
	panic("TODO: Listen")
}

// Addr возвращает адрес TCP-listener'а (nil, если Listen не вызывался).
func (g *Gateway) Addr() net.Addr {
	panic("TODO: Addr")
}

// SetForward задаёт UDP-адрес пересылки (net.ResolveUDPAddr).
// Возвращает ошибку, если адрес не резолвится.
func (g *Gateway) SetForward(udpAddr string) error {
	panic("TODO: SetForward")
}

// SetLog открывает path для записи лога (создаёт при отсутствии).
// Возвращает ошибку, если файл не открывается.
func (g *Gateway) SetLog(path string) error {
	panic("TODO: SetLog")
}

// Serve блокируется до отмены ctx: accept — в горутине, каждое соединение —
// в своей горутине; по каждой строке — запись в лог (с "\n") и
// UDP-датаграмма. По отмене ctx закрывает всё (Close).
func (g *Gateway) Serve(ctx context.Context) {
	panic("TODO: Serve")
}

// Close закрывает listener, лог и UDP-соединение. Идемпотентно:
// повторный вызов не ошибка.
func (g *Gateway) Close() {
	panic("TODO: Close")
}
