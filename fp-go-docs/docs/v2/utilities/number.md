---
title: Number
hide_title: true
description: Numeric utilities and type class instances for mathematical operations in fp-go
sidebar_position: 30
---

import { PageHeader, Section, CodeCard, ApiTable, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="v2 · Utilities"
  title="Number"
  titleAccent="Type Classes"
  lede="The number package provides type class instances and utilities for numeric types, enabling functional operations on numbers."
  meta={[
    { label: 'Package', value: 'number' },
    { label: 'Since', value: 'v2.0.0' },
    { label: 'Types', value: 'int, float64, etc.' }
  ]}
/>

---

<Section num="1" title="Overview">

The **Number** package provides type class instances for numeric types, enabling:
- **Eq**: Equality comparison
- **Ord**: Ordering and comparison
- **Semigroup**: Addition and multiplication
- **Monoid**: Addition and multiplication with identity elements

These instances work with Go's built-in numeric types: `int`, `int8`, `int16`, `int32`, `int64`, `float32`, `float64`, etc.

</Section>

---

<Section num="2" title="Equality">

<CodeCard file="number_eq.go" tag="example">
{`import N "github.com/IBM/fp-go/number"

// Compare numbers for equality
N.Eq.Equals(1, 1)    // true
N.Eq.Equals(1, 2)    // false
N.Eq.Equals(0, 0)    // true

// Works with different numeric types
N.Eq.Equals(3.14, 3.14)  // true (float64)
N.Eq.Equals(42, 42)      // true (int)`}
</CodeCard>

</Section>

---

<Section num="3" title="Ordering">

Compare and sort numbers:

<CodeCard file="number_ord.go" tag="example">
{`import N "github.com/IBM/fp-go/number"

// Compare numbers
N.Ord.Compare(1, 2)   // -1 (less than)
N.Ord.Compare(2, 1)   // 1  (greater than)
N.Ord.Compare(1, 1)   // 0  (equal)

// Derived operations
N.Ord.LessThan(1, 2)           // true
N.Ord.GreaterThan(2, 1)        // true
N.Ord.LessThanOrEqual(1, 1)    // true
N.Ord.GreaterThanOrEqual(2, 1) // true

// Min and Max
N.Ord.Min(5, 3)  // 3
N.Ord.Max(5, 3)  // 5`}
</CodeCard>

</Section>

---

<Section num="4" title="Semigroup Operations">

Combine numbers with addition or multiplication:

<CodeCard file="number_semigroup.go" tag="example">
{`import N "github.com/IBM/fp-go/number"

// SemigroupSum: Addition
N.SemigroupSum.Concat(1, 2)   // 3
N.SemigroupSum.Concat(10, 5)  // 15

// SemigroupProduct: Multiplication
N.SemigroupProduct.Concat(3, 4)   // 12
N.SemigroupProduct.Concat(2, 10)  // 20

// SemigroupMin: Minimum
N.SemigroupMin.Concat(5, 3)  // 3
N.SemigroupMin.Concat(1, 9)  // 1

// SemigroupMax: Maximum
N.SemigroupMax.Concat(5, 3)  // 5
N.SemigroupMax.Concat(1, 9)  // 9`}
</CodeCard>

</Section>

---

<Section num="5" title="Monoid Operations">

Semigroups with identity elements:

<CodeCard file="number_monoid.go" tag="example">
{`import N "github.com/IBM/fp-go/number"

// MonoidSum: Addition with identity 0
N.MonoidSum.Concat(1, 2)  // 3
N.MonoidSum.Empty()       // 0

// MonoidProduct: Multiplication with identity 1
N.MonoidProduct.Concat(3, 4)  // 12
N.MonoidProduct.Empty()       // 1

// Identity laws
val := 5
N.MonoidSum.Concat(N.MonoidSum.Empty(), val)  // 5
N.MonoidSum.Concat(val, N.MonoidSum.Empty())  // 5

N.MonoidProduct.Concat(N.MonoidProduct.Empty(), val)  // 5
N.MonoidProduct.Concat(val, N.MonoidProduct.Empty())  // 5`}
</CodeCard>

</Section>

---

<Section num="6" title="Summing Arrays">

Calculate sum and product of arrays:

<CodeCard file="number_sum.go" tag="example">
{`import (
    A "github.com/IBM/fp-go/array"
    F "github.com/IBM/fp-go/function"
    N "github.com/IBM/fp-go/number"
)

numbers := []int{1, 2, 3, 4, 5}

// Sum all numbers
sum := F.Pipe2(
    numbers,
    A.Fold(N.MonoidSum),
)
// 15

// Product of all numbers
product := F.Pipe2(
    numbers,
    A.Fold(N.MonoidProduct),
)
// 120

// Empty array uses identity
emptySum := F.Pipe2([]int{}, A.Fold(N.MonoidSum))
// 0 (identity for addition)

emptyProduct := F.Pipe2([]int{}, A.Fold(N.MonoidProduct))
// 1 (identity for multiplication)`}
</CodeCard>

</Section>

---

<Section num="7" title="Sorting Arrays">

Sort numbers using Ord instance:

<CodeCard file="number_sort.go" tag="example">
{`import (
    A "github.com/IBM/fp-go/array"
    F "github.com/IBM/fp-go/function"
    N "github.com/IBM/fp-go/number"
)

numbers := []int{5, 2, 8, 1, 9, 3}

// Sort ascending
sorted := F.Pipe2(
    numbers,
    A.Sort(N.Ord),
)
// []int{1, 2, 3, 5, 8, 9}

// Sort descending (reverse the Ord)
import O "github.com/IBM/fp-go/ord"

sortedDesc := F.Pipe2(
    numbers,
    A.Sort(O.Reverse(N.Ord)),
)
// []int{9, 8, 5, 3, 2, 1}`}
</CodeCard>

</Section>

---

<Section num="8" title="Finding Min/Max">

Find minimum and maximum values:

<CodeCard file="number_minmax.go" tag="example">
{`import (
    A "github.com/IBM/fp-go/array"
    F "github.com/IBM/fp-go/function"
    N "github.com/IBM/fp-go/number"
    O "github.com/IBM/fp-go/option"
)

numbers := []int{5, 2, 8, 1, 9, 3}

// Find minimum
minValue := F.Pipe2(
    numbers,
    A.Head[int],
)
// Some(5) - first element

// Better: use Fold with SemigroupMin
minValue := F.Pipe3(
    numbers,
    A.Reduce(func(acc, n int) int {
        return N.SemigroupMin.Concat(acc, n)
    }, numbers[0]),
)
// 1

// Find maximum
maxValue := F.Pipe3(
    numbers,
    A.Reduce(func(acc, n int) int {
        return N.SemigroupMax.Concat(acc, n)
    }, numbers[0]),
)
// 9`}
</CodeCard>

</Section>

---

<Section num="9" title="Statistical Operations">

Calculate statistics using monoids:

<CodeCard file="number_stats.go" tag="example">
{`type Stats struct {
    Count int
    Sum   int
    Min   int
    Max   int
}

numbers := []int{5, 2, 8, 1, 9, 3}

// Calculate all stats in one pass
stats := F.Pipe3(
    numbers,
    A.Reduce(func(acc Stats, n int) Stats {
        return Stats{
            Count: acc.Count + 1,
            Sum:   N.MonoidSum.Concat(acc.Sum, n),
            Min:   N.SemigroupMin.Concat(acc.Min, n),
            Max:   N.SemigroupMax.Concat(acc.Max, n),
        }
    }, Stats{
        Count: 0,
        Sum:   0,
        Min:   numbers[0],
        Max:   numbers[0],
    }),
)
// Stats{Count: 6, Sum: 28, Min: 1, Max: 9}

average := float64(stats.Sum) / float64(stats.Count)
// 4.666...`}
</CodeCard>

</Section>

---

<Section num="10" title="API Reference">

<ApiTable>
| Instance | Type | Description |
|----------|------|-------------|
| `Eq` | `Eq[T]` | Equality comparison |
| `Ord` | `Ord[T]` | Ordering and comparison |
| `SemigroupSum` | `Semigroup[T]` | Addition |
| `SemigroupProduct` | `Semigroup[T]` | Multiplication |
| `SemigroupMin` | `Semigroup[T]` | Minimum |
| `SemigroupMax` | `Semigroup[T]` | Maximum |
| `MonoidSum` | `Monoid[T]` | Addition with identity 0 |
| `MonoidProduct` | `Monoid[T]` | Multiplication with identity 1 |
</ApiTable>

**Supported Types:**
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`

</Section>

---

<Section num="11" title="Related Concepts">

**Common Use Cases:**
- Summing and aggregating numbers
- Finding min/max values
- Sorting numeric arrays
- Statistical calculations

**See Also:**
- [String](./string.md) - String type class instances
- [Eq](./eq.md) - Equality type class
- [Ord](./ord.md) - Ordering type class
- [Monoid](./monoid.md) - Understanding monoid operations

</Section>

---

<Pager
  prev={{ to: '/docs/v2/utilities/boolean', title: 'Boolean' }}
  next={{ to: '/docs/v2/utilities/string', title: 'String' }}
/>