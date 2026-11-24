---
title: Option
hide_title: true
description: Type-safe optional values with Some and None - eliminate nil pointer errors.
sidebar_position: 3
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="Option"
  lede="Represent optional values without nil pointers. Option is either Some(value) or None, providing type-safe handling of missing values."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/option' },
    { label: 'Type', value: 'Monad' }
  ]}
/>

<Section id="api" number="01" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Some` | `func Some[A any](value A) Option[A]` | Create an Option containing a value |
| `None` | `func None[A any]() Option[A]` | Create an empty Option |
| `Of` | `func Of[A any](value A) Option[A]` | Alias for Some (monadic pure) |
| `FromNillable` | `func FromNillable[A any](ptr *A) Option[A]` | Convert pointer to Option (nil → None) |
| `FromPredicate` | `func FromPredicate[A any](pred func(A) bool) func(A) Option[A]` | Create Option based on predicate |
</ApiTable>

### Predicates

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `IsSome` | `func IsSome[A any](opt Option[A]) bool` | Check if Option contains a value |
| `IsNone` | `func IsNone[A any](opt Option[A]) bool` | Check if Option is empty |
</ApiTable>

### Extractors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `GetOrElse` | `func GetOrElse[A any](defaultValue func() A) func(Option[A]) A` | Extract value or return default |
| `ToNillable` | `func ToNillable[A any](opt Option[A]) *A` | Convert to pointer (None → nil) |
| `Match` | `func Match[A, B any](onNone func() B, onSome func(A) B) func(Option[A]) B` | Pattern match on Option |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[A, B any](f func(A) B) func(Option[A]) Option[B]` | Transform the value if present |
| `Chain` | `func Chain[A, B any](f func(A) Option[B]) func(Option[A]) Option[B]` | FlatMap - chain optional operations |
| `Ap` | `func Ap[A, B any](fa Option[A]) func(Option[func(A) B]) Option[B]` | Apply wrapped function to wrapped value |
| `Flatten` | `func Flatten[A any](opt Option[Option[A]]) Option[A]` | Flatten nested Options |
</ApiTable>

### Filtering

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Filter` | `func Filter[A any](pred func(A) bool) func(Option[A]) Option[A]` | Keep value only if predicate holds |
| `FilterMap` | `func FilterMap[A, B any](f func(A) Option[B]) func(Option[A]) Option[B]` | Map and filter in one operation |
</ApiTable>

### Combinators

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Alt` | `func Alt[A any](second Option[A]) func(Option[A]) Option[A]` | Return first Some, or second if first is None |
| `OrElse` | `func OrElse[A any](alternative func() Option[A]) func(Option[A]) Option[A]` | Lazy alternative |
</ApiTable>

</Section>

<Section id="examples" number="02" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    "fmt"
    O "github.com/IBM/fp-go/v2/option"
)

func main() {
    // Create Options
    some := O.Some(42)
    none := O.None[int]()
    
    // Check for values
    fmt.Println(O.IsSome(some)) // true
    fmt.Println(O.IsNone(none)) // true
    
    // Extract values
    value := O.GetOrElse(func() int { return 0 })(some)
    fmt.Println(value) // 42
    
    defaultValue := O.GetOrElse(func() int { return 0 })(none)
    fmt.Println(defaultValue) // 0
}
`}
</CodeCard>

### Transformations

<CodeCard file="transformations.go">
{`package main

import (
    "fmt"
    "strings"
    O "github.com/IBM/fp-go/v2/option"
    F "github.com/IBM/fp-go/v2/function"
)

func main() {
    // Map: transform the value
    opt := O.Some("hello")
    upper := F.Pipe1(
        opt,
        O.Map(strings.ToUpper),
    )
    fmt.Println(O.GetOrElse(func() string { return "" })(upper)) // HELLO
    
    // Chain: flatMap operations
    divide := func(a, b int) O.Option[int] {
        if b == 0 {
            return O.None[int]()
        }
        return O.Some(a / b)
    }
    
    result := F.Pipe1(
        O.Some(10),
        O.Chain(func(n int) O.Option[int] {
            return divide(n, 2)
        }),
    )
    fmt.Println(O.GetOrElse(func() int { return 0 })(result)) // 5
    
    // Filter: keep value only if predicate holds
    filtered := F.Pipe1(
        O.Some(42),
        O.Filter(func(n int) bool { return n > 50 }),
    )
    fmt.Println(O.IsNone(filtered)) // true (42 is not > 50)
}
`}
</CodeCard>

### Pattern Matching

<CodeCard file="pattern_matching.go">
{`package main

import (
    "fmt"
    O "github.com/IBM/fp-go/v2/option"
)

func main() {
    opt := O.Some(42)
    
    // Match: handle both cases
    result := O.Match(
        func() string { return "No value" },
        func(n int) string { return fmt.Sprintf("Value: %d", n) },
    )(opt)
    
    fmt.Println(result) // Value: 42
}
`}
</CodeCard>

### Practical Example: Configuration

<CodeCard file="config.go">
{`package main

import (
    "fmt"
    O "github.com/IBM/fp-go/v2/option"
    F "github.com/IBM/fp-go/v2/function"
)

type Config struct {
    Port    O.Option[int]
    Host    O.Option[string]
    Timeout O.Option[int]
}

func getPort(config Config) int {
    return F.Pipe1(
        config.Port,
        O.GetOrElse(func() int { return 8080 }),
    )
}

func getHost(config Config) string {
    return F.Pipe1(
        config.Host,
        O.GetOrElse(func() string { return "localhost" }),
    )
}

func main() {
    // Config with explicit values
    config1 := Config{
        Port: O.Some(3000),
        Host: O.Some("0.0.0.0"),
        Timeout: O.Some(30),
    }
    fmt.Printf("Server: %s:%d\n", getHost(config1), getPort(config1))
    // Server: 0.0.0.0:3000
    
    // Config with defaults
    config2 := Config{
        Port: O.None[int](),
        Host: O.None[string](),
        Timeout: O.Some(60),
    }
    fmt.Printf("Server: %s:%d\n", getHost(config2), getPort(config2))
    // Server: localhost:8080
}
`}
</CodeCard>

### Working with Pointers

<CodeCard file="pointers.go">
{`package main

import (
    "fmt"
    O "github.com/IBM/fp-go/v2/option"
)

func findUser(id string) *User {
    if id == "123" {
        return &User{ID: "123", Name: "Alice"}
    }
    return nil
}

type User struct {
    ID   string
    Name string
}

func main() {
    // Convert pointer to Option
    userPtr := findUser("123")
    opt := O.FromNillable(userPtr)
    
    fmt.Println(O.IsSome(opt)) // true
    
    // Convert back to pointer
    ptr := O.ToNillable(opt)
    if ptr != nil {
        fmt.Println(ptr.Name) // Alice
    }
    
    // Handle missing user
    missingPtr := findUser("999")
    missingOpt := O.FromNillable(missingPtr)
    fmt.Println(O.IsNone(missingOpt)) // true
}
`}
</CodeCard>

</Section>

<Section id="why-option" number="03" title="Why Use" titleAccent="Option?">

<Checklist>
  <ChecklistItem status="required">
    **Type Safety** — Compiler ensures you handle both Some and None cases
  </ChecklistItem>
  <ChecklistItem status="required">
    **No Nil Panics** — Eliminates nil pointer dereference errors
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Explicit Intent** — Makes it clear when a value might be absent
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Composability** — Works seamlessly with pipe, map, chain, and other fp-go functions
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Functional Style** — Encourages declarative, expression-based code
  </ChecklistItem>
</Checklist>

<Compare>
<CompareCol kind="bad">
<CodeCard file="traditional.go">
{`// ❌ Traditional nil handling
func getUser(id string) *User {
    // Might return nil
    return db.FindUser(id)
}

user := getUser("123")
if user != nil {
    // Easy to forget nil check
    fmt.Println(user.Name) // Potential panic!
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="good">
<CodeCard file="option_based.go">
{`// ✅ Option-based handling
func getUser(id string) O.Option[User] {
    user := db.FindUser(id)
    return O.FromNillable(user)
}

// Compiler forces you to handle both cases
O.Match(
    func() { fmt.Println("User not found") },
    func(u User) { fmt.Println(u.Name) },
)(getUser("123"))
`}
</CodeCard>
</CompareCol>
</Compare>

</Section>
