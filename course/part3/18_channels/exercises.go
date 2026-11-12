// Lesson 18 exercises: channels. Read README.md first!
//
// Check one exercise at a time:   go test -race -run Ex1 -v
package lesson18

// Ex1: return a channel, and start a goroutine that sends "hello, " + name on it.
//
//	ch := Greeting("Ana")
//	<-ch == "hello, Ana"
//
// Hint:
//
//	ch := make(chan string)
//	go func() { ch <- ... }()
//	return ch
func Greeting(name string) <-chan string {
	// TODO
	return nil
}

// Ex2: return a channel that produces 1, 2, ..., n and is then CLOSED.
// The sending happens in a goroutine started inside Numbers. (README section 5.)
// For n <= 0, just close the channel without sending anything.
func Numbers(n int) <-chan int {
	// TODO
	return nil
}

// Ex3: receive every value from ch until it's closed, and return their sum.
// Hint: for v := range ch { ... }
func Sum(ch <-chan int) int {
	// TODO
	return 0
}

// Ex4: start one goroutine per number. Each sends n*n on a results channel.
// Receive exactly len(nums) values and return their sum. (README section 7.)
// SumOfSquares([]int{1, 2, 3}) == 14
func SumOfSquares(nums []int) int {
	// TODO
	return 0
}

// Ex5: a "stage": start a goroutine that receives every value from in, sends value*2 on a new
// channel, and closes that channel when in is closed. Return the new channel.
//
//	for v := range Doubler(Numbers(3)) { ... }   // 2, 4, 6
func Doubler(in <-chan int) <-chan int {
	// TODO
	return nil
}

// Ex6: return a channel that already holds all the items, in order, and is closed,
// WITHOUT starting any goroutine.
//
// Hint: a buffered channel with room for all the items: make(chan string, len(items)).
// Sends won't block because there's room. Then close it and return it.
func Queue(items []string) <-chan string {
	// TODO
	return nil
}
