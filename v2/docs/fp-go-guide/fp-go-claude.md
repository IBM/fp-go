# fp-go/v2 Reference for Claude Code

fp-go/v2 (`github.com/IBM/fp-go/v2`) is a typed functional programming library for Go 1.24+. It provides Option, Either/Result, IO, and Effect monads with data-last, curried APIs designed for pipeline composition via `Pipe` and `Flow`. The library follows Haskell/fp-ts conventions adapted for Go generics with explicit arity-numbered functions (e.g., `Pipe3`, `Flow2`). Two module families exist: **standard** (struct-based monads, full FP toolkit) and **idiomatic** (Go-native `(value, error)` tuples, zero-alloc, 2-32x faster).

## Import Conventions

| Alias | Package |
|-------|---------|
| `F`   | `github.com/IBM/fp-go/v2/function` |
| `O`   | `github.com/IBM/fp-go/v2/option` |
| `E`   | `github.com/IBM/fp-go/v2/either` |
| `R`   | `github.com/IBM/fp-go/v2/result` |
| `A`   | `github.com/IBM/fp-go/v2/array` |
| `IO`  | `github.com/IBM/fp-go/v2/io` |
| `IOR` | `github.com/IBM/fp-go/v2/ioresult` |
| `IOE` | `github.com/IBM/fp-go/v2/ioeither` |
| `RIO` | `github.com/IBM/fp-go/v2/context/readerioresult` |
| `EFF` | `github.com/IBM/fp-go/v2/effect` |
| `P`   | `github.com/IBM/fp-go/v2/pair` |
| `T`   | `github.com/IBM/fp-go/v2/tuple` |
| `N`   | `github.com/IBM/fp-go/v2/number` |
| `S`   | `github.com/IBM/fp-go/v2/string` |
| `B`   | `github.com/IBM/fp-go/v2/boolean` |
| `L`   | `github.com/IBM/fp-go/v2/optics/lens` |
| `PR`  | `github.com/IBM/fp-go/v2/optics/prism` |

**Idiomatic variants** (tuple-based, zero-alloc):

| Alias | Package |
|-------|---------|
| `IR`  | `github.com/IBM/fp-go/v2/idiomatic/result` |
| `IO_` | `github.com/IBM/fp-go/v2/idiomatic/option` |
| `IIR` | `github.com/IBM/fp-go/v2/idiomatic/ioresult` |
| `IRR` | `github.com/IBM/fp-go/v2/idiomatic/context/readerresult` |
| `IRO` | `github.com/IBM/fp-go/v2/idiomatic/readerioresult` |

## Monad Selection

- **Pure value** -- use the value directly, no wrapper needed
- **May be absent** -- `Option[A]` (struct-based) or `(A, bool)` (idiomatic)
- **Can fail with `error`** -- `Result[A]` = `Either[error, A]`
  - Need custom error type E -- use `Either[E, A]` instead
- **Lazy + can fail** -- `IOResult[A]` = `func() Either[error, A]`
  - Idiomatic: `func() (A, error)`
- **Needs `context.Context` + lazy + can fail** -- `ReaderIOResult[A]` via `context/readerioresult`
  - Type: `func(context.Context) func() Either[error, A]`
  - Idiomatic: `func(context.Context) (A, error)` via `idiomatic/context/readerresult`
- **Typed DI + context + lazy + can fail** -- `Effect[C, A]` via `effect` package
  - Type: `func(C) func(context.Context) func() Either[error, A]`
  - C is your dependency/config struct; context.Context is handled internally
- **Performance-critical** -- prefer `idiomatic/` variants throughout

## Standard vs Idiomatic

| Aspect | Standard | Idiomatic |
|--------|----------|-----------|
| Representation | `Either[error, A]` struct | `(A, error)` tuple |
| Performance | Baseline | 2-32x faster, zero allocs |
| Custom error types | `Either[E, A]` for any E | error only |
| Do-notation | Full support | Full support |
| FP toolkit | Complete | Complete |
| Go interop | Requires `Unwrap`/`Eitherize` | Native `(val, err)` |

**Rule of thumb**: Use idiomatic for production code and hot paths. Use standard when you need custom error types (`Either[E, A]`) or when composing with packages that use the standard types.

## Core Types

```go
// function package
type Void = struct{}
var VOID Void = struct{}{}

// option
type Option[A any] struct { /* Some/None */ }

// either
type Either[E, A any] struct { /* Left/Right */ }

// result (specialized Either)
type Result[A any] = Either[error, A]

// io
type IO[A any] = func() A

// ioresult
type IOResult[A any] = IO[Result[A]]  // = func() Either[error, A]

// context/readerioresult
type ReaderIOResult[A any] = func(context.Context) func() Either[error, A]

// effect
type Effect[C, A any] = func(C) func(context.Context) func() Either[error, A]
type Kleisli[C, A, B any] = func(A) Effect[C, B]

// idiomatic equivalents
type IOResult[A any] = func() (A, error)
type ReaderResult[A any] = func(context.Context) (A, error)
```

## Key Rules

1. **Data-last**: Configuration/behavior params come first, data comes last. This enables partial application and pipeline composition.

2. **Type parameter ordering**: Non-inferrable type params come first. Example: `Ap[B, E, A]` -- B cannot be inferred, so it leads. `Map[A, B]` -- both usually inferred.

3. **Composition direction**:
   - `F.Flow1/2/3/.../N` -- left-to-right (use this for pipelines)
   - `Compose` -- right-to-left (mathematical convention; avoid in pipelines)

4. **Pipe vs Flow**:
   - `F.Pipe3(value, f1, f2, f3)` -- apply data to a pipeline immediately
   - `F.Flow3(f1, f2, f3)` -- create a reusable pipeline (returns a function)

5. **Arity-numbered functions**: `Pipe1` through `Pipe20`, `Flow1` through `Flow20`. Choose the number matching your operation count.

6. **Naming conventions**:
   - `Chain` = flatMap/bind (`A -> F[B]`, flattens)
   - `Map` = fmap (`A -> B`, lifts into context)
   - `Ap` = applicative apply (apply wrapped function to wrapped value)
   - `ChainFirst` / `Tap` = execute for side effects, keep original value
   - `ChainEitherK` = lift pure `func(A) Either[E, B]` into monadic chain
   - `Of` = pure/return (lift value into monad)
   - `Fold` = catamorphism (handle both cases)
   - `Left` / `Right` = Either constructors
   - `Some` / `None` = Option constructors

7. **Prefer `result` over `either`** unless you need a custom error type E. `Result[A]` = `Either[error, A]`.

8. **Wrapping Go functions**:
   - `result.Eitherize1(fn)` wraps `func(X) (Y, error)` into `func(X) Result[Y]`
   - `result.Eitherize2(fn)` wraps `func(X, Y) (Z, error)` into `func(X, Y) Result[Z]`
   - Variants up to `Eitherize15`

9. **Use `function.Void` / `function.VOID`** instead of `struct{}` / `struct{}{}`.

10. **Go 1.24+ required** (generic type aliases).

## Common Patterns

### Pipeline with Pipe
```go
result := F.Pipe3(
    inputValue,
    R.Map(transform),
    R.Chain(validate),
    R.Fold(onError, onSuccess),
)
```

### Reusable pipeline with Flow
```go
pipeline := F.Flow3(
    R.Map(normalize),
    R.Chain(validate),
    R.Map(format),
)
output := pipeline(R.Of(input))
```

### Wrapping Go error functions
```go
safeParseInt := R.Eitherize1(strconv.Atoi)
// safeParseInt: func(string) Result[int]
result := safeParseInt("42") // Right(42)
```

### Effect with DI
```go
type Deps struct { DB *sql.DB }

fetchUser := EFF.Eitherize(func(deps Deps, ctx context.Context) (*User, error) {
    return deps.DB.QueryRowContext(ctx, "SELECT ...").Scan(...)
})
// fetchUser: Effect[Deps, *User]

// Execute:
val, err := EFF.RunSync(EFF.Provide[*User](myDeps)(fetchUser))(ctx)
```

### Effect composition
```go
pipeline := F.Pipe1(
    fetchUser,
    EFF.Map[Deps](func(u *User) string { return u.Name }),
)
```

### Do-notation (building up state)
```go
type State struct { X int; Y string }

result := F.Pipe3(
    R.Do(State{}),
    R.Bind(
        func(x int) func(State) State {
            return func(s State) State { s.X = x; return s }
        },
        func(s State) Result[int] { return R.Of(42) },
    ),
    R.Let(
        func(y string) func(State) State {
            return func(s State) State { s.Y = y; return s }
        },
        func(s State) string { return fmt.Sprintf("val=%d", s.X) },
    ),
)
```

### Optics (Lens)
```go
type Person struct { Name string; Age int }

nameLens := L.MakeLens(
    func(p Person) string { return p.Name },
    func(p Person, name string) Person { p.Name = name; return p },
)

name := nameLens.Get(person)              // get
updated := nameLens.Set("Bob")(person)    // set (returns new Person)
modified := L.Modify(strings.ToUpper)(nameLens)(person) // modify
```

### Option handling
```go
result := F.Pipe3(
    O.Some(42),
    O.Map(func(x int) int { return x * 2 }),
    O.GetOrElse(F.Constant(0)),
)
```

### Idiomatic IOResult
```go
readFile := func() ([]byte, error) { return os.ReadFile("config.json") }
// This IS an idiomatic IOResult[[]byte] -- just a func() ([]byte, error)

parsed := IIR.Map(parseConfig)(readFile)
config, err := parsed()
```

### ReaderIOResult (context-dependent IO)
```go
// Eitherize1 wraps func(context.Context, T0) (R, error) -> func(T0) ReaderIOResult[R]
fetchURL := RIO.Eitherize1(func(ctx context.Context, url string) ([]byte, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := http.DefaultClient.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    return iolib.ReadAll(resp.Body)
})
// fetchURL: func(string) ReaderIOResult[[]byte]
result := fetchURL("https://example.com")(ctx)() // execute
```

## Deeper Documentation

- `fp-go-cookbook.md` -- migration recipes and "how do I X in fp-go?"
- `fp-go-core-patterns.md` -- core types, operations, and composition details
- `fp-go-mastery.md` -- advanced FP techniques, architecture, and Effect system
- `fp-go-full-reference.md` -- complete API inventory across all packages
