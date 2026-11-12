// Reference solution for lesson 20.
package lesson20

import "sync"

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

func SumSquares(nums ...int) int {
	total := 0
	for v := range Square(Generate(nums...)) {
		total += v
	}
	return total
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
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func WorkerPool(jobs []int, workers int, work func(int) int) []int {
	if workers < 1 {
		workers = 1
	}
	results := make([]int, len(jobs))
	indexes := make(chan int)
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range indexes {
				results[i] = work(jobs[i])
			}
		}()
	}
	for i := range jobs {
		indexes <- i
	}
	close(indexes)
	wg.Wait()
	return results
}

func RunLimited(tasks []func(), limit int) {
	if limit < 1 {
		limit = 1
	}
	sem := make(chan struct{}, limit)
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			task()
		}()
	}
	wg.Wait()
}
