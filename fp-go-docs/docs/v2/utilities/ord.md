---
title: Ord (Ordering)
hide_title: true
description: Type-safe ordering and comparison for sorting and ordering operations.
sidebar_position: 26
---

import { PageHeader, Section, CodeCard, ApiTable, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="Reference · Utilities"
  title="Ord (Ordering)"
  lede="Type-safe ordering and comparison. Define custom ordering for any type using the Ord type class for sorting and comparison operations."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/ord' },
    { label: 'Type Class', value: 'Ord[A]' }
  ]}
/>

---

<Section num="1" title="Core API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `FromCompare` | `func FromCompare[A any](func(A, A) int) Ord[A]` | Create Ord from function |
| `Contramap` | `func Contramap[A, B any](func(B) A) func(Ord[A]) Ord[B]` | Derive Ord by mapping |
| `Reverse` | `func Reverse[A any](Ord[A]) Ord[A]` | Reverse ordering |
</ApiTable>

</Section>

---

<Section num="2" title="Usage Examples">

### Basic Usage

<CodeCard file="basic.go">
{`import (
    O "github.com/IBM/fp-go/v2/ord"
    N "github.com/IBM/fp-go/v2/number"
)

// Built-in orderings
N.Ord.Compare(1, 2)  // -1 (less than)
N.Ord.Compare(2, 2)  // 0 (equal)
N.Ord.Compare(3, 2)  // 1 (greater than)
`}
</CodeCard>

### FromCompare

<CodeCard file="from_compare.go">
{`type User struct {
    ID   int
    Name string
    Age  int
}

// Order by age
userOrd := O.FromCompare(func(a, b User) int {
    if a.Age < b.Age { return -1 }
    if a.Age > b.Age { return 1 }
    return 0
})

user1 := User{ID: 1, Name: "Alice", Age: 30}
user2 := User{ID: 2, Name: "Bob", Age: 25}

userOrd.Compare(user1, user2)  // 1 (Alice is older)
`}
</CodeCard>

### Contramap - Derive Ordering

<CodeCard file="contramap.go">
{`type User struct {
    ID   int
    Name string
}

// Order users by ID
userOrd := O.Contramap(
    func(u User) int { return u.ID },
)(N.Ord)

user1 := User{ID: 1, Name: "Alice"}
user2 := User{ID: 2, Name: "Bob"}

userOrd.Compare(user1, user2)  // -1 (1 < 2)
`}
</CodeCard>

### Reverse Ordering

<CodeCard file="reverse.go">
{`// Descending order
descending := O.Reverse(N.Ord)

descending.Compare(1, 2)  // 1 (reversed from -1)
descending.Compare(3, 2)  // -1 (reversed from 1)

// Use with sorting
import A "github.com/IBM/fp-go/v2/array"

numbers := []int{3, 1, 4, 1, 5, 9, 2, 6}
sorted := A.Sort(descending)(numbers)
// []int{9, 6, 5, 4, 3, 2, 1, 1}
`}
</CodeCard>

### String Ordering

<CodeCard file="string.go">
{`import S "github.com/IBM/fp-go/v2/string"

// Alphabetical order
S.Ord.Compare("apple", "banana")  // -1
S.Ord.Compare("banana", "apple")  // 1

// Case-insensitive ordering
caseInsensitiveOrd := O.FromCompare(func(a, b string) int {
    return strings.Compare(
        strings.ToLower(a),
        strings.ToLower(b),
    )
})

caseInsensitiveOrd.Compare("Apple", "banana")  // -1
`}
</CodeCard>

### Multi-Field Ordering

<CodeCard file="multifield.go">
{`type Product struct {
    Category string
    Name     string
    Price    float64
}

// Order by category, then by name
productOrd := O.FromCompare(func(a, b Product) int {
    // First compare category
    if a.Category < b.Category { return -1 }
    if a.Category > b.Category { return 1 }
    
    // If category equal, compare name
    if a.Name < b.Name { return -1 }
    if a.Name > b.Name { return 1 }
    
    return 0
})
`}
</CodeCard>

### Custom Priority Ordering

<CodeCard file="priority.go">
{`type Priority int

const (
    Low    Priority = 1
    Medium Priority = 2
    High   Priority = 3
)

// High priority first
priorityOrd := O.FromCompare(func(a, b Priority) int {
    // Reverse comparison for high-first
    if a > b { return -1 }
    if a < b { return 1 }
    return 0
})

// Or use Reverse
priorityOrd := O.Reverse(
    O.FromCompare(func(a, b Priority) int {
        if a < b { return -1 }
        if a > b { return 1 }
        return 0
    }),
)
`}
</CodeCard>

### Date/Time Ordering

<CodeCard file="datetime.go">
{`type Event struct {
    Name      string
    Timestamp time.Time
}

// Order by timestamp
eventOrd := O.Contramap(
    func(e Event) time.Time { return e.Timestamp },
)(O.FromCompare(func(a, b time.Time) int {
    if a.Before(b) { return -1 }
    if a.After(b) { return 1 }
    return 0
}))

// Or simpler with Unix timestamp
eventOrd := O.Contramap(
    func(e Event) int64 { return e.Timestamp.Unix() },
)(O.FromCompare(func(a, b int64) int {
    if a < b { return -1 }
    if a > b { return 1 }
    return 0
}))
`}
</CodeCard>

</Section>

---

<Section num="3" title="Common Patterns">

### Sorting with Ord

<CodeCard file="sorting.go">
{`import A "github.com/IBM/fp-go/v2/array"

users := []User{
    {ID: 3, Name: "Charlie", Age: 35},
    {ID: 1, Name: "Alice", Age: 30},
    {ID: 2, Name: "Bob", Age: 25},
}

// Sort by age
byAge := O.Contramap(func(u User) int { return u.Age })(N.Ord)
sorted := A.Sort(byAge)(users)
// Sorted by age: Bob(25), Alice(30), Charlie(35)

// Sort by name
byName := O.Contramap(func(u User) string { return u.Name })(S.Ord)
sorted := A.Sort(byName)(users)
// Sorted by name: Alice, Bob, Charlie
`}
</CodeCard>

### Min/Max Operations

<CodeCard file="minmax.go">
{`// Find minimum
min := func(ord O.Ord[A]) func(A, A) A {
    return func(a, b A) A {
        if ord.Compare(a, b) <= 0 {
            return a
        }
        return b
    }
}

// Find maximum
max := func(ord O.Ord[A]) func(A, A) A {
    return func(a, b A) A {
        if ord.Compare(a, b) >= 0 {
            return a
        }
        return b
    }
}

minValue := min(N.Ord)(5, 3)  // 3
maxValue := max(N.Ord)(5, 3)  // 5
`}
</CodeCard>

### Composite Ordering

<CodeCard file="composite.go">
{`// Order by multiple criteria
type Task struct {
    Priority int
    DueDate  time.Time
    Name     string
}

taskOrd := O.FromCompare(func(a, b Task) int {
    // First by priority (high first)
    if a.Priority > b.Priority { return -1 }
    if a.Priority < b.Priority { return 1 }
    
    // Then by due date (earliest first)
    if a.DueDate.Before(b.DueDate) { return -1 }
    if a.DueDate.After(b.DueDate) { return 1 }
    
    // Finally by name
    return strings.Compare(a.Name, b.Name)
})
`}
</CodeCard>

</Section>

---

<Callout type="info">

**Contramap**: Use `Contramap` to derive ordering from existing Ord instances. It's more composable and reusable than writing custom comparison functions.

</Callout>

<Callout type="info">

**Comparison Result**:
- Return `-1` if first argument is less than second
- Return `0` if arguments are equal
- Return `1` if first argument is greater than second

</Callout>


---

<Pager
  prev={{ to: '/docs/v2/utilities/magma', title: 'Magma' }}
  next={{ to: '/docs/v2/utilities/pipe-flow', title: 'Pipe & Flow' }}
/>

---
