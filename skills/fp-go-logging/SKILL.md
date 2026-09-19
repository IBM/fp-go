---
name: fp-go-logging
description: Use this skill when adding logging to fp-go functional pipelines (github.com/IBM/fp-go/v2). Trigger on logging inside Pipe/Flow, Tap/TapIOK/TapReaderIOK/TapLeft used for logging, TapSLog, SLogLeft/SLogRight, entry/exit logging (LogEntryExit, LogEntryExitWithCallback), correlation IDs, adding context to errors before they are logged (MapLeft with errors.OnError), request-scoped slog loggers in context.Context (logging.WithLogger, GetLoggerFromContext), or IO.Logf/Printf in IO/IOResult/ReaderIOResult.
---
# fp-go Logging

Logging is a **side effect**: it must never change the value or the error flowing through a
pipeline. fp-go expresses this with the `Tap*` family of operators. A `Tap` runs an effect on
the current value (or error) and hands the **original** computation downstream unchanged.

This skill covers three tasks:

1. [Add a simple log statement](#1-a-simple-log-statement-with-tap): `TapSLog`, `TapIOK`, `TapReaderIOK`, `TapLeft`
2. [Add entry and exit logs](#2-entry-and-exit-logs): `LogEntryExit`
3. [Add context to the error channel](#3-adding-context-to-the-error-channel): `MapLeft` + `errors.OnError`, then log once

Examples assume these imports; all live under `github.com/IBM/fp-go/v2`:

```go
import (
    "log/slog"

    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    ER  "github.com/IBM/fp-go/v2/errors"
    F   "github.com/IBM/fp-go/v2/function"
    IO  "github.com/IBM/fp-go/v2/io"
    "github.com/IBM/fp-go/v2/logging"
    R   "github.com/IBM/fp-go/v2/result"
)
```

Runnable versions of the snippets are in `v2/context/readerioresult/logging_example_test.go`
(`ExampleTapSLog`, `ExampleTapIOK_logging`, `ExampleTapReaderIOK_logging`, `ExampleSLogRight`,
`ExampleSLogLeft`, `ExampleLogEntryExit`, `ExampleLogEntryExit_nested`,
`ExampleLogEntryExitWithCallback`, `ExampleMapLeft`, `ExampleMapLeft_logging`).

## Quick reference

| Goal | Operator (in `F.Pipe` on a `ReaderIOResult[A]`) |
|---|---|
| Log value **or** error, structured | `RIO.TapSLog[A]("msg")` |
| Same at DEBUG level | `RIO.TapSLogDebug[A]("msg")` |
| Log success only | `RIO.Tap(RIO.SLogRight[A]("msg"))` |
| Log error only (ERROR level) | `RIO.TapLeft[A](RIO.SLogLeft("msg"))` |
| printf-style log of the value | `RIO.TapIOK(IO.Logf[A]("msg: %v"))` |
| Custom attributes / custom message | `RIO.TapReaderIOK(logX)` with `logX: func(A) ReaderIO[Void]` |
| Entry + exit logs with timing and correlation ID | `RIO.LogEntryExit[A]("name")` |
| Entry + exit at another level | `RIO.LogEntryExitWithCallback[A](slog.LevelDebug, logging.GetLoggerFromContext, "name")` |
| Add context to an error | `RIO.MapLeft[A](ER.OnError("loading user %d", id))` |
| Choose the logger for a pipeline | `RIO.Local[A](logging.WithLogger(l))` |

`Tap*` and `ChainFirst*` are aliases (`TapIOK` = `ChainFirstIOK`, `TapLeft` = `ChainFirstLeft`, …).
Prefer the `Tap` names for logging; they state the intent.

## 1. A simple log statement with `Tap`

### Default: `TapSLog`

`TapSLog[A](msg)` is an `Operator[A, A]`. Put it in the pipe **after** the step you want to
observe. It logs `value=<A>` on success or `error=<err>` on failure at INFO level, using the
logger from the context, and returns the computation unchanged.

```go
func greeting(id int) RIO.ReaderIOResult[string] {
    return F.Pipe2(
        fetchUser(id),
        RIO.TapSLog[User]("fetched user"),
        RIO.Map(getName),
    )
}
// level=INFO msg="fetched user" value="{ID:42 Name:Alice}"
// level=INFO msg="fetched user" error="not found"
```

- Logs happen when the pipeline **runs**, not when it is built.
- The type parameter is the type at that point in the pipe (`User` here, not `string`).
- `TapSLogInfo` / `TapSLogDebug` choose the level explicitly. Debug logs cost nothing when the
  logger's level filters them out.
- `context/readerresult` and `context/readerio` have their own `TapSLog[A]` with the same usage.

### printf-style: `TapIOK` with `IO.Logf`

Any `io.Kleisli[A, B]` (`func(A) IO[B]`) can be tapped with `TapIOK`. The `io` package
has ready-made ones:

```go
RIO.TapIOK(IO.Logf[User]("fetched user %+v"))       // log.Printf
RIO.TapIOK(IO.LogGo[User]("user {{.Name}}"))        // text/template, log.Println
RIO.TapIOK(IO.Printf[User]("fetched user %+v\n"))   // fmt.Printf (stdout, no timestamp)
RIO.TapIOK(IO.Logger[User](myStdLogger)("user"))    // custom *log.Logger
```

These use the standard `log` package, not `slog`, and log **only on success**; an error skips
the tap. Prefer `TapSLog` when the application uses structured logging.

The same operator exists for the other monads: `IOR.TapIOK` (ioresult), `RR.TapIOK`
(context/readerresult), `IO.ChainFirst` (plain IO).

### Custom attributes: `TapReaderIOK`

When you want a specific message and specific attributes, write a small named
`func(A) ReaderIO[Void]` that reads the logger from the context. Reading from the context means
the log uses the request-scoped logger and inherits the correlation ID of an enclosing
`LogEntryExit`.

```go
import RIOC "github.com/IBM/fp-go/v2/context/readerio"

func logUserLoaded(u User) RIOC.ReaderIO[F.Void] {
    return func(ctx context.Context) IO.IO[F.Void] {
        return func() F.Void {
            logging.GetLoggerFromContext(ctx).InfoContext(ctx, "user loaded", "id", u.ID, "name", u.Name)
            return F.VOID
        }
    }
}

pipeline := F.Pipe2(
    fetchUser(42),
    RIO.TapReaderIOK(logUserLoaded),
    RIO.Map(getName),
)
// level=INFO msg="user loaded" id=42 name=Alice
```

Do not call `slog.Info(...)` from inside `RIO.Map`: `Map` is for pure functions, and the global
logger ignores the logger scoped to the context.

### Success only / error only

`SLogRight` and `SLogLeft` are Kleisli arrows, so they go through `Tap` / `TapLeft`:

```go
RIO.Tap(RIO.SLogRight[User]("fetched user"))             // INFO, success only
RIO.TapLeft[User](RIO.SLogLeft("fetching user failed"))   // ERROR, error only; the error still propagates
```

`TapLeft` needs the explicit `[A]` because `A` cannot be inferred from the error handler.

### Common mistakes

```go
// ❌ SLog is Kleisli[Result[A], Void]: its input is the Result, not A, so this does not compile
RIO.Chain(RIO.SLog[User]("fetched"))
// ✅
RIO.TapSLog[User]("fetched")

// ❌ Chain replaces the value with Void
RIO.Chain(RIO.SLogRight[User]("fetched"))
// ✅
RIO.Tap(RIO.SLogRight[User]("fetched"))

// ❌ side effect inside a pure Map; bypasses the context logger
RIO.Map(func(u User) User { slog.Info("fetched", "user", u); return u })
// ✅
RIO.TapSLog[User]("fetched")
```

## 2. Entry and exit logs

`LogEntryExit[A](name)` is an `Operator[A, A]` that **wraps the computation it is applied to**.
It logs `[entering]` before that computation starts and `[exiting ]` (success) or `[throwing]`
(failure, with the error) after it ends. Both records carry the same correlation `ID`, and the
exit record carries the `duration`.

```go
func loadUser(id int) RIO.ReaderIOResult[User] {
    return F.Pipe1(
        fetchUser(id),
        RIO.LogEntryExit[User]("loadUser"),
    )
}
// level=INFO msg=[entering] ID=7 name=loadUser
// level=INFO msg="[exiting ]" ID=7 name=loadUser duration=1.2ms
//
// on failure:
// level=INFO msg=[entering] ID=8 name=loadUser
// level=INFO msg=[throwing] ID=8 name=loadUser duration=0.4ms error="not found"
```

### What gets wrapped

In a `F.Pipe`, `LogEntryExit` wraps **everything above it**. Put it last to time a whole
function, or apply it to a single step to time just that step:

```go
// Whole pipeline: the TapSLog runs between entry and exit and is logged with the same ID
loadName := F.Pipe3(
    fetchUser(42),
    RIO.TapSLog[User]("fetched user"),
    RIO.Map(getName),
    RIO.LogEntryExit[string]("loadName"),
)
// level=INFO msg=[entering] ID=3 name=loadName
// level=INFO msg="fetched user" ID=3 value="{ID:42 Name:Alice}"
// level=INFO msg="[exiting ]" ID=3 name=loadName duration=…

// One step inside a larger pipeline: wrap the Kleisli, not the pipeline
F.Pipe2(
    fetchUser(42),
    RIO.Chain(F.Flow2(fetchOrders, RIO.LogEntryExit[[]Order]("fetchOrders"))),
    RIO.Map(summarize),
)
```

Nested `LogEntryExit` calls each get their own ID. The ID is installed on the context logger
for the wrapped computation, so every context-aware log inside it (`TapSLog`, `SLogLeft`,
`TapReaderIOK` using `GetLoggerFromContext`) is correlated automatically.

### Why not `TapSLog("entering")`?

A `Tap` placed before a step runs after the **previous** step. It can mark "about to start",
but it cannot log the end of a failed step, measure the duration or correlate the two records.
Use `LogEntryExit` for entry/exit logging.

### Level and logger

```go
RIO.LogEntryExitWithCallback[User](slog.LevelDebug, logging.GetLoggerFromContext, "loadUser")
```

If the logger does not have the level enabled, the whole instrumentation (ID, timer, both
records) is skipped. The second argument is any `func(context.Context) *slog.Logger`.

## 3. Adding context to the error channel

An error usually surfaces far from where it happened. Logging it at every layer produces
duplicate records that lack context. The idiomatic approach:

1. **Enrich** the error where the context is known, with `MapLeft` and `errors.OnError`.
2. **Log once**, at the boundary (handler, command, job), with `TapLeft(SLogLeft(...))`,
   `TapSLog`, or the `[throwing]` record of `LogEntryExit`.

### `MapLeft` + `errors.OnError`

`RIO.MapLeft[A](f Endomorphism[error])` rewrites the error and leaves successes untouched.
`ER.OnError(format, args...)` builds such an endomorphism. It wraps the original error with
`%w`, so `errors.Is` / `errors.As` keep working.

```go
func loadGreeting(id int) RIO.ReaderIOResult[string] {
    return F.Pipe2(
        fetchUser(id),
        RIO.Map(getName),
        RIO.MapLeft[string](ER.OnError("loading greeting for user %d", id)),
    )
}
// err.Error() == "loading greeting for user 7, Caused By: not found"
// errors.Is(err, ErrNotFound) == true
```

Capture the inputs that explain the failure (`id` above) from the enclosing Kleisli's
parameters. The error then says *which* call failed, not just *that* one failed.

### Log once at the boundary

```go
func handle(id int) RIO.ReaderIOResult[string] {
    return F.Pipe2(
        loadGreeting(id),
        RIO.MapLeft[string](ER.OnError("handling request")),
        RIO.TapLeft[string](RIO.SLogLeft("request failed")),
    )
}
// level=ERROR msg="request failed" error="handling request, Caused By: loading greeting for user 7, Caused By: not found"
```

### Combine with `LogEntryExit`

`[throwing]` logs the error as it leaves the wrapped computation. Put `MapLeft` **above**
`LogEntryExit` so the exit record includes the context:

```go
F.Pipe2(
    fetchUser(id),
    RIO.MapLeft[User](ER.OnError("loading user %d", id)),  // enrich first …
    RIO.LogEntryExit[User]("loadUser"),                      // … then [throwing] logs the enriched error
)
```

### Rules

- Use `MapLeft` to change the error. Do not use `OrElse`/`ChainLeft` with `RIO.Left` just to
  rewrap; those are for recovery.
- Always wrap with `%w` (which `ER.OnError` does). `fmt.Errorf("…: %v", err)` breaks `errors.Is`.
- Enrich where the context is; log where the error is handled. A `TapLeft` log in every layer
  duplicates the record.
- `MapLeft` is available in `context/readerioresult`, `context/readerresult`,
  `context/readerreaderioresult`, `readerresult`, `result` and `ioresult`. In `result`/`ioresult`
  the function has type `func(error) E`; `ER.OnError(...)` fits with `E = error`.

## Choosing the logger

All context-aware operators (`TapSLog*`, `SLogLeft/Right`, `LogEntryExit*`) read the logger with
`logging.GetLoggerFromContext`, which falls back to the global logger.

```go
// Global logger (defaults to slog.Default())
old := logging.SetLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
defer logging.SetLogger(old)

// Request-scoped logger: install it around the whole pipeline with Local.
// logging.WithLogger(l) is the ctx → Pair[CancelFunc, ctx] arrow that Local expects.
reqLogger := logging.GetLogger().With("requestID", requestID)

value, err := R.Unwrap(F.Pipe1(
    pipeline,
    RIO.Local[Data](logging.WithLogger(reqLogger)),
)(ctx)())
```

`Local` exists for `context/readerio`, `context/readerresult` and `context/statereaderioresult`
as well. To get a plain enriched context outside a pipeline, use
`_, ctx2 := pair.Unpack(logging.WithLogger(reqLogger)(ctx))`. For a logger stored under your
own key, build the callback with
`F.Flow2(CR.AskValue[*slog.Logger](myKey), O.GetOrElse(slog.Default))`. For general context
handling see the `fp-go-context` skill.

## Package overview

| Package | Logging API |
|---|---|
| `context/readerioresult` | `TapSLog`, `TapSLogInfo`, `TapSLogDebug`, `SLog*`, `SLogLeft`, `SLogRight`, `SLogWithCallback`, `LogEntryExit`, `LogEntryExitWithCallback`, `MapLeft`, `Tap*`, `Local` |
| `context/readerresult` | `TapSLog`, `SLog` (`Kleisli[Result[A], A]`), `TapIOK`, `MapLeft`, `Local` |
| `context/readerio` | `TapSLog`, `SLog`, `Tap`, `TapIOK`, `Local` |
| `io` | `Logf`, `LogGo`, `Printf`, `PrintGo`, `Logger`: `func(A) IO[A]`, used with `TapIOK` |
| `readerio` | `Logf`, `LogGo`, `Printf`, `PrintGo` for the generic ReaderIO |
| `logging` | `GetLogger`, `SetLogger`, `WithLogger`, `GetLoggerFromContext`, `LoggingCallbacks` |
| `errors` | `OnError` (wrap with context), `OnNone`, `OnSome` (build errors) |

`SLogWithCallback[A](level, cb, msg)` has the shape of `SLog`. Turn it into an operator with
`RIOC.Tap(...)` (`context/readerio`) when you need `TapSLog` with a custom level and logger.

## Complete example

```go
func loadUser(id int) RIO.ReaderIOResult[User] {
    return F.Pipe3(
        fetchUser(id),
        RIO.TapSLogDebug[User]("fetched user"),              // simple log statement
        RIO.MapLeft[User](ER.OnError("loading user %d", id)), // add context to the error
        RIO.LogEntryExit[User]("loadUser"),                   // entry/exit around the three steps
    )
}

func handleGreeting(reqLogger *slog.Logger, id int) RIO.ReaderIOResult[string] {
    return F.Pipe3(
        loadUser(id),
        RIO.Map(getName),
        RIO.TapLeft[string](RIO.SLogLeft("greeting failed")), // log the error once, at the boundary
        RIO.Local[string](logging.WithLogger(reqLogger)),     // every log above uses reqLogger
    )
}

name, err := R.Unwrap(handleGreeting(reqLogger, 7)(ctx)())
```
