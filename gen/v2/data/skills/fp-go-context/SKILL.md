---
name: fp-go-context
description: Use this skill when working with Go's context.Context in fp-go code (github.com/IBM/fp-go/v2/context/...). Trigger on mentions of context.Context in fp-go pipelines, request-scoped values, ctx.Value, context.WithValue, context keys, AskValue, WithValue, timeouts or deadlines (WithTimeout, WithDeadline, context.WithTimeout, defer cancel), cancellation, Local / LocalIOK / LocalIOResultK, WithContext / WithContextK, request-scoped loggers, converting func(ctx, ...) (T, error) functions to ReaderIOResult, or reviewing code that threads ctx by hand.
---

# fp-go Context Handling

## Core Principle

In imperative Go the context is threaded by hand: every function takes `ctx` first, values are read with `ctx.Value(k).(T)`, scopes are opened with `ctx, cancel := context.WithTimeout(...); defer cancel()`.

In fp-go the context is the **Reader environment**. A computation is a *description* `func(context.Context) …` that runs only when a context is supplied. Therefore:

1. **Pipelines never mention `ctx`.** Build them from operators; the context flows implicitly.
2. **Supply the context exactly once, at the edge** — HTTP handler (`r.Context()`), `main`, or a test (`t.Context()`).
3. **Read the context with operators** (`Ask`, `FromReader`, `AskValue`), never with `ctx.Value(k).(T)`.
4. **Scope the context with operators** (`WithValue`, `WithTimeout`, `WithDeadline`, `Local`), never with `context.With*` + `defer cancel()`.
5. **Only request-scoped data goes into the context.** Dependencies (DB, clients, config) belong in `Effect[Deps, A]`.

## Packages and Availability

| Package | Shape | Use when |
|---------|-------|----------|
| `context/readerioresult` (`RIO`) | `func(ctx) func() Result[A]` | default for services: IO + errors + context |
| `context/readerresult` (`RR`) | `func(ctx) Result[A]` | synchronous, can fail |
| `context/readerio` (`RIOC`) | `func(ctx) func() A` | IO that cannot fail |
| `context/statereaderioresult` (`SRIO`) | `func(S) func(ctx) func() Result[Pair[S, A]]` | explicit state + context |
| `idiomatic/context/readerresult` (`IRR`) | `func(ctx) (A, error)` | high-performance, native `(A, error)` |
| `context/reader` (`CR`) | `func(ctx) A` | plain building blocks for use *outside* pipelines |
| `effect` (`EF`) | `func(C) RIO.ReaderIOResult[A]` | typed deps `C` **plus** the runtime context |

| Operator | RIO | RR | RIOC | SRIO | IRR | CR |
|----------|:---:|:--:|:----:|:----:|:---:|:--:|
| `Ask()` | ✓ | ✓ | ✓ |  | ✓ | |
| `FromReader(f)` / `Asks(f)` | ✓ | ✓ | ✓ | `Asks` | ✓ | |
| `AskValue[V](key)` | ✓ | ✓ | ✓ | ✓ (`[S, V]`) | ✓ | ✓ |
| `WithValue[A](key, v)` | ✓ | ✓ | ✓ | ✓ (`[S, A]`) | ✓ | Kleisli form |
| `WithTimeout[A](d)` / `WithDeadline[A](t)` | ✓ | ✓ | ✓ | ✓ (`[S, A]`) | ✓ | |
| `Local[A](f)` | ✓ | ✓ | ✓ | ✓ | ✓ | |
| `LocalIOK` / `LocalIOResultK` | ✓ / ✓ | | ✓ / | | | |
| `WithContext` / `WithContextK` | ✓ | ✓ | | | ✓ | |
| `NopCancel(ctx)` | | | | | | ✓ |

`Local`'s argument is `func(ctx) Pair[context.CancelFunc, context.Context]` in the standard packages and `func(ctx) (context.Context, context.CancelFunc)` in `IRR`.

## 1. Running at the Edge

```go
// HTTP handler: the request context is supplied once
func handler(w http.ResponseWriter, r *http.Request) {
    res := handleRequest(r)(r.Context())()   // Result[Response] — ONE value
    resp, err := R.Unwrap(res)               // bridge back to (A, error)
    // ...
}

// Tests: use t.Context(), cancelled automatically when the test ends
res := pipeline(t.Context())()
```

Never capture a `ctx` in a closure and never store it in a struct; pass it only when running.

## 2. Bridging `func(ctx, …) (T, error)`

Existing context-first Go functions are lifted, not wrapped by hand:

```go
// func(context.Context, string) ([]byte, error)  ->  func(string) RIO.ReaderIOResult[[]byte]
fetch := RIO.Eitherize1(fetchBytes)

// idiomatic package: the shape is already native
fetchI := IRR.From1(fetchBytes)                       // func(string) IRR.ReaderResult[[]byte]

// and back, for APIs that expect the Go shape
fetchGo := RIO.Uneitherize1(fetch)                    // func(context.Context, string) ([]byte, error)
```

The lifted function receives the pipeline's (possibly scoped) context — timeouts and values applied with the operators below reach it automatically.

## 3. Reading the Context

| Need | Operator | Result |
|------|----------|--------|
| whole context | `RIO.Ask()` | `ReaderIOResult[context.Context]` |
| pure projection | `RIO.FromReader(f)` | `ReaderIOResult[A]` |
| one typed value | `RIO.AskValue[V](key)` | `ReaderIOResult[Option[V]]` |

`AskValue` **never panics and never fails**: `Some(v)` if the key holds a `V`, `None` if it is absent *or holds another type*. The caller decides what "missing" means.

### Keys and accessors

Use an unexported key type — never plain strings or exported types — and keep one read and one write accessor next to each key, so the rest of the code never touches the key:

```go
type ctxKey int

const (
    userKey ctxKey = iota
    requestIDKey
)

var errNoUser = errors.New("no authenticated user in context")

// optional value: None -> default
func requestID() RIO.ReaderIOResult[string] {
    return F.Pipe1(
        RIO.AskValue[string](requestIDKey),
        RIO.Map(O.GetOrElse(F.Constant("-"))),
    )
}

// required value: None -> error
func requireUser() RIO.ReaderIOResult[User] {
    return F.Pipe1(
        RIO.AskValue[User](userKey),
        RIO.Chain(RIO.FromOption[User](F.Constant(errNoUser))),
    )
}

// writer: scopes the value to the wrapped computation
func withUser[A any](u User) RIO.Operator[A, A] {
    return RIO.WithValue[A](userKey, u)
}
```

In `RR` / `IRR`, turn `None` into an error with `ChainOptionK`:

```go
F.Pipe1(RR.AskValue[User](userKey),  RR.ChainOptionK[O.Option[User], User](F.Constant(errNoUser))(F.Identity[O.Option[User]]))
F.Pipe1(IRR.AskValue[User](userKey), IRR.ChainOptionK[O.Option[User], User](F.Constant(errNoUser))(O.Unwrap[User]))
```

## 4. Scoping the Context

Scoping operators run the wrapped computation with a **derived** context. The caller's context is never modified, and the derived context's cancel function is **always** released when the computation completes — no `defer cancel()`, no leaks.

| Operator | Derives with | Notes |
|----------|--------------|-------|
| `WithValue[A](key, v)` | `context.WithValue` | inner values shadow outer ones for the same key |
| `WithTimeout[A](d)` | `context.WithTimeout` | relative to when the computation *runs*, not when the operator is built |
| `WithDeadline[A](t)` | `context.WithDeadline` | an earlier parent deadline still wins |
| `Local[A](f)` | anything | general form; the others are built on it |

They are ordinary operators and compose in `Pipe`. The context flows **from the last operator inwards**, so the last one is applied first:

```go
func handleRequest(r *http.Request) RIO.ReaderIOResult[Response] {
    return F.Pipe4(
        requireUser(),
        RIO.Chain(loadProfile),
        RIO.WithTimeout[Response](5*time.Second),               // bounds the work above
        withUser[Response](userFromRequest(r)),                  // visible to everything above
        RIO.WithValue[Response](requestIDKey, r.Header.Get(HD.XRequestID)),
    )
}
```

Scope as narrowly as the requirement: put `WithTimeout` directly after the step it should bound, not around the whole handler, if only one call needs it.

### Deriving the context from an effect

When the new context value itself needs IO (generating an ID, loading a token), use `LocalIOK` / `LocalIOResultK` with the `context/reader` building blocks:

```go
// uuid.NewString is a func() string, i.e. already an IO[string]
func addRequestID(ctx context.Context) IO.IO[RIO.ContextCancel] {
    return F.Pipe1(
        uuid.NewString,
        IO.Map(F.Flow3(
            CR.WithValue[string](requestIDKey), // string -> Endomorphism[context.Context]
            RD.Read[context.Context](ctx),      // apply it to ctx
            CR.NopCancel,                       // -> Pair[CancelFunc, context.Context]
        )),
    )
}

scoped := RIO.LocalIOK[Response](addRequestID)(handler)
```

`LocalIOResultK` does the same with a fallible derivation; on failure the wrapped computation does not run.

## 5. Cancellation

- In `RIO`, `RR` and `IRR`, **`Chain` checks the context before each step**: once the context is cancelled, the next step does not run and the result is `Left(context.Cause(ctx))`. `Map` does not check (pure functions are cheap).
- `RIOC` (`context/readerio`) has no error channel, so it cannot short-circuit; use it only for IO that is fine to complete.
- `WithContext(ma)` / `WithContextK(f)` add the same check explicitly, e.g. before an expensive first step.
- `Delay`, retries, `Bracket` / `WithResource` and the HTTP client observe cancellation.
- **Leaf computations that block** (loops, `select`, custom IO) must watch `ctx.Done()` themselves and return `ctx.Err()` / `context.Cause(ctx)`:

```go
func waitForJob(id string) RIO.ReaderIOResult[Job] {
    return func(ctx context.Context) RIO.IOResult[Job] {
        return func() R.Result[Job] {
            select {
            case job := <-jobs(id):
                return R.Of(job)
            case <-ctx.Done():
                return R.Left[Job](context.Cause(ctx))
            }
        }
    }
}
```

## 6. Request-Scoped Logger

`logging.WithLogger(l)` already has `Local`'s argument shape, so a request logger is one operator. All context-aware logging (`TapSLog`, `SLog`, `LogEntryExit`) inside picks it up:

```go
F.Pipe1(
    pipeline,
    RIO.Local[Response](logging.WithLogger(slog.Default().With("requestID", id))),
)
```

Read it with `logging.GetLoggerFromContext` (falls back to the global logger). For a logger under your own key: `F.Flow2(CR.AskValue[*slog.Logger](myKey), O.GetOrElse(slog.Default))`.

## 7. Context in `Effect`

`Effect[C, A]` is `func(C) RIO.ReaderIOResult[A]`: `C` carries typed dependencies, `context.Context` is still the runtime context supplied by `RunSync(…)(ctx)`.

- `EF.Local`, `EF.Ask`, `EF.Asks` operate on **`C`**, not on `context.Context`.
- `EF.Eitherize(func(C, context.Context) (T, error))` receives both.
- To scope the runtime context of an `Effect`, lift the `RIO` operator over `C` with `RD.Map` (plain `reader.Map`):

```go
bounded := RD.Map[Deps](RIO.WithTimeout[User](2*time.Second))(fetchUserEffect)      // Effect[Deps, User]
tagged  := RD.Map[Deps](RIO.WithValue[User](requestIDKey, id))(fetchUserEffect)
```

(`RD` = `github.com/IBM/fp-go/v2/reader`.)

## 8. Outside Pipelines

When a plain `context.Context` is genuinely needed (handing it to a non-fp-go API, building a fixture), use the `context/reader` building blocks instead of the standard library calls, so reads and writes stay symmetric:

```go
ctx2 := CR.WithValue[User](userKey)(u)(ctx)       // context.WithValue, curried
user := CR.AskValue[User](userKey)(ctx2)          // Option[User]
cc   := CR.NopCancel(ctx2)                        // Pair[CancelFunc, ctx] for Local-shaped APIs
```

For optics composition over the context: `lenses.AtContext[V](key)` (`optics/lenses`) is a lens `context.Context → Option[V]`.

## 9. What Belongs in the Context

| Data | Where |
|------|-------|
| request / correlation / trace IDs | context (`WithValue`) |
| authenticated principal | context (`WithValue`) |
| request-scoped logger | context (`Local(logging.WithLogger(l))`) |
| deadlines, cancellation | context (`WithTimeout`, `WithDeadline`) |
| DB handles, HTTP clients, repositories | `Effect[Deps, A]` |
| configuration, feature flags | `Effect[Deps, A]` or function parameters |
| per-call inputs | function parameters (Kleisli arrows) |

A value that is required for the program to be correct is a dependency, not a context value — the compiler cannot check that it was provided.

## 10. Testing

```go
func TestRequireUser(t *testing.T) {
    // present
    res := F.Pipe1(requireUser(), withUser[User](alice))(t.Context())()
    assert.Equal(t, R.Of(alice), res)

    // missing -> error, no panic
    assert.Equal(t, R.Left[User](errNoUser), requireUser()(t.Context())())
}

func TestTimeout(t *testing.T) {
    res := F.Pipe1(waitForJob("slow"), RIO.WithTimeout[Job](10*time.Millisecond))(t.Context())()
    assert.Equal(t, R.Left[Job](context.DeadlineExceeded), res)
}
```

- Always run with `t.Context()`, not `context.Background()`.
- Inject context values with the same `with…` accessors production code uses.
- Test the missing-value path of every required value.
- Keep timeouts in tests short (milliseconds) and assert on `context.DeadlineExceeded`.

## Anti-Patterns

| ❌ Avoid | ✅ Prefer | Why |
|---------|----------|-----|
| `ctx.Value(k).(T)` | `RIO.AskValue[T](k)` + `GetOrElse` / `FromOption` | assertion panics when missing or mistyped |
| `v, _ := ctx.Value(k).(T)` then `if v == ""` | `AskValue` + `Option` | zero value is indistinguishable from "missing" |
| `context.WithValue(ctx, "user", u)` | unexported `ctxKey` + `WithValue[A](userKey, u)` | string keys collide across packages |
| `ctx, cancel := context.WithTimeout(ctx, d); defer cancel()` around `pipeline(ctx)()` | `F.Pipe1(pipeline, RIO.WithTimeout[A](d))` | scoping belongs to the description, not the call site |
| `Local` with `ctx2, _ := context.WithTimeout(…)` | `WithTimeout[A](d)` | discarded cancel leaks a timer until the deadline |
| `func(ctx) … { return pipeline(ctx) }` wrappers | return the pipeline value itself | the Reader already *is* that function |
| capturing `ctx` in a closure or struct field | supply it when running | the pipeline silently uses a stale context |
| DB / client stored in the context | `Effect[Deps, A]` | untyped, unchecked, hard to test |
| `pair.Unpack(logging.WithLogger(l)(ctx))` + `defer` | `RIO.Local[A](logging.WithLogger(l))` | already `Local`-shaped |
| `context.Background()` inside library code or tests | the caller's context / `t.Context()` | breaks cancellation and deadlines |
| blocking leaf that ignores `ctx.Done()` | `select` on `ctx.Done()` | timeouts and cancellation cannot take effect |

## Review Checklist

- [ ] No `ctx.Value(...)` type assertions; values are read with `AskValue`.
- [ ] Keys are unexported types; each key has a read and a write accessor.
- [ ] Required values are turned into errors explicitly (`FromOption` / `ChainOptionK`), optional ones defaulted.
- [ ] No `context.WithValue` / `WithTimeout` / `WithDeadline` / `WithCancel` inside or around pipelines — operators are used instead.
- [ ] No discarded cancel functions.
- [ ] The context is supplied once, at the edge; not captured or stored.
- [ ] Blocking leaf computations observe `ctx.Done()`.
- [ ] Dependencies live in `Effect[Deps, A]`, not in the context.
- [ ] Tests run with `t.Context()` and cover the missing-value and timeout paths.

## Import Reference

```go
import (
    RIO  "github.com/IBM/fp-go/v2/context/readerioresult"
    RR   "github.com/IBM/fp-go/v2/context/readerresult"
    RIOC "github.com/IBM/fp-go/v2/context/readerio"
    SRIO "github.com/IBM/fp-go/v2/context/statereaderioresult"
    IRR  "github.com/IBM/fp-go/v2/idiomatic/context/readerresult"
    CR   "github.com/IBM/fp-go/v2/context/reader"
    EF   "github.com/IBM/fp-go/v2/effect"
    F    "github.com/IBM/fp-go/v2/function"
    O    "github.com/IBM/fp-go/v2/option"
    R    "github.com/IBM/fp-go/v2/result"
    IO   "github.com/IBM/fp-go/v2/io"
    RD   "github.com/IBM/fp-go/v2/reader"
    HD   "github.com/IBM/fp-go/v2/http/headers"
    LS   "github.com/IBM/fp-go/v2/optics/lenses"
    "github.com/IBM/fp-go/v2/logging"
)
```

See also the `context` package documentation (`go doc github.com/IBM/fp-go/v2/context`) and the `fp-go`, `fp-go-logging` and `fp-go-http` skills.

Requires Go 1.24+.
