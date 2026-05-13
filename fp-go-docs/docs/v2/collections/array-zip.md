---
title: Array - Zip
hide_title: true
description: Combining arrays element-wise with Zip, ZipWith, and Unzip operations.
sidebar_position: 4
---

<PageHeader
  eyebrow="Reference · Collections"
  title="Array"
  titleAccent="Zip"
  lede="Combining arrays element-wise. Zip operations pair elements at the same index for parallel processing and data alignment."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/array' },
    { label: 'Operations', value: 'Zip, ZipWith, Unzip' }
  ]}
/>

<Section id="api" number="01" title="Core" titleAccent="API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Zip` | `func Zip[A, B any](bs []B) func([]A) []Tuple2[A, B]` | Combine into pairs |
| `ZipWith` | `func ZipWith[A, B, C any](as []A, bs []B, f func(A, B) C) []C` | Combine with function |
| `Unzip` | `func Unzip[A, B any]([]Tuple2[A, B]) Tuple2[[]A, []B]` | Split pairs into arrays |
</ApiTable>

</Section>

<Section id="examples" number="02" title="Usage" titleAccent="Examples">

### Zip - Basic

<CodeCard file="zip.go">
{`import (
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
    T "github.com/IBM/fp-go/v2/tuple"
)

names := []string{"Alice", "Bob", "Charlie"}
ages := []int{30, 25, 35}

// Zip into tuples
pairs := F.Pipe2(
    names,
    A.Zip(ages),
)
// []Tuple2[string, int]{
//   {Head: "Alice", Tail: 30},
//   {Head: "Bob", Tail: 25},
//   {Head: "Charlie", Tail: 35},
// }

// Shorter array determines result length
short := []string{"A", "B"}
long := []int{1, 2, 3, 4, 5}

result := F.Pipe2(short, A.Zip(long))
// []Tuple2[string, int]{{A, 1}, {B, 2}}
// Elements 3, 4, 5 are discarded
`}
</CodeCard>

### ZipWith - Custom Function

<CodeCard file="zipwith.go">
{`// Calculate totals
prices := []float64{10.0, 20.0, 30.0}
quantities := []int{2, 3, 1}

totals := A.ZipWith(
    prices,
    quantities,
    func(price float64, qty int) float64 {
        return price * float64(qty)
    },
)
// []float64{20.0, 60.0, 30.0}

// Combine names
firstNames := []string{"John", "Jane", "Bob"}
lastNames := []string{"Doe", "Smith", "Johnson"}

fullNames := A.ZipWith(
    firstNames,
    lastNames,
    func(first, last string) string {
        return first + " " + last
    },
)
// []string{"John Doe", "Jane Smith", "Bob Johnson"}
`}
</CodeCard>

### Create Structs

<CodeCard file="structs.go">
{`type User struct {
    Name  string
    Email string
}

names := []string{"Alice", "Bob"}
emails := []string{"alice@example.com", "bob@example.com"}

users := A.ZipWith(
    names,
    emails,
    func(name, email string) User {
        return User{Name: name, Email: email}
    },
)
// []User{
//   {Name: "Alice", Email: "alice@example.com"},
//   {Name: "Bob", Email: "bob@example.com"},
// }
`}
</CodeCard>

### Unzip

<CodeCard file="unzip.go">
{`pairs := []T.Tuple2[string, int]{
    {Head: "Alice", Tail: 30},
    {Head: "Bob", Tail: 25},
    {Head: "Charlie", Tail: 35},
}

result := A.Unzip(pairs)
names := result.Head   // []string{"Alice", "Bob", "Charlie"}
ages := result.Tail    // []int{30, 25, 35}
`}
</CodeCard>

### Parallel Data Processing

<CodeCard file="parallel.go">
{`type Product struct {
    ID    int
    Name  string
    Price float64
}

type Discount struct {
    ProductID int
    Percent   float64
}

products := []Product{
    {ID: 1, Name: "Laptop", Price: 1000},
    {ID: 2, Name: "Mouse", Price: 50},
}

discounts := []Discount{
    {ProductID: 1, Percent: 10},
    {ProductID: 2, Percent: 20},
}

// Apply discounts
discounted := A.ZipWith(
    products,
    discounts,
    func(p Product, d Discount) Product {
        return Product{
            ID:    p.ID,
            Name:  p.Name,
            Price: p.Price * (1 - d.Percent/100),
        }
    },
)
// Products with discounts applied
`}
</CodeCard>

### Comparing Arrays

<CodeCard file="compare.go">
{`expected := []int{1, 2, 3, 4, 5}
actual := []int{1, 2, 4, 4, 5}

// Find differences
differences := A.ZipWith(
    expected,
    actual,
    func(exp, act int) string {
        if exp == act {
            return "✓"
        }
        return fmt.Sprintf("✗ expected %d, got %d", exp, act)
    },
)
// []string{"✓", "✓", "✗ expected 3, got 4", "✓", "✓"}
`}
</CodeCard>

### Time Series

<CodeCard file="timeseries.go">
{`type DataPoint struct {
    Timestamp time.Time
    Value     float64
}

timestamps := []time.Time{
    time.Now(),
    time.Now().Add(1 * time.Hour),
    time.Now().Add(2 * time.Hour),
}

values := []float64{23.5, 24.1, 23.8}

// Create time series
timeSeries := A.ZipWith(
    timestamps,
    values,
    func(t time.Time, v float64) DataPoint {
        return DataPoint{Timestamp: t, Value: v}
    },
)
`}
</CodeCard>

</Section>

<Section id="patterns" number="03" title="Common" titleAccent="Patterns">

### Zip with Index

<CodeCard file="index.go">
{`// Add index to elements
values := []string{"apple", "banana", "cherry"}

indexed := F.Pipe2(
    values,
    A.MapWithIndex(T.MakeTuple2[int, string]),
)
// []Tuple2[int, string]{
//   {Head: 0, Tail: "apple"},
//   {Head: 1, Tail: "banana"},
//   {Head: 2, Tail: "cherry"},
// }
`}
</CodeCard>

### Three Arrays

<CodeCard file="three_arrays.go">
{`// Zip three arrays using nested ZipWith
first := []string{"A", "B", "C"}
second := []int{1, 2, 3}
third := []bool{true, false, true}

type Triple struct {
    S string
    I int
    B bool
}

// First zip two arrays
step1 := A.ZipWith(
    first,
    second,
    T.MakeTuple2[string, int],
)

// Then zip with third
result := A.ZipWith(
    step1,
    third,
    func(t T.Tuple2[string, int], b bool) Triple {
        return Triple{S: t.Head, I: t.Tail, B: b}
    },
)
`}
</CodeCard>

### Parallel Transformation

<CodeCard file="parallel_transform.go">
{`// Transform two related arrays in parallel
transformed := A.ZipWith(
    sources,
    targets,
    func(src Source, tgt Target) Result {
        return transform(src, tgt)
    },
)
`}
</CodeCard>

</Section>

<Callout type="info">

**Performance**: `ZipWith` is more efficient than `Zip` + `Map` as it avoids creating intermediate tuples. The shorter array determines the result length.

</Callout>
