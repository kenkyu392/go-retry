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
	i := 0
	start := time.Now()
	times := make([]time.Duration, 0)
	errs := Do(ExponentialBackoff(),
		func(ctx context.Context) error {
			if i < 3 {
				i++
				times = append(times, time.Since(start))
				return fmt.Errorf("retry %d", i)
			}
			return nil
		},
	)
	if got := len(errs); got != 3 {
		t.Fatalf("got: %v, want: 3", got)
	}
	for i := 0; i < len(times); i++ {
		if i+1 < len(times) && times[i]*2 > times[i+1] {
			t.Fatal("delay is too small")
		}
	}
}
