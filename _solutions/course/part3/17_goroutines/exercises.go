// Reference solution for lesson 17. Compare with yours AFTER your tests pass.
package lesson17

import "sync"

func RunBoth(a, b func()) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); a() }()
	go func() { defer wg.Done(); b() }()
	wg.Wait()
}

func RunAll(tasks []func()) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			task()
		}()
	}
	wg.Wait()
}

func SquareAll(nums []int) []int {
	results := make([]int, len(nums))
	var wg sync.WaitGroup
	for i, n := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = n * n
		}()
	}
	wg.Wait()
	return results
}

func FetchAll(urls []string, fetch func(url string) string) []string {
	results := make([]string, len(urls))
	var wg sync.WaitGroup
	for i, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = fetch(url)
		}()
	}
	wg.Wait()
	return results
}

func ParallelSum(nums []int, parts int) int {
	if parts < 1 {
		parts = 1
	}
	size := (len(nums) + parts - 1) / parts
	partial := make([]int, parts)
	var wg sync.WaitGroup
	for i := range parts {
		start := i * size
		if start >= len(nums) {
			break
		}
		end := min(start+size, len(nums))
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, n := range nums[start:end] {
				partial[i] += n
			}
		}()
	}
	wg.Wait()
	total := 0
	for _, p := range partial {
		total += p
	}
	return total
}
