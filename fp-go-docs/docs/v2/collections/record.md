---
title: Record (Map)
hide_title: true
description: Functional operations for Go maps with immutable transformations and type-safe lookups.
sidebar_position: 11
---

import { PageHeader, Section, CodeCard, ApiTable, Callout, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="Reference · Collections"
  title="Record (Map)"
  lede="Functional operations for Go maps. Treat maps as immutable data structures with type-safe lookups and transformations."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/record' },
    { label: 'Type', value: 'map[K]V' }
  ]}
/>

---

<Section num="1" title="Core API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Empty` | `func Empty[K comparable, V any]() map[K]V` | Create empty map |
| `Singleton` | `func Singleton[K comparable, V any](K, V) map[K]V` | Single entry map |
| `FromEntries` | `func FromEntries[K comparable, V any]([]Tuple2[K, V]) map[K]V` | From key-value pairs |
| `Lookup` | `func Lookup[V any](K) func(map[K]V) Option[V]` | Safe key access |
| `Has` | `func Has[K comparable, V any](K, map[K]V) bool` | Check key exists |
| `UpsertAt` | `func UpsertAt[K comparable, V any](K, V) func(map[K]V) map[K]V` | Add or update |
| `DeleteAt` | `func DeleteAt[V any](K) func(map[K]V) map[K]V` | Remove entry |
| `Map` | `func Map[K comparable, A, B any](func(A) B) func(map[K]A) map[K]B` | Transform values |
| `MapWithIndex` | `func MapWithIndex[K comparable, A, B any](func(K, A) B) func(map[K]A) map[K]B` | Transform with keys |
| `Filter` | `func Filter[K comparable, V any](func(V) bool) func(map[K]V) map[K]V` | Keep matching |
| `FilterWithIndex` | `func FilterWithIndex[K comparable, V any](func(K, V) bool) func(map[K]V) map[K]V` | Filter with keys |
| `FilterMap` | `func FilterMap[K comparable, A, B any](func(A) Option[B]) func(map[K]A) map[K]B` | Filter and transform |
| `Reduce` | `func Reduce[K comparable, A, B any](func(B, A) B, B) func(map[K]A) B` | Fold to value |
| `ReduceWithIndex` | `func ReduceWithIndex[K comparable, A, B any](func(K, B, A) B, B) func(map[K]A) B` | Reduce with keys |
| `Keys` | `func Keys[K comparable, V any](map[K]V) []K` | Get all keys |
| `Values` | `func Values[K comparable, V any](map[K]V) []V` | Get all values |
| `ToEntries` | `func ToEntries[K comparable, V any](map[K]V) []Tuple2[K, V]` | To key-value pairs |
| `Collect` | `func Collect[K comparable, A, B any](func(K, A) B) func(map[K]A) []B` | Transform to array |
| `Merge` | `func Merge[K comparable, V any](map[K]V) func(map[K]V) map[K]V` | Combine maps |
| `Union` | `func Union[K comparable, V any](Magma[V]) func(map[K]V) func(map[K]V) map[K]V` | Merge with function |
| `IsEmpty` | `func IsEmpty[K comparable, V any](map[K]V) bool` | Check if empty |
| `Size` | `func Size[K comparable, V any](map[K]V) int` | Get entry count |
</ApiTable>

</Section>

---

<Section num="2" title="Usage Examples">

### Creating Records

<CodeCard file="create.go">
{`import (
    R "github.com/IBM/fp-go/v2/record"
    T "github.com/IBM/fp-go/v2/tuple"
)

// Empty map
empty := R.Empty[string, int]()
// map[string]int{}

// Single entry
single := R.Singleton("key", 42)
// map[string]int{"key": 42}

// From entries
entries := []T.Tuple2[string, int]{
    T.MakeTuple2("a", 1),
    T.MakeTuple2("b", 2),
    T.MakeTuple2("c", 3),
}
m := R.FromEntries(entries)
// map[string]int{"a": 1, "b": 2, "c": 3}
`}
</CodeCard>

### Lookup and Access

<CodeCard file="lookup.go">
{`import (
    F "github.com/IBM/fp-go/v2/function"
    O "github.com/IBM/fp-go/v2/option"
)

m := map[string]int{"a": 1, "b": 2, "c": 3}

// Safe lookup - returns Option
value := F.Pipe2(
    m,
    R.Lookup[int]("b"),
)
// Some(2)

notFound := F.Pipe2(
    m,
    R.Lookup[int]("z"),
)
// None

// With default value
result := F.Pipe3(
    m,
    R.Lookup[int]("key"),
    O.GetOrElse(func() int { return 0 }),
)

// Check existence
R.Has("a", m)  // true
R.Has("z", m)  // false
`}
</CodeCard>

### Updating Records

<CodeCard file="update.go">
{`m := map[string]int{"a": 1, "b": 2}

// Add or update
updated := F.Pipe2(
    m,
    R.UpsertAt("c", 3),
)
// map[string]int{"a": 1, "b": 2, "c": 3}

// Replace existing
replaced := F.Pipe2(
    m,
    R.UpsertAt("a", 10),
)
// map[string]int{"a": 10, "b": 2}

// Remove entry
removed := F.Pipe2(
    m,
    R.DeleteAt[int]("b"),
)
// map[string]int{"a": 1}
`}
</CodeCard>

### Transformation

<CodeCard file="transform.go">
{`m := map[string]int{"a": 1, "b": 2, "c": 3}

// Map values
doubled := F.Pipe2(
    m,
    R.Map(func(v int) int { return v * 2 }),
)
// map[string]int{"a": 2, "b": 4, "c": 6}

// Map with keys
labeled := F.Pipe2(
    m,
    R.MapWithIndex(func(k string, v int) string {
        return fmt.Sprintf("%s=%d", k, v)
    }),
)
// map[string]string{"a": "a=1", "b": "b=2", "c": "c=3"}
`}
</CodeCard>

### Filtering

<CodeCard file="filter.go">
{`m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

// Filter by value
evens := F.Pipe2(
    m,
    R.Filter(func(v int) bool { return v%2 == 0 }),
)
// map[string]int{"b": 2, "d": 4}

// Filter with keys
filtered := F.Pipe2(
    m,
    R.FilterWithIndex(func(k string, v int) bool {
        return k != "a" && v > 1
    }),
)
// map[string]int{"b": 2, "c": 3, "d": 4}
`}
</CodeCard>

### FilterMap

<CodeCard file="filtermap.go">
{`m := map[string]string{
    "a": "123",
    "b": "abc",
    "c": "456",
}

// Parse numbers, skip invalid
numbers := F.Pipe2(
    m,
    R.FilterMap(func(s string) O.Option[int] {
        if n, err := strconv.Atoi(s); err == nil {
            return O.Some(n)
        }
        return O.None[int]()
    }),
)
// map[string]int{"a": 123, "c": 456}
`}
</CodeCard>

### Reduction

<CodeCard file="reduce.go">
{`m := map[string]int{"a": 1, "b": 2, "c": 3}

// Sum all values
sum := F.Pipe2(
    m,
    R.Reduce(func(acc, v int) int { return acc + v }, 0),
)
// 6

// Build string with keys
str := F.Pipe2(
    m,
    R.ReduceWithIndex(func(k string, acc string, v int) string {
        return acc + fmt.Sprintf("%s:%d ", k, v)
    }, ""),
)
// "a:1 b:2 c:3 " (order may vary)
`}
</CodeCard>

### Conversion

<CodeCard file="convert.go">
{`m := map[string]int{"a": 1, "b": 2, "c": 3}

// Get keys
keys := R.Keys(m)
// []string{"a", "b", "c"} (order may vary)

// Get values
values := R.Values(m)
// []int{1, 2, 3} (order may vary)

// To entries
entries := R.ToEntries(m)
// []Tuple2[string, int]{
//   {Head: "a", Tail: 1},
//   {Head: "b", Tail: 2},
//   {Head: "c", Tail: 3},
// }

// Collect to array
pairs := F.Pipe2(
    m,
    R.Collect(func(k string, v int) string {
        return fmt.Sprintf("%s=%d", k, v)
    }),
)
// []string{"a=1", "b=2", "c=3"} (order may vary)
`}
</CodeCard>

### Combining Records

<CodeCard file="combine.go">
{`m1 := map[string]int{"a": 1, "b": 2}
m2 := map[string]int{"b": 20, "c": 3}

// Merge - right wins
merged := F.Pipe2(
    m1,
    R.Merge(m2),
)
// map[string]int{"a": 1, "b": 20, "c": 3}

// Union with custom merge
import Mg "github.com/IBM/fp-go/v2/magma"

sumMagma := Mg.MakeMagma(func(x, y int) int { return x + y })

combined := F.Pipe2(
    m1,
    R.Union(sumMagma)(m2),
)
// map[string]int{"a": 1, "b": 22, "c": 3}
`}
</CodeCard>

</Section>

---

<Section num="3" title="Common Patterns">

### Configuration Management

<CodeCard file="config.go">
{`type Config struct {
    Host    string
    Port    int
    Timeout int
}

defaults := map[string]Config{
    "dev":  {Host: "localhost", Port: 8080, Timeout: 30},
    "prod": {Host: "api.example.com", Port: 443, Timeout: 60},
}

// Get config with fallback
config := F.Pipe3(
    defaults,
    R.Lookup[Config]("staging"),
    O.GetOrElse(func() Config { return defaults["dev"] }),
)
`}
</CodeCard>

### Grouping Data

<CodeCard file="grouping.go">
{`type User struct {
    ID   int
    Name string
    Role string
}

users := []User{
    {ID: 1, Name: "Alice", Role: "admin"},
    {ID: 2, Name: "Bob", Role: "user"},
    {ID: 3, Name: "Charlie", Role: "admin"},
}

// Group by role
import A "github.com/IBM/fp-go/v2/array"

byRole := F.Pipe2(
    users,
    A.Reduce(func(acc map[string][]User, u User) map[string][]User {
        acc[u.Role] = append(acc[u.Role], u)
        return acc
    }, map[string][]User{}),
)
// map[string][]User{
//   "admin": [{Alice}, {Charlie}],
//   "user": [{Bob}],
// }
`}
</CodeCard>

### Caching

<CodeCard file="cache.go">
{`type Cache map[string]interface{}

cache := R.Empty[string, interface{}]()

// Add to cache
cache = F.Pipe2(
    cache,
    R.UpsertAt("user:1", User{ID: 1, Name: "Alice"}),
)

// Lookup from cache
user := F.Pipe3(
    cache,
    R.Lookup[interface{}]("user:1"),
    O.Map(func(v interface{}) User { return v.(User) }),
    O.GetOrElse(func() User { return fetchUser(1) }),
)
`}
</CodeCard>

### Transform and Filter Chain

<CodeCard file="chain.go">
{`result := F.Pipe3(
    m,
    R.Map(transform),
    R.FilterWithIndex(predicate),
)
`}
</CodeCard>

### Merge Multiple Maps

<CodeCard file="merge_multiple.go">
{`merged := F.Pipe3(
    map1,
    R.Merge(map2),
    R.Merge(map3),
)
`}
</CodeCard>

### Checking Records

<CodeCard file="check.go">
{`m := map[string]int{"a": 1, "b": 2}
empty := map[string]int{}

R.IsEmpty(m)      // false
R.IsEmpty(empty)  // true
R.IsNonEmpty(m)   // true
R.Size(m)         // 2
`}
</CodeCard>

</Section>

---

<Callout type="warn">

**Map Iteration Order**: Go maps have undefined iteration order. Operations like `Keys`, `Values`, and `ToEntries` may return elements in different orders. Use `record-ord` for ordered operations.

</Callout>

<Callout type="info">

**Immutability**: All record operations return new maps without modifying the original. This ensures referential transparency and makes code easier to reason about.

</Callout>


---

<Pager
  prev={{ to: '/docs/v2/collections/sequence-traverse', title: 'Sequence & Traverse' }}
  next={{ to: '/docs/v2/collections/record-ord', title: 'Record Ordered' }}
/>

