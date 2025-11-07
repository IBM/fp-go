# Lens Generator Example

This example demonstrates the lens code generator for Go structs.

## Overview

The lens generator automatically creates lens types for Go structs annotated with `fp-go:Lens`. Lenses provide a functional way to access and update nested immutable data structures.

## Usage

### 1. Annotate Your Structs

Add the `fp-go:Lens` annotation in a comment above your struct declaration:

```go
// fp-go:Lens
type Person struct {
    Name  string
    Age   int
    Email string
    Phone *string  // Pointer fields generate LensO (optional lens)
}
```

### 2. Generate Lens Code

Run the generator command:

```bash
go run ../../main.go lens --dir . --filename gen.go
```

Or use it as a go generate directive:

```go
//go:generate go run ../../main.go lens --dir . --filename gen.go
```

### 3. Use the Generated Lenses

The generator creates:
- A `<TypeName>Lens` struct with a lens for each exported field
- A `Make<TypeName>Lens()` constructor function

```go
// Create lenses
lenses := MakePersonLens()

// Get a field value
name := lenses.Name.Get(person)

// Set a field value (returns a new instance)
updated := lenses.Name.Set("Bob")(person)

// Modify a field value
incremented := F.Pipe1(
    lenses.Age,
    L.Modify[Person](func(age int) int { return age + 1 }),
)(person)
```

## Features

### Optional Fields (LensO)

Pointer fields automatically generate `LensO` (optional lenses) that work with `Option[*T]`:

```go
// fp-go:Lens
type Person struct {
    Name  string
    Phone *string  // Generates LensO[Person, *string]
}

lenses := MakePersonLens()
person := Person{Name: "Alice", Phone: nil}

// Get returns Option[*string]
phoneOpt := lenses.Phone.Get(person)  // None

// Set with Some
phone := "555-1234"
updated := lenses.Phone.Set(O.Some(&phone))(person)

// Set with None (clears the field)
cleared := lenses.Phone.Set(O.None[*string]())(person)
```

### Immutable Updates

All lens operations return new instances, leaving the original unchanged:

```go
person := Person{Name: "Alice", Age: 30}
updated := lenses.Name.Set("Bob")(person)
// person.Name is still "Alice"
// updated.Name is "Bob"
```

### Lens Composition

Compose lenses to access deeply nested fields:

```go
// Access company.CEO.Name
ceoNameLens := F.Pipe1(
    companyLenses.CEO,
    L.Compose[Company](personLenses.Name),
)

name := ceoNameLens.Get(company)
updated := ceoNameLens.Set("Jane")(company)
```

### Type Safety

All operations are type-safe at compile time:

```go
// Compile error: type mismatch
lenses.Age.Set("not a number")(person)
```

## Generated Code Structure

For each annotated struct, the generator creates:

```go
// Lens struct with a lens for each field
type PersonLens struct {
    Name  L.Lens[Person, string]
    Age   L.Lens[Person, int]
    Email L.Lens[Person, string]
}

// Constructor function
func MakePersonLens() PersonLens {
    return PersonLens{
        Name: L.MakeLens(
            func(s Person) string { return s.Name },
            func(s Person, v string) Person { s.Name = v; return s },
        ),
        // ... other fields
    }
}
```

## Generated Code Structure

For each annotated struct, the generator creates:

```go
// Regular field generates Lens
type PersonLens struct {
    Name  L.Lens[Person, string]
    Phone LO.LensO[Person, *string]  // Pointer field generates LensO
}

// Constructor function
func MakePersonLens() PersonLens {
    return PersonLens{
        Name: L.MakeLens(
            func(s Person) string { return s.Name },
            func(s Person, v string) Person { s.Name = v; return s },
        ),
        Phone: L.MakeLens(
            func(s Person) O.Option[*string] { return O.FromNillable(s.Phone) },
            func(s Person, v O.Option[*string]) Person {
                s.Phone = O.GetOrElse(func() *string { return nil })(v)
                return s
            },
        ),
    }
}
```

## Command Options

- `--dir`: Directory to scan for Go files (default: ".")
- `--filename`: Name of the generated file (default: "gen.go")

## Notes

- Only pointer fields (`*T`) generate `LensO` (optional lenses)
- The `json:"...,omitempty"` tag alone does not make a field optional in the lens generator
- Pointer fields work with `Option[*T]` using `FromNillable` and `GetOrElse`

## Examples

See `example_test.go` for comprehensive examples including:
- Basic lens operations (Get, Set, Modify)
- Nested struct access
- Lens composition
- Complex data structure manipulation