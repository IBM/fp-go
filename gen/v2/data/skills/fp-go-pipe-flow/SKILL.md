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

The number suffix matches the number of transformation steps. `Pipe1`–`Pipe20` and `Flow1`–`Flow20` are generated; there is nothing above 20.

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
| `EM.Endomorphism[A]` | `github.com/IBM/fp-go/v2/endomorphism` | `func(A) A` |

Use these as return types for functions that act as predicates or
self-transformations — they communicate intent and enable direct use in
combinators like `A.Filter`, `A.Map`, `F.Ternary`.

```go
import (
    F "github.com/IBM/fp-go/v2/function"
    N "github.com/IBM/fp-go/v2/number"
    P "github.com/IBM/fp-go/v2/predicate"
    EM "github.com/IBM/fp-go/v2/endomorphism"
    A "github.com/IBM/fp-go/v2/array"
)

// Predicate — point-free using N.MoreThan
func isAdult() P.Predicate[User] {
    return F.Flow2(getAge, N.MoreThan(18))
}

// Endomorphism — self-transformation
func doubleAll() EM.Endomorphism[[]int] {
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
func adultNames(users []User) RD.Reader[Env, string] {
    return F.Pipe2(
        RD.Of[Env](users),
        RD.Map[Env](pureTransform),
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
computation needs to thread an environment (a config struct, a repository,
…). For `context.Context` plus IO and errors use `context/readerioresult`
(`RIO`) instead — see the **fp-go-context** skill.

```go
import (
    F  "github.com/IBM/fp-go/v2/function"
    RD "github.com/IBM/fp-go/v2/reader"
)

type Env struct {
    Users map[string]User
}

// Leaf accessor (or a generated lens' .Get)
func getUsers(e Env) map[string]User { return e.Users }

// lookupUser is curried: func(id string) func(map[string]User) User

// Kleisli arrow: string -> Reader[Env, User], built from a pure projection
func fetchUser(id string) RD.Reader[Env, User] {
    return RD.Asks(F.Flow2(getUsers, lookupUser(id)))
}
```

### When to use `reader.Map` vs full `Pipe` with reader operations

- **`reader.Map`** inside `Flow` — when the step is pure and the environment
  does not need to appear explicitly. This is the "abbreviation" pattern.
- **`Pipe` with `reader.Chain`, `reader.Bind`, `reader.ApS`** — when the
  sequence needs the environment (e.g. calls another kleisli arrow) or when
  do-notation makes the data flow clearer.

```go
// reader.Map inside Flow — no env name, clean point-free.
// NOTE: RD.Map returns an Operator over Reader values, so the PRECEDING step in the
// Flow must already produce a Reader. A plain func([]User) []string composed with
// RD.Map does not type-check.
func renderUsers() func(string) RD.Reader[Env, string] {
    return F.Flow2(
        fetchTeam,                   // string -> Reader[Env, []User]
        RD.Map[Env](F.Flow2(         // Reader[Env, []User] -> Reader[Env, string]
            A.Map(getName),
            A.Intercalate(S.Monoid)(","),
        )),
    )
}

// Pipe with reader monad — env access required
func enrichedUser(id string) RD.Reader[Env, EnrichedUser] {
    return F.Pipe3(
        fetchUser(id),
        RD.Chain(fetchProfile),
        RD.Chain(fetchPermissions),
        RD.Map[Env](combineToEnriched),
    )
}
```

---

## Do-notation: `Do` / `Bind` / `ApS` / `Let`

Do-notation is the idiomatic way to assemble multiple reader (or IO/result)
computations into a named-field record. Always use it inside a `Pipe`.

```go
import (
    F  "github.com/IBM/fp-go/v2/function"
    L  "github.com/IBM/fp-go/v2/optics/lens"
    RD "github.com/IBM/fp-go/v2/reader"
)

// Lenses — generated (`// fp-go:Lens`) or built once with L.MakeLens.
// lens.Set already has the setter shape func(T) func(S) S — no hand-written setters.
var (
    userIDLens = L.MakeLens(
        func(s RequestState) string { return s.UserID },
        func(s RequestState, v string) RequestState { s.UserID = v; return s },
    )
    profileLens = L.MakeLens(
        func(s RequestState) Profile { return s.Profile },
        func(s RequestState, v Profile) RequestState { s.Profile = v; return s },
    )
    permsLens = L.MakeLens(
        func(s RequestState) Perms { return s.Perms },
        func(s RequestState, v Perms) RequestState { s.Perms = v; return s },
    )
)

// Kleisli arrows — named functions, never inline:
//   fetchProfile: func(userID string) Reader[Env, Profile]
//   fetchPerms:   func(p Profile)     Reader[Env, Perms]

// Pipeline — returned as a function, not a var
func buildRequestState(userID string) RD.Reader[Env, RequestState] {
    return F.Pipe3(
        RD.Do[Env](RequestState{}),
        RD.LetTo[Env](userIDLens.Set, userID),
        RD.Bind(profileLens.Set, F.Flow2(userIDLens.Get, fetchProfile)),
        RD.Bind(permsLens.Set, F.Flow2(profileLens.Get, fetchPerms)),
    )
}
```

`F.Flow2(lens.Get, kleisli)` is the point-free way to feed one field of the
accumulated state into the next step.

### `Bind` vs `ApS` vs `Let`

| Combinator | When to use |
|-----------|-------------|
| `Bind(setter, kleisli)` | Result depends on accumulated state (sequential) |
| `ApS(setter, reader)` | Result is independent of other fields |
| `Let(setter, pureFunc)` | Pure transformation of accumulated state, no reader needed |
| `LetTo(setter, value)` | Attach a constant value to state |

Use `ApS` when values can be computed independently; `Bind` when a later step
depends on an earlier one. Mixing them in the same pipeline is normal. A `Bind`
whose Kleisli ignores the state (`func(_ S) M[T] { return m }`) is always an
`ApS(setter, m)`.

---

## Lenses for struct field access

Never access struct fields with inline functions inside a `Pipe`. Create a
lens (preferably generated with `// fp-go:Lens`, see the **fp-go-lens** skill)
or a named leaf accessor so the pipeline stays point-free.

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

Use `RD.ApSL(lens, reader)` / `RD.BindL(lens, kleisli)` as do-notation variants
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
    env := Env{Users: map[string]User{"user-42": {ID: "user-42"}}}
    state := buildRequestState("user-42")(env)
    assert.Equal(t, "user-42", state.UserID)
}
```

### Testing guidelines

- For pure `Flow`/`Pipe` functions: call the returned function with a concrete
  value and assert with `assert.Equal`.
- For reader pipelines: call the reader with a concrete environment struct.
- For `IOResult`/`ReaderIOResult`: call the innermost IO thunk and compare
  with `R.Of(expected)`; run a `ReaderIOResult` with `t.Context()`.
- Prefer table-driven tests for pipelines with multiple input/output pairs.
- Do not mock the environment — pass a real (but lightweight) struct.

---

## Common import aliases

These follow the canonical alias table in the **fp-go** skill.

```go
import (
    F   "github.com/IBM/fp-go/v2/function"
    A   "github.com/IBM/fp-go/v2/array"
    O   "github.com/IBM/fp-go/v2/option"
    E   "github.com/IBM/fp-go/v2/either"
    R   "github.com/IBM/fp-go/v2/result"
    IOR "github.com/IBM/fp-go/v2/ioresult"
    RD  "github.com/IBM/fp-go/v2/reader"
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    L   "github.com/IBM/fp-go/v2/optics/lens"
    N   "github.com/IBM/fp-go/v2/number"
    S   "github.com/IBM/fp-go/v2/string"
    P   "github.com/IBM/fp-go/v2/predicate"
    EM  "github.com/IBM/fp-go/v2/endomorphism"
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
| Lift a pure function into Reader | `RD.Map[Env](pureFunc)` |
| Chain kleisli arrows | `RD.Chain(kleisliFunc)` |
| Start do-notation block | `RD.Do[Env](emptyStruct)` |
| Add dependent field | `RD.Bind(lens.Set, F.Flow2(otherLens.Get, kleisliFunc))` |
| Add independent field | `RD.ApS(lens.Set, readerValue)` |
| Add pure derived field | `RD.Let[Env](lens.Set, F.Flow2(otherLens.Get, pureFunc))` |
| Lens getter in pipeline | `var getX = xLens.Get` |
| Do-notation with lens | `RD.ApSL(lens, readerValue)` |
| Access full environment | `RD.Ask[Env]()` |
| Access field of environment | `RD.Asks(getX)` |
| Read a `context.Context` value | `RIO.AskValue[V](key)` → `Option[V]` (not `ctx.Value(key).(V)`) |
| Scope a value / timeout to a step | `RIO.WithValue[A](key, v)`, `RIO.WithTimeout[A](d)` as the last `Pipe` step |
