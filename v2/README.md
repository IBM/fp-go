# fp-go V2: Enhanced Functional Programming for Go 1.24+

[![Go Reference](https://pkg.go.dev/badge/github.com/IBM/fp-go/v2.svg)](https://pkg.go.dev/github.com/IBM/fp-go/v2)
[![Coverage Status](https://coveralls.io/repos/github/IBM/fp-go/badge.svg?branch=main&flag=v2)](https://coveralls.io/github/IBM/fp-go?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/IBM/fp-go/v2)](https://goreportcard.com/report/github.com/IBM/fp-go/v2)

**fp-go** is a comprehensive functional programming library for Go, bringing type-safe functional patterns inspired by [fp-ts](https://gcanti.github.io/fp-ts/) to the Go ecosystem. Version 2 leverages [generic type aliases](https://github.com/golang/go/issues/46477) introduced in Go 1.24, providing a more ergonomic and streamlined API.

## 📚 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Requirements](#-requirements)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Breaking Changes](#️-breaking-changes)
- [Key Improvements](#-key-improvements)
- [Migration Guide](#-migration-guide)
- [What's New](#-whats-new)
- [Documentation](#-documentation)
- [Contributing](#-contributing)
- [License](#-license)

## 🎯 Overview

fp-go brings the power of functional programming to Go with:

- **Type-safe abstractions** - Monads, Functors, Applicatives, and more
- **Composable operations** - Build complex logic from simple, reusable functions
- **Error handling** - Elegant error management with `Either`, `Result`, and `IOEither`
- **Lazy evaluation** - Control when and how computations execute
- **Optics** - Powerful lens, prism, and traversal operations for immutable data manipulation

## ✨ Features

- 🔒 **Type Safety** - Leverage Go's generics for compile-time guarantees
- 🧩 **Composability** - Chain operations naturally with functional composition
- 📦 **Rich Type System** - `Option`, `Either`, `Result`, `IO`, `Reader`, and more
- 🎯 **Practical** - Designed for real-world Go applications
- 🚀 **Performance** - Zero-cost abstractions where possible
- 📖 **Well-documented** - Comprehensive API documentation and examples
- 🧪 **Battle-tested** - Extensive test coverage

## 🔧 Requirements

- **Go 1.24 or later** (for generic type alias support)

## 📦 Installation

```bash
go get github.com/IBM/fp-go/v2
```

## 🚀 Quick Start

### Working with Option

```go
package main

import (
    "fmt"
    "github.com/IBM/fp-go/v2/option"
    N "github.com/IBM/fp-go/v2/number"
)

func main() {
    // Create an Option
    some := option.Some(42)
    none := option.None[int]()
    
    // Map over values
    doubled := option.Map(N.Mul(2))(some)
    fmt.Println(option.GetOrElse(0)(doubled)) // Output: 84
    
    // Chain operations
    result := option.Chain(func(x int) option.Option[string] {
        if x > 0 {
            return option.Some(fmt.Sprintf("Positive: %d", x))
        }
        return option.None[string]()
    })(some)
    
    fmt.Println(option.GetOrElse("No value")(result)) // Output: Positive: 42
}
```

### Error Handling with Result

```go
package main

import (
    "errors"
    "fmt"
    "github.com/IBM/fp-go/v2/result"
)

func divide(a, b int) result.Result[int] {
    if b == 0 {
        return result.Error[int](errors.New("division by zero"))
    }
    return result.Ok(a / b)
}

func main() {
    res := divide(10, 2)
    
    // Pattern match on the result
    result.Fold(
        func(err error) { fmt.Println("Error:", err) },
        func(val int) { fmt.Println("Result:", val) },
    )(res)
    // Output: Result: 5
    
    // Or use GetOrElse for a default value
    value := result.GetOrElse(0)(divide(10, 0))
    fmt.Println("Value:", value) // Output: Value: 0
}
```

### Composing IO Operations

```go
package main

import (
    "fmt"
    "github.com/IBM/fp-go/v2/io"
)

func main() {
    // Define pure IO operations
    readInput := io.MakeIO(func() string {
        return "Hello, fp-go!"
    })
    
    // Transform the result
    uppercase := io.Map(func(s string) string {
        return fmt.Sprintf(">>> %s <<<", s)
    })(readInput)
    
    // Execute the IO operation
    result := uppercase()
    fmt.Println(result) // Output: >>> Hello, fp-go! <<<
}
```

## ⚠️ Breaking Changes

### From V1 to V2

#### 1. Generic Type Aliases

V2 uses [generic type aliases](https://github.com/golang/go/issues/46477) which require Go 1.24+. This is the most significant change and enables cleaner type definitions.

**V1:**
```go
type ReaderIOEither[R, E, A any] RD.Reader[R, IOE.IOEither[E, A]]
```

**V2:**
```go
type ReaderIOEither[R, E, A any] = RD.Reader[R, IOE.IOEither[E, A]]
```

#### 2. Generic Type Parameter Ordering

Type parameters that **cannot** be inferred from function arguments now come first, improving type inference.

**V1:**
```go
// Ap in V1 - less intuitive ordering
func Ap[R, E, A, B any](fa ReaderIOEither[R, E, A]) func(ReaderIOEither[R, E, func(A) B]) ReaderIOEither[R, E, B]
```

**V2:**
```go
// Ap in V2 - B comes first as it cannot be inferred
func Ap[B, R, E, A any](fa ReaderIOEither[R, E, A]) func(ReaderIOEither[R, E, func(A) B]) ReaderIOEither[R, E, B]
```

This change allows the Go compiler to infer more types automatically, reducing the need for explicit type parameters.

#### 3. Pair Monad Semantics

Monadic operations for `Pair` now operate on the **second argument** to align with the [Haskell definition](https://hackage.haskell.org/package/TypeCompose-0.9.14/docs/Data-Pair.html).

**V1:**
```go
// Operations on first element
pair := MakePair(1, "hello")
result := Map(N.Mul(2))(pair) // Pair(2, "hello")
```

**V2:**
```go
// Operations on second element (Haskell-compatible)
pair := MakePair(1, "hello")
result := Map(func(s string) string { return s + "!" })(pair) // Pair(1, "hello!")
```

#### 4. Endomorphism Compose Semantics

The `Compose` function for endomorphisms now follows **mathematical function composition** (right-to-left execution), aligning with standard functional programming conventions.

**V1:**
```go
// Compose executed left-to-right
double := N.Mul(2)
increment := N.Add(1)
composed := Compose(double, increment)
result := composed(5) // (5 * 2) + 1 = 11
```

**V2:**
```go
// Compose executes RIGHT-TO-LEFT (mathematical composition)
double := N.Mul(2)
increment := N.Add(1)
composed := Compose(double, increment)
result := composed(5) // (5 + 1) * 2 = 12

// Use MonadChain for LEFT-TO-RIGHT execution
chained := MonadChain(double, increment)
result2 := chained(5) // (5 * 2) + 1 = 11
```

**Key Difference:**
- `Compose(f, g)` now means `f ∘ g`, which applies `g` first, then `f` (right-to-left)
- `MonadChain(f, g)` applies `f` first, then `g` (left-to-right)

## ✨ Key Improvements

### 1. Simplified Type Declarations

Generic type aliases eliminate the need for namespace imports in type declarations.

**V1 Approach:**
```go
import (
    ET "github.com/IBM/fp-go/either"
    OPT "github.com/IBM/fp-go/option"
)

func processData(input string) ET.Either[error, OPT.Option[int]] {
    // implementation
}
```

**V2 Approach:**
```go
import (
    "github.com/IBM/fp-go/v2/result"
    "github.com/IBM/fp-go/v2/option"
)

// Define type aliases once
type Result[A any] = result.Result[A]
type Option[A any] = option.Option[A]

// Use them throughout your codebase
func processData(input string) Result[Option[int]] {
    // implementation
}
```

### 2. No More `generic` Subpackages

The library implementation no longer requires separate `generic` subpackages, making the codebase simpler and easier to understand.

**V1 Structure:**
```
either/
  either.go
  generic/
    either.go  // Generic implementation
```

**V2 Structure:**
```
either/
  either.go  // Single, clean implementation
```

### 3. Better Type Inference

The reordered type parameters allow the Go compiler to infer more types automatically:

**V1:**
```go
// Often need explicit type parameters
result := Map[Context, error, int, string](transform)(value)
```

**V2:**
```go
// Compiler can infer more types
result := Map(transform)(value)  // Cleaner!
```

## 🚀 Migration Guide

### Step 1: Update Go Version

Ensure you're using Go 1.24 or later:

```bash
go version  # Should show go1.24 or higher
```

### Step 2: Update Import Paths

Change all import paths from `github.com/IBM/fp-go` to `github.com/IBM/fp-go/v2`:

**Before:**
```go
import (
    "github.com/IBM/fp-go/either"
    "github.com/IBM/fp-go/option"
)
```

**After:**
```go
import (
    "github.com/IBM/fp-go/v2/either"
    "github.com/IBM/fp-go/v2/option"
)
```

### Step 3: Remove `generic` Subpackage Imports

If you were using generic subpackages, remove them:

**Before:**
```go
import (
    E "github.com/IBM/fp-go/either/generic"
)
```

**After:**
```go
import (
    "github.com/IBM/fp-go/v2/either"
)
```

### Step 4: Update Type Parameter Order

Review functions like `Ap` where type parameter order has changed. The compiler will help identify these:

**Before:**
```go
result := Ap[Context, error, int, string](value)(funcInContext)
```

**After:**
```go
result := Ap[string, Context, error, int](value)(funcInContext)
// Or better yet, let the compiler infer:
result := Ap(value)(funcInContext)
```

### Step 5: Update Pair Operations

If you're using `Pair`, update operations to work on the second element:

**Before (V1):**
```go
pair := MakePair(42, "data")
// Map operates on first element
result := Map(N.Mul(2))(pair)
```

**After (V2):**
```go
pair := MakePair(42, "data")
// Map operates on second element
result := Map(func(s string) string { return s + "!" })(pair)
```

### Step 6: Simplify Type Aliases

Create project-wide type aliases for common patterns:

```go
// types.go - Define once, use everywhere
package myapp

import (
    "github.com/IBM/fp-go/v2/result"
    "github.com/IBM/fp-go/v2/option"
    "github.com/IBM/fp-go/v2/ioresult"
)

type Result[A any] = result.Result[A]
type Option[A any] = option.Option[A]
type IOResult[A any] = ioresult.IOResult[A]
```

## 🆕 What's New

### Cleaner API Surface

The elimination of `generic` subpackages means:
- Fewer imports to manage
- Simpler package structure
- Easier to navigate documentation
- More intuitive API

### Example: Before and After

**V1 Complex Example:**
```go
import (
    ET "github.com/IBM/fp-go/either"
    EG "github.com/IBM/fp-go/either/generic"
    IOET "github.com/IBM/fp-go/ioeither"
    IOEG "github.com/IBM/fp-go/ioeither/generic"
)

func process() IOET.IOEither[error, string] {
    return IOEG.Map[error, int, string](
        strconv.Itoa,
    )(fetchData())
}
```

**V2 Simplified Example:**
```go
import (
    "strconv"
    "github.com/IBM/fp-go/v2/ioresult"
)

type IOResult[A any] = ioresult.IOResult[A]

func process() IOResult[string] {
    return ioresult.Map(
        strconv.Itoa,
    )(fetchData())
}
```

## 📚 Documentation

- **[Design Decisions](./DESIGN.md)** - Key design principles and patterns explained
- **[Functional I/O in Go](./FUNCTIONAL_IO.md)** - Understanding Context, errors, and the Reader pattern for I/O operations
- **[Idiomatic vs Standard Packages](./IDIOMATIC_COMPARISON.md)** - Performance comparison and when to use each approach
- **[API Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2)** - Complete API reference
- **[Code Samples](./samples/)** - Practical examples and use cases
- **[Go 1.24 Release Notes](https://tip.golang.org/doc/go1.24)** - Information about generic type aliases

### Core Modules

#### Standard Packages (Struct-based)
- **Option** - Represent optional values without nil
- **Either** - Type-safe error handling with left/right values
- **Result** - Simplified Either with error as left type (recommended for error handling)
- **IO** - Lazy evaluation and side effect management
- **IOOption** - Combine IO with Option for optional values with side effects
- **IOResult** - Combine IO with Result for error handling (recommended over IOEither)
- **Reader** - Dependency injection pattern
- **ReaderOption** - Combine Reader with Option for optional values with dependency injection
- **ReaderIOOption** - Combine Reader, IO, and Option for optional values with dependency injection and side effects
- **ReaderIOResult** - Combine Reader, IO, and Result for complex workflows
- **Array** - Functional array operations
- **Record** - Functional record/map operations
- **[Optics](./optics/README.md)** - Lens, Prism, Optional, and Traversal for immutable updates

#### Idiomatic Packages (Tuple-based, High Performance)
- **idiomatic/option** - Option monad using native Go `(value, bool)` tuples
- **idiomatic/result** - Result monad using native Go `(value, error)` tuples
- **idiomatic/ioresult** - IOResult monad using `func() (value, error)` for IO operations
- **idiomatic/readerresult** - Reader monad combined with Result pattern
- **idiomatic/readerioresult** - Reader monad combined with IOResult pattern

The idiomatic packages offer 2-10x performance improvements and zero allocations by using Go's native tuple patterns instead of struct wrappers. Use them for performance-critical code or when you prefer Go's native error handling style.

## 🤔 Should I Migrate?

**Migrate to V2 if:**
- ✅ You can use Go 1.24+
- ✅ You want cleaner, more maintainable code
- ✅ You want better type inference
- ✅ You're starting a new project

**Stay on V1 if:**
- ⚠️ You're locked to Go < 1.24
- ⚠️ Migration effort outweighs benefits for your project
- ⚠️ You need stability in production (V2 is newer)

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. **Report bugs** - Open an issue with a clear description and reproduction steps
2. **Suggest features** - Share your ideas for improvements
3. **Submit PRs** - Fix bugs or add features (please discuss major changes first)
4. **Improve docs** - Help make the documentation clearer and more comprehensive

Please read our contribution guidelines before submitting pull requests.

## 🐛 Issues and Feedback

Found a bug or have a suggestion? Please [open an issue](https://github.com/IBM/fp-go/issues) on GitHub.

## 📄 License

This project is licensed under the Apache License 2.0. See the [LICENSE](https://github.com/IBM/fp-go/blob/main/LICENSE) file for details.

---

**Made with ❤️ by IBM**