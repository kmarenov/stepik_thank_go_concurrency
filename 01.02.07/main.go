package main

import (
	"fmt"
	"strings"
	"unicode"
)

// nextFunc возвращает следующее слово из генератора
type nextFunc func() string

// counter хранит количество цифр в каждом слове.
// ключ карты - слово, а значение - количество цифр в слове.
type counter map[string]int

// pair хранит слово и количество цифр в нем
type pair struct {
	word  string
	count int
}

// countDigitsInWords считает количество цифр в словах,
// выбирая очередные слова с помощью next()
func countDigitsInWords(next nextFunc) counter {
	counted := make(chan pair)

	// начало решения

	go func() {
		// Пройдите по словам,
		// посчитайте количество цифр в каждом,
		// и запишите его в канал counted
		for {
			word := next()
			count := countDigits(word)
			counted <- pair{
				word:  word,
				count: count,
			}
			if word == "" {
				return
			}
		}
	}()

	// Считайте значения из канала counted
	// и заполните stats.
	stats := counter{}

	// В результате stats должна содержать слова
	// и количество цифр в каждом.
	for {
		p := <-counted
		if p.word == "" {
			break
		}
		stats[p.word] = p.count
	}

	// конец решения

	return stats
}

// countDigits возвращает количество цифр в строке
func countDigits(str string) int {
	count := 0
	for _, char := range str {
		if unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

// printStats печатает слова и количество цифр в каждом
func printStats(stats counter) {
	for word, count := range stats {
		fmt.Printf("%s: %d\n", word, count)
	}
}

// wordGenerator возвращает генератор, который выдает слова из фразы
func wordGenerator(phrase string) nextFunc {
	words := strings.Fields(phrase)
	idx := 0
	return func() string {
		if idx == len(words) {
			return ""
		}
		word := words[idx]
		idx++
		return word
	}
}

func main() {
	phrase := "0ne 1wo thr33 4068"
	next := wordGenerator(phrase)
	stats := countDigitsInWords(next)
	printStats(stats)
}
