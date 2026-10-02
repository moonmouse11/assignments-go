// Пакет letterboxed — глава 20. Словарь для тестов: только данные, без кода.
package letterboxed

// Words — учебный словарь: строчные слова из букв a–l. Каждое слово
// ВАЛИДНО для учебного пазла со сторонами ["abc","def","ghi","jkl"]:
// соседние буквы всегда лежат на разных сторонах.
// bdjceglfk — служебное «слово» для теста решения из двух слов.
var Words = []string{
	// 3 буквы
	"age", "aid", "ale", "beg", "bid", "did", "die", "eld", "gel", "hag",
	"ice", "jib", "keg", "kid", "lei", "leg", "lie", "lid",
	// 4 буквы
	"blah", "dale", "dice", "dike", "fail", "gale", "hadj", "haji", "hail",
	"jail", "jibe", "kale", "kick", "lake", "lick", "like", "glia", "glib",
	"flag", "flake", "bleak", "bilge",
	// 5 букв
	"ahead", "algae", "alibi", "alike", "galah", "gelid", "glebe", "iliad",
	"khadi", "khaki", "legal", "ladle", "algal", "field", "fjeld",
	// 6 букв
	"afield", "beadle", "behead", "beheld",
	// 9 букв — служебное слово для теста IsSolution
	"bdjceglfk",
}
