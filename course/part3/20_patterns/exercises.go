// Lesson 20 exercises: pipelines and worker pools. Read README.md first!
//
// Check one exercise at a time:   go test -race -run Ex1 -v
package lesson20

// Ex1: the first stage. Return a channel that sends each number in order, then closes.
// Do the sending in a goroutine. (README section 1, generate.)
func Generate(nums ...int) <-chan int {
	// TODO
	return nil
}

// Ex2: a middle stage. Receive every value from in, send its square, and close the output
// when in is closed.
func Square(in <-chan int) <-chan int {
	// TODO
	return nil
}

// Ex3: connect the pipeline: Generate(nums...) → Square → add everything up.
// SumSquares(1, 2, 3) == 14
func SumSquares(nums ...int) int {
	// TODO
	return 0
}

// Ex4: fan-in. Forward every value from every input channel to one output channel.
// Close the output once ALL inputs are closed and all their values were forwarded.
// The order of values doesn't matter. (README section 2.)
//
// You'll need "sync" for the WaitGroup.
func Merge(chans ...<-chan int) <-chan int {
	// TODO
	return nil
}

// Ex5: run work(job) for every job using exactly `workers` goroutines, and return the results
// in the SAME ORDER as jobs. (README section 3.)
//
//	WorkerPool([]int{1, 2, 3}, 2, double) == [2 4 6]
//
// If workers < 1, use 1.
func WorkerPool(jobs []int, workers int, work func(int) int) []int {
	// TODO
	return nil
}

// Ex6: run every task, starting them all as goroutines, but with at most `limit` running at the
// same time. Return when all tasks are done. (README section 4, semaphore.)
// If limit < 1, use 1.
func RunLimited(tasks []func(), limit int) {
	// TODO
	for _, task := range tasks {
		task()
	}
}
