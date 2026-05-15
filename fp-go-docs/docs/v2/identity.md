---
title: Identity
hide_title: true
description: The simplest monad - wraps a value with no effects. Perfect for teaching monad concepts and generic programming.
sidebar_position: 26
---

<PageHeader
  eyebrow="Reference · Core Type"
  title="Identity"
  lede="The simplest monad - just wraps a value with no effects. Identity is primarily used for teaching monad concepts, as a base case in generic programming, and for demonstrating monad laws."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/identity' },
    { label: 'Type', value: 'Monad' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

Identity is literally just the value itself - no wrapping, no effects, no context:
- **Simplest monad**: No additional structure
- **Teaching tool**: Perfect for understanding monad laws
- **Base case**: Used in generic programming

<CodeCard file="type_definition.go">
{`package identity

// Identity is just the value itself
type Identity[A any] = A
`}
</CodeCard>

### When to Use

<ApiTable>
| Use Case | Example |
|----------|---------|
| Teaching | Demonstrate monad laws |
| Generic programming | Base case for type classes |
| Testing | Simplest monad for unit tests |
| Type-level programming | Phantom types and constraints |
</ApiTable>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Of` | `func Of[A any](value A) Identity[A]` | Wrap value (no-op) |
| `MakeOf` | `func MakeOf[A any]() func(A) Identity[A]` | Constructor factory |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[A, B any](f func(A) B) func(Identity[A]) Identity[B]` | Transform value |
| `Chain` | `func Chain[A, B any](f func(A) Identity[B]) func(Identity[A]) Identity[B]` | FlatMap - sequence operations |
| `Ap` | `func Ap[A, B any](fa Identity[A]) func(Identity[func(A) B]) Identity[B]` | Apply wrapped function |
</ApiTable>

### Extraction

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Unwrap` | `func Unwrap[A any](id Identity[A]) A` | Extract value (no-op) |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Operations

<CodeCard file="basic.go">
{`package main

import (
    "fmt"
    I "github.com/IBM/fp-go/v2/identity"
)

func main() {
    // Create identity (just the value)
    id := I.Of(42)
    fmt.Println(id)  // 42
    
    // Map transforms the value
    doubled := I.Map(func(n int) int {
        return n * 2
    })(id)
    fmt.Println(doubled)  // 84
    
    // Chain sequences operations
    result := I.Chain(func(n int) I.Identity[string] {
        return I.Of(fmt.Sprintf("Value: %d", n))
    })(id)
    fmt.Println(result)  // "Value: 42"
    
    // Unwrap extracts the value (no-op)
    value := I.Unwrap(result)
    fmt.Println(value)  // "Value: 42"
}
`}
</CodeCard>

### Monad Laws

<CodeCard file="monad_laws.go">
{`package main

import (
    "fmt"
    I "github.com/IBM/fp-go/v2/identity"
)

func main() {
    // Left identity: Of(a).Chain(f) == f(a)
    a := 42
    f := func(n int) I.Identity[int] {
        return I.Of(n * 2)
    }
    
    left := I.Chain(f)(I.Of(a))
    right := f(a)
    fmt.Println(left == right)  // true
    
    // Right identity: m.Chain(Of) == m
    m := I.Of(42)
    result := I.Chain(I.Of[int])(m)
    fmt.Println(result == m)  // true
    
    // Associativity: m.Chain(f).Chain(g) == m.Chain(x => f(x).Chain(g))
    g := func(n int) I.Identity[int] {
        return I.Of(n + 10)
    }
    
    left2 := I.Chain(g)(I.Chain(f)(m))
    right2 := I.Chain(func(x int) I.Identity[int] {
        return I.Chain(g)(f(x))
    })(m)
    fmt.Println(left2 == right2)  // true
}
`}
</CodeCard>

### Functor Laws

<CodeCard file="functor_laws.go">
{`package main

import (
    "fmt"
    I "github.com/IBM/fp-go/v2/identity"
)

func main() {
    // Identity law: Map(id) == id
    id := func(x int) int { return x }
    value := I.Of(42)
    
    mapped := I.Map(id)(value)
    fmt.Println(mapped == value)  // true
    
    // Composition law: Map(f ∘ g) == Map(f) ∘ Map(g)
    f := func(n int) int { return n * 2 }
    g := func(n int) int { return n + 10 }
    
    // Map(f ∘ g)
    left := I.Map(func(x int) int {
        return f(g(x))
    })(value)
    
    // Map(f) ∘ Map(g)
    right := I.Map(f)(I.Map(g)(value))
    
    fmt.Println(left == right)  // true
}
`}
</CodeCard>

### Generic Programming

<CodeCard file="generic.go">
{`package main

import (
    "fmt"
    I "github.com/IBM/fp-go/v2/identity"
)

// Generic function that works with any monad
func Transform[M any, A, B any](
    mapFn func(func(A) B) func(M) M,
    f func(A) B,
) func(M) M {
    return mapFn(f)
}

func main() {
    // Use with Identity
    double := func(n int) int { return n * 2 }
    
    transform := Transform[I.Identity[int], int, int](
        I.Map[int, int],
        double,
    )
    
    result := transform(I.Of(21))
    fmt.Println(result)  // 42
}
`}
</CodeCard>

### Testing Helper

<CodeCard file="testing.go">
{`package main

import (
    "testing"
    I "github.com/IBM/fp-go/v2/identity"
)

// Test monad behavior with Identity (simplest case)
func TestMonadBehavior(t *testing.T) {
    // Test Map
    value := I.Of(10)
    doubled := I.Map(func(n int) int { return n * 2 })(value)
    
    if doubled != 20 {
        t.Errorf("Expected 20, got %d", doubled)
    }
    
    // Test Chain
    result := I.Chain(func(n int) I.Identity[int] {
        return I.Of(n + 5)
    })(value)
    
    if result != 15 {
        t.Errorf("Expected 15, got %d", result)
    }
}

// Use Identity to test pure functions
func TestPureFunction(t *testing.T) {
    add := func(a, b int) int { return a + b }
    
    // Wrap in Identity for testing
    result := I.Map(func(x int) int {
        return add(x, 10)
    })(I.Of(5))
    
    if result != 15 {
        t.Errorf("Expected 15, got %d", result)
    }
}
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Pattern 1: Type Class Base Case

<CodeCard file="typeclass.go">
{`package main

import (
    I "github.com/IBM/fp-go/v2/identity"
)

// Functor type class
type Functor[F any] interface {
    Map(f func(any) any) func(F) F
}

// Identity implements Functor (trivially)
type IdentityFunctor[A any] struct{}

func (IdentityFunctor[A]) Map(f func(A) A) func(I.Identity[A]) I.Identity[A] {
    return I.Map[A, A](f)
}
`}
</CodeCard>

### Pattern 2: Teaching Tool

<CodeCard file="teaching.go">
{`package main

import (
    "fmt"
    I "github.com/IBM/fp-go/v2/identity"
)

// Demonstrate that Identity is transparent
func demonstrateTransparency() {
    // These are equivalent
    value1 := 42
    value2 := I.Of(42)
    
    fmt.Println(value1 == value2)  // true
    
    // Map is just function application
    f := func(n int) int { return n * 2 }
    
    result1 := f(value1)
    result2 := I.Map(f)(value2)
    
    fmt.Println(result1 == result2)  // true
}
`}
</CodeCard>

### Pattern 3: Composition Study

<CodeCard file="composition.go">
{`package main

import (
    "fmt"
    I "github.com/IBM/fp-go/v2/identity"
    F "github.com/IBM/fp-go/v2/function"
)

func main() {
    // Study function composition with Identity
    f := func(n int) int { return n * 2 }
    g := func(n int) int { return n + 10 }
    h := func(n int) int { return n - 5 }
    
    // Compose functions
    composed := F.Pipe3(
        I.Of(5),
        I.Map(f),
        I.Map(g),
        I.Map(h),
    )
    
    fmt.Println(composed)  // ((5 * 2) + 10) - 5 = 15
}
`}
</CodeCard>

</Section>

<Callout type="info">

**Why Identity?** While Identity seems trivial (it's just the value!), it's invaluable for:
- **Teaching**: Understanding monad laws without complexity
- **Testing**: Verifying generic monad code with the simplest case
- **Type theory**: Serving as the identity element in category theory

</Callout>
