# Design Decisions

This document explains the key design decisions and principles behind fp-go's API design.

## Table of Contents

- [Data Last Principle](#data-last-principle)
- [Kleisli and Operator Types](#kleisli-and-operator-types)
- [Monadic Operations Comparison](#monadic-operations-comparison)
- [Type Parameter Ordering](#type-parameter-ordering)
- [Generic Type Aliases](#generic-type-aliases)

## Data Last Principle

fp-go follows the **"data last"** principle, where the data being operated on is always the last parameter in a function. This design choice enables powerful function composition and partial application patterns.

This principle is deeply rooted in functional programming tradition, particularly in **Haskell's design philosophy**. Haskell functions are automatically curried and follow the data-last convention, making function composition natural and elegant. For example, Haskell's `map` function has the signature `(a -> b) -> [a] -> [b]`, where the transformation function comes before the list.

### What is "Data Last"?

In the "data last" style, functions are structured so that:
1. Configuration parameters come first
2. The data to be transformed comes last

This is the opposite of the traditional object-oriented style where the data (receiver) comes first.

### Why "Data Last"?

The "data last" principle enables:

1. **Natural Currying**: Functions can be partially applied to create specialized transformations
2. **Function Composition**: Operations can be composed before applying them to data
3. **Point-Free Style**: Write transformations without explicitly mentioning the data
4. **Reusability**: Create reusable transformation pipelines

This design aligns with Haskell's approach where all functions are curried by default, enabling elegant composition patterns that have proven effective over decades of functional programming practice.

### Examples

#### Basic Transformation

```go
// Data last style (fp-go)
double := array.Map(number.Mul(2))
result := double([]int{1, 2, 3}) // [2, 4, 6]

// Compare with data first style (traditional)
result := array.Map([]int{1, 2, 3}, number.Mul(2))
```

#### Function Composition

```go
import (
    A "github.com/IBM/fp-go/v2/array"
    F "github.com/IBM/fp-go/v2/function"
    N "github.com/IBM/fp-go/v2/number"
)

// Create a pipeline of transformations
pipeline := F.Flow3(
    A.Filter(N.MoreThan(0)),     // Keep positive numbers
    A.Map(N.Mul(2)),                                  // Double each number
    A.Reduce(func(acc, x int) int { return acc + x }, 0), // Sum them up
)

// Apply the pipeline to different data
result1 := pipeline([]int{-1, 2, 3, -4, 5})  // (2 + 3 + 5) * 2 = 20
result2 := pipeline([]int{1, 2, 3})          // (1 + 2 + 3) * 2 = 12
```

#### Partial Application

```go
import (
    O "github.com/IBM/fp-go/v2/option"
)

// Create specialized functions by partial application
getOrZero := O.GetOrElse(func() int { return 0 })
getOrEmpty := O.GetOrElse(func() string { return "" })

// Use them with different data
value1 := getOrZero(O.Some(42))        // 42
value2 := getOrZero(O.None[int]())     // 0

text1 := getOrEmpty(O.Some("hello"))   // "hello"
text2 := getOrEmpty(O.None[string]())  // ""
```

#### Building Reusable Transformations

```go
import (
    E "github.com/IBM/fp-go/v2/either"
    O "github.com/IBM/fp-go/v2/option"
)

// Create a reusable validation pipeline
type User struct {
    Name  string
    Email string
    Age   int
}

validateAge := E.FromPredicate(
    func(u User) bool { return u.Age >= 18 },
    func(u User) error { return errors.New("must be 18 or older") },
)

validateEmail := E.FromPredicate(
    func(u User) bool { return strings.Contains(u.Email, "@") },
    func(u User) error { return errors.New("invalid email") },
)

// Compose validators
validateUser := F.Flow2(
    validateAge,
    E.Chain(validateEmail),
)

// Apply to different users
result1 := validateUser(User{Name: "Alice", Email: "alice@example.com", Age: 25})
result2 := validateUser(User{Name: "Bob", Email: "invalid", Age: 30})
```

#### Monadic Operations

```go
import (
    O "github.com/IBM/fp-go/v2/option"
)

// Data last enables clean monadic chains
parseAndDouble := F.Flow2(
    O.FromPredicate(func(s string) bool { return s != "" }),
    O.Chain(func(s string) O.Option[int] {
        n, err := strconv.Atoi(s)
        if err != nil {
            return O.None[int]()
        }
        return O.Some(n * 2)
    }),
)

result1 := parseAndDouble("21")  // Some(42)
result2 := parseAndDouble("")    // None
result3 := parseAndDouble("abc") // None
```

### Monadic vs Non-Monadic Forms

fp-go provides two forms for most operations:

1. **Curried form** (data last): Returns a function that can be composed
2. **Monadic form** (data first): Takes all parameters at once

```go
// Curried form - data last, returns a function
Map[A, B any](f func(A) B) func(Option[A]) Option[B]

// Monadic form - data first, direct execution
MonadMap[A, B any](fa Option[A], f func(A) B) Option[B]
```

**When to use each:**

- **Curried form**: When building pipelines, composing functions, or creating reusable transformations
- **Monadic form**: When you have all parameters available and want direct execution

```go
// Curried form - building a pipeline
transform := F.Flow3(
    O.Map(strings.ToUpper),
    O.Filter(func(s string) bool { return len(s) > 3 }),
    O.GetOrElse(func() string { return "DEFAULT" }),
)
result := transform(O.Some("hello"))

// Monadic form - direct execution
result := O.MonadMap(O.Some("hello"), strings.ToUpper)
```

### Further Reading on Data-Last Pattern

The data-last currying pattern is well-documented in the functional programming community:

#### Haskell Design Philosophy
- [Haskell Wiki - Currying](https://wiki.haskell.org/Currying) - Comprehensive explanation of currying in Haskell
- [Learn You a Haskell - Higher Order Functions](http://learnyouahaskell.com/higher-order-functions) - Introduction to currying and partial application
- [Haskell's Prelude](https://hackage.haskell.org/package/base/docs/Prelude.html) - Standard library showing data-last convention throughout

#### General Functional Programming
- [Mostly Adequate Guide - Ch. 4: Currying](https://mostly-adequate.gitbook.io/mostly-adequate-guide/ch04) - Excellent introduction with clear examples
- [Curry and Function Composition](https://medium.com/javascript-scene/curry-and-function-composition-2c208d774983) by Eric Elliott
- [Why Curry Helps](https://hughfdjackson.com/javascript/why-curry-helps/) - Practical benefits of currying

#### Related Libraries
- [fp-ts Documentation](https://gcanti.github.io/fp-ts/) - TypeScript library that inspired fp-go's design
- [fp-ts Issue #1238](https://github.com/gcanti/fp-ts/issues/1238) - Real-world examples of data-last refactoring

## Kleisli and Operator Types

fp-go uses consistent type aliases across all monads to make code more recognizable and composable. These types provide a common vocabulary that works across different monadic contexts.

### Type Definitions

```go
// Kleisli arrow - a function that returns a monadic value
type Kleisli[A, B any] = func(A) M[B]

// Operator - a function that transforms a monadic value
type Operator[A, B any] = func(M[A]) M[B]
```

Where `M` represents the specific monad (Option, Either, IO, etc.).

### Why These Types Matter

1. **Consistency**: The same type names appear across all monads
2. **Recognizability**: Experienced functional programmers immediately understand the intent
3. **Composability**: Functions with these types compose naturally
4. **Documentation**: Type signatures clearly communicate the operation's behavior

### Examples Across Monads

#### Option Monad

```go
// option/option.go
type Kleisli[A, B any] = func(A) Option[B]
type Operator[A, B any] = func(Option[A]) Option[B]

// Chain uses Kleisli
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]

// Map returns an Operator
func Map[A, B any](f func(A) B) Operator[A, B]
```

#### Either Monad

```go
// either/either.go
type Kleisli[E, A, B any] = func(A) Either[E, B]
type Operator[E, A, B any] = func(Either[E, A]) Either[E, B]

// Chain uses Kleisli
func Chain[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, B]

// Map returns an Operator
func Map[E, A, B any](f func(A) B) Operator[E, A, B]
```

#### IO Monad

```go
// io/io.go
type Kleisli[A, B any] = func(A) IO[B]
type Operator[A, B any] = func(IO[A]) IO[B]

// Chain uses Kleisli
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]

// Map returns an Operator
func Map[A, B any](f func(A) B) Operator[A, B]
```

#### Array (List Monad)

```go
// array/array.go
type Kleisli[A, B any] = func(A) []B
type Operator[A, B any] = func([]A) []B

// Chain uses Kleisli
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]

// Map returns an Operator
func Map[A, B any](f func(A) B) Operator[A, B]
```

### Pattern Recognition

Once you learn these patterns in one monad, you can apply them to all monads:

```go
// The pattern is always the same, just the monad changes

// Option
validateAge := option.Chain(func(user User) option.Option[User] {
    if user.Age >= 18 {
        return option.Some(user)
    }
    return option.None[User]()
})

// Either
validateAge := either.Chain(func(user User) either.Either[error, User] {
    if user.Age >= 18 {
        return either.Right[error](user)
    }
    return either.Left[User](errors.New("too young"))
})

// IO
validateAge := io.Chain(func(user User) io.IO[User] {
    return io.Of(user) // Always succeeds in IO
})

// Array
validateAge := array.Chain(func(user User) []User {
    if user.Age >= 18 {
        return []User{user}
    }
    return []User{} // Empty array = failure
})
```

### Composing Kleisli Arrows

Kleisli arrows compose naturally using monadic composition:

```go
import (
    O "github.com/IBM/fp-go/v2/option"
    F "github.com/IBM/fp-go/v2/function"
)

// Define Kleisli arrows
parseAge := func(s string) O.Option[int] {
    n, err := strconv.Atoi(s)
    if err != nil {
        return O.None[int]()
    }
    return O.Some(n)
}

validateAge := func(age int) O.Option[int] {
    if age >= 18 {
        return O.Some(age)
    }
    return O.None[int]()
}

formatAge := func(age int) O.Option[string] {
    return O.Some(fmt.Sprintf("Age: %d", age))
}

// Compose them using Flow and Chain
pipeline := F.Flow3(
    parseAge,
    O.Chain(validateAge),
    O.Chain(formatAge),
)

result := pipeline("25") // Some("Age: 25")
result := pipeline("15") // None (too young)
result := pipeline("abc") // None (parse error)
```

### Building Reusable Operators

Operators can be created once and reused across your codebase:

```go
import (
    E "github.com/IBM/fp-go/v2/either"
)

// Create reusable operators
type ValidationError struct {
    Field   string
    Message string
}

// Reusable validation operators
validateNonEmpty := E.Chain(func(s string) E.Either[ValidationError, string] {
    if s == "" {
        return E.Left[string](ValidationError{
            Field:   "input",
            Message: "cannot be empty",
        })
    }
    return E.Right[ValidationError](s)
})

validateEmail := E.Chain(func(s string) E.Either[ValidationError, string] {
    if !strings.Contains(s, "@") {
        return E.Left[string](ValidationError{
            Field:   "email",
            Message: "invalid format",
        })
    }
    return E.Right[ValidationError](s)
})

// Compose operators
validateEmailInput := F.Flow2(
    validateNonEmpty,
    validateEmail,
)

// Use across your application
result1 := validateEmailInput(E.Right[ValidationError]("user@example.com"))
result2 := validateEmailInput(E.Right[ValidationError](""))
result3 := validateEmailInput(E.Right[ValidationError]("invalid"))
```

### Benefits of Consistent Naming

1. **Cross-monad understanding**: Learn once, apply everywhere
2. **Easier refactoring**: Changing monads requires minimal code changes
3. **Better tooling**: IDEs can provide better suggestions
4. **Team communication**: Shared vocabulary across the team
5. **Library integration**: Third-party libraries follow the same patterns

### Identity Monad - The Simplest Case

The Identity monad shows these types in their simplest form:

```go
// identity/doc.go
type Operator[A, B any] = func(A) B

// In Identity, there's no wrapping, so:
// - Kleisli[A, B] is just func(A) B
// - Operator[A, B] is just func(A) B
// They're the same because Identity adds no context
```

This demonstrates that these type aliases represent fundamental functional programming concepts, not just arbitrary naming conventions.


## Monadic Operations Comparison

fp-go's monadic operations are inspired by functional programming languages and libraries. Here's how they compare:

| fp-go | fp-ts | Haskell | Scala | Description |
|-------|-------|---------|-------|-------------|
| `Map` | `map` | `fmap` | `map` | Functor mapping - transforms the value inside a context |
| `Chain` | `chain` | `>>=` (bind) | `flatMap` | Monadic bind - chains computations that return wrapped values |
| `Ap` | `ap` | `<*>` | `ap` | Applicative apply - applies a wrapped function to a wrapped value |
| `Of` | `of` | `return`/`pure` | `pure` | Lifts a pure value into a monadic context |
| `Fold` | `fold` | `either` | `fold` | Eliminates the context by providing handlers for each case |
| `Filter` | `filter` | `mfilter` | `filter` | Keeps values that satisfy a predicate |
| `Flatten` | `flatten` | `join` | `flatten` | Removes one level of nesting |
| `ChainFirst` | `chainFirst` | `>>` (then) | `tap` | Chains for side effects, keeping the original value |
| `Alt` | `alt` | `<\|>` | `orElse` | Provides an alternative value if the first fails |
| `GetOrElse` | `getOrElse` | `fromMaybe` | `getOrElse` | Extracts the value or provides a default |
| `FromPredicate` | `fromPredicate` | `guard` | `filter` | Creates a monadic value based on a predicate |
| `Sequence` | `sequence` | `sequence` | `sequence` | Transforms a collection of effects into an effect of a collection |
| `Traverse` | `traverse` | `traverse` | `traverse` | Maps and sequences in one operation |
| `Reduce` | `reduce` | `foldl` | `foldLeft` | Folds a structure from left to right |
| `ReduceRight` | `reduceRight` | `foldr` | `foldRight` | Folds a structure from right to left |

### Key Differences from Other Languages

#### Naming Conventions

- **Go conventions**: fp-go uses PascalCase for exported functions (e.g., `Map`, `Chain`) following Go's naming conventions
- **Type parameters first**: Non-inferrable type parameters come first (e.g., `Ap[B, E, A any]`)
- **Monadic prefix**: Direct execution forms use the `Monad` prefix (e.g., `MonadMap`, `MonadChain`)

#### Type System

```go
// fp-go (explicit type parameters when needed)
result := option.Map(transform)(value)
result := option.Map[string, int](transform)(value) // explicit when inference fails

// Haskell (type inference)
result = fmap transform value

// Scala (type inference with method syntax)
result = value.map(transform)

// fp-ts (TypeScript type inference)
const result = pipe(value, map(transform))
```

#### Currying

```go
// fp-go - explicit currying with data last
double := array.Map(number.Mul(2))
result := double(numbers)

// Haskell - automatic currying
double = fmap (*2)
result = double numbers

// Scala - method syntax
result = numbers.map(_ * 2)
```

## Type Parameter Ordering

fp-go v2 uses a specific ordering for type parameters to maximize type inference:

### Rule: Non-Inferrable Parameters First

Type parameters that **cannot be inferred** from function arguments come first. This allows the Go compiler to infer as many types as possible.

```go
// Ap - B cannot be inferred from arguments, so it comes first
func Ap[B, E, A any](fa Either[E, A]) func(Either[E, func(A) B]) Either[E, B]

// Usage - only B needs to be specified
result := either.Ap[string](value)(funcInEither)
```

### Examples

```go
// Map - all types can be inferred from arguments
func Map[E, A, B any](f func(A) B) func(Either[E, A]) Either[E, B]
// Usage - no type parameters needed
result := either.Map(transform)(value)

// Chain - all types can be inferred
func Chain[E, A, B any](f func(A) Either[E, B]) func(Either[E, A]) Either[E, B]
// Usage - no type parameters needed
result := either.Chain(validator)(value)

// Of - E cannot be inferred, comes first
func Of[E, A any](value A) Either[E, A]
// Usage - only E needs to be specified
result := either.Of[error](42)
```

### Benefits

1. **Less verbose code**: Most operations don't require explicit type parameters
2. **Better IDE support**: Type inference provides better autocomplete
3. **Clearer intent**: Only specify types that can't be inferred

## Generic Type Aliases

fp-go v2 leverages Go 1.24's generic type aliases for cleaner type definitions:

```go
// V2 - using generic type alias (requires Go 1.24+)
type ReaderIOEither[R, E, A any] = RD.Reader[R, IOE.IOEither[E, A]]

// V1 - using type definition (Go 1.18+)
type ReaderIOEither[R, E, A any] RD.Reader[R, IOE.IOEither[E, A]]
```

### Benefits

1. **True aliases**: The type is interchangeable with its definition
2. **No namespace imports needed**: Can use types directly without package prefixes
3. **Simpler codebase**: Eliminates the need for `generic` subpackages
4. **Better composability**: Types compose more naturally

### Migration Pattern

```go
// Define project-wide aliases once
package types

import (
    "github.com/IBM/fp-go/v2/option"
    "github.com/IBM/fp-go/v2/result"
    "github.com/IBM/fp-go/v2/ioresult"
)

type Option[A any] = option.Option[A]
type Result[A any] = result.Result[A]
type IOResult[A any] = ioresult.IOResult[A]

// Use throughout your codebase
package myapp

import "myproject/types"

func process(input string) types.Result[types.Option[int]] {
    // implementation
}
```

---

For more information, see:
- [README.md](./README.md) - Overview and quick start
- [FUNCTIONAL_IO.md](./FUNCTIONAL_IO.md) - Functional I/O patterns with Context and Reader
- [IDIOMATIC_COMPARISON.md](./IDIOMATIC_COMPARISON.md) - Performance comparison between standard and idiomatic packages
- [API Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2) - Complete API reference
- [Samples](./samples/) - Practical examples