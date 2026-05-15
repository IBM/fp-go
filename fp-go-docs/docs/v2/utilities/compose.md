---
title: Compose
hide_title: true
description: Right-to-left function composition for mathematical-style function combination in fp-go
sidebar_position: 22
---

import { PageHeader, Section, CodeCard, ApiTable, Compare, Callout, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="v2 · Utilities"
  title="Compose"
  titleAccent="Right-to-Left"
  lede="Compose creates new functions by combining existing ones, executing from right to left in mathematical composition style."
  meta={[
    { label: 'Package', value: 'function' },
    { label: 'Since', value: 'v2.0.0' },
    { label: 'Direction', value: 'Right-to-Left' }
  ]}
/>

---

<Section num="1" title="Overview">

**Compose** combines functions in right-to-left order, following mathematical function composition notation: `(f ∘ g)(x) = f(g(x))`.

**Key Characteristics:**
- **Right-to-left execution**: Inner function executes first
- **Mathematical notation**: Matches `f ∘ g` composition
- **Type-safe**: Ensures output of one function matches input of next

<Callout type="info">
For left-to-right composition (more readable in code), use [Flow](./pipe-flow.md) instead.
</Callout>

</Section>

---

<Section num="2" title="Basic Composition">

<CodeCard file="compose_basic.go" tag="example">
{`import F "github.com/IBM/fp-go/function"

// Define simple functions
double := func(n int) int {
    return n * 2
}

addTen := func(n int) int {
    return n + 10
}

// Compose: addTen(double(x))
// Executes right-to-left: double first, then addTen
transform := F.Compose2(addTen, double)

result := transform(5)
// 5 -> double -> 10 -> addTen -> 20

// Equivalent to:
result := addTen(double(5))  // 20`}
</CodeCard>

</Section>

---

<Section num="3" title="Compose vs Flow">

Understanding the difference between Compose and Flow:

<Compare>
<div slot="left">

**Compose (Right-to-Left)**
```go
// Mathematical: f ∘ g
c := F.Compose2(f, g)
// Executes: f(g(x))

double := func(n int) int { return n * 2 }
addTen := func(n int) int { return n + 10 }

// Read right-to-left
transform := F.Compose2(addTen, double)
transform(5)  // 20
// 5 -> double(5)=10 -> addTen(10)=20
```

</div>
<div slot="right">

**Flow (Left-to-Right)**
```go
// Pipeline style
f := F.Flow2(g, f)
// Executes: f(g(x))

double := func(n int) int { return n * 2 }
addTen := func(n int) int { return n + 10 }

// Read left-to-right
transform := F.Flow2(double, addTen)
transform(5)  // 20
// 5 -> double(5)=10 -> addTen(10)=20
```

</div>
</Compare>

<Callout type="warn">
Both produce the same result, but the order of arguments is reversed. Choose based on readability preference.
</Callout>

</Section>

---

<Section num="4" title="Multiple Functions">

Compose multiple functions together:

<CodeCard file="compose_multiple.go" tag="example">
{`import F "github.com/IBM/fp-go/function"

// Define transformation functions
double := func(n int) int { return n * 2 }
addTen := func(n int) int { return n + 10 }
square := func(n int) int { return n * n }

// Compose three functions
// Executes: square(addTen(double(x)))
transform := F.Compose3(square, addTen, double)

result := transform(5)
// 5 -> double -> 10 -> addTen -> 20 -> square -> 400

// Step by step:
// 1. double(5) = 10
// 2. addTen(10) = 20
// 3. square(20) = 400`}
</CodeCard>

</Section>

---

<Section num="5" title="String Transformations">

Compose string operations:

<CodeCard file="compose_strings.go" tag="example">
{`import (
    F "github.com/IBM/fp-go/function"
    "strings"
)

// Define string transformations
trim := func(s string) string {
    return strings.TrimSpace(s)
}

upper := func(s string) string {
    return strings.ToUpper(s)
}

addPrefix := func(s string) string {
    return ">>> " + s
}

// Compose: addPrefix(upper(trim(x)))
normalize := F.Compose3(addPrefix, upper, trim)

result := normalize("  hello world  ")
// "  hello world  " -> trim -> "hello world" 
//                   -> upper -> "HELLO WORLD"
//                   -> addPrefix -> ">>> HELLO WORLD"`}
</CodeCard>

</Section>

---

<Section num="6" title="Type Transformations">

Compose functions that change types:

<CodeCard file="compose_types.go" tag="example">
{`import F "github.com/IBM/fp-go/function"

// int -> string
toString := func(n int) string {
    return fmt.Sprintf("%d", n)
}

// string -> int (length)
length := func(s string) int {
    return len(s)
}

// int -> bool
isEven := func(n int) bool {
    return n%2 == 0
}

// Compose: isEven(length(toString(x)))
// int -> string -> int -> bool
check := F.Compose3(isEven, length, toString)

check(42)    // true  (toString="42", len=2, isEven=true)
check(123)   // false (toString="123", len=3, isEven=false)
check(1000)  // true  (toString="1000", len=4, isEven=true)`}
</CodeCard>

</Section>

---

<Section num="7" title="Validation Pipeline">

Build validation chains:

<CodeCard file="compose_validation.go" tag="example">
{`type User struct {
    Name  string
    Email string
    Age   int
}

// Validation functions
validateName := func(u User) User {
    if u.Name == "" {
        panic("Name required")
    }
    return u
}

validateEmail := func(u User) User {
    if !strings.Contains(u.Email, "@") {
        panic("Invalid email")
    }
    return u
}

validateAge := func(u User) User {
    if u.Age < 18 {
        panic("Must be 18+")
    }
    return u
}

// Compose validators (right-to-left)
// Executes: validateAge -> validateEmail -> validateName
validate := F.Compose3(validateName, validateEmail, validateAge)

user := User{Name: "Alice", Email: "alice@example.com", Age: 25}
validated := validate(user)  // All validations pass`}
</CodeCard>

</Section>

---

<Section num="8" title="When to Use Compose">

<Callout type="info">
**Use Compose when:**
- You prefer mathematical notation
- Working with mathematical concepts
- Porting code from languages with `∘` operator
- Building abstract function combinators

**Use Flow when:**
- You prefer readable left-to-right pipelines
- Working with data transformations
- Building business logic
- Most practical applications
</Callout>

</Section>

---

<Section num="9" title="API Reference">

<ApiTable>
| Function | Type | Description |
|----------|------|-------------|
| `Compose2[A, B, C]` | `(B -> C, A -> B) -> (A -> C)` | Compose two functions |
| `Compose3[A, B, C, D]` | `(C -> D, B -> C, A -> B) -> (A -> D)` | Compose three functions |
| `Compose4[A, B, C, D, E]` | `(D -> E, C -> D, B -> C, A -> B) -> (A -> E)` | Compose four functions |
</ApiTable>

**Execution Order:**
```
Compose2(f, g)(x) = f(g(x))
Compose3(f, g, h)(x) = f(g(h(x)))
```

</Section>

---

<Section num="10" title="Related Concepts">

**Function Composition Styles:**
- **Compose**: Right-to-left (mathematical)
- **Flow**: Left-to-right (pipeline)
- **Pipe**: Data-first left-to-right

**See Also:**
- [Pipe & Flow](./pipe-flow.md) - Left-to-right composition
- [Function](./function.md) - Core function utilities

</Section>

---

<Pager
  prev={{ to: '/docs/v2/utilities/tuple', title: 'Tuple' }}
  next={{ to: '/docs/v2/utilities/bind-curry', title: 'Bind & Curry' }}
/>