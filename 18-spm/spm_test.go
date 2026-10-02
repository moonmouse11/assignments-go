package spm

import (
	"crypto/ed25519"
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"
)

// fixtures создаёт srcDir с тремя файлами и ожидаемое содержимое.
func fixtures(t *testing.T) (srcDir string, want map[string]string) {
	t.Helper()
	srcDir = t.TempDir()
	want = map[string]string{
		"app.conf": "listen = 1024\n",
		"run.sh":   "#!/bin/sh\necho ok\n",
		"docs.txt": "сервисный пакет\n",
	}
	for name, content := range want {
		if err := os.WriteFile(filepath.Join(srcDir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("фикстура %s: %v", name, err)
		}
	}
	return srcDir, want
}

func keypair(t *testing.T) (ed25519.PublicKey, ed25519.PrivateKey) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	return pub, priv
}

func signedPkg(t *testing.T) (pkg, signed []byte, pub ed25519.PublicKey, want map[string]string) {
	t.Helper()
	srcDir, want := fixtures(t)
	pkg, err := Build(srcDir)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	pub, priv := keypair(t)
	signed, err = Sign(pkg, priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	return pkg, signed, pub, want
}

func TestBuildSignInstall(t *testing.T) {
	_, signed, pub, want := signedPkg(t)

	dst := t.TempDir()
	if err := Install(signed, pub, dst); err != nil {
		t.Fatalf("Install: %v", err)
	}
	for name, content := range want {
		got, err := os.ReadFile(filepath.Join(dst, name))
		if err != nil {
			t.Fatalf("прочитать %s после Install: %v", name, err)
		}
		if string(got) != content {
			t.Errorf("файл %s: получил %q, хотел %q", name, got, content)
		}
	}
	if err := Verify(signed, pub); err != nil {
		t.Errorf("Verify чистого пакета: %v", err)
	}
}

func TestVerifyCorruptPackage(t *testing.T) {
	pkg, _, pub, _ := signedPkg(t)
	_, priv := keypair(t)

	pkg[len(pkg)/2] ^= 0xFF // порча одного байта файла внутри пакета
	bad, err := Sign(pkg, priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if err := Verify(bad, pub); err == nil {
		t.Fatal("Verify: хотел ошибку для пакета с испорченным файлом, получил nil")
	}
}

func TestVerifyCorruptSignature(t *testing.T) {
	_, signed, pub, _ := signedPkg(t)

	bad := append([]byte{}, signed...)
	bad[len(bad)-1] ^= 0xFF // переворот последнего байта подписи
	if err := Verify(bad, pub); err == nil {
		t.Fatal("Verify: хотел ошибку для испорченной подписи, получил nil")
	}
}

func TestVerifyWrongKey(t *testing.T) {
	_, signed, _, _ := signedPkg(t)
	otherPub, _ := keypair(t)

	if err := Verify(signed, otherPub); err == nil {
		t.Fatal("Verify: хотел ошибку для чужого ключа, получил nil")
	}
}

func TestInstallRefusesCorrupt(t *testing.T) {
	_, signed, pub, _ := signedPkg(t)

	bad := append([]byte{}, signed...)
	bad[len(bad)/2] ^= 0xFF // порча байта внутри пакета
	dst := t.TempDir()
	if err := Install(bad, pub, dst); err == nil {
		t.Fatal("Install: хотел ошибку для испорченного пакета, получил nil")
	}
	// dst должен остаться без файлов распаковки.
	entries, _ := os.ReadDir(dst)
	if len(entries) != 0 {
		t.Errorf("Install распаковал файлы до ошибки: %d записей", len(entries))
	}
}
