# go-retry

[![test](https://github.com/kenkyu392/go-retry/workflows/test/badge.svg?branch=master)](https://github.com/kenkyu392/go-retry)
[![codecov](https://codecov.io/gh/kenkyu392/go-retry/branch/master/graph/badge.svg)](https://codecov.io/gh/kenkyu392/go-retry)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-00ADD8?logo=go)](https://pkg.go.dev/github.com/kenkyu392/go-retry)
[![go report card](https://goreportcard.com/badge/github.com/kenkyu392/go-retry)](https://goreportcard.com/report/github.com/kenkyu392/go-retry)
[![license](https://img.shields.io/github/license/kenkyu392/go-retry)](LICENSE)

Provides a retry function that allows for step-by-step execution.

## Installation

```
go get -u github.com/kenkyu392/go-retry
```

## Usage

```go
package main

import (
	"context"
	"log"

	"github.com/kenkyu392/go-retry"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Or you can use Do(retry.Duration(time.Second)...) to specify the delay before retrying.
	// Retry only the function in error, not all functions.
	errs := retry.DoWithContext(ctx,
		// Use exponential backoff with a maximum of 5 retries for the
		// delay time specification function.
		retry.ExponentialBackoff(5),
		// Execute the first step...
		func(ctx context.Context) error {
			// You can use Skipped or Canceled in a function to skip the current
			// function or cancel all remaining functions.
			return nil
		},
		// Execute the second step...
		func(ctx context.Context) error {
			return nil
		},
	)
	// All errors encountered during execution are returned in an array.
	for _, err := range errs {
		log.Print(err)
	}
}
```

## License

[MIT](LICENSE)
