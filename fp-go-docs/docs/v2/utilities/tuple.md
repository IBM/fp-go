---
title: Tuple
hide_title: true
description: Working with pairs and tuples for grouping multiple values in fp-go
sidebar_position: 24
---

import { PageHeader, Section, CodeCard, ApiTable, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="v2 · Utilities"
  title="Tuple"
  titleAccent="Pairs & Groups"
  lede="Tuples provide a way to group multiple values together. The tuple package provides utilities for creating and manipulating tuples."
  meta={[
    { label: 'Package', value: 'tuple' },
    { label: 'Since', value: 'v2.0.0' },
    { label: 'Type', value: 'Tuple2[A, B]' }
  ]}
/>

---

<Section num="1" title="Overview">

A **Tuple** is a fixed-size collection of values that can have different types. The most common is **Tuple2** (also called a pair), which holds two values.

**Key Features:**
- Type-safe grouping of values
- Immutable by design
- Useful for function returns and data pairing
- Supports mapping and transformation operations

</Section>

---

<Section num="2" title="Creating Tuples">

<CodeCard file="tuple_create.go" tag="example">
{`import T "github.com/IBM/fp-go/tuple"

// Create a tuple (pair)
pair := T.MakeTuple2("Alice", 30)
// Tuple2[string, int]{Head: "Alice", Tail: 30}

// Access elements
name := pair.Head   // "Alice" (first element)
age := pair.Tail    // 30 (second element)

// Different types
mixed := T.MakeTuple2(42, "answer")
// Tuple2[int, string]{Head: 42, Tail: "answer"}

// Nested tuples
nested := T.MakeTuple2(
    T.MakeTuple2(1, 2),
    T.MakeTuple2(3, 4),
)
// Tuple2[Tuple2[int, int], Tuple2[int, int]]`}
</CodeCard>

</Section>

---

<Section num="3" title="Mapping First Element">

Transform the first element of a tuple:

<CodeCard file="tuple_mapfst.go" tag="example">
{`import T "github.com/IBM/fp-go/tuple"

pair := T.MakeTuple2(5, "hello")

// Map first element (Head)
result := T.MapFst(func(n int) int {
    return n * 2
})(pair)
// Tuple2{Head: 10, Tail: "hello"}

// Chain transformations
result := T.MapFst(func(n int) string {
    return fmt.Sprintf("Number: %d", n)
})(pair)
// Tuple2{Head: "Number: 5", Tail: "hello"}`}
</CodeCard>

</Section>

---

<Section num="4" title="Mapping Second Element">

Transform the second element of a tuple:

<CodeCard file="tuple_mapsnd.go" tag="example">
{`import T "github.com/IBM/fp-go/tuple"

pair := T.MakeTuple2("hello", 10)

// Map second element (Tail)
result := T.MapSnd(func(n int) int {
    return n * 2
})(pair)
// Tuple2{Head: "hello", Tail: 20}

// Change type
result := T.MapSnd(func(n int) string {
    return fmt.Sprintf("%d items", n)
})(pair)
// Tuple2{Head: "hello", Tail: "10 items"}`}
</CodeCard>

</Section>

---

<Section num="5" title="Mapping Both Elements">

Transform both elements simultaneously:

<CodeCard file="tuple_bimap.go" tag="example">
{`import T "github.com/IBM/fp-go/tuple"

pair := T.MakeTuple2(5, 10)

// Map both elements
result := T.Bimap(
    func(a int) int { return a * 2 },
    func(b int) int { return b + 1 },
)(pair)
// Tuple2{Head: 10, Tail: 11}

// Change both types
result := T.Bimap(
    func(a int) string { return fmt.Sprintf("A=%d", a) },
    func(b int) bool { return b > 5 },
)(pair)
// Tuple2{Head: "A=5", Tail: true}`}
</CodeCard>

</Section>

---

<Section num="6" title="Swapping Elements">

Reverse the order of tuple elements:

<CodeCard file="tuple_swap.go" tag="example">
{`import T "github.com/IBM/fp-go/tuple"

pair := T.MakeTuple2("hello", 42)
// Tuple2[string, int]{Head: "hello", Tail: 42}

swapped := T.Swap(pair)
// Tuple2[int, string]{Head: 42, Tail: "hello"}

// Swap is its own inverse
original := T.Swap(swapped)
// Back to Tuple2[string, int]{Head: "hello", Tail: 42}`}
</CodeCard>

</Section>

---

<Section num="7" title="Function Results">

Return multiple values as tuples:

<CodeCard file="tuple_results.go" tag="example">
{`import T "github.com/IBM/fp-go/tuple"

// Return quotient and remainder
func divMod(a, b int) T.Tuple2[int, int] {
    return T.MakeTuple2(a/b, a%b)
}

result := divMod(17, 5)
quotient := result.Head   // 3
remainder := result.Tail  // 2

// Parse and validate
func parseAndValidate(s string) T.Tuple2[int, error] {
    val, err := strconv.Atoi(s)
    return T.MakeTuple2(val, err)
}

result := parseAndValidate("42")
value := result.Head  // 42
err := result.Tail    // nil`}
</CodeCard>

</Section>

---

<Section num="8" title="Key-Value Pairs">

Use tuples for key-value associations:

<CodeCard file="tuple_keyvalue.go" tag="example">
{`import (
    R "github.com/IBM/fp-go/record"
    T "github.com/IBM/fp-go/tuple"
)

// Convert map to array of tuples
myMap := map[string]int{
    "apples":  5,
    "bananas": 3,
    "oranges": 7,
}

entries := R.ToEntries(myMap)
// []Tuple2[string, int]{
//   {Head: "apples", Tail: 5},
//   {Head: "bananas", Tail: 3},
//   {Head: "oranges", Tail: 7},
// }

// Process pairs
formatted := F.Pipe2(
    entries,
    A.Map(func(t T.Tuple2[string, int]) string {
        return fmt.Sprintf("%s: %d", t.Head, t.Tail)
    }),
)
// []string{"apples: 5", "bananas: 3", "oranges: 7"}`}
</CodeCard>

</Section>

---

<Section num="9" title="Zipping Arrays">

Combine two arrays into tuples:

<CodeCard file="tuple_zip.go" tag="example">
{`import (
    A "github.com/IBM/fp-go/array"
    T "github.com/IBM/fp-go/tuple"
)

names := []string{"Alice", "Bob", "Charlie"}
ages := []int{30, 25, 35}

// Zip arrays into tuples
pairs := A.Zip(ages)(names)
// []Tuple2[string, int]{
//   {Head: "Alice", Tail: 30},
//   {Head: "Bob", Tail: 25},
//   {Head: "Charlie", Tail: 35},
// }

// Process zipped data
formatted := F.Pipe2(
    pairs,
    A.Map(func(t T.Tuple2[string, int]) string {
        return fmt.Sprintf("%s is %d years old", t.Head, t.Tail)
    }),
)
// []string{
//   "Alice is 30 years old",
//   "Bob is 25 years old",
//   "Charlie is 35 years old",
// }`}
</CodeCard>

</Section>

---

<Section num="10" title="Unzipping Tuples">

Split array of tuples back into separate arrays:

<CodeCard file="tuple_unzip.go" tag="example">
{`import (
    A "github.com/IBM/fp-go/array"
    T "github.com/IBM/fp-go/tuple"
)

pairs := []T.Tuple2[string, int]{
    T.MakeTuple2("Alice", 30),
    T.MakeTuple2("Bob", 25),
    T.MakeTuple2("Charlie", 35),
}

// Unzip into separate arrays
result := A.Unzip(pairs)
names := result.Head  // []string{"Alice", "Bob", "Charlie"}
ages := result.Tail   // []int{30, 25, 35}`}
</CodeCard>

</Section>

---

<Section num="11" title="API Reference">

<ApiTable>
| Function | Type | Description |
|----------|------|-------------|
| `MakeTuple2[A, B]` | `(A, B) -> Tuple2[A, B]` | Creates a tuple from two values |
| `MapFst[A, B, C]` | `(A -> C) -> Tuple2[A, B] -> Tuple2[C, B]` | Maps the first element |
| `MapSnd[A, B, C]` | `(B -> C) -> Tuple2[A, B] -> Tuple2[A, C]` | Maps the second element |
| `Bimap[A, B, C, D]` | `(A -> C, B -> D) -> Tuple2[A, B] -> Tuple2[C, D]` | Maps both elements |
| `Swap[A, B]` | `Tuple2[A, B] -> Tuple2[B, A]` | Swaps tuple elements |
</ApiTable>

**Type Definition:**
```go
type Tuple2[A, B any] struct {
    Head A  // First element
    Tail B  // Second element
}
```

</Section>

---

<Section num="12" title="Related Concepts">

**Common Use Cases:**
- Returning multiple values from functions
- Key-value pairs and associations
- Zipping and unzipping arrays
- Intermediate data structures in transformations

**See Also:**
- [Array Zip](../collections/array-zip.md) - Combining arrays into tuples
- [Record](../collections/record.md) - Key-value operations with maps

</Section>

---

<Pager
  prev={{ to: '/docs/v2/utilities/string', title: 'String' }}
  next={{ to: '/docs/v2/utilities/compose', title: 'Compose' }}
/>