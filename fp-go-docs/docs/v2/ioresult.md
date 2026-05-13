---
title: IOResult
hide_title: true
description: Lazy side effects with Go error handling - combine IO and Result for effectful computations.
sidebar_position: 13
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="IOResult"
  lede="Combine lazy evaluation (IO) with Go's error interface (Result). IOResult[A] represents a synchronous computation with side effects that may fail."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/ioresult' },
    { label: 'Type', value: 'Monad (IO[Result[A]])' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

IOResult is IOEither specialized for Go's error interface:

<CodeCard file="type_definition.go">
{`package ioresult

// IOResult is IO of Result
type IOResult[A any] = IO[Result[A]]
// Which expands to: func() Result[A]
// Which expands to: func() Either[error, A]
`}
</CodeCard>

### When to Use

<ApiTable>
| Use IOResult When | Use IOEither When |
|-------------------|-------------------|
| Standard Go error handling | Custom error types |
| Library interoperability | Domain-specific errors |
| Simpler type signatures (1 param) | Type-level error distinction |
| Working with existing Go code | Rich error information |
</ApiTable>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ok` | `func Ok[A any](value A) IOResult[A]` | Create successful IOResult |
| `Error` | `func Error[A any](err error) IOResult[A]` | Create failed IOResult |
| `Of` | `func Of[A any](value A) IOResult[A]` | Alias for Ok |
| `TryCatchError` | `func TryCatchError[A any](f func() (A, error)) IOResult[A]` | From function returning tuple |
| `FromIO` | `func FromIO[A any](io IO[A]) IOResult[A]` | Lift IO to IOResult |
| `FromResult` | `func FromResult[A any](r Result[A]) IOResult[A]` | Lift Result to IOResult |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[A, B any](f func(A) B) func(IOResult[A]) IOResult[B]` | Transform success value |
| `MapError` | `func MapError[A any](f func(error) error) func(IOResult[A]) IOResult[A]` | Transform error |
| `BiMap` | `func BiMap[A, B any](fe func(error) error, fa func(A) B) func(IOResult[A]) IOResult[B]` | Transform both |
| `Chain` | `func Chain[A, B any](f func(A) IOResult[B]) func(IOResult[A]) IOResult[B]` | Sequence operations |
</ApiTable>

### Combining

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ap` | `func Ap[A, B any](fa IOResult[A]) func(IOResult[func(A) B]) IOResult[B]` | Apply wrapped function |
| `SequenceArray` | `func SequenceArray[A any]([]IOResult[A]) IOResult[[]A]` | All-or-nothing |
| `TraverseArray` | `func TraverseArray[A, B any](f func(A) IOResult[B]) func([]A) IOResult[[]B]` | Map and sequence |
</ApiTable>

### Extraction

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Unwrap` | `func Unwrap[A any](IOResult[A]) IO[tuple.Tuple2[A, error]]` | Convert to IO of tuple |
| `ToOption` | `func ToOption[A any](IOResult[A]) IO[Option[A]]` | Convert to IO of Option |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    "errors"
    IOR "github.com/IBM/fp-go/v2/ioresult"
)

func main() {
    // Create IOResult values
    success := IOR.Ok(42)
    failure := IOR.Error[int](errors.New("something went wrong"))
    
    // Execute
    result := success()  // Result[int] = Ok(42)
    result = failure()   // Result[int] = Error("something went wrong")
}
`}
</CodeCard>

### File Operations

<CodeCard file="file_ops.go">
{`package main

import (
    "os"
    IOR "github.com/IBM/fp-go/v2/ioresult"
)

func readFile(path string) IOR.IOResult[[]byte] {
    return IOR.TryCatchError(func() ([]byte, error) {
        return os.ReadFile(path)
    })
}

func writeFile(path string, data []byte) IOR.IOResult[unit.Unit] {
    return IOR.TryCatchError(func() (unit.Unit, error) {
        err := os.WriteFile(path, data, 0644)
        return unit.Unit{}, err
    })
}

func main() {
    result := readFile("config.json")()
    // Result[[]byte]
}
`}
</CodeCard>

### Chaining Operations

<CodeCard file="chaining.go">
{`package main

import (
    IOR "github.com/IBM/fp-go/v2/ioresult"
    F "github.com/IBM/fp-go/v2/function"
)

func fetchUser(id string) IOR.IOResult[User] {
    return IOR.TryCatchError(func() (User, error) {
        return db.FindUser(id)
    })
}

func validateUser(user User) IOR.IOResult[User] {
    if user.Age < 18 {
        return IOR.Error[User](errors.New("user must be 18+"))
    }
    return IOR.Ok(user)
}

func saveUser(user User) IOR.IOResult[User] {
    return IOR.TryCatchError(func() (User, error) {
        return db.SaveUser(user)
    })
}

func processUser(id string) IOR.IOResult[User] {
    return F.Pipe3(
        fetchUser(id),
        IOR.Chain(validateUser),
        IOR.Chain(saveUser),
    )
}

func main() {
    result := processUser("user-123")()
    // Result[User]
}
`}
</CodeCard>

### HTTP Requests

<CodeCard file="http.go">
{`package main

import (
    "io"
    "net/http"
    IOR "github.com/IBM/fp-go/v2/ioresult"
)

func httpGet(url string) IOR.IOResult[[]byte] {
    return IOR.TryCatchError(func() ([]byte, error) {
        resp, err := http.Get(url)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != 200 {
            return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
        }
        
        return io.ReadAll(resp.Body)
    })
}

func main() {
    result := httpGet("https://api.example.com/data")()
    // Result[[]byte]
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
    IOR "github.com/IBM/fp-go/v2/ioresult"
)

func processFile(path string) IOR.IOResult[Data] {
    return IOR.Bracket(
        // Acquire
        func() IOR.IOResult[*os.File] {
            return IOR.TryCatchError(func() (*os.File, error) {
                return os.Open(path)
            })
        },
        // Use
        func(f *os.File) IOR.IOResult[Data] {
            return parseData(f)
        },
        // Release (always runs)
        func(f *os.File, _ IOR.IOResult[Data]) IOR.IOResult[unit.Unit] {
            f.Close()
            return IOR.Ok(unit.Unit{})
        },
    )
}
`}
</CodeCard>

### Pattern 2: Error Context

<CodeCard file="error_context.go">
{`package main

import (
    IOR "github.com/IBM/fp-go/v2/ioresult"
    F "github.com/IBM/fp-go/v2/function"
)

func processData(id string) IOR.IOResult[Data] {
    return F.Pipe3(
        fetchData(id),
        IOR.MapError(func(err error) error {
            return fmt.Errorf("fetch failed: %w", err)
        }),
        IOR.Chain(validateData),
        IOR.MapError(func(err error) error {
            return fmt.Errorf("validation failed: %w", err)
        }),
    )
}
`}
</CodeCard>

### Comparison: IOResult vs IOEither

<Compare>
<CompareCol kind="good">
<CodeCard file="ioresult.go">
{`// ✅ Use IOResult for standard errors
func readFile(path string) IOR.IOResult[[]byte] {
    return IOR.TryCatchError(os.ReadFile(path))
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="good">
<CodeCard file="ioeither.go">
{`// ✅ Use IOEither for custom errors
type FileError struct {
    Path string
    Code int
}

func readFile(path string) IOE.IOEither[FileError, []byte] {
    // Rich error information
}
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>
