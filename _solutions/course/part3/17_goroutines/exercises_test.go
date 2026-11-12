package lesson17

import (
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// Always run these with:  go test -race -v

func TestEx1RunBoth(t *testing.T) {
	var ran atomic.Int32
	slow := func() { time.Sleep(60 * time.Millisecond); ran.Add(1) }

	start := time.Now()
	RunBoth(slow, slow)
	elapsed := time.Since(start)

	if ran.Load() != 2 {
		t.Fatalf("Ex1: %d of the 2 functions had finished when RunBoth returned. Did you wg.Wait()?", ran.Load())
	}
	if elapsed > 110*time.Millisecond {
		t.Errorf("Ex1: took %v. Two 60ms functions at the same time should take about 60ms, not 120ms", elapsed.Round(time.Millisecond))
	}
}

func TestEx2RunAll(t *testing.T) {
	var ran atomic.Int32
	tasks := make([]func(), 10)
	for i := range tasks {
		tasks[i] = func() { time.Sleep(40 * time.Millisecond); ran.Add(1) }
	}

	start := time.Now()
	RunAll(tasks)
	elapsed := time.Since(start)

	if ran.Load() != 10 {
		t.Fatalf("Ex2: only %d of 10 tasks had finished when RunAll returned", ran.Load())
	}
	if elapsed > 150*time.Millisecond {
		t.Errorf("Ex2: took %v. Ten 40ms tasks at the same time should take about 40ms", elapsed.Round(time.Millisecond))
	}
	RunAll(nil) // must not hang or crash
}

func TestEx3SquareAll(t *testing.T) {
	got := SquareAll([]int{1, 2, 3, 4})
	if want := []int{1, 4, 9, 16}; !slices.Equal(got, want) {
		t.Errorf("Ex3: SquareAll([1 2 3 4]) = %v, want %v", got, want)
	}
	if got := SquareAll(nil); len(got) != 0 {
		t.Errorf("Ex3: SquareAll(empty) = %v, want empty", got)
	}
}

func TestEx4FetchAll(t *testing.T) {
	fetch := func(url string) string {
		time.Sleep(50 * time.Millisecond)
		return "<page " + strings.TrimPrefix(url, "https://") + ">"
	}
	urls := []string{"https://a.com", "https://b.com", "https://c.com", "https://d.com", "https://e.com"}

	start := time.Now()
	got := FetchAll(urls, fetch)
	elapsed := time.Since(start)

	want := []string{"<page a.com>", "<page b.com>", "<page c.com>", "<page d.com>", "<page e.com>"}
	if !slices.Equal(got, want) {
		t.Fatalf("Ex4: FetchAll = %q, want %q", got, want)
	}
	if elapsed > 150*time.Millisecond {
		t.Errorf("Ex4: took %v. Five 50ms fetches at the same time should take about 50ms", elapsed.Round(time.Millisecond))
	}
}

func TestEx5ParallelSum(t *testing.T) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i + 1
	}
	for _, parts := range []int{1, 3, 4, 7, 1000, 2000, 0} {
		if got := ParallelSum(nums, parts); got != 500500 {
			t.Errorf("Ex5: ParallelSum(1..1000, parts=%d) = %d, want 500500", parts, got)
		}
	}
	if got := ParallelSum([]int{1, 2, 3, 4, 5, 6}, 3); got != 21 {
		t.Errorf("Ex5: ParallelSum([1..6], 3) = %d, want 21", got)
	}
	if got := ParallelSum(nil, 4); got != 0 {
		t.Errorf("Ex5: ParallelSum(empty, 4) = %d, want 0", got)
	}
}
