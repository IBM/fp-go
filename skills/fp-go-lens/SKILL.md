---
name: fp-go-lens
description: Use this skill when working with lenses and optics in Go using the fp-go library (github.com/IBM/fp-go/v2/optics/lens). Trigger on mentions of lenses, optics, MakeLens, MakeLensRef, lens composition, immutable updates to nested structs, accessing nested data structures, Compose, ComposeOption, FromNillable, Modify, getter/setter patterns, or functional updates to Go structs. Also trigger when the user needs to update deeply nested fields immutably or work with optional fields in struct hierarchies.
---

# fp-go Lenses for Structs

## Overview

Lenses are functional optics that provide immutable access to nested data structures. In fp-go, a [`Lens[S, A]`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#Lens) focuses on a property `A` within a structure `S`, allowing you to read and update that property without mutating the original structure.

```go
import L "github.com/IBM/fp-go/v2/optics/lens"
```

## Core Concept

A lens consists of two operations:
- **Get**: Extract a value from a structure
- **Set**: Create a new structure with an updated value

```go
type Lens[S, A any] struct {
    Get func(s S) A
    Set func(a A) Endomorphism[S]  // func(a A) func(S) S
}
```

## Creating Lenses: MakeLens vs MakeLensRef

### MakeLens - For Value Types

Use [`MakeLens`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLens) when working with **structs passed by value**. The setter must create a copy of the data, which happens automatically for value types.

```go
type Person struct {
    Name string
    Age  int
}

// Getter and setter for value types
func (p Person) GetName() string {
    return p.Name
}

func (p Person) SetName(name string) Person {
    p.Name = name  // Automatic copy because Person is passed by value
    return p
}

// Create lens for value type
nameLens := L.MakeLens(Person.GetName, Person.SetName)

// Usage
person := Person{Name: "Alice", Age: 30}
updated := nameLens.Set("Bob")(person)

// person.Name is still "Alice" (immutable)
// updated.Name is "Bob"
```

**Key Point**: When you pass a struct by value in Go, modifications to the parameter don't affect the original. This automatic copying makes [`MakeLens`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLens) safe for value types.

### MakeLensRef - For Pointer Types

Use [`MakeLensRef`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLensRef) when working with **pointers to structs**. The setter does NOT need to create a copy—[`MakeLensRef`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLensRef) automatically wraps your setter to copy the pointer before modification.

```go
type Address struct {
    street string
    city   string
}

// Getter and setter for pointer types
func (a *Address) GetStreet() string {
    return a.street
}

func (a *Address) SetStreet(street string) *Address {
    a.street = street  // Direct modification - MakeLensRef will handle copying
    return a
}

// Create lens for pointer type
streetLens := L.MakeLensRef((*Address).GetStreet, (*Address).SetStreet)

// Usage
addr := &Address{street: "Main St", city: "Boston"}
updated := streetLens.Set("Oak Ave")(addr)

// addr.street is still "Main St" (original unchanged)
// updated.street is "Oak Ave" (new copy)
```

**Key Point**: [`MakeLensRef`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLensRef) internally calls `setCopy` which creates a shallow copy (`cpy := *s`) before applying your setter. This ensures immutability even though your setter modifies the pointer directly.

## Comparison Table

| Aspect | MakeLens | MakeLensRef |
|--------|----------|-------------|
| **Input Type** | Value type `S` | Pointer type `*S` |
| **Setter Signature** | `func(S, A) S` | `func(*S, A) *S` |
| **Copy Responsibility** | Caller must copy | Automatic (handled by framework) |
| **Use When** | Struct passed by value | Struct passed by pointer |
| **Example** | `Person` struct | `*Address` pointer |

## Composing Lenses

Lenses can be composed to access deeply nested properties:

```go
type Street struct {
    num  int
    name string
}

type Address struct {
    city   string
    street *Street
}

// Lenses for each level
streetNameLens := L.MakeLensRef((*Street).GetName, (*Street).SetName)
addressStreetLens := L.MakeLensRef((*Address).GetStreet, (*Address).SetStreet)

// Compose to access street name through address
streetNameInAddress := L.Compose[*Address](streetNameLens)(addressStreetLens)

// Usage
addr := &Address{
    city: "Boston",
    street: &Street{num: 123, name: "Main St"},
}

updated := streetNameInAddress.Set("Oak Ave")(addr)
// addr.street.name is still "Main St"
// updated.street.name is "Oak Ave"
```

## Working with Optional Fields

### FromNillable - Pointer to Option

Convert a lens for a pointer field into a lens for an [`Option`](https://pkg.go.dev/github.com/IBM/fp-go/v2/option#Option):

```go
type Company struct {
    name    string
    address *Address  // Optional field
}

addressLens := L.MakeLens(Company.GetAddress, Company.SetAddress)
optionalAddressLens := L.FromNillable(addressLens)

// Get returns Option[*Address]
company := Company{name: "Acme"}
result := optionalAddressLens.Get(company)  // None

// Set with Some creates the address
updated := optionalAddressLens.Set(O.Some(&Address{...}))(company)

// Set with None removes the address
cleared := optionalAddressLens.Set(O.None[*Address]())(updated)
```

### ComposeOption - Compose with Optional Parent

When the parent structure is optional, use [`ComposeOption`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#ComposeOption) with a default value:

```go
type Person struct {
    name    string
    address *Address  // Optional
}

defaultAddress := &Address{city: "Unknown", street: &Street{}}

addressLens := L.FromNillable(L.MakeLens(Person.GetAddress, Person.SetAddress))
streetLens := L.MakeLensRef((*Address).GetStreet, (*Address).SetStreet)

// Compose with default for missing address
streetInPerson := F.Pipe1(
    addressLens,
    L.ComposeOption[Person, *Street](defaultAddress)(streetLens),
)

// If person has no address, uses defaultAddress
person := Person{name: "Alice"}
updated := streetInPerson.Set(O.Some(&Street{name: "Main St"}))(person)
// Creates address with default values and sets the street
```

## Modifying Values

Use [`Modify`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#Modify) to transform a value through a lens:

```go
type Counter struct {
    value int
}

valueLens := L.MakeLens(
    func(c Counter) int { return c.value },
    func(c Counter, v int) Counter { c.value = v; return c },
)

counter := Counter{value: 5}

// Increment by applying a function
incremented := L.Modify[Counter](func(v int) int { 
    return v + 1 
})(valueLens)(counter)

// counter.value is still 5
// incremented.value is 6
```

## Best Practices

1. **Choose the Right Constructor**:
   - Use [`MakeLens`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLens) for value types (structs passed by value)
   - Use [`MakeLensRef`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLensRef) for pointer types (structs passed by pointer)

2. **Setter Implementation**:
   - With [`MakeLens`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLens): Modify the value parameter directly (it's already a copy)
   - With [`MakeLensRef`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLensRef): Modify the pointer parameter directly (framework handles copying)

3. **Immutability**:
   - Lenses always return new structures
   - Original data is never modified
   - Safe for concurrent access

4. **Composition**:
   - Build complex lenses from simple ones using [`Compose`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#Compose)
   - Use [`ComposeOption`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#ComposeOption) when dealing with optional nested structures

## Complete Example

```go
package main

import (
    F "github.com/IBM/fp-go/v2/function"
    L "github.com/IBM/fp-go/v2/optics/lens"
    O "github.com/IBM/fp-go/v2/option"
)

type Street struct {
    num  int
    name string
}

type Address struct {
    city   string
    street *Street
}

type Person struct {
    name    string
    age     int
    address *Address
}

// Getters and setters
func (s *Street) GetName() string { return s.name }
func (s *Street) SetName(name string) *Street { s.name = name; return s }

func (a *Address) GetStreet() *Street { return a.street }
func (a *Address) SetStreet(s *Street) *Address { a.street = s; return a }

func (p Person) GetAddress() *Address { return p.address }
func (p Person) SetAddress(a *Address) Person { p.address = a; return p }

func main() {
    // Create lenses
    streetNameLens := L.MakeLensRef((*Street).GetName, (*Street).SetName)
    addressStreetLens := L.MakeLensRef((*Address).GetStreet, (*Address).SetStreet)
    personAddressLens := L.FromNillable(L.MakeLens(Person.GetAddress, Person.SetAddress))
    
    // Compose to access street name through person
    defaultAddress := &Address{city: "Unknown", street: &Street{}}
    streetNameInPerson := F.Pipe2(
        personAddressLens,
        L.ComposeOption[Person, *Street](defaultAddress)(addressStreetLens),
        L.ComposeOption[Person, string](&Street{})(streetNameLens),
    )
    
    // Usage
    person := Person{
        name: "Alice",
        age:  30,
        address: &Address{
            city: "Boston",
            street: &Street{num: 123, name: "Main St"},
        },
    }
    
    // Update street name immutably
    updated := streetNameInPerson.Set(O.Some("Oak Ave"))(person)
    
    // person.address.street.name is still "Main St"
    // updated.address.street.name is "Oak Ave"
}
```

## Import Reference

```go
import (
    L "github.com/IBM/fp-go/v2/optics/lens"
    F "github.com/IBM/fp-go/v2/function"
    O "github.com/IBM/fp-go/v2/option"
)
```

## Further Reading

- [Introduction to optics: lenses and prisms](https://medium.com/@gcanti/introduction-to-optics-lenses-and-prisms-3230e73bfcfe)
- [Lens Laws](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens/testing) - Property-based tests for lens correctness