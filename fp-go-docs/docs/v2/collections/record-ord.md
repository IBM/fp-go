---
title: Record - Ordered
hide_title: true
description: Working with maps in a specific key order using the Ord type class.
sidebar_position: 12
---

import { PageHeader, Section, CodeCard, ApiTable, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="Reference · Collections"
  title="Record Ordered"
  lede="Working with maps in a specific key order. Process map entries in sorted order using the Ord type class."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/record' },
    { label: 'Operations', value: 'KeysOrd, ValuesOrd, ReduceOrd, CollectOrd' }
  ]}
/>

---

<Section num="1" title="Core API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `KeysOrd` | `func KeysOrd[K, V any](Ord[K]) func(map[K]V) []K` | Get keys in order |
| `ValuesOrd` | `func ValuesOrd[K, V any](Ord[K]) func(map[K]V) []V` | Get values by key order |
| `ReduceOrd` | `func ReduceOrd[K, A, B any](Ord[K]) func(func(B, A) B, B) func(map[K]A) B` | Reduce in order |
| `ReduceOrdWithIndex` | `func ReduceOrdWithIndex[K, A, B any](Ord[K]) func(func(K, B, A) B, B) func(map[K]A) B` | Reduce with keys in order |
| `CollectOrd` | `func CollectOrd[K, A, B any](Ord[K]) func(func(K, A) B) func(map[K]A) []B` | Transform to array in order |
</ApiTable>

</Section>

---

<Section num="2" title="Usage Examples">

### KeysOrd - Sorted Keys

<CodeCard file="keys_ord.go">
{`import (
    R "github.com/IBM/fp-go/v2/record"
    S "github.com/IBM/fp-go/v2/string"
    F "github.com/IBM/fp-go/v2/function"
)

m := map[string]int{"c": 3, "a": 1, "b": 2}

// Get keys in alphabetical order
keys := F.Pipe2(
    m,
    R.KeysOrd(S.Ord),
)
// []string{"a", "b", "c"}
`}
</CodeCard>

### ValuesOrd - Values by Key Order

<CodeCard file="values_ord.go">
{`m := map[string]int{"c": 3, "a": 1, "b": 2}

// Get values ordered by their keys
values := F.Pipe2(
    m,
    R.ValuesOrd(S.Ord),
)
// []int{1, 2, 3}  // Ordered by keys: a, b, c
`}
</CodeCard>

### ReduceOrd - Ordered Reduction

<CodeCard file="reduce_ord.go">
{`m := map[string]int{"c": 3, "a": 1, "b": 2}

// Build string in alphabetical order
str := F.Pipe2(
    m,
    R.ReduceOrdWithIndex(S.Ord)(
        func(k string, acc string, v int) string {
            return acc + fmt.Sprintf("%s:%d ", k, v)
        },
        "",
    ),
)
// "a:1 b:2 c:3 "
`}
</CodeCard>

### CollectOrd - Transform in Order

<CodeCard file="collect_ord.go">
{`m := map[string]int{"c": 3, "a": 1, "b": 2}

// Transform to array in key order
pairs := F.Pipe2(
    m,
    R.CollectOrd(S.Ord)(func(k string, v int) string {
        return fmt.Sprintf("%s=%d", k, v)
    }),
)
// []string{"a=1", "b=2", "c=3"}
`}
</CodeCard>

### Numeric Key Ordering

<CodeCard file="numeric.go">
{`import N "github.com/IBM/fp-go/v2/number"

m := map[int]string{
    3: "three",
    1: "one",
    2: "two",
}

// Get keys in numeric order
keys := F.Pipe2(
    m,
    R.KeysOrd(N.Ord),
)
// []int{1, 2, 3}

// Get values in key order
values := F.Pipe2(
    m,
    R.ValuesOrd(N.Ord),
)
// []string{"one", "two", "three"}
`}
</CodeCard>

### Custom Ordering

<CodeCard file="custom_ord.go">
{`import O "github.com/IBM/fp-go/v2/ord"

type Priority int

const (
    Low    Priority = 1
    Medium Priority = 2
    High   Priority = 3
)

// Custom ordering: High > Medium > Low
priorityOrd := O.FromCompare(func(a, b Priority) int {
    return int(b - a)  // Reverse order
})

tasks := map[Priority]string{
    Low:    "Documentation",
    High:   "Critical Bug",
    Medium: "Feature Request",
}

// Get tasks in priority order
ordered := F.Pipe2(
    tasks,
    R.ValuesOrd(priorityOrd),
)
// []string{"Critical Bug", "Feature Request", "Documentation"}
`}
</CodeCard>

### Configuration Processing

<CodeCard file="config.go">
{`type Config struct {
    Key   string
    Value string
}

configs := map[string]Config{
    "database": {Key: "db", Value: "postgres"},
    "cache":    {Key: "cache", Value: "redis"},
    "api":      {Key: "api", Value: "rest"},
}

// Process in alphabetical order
ordered := F.Pipe2(
    configs,
    R.CollectOrd(S.Ord)(func(k string, c Config) string {
        return fmt.Sprintf("%s: %s=%s", k, c.Key, c.Value)
    }),
)
// []string{
//   "api: api=rest",
//   "cache: cache=redis",
//   "database: db=postgres",
// }
`}
</CodeCard>

</Section>

---

<Section num="3" title="Common Patterns">

### Deterministic Output

<CodeCard file="deterministic.go">
{`// Generate consistent output regardless of map iteration order
func FormatConfig(config map[string]string) string {
    return F.Pipe2(
        config,
        R.ReduceOrdWithIndex(S.Ord)(
            func(k string, acc string, v string) string {
                return acc + fmt.Sprintf("%s=%s\n", k, v)
            },
            "",
        ),
    )
}

config := map[string]string{
    "port":    "8080",
    "host":    "localhost",
    "timeout": "30",
}

output := FormatConfig(config)
// Always produces:
// "host=localhost\nport=8080\ntimeout=30\n"
`}
</CodeCard>

### Sorted Aggregation

<CodeCard file="aggregate.go">
{`type Stats struct {
    Count int
    Total float64
}

data := map[string]Stats{
    "2023-03": {Count: 100, Total: 1500},
    "2023-01": {Count: 80, Total: 1200},
    "2023-02": {Count: 90, Total: 1350},
}

// Process in chronological order
report := F.Pipe2(
    data,
    R.CollectOrd(S.Ord)(func(month string, s Stats) string {
        avg := s.Total / float64(s.Count)
        return fmt.Sprintf("%s: %.2f avg", month, avg)
    }),
)
// []string{
//   "2023-01: 15.00 avg",
//   "2023-02: 15.00 avg",
//   "2023-03: 15.00 avg",
// }
`}
</CodeCard>

### Priority Queue Simulation

<CodeCard file="priority.go">
{`type Task struct {
    Name     string
    Priority int
}

tasks := map[int]Task{
    1: {Name: "Low priority", Priority: 1},
    5: {Name: "High priority", Priority: 5},
    3: {Name: "Medium priority", Priority: 3},
}

// Process by priority (highest first)
import O "github.com/IBM/fp-go/v2/ord"

reverseOrd := O.Reverse(N.Ord)

ordered := F.Pipe2(
    tasks,
    R.ValuesOrd(reverseOrd),
)
// []Task{
//   {Name: "High priority", Priority: 5},
//   {Name: "Medium priority", Priority: 3},
//   {Name: "Low priority", Priority: 1},
// }
`}
</CodeCard>

### Building Sorted Index

<CodeCard file="index.go">
{`type Document struct {
    ID    int
    Title string
}

docs := map[int]Document{
    3: {ID: 3, Title: "Third"},
    1: {ID: 1, Title: "First"},
    2: {ID: 2, Title: "Second"},
}

// Create sorted index
index := F.Pipe2(
    docs,
    R.CollectOrd(N.Ord)(func(id int, doc Document) string {
        return fmt.Sprintf("[%d] %s", id, doc.Title)
    }),
)
// []string{"[1] First", "[2] Second", "[3] Third"}
`}
</CodeCard>

</Section>

---

<Callout type="info">

**Deterministic Output**: Use ordered operations when you need consistent, reproducible output. This is especially important for testing, logging, and user-facing displays.

</Callout>

<Callout type="info">

**Performance**: Ordered operations require sorting, which adds O(n log n) complexity. Use regular operations when order doesn't matter.

</Callout>


---

<Pager
  prev={{ to: '/docs/v2/collections/record', title: 'Record (Map)' }}
  next={{ to: '/docs/v2/collections/record-monoid', title: 'Record Monoid' }}
/>

---
