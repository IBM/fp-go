---
title: Eq (Equality)
hide_title: true
description: Type-safe equality checking with custom equality semantics.
sidebar_position: 25
---

import { PageHeader, Section, CodeCard, ApiTable, Callout, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="Reference · Utilities"
  title="Eq (Equality)"
  lede="Type-safe equality checking. Define custom equality semantics for any type using the Eq type class."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/eq' },
    { label: 'Type Class', value: 'Eq[A]' }
  ]}
/>

---

<Section num="1" title="Core API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `FromEquals` | `func FromEquals[A any](func(A, A) bool) Eq[A]` | Create Eq from function |
| `Contramap` | `func Contramap[A, B any](func(B) A) func(Eq[A]) Eq[B]` | Derive Eq by mapping |
</ApiTable>

</Section>

---

<Section num="2" title="Usage Examples">

### Basic Usage

<CodeCard file="basic.go">
{`import E "github.com/IBM/fp-go/v2/eq"

// Create equality from function
type User struct {
    ID   int
    Name string
}

userEq := E.FromEquals(func(a, b User) bool {
    return a.ID == b.ID
})

user1 := User{ID: 1, Name: "Alice"}
user2 := User{ID: 1, Name: "Alice Updated"}

userEq.Equals(user1, user2)  // true - same ID
`}
</CodeCard>

### Built-in Eq Instances

<CodeCard file="builtin.go">
{`import (
    N "github.com/IBM/fp-go/v2/number"
    S "github.com/IBM/fp-go/v2/string"
)

// Number equality
N.Eq.Equals(1, 1)  // true
N.Eq.Equals(1, 2)  // false

// String equality
S.Eq.Equals("hello", "hello")  // true
S.Eq.Equals("hello", "world")  // false
`}
</CodeCard>

### Contramap - Derive Equality

<CodeCard file="contramap.go">
{`type User struct {
    ID   int
    Name string
}

// Equality based on ID
userEq := E.Contramap(
    func(u User) int { return u.ID },
)(N.Eq)

userEq.Equals(
    User{ID: 1, Name: "Alice"},
    User{ID: 1, Name: "Alice Updated"},
)  // true - same ID

userEq.Equals(
    User{ID: 1, Name: "Alice"},
    User{ID: 2, Name: "Bob"},
)  // false - different IDs
`}
</CodeCard>

### Case-Insensitive Equality

<CodeCard file="case_insensitive.go">
{`// Case-insensitive string equality
caseInsensitiveEq := E.FromEquals(func(a, b string) bool {
    return strings.ToLower(a) == strings.ToLower(b)
})

caseInsensitiveEq.Equals("Hello", "HELLO")  // true
caseInsensitiveEq.Equals("Hello", "World")  // false
`}
</CodeCard>

### Struct Field Equality

<CodeCard file="field.go">
{`type Product struct {
    SKU   string
    Name  string
    Price float64
}

// Compare by SKU only
productEq := E.Contramap(
    func(p Product) string { return p.SKU },
)(S.Eq)

p1 := Product{SKU: "A123", Name: "Laptop", Price: 999}
p2 := Product{SKU: "A123", Name: "Laptop Pro", Price: 1299}

productEq.Equals(p1, p2)  // true - same SKU
`}
</CodeCard>

### Approximate Float Equality

<CodeCard file="float.go">
{`// Approximate equality for floats
const epsilon = 0.0001

floatEq := E.FromEquals(func(a, b float64) bool {
    return math.Abs(a-b) < epsilon
})

floatEq.Equals(3.14159, 3.14160)  // true - within epsilon
floatEq.Equals(3.14159, 3.15000)  // false - outside epsilon
`}
</CodeCard>

### Array Equality

<CodeCard file="array.go">
{`import A "github.com/IBM/fp-go/v2/array"

// Equality for arrays
arrayEq := A.Eq(N.Eq)

arr1 := []int{1, 2, 3}
arr2 := []int{1, 2, 3}
arr3 := []int{1, 2, 4}

arrayEq.Equals(arr1, arr2)  // true
arrayEq.Equals(arr1, arr3)  // false
`}
</CodeCard>

</Section>

---

<Section num="3" title="Common Patterns">

### Testing Helper

<CodeCard file="testing.go">
{`func AssertEqual[A any](
    t *testing.T,
    eq E.Eq[A],
    expected, actual A,
) {
    if !eq.Equals(expected, actual) {
        t.Errorf("Not equal:\nExpected: %v\nActual: %v",
            expected, actual)
    }
}

// Usage in tests
func TestSomething(t *testing.T) {
    expected := User{ID: 1, Name: "Alice"}
    actual := fetchUser(1)
    
    AssertEqual(t, userEq, expected, actual)
}
`}
</CodeCard>

### Semantic Equality

<CodeCard file="semantic.go">
{`type Status string

const (
    Active   Status = "active"
    Inactive Status = "inactive"
    Enabled  Status = "enabled"
    Disabled Status = "disabled"
)

// Treat active/enabled and inactive/disabled as equal
statusEq := E.FromEquals(func(a, b Status) bool {
    normalize := func(s Status) string {
        if s == Active || s == Enabled {
            return "active"
        }
        return "inactive"
    }
    return normalize(a) == normalize(b)
})

statusEq.Equals(Active, Enabled)    // true
statusEq.Equals(Inactive, Disabled) // true
statusEq.Equals(Active, Inactive)   // false
`}
</CodeCard>

### Composite Equality

<CodeCard file="composite.go">
{`type Address struct {
    Street string
    City   string
    Zip    string
}

// Equality based on city and zip only
addressEq := E.FromEquals(func(a, b Address) bool {
    return a.City == b.City && a.Zip == b.Zip
})

addr1 := Address{Street: "123 Main St", City: "NYC", Zip: "10001"}
addr2 := Address{Street: "456 Oak Ave", City: "NYC", Zip: "10001"}

addressEq.Equals(addr1, addr2)  // true - same city and zip
`}
</CodeCard>

</Section>

---

<Callout type="info">

**Contramap**: Use `Contramap` to derive equality from existing Eq instances. It's more composable than writing custom equality functions.

</Callout>

<Callout type="info">

**Use Cases**: Custom Eq instances are useful for:
- Domain-specific equality (e.g., case-insensitive strings)
- Approximate numeric equality
- Comparing by specific fields
- Testing and assertions

</Callout>

---

<Pager
  prev={{ to: '/docs/v2/utilities/compose', title: 'Compose' }}
  next={{ to: '/docs/v2/utilities/function', title: 'Function Utilities' }}
/>
