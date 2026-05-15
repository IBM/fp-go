---
title: Record - Conversion
hide_title: true
description: Converting between maps, arrays, and other data structures with flexible merge strategies.
sidebar_position: 18
---

import { PageHeader, Section, CodeCard, ApiTable, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="Reference · Collections"
  title="Record Conversion"
  lede="Converting between maps and other data structures. Build maps from arrays with custom merge strategies for duplicate keys."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/record' },
    { label: 'Operations', value: 'FromArray, FromArrayMap, ToArray, ToEntries' }
  ]}
/>

---

<Section num="1" title="Core API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `FromArray` | `func FromArray[K comparable, V any](Magma[V]) func([]Tuple2[K, V]) map[K]V` | Array to map |
| `FromArrayMap` | `func FromArrayMap[K comparable, A, V any](Magma[V]) func(func(A) Tuple2[K, V]) func([]A) map[K]V` | Array to map with transform |
| `ToArray` | `func ToArray[K comparable, V any](map[K]V) []Tuple2[K, V]` | Map to array |
| `ToEntries` | `func ToEntries[K comparable, V any](map[K]V) []Tuple2[K, V]` | Map to entries |
</ApiTable>

</Section>

---

<Section num="2" title="Usage Examples">

### FromArray - Basic

<CodeCard file="from_array.go">
{`import (
    R "github.com/IBM/fp-go/v2/record"
    T "github.com/IBM/fp-go/v2/tuple"
    Mg "github.com/IBM/fp-go/v2/magma"
)

entries := []T.Tuple2[string, int]{
    T.MakeTuple2("a", 1),
    T.MakeTuple2("b", 2),
    T.MakeTuple2("a", 10),  // Duplicate key
}

// Last value wins
lastWins := Mg.MakeMagma(func(_, y int) int { return y })
m := R.FromArray(lastWins)(entries)
// map[string]int{"a": 10, "b": 2}

// First value wins
firstWins := Mg.MakeMagma(func(x, _ int) int { return x })
m2 := R.FromArray(firstWins)(entries)
// map[string]int{"a": 1, "b": 2}

// Sum duplicates
sumMagma := Mg.MakeMagma(func(x, y int) int { return x + y })
m3 := R.FromArray(sumMagma)(entries)
// map[string]int{"a": 11, "b": 2}
`}
</CodeCard>

### FromArrayMap - Transform and Build

<CodeCard file="from_array_map.go">
{`type User struct {
    ID   int
    Name string
}

users := []User{
    {ID: 1, Name: "Alice"},
    {ID: 2, Name: "Bob"},
    {ID: 1, Name: "Alice Updated"},
}

// Convert to map by ID (last wins)
lastWins := Mg.MakeMagma(func(_, y User) User { return y })
userMap := R.FromArrayMap(lastWins)(
    func(u User) T.Tuple2[int, User] {
        return T.MakeTuple2(u.ID, u)
    },
)(users)
// map[int]User{
//   1: {ID: 1, Name: "Alice Updated"},
//   2: {ID: 2, Name: "Bob"},
// }
`}
</CodeCard>

### ToArray / ToEntries

<CodeCard file="to_array.go">
{`m := map[string]int{"a": 1, "b": 2, "c": 3}

// Convert to array of tuples
entries := R.ToEntries(m)
// []Tuple2[string, int]{
//   {Head: "a", Tail: 1},
//   {Head: "b", Tail: 2},
//   {Head: "c", Tail: 3},
// } (order may vary)

// ToArray is an alias for ToEntries
arr := R.ToArray(m)
// Same result
`}
</CodeCard>

### Grouping Data

<CodeCard file="grouping.go">
{`type Order struct {
    ID         int
    CustomerID int
    Amount     float64
}

orders := []Order{
    {ID: 1, CustomerID: 1, Amount: 100},
    {ID: 2, CustomerID: 2, Amount: 200},
    {ID: 3, CustomerID: 1, Amount: 150},
}

// Group by customer, sum amounts
sumMagma := Mg.MakeMagma(func(x, y float64) float64 { return x + y })
byCustomer := R.FromArrayMap(sumMagma)(
    func(o Order) T.Tuple2[int, float64] {
        return T.MakeTuple2(o.CustomerID, o.Amount)
    },
)(orders)
// map[int]float64{1: 250, 2: 200}
`}
</CodeCard>

### Building Indexes

<CodeCard file="index.go">
{`type Product struct {
    SKU  string
    Name string
}

products := []Product{
    {SKU: "A123", Name: "Laptop"},
    {SKU: "B456", Name: "Mouse"},
    {SKU: "A123", Name: "Laptop Pro"},  // Duplicate SKU
}

// Build index (first wins)
firstWins := Mg.MakeMagma(func(x, _ Product) Product { return x })
index := R.FromArrayMap(firstWins)(
    func(p Product) T.Tuple2[string, Product] {
        return T.MakeTuple2(p.SKU, p)
    },
)(products)
// map[string]Product{
//   "A123": {SKU: "A123", Name: "Laptop"},
//   "B456": {SKU: "B456", Name: "Mouse"},
// }
`}
</CodeCard>

### Collecting Lists

<CodeCard file="collect.go">
{`type Event struct {
    Type string
    Data string
}

events := []Event{
    {Type: "click", Data: "button1"},
    {Type: "hover", Data: "link1"},
    {Type: "click", Data: "button2"},
}

// Collect events by type
appendMagma := Mg.MakeMagma(func(x, y []string) []string {
    return append(x, y...)
})

byType := R.FromArrayMap(appendMagma)(
    func(e Event) T.Tuple2[string, []string] {
        return T.MakeTuple2(e.Type, []string{e.Data})
    },
)(events)
// map[string][]string{
//   "click": ["button1", "button2"],
//   "hover": ["link1"],
// }
`}
</CodeCard>

### Frequency Counting

<CodeCard file="frequency.go">
{`words := []string{"apple", "banana", "apple", "cherry", "banana", "apple"}

// Count occurrences
countMagma := Mg.MakeMagma(func(x, y int) int { return x + y })
frequencies := R.FromArrayMap(countMagma)(
    func(word string) T.Tuple2[string, int] {
        return T.MakeTuple2(word, 1)
    },
)(words)
// map[string]int{
//   "apple": 3,
//   "banana": 2,
//   "cherry": 1,
// }
`}
</CodeCard>

</Section>

---

<Section num="3" title="Common Patterns">

### Deduplication

<CodeCard file="dedup.go">
{`type Item struct {
    ID   string
    Data string
}

items := []Item{
    {ID: "1", Data: "first"},
    {ID: "2", Data: "second"},
    {ID: "1", Data: "duplicate"},
}

// Keep first occurrence
firstWins := Mg.MakeMagma(func(x, _ Item) Item { return x })
unique := R.FromArrayMap(firstWins)(
    func(i Item) T.Tuple2[string, Item] {
        return T.MakeTuple2(i.ID, i)
    },
)(items)

// Convert back to array
result := R.Values(unique)
// []Item{{ID: "1", Data: "first"}, {ID: "2", Data: "second"}}
`}
</CodeCard>

### Aggregation

<CodeCard file="aggregate.go">
{`type Sale struct {
    Product string
    Amount  float64
}

sales := []Sale{
    {Product: "laptop", Amount: 999},
    {Product: "mouse", Amount: 29},
    {Product: "laptop", Amount: 899},
}

// Total sales by product
sumMagma := Mg.MakeMagma(func(x, y float64) float64 { return x + y })
totals := R.FromArrayMap(sumMagma)(
    func(s Sale) T.Tuple2[string, float64] {
        return T.MakeTuple2(s.Product, s.Amount)
    },
)(sales)
// map[string]float64{"laptop": 1898, "mouse": 29}
`}
</CodeCard>

### Multi-value Mapping

<CodeCard file="multivalue.go">
{`type Tag struct {
    Resource string
    Tag      string
}

tags := []Tag{
    {Resource: "server1", Tag: "prod"},
    {Resource: "server1", Tag: "web"},
    {Resource: "server2", Tag: "prod"},
}

// Collect all tags per resource
appendMagma := Mg.MakeMagma(func(x, y []string) []string {
    return append(x, y...)
})

resourceTags := R.FromArrayMap(appendMagma)(
    func(t Tag) T.Tuple2[string, []string] {
        return T.MakeTuple2(t.Resource, []string{t.Tag})
    },
)(tags)
// map[string][]string{
//   "server1": ["prod", "web"],
//   "server2": ["prod"],
// }
`}
</CodeCard>

</Section>

---

<Callout type="info">

**Magma Strategy**: The magma parameter determines how duplicate keys are handled:
- **Last wins**: `func(_, y T) T { return y }`
- **First wins**: `func(x, _ T) T { return x }`
- **Sum**: `func(x, y int) int { return x + y }`
- **Append**: `func(x, y []T) []T { return append(x, y...) }`

</Callout>

<Callout type="info">

**Order**: `ToArray` and `ToEntries` return elements in undefined order since Go maps don't maintain insertion order. Use `record-ord` for ordered operations.

</Callout>


---

<Pager
  prev={{ to: '/docs/v2/collections/record-eq', title: 'Record Equality' }}
  next={{ to: '/docs/advanced/architecture', title: 'Architecture' }}
/>

---
