# fp-go/v2 Advanced Techniques

> LLM context document. Import path: `github.com/IBM/fp-go/v2`.
> Covers advanced FP patterns for experienced practitioners working with fp-go v2.

---

## 1. Do Notation Deep Dive

Do notation simulates Haskell/Scala for-comprehension using struct embedding and curried setters. It exists in every monad package: `result`, `either`, `stateio`, `statereaderioeither`, `state`, `record`, `iterator/iter`, and more.

### Core Operations

All do-notation functions share the same pattern across packages. Below uses `result` (which delegates to `either`) as the reference.

```go
// Start a do-notation chain with an initial empty state
func Do[S any](empty S) Result[S]

// Bind: run an effectful computation, inject result into state via setter
func Bind[S1, S2, T any](
    setter func(T) func(S1) S2,
    f Kleisli[S1, T],          // func(S1) Result[T]
) Operator[S1, S2]             // func(Result[S1]) Result[S2]

// Let: run a pure computation (no effect), inject into state
func Let[S1, S2, T any](
    key func(T) func(S1) S2,
    f func(S1) T,
) Operator[S1, S2]

// LetTo: inject a constant value into state
func LetTo[S1, S2, T any](
    key func(T) func(S1) S2,
    b T,
) Operator[S1, S2]

// BindTo: start the chain by wrapping a value into an initial state
func BindTo[S1, T any](
    setter func(T) S1,
) Operator[T, S1]

// ApS: applicative version of Bind (independent computation)
func ApS[S1, S2, T any](
    setter func(T) func(S1) S2,
    fa Result[T],
) Operator[S1, S2]
```

### Lens-Based Variants

When working with struct fields that have optics, use the `L` suffixed variants. These take a `Lens[S, T]` instead of a manual setter, operating on the same struct type `S`:

```go
// BindL: monadic bind using a lens to update a field
func BindL[S, T any](
    lens Lens[S, T],
    f Kleisli[T, T],     // func(T) Result[T]
) Operator[S, S]

// LetL: pure transformation of a field via lens
func LetL[S, T any](
    lens Lens[S, T],
    f Endomorphism[T],   // func(T) T
) Operator[S, S]

// LetToL: set a field to a constant via lens
func LetToL[S, T any](
    lens Lens[S, T],
    b T,
) Operator[S, S]

// ApSL: applicative bind using a lens
func ApSL[S, T any](
    lens Lens[S, T],
    fa Result[T],
) Operator[S, S]
```

### Building a Complex Struct Step-by-Step

The pattern uses struct embedding to grow the state type at each step:

```go
import (
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
)

// Define states by embedding. Each step adds a field.
type Initial struct{}

type WithHost struct {
    Initial
    Host string
}

type WithPort struct {
    WithHost
    Port int
}

type WithDB struct {
    WithPort
    DBConn string
}

// Curried setters: func(T) func(S1) S2
var SetHost = F.Curry2(func(host string, s Initial) WithHost {
    return WithHost{Initial: s, Host: host}
})

var SetPort = F.Curry2(func(port int, s WithHost) WithPort {
    return WithPort{WithHost: s, Port: port}
})

var SetDB = F.Curry2(func(db string, s WithPort) WithDB {
    return WithPort{WithPort: s, DBConn: db}
})

// Effectful lookups
func lookupHost() R.Result[string]          { return R.Of("localhost") }
func lookupPort(host string) R.Result[int]  { return R.Of(8080) }
func connectDB(host string, port int) R.Result[string] {
    return R.Of(fmt.Sprintf("%s:%d", host, port))
}

// Compose with do-notation
config := F.Pipe4(
    R.Do(Initial{}),
    R.Bind(SetHost, func(_ Initial) R.Result[string] {
        return lookupHost()
    }),
    R.Bind(SetPort, func(s WithHost) R.Result[int] {
        return lookupPort(s.Host)
    }),
    R.Bind(SetDB, func(s WithPort) R.Result[string] {
        return connectDB(s.Host, s.Port)
    }),
)
// config: Result[WithDB] with all fields populated
```

### StateIO Do Notation

`stateio` provides the same pattern but each step can perform IO:

```go
func Do[ST, A any](s ST, a A) StateIO[ST, A]
func Bind[ST, S1, S2, T any](setter, effect) Operator[ST, S1, S2]
```

The `ST` type parameter is the mutable state threaded through the IO computation, while `S1`/`S2` are the do-notation accumulator types.

---

## 2. Kleisli Composition

### What Kleisli Arrows Are

A Kleisli arrow is a function `func(A) M[B]` where `M` is a monad. In fp-go, every monad package defines:

```go
type Kleisli[A, B any] = func(A) M[B]
```

For `result`: `type Kleisli[A, B any] = func(A) Result[B]`
For `effect`: `type Kleisli[C, A, B any] = func(A) Effect[C, B]`
For `reader`: `type Kleisli[R, A, B any] = func(A) Reader[R, B]`

### Chain Composes Kleisli Arrows

`Chain` is monadic bind in curried form. It sequences two computations where the second depends on the first:

```go
// Chain for result
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]

// Usage: pipe a Result[A] through Chain to get Result[B]
result := F.Pipe2(
    getUser(id),                    // Result[User]
    result.Chain(func(u User) Result[Order] {
        return getLatestOrder(u.ID)
    }),
    result.Chain(func(o Order) Result[string] {
        return formatOrder(o)
    }),
)
```

### ChainK Variants: Lifting Between Layers

`ChainK` variants lift a Kleisli arrow from a "lower" monad into a "higher" one. This avoids manual `FromX` + `Chain` composition:

```go
// effect package (Effect[C, A] = ReaderReaderIOResult[C, A])
func ChainIOK[C, A, B any](f io.Kleisli[A, B]) Operator[C, A, B]
func ChainResultK[C, A, B any](f result.Kleisli[A, B]) Operator[C, A, B]
func ChainReaderK[C, A, B any](f reader.Kleisli[C, A, B]) Operator[C, A, B]
func ChainThunkK[C, A, B any](f thunk.Kleisli[A, B]) Operator[C, A, B]

// ioresult package
func ChainIOK[A, B any](f io.Kleisli[A, B]) Operator[A, B]
func ChainEitherK[A, B any](f result.Kleisli[A, B]) Operator[A, B]
func ChainResultK[A, B any](f result.Kleisli[A, B]) Operator[A, B]

// context/readerioresult package
func ChainEitherK[A, B any](f either.Kleisli[error, A, B]) Operator[A, B]
func ChainResultK[A, B any](f either.Kleisli[error, A, B]) Operator[A, B]
func ChainIOK[A, B any](f io.Kleisli[A, B]) Operator[A, B]
func ChainIOResultK[A, B any](f ioresult.Kleisli[A, B]) Operator[A, B]
func ChainReaderK[A, B any](f reader.Kleisli[context.Context, A, B]) Operator[A, B]
```

### Composing Kleisli Arrows with Flow

Use `F.Flow2`/`F.Flow3` etc. to compose Kleisli arrows into reusable pipelines:

```go
// Compose two Kleisli arrows: func(A) Result[B] + func(B) Result[C] => func(A) Result[C]
getUserOrders := F.Flow2(
    getUser,        // func(ID) Result[User]
    result.Chain(getOrders),  // func(User) Result[[]Order]
)
// getUserOrders: func(ID) Result[[]Order]
```

### When to Use Chain vs ChainK

- **Chain**: both computations live in the same monad
- **ChainK**: the inner computation lives in a "lower" monad (IO, Either, Result, Reader) and needs lifting into the outer monad (Effect, IOResult, ReaderIOResult)
- Rule: if your function returns `IO[B]` but you're in an `Effect[C, _]` pipeline, use `ChainIOK`

---

## 3. Monad Transformer Stacks

### Layer Diagram

```
Effect[C, A]
  = ReaderReaderIOResult[C, A]     -- adds an extra reader layer (C)
  = func(C) ReaderIOResult[A]      -- strips outer reader
  = func(C) func(context.Context) IOResult[A]
  = func(C) func(context.Context) func() Either[error, A]

                    Layer Stack
    +-----------------------------------------+
    |  Reader[C, _]        -- business config  |
    |  Reader[context.Context, _] -- Go ctx    |
    |  IO[_]               -- side effects     |
    |  Either[error, _]    -- error handling   |
    +-----------------------------------------+
```

### Key Type Aliases

```
Reader[R, A]          = func(R) A
IO[A]                 = func() A
Either[E, A]          = tagged union (Left E | Right A)
Result[A]             = Either[error, A]
IOResult[A]           = func() Either[error, A]
ReaderIOResult[A]     = func(context.Context) func() Either[error, A]
Effect[C, A]          = func(C) func(context.Context) func() Either[error, A]
Thunk[A]              = ReaderIOResult[A]   -- an Effect with context already provided
```

### Choosing the Right Stack

| Need | Use |
|---|---|
| Pure computation, may fail | `Result[A]` / `Either[E, A]` |
| Side effects, may fail | `IOResult[A]` |
| Side effects + context.Context | `context/readerioresult.ReaderIOResult[A]` |
| Side effects + context.Context + config | `Effect[C, A]` |
| Pure state threading | `State[S, A]` |
| State + IO | `StateIO[S, A]` |
| Optional value | `Option[A]` |
| Lazy optional with IO | `IOOption[A]` |

### Lifting Between Layers

From the `effect` package:

```go
// Lift a pure Result into an Effect
func FromResult[C, A any](r Result[A]) Effect[C, A]

// Lift a side-effectful thunk into an Effect
func FromThunk[C, A any](f Thunk[A]) Effect[C, A]
```

Each monad package provides `FromIO`, `FromEither`, `FromResult`, etc. to lift lower layers up. The pattern is consistent:

```go
// In ioresult:
ioresult.FromIO[A](ioAction)        // IO[A] -> IOResult[A]
ioresult.FromEither[A](either)      // Either[error, A] -> IOResult[A]

// In context/readerioresult:
readerioresult.FromIOResult(ior)     // IOResult[A] -> ReaderIOResult[A]
readerioresult.FromIO(io)            // IO[A] -> ReaderIOResult[A]
readerioresult.FromEither(e)         // Either[error, A] -> ReaderIOResult[A]
```

### Context Specializations

`context/readerioresult` fixes the reader environment to `context.Context`. Generic `readerioresult` allows any `R`:

```go
// context/readerioresult -- R is always context.Context
type ReaderIOResult[A any] = func(context.Context) func() Either[error, A]

// readerioresult -- R is generic
type ReaderIOResult[R, A any] = func(R) func() Either[error, A]
```

Use `context/` variants when you need Go's `context.Context` for cancellation/deadlines. Use generic variants when the reader environment is application-specific configuration.

### Performance Implications

Deeper stacks add function call overhead. Each layer is a closure allocation. For hot paths:
- Prefer `IOResult` over `Effect` if you don't need the reader layers
- Pre-bind context outside tight loops
- Consider `idiomatic/` variants for performance-critical code (see section 10)

---

## 4. Profunctor Mapping

### Promap

`Promap` transforms both the input (context/environment) and the output (value) of an effect simultaneously. It is contravariant on input and covariant on output:

```go
// effect/profunctor.go
func Promap[E, A, D, B any](
    f Reader[D, E],    // contravariant: transform input D -> E
    g Reader[A, B],    // covariant: transform output A -> B
) Kleisli[D, Effect[E, A], B]
```

Implementation: `F.Flow2(Local[A](f), Map[D](g))`

### Local and LocalReaderK

`Local` adapts the reader environment before executing an effect:

```go
// effect/dependencies.go
func Local[A, C1, C2 any](acc Reader[C1, C2]) Kleisli[C1, Effect[C2, A], A]
func LocalReaderK[A, C1, C2 any](f reader.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A]
func LocalIOK[A, C1, C2 any](f io.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A]
func LocalResultK[A, C1, C2 any](f result.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A]
func LocalThunkK[A, C1, C2 any](f thunk.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A]
func LocalEffectK[A, C1, C2 any](f Kleisli[C2, C2, C1]) func(Effect[C1, A]) Effect[C2, A]
```

### Use Cases

**Adapter pattern**: You have `Effect[DBConfig, User]` but your app provides `AppConfig`. Use `Promap` or `Local`:

```go
type AppConfig struct { DB DBConfig; Cache CacheConfig }

extractDB := func(app AppConfig) DBConfig { return app.DB }

// Adapt Effect[DBConfig, User] to work with AppConfig
adapted := effect.Local[User](extractDB)(getUserEffect)
// adapted: Effect[AppConfig, User]
```

**Context transformation**: `LocalReaderK` is pure-function based. `LocalIOK`, `LocalResultK`, `LocalThunkK`, `LocalEffectK` allow the transformation itself to perform IO, fail, etc.

---

## 5. DI Architecture with `di` Package

### Core Types

```go
// di/token.go
type InjectionToken[T any] interface {
    Identity() Dependency[T]             // required, eager
    Option()   Dependency[Option[T]]     // optional, eager
    IOEither() Dependency[IOResult[T]]   // required, lazy (memoized singleton)
    IOOption() Dependency[IOOption[T]]   // optional, lazy (memoized singleton)
}

type MultiInjectionToken[T any] interface {
    Container() InjectionToken[[]T]  // resolve all implementations
    Item()      InjectionToken[T]    // provide one implementation
}

// Create tokens
func MakeToken[T any](name string) InjectionToken[T]
func MakeMultiToken[T any](name string) MultiInjectionToken[T]
```

### Providers

```go
// di/provider.go - 0 dependencies
func MakeProvider0[R any](token InjectionToken[R], fct IOResult[R]) DIE.Provider
func ConstProvider[R any](token InjectionToken[R], value R) DIE.Provider

// di/gen.go - 1 to N dependencies (generated, up to 15)
func MakeProvider1[T1, R any](
    token InjectionToken[R],
    dep1 Dependency[T1],
    fct func(T1) IOResult[R],
) DIE.Provider

func MakeProvider2[T1, T2, R any](
    token InjectionToken[R],
    dep1 Dependency[T1],
    dep2 Dependency[T2],
    fct func(T1) func(T2) IOResult[R],
) DIE.Provider

// ... up to MakeProvider15
```

**Note**: Factory functions for MakeProvider2+ are curried: `func(T1) func(T2) IOResult[R]`.

### Tokens with Defaults

```go
func MakeTokenWithDefault0[R any](name string, fct IOResult[R]) InjectionToken[R]

func MakeTokenWithDefault1[T1, R any](
    name string,
    dep1 Dependency[T1],
    fct func(T1) IOResult[R],
) InjectionToken[R]

// ... up to MakeTokenWithDefault15
```

These create tokens that have a built-in provider factory, so they resolve even without explicit provider registration.

### Injector and Resolution

```go
// di/erasure/injector.go
func MakeInjector(providers []Provider) InjectableFactory

// di/injector.go - type-safe resolution
func Resolve[T any](token InjectionToken[T]) RIOR.ReaderIOResult[DIE.InjectableFactory, T]
```

### Application Entry Point

```go
// di/app.go
var InjMain = MakeToken[any]("APP")
var Main = Resolve(InjMain)

// RunMain: create injector from providers, resolve InjMain, run it
var RunMain = F.Flow3(
    DIE.MakeInjector,
    Main,
    IOR.Fold(IO.Of[error], F.Constant1[any](IO.Of[error](nil))),
)
```

### Dependency Resolution Modes

Each `InjectionToken[T]` supports four resolution modes via methods:

| Method | Type | Behavior |
|---|---|---|
| `Identity()` | `Dependency[T]` | Required, eager. Fails if not provided. |
| `Option()` | `Dependency[Option[T]]` | Optional, eager. Returns `None` if missing. |
| `IOEither()` | `Dependency[IOResult[T]]` | Required, lazy. Memoized singleton. |
| `IOOption()` | `Dependency[IOOption[T]]` | Optional, lazy. Memoized singleton. |

### Full Working Example

```go
package main

import (
    "fmt"
    "github.com/IBM/fp-go/v2/di"
    DIE "github.com/IBM/fp-go/v2/di/erasure"
    "github.com/IBM/fp-go/v2/ioresult"
)

// Define tokens
var (
    TokenDB     = di.MakeToken[*DBService]("DBService")
    TokenCache  = di.MakeToken[*CacheService]("CacheService")
    TokenApp    = di.MakeToken[*App]("App")
)

// Providers
func dbProvider() DIE.Provider {
    return di.MakeProvider0(TokenDB, ioresult.Of(&DBService{url: "postgres://..."}))
}

func cacheProvider() DIE.Provider {
    return di.MakeProvider0(TokenCache, ioresult.Of(&CacheService{ttl: 300}))
}

func appProvider() DIE.Provider {
    return di.MakeProvider2(
        di.InjMain,              // register as main
        TokenDB.Identity(),      // required dependency
        TokenCache.Identity(),   // required dependency
        func(db *DBService) func(cache *CacheService) di.IOResult[any] {
            return func(cache *CacheService) di.IOResult[any] {
                return ioresult.Of[any](&App{db: db, cache: cache})
            }
        },
    )
}

func main() {
    err := di.RunMain([]DIE.Provider{
        dbProvider(),
        cacheProvider(),
        appProvider(),
    })()
    if err != nil {
        fmt.Println("Error:", err)
    }
}

// For testing: swap providers
func testProviders() []DIE.Provider {
    return []DIE.Provider{
        di.ConstProvider(TokenDB, &DBService{url: "sqlite://test"}),
        di.ConstProvider(TokenCache, &CacheService{ttl: 0}),
        appProvider(),
    }
}
```

---

## 6. Circuit Breaker + Retry

### Retry Policies

```go
// retry/retry.go
type RetryStatus struct {
    IterNumber      uint
    CumulativeDelay time.Duration
    PreviousDelay   Option[time.Duration]
}

type RetryPolicy = func(RetryStatus) Option[time.Duration]

var DefaultRetryStatus = RetryStatus{IterNumber: 0, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()}

// Policy constructors
func LimitRetries(i uint) RetryPolicy
func ConstantDelay(delay time.Duration) RetryPolicy
func ExponentialBackoff(delay time.Duration) RetryPolicy
func CapDelay(maxDelay time.Duration, policy RetryPolicy) RetryPolicy

// Combine policies via Monoid
// Both must return Some for retry to continue; larger delay wins
var Monoid = M.FunctionMonoid[RetryStatus](...)

func ApplyPolicy(policy RetryPolicy, status RetryStatus) RetryStatus
func Always[A any](a A) func(RetryStatus) A
```

**Composing policies**:

```go
import M "github.com/IBM/fp-go/v2/monoid"

// Retry up to 5 times with exponential backoff capped at 10s
policy := M.ConcatAll(retry.Monoid)([]retry.RetryPolicy{
    retry.LimitRetries(5),
    retry.CapDelay(10*time.Second, retry.ExponentialBackoff(100*time.Millisecond)),
})
```

### Retrying

Available in every monad package (`effect`, `io`, `ioresult`, `readerio`, `context/readerioresult`, etc.):

```go
// effect/retry.go
func Retrying[C, A any](
    policy retry.RetryPolicy,
    action Kleisli[C, retry.RetryStatus, A],  // func(RetryStatus) Effect[C, A]
    check  Predicate[Result[A]],              // should we retry?
) Effect[C, A]

// io/retry.go
func Retrying[A any](
    policy retry.RetryPolicy,
    action func(retry.RetryStatus) IO[A],
    check  func(A) bool,
) IO[A]
```

### Circuit Breaker States

```go
// circuitbreaker/types.go
type BreakerState = Either[openState, ClosedState]

// ClosedState interface -- two implementations
type ClosedState interface {
    Empty() ClosedState
    AddError(time.Time) ClosedState
    AddSuccess(time.Time) ClosedState
    Check(time.Time) Option[ClosedState]
}

// Counter-based: opens after N consecutive failures
func MakeClosedStateCounter(maxFailures uint) ClosedState

// History-based: opens after N failures within a time window
func MakeClosedStateHistory(timeWindow time.Duration, maxFailures uint) ClosedState
```

State transitions:
- **Closed -> Open**: failure threshold exceeded (`Check` returns `None`)
- **Open -> Half-Open**: reset time exceeded (canary allowed)
- **Half-Open -> Closed**: canary succeeds
- **Half-Open -> Open**: canary fails (extended backoff)

### Metrics

```go
type Metrics interface {
    Accept(time.Time) IO[Void]
    Reject(time.Time) IO[Void]
    Open(time.Time) IO[Void]
    Close(time.Time) IO[Void]
    Canary(time.Time) IO[Void]
}

func MakeMetricsFromLogger(name string, logger *log.Logger) Metrics
func MakeVoidMetrics() Metrics  // no-op, for testing/benchmarks
```

### Creating a Circuit Breaker

```go
func MakeCircuitBreaker[E, T, HKTT, HKTOP, HKTHKTT any](
    left func(E) HKTT,
    chainFirstIOK func(io.Kleisli[T, BreakerState]) func(HKTT) HKTT,
    chainFirstLeftIOK func(io.Kleisli[E, BreakerState]) func(HKTT) HKTT,
    chainFirstIOK2 func(io.Kleisli[Either[E, T], Void]) func(HKTT) HKTT,
    fromIO func(IO[func(HKTT) HKTT]) HKTOP,
    flap func(HKTT) func(HKTOP) HKTHKTT,
    flatten func(HKTHKTT) HKTT,
    currentTime IO[time.Time],
    closedState ClosedState,
    makeError Reader[time.Time, E],
    checkError option.Kleisli[E, E],
    policy retry.RetryPolicy,
    metrics Metrics,
) State[Pair[IORef[BreakerState], HKTT], HKTT]

func MakeSingletonBreaker[HKTT any](
    cb State[Pair[IORef[BreakerState], HKTT], HKTT],
    closedState ClosedState,
) func(HKTT) HKTT

var MakeClosedIORef = F.Flow2(createClosedCircuit, ioref.MakeIORef)
```

`MakeCircuitBreaker` is highly generic (works with any HKT). The typical usage is via a pre-built wrapper for your specific monad. `MakeSingletonBreaker` wraps it into a simple `func(HKTT) HKTT` operator.

### Combined Retry + Circuit Breaker Pattern

```go
// 1. Define retry policy
policy := M.ConcatAll(retry.Monoid)([]retry.RetryPolicy{
    retry.LimitRetries(3),
    retry.CapDelay(5*time.Second, retry.ExponentialBackoff(200*time.Millisecond)),
})

// 2. Create circuit breaker (singleton, wraps your monad)
breaker := MakeSingletonBreaker(
    MakeCircuitBreaker(/* ... monad-specific args ... */),
    MakeClosedStateCounter(5),
)

// 3. Combine: circuit breaker wraps the retried operation
protectedCall := breaker(
    effect.Retrying(policy, action, check),
)
```

---

## 7. Tail Recursion (Trampolining)

### The Trampoline Type

```go
// tailrec/types.go
type Trampoline[B, L any] struct {
    Land    L       // final result (valid when Landed == true)
    Bounce  B       // intermediate state (valid when Landed == false)
    Landed  bool    // true = done, false = continue
}

// tailrec/trampoline.go
func Bounce[L, B any](b B) Trampoline[B, L]   // continue with new state
func Land[B, L any](l L) Trampoline[B, L]     // terminate with result
```

Note the type parameter order: `Bounce[L, B]` takes `L` first (the unused type), `B` second (the value being set). Similarly `Land[B, L]` takes `B` first (unused), `L` second (result).

### TailRec in Monad Packages

Every monad provides a `TailRec` that executes a trampolined computation within that monad:

```go
// result/rec.go
func TailRec[A, B any](f Kleisli[A, Trampoline[A, B]]) Kleisli[A, B]

// io/rec.go
func TailRec[A, B any](f Kleisli[A, Trampoline[A, B]]) Kleisli[A, B]

// readerio/rec.go
func TailRec[R, A, B any](f Kleisli[R, A, Trampoline[A, B]]) Kleisli[R, A, B]

// context/readerioresult/rec.go
func TailRec[A, B any](f Kleisli[A, Trampoline[A, B]]) Kleisli[A, B]
```

### Example: Stack-Safe Factorial

```go
import (
    "github.com/IBM/fp-go/v2/tailrec"
    "github.com/IBM/fp-go/v2/result"
)

type FactState struct {
    N   int
    Acc int
}

factStep := func(s FactState) result.Result[tailrec.Trampoline[FactState, int]] {
    if s.N <= 1 {
        return result.Of(tailrec.Land[FactState](s.Acc))
    }
    return result.Of(tailrec.Bounce[int](FactState{N: s.N - 1, Acc: s.Acc * s.N}))
}

factorial := result.TailRec(factStep)
// factorial(FactState{N: 10000, Acc: 1}) -- won't overflow the stack
```

### When to Use

- Deep recursion (thousands+ frames) that would overflow Go's goroutine stack
- Recursive algorithms within monadic contexts (IO, Result, ReaderIO)
- Processing deeply nested data structures (trees, linked lists)
- Go's goroutine stacks grow dynamically but have limits; trampolining guarantees constant stack usage

---

## 8. State + StateIO

### State Monad (Pure)

```go
// state/types.go
type State[S, A any] = Reader[S, Pair[S, A]]  // func(S) Pair[S, A]
type Kleisli[S, A, B any] = Reader[A, State[S, B]]
type Operator[S, A, B any] = Kleisli[S, State[S, A], B]

// state/state.go - Core operations
func Of[S, A any](a A) State[S, A]              // pure value, state unchanged
func Get[S any]() State[S, S]                    // read current state
func Gets[FCT ~func(S) A, A, S any](f FCT) State[S, A]  // extract from state
func Put[S any]() State[S, Void]                 // replace state
func Modify[FCT ~func(S) S, S any](f FCT) State[S, Void] // transform state

func Map[S any, FCT ~func(A) B, A, B any](f FCT) Operator[S, A, B]
func Chain[S any, FCT ~func(A) State[S, B], A, B any](f FCT) Operator[S, A, B]
func Ap[B, S, A any](ga State[S, A]) Operator[S, func(A) B, B]
func Flatten[S, A any](mma State[S, State[S, A]]) State[S, A]

func Execute[A, S any](s S) func(State[S, A]) S   // run, return final state
func Evaluate[A, S any](s S) func(State[S, A]) A  // run, return value
func MonadChainFirst[S any, FCT ~func(A) State[S, B], A, B any](ma State[S, A], f FCT) State[S, A]
func ChainFirst[S any, FCT ~func(A) State[S, B], A, B any](f FCT) Operator[S, A, A]
func Flap[S, A, B any](a A) Operator[S, func(A) B, B]
```

### StateIO (State + IO)

```go
// stateio/type.go
type StateIO[S, A any] = Reader[S, IO[Pair[S, A]]]  // func(S) func() Pair[S, A]
type Kleisli[S, A, B any] = Reader[A, StateIO[S, B]]
type Operator[S, A, B any] = Reader[StateIO[S, A], StateIO[S, B]]

// stateio/state.go
func Of[S, A any](a A) StateIO[S, A]
func Map[S, A, B any](f func(A) B) Operator[S, A, B]
func Chain[S, A, B any](f Kleisli[S, A, B]) Operator[S, A, B]
func Ap[B, S, A any](fa StateIO[S, A]) Operator[S, func(A) B, B]
func FromIO[S, A any](fa IO[A]) StateIO[S, A]
func FromIOK[S, A, B any](f func(A) IO[B]) Kleisli[S, A, B]
```

### Use Cases

**Counter/Accumulator**:

```go
type Counter struct { Value int }

increment := state.Modify(func(c Counter) Counter {
    return Counter{Value: c.Value + 1}
})

getCount := state.Gets(func(c Counter) int { return c.Value })

program := F.Pipe3(
    increment,
    state.Chain(func(_ function.Void) state.State[Counter, function.Void] {
        return increment
    }),
    state.Chain(func(_ function.Void) state.State[Counter, int] {
        return getCount
    }),
)

result := state.Evaluate[int](Counter{Value: 0})(program) // 2
```

**StateIO for stateful IO**: threading mutable state through IO operations (file handles, connection pools, accumulators that also write logs).

---

## 9. Iterator Integration

### Two Iterator Models

**Stateless iterators** (`iterator/stateless`): Pure functional iterators using `Lazy[Option[Pair[Iterator[U], U]]]`. Immutable, safe for concurrent use, but create allocations per step.

```go
// iterator/stateless/types.go
type Iterator[U any] Lazy[Option[Pair[Iterator[U], U]]]
type Kleisli[A, B any] = Reader[A, Iterator[B]]
type Operator[A, B any] = Kleisli[Iterator[A], B]
```

**`iter.Seq` iterators** (`iterator/iter`): Wraps Go 1.23+ `iter.Seq[T]` and `iter.Seq2[K, V]` with functional operations. Lazy, native Go range-over-func compatible.

```go
// iterator/iter/types.go
type Seq[T any] = iter.Seq[T]
type Seq2[K, V any] = iter.Seq2[K, V]
type Kleisli[A, B any] = func(A) Seq[B]
type Operator[A, B any] = Kleisli[Seq[A], B]
```

### Stateless Iterator Operations

```go
// iterator/stateless/iterator.go
func Empty[U any]() Iterator[U]
func Of[U any](a U) Iterator[U]
func From[U any](data ...U) Iterator[U]
func FromArray[U any](as []U) Iterator[U]
func ToArray[U any](u Iterator[U]) []U

func Map[U, V any](f func(U) V) Operator[U, V]
func Chain[U, V any](f Kleisli[U, V]) Kleisli[Iterator[U], V]
func Filter[U any](f Predicate[U]) Operator[U, U]
func FilterMap[U, V any](f func(U) Option[V]) Operator[U, V]
func Flatten[U any](ma Iterator[Iterator[U]]) Iterator[U]
func Reduce[U, V any](f func(V, U) V, initial V) func(Iterator[U]) V
func Ap[V, U any](ma Iterator[U]) Operator[func(U) V, V]

func MakeBy[FCT ~func(int) U, U any](f FCT) Iterator[U]  // infinite
func Replicate[U any](a U) Iterator[U]                     // infinite
func Repeat[U any](n int, a U) Iterator[U]                 // n copies
func Count(start int) Iterator[int]                         // infinite counter

func Fold[U any](m M.Monoid[U]) func(Iterator[U]) U
func FoldMap[U, V any](m M.Monoid[V]) func(func(U) V) func(Iterator[U]) V
func ChainFirst[U, V any](f Kleisli[U, V]) Operator[U, U]
func FilterChain[U, V any](f func(U) Option[Iterator[V]]) Operator[U, V]
```

### iter.Seq Operations

```go
// iterator/iter/iter.go
func Of[A any](a A) Seq[A]
func Of2[K, A any](k K, a A) Seq2[K, A]
func From[A any](data ...A) Seq[A]
func Map[A, B any](f func(A) B) Operator[A, B]
func MapWithIndex[A, B any](f func(int, A) B) Operator[A, B]
func MapWithKey[K, A, B any](f func(K, A) B) Operator2[K, A, B]
func Filter[A any](pred func(A) bool) Operator[A, A]
func FilterWithIndex[A any](pred func(int, A) bool) Operator[A, A]
func FilterMap[A, B any](f option.Kleisli[A, B]) Operator[A, B]
func FilterMapWithIndex[A, B any](f func(int, A) Option[B]) Operator[A, B]
func Chain[A, B any](f func(A) Seq[B]) Operator[A, B]
func Flatten[A any](mma Seq[Seq[A]]) Seq[A]
func Reduce[A, B any](f func(B, A) B, initial B) func(Seq[A]) B
func ReduceWithIndex[A, B any](f func(int, B, A) B, initial B) func(Seq[A]) B
func ReduceWithKey[K, A, B any](f func(K, B, A) B, initial B) func(Seq2[K, A]) B
func Fold[A any](m M.Monoid[A]) func(Seq[A]) A
func FoldMap[A, B any](m M.Monoid[B]) func(func(A) B) func(Seq[A]) B
func Prepend[A any](head A) Operator[A, A]
func Append[A any](tail A) Operator[A, A]
func MonadZip[A, B any](fa Seq[A], fb Seq[B]) Seq2[A, B]
func Zip[A, B any](fb Seq[B]) func(Seq[A]) Seq2[A, B]
func ToSeqPair[A, B any](as Seq2[A, B]) Seq[Pair[A, B]]
func FromSeqPair[A, B any](as Seq[Pair[A, B]]) Seq2[A, B]
func MapToArray[A, B any](f func(A) B) func(Seq[A]) []B

// iterator/iter/take.go
func Take[U any](n int) Operator[U, U]
func TakeWhile[U any](p Predicate[U]) Operator[U, U]

// iterator/iter/cycle.go
func Cycle[U any](ma Seq[U]) Seq[U]     // infinite cycle

// iterator/iter/scan.go
func Scan[FCT ~func(V, U) V, U, V any](f FCT, initial V) Operator[U, V]

// iterator/iter/uniq.go
func Uniq[A any, K comparable](f func(A) K) Operator[A, A]
```

### Interop: iter.Seq <-> Stateless Iterator

The `iterator/iter` package also has `Iterator[T]` as an alias for `stateless.Iterator[T]`, allowing conversion between the two models.

### Do-Notation for iter.Seq

```go
// iterator/iter/bind.go
func Bind[S1, S2, T any](
    setter func(T) func(S1) S2,
    f func(S1) Seq[T],
) func(Seq[S1]) Seq[S2]

func BindTo[S1, T any](setter func(T) S1) func(Seq[T]) Seq[S1]

func BindL[S, T any](
    l Lens[S, T],
    f func(S) Seq[T],
) func(Seq[S]) Seq[S]
```

---

## 10. Performance Optimization

### Idiomatic vs Standard Packages

fp-go provides `idiomatic/` variants for performance-critical code. The key difference is representation:

| Aspect | Standard | Idiomatic |
|---|---|---|
| Either representation | Tagged struct (24 bytes) | Go tuple `(A, error)` (16 bytes) |
| Stack allocation | 24 bytes per Either | 16 bytes per tuple |
| ChainFirst | ~87 ns | ~2.7 ns (32x faster) |
| Pipeline overhead | Higher (struct wrapping) | Lower (native tuples) |
| Do-notation | Yes | Limited |
| Flatten/Swap | Yes | Not available |
| Custom error types | `Either[E, A]` for any E | Fixed to `error` |

### When to Use Each

**Use standard** (default choice):
- When you need `Either[E, A]` with custom error types (not just `error`)
- When you need Do-notation (`Bind`, `Let`, `ApS`)
- When you need `Flatten`, `Swap`, or other structural operations
- For business logic where clarity matters more than nanoseconds

**Use idiomatic**:
- Hot paths with millions of operations per second
- Simple pipelines: `Map`, `Chain`, `ChainFirst`
- When error type is always `error`
- Performance-sensitive middleware or data processing

### Import Pattern

```go
// Standard
import "github.com/IBM/fp-go/v2/ioresult"
import "github.com/IBM/fp-go/v2/context/readerioresult"

// Idiomatic equivalents
import "github.com/IBM/fp-go/v2/idiomatic/ioresult"
import "github.com/IBM/fp-go/v2/idiomatic/context/readerresult"
```

### Allocation Analysis

Standard `Either[error, A]`:
```
Either struct { value any; isRight bool }  // 8 + 8 + 8 = 24 bytes
```

Idiomatic result:
```
(A, error)  // Go multi-return: 8 + 8 = 16 bytes, stack-allocated
```

For pipelines with many intermediate Either values, idiomatic saves 33% memory per step and avoids interface boxing.

---

## 11. Builder Pattern

### Generic Builder Interface

```go
// builder/builder.go
type Builder[T any] interface {
    Build() Result[T]
}
```

Any struct implementing `Build() Result[T]` can participate in the builder pattern with optics integration.

### BuilderPrism

```go
// builder/prism.go
func BuilderPrism[T any, B Builder[T]](creator func(T) B) Prism[B, T]
```

Creates a `Prism` bidirectional conversion between a builder and its product:
- **Extract** (builder -> product): calls `Build()`, converts `Result` to `Option`
- **Construct** (product -> builder): uses the provided `creator` function

```go
type PersonBuilder struct {
    name string
    age  int
}

func (b PersonBuilder) Build() result.Result[Person] {
    if b.name == "" {
        return result.Error[Person](errors.New("name required"))
    }
    return result.Of(Person{Name: b.name, Age: b.age})
}

prism := builder.BuilderPrism(func(p Person) PersonBuilder {
    return PersonBuilder{name: p.Name, age: p.Age}
})

// prism.GetOption(PersonBuilder{name: "Alice", age: 30})
// => Some(Person{Name: "Alice", Age: 30})

// prism.ReverseGet(Person{Name: "Bob", Age: 25})
// => PersonBuilder{name: "Bob", age: 25}
```

### Endomorphism as Builder

The `endomorphism` package provides an alternative builder pattern using function composition:

```go
// endomorphism/builder.go
func Build[A any](e Endomorphism[A]) A     // apply to zero value
func ConcatAll[T any](es []Endomorphism[T]) Endomorphism[T]  // compose (right-to-left)
func Reduce[T any](es []Endomorphism[T]) T  // apply left-to-right from zero value
```

Example:

```go
type Config struct {
    Host    string
    Port    int
    Debug   bool
}

withHost := func(h string) endomorphism.Endomorphism[Config] {
    return func(c Config) Config { c.Host = h; return c }
}
withPort := func(p int) endomorphism.Endomorphism[Config] {
    return func(c Config) Config { c.Port = p; return c }
}
withDebug := func(d bool) endomorphism.Endomorphism[Config] {
    return func(c Config) Config { c.Debug = d; return c }
}

config := endomorphism.Reduce([]endomorphism.Endomorphism[Config]{
    withHost("localhost"),
    withPort(8080),
    withDebug(true),
})
// Config{Host: "localhost", Port: 8080, Debug: true}
```

---

## 12. Endomorphism

### Core Types

```go
// endomorphism/types.go
type Endomorphism[A any] = func(A) A
type Kleisli[A any] = func(A) Endomorphism[A]
type Operator[A any] = Endomorphism[Endomorphism[A]]
```

### Composition: Right-to-Left vs Left-to-Right

```go
// endomorphism/endo.go

// RIGHT-TO-LEFT (mathematical composition: f . g)
func MonadCompose[A any](f, g Endomorphism[A]) Endomorphism[A]
// MonadCompose(double, increment)(5) = double(increment(5)) = 12

func Compose[A any](g Endomorphism[A]) Operator[A]
// Compose(increment)(double)(5) = double(increment(5)) = 12

// LEFT-TO-RIGHT (pipeline order)
func MonadChain[A any](ma, f Endomorphism[A]) Endomorphism[A]
// MonadChain(double, increment)(5) = increment(double(5)) = 11

func Chain[A any](f Endomorphism[A]) Operator[A]
// Chain(increment)(double)(5) = increment(double(5)) = 11

// For effect, preserving first result
func MonadChainFirst[A any](ma, f Endomorphism[A]) Endomorphism[A]
func ChainFirst[A any](f Endomorphism[A]) Operator[A]

// Map = Compose (same thing for endomorphisms)
func MonadMap[A any](f, ma Endomorphism[A]) Endomorphism[A]
func Map[A any](f Endomorphism[A]) Operator[A]

// Ap = Compose (same thing for endomorphisms)
func MonadAp[A any](fab, fa Endomorphism[A]) Endomorphism[A]
func Ap[A any](fa Endomorphism[A]) Operator[A]

func Flatten[A any](mma Endomorphism[Endomorphism[A]]) Endomorphism[A]
func Join[A any](f Kleisli[A]) Endomorphism[A]  // W combinator: f(a)(a)
```

### Read

```go
func Read[A any](a A) func(Endomorphism[A]) A
```

Captures a value and returns a function that applies endomorphisms to it. Useful for evaluating multiple transformations on the same input:

```go
applyTo5 := endomorphism.Read(5)
applyTo5(N.Mul(2))   // 10
applyTo5(N.Add(1))   // 6
```

### Monoid/Semigroup

```go
func Semigroup[A any]() S.Semigroup[Endomorphism[A]]
func Monoid[A any]() M.Monoid[Endomorphism[A]]
func Identity[A any]() Endomorphism[A]
func Of[F ~func(A) A, A any](f F) Endomorphism[A]
```

The monoid uses **right-to-left** composition for `Concat` and `Identity` as the empty element:

```go
m := endomorphism.Monoid[int]()
combined := M.ConcatAll(m)([]endomorphism.Endomorphism[int]{
    N.Mul(2),
    N.Add(1),
    func(x int) int { return x * x },
})
// combined(5) = Mul2(Add1(Square(5))) = Mul2(Add1(25)) = Mul2(26) = 52
```

### FromSemigroup

```go
func FromSemigroup[A any](s S.Semigroup[A]) Kleisli[A]
```

Converts a semigroup into a Kleisli arrow. Follows "data last": `FromSemigroup(addSg)(5)` creates `func(x) { return x + 5 }`.

### Builder Functions

```go
func Build[A any](e Endomorphism[A]) A
func ConcatAll[T any](es []Endomorphism[T]) Endomorphism[T]  // right-to-left
func Reduce[T any](es []Endomorphism[T]) T                    // left-to-right from zero
```

### Use Cases

**Middleware chains**:
```go
type Handler = func(http.Request) http.Response
type Middleware = endomorphism.Endomorphism[Handler]

logging := func(h Handler) Handler { /* wrap with logging */ }
auth := func(h Handler) Handler { /* wrap with auth */ }
cors := func(h Handler) Handler { /* wrap with CORS */ }

// Compose middlewares (right-to-left: cors wraps auth wraps logging)
stack := endomorphism.ConcatAll([]Middleware{cors, auth, logging})
finalHandler := stack(baseHandler)
```

**Transformation pipelines**:
```go
// Left-to-right pipeline execution
result := endomorphism.Reduce([]endomorphism.Endomorphism[string]{
    strings.TrimSpace,
    strings.ToLower,
    func(s string) string { return strings.ReplaceAll(s, " ", "-") },
})
// "" -> "" -> "" -> "" (from zero value)
// Use Read or Build for specific starting values
```

---

## 13. Advanced Array Operations

### Indexed Operations

```go
// array/array.go
func FilterMapWithIndex[A, B any](f func(int, A) Option[B]) Operator[A, B]
func ReduceWithIndex[A, B any](f func(int, B, A) B, initial B) func([]A) B
```

`FilterMapWithIndex` combines filtering and mapping with access to the element index. `ReduceWithIndex` is a fold that also passes the index.

### Sorting

```go
// array/sort.go
func SortBy[T any](ord []O.Ord[T]) Operator[T, T]
func SortByKey[K, T any](ord O.Ord[K], f func(T) K) Operator[T, T]
```

`SortBy` accepts multiple `Ord` instances for multi-criteria stable sorting. `SortByKey` sorts by extracting a comparable key.

```go
import (
    A "github.com/IBM/fp-go/v2/array"
    "github.com/IBM/fp-go/v2/ord"
)

type Person struct { Name string; Age int }

byAge := ord.Contramap(func(p Person) int { return p.Age })(ord.FromStrictCompare[int]())
byName := ord.Contramap(func(p Person) string { return p.Name })(ord.FromStrictCompare[string]())

sorted := A.SortBy([]ord.Ord[Person]{byAge, byName})(people)
```

### Intercalate

```go
func Intercalate[A any](m M.Monoid[A]) func(A) func([]A) A
```

Inserts a separator between elements and concatenates using the monoid:

```go
import (
    A "github.com/IBM/fp-go/v2/array"
    S "github.com/IBM/fp-go/v2/string"
)

joined := A.Intercalate(S.Monoid)(", ")([]string{"a", "b", "c"})
// "a, b, c"
```

### Traverse and Sequence

```go
// array/traverse.go
func Traverse[A, B, HKTB, HKTAB, HKTRB any](
    fof pointed.OfType[[]B, HKTRB],
    fmap functor.MapType[[]B, func(B) []B, HKTRB, HKTAB],
    fap apply.ApType[HKTB, HKTRB, HKTAB],
    f func(A) HKTB,
) func([]A) HKTRB

func TraverseWithIndex[A, B, HKTB, HKTAB, HKTRB any](
    fof pointed.OfType[[]B, HKTRB],
    fmap func(func([]B) func(B) []B) func(HKTRB) HKTAB,
    fap func(HKTB) func(HKTAB) HKTRB,
    f func(int, A) HKTB,
) func([]A) HKTRB
```

Traverse maps each array element to an effect, then sequences the effects into one effect containing an array. The HKT type parameters are necessary because Go lacks higher-kinded types.

For `Option`:
```go
parseAll := array.Traverse(
    option.Of[[]int],
    option.Map[[]int, func(int) []int],
    option.Ap[[]int, int],
    func(s string) option.Option[int] {
        n, err := strconv.Atoi(s)
        if err != nil { return option.None[int]() }
        return option.Some(n)
    },
)

parseAll([]string{"1", "2", "3"}) // Some([1, 2, 3])
parseAll([]string{"1", "x", "3"}) // None
```

### Other Array Operations

```go
// Construction
func Of[A any](a ...A) []A
func MakeBy[A any](n int, f func(int) A) []A
func Replicate[A any](n int, a A) []A

// Transformation
func Map[A, B any](f func(A) B) Operator[A, B]
func Chain[A, B any](f func(A) []B) Operator[A, B]
func Filter[A any](f Predicate[A]) Operator[A, A]
func FilterMap[A, B any](f func(A) Option[B]) Operator[A, B]
func Flatten[A any](mma [][]A) []A
func Intersperse[A any](middle A) Operator[A, A]

// Reduction
func Reduce[A, B any](f func(B, A) B, initial B) func([]A) B
func Fold[A any](m M.Monoid[A]) func([]A) A
func FoldMap[A, B any](m M.Monoid[B]) func(func(A) B) func([]A) B

// Lookup and access
func Lookup[A any](idx int) func([]A) Option[A]
func Head[A any](as []A) Option[A]
func Last[A any](as []A) Option[A]
func Tail[A any](as []A) Option[[]A]
func Init[A any](as []A) Option[[]A]

// Size and predicates
func Size[A any](as []A) int
func IsEmpty[A any](as []A) bool
func IsNonEmpty[A any](as []A) bool

// Set operations (with Eq)
func Uniq[A any](eq Eq[A]) Operator[A, A]
func Difference[A any](eq Eq[A]) func([]A) Operator[A, A]
func Intersection[A any](eq Eq[A]) func([]A) Operator[A, A]
func Union[A any](eq Eq[A]) func([]A) Operator[A, A]
```

---

## Quick Reference: Import Aliases Convention

fp-go code follows a standard import alias convention. Learn these to read any fp-go source:

```go
import (
    A  "github.com/IBM/fp-go/v2/array"
    E  "github.com/IBM/fp-go/v2/either"
    F  "github.com/IBM/fp-go/v2/function"
    IO "github.com/IBM/fp-go/v2/io"
    IOE "github.com/IBM/fp-go/v2/ioeither"
    IOR "github.com/IBM/fp-go/v2/ioresult"
    L  "github.com/IBM/fp-go/v2/lazy"
    M  "github.com/IBM/fp-go/v2/monoid"
    N  "github.com/IBM/fp-go/v2/number"
    O  "github.com/IBM/fp-go/v2/option"
    R  "github.com/IBM/fp-go/v2/result"
    S  "github.com/IBM/fp-go/v2/semigroup"
    RIOR "github.com/IBM/fp-go/v2/readerioresult"
    DIE "github.com/IBM/fp-go/v2/di/erasure"
)
```

---

## Pattern: Composing Multiple Concerns

A complete example combining several advanced techniques:

```go
// Effect with retry, using do-notation to build result
func fetchUserProfile[C any](userID string) effect.Effect[C, UserProfile] {
    policy := M.ConcatAll(retry.Monoid)([]retry.RetryPolicy{
        retry.LimitRetries(3),
        retry.CapDelay(5*time.Second, retry.ExponentialBackoff(200*time.Millisecond)),
    })

    return effect.Retrying[C, UserProfile](
        policy,
        func(status retry.RetryStatus) effect.Effect[C, UserProfile] {
            return F.Pipe4(
                R.Do(Initial{}),
                R.Bind(SetUser, func(_ Initial) R.Result[User] {
                    return lookupUser(userID)
                }),
                R.Bind(SetPrefs, func(s WithUser) R.Result[Prefs] {
                    return lookupPrefs(s.User.ID)
                }),
                R.Let(SetProfile, func(s WithPrefs) UserProfile {
                    return UserProfile{User: s.User, Prefs: s.Prefs}
                }),
                R.Map(func(s WithProfile) UserProfile { return s.Profile }),
                // Lift Result into Effect
                effect.FromResult[C, UserProfile],
            )
        },
        func(r result.Result[UserProfile]) bool {
            return either.IsLeft(r) // retry on any error
        },
    )
}
```

This combines:
1. **Retry** with exponential backoff
2. **Do-notation** for sequential, dependent computations
3. **Effect** for the full reader+IO+error stack
4. **Kleisli composition** via `Bind` chains
