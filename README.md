# goPool
goPool is a lightweight goroutine pool manager that allows you to limit the number of concurrently running goroutines.

## Features

- Limit number of goroutines with a single line
- Familiar API: `Add()`, `Done()`, and `Wait()`
- Context-aware cancellation via `AddWithContext()`
- No external dependencies

## Installation

```bash
go get github.com/SandeepMahapatra/goPool
```

## Examples

### Example 1: Basic Usage

```go
package main

import (
	"fmt"
	"time"

	gopool "github.com/SandeepMahapatra/goPool"
)

func main() {
	pool := gopool.New(3) // allow only 3 goroutines at a time

	for i := 0; i < 10; i++ {
		pool.Add()
		go func(i int) {
			defer pool.Done()
			fmt.Printf("Working on task %d\n", i)
			time.Sleep(2 * time.Second)
		}(i)
	}

	pool.Wait()
	fmt.Println("All tasks completed")
}
```

### Example 2: Using `AddWithContext`
```go
package main

import (
    "context"
    "fmt"
    "time"

	gopool "github.com/SandeepMahapatra/goPool"
)

func main() {
pool := gopool.New(2)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	for i := 0; i < 5; i++ {
		if err := pool.AddWithContext(ctx); err != nil {
			fmt.Printf("Skipped task %d due to timeout\n", i)
			continue
		}

		go func(i int) {
			defer pool.Done()
			fmt.Printf("Processing %d\n", i)
			time.Sleep(2 * time.Second)
		}(i)
	}

	pool.Wait()
}
```

