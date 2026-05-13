---
sidebar_position: 7
title: Filtering and Mapping
description: Working with collections
hide_title: true
---

<PageHeader
  eyebrow="Recipes · 07 / 17"
  title="Filtering and"
  titleAccent="Mapping"
  lede="Filter and map collections effectively using fp-go's array operations for clean, composable data processing."
  meta={[
    { label: 'Difficulty', value: 'Beginner' },
    { label: 'Patterns', value: '6' },
    { label: 'Use Cases', value: 'Collections, Data Processing' }
  ]}
/>

<TLDR>
  <TLDRCard title="Filter Before Map" icon="filter">
    Reduce data early to improve performance—filter out unwanted elements before transforming.
  </TLDRCard>
  <TLDRCard title="Use FilterMap" icon="zap">
    When filtering and mapping, use FilterMap for efficiency—combines both operations in one pass.
  </TLDRCard>
  <TLDRCard title="Keep Functions Pure" icon="shield">
    No side effects in filter or map functions—makes code predictable and composable.
  </TLDRCard>
</TLDR>

<Section id="basic-filtering" number="01" title="Basic" titleAccent="Filtering">

Filter elements based on predicates to select only the items you need.

<CodeCard file="basic-filtering.go">
{`package main

import (
    "fmt"
    
    A "github.com/IBM/fp-go/v2/array"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // Filter even numbers
    evens := A.Filter(func(n int) bool {
        return n%2 == 0
    })(numbers)
    
    fmt.Printf("Even numbers: %v\\n", evens)
    // Output: Even numbers: [2 4 6 8 10]
    
    // Filter numbers greater than 5
    greaterThan5 := A.Filter(func(n int) bool {
        return n > 5
    })(numbers)
    
    fmt.Printf("Greater than 5: %v\\n", greaterThan5)
    // Output: Greater than 5: [6 7 8 9 10]
}`}
</CodeCard>

</Section>

<Section id="basic-mapping" number="02" title="Basic" titleAccent="Mapping">

Transform elements in a collection by applying a function to each item.

<CodeCard file="basic-mapping.go">
{`package main

import (
    "fmt"
    "strings"
    
    A "github.com/IBM/fp-go/v2/array"
)

func main() {
    words := []string{"hello", "world", "functional", "programming"}
    
    // Convert to uppercase
    uppercase := A.Map(strings.ToUpper)(words)
    fmt.Printf("Uppercase: %v\\n", uppercase)
    // Output: Uppercase: [HELLO WORLD FUNCTIONAL PROGRAMMING]
    
    // Get lengths
    lengths := A.Map(func(s string) int {
        return len(s)
    })(words)
    fmt.Printf("Lengths: %v\\n", lengths)
    // Output: Lengths: [5 5 10 11]
}`}
</CodeCard>

</Section>

<Section id="filter-map-combined" number="03" title="Filter and Map" titleAccent="Combined">

Chain filtering and mapping operations for powerful data transformations.

<CodeCard file="filter-map-combined.go">
{`package main

import (
    "fmt"
    "strings"
    
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
)

type Product struct {
    Name  string
    Price float64
    Stock int
}

func main() {
    products := []Product{
        {Name: "Apple", Price: 1.50, Stock: 10},
        {Name: "Banana", Price: 0.75, Stock: 0},
        {Name: "Cherry", Price: 2.00, Stock: 5},
        {Name: "Date", Price: 3.00, Stock: 0},
    }
    
    // Get names of in-stock products under $2
    result := F.Pipe2(
        products,
        A.Filter(func(p Product) bool {
            return p.Stock > 0 && p.Price < 2.00
        }),
        A.Map(func(p Product) string {
            return strings.ToUpper(p.Name)
        }),
    )
    
    fmt.Printf("Available affordable products: %v\\n", result)
    // Output: Available affordable products: [APPLE]
}`}
</CodeCard>

</Section>

<Section id="filtermap" number="04" title="FilterMap - Filter and" titleAccent="Transform">

Filter and transform in a single operation for better performance.

<CodeCard file="filtermap.go">
{`package main

import (
    "fmt"
    "strconv"
    
    A "github.com/IBM/fp-go/v2/array"
    O "github.com/IBM/fp-go/v2/option"
)

// Parse string to int, returning None for invalid strings
func parseInt(s string) O.Option[int] {
    val, err := strconv.Atoi(s)
    if err != nil {
        return O.None[int]()
    }
    return O.Some(val)
}

func main() {
    inputs := []string{"1", "abc", "2", "def", "3", "4"}
    
    // Parse and filter valid integers
    numbers := A.FilterMap(parseInt)(inputs)
    
    fmt.Printf("Valid numbers: %v\\n", numbers)
    // Output: Valid numbers: [1 2 3 4]
}`}
</CodeCard>

</Section>

<Section id="partition" number="05" title="Partition - Split by" titleAccent="Predicate">

Split a collection into two based on a predicate—one for matches, one for non-matches.

<CodeCard file="partition.go">
{`package main

import (
    "fmt"
    
    A "github.com/IBM/fp-go/v2/array"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // Partition into evens and odds
    evens, odds := A.Partition(func(n int) bool {
        return n%2 == 0
    })(numbers)
    
    fmt.Printf("Evens: %v\\n", evens)
    fmt.Printf("Odds: %v\\n", odds)
    // Output:
    // Evens: [2 4 6 8 10]
    // Odds: [1 3 5 7 9]
}`}
</CodeCard>

</Section>

<Section id="complex-filtering" number="06" title="Complex" titleAccent="Filtering">

Use multiple conditions and custom logic for sophisticated filtering.

<CodeCard file="complex-filtering.go">
{`package main

import (
    "fmt"
    "strings"
    
    A "github.com/IBM/fp-go/v2/array"
)

type User struct {
    Name   string
    Age    int
    Active bool
    Role   string
}

func main() {
    users := []User{
        {Name: "Alice", Age: 25, Active: true, Role: "admin"},
        {Name: "Bob", Age: 17, Active: true, Role: "user"},
        {Name: "Charlie", Age: 30, Active: false, Role: "user"},
        {Name: "Diana", Age: 22, Active: true, Role: "user"},
    }
    
    // Find active adult users
    activeAdults := A.Filter(func(u User) bool {
        return u.Active && u.Age >= 18
    })(users)
    
    fmt.Println("Active adult users:")
    for _, u := range activeAdults {
        fmt.Printf("  - %s (%d)\\n", u.Name, u.Age)
    }
    
    // Find admins or users over 25
    privileged := A.Filter(func(u User) bool {
        return u.Role == "admin" || u.Age > 25
    })(users)
    
    fmt.Println("\\nPrivileged users:")
    for _, u := range privileged {
        fmt.Printf("  - %s (%s)\\n", u.Name, u.Role)
    }
}`}
</CodeCard>

</Section>

<Section id="best-practices" number="07" title="Best" titleAccent="Practices">

<Checklist>
  <ChecklistItem status="required">
    **Use Filter before Map** — Reduce data early to improve performance
  </ChecklistItem>
  <ChecklistItem status="required">
    **Use FilterMap** — When filtering and mapping, use FilterMap for efficiency
  </ChecklistItem>
  <ChecklistItem status="required">
    **Keep predicates pure** — No side effects in filter functions
  </ChecklistItem>
  <ChecklistItem status="required">
    **Keep mappers pure** — No side effects in map functions
  </ChecklistItem>
  <ChecklistItem status="recommended">
    **Use Partition** — When you need both filtered and rejected items
  </ChecklistItem>
  <ChecklistItem status="optional">
    **Compose operations** — Chain multiple filters and maps for clarity
  </ChecklistItem>
</Checklist>

</Section>
