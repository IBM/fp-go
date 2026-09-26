---
name: fp-go-effect
description: >-
  Use this skill when writing, refactoring, or reviewing fp-go v2 service code
  built on the effect package (github.com/IBM/fp-go/v2/effect), and whenever a
  service needs dependencies such as database or HTTP clients, repositories,
  configuration, environment access or tool registries. Explains Effect[C, A],
  why the type parameter C carries the dependencies (and context.Context does
  not), how to design capability interfaces (XxxDeps with getters, MakeXxxDeps,
  AsXxxDeps), how to narrow and widen them with Local, how to lift Go
  functions and interface methods with Eitherize, Eitherize1, Asks and
  FromIdiomatic, how to build dependencies with an Effect, and how to Provide
  and RunSync at the edge and test with fakes. Trigger on mentions of Effect,
  effect.Effect, dependency injection in fp-go, Provide, RunSync, Asks, Local,
  LocalEffectK, Eitherize1, typed dependencies, ReaderReaderIOResult, or on
  code that stores clients in context.Context or threads them through
  function parameters.
---

# fp-go Effect: Typed Dependencies

## Core Principle: `C` Is for Dependencies

```go
// Effect[C, A] = func(C) func(context.Context) func() Result[A]
type Effect[C, A any] = readerreaderioresult.ReaderReaderIOResult[C, A]
```

An `Effect` takes three inputs, and each input has its own job:

| Input | Carries | Supplied |
|---|---|---|
| `C` | **dependencies**: long-lived collaborators such as repositories, DB and HTTP clients, SDK clients, configuration, environment access, tool registries, clocks | once, at the composition root, with `EF.Provide` |
| `context.Context` | **request scope**: cancellation, deadlines, request / trace IDs, authenticated principal, request logger | once per run, with `EF.RunSync(...)(ctx)` |
| Kleisli argument (`func(A) Effect[C, B]`) | **per-call input**: IDs, payloads, the value flowing through the pipeline | as a function argument |

Rules that follow from this:

1. **Every dependency goes into `C`.** Never store it in `context.Context` (untyped, unchecked, panics on a missing key). Never pass it as a parameter through every layer, and never keep it in a package-level global.
2. **`C` is a type, so the compiler checks the wiring.** A function's signature `Effect[UserDeps, User]` states exactly what it needs. Code that forgets to provide something does not compile.
3. **Declare the narrowest `C` per function.** A function that only reads users asks for `UserDeps`, not for the whole application. Widen it with `Local` where it is composed.
4. **Provide once, run many times.** Build the dependency graph at startup, `Provide` it, and `RunSync` per request with that request's context.
5. **Request-scoped data stays in `context.Context`** (see the **fp-go-context** skill). If a value differs per request, it is not a dependency.

## Before You Generate

fp-go is low-frequency in training data, so signatures are easy to misremember.
For any combinator not shown below, look it up via the fp-go MCP server's
`search_examples` / `get_example` tools (see the **fp-go-mcp** skill) instead of
guessing. After writing code, run `go build ./...` and `go vet ./...` and fix any
type-parameter or argument-order errors before presenting it.

## When to Use Effect

| Situation | Use |
|---|---|
| Needs only `context.Context` (no collaborators) | `context/readerioresult` (`RIO`) |
| Needs any collaborator: client, repository, config, env, registry | **`Effect[C, A]`** |
| Pure computation | a plain function or `Flow`, no monad |

`Effect` is the recommended top-level monad for service code. Every `RIO` operator still works inside it (see "Runtime Context Inside an Effect").

## Designing the Dependency Type

### One Capability Interface per Package

Each package that needs something from outside declares a small **capability interface** with getter methods, a private implementation, a constructor and a generic upcast helper:

```go
package users

// UserRepo is the port; the adapter (Postgres, in-memory, fake) implements it.
type UserRepo interface {
    FindUser(ctx context.Context, id int) (User, error)
}

// UserDeps is what this package needs from its environment.
type UserDeps interface {
    GetUserRepo() UserRepo
}

type userDeps struct{ repo UserRepo }

func (d *userDeps) GetUserRepo() UserRepo { return d.repo }

// MakeUserDeps is the only way to build the dependency.
func MakeUserDeps(repo UserRepo) UserDeps { return &userDeps{repo} }

// AsUserDeps upcasts any wider dependency that embeds UserDeps.
// It is used with EF.Local to run a UserDeps effect inside a wider C.
func AsUserDeps[R UserDeps](r R) UserDeps { return r }
```

- **Getters, not fields.** An interface lets callers supply any implementation and lets tests supply fakes. The implementation stays private.
- **Getters may return capabilities as functions**, for example `GetLookupEnv() IOR.Kleisli[string, string]` or `GetRequestOptions() EF.Thunk[[]Option]`. A capability that is a function is trivial to fake and is only evaluated when the effect runs.
- **Keep it small.** A `XxxDeps` interface lists what *this* package uses, typically one to three getters.

A plain struct of fields (`type Deps struct{ DB DB; Config Config }`) also works for a small program. Capability interfaces scale better, because every package keeps its own narrow view and the composite is assembled by embedding.

### Composite Dependencies by Embedding

Higher layers combine the capabilities of the packages they use. The interface embeds interfaces, and the implementation embeds the implementations:

```go
package app

type AppDeps interface {
    users.UserDeps
    mail.MailDeps
}

type appDeps struct {
    users.UserDeps
    mail.MailDeps
}

func MakeAppDeps(u users.UserDeps, m mail.MailDeps) AppDeps { return &appDeps{u, m} }
```

Because Go generics are invariant, an `Effect[UserDeps, A]` is **not** an `Effect[AppDeps, A]`, even though `AppDeps` embeds `UserDeps`. `Local` performs that conversion (next sections).

## Lifting Leaves

The leaves are where Go code meets the effect. Lift, never hand-write the nested closures:

| Go shape | Lift with | Result |
|---|---|---|
| interface method `func(Repo, ctx, A) (B, error)` (method expression `Repo.Find`) | `EF.Eitherize1(Repo.Find)` | `Kleisli[Repo, A, B]` |
| `func(C, ctx) (B, error)` | `EF.Eitherize(f)` | `Effect[C, B]` |
| `func(C, ctx, A) (B, error)` | `EF.Eitherize1(f)` | `Kleisli[C, A, B]` |
| service method `func(A) func(ctx, C) (B, error)` | `EF.FromIdiomatic(f)` | `Kleisli[C, A, B]` |
| **pure** getter / projection `func(C) A` | `EF.Asks(f)`, e.g. `EF.Asks(MailDeps.GetSender)` | `Effect[C, A]` |
| the whole dependency | `EF.Ask[C]()` | `Effect[C, C]` |
| dependency-free `ReaderIOResult[A]` | `EF.FromThunk[C](t)` | `Effect[C, A]` |
| `Result[A]`, `IO[A]`, value, error | `EF.FromResult[C]`, `EF.FromIO[C]`, `EF.Of[C]`, `EF.Fail[C, A]` | `Effect[C, A]` |

Interface method expressions make the leaves point-free. The receiver becomes the dependency, and `Local` with the getter fetches it from the capability interface:

```go
func FindUser() EF.Kleisli[UserDeps, int, User] {
    return F.Flow2(
        EF.Eitherize1(UserRepo.FindUser),     // Kleisli[UserRepo, int, User]
        EF.Local[User](UserDeps.GetUserRepo), // run it on UserDeps
    )
}
```

**`Asks` is for pure projections only.** `Effect[C, A]` *is* `func(C) ReaderIOResult[A]`, so `EF.Asks(func(c C) ReaderIOResult[A] {...})` silently yields the nested `Effect[C, ReaderIOResult[A]]`. For a getter that returns an effectful capability, use `Asks` followed by `ChainThunkK`:

```go
func LookupEnv(key string) EF.Effect[EnvDeps, string] {
    return F.Pipe1(
        EF.Asks(EnvDeps.GetLookupEnv), // Effect[EnvDeps, IOR.Kleisli[string, string]]
        EF.ChainThunkK[EnvDeps](F.Flow2(
            RD.Read[IOR.IOResult[string]](key), // apply the capability to key
            RIO.FromIOResult[string],
        )),
    )
}
```

## Narrowing with `Local`

`EF.Local[A](f)` takes `f: C1 → C2` (from what you have to what the effect needs) and turns an `Effect[C2, A]` into an `Effect[C1, A]`. With the `AsXxxDeps` helpers this is how a narrow effect runs inside a wide dependency:

```go
func NotifyUser() EF.Kleisli[AppDeps, int, MessageID] {
    findUser := F.Flow2(users.FindUser(), EF.Local[User](users.AsUserDeps[AppDeps]))
    sendMail := F.Flow2(mail.SendMail(), EF.Local[MessageID](mail.AsMailDeps[AppDeps]))

    return F.Flow3(
        findUser,
        EF.Map[AppDeps](F.Flow2(getEmail, welcome)),
        EF.Chain(sendMail),
    )
}
```

- `Local` needs only the value type `[A]`. `C1` and `C2` are inferred from `f`, and `AsXxxDeps[Wide]` fixes the generic helper to the wide type.
- **`Asks` + `Map` versus `Local`:** build a *new* effect from a getter and a pure function with `EF.Asks(getter)` plus `EF.Map[C](f)`. Use `Local` when you already have an `Effect[C2, A]` (a Kleisli defined elsewhere) and want to reuse it unchanged under another `C`.
- `EF.ContraMap` is an alias of `Local`.
- Compose with the `effect` combinators (`Asks`, `Map`, `Ap`, `Chain`, `ChainThunkK`, `FromThunk`, `Local`). Do not rebuild `Effect`'s nested reader shape by hand with `reader` / `context/readerioresult`. That also type-checks, but it depends on the internal representation and is hard to read.

When deriving the inner dependency needs more than a pure function, use the `Local*K` variants (`f` goes from the outer to the inner dependency):

| `f` | Operator |
|---|---|
| `C1 → C2`, pure | `Local` / `ContraMap` |
| `C1 → Reader[context.Context, C2]` | `LocalReaderK` |
| `C1 → IO[C2]` | `LocalIOK` |
| `C1 → Result[C2]` (validation) | `LocalResultK` |
| `C1 → IOResult[C2]` | `LocalIOResultK` |
| `C1 → ReaderIOResult[C2]` | `LocalThunkK` |
| `C1 → Effect[C1, C2]` (most general) | `LocalEffectK` |

## Composing

`Effect` has the full API of the other monads. Mind the leading type parameters, which `C` often cannot be inferred for:

| Operation | Form |
|---|---|
| transform the value | `EF.Map[C](f)` (`Map[C, A, B]`: annotate `C`) |
| sequence a Kleisli | `EF.Chain(k)` (inferred from `k`) |
| combine independent effects | `EF.Ap[B](fa)` on an `Effect[C, func(A) B]` |
| do-notation | `EF.Do[C](S{})`, `EF.Bind(lens.Set, k)`, `EF.ApS(lens.Set, eff)`, `EF.ApSL(lens, eff)`, `EF.Let` |
| lift other monads mid-pipeline | `EF.ChainResultK`, `EF.ChainIOK`, `EF.ChainThunkK[C]`, `EF.ChainReaderK` |
| side effects | `EF.Tap(k)`, `EF.TapIOK[C](f)`, `EF.TapThunkK[C](f)` |
| recover / alternative | `EF.ChainLeft(k)`, `EF.Alt(LZ.Of(other))`, `EF.TapLeft[A](k)` |
| slices | `EF.TraverseArray(k)` |
| retries | `EF.Retrying(policy, action, check)` |
| run | `EF.Provide[A](deps)`, then `EF.RunSync(thunk)(ctx)` |

`Ap[B, C, A]` and `Provide[A, C]` lead with the type that cannot be inferred. Write `EF.Ap[ChatDeps](apiKey)` and `EF.Provide[User](deps)`.

## Building Dependencies with an Effect

Constructing a dependency is often effectful itself: reading an environment variable, loading a file, opening a connection. Model the constructor as an `Effect` over the *bootstrap* capabilities:

```go
type BootDeps interface {
    env.EnvDeps
    GetMailer() Mailer
}

// The sender address is read lazily from the environment when the effect runs.
func MakeMailDepsFromEnv() EF.Effect[BootDeps, MailDeps] {
    sender := F.Pipe1(env.LookupEnv("MAIL_SENDER"), EF.Local[string](env.AsEnvDeps[BootDeps]))

    return F.Pipe1(
        EF.Asks(F.Flow2(BootDeps.GetMailer, F.Curry2(MakeMailDeps))), // Effect[BootDeps, func(string) MailDeps]
        EF.Ap[MailDeps](sender),
    )
}
```

Run it once at startup and provide the result. Alternatively, `EF.LocalEffectK[A](F.Constant1[BootDeps](MakeMailDepsFromEnv()))` runs a `MailDeps` effect directly on `BootDeps`, building the dependency on every run. A missing variable fails the effect with a normal error. Nothing panics.

## Running at the Edge

The composition root (`main`, server setup, CLI command) is the only place that builds concrete dependencies and runs effects:

```go
func main() {
    deps := app.MakeAppDeps(
        users.MakeUserDeps(postgres.NewUserRepo(db)),
        mail.MakeMailDeps(smtp.NewMailer(cfg), cfg.Sender),
    )

    http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
        id, err := EF.RunSync(EF.Provide[MessageID](deps)(app.NotifyUser()(userID(r))))(r.Context())
        // write id or err
    })
}
```

- `Provide` eliminates `C` and yields a `ReaderIOResult[A]`. `RunSync` runs it with a `context.Context` and returns an idiomatic `(A, error)`. Alternatively, `thunk(ctx)()` yields a `Result[A]`.
- Build the dependencies **once**, not per request. Pass the request's context, not `context.Background()`.
- Library code returns `Effect` or `Kleisli` values. Only the edge calls `Provide` and `RunSync`.

## Runtime Context Inside an Effect

`EF.Local`, `EF.Ask` and `EF.Asks` operate on `C`. To scope the **runtime** `context.Context`, lift any `context/readerioresult` operator over `C` with `reader.Map`:

```go
F.Pipe1(
    app.NotifyUser()(42),
    RD.Map[AppDeps](RIO.WithTimeout[MessageID](2*time.Second)), // Effect[AppDeps, MessageID]
)
RD.Map[AppDeps](RIO.WithValue[MessageID](requestIDKey, id))
RD.Map[AppDeps](RIO.LogEntryExit[MessageID]("notifyUser"))
```

## Testing

Tests supply fakes through the same constructors production code uses. No mocking framework and no global state are needed:

```go
type fakeRepo map[int]User

func (r fakeRepo) FindUser(_ context.Context, id int) (User, error) {
    if u, ok := r[id]; ok {
        return u, nil
    }
    return User{}, errNotFound
}

func TestFindUser(t *testing.T) {
    deps := users.MakeUserDeps(fakeRepo{1: {ID: 1, Email: "a@x"}}) // only what FindUser declares

    u, err := EF.RunSync(EF.Provide[User](deps)(users.FindUser()(1)))(t.Context())
    assert.NoError(t, err)
    assert.Equal(t, "a@x", u.Email)

    _, err = EF.RunSync(EF.Provide[User](deps)(users.FindUser()(2)))(t.Context())
    assert.ErrorIs(t, err, errNotFound)
}
```

- Test each function with the narrowest dependency it declares. That is the pay-off of narrow `C`s.
- Run with `t.Context()`.
- A fake capability can also be a function value, e.g. `MakeToolDeps(fakeCaller)` where the getter returns a lookup function built from a map.

## Common Mistakes

| Mistake | Fix |
|---|---|
| `ctx.Value("db").(DB)` or `context.WithValue(ctx, dbKey, db)` | The DB is a dependency: put it in `C` behind a getter |
| Passing `db`, `client`, `cfg` as parameters through every function | Return `Effect[XxxDeps, A]` / `Kleisli`; the dependency travels in `C` |
| Package-level `var client = ...` singletons | Build in `MakeXxxDeps`, provide at the edge |
| Request IDs, principal or deadlines in `C` | They vary per request: `context.Context` (`RIO.WithValue`, `RIO.WithTimeout` via `RD.Map[C]`) |
| One huge `AppDeps` as `C` of every function | Declare the narrowest `XxxDeps`; widen with `EF.Local[A](AsXxxDeps[Wide])` |
| `Effect[AppDeps, A]` expected, `Effect[UserDeps, A]` given, "fixed" with a type assertion | Generics are invariant: use `EF.Local` with the `AsXxxDeps` helper |
| `EF.Asks(func(c C) RIO.ReaderIOResult[A] {...})` | Yields `Effect[C, ReaderIOResult[A]]`. Use `EF.Eitherize` / `Eitherize1` / `FromIdiomatic`, or `Asks` + `ChainThunkK` |
| Hand-written `func(c C) func(ctx) func() Result[A]` around a Go method | `EF.Eitherize1(Iface.Method)` + `EF.Local[B](Deps.GetIface)` |
| Rebuilding `Effect` with `reader.Map` / `RIO` internals | Compose with `EF.Asks`, `Map`, `Ap`, `Chain`, `ChainThunkK`, `FromThunk`, `Local` |
| `EF.Map(f)`, `EF.Provide(deps)` without annotation | `EF.Map[C](f)`, `EF.Provide[A](deps)` |
| Building dependencies inside the request handler | Build once at startup; provide per run |
| `Provide` / `RunSync` inside library code | Return the `Effect`; only the edge runs it |
| `C = context.Context` | That is `RIO`; use `context/readerioresult` directly |

## Review Checklist

- [ ] Every collaborator (client, repository, config, env access, registry) comes from `C`, not from `context.Context`, globals or threaded parameters.
- [ ] Each package defines a small `XxxDeps` interface with getters, a private implementation, `MakeXxxDeps` and `AsXxxDeps`.
- [ ] Functions declare the narrowest `C` and are widened with `EF.Local[A](AsXxxDeps[Wide])`.
- [ ] Leaves are lifted with `Eitherize` / `Eitherize1` / `FromIdiomatic` / `Asks`; no hand-written nested closures.
- [ ] `Asks` receives only pure projections.
- [ ] Request-scoped data stays in `context.Context`; runtime scoping uses `RD.Map[C](RIO.…)`.
- [ ] Dependencies are built once at the composition root; `Provide` and `RunSync` appear only there and in tests.
- [ ] Tests provide fakes through `MakeXxxDeps` and run with `t.Context()`.

## Import Reference

```go
import (
    EF  "github.com/IBM/fp-go/v2/effect"
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    RD  "github.com/IBM/fp-go/v2/reader"
    IOR "github.com/IBM/fp-go/v2/ioresult"
    F   "github.com/IBM/fp-go/v2/function"
    LZ  "github.com/IBM/fp-go/v2/lazy"
    R   "github.com/IBM/fp-go/v2/result"
)
```

See also the `fp-go` skill (core types, aliases), `fp-go-context` (request-scoped values, timeouts) and `fp-go-logging`.
