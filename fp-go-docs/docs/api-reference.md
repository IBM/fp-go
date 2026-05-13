---
title: API Reference
hide_title: true
description: Complete API reference for fp-go v2 - core types, collections, utilities, and advanced patterns.
sidebar_position: 1
---

<PageHeader
  eyebrow="Reference"
  title="API"
  titleAccent="Reference"
  lede="Complete API reference for fp-go v2 covering core types, collections, function utilities, and advanced patterns."
  meta={[
    { label: 'Version', value: 'v2.x' },
    { label: 'Modules', value: '54+' }
  ]}
/>

<Section id="core-types" number="01" title="Core" titleAccent="Types">

The fundamental building blocks for functional programming in Go.

### Essential Types

<ApiTable>
| Symbol | Description |
|--------|-------------|
| **[Option](./option)** | Represent optional values without null pointers |
| **[Result](./v2/result)** | Type-safe error handling with Ok/Err variants |
| **[Either](./v2/either)** | Represent a value of one of two possible types (Left/Right) |
| **[IO](./v2/io)** | Lazy side effects that execute when called |
| **[IOResult](./v2/ioresult)** | Combine IO with Result for effectful error handling |
| **[IOEither](./v2/ioeither)** | Combine IO with Either for lazy error handling |
| **[IOOption](./v2/iooption)** | Combine IO with Option for optional effects |
</ApiTable>

### Reader Types

Dependency injection and environment management:

<ApiTable>
| Symbol | Description |
|--------|-------------|
| **[Reader](./v2/reader)** | Computations that depend on a shared environment |
| **[ReaderEither](./v2/readereither)** | Reader with Either for error handling |
| **[ReaderIO](./v2/readerio)** | Reader with IO for side effects |
| **[ReaderIOEither](./v2/readerioeither)** | Reader with IO and Either (most powerful combination) |
| **[ReaderIOResult](./v2/readerioresult)** | Reader with IO and Result |
| **[ReaderOption](./v2/readeroption)** | Reader with Option for optional dependencies |
</ApiTable>

### State & Advanced Types

Stateful computations and advanced patterns:

<ApiTable>
| Symbol | Description |
|--------|-------------|
| **[State](./v2/state)** | Stateful computations with get/put operations |
| **[StateReaderIOEither](./v2/statereaderioeither)** | State with Reader, IO, and Either (ultimate monad stack) |
| **[Lazy](./v2/lazy)** | Lazy evaluation with memoization |
| **[Identity](./v2/identity)** | Simplest monad wrapper |
| **[Constant](./v2/constant)** | Constant functor |
| **[Endomorphism](./v2/endomorphism)** | Functions from type to itself |
</ApiTable>

<CodeCard file="core_example.go">
{`package main

import (
    "fmt"
    O "github.com/IBM/fp-go/option"
    E "github.com/IBM/fp-go/either"
    R "github.com/IBM/fp-go/result"
)

func main() {
    // Option for optional values
    opt := O.Some(42)
    fmt.Println(O.IsSome(opt)) // true
    
    // Either for errors
    either := E.Right[error](42)
    fmt.Println(E.IsRight(either)) // true
    
    // Result for Ok/Err
    result := R.Ok[error](42)
    fmt.Println(R.IsOk(result)) // true
}
`}
</CodeCard>

</Section>

<Section id="collections" number="02" title="Collections" titleAccent="& Arrays">

Functional operations on slices and maps.

### Array Operations

<ApiTable>
| Symbol | Description |
|--------|-------------|
| **[Array](./v2/collections/array)** | Core array operations (Map, Filter, Reduce) |
| **[Array.Ap](./v2/collections/array-ap)** | Applicative operations for arrays |
| **[Array.Eq](./v2/collections/array-eq)** | Equality checking for arrays |
| **[Array.Find](./v2/collections/array-find)** | Search and lookup operations |
| **[Array.Monoid](./v2/collections/array-monoid)** | Monoid instance for arrays |
| **[Array.Sort](./v2/collections/array-sort)** | Sorting with custom comparators |
| **[Array.Uniq](./v2/collections/array-uniq)** | Remove duplicates |
| **[Array.Zip](./v2/collections/array-zip)** | Combine multiple arrays |
| **[NonEmptyArray](./v2/collections/nonempty-array)** | Arrays guaranteed to have at least one element |
</ApiTable>

### Record Operations

<ApiTable>
| Symbol | Description |
|--------|-------------|
| **[Record](./v2/collections/record)** | Core map operations (Map, Filter, Keys, Values) |
| **[Record.Ap](./v2/collections/record-ap)** | Applicative operations for records |
| **[Record.Chain](./v2/collections/record-chain)** | Monadic chaining for records |
| **[Record.Conversion](./v2/collections/record-conversion)** | Convert between records and arrays |
| **[Record.Eq](./v2/collections/record-eq)** | Equality checking for records |
| **[Record.Monoid](./v2/collections/record-monoid)** | Monoid instance for records |
| **[Record.Ord](./v2/collections/record-ord)** | Ordering for records |
| **[Record.Traverse](./v2/collections/record-traverse)** | Traverse records with effects |
| **[Sequence.Traverse](./v2/collections/sequence-traverse)** | Generic sequence traversal |
</ApiTable>

<CodeCard file="collections_example.go">
{`package main

import (
    "fmt"
    A "github.com/IBM/fp-go/array"
    R "github.com/IBM/fp-go/record"
)

func main() {
    // Array operations
    numbers := []int{1, 2, 3, 4, 5}
    doubled := A.Map(func(n int) int { return n * 2 })(numbers)
    fmt.Println(doubled) // [2 4 6 8 10]
    
    // Record operations
    ages := map[string]int{"Alice": 30, "Bob": 25}
    names := R.Keys(ages)
    fmt.Println(names) // [Alice Bob]
}
`}
</CodeCard>

</Section>

<Section id="utilities" number="03" title="Function" titleAccent="Utilities">

Core function manipulation and composition utilities.

### Composition & Flow

<ApiTable>
| Symbol | Description |
|--------|-------------|
| **[Pipe & Flow](./v2/utilities/pipe-flow)** | Left-to-right function composition |
| **[Compose](./v2/utilities/compose)** | Right-to-left function composition |
| **[Function](./v2/utilities/function)** | Core function utilities (Identity, Constant, Flip) |
</ApiTable>

### Currying & Binding

<ApiTable>
| Symbol | Description |
|--------|-------------|
| **[Bind & Curry](./v2/utilities/bind-curry)** | Currying and partial application |
</ApiTable>

### Type Classes

<ApiTable>
| Symbol | Description |
|--------|-------------|
| **[Eq](./v2/utilities/eq)** | Equality type class |
| **[Ord](./v2/utilities/ord)** | Ordering type class |
| **[Semigroup](./v2/utilities/semigroup)** | Associative binary operation |
| **[Monoid](./v2/utilities/monoid)** | Semigroup with identity element |
| **[Magma](./v2/utilities/magma)** | Binary operation without laws |
</ApiTable>

### Primitive Utilities

<ApiTable>
| Symbol | Description |
|--------|-------------|
| **[Boolean](./v2/utilities/boolean)** | Boolean operations and combinators |
| **[Number](./v2/utilities/number)** | Numeric operations and instances |
| **[String](./v2/utilities/string)** | String operations and instances |
| **[Predicate](./v2/utilities/predicate)** | Predicate combinators (And, Or, Not) |
| **[Tuple](./v2/utilities/tuple)** | Tuple operations and utilities |
</ApiTable>

<CodeCard file="utilities_example.go">
{`package main

import (
    "fmt"
    F "github.com/IBM/fp-go/function"
)

func main() {
    // Pipe: left-to-right composition
    add10 := func(n int) int { return n + 10 }
    double := func(n int) int { return n * 2 }
    
    result := F.Pipe2(
        5,
        add10,   // 5 + 10 = 15
        double,  // 15 * 2 = 30
    )
    fmt.Println(result) // 30
}
`}
</CodeCard>

</Section>

<Section id="type-hierarchy" number="04" title="Type" titleAccent="Hierarchy">

Understanding the relationships between types:

<CodeCard file="hierarchy.txt">
{`Functor
  ├─ Applicative
  │   └─ Monad
  │       ├─ Option
  │       ├─ Either
  │       ├─ Result
  │       ├─ IO
  │       ├─ IOEither
  │       ├─ IOResult
  │       ├─ IOOption
  │       ├─ Reader
  │       ├─ ReaderEither
  │       ├─ ReaderIO
  │       ├─ ReaderIOEither
  │       ├─ ReaderIOResult
  │       ├─ ReaderOption
  │       ├─ State
  │       └─ StateReaderIOEither
  └─ Array
`}
</CodeCard>

**Key Concepts:**

- **Functor** - Types that can be mapped over (`Map`)
- **Applicative** - Functors with `Of` (pure) and `Ap` (apply)
- **Monad** - Applicatives with `Chain` (flatMap/bind)
- **Array** - Functor but not a monad (multiple values)

</Section>

<Section id="common-patterns" number="05" title="Common" titleAccent="Patterns">

Frequently used patterns and idioms.

### Error Handling

<Compare>
<CompareCol kind="bad">
<CodeCard file="traditional.go">
{`// ❌ Traditional Go error handling
func getUser(id string) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, err
    }
    return user, nil
}
`}
</CodeCard>
</CompareCol>

<CompareCol kind="good">
<CodeCard file="functional.go">
{`// ✅ Functional error handling
func getUser(id string) either.Either[error, User] {
    return db.FindUser(id)
}

// Chain operations
result := pipe.Pipe2(
    getUser("123"),
    E.Chain(validateUser),
    E.Map(enrichUser),
)
`}
</CodeCard>
</CompareCol>
</Compare>

### Side Effects

<CodeCard file="effects.go">
{`// IO for lazy effects
io := IO.Of(func() int {
    fmt.Println("Computing...")
    return 42
})

// Execute when ready
result := io() // Prints "Computing..." and returns 42

// IOEither for effects with errors
ioe := IOE.TryCatch(func() (int, error) {
    return readFile("config.json")
})
`}
</CodeCard>

### Dependency Injection

<CodeCard file="reader.go">
{`type Config struct {
    DBUrl string
    Port  int
}

// Reader for dependencies
getDBUrl := R.Asks(func(c Config) string {
    return c.DBUrl
})

// ReaderIOEither for full power
program := RIOE.Chain(
    getDBUrl,
    func(url string) RIOE.ReaderIOEither[Config, error, *DB] {
        return RIOE.FromIO(connectDB(url))
    },
)

// Run with config
config := Config{DBUrl: "localhost:5432", Port: 8080}
result := program(config)()
`}
</CodeCard>

</Section>

<Section id="quick-links" number="06" title="Quick" titleAccent="Links">

### Getting Started
- [Introduction](./intro) - What is fp-go?
- [Installation](./installation) - Get up and running
- [Quick Start](./quickstart) - Your first functional program

### Learning Resources
- [Core Concepts](./concepts/index) - Fundamental FP concepts
- [Recipes](./recipes/index) - Practical examples
- [Migration Guide](./migration/v1-to-v2) - Upgrade from v1

### Advanced Topics
- [Patterns](./advanced/patterns) - Monad transformers, Free monads, Tagless final
- [Type Theory](./advanced/type-theory) - Category theory foundations
- [Performance](./advanced/performance) - Optimization techniques
- [Architecture](./advanced/architecture) - Application design patterns

### Legacy Version
- [v1.x Documentation](./1.0.0/intro) - Previous version docs

</Section>
