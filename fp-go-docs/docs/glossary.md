---
sidebar_position: 6
title: Glossary
hide_title: true
description: Functional programming terminology explained with Go examples using fp-go.
---

<PageHeader
  eyebrow="Reference · Glossary"
  title="Functional Go"
  titleAccent="glossary."
  lede="Every functional programming term you'll see in this documentation — defined with a real fp-go example."
  meta={[
    {label: '// Version', value: <>v2.2.82 <MetaPill>LATEST</MetaPill></>},
    {label: '// Terms', value: '40+'},
    {label: '// Updated', value: 'today'},
  ]}
/>

<Section id="a" number="A" title="A">

### Applicative

A type class that allows applying a function wrapped in a context to a value wrapped in a context.

<CodeCard file="applicative.go">
{`// Apply a function in Result to a value in Result
funcResult := result.Ok(func(x int) int { return x * 2 })
valueResult := result.Ok(21)
applied := result.Ap(valueResult)(funcResult) // Result[42]`}
</CodeCard>

**See also:** [Functor](#f), [Monad](#m)

### Arity

The number of arguments a function takes.

<CodeCard file="arity.go">
{`// Arity 0 (nullary)
func getValue() int { return 42 }

// Arity 1 (unary)
func double(x int) int { return x * 2 }

// Arity 2 (binary)
func add(x, y int) int { return x + y }

// Arity 3 (ternary)
func sum3(x, y, z int) int { return x + y + z }`}
</CodeCard>

</Section>

<Section id="c" number="C" title="C">

### Chain

Also known as **FlatMap** or **Bind**. Transforms a value in a context and flattens the result.

<CodeCard file="chain.go">
{`// Chain flattens nested Results
func divide(a, b int) result.Result[int] {
    if b == 0 {
        return result.Err[int](errors.New("division by zero"))
    }
    return result.Ok(a / b)
}

result := function.Pipe2(
    result.Ok(10),
    result.Chain(func(x int) result.Result[int] {
        return divide(x, 2) // Returns Result[int], not Result[Result[int]]
    }),
)`}
</CodeCard>

**See also:** [Map](#m), FlatMap

### Compose

Combines functions right-to-left (mathematical composition).

<CodeCard file="compose.go">
{`// f ∘ g = f(g(x))
composed := function.Compose2(
    func(x int) int { return x * 2 },    // Applied second
    func(x int) int { return x + 1 },    // Applied first
)
result := composed(5) // (5 + 1) * 2 = 12`}
</CodeCard>

**See also:** [Flow](#f), [Pipe](#p)

### Currying

Transforming a function with multiple arguments into a sequence of functions each taking a single argument.

<CodeCard file="currying.go">
{`// Uncurried
func add(x, y int) int {
    return x + y
}

// Curried
func addCurried(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}

add5 := addCurried(5)
result := add5(3) // 8`}
</CodeCard>

</Section>

<Section id="e" number="E" title="E">

### Effect

A computation that interacts with the outside world (I/O, network, database, etc.).

<CodeCard file="effect.go">
{`// Pure function - no effects
func add(x, y int) int {
    return x + y
}

// Effectful function - reads from disk
func readFile(path string) ioresult.IOResult[[]byte] {
    return func() result.Result[[]byte] {
        data, err := os.ReadFile(path)
        return result.FromGoError(data, err)
    }
}`}
</CodeCard>

**See also:** [IO](#i), [Pure Function](#p), [Side Effect](#s)

### Either

A type representing a value that can be one of two types: Left (typically error) or Right (typically success).

<CodeCard file="either.go">
{`// Either[error, int] - Left for errors, Right for success
func divide(a, b int) either.Either[error, int] {
    if b == 0 {
        return either.Left[int](errors.New("division by zero"))
    }
    return either.Right[error](a / b)
}`}
</CodeCard>

<Callout title="In v2.">
  Prefer <a href="#r">Result</a> for error handling.
</Callout>

**See also:** [Result](#r), [Option](#o)

### Endomorphism

A function where the input and output types are the same.

<CodeCard file="endomorphism.go">
{`// Endomorphism: int -> int
func double(x int) int {
    return x * 2
}

// Endomorphism: string -> string
func uppercase(s string) string {
    return strings.ToUpper(s)
}

// Compose endomorphisms
composed := function.Compose2(double, increment)`}
</CodeCard>

</Section>

<Section id="f" number="F" title="F">

### FlatMap

See [Chain](#c).

### Flow

Combines functions left-to-right (pipeline composition).

<CodeCard file="flow.go">
{`// g → f = f(g(x))
pipeline := function.Flow2(
    func(x int) int { return x + 1 },    // Applied first
    func(x int) int { return x * 2 },    // Applied second
)
result := pipeline(5) // (5 + 1) * 2 = 12`}
</CodeCard>

**See also:** [Pipe](#p), [Compose](#c)

### Fold

Reduces a structure to a single value by applying a function.

<CodeCard file="fold.go">
{`// Fold Result to extract value
result := divide(10, 2)
value := result.Fold(
    func(err error) int { return 0 },      // Handle error
    func(val int) int { return val },      // Handle success
)

// Fold array
sum := array.Reduce(
    func(acc, x int) int { return acc + x },
    0,
)(numbers)`}
</CodeCard>

**See also:** [Reduce](#r)

### Functor

A type that can be mapped over. Implements `Map` operation.

<CodeCard file="functor.go">
{`// Result is a Functor
result := result.Ok(21)
doubled := result.Map(func(x int) int { return x * 2 }) // Result[42]

// Array is a Functor
numbers := []int{1, 2, 3}
doubled := array.Map(func(x int) int { return x * 2 })(numbers) // [2, 4, 6]`}
</CodeCard>

**See also:** [Map](#m), [Applicative](#a), [Monad](#m)

</Section>

<Section id="h" number="H" title="H">

### Higher-Order Function

A function that takes functions as arguments or returns functions.

<CodeCard file="hof.go">
{`// Takes function as argument
func applyTwice(f func(int) int, x int) int {
    return f(f(x))
}

// Returns function
func makeAdder(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}`}
</CodeCard>

### Higher-Kinded Type (HKT)

A type that abstracts over type constructors. Go doesn't support HKT natively, but fp-go works around this limitation.

<CodeCard file="hkt.go">
{`// HKT concept: F[_] where F is a type constructor
// Result[_], Option[_], Array[_] are all type constructors

// fp-go uses workarounds for HKT-like behavior
// See advanced documentation for details`}
</CodeCard>

**See also:** [Type Constructor](#t)

</Section>

<Section id="i" number="I" title="I">

### Identity

A monad that wraps a value without adding any computational context. Mainly used for teaching monad laws.

<CodeCard file="identity.go">
{`// Identity just wraps a value
id := identity.Of(42)
doubled := identity.Map(func(x int) int { return x * 2 })(id)
value := identity.Unwrap(doubled) // 42`}
</CodeCard>

### Immutability

Data that cannot be changed after creation. All transformations create new values.

<CodeCard file="immutability.go">
{`// Immutable transformation
original := []int{1, 2, 3}
doubled := array.Map(func(x int) int { return x * 2 })(original)
// original is unchanged: [1, 2, 3]
// doubled is new: [2, 4, 6]`}
</CodeCard>

### IO

A type representing a lazy computation that performs side effects.

<CodeCard file="io.go">
{`// IO wraps side effects
func readFile(path string) io.IO[[]byte] {
    return func() []byte {
        data, _ := os.ReadFile(path)
        return data
    }
}

// Compose IO operations
program := function.Pipe2(
    readFile("config.json"),
    io.Map(parseJSON),
    io.Map(validateConfig),
)

// Execute when ready (side effect happens here)
config := program()`}
</CodeCard>

**See also:** [IOResult](#i), [Effect](#e)

### IOResult

Combines IO (lazy evaluation) with Result (error handling).

<CodeCard file="ioresult.go">
{`func readFile(path string) ioresult.IOResult[[]byte] {
    return func() result.Result[[]byte] {
        data, err := os.ReadFile(path)
        return result.FromGoError(data, err)
    }
}

// Compose with automatic error handling
program := function.Pipe2(
    readFile("config.json"),
    ioresult.Chain(parseJSON),
    ioresult.Chain(validateConfig),
)

// Execute and get Result
result := program()`}
</CodeCard>

**See also:** [IO](#i), [Result](#r)

</Section>

<Section id="l" number="L" title="L">

### Lazy Evaluation

Delaying computation until the result is needed.

<CodeCard file="lazy.go">
{`// Eager - computed immediately
result := expensiveComputation()

// Lazy - computed when called
lazyResult := lazy.Of(func() int {
    return expensiveComputation()
})

// Computation happens here
value := lazyResult()`}
</CodeCard>

**See also:** [IO](#i), Lazy Type

### Lazy Type

A type that defers computation and memoizes the result.

<CodeCard file="lazy-type.go">
{`// Create lazy value
lazyValue := lazy.Of(func() int {
    fmt.Println("Computing...")
    return 42
})

// First call - computes and caches
value1 := lazyValue() // Prints "Computing...", returns 42

// Second call - returns cached value
value2 := lazyValue() // Returns 42 (no print)`}
</CodeCard>

### Lens

An optic for focusing on a part of a data structure for immutable updates.

<CodeCard file="lens.go">
{`// Lens for accessing nested fields
type User struct {
    Name    string
    Address Address
}

type Address struct {
    Street string
    City   string
}

// Create lens for Address.City
cityLens := lens.Compose(
    userAddressLens,
    addressCityLens,
)

// Update city immutably
updated := lens.Set(cityLens, "New York")(user)
// user is unchanged, updated is new User with new city`}
</CodeCard>

**See also:** [Optics](#o), [Prism](#p)

</Section>

<Section id="m" number="M" title="M">

### Map

Transforms the value inside a context without changing the context.

<CodeCard file="map.go">
{`// Map over Result
result := result.Ok(21)
doubled := result.Map(func(x int) int { return x * 2 }) // Result[42]

// Map over Array
numbers := []int{1, 2, 3}
doubled := array.Map(func(x int) int { return x * 2 })(numbers) // [2, 4, 6]

// Map over Option
opt := option.Some(21)
doubled := option.Map(func(x int) int { return x * 2 })(opt) // Some(42)`}
</CodeCard>

**See also:** [Functor](#f), [Chain](#c)

### Monad

A type class that supports `Map` (Functor), `Ap` (Applicative), and `Chain` (Monad) operations.

<CodeCard file="monad.go">
{`// Result is a Monad
result := result.Ok(10)

// Map (Functor)
mapped := result.Map(func(x int) int { return x * 2 })

// Chain (Monad)
chained := result.Chain(func(x int) result.Result[int] {
    return divide(x, 2)
})`}
</CodeCard>

<Callout title="Monad laws.">
  <ol>
    <li><strong>Left identity:</strong> <code>return a &gt;&gt;= f</code> ≡ <code>f a</code></li>
    <li><strong>Right identity:</strong> <code>m &gt;&gt;= return</code> ≡ <code>m</code></li>
    <li><strong>Associativity:</strong> <code>(m &gt;&gt;= f) &gt;&gt;= g</code> ≡ <code>m &gt;&gt;= (\x -&gt; f x &gt;&gt;= g)</code></li>
  </ol>
</Callout>

**See also:** [Functor](#f), [Applicative](#a)

### Monoid

A type with an associative binary operation and an identity element.

<CodeCard file="monoid.go">
{`// Monoid for addition
// Identity: 0
// Operation: +
sum := monoid.Concat(
    monoid.Sum,
    []int{1, 2, 3, 4},
) // 10

// Monoid for string concatenation
// Identity: ""
// Operation: +
combined := monoid.Concat(
    monoid.String,
    []string{"Hello", " ", "World"},
) // "Hello World"`}
</CodeCard>

**See also:** [Semigroup](#s)

</Section>

<Section id="o" number="O" title="O">

### Optics

A family of composable tools for accessing and updating immutable data structures.

<Callout title="Types of optics.">
  <ul>
    <li><strong>Lens:</strong> focus on a field</li>
    <li><strong>Prism:</strong> focus on a variant of a sum type</li>
    <li><strong>Optional:</strong> focus on an optional field</li>
    <li><strong>Traversal:</strong> focus on multiple elements</li>
  </ul>
</Callout>

<CodeCard file="optics.go">
{`// Lens example
updated := lens.Set(nameLens, "Alice")(user)

// Prism example
value := prism.GetOption(rightPrism)(either)

// Traversal example
updated := traversal.Modify(arrayTraversal, double)(numbers)`}
</CodeCard>

**See also:** [Lens](#l), [Prism](#p)

### Option

A type representing an optional value: `Some(value)` or `None`.

<CodeCard file="option.go">
{`// Some - has a value
some := option.Some(42)

// None - no value
none := option.None[int]()

// Safe operations
value := option.GetOrElse(func() int { return 0 })(some) // 42
value = option.GetOrElse(func() int { return 0 })(none)  // 0`}
</CodeCard>

**See also:** [Result](#r), [Either](#e)

</Section>

<Section id="p" number="P" title="P">

### Partial Application

Fixing some arguments of a function to create a new function with fewer arguments.

<CodeCard file="partial.go">
{`// Original function
func add(x, y int) int {
    return x + y
}

// Partially applied
add5 := func(y int) int {
    return add(5, y)
}

result := add5(3) // 8`}
</CodeCard>

**See also:** [Currying](#c)

### Pipe

Applies a value to a sequence of functions left-to-right.

<CodeCard file="pipe.go">
{`// x |> f |> g = g(f(x))
result := function.Pipe3(
    10,                                    // Start with 10
    func(x int) int { return x + 1 },     // 11
    func(x int) int { return x * 2 },     // 22
    func(x int) int { return x - 2 },     // 20
)`}
</CodeCard>

**See also:** [Flow](#f), [Compose](#c)

### Predicate

A function that returns a boolean.

<CodeCard file="predicate.go">
{`// Predicate: int -> bool
isPositive := func(x int) bool {
    return x > 0
}

// Compose predicates
isEvenAndPositive := predicate.And(isEven, isPositive)

// Use in filter
filtered := array.Filter(isPositive)(numbers)`}
</CodeCard>

### Prism

An optic for focusing on a variant of a sum type.

<CodeCard file="prism.go">
{`// Prism for Either
rightPrism := prism.Right[error, int]()

// Get value if Right
value := prism.GetOption(rightPrism)(either.Right[error](42))
// Some(42)

// Returns None if Left
value = prism.GetOption(rightPrism)(either.Left[int](err))
// None`}
</CodeCard>

**See also:** [Optics](#o), [Lens](#l)

### Pure Function

A function that:
1. Always returns the same output for the same input
2. Has no side effects

<CodeCard file="pure.go">
{`// Pure - deterministic, no side effects
func add(x, y int) int {
    return x + y
}

// Impure - depends on external state
var counter int
func increment() int {
    counter++
    return counter
}

// Impure - has side effects
func logAndAdd(x, y int) int {
    fmt.Println("Adding", x, y) // Side effect!
    return x + y
}`}
</CodeCard>

**See also:** [Side Effect](#s), [Referential Transparency](#r)

</Section>

<Section id="r" number="R" title="R">

### Reader

A monad for dependency injection. Represents a computation that depends on an environment.

<CodeCard file="reader.go">
{`// Reader[Config, User] - computation that needs Config to produce User
func getUser(id string) reader.Reader[Config, User] {
    return func(config Config) User {
        return config.DB.QueryUser(id)
    }
}

// Compose readers
program := function.Pipe2(
    getUser("123"),
    reader.Map(enrichUser),
    reader.Chain(validateUser),
)

// Provide environment and run
user := program(config)`}
</CodeCard>

**See also:** ReaderIOResult

### ReaderIOResult

Combines Reader (dependency injection), IO (lazy evaluation), and Result (error handling).

<CodeCard file="readerioresult.go">
{`func fetchUser(id string) readerioresult.ReaderIOResult[Database, error, User] {
    return func(db Database) ioresult.IOResult[User] {
        return func() result.Result[User] {
            user, err := db.Query(id)
            return result.FromGoError(user, err)
        }
    }
}

// Compose with automatic error handling and dependency injection
program := function.Pipe2(
    fetchUser("123"),
    readerioresult.Chain(enrichUser),
    readerioresult.Chain(validateUser),
)

// Provide database and execute
result := program(db)()`}
</CodeCard>

**See also:** [Reader](#r), [IOResult](#i)

### Reduce

Combines elements of a collection into a single value.

<CodeCard file="reduce.go">
{`// Sum array
sum := array.Reduce(
    func(acc, x int) int { return acc + x },
    0,  // Initial value
)(numbers)

// Product array
product := array.Reduce(
    func(acc, x int) int { return acc * x },
    1,  // Initial value
)(numbers)`}
</CodeCard>

**See also:** [Fold](#f)

### Referential Transparency

An expression is referentially transparent if it can be replaced with its value without changing program behavior.

<CodeCard file="ref-transparency.go">
{`// Referentially transparent
func add(x, y int) int {
    return x + y
}
// add(2, 3) can always be replaced with 5

// Not referentially transparent
func random() int {
    return rand.Int()
}
// random() cannot be replaced with a fixed value`}
</CodeCard>

**See also:** [Pure Function](#p)

### Result

A type representing success (`Ok`) or failure (`Err`). Recommended for error handling in v2.

<CodeCard file="result.go">
{`// Ok - success
success := result.Ok(42)

// Err - failure
failure := result.Err[int](errors.New("error"))

// Safe operations
value := result.GetOrElse(func() int { return 0 })(success) // 42
value = result.GetOrElse(func() int { return 0 })(failure)  // 0`}
</CodeCard>

**See also:** [Either](#e), [Option](#o)

</Section>

<Section id="s" number="S" title="S">

### Semigroup

A type with an associative binary operation.

<CodeCard file="semigroup.go">
{`// Semigroup for addition
sum := semigroup.Concat(
    semigroup.Sum,
    1, 2, 3,
) // 6

// Semigroup for string concatenation
combined := semigroup.Concat(
    semigroup.String,
    "Hello", " ", "World",
) // "Hello World"`}
</CodeCard>

**See also:** [Monoid](#m)

### Side Effect

An observable interaction with the outside world.

<Callout title="Examples of side effects.">
  Reading/writing files. Network requests. Database queries. Printing to console. Modifying global state. Random number generation.
</Callout>

<CodeCard file="side-effect.go">
{`// Has side effects
func saveUser(user User) error {
    return db.Save(user) // Side effect: database write
}

// Wrap in IO to make explicit
func saveUser(user User) ioresult.IOResult[Unit] {
    return func() result.Result[Unit] {
        err := db.Save(user)
        return result.FromGoError(unit.Unit, err)
    }
}`}
</CodeCard>

**See also:** [Pure Function](#p), [IO](#i)

### State

A monad for stateful computations.

<CodeCard file="state.go">
{`// State[S, A] - computation that transforms state S and produces value A
func increment() state.State[int, int] {
    return func(s int) (int, int) {
        newState := s + 1
        return newState, newState // (new state, value)
    }
}

// Compose stateful computations
program := function.Pipe2(
    increment(),
    state.Chain(func(x int) state.State[int, int] {
        return increment()
    }),
)

// Run with initial state
finalState, value := program(0) // (2, 2)`}
</CodeCard>

</Section>

<Section id="t" number="T" title="T">

### Traverse

Transforms a collection of values in a context, collecting the results.

<CodeCard file="traverse.go">
{`// Traverse array with Result
results := array.TraverseResult(func(x int) result.Result[int] {
    if x < 0 {
        return result.Err[int](errors.New("negative"))
    }
    return result.Ok(x * 2)
})(numbers)
// Returns Result[[]int] - all or nothing`}
</CodeCard>

**See also:** Sequence

### Type Constructor

A type that takes type parameters to produce concrete types.

<CodeCard file="type-constructor.go">
{`// Result is a type constructor
// Result[_] takes a type parameter
type Result[A any] interface {
    // ...
}

// Concrete types:
// Result[int]
// Result[string]
// Result[User]`}
</CodeCard>

**See also:** [Higher-Kinded Type](#h)

</Section>

<Section id="u" number="U" title="U">

### Unit

A type with only one value, representing "no meaningful value".

<CodeCard file="unit.go">
{`// Unit type
type Unit struct{}

// Used when you need a value but don't care what it is
func doSomething() result.Result[Unit] {
    // Do something...
    return result.Ok(Unit{})
}`}
</CodeCard>

</Section>

<Section id="see-also" number="∞" title="See also">

<ul>
  <li><a href="./concepts">Core Concepts</a></li>
  <li><a href="./concepts/monads">Monads Explained</a></li>
  <li><a href="./concepts/pure-functions">Pure Functions</a></li>
  <li><a href="./concepts/composition">Composition</a></li>
  <li><a href="https://pkg.go.dev/github.com/IBM/fp-go/v2">API Reference</a></li>
</ul>

</Section>
