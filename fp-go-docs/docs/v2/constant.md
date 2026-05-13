---
title: Constant
hide_title: true
description: Constant functor for advanced functional patterns - always returns the same value, ignoring transformations.
sidebar_position: 27
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="Constant"
  lede="Constant functor for advanced functional patterns. Always returns the same value, ignoring transformations. Used primarily with optics and traversals."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/constant' },
    { label: 'Type', value: 'Functor' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

Constant is a functor that ignores transformations and always holds the same value:
- **Phantom type**: Type parameter `A` is not used
- **Immutable**: Value never changes
- **Optics**: Essential for traversals and folds

<CodeCard file="type_definition.go">
{`package constant

// Constant ignores type parameter A
type Constant[C, A any] = C
`}
</CodeCard>

### When to Use

<ApiTable>
| Use Case | Example |
|----------|---------|
| Optics | Traversals, folds, getters |
| Accumulation | Collecting values during traversal |
| Phantom types | Type-level programming |
| Advanced patterns | Applicative functors |
</ApiTable>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Of` | `func Of[C, A any](value C) Constant[C, A]` | Create constant with phantom type |
| `MakeOf` | `func MakeOf[C, A any]() func(C) Constant[C, A]` | Constructor factory |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[C, A, B any](f func(A) B) func(Constant[C, A]) Constant[C, B]` | No-op transformation (value unchanged) |
| `Ap` | `func Ap[C, A, B any](fa Constant[C, A]) func(Constant[C, func(A) B]) Constant[C, B]` | Applicative apply (combines C values) |
</ApiTable>

### Extraction

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Unwrap` | `func Unwrap[C, A any](c Constant[C, A]) C` | Extract the constant value |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    "fmt"
    C "github.com/IBM/fp-go/v2/constant"
)

func main() {
    // Create constant with phantom type
    c := C.Of[string, int]("hello")
    // c is "hello", int type is phantom
    
    // Map does nothing (constant stays the same)
    mapped := C.Map[string, int, bool](func(n int) bool {
        return n > 0
    })(c)
    
    // Value is still "hello"
    fmt.Println(C.Unwrap(mapped)) // "hello"
}
`}
</CodeCard>

### With Monoids

<CodeCard file="monoid.go">
{`package main

import (
    C "github.com/IBM/fp-go/v2/constant"
    M "github.com/IBM/fp-go/v2/monoid"
)

func main() {
    // Combine constants using monoid
    c1 := C.Of[string, int]("hello")
    c2 := C.Of[string, int](" world")
    
    // Use string monoid to combine
    combined := M.Concat(M.String)(
        C.Unwrap(c1),
        C.Unwrap(c2),
    )
    
    fmt.Println(combined) // "hello world"
}
`}
</CodeCard>

### Accumulation Pattern

<CodeCard file="accumulation.go">
{`package main

import (
    C "github.com/IBM/fp-go/v2/constant"
    M "github.com/IBM/fp-go/v2/monoid"
)

type User struct {
    Name  string
    Email string
    Age   int
}

// Collect all string fields
func collectStrings(u User) C.Constant[[]string, User] {
    return C.Of[[]string, User]([]string{u.Name, u.Email})
}

func main() {
    users := []User{
        {Name: "Alice", Email: "alice@example.com", Age: 30},
        {Name: "Bob", Email: "bob@example.com", Age: 25},
    }
    
    // Accumulate all strings
    var allStrings []string
    for _, u := range users {
        c := collectStrings(u)
        allStrings = append(allStrings, C.Unwrap(c)...)
    }
    
    fmt.Println(allStrings)
    // ["Alice", "alice@example.com", "Bob", "bob@example.com"]
}
`}
</CodeCard>

### With Optics

<CodeCard file="optics.go">
{`package main

import (
    C "github.com/IBM/fp-go/v2/constant"
    M "github.com/IBM/fp-go/v2/monoid"
)

type Config struct {
    Host string
    Port int
    SSL  bool
}

// Fold over structure, collecting values
func foldConfig(cfg Config) C.Constant[string, Config] {
    summary := fmt.Sprintf("Host: %s, Port: %d", cfg.Host, cfg.Port)
    return C.Of[string, Config](summary)
}

func main() {
    cfg := Config{Host: "localhost", Port: 8080, SSL: true}
    
    result := foldConfig(cfg)
    fmt.Println(C.Unwrap(result))
    // "Host: localhost, Port: 8080"
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern 1: Traversal Accumulation

<CodeCard file="traversal.go">
{`package main

import (
    C "github.com/IBM/fp-go/v2/constant"
    M "github.com/IBM/fp-go/v2/monoid"
)

// Traverse and accumulate
func sumFields(data []int) C.Constant[int, []int] {
    sum := 0
    for _, n := range data {
        sum += n
    }
    return C.Of[int, []int](sum)
}

func main() {
    data := []int{1, 2, 3, 4, 5}
    result := sumFields(data)
    
    fmt.Println(C.Unwrap(result)) // 15
}
`}
</CodeCard>

### Pattern 2: Type-Safe Getters

<CodeCard file="getters.go">
{`package main

import (
    C "github.com/IBM/fp-go/v2/constant"
)

type Person struct {
    Name string
    Age  int
}

// Getter using Constant
func getName(p Person) C.Constant[string, Person] {
    return C.Of[string, Person](p.Name)
}

func getAge(p Person) C.Constant[int, Person] {
    return C.Of[int, Person](p.Age)
}

func main() {
    person := Person{Name: "Alice", Age: 30}
    
    name := C.Unwrap(getName(person))
    age := C.Unwrap(getAge(person))
    
    fmt.Printf("%s is %d years old\n", name, age)
}
`}
</CodeCard>

### Pattern 3: Validation Accumulation

<CodeCard file="validation.go">
{`package main

import (
    C "github.com/IBM/fp-go/v2/constant"
)

type ValidationError struct {
    Field   string
    Message string
}

func validateUser(u User) C.Constant[[]ValidationError, User] {
    var errors []ValidationError
    
    if u.Name == "" {
        errors = append(errors, ValidationError{
            Field:   "name",
            Message: "name is required",
        })
    }
    
    if u.Age < 18 {
        errors = append(errors, ValidationError{
            Field:   "age",
            Message: "must be 18 or older",
        })
    }
    
    return C.Of[[]ValidationError, User](errors)
}
`}
</CodeCard>

</Section>

<Callout type="info">

**Advanced Usage**: Constant is primarily used in advanced functional programming patterns, particularly with optics (lenses, prisms, traversals). For most use cases, consider simpler types like Option, Result, or Either.

</Callout>
