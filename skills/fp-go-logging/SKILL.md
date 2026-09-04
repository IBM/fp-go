---
name: fp-go-logging
description: Use this skill when working with logging in fp-go functional pipelines (github.com/IBM/fp-go/v2/logging). Trigger on mentions of logging in monadic contexts, ChainFirstIOK, ChainFirst with logging, side-effect logging, structured logging (SLog), context-embedded loggers, LoggingCallbacks, TapSLog, LogEntryExit, Printf/Logf in IO/ReaderIO/ReaderIOResult, or when the user needs to add logging to functional pipelines without breaking the data flow.
---
# fp-go Logging

## Overview

fp-go provides logging utilities that integrate naturally with functional pipelines. Logging is always a **side effect** — it should not change the value being processed. The library achieves this through `ChainFirst`-style combinators that thread the original value through unchanged while performing the log.

## Packages

| Package | Purpose |
|---------|---------|
| `github.com/IBM/fp-go/v2/logging` | Global logger, context-embedded logger, `LoggingCallbacks` |
| `github.com/IBM/fp-go/v2/io` | `Logf`, `Logger`, `LogGo`, `Printf`, `PrintGo` — IO-level logging helpers |
| `github.com/IBM/fp-go/v2/readerio` | `Logf`, `Printf`, `LogGo`, `PrintGo` — printf/template logging for ReaderIO (**no** `SLog` here) |
| `github.com/IBM/fp-go/v2/context/readerio` | `SLog`, `TapSLog`, `SLogWithCallback` — structured logging for context ReaderIO |
| `github.com/IBM/fp-go/v2/context/readerresult` | `SLog`, `TapSLog`, `SLogWithCallback` — structured logging for ReaderResult |
| `github.com/IBM/fp-go/v2/context/readerioresult` | `TapSLog`/`TapSLogInfo`/`TapSLogDebug`, `SLog`/`SLogInfo`/`SLogDebug`, `SLogLeft`, `SLogRight`, `SLogWithCallback`, `LogEntryExit`, `LogEntryExitWithCallback` — full suite for ReaderIOResult |

## Logging Inside Pipelines

The idiomatic way to log inside a monadic pipeline is `ChainFirstIOK` (or `ChainFirst` where the monad is already IO). These combinators execute a side-effecting function and pass the **original value** downstream unchanged.

### With `IOResult` / `ReaderIOResult` — printf-style

```go
import (
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    IO  "github.com/IBM/fp-go/v2/io"
    F   "github.com/IBM/fp-go/v2/function"
)

pipeline := F.Pipe3(
    fetchUser(42),
    RIO.ChainResultK(result.Eitherize1(validateUser)), // Kleisli User → Result[User]
    // Log after validation — value flows through unchanged
    RIO.ChainFirstIOK(IO.Logf[User]("Validated user: %v")),
    RIO.Map(enrichUser),
)
```

`IO.Logf[A](format string) func(A) IO[A]` logs using `log.Printf` and returns the value unchanged. It's a Kleisli arrow suitable for `ChainFirst` and `ChainFirstIOK`.

### With `IOResult` / plain `IO`

```go
import (
    IOR  "github.com/IBM/fp-go/v2/ioresult"
    IOF  "github.com/IBM/fp-go/v2/ioresult/file"   // ReadFile, WriteFile, …
    IO   "github.com/IBM/fp-go/v2/io"
    J    "github.com/IBM/fp-go/v2/json"
    F    "github.com/IBM/fp-go/v2/function"
)

pipeline := F.Pipe3(
    IOF.ReadFile("config.json"),                        // IOResult[[]byte]
    IOR.ChainResultK(J.Unmarshal[Config]),              // []byte → Result[Config]
    IOR.ChainFirstIOK(IO.Logf[Config]("Loaded config: %v")),
    IOR.Map(processConfig),
)
```

(`ioeither/file` and `ioeither` are the `Either[E, A]` equivalents; reach for them only when
the error type is not `error`. `IOR.Map` needs no leading type parameter — `IOE.Map[error]` does.)

### Logging Arrays in TraverseArray

```go
import (
    A   "github.com/IBM/fp-go/v2/array"
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    IO  "github.com/IBM/fp-go/v2/io"
    F   "github.com/IBM/fp-go/v2/function"
)

// Log each item individually, then log the final slice
pipeline := F.Pipe2(
    A.MakeBy(3, idxToFilename),
    RIO.TraverseArray(F.Flow3(
        RIO.Eitherize1(readFileCtx),                  // func(ctx, string) ([]byte, error)
        RIO.ChainResultK(J.Unmarshal[Record]),
        RIO.ChainFirstIOK(IO.Logf[Record]("Parsed record: %v")),
    )),
    RIO.ChainFirstIOK(IO.Logf[[]Record]("All records: %v")),
)
```

## IO Logging Functions

All live in `github.com/IBM/fp-go/v2/io`:

### `Logf` — printf-style

```go
IO.Logf[A any](format string) func(A) IO[A]
```

Uses `log.Printf`. The format string works like `fmt.Sprintf`.

```go
IO.Logf[User]("Processing user: %+v")
IO.Logf[int]("Count: %d")
```

### `Logger` — with custom `*log.Logger`

```go
IO.Logger[A any](loggers ...*log.Logger) func(prefix string) func(A) IO[A]
```

Uses `logger.Printf(prefix+": %v", value)`. Pass your own `*log.Logger` instance.

```go
customLog := log.New(os.Stderr, "APP ", log.LstdFlags)
logUser := IO.Logger[User](customLog)("user")
// logs: "APP user: {ID:42 Name:Alice}"
```

### `LogGo` — Go template syntax

```go
IO.LogGo[A any](tmpl string) func(A) IO[A]
```

Uses Go's `text/template`. The template receives the value as `.`.

```go
type User struct{ Name string; Age int }
IO.LogGo[User]("User {{.Name}} is {{.Age}} years old")
```

### `Printf` / `PrintGo` — stdout instead of log

Same signatures as `Logf` / `LogGo` but use `fmt.Printf`/`fmt.Println` (no log prefix, no timestamp).

```go
IO.Printf[Result]("Result: %v\n")
IO.PrintGo[User]("Name: {{.Name}}")
```

## Structured Logging in the `context` Package

The `context/readerioresult`, `context/readerresult`, and `context/readerio` packages provide structured `slog`-based logging functions that are context-aware: they retrieve the logger from the context (via `logging.GetLoggerFromContext`) rather than using a fixed logger instance.

### `TapSLog` — inline structured logging in a ReaderIOResult pipeline

`TapSLog` is an **Operator** (`func(ReaderIOResult[A]) ReaderIOResult[A]`). It sits directly in a `F.Pipe` call on a `ReaderIOResult`, logs the current value or error using `slog`, and passes the result through unchanged.

```go
import (
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    F   "github.com/IBM/fp-go/v2/function"
)

pipeline := F.Pipe4(
    fetchOrder(orderID),
    RIO.TapSLog[Order]("Order fetched"),        // logs value=<Order> or error=<err>
    RIO.Chain(validateOrder),
    RIO.TapSLog[Order]("Order validated"),
    RIO.Chain(processPayment),
)

order, err := result.Unwrap(pipeline(ctx)())
```

- Logs **both** success values (`value=<A>`) and errors (`error=<err>`) using `slog` structured attributes.
- Respects the logger level — if the logger is configured to discard Info-level logs, nothing is written.
- Available in both `context/readerioresult` and `context/readerresult`.

### `SLog`, `SLogRight`, `SLogLeft` — the underlying Kleisli arrows

In `context/readerioresult`, `SLog[A](msg)` is `Kleisli[Result[A], Void]` — it takes the **raw `Result[A]`** and returns `Void`. It is the building block `TapSLog` is made from (`TapSLog = readerio.Tap ∘ SLog`), and it **cannot** be dropped into `RIO.Chain` or `RIO.ChainFirst` on a `ReaderIOResult[A]`: those expect a Kleisli over `A`, not over `Result[A]`.

```go
// ❌ WRONG — does not compile; SLog's input is Result[Data], its output is Void
RIO.Chain(RIO.SLog[Data]("Data fetched"))

// ✅ Log value and error, pass the value through — this is what you almost always want
RIO.TapSLog[Data]("Data fetched")

// ✅ Log only the success branch via ChainFirst — SLogRight is Kleisli[A, Void]
RIO.ChainFirst(RIO.SLogRight[Data]("Data fetched"))
```

| Function | Shape | Use |
|---|---|---|
| `TapSLog[A](msg)` | `Operator[A, A]` | drop directly into `F.Pipe` — logs value **or** error, value flows through |
| `TapSLogInfo[A]` / `TapSLogDebug[A]` | `Operator[A, A]` | same, at an explicit level |
| `SLogRight[A](msg)` | `Kleisli[A, Void]` | usable with `RIO.ChainFirst` — success branch only |
| `SLogLeft(msg)` | `Kleisli[error, Void]` | usable with `RIO.ChainFirstLeft[A]` — error branch only (annotate `A`) |
| `SLog[A](msg)` / `SLogInfo` / `SLogDebug` | `Kleisli[Result[A], Void]` | low-level; prefer `TapSLog` |

**In `context/readerresult` the shape differs**: there `SLog[A](msg)` is `Kleisli[Result[A], A]` (it passes the value on rather than returning `Void`). `TapSLog[A](msg) Operator[A, A]` is identical in both packages — prefer it and the difference never bites.

### `SLogWithCallback` — custom log level and logger source

```go
import (
    RIO  "github.com/IBM/fp-go/v2/context/readerioresult"
    "log/slog"
)

// Log at DEBUG level with a custom logger extracted from context
debugLog := RIO.SLogWithCallback[User](
    slog.LevelDebug,
    logging.GetLoggerFromContext, // or any func(context.Context) *slog.Logger
    "Fetched user",
)

pipeline := F.Pipe2(
    fetchUser(123),
    RIO.Chain(debugLog),
    RIO.Map(func(u User) string { return u.Name }),
)
```

### `LogEntryExit` — automatic entry/exit timing with correlation IDs

`LogEntryExit` wraps a `ReaderIOResult` computation with structured entry and exit log messages. It assigns a unique **correlation ID** (`ID=<n>`) to each invocation so concurrent or nested operations can be correlated in logs.

```go
import (
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    F   "github.com/IBM/fp-go/v2/function"
)

pipeline := F.Pipe3(
    fetchUser(123),
    RIO.LogEntryExit[User]("fetchUser"),   // wraps the operation
    RIO.Chain(func(user User) RIO.ReaderIOResult[[]Order] {
        return F.Pipe1(
            fetchOrders(user.ID),
            RIO.LogEntryExit[[]Order]("fetchOrders"),
        )
    }),
)

orders, err := result.Unwrap(pipeline(ctx)())
// Logs:
// level=INFO msg="[entering]" name=fetchUser ID=1
// level=INFO msg="[exiting ]" name=fetchUser ID=1 duration=42ms
// level=INFO msg="[entering]" name=fetchOrders ID=2
// level=INFO msg="[exiting ]" name=fetchOrders ID=2 duration=18ms
```

On error, the exit log changes to `[throwing]` and includes the error:

```
level=INFO msg="[throwing]" name=fetchUser ID=3 duration=5ms error="user not found"
```

Key properties:
- **Correlation ID** (`ID=`) is unique per operation, monotonically increasing, and stored in the context so nested operations can access the parent's ID.
- **Duration** (`duration=`) is measured from entry to exit.
- **Logger is taken from the context** — embed a request-scoped logger with `logging.WithLogger` before executing the pipeline and `LogEntryExit` picks it up automatically.
- **Level-aware** — if the logger does not have the log level enabled, the entire entry/exit instrumentation is skipped (zero overhead).
- The original `ReaderIOResult[A]` value flows through **unchanged**.

```go
// Use a context logger so all log messages carry request metadata
cancelFn, ctxWithLogger := pair.Unpack(
    logging.WithLogger(
        slog.Default().With("requestID", r.Header.Get("X-Request-ID")),
    )(r.Context()),
)
defer cancelFn()

value, err := result.Unwrap(pipeline(ctxWithLogger)())
```

### `LogEntryExitWithCallback` — custom log level

```go
import (
    RIO  "github.com/IBM/fp-go/v2/context/readerioresult"
    "log/slog"
)

// Log at DEBUG level instead of INFO
debugPipeline := F.Pipe1(
    expensiveComputation(),
    RIO.LogEntryExitWithCallback[Result](
        slog.LevelDebug,
        logging.GetLoggerFromContext,
        "expensiveComputation",
    ),
)
```

### `SLog` / `SLogWithCallback` in `context/readerresult`

The same `SLog` and `TapSLog` functions are also available in `context/readerresult` for use with the synchronous `ReaderResult[A] = func(context.Context) (A, error)` monad:

```go
import RR "github.com/IBM/fp-go/v2/context/readerresult"

pipeline := F.Pipe3(
    queryDB(id),
    RR.TapSLog[Row]("Row fetched"),
    RR.Chain(parseRow),
    RR.TapSLog[Record]("Record parsed"),
)
```

## Global Logger (`logging` package)

The `logging` package manages a global `*slog.Logger` (structured logging, Go 1.21+).

```go
import "github.com/IBM/fp-go/v2/logging"

// Get the current global logger (defaults to slog.Default())
logger := logging.GetLogger()
logger.Info("application started", "version", "1.0")

// Replace the global logger; returns the old one for deferred restore
old := logging.SetLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
defer logging.SetLogger(old)
```

## Context-Embedded Logger

Embed a `*slog.Logger` in a `context.Context` to carry request-scoped loggers across the call stack. All context-package logging functions (`TapSLog`, `SLog`, `LogEntryExit`) pick up this logger automatically.

```go
import (
    "github.com/IBM/fp-go/v2/logging"
    "github.com/IBM/fp-go/v2/pair"
    "log/slog"
)

// Create a request-scoped logger
reqLogger := slog.Default().With("requestID", "abc-123")

// Embed it into a context using the Kleisli arrow WithLogger
cancelFn, ctxWithLogger := pair.Unpack(logging.WithLogger(reqLogger)(ctx))
defer cancelFn()

// All downstream logging (TapSLog, LogEntryExit, etc.) uses reqLogger
value, err := result.Unwrap(pipeline(ctxWithLogger)())
```

`WithLogger` returns a `ContextCancel = Pair[context.CancelFunc, context.Context]`. The cancel function is a no-op — the context is only enriched, not made cancellable.

`GetLoggerFromContext` falls back to the global logger if no logger is found in the context.

## `LoggingCallbacks` — Dual-Logger Pattern

```go
import "github.com/IBM/fp-go/v2/logging"

// Returns (infoCallback, errorCallback) — both are func(string, ...any)
infoLog, errLog := logging.LoggingCallbacks()                    // use log.Default() for both
infoLog, errLog := logging.LoggingCallbacks(myLogger)            // same logger for both
infoLog, errLog := logging.LoggingCallbacks(infoLog, errorLog)   // separate loggers
```

Used internally by `io.Logger` and by packages that need separate info/error sinks.

## Choosing the Right Logging Function

| Situation | Use |
|-----------|-----|
| Quick printf logging mid-pipeline | `IO.Logf[A]("fmt")` with `ChainFirstIOK` |
| Go template formatting mid-pipeline | `IO.LogGo[A]("tmpl")` with `ChainFirstIOK` |
| Print to stdout (no log prefix) | `IO.Printf[A]("fmt")` with `ChainFirstIOK` |
| Structured slog — log value or error inline | `RIO.TapSLog[A]("msg")` (Operator, used in Pipe) |
| Structured slog — success branch only | `RIO.ChainFirst(RIO.SLogRight[A]("msg"))` |
| Structured slog — error branch only | `RIO.ChainFirstLeft(RIO.SLogLeft("msg"))` |
| Structured slog — custom log level | `RIO.SLogWithCallback[A](level, cb, "msg")` |
| Entry/exit timing + correlation IDs | `RIO.LogEntryExit[A]("name")` |
| Entry/exit at custom log level | `RIO.LogEntryExitWithCallback[A](level, cb, "name")` |
| Structured logging globally | `logging.GetLogger()` / `logging.SetLogger()` |
| Request-scoped logger in context | `logging.WithLogger(logger)` + `logging.GetLoggerFromContext(ctx)` |
| Custom `*log.Logger` in pipeline | `IO.Logger[A](logger)("prefix")` with `ChainFirstIOK` |

## Complete Example

```go
package main

import (
    "context"
    "log/slog"
    "os"

    F   "github.com/IBM/fp-go/v2/function"
    IO  "github.com/IBM/fp-go/v2/io"
    L   "github.com/IBM/fp-go/v2/logging"
    P   "github.com/IBM/fp-go/v2/pair"
    RIO "github.com/IBM/fp-go/v2/context/readerioresult"
    "github.com/IBM/fp-go/v2/result"
)

func main() {
    // Configure JSON structured logging globally
    L.SetLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

    // Embed a request-scoped logger into the context
    _, ctx := P.Unpack(L.WithLogger(
        L.GetLogger().With("requestID", "req-001"),
    )(context.Background()))

    pipeline := F.Pipe5(
        fetchData(42),
        RIO.LogEntryExit[Data]("fetchData"),                  // entry/exit with timing + ID
        RIO.TapSLog[Data]("raw data"),                        // inline structured value log
        RIO.ChainResultK(result.Eitherize1(transformData)),   // func(Data) (Out, error)
        RIO.LogEntryExit[Out]("transformData"),
        RIO.ChainFirstIOK(IO.LogGo[Out]("result: {{.Value}}")), // template log
    )

    // pipeline(ctx)() is a single Result[Out]; unwrap it for idiomatic Go
    value, err := result.Unwrap(pipeline(ctx)())
    if err != nil {
        L.GetLogger().Error("pipeline failed", "error", err)
    }
    _ = value
}
```
