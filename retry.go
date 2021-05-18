package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

var (
	// Canceled is used to cancel a retry.
	Canceled = errors.New("canceled")
	// Skipped is used to skip a function.
	Skipped = errors.New("skipped")
)

// DurationFunc ...
type DurationFunc func(retries int) time.Duration

// Duration creates a DurationFunc that returns a Duration.
func Duration(d time.Duration) DurationFunc {
	return func(retries int) time.Duration {
		return d
	}
}

// ExponentialBackoff creates and returns a DurationFunc that exponentially
// backoff the retry interval.
// If maxRetries is negative, retry without limit.
func ExponentialBackoff(maxRetries int) DurationFunc {
	return func(retries int) time.Duration {
		if maxRetries >= 0 && maxRetries < retries {
			return -1
		}
		return time.Duration(math.Pow(2, float64(retries))*100) * time.Millisecond
	}
}

// Do wraps DoWithContext using the background context.
func Do(delayFn DurationFunc, fns ...func(context.Context) error) []error {
	return DoWithContext(context.Background(), delayFn, fns...)
}

// DoWithContext executes the given functions in order.
// If the function returns an error, the function will be executed again after
// the time specified by delay.
// If the result of delayFn is negative, execute the next function without retrying.
func DoWithContext(ctx context.Context, delayFn DurationFunc, fns ...func(context.Context) error) []error {
	var errs = make([]error, 0)
	for i, retries := 0, 0; i < len(fns); {
		if err := fns[i](ctx); err != nil {
			switch {
			default:
				errs = append(errs, err)
				retries++
				if delay := delayFn(retries); delay >= 0 {
					select {
					case <-ctx.Done():
						errs = append(errs, ctx.Err())
						return errs
					case <-time.After(delayFn(retries)):
						continue
					}
				}
			case errors.Is(err, Canceled):
				// If Canceled is received, the function will exit.
				return errs
			case errors.Is(err, Skipped):
				// If Skipped is received, proceed to the next function.
			}
		}
		select {
		case <-ctx.Done():
			errs = append(errs, ctx.Err())
			return errs
		default:
			retries = 0
			i++
		}
	}
	return errs
}
