# 🔍 Optics

Functional optics for composable data access and manipulation in Go.

## 📖 Overview

Optics are first-class, composable references to parts of data structures. They provide a uniform interface for reading, writing, and transforming nested immutable data without verbose boilerplate code.

## ✨ Why Use Optics?

Optics bring powerful benefits to your Go code:

- **🎯 Composability**: Optics naturally compose with each other and with monadic operations, enabling elegant data transformations through function composition
- **🔒 Immutability**: Work with immutable data structures without manual copying and updating
- **🧩 Type Safety**: Leverage Go's type system to catch errors at compile time
- **📦 Reusability**: Define data access patterns once and reuse them throughout your codebase
- **🎨 Expressiveness**: Write declarative code that clearly expresses intent
- **🔄 Bidirectionality**: Read and write through the same abstraction
- **🚀 Productivity**: Eliminate boilerplate for nested data access and updates
- **🧪 Testability**: Optics are pure functions, making them easy to test and reason about

### 🔗 Composition with Monadic Operations

One of the most powerful features of optics is their natural composition with monadic operations. Optics integrate seamlessly with `fp-go`'s monadic types like [`Option`](https://pkg.go.dev/github.com/IBM/fp-go/v2/option), [`Either`](https://pkg.go.dev/github.com/IBM/fp-go/v2/either), [`Result`](https://pkg.go.dev/github.com/IBM/fp-go/v2/result), and [`IO`](https://pkg.go.dev/github.com/IBM/fp-go/v2/io), allowing you to:

- Chain optional field access with [`Option`](https://pkg.go.dev/github.com/IBM/fp-go/v2/option) monads
- Handle errors gracefully with [`Either`](https://pkg.go.dev/github.com/IBM/fp-go/v2/either) or [`Result`](https://pkg.go.dev/github.com/IBM/fp-go/v2/result) monads
- Perform side effects with [`IO`](https://pkg.go.dev/github.com/IBM/fp-go/v2/io) monads
- Combine multiple optics in a single pipeline using [`Pipe`](https://pkg.go.dev/github.com/IBM/fp-go/v2/function#Pipe1)

This composability enables you to build complex data transformations from simple, reusable building blocks.

## 🚀 Quick Start

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

## 🛠️ Core Optics Types

### 🔎 [Lens](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens) - Product Types (Structs)
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

### 🔀 [Prism](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/prism) - Sum Types (Variants)
Focus on one variant of a sum type. Provides optional get and definite set.

**Use when:** Working with [`Either`](https://pkg.go.dev/github.com/IBM/fp-go/v2/either), [`Result`](https://pkg.go.dev/github.com/IBM/fp-go/v2/result), or custom sum types.

**💡 Important Use Case - Generalized Constructors for Do Notation:**

[Prisms](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/prism) act as generalized constructors, making them invaluable for `Do` notation workflows. The prism's `ReverseGet` function serves as a constructor that creates a value of the sum type from a specific variant. This is particularly useful when building up complex data structures step-by-step in monadic contexts:

```go
import "github.com/IBM/fp-go/v2/optics/prism"

// Prism for the Success variant
successPrism := prism.MakePrism(
    func(r Result) option.Option[int] {
        if s, ok := r.(Success); ok {
            return option.Some(s.Value)
        }
        return option.None[int]()
    },
    func(v int) Result { return Success{Value: v} }, // Constructor!
)

// Use in Do notation to construct values
result := F.Pipe2(
    computeValue(),
    option.Map(func(v int) int { return v * 2 }),
    option.Map(successPrism.ReverseGet), // Construct Result from int
)
```

### 🔄 [Iso](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/iso) - Isomorphisms
Bidirectional transformation between equivalent types with no information loss.

**Use when:** Converting between equivalent representations (e.g., Celsius ↔ Fahrenheit).

```go
import "github.com/IBM/fp-go/v2/optics/iso"

celsiusToFahrenheit := iso.MakeIso(
    func(c float64) float64 { return c*9/5 + 32 },
    func(f float64) float64 { return (f - 32) * 5 / 9 },
)
```

### ❓ [Optional](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/optional) - Maybe Values
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

### 🔢 [Traversal](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/traversal) - Multiple Values
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
    traversal.Modify[[]int, int](N.Mul(2)),
)
// Result: [2, 4, 6, 8, 10]
```

## 🔗 Composition

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

## ⚙️ Auto-Generation with [`go generate`](https://go.dev/blog/generate)

Lenses can be automatically generated using the `fp-go` CLI tool and a simple annotation. This eliminates boilerplate and ensures consistency.

### 📝 How to Use

1. **Annotate your struct** with the `fp-go:Lens` comment:

```go
//go:generate go run github.com/IBM/fp-go/v2 lens --dir . --filename gen_lens.go

// fp-go:Lens
type Person struct {
    Name  string
    Age   int
    Email string
    Phone *string  // Optional field
}

// fp-go:Lens
type Config struct {
    PublicField  string
    privateField int  // Unexported fields are supported!
}
```

**Note:** The generator supports both exported (uppercase) and unexported (lowercase) fields. Generated lenses for unexported fields will have lowercase names and can only be used within the same package as the struct.

2. **Run `go generate`**:

```bash
go generate ./...
```

3. **Use the generated lenses**:

```go
// Generated code creates PersonLenses, PersonRefLenses, and PersonPrisms
lenses := MakePersonLenses()

person := Person{Name: "Alice", Age: 30, Email: "alice@example.com"}

// Use the generated lenses
updatedPerson := lenses.Age.Set(31)(person)
name := lenses.Name.Get(person)

// Optional lenses for zero-value handling
personWithEmail := lenses.EmailO.Set(option.Some("new@example.com"))(person)
```

### 🎁 What Gets Generated

For each annotated struct, the generator creates:

- **`StructNameLenses`**: Lenses for value types with optional variants (`LensO`) for comparable fields
- **`StructNameRefLenses`**: Lenses for pointer types with prisms for constructing values
- **`StructNamePrisms`**: Prisms for all fields, useful for partial construction
- Constructor functions: `MakeStructNameLenses()`, `MakeStructNameRefLenses()`, `MakeStructNamePrisms()`

The generator supports:
- ✅ Generic types with type parameters
- ✅ Embedded structs (fields are promoted)
- ✅ Optional fields (pointers and `omitempty` tags)
- ✅ Custom package imports
- ✅ **Unexported fields** (lowercase names) - lenses will have lowercase names matching the field names

See [samples/lens](../samples/lens) for complete examples.

## 📊 Optics Hierarchy

```
[Iso](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/iso)[S, A]
    ↓
[Lens](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens)[S, A]
    ↓
[Optional](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/optional)[S, A]
    ↓
[Traversal](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/traversal)[S, A]

[Prism](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/prism)[S, A]
    ↓
[Optional](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/optional)[S, A]
    ↓
[Traversal](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/traversal)[S, A]
```

More specific optics can be converted to more general ones.

## 📦 Package Structure

### Core Optics
- **[optics/lens](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens)**: Lenses for product types (structs)
- **[optics/prism](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/prism)**: Prisms for sum types ([`Either`](https://pkg.go.dev/github.com/IBM/fp-go/v2/either), [`Result`](https://pkg.go.dev/github.com/IBM/fp-go/v2/result), etc.)
- **[optics/iso](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/iso)**: Isomorphisms for equivalent types
- **[optics/optional](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/optional)**: Optional optics for maybe values
- **[optics/traversal](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/traversal)**: Traversals for multiple values

### Utilities
- **[optics/builder](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/builder)**: Builder pattern for constructing complex optics
- **[optics/codec](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/codec)**: Type-safe encoding/decoding with validation
  - Provides `Type[A, O, I]` for bidirectional transformations with validation
  - Includes codecs for primitives (String, Int, Bool), collections (Array), and sum types (Either)
  - Supports refinement types and codec composition via `Pipe`
  - Integrates validation errors with context tracking

### Specialized Sub-packages
Each core optics package includes specialized sub-packages for common patterns:
- **array**: Optics for arrays/slices
- **either**: Optics for [`Either`](https://pkg.go.dev/github.com/IBM/fp-go/v2/either) types
- **option**: Optics for [`Option`](https://pkg.go.dev/github.com/IBM/fp-go/v2/option) types
- **record**: Optics for maps

## 📚 Documentation

For detailed documentation on each optic type, see:
- [Main Package Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics)
- [Lens Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens)
- [Prism Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/prism)
- [Iso Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/iso)
- [Optional Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/optional)
- [Traversal Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/traversal)

## 🌐 Further Reading

### Haskell Lens Library
The concepts in this library are inspired by the powerful [Haskell lens library](https://hackage.haskell.org/package/lens), which pioneered many of these abstractions.

### Articles and Resources
- [Introduction to optics: lenses and prisms](https://medium.com/@gcanti/introduction-to-optics-lenses-and-prisms-3230e73bfcfe) by Giulio Canti - Excellent introduction to optics concepts
- [Lenses in Functional Programming](https://www.schoolofhaskell.com/school/to-infinity-and-beyond/pick-of-the-week/a-little-lens-starter-tutorial) - Tutorial on lens fundamentals
- [Profunctor Optics: The Categorical View](https://bartoszmilewski.com/2017/07/07/profunctor-optics-the-categorical-view/) by Bartosz Milewski - Deep dive into the theory
- [Why Optics?](https://www.tweag.io/blog/2022-01-06-optics-vs-lenses/) - Discussion of benefits and use cases

### Why Functional Optics?
Functional optics solve real problems in software development:
- **Nested Updates**: Eliminate deeply nested field access patterns
- **Immutability**: Make working with immutable data practical and ergonomic
- **Abstraction**: Separate data access patterns from business logic
- **Composition**: Build complex operations from simple, reusable pieces
- **Type Safety**: Catch errors at compile time rather than runtime

## 💡 Examples

See the [samples/lens](../samples/lens) directory for complete working examples, including:
- Basic lens usage
- Lens composition
- Auto-generated lenses
- Prism usage for sum types
- Integration with monadic operations

## 📄 License

Apache License 2.0 - See LICENSE file for details.
