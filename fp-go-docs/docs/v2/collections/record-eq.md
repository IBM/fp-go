---
title: Record - Equality
hide_title: true
description: Comparing maps for equality using the Eq type class.
sidebar_position: 15
---

import { PageHeader, Section, CodeCard, ApiTable, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="Reference · Collections"
  title="Record Equality"
  lede="Comparing maps for equality. Use the Eq type class to define custom equality semantics for map values."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/record' },
    { label: 'Type Class', value: 'Eq' }
  ]}
/>

---

<Section num="1" title="Core API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Eq` | `func Eq[K comparable, V any](Eq[V]) Eq[map[K]V]` | Create map equality |
</ApiTable>

</Section>

---

<Section num="2" title="Usage Examples">

### Basic Equality

<CodeCard file="basic.go">
{`import (
    R "github.com/IBM/fp-go/v2/record"
    E "github.com/IBM/fp-go/v2/eq"
    N "github.com/IBM/fp-go/v2/number"
)

// Create record equality from value equality
recordEq := R.Eq(N.Eq)

m1 := map[string]int{"a": 1, "b": 2}
m2 := map[string]int{"a": 1, "b": 2}
m3 := map[string]int{"a": 1, "b": 3}
m4 := map[string]int{"a": 1}  // Different size

recordEq.Equals(m1, m2)  // true - same keys and values
recordEq.Equals(m1, m3)  // false - different values
recordEq.Equals(m1, m4)  // false - different keys
`}
</CodeCard>

### String Equality

<CodeCard file="string.go">
{`import S "github.com/IBM/fp-go/v2/string"

stringRecordEq := R.Eq(S.Eq)

m1 := map[string]string{"name": "Alice", "role": "admin"}
m2 := map[string]string{"name": "Alice", "role": "admin"}
m3 := map[string]string{"name": "alice", "role": "admin"}

stringRecordEq.Equals(m1, m2)  // true
stringRecordEq.Equals(m1, m3)  // false - case sensitive
`}
</CodeCard>

### Custom Struct Equality

<CodeCard file="struct.go">
{`type User struct {
    ID   int
    Name string
    Age  int
}

// Compare by ID only
userEq := E.FromEquals(func(a, b User) bool {
    return a.ID == b.ID
})

recordUserEq := R.Eq(userEq)

m1 := map[string]User{
    "alice": {ID: 1, Name: "Alice", Age: 30},
}
m2 := map[string]User{
    "alice": {ID: 1, Name: "Alice Updated", Age: 31},
}

recordUserEq.Equals(m1, m2)  // true - same ID
`}
</CodeCard>

### Case-Insensitive Equality

<CodeCard file="case_insensitive.go">
{`// Case-insensitive string equality
caseInsensitiveEq := E.FromEquals(func(a, b string) bool {
    return strings.ToLower(a) == strings.ToLower(b)
})

recordEq := R.Eq(caseInsensitiveEq)

m1 := map[string]string{"name": "Alice"}
m2 := map[string]string{"name": "ALICE"}

recordEq.Equals(m1, m2)  // true
`}
</CodeCard>

### Nested Map Equality

<CodeCard file="nested.go">
{`// Equality for nested maps
innerEq := R.Eq(N.Eq)
outerEq := R.Eq(innerEq)

m1 := map[string]map[string]int{
    "group1": {"a": 1, "b": 2},
    "group2": {"c": 3},
}
m2 := map[string]map[string]int{
    "group1": {"a": 1, "b": 2},
    "group2": {"c": 3},
}

outerEq.Equals(m1, m2)  // true
`}
</CodeCard>

### Array Value Equality

<CodeCard file="array.go">
{`import A "github.com/IBM/fp-go/v2/array"

// Equality for maps with array values
arrayEq := A.Eq(N.Eq)
recordArrayEq := R.Eq(arrayEq)

m1 := map[string][]int{
    "nums1": {1, 2, 3},
    "nums2": {4, 5},
}
m2 := map[string][]int{
    "nums1": {1, 2, 3},
    "nums2": {4, 5},
}

recordArrayEq.Equals(m1, m2)  // true
`}
</CodeCard>

### Approximate Float Equality

<CodeCard file="float.go">
{`// Approximate equality for floats
const epsilon = 0.0001

floatEq := E.FromEquals(func(a, b float64) bool {
    return math.Abs(a-b) < epsilon
})

recordFloatEq := R.Eq(floatEq)

m1 := map[string]float64{"pi": 3.14159}
m2 := map[string]float64{"pi": 3.14160}

recordFloatEq.Equals(m1, m2)  // true - within epsilon
`}
</CodeCard>

</Section>

---

<Section num="3" title="Common Patterns">

### Configuration Comparison

<CodeCard file="config.go">
{`type Config struct {
    Host    string
    Port    int
    Timeout int
}

configEq := E.FromEquals(func(a, b Config) bool {
    return a.Host == b.Host && a.Port == b.Port
    // Ignore Timeout in comparison
})

recordConfigEq := R.Eq(configEq)

configs1 := map[string]Config{
    "prod": {Host: "api.example.com", Port: 443, Timeout: 30},
}
configs2 := map[string]Config{
    "prod": {Host: "api.example.com", Port: 443, Timeout: 60},
}

recordConfigEq.Equals(configs1, configs2)  // true - timeout ignored
`}
</CodeCard>

### Testing Helper

<CodeCard file="testing.go">
{`func AssertMapsEqual[K comparable, V any](
    t *testing.T,
    eq E.Eq[V],
    expected, actual map[K]V,
) {
    recordEq := R.Eq(eq)
    if !recordEq.Equals(expected, actual) {
        t.Errorf("Maps not equal:\nExpected: %v\nActual: %v",
            expected, actual)
    }
}

// Usage in tests
func TestSomething(t *testing.T) {
    expected := map[string]int{"a": 1, "b": 2}
    actual := processData()
    
    AssertMapsEqual(t, N.Eq, expected, actual)
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

recordStatusEq := R.Eq(statusEq)

m1 := map[string]Status{"service1": Active}
m2 := map[string]Status{"service1": Enabled}

recordStatusEq.Equals(m1, m2)  // true - semantically equal
`}
</CodeCard>

</Section>

---

<Callout type="info">

**Key Equality**: Map keys must be comparable types in Go. The Eq instance only applies to values, not keys.

</Callout>

<Callout type="info">

**Custom Equality**: Define custom Eq instances to implement domain-specific equality semantics, such as case-insensitive comparison or approximate numeric equality.

</Callout>


---

<Pager
  prev={{ to: '/docs/v2/collections/record-traverse', title: 'Record Traverse' }}
  next={{ to: '/docs/v2/collections/record-conversion', title: 'Record Conversion' }}
/>

---
