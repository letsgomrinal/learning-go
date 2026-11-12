// Lesson 33 exercises: the scheduler, leaks and deadlocks. Read README.md first!
//
// Check one exercise at a time:   go test -race -run Ex1 -v
package lesson33

import (
	"sync"
	"time"
)

// Ex1: return any item for which match returns true ("" if none), checking all items concurrently.
// This version works, but it LEAKS goroutines (README section 5). Fix the leak.
// Also make sure it returns "" (and doesn't hang) when nothing matches.
//
// Hint: a channel buffered with room for every goroutine, plus a WaitGroup so you know when
// all of them are done (for the "no match" case: close the channel after wg.Wait() in a goroutine).
func FirstMatch(items []string, match func(string) bool) string {
	results := make(chan string)
	for _, it := range items {
		go func() {
			if match(it) {
				results <- it
			}
		}()
	}
	return <-results
}

// Ex2: start a goroutine that calls beat() every interval, and return a stop function.
// After stop() returns, beat must never be called again and the goroutine must have exited.
// Calling stop() twice must not panic.
//
// Hints:
//   - a done channel (chan struct{}) and a time.Ticker (remember ticker.Stop())
//   - the goroutine: for { select { case <-ticker.C: beat()  case <-done: return } }
//   - to wait until the goroutine has really exited, use a second channel it closes on the way out
//   - sync.Once makes stop() safe to call twice
func StartHeartbeat(interval time.Duration, beat func()) (stop func()) {
	// TODO
	return func() {}
}

// Account is a bank account protected by its own mutex.
type Account struct {
	ID      int
	mu      sync.Mutex
	Balance int
}

// Ex3: move amount from one account to another. Two transfers in OPPOSITE directions at the same
// time can deadlock (README section 6). Fix it by always locking the account with the lower ID first.
func Transfer(from, to *Account, amount int) {
	from.mu.Lock()
	defer from.mu.Unlock()
	time.Sleep(time.Microsecond) // makes the deadlock much easier to hit
	to.mu.Lock()
	defer to.mu.Unlock()
	from.Balance -= amount
	to.Balance += amount
}

// IsPrime is a deliberately simple CPU-bound function.
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

// Ex4: count the primes below limit, splitting the range 0..limit between runtime.GOMAXPROCS(0)
// worker goroutines (CPU-bound work: one worker per core, README section 4).
// Each worker counts its own part into its own slot; add them up at the end.
// It must give the same answer as a simple loop, and be faster on a multi-core machine.
func CountPrimes(limit int) int {
	// TODO (this is the simple single-goroutine version)
	n := 0
	for i := range limit {
		if IsPrime(i) {
			n++
		}
	}
	return n
}
