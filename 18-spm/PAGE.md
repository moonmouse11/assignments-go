# Глава 18. SPM — защищённый менеджер пакетов

> Пакет = tar.gz с файлами и манифестом, подписанный
> Ed25519: аутентичность и целостность при распространении.

## Цель главы

Собрать конвейер распространения: архивирование (`archive/tar`,
`compress/gzip`), манифест с хешами (`crypto/sha256`), цифровая подпись
(`crypto/ed25519`) и контрольная распаковка.

## Теория

- Пакет — tar.gz: файлы из дерева исходников + `MANIFEST.json` с именем,
  размером и SHA-256 каждого файла.
- Ed25519 подписывает **весь** пакет: байт подписи (64) прикладывается
  в конец — `pkg || sig`. Любая порча пакета ломает подпись.
- `Verify` проверяет и подпись, и каждый файл против манифеста:
  двойная защита, как в реальных пакетных менеджерах.

## Условие

1. `Build(srcDir)` собирает tar.gz (рекурсивно, имена относительно srcDir,
   прямой слэш) с `MANIFEST.json`; файлы в архиве упорядочены по имени.
2. `Sign(pkg, priv)` возвращает `pkg || sig` (Ed25519, 64 байта).
3. `Verify(signedPkg, pub)` проверяет подпись и затем хеши всех файлов
   против манифеста; любая порча — ошибка.
4. `Install(signedPkg, pub, dstDir)` — `Verify` и распаковка файлов
   в dstDir (кроме манифеста).

## API для реализации

```go
func Build(srcDir string) ([]byte, error)
func Sign(pkg []byte, priv ed25519.PrivateKey) ([]byte, error)
func Verify(signedPkg []byte, pub ed25519.PublicKey) error
func Install(signedPkg []byte, pub ed25519.PublicKey, dstDir string) error
```

Договоры — в `spm.go`. Тесты — в `spm_test.go`.

## Подсказки

- `ed25519.GenerateKey(rand.Reader)` в тесте; `filepath.WalkDir` в Build.
- `tar.Header{Name, Mode: 0o644, Size, Typeflag: tar.TypeReg}`; MANIFEST
  пишите первым файлом — так проще проверять.
- В `Verify` соберите `map[string][]byte` из архива и сравните
  с манифестом в обе стороны.

## Проверка

```
go test ./18-spm/                 # красный: TODO
go test -tags solution ./18-spm/ # зелёный: эталон
```

## Эталон

`solution.go` собирается только с тегом `solution`.

## Углубление

- YAML-конфигурация сборки (`gopkg.in/yaml.v3`, внешняя зависимость): имя
  пакета, версия, список `from → to` и скрипты `preinstall`/`postinstall`.
- Подпись каждого файла отдельно (`file` + `file.sig`) и подпись манифеста
  `Packagefile`; публичный ключ лежит внутри контейнера — в главе вместо
  этого одна подпись Ed25519 всего пакета.
- Контейнер целиком шифруется паролем из переменной окружения
  `SPM_PKGPASSWORD` (`os.Getenv`) — реюз защищённого формата из главы 15.
- CLI с подкомандами `build`/`install`/`version`, флаги `--config`
  (по умолчанию `$HOME/.spm.yaml`), `--keep`, `install --show`;
  preinstall/postinstall запускаются через `os/exec`.
- Манифест несёт UUID пакета (v4 собери из `crypto/rand`); оригинал также
  использует PEM-ключи RSA/OpenSSL и публикацию артефактов на distribution
  server (сеть/cloud, внешние зависимости).

---

## Исходники на этой странице

Ниже — те же файлы, что в репозитории. Локально: клонируйте курс и правьте код в этой папке.

```bash
git clone https://github.com/moonmouse11/assignments-go.git
cd assignments-go
go test ./18-spm/
```

## Заготовки (реализуйте)

### `spm.go`

{% raw %}
```go
//go:build !solution

// Пакет spm — глава 18. Реализуй функции так, чтобы прошли тесты
// в spm_test.go. Условие — в README.md этой папки.
package spm

import "crypto/ed25519"

// Структура пакета: tar.gz со всеми файлами srcDir (имена относительно
// srcDir, прямой слэш) + файл MANIFEST.json — JSON со списком записей
// {name, size, sha256} каждого файла, упорядоченных по имени.
//
// Подписанный пакет: pkg || sig, где sig — подпись Ed25519 (64 байта)
// от байтов pkg.

// Build собирает tar.gz-пакет из дерева srcDir: рекурсивный обход,
// MANIFEST.json генерируется и кладётся в архив первым файлом.
// Пустой srcDir (без файлов) — ошибка.
func Build(srcDir string) ([]byte, error) {
	panic("TODO: Build")
}

// Sign подписывает пакет pkg приватным ключом priv (Ed25519)
// и возвращает pkg || sig.
func Sign(pkg []byte, priv ed25519.PrivateKey) ([]byte, error) {
	panic("TODO: Sign")
}

// Verify проверяет signedPkg публичным ключом pub: сначала подпись
// Ed25519 поверх байтов пакета, затем манифест (размеры и sha256 всех
// файлов архива). Ошибки: короткий пакет, неверная подпись, нет
// манифеста, несовпадение файла с манифестом, лишний файл в архиве.
func Verify(signedPkg []byte, pub ed25519.PublicKey) error {
	panic("TODO: Verify")
}

// Install проверяет signedPkg (как Verify) и распаковывает все файлы
// пакета, кроме MANIFEST.json, в dstDir (подкаталоги создаются).
func Install(signedPkg []byte, pub ed25519.PublicKey, dstDir string) error {
	panic("TODO: Install")
}
```
{% endraw %}

## Тесты

### `spm_test.go`

{% raw %}
```go
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
```
{% endraw %}

## Эталонное решение

Сначала попробуйте сами. Решение скрыто — раскрывайте только для сверки.

<details>
<summary><strong>Эталон (`solution.go`)</strong> — нажмите, чтобы показать</summary>

{% raw %}
```go
//go:build solution

// Эталонное решение главы 18. Сравни со своей реализацией в spm.go.
package spm

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const manifestName = "MANIFEST.json"

// fileRec — запись манифеста: имя, размер и sha256 файла.
type fileRec struct {
	Name   string
	Size   int64
	SHA256 string
}

// manifest — JSON-манифест пакета.
type manifest struct {
	Files []fileRec `json:"files"`
}

func Build(srcDir string) ([]byte, error) {
	var names []string
	files := map[string][]byte{}
	err := filepath.WalkDir(srcDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, e := os.ReadFile(p)
		if e != nil {
			return e
		}
		rel, e := filepath.Rel(srcDir, p)
		if e != nil {
			return e
		}
		files[filepath.ToSlash(rel)] = data
		names = append(names, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("в %s нет файлов для пакета", srcDir)
	}
	recs := make([]fileRec, len(names))
	for i, n := range names {
		s := sha256.Sum256(files[n])
		recs[i] = fileRec{n, int64(len(files[n])), hex.EncodeToString(s[:])}
	}
	man, err := json.Marshal(manifest{recs})
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(zw)
	write := func(n string, d []byte) error {
		if err := tw.WriteHeader(&tar.Header{Typeflag: tar.TypeReg, Name: n, Mode: 0o644, Size: int64(len(d))}); err != nil {
			return err
		}
		_, err := tw.Write(d)
		return err
	}
	if err := write(manifestName, man); err != nil {
		return nil, err
	}
	for _, n := range names {
		if err := write(n, files[n]); err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Sign(pkg []byte, priv ed25519.PrivateKey) ([]byte, error) {
	if len(priv) != ed25519.PrivateKeySize {
		return nil, errors.New("неверная длина приватного ключа")
	}
	return append(append([]byte{}, pkg...), ed25519.Sign(priv, pkg)...), nil
}

// verifiedFiles проверяет подпись и манифест, возвращает файлы пакета.
func verifiedFiles(signedPkg []byte, pub ed25519.PublicKey) (map[string][]byte, error) {
	if len(signedPkg) < ed25519.SignatureSize {
		return nil, errors.New("пакет короче подписи")
	}
	cut := len(signedPkg) - ed25519.SignatureSize
	if !ed25519.Verify(pub, signedPkg[:cut], signedPkg[cut:]) {
		return nil, errors.New("неверная подпись пакета")
	}
	zr, err := gzip.NewReader(bytes.NewReader(signedPkg[:cut]))
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{}
	tr := tar.NewReader(zr)
	for {
		h, err := tr.Next()
		if err != nil {
			break // io.EOF — конец архива
		}
		files[h.Name], err = io.ReadAll(tr)
		if err != nil {
			return nil, err
		}
	}
	var mf manifest
	if err := json.Unmarshal(files[manifestName], &mf); err != nil {
		return nil, fmt.Errorf("нет манифеста или не читается: %w", err)
	}
	for _, r := range mf.Files {
		data, ok := files[r.Name]
		if !ok {
			return nil, fmt.Errorf("в архиве нет файла %q из манифеста", r.Name)
		}
		s := sha256.Sum256(data)
		if int64(len(data)) != r.Size || hex.EncodeToString(s[:]) != r.SHA256 {
			return nil, fmt.Errorf("файл %q не сходится с манифестом (размер или sha256)", r.Name)
		}
	}
	return files, zr.Close()
}

func Verify(signedPkg []byte, pub ed25519.PublicKey) error {
	_, err := verifiedFiles(signedPkg, pub)
	return err
}

func Install(signedPkg []byte, pub ed25519.PublicKey, dstDir string) error {
	files, err := verifiedFiles(signedPkg, pub)
	if err != nil {
		return err
	}
	for name, data := range files {
		if name == manifestName {
			continue
		}
		dst := filepath.Join(dstDir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
```
{% endraw %}


</details>

_Пакет `assignments-go/18-spm`, упражнение **spm**. Локально: `go test -tags solution ./18-spm/`._
