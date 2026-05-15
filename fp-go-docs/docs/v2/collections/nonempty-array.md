---
title: NonEmpty Array
hide_title: true
description: Arrays guaranteed to have at least one element, eliminating empty array edge cases.
sidebar_position: 5
---

import { PageHeader, Section, CodeCard, ApiTable, Callout, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="Reference · Collections"
  title="NonEmpty Array"
  lede="Arrays guaranteed to have at least one element. Eliminates empty array edge cases and provides safe head/tail operations."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/array' },
    { label: 'Type', value: 'NonEmptyArray[A]' }
  ]}
/>

---

<Section num="1" title="Core API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Of` | `func Of[A any](head A, tail ...A) NonEmptyArray[A]` | Create from elements |
| `From` | `func From[A any]([]A) Option[NonEmptyArray[A]]` | Safe creation from slice |
| `Head` | `func Head[A any](NonEmptyArray[A]) A` | Get first element |
| `Tail` | `func Tail[A any](NonEmptyArray[A]) []A` | Get remaining elements |
| `Last` | `func Last[A any](NonEmptyArray[A]) A` | Get last element |
| `Init` | `func Init[A any](NonEmptyArray[A]) []A` | All but last |
| `Map` | `func Map[A, B any](f func(A) B) func(NonEmptyArray[A]) NonEmptyArray[B]` | Transform elements |
| `FlatMap` | `func FlatMap[A, B any](f func(A) NonEmptyArray[B]) func(NonEmptyArray[A]) NonEmptyArray[B]` | Map and flatten |
| `Reduce` | `func Reduce[A, B any](f func(B, A) B, initial B) func(NonEmptyArray[A]) B` | Fold left |
| `Concat` | `func Concat[A any](second NonEmptyArray[A]) func(NonEmptyArray[A]) NonEmptyArray[A]` | Combine arrays |
</ApiTable>

</Section>

---

<Section num="2" title="Usage Examples">

### Creating NonEmpty Arrays

<CodeCard file="create.go">
{`import (
    A "github.com/IBM/fp-go/v2/array"
    O "github.com/IBM/fp-go/v2/option"
)

// Direct creation - always succeeds
arr1 := A.Of(1, 2, 3, 4, 5)
// NonEmptyArray[int] with 5 elements

// Single element
arr2 := A.Of("hello")
// NonEmptyArray[string] with 1 element

// From slice - returns Option
slice := []int{1, 2, 3}
result := A.From(slice)
// Some(NonEmptyArray[int]{1, 2, 3})

emptySlice := []int{}
result2 := A.From(emptySlice)
// None - empty slice cannot be NonEmpty
`}
</CodeCard>

### Safe Head/Tail Operations

<CodeCard file="head_tail.go">
{`arr := A.Of(1, 2, 3, 4, 5)

// Head - always safe, no Option needed
first := A.Head(arr)  // 1

// Tail - remaining elements
rest := A.Tail(arr)   // []int{2, 3, 4, 5}

// Last - always safe
last := A.Last(arr)   // 5

// Init - all but last
initial := A.Init(arr) // []int{1, 2, 3, 4}

// Single element array
single := A.Of(42)
A.Head(single)  // 42
A.Tail(single)  // []int{} - empty slice
A.Last(single)  // 42
A.Init(single)  // []int{} - empty slice
`}
</CodeCard>

### Transformation

<CodeCard file="transform.go">
{`import F "github.com/IBM/fp-go/v2/function"

numbers := A.Of(1, 2, 3, 4, 5)

// Map - result is also NonEmpty
doubled := F.Pipe2(
    numbers,
    A.Map(func(n int) int { return n * 2 }),
)
// NonEmptyArray[int]{2, 4, 6, 8, 10}

// FlatMap
result := F.Pipe2(
    A.Of(1, 2, 3),
    A.FlatMap(func(n int) A.NonEmptyArray[int] {
        return A.Of(n, n*10)
    }),
)
// NonEmptyArray[int]{1, 10, 2, 20, 3, 30}
`}
</CodeCard>

### Reduction

<CodeCard file="reduce.go">
{`numbers := A.Of(1, 2, 3, 4, 5)

// Sum
sum := F.Pipe2(
    numbers,
    A.Reduce(
        func(acc, n int) int { return acc + n },
        0,
    ),
)
// 15

// Product
product := F.Pipe2(
    numbers,
    A.Reduce(
        func(acc, n int) int { return acc * n },
        1,
    ),
)
// 120

// Build string
words := A.Of("Hello", "functional", "world")
sentence := F.Pipe2(
    words,
    A.Reduce(
        func(acc, word string) string {
            if acc == "" {
                return word
            }
            return acc + " " + word
        },
        "",
    ),
)
// "Hello functional world"
`}
</CodeCard>

### Concatenation

<CodeCard file="concat.go">
{`arr1 := A.Of(1, 2, 3)
arr2 := A.Of(4, 5, 6)

// Combine - result is NonEmpty
combined := F.Pipe2(
    arr1,
    A.Concat(arr2),
)
// NonEmptyArray[int]{1, 2, 3, 4, 5, 6}

// Chain multiple
result := F.Pipe3(
    A.Of(1),
    A.Concat(A.Of(2, 3)),
    A.Concat(A.Of(4, 5)),
)
// NonEmptyArray[int]{1, 2, 3, 4, 5}
`}
</CodeCard>

### Safe Minimum/Maximum

<CodeCard file="min_max.go">
{`import "github.com/IBM/fp-go/v2/ord"

numbers := A.Of(5, 2, 8, 1, 9, 3)

// Find minimum - always succeeds
min := F.Pipe2(
    numbers,
    A.Reduce(
        func(min, n int) int {
            if n < min {
                return n
            }
            return min
        },
        A.Head(numbers),
    ),
)
// 1

// Find maximum
max := F.Pipe2(
    numbers,
    A.Reduce(
        func(max, n int) int {
            if n > max {
                return n
            }
            return max
        },
        A.Head(numbers),
    ),
)
// 9
`}
</CodeCard>

### Configuration Lists

<CodeCard file="config.go">
{`type Config struct {
    Servers NonEmptyArray[string]
    Ports   NonEmptyArray[int]
}

// Guarantee at least one server
config := Config{
    Servers: A.Of(
        "server1.example.com",
        "server2.example.com",
    ),
    Ports: A.Of(8080, 8081, 8082),
}

// Primary server is always available
primary := A.Head(config.Servers)
// "server1.example.com"

// Fallback servers
fallbacks := A.Tail(config.Servers)
// []string{"server2.example.com"}
`}
</CodeCard>

### Validation Results

<CodeCard file="validation.go">
{`type ValidationError struct {
    Field   string
    Message string
}

type ValidationResult struct {
    Errors NonEmptyArray[ValidationError]
}

// At least one error
result := ValidationResult{
    Errors: A.Of(
        ValidationError{
            Field:   "email",
            Message: "Invalid email format",
        },
        ValidationError{
            Field:   "password",
            Message: "Password too short",
        },
    ),
}

// First error
firstError := A.Head(result.Errors)

// All errors
allErrors := F.Pipe2(
    result.Errors,
    A.Map(func(e ValidationError) string {
        return e.Field + ": " + e.Message
    }),
)
`}
</CodeCard>

</Section>

---

<Section num="3" title="Common Patterns">

### Safe List Processing

<CodeCard file="safe_list.go">
{`// Process list that must have items
func ProcessOrders(orders NonEmptyArray[Order]) Result {
    // No need to check for empty
    firstOrder := A.Head(orders)
    
    // Process all orders
    return F.Pipe2(
        orders,
        A.Map(processOrder),
        combineResults,
    )
}
`}
</CodeCard>

### Builder Pattern

<CodeCard file="builder.go">
{`type QueryBuilder struct {
    conditions NonEmptyArray[string]
}

func NewQuery(first string, rest ...string) QueryBuilder {
    return QueryBuilder{
        conditions: A.Of(first, rest...),
    }
}

func (q QueryBuilder) And(condition string) QueryBuilder {
    return QueryBuilder{
        conditions: F.Pipe2(
            q.conditions,
            A.Concat(A.Of(condition)),
        ),
    }
}

func (q QueryBuilder) Build() string {
    return F.Pipe2(
        q.conditions,
        A.Reduce(
            func(acc, cond string) string {
                return acc + " AND " + cond
            },
            A.Head(q.conditions),
        ),
    )
}

// Usage
query := NewQuery("age > 18").
    And("status = 'active'").
    And("verified = true").
    Build()
// "age > 18 AND status = 'active' AND verified = true"
`}
</CodeCard>

### Converting to/from Regular Arrays

<CodeCard file="conversion.go">
{`// From regular array
regularArray := []int{1, 2, 3}

nonEmpty := F.Pipe2(
    regularArray,
    A.From,
    O.Match(
        func() NonEmptyArray[int] {
            // Provide default
            return A.Of(0)
        },
        F.Identity[NonEmptyArray[int]],
    ),
)

// To regular array
toRegular := func(nea NonEmptyArray[int]) []int {
    return append([]int{A.Head(nea)}, A.Tail(nea)...)
}
`}
</CodeCard>

</Section>

---

<Callout type="info">

**Type Safety**: NonEmptyArray eliminates the need for runtime empty checks. Operations like `Head` and `Last` are always safe and return values directly, not Options.

</Callout>

<Callout type="info">

**When to Use**: Use NonEmptyArray when your domain logic requires at least one element - configuration lists, validation errors, search results that must have matches, etc.

</Callout>


---

<Pager
  prev={{ to: '/docs/v2/utilities/pipe-flow', title: 'Pipe & Flow' }}
  next={{ to: '/docs/v2/collections/sequence-traverse', title: 'Sequence & Traverse' }}
/>

---
