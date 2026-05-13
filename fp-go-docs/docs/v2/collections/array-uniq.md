---
title: Array - Uniq
hide_title: true
description: Removing duplicates from arrays with StrictUniq and Uniq operations.
sidebar_position: 5
---

<PageHeader
  eyebrow="Reference · Collections"
  title="Array"
  titleAccent="Uniqueness"
  lede="Removing duplicates from arrays. Use StrictUniq for comparable types or Uniq with custom key extraction."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/array' },
    { label: 'Operations', value: 'StrictUniq, Uniq' }
  ]}
/>

<Section id="api" number="01" title="Core" titleAccent="API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `StrictUniq` | `func StrictUniq[A comparable]([]A) []A` | Remove duplicates (comparable types) |
| `Uniq` | `func Uniq[A, K comparable](f func(A) K) func([]A) []A` | Remove duplicates by key |
</ApiTable>

</Section>

<Section id="examples" number="02" title="Usage" titleAccent="Examples">

### StrictUniq

<CodeCard file="strict_uniq.go">
{`import (
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
)

// Numbers
numbers := []int{1, 2, 2, 3, 1, 4, 3, 5}
unique := A.StrictUniq(numbers)
// []int{1, 2, 3, 4, 5}

// Strings
words := []string{"apple", "banana", "apple", "cherry", "banana"}
unique := A.StrictUniq(words)
// []string{"apple", "banana", "cherry"}

// Order preserved (first occurrence kept)
items := []int{3, 1, 2, 1, 3, 4, 2}
unique := A.StrictUniq(items)
// []int{3, 1, 2, 4}
`}
</CodeCard>

### Uniq by Field

<CodeCard file="uniq_by_field.go">
{`type User struct {
    ID   int
    Name string
    Age  int
}

users := []User{
    {ID: 1, Name: "Alice", Age: 30},
    {ID: 2, Name: "Bob", Age: 25},
    {ID: 1, Name: "Alice Updated", Age: 31},  // Duplicate ID
    {ID: 3, Name: "Charlie", Age: 35},
}

// Unique by ID
uniqueByID := F.Pipe2(
    users,
    A.Uniq(func(u User) int { return u.ID }),
)
// []User{
//   {ID: 1, Name: "Alice", Age: 30},      // First occurrence kept
//   {ID: 2, Name: "Bob", Age: 25},
//   {ID: 3, Name: "Charlie", Age: 35},
// }

// Unique by name (case-insensitive)
uniqueByName := F.Pipe2(
    users,
    A.Uniq(func(u User) string { 
        return strings.ToLower(u.Name) 
    }),
)
`}
</CodeCard>

### Uniq by Category

<CodeCard file="uniq_category.go">
{`type Product struct {
    Name     string
    Category string
    Price    float64
}

products := []Product{
    {Name: "Laptop", Category: "Electronics", Price: 999},
    {Name: "Mouse", Category: "Electronics", Price: 29},
    {Name: "Desk", Category: "Furniture", Price: 299},
    {Name: "Keyboard", Category: "Electronics", Price: 79},
}

// Unique by category
uniqueCategories := F.Pipe2(
    products,
    A.Uniq(func(p Product) string { return p.Category }),
)
// []Product{
//   {Name: "Laptop", Category: "Electronics", ...},
//   {Name: "Desk", Category: "Furniture", ...},
// }
`}
</CodeCard>

### Deduplicating Tags

<CodeCard file="tags.go">
{`type Article struct {
    Title string
    Tags  []string
}

articles := []Article{
    {Title: "Go Basics", Tags: []string{"go", "tutorial", "basics"}},
    {Title: "Advanced Go", Tags: []string{"go", "advanced", "patterns"}},
    {Title: "FP in Go", Tags: []string{"go", "functional", "fp"}},
}

// Get all unique tags
allTags := F.Pipe3(
    articles,
    A.Chain(func(a Article) []string { return a.Tags }),
    A.StrictUniq,
)
// []string{"go", "tutorial", "basics", "advanced", "patterns", "functional", "fp"}
`}
</CodeCard>

### Unique Combinations

<CodeCard file="combinations.go">
{`type Event struct {
    UserID int
    Action string
    Target string
}

events := []Event{
    {UserID: 1, Action: "view", Target: "page1"},
    {UserID: 1, Action: "view", Target: "page2"},
    {UserID: 1, Action: "view", Target: "page1"},  // Duplicate
    {UserID: 2, Action: "view", Target: "page1"},
}

// Unique by user-target combination
uniqueViews := F.Pipe2(
    events,
    A.Uniq(func(e Event) string {
        return fmt.Sprintf("%d-%s", e.UserID, e.Target)
    }),
)
// 3 unique user-target combinations
`}
</CodeCard>

</Section>

<Section id="patterns" number="03" title="Common" titleAccent="Patterns">

### Unique and Sort

<CodeCard file="unique_sort.go">
{`numbers := []int{5, 2, 8, 2, 1, 5, 3}

// Get unique values and sort
result := F.Pipe3(
    numbers,
    A.StrictUniq,
    A.Sort(N.Ord),
)
// []int{1, 2, 3, 5, 8}
`}
</CodeCard>

### Unique by Multiple Fields

<CodeCard file="multiple_fields.go">
{`type Record struct {
    Year  int
    Month int
    Value float64
}

records := []Record{
    {Year: 2024, Month: 1, Value: 100},
    {Year: 2024, Month: 2, Value: 200},
    {Year: 2024, Month: 1, Value: 150},  // Duplicate year-month
}

// Unique by year-month combination
uniqueByPeriod := F.Pipe2(
    records,
    A.Uniq(func(r Record) string {
        return fmt.Sprintf("%d-%02d", r.Year, r.Month)
    }),
)
// 2 records (2024-01 and 2024-02)
`}
</CodeCard>

### Union via Unique

<CodeCard file="union.go">
{`// Combine and deduplicate
func union[T comparable](arrays ...[]T) []T {
    combined := F.Pipe2(
        arrays,
        A.Flatten,
    )
    return A.StrictUniq(combined)
}

arr1 := []int{1, 2, 3}
arr2 := []int{3, 4, 5}
result := union(arr1, arr2)
// []int{1, 2, 3, 4, 5}
`}
</CodeCard>

</Section>

<Callout type="info">

**Performance**: For comparable types, use `StrictUniq` which is faster. For custom key extraction, use `Uniq`. First occurrence is always kept.

</Callout>
