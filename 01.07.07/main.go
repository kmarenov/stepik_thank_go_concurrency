package main

import (
	"context"
	"fmt"
	"strings"
	"unicode"
)

// информация о количестве цифр в каждом слове
type counter map[string]int

// слово и количество цифр в нем
type pair struct {
	word  string
	count int
}

// начало решения

// считает количество цифр в словах
func countDigitsInWords(ctx context.Context, words []string) counter {
	pending := submitWords(ctx, words)
	counted := countWords(ctx, pending)
	return fillStats(ctx, counted)
}

// отправляет слова на подсчет
func submitWords(ctx context.Context, words []string) <-chan string {
	out := make(chan string)
	if len(words) == 0 {
		close(out)
		return out
	}
	i := 0
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			case out <- words[i]:
				i++
				if i == len(words) {
					close(out)
					return
				}
			}
		}
	}()
	return out
}

// считает цифры в словах
func countWords(ctx context.Context, in <-chan string) <-chan pair {
	out := make(chan pair)
	go func() {
		for {
			select {
			case word, ok := <-in:
				if !ok {
					close(out)
					return
				}
				select {
				case out <- pair{word, countDigits(word)}:
				case <-ctx.Done():
					close(out)
					return
				}
			case <-ctx.Done():
				close(out)
				return
			}
		}
	}()
	return out
}

// готовит итоговую статистику
func fillStats(ctx context.Context, in <-chan pair) counter {
	stats := counter{}
	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		for {
			select {
			case p, ok := <-in:
				if !ok {
					return
				}
				stats[p.word] = p.count
			case <-ctx.Done():
				return
			}
		}
	}()
	<-done
	return stats
}

// конец решения

// считает количество цифр в слове
func countDigits(str string) int {
	count := 0
	for _, char := range str {
		if unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

func main() {
	phrase := "0ne 1wo thr33 4068"
	words := strings.Fields(phrase)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stats := countDigitsInWords(ctx, words)
	fmt.Println(stats)
}
