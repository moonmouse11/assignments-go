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
