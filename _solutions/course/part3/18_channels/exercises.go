// Reference solution for lesson 18.
package lesson18

func Greeting(name string) <-chan string {
	ch := make(chan string)
	go func() {
		ch <- "hello, " + name
	}()
	return ch
}

func Numbers(n int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 1; i <= n; i++ {
			ch <- i
		}
	}()
	return ch
}

func Sum(ch <-chan int) int {
	total := 0
	for v := range ch {
		total += v
	}
	return total
}

func SumOfSquares(nums []int) int {
	results := make(chan int)
	for _, n := range nums {
		go func() {
			results <- n * n
		}()
	}
	total := 0
	for range nums {
		total += <-results
	}
	return total
}

func Doubler(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			out <- v * 2
		}
	}()
	return out
}

func Queue(items []string) <-chan string {
	ch := make(chan string, len(items))
	for _, item := range items {
		ch <- item
	}
	close(ch)
	return ch
}
