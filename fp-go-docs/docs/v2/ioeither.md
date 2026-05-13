---
title: IOEither
hide_title: true
description: Lazy side effects with custom error handling - combine IO and Either for effectful computations.
sidebar_position: 14
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="IOEither"
  lede="Combine lazy evaluation (IO) with custom error handling (Either). IOEither[E, A] represents a synchronous computation with side effects that may fail with custom error type E."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/ioeither' },
    { label: 'Type', value: 'Monad (IO[Either[E, A]])' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

IOEither combines two powerful abstractions:

<CodeCard file="type_definition.go">
{`package ioeither

// IOEither is IO of Either
type IOEither[E, A any] = IO[Either[E, A]]
// Which expands to: func() Either[E, A]
`}
</CodeCard>

### When to Use

<ApiTable>
| Use IOEither When | Use IOResult When |
|-------------------|-------------------|
| Custom error types with rich information | Standard Go `error` interface |
| Domain-specific error handling | Library interoperability |
| Type-level error distinction | Simpler type signatures |
| Migration from v1 (used IOEither extensively) | Working with existing Go code |
</ApiTable>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Right` | `func Right[E, A any](value A) IOEither[E, A]` | Create successful IOEither |
| `Left` | `func Left[E, A any](err E) IOEither[E, A]` | Create failed IOEither |
| `Of` | `func Of[E, A any](value A) IOEither[E, A]` | Alias for Right |
| `TryCatch` | `func TryCatch[E, A any](f func() (A, E)) IOEither[E, A]` | From function returning tuple |
| `FromIO` | `func FromIO[E, A any](io IO[A]) IOEither[E, A]` | Lift IO to IOEither |
| `FromEither` | `func FromEither[E, A any](e Either[E, A]) IOEither[E, A]` | Lift Either to IOEither |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[E, A, B any](f func(A) B) func(IOEither[E, A]) IOEither[E, B]` | Transform success value |
| `MapLeft` | `func MapLeft[A, E1, E2 any](f func(E1) E2) func(IOEither[E1, A]) IOEither[E2, A]` | Transform error value |
| `BiMap` | `func BiMap[E1, E2, A, B any](fe func(E1) E2, fa func(A) B) func(IOEither[E1, A]) IOEither[E2, B]` | Transform both sides |
| `Chain` | `func Chain[E, A, B any](f func(A) IOEither[E, B]) func(IOEither[E, A]) IOEither[E, B]` | Sequence operations |
</ApiTable>

### Combining

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ap` | `func Ap[E, A, B any](fa IOEither[E, A]) func(IOEither[E, func(A) B]) IOEither[E, B]` | Apply wrapped function |
| `SequenceArray` | `func SequenceArray[E, A any]([]IOEither[E, A]) IOEither[E, []A]` | All-or-nothing |
| `TraverseArray` | `func TraverseArray[E, A, B any](f func(A) IOEither[E, B]) func([]A) IOEither[E, []B]` | Map and sequence |
</ApiTable>

### Pattern Matching

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Fold` | `func Fold[E, A, B any](onLeft func(E) IO[B], onRight func(A) IO[B]) func(IOEither[E, A]) IO[B]` | Extract to IO |
| `Match` | `func Match[E, A, B any](onLeft func(E) IOEither[E, B], onRight func(A) IOEither[E, B]) func(IOEither[E, A]) IOEither[E, B]` | Pattern match |
</ApiTable>

### Utilities

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `ChainFirst` | `func ChainFirst[E, A, B any](f func(A) IOEither[E, B]) func(IOEither[E, A]) IOEither[E, A]` | Side effect, keep value |
| `Alt` | `func Alt[E, A any](second IOEither[E, A]) func(IOEither[E, A]) IOEither[E, A]` | Fallback on Left |
| `OrElse` | `func OrElse[E, A any](f func(E) IOEither[E, A]) func(IOEither[E, A]) IOEither[E, A]` | Lazy fallback |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    "fmt"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

type AppError struct {
    Code    int
    Message string
}

func main() {
    // Create IOEither values
    success := IOE.Right[AppError](42)
    failure := IOE.Left[int](AppError{
        Code:    404,
        Message: "Not found",
    })
    
    // Execute
    result := success()  // Either[AppError, int] = Right(42)
    result = failure()   // Either[AppError, int] = Left(AppError{...})
}
`}
</CodeCard>

### File Operations

<CodeCard file="file_ops.go">
{`package main

import (
    "os"
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

type FileError struct {
    Path    string
    Message string
}

func readConfig(path string) IOE.IOEither[FileError, Config] {
    return IOE.TryCatch(func() (Config, FileError) {
        data, err := os.ReadFile(path)
        if err != nil {
            return Config{}, FileError{
                Path:    path,
                Message: err.Error(),
            }
        }
        
        var cfg Config
        if err := json.Unmarshal(data, &cfg); err != nil {
            return Config{}, FileError{
                Path:    path,
                Message: "invalid JSON: " + err.Error(),
            }
        }
        
        return cfg, FileError{}
    })
}

func main() {
    result := readConfig("config.json")()
    // Either[FileError, Config]
}
`}
</CodeCard>

### Chaining Operations

<CodeCard file="chaining.go">
{`package main

import (
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

type ValidationError struct {
    Field   string
    Message string
}

func fetchUser(id string) IOE.IOEither[ValidationError, User] {
    return IOE.TryCatch(func() (User, ValidationError) {
        // Fetch from database
        user, err := db.FindUser(id)
        if err != nil {
            return User{}, ValidationError{
                Field:   "id",
                Message: err.Error(),
            }
        }
        return user, ValidationError{}
    })
}

func validateUser(user User) IOE.IOEither[ValidationError, User] {
    if user.Age < 18 {
        return IOE.Left[User](ValidationError{
            Field:   "age",
            Message: "must be 18 or older",
        })
    }
    return IOE.Right[ValidationError](user)
}

func saveUser(user User) IOE.IOEither[ValidationError, User] {
    return IOE.TryCatch(func() (User, ValidationError) {
        err := db.SaveUser(user)
        if err != nil {
            return User{}, ValidationError{
                Field:   "save",
                Message: err.Error(),
            }
        }
        return user, ValidationError{}
    })
}

func processUser(id string) IOE.IOEither[ValidationError, User] {
    return F.Pipe3(
        fetchUser(id),
        IOE.Chain(validateUser),
        IOE.Chain(saveUser),
    )
}

func main() {
    result := processUser("user-123")()
    // Either[ValidationError, User]
}
`}
</CodeCard>

### Error Context

<CodeCard file="error_context.go">
{`package main

import (
    IOE "github.com/IBM/fp-go/v2/ioeither"
    F "github.com/IBM/fp-go/v2/function"
)

func enrichError(operation string) func(AppError) AppError {
    return func(err AppError) AppError {
        return AppError{
            Code:    err.Code,
            Message: fmt.Sprintf("%s: %s", operation, err.Message),
        }
    }
}

func processData(id string) IOE.IOEither[AppError, Data] {
    return F.Pipe3(
        fetchData(id),
        IOE.MapLeft(enrichError("fetch")),
        IOE.Chain(validateData),
        IOE.MapLeft(enrichError("validate")),
    )
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
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func processFile(path string) IOE.IOEither[AppError, Data] {
    return IOE.Bracket(
        // Acquire
        func() IOE.IOEither[AppError, *os.File] {
            return IOE.TryCatch(func() (*os.File, AppError) {
                f, err := os.Open(path)
                if err != nil {
                    return nil, AppError{Code: 500, Message: err.Error()}
                }
                return f, AppError{}
            })
        },
        // Use
        func(f *os.File) IOE.IOEither[AppError, Data] {
            return parseData(f)
        },
        // Release (always runs)
        func(f *os.File, _ IOE.IOEither[AppError, Data]) IOE.IOEither[AppError, unit.Unit] {
            f.Close()
            return IOE.Right[AppError](unit.Unit{})
        },
    )
}
`}
</CodeCard>

### Pattern 2: Fallback on Error

<CodeCard file="fallback.go">
{`package main

import (
    IOE "github.com/IBM/fp-go/v2/ioeither"
)

func fetchWithFallback(id string) IOE.IOEither[AppError, Data] {
    return F.Pipe1(
        fetchFromCache(id),
        IOE.OrElse(func(err AppError) IOE.IOEither[AppError, Data] {
            // Fallback to database on cache miss
            return fetchFromDB(id)
        }),
    )
}
`}
</CodeCard>

### Comparison: IOEither vs IOResult

<Compare>
<CompareCol kind="bad">
<CodeCard file="ioeither_for_error.go">
{`// ❌ Don't use IOEither[error, A]
func readFile(path string) IOE.IOEither[error, []byte] {
    // Use IOResult instead
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="good">
<CodeCard file="ioeither_for_custom.go">
{`// ✅ Use IOEither for custom errors
func readFile(path string) IOE.IOEither[FileError, []byte] {
    // Rich error information
}

// ✅ Use IOResult for standard errors
func readFile(path string) IOR.IOResult[[]byte] {
    // Simpler, idiomatic
}
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>
