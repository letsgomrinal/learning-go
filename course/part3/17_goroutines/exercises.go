// Lesson 17 exercises: goroutines and WaitGroup. Read README.md first!
//
// Check one exercise at a time:   go test -race -run Ex1 -v
package lesson17

// Ex1: run a and b AT THE SAME TIME, and return only when both have finished.
//
// Hint (README section 4):
//
//	var wg sync.WaitGroup
//	wg.Add(2)
//	go func() { defer wg.Done(); a() }()
//	... the same for b ...
//	wg.Wait()
//
// You'll need to import "sync".
func RunBoth(a, b func()) {
	// TODO (right now they run one after the other, and the test sees it's too slow)
	a()
	b()
}

// Ex2: run every task at the same time, and return when all are done.
// Use wg.Add(1) before each go statement.
func RunAll(tasks []func()) {
	// TODO
	for _, task := range tasks {
		task()
	}
}

// Ex3: return a slice with the square of each number. Compute each square in its own goroutine.
// Each goroutine writes to its own index of the results slice (README section 5).
// SquareAll([]int{1, 2, 3}) == [1 4 9]
func SquareAll(nums []int) []int {
	// TODO
	return nil
}

// Ex4: "download" every url by calling fetch(url), all at the same time.
// Return the results in the same order as urls.
//
// fetch is slow (the test makes it sleep), so doing them one by one is too slow.
func FetchAll(urls []string, fetch func(url string) string) []string {
	// TODO
	return nil
}

// Ex5: add up nums using `parts` goroutines. Each goroutine sums one chunk of the slice
// into its own slot of a `partial` slice. After wg.Wait(), add up the partial sums.
//
// ParallelSum([]int{1, 2, 3, 4, 5, 6}, 3) == 21
// (chunks: [1 2], [3 4], [5 6]; partial sums: 3, 7, 11)
//
// Hints:
//   - if parts < 1, use 1
//   - chunk size: size := (len(nums) + parts - 1) / parts   (rounds up)
//   - chunk i is nums[start:end] with start := i*size and end := min(start+size, len(nums))
//   - skip a chunk if start >= len(nums)
func ParallelSum(nums []int, parts int) int {
	// TODO
	return 0
}
