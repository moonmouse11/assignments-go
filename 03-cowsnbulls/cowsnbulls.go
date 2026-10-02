//go:build !solution

// Пакет cowsnbulls — глава 3. Реализуй Score и решатель так,
// чтобы прошли тесты. Условие — в README.md.
package cowsnbulls

// Score оценивает попытку guess против секрета secret
// (4 цифры, без повторов): bulls — цифра на своём месте,
// cows — правильная цифра не на своём месте.
func Score(secret, guess string) (bulls, cows int) {
	panic("TODO: Score")
}

// GuessingSolver — компьютер угадывает задуманное число,
// вычёркивая кандидатов по ответам человека.
type GuessingSolver struct {
	// TODO: твои поля (например, слайс живых кандидатов).
}

// NewGuessingSolver создаёт решателя со всеми 4-значными
// числами с неповторяющимися цифрами (5040 кандидатов).
func NewGuessingSolver() *GuessingSolver {
	panic("TODO: NewGuessingSolver")
}

// Guess возвращает текущую гипотезу решателя.
func (g *GuessingSolver) Guess() string {
	panic("TODO: Guess")
}

// Learn принимает ответ человека на попытку guess и вычёркивает
// всех кандидатов, дающих другой счёт.
func (g *GuessingSolver) Learn(guess string, bulls, cows int) {
	panic("TODO: Learn")
}
