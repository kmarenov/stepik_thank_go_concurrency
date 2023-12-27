package main

import (
	"fmt"
	"math/rand"
	"time"
)

// начало решения
func delay(dur time.Duration, fn func()) func() {
	ch := make(chan time.Time)
	end := time.Now().Add(dur)
	exit := make(chan struct{})
	closed := make(chan struct{}, 1)
	go func(c chan time.Time) {
		for {
			select {
			case <-exit:
				return
			default:
				if time.Now().After(end) {
					ch <- time.Now()
					return
				}
			}
		}
	}(ch)
	go func(c chan time.Time) {
		select {
		case <-ch:
			fn()
			return
		case <-exit:
			return
		}
	}(ch)
	return func() {
		select {
		case <-closed:
			closed <- struct{}{}
			return
		default:
			close(exit)
			closed <- struct{}{}
			return
		}
	}
}

// конец решения

func main() {
	rand.Seed(time.Now().Unix())

	work := func() {
		fmt.Println("work done")
	}

	cancel := delay(100*time.Millisecond, work)

	time.Sleep(10 * time.Millisecond)
	if rand.Float32() < 0.5 {
		cancel()
		fmt.Println("delayed function canceled")
	}
	time.Sleep(100 * time.Millisecond)
}
