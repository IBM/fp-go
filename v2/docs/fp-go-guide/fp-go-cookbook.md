# fp-go/v2 Cookbook

Task-oriented recipe guide for migrating Go code to fp-go and building new features.
All import paths use `github.com/IBM/fp-go/v2/`. Module requires Go 1.24+.

## Convention: import aliases used throughout

```go
import (
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
    O "github.com/IBM/fp-go/v2/option"
    E "github.com/IBM/fp-go/v2/either"
    A "github.com/IBM/fp-go/v2/array"
)
```

---

## 1. Wrapping Go Functions

### Recipe: Convert `func(X) (Y, error)` with Eitherize1

**Problem**: You have an idiomatic Go function returning `(Y, error)` and need it to return `result.Result[Y]`.
**Solution**: Use `result.Eitherize1` to automatically wrap the tuple return into `Result`.

```go
package main

import (
    "strconv"
    R "github.com/IBM/fp-go/v2/result"
)

// Original Go function:
//   func strconv.Atoi(s string) (int, error)

// Wrap it:
var parseInt = R.Eitherize1(strconv.Atoi)

// parseInt has type: func(string) R.Result[int]
// Usage:
//   parseInt("42")    => Right(42)
//   parseInt("oops")  => Left(error)
```

### Recipe: Convert inline `(value, error)` with TryCatchError

**Problem**: You have a `(value, error)` pair already in hand (not a function) and need a `Result`.
**Solution**: Use `result.TryCatchError` to wrap the pair directly.

```go
package main

import (
    "os"
    R "github.com/IBM/fp-go/v2/result"
)

func readConfig() R.Result[[]byte] {
    data, err := os.ReadFile("/etc/app.conf")
    return R.TryCatchError(data, err)
}

// Right(data) if err == nil, Left(err) otherwise
```

### Recipe: Eitherize for zero-arg and multi-arg functions

**Problem**: The function has 0, 2, or more parameters.
**Solution**: Use `Eitherize0` through `Eitherize15`. E.g., `R.Eitherize0(os.Getwd)` gives `func() Result[string]`. `R.Eitherize2(f)` wraps a 2-param function.

---

## 2. Composing Fallible Operations

### Recipe: Pipe + Chain for sequential fallible steps

**Problem**: You have multiple operations that each can fail, where each step depends on the previous result.
**Solution**: Use `F.Pipe3` (or PipeN) with `R.Chain` to sequence them.

```go
package main

import (
    "fmt"
    "strconv"
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
)

func parseAndDouble(input string) R.Result[string] {
    return F.Pipe3(
        input,
        R.Eitherize1(strconv.Atoi),                        // string -> Result[int]
        R.Chain(func(n int) R.Result[int] {                 // int -> Result[int]
            if n < 0 {
                return R.Left[int](fmt.Errorf("negative: %d", n))
            }
            return R.Right(n * 2)
        }),
        R.Map(strconv.Itoa),                                // int -> string
    )
}

// parseAndDouble("21") => Right("42")
// parseAndDouble("-1") => Left(error: "negative: -1")
// parseAndDouble("xx") => Left(error: strconv parse error)
```

### Recipe: Map vs Chain rule of thumb

**Problem**: When to use `R.Map` vs `R.Chain`.
**Solution**: Use `R.Map(f)` when `f` is pure (`A -> B`). Use `R.Chain(f)` when `f` returns a `Result` (`A -> Result[B]`).

---

## 3. Building Reusable Pipelines

### Recipe: Flow for named, reusable transformation pipelines

**Problem**: You want to define a reusable pipeline as a named function without immediately applying it to data.
**Solution**: Use `F.Flow3` (or FlowN) to compose functions into a pipeline.

```go
package main

import (
    "strconv"
    "strings"
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
)

// Pipeline: string -> Result[int]
// 1. Trim whitespace
// 2. Parse int
// 3. Double the value
var cleanParseDouble = F.Flow3(
    strings.TrimSpace,                                 // string -> string
    R.Eitherize1(strconv.Atoi),                        // string -> Result[int]
    R.Map(func(n int) int { return n * 2 }),           // Result[int] -> Result[int]
)

// cleanParseDouble(" 21 ") => Right(42)
// cleanParseDouble("abc")  => Left(error)
```

### Recipe: Pipe vs Flow decision

**Problem**: Choosing between Pipe and Flow.
**Solution**: `F.PipeN(value, f1, f2, ...)` applies now. `F.FlowN(f1, f2, ...)` builds a reusable function for later.

---

## 4. Optional Value Handling

### Recipe: Creating and extracting Options

**Problem**: Represent values that may or may not exist, extract them safely.
**Solution**: Use `O.Some`, `O.None`, `O.GetOrElse`, `O.Fold`.

```go
package main

import (
    "fmt"
    O "github.com/IBM/fp-go/v2/option"
    F "github.com/IBM/fp-go/v2/function"
)

func example() {
    // Create Options
    some := O.Some(42)
    none := O.None[int]()

    // Extract with default
    val1 := F.Pipe1(some, O.GetOrElse(F.Constant(0)))  // 42
    val2 := F.Pipe1(none, O.GetOrElse(F.Constant(0)))  // 0

    // Pattern match with Fold
    msg := F.Pipe1(some, O.Fold(
        func() string { return "nothing" },
        func(n int) string { return fmt.Sprintf("got %d", n) },
    ))
    // msg == "got 42"

    _ = val1
    _ = val2
    _ = msg
}
```

### Recipe: FromPredicate to create Options conditionally

**Problem**: You have a value but want None if it does not meet a condition.
**Solution**: Use `O.FromPredicate` to create a guard.

```go
package main

import (
    O "github.com/IBM/fp-go/v2/option"
    F "github.com/IBM/fp-go/v2/function"
)

var nonEmpty = O.FromPredicate(func(s string) bool { return s != "" })

func example() {
    r1 := nonEmpty("hello")  // Some("hello")
    r2 := nonEmpty("")       // None
    _ = r1
    _ = r2

    // Chain: apply further Option-returning operations
    result := F.Pipe2(
        "hello",
        nonEmpty,
        O.Chain(func(s string) O.Option[int] {
            if len(s) > 3 {
                return O.Some(len(s))
            }
            return O.None[int]()
        }),
    )
    // result == Some(5)
    _ = result
}
```

### Recipe: FromNillable and Chain for pointer/optional flows

**Problem**: Convert pointers to Options and chain optional steps.
**Solution**: `O.FromNillable(ptr)` converts `*A` to `Option[*A]`. `O.Chain(f)` sequences optional ops.

```go
package main

import (
    "strings"
    F "github.com/IBM/fp-go/v2/function"
    O "github.com/IBM/fp-go/v2/option"
)

func findDomain(email string) O.Option[string] {
    return F.Pipe2(
        email,
        O.FromPredicate(func(e string) bool { return strings.Contains(e, "@") }),
        O.Map(func(e string) string { return strings.Split(e, "@")[1] }),
    )
}
// findDomain("a@b.com") => Some("b.com")
// findDomain("invalid") => None
```

---

## 5. Error Recovery

### Recipe: Alt for fallback to alternative Result

**Problem**: If a computation fails, try an alternative.
**Solution**: Use `R.Alt` which takes a lazy alternative.

```go
package main

import (
    "os"
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
)

var readFile = R.Eitherize1(os.ReadFile)

func readConfigWithFallback() R.Result[[]byte] {
    return F.Pipe1(
        readFile("/etc/app/config.yaml"),
        R.Alt(func() R.Result[[]byte] {
            return readFile("/etc/app/config.default.yaml")
        }),
    )
}
```

### Recipe: GetOrElse for extracting with a default

**Problem**: Extract the success value or compute a default from the error.
**Solution**: Use `R.GetOrElse(func(error) A)` which takes the error and returns a fallback value.

### Recipe: OrElse for conditional error recovery

**Problem**: Recover from specific errors while propagating others.
**Solution**: Use `R.OrElse` to inspect the error and decide.

```go
package main

import (
    "errors"
    "fmt"
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
)

var ErrNotFound = fmt.Errorf("not found")

func recoverNotFound(res R.Result[string]) R.Result[string] {
    return F.Pipe1(
        res,
        R.OrElse(func(err error) R.Result[string] {
            if errors.Is(err, ErrNotFound) {
                return R.Right("default_value")
            }
            return R.Left[string](err) // propagate other errors
        }),
    )
}
```

---

## 6. Adding Dependency Injection (Reader)

### Recipe: ReaderIOResult for context-dependent operations

**Problem**: Operations need `context.Context` and can fail.
**Solution**: Use `context/readerioresult`. Type: `func(context.Context) func() Result[A]`.

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
)

func fetchStatus(url string) RIOE.ReaderIOResult[int] {
    return func(ctx context.Context) func() R.Result[int] {
        return func() R.Result[int] {
            req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
            if err != nil { return R.Left[int](err) }
            resp, err := http.DefaultClient.Do(req)
            if err != nil { return R.Left[int](err) }
            defer resp.Body.Close()
            return R.Right(resp.StatusCode)
        }
    }
}

func example() {
    pipeline := F.Pipe1(
        fetchStatus("https://example.com"),
        RIOE.Map(func(code int) string { return fmt.Sprintf("status: %d", code) }),
    )
    result := pipeline(context.Background())()
    _ = result
}
```

---

## 7. Adding Typed DI (Effect)

### Recipe: Eitherize to convert Go functions to Effects

**Problem**: You have a function `func(Deps, context.Context) (T, error)` and want a typed Effect.
**Solution**: Use `effect.Eitherize` to lift it into `Effect[Deps, T]`.

```go
package main

import (
    "context"
    "fmt"
    F "github.com/IBM/fp-go/v2/function"
    "github.com/IBM/fp-go/v2/effect"
)

type AppConfig struct {
    DBURL string
}

type User struct {
    Name string
}

// Standard Go function
func fetchUser(cfg AppConfig, ctx context.Context) (*User, error) {
    // Imagine a real DB call here
    return &User{Name: "Alice"}, nil
}

// Convert to Effect[AppConfig, *User]
var fetchUserEffect = effect.Eitherize(fetchUser)

// Compose
var pipeline = F.Pipe1(
    fetchUserEffect,
    effect.Map[AppConfig](func(u *User) string { return u.Name }),
)

func main() {
    cfg := AppConfig{DBURL: "postgres://localhost"}
    thunk := effect.Provide[string](cfg)(pipeline)
    result := effect.RunSync(thunk)
    val, err := result(context.Background())
    fmt.Println(val, err) // "Alice" <nil>
}
```

### Recipe: Eitherize1 for parameterized Effect Kleisli arrows

**Problem**: You have `func(Deps, context.Context, Arg) (T, error)` and need a Kleisli arrow.
**Solution**: Use `effect.Eitherize1` to get `func(Arg) Effect[Deps, T]`.

```go
package main

import (
    "context"
    "fmt"
    F "github.com/IBM/fp-go/v2/function"
    "github.com/IBM/fp-go/v2/effect"
)

type DB struct{ connStr string }
type User struct{ ID int; Name string }

func getUserByID(db DB, ctx context.Context, id int) (*User, error) {
    return &User{ID: id, Name: "Bob"}, nil
}

// Kleisli[DB, int, *User] = func(int) Effect[DB, *User]
var getUserK = effect.Eitherize1(getUserByID)

var pipeline = F.Pipe1(
    effect.Succeed[DB](42),                // Effect[DB, int]
    effect.Chain[DB](getUserK),            // Effect[DB, *User]
)

func main() {
    db := DB{connStr: "postgres://..."}
    thunk := effect.Provide[*User](db)(pipeline)
    val, err := effect.RunSync(thunk)(context.Background())
    fmt.Println(val, err)
}
```

### Recipe: Provide + RunSync to execute Effects

**Problem**: You have an `Effect[C, A]` and need to run it.
**Solution**: `effect.Provide[A](c)(eff)` supplies context C returning a `ReaderIOResult[A]`. Then `effect.RunSync(thunk)` returns `func(context.Context) (A, error)`.

---

## 8. Functional Array Operations

### Recipe: Map, Filter, Reduce on arrays

**Problem**: Transform, filter, or fold arrays functionally.
**Solution**: Use `array.Map`, `array.Filter`, `array.Reduce`.

```go
package main

import (
    "fmt"
    F "github.com/IBM/fp-go/v2/function"
    A "github.com/IBM/fp-go/v2/array"
)

func example() {
    nums := []int{1, 2, 3, 4, 5}

    // Map: double each element
    doubled := F.Pipe1(nums, A.Map(func(n int) int { return n * 2 }))
    // [2, 4, 6, 8, 10]

    // Filter: keep evens
    evens := F.Pipe1(nums, A.Filter(func(n int) bool { return n%2 == 0 }))
    // [2, 4]

    // Reduce: sum
    sum := F.Pipe1(nums, A.Reduce(func(acc, n int) int { return acc + n }, 0))
    // 15

    fmt.Println(doubled, evens, sum)
}
```

### Recipe: TraverseArray for fallible array operations

**Problem**: Apply a function that returns `Result` to each element, short-circuiting on first error.
**Solution**: Use `result.TraverseArray`.

```go
package main

import (
    "strconv"
    "fmt"
    R "github.com/IBM/fp-go/v2/result"
)

var parseInt = R.Eitherize1(strconv.Atoi)

func parseAll(inputs []string) R.Result[[]int] {
    return R.TraverseArray(parseInt)(inputs)
}

func main() {
    r1 := parseAll([]string{"1", "2", "3"})
    fmt.Println(r1) // Right([1, 2, 3])

    r2 := parseAll([]string{"1", "bad", "3"})
    fmt.Println(r2) // Left(error)
}
```

### Recipe: SequenceArray to collect Results

**Problem**: You have `[]Result[A]` and want `Result[[]A]`.
**Solution**: `R.SequenceArray([]Result[A])` returns `Result[[]A]`. Fails on first Left.

---

## 9. Immutable Struct Updates with Lenses

### Recipe: Auto-generate lenses with `//go:generate`

**Problem**: You need lenses for struct fields but writing them by hand is tedious.
**Solution**: Annotate your struct with `// fp-go:Lens` and use `go generate`.

```go
// file: types.go
package mypackage

//go:generate go run github.com/IBM/fp-go/v2/main.go lens --dir . --filename gen_lens.go

// fp-go:Lens
type Person struct {
    Name  string
    Age   int
    Email string
}
```

Run `go generate ./...` to produce `gen_lens.go` containing:
- `PersonLenses` struct with a `Lens[Person, T]` for each field
- `MakePersonLenses()` constructor
- `PersonRefLenses` for pointer-based access
- `MakePersonRefLenses()` constructor

### Recipe: Using lenses for immutable updates

**Problem**: Update a field on a struct immutably.
**Solution**: Use `lens.Get`, `lens.Set`, `lens.Modify`.

```go
package main

import (
    L "github.com/IBM/fp-go/v2/optics/lens"
)

type Person struct { Name string; Age int }

var nameLens = L.MakeLens(
    func(p Person) string { return p.Name },
    func(p Person, name string) Person { p.Name = name; return p },
)
var ageLens = L.MakeLens(
    func(p Person) int { return p.Age },
    func(p Person, age int) Person { p.Age = age; return p },
)

func example() {
    alice := Person{Name: "Alice", Age: 30}
    name := nameLens.Get(alice)                               // "Alice"
    bob := nameLens.Set("Bob")(alice)                         // Person{Name: "Bob", Age: 30}
    older := L.Modify(ageLens, func(a int) int { return a+1 })(alice) // Age: 31
    _ = name; _ = bob; _ = older
}
```

---

## 10. Parallel Execution

### Recipe: TraverseArrayPar for concurrent operations

**Problem**: Execute an array of IO operations in parallel with `context.Context`.
**Solution**: Use `readerioresult.TraverseArrayPar`.

```go
package main

import (
    "context"
    "net/http"
    RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
    R "github.com/IBM/fp-go/v2/result"
)

func fetchStatus(url string) RIOE.ReaderIOResult[int] {
    return func(ctx context.Context) func() R.Result[int] {
        return func() R.Result[int] {
            req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
            if err != nil { return R.Left[int](err) }
            resp, err := http.DefaultClient.Do(req)
            if err != nil { return R.Left[int](err) }
            defer resp.Body.Close()
            return R.Right(resp.StatusCode)
        }
    }
}

// Parallel execution: all URLs fetched concurrently
var fetchAllStatuses = RIOE.TraverseArrayPar(fetchStatus)
// fetchAllStatuses(urls)(ctx)() => Result[[]int]
```

---

## 11. Retry with Backoff

### Recipe: Retry an Effect with exponential backoff

**Problem**: An operation may fail transiently and should be retried with configurable policy.
**Solution**: Use `retry.LimitRetries`, `retry.ExponentialBackoff`, and `effect.Retrying`.

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/IBM/fp-go/v2/effect"
    "github.com/IBM/fp-go/v2/either"
    M "github.com/IBM/fp-go/v2/monoid"
    "github.com/IBM/fp-go/v2/result"
    "github.com/IBM/fp-go/v2/retry"
)

type Deps struct{}

func unreliableOp(status retry.RetryStatus) effect.Effect[Deps, string] {
    // Simulate: fail on first 2 tries, succeed on 3rd
    if status.IterNumber < 2 {
        return effect.Fail[Deps, string](fmt.Errorf("transient error (attempt %d)", status.IterNumber))
    }
    return effect.Succeed[Deps]("success!")
}

func main() {
    // Policy: max 5 retries with exponential backoff starting at 100ms
    policy := M.Concat(
        retry.LimitRetries(5),
        retry.ExponentialBackoff(100*time.Millisecond),
    )(retry.Monoid)

    eff := effect.Retrying[Deps, string](
        policy,
        unreliableOp,
        func(res result.Result[string]) bool {
            return either.IsLeft(res) // retry on any error
        },
    )

    thunk := effect.Provide[string](Deps{})(eff)
    val, err := effect.RunSync(thunk)(context.Background())
    fmt.Println(val, err) // "success!" <nil>
}
```

### Recipe: CapDelay to limit maximum retry delay

**Problem**: Exponential backoff can grow too large.
**Solution**: Wrap with `retry.CapDelay(maxDelay, policy)` before combining.

---

## 12. HTTP Client Pipelines

### Recipe: Build HTTP requests with the builder pattern

**Problem**: Construct HTTP requests functionally with composable configuration.
**Solution**: Use `http/builder` package. Start with `B.Default`, chain `B.WithURL`, `B.WithMethod`, `B.WithHeader`, `B.WithJSON`. Convert to `ReaderIOResult[*http.Request]` via `RB.Requester(builder)`.

```go
builder := F.Pipe3(
    B.Default,
    B.WithURL("https://api.example.com/users"),
    B.WithMethod("POST"),
    B.WithHeader("Authorization")("Bearer my-token"),
)
requester := RB.Requester(builder) // ReaderIOResult[*http.Request]
```

### Recipe: Full HTTP pipeline with JSON parsing

**Problem**: Build request, send it, parse JSON response.
**Solution**: Compose `RB.Requester(builder)` with `RH.ReadJSON[T](client)`.

```go
package main

import (
    "net/http"
    F "github.com/IBM/fp-go/v2/function"
    B "github.com/IBM/fp-go/v2/http/builder"
    RB "github.com/IBM/fp-go/v2/context/readerioresult/http/builder"
    RH "github.com/IBM/fp-go/v2/context/readerioresult/http"
    RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
)

type APIResponse struct { ID int `json:"id"`; Name string `json:"name"` }

func fetchUser(id string) RIOE.ReaderIOResult[APIResponse] {
    builder := F.Pipe2(
        B.Default,
        B.WithURL("https://api.example.com/users/"+id),
        B.WithHeader("Accept")("application/json"),
    )
    client := RH.MakeClient(http.DefaultClient)
    return F.Pipe1(RB.Requester(builder), RH.ReadJSON[APIResponse](client))
}
```

---

## 13. Folding / Extracting Values

### Recipe: R.Fold to pattern-match on Results

**Problem**: You need to handle both success and error cases, producing a common type.
**Solution**: Use `R.Fold` (or `E.Fold` for Either, `O.Fold` for Option).

```go
package main

import (
    "fmt"
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
)

func resultToMessage(res R.Result[int]) string {
    return F.Pipe1(res, R.Fold(
        func(err error) string { return "Error: " + err.Error() },
        func(n int) string { return fmt.Sprintf("Value: %d", n) },
    ))
}
```

### Recipe: O.Fold for Option pattern matching

**Problem**: Handle Some and None cases, producing a value.
**Solution**: `O.Fold(onNone, onSome)` -- `onNone` is `func() B`, `onSome` is `func(A) B`.

### Recipe: E.Fold for general Either

**Problem**: Pattern match on an Either with custom left type (not error).
**Solution**: `E.Fold(onLeft, onRight)` works the same as `R.Fold` but for `Either[E, A]` with any left type E.

---

## 14. Creating Project-Wide Type Aliases

### Recipe: Generic type aliases for cleaner signatures (Go 1.24+)

**Problem**: `result.Result[A]`, `effect.Effect[C, A]` are verbose in signatures.
**Solution**: Create project-level generic type aliases.

```go
package types

import (
    "github.com/IBM/fp-go/v2/effect"
    "github.com/IBM/fp-go/v2/result"
)

type (
    Result[A any]    = result.Result[A]
    AppEffect[A any] = effect.Effect[AppConfig, A]
    AppConfig struct { DBURL string; APIKey string }
)
```

---

## 15. Converting Between Standard and Idiomatic

### Recipe: Unwrap Result to (value, error) tuple

**Problem**: You need to call a standard Go API that expects `(T, error)`.
**Solution**: `val, err := R.Unwrap(result)` converts `Result[A]` to `(A, error)`.

### Recipe: Uneitherize to convert back to Go function signatures

**Problem**: You have a function returning `Result[A]` and need `func(X) (A, error)`.
**Solution**: `R.Uneitherize1(f)` is the inverse of `R.Eitherize1`. Also available: `Uneitherize0` through `Uneitherize15`.

### Recipe: Convert between Result and Option

**Problem**: Discard error to get Option, or convert Option to Result.
**Solution**: `R.ToOption(result)` discards Left, gives `Option`. `R.FromOption[A](onNone)(opt)` gives `Result`, using `onNone()` for the error when None.

---

## 16. Validation with Codecs

### Recipe: Using built-in codecs for type validation

**Problem**: Runtime type validation with detailed error reporting.
**Solution**: Use `optics/codec` pre-built codecs.

```go
intCodec := codec.Int()
r1 := intCodec.Decode(42)       // Right(42) -- Validation[int]
r2 := intCodec.Decode("hello")  // Left(validation errors)
s := intCodec.Encode(42)        // 42 (always succeeds)
check := intCodec.Is(42)        // Right(42) -- type check
```

### Recipe: URL and other built-in codecs

Available codecs: `codec.Int()`, `codec.String()`, `codec.Bool()`, `codec.URL()`, `codec.Array(itemCodec)`.
Each provides `.Decode(input)` returning `Validation[A]`, `.Encode(a)` returning the output type, and `.Is(any)` for type checking.

---

## 17. Do-Notation for Complex Flows

### Recipe: Do + Bind + Let for readable monadic chains in Effect

**Problem**: Complex multi-step effects are hard to read with nested Chain calls.
**Solution**: Use `effect.Do`, `effect.Bind`, `effect.Let` with a state struct.

```go
package main

import (
    "context"
    "fmt"
    F "github.com/IBM/fp-go/v2/function"
    "github.com/IBM/fp-go/v2/effect"
)

type Deps struct{}

type PipelineState struct {
    UserID   int
    UserName string
    Greeting string
}

func lookupUser(id int) effect.Effect[Deps, string] {
    return effect.Succeed[Deps](fmt.Sprintf("User_%d", id))
}

func pipeline() effect.Effect[Deps, PipelineState] {
    return F.Pipe3(
        // Start with initial state
        effect.Do[Deps](PipelineState{UserID: 42}),

        // Bind: run an effect and merge result into state
        effect.Bind(
            func(name string) func(PipelineState) PipelineState {
                return func(s PipelineState) PipelineState {
                    s.UserName = name
                    return s
                }
            },
            func(s PipelineState) effect.Effect[Deps, string] {
                return lookupUser(s.UserID)
            },
        ),

        // Let: pure computation on state
        effect.Let[Deps](
            func(greeting string) func(PipelineState) PipelineState {
                return func(s PipelineState) PipelineState {
                    s.Greeting = greeting
                    return s
                }
            },
            func(s PipelineState) string {
                return fmt.Sprintf("Hello, %s!", s.UserName)
            },
        ),
    )
}

func main() {
    eff := pipeline()
    thunk := effect.Provide[PipelineState](Deps{})(eff)
    val, err := effect.RunSync(thunk)(context.Background())
    fmt.Println(val, err)
    // {42 User_42 Hello, User_42!} <nil>
}
```

### Recipe: Do-notation with lenses (BindL, LetL, ApSL)

**Problem**: Manual setter functions in Bind/Let are verbose.
**Solution**: Use lens-based variants with auto-generated lenses.

```go
// Given lenses countLens and messageLens for a State struct:
eff := F.Pipe3(
    effect.Do[Ctx](State{}),
    effect.ApSL(countLens, effect.Succeed[Ctx](10)),     // set Count = 10
    effect.LetL[Ctx](countLens, func(n int) int { return n * 2 }), // Count *= 2
    // effect.BindL(msgLens, func(s string) effect.Effect[Ctx, string] { ... }),
)
```

---

## 18. ChainFirst / Tap for Side Effects

### Recipe: ChainFirst for logging without disrupting pipeline

**Problem**: You need to log intermediate values in a Result pipeline without changing the data flow.
**Solution**: Use `R.ChainFirst` which executes a side-effect but preserves the original value.

```go
package main

import (
    "fmt"
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
)

func processWithLogging(input string) R.Result[int] {
    return F.Pipe3(
        R.Of(input),
        R.Map(func(s string) int { return len(s) }),

        // Log the value without changing it
        R.ChainFirst(func(n int) R.Result[string] {
            fmt.Printf("DEBUG: length = %d\n", n)
            return R.Of("logged") // return value is discarded
        }),

        R.Map(func(n int) int { return n * 2 }),
    )
}
```

### Recipe: ChainFirst / TapThunkK in Option and Effect

**Problem**: Side effects at the Option or Effect level.
**Solution**: `O.ChainFirst(f)` for Options; `effect.ChainFirst[C](f)` for full effects; `effect.TapThunkK[C](f)` for context-independent IO side effects that preserve the original value.

---

## 19. FromPredicate for Guards

### Recipe: Option.FromPredicate for conditional wrapping

**Problem**: Turn a value into None if it does not satisfy a predicate.
**Solution**: `O.FromPredicate(pred)` returns `func(A) Option[A]`. Returns `Some(a)` if `pred(a)`, else `None`.

### Recipe: Result.FromPredicate for guards with error messages

**Problem**: Turn a value into an error if it does not satisfy a predicate, with a specific error.
**Solution**: Use `R.FromPredicate`.

```go
package main

import (
    "fmt"
    R "github.com/IBM/fp-go/v2/result"
)

var validatePort = R.FromPredicate(
    func(port int) bool { return port > 0 && port <= 65535 },
    func(port int) error { return fmt.Errorf("invalid port: %d", port) },
)

func example() {
    r1 := validatePort(8080)   // Right(8080)
    r2 := validatePort(-1)     // Left(error: "invalid port: -1")
    r3 := validatePort(99999)  // Left(error: "invalid port: 99999")
    _ = r1
    _ = r2
    _ = r3
}
```

### Recipe: Chaining multiple guards in a pipeline

**Problem**: Apply multiple validation guards in sequence.
**Solution**: Chain `R.FromPredicate` calls with `R.Chain`. Each guard short-circuits on failure.

```go
package main

import (
    "fmt"
    "strings"
    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
)

func validateUsername(input string) R.Result[string] {
    return F.Pipe3(
        R.Of(input),
        R.Chain(R.FromPredicate(
            func(s string) bool { return len(s) >= 3 },
            func(s string) error { return fmt.Errorf("too short: %q", s) },
        )),
        R.Chain(R.FromPredicate(
            func(s string) bool { return len(s) <= 20 },
            func(s string) error { return fmt.Errorf("too long: %q", s) },
        )),
        R.Chain(R.FromPredicate(
            func(s string) bool { return !strings.Contains(s, " ") },
            func(s string) error { return fmt.Errorf("has spaces: %q", s) },
        )),
    )
}
```

---

## 20. Step-by-Step Migration

### Recipe: Converting imperative fetch-validate-save to fp-go pipeline

**Problem**: You have typical imperative Go code doing fetch -> validate -> save with error checking at each step.
**Solution**: Wrap each step with `Eitherize1` or manual Result construction, compose with `Pipe` + `Chain`.

#### fp-go pipeline version:

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"

    F "github.com/IBM/fp-go/v2/function"
    R "github.com/IBM/fp-go/v2/result"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Step 1: Wrap each Go function that returns (T, error)
var fetchUserFP = R.Eitherize1(func(id int) (*User, error) {
    resp, err := http.Get(fmt.Sprintf("https://api.example.com/users/%d", id))
    if err != nil {
        return nil, fmt.Errorf("fetch failed: %w", err)
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read body failed: %w", err)
    }
    var user User
    if err := json.Unmarshal(body, &user); err != nil {
        return nil, fmt.Errorf("parse failed: %w", err)
    }
    return &user, nil
})

// Step 2: Convert validation to a Kleisli arrow using FromError
var validateUserFP = R.FromError(func(user *User) error {
    if user.Name == "" {
        return fmt.Errorf("name is required")
    }
    if user.Email == "" {
        return fmt.Errorf("email is required")
    }
    return nil
})

// Step 3: Wrap the save function
var saveUserFP = R.Eitherize1(func(user *User) (*User, error) {
    // Imagine DB save
    return user, nil
})

// Step 4: Compose into a pipeline
func processUser(id int) R.Result[*User] {
    return F.Pipe3(
        id,
        fetchUserFP,                // int -> Result[*User]
        R.Chain(validateUserFP),     // Result[*User] -> Result[*User]
        R.Chain(saveUserFP),         // Result[*User] -> Result[*User]
    )
}

// Step 5: At the boundary, unwrap back to Go if needed
func processUserGo(id int) (*User, error) {
    return R.Unwrap(processUser(id))
}
```

**Migration checklist**: (1) Wrap `(T,error)` funcs with `R.EitherizeN` (2) Wrap `func(T) error` validators with `R.FromError` (3) Replace `if err != nil` chains with `R.Chain`/`R.Map` (4) Add logging with `R.ChainFirst` (5) Add fallbacks with `R.OrElse`/`R.Alt` (6) Unwrap at boundaries with `R.Unwrap`.

### Recipe: Incremental migration strategy

**Problem**: You cannot convert the entire codebase at once.
**Solution**: Migrate inside-out: wrap leaf functions with `Eitherize`, compose internally with Pipe/Chain, and `Unwrap` at public API boundaries so callers see standard `(T, error)`.

---

## Quick Reference: Key Functions by Package

### function (alias F)

| Function | Signature | Purpose |
|----------|-----------|---------|
| `Pipe1` | `(T0, func(T0)T1) T1` | Apply 1 transformation |
| `Pipe2` | `(T0, f1, f2) T2` | Apply 2 transformations |
| `Pipe3` | `(T0, f1, f2, f3) T3` | Apply 3 transformations |
| `Flow2` | `(f1, f2) func(T0)T2` | Compose 2 functions |
| `Flow3` | `(f1, f2, f3) func(T0)T3` | Compose 3 functions |
| `Identity` | `(A) A` | Return argument unchanged |
| `Constant` | `(A) func()A` | Nullary function returning A |
| `Constant1` | `(A) func(B)A` | Ignore input, return A |

### result (alias R) -- Result[A] = Either[error, A]

| Function | Purpose |
|----------|---------|
| `Right(a)` | Wrap success value |
| `Left[A](err)` | Wrap error |
| `Of(a)` | Alias for Right |
| `Eitherize1(f)` | Wrap `func(X)(Y,error)` |
| `TryCatchError(val,err)` | Wrap `(value, error)` pair |
| `Map(f)` | Transform success value |
| `Chain(f)` | Sequence fallible operations |
| `ChainFirst(f)` | Side effect, keep original |
| `Fold(onErr, onOk)` | Pattern match |
| `GetOrElse(onErr)` | Extract with default |
| `Alt(lazy)` | Try alternative on failure |
| `OrElse(f)` | Conditional error recovery |
| `FromPredicate(pred, onFalse)` | Guard with error |
| `FromOption[A](onNone)` | Option -> Result |
| `ToOption(r)` | Result -> Option |
| `Unwrap(r)` | Result -> `(A, error)` |
| `TraverseArray(f)` | Map array with fallible f |
| `SequenceArray(rs)` | `[]Result[A]` -> `Result[[]A]` |

### option (alias O) -- Option[A]

| Function | Purpose |
|----------|---------|
| `Some(a)` | Wrap a value |
| `None[A]()` | Empty option |
| `Map(f)` | Transform Some value |
| `Chain(f)` | Sequence optional ops |
| `ChainFirst(f)` | Side effect, keep original |
| `Fold(onNone, onSome)` | Pattern match |
| `GetOrElse(onNone)` | Extract with default |
| `Alt(lazy)` | Try alternative on None |
| `Filter(pred)` | Keep if predicate passes |
| `FromPredicate(pred)` | Guard to Option |
| `FromNillable(ptr)` | Pointer -> Option |

### effect -- Effect[C, A]

| Function | Purpose |
|----------|---------|
| `Succeed[C](a)` | Lift pure value |
| `Fail[C,A](err)` | Create failed effect |
| `Map[C](f)` | Transform success |
| `Chain[C](f)` | Sequence effects |
| `ChainFirst[C](f)` | Side effect |
| `Eitherize(f)` | Wrap `func(C, ctx)(T, error)` |
| `Eitherize1(f)` | Wrap `func(C, ctx, A)(T, error)` |
| `Provide[A](c)` | Supply context C |
| `RunSync(thunk)` | Execute -> `func(ctx)(A, error)` |
| `Do[C](s)` | Start do-notation |
| `Bind(setter, f)` | Bind effect result |
| `Let[C](setter, f)` | Bind pure computation |
| `Retrying[C](policy, action, check)` | Retry with policy |
| `TraverseArray[C](f)` | Map array with effects |
| `Ask[C]()` | Access own context |
| `Local[A,C1,C2](f)` | Transform context |
