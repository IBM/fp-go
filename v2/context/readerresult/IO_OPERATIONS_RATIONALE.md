# Why Combining IO Operations with ReaderResult Makes Sense

## Overview

The `context/readerresult` package provides functions that combine IO operations (like `FromIO`, `ChainIOK`, `TapIOK`, etc.) with ReaderResult computations. This document explains why this combination is natural and appropriate, despite IO operations being side-effectful.

## Key Insight: ReaderResult is Already Effectful

**IMPORTANT**: Unlike pure functional Reader monads, `ReaderResult[A]` in this package is **already side-effectful** because it depends on `context.Context`.

### Why context.Context is Effectful

The `context.Context` type in Go is inherently effectful because it:

1. **Can be cancelled**: `ctx.Done()` returns a channel that closes when the context is cancelled
2. **Has deadlines**: `ctx.Deadline()` returns a time when the context expires
3. **Carries values**: `ctx.Value(key)` retrieves request-scoped values
4. **Propagates signals**: Cancellation signals propagate across goroutines
5. **Has observable state**: The context's state can change over time (e.g., when cancelled)

### Type Definition

```go
type ReaderResult[A any] = func(context.Context) Result[A]
```

This is **not** a pure function because:
- The behavior can change based on the context's state
- The context can be cancelled during execution
- The context carries mutable, observable state

## Comparison with Pure Reader Monads

### Pure Reader (from `readerresult` package)

```go
type ReaderResult[R, A any] = func(R) Result[A]
```

- `R` can be any type (config, state, etc.)
- The function is **pure** if `R` is immutable
- No side effects unless explicitly introduced

### Effectful Reader (from `context/readerresult` package)

```go
type ReaderResult[A any] = func(context.Context) Result[A]
```

- Always depends on `context.Context`
- **Inherently effectful** due to context's nature
- Side effects are part of the design

## Why IO Operations Fit Naturally

Since `ReaderResult` is already effectful, combining it with IO operations is a natural fit:

### 1. Both Represent Side Effects

```go
// IO operation - side effectful
io := func() int {
    fmt.Println("Performing IO")
    return 42
}

// ReaderResult - also side effectful (depends on context)
rr := func(ctx context.Context) Result[int] {
    // Can check if context is cancelled (side effect)
    if ctx.Err() != nil {
        return result.Error[int](ctx.Err())
    }
    return result.Of(42)
}

// Combining them is natural
combined := FromIO(io)
```

### 2. Context-Aware IO Operations

The combination allows IO operations to respect context cancellation:

```go
// IO operation that should respect cancellation
readFile := func(path string) ReaderResult[[]byte] {
    return func(ctx context.Context) Result[[]byte] {
        // Check cancellation before expensive IO
        if ctx.Err() != nil {
            return result.Error[[]byte](ctx.Err())
        }
        
        // Perform IO operation
        data, err := os.ReadFile(path)
        if err != nil {
            return result.Error[[]byte](err)
        }
        return result.Of(data)
    }
}
```

### 3. Practical Use Cases

#### Logging with Side Effects

```go
// Log to external system (IO operation)
logMetric := func(value int) func() string {
    return func() string {
        // Side effect: write to metrics system
        metrics.Record("value", value)
        return "logged"
    }
}

// Use with ReaderResult
pipeline := F.Pipe1(
    readerresult.Of(42),
    readerresult.TapIOK(logMetric),
)
```

#### Database Operations

```go
// Database query (IO operation with context)
queryDB := func(id int) ReaderResult[User] {
    return func(ctx context.Context) Result[User] {
        // Context used for timeout/cancellation
        user, err := db.QueryContext(ctx, "SELECT * FROM users WHERE id = ?", id)
        if err != nil {
            return result.Error[User](err)
        }
        return result.Of(user)
    }
}

// Chain with other operations
pipeline := F.Pipe2(
    readerresult.Of(123),
    readerresult.Chain(queryDB),
    readerresult.TapIOK(func(user User) func() string {
        return func() string {
            log.Printf("Retrieved user: %s", user.Name)
            return "logged"
        }
    }),
)
```

#### HTTP Requests

```go
// HTTP request (IO operation)
fetchData := func(url string) ReaderResult[Response] {
    return func(ctx context.Context) Result[Response] {
        req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            return result.Error[Response](err)
        }
        return result.Of(resp)
    }
}
```

## Functions That Combine IO with ReaderResult

### Lifting Functions

- **`FromIO[A]`**: Lifts a pure IO computation into ReaderResult
- **`FromIOResult[A]`**: Lifts an IOResult (IO with error handling) into ReaderResult

### Chaining Functions

- **`ChainIOK[A, B]`**: Sequences a ReaderResult with an IO computation
- **`ChainIOEitherK[A, B]`**: Sequences with an IOResult computation
- **`ChainIOResultK[A, B]`**: Alias for ChainIOEitherK

### Tapping Functions (Side Effects)

- **`TapIOK[A, B]`**: Executes IO for side effects, preserves original value
- **`ChainFirstIOK[A, B]`**: Same as TapIOK
- **`MonadTapIOK[A, B]`**: Monadic version of TapIOK
- **`MonadChainFirstIOK[A, B]`**: Monadic version of ChainFirstIOK

### Error Handling with IO

- **`TapLeftIOK[A, B]`**: Executes IO on error for side effects (logging, metrics)
- **`ChainFirstLeftIOK[A, B]`**: Same as TapLeftIOK

### Reading Context from IO

- **`ReadIO[A]`**: Executes ReaderResult with context from IO
- **`ReadIOEither[A]`**: Executes with context from IOResult
- **`ReadIOResult[A]`**: Alias for ReadIOEither

## Design Philosophy

### Embrace Effectfulness

Rather than trying to maintain purity (which is impossible with `context.Context`), this package embraces the effectful nature of Go's context and provides tools to work with it safely and composably.

### Composition Over Isolation

The package allows you to compose effectful operations (ReaderResult + IO) in a type-safe, functional way, rather than isolating them.

### Practical Go Idioms

This approach aligns with Go's pragmatic philosophy:
- Context is used everywhere in Go for cancellation and timeouts
- IO operations are common and necessary
- Combining them in a type-safe way improves code quality

## Contrast with Pure Functional Packages

### When to Use `context/readerresult` (This Package)

Use when you need:
- ✅ Context cancellation and timeouts
- ✅ Request-scoped values
- ✅ Integration with Go's standard library (http, database/sql, etc.)
- ✅ IO operations with error handling
- ✅ Practical, idiomatic Go code

### When to Use `readerresult` (Pure Package)

Use when you need:
- ✅ Pure dependency injection
- ✅ Testable computations with simple config objects
- ✅ No context propagation
- ✅ Generic environment types (not limited to context.Context)
- ✅ Purely functional composition

## Conclusion

Combining IO operations with ReaderResult in the `context/readerresult` package makes sense because:

1. **ReaderResult is already effectful** due to its dependency on `context.Context`
2. **IO operations are also effectful**, making them a natural fit
3. **The combination provides practical benefits** for real-world Go applications
4. **It aligns with Go's pragmatic philosophy** of embracing side effects when necessary
5. **It enables type-safe composition** of effectful operations

The key insight is that `context.Context` itself is a side effect, so adding more side effects (IO operations) doesn't violate any purity constraints—because there were none to begin with. This package provides tools to work with these side effects in a safe, composable, and type-safe manner.