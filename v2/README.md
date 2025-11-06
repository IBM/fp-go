# fp-go V2: Enhanced Functional Programming for Go 1.24+

[![Go Reference](https://pkg.go.dev/badge/github.com/IBM/fp-go/v2.svg)](https://pkg.go.dev/github.com/IBM/fp-go/v2)

Version 2 of fp-go leverages [generic type aliases](https://github.com/golang/go/issues/46477) introduced in Go 1.24, providing a more ergonomic and streamlined API.

## 📚 Table of Contents

- [Requirements](#-requirements)
- [Breaking Changes](#-breaking-changes)
- [Key Improvements](#-key-improvements)
- [Migration Guide](#-migration-guide)
- [Installation](#-installation)
- [What's New](#-whats-new)

## 🔧 Requirements

- **Go 1.24 or later** (for generic type alias support)

## ⚠️ Breaking Changes

### 1. Generic Type Aliases

V2 uses [generic type aliases](https://github.com/golang/go/issues/46477) which require Go 1.24+. This is the most significant change and enables cleaner type definitions.

**V1:**
```go
type ReaderIOEither[R, E, A any] RD.Reader[R, IOE.IOEither[E, A]]
```

**V2:**
```go
type ReaderIOEither[R, E, A any] = RD.Reader[R, IOE.IOEither[E, A]]
```

### 2. Generic Type Parameter Ordering

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

### 3. Pair Monad Semantics

Monadic operations for `Pair` now operate on the **second argument** to align with the [Haskell definition](https://hackage.haskell.org/package/TypeCompose-0.9.14/docs/Data-Pair.html).

**V1:**
```go
// Operations on first element
pair := MakePair(1, "hello")
result := Map(func(x int) int { return x * 2 })(pair) // Pair(2, "hello")
```

**V2:**
```go
// Operations on second element (Haskell-compatible)
pair := MakePair(1, "hello")
result := Map(func(s string) string { return s + "!" })(pair) // Pair(1, "hello!")
```

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
    "github.com/IBM/fp-go/v2/either"
    "github.com/IBM/fp-go/v2/option"
)

// Define type aliases once
type Either[A any] = either.Either[error, A]
type Option[A any] = option.Option[A]

// Use them throughout your codebase
func processData(input string) Either[Option[int]] {
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
result := Map(func(x int) int { return x * 2 })(pair)
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
    "github.com/IBM/fp-go/v2/either"
    "github.com/IBM/fp-go/v2/option"
    "github.com/IBM/fp-go/v2/ioeither"
)

type Either[A any] = either.Either[error, A]
type Option[A any] = option.Option[A]
type IOEither[A any] = ioeither.IOEither[error, A]
```

## 📦 Installation

```bash
go get github.com/IBM/fp-go/v2
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
        func(x int) string { return fmt.Sprintf("%d", x) },
    )(fetchData())
}
```

**V2 Simplified Example:**
```go
import (
    "github.com/IBM/fp-go/v2/either"
    "github.com/IBM/fp-go/v2/ioeither"
)

type IOEither[A any] = ioeither.IOEither[error, A]

func process() IOEither[string] {
    return ioeither.Map(
        func(x int) string { return fmt.Sprintf("%d", x) },
    )(fetchData())
}
```

## 📚 Additional Resources

- [Main README](../README.md) - Core concepts and design philosophy
- [API Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2)
- [Code Samples](../samples/)
- [Go 1.24 Release Notes](https://tip.golang.org/doc/go1.24)

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

## 🐛 Issues and Feedback

Found a bug or have a suggestion? Please [open an issue](https://github.com/IBM/fp-go/issues) on GitHub.

## 📄 License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.