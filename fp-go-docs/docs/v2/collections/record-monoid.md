---
title: Record - Monoid
hide_title: true
description: Combining maps using monoid operations for powerful composition patterns.
sidebar_position: 13
---

import { PageHeader, Section, CodeCard, ApiTable, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="Reference · Collections"
  title="Record Monoid"
  lede="Combining maps using monoid operations. Records form monoids under various operations, enabling powerful composition patterns."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/record' },
    { label: 'Operations', value: 'Fold, FoldMap, Union' }
  ]}
/>

---

<Section num="1" title="Core API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Fold` | `func Fold[K comparable, V any](Monoid[V]) func(map[K]V) V` | Fold map to value |
| `FoldMap` | `func FoldMap[K comparable, A, B any](Monoid[B]) func(func(A) B) func(map[K]A) B` | Map then fold |
| `Union` | `func Union[K comparable, V any](Magma[V]) func(map[K]V) func(map[K]V) map[K]V` | Merge with function |
</ApiTable>

</Section>

---

<Section num="2" title="Usage Examples">

### Fold - Basic

<CodeCard file="fold.go">
{`import (
    R "github.com/IBM/fp-go/v2/record"
    M "github.com/IBM/fp-go/v2/monoid"
    F "github.com/IBM/fp-go/v2/function"
)

m := map[string]int{"a": 1, "b": 2, "c": 3}

// Sum all values
sum := F.Pipe2(
    m,
    R.Fold(M.MonoidSum[int]()),
)
// 6

// Product of all values
product := F.Pipe2(
    m,
    R.Fold(M.MonoidProduct[int]()),
)
// 6
`}
</CodeCard>

### FoldMap - Transform and Fold

<CodeCard file="foldmap.go">
{`type Product struct {
    Name  string
    Price float64
}

products := map[string]Product{
    "laptop": {Name: "Laptop", Price: 999},
    "mouse":  {Name: "Mouse", Price: 29},
    "keyboard": {Name: "Keyboard", Price: 79},
}

// Calculate total price
total := F.Pipe2(
    products,
    R.FoldMap(M.MonoidSum[float64]())(func(p Product) float64 {
        return p.Price
    }),
)
// 1107.0
`}
</CodeCard>

### String Concatenation

<CodeCard file="concat.go">
{`m := map[string]string{
    "first":  "Hello",
    "second": "World",
    "third":  "!",
}

// Concatenate all values
result := F.Pipe2(
    m,
    R.Fold(M.MonoidString),
)
// "HelloWorld!" (order may vary)

// With separator
import S "github.com/IBM/fp-go/v2/string"

separated := F.Pipe2(
    m,
    R.FoldMap(M.MonoidString)(func(s string) string {
        return s + " "
    }),
)
// "Hello World ! " (order may vary)
`}
</CodeCard>

### Union - Merge Maps

<CodeCard file="union.go">
{`m1 := map[string]int{"a": 1, "b": 2}
m2 := map[string]int{"b": 3, "c": 4}

// Merge with sum
merged := R.Union(M.MonoidSum[int]())(m2)(m1)
// map[string]int{"a": 1, "b": 5, "c": 4}

// Merge with max
import Mg "github.com/IBM/fp-go/v2/magma"

maxMagma := Mg.MakeMagma(func(x, y int) int {
    if x > y {
        return x
    }
    return y
})

maxMerged := R.Union(maxMagma)(m2)(m1)
// map[string]int{"a": 1, "b": 3, "c": 4}
`}
</CodeCard>

### Collecting Arrays

<CodeCard file="arrays.go">
{`import A "github.com/IBM/fp-go/v2/array"

m := map[string][]int{
    "group1": {1, 2, 3},
    "group2": {4, 5},
    "group3": {6},
}

// Concatenate all arrays
allValues := F.Pipe2(
    m,
    R.Fold(M.MonoidArray[int]()),
)
// []int{1, 2, 3, 4, 5, 6} (order may vary)
`}
</CodeCard>

### Combining Configurations

<CodeCard file="config.go">
{`type Config struct {
    MaxRetries int
    Timeout    int
}

configs := map[string]Config{
    "service1": {MaxRetries: 3, Timeout: 30},
    "service2": {MaxRetries: 5, Timeout: 60},
}

// Find max values
maxConfig := F.Pipe2(
    configs,
    R.FoldMap(
        M.MakeMonoid(
            func(a, b Config) Config {
                return Config{
                    MaxRetries: max(a.MaxRetries, b.MaxRetries),
                    Timeout:    max(a.Timeout, b.Timeout),
                }
            },
            Config{MaxRetries: 0, Timeout: 0},
        ),
    )(F.Identity[Config]),
)
// Config{MaxRetries: 5, Timeout: 60}
`}
</CodeCard>

### Aggregating Statistics

<CodeCard file="stats.go">
{`type Stats struct {
    Count int
    Sum   float64
}

data := map[string][]float64{
    "group1": {1.0, 2.0, 3.0},
    "group2": {4.0, 5.0},
}

// Calculate combined statistics
statsMonoid := M.MakeMonoid(
    func(a, b Stats) Stats {
        return Stats{
            Count: a.Count + b.Count,
            Sum:   a.Sum + b.Sum,
        }
    },
    Stats{Count: 0, Sum: 0},
)

combined := F.Pipe2(
    data,
    R.FoldMap(statsMonoid)(func(values []float64) Stats {
        sum := 0.0
        for _, v := range values {
            sum += v
        }
        return Stats{Count: len(values), Sum: sum}
    }),
)
// Stats{Count: 5, Sum: 15.0}
`}
</CodeCard>

</Section>

---

<Section num="3" title="Common Patterns">

### Merging Multiple Maps

<CodeCard file="merge_multiple.go">
{`maps := []map[string]int{
    {"a": 1, "b": 2},
    {"b": 3, "c": 4},
    {"c": 5, "d": 6},
}

// Merge all with sum
import A "github.com/IBM/fp-go/v2/array"

result := F.Pipe2(
    maps,
    A.Reduce(
        func(acc, m map[string]int) map[string]int {
            return R.Union(M.MonoidSum[int]())(m)(acc)
        },
        map[string]int{},
    ),
)
// map[string]int{"a": 1, "b": 5, "c": 9, "d": 6}
`}
</CodeCard>

### Accumulating Results

<CodeCard file="accumulate.go">
{`type Result struct {
    Success int
    Failed  int
}

results := map[string]Result{
    "batch1": {Success: 10, Failed: 2},
    "batch2": {Success: 15, Failed: 1},
    "batch3": {Success: 8, Failed: 3},
}

resultMonoid := M.MakeMonoid(
    func(a, b Result) Result {
        return Result{
            Success: a.Success + b.Success,
            Failed:  a.Failed + b.Failed,
        }
    },
    Result{Success: 0, Failed: 0},
)

total := F.Pipe2(
    results,
    R.Fold(resultMonoid),
)
// Result{Success: 33, Failed: 6}
`}
</CodeCard>

### Building Indexes

<CodeCard file="index.go">
{`type Index map[string][]string

indexes := map[string]Index{
    "doc1": {"word1": {"doc1"}, "word2": {"doc1"}},
    "doc2": {"word1": {"doc2"}, "word3": {"doc2"}},
}

// Merge indexes
appendMagma := Mg.MakeMagma(func(x, y []string) []string {
    return append(x, y...)
})

combined := F.Pipe2(
    indexes,
    R.Fold(R.Union(appendMagma)),
)
// Index{
//   "word1": ["doc1", "doc2"],
//   "word2": ["doc1"],
//   "word3": ["doc2"],
// }
`}
</CodeCard>

</Section>

---

<Callout type="info">

**Monoid Laws**: Monoids must satisfy associativity and have an identity element. This ensures that folding operations produce consistent results regardless of evaluation order.

</Callout>

<Callout type="info">

**Order Independence**: Since Go maps have undefined iteration order, fold operations should use commutative monoids for predictable results.

</Callout>


---

<Pager
  prev={{ to: '/docs/v2/collections/record-ord', title: 'Record Ordered' }}
  next={{ to: '/docs/v2/collections/record-traverse', title: 'Record Traverse' }}
/>

---
