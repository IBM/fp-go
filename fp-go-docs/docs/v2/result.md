---
title: Result
hide_title: true
description: Type-safe error handling with Ok/Err - specialized Either for Go's error interface.
sidebar_position: 10
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="Result"
  lede="Type-safe error handling with Ok and Error variants. Result[A] is a specialized Either[error, A] designed for idiomatic Go error handling."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/result' },
    { label: 'Type', value: 'Monad (Either[error, A])' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

Result represents a computation that may fail with an error:
- **Ok** (Right): Success value of type `A`
- **Error** (Left): Failure with Go's `error` interface

<CodeCard file="type_definition.go">
{`package result

// Result is Either specialized for error handling
type Result[A any] = Either[error, A]
`}
</CodeCard>

### Why Result over Either?

<ApiTable>
| Feature | Result[A] | Either[E, A] |
|---------|-----------|--------------|
| Error type | Always `error` | Generic `E` |
| Type params | 1 (simpler) | 2 (more flexible) |
| Go interop | Native `(value, error)` | Requires conversion |
| Use case | Error handling | General sum types |
</ApiTable>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ok` | `func Ok[A any](value A) Result[A]` | Create success value |
| `Error` | `func Error[A any](err error) Result[A]` | Create error value |
| `Of` | `func Of[A any](value A) Result[A]` | Alias for Ok (monadic pure) |
| `TryCatchError` | `func TryCatchError[A any](value A, err error) Result[A]` | From `(value, error)` tuple |
| `FromError` | `func FromError[A any](f func(...) (A, error)) func(...) Result[A]` | Wrap error-returning function |
| `Try` | `func Try[A any](f func() A) Result[A]` | Catch panics as errors |
</ApiTable>

### Pattern Matching

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Fold` | `func Fold[A, B any](onError func(error) B, onOk func(A) B) func(Result[A]) B` | Extract value from both cases |
| `Match` | `func Match[A, B any](onError func(error) Result[B], onOk func(A) Result[B]) func(Result[A]) Result[B]` | Pattern match with Result return |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[A, B any](f func(A) B) func(Result[A]) Result[B]` | Transform success value |
| `Chain` | `func Chain[A, B any](f func(A) Result[B]) func(Result[A]) Result[B]` | FlatMap - sequence operations |
| `MapError` | `func MapError[A any](f func(error) error) func(Result[A]) Result[A]` | Transform error value |
| `BiMap` | `func BiMap[A, B any](fe func(error) error, fa func(A) B) func(Result[A]) Result[B]` | Transform both sides |
</ApiTable>

### Combining

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ap` | `func Ap[A, B any](fa Result[A]) func(Result[func(A) B]) Result[B]` | Apply wrapped function |
| `SequenceArray` | `func SequenceArray[A any]([]Result[A]) Result[[]A]` | All-or-nothing for arrays |
| `TraverseArray` | `func TraverseArray[A, B any](f func(A) Result[B]) func([]A) Result[[]B]` | Map and sequence |
</ApiTable>

### Extraction

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Unwrap` | `func Unwrap[A any](Result[A]) (A, error)` | Convert to `(value, error)` tuple |
| `GetOrElse` | `func GetOrElse[A any](f func() A) func(Result[A]) A` | Extract with default |
| `ToOption` | `func ToOption[A any](Result[A]) Option[A]` | Convert to Option (discards error) |
</ApiTable>

### Testing

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `IsOk` | `func IsOk[A any](Result[A]) bool` | Test for success |
| `IsError` | `func IsError[A any](Result[A]) bool` | Test for error |
| `Exists` | `func Exists[A any](pred func(A) bool) func(Result[A]) bool` | Test success value |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    "errors"
    "fmt"
    R "github.com/IBM/fp-go/v2/result"
)

func main() {
    // Create Results
    success := R.Ok(42)
    failure := R.Error[int](errors.New("something went wrong"))
    
    // Check values
    fmt.Println(R.IsOk(success))    // true
    fmt.Println(R.IsError(failure)) // true
    
    // Extract with default
    value := R.GetOrElse(func() int { return 0 })(success)
    fmt.Println(value) // 42
    
    // Convert to Go idiom
    val, err := R.Unwrap(success)
    fmt.Println(val, err) // 42 <nil>
}
`}
</CodeCard>

### From Go Idioms

<CodeCard file="from_go.go">
{`package main

import (
    "os"
    R "github.com/IBM/fp-go/v2/result"
)

// Convert (value, error) tuple
func readFile(path string) R.Result[[]byte] {
    return R.TryCatchError(os.ReadFile(path))
}

// Wrap error-returning function
var readFileFunc = R.FromError(os.ReadFile)

func main() {
    // Use wrapped function
    content := readFileFunc("config.json")
    
    // Handle result
    R.Match(
        func(err error) R.Result[string] {
            return R.Error[string](err)
        },
        func(data []byte) R.Result[string] {
            return R.Ok(string(data))
        },
    )(content)
}
`}
</CodeCard>

### Chaining Operations

<CodeCard file="chaining.go">
{`package main

import (
    "errors"
    "strconv"
    R "github.com/IBM/fp-go/v2/result"
    F "github.com/IBM/fp-go/v2/function"
)

func parseInt(s string) R.Result[int] {
    n, err := strconv.Atoi(s)
    return R.TryCatchError(n, err)
}

func validatePositive(n int) R.Result[int] {
    if n > 0 {
        return R.Ok(n)
    }
    return R.Error[int](errors.New("must be positive"))
}

func double(n int) R.Result[int] {
    return R.Ok(n * 2)
}

func main() {
    // Chain operations - short-circuits on first error
    result := F.Pipe3(
        parseInt("42"),
        R.Chain(validatePositive),
        R.Chain(double),
    )
    
    // Extract value
    value := R.GetOrElse(func() int { return 0 })(result)
    fmt.Println(value) // 84
}
`}
</CodeCard>

### Error Context

<CodeCard file="error_context.go">
{`package main

import (
    "fmt"
    R "github.com/IBM/fp-go/v2/result"
    F "github.com/IBM/fp-go/v2/function"
)

func fetchUser(id string) R.Result[User] {
    // ... fetch logic
}

func processUser(id string) R.Result[User] {
    return F.Pipe2(
        fetchUser(id),
        R.Chain(validateUser),
        R.MapError(func(err error) error {
            return fmt.Errorf("failed to process user %s: %w", id, err)
        }),
    )
}
`}
</CodeCard>

### Sequence Operations

<CodeCard file="sequence.go">
{`package main

import (
    "strconv"
    R "github.com/IBM/fp-go/v2/result"
)

func main() {
    // Parse all strings - fails if any fail
    parseAll := R.TraverseArray(func(s string) R.Result[int] {
        n, err := strconv.Atoi(s)
        return R.TryCatchError(n, err)
    })
    
    result := parseAll([]string{"1", "2", "3"})
    // Ok([1, 2, 3])
    
    result = parseAll([]string{"1", "bad", "3"})
    // Error("invalid syntax")
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern 1: Resource Management

<CodeCard file="resource.go">
{`package main

import (
    "os"
    R "github.com/IBM/fp-go/v2/result"
)

func processFile(path string) R.Result[Data] {
    return R.Bracket(
        // Acquire
        func() R.Result[*os.File] {
            return R.TryCatchError(os.Open(path))
        },
        // Use
        func(f *os.File) R.Result[Data] {
            return parseData(f)
        },
        // Release (always runs)
        func(f *os.File, _ R.Result[Data]) R.Result[unit.Unit] {
            f.Close()
            return R.Ok(unit.Unit{})
        },
    )
}
`}
</CodeCard>

### Pattern 2: Validation

<Compare>
<CompareCol kind="bad">
<CodeCard file="traditional.go">
{`// ❌ Traditional Go
func validateUser(u User) (User, error) {
    if u.Name == "" {
        return User{}, errors.New("name required")
    }
    if u.Age < 18 {
        return User{}, errors.New("must be 18+")
    }
    return u, nil
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="good">
<CodeCard file="functional.go">
{`// ✅ With Result
func validateUser(u User) R.Result[User] {
    return F.Pipe2(
        validateName(u),
        R.Chain(validateAge),
    )
}

func validateName(u User) R.Result[User] {
    if u.Name == "" {
        return R.Error[User](errors.New("name required"))
    }
    return R.Ok(u)
}
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>
