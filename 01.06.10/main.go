package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrCanceled error = errors.New("canceled")

// начало решения
func withRateLimit(limit int, fn func()) (handle func() error, cancel func()) {
	dur := time.Second / time.Duration(limit)
	ticker := time.NewTicker(dur)

	closed := make(chan struct{}, 1)
	calls := make(chan struct{})
	done := make(chan struct{})

	go func() {
		defer func() {
			close(done)
		}()
		for range calls {
			<-ticker.C
			go fn()
		}
		ticker.Stop()
		return
	}()

	handle = func() error {
		select {
		case <-closed:
			closed <- struct{}{}
			return ErrCanceled
		default:
			calls <- struct{}{}
			return nil
		}
	}

	cancel = func() {
		select {
		case <-closed:
			closed <- struct{}{}
			return
		default:
			close(calls)
			<-done
			closed <- struct{}{}
			return
		}
	}

	return handle, cancel
}

// конец решения

func main() {
	work := func() {
		fmt.Print(".")
	}

	handle, cancel := withRateLimit(5, work)
	defer cancel()

	start := time.Now()
	const n = 10
	for i := 0; i < n; i++ {
		handle()
	}
	fmt.Println()
	fmt.Printf("%d queries took %v\n", n, time.Since(start))
}
