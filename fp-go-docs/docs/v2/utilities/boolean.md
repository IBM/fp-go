---
title: Boolean
hide_title: true
description: Boolean utilities and type class instances for logical operations in fp-go
sidebar_position: 32
---

import { PageHeader, Section, CodeCard, ApiTable, Pager } from '@site/src/components/content';

<PageHeader
  eyebrow="v2 · Utilities"
  title="Boolean"
  titleAccent="Type Classes"
  lede="The boolean package provides type class instances and utilities for boolean values, including equality, ordering, and monoid operations."
  meta={[
    { label: 'Package', value: 'boolean' },
    { label: 'Since', value: 'v2.0.0' },
    { label: 'Type', value: 'bool' }
  ]}
/>

---

<Section num="1" title="Overview">

The **Boolean** package provides type class instances for the `bool` type, enabling functional operations on boolean values:
- **Eq**: Equality comparison
- **Ord**: Ordering (false < true)
- **Semigroup**: AND and OR operations
- **Monoid**: AND and OR with identity elements

</Section>

---

<Section num="2" title="Equality">

<CodeCard file="boolean_eq.go" tag="example">
{`import B "github.com/IBM/fp-go/boolean"

// Compare booleans for equality
B.Eq.Equals(true, true)    // true
B.Eq.Equals(true, false)   // false
B.Eq.Equals(false, false)  // true`}
</CodeCard>

</Section>

---

<Section num="3" title="Ordering">

Booleans have a natural ordering where `false < true`:

<CodeCard file="boolean_ord.go" tag="example">
{`import B "github.com/IBM/fp-go/boolean"

// Compare booleans (false < true)
B.Ord.Compare(false, true)   // -1 (less than)
B.Ord.Compare(true, false)   // 1  (greater than)
B.Ord.Compare(true, true)    // 0  (equal)
B.Ord.Compare(false, false)  // 0  (equal)

// Derived operations
B.Ord.LessThan(false, true)      // true
B.Ord.GreaterThan(true, false)   // true
B.Ord.LessThanOrEqual(false, false)  // true`}
</CodeCard>

</Section>

---

<Section num="4" title="Semigroup Operations">

Boolean semigroups provide AND and OR operations:

<CodeCard file="boolean_semigroup.go" tag="example">
{`import B "github.com/IBM/fp-go/boolean"

// SemigroupAll: Logical AND
B.SemigroupAll.Concat(true, true)    // true
B.SemigroupAll.Concat(true, false)   // false
B.SemigroupAll.Concat(false, false)  // false

// SemigroupAny: Logical OR
B.SemigroupAny.Concat(false, true)    // true
B.SemigroupAny.Concat(false, false)   // false
B.SemigroupAny.Concat(true, true)     // true`}
</CodeCard>

</Section>

---

<Section num="5" title="Monoid Operations">

Boolean monoids add identity elements to semigroups:

<CodeCard file="boolean_monoid.go" tag="example">
{`import B "github.com/IBM/fp-go/boolean"

// MonoidAll: AND with identity true
B.MonoidAll.Concat(true, true)   // true
B.MonoidAll.Concat(true, false)  // false
B.MonoidAll.Empty()              // true (identity)

// MonoidAny: OR with identity false
B.MonoidAny.Concat(false, true)   // true
B.MonoidAny.Concat(false, false)  // false
B.MonoidAny.Empty()               // false (identity)

// Identity laws
val := true
B.MonoidAll.Concat(B.MonoidAll.Empty(), val)  // true (left identity)
B.MonoidAll.Concat(val, B.MonoidAll.Empty())  // true (right identity)`}
</CodeCard>

</Section>

---

<Section num="6" title="Folding Boolean Arrays">

Use monoids to check if all or any conditions are true:

<CodeCard file="boolean_fold.go" tag="example">
{`import (
    A "github.com/IBM/fp-go/array"
    B "github.com/IBM/fp-go/boolean"
    F "github.com/IBM/fp-go/function"
)

conditions := []bool{true, true, false, true}

// Check if ALL are true
allTrue := F.Pipe2(
    conditions,
    A.Fold(B.MonoidAll),
)
// false (one is false)

// Check if ANY are true
anyTrue := F.Pipe2(
    conditions,
    A.Fold(B.MonoidAny),
)
// true (at least one is true)

// Empty array cases
emptyConditions := []bool{}

allEmpty := F.Pipe2(emptyConditions, A.Fold(B.MonoidAll))
// true (identity for AND)

anyEmpty := F.Pipe2(emptyConditions, A.Fold(B.MonoidAny))
// false (identity for OR)`}
</CodeCard>

</Section>

---

<Section num="7" title="Validation Example">

Check if all validations pass:

<CodeCard file="boolean_validation.go" tag="example">
{`type Validation struct {
    Field   string
    IsValid bool
}

validations := []Validation{
    {Field: "email", IsValid: true},
    {Field: "age", IsValid: true},
    {Field: "name", IsValid: false},
}

// Check if all validations pass
allValid := F.Pipe3(
    validations,
    A.Map(func(v Validation) bool { return v.IsValid }),
    A.Fold(B.MonoidAll),
)
// false (name validation failed)

// Check if any validation passed
anyValid := F.Pipe3(
    validations,
    A.Map(func(v Validation) bool { return v.IsValid }),
    A.Fold(B.MonoidAny),
)
// true (email and age passed)

// Find which validations failed
failed := F.Pipe2(
    validations,
    A.Filter(func(v Validation) bool { return !v.IsValid }),
)
// []Validation{{Field: "name", IsValid: false}}`}
</CodeCard>

</Section>

---

<Section num="8" title="Permission Checking">

Combine permission checks:

<CodeCard file="boolean_permissions.go" tag="example">
{`type User struct {
    IsAdmin     bool
    IsVerified  bool
    HasAccess   bool
}

user := User{
    IsAdmin:    false,
    IsVerified: true,
    HasAccess:  true,
}

// Check if user has all required permissions
permissions := []bool{
    user.IsVerified,
    user.HasAccess,
}

hasAllPermissions := F.Pipe2(
    permissions,
    A.Fold(B.MonoidAll),
)
// true

// Check if user has any admin privilege
adminPrivileges := []bool{
    user.IsAdmin,
    user.IsVerified && user.HasAccess,
}

hasAnyAdminPrivilege := F.Pipe2(
    adminPrivileges,
    A.Fold(B.MonoidAny),
)
// true (second condition is true)`}
</CodeCard>

</Section>

---

<Section num="9" title="API Reference">

<ApiTable>
| Instance | Type | Description |
|----------|------|-------------|
| `Eq` | `Eq[bool]` | Equality comparison |
| `Ord` | `Ord[bool]` | Ordering (false < true) |
| `SemigroupAll` | `Semigroup[bool]` | Logical AND |
| `SemigroupAny` | `Semigroup[bool]` | Logical OR |
| `MonoidAll` | `Monoid[bool]` | AND with identity true |
| `MonoidAny` | `Monoid[bool]` | OR with identity false |
</ApiTable>

**Monoid Identities:**
- `MonoidAll.Empty()` returns `true` (AND identity)
- `MonoidAny.Empty()` returns `false` (OR identity)

</Section>

---

<Section num="10" title="Related Concepts">

**Common Use Cases:**
- Validation aggregation
- Permission checking
- Condition evaluation
- Boolean algebra operations

**See Also:**
- [Predicate](./predicate.md) - Boolean-valued functions
- [Monoid](./monoid.md) - Understanding monoid operations
- [Semigroup](./semigroup.md) - Understanding semigroup operations

</Section>

---

<Pager
  prev={{ to: '/docs/v2/utilities/predicate', title: 'Predicate' }}
  next={{ to: '/docs/v2/utilities/number', title: 'Number' }}
/>