# 🏗️ Builder Pattern with fp-go

This package demonstrates a functional builder pattern using fp-go's optics library. It shows how to construct and validate objects using lenses, prisms, and codecs, separating the building phase from validation.

## 📋 Overview

The builder pattern here uses two key types:

- **`PartialPerson`** 🚧: An intermediate type with unvalidated fields (raw `string` and `int`)
- **`Person`** ✅: A validated type with refined fields (`NonEmptyString` and `AdultAge`)

The pattern provides two approaches for validation:

1. **Prism-based validation** 🔍 (simple, no error messages)
2. **Codec-based validation** 📝 (detailed error reporting)

## 🎯 Core Concepts

### 1. 🔧 Auto-Generated Lenses

The `fp-go:Lens` directive in `types.go` generates lens accessors for both types:

```go
// fp-go:Lens
type PartialPerson struct {
    name string
    age  int
}

// fp-go:Lens
type Person struct {
    Name NonEmptyString
    Age  AdultAge
}
```

This generates:
- `partialPersonLenses` with `.name` and `.age` lenses
- `personLenses` with `.Name` and `.Age` lenses

### 2. 🎁 Exporting Setters as `WithXXX` Methods

The lens setters are exported as builder methods:

```go
// WithName sets the Name field of a PartialPerson
WithName = partialPersonLenses.name.Set

// WithAge sets the Age field of a PartialPerson
WithAge = partialPersonLenses.age.Set
```

These return `Endomorphism[*PartialPerson]` functions that can be composed:

```go
builder := F.Pipe1(
    A.From(
        WithName("Alice"),
        WithAge(25),
    ),
    allOfPartialPerson,
)
partial := builder(&PartialPerson{})
```

Or use the convenience function:

```go
builder := MakePerson("Alice", 25)
```

## 🔍 Approach 1: Prism-Based Validation (No Error Messages)

### Creating Validation Prisms

Define prisms that validate individual fields:

> 💡 **Tip**: The `optics/prism` package provides many helpful out-of-the-box prisms for common validations, including:
> - `NonEmptyString()` - validates non-empty strings
> - `ParseInt()`, `ParseInt64()` - parses integers from strings
> - `ParseFloat32()`, `ParseFloat64()` - parses floats from strings
> - `ParseBool()` - parses booleans from strings
> - `ParseDate(layout)` - parses dates with custom layouts
> - `ParseURL()` - parses URLs
> - `FromZero()`, `FromNonZero()` - validates zero/non-zero values
> - `RegexMatcher()`, `RegexNamedMatcher()` - regex-based validation
> - `FromOption()`, `FromEither()`, `FromResult()` - extracts from monadic types
> - And many more! Check `optics/prism/prisms.go` for the full list.
>
> For custom validation logic, create your own prisms:

```go
namePrism = prism.MakePrismWithName(
    func(s string) Option[NonEmptyString] {
        if S.IsEmpty(s) {
            return option.None[NonEmptyString]()
        }
        return option.Of(NonEmptyString(s))
    },
    func(ns NonEmptyString) string {
        return string(ns)
    },
    "NonEmptyString",
)

agePrism = prism.MakePrismWithName(
    func(a int) Option[AdultAge] {
        if a < 18 {
            return option.None[AdultAge]()
        }
        return option.Of(AdultAge(a))
    },
    func(aa AdultAge) int {
        return int(aa)
    },
    "AdultAge",
)
```

### 🎭 Creating the PersonPrism

The `PersonPrism` converts between a builder and a validated `Person`:

```go
PersonPrism = prism.MakePrismWithName(
    buildPerson(),      // Forward: builder -> Option[*Person]
    buildEndomorphism(), // Reverse: *Person -> builder
    "Person",
)
```

**Forward direction** ➡️ (`buildPerson`):
1. Applies the builder to an empty `PartialPerson`
2. Validates each field using field prisms
3. Returns `Some(*Person)` if all validations pass, `None` otherwise

**Reverse direction** ⬅️ (`buildEndomorphism`):
1. Extracts validated fields from `Person`
2. Converts them back to raw types
3. Returns a builder that reconstructs the `PartialPerson`

### 💡 Usage Example

```go
// Create a builder
builder := MakePerson("Alice", 25)

// Validate and convert to Person
maybePerson := PersonPrism.GetOption(builder)

// maybePerson is Option[*Person]
// - Some(*Person) if validation succeeds ✅
// - None if validation fails (no error details) ❌
```

## 📝 Approach 2: Codec-Based Validation (With Error Messages)

### Creating Field Codecs

Convert prisms to codecs for detailed validation:

```go
nameCodec = codec.FromRefinement(namePrism)
ageCodec = codec.FromRefinement(agePrism)
```

### 🎯 Creating the PersonCodec

The `PersonCodec` provides bidirectional transformation with validation:

```go
func makePersonCodec() PersonCodec {
    return codec.MakeType(
        "Person",
        codec.Is[*Person](),
        makePersonValidate(),  // Validation with error reporting
        buildEndomorphism(),   // Encoding (same as prism)
    )
}
```

The `makePersonValidate` function:
1. Applies the builder to an empty `PartialPerson`
2. Validates each field using field codecs
3. Accumulates validation errors using applicative composition 📚
4. Returns `Validation[*Person]` (either errors or a valid `Person`)

### 💡 Usage Example

```go
// Create a builder
builder := MakePerson("", 15) // Invalid: empty name, age < 18

// Validate with detailed errors
personCodec := makePersonCodec()
validation := personCodec.Validate(builder)

// validation is Validation[*Person]
// - Right(*Person) if validation succeeds ✅
// - Left(ValidationErrors) with detailed error messages if validation fails ❌
```

## ⚖️ Key Differences

| Feature | Prism-Based 🔍 | Codec-Based 📝 |
|---------|-------------|-------------|
| Error Messages | No (returns `None`) ❌ | Yes (returns detailed errors) ✅ |
| Complexity | Simpler 🟢 | More complex 🟡 |
| Use Case | Simple validation | Production validation with user feedback |
| Return Type | `Option[*Person]` | `Validation[*Person]` |

## 📝 Pattern Summary

1. **Define types** 📐: Create `PartialPerson` (unvalidated) and `Person` (validated)
2. **Generate lenses** 🔧: Use `fp-go:Lens` directive
3. **Export setters** 🎁: Create `WithXXX` methods from lens setters
4. **Create validation prisms** 🎭: Define validation rules for each field
5. **Choose validation approach** ⚖️:
   - **Simple** 🔍: Create a `Prism` for quick validation without errors
   - **Detailed** 📝: Create a `Codec` for validation with error reporting

## ✨ Benefits

- **Type Safety** 🛡️: Validated types guarantee business rules at compile time
- **Composability** 🧩: Builders can be composed using monoid operations
- **Bidirectional** ↔️: Both prisms and codecs support encoding and decoding
- **Separation of Concerns** 🎯: Building and validation are separate phases
- **Functional** 🔄: Pure functions, no mutation, easy to test

## 📁 Files

- `types.go`: Type definitions and lens generation directives
- `builder.go`: Prism-based builder implementation
- `codec.go`: Codec-based validation implementation
- `codec_test.go`: Tests demonstrating usage patterns