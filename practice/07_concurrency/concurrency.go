// Lesson 07: Goroutines and channels. Go's superpower.
//
// Run:  go test -race -v ./practice/07_concurrency      ← ALWAYS use -race from now on
//
// ─── CONCEPTS ──────────────────────────────────────────────────────────────
//
// 1. GOROUTINES are cheap concurrent functions (~2KB stack, millions are fine):
//      go doWork(x)
//    Unlike Python threads there's NO GIL: goroutines run truly in PARALLEL on all cores.
//    Unlike asyncio there's NO async/await "coloring": any function can block, and the
//    runtime scheduler parks the goroutine and runs another one.
//    main() returning kills all goroutines. Nobody waits for them automatically.
//
// 2. CHANNELS are typed pipes between goroutines:
//      ch := make(chan int)      // UNBUFFERED: send blocks until someone receives
//      ch := make(chan int, 10)  // BUFFERED: send blocks only when the buffer is full
//      ch <- 5                   // send
//      v := <-ch                 // receive
//      v, ok := <-ch             // ok == false → channel closed and drained
//      close(ch)                 // "no more values". Only the SENDER closes.
//      for v := range ch {}      // receive until closed
//    Directional types document intent:  chan<- int (send-only), <-chan int (receive-only)
//
// 3. sync.WaitGroup waits for a group of goroutines:
//      var wg sync.WaitGroup
//      for _, x := range items {
//          wg.Add(1)
//          go func() { defer wg.Done(); work(x) }()
//      }
//      wg.Wait()
//    (Go 1.25+: wg.Go(func(){ work(x) }) does the Add and Done for you.)
//    Since Go 1.22 each loop iteration has its own `x`, so capturing it in a closure is safe.
//
// 4. select waits on several channel operations at once:
//      select {
//      case v := <-results:  ...
//      case <-time.After(time.Second): ... // timeout
//      case <-ctx.Done(): ...              // cancellation (lesson 08)
//      }
//
// 5. PATTERNS you'll implement below:
//    - Generator: a function returns <-chan T and fills it from its own goroutine
//    - Pipeline: stage → stage → stage, each a goroutine reading one channel and writing another
//    - Fan-out / fan-in: many workers read one channel; merge many channels into one
//    - Worker pool: a fixed number of workers processes a job queue
//
// 6. DEADLOCK AND LEAKS
//    "fatal error: all goroutines are asleep - deadlock!" means everyone is waiting.
//    A goroutine blocked forever on a channel nobody reads is a LEAK. Always ask:
//    "how does this goroutine end?"
// ───────────────────────────────────────────────────────────────────────────

package concurrency

import (
	"time"
)

// ParallelSum sums nums by splitting it into `workers` parts, summing each part in
// its own goroutine, and combining the partial sums.
//
// Two valid approaches, so try both:
//
//	a) each goroutine sends its partial sum on a channel; main receives `workers` values
//	b) each goroutine writes partial[i] in a []int (no lock needed: different indices),
//	   then a WaitGroup waits for all of them
//
// Handle workers > len(nums) and workers <= 0 (treat as 1).
func ParallelSum(nums []int, workers int) int {
	// TODO
	return 0
}

// Generate returns a channel that emits nums in order, then is CLOSED.
// The sending happens in a goroutine started inside Generate.
//
//	out := make(chan int)
//	go func() { defer close(out); for ... { out <- n } }()
//	return out
func Generate(nums ...int) <-chan int {
	// TODO
	return nil
}

// Square is a pipeline stage: it reads ints from in, emits their squares, and closes
// its output when in is closed.
//
//	for v := range Square(Generate(1, 2, 3)) { ... }  // 1, 4, 9
func Square(in <-chan int) <-chan int {
	// TODO
	return nil
}

// Merge fans in: it forwards every value from all input channels into ONE output
// channel, and closes the output once ALL inputs are closed.
//
// Hint: start one goroutine per input that forwards values, a WaitGroup to track them,
// and one more goroutine that does wg.Wait(); close(out).
// Why can't you just close(out) after starting the forwarders?
func Merge(chans ...<-chan int) <-chan int {
	// TODO
	return nil
}

// ProcessOrdered runs fn on every input using exactly `workers` goroutines (a worker
// pool) and returns the results IN THE SAME ORDER as the inputs.
//
// Hint: send job indices on a `jobs` channel; workers read an index i and write
// results[i] = fn(inputs[i]). Each worker writes a different index, so no lock is
// needed. Close `jobs` when done sending, and WaitGroup the workers.
func ProcessOrdered(inputs []string, workers int, fn func(string) string) []string {
	// TODO
	return nil
}

// WithTimeout runs fn in a goroutine and returns (result, true) if it finishes
// within d, or (0, false) if the timeout fires first.
//
// Hint: select on a result channel and time.After(d).
// LEAK TRAP: if the result channel is UNBUFFERED and the timeout wins, the goroutine
// blocks forever trying to send. Use make(chan int, 1) so it can always finish.
func WithTimeout(fn func() int, d time.Duration) (int, bool) {
	// TODO
	return 0, false
}
