---
name: fp-go-lens
description: Use this skill when working with lenses and optics in Go using the fp-go library (github.com/IBM/fp-go/v2/optics/lens). Trigger on mentions of lenses, optics, MakeLens, MakeLensRef, MakeLensStrict, lens composition, immutable updates to nested structs, accessing nested data structures, Compose, ComposeRef, ComposeOption, FromNillable, FromNillableRef, Modify, getter/setter patterns, or functional updates to Go structs. Also trigger when the user needs to update deeply nested fields immutably or work with optional fields in struct hierarchies. Also trigger for `// fp-go:Lens` annotation or go generate for lens code generation.
---

# fp-go Lenses for Structs

## Overview

Lenses are functional optics that provide immutable access to nested data structures. A [`Lens[S, A]`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#Lens) focuses on a field `A` within a structure `S`, providing `Get` (read) and `Set` (immutable update) operations.

```go
import (
    L  "github.com/IBM/fp-go/v2/optics/lens"
    LO "github.com/IBM/fp-go/v2/optics/lens/option"
    F  "github.com/IBM/fp-go/v2/function"
    O  "github.com/IBM/fp-go/v2/option"
)
```

---

## Code Generation (Recommended)

The fastest way to create lenses is via **code generation**. Annotate any struct with `// fp-go:Lens` and run `go generate ./...`.

```go
//go:generate go run github.com/IBM/fp-go/v2/main lens --dir . --filename gen_lens.go

// fp-go:Lens
type Person struct {
    Name  string
    Age   int
    Phone *string  // pointer → optional lens auto-generated
}
```

Running `go generate ./...` produces a `gen_lens.go` file with:

| Generated type | Description |
|---|---|
| `PersonLenses` | Lenses for `Person` (value type) |
| `PersonRefLenses` | Lenses for `*Person` (pointer type) |
| `PersonPrisms` | Prisms for `Person` |
| `PersonRefPrisms` | Prisms for `*Person` |

Each type has fields for **every struct field**, with both a required lens and an optional lens:

```go
type PersonLenses struct {
    Name  L.Lens[Person, string]            // required lens
    NameO LO.LensO[Person, string]          // optional lens (LensO = Lens[S, Option[A]])
    Age   L.Lens[Person, int]
    AgeO  LO.LensO[Person, int]
    Phone L.Lens[Person, *string]
    PhoneO LO.LensO[Person, *string]        // pointer fields auto-get optional variant
}
```

Constructor functions:

```go
lenses    := MakePersonLenses()             // PersonLenses
refLenses := MakePersonRefLenses()          // PersonRefLenses

// Usage
person  := Person{Name: "Alice", Age: 30}
updated := lenses.Name.Set("Bob")(person)   // Person{Name: "Bob", Age: 30}
name    := lenses.Name.Get(person)          // "Alice"
```

**Pointer fields** (`*string`, `*SomeStruct`) automatically generate optional lenses using `LO.FromNillable`.

**Embedded structs** and **generic types** are also supported — for embedded structs, lenses are generated for each promoted field.

---

## Manual Lens Creation

Use manual creation when code generation is not suitable or for one-off lenses.

### MakeLens — Value Types

Use [`MakeLens`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLens) for structs passed **by value**. The setter receives a copy automatically.

```go
type Person struct {
    Name string
    Age  int
}

nameLens := L.MakeLens(
    func(p Person) string          { return p.Name },
    func(p Person, name string) Person { p.Name = name; return p },
)

person  := Person{Name: "Alice", Age: 30}
updated := nameLens.Set("Bob")(person)   // Person{Name: "Bob", Age: 30}
name    := nameLens.Get(person)          // "Alice"
```

Method expressions also work (receiver becomes first argument):

```go
func (p Person) GetName() string            { return p.Name }
func (p Person) SetName(name string) Person { p.Name = name; return p }

nameLens := L.MakeLens(Person.GetName, Person.SetName)
```

**Setter signature**: `func(S, A) S` — struct first, value second.

### MakeLensRef — Pointer Types

Use [`MakeLensRef`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLensRef) for structs passed **by pointer**. The framework handles copying automatically — your setter modifies the pointer directly.

```go
type Address struct {
    Street string
    City   string
}

streetLens := L.MakeLensRef(
    func(a *Address) string                  { return a.Street },
    func(a *Address, s string) *Address      { a.Street = s; return a },
)

addr    := &Address{Street: "Main St", City: "Boston"}
updated := streetLens.Set("Oak Ave")(addr)   // new *Address, original unchanged
```

**Setter signature**: `func(*S, A) *S` — pointer first, value second. No manual copy needed.

### MakeLensStrict — Pointer Types with Comparable Fields (Optimization)

Use [`MakeLensStrict`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLensStrict) for pointer structs when the field type is **comparable** (`string`, `int`, pointers, etc.). If the new value equals the current value, the original pointer is returned unchanged (no allocation).

```go
nameLens := L.MakeLensStrict(
    func(p *Person) string             { return p.Name },
    func(p *Person, name string) *Person { p.Name = name; return p },
)

// Setting the same value → returns original pointer (no copy)
same    := nameLens.Set("Alice")(person)  // same == person
// Setting a different value → creates a new copy
updated := nameLens.Set("Bob")(person)
```

For non-comparable types use [`MakeLensWithEq`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLensWithEq) with a custom `Eq[A]`.

### MakeLensWithEq — Pointer Types with Custom Equality

Use [`MakeLensWithEq`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#MakeLensWithEq) when the field type is **not comparable** (slices, maps, structs containing those). Provide an `Eq[A]` to determine equality; if values are equal, the original pointer is returned unchanged without copying.

```go
import A "github.com/IBM/fp-go/v2/array"

// []string is not comparable, so MakeLensStrict cannot be used.
// array.StrictEquals[string]() compares slices element-by-element using slices.Equal.
tagsLens := L.MakeLensWithEq(
    A.StrictEquals[string](),
    func(p *Person) []string            { return p.Tags },
    func(p *Person, t []string) *Person { p.Tags = t; return p },
)
```

For custom element equality (e.g. case-insensitive strings or struct slices), use `EQ.FromEquals(func(a, b []T) bool { ... })` instead. For comparable types, `EQ.FromStrictEquals[T]()` is equivalent to using `MakeLensStrict`.

### Comparison

| Constructor | Source type | Copy responsibility | Best for |
|---|---|---|---|
| `MakeLens` | `S` (value) | Automatic (value copy) | Structs by value |
| `MakeLensRef` | `*S` (pointer) | Automatic (framework) | Structs by pointer |
| `MakeLensStrict` | `*S` (pointer) | Automatic + equality skip | Comparable fields, performance-sensitive |
| `MakeLensWithEq` | `*S` (pointer) | Automatic + custom Eq | Non-comparable fields with equality |

---

## Composing Lenses

### Compose — Value Outer Structure

[`Compose[S](ab)(sa)`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#Compose) combines a `Lens[S, A]` and a `Lens[A, B]` into a `Lens[S, B]`.

```go
type Street  struct { Name string }
type Address struct { Street Street }

addressLens := L.MakeLens(
    func(p Person) Address          { return p.Address },
    func(p Person, a Address) Person { p.Address = a; return p },
)
streetNameLens := L.MakeLens(
    func(a Address) string            { return a.Street.Name },
    func(a Address, n string) Address { a.Street.Name = n; return a },
)

// Compose: Person → Address → string
personStreetNameLens := F.Pipe1(addressLens, L.Compose[Person](streetNameLens))

updated := personStreetNameLens.Set("Oak Ave")(person)
```

### ComposeRef — Pointer Outer Structure

[`ComposeRef[S](ab)(sa)`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#ComposeRef) is the pointer version — use when the outer lens is `Lens[*S, A]`.

```go
// Composes Lens[*Person, Address] with Lens[Address, string] → Lens[*Person, string]
personStreetNameLens := F.Pipe1(addressRefLens, L.ComposeRef[Person](streetNameLens))
```

---

## Working with Optional Fields

Optional field lenses have type `LensO[S, A]` (alias for `Lens[S, Option[A]]`). Get returns `Option[A]`; Set takes `Option[A]`.

All functions below are in the **`optics/lens/option`** package (imported as `LO`).

### FromNillable — Pointer Field to Option

[`LO.FromNillable`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens/option#FromNillable) converts a `Lens[S, *A]` to a `LensO[S, *A]`. Get returns `None` when the pointer is nil.

```go
type Company struct {
    Name    string
    Address *Address  // optional
}

addressPtrLens := L.MakeLens(
    func(c Company) *Address             { return c.Address },
    func(c Company, a *Address) Company  { c.Address = a; return c },
)

// LensO[Company, *Address]
optAddressLens := LO.FromNillable(addressPtrLens)

company := Company{Name: "Acme"}
result  := optAddressLens.Get(company)              // None[*Address]

withAddr := optAddressLens.Set(O.Some(&Address{City: "Boston"}))(company)
cleared  := optAddressLens.Set(O.None[*Address]())(withAddr)
```

For pointer outer structs use [`LO.FromNillableRef`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens/option#FromNillableRef) (`Lens[*S, *A]` → `LensO[*S, *A]`).

### ComposeOption — Optional Container, Required Field

[`LO.ComposeOption[S, B](defaultA)(ab)`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens/option#ComposeOption) composes a `LensO[S, A]` (optional container) with a `Lens[A, B]` (required field) into a `LensO[S, B]`.

- **Get**: returns `None[B]` when A is absent
- **Set(Some[B])**: updates B in A, creating A from `defaultA` if absent
- **Set(None[B])**: removes A entirely

```go
type Config struct { Database *Database }

dbLens := LO.FromNillable(L.MakeLens(
    func(c Config) *Database             { return c.Database },
    func(c Config, db *Database) Config  { c.Database = db; return c },
))

portLens := L.MakeLensRef(
    func(db *Database) int               { return db.Port },
    func(db *Database, p int) *Database  { db.Port = p; return db },
)

defaultDB := &Database{Host: "localhost", Port: 5432}
configPortLens := F.Pipe1(dbLens, LO.ComposeOption[Config, int](defaultDB)(portLens))

config  := Config{}
port    := configPortLens.Get(config)                         // None[int]
updated := configPortLens.Set(O.Some(3306))(config)           // Database created from defaultDB
```

### Compose (option package) — Optional Container, Optional Field

[`LO.Compose[S, B](defaultA)(ab)`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens/option#Compose) is like `ComposeOption` but `ab` is a `LensO[A, B]` (the inner field is also optional).

```go
// Both the Settings container and MaxRetries field are optional pointers
settingsLens := LO.FromNillable(settingsPtrLens)   // LensO[Config, *Settings]
retriesLens  := LO.FromNillable(retriesPtrLens)    // LensO[*Settings, *int]

defaultSettings := &Settings{}
configRetriesLens := F.Pipe1(settingsLens, LO.Compose[Config, *int](defaultSettings)(retriesLens))
```

### Chaining Multiple Optional Levels

Use `F.Pipe2` / `F.Pipe3` to chain compose steps. Choose the compose function based on the type of the inner lens (`ab`):

| Function | `ab` type | Use when |
|---|---|---|
| `LO.ComposeOption` | `Lens[A, B]` — required field | The inner field always exists once A is present |
| `LO.Compose` | `LensO[A, B]` — optional field | The inner field is itself optional |

```go
// Person.Address (*Address, optional) → Address.Street (*Street, optional) → Street.Name (string, required)
defaultAddress := &Address{}
defaultStreet  := &Street{}

streetNameLens := L.MakeLensStrict(
    func(s *Street) string               { return s.Name },
    func(s *Street, n string) *Street    { s.Name = n; return s },
)
// LensO[*Address, *Street] — Street pointer is optional
addressStreetLens := LO.FromNillableRef(L.MakeLensRef(
    func(a *Address) *Street             { return a.Street },
    func(a *Address, s *Street) *Address { a.Street = s; return a },
))
// LensO[Person, *Address] — Address pointer is optional
personAddressLens := LO.FromNillable(L.MakeLens(
    func(p Person) *Address              { return p.Address },
    func(p Person, a *Address) Person    { p.Address = a; return p },
))

streetNameInPerson := F.Pipe2(
    personAddressLens,                                                      // LensO[Person, *Address]
    LO.Compose[Person, *Street](defaultAddress)(addressStreetLens),         // ab is LensO → LO.Compose
    LO.ComposeOption[Person, string](defaultStreet)(streetNameLens),        // ab is Lens  → LO.ComposeOption
)
// streetNameInPerson is LensO[Person, string]
```

---

## Modifying Values

[`Modify`](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens#Modify) applies a transformation function to the focused value:

```go
type Counter struct { Value int }

valueLens := L.MakeLens(
    func(c Counter) int             { return c.Value },
    func(c Counter, v int) Counter  { c.Value = v; return c },
)

counter     := Counter{Value: 5}
incremented := valueLens.Modify(N.Add(1))(counter)
// incremented.Value == 6, counter.Value == 5 (unchanged)
```

Or using the package-level function in pipelines:

```go
incremented := F.Pipe1(valueLens, L.Modify[Counter](N.Add(1)))(counter)
```

---

## Import Reference

```go
import (
    L  "github.com/IBM/fp-go/v2/optics/lens"            // MakeLens, MakeLensRef, MakeLensStrict, MakeLensWithEq, Compose, ComposeRef, Modify
    LO "github.com/IBM/fp-go/v2/optics/lens/option"     // LensO, FromNillable, FromNillableRef, ComposeOption, Compose
    F  "github.com/IBM/fp-go/v2/function"               // Pipe1, Pipe2, Pipe3, ...
    O  "github.com/IBM/fp-go/v2/option"                 // Some, None, Option
    EQ "github.com/IBM/fp-go/v2/eq"                     // FromStrictEquals, FromEquals (for MakeLensWithEq)
    A  "github.com/IBM/fp-go/v2/array"                  // StrictEquals, Eq (for MakeLensWithEq with slice fields)
    N  "github.com/IBM/fp-go/v2/number"                 // Add, Sub, Mul, ... (arithmetic on numeric fields)
)
```

---

## Complete Manual Example

```go
package main

import (
    F  "github.com/IBM/fp-go/v2/function"
    L  "github.com/IBM/fp-go/v2/optics/lens"
    LO "github.com/IBM/fp-go/v2/optics/lens/option"
    O  "github.com/IBM/fp-go/v2/option"
)

type Street struct {
    Name string
}

type Address struct {
    City   string
    Street *Street  // optional
}

type Person struct {
    Name    string
    Age     int
    Address *Address  // optional
}

func main() {
    // Leaf lens: *Street.Name (comparable → use Strict for optimization)
    streetNameLens := L.MakeLensStrict(
        func(s *Street) string              { return s.Name },
        func(s *Street, n string) *Street   { s.Name = n; return s },
    )

    // Mid lens: *Address.Street (pointer field → optional)
    addressStreetLens := LO.FromNillableRef(L.MakeLensRef(
        func(a *Address) *Street              { return a.Street },
        func(a *Address, s *Street) *Address  { a.Street = s; return a },
    ))

    // Root lens: Person.Address (pointer field → optional)
    personAddressLens := LO.FromNillable(L.MakeLens(
        func(p Person) *Address              { return p.Address },
        func(p Person, a *Address) Person    { p.Address = a; return p },
    ))

    // Compose all levels
    defaultAddress := &Address{City: "Unknown"}
    defaultStreet  := &Street{}

    streetNameInPerson := F.Pipe2(
        personAddressLens,
        LO.Compose[Person, *Street](defaultAddress)(addressStreetLens),
        LO.ComposeOption[Person, string](defaultStreet)(streetNameLens),
    )

    // Update deeply nested field immutably
    person := Person{
        Name: "Alice",
        Age:  30,
        Address: &Address{
            City:   "Boston",
            Street: &Street{Name: "Main St"},
        },
    }

    updated := streetNameInPerson.Set(O.Some("Oak Ave"))(person)
    // person.Address.Street.Name == "Main St" (original unchanged)
    // updated.Address.Street.Name == "Oak Ave"

    // Person with no address: default values are used
    noAddr   := Person{Name: "Bob"}
    withAddr := streetNameInPerson.Set(O.Some("Elm St"))(noAddr)
    // withAddr.Address == &Address{City: "Unknown", Street: &Street{Name: "Elm St"}}
}
```

---

## Further Reading

- [Introduction to optics: lenses and prisms](https://medium.com/@gcanti/introduction-to-optics-lenses-and-prisms-3230e73bfcfe)
- [Lens Laws](https://pkg.go.dev/github.com/IBM/fp-go/v2/optics/lens/testing) — Property-based tests for lens correctness
- [Code generation CLI](https://pkg.go.dev/github.com/IBM/fp-go/v2/cli) — `// fp-go:Lens` annotation details
