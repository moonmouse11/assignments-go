//go:build !solution

// Пакет password — глава 16. Реализуй типы и методы так, чтобы прошли
// тесты в password_test.go. Условие — в README.md этой папки.
package password

// Entry — одна запись менеджера паролей: контекст (сайт/сервис),
// имя пользователя и пароль для него. Ключ — Context.
type Entry struct {
	Context  string
	Username string
	Password string
}

// Vault — хранилище записей. Ключ — Context записи.
type Vault struct {
	// TODO: поле для хранения записей (map по Context).
}

// NewVault возвращает пустое хранилище.
func NewVault() *Vault {
	panic("TODO: NewVault")
}

// Add добавляет запись e. Пустой Context — ошибка. Повторный Add
// того же Context обновляет запись (UN-4, UN-5).
func (v *Vault) Add(e Entry) error {
	panic("TODO: Add")
}

// Get возвращает запись по контексту context; ok=false, если записи нет.
func (v *Vault) Get(context string) (Entry, bool) {
	panic("TODO: Get")
}

// List возвращает все записи, отсортированные по алфавиту контекстов.
func (v *Vault) List() []Entry {
	panic("TODO: List")
}

// SaveTo сериализует хранилище в JSON и записывает в path зашифрованным
// (pbkdf2-hmac-sha256 + AES-GCM, реализация локальная в пакете) файлом
// под мастер-паролем master. Заголовок файла: magic "PSWD", версия,
// salt[16], nonce[12], sha256 заголовка; затем шифротекст.
func (v *Vault) SaveTo(path, master string) error {
	panic("TODO: SaveTo")
}

// LoadFrom читает зашифрованный файл path, расшифровывает его мастер-паролем
// master и возвращает хранилище. Ошибки: неверный формат/портченый файл,
// неверный мастер-пароль.
func LoadFrom(path, master string) (*Vault, error) {
	panic("TODO: LoadFrom")
}
