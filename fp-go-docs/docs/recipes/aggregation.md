---
sidebar_position: 8
title: Data Aggregation
description: Reducing and summarizing data
hide_title: true
---

<PageHeader
  eyebrow="Recipes · 08 / 17"
  title="Data"
  titleAccent="Aggregation"
  lede="Aggregate and summarize data using functional patterns with monoids, fold operations, and traverse for powerful data reduction."
  meta={[
    { label: 'Difficulty', value: 'Intermediate' },
    { label: 'Patterns', value: '6' },
    { label: 'Use Cases', value: 'Statistics, Summaries, Analytics' }
  ]}
/>

<TLDR>
  <TLDRCard title="Use Built-in Monoids" icon="package">
    Leverage fp-go's monoids for common operations—sum, product, string concatenation all have optimized implementations.
  </TLDRCard>
  <TLDRCard title="Leverage FoldMap" icon="zap">
    Transform and aggregate in one pass—more efficient than separate map and reduce operations.
  </TLDRCard>
  <TLDRCard title="Handle Empty Collections" icon="alert-circle">
    Use Option for potentially empty results—prevents panics and makes edge cases explicit.
  </TLDRCard>
</TLDR>

<Section id="basic-fold" number="01" title="Basic Aggregation with" titleAccent="Fold">

Use fold operations to reduce collections to single values.

<CodeCard file="basic-fold.go">
{`package main

import (
    "fmt"
    A "github.com/IBM/fp-go/v2/array"
    N "github.com/IBM/fp-go/v2/number"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // Sum using MonoidSum
    sum := A.Reduce(N.MonoidSum[int]())(numbers)
    fmt.Println("Sum:", sum) // Sum: 15
    
    // Product using MonoidProduct
    product := A.Reduce(N.MonoidProduct[int]())(numbers)
    fmt.Println("Product:", product) // Product: 120
}`}
</CodeCard>

<CodeCard file="string-concat.go">
{`package main

import (
    "fmt"
    A "github.com/IBM/fp-go/v2/array"
    S "github.com/IBM/fp-go/v2/string"
)

func main() {
    words := []string{"Hello", " ", "functional", " ", "world"}
    
    // Concatenate strings
    result := A.Reduce(S.Monoid)(words)
    fmt.Println(result) // Hello functional world
    
    // With separator
    items := []string{"apple", "banana", "cherry"}
    joined := A.Intercalate(S.Monoid)(", ")(items)
    fmt.Println(joined) // apple, banana, cherry
}`}
</CodeCard>

</Section>

<Section id="complex-data" number="02" title="Aggregating Complex" titleAccent="Data">

Group and aggregate complex data structures for summaries and analytics.

<CodeCard file="grouping-counting.go">
{`package main

import (
    "fmt"
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/record"
)

type Product struct {
    Name     string
    Category string
    Price    float64
}

func main() {
    products := []Product{
        {Name: "Laptop", Category: "Electronics", Price: 999.99},
        {Name: "Mouse", Category: "Electronics", Price: 29.99},
        {Name: "Desk", Category: "Furniture", Price: 299.99},
        {Name: "Chair", Category: "Furniture", Price: 199.99},
    }
    
    // Count by category
    countByCategory := F.Pipe2(
        products,
        A.GroupBy(func(p Product) string { return p.Category }),
        R.Map(func(items []Product) int { return len(items) }),
    )
    
    fmt.Println("Count by category:", countByCategory)
    // Count by category: map[Electronics:2 Furniture:2]
    
    // Sum prices by category
    sumByCategory := F.Pipe2(
        products,
        A.GroupBy(func(p Product) string { return p.Category }),
        R.Map(func(items []Product) float64 {
            return A.Reduce(func(acc, p float64) float64 {
                return acc + p
            })(0.0)(A.Map(func(p Product) float64 {
                return p.Price
            })(items))
        }),
    )
    
    fmt.Println("Sum by category:", sumByCategory)
    // Sum by category: map[Electronics:1029.98 Furniture:499.98]
}`}
</CodeCard>

</Section>

<Section id="statistics" number="03" title="Statistical" titleAccent="Aggregations">

Calculate statistics like min, max, average, and more from collections.

<CodeCard file="statistics.go">
{`package main

import (
    "fmt"
    "math"
    A "github.com/IBM/fp-go/v2/array"
    O "github.com/IBM/fp-go/v2/option"
)

type Stats struct {
    Count   int
    Sum     float64
    Min     float64
    Max     float64
    Average float64
}

func calculateStats(numbers []float64) O.Option[Stats] {
    if len(numbers) == 0 {
        return O.None[Stats]()
    }
    
    count := len(numbers)
    sum := A.Reduce(func(acc, n float64) float64 {
        return acc + n
    })(0.0)(numbers)
    
    min := A.Reduce(func(acc, n float64) float64 {
        return math.Min(acc, n)
    })(numbers[0])(numbers[1:])
    
    max := A.Reduce(func(acc, n float64) float64 {
        return math.Max(acc, n)
    })(numbers[0])(numbers[1:])
    
    return O.Some(Stats{
        Count:   count,
        Sum:     sum,
        Min:     min,
        Max:     max,
        Average: sum / float64(count),
    })
}

func main() {
    numbers := []float64{10.5, 20.3, 15.7, 30.2, 25.8}
    
    stats := calculateStats(numbers)
    
    O.Match(
        func() { fmt.Println("No data") },
        func(s Stats) {
            fmt.Printf("Count: %d\\n", s.Count)
            fmt.Printf("Sum: %.2f\\n", s.Sum)
            fmt.Printf("Min: %.2f\\n", s.Min)
            fmt.Printf("Max: %.2f\\n", s.Max)
            fmt.Printf("Average: %.2f\\n", s.Average)
        },
    )(stats)
}`}
</CodeCard>

</Section>

<Section id="effects" number="04" title="Aggregating with" titleAccent="Effects">

Aggregate operations that may fail or have side effects using traverse.

<CodeCard file="collecting-results.go">
{`package main

import (
    "fmt"
    "strconv"
    A "github.com/IBM/fp-go/v2/array"
    E "github.com/IBM/fp-go/v2/either"
    F "github.com/IBM/fp-go/v2/function"
)

func parseNumbers(strings []string) E.Either[error, []int] {
    return A.Traverse[string](E.Applicative[error, int]())(
        func(s string) E.Either[error, int] {
            n, err := strconv.Atoi(s)
            if err != nil {
                return E.Left[int](err)
            }
            return E.Right[error](n)
        },
    )(strings)
}

func main() {
    // Valid input
    validStrings := []string{"1", "2", "3", "4", "5"}
    result1 := parseNumbers(validStrings)
    
    E.Match(
        func(err error) { fmt.Println("Error:", err) },
        func(numbers []int) {
            sum := A.Reduce(func(acc, n int) int {
                return acc + n
            })(0)(numbers)
            fmt.Println("Sum:", sum) // Sum: 15
        },
    )(result1)
    
    // Invalid input
    invalidStrings := []string{"1", "2", "invalid", "4"}
    result2 := parseNumbers(invalidStrings)
    
    E.Match(
        func(err error) { fmt.Println("Error:", err) },
        func(numbers []int) { fmt.Println("Sum:", numbers) },
    )(result2)
}`}
</CodeCard>

</Section>

<Section id="custom-monoids" number="05" title="Custom" titleAccent="Monoids">

Build custom monoids for domain-specific aggregations.

<CodeCard file="custom-monoid.go">
{`package main

import (
    "fmt"
    A "github.com/IBM/fp-go/v2/array"
    M "github.com/IBM/fp-go/v2/monoid"
)

type MinMax struct {
    Min int
    Max int
}

func minMaxMonoid() M.Monoid[MinMax] {
    return M.MakeMonoid(
        // Empty value
        func() MinMax {
            return MinMax{Min: int(^uint(0) >> 1), Max: -int(^uint(0)>>1) - 1}
        },
        // Concat operation
        func(a, b MinMax) MinMax {
            min := a.Min
            if b.Min < min {
                min = b.Min
            }
            max := a.Max
            if b.Max > max {
                max = b.Max
            }
            return MinMax{Min: min, Max: max}
        },
    )
}

func toMinMax(n int) MinMax {
    return MinMax{Min: n, Max: n}
}

func main() {
    numbers := []int{5, 2, 8, 1, 9, 3}
    
    result := A.FoldMap(minMaxMonoid())(toMinMax)(numbers)
    
    fmt.Printf("Min: %d, Max: %d\\n", result.Min, result.Max)
    // Min: 1, Max: 9
}`}
</CodeCard>

</Section>

<Section id="parallel" number="06" title="Parallel" titleAccent="Aggregation">

Aggregate data in parallel for improved performance on large datasets.

<CodeCard file="parallel-sum.go">
{`package main

import (
    "fmt"
    "sync"
    A "github.com/IBM/fp-go/v2/array"
)

func parallelSum(numbers []int, workers int) int {
    if len(numbers) == 0 {
        return 0
    }
    
    chunkSize := (len(numbers) + workers - 1) / workers
    chunks := A.Chunksof(chunkSize)(numbers)
    
    results := make([]int, len(chunks))
    var wg sync.WaitGroup
    
    for i, chunk := range chunks {
        wg.Add(1)
        go func(idx int, data []int) {
            defer wg.Done()
            results[idx] = A.Reduce(func(acc, n int) int {
                return acc + n
            })(0)(data)
        }(i, chunk)
    }
    
    wg.Wait()
    
    return A.Reduce(func(acc, n int) int {
        return acc + n
    })(0)(results)
}

func main() {
    numbers := A.MakeBy(1000)(func(i int) int { return i + 1 })
    
    sum := parallelSum(numbers, 4)
    fmt.Println("Sum:", sum) // Sum: 500500
}`}
</CodeCard>

</Section>

<Section id="best-practices" number="07" title="Best" titleAccent="Practices">

<Checklist>
  <ChecklistItem status="required">
    **Choose the right monoid** — Use built-in monoids when available
  </ChecklistItem>
  <ChecklistItem status="required">
    **Leverage FoldMap** — Transform and aggregate in one pass for efficiency
  </ChecklistItem>
  <ChecklistItem status="required">
    **Handle empty collections** — Use Option for potentially empty results
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Build custom monoids** — Create reusable aggregation patterns for your domain
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Consider parallelization** — For large datasets, parallel aggregation can improve performance
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Use traverse for effects** — When aggregating operations that may fail
  </ChecklistItem>
</Checklist>

</Section>
