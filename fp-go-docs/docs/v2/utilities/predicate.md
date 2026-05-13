---
title: Predicate
hide_title: true
description: Boolean-valued functions and combinators for filtering and validation in fp-go
sidebar_position: 23
---

import { PageHeader, Section, CodeCard, ApiTable, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="v2 · Utilities"
  title="Predicate"
  titleAccent="Boolean Logic"
  lede="Predicates are functions that return boolean values. The predicate package provides utilities for combining and manipulating them."
  meta={[
    { label: 'Package', value: 'predicate' },
    { label: 'Since', value: 'v2.0.0' },
    { label: 'Type', value: 'func(A) bool' }
  ]}
/>

---

<Section num="1" title="Overview">

A **Predicate** is a function that takes a value and returns a boolean. Predicates are fundamental for:
- Filtering collections
- Validation logic
- Conditional operations
- Boolean algebra

The predicate package provides combinators to compose complex predicates from simple ones.

</Section>

---

<Section num="2" title="Basic Predicates">

<CodeCard file="predicate_basic.go" tag="example">
{`import P "github.com/IBM/fp-go/predicate"

// Define simple predicates
isEven := func(n int) bool {
    return n%2 == 0
}

isPositive := func(n int) bool {
    return n > 0
}

isGreaterThan10 := func(n int) bool {
    return n > 10
}

// Use predicates
isEven(4)        // true
isPositive(-5)   // false
isGreaterThan10(15)  // true`}
</CodeCard>

</Section>

---

<Section num="3" title="Not Combinator">

Negate a predicate:

<CodeCard file="predicate_not.go" tag="example">
{`import P "github.com/IBM/fp-go/predicate"

isEven := func(n int) bool { return n%2 == 0 }

// Negate the predicate
isOdd := P.Not(isEven)

isOdd(3)  // true
isOdd(4)  // false

// Double negation
isEvenAgain := P.Not(isOdd)
isEvenAgain(4)  // true`}
</CodeCard>

</Section>

---

<Section num="4" title="And Combinator">

Combine predicates with logical AND:

<CodeCard file="predicate_and.go" tag="example">
{`import P "github.com/IBM/fp-go/predicate"

isEven := func(n int) bool { return n%2 == 0 }
isPositive := func(n int) bool { return n > 0 }

// Both conditions must be true
isEvenAndPositive := P.And(isEven, isPositive)

isEvenAndPositive(4)   // true (even AND positive)
isEvenAndPositive(-2)  // false (even but NOT positive)
isEvenAndPositive(3)   // false (positive but NOT even)
isEvenAndPositive(-3)  // false (neither)

// Chain multiple conditions
isGreaterThan10 := func(n int) bool { return n > 10 }
complexPredicate := P.And(
    isEven,
    P.And(isPositive, isGreaterThan10),
)
complexPredicate(12)  // true
complexPredicate(8)   // false (not > 10)`}
</CodeCard>

</Section>

---

<Section num="5" title="Or Combinator">

Combine predicates with logical OR:

<CodeCard file="predicate_or.go" tag="example">
{`import P "github.com/IBM/fp-go/predicate"

isEven := func(n int) bool { return n%2 == 0 }
isPositive := func(n int) bool { return n > 0 }

// At least one condition must be true
isEvenOrPositive := P.Or(isEven, isPositive)

isEvenOrPositive(4)   // true (even OR positive - both)
isEvenOrPositive(-2)  // true (even OR positive - even only)
isEvenOrPositive(3)   // true (even OR positive - positive only)
isEvenOrPositive(-3)  // false (neither)

// Complex combinations
isSpecial := P.Or(
    P.And(isEven, isPositive),
    func(n int) bool { return n == 0 },
)
isSpecial(4)   // true (even and positive)
isSpecial(0)   // true (special case)
isSpecial(-2)  // false`}
</CodeCard>

</Section>

---

<Section num="6" title="Validation Example">

Use predicates for complex validation:

<CodeCard file="predicate_validation.go" tag="example">
{`type User struct {
    Name  string
    Age   int
    Email string
}

// Define validation predicates
hasName := func(u User) bool {
    return u.Name != ""
}

isAdult := func(u User) bool {
    return u.Age >= 18
}

hasEmail := func(u User) bool {
    return u.Email != "" && strings.Contains(u.Email, "@")
}

// Combine validators
isValidUser := P.And(
    hasName,
    P.And(isAdult, hasEmail),
)

user1 := User{Name: "Alice", Age: 25, Email: "alice@example.com"}
user2 := User{Name: "Bob", Age: 16, Email: "bob@example.com"}
user3 := User{Name: "", Age: 30, Email: "test@example.com"}

isValidUser(user1)  // true
isValidUser(user2)  // false (not adult)
isValidUser(user3)  // false (no name)`}
</CodeCard>

</Section>

---

<Section num="7" title="Filtering with Predicates">

Use predicates to filter collections:

<CodeCard file="predicate_filter.go" tag="example">
{`import (
    A "github.com/IBM/fp-go/array"
    F "github.com/IBM/fp-go/function"
    P "github.com/IBM/fp-go/predicate"
)

numbers := []int{-2, -1, 0, 1, 2, 3, 4, 5, 6}

isEven := func(n int) bool { return n%2 == 0 }
isPositive := func(n int) bool { return n > 0 }

// Filter with combined predicates
evenAndPositive := F.Pipe2(
    numbers,
    A.Filter(P.And(isEven, isPositive)),
)
// []int{2, 4, 6}

// Filter with OR
evenOrPositive := F.Pipe2(
    numbers,
    A.Filter(P.Or(isEven, isPositive)),
)
// []int{-2, 0, 1, 2, 3, 4, 5, 6}

// Filter with NOT
notEven := F.Pipe2(
    numbers,
    A.Filter(P.Not(isEven)),
)
// []int{-1, 1, 3, 5}`}
</CodeCard>

</Section>

---

<Section num="8" title="API Reference">

<ApiTable>
| Function | Type | Description |
|----------|------|-------------|
| `Not[A]` | `Predicate[A] -> Predicate[A]` | Negates a predicate |
| `And[A]` | `(Predicate[A], Predicate[A]) -> Predicate[A]` | Logical AND of two predicates |
| `Or[A]` | `(Predicate[A], Predicate[A]) -> Predicate[A]` | Logical OR of two predicates |
</ApiTable>

**Type Definition:**
```go
type Predicate[A any] func(A) bool
```

</Section>

---

<Section num="9" title="Related Concepts">

**Common Use Cases:**
- Filtering arrays and collections
- Validation logic
- Conditional branching
- Boolean algebra operations

**See Also:**
- [Function](./function.md) - Core function utilities
- [Boolean](./boolean.md) - Boolean type class instances
- [Array Filter](../collections/array.md) - Using predicates with arrays

</Section>

---

<Pager
  prev={{ to: '/docs/v2/utilities/monoid', title: 'Monoid' }}
  next={{ to: '/docs/v2/utilities/boolean', title: 'Boolean' }}
/>