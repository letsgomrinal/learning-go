// Reference solution for lesson 07.
package concurrency

import (
	"sync"
	"time"
)

func ParallelSum(nums []int, workers int) int {
	if workers <= 0 {
		workers = 1
	}
	workers = min(workers, max(len(nums), 1))
	chunk := (len(nums) + workers - 1) / workers

	results := make(chan int, workers) // buffered: workers never block on send
	launched := 0
	for start := 0; start < len(nums); start += chunk {
		part := nums[start:min(start+chunk, len(nums))]
		launched++
		go func() {
			s := 0
			for _, n := range part {
				s += n
			}
			results <- s
		}()
	}

	total := 0
	for range launched {
		total += <-results
	}
	return total
}

func Generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func Square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func Merge(chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, ch := range chans {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}()
	}
	// Close only after every forwarder is done; closing earlier would make
	// a forwarder panic with "send on closed channel".
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func ProcessOrdered(inputs []string, workers int, fn func(string) string) []string {
	results := make([]string, len(inputs))
	jobs := make(chan int)

	var wg sync.WaitGroup
	for range max(workers, 1) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				results[i] = fn(inputs[i]) // distinct index per job: no race
			}
		}()
	}

	for i := range inputs {
		jobs <- i
	}
	close(jobs) // lets the workers' range loops finish
	wg.Wait()
	return results
}

func WithTimeout(fn func() int, d time.Duration) (int, bool) {
	ch := make(chan int, 1) // buffered so the goroutine can exit even if nobody reads
	go func() { ch <- fn() }()

	select {
	case v := <-ch:
		return v, true
	case <-time.After(d):
		return 0, false
	}
}
