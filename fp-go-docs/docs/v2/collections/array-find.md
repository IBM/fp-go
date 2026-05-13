---
title: Array - Find
hide_title: true
description: Finding elements in arrays with type-safe operations - FindFirst, FindLast, FindFirstMap.
sidebar_position: 2
---

<PageHeader
  eyebrow="Reference · Collections"
  title="Array"
  titleAccent="Find"
  lede="Finding elements in arrays with type-safe operations. All find operations return Option[T] to safely handle cases where no element is found."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/array' },
    { label: 'Returns', value: 'Option[T]' }
  ]}
/>

<Section id="api" number="01" title="Core" titleAccent="API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `FindFirst` | `func FindFirst[A any](pred func(A) bool) func([]A) Option[A]` | Find first matching element |
| `FindFirstWithIndex` | `func FindFirstWithIndex[A any](pred func(int, A) bool) func([]A) Option[A]` | Find first with index |
| `FindLast` | `func FindLast[A any](pred func(A) bool) func([]A) Option[A]` | Find last matching element |
| `FindLastWithIndex` | `func FindLastWithIndex[A any](pred func(int, A) bool) func([]A) Option[A]` | Find last with index |
| `FindFirstMap` | `func FindFirstMap[A, B any](f func(A) Option[B]) func([]A) Option[B]` | Find and transform |
| `FindLastMap` | `func FindLastMap[A, B any](f func(A) Option[B]) func([]A) Option[B]` | Find last and transform |
</ApiTable>

</Section>

<Section id="examples" number="02" title="Usage" titleAccent="Examples">

### FindFirst

<CodeCard file="find_first.go">
{`import (
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
    O "github.com/IBM/fp-go/v2/option"
)

numbers := []int{1, 2, 3, 4, 5, 6}

// Find first even number
firstEven := F.Pipe2(
    numbers,
    A.FindFirst(func(n int) bool { return n%2 == 0 }),
)
// Some(2)

// Find first number > 10
notFound := F.Pipe2(
    numbers,
    A.FindFirst(func(n int) bool { return n > 10 }),
)
// None
`}
</CodeCard>

### FindLast

<CodeCard file="find_last.go">
{`numbers := []int{1, 2, 3, 4, 5, 6}

// Find last even number
lastEven := F.Pipe2(
    numbers,
    A.FindLast(func(n int) bool { return n%2 == 0 }),
)
// Some(6)

// Find last number < 3
lastSmall := F.Pipe2(
    numbers,
    A.FindLast(func(n int) bool { return n < 3 }),
)
// Some(2)
`}
</CodeCard>

### FindFirstMap

<CodeCard file="find_first_map.go">
{`type User struct {
    ID    int
    Name  string
    Email string
}

users := []User{
    {ID: 1, Name: "Alice", Email: ""},
    {ID: 2, Name: "Bob", Email: "bob@example.com"},
    {ID: 3, Name: "Charlie", Email: "charlie@example.com"},
}

// Find first user with email
firstWithEmail := F.Pipe2(
    users,
    A.FindFirstMap(func(u User) O.Option[string] {
        if u.Email != "" {
            return O.Some(u.Email)
        }
        return O.None[string]()
    }),
)
// Some("bob@example.com")
`}
</CodeCard>

### Parsing Example

<CodeCard file="parsing.go">
{`strings := []string{"abc", "123", "def", "456"}

// Find first valid number
firstNumber := F.Pipe2(
    strings,
    A.FindFirstMap(func(s string) O.Option[int] {
        if n, err := strconv.Atoi(s); err == nil {
            return O.Some(n)
        }
        return O.None[int]()
    }),
)
// Some(123)
`}
</CodeCard>

### Real-World: Configuration

<CodeCard file="config.go">
{`type Config struct {
    Key   string
    Value string
    Env   string
}

configs := []Config{
    {Key: "api_url", Value: "http://dev.api", Env: "dev"},
    {Key: "api_url", Value: "http://prod.api", Env: "prod"},
    {Key: "timeout", Value: "30", Env: "prod"},
}

// Find production API URL
prodAPI := F.Pipe3(
    configs,
    A.FindFirst(func(c Config) bool {
        return c.Key == "api_url" && c.Env == "prod"
    }),
    O.Map(func(c Config) string { return c.Value }),
    O.GetOrElse(func() string { return "http://default.api" }),
)
// "http://prod.api"
`}
</CodeCard>

### Real-World: Products

<CodeCard file="products.go">
{`type Product struct {
    ID       int
    Name     string
    Price    float64
    InStock  bool
    Category string
}

products := []Product{
    {ID: 1, Name: "Laptop", Price: 999, InStock: false, Category: "Electronics"},
    {ID: 2, Name: "Mouse", Price: 29, InStock: true, Category: "Electronics"},
    {ID: 3, Name: "Keyboard", Price: 79, InStock: true, Category: "Electronics"},
}

// Find first affordable in-stock electronics
affordable := F.Pipe2(
    products,
    A.FindFirst(func(p Product) bool {
        return p.InStock && 
               p.Category == "Electronics" && 
               p.Price < 100
    }),
)
// Some(Product{ID: 2, Name: "Mouse", ...})
`}
</CodeCard>

</Section>

<Section id="patterns" number="03" title="Common" titleAccent="Patterns">

### Find with Default

<CodeCard file="default.go">
{`// Find with fallback
result := F.Pipe2(
    numbers,
    A.FindFirst(func(n int) bool { return n > 100 }),
    O.GetOrElse(func() int { return -1 }),
)
// -1
`}
</CodeCard>

### Find and Transform

<CodeCard file="transform.go">
{`// Find, transform, and provide default
email := F.Pipe3(
    users,
    A.FindFirst(func(u User) bool { return u.ID == targetID }),
    O.Map(func(u User) string { return u.Email }),
    O.GetOrElse(func() string { return "unknown@example.com" }),
)
`}
</CodeCard>

### Validation and Find

<CodeCard file="validation.go">
{`// Find first valid item
valid := F.Pipe2(
    items,
    A.FindFirstMap(func(item Item) O.Option[Item] {
        if err := item.Validate(); err == nil {
            return O.Some(item)
        }
        return O.None[Item]()
    }),
)
`}
</CodeCard>

</Section>

<Callout type="info">

**Performance**: Find operations stop at the first match, making them more efficient than Filter + Head for finding single elements.

</Callout>
