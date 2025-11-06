# fp-go: Functional Programming Library for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/IBM/fp-go.svg)](https://pkg.go.dev/github.com/IBM/fp-go)
[![Coverage Status](https://coveralls.io/repos/github/IBM/fp-go/badge.svg?branch=main)](https://coveralls.io/github/IBM/fp-go?branch=main)

**🚧 Work in progress! 🚧** Despite major version 1 (due to [semantic-release limitations](https://github.com/semantic-release/semantic-release/issues/1507)), we're working to minimize breaking changes.

![logo](resources/images/logo.png)

A comprehensive functional programming library for Go, strongly influenced by the excellent [fp-ts](https://github.com/gcanti/fp-ts) library for TypeScript.

## 📚 Table of Contents

- [Getting Started](#-getting-started)
- [Design Goals](#-design-goals)
- [Core Concepts](#-core-concepts)
- [Comparison to Idiomatic Go](#comparison-to-idiomatic-go)
- [Implementation Notes](#implementation-notes)
- [Common Operations](#common-operations)
- [Resources](#-resources)

## 🚀 Getting Started

### Installation

```bash
go get github.com/IBM/fp-go
```

### Quick Example

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

### Resources

- 📖 [API Documentation](https://pkg.go.dev/github.com/IBM/fp-go)
- 💡 [Code Samples](./samples/)
- 🆕 [V2 Documentation](./v2/README.md) (requires Go 1.24+)

## 🎯 Design Goals

This library aims to provide a set of data types and functions that make it easy and fun to write maintainable and testable code in Go. It encourages the following patterns:

### Core Principles

- **Pure Functions**: Write many small, testable, and pure functions that produce output only depending on their input and execute no side effects
- **Side Effect Isolation**: Isolate side effects into lazily executed functions using the `IO` monad
- **Consistent Composition**: Expose a consistent set of composition functions across all data types
  - Each data type has a small set of composition functions
  - Functions are named consistently across all data types
  - Semantics of same-named functions are consistent across data types

### 🧘🏽 Alignment with the Zen of Go

This library respects and aligns with [The Zen of Go](https://the-zen-of-go.netlify.app/):

| Principle | Alignment | Explanation |
|-----------|-----------|-------------|
| 🧘🏽 Each package fulfills a single purpose | ✔️ | Each top-level package (Option, Either, ReaderIOEither, etc.) defines one data type and its operations |
| 🧘🏽 Handle errors explicitly | ✔️ | Clear distinction between operations that can/cannot fail; failures represented via `Either` type |
| 🧘🏽 Return early rather than nesting deeply | ✔️ | Small, focused functions composed together; `Either` monad handles error paths automatically |
| 🧘🏽 Leave concurrency to the caller | ✔️ | Pure functions are synchronous; I/O operations are asynchronous by default |
| 🧘🏽 Before you launch a goroutine, know when it will stop | 🤷🏽 | Library doesn't start goroutines; Task monad supports cancellation via context |
| 🧘🏽 Avoid package level state | ✔️ | No package-level state anywhere |
| 🧘🏽 Simplicity matters | ✔️ | Small, consistent interface across data types; focus on business logic |
| 🧘🏽 Write tests to lock in behaviour | 🟡 | Programming pattern encourages testing; library has growing test coverage |
| 🧘🏽 If you think it's slow, first prove it with a benchmark | ✔️ | Performance claims should be backed by benchmarks |
| 🧘🏽 Moderation is a virtue | ✔️ | No custom goroutines or expensive synchronization; atomic counters for coordination |
| 🧘🏽 Maintainability counts | ✔️ | Small, concise operations; pure functions with clear type signatures |

## 💡 Core Concepts

### Data Types

The library provides several key functional data types:

- **`Option[A]`**: Represents an optional value (Some or None)
- **`Either[E, A]`**: Represents a value that can be one of two types (Left for errors, Right for success)
- **`IO[A]`**: Represents a lazy computation that produces a value
- **`IOEither[E, A]`**: Represents a lazy computation that can fail
- **`Reader[R, A]`**: Represents a computation that depends on an environment
- **`ReaderIOEither[R, E, A]`**: Combines Reader, IO, and Either for effectful computations with dependencies
- **`Task[A]`**: Represents an asynchronous computation
- **`State[S, A]`**: Represents a stateful computation

### Monadic Operations

All data types support common monadic operations:

- **`Map`**: Transform the value inside a context
- **`Chain`** (FlatMap): Transform and flatten nested contexts
- **`Ap`**: Apply a function in a context to a value in a context
- **`Of`**: Wrap a value in a context
- **`Fold`**: Extract a value from a context

## Comparison to Idiomatic Go

This section explains how functional APIs differ from idiomatic Go and how to convert between them.

### Pure Functions

Pure functions take input parameters and compute output without changing global state or mutating inputs. They always return the same output for the same input.

#### Without Errors

If your pure function doesn't return an error, the idiomatic signature works as-is:

```go
func add(a, b int) int {
    return a + b
}
```

#### With Errors

**Idiomatic Go:**
```go
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

**Functional Style:**
```go
func divide(a, b int) either.Either[error, int] {
    if b == 0 {
        return either.Left[int](errors.New("division by zero"))
    }
    return either.Right[error](a / b)
}
```

**Conversion:**
- Use `either.EitherizeXXX` to convert from idiomatic to functional style
- Use `either.UneitherizeXXX` to convert from functional to idiomatic style

### Effectful Functions

An effectful function changes data outside its scope or doesn't always produce the same output for the same input.

#### Without Errors

**Functional signature:** `IO[T]`

```go
func getCurrentTime() io.IO[time.Time] {
    return func() time.Time {
        return time.Now()
    }
}
```

#### With Errors

**Functional signature:** `IOEither[error, T]`

```go
func readFile(path string) ioeither.IOEither[error, []byte] {
    return func() either.Either[error, []byte] {
        data, err := os.ReadFile(path)
        if err != nil {
            return either.Left[[]byte](err)
        }
        return either.Right[error](data)
    }
}
```

**Conversion:**
- Use `ioeither.EitherizeXXX` to convert idiomatic Go functions to functional style

### Go Context

Functions that take a `context.Context` are effectful because they depend on mutable context.

**Idiomatic Go:**
```go
func fetchData(ctx context.Context, url string) ([]byte, error) {
    // implementation
}
```

**Functional Style:**
```go
func fetchData(url string) readerioeither.ReaderIOEither[context.Context, error, []byte] {
    return func(ctx context.Context) ioeither.IOEither[error, []byte] {
        return func() either.Either[error, []byte] {
            // implementation
        }
    }
}
```

**Conversion:**
- Use `readerioeither.EitherizeXXX` to convert idiomatic Go functions with context to functional style

## Implementation Notes

### Generics

All monadic operations use Go generics for type safety:

- ✅ **Pros**: Type-safe composition, IDE support, compile-time correctness
- ⚠️ **Cons**: May result in larger binaries (different versions per type)
- 💡 **Tip**: For binary size concerns, use type erasure with `any` type

### Ordering of Generic Type Parameters

Go requires all type parameters on the global function definition. Parameters that cannot be auto-detected come first:

```go
// Map: B cannot be auto-detected, so it comes first
func Map[R, E, A, B any](f func(A) B) func(ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B]

// Ap: B cannot be auto-detected from the argument
func Ap[B, R, E, A any](fa ReaderIOEither[R, E, A]) func(ReaderIOEither[R, E, func(A) B]) ReaderIOEither[R, E, B]
```

This ordering maximizes type inference where possible.

### Use of the ~ Operator

Go doesn't support generic type aliases (until Go 1.24), only type definitions. The `~` operator allows generic implementations to work with compatible types:

```go
type ReaderIOEither[R, E, A any] RD.Reader[R, IOE.IOEither[E, A]]
```

**Generic Subpackages:**
- Each higher-level type has a `generic` subpackage with fully generic implementations
- These are for library extensions, not end-users
- Main packages specialize generic implementations for convenience

### Higher Kinded Types (HKT)

Go doesn't support HKT natively. This library addresses this by:

- Introducing HKTs as individual types (e.g., `HKTA` for `HKT[A]`)
- Implementing generic algorithms in the `internal` package
- Keeping complexity hidden from end-users

## Common Operations

### Map/Chain/Ap/Flap

| Operator | Parameter        | Monad           | Result   | Use Case |
| -------- | ---------------- | --------------- | -------- | -------- |
| Map      | `func(A) B`      | `HKT[A]`        | `HKT[B]` | Transform value in context |
| Chain    | `func(A) HKT[B]` | `HKT[A]`        | `HKT[B]` | Transform and flatten |
| Ap       | `HKT[A]`         | `HKT[func(A)B]` | `HKT[B]` | Apply function in context |
| Flap     | `A`              | `HKT[func(A)B]` | `HKT[B]` | Apply value to function in context |

### Example: Chaining Operations

```go
import (
    "github.com/IBM/fp-go/either"
    "github.com/IBM/fp-go/function"
)

result := function.Pipe3(
    either.Right[error](10),
    either.Map(func(x int) int { return x * 2 }),
    either.Chain(func(x int) either.Either[error, int] {
        if x > 15 {
            return either.Right[error](x)
        }
        return either.Left[int](errors.New("too small"))
    }),
    either.GetOrElse(func() int { return 0 }),
)
```

## 📚 Resources

- [API Documentation](https://pkg.go.dev/github.com/IBM/fp-go)
- [Code Samples](./samples/)
- [V2 Documentation](./v2/README.md) - New features in Go 1.24+
- [fp-ts](https://github.com/gcanti/fp-ts) - Original TypeScript inspiration

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## 📄 License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.
