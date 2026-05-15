---
title: Array - Sort
hide_title: true
description: Sorting arrays with type-safe ordering using the Ord type class.
sidebar_position: 3
---

<PageHeader
  eyebrow="Reference · Collections"
  title="Array"
  titleAccent="Sort"
  lede="Sorting arrays with type-safe ordering. Stable sorting operations using the Ord type class. All operations return new arrays without modifying the original."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/array' },
    { label: 'Type Class', value: 'Ord' }
  ]}
/>

<Section id="api" number="01" title="Core" titleAccent="API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Sort` | `func Sort[A any](ord Ord[A]) func([]A) []A` | Sort using Ord instance |
| `SortByKey` | `func SortByKey[A, B any](ord Ord[B], f func(A) B) func([]A) []A` | Sort by extracting key |
| `SortBy` | `func SortBy[A any](ords []Ord[A]) func([]A) []A` | Sort using multiple orderings |
</ApiTable>

</Section>

<Section id="examples" number="02" title="Usage" titleAccent="Examples">

### Basic Sorting

<CodeCard file="basic.go">
{`import (
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
    N "github.com/IBM/fp-go/v2/number"
    O "github.com/IBM/fp-go/v2/ord"
)

numbers := []int{5, 2, 8, 1, 9, 3}

// Sort ascending
sorted := F.Pipe2(
    numbers,
    A.Sort(N.Ord),
)
// []int{1, 2, 3, 5, 8, 9}

// Sort descending
descending := F.Pipe2(
    numbers,
    A.Sort(O.Reverse(N.Ord)),
)
// []int{9, 8, 5, 3, 2, 1}
`}
</CodeCard>

### String Sorting

<CodeCard file="strings.go">
{`import S "github.com/IBM/fp-go/v2/string"

words := []string{"zebra", "apple", "mango", "banana"}

// Sort alphabetically
sorted := F.Pipe2(
    words,
    A.Sort(S.Ord),
)
// []string{"apple", "banana", "mango", "zebra"}

// Case-insensitive sort
caseInsensitive := O.FromCompare(func(a, b string) int {
    return strings.Compare(
        strings.ToLower(a),
        strings.ToLower(b),
    )
})

sorted := F.Pipe2(
    []string{"Zebra", "apple", "Mango", "banana"},
    A.Sort(caseInsensitive),
)
// []string{"apple", "banana", "Mango", "Zebra"}
`}
</CodeCard>

### SortByKey

<CodeCard file="sort_by_key.go">
{`type Person struct {
    Name string
    Age  int
}

people := []Person{
    {Name: "Alice", Age: 30},
    {Name: "Bob", Age: 25},
    {Name: "Charlie", Age: 35},
}

// Sort by age
byAge := F.Pipe2(
    people,
    A.SortByKey(N.Ord, func(p Person) int { return p.Age }),
)
// [{Bob 25} {Alice 30} {Charlie 35}]

// Sort by name
byName := F.Pipe2(
    people,
    A.SortByKey(S.Ord, func(p Person) string { return p.Name }),
)
// [{Alice 30} {Bob 25} {Charlie 35}]
`}
</CodeCard>

### Multiple Orderings

<CodeCard file="multiple.go">
{`type Employee struct {
    Department string
    Name       string
    Salary     int
}

employees := []Employee{
    {Department: "IT", Name: "Alice", Salary: 80000},
    {Department: "HR", Name: "Bob", Salary: 60000},
    {Department: "IT", Name: "Charlie", Salary: 75000},
    {Department: "HR", Name: "Diana", Salary: 65000},
}

// Sort by department, then by salary (descending)
sorted := F.Pipe2(
    employees,
    A.SortBy([]O.Ord[Employee]{
        // First by department
        O.FromCompare(func(a, b Employee) int {
            return strings.Compare(a.Department, b.Department)
        }),
        // Then by salary (descending)
        O.Reverse(O.FromCompare(func(a, b Employee) int {
            return a.Salary - b.Salary
        })),
    }),
)
// [{HR Diana 65000} {HR Bob 60000} {IT Alice 80000} {IT Charlie 75000}]
`}
</CodeCard>

### Custom Ordering

<CodeCard file="custom.go">
{`// Custom comparison function
customOrd := O.FromCompare(func(a, b MyType) int {
    // Return: -1 if a < b, 0 if a == b, 1 if a > b
    if a.Value < b.Value {
        return -1
    }
    if a.Value > b.Value {
        return 1
    }
    return 0
})

// Derive ordering from existing one
type User struct {
    ID   int
    Name string
}

// Order users by ID using number ordering
userOrd := O.Contramap(
    func(u User) int { return u.ID },
)(N.Ord)

sorted := F.Pipe2(
    users,
    A.Sort(userOrd),
)
`}
</CodeCard>

### Sorting Events by Time

<CodeCard file="events.go">
{`type Event struct {
    Name      string
    Timestamp time.Time
    Priority  int
}

events := []Event{
    {Name: "Meeting", Timestamp: time.Now().Add(2 * time.Hour), Priority: 2},
    {Name: "Deadline", Timestamp: time.Now().Add(1 * time.Hour), Priority: 1},
    {Name: "Review", Timestamp: time.Now().Add(3 * time.Hour), Priority: 2},
}

// Sort by priority (ascending), then by time (ascending)
timeOrd := O.FromCompare(func(a, b time.Time) int {
    if a.Before(b) { return -1 }
    if a.After(b) { return 1 }
    return 0
})

sorted := F.Pipe2(
    events,
    A.SortBy([]O.Ord[Event]{
        O.FromCompare(func(a, b Event) int {
            return a.Priority - b.Priority
        }),
        O.Contramap(func(e Event) time.Time { return e.Timestamp })(timeOrd),
    }),
)
// Events sorted by priority first, then by time
`}
</CodeCard>

### Sorting with Nullables

<CodeCard file="nullables.go">
{`type Record struct {
    ID        int
    UpdatedAt *time.Time  // nullable
}

records := []Record{
    {ID: 1, UpdatedAt: nil},
    {ID: 2, UpdatedAt: ptr(time.Now().Add(-1 * time.Hour))},
    {ID: 3, UpdatedAt: ptr(time.Now())},
}

// Sort with nulls last
nullsLastOrd := O.FromCompare(func(a, b Record) int {
    if a.UpdatedAt == nil && b.UpdatedAt == nil {
        return 0
    }
    if a.UpdatedAt == nil {
        return 1  // a goes after b
    }
    if b.UpdatedAt == nil {
        return -1  // a goes before b
    }
    if a.UpdatedAt.Before(*b.UpdatedAt) {
        return -1
    }
    if a.UpdatedAt.After(*b.UpdatedAt) {
        return 1
    }
    return 0
})

sorted := F.Pipe2(
    records,
    A.Sort(nullsLastOrd),
)
// [{2 ...} {3 ...} {1 nil}]
`}
</CodeCard>

</Section>

<Section id="patterns" number="03" title="Common" titleAccent="Patterns">

### Sort and Take Top N

<CodeCard file="top_n.go">
{`// Get top 5 by score
top5 := F.Pipe3(
    items,
    A.SortByKey(O.Reverse(N.Ord), func(i Item) int { return i.Score }),
    A.Slice(0, 5),
)
`}
</CodeCard>

### Conditional Sorting

<CodeCard file="conditional.go">
{`// Sort ascending or descending based on flag
ord := func(ascending bool) O.Ord[int] {
    if ascending {
        return N.Ord
    }
    return O.Reverse(N.Ord)
}

sorted := F.Pipe2(
    numbers,
    A.Sort(ord(isAscending)),
)
`}
</CodeCard>

### Precompute Expensive Keys

<CodeCard file="precompute.go">
{`type ItemWithKey struct {
    Item Item
    Key  float64
}

// Precompute expensive keys
withKeys := F.Pipe2(
    items,
    A.Map(func(item Item) ItemWithKey {
        return ItemWithKey{
            Item: item,
            Key:  expensiveComputation(item),
        }
    }),
)

// Sort by precomputed key
sorted := F.Pipe3(
    withKeys,
    A.SortByKey(floatOrd, func(iwk ItemWithKey) float64 { return iwk.Key }),
    A.Map(func(iwk ItemWithKey) Item { return iwk.Item }),
)
`}
</CodeCard>

</Section>

<Callout type="info">

**Stability**: All sort operations are stable - elements with equal keys maintain their relative order. This is useful for multi-level sorting.

</Callout>
