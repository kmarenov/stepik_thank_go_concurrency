package main

import (
	"fmt"
	"math/rand"
)

// начало решения

// генерит случайные слова из 5 букв
// с помощью randomWord(5)
func generate(cancel <-chan struct{}) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case out <- randomWord(5):
			case <-cancel:
				return
			}
		}
	}()
	return out
}

func isUnique(word string) bool {
	chars := make(map[rune]struct{})

	for _, char := range word {
		if _, ok := chars[char]; ok {
			return false
		}
		chars[char] = struct{}{}
	}

	return true
}

// выбирает слова, в которых не повторяются буквы,
// abcde - подходит
// abcda - не подходит
func takeUnique(cancel <-chan struct{}, in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for in != nil {
			select {
			case word, ok := <-in:
				if ok {
					if isUnique(word) {
						select {
						case out <- word:
						case <-cancel:
							return
						}
					}
				} else {
					in = nil
				}
			case <-cancel:
				return
			}
		}
	}()
	return out
}

func rev(word string) string {
	var r string
	for _, c := range word {
		r = string(c) + r
	}
	return r
}

type pair struct {
	forv string
	revt string
}

// переворачивает слова
// abcde -> edcba
func reverse(cancel <-chan struct{}, in <-chan string) <-chan pair {
	out := make(chan pair)
	go func() {
		defer close(out)
		for in != nil {
			select {
			case word, ok := <-in:
				if ok {
					select {
					case out <- pair{forv: word, revt: rev(word)}:
					case <-cancel:
						return
					}
				} else {
					in = nil
				}
			case <-cancel:
				return
			}
		}
	}()
	return out
}

// объединяет c1 и c2 в общий канал
func merge(cancel <-chan struct{}, c1, c2 <-chan pair) <-chan pair {
	out := make(chan pair)

	go func() {
		defer close(out)
		for c1 != nil || c2 != nil {
			select {
			case p1, ok := <-c1:
				if ok {
					select {
					case out <- p1:
					case <-cancel:
						return
					}
				} else {
					c1 = nil
				}
			case p2, ok := <-c2:
				if ok {
					select {
					case out <- p2:
					case <-cancel:
						return
					}
				} else {
					c2 = nil
				}
			case <-cancel:
				return
			}
		}
	}()

	return out
}

// печатает первые n результатов
func print(cancel <-chan struct{}, in <-chan pair, n int) {
	for i := 0; i < n; i++ {
		if in != nil {
			select {
			case p, ok := <-in:
				if ok {
					fmt.Println(p.forv + " -> " + p.revt)
				} else {
					in = nil
				}
			case <-cancel:
				return
			}
		}
	}
}

// конец решения

// генерит случайное слово из n букв
func randomWord(n int) string {
	const letters = "aeiourtnsl"
	chars := make([]byte, n)
	for i := range chars {
		chars[i] = letters[rand.Intn(len(letters))]
	}
	return string(chars)
}

func main() {
	cancel := make(chan struct{})
	defer close(cancel)

	c1 := generate(cancel)
	c2 := takeUnique(cancel, c1)
	c3_1 := reverse(cancel, c2)
	c3_2 := reverse(cancel, c2)
	c4 := merge(cancel, c3_1, c3_2)
	print(cancel, c4, 10)
}
