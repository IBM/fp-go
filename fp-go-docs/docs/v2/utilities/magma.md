---
title: Magma
hide_title: true
description: Binary operations without constraints - the most basic algebraic structure.
sidebar_position: 29
---

import { PageHeader, Section, CodeCard, ApiTable, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="Reference · Utilities"
  title="Magma"
  lede="Binary operations without constraints. The most basic algebraic structure - just a way to combine two values."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/magma' },
    { label: 'Type Class', value: 'Magma[A]' }
  ]}
/>

---

<Section num="1" title="Core API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `MakeMagma` | `func MakeMagma[A any](func(A, A) A) Magma[A]` | Create magma from function |
</ApiTable>

</Section>

---

<Section num="2" title="Usage Examples">

### Basic Usage

<CodeCard file="basic.go">
{`import Mg "github.com/IBM/fp-go/v2/magma"

// Create magma for averaging
avgMagma := Mg.MakeMagma(func(a, b float64) float64 {
    return (a + b) / 2
})

result := avgMagma.Concat(10, 20)
// 15.0
`}
</CodeCard>

### Conflict Resolution Strategies

<CodeCard file="strategies.go">
{`// Last-wins strategy
lastWins := Mg.MakeMagma(func(a, b T) T {
    return b
})

// First-wins strategy
firstWins := Mg.MakeMagma(func(a, b T) T {
    return a
})

// Max strategy
maxMagma := Mg.MakeMagma(func(a, b int) int {
    if a > b { return a }
    return b
})

// Min strategy
minMagma := Mg.MakeMagma(func(a, b int) int {
    if a < b { return a }
    return b
})
`}
</CodeCard>

### Custom Merging

<CodeCard file="merge.go">
{`type User struct {
    ID        int
    Name      string
    UpdatedAt time.Time
}

// Merge users (keep most recent)
userMagma := Mg.MakeMagma(func(a, b User) User {
    if a.UpdatedAt.After(b.UpdatedAt) {
        return a
    }
    return b
})

merged := userMagma.Concat(oldUser, newUser)
// Returns the user with the later UpdatedAt
`}
</CodeCard>

### Array Merging

<CodeCard file="arrays.go">
{`// Concatenate arrays
appendMagma := Mg.MakeMagma(func(a, b []int) []int {
    return append(a, b...)
})

result := appendMagma.Concat([]int{1, 2}, []int{3, 4})
// []int{1, 2, 3, 4}

// Union (deduplicate)
unionMagma := Mg.MakeMagma(func(a, b []int) []int {
    seen := make(map[int]bool)
    result := make([]int, 0)
    
    for _, v := range append(a, b...) {
        if !seen[v] {
            seen[v] = true
            result = append(result, v)
        }
    }
    return result
})
`}
</CodeCard>

### Record Conversion

<CodeCard file="record.go">
{`import R "github.com/IBM/fp-go/v2/record"

// Use with FromArray for duplicate key handling
entries := []T.Tuple2[string, int]{
    T.MakeTuple2("a", 1),
    T.MakeTuple2("b", 2),
    T.MakeTuple2("a", 10),  // Duplicate
}

// Sum duplicates
sumMagma := Mg.MakeMagma(func(x, y int) int { return x + y })
m := R.FromArray(sumMagma)(entries)
// map[string]int{"a": 11, "b": 2}
`}
</CodeCard>

</Section>

---

<Section num="3" title="Common Patterns">

### Configuration Merging

<CodeCard file="config.go">
{`type Config struct {
    Host    string
    Port    int
    Timeout int
}

// Merge configs (right wins for non-zero values)
configMagma := Mg.MakeMagma(func(a, b Config) Config {
    return Config{
        Host:    if b.Host != "" { b.Host } else { a.Host },
        Port:    if b.Port != 0 { b.Port } else { a.Port },
        Timeout: if b.Timeout != 0 { b.Timeout } else { a.Timeout },
    }
})
`}
</CodeCard>

</Section>

---

<Callout type="info">

**Magma vs Semigroup**: Magma doesn't require associativity. Use Magma when the order of operations matters (like averaging), and Semigroup when it doesn't (like addition).

</Callout>

<Callout type="info">

**Use Cases**: Magma is perfect for:
- Conflict resolution strategies
- Custom merge logic
- Non-associative operations
- Building blocks for more complex structures

</Callout>


---

<Pager
  prev={{ to: '/docs/v2/utilities/function', title: 'Function Utilities' }}
  next={{ to: '/docs/v2/utilities/ord', title: 'Ord (Ordering)' }}
/>

---
