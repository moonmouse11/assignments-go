# Глава 18. SPM — защищённый менеджер пакетов

> По мотивам `spm.pdf`. Пакет = tar.gz с файлами и манифестом, подписанный
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
