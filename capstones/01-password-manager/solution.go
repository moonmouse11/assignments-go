//go:build solution

// Эталон капстоуна 1.
package passwordmanager

import (
	"fmt"

	"assignments-go/10-pwdgen"
	"assignments-go/16-password"
)

type Manager struct {
	vault  *password.Vault
	gen    *pwdgen.Gen
	master string
}

func New(seed uint64, words []string, master string) *Manager {
	return &Manager{
		vault:  password.NewVault(),
		gen:    pwdgen.NewGen(seed, words, "-"),
		master: master,
	}
}

func (m *Manager) Generate(context, username string, nSeg int) (string, error) {
	if m.gen == nil {
		return "", fmt.Errorf("генератор не настроен (менеджер открыт из файла)")
	}
	pw := m.gen.Password(nSeg)
	err := m.vault.Add(password.Entry{Context: context, Username: username, Password: pw})
	if err != nil {
		return "", err
	}
	return pw, nil
}

func (m *Manager) Get(context string) (string, bool) {
	e, ok := m.vault.Get(context)
	if !ok {
		return "", false
	}
	return e.Password, true
}

func (m *Manager) Save(path string) error {
	return m.vault.SaveTo(path, m.master)
}

func Open(path, master string) (*Manager, error) {
	v, err := password.LoadFrom(path, master)
	if err != nil {
		return nil, err
	}
	return &Manager{vault: v, master: master}, nil
}

func (m *Manager) List() []password.Entry {
	return m.vault.List()
}
