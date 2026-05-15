---
title: Array - Equality
hide_title: true
description: Equality checking for arrays using the Eq type class.
sidebar_position: 8
---

<PageHeader
  eyebrow="Reference · Collections"
  title="Array"
  titleAccent="Equality"
  lede="Equality checking for arrays using the Eq type class. Compare arrays element-wise with custom equality functions."
  meta={[
    { label: 'Package', value: 'github.com/IBM/fp-go/v2/array' },
    { label: 'Type Class', value: 'Eq' }
  ]}
/>

<Section id="api" number="01" title="Core" titleAccent="API">

<ApiTable>
| Function | Signature | Description |
|----------|-----------|-------------|
| `Eq` | `func Eq[A any](eq Eq[A]) Eq[[]A]` | Create array equality from element equality |
</ApiTable>

</Section>

<Section id="examples" number="02" title="Usage" titleAccent="Examples">

### Basic Equality

<CodeCard file="basic.go">
{`import (
    A "github.com/IBM/fp-go/v2/array"
    E "github.com/IBM/fp-go/v2/eq"
    N "github.com/IBM/fp-go/v2/number"
)

// Create array equality from element equality
arrayEq := A.Eq(N.Eq)

arr1 := []int{1, 2, 3}
arr2 := []int{1, 2, 3}
arr3 := []int{1, 2, 4}

arrayEq.Equals(arr1, arr2)  // true
arrayEq.Equals(arr1, arr3)  // false
`}
</CodeCard>

### Custom Equality

<CodeCard file="custom.go">
{`type User struct {
    ID   int
    Name string
}

// Compare by ID only
userEq := E.FromEquals(func(a, b User) bool {
    return a.ID == b.ID
})

arrayUserEq := A.Eq(userEq)

users1 := []User{{ID: 1, Name: "Alice"}}
users2 := []User{{ID: 1, Name: "Alice Updated"}}

arrayUserEq.Equals(users1, users2)  // true (same ID)
`}
</CodeCard>

### String Arrays

<CodeCard file="strings.go">
{`import S "github.com/IBM/fp-go/v2/string"

arrayStrEq := A.Eq(S.Eq)

arr1 := []string{"hello", "world"}
arr2 := []string{"hello", "world"}
arr3 := []string{"hello", "go"}

arrayStrEq.Equals(arr1, arr2)  // true
arrayStrEq.Equals(arr1, arr3)  // false
`}
</CodeCard>

</Section>
