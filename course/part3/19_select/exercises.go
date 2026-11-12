// Lesson 19 exercises: select and timeouts. Read README.md first!
//
// Check one exercise at a time:   go test -race -run Ex1 -v
package lesson19

import "time"

// Ex1: wait for a value from a or b, and return whichever arrives first.
// Use a select with two cases.
func First(a, b <-chan string) string {
	// TODO
	return ""
}

// Ex2: wait for a value from ch for at most d.
// Return the value and true if it arrives in time, or 0 and false if the time runs out.
// (Also return 0 and false if ch is closed.)
//
// Hint: select with case v, ok := <-ch and case <-time.After(d)
func ReceiveWithTimeout(ch <-chan int, d time.Duration) (int, bool) {
	// TODO
	return 0, false
}

// Ex3: try to send v on ch WITHOUT waiting.
// Return true if it was sent, false if it would have had to wait.
// Hint: select with a default case (README section 3).
func TrySend(ch chan<- int, v int) bool {
	// TODO
	return false
}

// Ex4: try to receive from ch WITHOUT waiting.
// Return the value and true if one was ready (and ch isn't closed), or 0 and false otherwise.
func TryReceive(ch <-chan int) (int, bool) {
	// TODO
	return 0, false
}

// Ex5: run work() in a goroutine. Return its result and true if it finishes within d,
// or 0 and false if it doesn't.
//
// The test checks that you don't leak the goroutine when the timeout wins. (README section 2.)
func RunWithTimeout(work func() int, d time.Duration) (int, bool) {
	// TODO
	return 0, false
}

// Ex6: count how many times a ticker with the given interval ticks before `total` has passed.
// CountTicks(10*time.Millisecond, 55*time.Millisecond) is about 5.
//
// Hint (README section 4):
//
//	ticker := time.NewTicker(interval)
//	defer ticker.Stop()
//	timeout := time.After(total)
//	for { select { case <-ticker.C: count++  case <-timeout: return count } }
func CountTicks(interval, total time.Duration) int {
	// TODO
	return 0
}

// Ex7: return a channel that produces 0, 1, 2, 3, ... until done is closed.
// When done is closed, the goroutine must stop AND close the returned channel.
// (README section 5.)
func Counter(done <-chan struct{}) <-chan int {
	// TODO
	return nil
}
