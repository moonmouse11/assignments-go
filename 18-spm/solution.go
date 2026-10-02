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
