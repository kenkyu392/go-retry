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
	// 1 7 7 9 1 8 5 0 6 0
	errs := Do(time.Millisecond*500,
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
		t.Fatal("5 errors need to occur")
	}
}

func TestDoWithContext(t *testing.T) {
	t.Run("case=cancel-1", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		errs := DoWithContext(ctx, time.Second,
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

		errs := DoWithContext(ctx, time.Second,
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
