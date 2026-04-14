# fp-go/v2 Core Patterns Reference

> Module: `github.com/IBM/fp-go/v2`
> This document covers all core types, operations, and composition patterns.
> Audience: LLM agents and experienced FP practitioners.

---

## 1. Type Hierarchy

```
Option[A]                       -- may be absent
Either[E, A]                    -- can fail with E
Result[A]                       -- Either[error, A]
IO[A]                           -- lazy synchronous computation (func() A)
IOResult[A]                     -- IO[Result[A]] = func() Either[error, A]
Reader[R, A]                    -- depends on R (func(R) A)
ReaderIOResult[R, A]            -- Reader[R, IOResult[A]] = func(R) func() Either[error, A]
Effect[C, A]                    -- Reader[C, ReaderIOResult[context.Context, A]]
                                   typed DI + context.Context + IO + error

Layers (bottom to top):
  Pure values
    |
  Option[A]              -- optionality
    |
  Either[E, A]           -- typed error channel
    |
  Result[A]              -- Either[error, A]  (Go-idiomatic error)
    |
  IO[A]                  -- deferred side effects
    |
  IOResult[A]            -- IO + error
    |
  Reader[R, A]           -- dependency injection
    |
  ReaderIOResult[R, A]   -- Reader + IO + error
    |
  Effect[C, A]           -- Reader[C, ctx->IO[Result[A]]]
                            full effect system with DI, context.Context, IO, error
```

The `effect` package is the recommended top-level type for application code.
Use `Option`/`Either`/`Result` for pure data transformations.

---

## 2. Core Types

### 2.1 Option

**Package**: `github.com/IBM/fp-go/v2/option`

**Type definition**:
```go
type Option[A any] struct {
    value  A
    isSome bool
}
```

**Kleisli and Operator types**:
```go
type Kleisli[A, B any]  = func(A) Option[B]
type Operator[A, B any] = Kleisli[Option[A], B]   // func(Option[A]) Option[B]
```

**Construction**:
```go
func Some[T any](value T) Option[T]
func None[T any]() Option[T]
func Of[T any](value T) Option[T]                         // alias for Some
func Zero[A any]() Option[A]                               // alias for None
func FromPredicate[A any](pred func(A) bool) Kleisli[A, A]
func FromNonZero[A comparable]() Kleisli[A, A]
func FromNillable[A any](a *A) Option[*A]
func FromValidation[A, B any](f func(A) (B, bool)) Kleisli[A, B]
func TryCatch[A any](f func() (A, error)) Option[A]
```

**Inspection**:
```go
func IsSome[T any](val Option[T]) bool
func IsNone[T any](val Option[T]) bool
func Unwrap[A any](ma Option[A]) (A, bool)
```

**Transformation**:
```go
func Map[A, B any](f func(A) B) Operator[A, B]
func MapTo[A, B any](b B) Operator[A, B]
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainTo[A, B any](mb Option[B]) Operator[A, B]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func Ap[B, A any](fa Option[A]) Operator[func(A) B, B]
func Flatten[A any](mma Option[Option[A]]) Option[A]
func Flap[B, A any](a A) Operator[func(A) B, B]
```

**Extraction**:
```go
func Fold[A, B any](onNone func() B, onSome func(A) B) func(Option[A]) B
func MonadFold[A, B any](ma Option[A], onNone func() B, onSome func(A) B) B
func GetOrElse[A any](onNone func() A) func(Option[A]) A
func Reduce[A, B any](f func(B, A) B, initial B) func(Option[A]) B
```

**Filtering**:
```go
func Filter[A any](pred func(A) bool) Operator[A, A]
```

**Alternative**:
```go
func Alt[A any](that func() Option[A]) Operator[A, A]
```

**Sequence**:
```go
func Sequence2[T1, T2, R any](f func(T1, T2) Option[R]) func(Option[T1], Option[T2]) Option[R]
```

**Example**:
```go
import (
    O "github.com/IBM/fp-go/v2/option"
    F "github.com/IBM/fp-go/v2/function"
)

result := F.Pipe3(
    O.Some(42),
    O.Filter(func(x int) bool { return x > 0 }),
    O.Map(func(x int) string { return fmt.Sprintf("value: %d", x) }),
    O.GetOrElse(func() string { return "no value" }),
)
// result == "value: 42"
```

---

### 2.2 Either

**Package**: `github.com/IBM/fp-go/v2/either`

**Type definition**:
```go
type Either[E, A any] struct {
    r      A
    l      E
    isLeft bool
}
```

**Kleisli and Operator types**:
```go
type Kleisli[E, A, B any]  = func(A) Either[E, B]
type Operator[E, A, B any] = Kleisli[E, Either[E, A], B]
```

**Construction**:
```go
func Right[E, A any](value A) Either[E, A]
func Left[A, E any](value E) Either[E, A]
func Of[E, A any](value A) Either[E, A]              // alias for Right
func Zero[E, A any]() Either[E, A]                    // Right with zero A
func FromPredicate[E, A any](pred func(A) bool, onFalse func(A) E) Kleisli[E, A, A]
func FromNillable[A, E any](e E) Kleisli[E, *A, *A]
func FromOption[A, E any](onNone func() E) func(Option[A]) Either[E, A]
func FromError[A any](f func(A) error) func(A) Either[error, A]
func TryCatch[FE func(error) E, E, A any](val A, err error, onThrow FE) Either[E, A]
func TryCatchError[A any](val A, err error) Either[error, A]
```

**Inspection**:
```go
func IsLeft[E, A any](val Either[E, A]) bool
func IsRight[E, A any](val Either[E, A]) bool
func Unwrap[E, A any](ma Either[E, A]) (A, E)
func UnwrapError[A any](ma Either[error, A]) (A, error)
```

**Transformation**:
```go
func Map[E, A, B any](f func(A) B) Operator[E, A, B]
func MapTo[E, A, B any](b B) Operator[E, A, B]
func MapLeft[A, E1, E2 any](f func(E1) E2) func(Either[E1, A]) Either[E2, A]
func BiMap[E1, E2, A, B any](f func(E1) E2, g func(A) B) func(Either[E1, A]) Either[E2, B]
func Chain[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, B]
func ChainLeft[EA, EB, A any](f Kleisli[EB, EA, A]) Kleisli[EB, Either[EA, A], A]
func ChainFirst[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, A]
func ChainTo[A, E, B any](mb Either[E, B]) Operator[E, A, B]
func ChainOptionK[A, B, E any](onNone func() E) func(func(A) Option[B]) Operator[E, A, B]
func Ap[B, E, A any](fa Either[E, A]) Operator[E, func(A) B, B]
func Flatten[E, A any](mma Either[E, Either[E, A]]) Either[E, A]
func Swap[E, A any](val Either[E, A]) Either[A, E]
func Flap[E, B, A any](a A) Operator[E, func(A) B, B]
```

**Extraction**:
```go
func Fold[E, A, B any](onLeft func(E) B, onRight func(A) B) func(Either[E, A]) B
func MonadFold[E, A, B any](ma Either[E, A], onLeft func(E) B, onRight func(A) B) B
func GetOrElse[E, A any](onLeft func(E) A) func(Either[E, A]) A
func Reduce[E, A, B any](f func(B, A) B, initial B) func(Either[E, A]) B
func ToOption[E, A any](ma Either[E, A]) Option[A]
func ToError[A any](e Either[error, A]) error
```

**Alternative**:
```go
func Alt[E, A any](that func() Either[E, A]) Operator[E, A, A]
func AltW[E, E1, A any](that func() Either[E1, A]) Kleisli[E1, Either[E, A], A]
func OrElse[E1, E2, A any](onLeft Kleisli[E2, E1, A]) Kleisli[E2, Either[E1, A], A]
```

**Sequence**:
```go
func Sequence2[E, T1, T2, R any](f func(T1, T2) Either[E, R]) func(Either[E, T1], Either[E, T2]) Either[E, R]
func Sequence3[E, T1, T2, T3, R any](f func(T1, T2, T3) Either[E, R]) func(...) Either[E, R]
```

**Example**:
```go
import (
    E "github.com/IBM/fp-go/v2/either"
    F "github.com/IBM/fp-go/v2/function"
)

result := F.Pipe3(
    E.Right[error](42),
    E.Map[error](func(x int) int { return x * 2 }),
    E.Chain(func(x int) E.Either[error, string] {
        if x > 0 {
            return E.Right[error](fmt.Sprintf("%d", x))
        }
        return E.Left[string](errors.New("non-positive"))
    }),
    E.GetOrElse(func(err error) string { return "error: " + err.Error() }),
)
// result == "84"
```

---

### 2.3 Result

**Package**: `github.com/IBM/fp-go/v2/result`

**Type definition** (type alias):
```go
type Result[T any] = Either[error, T]
type Kleisli[A, B any] = func(A) Result[B]
type Operator[A, B any] = Kleisli[Result[A], B]
```

Result is `Either[error, A]` -- every function mirrors `either` with `E` fixed to `error`.

**Key functions** (all delegate to `either`):
```go
func Of[A any](value A) Result[A]
func Right[A any](value A) Result[A]
func Left[A any](value error) Result[A]
func Map[A, B any](f func(A) B) Operator[A, B]
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func ChainLeft[A any](f Kleisli[error, A]) Operator[A, A]
func Ap[B, A any](fa Result[A]) Operator[func(A) B, B]
func Fold[A, B any](onLeft func(error) B, onRight func(A) B) func(Result[A]) B
func GetOrElse[A any](onLeft func(error) A) func(Result[A]) A
func Unwrap[A any](ma Result[A]) (A, error)
func UnwrapError[A any](ma Result[A]) (A, error)
func TryCatchError[A any](val A, err error) Result[A]
func FromError[A any](f func(A) error) Kleisli[A, A]
func FromOption[A any](onNone func() error) func(Option[A]) Result[A]
func ToOption[A any](ma Result[A]) Option[A]
func ToError[A any](e Result[A]) error
func Alt[A any](that func() Result[A]) Operator[A, A]
func OrElse[A any](onLeft Kleisli[error, A]) Operator[A, A]
func Flatten[A any](mma Result[Result[A]]) Result[A]
func InstanceOf[A any](a any) Result[A]
```

**Example**:
```go
import (
    R "github.com/IBM/fp-go/v2/result"
    F "github.com/IBM/fp-go/v2/function"
)

val, err := R.Unwrap(F.Pipe2(
    R.TryCatchError(strconv.Atoi("42")),
    R.Map(func(n int) int { return n * 2 }),
))
// val == 84, err == nil
```

---

### 2.4 IO

**Package**: `github.com/IBM/fp-go/v2/io`

**Type definition** (type alias):
```go
type IO[A any] = func() A
type Kleisli[A, B any] = func(A) IO[B]
type Operator[A, B any] = Kleisli[IO[A], B]
```

IO represents a lazy synchronous computation. Nothing executes until `io()` is called.

**Construction**:
```go
func Of[A any](a A) IO[A]                            // wraps pure value
func FromIO[A any](a IO[A]) IO[A]                    // identity
func FromImpure[ANY ~func()](f ANY) IO[Void]          // void side effect
func Defer[A any](gen func() IO[A]) IO[A]            // lazy IO creation
func Memoize[A any](ma IO[A]) IO[A]                  // compute once
var Now IO[time.Time]                                  // current timestamp
```

**Transformation**:
```go
func Map[A, B any](f func(A) B) Operator[A, B]
func MapTo[A, B any](b B) Operator[A, B]
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func ChainTo[A, B any](fb IO[B]) Operator[A, B]
func Ap[B, A any](ma IO[A]) Operator[func(A) B, B]          // parallel by default
func ApSeq[B, A any](ma IO[A]) Operator[func(A) B, B]       // sequential
func ApPar[B, A any](ma IO[A]) Operator[func(A) B, B]       // explicitly parallel
func Flatten[A any](mma IO[IO[A]]) IO[A]
func Flap[B, A any](a A) Operator[func(A) B, B]
```

**Timing**:
```go
func Delay[A any](delay time.Duration) Operator[A, A]
func After[A any](timestamp time.Time) Operator[A, A]
func WithTime[A any](a IO[A]) IO[Pair[Pair[time.Time, time.Time], A]]
func WithDuration[A any](a IO[A]) IO[Pair[time.Duration, A]]
```

**Example**:
```go
import (
    IO "github.com/IBM/fp-go/v2/io"
    F "github.com/IBM/fp-go/v2/function"
)

greeting := F.Pipe2(
    IO.Of("World"),
    IO.Map(func(s string) string { return "Hello, " + s + "!" }),
    IO.Delay[string](time.Second),
)
result := greeting() // waits 1 second, returns "Hello, World!"
```

---

### 2.5 IOResult

**Package**: `github.com/IBM/fp-go/v2/ioresult`

**Type definition** (type alias):
```go
type IOResult[A any] = IO[Result[A]]    // func() Either[error, A]
type Kleisli[A, B any] = func(A) IOResult[B]
type Operator[A, B any] = Kleisli[IOResult[A], B]
```

IOResult combines IO with Result. The computation is deferred and may fail.

**Key functions** (same pattern as other monads):
```go
func Of[A any](a A) IOResult[A]
func Left[A any](e error) IOResult[A]
func Right[A any](a A) IOResult[A]
func FromResult[A any](e Result[A]) IOResult[A]
func TryCatchError[A any](f func() (A, error)) IOResult[A]
func Map[A, B any](f func(A) B) Operator[A, B]
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B]
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A]
func Ap[B, A any](fa IOResult[A]) Operator[func(A) B, B]
func Fold[A, B any](onLeft func(error) IO[B], onRight func(A) IO[B]) func(IOResult[A]) IO[B]
func Alt[A any](that func() IOResult[A]) Operator[A, A]
func Flatten[A any](mma IOResult[IOResult[A]]) IOResult[A]
func Memoize[A any](ma IOResult[A]) IOResult[A]
```

---

### 2.6 Reader

**Package**: `github.com/IBM/fp-go/v2/reader`

**Type definition** (type alias):
```go
type Reader[R, A any] = func(R) A
type Kleisli[R, A, B any] = func(A) Reader[R, B]
type Operator[R, A, B any] = Kleisli[R, Reader[R, A], B]
```

Reader encodes dependency injection. `R` is the shared environment.

**Construction**:
```go
func Ask[R any]() Reader[R, R]                           // returns environment itself
func Asks[R, A any](f Reader[R, A]) Reader[R, A]         // project from environment
func Of[R, A any](a A) Reader[R, A]                      // constant, ignores environment
func OfLazy[R, A any](fa func() A) Reader[R, A]          // deferred constant
```

**Transformation**:
```go
func Map[E, A, B any](f func(A) B) Operator[E, A, B]
func MapTo[E, A, B any](b B) Operator[E, A, B]
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B]
func ChainTo[A, R, B any](b Reader[R, B]) Operator[R, A, B]
func Ap[B, R, A any](fa Reader[R, A]) Operator[R, func(A) B, B]
func Flatten[R, A any](mma Reader[R, Reader[R, A]]) Reader[R, A]
func Compose[C, R, B any](ab Reader[R, B]) Kleisli[R, Reader[B, C], C]
func Flap[R, B, A any](a A) Operator[R, func(A) B, B]
```

**Execution**:
```go
func Read[A, E any](e E) func(Reader[E, A]) A            // run with environment
```

**Key distinction -- Chain vs Compose**:
- `Chain`: Both readers share the SAME environment R. Second depends on the VALUE of the first.
- `Compose`: First reader's output becomes the INPUT environment of the second.

---

### 2.7 ReaderIOResult

**Package**: `github.com/IBM/fp-go/v2/readerioresult`

**Type definition**:
```go
type ReaderIOResult[R, A any] = Reader[R, IOResult[A]]
// Expanded: func(R) func() Either[error, A]

type Kleisli[R, A, B any] = func(A) ReaderIOResult[R, B]
type Operator[R, A, B any] = Kleisli[R, ReaderIOResult[R, A], B]
```

Combines Reader (dependency injection) + IO (side effects) + Result (error handling).
Same API pattern: Of, Map, Chain, ChainFirst, Ap, Fold, Alt, etc.

---

### 2.8 Effect

**Package**: `github.com/IBM/fp-go/v2/effect`

**Type definition**:
```go
// Effect[C, A] = Reader[C, Reader[context.Context, IO[Result[A]]]]
// Expanded: func(C) func(context.Context) func() Either[error, A]
type Effect[C, A any] = readerreaderioresult.ReaderReaderIOResult[C, A]

type Thunk[A any] = ReaderIOResult[A]   // context.Context -> IO[Result[A]]
type Kleisli[C, A, B any] = func(A) Effect[C, B]
type Operator[C, A, B any] = func(Effect[C, A]) Effect[C, B]
```

Effect is the recommended top-level monad. `C` is the application context (config, services), and `context.Context` is automatically threaded.

**Construction**:
```go
func Of[C, A any](a A) Effect[C, A]
func Succeed[C, A any](a A) Effect[C, A]              // alias for Of
func Fail[C, A any](err error) Effect[C, A]
func FromThunk[C, A any](f Thunk[A]) Effect[C, A]
func FromResult[C, A any](r Result[A]) Effect[C, A]
func Asks[C, A any](r Reader[C, A]) Effect[C, A]      // project from context C
func Suspend[C, A any](fa func() Effect[C, A]) Effect[C, A]
```

**Transformation**:
```go
func Map[C, A, B any](f func(A) B) Operator[C, A, B]
func Chain[C, A, B any](f Kleisli[C, A, B]) Operator[C, A, B]
func ChainFirst[C, A, B any](f Kleisli[C, A, B]) Operator[C, A, A]
func Tap[C, A, ANY any](f Kleisli[C, A, ANY]) Operator[C, A, A]
func Ap[B, C, A any](fa Effect[C, A]) Operator[C, func(A) B, B]
```

**Lifting from other monads**:
```go
func ChainIOK[C, A, B any](f io.Kleisli[A, B]) Operator[C, A, B]
func ChainFirstIOK[C, A, B any](f io.Kleisli[A, B]) Operator[C, A, A]
func TapIOK[C, A, B any](f io.Kleisli[A, B]) Operator[C, A, A]
func ChainResultK[C, A, B any](f result.Kleisli[A, B]) Operator[C, A, B]
func ChainReaderK[C, A, B any](f reader.Kleisli[C, A, B]) Operator[C, A, B]
func ChainReaderIOK[C, A, B any](f readerio.Kleisli[C, A, B]) Operator[C, A, B]
func ChainThunkK[C, A, B any](f thunk.Kleisli[A, B]) Operator[C, A, B]
func ChainFirstThunkK[C, A, B any](f thunk.Kleisli[A, B]) Operator[C, A, A]
func TapThunkK[C, A, B any](f thunk.Kleisli[A, B]) Operator[C, A, A]
```

**Context manipulation (dependency injection)**:
```go
func Ask[C any]() Effect[C, C]                                           // read entire context
func Asks[C, A any](r Reader[C, A]) Effect[C, A]                        // project from context
func Local[A, C1, C2 any](acc Reader[C1, C2]) Kleisli[C1, Effect[C2, A], A]   // transform context (pure)
func Contramap[A, C1, C2 any](acc Reader[C1, C2]) Kleisli[C1, Effect[C2, A], A] // alias for Local
func LocalIOK[A, C1, C2 any](f io.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A]         // IO-based transform
func LocalIOResultK[A, C1, C2 any](f ioresult.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A] // IOResult-based
func LocalResultK[A, C1, C2 any](f result.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A] // Result-based
func LocalThunkK[A, C1, C2 any](f thunk.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A]   // Thunk-based
func LocalEffectK[A, C1, C2 any](f Kleisli[C2, C2, C1]) func(Effect[C1, A]) Effect[C2, A]    // Effect-based
func LocalReaderK[A, C1, C2 any](f reader.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A] // Reader-based
```

Context transformation strength (weakest to strongest):
1. `Local` / `Contramap` -- pure function `C2 -> C1`
2. `LocalResultK` -- may fail `C2 -> Result[C1]`
3. `LocalIOK` -- IO side effects `C2 -> IO[C1]`
4. `LocalIOResultK` -- IO + error `C2 -> IOResult[C1]`
5. `LocalReaderK` -- pure + runtime context `C2 -> Reader[ctx, C1]`
6. `LocalThunkK` -- runtime context + IO + error `C2 -> Thunk[C1]`
7. `LocalEffectK` -- full effect `C2 -> Effect[C2, C1]`

**Filtering within Effects**:
```go
func Filter[C, HKTA, A any](filter func(Predicate[A]) Endomorphism[HKTA]) func(Predicate[A]) Operator[C, HKTA, HKTA]
func FilterArray[C, A any](p Predicate[A]) Operator[C, []A, []A]
func FilterIter[C, A any](p Predicate[A]) Operator[C, Seq[A], Seq[A]]
func FilterMap[C, HKTA, HKTB, A, B any](filter func(option.Kleisli[A, B]) Reader[HKTA, HKTB]) func(option.Kleisli[A, B]) Operator[C, HKTA, HKTB]
func FilterMapArray[C, A, B any](p option.Kleisli[A, B]) Operator[C, []A, []B]
func FilterMapIter[C, A, B any](p option.Kleisli[A, B]) Operator[C, Seq[A], Seq[B]]
```

**Branching**:
```go
func Ternary[C, A, B any](pred func(A) bool, onTrue, onFalse Kleisli[C, A, B]) Kleisli[C, A, B]
```

**Eitherize -- converting Go functions to Effects**:
```go
func Eitherize[C, T any](f func(C, context.Context) (T, error)) Effect[C, T]
func Eitherize1[C, A, T any](f func(C, context.Context, A) (T, error)) Kleisli[C, A, T]
```

**Running Effects**:
```go
func Provide[A, C any](c C) func(Effect[C, A]) ReaderIOResult[A]
func Read[A, C any](c C) func(Effect[C, A]) Thunk[A]
func RunSync[A any](fa ReaderIOResult[A]) func(context.Context) (A, error)
```

**Do notation** (imperative-style effect composition):
```go
func Do[C, S any](empty S) Effect[C, S]
```

Bind operations for do-notation:
```go
// Bind: run effectful computation, set result into state
func Bind[C, S1, S2, T any](
    setter func(T) func(S1) S2,
    f func(S1) Effect[C, T],
) Operator[C, S1, S2]

// Let: set pure value into state
func Let[C, S1, S2, T any](
    setter func(T) func(S1) S2,
    f func(S1) T,
) Operator[C, S1, S2]

// Bind variants for lifting from other monads:
func BindIOK[C, S1, S2, T any](setter, f) Operator[C, S1, S2]
func BindResultK[C, S1, S2, T any](setter, f) Operator[C, S1, S2]
func BindReaderK[C, S1, S2, T any](setter, f) Operator[C, S1, S2]
func BindReaderIOK[C, S1, S2, T any](setter, f) Operator[C, S1, S2]
func BindEitherK[C, S1, S2, T any](setter, f) Operator[C, S1, S2]
```

Do-notation example:
```go
type State struct {
    User    User
    Posts   []Post
    Count   int
}

pipeline := F.Pipe3(
    EFF.Do[AppConfig](State{}),
    EFF.Bind(
        func(u User) func(State) State {
            return func(s State) State { s.User = u; return s }
        },
        func(_ State) EFF.Effect[AppConfig, User] {
            return fetchUserEff(123)
        },
    ),
    EFF.Bind(
        func(ps []Post) func(State) State {
            return func(s State) State { s.Posts = ps; return s }
        },
        func(s State) EFF.Effect[AppConfig, []Post] {
            return fetchPostsEff(s.User.ID)
        },
    ),
    EFF.Let(
        func(c int) func(State) State {
            return func(s State) State { s.Count = c; return s }
        },
        func(s State) int { return len(s.Posts) },
    ),
)
```

**Complete usage pattern**:
```go
import (
    "context"
    E "github.com/IBM/fp-go/v2/effect"
    F "github.com/IBM/fp-go/v2/function"
)

type AppConfig struct {
    APIKey string
}

func fetchData(cfg AppConfig, ctx context.Context) (string, error) {
    return "data from " + cfg.APIKey, nil
}

var fetchDataEff = E.Eitherize(fetchData)

pipeline := F.Pipe1(
    fetchDataEff,
    E.Map[AppConfig](func(s string) string { return "got: " + s }),
)

cfg := AppConfig{APIKey: "secret"}
value, err := E.RunSync(E.Provide[string](cfg)(pipeline))(context.Background())
// value == "got: data from secret", err == nil
```

---

## 3. Function Composition

**Package**: `github.com/IBM/fp-go/v2/function`

Import convention: `F "github.com/IBM/fp-go/v2/function"`

### 3.1 Pipe

`PipeN` takes an initial value and applies N functions left-to-right.

```go
func Pipe0[T0 any](t0 T0) T0
func Pipe1[F1 ~func(T0) T1, T0, T1 any](t0 T0, f1 F1) T1
func Pipe2[F1 ~func(T0) T1, F2 ~func(T1) T2, T0, T1, T2 any](t0 T0, f1 F1, f2 F2) T2
func Pipe3[...](t0 T0, f1 F1, f2 F2, f3 F3) T3
// ... up to Pipe10+
```

Usage -- value-first pipeline:
```go
result := F.Pipe3(
    someValue,          // initial value
    firstTransform,     // T0 -> T1
    secondTransform,    // T1 -> T2
    thirdTransform,     // T2 -> T3
)
```

### 3.2 Flow

`FlowN` creates a composed function (no initial value, returns a function).

```go
func Flow1[F1 ~func(T0) T1, T0, T1 any](f1 F1) func(T0) T1
func Flow2[F1 ~func(T0) T1, F2 ~func(T1) T2, T0, T1, T2 any](f1 F1, f2 F2) func(T0) T2
func Flow3[...](f1 F1, f2 F2, f3 F3) func(T0) T3
// ... up to Flow10+
```

Usage -- point-free composition:
```go
transform := F.Flow3(
    firstTransform,     // T0 -> T1
    secondTransform,    // T1 -> T2
    thirdTransform,     // T2 -> T3
)
result := transform(someValue)
```

### 3.3 Nullary

`NullaryN` creates a nullary function from a nullary producer and N transform functions:

```go
func Nullary1[F1 ~func() T1, T1 any](f1 F1) func() T1
func Nullary2[F1 ~func() T1, F2 ~func(T1) T2, T1, T2 any](f1 F1, f2 F2) func() T2
func Nullary3[...](f1, f2, f3) func() T3
```

### 3.4 Curry / Uncurry

```go
func Curry2[FCT ~func(T0, T1) T2, T0, T1, T2 any](f FCT) func(T0) func(T1) T2
func Curry3[FCT ~func(T0, T1, T2) T3, ...](f FCT) func(T0) func(T1) func(T2) T3
// up to Curry10+

func Uncurry2[FCT ~func(T0) func(T1) T2, T0, T1, T2 any](f FCT) func(T0, T1) T2
func Uncurry3[...](f FCT) func(T0, T1, T2) T3
// up to Uncurry10+
```

### 3.5 Utility Functions

```go
func Identity[A any](a A) A                            // id :: a -> a
func Constant[A any](a A) func() A                     // const :: a -> () -> a
func Constant1[B, A any](a A) func(B) A                // always returns a, ignores input
func Constant2[B, C, A any](a A) func(B, C) A          // ignores both inputs
func Flip[T1, T2, R any](f func(T1) func(T2) R) func(T2) func(T1) R
func Swap[T1, T2, R any](f func(T1, T2) R) func(T2, T1) R
func Bind1st[T1, T2, R any](f func(T1, T2) R, t1 T1) func(T2) R
func Bind2nd[T1, T2, R any](f func(T1, T2) R, t2 T2) func(T1) R
func First[T1, T2 any](t1 T1, _ T2) T1                // projection
func Second[T1, T2 any](_ T1, t2 T2) T2                // projection
func SK[T1, T2 any](_ T1, t2 T2) T2                    // SKI combinator
func Ternary[A, B any](pred func(A) bool, onTrue, onFalse func(A) B) func(A) B
func Zero[A any]() A                                    // zero value of type
func ToAny[A any](a A) any                              // upcast to any
```

**Type**:
```go
type Void = struct{}
var VOID Void = struct{}{}
```

### 3.6 When to Use Pipe vs Flow

| Pattern | Use When |
|---------|----------|
| `Pipe` | You have a concrete value and want to transform it through a pipeline |
| `Flow` | You want to create a reusable composed function (point-free style) |

```go
// Pipe: immediate evaluation with a value
result := F.Pipe2(42, double, toString)

// Flow: create a reusable function
transform := F.Flow2(double, toString)
result1 := transform(42)
result2 := transform(99)
```

---

## 4. Optics

**Package**: `github.com/IBM/fp-go/v2/optics/...`

### 4.1 Lens

**Package**: `github.com/IBM/fp-go/v2/optics/lens`

A Lens focuses on a field inside a product type (struct). It provides immutable get/set.

**Type definition**:
```go
type Lens[S, A any] struct {
    Get func(S) A
    Set func(S, A) S
}
```

**Construction**:
```go
func MakeLens[GET ~func(S) A, SET ~func(S, A) S, S, A any](get GET, set SET) Lens[S, A]
func MakeLensRef[GET ~func(*S) A, SET ~func(*S, A) *S, S, A any](get GET, set SET) Lens[*S, A]
```

**Operations**:
```go
func Get[S, A any](sa Lens[S, A]) func(S) A
func Set[S, A any](sa Lens[S, A]) func(A) func(S) S
func Modify[S any, FCT ~func(A) A, A any](f FCT) func(Lens[S, A]) func(S) S
func Compose[S, A, B any](ab Lens[A, B]) func(Lens[S, A]) Lens[S, B]
```

**Auto-generation**: Place in your Go file:
```go
//go:generate go run github.com/IBM/fp-go/v2 lens --dir . --filename gen_lens.go
```
This scans struct types and generates lenses for exported fields.

**Example**:
```go
type Address struct {
    Street string
    City   string
}

type Person struct {
    Name    string
    Address Address
}

nameLens := lens.MakeLens(
    func(p Person) string { return p.Name },
    func(p Person, n string) Person { p.Name = n; return p },
)

// Get
name := lens.Get(nameLens)(person)

// Set
updated := lens.Set(nameLens)("Alice")(person)

// Modify
uppered := lens.Modify(strings.ToUpper)(nameLens)(person)

// Compose lenses
addressLens := lens.MakeLens(
    func(p Person) Address { return p.Address },
    func(p Person, a Address) Person { p.Address = a; return p },
)
cityLens := lens.MakeLens(
    func(a Address) string { return a.City },
    func(a Address, c string) Address { a.City = c; return a },
)
personCityLens := F.Pipe1(addressLens, lens.Compose[Person](cityLens))
```

### 4.2 Prism

**Package**: `github.com/IBM/fp-go/v2/optics/prism`

A Prism focuses on a case of a sum type (variant). It may fail to match.

**Type definition**:
```go
type Prism[S, A any] struct {
    GetOption  func(S) Option[A]
    ReverseGet func(A) S
}
```

**Construction**:
```go
func MakePrism[S, A any](getOption func(S) Option[A], reverseGet func(A) S) Prism[S, A]
```

### 4.3 Iso

**Package**: `github.com/IBM/fp-go/v2/optics/iso`

An Iso is a lossless bidirectional transformation between types.

**Type definition**:
```go
type Iso[S, A any] struct {
    Get        func(S) A
    ReverseGet func(A) S
}
```

**Construction and operations**:
```go
func MakeIso[S, A any](get func(S) A, reverse func(A) S) Iso[S, A]
func Id[S any]() Iso[S, S]
func Compose[S, A, B any](ab Iso[A, B]) func(Iso[S, A]) Iso[S, B]
func Reverse[S, A any](sa Iso[S, A]) Iso[A, S]
func Modify[S any, FCT ~func(A) A, A any](f FCT) func(Iso[S, A]) Endomorphism[S]
func Unwrap[A, S any](s S) func(Iso[S, A]) A               // alias: To
func Wrap[S, A any](a A) func(Iso[S, A]) S                 // alias: From
func IMap[S, A, B any](ab func(A) B, ba func(B) A) func(Iso[S, A]) Iso[S, B]
```

**Example**:
```go
import "github.com/IBM/fp-go/v2/optics/iso"

celsiusToFahrenheit := iso.MakeIso(
    func(c float64) float64 { return c*9/5 + 32 },
    func(f float64) float64 { return (f - 32) * 5 / 9 },
)
f := celsiusToFahrenheit.Get(100.0)        // 212.0
c := celsiusToFahrenheit.ReverseGet(212.0) // 100.0
```

### 4.4 Optional (Lens + Option)

**Package**: `github.com/IBM/fp-go/v2/optics/lens/option`

An optional lens focuses on a field that may not exist.

```go
type OptionalLens[S, A any] struct {
    GetOption func(S) Option[A]
    Set       func(S, A) S
}
```

### 4.5 Traversal

**Package**: `github.com/IBM/fp-go/v2/optics/traversal`

Traversals focus on multiple targets within a structure.

### 4.6 Codec

**Package**: `github.com/IBM/fp-go/v2/optics/codec`

Codecs combine encoding and decoding with validation, built on top of optics.

Sub-packages:
- `optics/codec/decode` -- decoding monads
- `optics/codec/validate` -- validation with accumulated errors
- `optics/codec/validation` -- validation monad

---

## 5. Algebraic Structures

### 5.1 Eq

**Package**: `github.com/IBM/fp-go/v2/eq`

```go
type Eq[T any] interface {
    Equals(x, y T) bool
}

func FromStrictEquals[T comparable]() Eq[T]
func FromEquals[T any](c func(x, y T) bool) Eq[T]
func Empty[T any]() Eq[T]                              // always true
func Equals[T any](eq Eq[T]) func(T) func(T) bool      // curried
```

Laws: reflexive, symmetric, transitive.

**Usage with other types**: `option.Eq`, `either.Eq`, `array.Eq` create Eq instances for container types given an Eq for their elements.

### 5.2 Ord

**Package**: `github.com/IBM/fp-go/v2/ord`

```go
type Ord[T any] interface {
    Eq[T]
    Compare(x, y T) int  // -1, 0, 1
}

func FromStrictCompare[A constraints.Ordered]() Ord[A]
func FromCompare[T any](compare func(T, T) int) Ord[T]
func MakeOrd[T any](c func(x, y T) int, e func(x, y T) bool) Ord[T]
func Reverse[T any](o Ord[T]) Ord[T]
func Contramap[A, B any](f func(B) A) func(Ord[A]) Ord[B]

// Predicates
func Lt[A any](o Ord[A]) func(A) func(A) bool
func Leq[A any](o Ord[A]) func(A) func(A) bool
func Gt[A any](o Ord[A]) func(A) func(A) bool
func Geq[A any](o Ord[A]) func(A) func(A) bool
func Between[A any](o Ord[A]) func(A, A) func(A) bool

// Selection
func Min[A any](o Ord[A]) func(A, A) A
func Max[A any](o Ord[A]) func(A, A) A
func Clamp[A any](o Ord[A]) func(A, A) func(A) A

// Built-in
func OrdTime() Ord[time.Time]
```

**Example**:
```go
type Person struct { Name string; Age int }
personByAge := ord.Contramap(func(p Person) int { return p.Age })(ord.FromStrictCompare[int]())
```

### 5.3 Semigroup

**Package**: `github.com/IBM/fp-go/v2/semigroup`

```go
type Semigroup[A any] interface {
    Concat(x, y A) A
}

func MakeSemigroup[A any](c func(A, A) A) Semigroup[A]
func Reverse[A any](m Semigroup[A]) Semigroup[A]
func First[A any]() Semigroup[A]                            // always first argument
func Last[A any]() Semigroup[A]                             // always second argument
func FunctionSemigroup[A, B any](s Semigroup[B]) Semigroup[func(A) B]
func ConcatWith[A any](s Semigroup[A]) func(A) func(A) A   // curried, left first
func AppendTo[A any](s Semigroup[A]) func(A) func(A) A     // curried, right first
```

Law: associativity -- `Concat(Concat(x, y), z) == Concat(x, Concat(y, z))`

### 5.4 Monoid

**Package**: `github.com/IBM/fp-go/v2/monoid`

```go
type Monoid[A any] interface {
    Semigroup[A]
    Empty() A
}

func MakeMonoid[A any](c func(A, A) A, e A) Monoid[A]
func Reverse[A any](m Monoid[A]) Monoid[A]
func ToSemigroup[A any](m Monoid[A]) Semigroup[A]
```

Laws: associativity + left/right identity (`Concat(Empty(), x) == x`).

**Common monoids** (in `number`, `string`, `boolean` packages):
- `number.MonoidSum[int]()` -- addition with 0
- `number.MonoidProduct[int]()` -- multiplication with 1
- `string.Monoid` -- concatenation with ""
- `boolean.MonoidAll` -- AND with true
- `boolean.MonoidAny` -- OR with false

### 5.5 Composing with Monads

Each monad package provides `ApplicativeMonoid` and/or `AlternativeMonoid`:

```go
// Combine Options using a monoid on their values
option.ApplicativeMonoid[A](m monoid.Monoid[A]) monoid.Monoid[Option[A]]

// Combine Eithers using a monoid on their Right values
either.ApplicativeMonoid[E, A](m monoid.Monoid[A]) monoid.Monoid[Either[E, A]]

// First-success alternative semantics
effect.AlternativeMonoid[C, A](m monoid.Monoid[A]) monoid.Monoid[Effect[C, A]]
```

---

## 6. Array Operations

**Package**: `github.com/IBM/fp-go/v2/array`

Import convention: `A "github.com/IBM/fp-go/v2/array"`

### 6.1 Construction

```go
func From[A any](data ...A) []A
func Of[A any](a A) []A                              // single-element array
func Empty[A any]() []A                               // empty array
```

### 6.2 Transformation

```go
func Map[A, B any](f func(A) B) Operator[A, B]                    // Operator[A, B] = func([]A) []B
func MapWithIndex[A, B any](f func(int, A) B) Operator[A, B]
func MapRef[A, B any](f func(*A) B) Operator[A, B]                // avoids copying
```

### 6.3 Filtering

```go
func Filter[A any](pred func(A) bool) Operator[A, A]
func FilterWithIndex[A any](pred func(int, A) bool) Operator[A, A]
func FilterRef[A any](pred func(*A) bool) Operator[A, A]
func FilterMap[A, B any](f option.Kleisli[A, B]) Operator[A, B]   // map + filter in one
func FilterMapWithIndex[A, B any](f func(int, A) Option[B]) Operator[A, B]
```

### 6.4 Folding

```go
func Reduce[A, B any](f func(B, A) B, initial B) func([]A) B
func ReduceWithIndex[A, B any](f func(int, B, A) B, initial B) func([]A) B
func ReduceRight[A, B any](f func(A, B) B, initial B) func([]A) B
func ReduceRightWithIndex[A, B any](f func(int, A, B) B, initial B) func([]A) B
```

### 6.5 Element Access

```go
func Head[A any](as []A) Option[A]
func Last[A any](as []A) Option[A]
func Tail[A any](as []A) Option[[]A]
```

### 6.6 Building

```go
func Append[A any](as []A, a A) []A
func Prepend[A any](head A) Operator[A, A]             // func([]A) []A
func Concat[A any](suffix []A) Operator[A, A]           // func([]A) []A
func PrependAll[A any](middle A) Operator[A, A]         // intersperse
```

### 6.7 Sorting

```go
func Sort[T any](ord O.Ord[T]) Operator[T, T]
func SortByKey[K, T any](ord O.Ord[K], f func(T) K) Operator[T, T]
func SortBy[T any](ord []O.Ord[T]) Operator[T, T]
```

### 6.8 Zipping

```go
func Zip[A, B any](fb []B) func([]A) []pair.Pair[A, B]
func ZipWith[FCT ~func(A, B) C, A, B, C any](fa []A, fb []B, f FCT) []C
```

### 6.9 Traversal / Sequence

Each monad package (option, result, ioresult, effect, etc.) provides `TraverseArray` and `SequenceArray` in its own package:

```go
// In option package:
option.TraverseArray[A, B any](f func(A) Option[B]) func([]A) Option[[]B]
option.SequenceArray[A any](arr []Option[A]) Option[[]A]
option.CompactArray[A any](arr []Option[A]) []A

// In result package:
result.TraverseArray[A, B any](f func(A) Result[B]) func([]A) Result[[]B]
result.SequenceArray[A any](arr []Result[A]) Result[[]A]

// In effect package:
effect.TraverseArray[C, A, B any](f func(A) Effect[C, B]) func([]A) Effect[C, []B]
```

### 6.10 Monoid for arrays

```go
// In array/magma package:
func ConcatAll[A any](m monoid.Monoid[A]) func([]A) A
```

### 6.11 NonEmptyArray

**Package**: `github.com/IBM/fp-go/v2/array/nonempty`

A non-empty array guarantees at least one element.

```go
type NonEmptyArray[A any] = []A   // invariant: len >= 1

func Of[A any](first A) NonEmptyArray[A]
func From[A any](first A, data ...A) NonEmptyArray[A]
func Head[A any](as NonEmptyArray[A]) A               // always succeeds
func Tail[A any](as NonEmptyArray[A]) []A
func Map[A, B any](f func(A) B) Operator[A, B]
func Reduce[A, B any](f func(B, A) B, initial B) func(NonEmptyArray[A]) B
```

---

## 7. Idiomatic Packages

**Packages**: `github.com/IBM/fp-go/v2/idiomatic/option`, `github.com/IBM/fp-go/v2/idiomatic/result`

These provide alternative representations using Go's native tuple returns instead of wrapper structs.

### 7.1 Idiomatic Option

**Package**: `github.com/IBM/fp-go/v2/idiomatic/option`

Instead of `Option[A]` struct, uses `(A, bool)` tuples:

```go
type Kleisli[A, B any] = func(A) (B, bool)       // was: func(A) Option[B]
type Operator[A, B any] = func(A, bool) (B, bool) // was: func(Option[A]) Option[B]
```

**Construction**:
```go
func Some[A any](a A) (A, bool)                   // returns (a, true)
func None[A any]() (A, bool)                       // returns (zero, false)
func Of[A any](a A) (A, bool)                     // alias for Some
func FromPredicate[A any](pred func(A) bool) Kleisli[A, A]
func FromNillable[A any](a *A) (*A, bool)
```

**Transformation**:
```go
func Map[A, B any](f func(A) B) Operator[A, B]
func MapTo[A, B any](b B) Operator[A, B]
func Ap[B, A any](fa A, faok bool) Operator[func(A) B, B]
func Fold[A, B any](onNone func() B, onSome func(A) B) func(A, bool) B
func GetOrElse[A any](onNone func() A) func(A, bool) A
```

### 7.2 Idiomatic Result

**Package**: `github.com/IBM/fp-go/v2/idiomatic/result`

Instead of `Either[error, A]` struct, uses `(A, error)` tuples:

```go
type Kleisli[A, B any] = func(A) (B, error)        // standard Go pattern
type Operator[A, B any] = func(A, error) (B, error) // transforms (value, error) pairs
```

**Construction**:
```go
func Left[A any](err error) (A, error)              // returns (zero, err)
func Right[A any](a A) (A, error)                    // returns (a, nil)
func IsLeft[A any](_ A, err error) bool              // err != nil
func IsRight[A any](_ A, err error) bool             // err == nil
```

### 7.3 Differences Between Standard and Idiomatic

| Aspect | Standard (`option`, `result`) | Idiomatic (`idiomatic/option`, `idiomatic/result`) |
|--------|-------------------------------|--------------------------------------------------|
| Type | Wrapper struct `Option[A]`, `Either[E, A]` | Go tuples `(A, bool)`, `(A, error)` |
| JSON | Implements `MarshalJSON`/`UnmarshalJSON` | No automatic serialization |
| Performance | One allocation for struct | Zero allocation (stack values) |
| Composability | Full monad stack integration | Limited to same-style composition |
| Interop | Needs `Unwrap` to get Go values | Already in Go-native form |
| Use when | Building functional pipelines, monad stacks | Performance-critical paths, Go-native APIs |

### 7.4 Converting Between Representations

```go
// Standard Option -> Idiomatic: use Unwrap
val, ok := option.Unwrap(opt)

// Idiomatic -> Standard Option
opt := option.FromValidation(func(a A) (A, bool) { return val, ok })

// Standard Result -> Go tuple: use Unwrap/UnwrapError
val, err := result.Unwrap(res)

// Go tuple -> Standard Result: use TryCatchError
res := result.TryCatchError(val, err)
```

---

## 8. Naming Conventions

### 8.1 Function Name Patterns

| Suffix | Meaning |
|--------|---------|
| `Monad*` | Takes the container as the first argument (uncurried) |
| (none) | Curried form, suitable for `Pipe` / `Flow` |
| `*K` | Converts from a different monad kind (Kleisli lift) |
| `*L` | Lens-based variant for do-notation |
| `*W` | "Widening" -- allows different type parameters |
| `*First` | Performs side effect but returns original value |
| `Tap*` | Alias for `*First` variants |

### 8.2 Import Conventions

```go
import (
    F "github.com/IBM/fp-go/v2/function"
    O "github.com/IBM/fp-go/v2/option"
    E "github.com/IBM/fp-go/v2/either"
    R "github.com/IBM/fp-go/v2/result"
    IO "github.com/IBM/fp-go/v2/io"
    IOR "github.com/IBM/fp-go/v2/ioresult"
    RD "github.com/IBM/fp-go/v2/reader"
    RIR "github.com/IBM/fp-go/v2/readerioresult"
    EFF "github.com/IBM/fp-go/v2/effect"
    A "github.com/IBM/fp-go/v2/array"
    N "github.com/IBM/fp-go/v2/number"
    S "github.com/IBM/fp-go/v2/string"
)
```

### 8.3 The `MonadX` vs `X` Pattern

Every monad operation exists in two forms:

```go
// Uncurried -- takes container + function
result := option.MonadMap(someOpt, transform)

// Curried -- returns an operator for Pipe/Flow
result := F.Pipe1(someOpt, option.Map(transform))
```

Always prefer the curried form with `Pipe`/`Flow` for composition.

---

## 9. Common Patterns

### 9.1 Error Recovery with OrElse / ChainLeft

```go
// Try primary, fall back to secondary
result := F.Pipe2(
    primaryAction(),
    R.OrElse(func(err error) R.Result[string] {
        if errors.Is(err, ErrNotFound) {
            return secondaryAction()
        }
        return R.Left[string](err)
    }),
)
```

### 9.2 Validation with Either

```go
// Collect all validation errors (not short-circuit)
// Use either/validation for applicative validation
import V "github.com/IBM/fp-go/v2/either/validation"

validated := V.SequenceT2(
    semigroup.MakeSemigroup(func(a, b []string) []string { return append(a, b...) }),
)(
    validateName(input),
    validateAge(input),
)
```

### 9.3 Dependency Injection with Effect

```go
type Services struct {
    DB     *sql.DB
    Logger *slog.Logger
}

func getUser(id int) EFF.Effect[Services, User] {
    return EFF.Chain[Services](func(svc Services) EFF.Effect[Services, User] {
        return EFF.Eitherize(func(svc Services, ctx context.Context) (User, error) {
            return svc.DB.QueryRowContext(ctx, "SELECT ...").Scan(...)
        })
    })(EFF.Asks[Services](F.Identity[Services]))
}

// Or more concisely with Eitherize1:
var getUserEff = EFF.Eitherize1(func(svc Services, ctx context.Context, id int) (User, error) {
    return svc.DB.QueryRowContext(ctx, "SELECT ...").Scan(...)
})
```

### 9.4 Array Pipeline

```go
result := F.Pipe4(
    users,
    A.Filter(func(u User) bool { return u.Active }),
    A.Map(func(u User) string { return u.Email }),
    A.Sort(ord.FromStrictCompare[string]()),
    A.Head,
)
// result: Option[string]
```

### 9.5 Traverse -- Fail-Fast on Collections

```go
// Parse all strings to ints, fail on first error
parsed := F.Pipe1(
    []string{"1", "2", "3"},
    R.TraverseArray(func(s string) R.Result[int] {
        return R.TryCatchError(strconv.Atoi(s))
    }),
)
// parsed: Result[[]int] = Right([1, 2, 3])
```

### 9.6 Memoization

```go
// IO computations can be memoized (computed once)
expensiveIO := IO.Memoize(func() int {
    // expensive computation
    return computeExpensiveValue()
})
// First call computes, subsequent calls return cached result
```

---

## 10. Code Generation

fp-go uses `//go:generate` annotations for auto-generating boilerplate.

### 10.1 Lens Generation

```go
//go:generate go run github.com/IBM/fp-go/v2 lens --dir . --filename gen_lens.go
```

Generates lenses for all exported struct fields in the package.
Add `--include-test-files` to also scan test files.

### 10.2 Monad Operation Generation

```go
//go:generate go run github.com/IBM/fp-go/v2 option --count 10 --filename gen.go
```

Generates `Pipe`, `Flow`, `Curry`, `Uncurry`, `Sequence`, and other arity-dependent
functions up to the specified count.

---

## 11. Quick Reference: Choosing the Right Type

```
Do I need error handling?
  No  -> Option[A] (presence/absence)
  Yes -> Do I need a typed error?
           No  -> Result[A]  (error channel)
           Yes -> Either[E, A]

Do I need side effects?
  No  -> Use pure types above
  Yes -> Do I need dependency injection?
           No  -> IOResult[A]
           Yes -> Do I need context.Context?
                    No  -> ReaderIOResult[R, A]
                    Yes -> Effect[C, A]  (recommended)
```

**Rule of thumb**: Start with `Effect` for application code. Use `Option`/`Result` for pure data transformations. Use `IO` only when you need lazy computation without errors.

---

## 12. Intermediate Monad Types

These are available but less commonly used directly. They exist as layers in the monad stack.

### 12.1 IOEither

**Package**: `github.com/IBM/fp-go/v2/ioeither`

```go
type IOEither[E, A any] = IO[Either[E, A]]   // func() Either[E, A]
```

IOEither with typed error E. IOResult is the `E = error` specialization.

### 12.2 IOOption

**Package**: `github.com/IBM/fp-go/v2/iooption`

```go
type IOOption[A any] = IO[Option[A]]   // func() Option[A]
```

IO with optional result. Useful when absence is expected (not an error).

### 12.3 ReaderResult

**Package**: `github.com/IBM/fp-go/v2/readerresult`

```go
type ReaderResult[R, A any] = Reader[R, Result[A]]   // func(R) Either[error, A]
```

Reader with error handling but no IO.

### 12.4 ReaderIO

**Package**: `github.com/IBM/fp-go/v2/readerio`

```go
type ReaderIO[R, A any] = Reader[R, IO[A]]   // func(R) func() A
```

Reader with IO but no error handling.

### 12.5 ReaderOption

**Package**: `github.com/IBM/fp-go/v2/readeroption`

```go
type ReaderOption[R, A any] = Reader[R, Option[A]]   // func(R) Option[A]
```

Reader with optional result.

### 12.6 ReaderIOEither

**Package**: `github.com/IBM/fp-go/v2/readerioeither`

```go
type ReaderIOEither[R, E, A any] = Reader[R, IOEither[E, A]]  // func(R) func() Either[E, A]
```

Full reader + IO + typed error. ReaderIOResult is the `E = error` specialization.

### 12.7 State / StateIO

**Package**: `github.com/IBM/fp-go/v2/state`, `github.com/IBM/fp-go/v2/stateio`

```go
type State[S, A any] = func(S) Pair[A, S]         // pure stateful computation
type StateIO[S, A any] = func(S) IO[Pair[A, S]]   // stateful IO computation
```

State monads thread mutable state through a computation.

---

## 13. Pair (Tuple)

**Package**: `github.com/IBM/fp-go/v2/pair`

```go
type Pair[L, R any] struct { /* private fields */ }

func MakePair[L, R any](l L, r R) Pair[L, R]
func Head[L, R any](p Pair[L, R]) L
func Tail[L, R any](p Pair[L, R]) R
```

Pair is a functor over the Tail (right) element:
```go
func Map[L, A, B any](f func(A) B) func(Pair[L, A]) Pair[L, B]
func BiMap[L1, L2, R1, R2 any](f func(L1) L2, g func(R1) R2) func(Pair[L1, R1]) Pair[L2, R2]
```

---

## 14. Predicate

**Package**: `github.com/IBM/fp-go/v2/predicate`

```go
type Predicate[A any] = func(A) bool

func Not[A any](pred Predicate[A]) Predicate[A]
func And[A any](second Predicate[A]) func(Predicate[A]) Predicate[A]
func Or[A any](second Predicate[A]) func(Predicate[A]) Predicate[A]
func IsZero[A comparable]() Predicate[A]
func IsNonZero[A comparable]() Predicate[A]
func IsEqual[A any](eq Eq[A]) func(A) Predicate[A]
```

---

## 15. Endomorphism

**Package**: `github.com/IBM/fp-go/v2/endomorphism`

```go
type Endomorphism[A any] = func(A) A

func Identity[A any]() Endomorphism[A]
```

Endomorphisms form a monoid under composition:
```go
func Monoid[A any]() monoid.Monoid[Endomorphism[A]]  // concat = compose, empty = identity
```

---

## 16. Retry

**Package**: `github.com/IBM/fp-go/v2/retry`

Used with Effect and IOResult for retrying failed operations.

```go
// In effect package:
func Retrying[C, A any](
    policy retry.RetryPolicy,
    action func(retry.RetryStatus) Effect[C, A],
    check func(Result[A]) bool,
) Effect[C, A]
```

---

## 17. Type Expansion Table

For quick lookup of what each monad is made of:

| Type | Expanded | Has DI | Has IO | Has Error | Error Type |
|------|----------|--------|--------|-----------|------------|
| `Option[A]` | struct{value A, isSome bool} | No | No | No | n/a |
| `Either[E,A]` | struct{r A, l E, isLeft bool} | No | No | Yes | E |
| `Result[A]` | `Either[error, A]` | No | No | Yes | error |
| `IO[A]` | `func() A` | No | Yes | No | n/a |
| `IOResult[A]` | `func() Either[error, A]` | No | Yes | Yes | error |
| `IOEither[E,A]` | `func() Either[E, A]` | No | Yes | Yes | E |
| `Reader[R,A]` | `func(R) A` | Yes(R) | No | No | n/a |
| `ReaderIOResult[R,A]` | `func(R) func() Either[error, A]` | Yes(R) | Yes | Yes | error |
| `Effect[C,A]` | `func(C) func(ctx) func() Either[error, A]` | Yes(C+ctx) | Yes | Yes | error |
