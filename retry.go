package retry

import (
	"context"
	"errors"
	"time"
)

var (
	// Canceled is used to cancel a retry.
	Canceled = errors.New("canceled")
	// Skipped is used to skip a function.
	Skipped = errors.New("skipped")
)

// Do wraps DoWithContext using the background context.
func Do(delay time.Duration, fns ...func(context.Context) error) []error {
	return DoWithContext(context.Background(), delay, fns...)
}

// DoWithContext executes the given functions in order.
// If the function returns an error, the function will be executed again after
// the time specified by delay.
func DoWithContext(ctx context.Context, delay time.Duration, fns ...func(context.Context) error) []error {
	var errs = make([]error, 0)
	for i := 0; i < len(fns); {
		if err := fns[i](ctx); err != nil {
			switch {
			default:
				errs = append(errs, err)
				select {
				case <-ctx.Done():
					errs = append(errs, ctx.Err())
					return errs
				case <-time.After(delay):
					continue
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
			i++
		}
	}
	return errs
}
