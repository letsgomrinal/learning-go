// Reference solution for lesson 33.
package lesson33

import (
	"runtime"
	"sync"
	"time"
)

func FirstMatch(items []string, match func(string) bool) string {
	results := make(chan string, len(items)) // room for everyone: no goroutine can get stuck
	var wg sync.WaitGroup
	for _, it := range items {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if match(it) {
				results <- it
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results) // no match at all: the receive below gets "" instead of hanging
	}()
	return <-results
}

func StartHeartbeat(interval time.Duration, beat func()) (stop func()) {
	done := make(chan struct{})
	exited := make(chan struct{})
	go func() {
		defer close(exited)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				beat()
			case <-done:
				return
			}
		}
	}()
	var once sync.Once
	return func() {
		once.Do(func() { close(done) })
		<-exited // wait until the goroutine has really finished
	}
}

type Account struct {
	ID      int
	mu      sync.Mutex
	Balance int
}

func Transfer(from, to *Account, amount int) {
	first, second := from, to
	if to.ID < from.ID {
		first, second = to, from
	}
	first.mu.Lock()
	defer first.mu.Unlock()
	time.Sleep(time.Microsecond)
	second.mu.Lock()
	defer second.mu.Unlock()
	from.Balance -= amount
	to.Balance += amount
}

func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func CountPrimes(limit int) int {
	workers := runtime.GOMAXPROCS(0)
	counts := make([]int, workers)
	var wg sync.WaitGroup
	for w := range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Interleave numbers (w, w+workers, ...) so every worker gets a fair share of the
			// expensive large numbers.
			for i := w; i < limit; i += workers {
				if IsPrime(i) {
					counts[w]++
				}
			}
		}()
	}
	wg.Wait()
	total := 0
	for _, c := range counts {
		total += c
	}
	return total
}
