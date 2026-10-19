# Effect Tier 2: Implementation Plan

This document plans the second set of additions to the `effect` package, modelled on the runtime features of effect-ts v4. Tier 1, merged in #278, exposed existing combinators (`Bracket`, `TraverseArrayParN`, `WithTimeout`, `AnnotateLogs`, ...). Tier 2 adds what fp-go cannot express today:

1. **Scope**: finalizers for resources that outlive a single `Bracket`.
2. **Layers with `di`**: services built lazily on first use, memoized per scope and released when the scope closes, by extending the existing `di` package.
3. **Fiber, Fork and Race**: structured concurrency with no leaked goroutines.
4. **Error handling by type**: `errors.As` based catching and recognisable interruption.
5. **Schedule**: composable policies for `Retry` and `Repeat`.
6. **Clock**: an injectable time source, so that time-based code can be tested without sleeping.

All changes go into `v2` only and follow [AGENTS.md](../AGENTS.md): doc comments without Markdown, `Example` functions for every public function, tests by category, and signed, signed-off commits.

## Table of Contents

- [Design Principles](#design-principles)
- [Where Runtime Services Live](#where-runtime-services-live)
- [Phase 0: Prerequisites](#phase-0-prerequisites)
  - [0.1 Finalizers Must Survive Cancellation](#01-finalizers-must-survive-cancellation)
  - [0.2 Clock](#02-clock)
  - [0.3 Catching Errors by Type](#03-catching-errors-by-type)
- [Phase 1: Scope](#phase-1-scope)
- [Phase 2: Layers with `di`](#phase-2-layers-with-di)
- [Phase 3: Fiber, Fork and Race](#phase-3-fiber-fork-and-race)
- [Phase 4: Interruption](#phase-4-interruption)
- [Phase 5: Schedule](#phase-5-schedule)
- [Delivery Plan](#delivery-plan)
- [Testing Strategy](#testing-strategy)
- [Risks](#risks)
- [Decisions](#decisions)
- [Out of Scope](#out-of-scope)

## Design Principles

- **The `Effect` type does not change.** `Effect[C, A] = func(C) func(context.Context) func() Result[A]` stays as it is. Every Tier 2 feature is a function over this type, so all existing combinators keep working and nothing breaks.
- **Errors stay `error`.** effect-ts tracks the error type `E`. Go cannot express unions of error types, so fp-go keeps a single `error` channel and adds type-directed helpers on top of `errors.As` and `errors.Is`.
- **Port semantics, not syntax.** `Effect.gen`, `ServiceMap` and `Context.Tag` rely on generators and type-level sets that Go does not have. Do-notation with lenses and an explicit struct `C` remain the Go equivalents.
- **Panics are bugs, not errors.** Code built with fp-go must not panic, and the library does not recover panics, neither in combinators nor in goroutines it starts (`ApPar`, `TraverseArrayPar`, `TraverseArrayParN`, `Fork`, `Race`). A panic crashes the process, as it does in plain Go. This is a deliberate difference from effect-ts, which captures panics as defects in its `Cause`.
- **No general data structures.** Queue, PubSub, STM, Cache and Pool stay out of scope. Channels and existing packages cover them.
- **Implement the primitive in `context/readerioresult` first, then lift it.** This is the pattern Tier 1 used for `TraverseArrayParN`, `WithSemaphore` and `AnnotateLogs`. Features that need the dependency type `C` (Layer) live in `effect` only.
- **One concept per problem.** Where fp-go already has a model for something, Tier 2 extends it instead of adding a second one. Layers extend `di`, and Schedule bridges to `retry`.
- **Naming.** Use the effect-ts name when the concept is the same (`Scoped`, `Fork`, `Race`, `Retry`, `Repeat`, `AnnotateLogs`). Use a Go name when the Go concept differs: `CatchAs` / `CatchIs` match error types and values with `errors.As` / `errors.Is` instead of tags, and `WithSemaphore` takes a channel instead of a semaphore object.
- **Keep it simple.** Each runtime service is supplied in exactly one way. The clock, for example, comes only from the `context.Context`.
- **No new module dependencies.** The module depends only on `testify` and `urfave/cli`, and it should stay that way.

## Where Runtime Services Live

Scope and Clock are runtime services, and they need a home. There are two options:

| | In `C` (dependencies) | In `context.Context` |
|---|---|---|
| Compile-time checked | yes | no, falls back to a default |
| Works with existing combinators | no, every `C` would need a field | yes |
| Precedent in fp-go | dependencies such as repositories and clients | the logger (`logging.WithLogger`) and cancellation |
| effect-ts | `Scope` appears in `R` | n/a |

**Recommendation: use `context.Context`** for Scope and Clock. Both are infrastructure scoped to a run, like the logger and the deadline, not business dependencies. The package doc of `effect` ("do not store dependencies in the context") needs a sentence that makes this distinction explicit.

Services are the opposite case: they are business dependencies, so the `di` injector that builds them is passed as `C` (see [Phase 2](#phase-2-layers-with-di)). The scope it registers finalizers in still comes from the context.

The cost is that a missing scope is a runtime condition, not a compile error. Operations that need a scope (`AcquireRelease`, `Fork`) fail with `ErrNoScope` instead of guessing, and `RunScoped` installs a root scope for the common case (see [Phase 1](#phase-1-scope)).

## Phase 0: Prerequisites

These are small, independent changes that the later phases build on.

### 0.1 Finalizers Must Survive Cancellation

**Problem (verified, tracked in [#279](https://github.com/IBM/fp-go/issues/279)).** When the context is cancelled while `use` runs, `Bracket` runs `release` with the cancelled context. A context-aware release (anything built with `WithContext`, `Eitherize` or a `Chain` that checks the context) then fails immediately and **does not release the resource**:

```go
// use cancels the context; release is context-aware
res := effect.Bracket(acquire, use, release)(deps)(ctx)()
// result: Left(context canceled), released: false
```

This is the most common way a resource is released (timeouts, client disconnects, a sibling failing in `TraverseArrayPar`). Scope finalizers would have the same problem.

**Fix.** Run `release` with `context.WithoutCancel(ctx)`. This keeps the context values (logger, request ID, clock, scope) and drops cancellation and deadline. This corresponds to effect-ts, where finalizers are uninterruptible.

- Apply it in `context/readerioresult.Bracket` and `WithResource`, and in `context/readerreaderioresult.Bracket`. `effect.Bracket`, `WithResource` and `WithLock` inherit the fix.
- Add `Uninterruptible[A](ma) ReaderIOResult[A]` to `context/readerioresult` and `effect.Uninterruptible[C, A]`. These are the public form of the same idea.
- A finalizer that must be bounded uses `WithTimeout` inside the release explicitly.

**Acceptance.** A regression test cancels the context during `use` and asserts that a context-aware release runs. All existing `Bracket` tests stay green.

### 0.2 Clock

**Goal.** One time source for `Delay`, `Timer`, `WithTimeout`, `Retrying`, the circuit breaker and the new `Schedule`, so that tests can advance time instead of sleeping.

**Status today.** `Delay`, `Timer` and `Retrying` call `time.NewTimer` / `context.WithTimeout` directly. The circuit breaker already takes an `IO[time.Time]`. A `VirtualTimer` exists, but only in test files (`circuitbreaker/circuitbreaker_test.go`, `context/readerioresult/circuitbreaker_test.go`).

**New package `v2/clock`:**

```go
// Clock is a source of time.
type Clock interface {
    Now() time.Time
    // After returns a channel that receives the time once d has elapsed, and a
    // function that stops the timer and releases its resources.
    After(d time.Duration) (<-chan time.Time, func() bool)
}

func Real() Clock                                          // backed by package time
func WithClock(c Clock) pair.Kleisli[context.CancelFunc, context.Context, context.Context]
func FromContext(ctx context.Context) Clock                // falls back to Real()
func Now(ctx context.Context) time.Time                    // convenience
func Sleep(ctx context.Context, d time.Duration) error     // honours cancellation
```

**New package `v2/clock/testclock`:**

```go
func New(start time.Time) *TestClock
func (c *TestClock) Advance(d time.Duration)   // fires every timer that is due, in order
func (c *TestClock) Set(t time.Time)
func (c *TestClock) BlockUntil(n int)          // waits until n timers are pending
```

`BlockUntil` is what makes tests deterministic: the test waits until the code under test is sleeping before it advances the clock.

**Refit existing code** to read the clock from the context: `context/readerioresult.Delay`, `Timer`, `Retrying` (its `delayWithCancel`) and `WithTimeout`.

`WithTimeout` keeps `context.WithTimeout` when the clock is the real clock, so its semantics stay exactly the same. With a test clock it uses `context.WithCancelCause` with cause `context.DeadlineExceeded`, which a timer on the test clock triggers. In that case `ctx.Err()` is `context.Canceled` and `context.Cause(ctx)` is `context.DeadlineExceeded`. Document this, and have the fp-go combinators report `context.Cause`, which they already do.

Add `clock.NowIO(clock) IO[time.Time]` so the circuit breaker constructors can be fed from a `Clock` without changing their signatures.

**Acceptance.** `Delay`, `Timer`, `WithTimeout` and `Retrying` with exponential backoff are each tested with `testclock`, without a single real sleep. The real-clock benchmarks of `Delay` and `Timer` show no regression.

### 0.3 Catching Errors by Type

effect-ts dispatches on error tags (`catchTag`, `catchTags`). The Go counterpart dispatches on error types and sentinel values.

```go
// context/readerioresult
func CatchAs[E error, A any](f func(E) ReaderIOResult[A]) Operator[A, A]
func CatchIs[A any](target error, f func(error) ReaderIOResult[A]) Operator[A, A]
func MapErrorAs[E error, A any](f func(E) error) Operator[A, A]
func TapErrorAs[E error, A, B any](f func(E) ReaderIOResult[B]) Operator[A, A]

// effect
func CatchAs[E error, C, A any](f func(E) Effect[C, A]) Operator[C, A, A]
func CatchIs[C, A any](target error, f Kleisli[C, error, A]) Operator[C, A, A]
func MapErrorAs[E error, C, A any](f func(E) error) Operator[C, A, A]
func TapErrorAs[E error, C, A, B any](f func(E) Effect[C, B]) Operator[C, A, A]
```

- `E` comes first, so call sites read `effect.CatchAs[*NotFoundError](handler)`. In most cases `C` and `A` are inferred from the handler.
- Errors that do not match pass through unchanged, including wrapped errors that `errors.As` cannot unwrap to `E`.
- These are thin wrappers over `ChainLeft` / `TapLeft`, so they need no new machinery.

**Acceptance.** Tests cover a matching type, a non-matching type, a wrapped error (`fmt.Errorf("%w")`), `errors.Join` and a pointer versus value receiver. There is one `Example` per function.

## Phase 1: Scope

**Goal.** Resources whose lifetime is not one lexical `Bracket`: a connection pool opened at startup and closed at shutdown, a temp directory for a whole test, a background goroutine tied to a request.

**Semantics** (aligned with effect-ts `Scope`):

- A scope collects finalizers. Closing it runs them in reverse order of registration (LIFO).
- Every finalizer runs, even if earlier ones fail. Their errors are combined with `errors.Join`.
- Finalizers run uninterruptibly (see [0.1](#01-finalizers-must-survive-cancellation)).
- A finalizer receives the outcome of the scope (`Result[Void]`: success, failure or interruption), like `Bracket`'s release.
- Closing is idempotent. Adding a finalizer to a closed scope runs it immediately.
- Scopes nest. A child scope registers its own close as a finalizer of its parent.

**API (`context/readerioresult`, then lifted to `effect`):**

```go
type Scope struct { /* mutex, finalizer stack, closed flag */ }

func NewScope() *Scope
func (s *Scope) Close(outcome Result[Void]) error
func WithScope(s *Scope) pair.Kleisli[context.CancelFunc, context.Context, context.Context]
func ScopeFromContext(ctx context.Context) option.Option[*Scope]

// AddFinalizer registers a finalizer in the ambient scope.
func AddFinalizer(f func(Result[Void]) ReaderIOResult[Void]) ReaderIOResult[Void]

// AcquireRelease acquires a resource and registers its release in the ambient scope.
func AcquireRelease[A, ANY any](acquire ReaderIOResult[A], release func(A, Result[Void]) ReaderIOResult[ANY]) ReaderIOResult[A]

// Scoped runs a computation in a fresh child scope and closes it afterwards.
func Scoped[A any](ma ReaderIOResult[A]) ReaderIOResult[A]
```

`effect` gets `AddFinalizer`, `AcquireRelease` and `Scoped` with the extra type parameter `C`, delegating as in Tier 1.

**Missing ambient scope.** `AcquireRelease` without a scope in the context fails with a sentinel error, `ErrNoScope`. It neither leaks the resource silently nor panics. `RunSync` gets a sibling, `RunScoped`, that installs a root scope and closes it when the run ends. That makes the common case work without ceremony.

**Shape of the scope type.** `Scope` is a pointer with methods because it is mutable runtime state, like `ioref.IORef`. The functional API (`AcquireRelease`, `Scoped`) is what user code calls. The methods exist for the application boundary (`main`, test setup).

**Acceptance:**

- Tests cover LIFO order, all finalizers running despite failures, the joined errors, closing twice, adding after close, nested scopes and the outcome passed to finalizers.
- A test cancels the context in the middle of a scoped computation and checks that all finalizers run.
- A test shows `Scoped` around `TraverseArrayPar`: each element acquires a resource, and all of them are released, including when one element fails.

## Phase 2: Layers with `di`

**Goal.** Build the services an application needs lazily, on first use, with resource safety and sharing. effect-ts calls this a Layer. fp-go already has the right model for it: the `di` package resolves dependencies by token, lazily, and memoizes them as singletons per injector. Instead of introducing a second dependency system, Phase 2 extends `di` with what effect-ts layers add: context-aware construction and finalizers bound to a Scope.

### What `di` Already Provides

- **Tokens** (`di.MakeToken[T]`) identify a service. Providers (`di.MakeProvider0` … `MakeProvider15`) declare the tokens they depend on and how to build the service.
- **Lazy resolution.** `erasure.MakeInjector` builds a service only when it is first requested, through `di.Resolve(token)` or as a dependency of another service. Services that are never requested are never built.
- **Memoization.** Each service is built at most once per injector (`sync.Map` + `lazy.Memoize`), also when it is requested concurrently. The result is memoized whether it is a success or a failure.
- **Cycle detection and missing-provider errors** at resolution time.
- **Overrides.** When providers are merged, the last provider registered for a token wins. Tests replace a service by appending a provider for its token.
- **Defaults, optional, lazy and multi-value dependencies** (`MakeTokenWithDefault`, `Option()`, `IOEither()`, `MakeMultiToken`).

What is missing: provider factories are `func(T1, …) IOResult[R]`. They see neither a `context.Context` nor a Scope, so they cannot log with the context logger, read the clock, or register the release of the resource they open.

### Extension 1: Scoped Providers (`di`)

A second family of provider constructors whose factories are context-aware computations:

```go
// di, generated alongside MakeProviderN (go run .. di --count 15)
func MakeScopedProvider0[R any](token InjectionToken[R], f ReaderIOResult[R]) DIE.Provider
func MakeScopedProvider1[T1, R any](token InjectionToken[R], d1 Dependency[T1], f func(T1) ReaderIOResult[R]) DIE.Provider
// … up to MakeScopedProvider15

func MakeScopedTokenWithDefault0[R any](name string, f ReaderIOResult[R]) InjectionToken[R]
// … and the remaining default variants
```

`ReaderIOResult` is `context/readerioresult.ReaderIOResult`. A factory typically uses `AcquireRelease` from Phase 1, so the release of the resource is registered in the scope of the injector:

```go
var DBToken = di.MakeToken[*sql.DB]("DB")

var DBProvider = di.MakeScopedProvider1(
    DBToken,
    ConfigToken.Identity(),
    func(cfg Config) readerioresult.ReaderIOResult[*sql.DB] {
        return readerioresult.AcquireRelease(openDB(cfg), closeDB)
    },
)
```

Existing `MakeProviderN` providers keep working unchanged in the same injector. They are simply built without a context and register no finalizers.

**Import graph.** `di` gains a dependency on `context/readerioresult`. That package does not import `di`, so there is no cycle. `effect` gains a dependency on `di`.

### Extension 2: Injector Bound to a Scope

```go
// di/erasure
func MakeScopedInjector(ctx context.Context, providers []Provider) InjectableFactory
```

- The injector builds every service with the **context of the scope it was created in**: `context.WithoutCancel(ctx)` plus the scope. Context values such as the logger, the clock and the scope are kept. The cancellation of the request that happened to trigger the build is not. Without this, a request that times out while a shared singleton is being built would leave a failed singleton behind for every later request.
- A caller that waits for a service that another caller is building honours its own context. If the waiting caller is cancelled, it stops waiting, but the build continues for the others.
- **Only successful builds are memoized.** Concurrent requests share one build in flight. If that build fails, all of them get the error, the entry is removed, and the next request starts a new build. A transient failure on first use (the database is briefly unavailable) therefore does not fail every later request until restart. Plain `MakeProviderN` providers keep the `di` semantics and memoize their result, failure included.
- Finalizers registered during a build go into the scope from `ctx`. Because a service finishes building after its dependencies, closing the scope (LIFO) releases dependents before their dependencies.
- Without a scope in `ctx`, the first scoped build fails with `ErrNoScope`, consistent with `AcquireRelease` and `Fork`.

`MakeInjector` stays as it is.

### Extension 3: Using the Injector from `effect`

Effects that use services take the injector as their dependency type. This keeps the rule from the package doc: dependencies live in `C`.

```go
// effect
type Injector = erasure.InjectableFactory

// Resolve returns the service for a token, building it on first use.
func Resolve[T any](token di.InjectionToken[T]) Effect[Injector, T]

// ProvideServices runs an effect with a scoped injector over the given providers.
// The scope is closed, and all services that were built are released, when the effect ends.
func ProvideServices[A any](providers ...erasure.Provider) func(Effect[Injector, A]) Thunk[A]

// AcquireServices creates a scoped injector for the lifetime of an application and
// returns it with the function that closes its scope.
func AcquireServices(providers ...erasure.Provider) Thunk[Pair[Injector, func(context.Context) error]]

// Preload resolves the given tokens eagerly, for example at startup.
func Preload(tokens ...erasure.Dependency) Effect[Injector, Void]
```

A typical service:

```go
func FindUser(id string) effect.Effect[effect.Injector, User] {
    return F.Pipe1(
        effect.Resolve(UserRepoToken),
        effect.Chain(func(repo UserRepo) effect.Effect[effect.Injector, User] {
            return effect.FromThunk[effect.Injector](repo.Find(id))
        }),
    )
}
```

Code that prefers a typed dependency struct converts at the boundary with the existing `LocalEffectK`, building the struct from `Resolve` calls with `Do` and `ApSL`. This resolves the services in the struct when the effect starts, not when they are first used.

### Application Lifetime

A server calls `AcquireServices` once in `main`, serves each request with `Provide[A](injector)`, and calls the close function on shutdown. Services are built on the first request that needs them and shared by all later requests. `Preload` is optional. It moves the cost of building, and the discovery of a missing provider, to startup.

**Acceptance:**

- A service graph config → DB → repository → service is built lazily: a service that is never resolved is never built (counter in its factory).
- A service shared by two others is built once, also when both are resolved concurrently.
- Closing the scope releases the services in reverse order of construction, and services that were never built are not released.
- A request that is cancelled while it waits for a shared service fails with its cancellation cause, and the build still completes for other requests.
- A scoped provider that fails on its first build is built again on the next request, and succeeds. Concurrent requests during the failing build share it and all get the error.
- A test overrides a provider by appending another one for the same token.
- A missing provider and a cycle are reported as errors by `Resolve`. `Preload` reports them at startup.
- `MakeProviderN` and `MakeScopedProviderN` providers are mixed in one injector.
- `Example_services` shows the full pattern: providers, `AcquireServices`, request handling and shutdown.

## Phase 3: Fiber, Fork and Race

**Goal.** Structured concurrency. Every goroutine that an effect starts has an owner, can be interrupted, and is awaited, so none outlive their scope.

### API

```go
// Fiber is a handle to an effect running in its own goroutine.
type Fiber[A any] struct { /* done channel, result, cancel func */ }

// context/readerioresult, lifted to effect
func Fork[A any](ma ReaderIOResult[A]) ReaderIOResult[*Fiber[A]]       // child of the ambient scope; ErrNoScope without one
func ForkDaemon[A any](ma ReaderIOResult[A]) ReaderIOResult[*Fiber[A]] // not tied to any scope
func ForkIn[A any](s *Scope) func(ReaderIOResult[A]) ReaderIOResult[*Fiber[A]]

func Join[A any](f *Fiber[A]) ReaderIOResult[A]               // waits, honours the joiner's context
func Await[A any](f *Fiber[A]) ReaderIOResult[Result[A]]      // waits, never fails
func Interrupt[A any](f *Fiber[A]) ReaderIOResult[Result[A]]  // cancels with ErrInterrupted, then waits
func Poll[A any](f *Fiber[A]) IO[option.Option[Result[A]]]    // non-blocking

func Race[A any](ms ...ReaderIOResult[A]) ReaderIOResult[A]      // first success wins; all fail → errors.Join
func RaceFirst[A any](ms ...ReaderIOResult[A]) ReaderIOResult[A] // first completion wins, success or failure
```

### Semantics

- **Context.** A forked effect runs with a child of the forking effect's context (`context.WithCancelCause`), so it inherits values, cancellation and deadline. `ForkDaemon` uses `context.WithoutCancel`.
- **Ownership.** `Fork` registers a finalizer in the ambient scope that interrupts the fiber and waits for it. When the scope closes, no fiber is left running. Without an ambient scope, `Fork` fails with `ErrNoScope` and does not start the goroutine. A fiber that must outlive every scope is started explicitly with `ForkDaemon`, and one that belongs to a specific scope with `ForkIn`.
- **Interruption.** Interruption is cancellation with the cause `ErrInterrupted`, a sentinel that wraps `context.Canceled`, so `errors.Is(err, context.Canceled)` still holds. Interrupting a fiber that has already completed returns its result.
- **Losers are always awaited.** `Race` and `RaceFirst` interrupt the losers and wait for them before returning, so the losers' finalizers have run when the race returns. effect-ts guarantees the same thing, and it is the main difference from a hand-written `select`.
- **Panics.** Fibers do not recover panics, following the [design principles](#design-principles). A panic in a forked effect crashes the process, as a panic in any goroutine does.

**Acceptance:**

- Goroutine leak tests: `runtime.NumGoroutine` before and after, with a short settle loop, because the module avoids a `goleak` dependency.
- `Race`: the losers' finalizers have run when it returns, all failures are joined, and the first success wins even if a failure completes earlier.
- `Scoped(Fork(...))`: closing the scope interrupts and awaits the fiber.
- `Fork` without an ambient scope fails with `ErrNoScope` and starts no goroutine.

## Phase 4: Interruption

effect-ts distinguishes three kinds of failure in its `Cause`: expected failures, defects and interruption. fp-go has no defects, because it does not recover panics (see [Design Principles](#design-principles)). It keeps `error` as the only failure channel and makes the remaining kinds recognisable:

| Kind | Representation | Predicate |
|---|---|---|
| Failure | any `error` | everything else |
| Interruption | `context.Canceled`, `context.DeadlineExceeded`, `ErrInterrupted` | `IsInterrupted(err)` |
| Parallel failures | `errors.Join(...)`, as already returned by Scope and Race | standard `errors.Is` / `errors.As` |

**API (`context/readerioresult`, lifted to `effect`):**

```go
// ErrInterrupted is the cancellation cause used by Interrupt. It wraps context.Canceled.
var ErrInterrupted error

func IsInterrupted(err error) bool

// OnInterrupt runs f if the computation fails because it was interrupted.
// f runs uninterruptibly, like a finalizer.
func OnInterrupt[A any](f func(error) ReaderIOResult[Void]) Operator[A, A]
```

`IsInterrupted` lets error handlers tell "the caller gave up" apart from "the operation failed". For example, a `Retry` schedule built with `WhileInput` should not retry an interrupted effect.

**Acceptance.** There are tests for `IsInterrupted` with direct, wrapped and joined errors, and for `OnInterrupt` on interruption, on an ordinary failure (not called) and on success (not called).

## Phase 5: Schedule

**Goal.** One composable description of "when to try again", which covers `Retry`, `Repeat` and polling. `retry.RetryPolicy` covers only part of it: it is stateless, delay-only, and does not see the value or error.

### Model

A schedule is a pure state machine. Each step consumes an input (the error for `Retry`, the value for `Repeat`) and the current time, and either stops or produces an output, a delay and the next step:

```go
type Decision[In, Out any] struct {
    Out   Out
    Delay time.Duration
    Next  Schedule[In, Out]
}

// Schedule is a pure, recursive step function. None means "stop".
type Schedule[In, Out any] func(in In, now time.Time) option.Option[Decision[In, Out]]
```

Encoding the state in the returned `Next` keeps schedules pure and reusable. The same value can drive any number of independent runs, without the `Lazy` factory a closure-with-state design would need.

### Constructors and Combinators (package `v2/schedule`)

| effect-ts | fp-go | Output |
|---|---|---|
| `Schedule.recurs(n)` | `Recurs[In](n)` | iteration count |
| `Schedule.spaced(d)` | `Spaced[In](d)` | iteration count |
| `Schedule.fixed(d)` | `Fixed[In](d)` (fixed rate, uses `now`) | iteration count |
| `Schedule.exponential(base, factor)` | `Exponential[In](base, factor)` | current delay |
| `Schedule.fibonacci(one)` | `Fibonacci[In](one)` | current delay |
| `Schedule.upTo(d)` / `during` | `UpTo(d)` | unchanged |
| `Schedule.jittered` | `Jittered(min, max, rnd)` | unchanged |
| `Schedule.whileInput` / `untilInput` | `WhileInput(pred)`, `UntilInput(pred)` | unchanged |
| `Schedule.whileOutput` / `untilOutput` | `WhileOutput(pred)`, `UntilOutput(pred)` | unchanged |
| `Schedule.union` | `Union(a, b)`: continue while either does, shorter delay | `Pair` of outputs |
| `Schedule.intersect` | `Intersect(a, b)`: continue while both do, longer delay | `Pair` of outputs |
| `Schedule.andThen` | `AndThen(a, b)` | `Either` of outputs |
| `Schedule.map` | `Map(f)` | `f(out)` |
| `Schedule.modifyDelay` | `ModifyDelay(f)` / `CapDelay(max)` | unchanged |

`Jittered` takes its randomness as a parameter (`func() float64`) instead of calling `math/rand`, so that it stays pure and testable.

Bridges to the existing retry package: `FromRetryPolicy(retry.RetryPolicy) Schedule[error, retry.RetryStatus]`. The existing `retry` functions stay as they are, and nothing is deprecated in this plan.

### Effect Integration

```go
// effect
func Retry[C, A, Out any](s schedule.Schedule[error, Out]) Operator[C, A, A]
func RetryOrElse[C, A, Out any](s schedule.Schedule[error, Out], orElse func(error, Out) Effect[C, A]) Operator[C, A, A]
func Repeat[C, A, Out any](s schedule.Schedule[A, Out]) Operator[C, A, Out]
func RepeatOrElse[C, A, Out any](s schedule.Schedule[A, Out], orElse func(error, option.Option[Out]) Effect[C, Out]) Operator[C, A, Out]
```

- Delays use the context clock ([0.2](#02-clock)), so a full retry sequence with exponential backoff runs in microseconds under `testclock`.
- Loops use `TailRec`, as `Retrying` does today, so long-running repeats do not grow the stack.
- The names `Retry` and `Repeat` are free in `effect`. The existing `Retrying` / `RetryingI` stay for policy-based code.

**Acceptance.**

- Each constructor and combinator has property-style tests: the delay sequence for N steps, the stop condition and the output.
- `Union` / `Intersect` laws hold on samples: commutativity of the delays, and `Recurs(0)` as the identity of `Union` on termination.
- `Retry` with `Exponential` + `Recurs(5)` + `WhileInput(isTransient)` is tested end to end with `testclock`.
- `Repeat` is tested with `Spaced` for a polling loop that stops on a predicate.

## Delivery Plan

Each row is one issue and one pull request, merged in order. The phases depend on each other as shown.

| # | Pull request | Depends on | Size |
|---|---|---|---|
| 1 | Finalizers survive cancellation (`Bracket` fix, `Uninterruptible`), #279 | none | S |
| 2 | `CatchAs`, `CatchIs`, `MapErrorAs`, `TapErrorAs` | none | S |
| 3 | `clock` and `clock/testclock`, refit `Delay`, `Timer`, `WithTimeout`, `Retrying` | none | M |
| 4 | `Scope`, `AcquireRelease`, `Scoped`, `RunScoped` | 1 | M |
| 5 | `di`: `MakeScopedProviderN` (generator), `MakeScopedInjector`; `effect`: `Resolve`, `ProvideServices`, `AcquireServices`, `Preload` | 4 | L |
| 6 | `ErrInterrupted`, `IsInterrupted`, `OnInterrupt` | 1 | S |
| 7 | `Fiber`, `Fork`, `Join`, `Interrupt`, `Race`, `RaceFirst` | 4, 6 | L |
| 8 | `schedule` package | 3 | M |
| 9 | `Retry`, `Repeat` and the `OrElse` variants in `effect` | 8 | S |
| 10 | Guide: "Building services with effect" in `docs/fp-go-guide`, updated fp-go skills (`skills/` at the repository root, then `go generate` for the copy under `gen/v2/data`) | 5, 7, 9 | M |

Pull requests 1–3 are independent and can be worked on in parallel. Pull request 5 (services with `di`) is the one with the highest value for users and should not wait for the concurrency work.

## Testing Strategy

- **Determinism first.** Everything time-related is tested with `testclock` after pull request 3. Existing real-time tests stay, but new tests must not sleep. The flaky `context/readerresult` timeout test is a reminder of why.
- **Leak checks.** Concurrency tests compare `runtime.NumGoroutine` before and after, with a bounded settle loop. A shared helper lives in an internal test package.
- **Race detector.** CI already runs `go test -race ./...`. New tests must be race-clean, which matters most for Scope, the scoped injector and Fiber, the shared mutable state in this plan.
- **Cancellation matrix.** For every new combinator, test a context that is cancelled before start, during the effect, during a finalizer, and after completion.
- **Benchmarks** for the hot paths: `Bracket` (before and after the fix), `Scoped` with N finalizers, `Build` on a cache hit, `Fork` + `Join`, and one `Schedule` step.
- **Examples.** Every public function has an `Example` with `// Output:` that follows the AGENTS.md rules (`F.Pipe`/`F.Flow`, `lazy.Of`, `A.From`).

## Risks

| Risk | Mitigation |
|---|---|
| Running finalizers without cancellation changes behaviour for users who relied on releases being skipped | Treat it as a bug fix and call it out in the release notes. A release that must stay cancellable can opt back in with an explicit `WithTimeout`. |
| Runtime services in `context.Context` fail at runtime instead of compile time | `RunScoped` installs a root scope. `AcquireRelease` and `Fork` fail with an explicit `ErrNoScope`. The only fallback, the real clock, is documented. |
| Services are built lazily, so a missing provider or a failing build shows up on the first request that needs the service, possibly late in production | `Preload` resolves chosen tokens at startup, and tests run the same provider graph with `ProvideServices`. Errors name the token and the dependency path, as `di` does today. |
| A panic in a forked effect or a parallel traversal crashes the process | This is intended: panics are bugs (see [Design Principles](#design-principles)). Doc comments of `Fork`, `Race` and the parallel traversals state it explicitly. |
| The size of the `effect` API grows | New concepts get their own packages (`clock`, `schedule`). `effect` only holds the lifted combinators, grouped in files by topic, as in Tier 1. |

## Decisions

Decisions taken while reviewing this plan:

| Question | Decision |
|---|---|
| Recover panics in `ApPar`, `TraverseArrayPar` and new goroutines? | No. Panics are bugs and must not occur; the library recovers none (see [Design Principles](#design-principles)). |
| `Fork` without an ambient scope? | Fails with `ErrNoScope`. Detached fibers are explicit, with `ForkDaemon`. |
| Clock also injectable through `C`? | No. The clock comes only from the `context.Context`, to keep a single way. |
| Naming? | effect-ts names for the same concepts, Go names where Go differs (see [Design Principles](#design-principles)). |
| Layers: separate typed model or `di`? | Extend `di`. Services are built lazily on first use and memoized per scope. |
| Cache failed builds? | Not for scoped providers. A failed build is forgotten, and the next request tries again. `MakeProviderN` keeps caching its result, failure included. |

## Out of Scope

- **Tier 3:** `Effect.fn` / tracing spans, `Config`, request batching (`RequestResolver`), and `Stream`.
- **Data structures:** Queue, PubSub, STM, Cache and Pool (see [Design Principles](#design-principles)).
- **Changes to v1** (the repository root outside `v2/`), which is frozen.
