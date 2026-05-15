---
title: Array
hide_title: true
description: Functional operations for Go slices - map, filter, reduce, and more with immutable data structures.
sidebar_position: 1
---

<PageHeader
  eyebrow="Reference · Collections"
  title="Array"
  lede="Functional operations for Go slices. Comprehensive set of operations treating slices as immutable data structures with type-safe transformations."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/array' },
    { label: 'Type', value: '[]T' }
  ]}
/>

<Section id="overview" number="01" title="Overview">

The array package provides functional operations for Go slices:
- **Immutable**: Always returns new slices
- **Type-safe**: Leverages Go generics
- **Composable**: Chain operations together

<CodeCard file="basic.go">
{`import A "github.com/IBM/fp-go/v2/array"

// Create arrays
numbers := A.From(1, 2, 3, 4, 5)
single := A.Of(42)
empty := A.Empty[int]()

// Transform
doubled := A.Map(func(n int) int { return n * 2 })(numbers)
// []int{2, 4, 6, 8, 10}

// Filter
evens := A.Filter(func(n int) bool { return n%2 == 0 })(numbers)
// []int{2, 4}

// Reduce
sum := A.Reduce(func(acc, n int) int { return acc + n }, 0)(numbers)
// 15
`}
</CodeCard>

</Section>

<Section id="api" number="02" title="Core" titleAccent="API">

### Constructors

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `From` | `func From[A any](values ...A) []A` | Create from variadic args |
| `Of` | `func Of[A any](value A) []A` | Single element array |
| `Empty` | `func Empty[A any]() []A` | Empty array |
| `MakeBy` | `func MakeBy[A any](n int, f func(int) A) []A` | Generate with function |
| `Replicate` | `func Replicate[A any](n int, value A) []A` | Repeat value n times |
</ApiTable>

### Transformations

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Map` | `func Map[A, B any](f func(A) B) func([]A) []B` | Transform each element |
| `MapWithIndex` | `func MapWithIndex[A, B any](f func(int, A) B) func([]A) []B` | Map with index |
| `Chain` | `func Chain[A, B any](f func(A) []B) func([]A) []B` | FlatMap |
| `FilterMap` | `func FilterMap[A, B any](f func(A) Option[B]) func([]A) []B` | Filter and map |
</ApiTable>

### Filtering

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Filter` | `func Filter[A any](pred func(A) bool) func([]A) []A` | Keep matching elements |
| `FilterWithIndex` | `func FilterWithIndex[A any](pred func(int, A) bool) func([]A) []A` | Filter with index |
| `Partition` | `func Partition[A any](pred func(A) bool) func([]A) Pair[[]A, []A]` | Split into two arrays |
</ApiTable>

### Reduction

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Reduce` | `func Reduce[A, B any](f func(B, A) B, initial B) func([]A) B` | Fold left |
| `ReduceRight` | `func ReduceRight[A, B any](f func(A, B) B, initial B) func([]A) B` | Fold right |
| `ReduceWithIndex` | `func ReduceWithIndex[A, B any](f func(int, B, A) B, initial B) func([]A) B` | Reduce with index |
</ApiTable>

### Access

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Head` | `func Head[A any]([]A) Option[A]` | First element |
| `Last` | `func Last[A any]([]A) Option[A]` | Last element |
| `Tail` | `func Tail[A any]([]A) Option[[]A]` | All but first |
| `Lookup` | `func Lookup[A any](index int) func([]A) Option[A]` | Element at index |
</ApiTable>

### Combining

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Append` | `func Append[A any](arr []A, value A) []A` | Add to end |
| `Prepend` | `func Prepend[A any](value A) func([]A) []A` | Add to start |
| `Flatten` | `func Flatten[A any]([][]A) []A` | Flatten nested arrays |
| `Concat` | `func Concat[A any](arrays ...[]A) []A` | Concatenate arrays |
</ApiTable>

### Checking

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `IsEmpty` | `func IsEmpty[A any]([]A) bool` | Check if empty |
| `IsNonEmpty` | `func IsNonEmpty[A any]([]A) bool` | Check if non-empty |
| `Size` | `func Size[A any]([]A) int` | Get length |
| `Elem` | `func Elem[A any](eq Eq[A]) func(A) func([]A) bool` | Check if contains |
</ApiTable>

</Section>

<Section id="examples" number="03" title="Usage" titleAccent="Examples">

### Basic Transformations

<CodeCard file="transformations.go">
{`import (
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
)

numbers := []int{1, 2, 3, 4, 5}

// Map
doubled := F.Pipe2(
    numbers,
    A.Map(func(n int) int { return n * 2 }),
)
// []int{2, 4, 6, 8, 10}

// Filter
evens := F.Pipe2(
    numbers,
    A.Filter(func(n int) bool { return n%2 == 0 }),
)
// []int{2, 4}

// Reduce
sum := F.Pipe2(
    numbers,
    A.Reduce(func(acc, n int) int { return acc + n }, 0),
)
// 15
`}
</CodeCard>

### Chain (FlatMap)

<CodeCard file="chain.go">
{`type User struct {
    Name  string
    Roles []string
}

users := []User{
    {Name: "Alice", Roles: []string{"admin", "user"}},
    {Name: "Bob", Roles: []string{"user"}},
}

// Get all roles
allRoles := F.Pipe2(
    users,
    A.Chain(func(u User) []string { return u.Roles }),
)
// []string{"admin", "user", "user"}
`}
</CodeCard>

### FilterMap

<CodeCard file="filtermap.go">
{`import O "github.com/IBM/fp-go/v2/option"

numbers := []int{1, 2, 3, 4, 5}

// Keep even numbers and double them
result := F.Pipe2(
    numbers,
    A.FilterMap(func(n int) O.Option[int] {
        if n%2 == 0 {
            return O.Some(n * 2)
        }
        return O.None[int]()
    }),
)
// []int{4, 8}
`}
</CodeCard>

### Partition

<CodeCard file="partition.go">
{`numbers := []int{1, 2, 3, 4, 5, 6}

// Separate evens and odds
result := F.Pipe2(
    numbers,
    A.Partition(func(n int) bool { return n%2 == 0 }),
)

odds := result.Head   // []int{1, 3, 5}
evens := result.Tail  // []int{2, 4, 6}
`}
</CodeCard>

### Real-World Example

<CodeCard file="products.go">
{`type Product struct {
    ID       int
    Name     string
    Price    float64
    Category string
    InStock  bool
}

products := []Product{
    {ID: 1, Name: "Laptop", Price: 999.99, Category: "Electronics", InStock: true},
    {ID: 2, Name: "Mouse", Price: 29.99, Category: "Electronics", InStock: true},
    {ID: 3, Name: "Desk", Price: 299.99, Category: "Furniture", InStock: false},
    {ID: 4, Name: "Chair", Price: 199.99, Category: "Furniture", InStock: true},
}

// Get names of in-stock electronics under $500
result := F.Pipe3(
    products,
    A.Filter(func(p Product) bool {
        return p.InStock && p.Category == "Electronics" && p.Price < 500
    }),
    A.Map(func(p Product) string { return p.Name }),
    A.Sort(ord.FromCompare(strings.Compare)),
)
// []string{"Mouse"}

// Calculate total value of in-stock items
totalValue := F.Pipe3(
    products,
    A.Filter(func(p Product) bool { return p.InStock }),
    A.Reduce(func(sum float64, p Product) float64 {
        return sum + p.Price
    }, 0.0),
)
// 1429.97
`}
</CodeCard>

</Section>

<Section id="patterns" number="04" title="Common" titleAccent="Patterns">

### Safe Head

<CodeCard file="safe_head.go">
{`// Instead of: arr[0] (panics if empty)
head := F.Pipe2(
    arr,
    A.Head,
    O.GetOrElse(func() int { return 0 }),
)
`}
</CodeCard>

### Transform and Aggregate

<CodeCard file="aggregate.go">
{`// Map then reduce
result := F.Pipe4(
    numbers,
    A.Map(func(n int) int { return n * n }),
    A.Filter(func(n int) bool { return n > 10 }),
    A.Reduce(func(sum, n int) int { return sum + n }, 0),
)
`}
</CodeCard>

### Conditional Transformation

<CodeCard file="conditional.go">
{`// Use FilterMap for conditional transforms
result := F.Pipe2(
    items,
    A.FilterMap(func(item Item) O.Option[Result] {
        if item.IsValid() {
            return O.Some(item.Transform())
        }
        return O.None[Result]()
    }),
)
`}
</CodeCard>

</Section>

<Callout type="info">

**Best Practices**:
- Use `Pipe` for readability when chaining operations
- Prefer immutability - array functions return new slices
- Combine operations with `FilterMap` instead of separate `Filter` + `Map`
- Check for empty with `IsEmpty` before accessing elements

</Callout>
