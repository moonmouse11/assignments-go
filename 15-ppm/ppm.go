//go:build !solution

// Пакет ppm — глава 15. Реализуй функции так, чтобы прошли тесты
// в ppm_test.go. Условие — в README.md этой папки.
package ppm

import "io"

// Формат файла PPMF: заголовок 65 байт, затем AES-GCM-шифротекст.
//
//	[0:4]    magic "PPMF"
//	[4:5]    версия формата, 1
//	[5:21]   salt (16 байт, PBKDF2)
//	[21:33]  nonce GCM (12 байт)
//	[33:65]  sha256 от байтов [0:33]
//	[65:]    AES-GCM(payload), ключ — pbkdf2(sha256, пароль, salt, iter, 32)

// iterations — число итераций PBKDF2. В учебных целях маленькое:
// тесты обязаны выполняться быстрее секунды.
const iterations = 1000

// Save шифрует payload паролем password (pbkdf2 + AES-GCM, случайные
// salt и nonce) и записывает файл формата PPMF в w. Возвращает ошибку
// при сбое записи или ГСЧ.
func Save(payload []byte, password string, w io.Writer) error {
	panic("TODO: Save")
}

// Load читает из r файл формата PPMF и расшифровывает payload паролем
// password. Ошибки: короткий файл, неверный magic/версия, порча
// заголовка (sha256 не сходится), неверный пароль или порча payload
// (ошибка расшифровки GCM).
func Load(r io.Reader, password string) ([]byte, error) {
	panic("TODO: Load")
}
