//go:build !solution

// Капстоун 1: менеджер паролей. Реализуй клей поверх 10-pwdgen и 16-password
// так, чтобы прошли тесты в manager_test.go. Условие — в README.md.
package passwordmanager

import (
	"assignments-go/10-pwdgen"
	"assignments-go/16-password"
)

// var _ — держит импорт pwdgen: понадобится в New (см. решение).
var _ = pwdgen.NewGen

// Manager — менеджер паролей: хранилище (16-password) + генератор (10-pwdgen)
// + master-пароль. У открытого из файла генератора нет.
type Manager struct {
	// TODO: твои поля.
}

// New создаёт менеджер: свежее хранилище и генератор (разделитель "-").
func New(seed uint64, words []string, master string) *Manager {
	panic("TODO: New")
}

// Generate генерирует пароль из nSeg сегментов, кладёт его в хранилище
// под context/username и возвращает. Пустой context — ошибка.
// У менеджера без генератора (Open) — ошибка.
func (m *Manager) Generate(context, username string, nSeg int) (string, error) {
	panic("TODO: Generate")
}

// Get возвращает пароль контекста.
func (m *Manager) Get(context string) (string, bool) {
	panic("TODO: Get")
}

// Save пишет хранилище в защищённый файл path с master-паролем менеджера.
func (m *Manager) Save(path string) error {
	panic("TODO: Save")
}

// Open читает менеджер из защищённого файла. Неверный master — ошибка.
func Open(path, master string) (*Manager, error) {
	panic("TODO: Open")
}

// List возвращает записи по алфавиту контекстов.
func (m *Manager) List() []password.Entry {
	panic("TODO: List")
}
