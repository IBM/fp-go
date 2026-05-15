---
title: Either
hide_title: true
description: Generic disjoint union type with Left and Right - for custom error types and domain modeling.
sidebar_position: 11
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="Either"
  lede="Generic disjoint union representing a value of one of two possible types. Either[E, A] can be Left (typically error) or Right (typically success)."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/either' },
    { label: 'Type', value: 'Monad' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

Either represents a choice between two types:
- **Left**: Value of type `E` (typically error or failure)
- **Right**: Value of type `A` (typically success)

Unlike Result[A] which fixes the error type to `error`, Either[E, A] is generic over both types.

### When to Use

<ApiTable>
| Use Either When | Use Result When |
|-----------------|-----------------|
| Custom error types with rich information | Standard Go error handling |
| Domain modeling with sum types | Library interoperability |
| Non-error alternatives (cache miss, etc.) | Simpler type signatures |
| Type-level distinction between error types | Working with existing Go code |
</ApiTable>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Left` | `func Left[A, E any](value E) Either[E, A]` | Create Left value (failure) |
| `Right` | `func Right[E, A any](value A) Either[E, A]` | Create Right value (success) |
| `Of` | `func Of[E, A any](value A) Either[E, A]` | Alias for Right (monadic pure) |
| `TryCatch` | `func TryCatch[E, A any](f func() (A, E)) Either[E, A]` | From function returning tuple |
| `FromPredicate` | `func FromPredicate[E, A any](pred func(A) bool, onFalse func(A) E) func(A) Either[E, A]` | Conditional Either |
</ApiTable>

### Pattern Matching

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Fold` | `func Fold[E, A, B any](onLeft func(E) B, onRight func(A) B) func(Either[E, A]) B` | Extract value from both cases |
| `Match` | `func Match[E, A, B any](onLeft func(E) Either[E, B], onRight func(A) Either[E, B]) func(Either[E, A]) Either[E, B]` | Pattern match with Either return |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[E, A, B any](f func(A) B) func(Either[E, A]) Either[E, B]` | Transform Right value |
| `MapLeft` | `func MapLeft[A, E1, E2 any](f func(E1) E2) func(Either[E1, A]) Either[E2, A]` | Transform Left value |
| `BiMap` | `func BiMap[E1, E2, A, B any](fe func(E1) E2, fa func(A) B) func(Either[E1, A]) Either[E2, B]` | Transform both sides |
| `Chain` | `func Chain[E, A, B any](f func(A) Either[E, B]) func(Either[E, A]) Either[E, B]` | FlatMap - sequence operations |
</ApiTable>

### Combining

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Ap` | `func Ap[E, A, B any](fa Either[E, A]) func(Either[E, func(A) B]) Either[E, B]` | Apply wrapped function |
| `SequenceArray` | `func SequenceArray[E, A any]([]Either[E, A]) Either[E, []A]` | All-or-nothing for arrays |
| `TraverseArray` | `func TraverseArray[E, A, B any](f func(A) Either[E, B]) func([]A) Either[E, []B]` | Map and sequence |
</ApiTable>

### Extraction

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Unwrap` | `func Unwrap[A any](Either[error, A]) (A, error)` | Convert to tuple (only for Either[error, A]) |
| `GetOrElse` | `func GetOrElse[E, A any](f func() A) func(Either[E, A]) A` | Extract with default |
| `ToOption` | `func ToOption[E, A any](Either[E, A]) Option[A]` | Convert to Option (discards Left) |
</ApiTable>

### Testing

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `IsLeft` | `func IsLeft[E, A any](Either[E, A]) bool` | Test for Left |
| `IsRight` | `func IsRight[E, A any](Either[E, A]) bool` | Test for Right |
| `Exists` | `func Exists[E, A any](pred func(A) bool) func(Either[E, A]) bool` | Test Right value |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    "fmt"
    E "github.com/IBM/fp-go/v2/either"
)

func main() {
    // Create Either values
    success := E.Right[string](42)
    failure := E.Left[int]("error message")
    
    // Check values
    fmt.Println(E.IsRight(success)) // true
    fmt.Println(E.IsLeft(failure))  // true
    
    // Pattern match
    result := E.Fold(
        func(err string) string { return "Error: " + err },
        func(n int) string { return fmt.Sprintf("Value: %d", n) },
    )(success)
    fmt.Println(result) // "Value: 42"
}
`}
</CodeCard>

### Custom Error Types

<CodeCard file="custom_errors.go">
{`package main

import (
    "time"
    E "github.com/IBM/fp-go/v2/either"
)

type ValidationError struct {
    Field   string
    Message string
}

type AppError struct {
    Code      int
    Message   string
    Timestamp time.Time
}

func validateEmail(email string) E.Either[ValidationError, string] {
    if !strings.Contains(email, "@") {
        return E.Left[string](ValidationError{
            Field:   "email",
            Message: "must contain @",
        })
    }
    return E.Right[ValidationError](email)
}

func fetchUser(id string) E.Either[AppError, User] {
    if id == "" {
        return E.Left[User](AppError{
            Code:      400,
            Message:   "Invalid user ID",
            Timestamp: time.Now(),
        })
    }
    // ... fetch logic
    return E.Right[AppError](user)
}
`}
</CodeCard>

### Transformations

<CodeCard file="transformations.go">
{`package main

import (
    E "github.com/IBM/fp-go/v2/either"
    F "github.com/IBM/fp-go/v2/function"
)

func main() {
    // Map: transform Right value
    doubled := F.Pipe1(
        E.Right[string](21),
        E.Map(func(n int) int { return n * 2 }),
    )
    // Right(42)
    
    // MapLeft: transform Left value
    withContext := F.Pipe1(
        E.Left[int]("error"),
        E.MapLeft(func(msg string) AppError {
            return AppError{Code: 500, Message: msg}
        }),
    )
    // Left(AppError{Code: 500, Message: "error"})
    
    // Chain: sequence operations
    result := F.Pipe2(
        parseInput("42"),
        E.Chain(validateInput),
        E.Chain(processInput),
    )
}
`}
</CodeCard>

### Domain Modeling

<CodeCard file="domain_modeling.go">
{`package main

import (
    E "github.com/IBM/fp-go/v2/either"
)

// Payment methods as sum type
type CreditCard struct {
    Number string
    CVV    string
}

type BankTransfer struct {
    AccountNumber string
    RoutingNumber string
}

type PaymentMethod = E.Either[CreditCard, BankTransfer]

func processPayment(method PaymentMethod, amount float64) E.Either[string, Receipt] {
    return E.Fold(
        func(cc CreditCard) E.Either[string, Receipt] {
            return processCreditCard(cc, amount)
        },
        func(bt BankTransfer) E.Either[string, Receipt] {
            return processBankTransfer(bt, amount)
        },
    )(method)
}

// Cache results
type CacheMiss struct {
    Key    string
    Reason string
}

type CacheResult[A any] = E.Either[CacheMiss, A]

func getFromCache[A any](key string) CacheResult[A] {
    if val, ok := cache.Get(key); ok {
        return E.Right[CacheMiss](val.(A))
    }
    return E.Left[A](CacheMiss{
        Key:    key,
        Reason: "not found",
    })
}
`}
</CodeCard>

### Sequence Operations

<CodeCard file="sequence.go">
{`package main

import (
    E "github.com/IBM/fp-go/v2/either"
)

func main() {
    // Parse all - fails on first error
    parseAll := E.TraverseArray(func(s string) E.Either[string, int] {
        n, err := strconv.Atoi(s)
        if err != nil {
            return E.Left[int](err.Error())
        }
        return E.Right[string](n)
    })
    
    result := parseAll([]string{"1", "2", "3"})
    // Right([1, 2, 3])
    
    result = parseAll([]string{"1", "bad", "3"})
    // Left("invalid syntax")
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern 1: Error Type Conversion

<CodeCard file="error_conversion.go">
{`package main

import (
    E "github.com/IBM/fp-go/v2/either"
    F "github.com/IBM/fp-go/v2/function"
)

type DBError struct{ Msg string }
type AppError struct{ Code int; Msg string }

func processOrder(id string) E.Either[AppError, Order] {
    return F.Pipe3(
        fetchOrder(id),  // Either[DBError, Order]
        E.MapLeft(func(e DBError) AppError {
            return AppError{Code: 500, Msg: e.Msg}
        }),
        E.Chain(validateOrder),  // Either[AppError, Order]
        E.Chain(saveOrder),
    )
}
`}
</CodeCard>

### Pattern 2: Validation with Multiple Errors

<CodeCard file="validation.go">
{`package main

import (
    "strings"
    E "github.com/IBM/fp-go/v2/either"
)

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
    var msgs []string
    for _, e := range ve {
        msgs = append(msgs, fmt.Sprintf("%s: %s", e.Field, e.Message))
    }
    return strings.Join(msgs, "; ")
}

func validateUser(user User) E.Either[ValidationErrors, User] {
    var errors ValidationErrors
    
    if user.Name == "" {
        errors = append(errors, ValidationError{
            Field:   "name",
            Message: "required",
        })
    }
    
    if !strings.Contains(user.Email, "@") {
        errors = append(errors, ValidationError{
            Field:   "email",
            Message: "invalid format",
        })
    }
    
    if len(errors) > 0 {
        return E.Left[User](errors)
    }
    
    return E.Right[ValidationErrors](user)
}
`}
</CodeCard>

### Comparison: Either vs Result

<Compare>
<CompareCol kind="bad">
<CodeCard file="either_for_errors.go">
{`// ❌ Don't use Either[error, A] for standard errors
func fetchUser(id string) E.Either[error, User] {
    // Use Result[User] instead
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="good">
<CodeCard file="either_for_custom.go">
{`// ✅ Use Either for custom error types
func fetchUser(id string) E.Either[AppError, User] {
    // Rich error information
}

// ✅ Use Result for standard errors
func fetchUser(id string) R.Result[User] {
    // Simpler, idiomatic Go
}
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>
