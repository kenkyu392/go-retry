package retry

import (
	"context"
	"time"
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
			errs = append(errs, err)
			select {
			case <-ctx.Done():
				errs = append(errs, ctx.Err())
				return errs
			case <-time.After(delay):
				continue
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
