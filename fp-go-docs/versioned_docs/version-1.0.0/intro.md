---
sidebar_position: 1
---

# What is fp-go?

fp-go is a comprehensive functional programming library for Go, strongly influenced by the excellent [fp-ts](https://github.com/gcanti/fp-ts) library for TypeScript. It provides a set of data types and functions that make it easy and fun to write maintainable and testable code in Go.

The library encourages writing many small, testable, and pure functions that produce output only depending on their input and execute no side effects. It also helps isolate side effects into lazily executed functions and provides a consistent set of composition functions across all data types.

## Quick Example

Here's a simple example showing how to handle errors functionally using the `Either` type:

```go
import (
    "errors"
    "github.com/IBM/fp-go/either"
    "github.com/IBM/fp-go/function"
)

// Pure function that can fail
func divide(a, b int) either.Either[error, int] {
    if b == 0 {
        return either.Left[int](errors.New("division by zero"))
    }
    return either.Right[error](a / b)
}

// Compose operations safely
result := function.Pipe2(
    divide(10, 2),
    either.Map(func(x int) int { return x * 2 }),
    either.GetOrElse(func() int { return 0 }),
)
// result = 10
```

This example demonstrates how fp-go allows you to:
- Express operations that can fail using the `Either` type
- Chain operations together safely
- Handle errors explicitly without nested if statements
- Write pure, composable functions