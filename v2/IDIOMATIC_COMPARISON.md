# Idiomatic vs Standard Package Comparison

> **Latest Update:** 2025-11-18 - Updated with fresh benchmarks after `either` package optimizations

This document provides a comprehensive comparison between the `idiomatic` packages and the standard fp-go packages (`result` and `option`).

**See also:** [BENCHMARK_COMPARISON.md](./BENCHMARK_COMPARISON.md) for detailed performance analysis.

## Table of Contents

1. [Overview](#overview)
2. [Design Differences](#design-differences)
3. [Performance Comparison](#performance-comparison)
4. [API Comparison](#api-comparison)
5. [When to Use Each](#when-to-use-each)

## Overview

The fp-go library provides two approaches to functional programming patterns in Go:

- **Standard Packages** (`result`, `either`, `option`): Use struct wrappers for algebraic data types
- **Idiomatic Packages** (`idiomatic/result`, `idiomatic/option`): Use native Go tuples for the same patterns

### Key Insight

After recent optimizations to the `either` package, both approaches now offer excellent performance:

- **Simple operations** (~1-5 ns/op): Both packages perform comparably
- **Core transformations**: Idiomatic is **1.2-2.3x faster**
- **Complex operations**: Idiomatic is **2-32x faster** with significantly fewer allocations
- **Real-world pipelines**: Idiomatic shows **2-3.4x speedup**

The idiomatic packages provide:
- Consistently better performance across most operations
- Zero allocations for complex operations (ChainFirst: 72 B → 0 B)
- More familiar Go idioms
- Seamless integration with existing Go code

## Design Differences

### Data Representation

#### Standard Result Package

```go
// Uses Either[error, A] which is a struct wrapper
type Result[A any] = Either[error, A]
type Either[E, A any] struct {
    r      A
    l      E
    isLeft bool
}

// Creating values - ZERO heap allocations (struct returned by value)
success := result.Right[error](42)  // Returns Either struct by value (0 B/op)
failure := result.Left[int](err)    // Returns Either struct by value (0 B/op)

// Benchmarks confirm:
// BenchmarkRight-16    871258489    1.384 ns/op    0 B/op    0 allocs/op
// BenchmarkLeft-16     683089270    1.761 ns/op    0 B/op    0 allocs/op
```

#### Idiomatic Result Package

```go
// Uses native Go tuples (value, error)
type Kleisli[A, B any]  = func(A) (B, error)
type Operator[A, B any] = func(A, error) (B, error)

// Creating values - ZERO allocations (tuples on stack)
success := result.Right(42)         // Returns (42, nil) - 0 B/op
failure := result.Left[int](err)    // Returns (0, err) - 0 B/op

// Benchmarks confirm:
// BenchmarkRight-16    789879016    1.427 ns/op    0 B/op    0 allocs/op
// BenchmarkLeft-16     895412131    1.349 ns/op    0 B/op    0 allocs/op
```

### Type Signatures

#### Standard Result

```go
// Functions take and return Result[T] structs
func Map[A, B any](f func(A) B) func(Result[A]) Result[B]
func Chain[A, B any](f Kleisli[A, B]) func(Result[A]) Result[B]
func Fold[A, B any](onLeft func(error) B, onRight func(A) B) func(Result[A]) B

// Usage requires wrapping/unwrapping
result := result.Right[error](42)
mapped := result.Map(double)(result)
value, err := result.UnwrapError(mapped)
```

#### Idiomatic Result

```go
// Functions work directly with tuples
func Map[A, B any](f func(A) B) func(A, error) (B, error)
func Chain[A, B any](f Kleisli[A, B]) func(A, error) (B, error)
func Fold[A, B any](onLeft func(error) B, onRight func(A) B) func(A, error) B

// Usage works naturally with Go's error handling
value, err := result.Right(42)
value, err = result.Map(double)(value, err)
// Can use directly: if err != nil { ... }
```

### Memory Layout

#### Standard Result (struct-based)

```
Either[error, int] struct (returned by value):
┌─────────────────────┐
│ r: int (8B)         │  Stack allocation: 24 bytes
│ l: error (8B)       │  NO heap allocation when returned by value
│ isLeft: bool (1B)   │  Benchmarks show 0 B/op, 0 allocs/op
│ padding (7B)        │
└─────────────────────┘

Key insight: Go returns small structs (<= ~64 bytes) by value on the stack.
The Either struct (24 bytes) does NOT escape to heap in normal usage.
```

#### Idiomatic Result (tuple-based)

```
(int, error) tuple:
┌─────────────────────┐
│ int: 8 bytes        │  Stack allocation: 16 bytes
│ error: 8 bytes      │  NO heap allocation
└─────────────────────┘

Both approaches achieve zero heap allocations for constructor operations!
```

### Why Both Have Zero Allocations

Both packages avoid heap allocations for simple operations:

**Standard Either/Result:**
- `Either` struct is small (24 bytes)
- Go returns by value on the stack
- Inlining eliminates function call overhead
- Result: `0 B/op, 0 allocs/op`

**Idiomatic Result:**
- Tuples are native Go multi-value returns
- Always on stack, never heap
- Even simpler than structs
- Result: `0 B/op, 0 allocs/op`

**When Either WOULD escape to heap:**
```go
// Taking address of local Either
func bad1() *Either[error, int] {
    e := Right[error](42)
    return &e  // ESCAPES: pointer to local
}

// Storing in interface
func bad2() interface{} {
    return Right[error](42)  // ESCAPES: interface boxing
}

// Closure capture with pointer receiver
func bad3() func() Either[error, int] {
    e := Right[error](42)
    return func() Either[error, int] {
        return e  // May escape depending on usage
    }
}
```

In normal functional composition (Map, Chain, Fold), neither package causes heap allocations for simple operations.

## Performance Comparison

> **Latest benchmarks:** 2025-11-18 after `either` package optimizations
>
> For detailed analysis, see [BENCHMARK_COMPARISON.md](./BENCHMARK_COMPARISON.md)

### Quick Summary (Either vs Idiomatic)

Both packages now show **excellent performance** after optimizations:

| Category | Either | Idiomatic | Winner | Speedup |
|----------|--------|-----------|--------|---------|
| **Constructors** | 1.4-1.8 ns/op | 1.2-1.4 ns/op | **TIE** | ~1.0-1.3x |
| **Predicates** | 1.5 ns/op | 1.3-1.5 ns/op | **TIE** | ~1.0x |
| **Map Operations** | 4.2-7.2 ns/op | 2.5-4.3 ns/op | **Idiomatic** | 1.2-2.1x |
| **Chain Operations** | 4.4-5.4 ns/op | 2.3-2.5 ns/op | **Idiomatic** | 1.8-2.3x |
| **ChainFirst** | **87.6 ns/op** (72 B) | **2.7 ns/op** (0 B) | **Idiomatic** | **32.4x** ✓✓✓ |
| **BiMap** | 11.5-16.8 ns/op | 3.5-3.8 ns/op | **Idiomatic** | 3.3-4.4x |
| **Alt/OrElse** | 4.0-5.7 ns/op | 2.4 ns/op | **Idiomatic** | 1.6-2.4x |
| **GetOrElse** | 6.3-9.0 ns/op | 1.5-2.1 ns/op | **Idiomatic** | 3.1-6.1x |
| **Pipelines** | 75-280 ns/op | 26-116 ns/op | **Idiomatic** | 2.4-3.4x |

### Constructor Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Winner |
|-----------|----------------|-------------------|---------|--------|
| Left      | 1.76           | **1.35**          | 1.3x    | Idiomatic ✓ |
| Right     | 1.38           | 1.43              | ~1.0x   | Tie |
| Of        | 1.68           | **1.22**          | 1.4x    | Idiomatic ✓ |

**Analysis:** After optimizations, both packages have comparable constructor performance.

### Core Transformation Operations

| Operation        | Either (ns/op) | Idiomatic (ns/op) | Speedup | Winner |
|------------------|----------------|-------------------|---------|--------|
| Map (Right)      | 5.13           | **4.34**          | 1.2x    | Idiomatic ✓ |
| Map (Left)       | 4.19           | **2.48**          | 1.7x    | Idiomatic ✓ |
| MapLeft (Right)  | 3.93           | **2.22**          | 1.8x    | Idiomatic ✓ |
| MapLeft (Left)   | 7.22           | **3.51**          | 2.1x    | Idiomatic ✓ |
| Chain (Right)    | 5.44           | **2.34**          | 2.3x    | Idiomatic ✓ |
| Chain (Left)     | 4.44           | **2.53**          | 1.8x    | Idiomatic ✓ |

### Complex Operations - The Big Difference

| Operation             | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------------------|----------------|-------------------|---------|---------------|-------------|
| **ChainFirst (Right)** | **87.62**     | **2.71**          | **32.4x** ✓✓✓ | 72 B, 3 allocs | **0 B, 0 allocs** |
| ChainFirst (Left)     | 3.94           | 2.48              | 1.6x    | 0 B | 0 B |
| BiMap (Right)         | 16.79          | **3.82**          | 4.4x    | 0 B | 0 B |
| BiMap (Left)          | 11.47          | **3.47**          | 3.3x    | 0 B | 0 B |

**Critical Insight:** ChainFirst shows the most dramatic difference - **32x faster** with **zero allocations** in idiomatic.

### Pipeline Benchmarks (Real-World Scenarios)

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Either Allocs | Idio Allocs |
|-----------|----------------|-------------------|---------|---------------|-------------|
| Pipeline Map (Right)    | 112.7  | **46.5**  | **2.4x** ✓ | 72 B, 3 allocs | 48 B, 2 allocs |
| Pipeline Chain (Right)  | 74.4   | **26.1**  | **2.9x** ✓ | 48 B, 2 allocs | 24 B, 1 alloc |
| Pipeline Complex (Right)| 279.8  | **116.3** | **2.4x** ✓ | 192 B, 8 allocs | 120 B, 5 allocs |

**Analysis:** In realistic composition scenarios, idiomatic is consistently 2-3x faster with fewer allocations.

### Extraction Operations

| Operation | Either (ns/op) | Idiomatic (ns/op) | Speedup | Winner |
|-----------|----------------|-------------------|---------|--------|
| GetOrElse (Right) | 9.01   | **1.49**  | **6.1x** ✓✓ | Idiomatic |
| GetOrElse (Left)  | 6.35   | **2.08**  | **3.1x** ✓✓ | Idiomatic |
| Alt (Right)       | 5.72   | **2.40**  | **2.4x** ✓  | Idiomatic |
| Alt (Left)        | 4.89   | **2.39**  | **2.0x** ✓  | Idiomatic |
| Fold (Right)      | 4.03   | **2.75**  | **1.5x** ✓  | Idiomatic |
| Fold (Left)       | 3.69   | **2.40**  | **1.5x** ✓  | Idiomatic |

**Analysis:** Idiomatic shows significant advantages (1.5-6x) for value extraction operations.

### Key Findings After Optimizations

1. **Both packages are now fast** - Simple operations are in the 1-5 ns/op range for both
2. **Idiomatic leads in most operations** - 1.2-2.3x faster for common transformations
3. **ChainFirst is the standout** - 32x faster with zero allocations in idiomatic
4. **Pipelines favor idiomatic** - 2-3.4x faster in realistic composition scenarios
5. **Memory efficiency** - Idiomatic consistently uses fewer allocations

### Performance Summary

**Idiomatic Advantages:**
- **Core operations**: 1.2-2.3x faster for Map, Chain, Fold
- **Complex operations**: 3-32x faster with zero allocations
- **Pipelines**: 2-3.4x faster with significantly fewer allocations
- **Extraction**: 1.5-6x faster for GetOrElse, Alt, Fold
- **Consistency**: Predictable, fast performance across all operations

**Either Advantages:**
- **Comparable performance**: After optimizations, matches idiomatic for simple operations
- **Feature richness**: More operations (Do-notation, Bind, Let, Flatten, Swap)
- **Type flexibility**: Full Either[E, A] with custom error types
- **Zero allocations**: Most simple operations have zero allocations

## API Comparison

### Creating Values

#### Standard Result

```go
import "github.com/IBM/fp-go/v2/result"

// Create success/failure
success := result.Right[error](42)
failure := result.Left[int](errors.New("oops"))

// Type annotation required
var r result.Result[int] = result.Right[error](42)
```

#### Idiomatic Result

```go
import "github.com/IBM/fp-go/v2/idiomatic/result"

// Create success/failure (more concise)
success := result.Right(42)           // (42, nil)
failure := result.Left[int](errors.New("oops"))  // (0, error)

// Native Go pattern
value, err := result.Right(42)
if err != nil {
    // handle error
}
```

### Transforming Values

#### Standard Result

```go
// Map transforms the success value
double := result.Map(N.Mul(2))
result := double(result.Right[error](21))  // Right(42)

// Chain sequences operations
validate := result.Chain(func(x int) result.Result[int] {
    if x > 0 {
        return result.Right[error](x * 2)
    }
    return result.Left[int](errors.New("negative"))
})
```

#### Idiomatic Result

```go
// Map transforms the success value
double := result.Map(N.Mul(2))
value, err := double(21, nil)  // (42, nil)

// Chain sequences operations
validate := result.Chain(func(x int) (int, error) {
    if x > 0 {
        return x * 2, nil
    }
    return 0, errors.New("negative")
})
```

### Pattern Matching

#### Standard Result

```go
// Fold extracts the value
output := result.Fold(
    func(err error) string { return "Error: " + err.Error() },
    func(n int) string { return fmt.Sprintf("Value: %d", n) },
)(myResult)

// GetOrElse with default
value := result.GetOrElse(func(err error) int { return 0 })(myResult)
```

#### Idiomatic Result

```go
// Fold extracts the value (same API, different input)
output := result.Fold(
    func(err error) string { return "Error: " + err.Error() },
    func(n int) string { return fmt.Sprintf("Value: %d", n) },
)(value, err)

// GetOrElse with default
value := result.GetOrElse(func(err error) int { return 0 })(value, err)

// Or use native Go pattern
if err != nil {
    value = 0
}
```

### Integration with Existing Code

#### Standard Result

```go
// Converting from (value, error) to Result
func doSomething() (int, error) {
    return 42, nil
}

result := result.TryCatchError(doSomething())

// Converting back to (value, error)
value, err := result.UnwrapError(result)
```

#### Idiomatic Result

```go
// Direct compatibility with (value, error)
func doSomething() (int, error) {
    return 42, nil
}

// No conversion needed!
value, err := doSomething()
value, err = result.Map(double)(value, err)
```

### Pipeline Composition

#### Standard Result

```go
import F "github.com/IBM/fp-go/v2/function"

output := F.Pipe3(
    result.Right[error](10),
    result.Map(double),
    result.Chain(validate),
    result.Map(format),
)

// Need to unwrap at the end
value, err := result.UnwrapError(output)
```

#### Idiomatic Result

```go
import F "github.com/IBM/fp-go/v2/function"

value, err := F.Pipe3(
    result.Right(10),
    result.Map(double),
    result.Chain(validate),
    result.Map(format),
)

// Already in (value, error) form
if err != nil {
    // handle error
}
```

## Detailed Design Comparison

### Type System

#### Standard Result

**Strengths:**
- Full algebraic data type semantics
- Explicit Either[E, A] allows custom error types
- Type-safe by construction
- Clear separation of error and success channels

**Weaknesses:**
- Requires wrapper structs (memory overhead)
- Less familiar to Go developers
- Needs conversion functions for Go's standard library
- More verbose type annotations

#### Idiomatic Result

**Strengths:**
- Native Go idioms (value, error) pattern
- Zero wrapper overhead
- Seamless stdlib integration
- Familiar to all Go developers
- Terser syntax

**Weaknesses:**
- Error type fixed to `error`
- Less explicit about Either semantics
- Cannot use custom error types without conversion
- Slightly less type-safe (can accidentally ignore bool/error)

### Monad Laws

Both packages satisfy the monad laws, but enforce them differently:

#### Standard Result

```go
// Left identity: return a >>= f  ≡  f a
assert.Equal(
    result.Chain(f)(result.Of(a)),
    f(a),
)

// Right identity: m >>= return  ≡  m
assert.Equal(
    result.Chain(result.Of[int])(m),
    m,
)

// Associativity: (m >>= f) >>= g  ≡  m >>= (\x -> f x >>= g)
assert.Equal(
    result.Chain(g)(result.Chain(f)(m)),
    result.Chain(func(x int) result.Result[int] {
        return result.Chain(g)(f(x))
    })(m),
)
```

#### Idiomatic Result

```go
// Same laws, different syntax
// Left identity
a, aerr := result.Of(val)
b, berr := result.Chain(f)(a, aerr)
c, cerr := f(val)
assert.Equal((b, berr), (c, cerr))

// Right identity
value, err := m()
identity := result.Chain(result.Of[int])
assert.Equal(identity(value, err), (value, err))

// Associativity (same structure, tuple-based)
```

### Error Handling Philosophy

#### Standard Result

```go
// Explicit error handling through types
func processUser(id int) result.Result[User] {
    user := fetchUser(id)  // Returns Result[User]

    return F.Pipe2(
        user,
        result.Chain(validateUser),
        result.Chain(enrichUser),
    )
}

// Must explicitly unwrap
user, err := result.UnwrapError(processUser(42))
if err != nil {
    log.Error(err)
}
```

#### Idiomatic Result

```go
// Natural Go error handling
func processUser(id int) (User, error) {
    user, err := fetchUser(id)  // Returns (User, error)

    return F.Pipe2(
        (user, err),
        result.Chain(validateUser),
        result.Chain(enrichUser),
    )
}

// Already in Go form
user, err := processUser(42)
if err != nil {
    log.Error(err)
}
```

### Composition Patterns

#### Standard Result

```go
// Applicative composition
import A "github.com/IBM/fp-go/v2/apply"

type Config struct {
    Host string
    Port int
    DB   string
}

config := A.SequenceT3(
    result.FromPredicate(validHost, hostError)(host),
    result.FromPredicate(validPort, portError)(port),
    result.FromPredicate(validDB, dbError)(db),
)(func(h string, p int, d string) Config {
    return Config{h, p, d}
})
```

#### Idiomatic Result

```go
// Direct tuple composition
config, err := func() (Config, error) {
    host, err := result.FromPredicate(validHost, hostError)(host)
    if err != nil {
        return Config{}, err
    }

    port, err := result.FromPredicate(validPort, portError)(port)
    if err != nil {
        return Config{}, err
    }

    db, err := result.FromPredicate(validDB, dbError)(db)
    if err != nil {
        return Config{}, err
    }

    return Config{host, port, db}, nil
}()
```

## When to Use Each

### Use Idiomatic Result When (Recommended for Most Cases):

1. **Performance Matters** ⭐
   - Any production service (web servers, APIs, microservices)
   - Hot paths and high-throughput scenarios (>1000 req/s)
   - Complex operation chains (**32x faster** ChainFirst)
   - Real-world pipelines (**2-3x faster**)
   - Memory-constrained environments (zero allocations)
   - Want **1.2-6x speedup** across most operations

2. **Go Integration** ⭐⭐
   - Working with existing Go codebases
   - Interfacing with standard library (native (value, error))
   - Team familiar with Go, new to FP
   - Want zero-cost functional abstractions
   - Seamless error handling patterns

3. **Pragmatic Functional Programming**
   - Value performance AND functional patterns
   - Prefer Go idioms over FP terminology
   - Simpler function signatures
   - Lower cognitive overhead
   - Production-ready patterns

4. **Real-World Applications**
   - Web servers, REST APIs, gRPC services
   - CLI tools and command-line applications
   - Data processing pipelines
   - Any latency-sensitive application
   - Systems with tight performance budgets

**Performance Gains:** Use idiomatic for 1.2-32x speedup depending on operation, with consistently lower allocations.

### Use Standard Either/Result When:

1. **Type Safety & Flexibility**
   - Need explicit Either[E, A] with **custom error types**
   - Building domain-specific error hierarchies
   - Want to distinguish different error categories at type level
   - Type system enforcement is critical

2. **Advanced FP Features**
   - Using Do-notation for complex monadic compositions
   - Need operations like Flatten, Swap, Bind, Let
   - Leveraging advanced type classes (Semigroup, Monoid)
   - Want the complete FP toolkit

3. **FP Expertise & Education**
   - Porting code from other FP languages (Scala, Haskell)
   - Teaching functional programming concepts
   - Team has strong FP background
   - Explicit algebraic data types preferred
   - Code review benefits from FP terminology

4. **Performance is Acceptable**
   - After optimizations, Either is **quite fast** (1-5 ns/op for simple operations)
   - Difference matters mainly at high scale (millions of operations)
   - Code clarity > micro-optimizations
   - Simple operations dominate your workload

**Note:** Either package is now performant enough for most use cases. Choose it for features, not performance concerns.

### Hybrid Approach

You can use both packages together:

```go
import (
    stdResult "github.com/IBM/fp-go/v2/result"
    "github.com/IBM/fp-go/v2/idiomatic/result"
)

// Use standard for complex types
type ValidationError struct {
    Field string
    Error string
}

func validateInput(input string) stdResult.Either[ValidationError, Input] {
    // ... validation logic
}

// Convert to idiomatic for performance
func processInput(input string) (Output, error) {
    validated := validateInput(input)
    value, err := stdResult.UnwrapError(
        stdResult.MapLeft(toError)(validated),
    )

    // Use idiomatic for hot path
    return result.Chain(heavyProcessing)(value, err)
}
```

## Migration Guide

### From Standard to Idiomatic

```go
// Before (standard)
import "github.com/IBM/fp-go/v2/result"

func process(x int) result.Result[int] {
    return F.Pipe2(
        result.Right[error](x),
        result.Map(double),
        result.Chain(validate),
    )
}

// After (idiomatic)
import "github.com/IBM/fp-go/v2/idiomatic/result"

func process(x int) (int, error) {
    return F.Pipe2(
        result.Right(x),
        result.Map(double),
        result.Chain(validate),
    )
}
```

### Key Changes

1. **Type signatures**: `Result[T]` → `(T, error)`
2. **Kleisli**: `func(A) Result[B]` → `func(A) (B, error)`
3. **Operator**: `func(Result[A]) Result[B]` → `func(A, error) (B, error)`
4. **Return values**: Function calls return tuples, not wrapped values
5. **Pattern matching**: Same Fold/GetOrElse API, different inputs

## Conclusion

### Performance Summary (After Either Optimizations)

The latest benchmark results show a clear pattern:

**Both packages are now fast**, but idiomatic consistently leads:

- **Constructors & Predicates**: Both ~1-2 ns/op (essentially tied)
- **Core transformations**: Idiomatic **1.2-2.3x faster** (Map, Chain, Fold)
- **Complex operations**: Idiomatic **3-32x faster** (BiMap, ChainFirst)
- **Pipelines**: Idiomatic **2-3.4x faster** with fewer allocations
- **Extraction**: Idiomatic **1.5-6x faster** (GetOrElse, Alt)

**Key Insight:** The idiomatic package delivers **consistently better performance** across the board while maintaining zero-cost abstractions. The Either package is now fast enough for most use cases, but idiomatic is the performance winner.

### Updated Recommendation Matrix

| Scenario | Recommendation | Reason |
|----------|---------------|--------|
| **New Go project** | **Idiomatic** ⭐ | Natural Go patterns, 1.2-6x faster, better integration |
| **Production services** | **Idiomatic** ⭐⭐ | 2-3x faster pipelines, zero allocations, proven performance |
| **Performance critical** | **Idiomatic** ⭐⭐⭐ | 32x faster complex ops, minimal allocations |
| **Microservices/APIs** | **Idiomatic** ⭐⭐ | High throughput, familiar patterns, better performance |
| **CLI Tools** | **Idiomatic** ⭐ | Low overhead, Go idioms, fast |
| Custom error types | Standard/Either | Need Either[E, A] with domain types |
| Learning FP | Standard/Either | Clearer ADT semantics, educational |
| FP-heavy codebase | Standard/Either | Consistency, Do-notation, full FP toolkit |
| Library/Framework | Either way | Both are good; choose based on API style |

### Real-World Impact

For a service handling 10,000 requests/second with typical pipeline operations:

```
Either package:     280 ns/op × 10M req/day = 2,800 seconds = 46.7 minutes
Idiomatic package:  116 ns/op × 10M req/day = 1,160 seconds = 19.3 minutes
Time saved: 27.4 minutes of CPU time per day
```

At scale, this translates to:
- Lower latency (2-3x faster response times for FP operations)
- Reduced CPU usage (fewer cores needed)
- Lower memory pressure (significantly fewer allocations)
- Better resource utilization

### Final Recommendation

**For most Go projects:** Use **idiomatic packages**
- 1.2-32x faster across operations
- Native Go idioms
- Zero-cost abstractions
- Production-proven performance
- Easier integration

**For specialized needs:** Use **standard Either/Result**
- Need custom error types Either[E, A]
- Want Do-notation and advanced FP features
- Porting from FP languages
- Educational/learning context
- FP-heavy existing codebase

### Bottom Line

After optimizations, both packages are excellent:

- **Either/Result**: Fast enough for most use cases, feature-rich, type-safe
- **Idiomatic**: **Faster in practice** (1.2-32x), native Go, zero-cost FP

The idiomatic packages now represent the **best of both worlds**: full functional programming capabilities with Go's native performance and idioms. Unless you specifically need Either[E, A]'s custom error types or advanced FP features, **idiomatic is the recommended choice** for production Go services.

Both maintain the core benefits of functional programming—choose based on whether you prioritize performance & Go integration (idiomatic) or type flexibility & FP features (either).
