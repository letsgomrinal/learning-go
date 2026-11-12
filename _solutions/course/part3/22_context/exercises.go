// Reference solution for lesson 22.
package lesson22

import (
	"context"
	"time"
)

func Sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func SumUntilCancelled(ctx context.Context, nums []int, step func(int) int) (int, error) {
	sum := 0
	for _, n := range nums {
		if err := ctx.Err(); err != nil {
			return sum, err
		}
		sum += step(n)
	}
	return sum, nil
}

func Ticker(ctx context.Context, every time.Duration) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		ticker := time.NewTicker(every)
		defer ticker.Stop()
		for i := 1; ; i++ {
			select {
			case <-ticker.C:
				select {
				case out <- i:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func WithTimeout(d time.Duration, work func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	return work(ctx)
}

func FirstResult(ctx context.Context, fns []func(ctx context.Context) string) string {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan string, len(fns))
	for _, fn := range fns {
		go func() {
			results <- fn(ctx)
		}()
	}
	select {
	case r := <-results:
		return r
	case <-ctx.Done():
		return ""
	}
}

func FetchAll(ctx context.Context, ids []int, fetch func(ctx context.Context, id int) (string, error)) ([]string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]string, len(ids))
	errs := make(chan error, len(ids))
	for i, id := range ids {
		go func() {
			res, err := fetch(ctx, id)
			if err == nil {
				results[i] = res
			}
			errs <- err
		}()
	}
	for range ids {
		if err := <-errs; err != nil {
			cancel()
			return nil, err
		}
	}
	return results, nil
}
