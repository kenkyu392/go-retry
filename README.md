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
	"time"

	"github.com/kenkyu392/go-retry"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Or use Do(time.Second...
	// Retry the function that encountered the error, not all functions.
	errs := retry.DoWithContext(ctx, time.Second,
		func(ctx context.Context) error {
			// Execute the first step...
			return nil
		},
		func(ctx context.Context) error {
			// Execute the second step...
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
