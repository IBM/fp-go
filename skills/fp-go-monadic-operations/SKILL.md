# fp-go Monadic Operations

## Overview

`fp-go` (import path `github.com/IBM/fp-go/v2`) brings type-safe functional programming to Go using generics. Every monad follows a **consistent interface**: once you know the pattern in one monad, it transfers to all others.

All functions use the **data-last** principle: the data being transformed is always the last argument, enabling partial application and pipeline composition.

## Core Types

| Type | Package | Represents |
|------|---------|------------|
| `Option[A]` | `option` | A value that may or may not be present (replaces nil) |
| `Either[E, A]` | `either` | A value that is either a left error `E` or a right success `A` |
| `Result[A]` | `result` | `Either[error, A]` — shorthand for the common case |
| `IO[A]` | `io` | A lazy computation that produces `A` (possibly with side effects) |
| `IOResult[A]` | `ioresult` | `IO[Result[A]]` — lazy computation that can fail |
| `ReaderIOResult[A]` | `context/readerioresult` | `func(context.Context) IOResult[A]` — context-aware IO with errors |
| `Effect[C, A]` | `effect` | `func(C) ReaderIOResult[A]` — typed dependency injection + IO + errors |

Idiomatic (high-performance, tuple-based) equivalents live in `idiomatic/`:
- `idiomatic/option` — `(A, bool)` tuples
- `idiomatic/result` — `(A, error)` tuples
- `idiomatic/ioresult` — `func() (A, error)`
- `idiomatic/context/readerresult` — `func(context.Context) (A, error)`

## Standard Operations

Every monad exports these operations (PascalCase for exported Go names):

| fp-go | fp-ts / Haskell | Description |
|-------|----------------|-------------|
| `Of` | `of` / `pure` | Lift a pure value into the monad |
| `Map` | `map` / `fmap` | Transform the value inside without changing the context |
| `Chain` | `chain` / `>>=` | Sequence a computation that itself returns a monadic value |
| `Ap` | `ap` / `<*>` | Apply a wrapped function to a wrapped value |
| `Fold` | `fold` / `either` | Eliminate the context — handle every case and extract a plain value |
| `GetOrElse` | `getOrElse` / `fromMaybe` | Extract the value or use a default (Option/Result) |
| `Filter` | `filter` / `mfilter` | Keep only values satisfying a predicate |
| `Flatten` | `flatten` / `join` | Remove one level of nesting (`M[M[A]]` → `M[A]`) |
| `ChainFirst` | `chainFirst` / `>>` | Sequence for side effects; keeps the original value |
| `Alt` | `alt` / `<\|>` | Provide an alternative when the first computation fails |
| `FromPredicate` | `fromPredicate` / `guard` | Build a monadic value from a predicate |
| `Sequence` | `sequence` | Turn `[]M[A]` into `M[[]A]` |
| `Traverse` | `traverse` | Map and sequence in one step |

Curried (composable) vs. monadic (direct) form:

```go
// Curried — data last, returns a transformer function
option.Map(strings.ToUpper)              // func(Option[string]) Option[string]

// Monadic — data first, immediate execution
option.MonadMap(option.Some("hello"), strings.ToUpper)
```

Use curried form for pipelines; use `Monad*` form when you already have all arguments.

## Key Type Aliases (defined per monad)

```go
// A Kleisli arrow: a function from A to a monadic B
type Kleisli[A, B any] = func(A) M[B]

// An operator: transforms one monadic value into another
type Operator[A, B any] = func(M[A]) M[B]
```

`Chain` takes a `Kleisli`, `Map` returns an `Operator`. The naming is consistent across all monads.

## Examples

### Option — nullable values without nil

```go
import (
    O "github.com/IBM/fp-go/v2/option"
    F "github.com/IBM/fp-go/v2/function"
    "strconv"
)

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

parseAndDouble("21")  // Some(42)
parseAndDouble("")    // None
parseAndDouble("abc") // None
```

### Result — error handling without if-err boilerplate

```go
import (
    R  "github.com/IBM/fp-go/v2/result"
    F  "github.com/IBM/fp-go/v2/function"
    "strconv"
    "errors"
)

parse := R.Eitherize1(strconv.Atoi)  // lifts (int, error) → Result[int]

validate := func(n int) R.Result[int] {
    if n < 0 {
        return R.Error[int](errors.New("must be non-negative"))
    }
    return R.Of(n)
}

pipeline := F.Flow2(parse, R.Chain(validate))

pipeline("42")   // Ok(42)
pipeline("-1")   // Error("must be non-negative")
pipeline("abc")  // Error(strconv parse error)
```

### IOResult — lazy IO with error handling

```go
import (
    IOE "github.com/IBM/fp-go/v2/ioresult"
    F   "github.com/IBM/fp-go/v2/function"
    J   "github.com/IBM/fp-go/v2/json"
    "os"
)

readConfig := F.Flow2(
    IOE.Eitherize1(os.ReadFile),           // func(string) IOResult[[]byte]
    IOE.ChainEitherK(J.Unmarshal[Config]), // parse JSON, propagate errors
)

result := readConfig("config.json")() // execute lazily
```

### ReaderIOResult — context-aware pipelines (recommended for services)

```go
import (
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    F   "github.com/IBM/fp-go/v2/function"
    "context"
)

// type ReaderIOResult[A any] = func(context.Context) func() result.Result[A]

fetchUser := func(id int) RIO.ReaderIOResult[User] {
    return func(ctx context.Context) func() result.Result[User] {
        return func() result.Result[User] {
            // perform IO here
        }
    }
}

pipeline := F.Pipe3(
    fetchUser(42),
    RIO.ChainEitherK(validateUser),    // lift pure (User, error) function
    RIO.Map(enrichUser),               // lift pure User → User function
    RIO.ChainFirstIOK(IO.Logf[User]("Fetched: %v")), // side-effect logging
)

user, err := pipeline(ctx)() // provide context once, execute
```

### Traversal — process slices monadically

```go
import (
    A   "github.com/IBM/fp-go/v2/array"
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    F   "github.com/IBM/fp-go/v2/function"
)

// Fetch all users, stop on first error
fetchAll := F.Pipe1(
    A.MakeBy(10, userID),
    RIO.TraverseArray(fetchUser),  // []ReaderIOResult[User] → ReaderIOResult[[]User]
)
```

## Function Composition with Flow and Pipe

```go
import F "github.com/IBM/fp-go/v2/function"

// Flow: compose functions left-to-right, returns a new function
transform := F.Flow3(
    option.Map(strings.TrimSpace),
    option.Filter(func(s string) bool { return s != "" }),
    option.GetOrElse(func() string { return "default" }),
)
result := transform(option.Some("  hello  ")) // "hello"

// Pipe: apply a value through a pipeline immediately
result := F.Pipe3(
    option.Some("  hello  "),
    option.Map(strings.TrimSpace),
    option.Filter(func(s string) bool { return s != "" }),
    option.GetOrElse(func() string { return "default" }),
)
```

## Lifting Pure Functions into Monadic Context

fp-go provides helpers to promote non-monadic functions:

| Helper | Lifts |
|--------|-------|
| `ChainEitherK` | `func(A) (B, error)` → works inside the monad |
| `ChainOptionK` | `func(A) Option[B]` → works inside the monad |
| `ChainFirstIOK` | `func(A) IO[B]` for side effects, keeps original value |
| `Eitherize1..N` | `func(A) (B, error)` → `func(A) Result[B]` |
| `FromPredicate` | `func(A) bool` + error builder → `func(A) Result[A]` |

## Type Parameter Ordering Rule (V2)

Non-inferrable type parameters come **first**, so the compiler can infer the rest:

```go
// B cannot be inferred from the argument — it comes first
result := either.Ap[string](value)(funcInEither)

// All types inferrable — no explicit params needed
result := either.Map(transform)(value)
result := either.Chain(validator)(value)
```

## When to Use Which Monad

| Situation | Use |
|-----------|-----|
| Value that might be absent | `Option[A]` |
| Operation that can fail with custom error type | `Either[E, A]` |
| Operation that can fail with `error` | `Result[A]` |
| Lazy IO, side effects | `IO[A]` |
| IO that can fail | `IOResult[A]` |
| IO + context (cancellation, deadlines) | `ReaderIOResult[A]` from `context/readerioresult` |
| IO + context + typed dependencies | `Effect[C, A]` |
| High-performance services | Idiomatic packages in `idiomatic/` |

## Do-Notation: Accumulating State with `Bind` and `ApS`

When a pipeline needs to carry **multiple intermediate results** forward — not just a single value — the `Chain`/`Map` style becomes unwieldy because each step only threads one value and prior results are lost. Do-notation solves this by accumulating results into a growing struct (the "state") at each step.

Every monad that supports do-notation exports the same family of functions. The examples below use `context/readerioresult` (`RIO`), but the identical API is available in `result`, `option`, `ioresult`, `readerioresult`, and others.

### The Function Family

| Function | Kind | What it does |
|----------|------|-------------|
| `Do(empty S)` | — | Lift an empty struct into the monad; starting point |
| `BindTo(setter)` | monadic | Convert an existing `M[T]` into `M[S]`; alternative start |
| `Bind(setter, f)` | monadic | Add a result; `f` receives the **current state** and returns `M[T]` |
| `ApS(setter, fa)` | applicative | Add a result; `fa` is **independent** of the current state |
| `Let(setter, f)` | pure | Add a value computed by a **pure function** of the state |
| `LetTo(setter, value)` | pure | Add a **constant** value |

Lens variants (`BindL`, `ApSL`, `LetL`, `LetToL`) accept a `Lens[S, T]` instead of a manual setter, integrating naturally with the optics system.

### `Bind` — Sequential, Dependent Steps

`Bind` sequences two monadic computations. The function `f` receives the **full accumulated state** so it can read anything gathered so far. Errors short-circuit automatically.

```go
import (
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    F   "github.com/IBM/fp-go/v2/function"
    L   "github.com/IBM/fp-go/v2/optics/lens"
    "context"
)

type Pipeline struct {
    User   User
    Config Config
    Posts  []Post
}

// Lenses — focus on individual fields; .Set is already func(T) func(S) S
var (
    userLens   = L.MakeLens(func(s Pipeline) User   { return s.User },   func(s Pipeline, u User)   Pipeline { s.User = u; return s })
    configLens = L.MakeLens(func(s Pipeline) Config { return s.Config }, func(s Pipeline, c Config) Pipeline { s.Config = c; return s })
    postsLens  = L.MakeLens(func(s Pipeline) []Post { return s.Posts },  func(s Pipeline, p []Post) Pipeline { s.Posts = p; return s })
)

result := F.Pipe3(
    RIO.Do(Pipeline{}),                                                   // lift empty struct
    RIO.Bind(userLens.Set,   func(_ Pipeline) RIO.ReaderIOResult[User] { return fetchUser(42) }),
    RIO.Bind(configLens.Set, F.Flow2(userLens.Get, fetchConfigForUser)), // read s.User, pass to fetcher
    RIO.Bind(postsLens.Set,  F.Flow2(userLens.Get, fetchPostsForUser)),  // read s.User, pass to fetcher
)

pipeline, err := result(context.Background())()
// pipeline.User, pipeline.Config, pipeline.Posts are all populated
```

The setter signature is `func(T) func(S1) S2` — it takes the new value and returns a state transformer. `lens.Set` already has this shape, so no manual setter functions are needed. `F.Flow2(lens.Get, f)` composes the field getter with any Kleisli arrow `f` point-free.

### `ApS` — Independent, Applicative Steps

`ApS` uses **applicative** semantics: `fa` is evaluated without any access to the current state. Use it when steps have no dependency on each other — the library can choose to execute them concurrently.

```go
import (
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    F   "github.com/IBM/fp-go/v2/function"
    L   "github.com/IBM/fp-go/v2/optics/lens"
)

type Summary struct {
    User    User
    Weather Weather
}

var (
    userLens    = L.MakeLens(func(s Summary) User    { return s.User },    func(s Summary, u User)    Summary { s.User = u; return s })
    weatherLens = L.MakeLens(func(s Summary) Weather { return s.Weather }, func(s Summary, w Weather) Summary { s.Weather = w; return s })
)

// Both are independent — neither needs the other's result
result := F.Pipe2(
    RIO.Do(Summary{}),
    RIO.ApS(userLens.Set,    fetchUser(42)),
    RIO.ApS(weatherLens.Set, fetchWeather("NYC")),
)
```

**Key difference from `Bind`:**

| | `Bind(setter, f)` | `ApS(setter, fa)` |
|-|---|---|
| Second argument | `func(S1) M[T]` — a **function** of state | `M[T]` — a **fixed** monadic value |
| Can read prior state? | Yes — receives `S1` | No — no access to state |
| Semantics | Monadic (sequential) | Applicative (independent) |

### `Let` and `LetTo` — Pure Additions

`Let` adds a value computed by a **pure function** of the current state (no monad, cannot fail):

```go
import (
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    F   "github.com/IBM/fp-go/v2/function"
    L   "github.com/IBM/fp-go/v2/optics/lens"
)

type Enriched struct {
    User     User
    FullName string
}

var (
    userLens     = L.MakeLens(func(s Enriched) User   { return s.User },     func(s Enriched, u User)   Enriched { s.User = u; return s })
    fullNameLens  = L.MakeLens(func(s Enriched) string { return s.FullName }, func(s Enriched, n string) Enriched { s.FullName = n; return s })
)

fullName := func(u User) string { return u.FirstName + " " + u.LastName }

result := F.Pipe2(
    RIO.Do(Enriched{}),
    RIO.Bind(userLens.Set,     func(_ Enriched) RIO.ReaderIOResult[User] { return fetchUser(42) }),
    RIO.Let(fullNameLens.Set,  F.Flow2(userLens.Get, fullName)), // read s.User, compute pure string
)
```

`LetTo` adds a **constant** with no computation:

```go
RIO.LetTo(setVersion, "v1.2.3")
```

### `BindTo` — Starting from an Existing Value

When you have an existing `M[T]` and want to project it into a state struct rather than starting from `Do(empty)`:

```go
type State struct{ User User }

result := F.Pipe1(
    fetchUser(42),                                           // ReaderIOResult[User]
    RIO.BindTo(func(u User) State { return State{User: u} }),// ReaderIOResult[State]
)
```

### Lens Variants (`ApSL`, `BindL`, `LetL`, `LetToL`)

If you have a `Lens[S, T]` (from the optics system or code generation), you can skip writing the setter function entirely:

```go
import (
    RO "github.com/IBM/fp-go/v2/readeroption"
    F  "github.com/IBM/fp-go/v2/function"
)

// Lenses generated by go:generate (see optics/README.md)
// personLenses.Name : Lens[*Person, Name]
// personLenses.Age  : Lens[*Person, Age]

makePerson := F.Pipe2(
    RO.Do[*PartialPerson](emptyPerson),
    RO.ApSL(personLenses.Name, maybeName), // replaces: ApS(personLenses.Name.Set, maybeName)
    RO.ApSL(personLenses.Age,  maybeAge),
)
```

This exact pattern is used in [`samples/builder`](samples/builder/builder.go) to validate and construct a `Person` from an unvalidated `PartialPerson`.

### Lifted Variants for Mixed Monads

`context/readerioresult` provides `Bind*K` helpers that lift simpler computations directly into the do-chain:

| Helper | Lifts |
|--------|-------|
| `BindResultK` / `BindEitherK` | `func(S1) (T, error)` — pure result |
| `BindIOResultK` / `BindIOEitherK` | `func(S1) func() (T, error)` — lazy IO result |
| `BindIOK` | `func(S1) func() T` — infallible IO |
| `BindReaderK` | `func(S1) func(ctx) T` — context reader |

```go
RIO.BindResultK(setUser, func(s Pipeline) (User, error) {
    return validateAndBuild(s)   // plain (value, error) function, no wrapping needed
})
```

### Decision Guide

```
Does the new step need to read prior accumulated state?
    YES  →  Bind   (monadic, sequential; f receives current S)
    NO   →  ApS    (applicative, independent; fa is a fixed M[T])

Is the new value derived purely from state, with no monad?
    YES  →  Let    (pure function of S)

Is the new value a compile-time or runtime constant?
    YES  →  LetTo

Starting from an existing M[T] rather than an empty struct?
    YES  →  BindTo
```

### Complete Example — `result` Monad

The same pattern works with simpler monads. Here with `result.Result[A]`:

`Eitherize1` converts any standard `func(A) (B, error)` into `func(A) Result[B]`. Define these lifted functions once as variables. Then use lenses to focus on individual struct fields and compose with `F.Flow2(lens.Get, f)` — no inline lambdas, no manual error handling.

```go
import (
    R    "github.com/IBM/fp-go/v2/result"
    F    "github.com/IBM/fp-go/v2/function"
    L    "github.com/IBM/fp-go/v2/optics/lens"
    N    "github.com/IBM/fp-go/v2/number"
    "strconv"
)

type Parsed struct {
    Raw    string
    Number int
    Double int
}

// Lenses — focus on individual fields of Parsed.
var (
    rawLens    = L.MakeLens(
        func(s Parsed) string { return s.Raw },
        func(s Parsed, v string) Parsed { s.Raw = v; return s },
    )
    numberLens = L.MakeLens(
        func(s Parsed) int { return s.Number },
        func(s Parsed, v int) Parsed { s.Number = v; return s },
    )
    doubleLens = L.MakeLens(
        func(s Parsed) int { return s.Double },
        func(s Parsed, v int) Parsed { s.Double = v; return s },
    )
)

// Lifted functions — convert standard (value, error) functions into Result-returning ones.
var (
    atoi = R.Eitherize1(strconv.Atoi) // func(string) Result[int]
)

parse := func(input string) R.Result[Parsed] {
    return F.Pipe3(
        R.Do(Parsed{}),
        R.LetTo(rawLens.Set, input),                      // set Raw to constant input
        R.Bind(numberLens.Set, F.Flow2(rawLens.Get, atoi)),           // get Raw, parse → Result[int]
        R.Let(doubleLens.Set, F.Flow2(numberLens.Get, N.Mul(2))),   // get Number, multiply → int
    )
}

parse("21")  // Ok(Parsed{Raw:"21", Number:21, Double:42})
parse("abc") // Error(strconv parse error)
```

`rawLens.Set` is already `func(string) func(Parsed) Parsed`, matching the setter signature `Bind` and `LetTo` expect — no manual setter functions to write. `F.Flow2(rawLens.Get, atoi)` composes the field getter with the eitherized parse function into a `Kleisli[Parsed, int]` without any intermediate lambda.

## Import Paths

```go
import (
    "github.com/IBM/fp-go/v2/option"
    "github.com/IBM/fp-go/v2/result"
    "github.com/IBM/fp-go/v2/either"
    "github.com/IBM/fp-go/v2/io"
    "github.com/IBM/fp-go/v2/ioresult"
    "github.com/IBM/fp-go/v2/context/readerioresult"
    "github.com/IBM/fp-go/v2/effect"
    F "github.com/IBM/fp-go/v2/function"
    A "github.com/IBM/fp-go/v2/array"
)
```

Requires Go 1.24+ (generic type aliases).
