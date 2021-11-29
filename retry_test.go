package retry

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	rand.Seed(1) // 1 7 7 9 1 8 5 0 6 0
	errs := Do(Duration(time.Millisecond*500),
		func(ctx context.Context) error {
			if n := rand.Intn(10); n%2 != 0 {
				return fmt.Errorf("step 1: %d is not an even number", n)
			}
			return nil
		},
		func(ctx context.Context) error {
			if n := rand.Intn(10); n%2 != 0 {
				return fmt.Errorf("step 2: %d is not an even number", n)
			}
			return nil
		},
	)
	for _, err := range errs {
		t.Logf("%#v", err)
	}
	if len(errs) != 6 {
		t.Fatal("6 errors need to occur")
	}
}

func TestDoWithContext(t *testing.T) {
	t.Run("case=cancel-1", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		errs := DoWithContext(ctx, Duration(time.Second),
			func(ctx context.Context) error {
				cancel()
				return nil
			},
			func(ctx context.Context) error {
				return errors.New("cancel")
			},
		)
		for _, err := range errs {
			t.Logf("%#v", err)
		}
		if len(errs) != 1 {
			t.Fatal("1 error need to occur")
		}
	})

	t.Run("case=cancel-2", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		errs := DoWithContext(ctx, Duration(time.Second),
			func(ctx context.Context) error {
				return nil
			},
			func(ctx context.Context) error {
				cancel()
				return errors.New("cancel")
			},
		)
		for _, err := range errs {
			t.Logf("%#v", err)
		}
		if len(errs) != 2 {
			t.Fatal("2 errors need to occur")
		}
	})
}

func TestErr(t *testing.T) {
	rand.Seed(1) // 1 7 7 9 1 8 5 0 6 0
	errs := Do(Duration(time.Millisecond*500),
		func(ctx context.Context) error {
			if n := rand.Intn(10); n%2 != 0 {
				return Skipped
			}
			return nil
		},
		func(ctx context.Context) error {
			return Canceled
		},
		func(ctx context.Context) error {
			if n := rand.Intn(10); n%2 != 0 {
				return fmt.Errorf("step 2: %d is not an even number", n)
			}
			return nil
		},
	)
	for _, err := range errs {
		t.Logf("%#v", err)
	}
	if got := rand.Intn(10); got != 7 {
		t.Errorf("got: %v, want: 7", got)
	}
	if len(errs) != 0 {
		t.Fatal("no error should occur")
	}
}

func TestExponentialBackoff(t *testing.T) {
	// 2ms + 4ms + 8ms + 1.6s = min:3.0s max:3.4s
	const (
		maxRetries = 4
		bufferTime = time.Millisecond * 100
		minTime    = time.Second * 3
		maxTime    = minTime + (time.Millisecond * 400) + bufferTime
	)

	t.Run("case=limited", func(t *testing.T) {
		start := time.Now()
		errs := Do(ExponentialBackoff(maxRetries),
			func(ctx context.Context) error {
				return errors.New("retry")
			},
		)
		since := time.Since(start)
		if got := len(errs); got != maxRetries+1 {
			t.Fatalf("got: %v, want: %d", got, maxRetries+1)
		}
		if since < minTime || since > maxTime {
			t.Fatalf("delay is too short or too long (%v)", since)
		}
	})

	t.Run("case=infinity", func(t *testing.T) {
		var retries = 0
		start := time.Now()
		errs := Do(ExponentialBackoff(-1),
			func(ctx context.Context) error {
				if retries++; retries > maxRetries {
					return nil
				}
				return errors.New("retry")
			},
		)
		since := time.Since(start)
		if got := len(errs); got != maxRetries {
			t.Fatalf("got: %v, want: %d", got, maxRetries)
		}
		if since < minTime || since > maxTime {
			t.Fatalf("delay is too short or too long (%v)", since)
		}
	})
}
