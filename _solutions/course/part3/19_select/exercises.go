// Reference solution for lesson 19.
package lesson19

import "time"

func First(a, b <-chan string) string {
	select {
	case v := <-a:
		return v
	case v := <-b:
		return v
	}
}

func ReceiveWithTimeout(ch <-chan int, d time.Duration) (int, bool) {
	select {
	case v, ok := <-ch:
		if !ok {
			return 0, false
		}
		return v, true
	case <-time.After(d):
		return 0, false
	}
}

func TrySend(ch chan<- int, v int) bool {
	select {
	case ch <- v:
		return true
	default:
		return false
	}
}

func TryReceive(ch <-chan int) (int, bool) {
	select {
	case v, ok := <-ch:
		return v, ok
	default:
		return 0, false
	}
}

func RunWithTimeout(work func() int, d time.Duration) (int, bool) {
	ch := make(chan int, 1) // buffered: the goroutine can always finish
	go func() { ch <- work() }()
	select {
	case v := <-ch:
		return v, true
	case <-time.After(d):
		return 0, false
	}
}

func CountTicks(interval, total time.Duration) int {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	timeout := time.After(total)
	count := 0
	for {
		select {
		case <-ticker.C:
			count++
		case <-timeout:
			return count
		}
	}
}

func Counter(done <-chan struct{}) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; ; i++ {
			select {
			case out <- i:
			case <-done:
				return
			}
		}
	}()
	return out
}
