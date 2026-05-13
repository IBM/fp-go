---
title: Semigroup
hide_title: true
description: Combining values associatively with binary operations in fp-go
sidebar_position: 27
---

import { PageHeader, Section, CodeCard, ApiTable, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="v2 · Utilities"
  title="Semigroup"
  titleAccent="Combining Values"
  lede="A semigroup is a type with an associative binary operation that combines two values of the same type."
  meta={[
    { label: 'Package', value: 'semigroup' },
    { label: 'Since', value: 'v2.0.0' },
    { label: 'Extends', value: 'Magma' }
  ]}
/>

---

<Section num="1" title="Overview">

A **Semigroup** is an algebraic structure consisting of a set and an associative binary operation. The key property is **associativity**: `(a • b) • c = a • (b • c)`.

**Key Properties:**
- **Associative**: Order of operations doesn't matter
- **Closed**: Combining two values produces another value of the same type
- **Type-safe**: Enforced through Go's type system

</Section>

---

<Section num="2" title="Basic Usage">

<CodeCard file="semigroup_basic.go" tag="example">
{`import S "github.com/IBM/fp-go/semigroup"

// Create a string concatenation semigroup
stringSemigroup := S.MakeSemigroup(func(a, b string) string {
    return a + b
})

result := stringSemigroup.Concat("Hello", " World")
// "Hello World"

// Associativity holds
s1 := stringSemigroup.Concat(
    stringSemigroup.Concat("a", "b"),
    "c",
)
s2 := stringSemigroup.Concat(
    "a",
    stringSemigroup.Concat("b", "c"),
)
// s1 == s2 == "abc"`}
</CodeCard>

</Section>

---

<Section num="3" title="Built-in Semigroups">

fp-go provides several built-in semigroups for common types:

<CodeCard file="semigroup_builtin.go" tag="example">
{`import (
    N "github.com/IBM/fp-go/number"
    S "github.com/IBM/fp-go/string"
)

// Number addition semigroup
sum := N.SemigroupSum.Concat(1, 2)  // 3

// Number multiplication semigroup
product := N.SemigroupProduct.Concat(3, 4)  // 12

// String concatenation semigroup
text := S.Semigroup.Concat("Hello", " World")  // "Hello World"

// Min/Max semigroups
min := N.SemigroupMin.Concat(5, 3)  // 3
max := N.SemigroupMax.Concat(5, 3)  // 5`}
</CodeCard>

</Section>

---

<Section num="4" title="Custom Semigroups">

Create semigroups for your own types:

<CodeCard file="semigroup_custom.go" tag="example">
{`type Config struct {
    Timeout int
    Retries int
    Debug   bool
}

// Last-wins merge strategy
configSemigroup := S.MakeSemigroup(func(a, b Config) Config {
    return Config{
        Timeout: b.Timeout,
        Retries: b.Retries,
        Debug:   b.Debug,
    }
})

defaults := Config{Timeout: 30, Retries: 3, Debug: false}
userConfig := Config{Timeout: 60, Retries: 5, Debug: true}

merged := configSemigroup.Concat(defaults, userConfig)
// Config{Timeout: 60, Retries: 5, Debug: true}

// Field-wise merge strategy
fieldwiseSemigroup := S.MakeSemigroup(func(a, b Config) Config {
    return Config{
        Timeout: max(a.Timeout, b.Timeout),
        Retries: max(a.Retries, b.Retries),
        Debug:   a.Debug || b.Debug,
    }
})`}
</CodeCard>

</Section>

---

<Section num="5" title="Combining Multiple Values">

Use semigroups to combine multiple values:

<CodeCard file="semigroup_multiple.go" tag="example">
{`import (
    A "github.com/IBM/fp-go/array"
    S "github.com/IBM/fp-go/semigroup"
)

// Combine array of strings
stringSemigroup := S.MakeSemigroup(func(a, b string) string {
    return a + b
})

words := []string{"Hello", " ", "functional", " ", "world"}
result := A.Reduce(
    func(acc, word string) string {
        return stringSemigroup.Concat(acc, word)
    },
    "",
)(words)
// "Hello functional world"

// Combine array of numbers
numbers := []int{1, 2, 3, 4, 5}
sum := A.Reduce(
    func(acc, n int) int {
        return N.SemigroupSum.Concat(acc, n)
    },
    0,
)(numbers)
// 15`}
</CodeCard>

</Section>

---

<Section num="6" title="API Reference">

<ApiTable>
| Function | Type | Description |
|----------|------|-------------|
| `MakeSemigroup[A]` | `(func(A, A) A) -> Semigroup[A]` | Creates a semigroup from a binary operation |
| `Concat` | `(A, A) -> A` | Combines two values associatively |
</ApiTable>

**Built-in Semigroups:**

<ApiTable>
| Semigroup | Package | Operation |
|-----------|---------|-----------|
| `SemigroupSum` | `number` | Addition |
| `SemigroupProduct` | `number` | Multiplication |
| `SemigroupMin` | `number` | Minimum |
| `SemigroupMax` | `number` | Maximum |
| `Semigroup` | `string` | Concatenation |
</ApiTable>

</Section>

---

<Section num="7" title="Related Concepts">

**Algebraic Hierarchy:**
- **Magma** → **Semigroup** → **Monoid** → **Group**
- Semigroup adds associativity to Magma
- Monoid adds identity element to Semigroup

**See Also:**
- [Monoid](./monoid.md) - Semigroup with identity element
- [Magma](./magma.md) - Binary operation without associativity requirement

</Section>

---

<Pager
  prev={{ to: '/docs/v2/utilities/predicate', title: 'Predicate' }}
  next={{ to: '/docs/v2/utilities/monoid', title: 'Monoid' }}
/>