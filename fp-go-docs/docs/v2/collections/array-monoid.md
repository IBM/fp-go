---
title: Array - Monoid
hide_title: true
description: Combining arrays using monoid operations - Fold, FoldMap, ConcatAll.
sidebar_position: 6
---

<PageHeader
  eyebrow="Reference · Collections"
  title="Array"
  titleAccent="Monoid"
  lede="Combining arrays using monoid operations. Arrays form a monoid under concatenation, enabling powerful composition patterns."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/array' },
    { label: 'Type Class', value: 'Monoid' }
  ]}
/>

<Section id="api" number="01" title="Core" titleAccent="API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Fold` | `func Fold[A any](m Monoid[A]) func([]A) A` | Fold array using monoid |
| `FoldMap` | `func FoldMap[A, B any](m Monoid[B]) func(func(A) B) func([]A) B` | Map then fold |
| `ConcatAll` | `func ConcatAll[A any](m Monoid[[]A]) func([][]A) []A` | Concatenate all elements |
| `Monoid` | `func Monoid[A any]() Monoid[[]A]` | Array monoid instance |
</ApiTable>

</Section>

<Section id="examples" number="02" title="Usage" titleAccent="Examples">

### Fold

<CodeCard file="fold.go">
{`import (
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
    M "github.com/IBM/fp-go/v2/monoid"
    N "github.com/IBM/fp-go/v2/number"
    S "github.com/IBM/fp-go/v2/string"
)

// Sum numbers
numbers := []int{1, 2, 3, 4, 5}
sum := F.Pipe2(
    numbers,
    A.Fold(N.MonoidSum),
)
// 15

// Concatenate strings
words := []string{"Hello", " ", "World"}
result := F.Pipe2(
    words,
    A.Fold(S.Monoid),
)
// "Hello World"
`}
</CodeCard>

### FoldMap

<CodeCard file="foldmap.go">
{`type Product struct {
    Name  string
    Price float64
}

products := []Product{
    {Name: "Laptop", Price: 999},
    {Name: "Mouse", Price: 29},
    {Name: "Keyboard", Price: 79},
}

// Calculate total price
total := F.Pipe2(
    products,
    A.FoldMap(M.MonoidSum[float64]())(func(p Product) float64 {
        return p.Price
    }),
)
// 1107.0
`}
</CodeCard>

### ConcatAll

<CodeCard file="concat.go">
{`// Flatten nested arrays
nested := [][]int{{1, 2}, {3, 4}, {5}}
flat := A.ConcatAll(A.Monoid[int]())(nested)
// []int{1, 2, 3, 4, 5}
`}
</CodeCard>

### Aggregating Data

<CodeCard file="aggregate.go">
{`type Sale struct {
    Amount   float64
    Quantity int
}

sales := []Sale{
    {Amount: 100, Quantity: 2},
    {Amount: 200, Quantity: 1},
    {Amount: 50, Quantity: 5},
}

totalAmount := F.Pipe2(
    sales,
    A.FoldMap(M.MonoidSum[float64]())(func(s Sale) float64 {
        return s.Amount
    }),
)
// 350.0

totalQuantity := F.Pipe2(
    sales,
    A.FoldMap(M.MonoidSum[int]())(func(s Sale) int {
        return s.Quantity
    }),
)
// 8
`}
</CodeCard>

### Custom Monoid

<CodeCard file="custom.go">
{`type ValidationResult struct {
    Errors []string
    Valid  bool
}

// Monoid for validation results
validationMonoid := M.MakeMonoid(
    func(a, b ValidationResult) ValidationResult {
        return ValidationResult{
            Errors: append(a.Errors, b.Errors...),
            Valid:  a.Valid && b.Valid,
        }
    },
    ValidationResult{Errors: []string{}, Valid: true},
)

results := []ValidationResult{
    {Errors: []string{}, Valid: true},
    {Errors: []string{"Invalid email"}, Valid: false},
    {Errors: []string{"Password too short"}, Valid: false},
}

combined := F.Pipe2(
    results,
    A.Fold(validationMonoid),
)
// ValidationResult{
//   Errors: []string{"Invalid email", "Password too short"},
//   Valid: false,
// }
`}
</CodeCard>

</Section>
