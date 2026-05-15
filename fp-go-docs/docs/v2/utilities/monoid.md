---
title: Monoid
hide_title: true
description: Semigroup with identity element for combining values with a default in fp-go
sidebar_position: 28
---

import { PageHeader, Section, CodeCard, ApiTable, Callout, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="v2 · Utilities"
  title="Monoid"
  titleAccent="Identity & Combination"
  lede="A monoid is a semigroup with an identity element, providing both a way to combine values and a default/empty value."
  meta={[
    { label: 'Package', value: 'monoid' },
    { label: 'Since', value: 'v2.0.0' },
    { label: 'Extends', value: 'Semigroup' }
  ]}
/>

---

<Section num="1" title="Overview">

A **Monoid** extends Semigroup by adding an **identity element** (also called "empty" or "zero"). This identity element acts as a neutral value for the binary operation.

**Key Properties:**
- **Associative**: `(a • b) • c = a • (b • c)` (from Semigroup)
- **Left Identity**: `empty • a = a`
- **Right Identity**: `a • empty = a`

<Callout type="info">
The identity element makes monoids perfect for folding/reducing operations, as you always have a safe starting value.
</Callout>

</Section>

---

<Section num="2" title="Basic Usage">

<CodeCard file="monoid_basic.go" tag="example">
{`import M "github.com/IBM/fp-go/monoid"

// Create a string concatenation monoid
stringMonoid := M.MakeMonoid(
    func(a, b string) string { return a + b },  // Concat operation
    "",  // Identity element (empty string)
)

result := stringMonoid.Concat("Hello", " World")
// "Hello World"

empty := stringMonoid.Empty()
// ""

// Identity laws hold
s1 := stringMonoid.Concat(empty, "test")  // "test"
s2 := stringMonoid.Concat("test", empty)  // "test"
// s1 == s2 == "test"`}
</CodeCard>

</Section>

---

<Section num="3" title="Built-in Monoids">

fp-go provides several built-in monoids for common types:

<CodeCard file="monoid_builtin.go" tag="example">
{`import (
    N "github.com/IBM/fp-go/number"
    S "github.com/IBM/fp-go/string"
)

// Number addition monoid (identity: 0)
sum := N.MonoidSum.Concat(1, 2)  // 3
zero := N.MonoidSum.Empty()      // 0

// Number multiplication monoid (identity: 1)
product := N.MonoidProduct.Concat(3, 4)  // 12
one := N.MonoidProduct.Empty()           // 1

// String concatenation monoid (identity: "")
text := S.Monoid.Concat("Hello", " World")  // "Hello World"
emptyStr := S.Monoid.Empty()                // ""

// Min/Max monoids
minVal := N.MonoidMin.Concat(5, 3)  // 3
maxVal := N.MonoidMax.Concat(5, 3)  // 5`}
</CodeCard>

</Section>

---

<Section num="4" title="Folding with Monoids">

Monoids are perfect for folding/reducing collections:

<CodeCard file="monoid_fold.go" tag="example">
{`import (
    A "github.com/IBM/fp-go/array"
    F "github.com/IBM/fp-go/function"
    N "github.com/IBM/fp-go/number"
)

// Sum array of numbers
numbers := []int{1, 2, 3, 4, 5}
sum := F.Pipe2(
    numbers,
    A.Fold(N.MonoidSum),
)
// 15

// Product of numbers
product := F.Pipe2(
    numbers,
    A.Fold(N.MonoidProduct),
)
// 120

// Empty array uses identity
emptySum := F.Pipe2(
    []int{},
    A.Fold(N.MonoidSum),
)
// 0 (identity element)`}
</CodeCard>

</Section>

---

<Section num="5" title="Configuration Merging">

Use monoids to merge configurations:

<CodeCard file="monoid_config.go" tag="example">
{`type Config struct {
    Timeout int
    Retries int
    Debug   bool
}

// Max-merge strategy with defaults
configMonoid := M.MakeMonoid(
    func(a, b Config) Config {
        return Config{
            Timeout: max(a.Timeout, b.Timeout),
            Retries: max(a.Retries, b.Retries),
            Debug:   a.Debug || b.Debug,
        }
    },
    Config{Timeout: 0, Retries: 0, Debug: false},  // Identity
)

// Merge multiple configs
configs := []Config{
    {Timeout: 30, Retries: 3, Debug: false},  // defaults
    {Timeout: 0, Retries: 5, Debug: false},   // env config
    {Timeout: 60, Retries: 0, Debug: true},   // user config
}

final := F.Pipe2(
    configs,
    A.Fold(configMonoid),
)
// Config{Timeout: 60, Retries: 5, Debug: true}`}
</CodeCard>

</Section>

---

<Section num="6" title="Aggregating Statistics">

Combine statistics from multiple sources:

<CodeCard file="monoid_stats.go" tag="example">
{`type Stats struct {
    Count int
    Total float64
    Min   float64
    Max   float64
}

statsMonoid := M.MakeMonoid(
    func(a, b Stats) Stats {
        return Stats{
            Count: a.Count + b.Count,
            Total: a.Total + b.Total,
            Min:   min(a.Min, b.Min),
            Max:   max(a.Max, b.Max),
        }
    },
    Stats{
        Count: 0,
        Total: 0.0,
        Min:   math.MaxFloat64,
        Max:   -math.MaxFloat64,
    },
)

// Aggregate stats from multiple sources
allStats := []Stats{
    {Count: 10, Total: 100.0, Min: 5.0, Max: 20.0},
    {Count: 15, Total: 225.0, Min: 3.0, Max: 25.0},
    {Count: 8, Total: 80.0, Min: 7.0, Max: 15.0},
}

combined := F.Pipe2(
    allStats,
    A.Fold(statsMonoid),
)
// Stats{Count: 33, Total: 405.0, Min: 3.0, Max: 25.0}`}
</CodeCard>

</Section>

---

<Section num="7" title="API Reference">

<ApiTable>
| Function | Type | Description |
|----------|------|-------------|
| `MakeMonoid[A]` | `(func(A, A) A, A) -> Monoid[A]` | Creates a monoid from concat and empty |
| `Concat` | `(A, A) -> A` | Combines two values (from Semigroup) |
| `Empty` | `() -> A` | Returns the identity element |
</ApiTable>

**Built-in Monoids:**

<ApiTable>
| Monoid | Package | Operation | Identity |
|--------|---------|-----------|----------|
| `MonoidSum` | `number` | Addition | `0` |
| `MonoidProduct` | `number` | Multiplication | `1` |
| `MonoidMin` | `number` | Minimum | `MaxInt` |
| `MonoidMax` | `number` | Maximum | `MinInt` |
| `Monoid` | `string` | Concatenation | `""` |
</ApiTable>

</Section>

---

<Section num="8" title="Related Concepts">

**Algebraic Hierarchy:**
- **Magma** → **Semigroup** → **Monoid** → **Group**
- Monoid adds identity element to Semigroup
- Group adds inverse operation to Monoid

**Common Use Cases:**
- Folding/reducing collections
- Merging configurations
- Aggregating statistics
- Combining results from parallel operations

**See Also:**
- [Semigroup](./semigroup.md) - Monoid without identity element
- [Array Monoid](../collections/array-monoid.md) - Array folding operations
- [Record Monoid](../collections/record-monoid.md) - Map folding operations

</Section>

---

<Pager
  prev={{ to: '/docs/v2/utilities/semigroup', title: 'Semigroup' }}
  next={{ to: '/docs/v2/utilities/predicate', title: 'Predicate' }}
/>