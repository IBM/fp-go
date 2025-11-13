# Optics

Functional optics for composable data access and manipulation in Go.

## Overview

Optics are first-class, composable references to parts of data structures. They provide a uniform interface for reading, writing, and transforming nested immutable data without verbose boilerplate code.

## Quick Start

```go
import (
    "github.com/IBM/fp-go/v2/optics/lens"
    F "github.com/IBM/fp-go/v2/function"
)

type Person struct {
    Name string
    Age  int
}

// Create a lens for the Name field
nameLens := lens.MakeLens(
    func(p Person) string { return p.Name },
    func(p Person, name string) Person {
        p.Name = name
        return p
    },
)

person := Person{Name: "Alice", Age: 30}

// Get the name
name := nameLens.Get(person) // "Alice"

// Set a new name (returns a new Person)
updated := nameLens.Set("Bob")(person)
// person.Name is still "Alice", updated.Name is "Bob"
```

## Core Optics Types

### Lens - Product Types (Structs)
Focus on a single field within a struct. Provides get and set operations.

**Use when:** Working with struct fields that always exist.

```go
ageLens := lens.MakeLens(
    func(p Person) int { return p.Age },
    func(p Person, age int) Person {
        p.Age = age
        return p
    },
)
```

### Prism - Sum Types (Variants)
Focus on one variant of a sum type. Provides optional get and definite set.

**Use when:** Working with Either, Result, or custom sum types.

```go
import "github.com/IBM/fp-go/v2/optics/prism"

successPrism := prism.MakePrism(
    func(r Result) option.Option[int] {
        if s, ok := r.(Success); ok {
            return option.Some(s.Value)
        }
        return option.None[int]()
    },
    func(v int) Result { return Success{Value: v} },
)
```

### Iso - Isomorphisms
Bidirectional transformation between equivalent types with no information loss.

**Use when:** Converting between equivalent representations (e.g., Celsius ↔ Fahrenheit).

```go
import "github.com/IBM/fp-go/v2/optics/iso"

celsiusToFahrenheit := iso.MakeIso(
    func(c float64) float64 { return c*9/5 + 32 },
    func(f float64) float64 { return (f - 32) * 5 / 9 },
)
```

### Optional - Maybe Values
Focus on a value that may or may not exist.

**Use when:** Working with nullable fields or values that may be absent.

```go
import "github.com/IBM/fp-go/v2/optics/optional"

timeoutOptional := optional.MakeOptional(
    func(c Config) option.Option[*int] {
        return option.FromNillable(c.Timeout)
    },
    func(c Config, t *int) Config {
        c.Timeout = t
        return c
    },
)
```

### Traversal - Multiple Values
Focus on multiple values simultaneously, allowing batch operations.

**Use when:** Working with collections or updating multiple fields at once.

```go
import (
    "github.com/IBM/fp-go/v2/optics/traversal"
    TA "github.com/IBM/fp-go/v2/optics/traversal/array"
)

numbers := []int{1, 2, 3, 4, 5}

// Double all elements
doubled := F.Pipe2(
    numbers,
    TA.Traversal[int](),
    traversal.Modify[[]int, int](func(n int) int { return n * 2 }),
)
// Result: [2, 4, 6, 8, 10]
```

## Composition

The real power of optics comes from composition:

```go
type Company struct {
    Name    string
    Address Address
}

type Address struct {
    Street string
    City   string
}

// Individual lenses
addressLens := lens.MakeLens(
    func(c Company) Address { return c.Address },
    func(c Company, a Address) Company {
        c.Address = a
        return c
    },
)

cityLens := lens.MakeLens(
    func(a Address) string { return a.City },
    func(a Address, city string) Address {
        a.City = city
        return a
    },
)

// Compose to access city directly from company
companyCityLens := F.Pipe1(
    addressLens,
    lens.Compose[Company](cityLens),
)

company := Company{
    Name: "Acme Corp",
    Address: Address{Street: "Main St", City: "NYC"},
}

city := companyCityLens.Get(company)           // "NYC"
updated := companyCityLens.Set("Boston")(company)
```

## Optics Hierarchy

```
Iso[S, A]
    ↓
Lens[S, A]
    ↓
Optional[S, A]
    ↓
Traversal[S, A]

Prism[S, A]
    ↓
Optional[S, A]
    ↓
Traversal[S, A]
```

More specific optics can be converted to more general ones.

## Package Structure

- **optics/lens**: Lenses for product types (structs)
- **optics/prism**: Prisms for sum types (Either, Result, etc.)
- **optics/iso**: Isomorphisms for equivalent types
- **optics/optional**: Optional optics for maybe values
- **optics/traversal**: Traversals for multiple values

Each package includes specialized sub-packages for common patterns:
- **array**: Optics for arrays/slices
- **either**: Optics for Either types
- **option**: Optics for Option types
- **record**: Optics for maps

## Documentation

For detailed documentation on each optic type, see:
- [Main Package Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics)
- [Lens Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens)
- [Prism Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/prism)
- [Iso Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/iso)
- [Optional Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/optional)
- [Traversal Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/traversal)

## Further Reading

For an introduction to functional optics concepts:
- [Introduction to optics: lenses and prisms](https://medium.com/@gcanti/introduction-to-optics-lenses-and-prisms-3230e73bfcfe) by Giulio Canti

## Examples

See the [samples/lens](../samples/lens) directory for complete working examples.

## License

Apache License 2.0 - See LICENSE file for details.
