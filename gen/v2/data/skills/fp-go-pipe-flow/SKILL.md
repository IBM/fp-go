---
name: fp-go-pipe-flow
description: >
  Guides writing, refactoring, and reviewing fp-go v2 code that uses functional
  composition via Pipe and Flow. Apply this skill whenever the user asks to write
  new fp-go code, refactor imperative Go into functional style, explain a Pipe/Flow
  pipeline, add do-notation (Do/Bind/ApS), use the reader monad, create lenses for
  struct fields, or generate unit tests for functional pipelines. Trigger on any
  mention of: Pipe, Flow, reader monad, kleisli, do-notation, Bind, ApS, fp-go
  pipeline, point-free style, lens composition in Go.
---

# fp-go Pipe and Flow Patterns

All imports **must** come from `github.com/IBM/fp-go/v2`, never from
`github.com/IBM/fp-go` (the v1 path).

---

## Before you generate

fp-go is low-frequency in training data, so signatures are easy to misremember.
For any combinator not shown below, look it up via the fp-go MCP server's
`search_examples` / `get_example` tools (see the **fp-go-mcp** skill) instead of
guessing. After writing code, run `go build ./...` and `go vet ./...` and fix any
type-parameter or argument-order errors before presenting it.

---

## Core concepts

### Pipe — data-first composition

`Pipe` takes an initial value and threads it through a sequence of functions.
Use it when you already have a value to start from.

```go
import F "github.com/IBM/fp-go/v2/function"

// PipeN threads a value through N functions
result := F.Pipe3(initialValue, step1, step2, step3)
```

The number suffix matches the number of transformation steps (Pipe1 … Pipe26).

### Flow — function-first composition

`Flow` composes N functions into a single function that awaits its input.
Use it to build reusable pipeline functions, especially as arguments to `Map`,
`Chain`, or `TraverseArray`.

```go
// FlowN returns func(T0) TN
pipeline := F.Flow3(step1, step2, step3)
result := pipeline(initialValue)
```

**Rule of thumb: prefer `Pipe` when you have the starting value; use `Flow`
when you are building a reusable function.**

---

## Prefer functions over variables

Go does not eliminate dead variables, but unused functions are zero-cost.
Always wrap a `Pipe`/`Flow` result in a named function rather than storing it
in a package-level `var`.

```go
// WRONG — var is allocated even if never called
var processUser = F.Flow2(getName, strings.ToUpper)

// CORRECT — zero cost until called; also more composable
func processUser() func(User) string {
    return F.Flow2(getName, strings.ToUpper)
}
```

Use `var` only for lenses and pre-bound combinator helpers (like `lens.Get`
assigned to a named getter), not for full pipeline results.

---

## Point-free style

Avoid explicit argument names wherever a named combinator or `Flow` can express
the same thing.

```go
// WRONG — explicit argument
func isAdult(u User) bool { return getAge(u) > 18 }

// CORRECT — point-free, returns typed Predicate
func isAdult() P.Predicate[User] {
    return F.Flow2(getAge, N.MoreThan(18))
}
```

### Type aliases to use

| Type | Package | Meaning |
|------|---------|---------|
| `P.Predicate[A]` | `github.com/IBM/fp-go/v2/predicate` | `func(A) bool` |
| `E.Endomorphism[A]` | `github.com/IBM/fp-go/v2/endomorphism` | `func(A) A` |

Use these as return types for functions that act as predicates or
self-transformations — they communicate intent and enable direct use in
combinators like `A.Filter`, `A.Map`, `F.Ternary`.

```go
import (
    F "github.com/IBM/fp-go/v2/function"
    N "github.com/IBM/fp-go/v2/number"
    P "github.com/IBM/fp-go/v2/predicate"
    E "github.com/IBM/fp-go/v2/endomorphism"
    A "github.com/IBM/fp-go/v2/array"
)

// Predicate — point-free using N.MoreThan
func isAdult() P.Predicate[User] {
    return F.Flow2(getAge, N.MoreThan(18))
}

// Endomorphism — self-transformation
func doubleAll() E.Endomorphism[[]int] {
    return A.Map[int, int](N.Mul(2))
}
```

### Numeric combinators

Prefer `N.MoreThan`, `N.LessThan`, `N.Mul`, `N.Add` etc. over inline
comparisons or arithmetic in lambdas:

```go
N.MoreThan(18)   // func(int) bool   — x > 18
N.LessThan(100)  // func(int) bool   — x < 100
N.Mul(2)         // func(int) int    — x * 2
N.Add(1)         // func(int) int    — x + 1
```

---

## Pure pipelines vs the reader monad

**Only use the reader monad when the computation genuinely needs an environment
(context, config, DB, logger, etc.).** For pure transformations that don't
need external input, use `Flow` or `Pipe` directly — no reader wrapping needed.

```go
// WRONG — forces reader monad on a pure computation
func adultNames(users []User) R.Reader[context.Context, string] {
    return F.Pipe2(
        R.Of[context.Context](users),
        R.Map[context.Context, []User, string](pureTransform),
    )
}

// CORRECT — pure; no environment needed
func adultNames() func([]User) string {
    return F.Flow2(
        A.FilterMap(toAdultName()),
        A.Intercalate(S.Monoid)(","),
    )
}
```

### Per-element filter+map: use `A.FilterMap`

When filtering and then extracting a field, combine both into a single pass
with `A.FilterMap` and `O.FromPredicate`:

```go
import (
    F "github.com/IBM/fp-go/v2/function"
    A "github.com/IBM/fp-go/v2/array"
    O "github.com/IBM/fp-go/v2/option"
    N "github.com/IBM/fp-go/v2/number"
    P "github.com/IBM/fp-go/v2/predicate"
    S "github.com/IBM/fp-go/v2/string"
    E "github.com/IBM/fp-go/v2/endomorphism"
)

// isAdult — point-free predicate
func isAdult() P.Predicate[User] {
    return F.Flow2(getAge, N.MoreThan(18))
}

// toAdultName — User -> Option[string]: Some(name) if adult, None otherwise
func toAdultName() func(User) O.Option[string] {
    return F.Flow2(
        O.FromPredicate(isAdult()),  // User -> Option[User]
        O.Map(getName),              // Option[User] -> Option[string]
    )
}

// adultNames — pure pipeline, no reader monad needed
func adultNames() func([]User) string {
    return F.Flow2(
        A.FilterMap(toAdultName()),       // []User -> []string
        A.Intercalate(S.Monoid)(","),     // []string -> string
    )
}
```

---

## Reader monad

The reader monad `Reader[R, A]` is `func(R) A` — a computation that reads
from an environment `R` and produces `A`. Only reach for it when the
computation needs to thread an environment (e.g. `context.Context`, a
config struct, a DB handle).

```go
import (
    F   "github.com/IBM/fp-go/v2/function"
    R   "github.com/IBM/fp-go/v2/reader"
    "context"
)

// Kleisli arrow: A -> Reader[Env, B]
func fetchUser(id string) R.Reader[context.Context, User] {
    return R.Asks(func(ctx context.Context) User {
        return User{ID: id}
    })
}
```

### When to use `reader.Map` vs full `Pipe` with reader operations

- **`reader.Map`** inside `Flow` — when the step is pure and the environment
  does not need to appear explicitly. This is the "abbreviation" pattern.
- **`Pipe` with `reader.Chain`, `reader.Bind`, `reader.ApS`** — when the
  sequence needs the context (e.g. calls another kleisli arrow) or when
  do-notation makes the data flow clearer.

```go
// reader.Map inside Flow — no env name, clean point-free
func renderUsers() func([]User) R.Reader[context.Context, string] {
    return F.Flow2(
        usersToNames,
        R.Map[context.Context](strings.Join),
    )
}

// Pipe with reader monad — env access required
func enrichedUser(id string) R.Reader[context.Context, EnrichedUser] {
    return F.Pipe3(
        fetchUser(id),
        R.Chain(fetchProfile),
        R.Chain(fetchPermissions),
        R.Map[context.Context](combineToEnriched),
    )
}
```

---

## Do-notation: `Do` / `Bind` / `ApS` / `Let`

Do-notation is the idiomatic way to assemble multiple reader (or IO/result)
computations into a named-field record. Always use it inside a `Pipe`.

```go
import (
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/reader"
    "context"
)

// Setters — hand-written or generated, kept outside the pipe
func setProfile(p Profile) func(RequestState) RequestState {
    return func(s RequestState) RequestState { s.Profile = p; return s }
}

// Kleisli arrows — named functions, never inline
func fetchProfile(s RequestState) R.Reader[context.Context, Profile] {
    return R.Asks(func(ctx context.Context) Profile { /* … */ return Profile{} })
}

// Pipeline — returned as a function, not a var
func buildRequestState(userID string) R.Reader[context.Context, RequestState] {
    return F.Pipe3(
        R.Do[context.Context](RequestState{UserID: userID}),
        R.Bind(setProfile, fetchProfile),
        R.Bind(setPerms, fetchPerms),
        R.Map[context.Context](F.Identity[RequestState]),
    )
}
```

### `Bind` vs `ApS` vs `Let`

| Combinator | When to use |
|-----------|-------------|
| `Bind(setter, kleisli)` | Result depends on accumulated state (sequential) |
| `ApS(setter, reader)` | Result is independent of other fields |
| `Let(setter, pureFunc)` | Pure transformation of accumulated state, no reader needed |
| `LetTo(setter, value)` | Attach a constant value to state |

Use `ApS` when values can be computed independently; `Bind` when a later step
depends on an earlier one. Mixing them in the same pipeline is normal.

---

## Lenses for struct field access

Never access struct fields with inline functions inside a `Pipe`. Create a
lens or a dedicated helper so the pipeline stays point-free.

```go
import (
    L "github.com/IBM/fp-go/v2/optics/lens"
)

var hostLens = L.MakeLens(
    func(c Config) string { return c.Host },
    func(c Config, v string) Config { c.Host = v; return c },
)

// Assign lens.Get to a named var — then pass it anywhere point-free
var getHost = hostLens.Get   // func(Config) string
var getPort = portLens.Get   // func(Config) int
```

Use `R.ApSL(lens, reader)` / `R.BindL(lens, kleisli)` as do-notation variants
that take a lens directly instead of a setter function.

---

## Unit tests

Generate a `_test.go` for every non-trivial pipeline or flow.

```go
func TestAdultNames(t *testing.T) {
    users := []User{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 16}}
    assert.Equal(t, "Alice", adultNames()(users))
}

func TestBuildRequestState(t *testing.T) {
    ctx := context.Background()
    state := buildRequestState("user-42")(ctx)
    assert.Equal(t, "user-42", state.UserID)
}
```

### Testing guidelines

- For pure `Flow`/`Pipe` functions: call the returned function with a concrete
  value and assert with `assert.Equal`.
- For reader pipelines: call the reader with a concrete environment struct.
- For `IOResult`/`ReaderIOResult`: call the innermost IO thunk and compare
  with `result.Of(expected)`.
- Prefer table-driven tests for pipelines with multiple input/output pairs.
- Do not mock the environment — pass a real (but lightweight) struct.

---

## Common import aliases

```go
import (
    F  "github.com/IBM/fp-go/v2/function"
    R  "github.com/IBM/fp-go/v2/reader"
    RR "github.com/IBM/fp-go/v2/readerresult"
    IO "github.com/IBM/fp-go/v2/ioresult"
    E  "github.com/IBM/fp-go/v2/either"
    O  "github.com/IBM/fp-go/v2/option"
    L  "github.com/IBM/fp-go/v2/optics/lens"
    A  "github.com/IBM/fp-go/v2/array"
    N  "github.com/IBM/fp-go/v2/number"
    S  "github.com/IBM/fp-go/v2/string"
    P  "github.com/IBM/fp-go/v2/predicate"
    En "github.com/IBM/fp-go/v2/endomorphism"
)
```

---

## Quick reference

| Goal | Pattern |
|------|---------|
| Thread a value through N steps | `F.PipeN(value, f1, f2, …)` |
| Build a reusable function | `F.FlowN(f1, f2, …)` |
| Point-free numeric predicate | `F.Flow2(getField, N.MoreThan(n))` returning `P.Predicate[T]` |
| Filter+map in one pass | `A.FilterMap(F.Flow2(O.FromPredicate(pred), O.Map(f)))` |
| Lift a pure function into Reader | `R.Map[Env](pureFunc)` |
| Chain kleisli arrows | `R.Chain(kleisliFunc)` |
| Start do-notation block | `R.Do[Env](emptyStruct)` |
| Add dependent field | `R.Bind(setter, kleisliFunc)` |
| Add independent field | `R.ApS(setter, readerValue)` |
| Add pure derived field | `R.Let[Env](setter, pureFunc)` |
| Lens getter in pipeline | `var getX = xLens.Get` |
| Do-notation with lens | `R.ApSL(lens, readerValue)` |
| Access full environment | `R.Ask[Env]()` |
| Access field of environment | `R.Asks(getX)` |
