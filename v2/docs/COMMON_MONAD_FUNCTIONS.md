# Common Monad Functions in fp-go/v2

This document catalogs functions that exist across multiple monads in the fp-go/v2 library, explaining their conceptual purpose and listing available implementations.

## Table of Contents

- [Core Type Class Operations](#core-type-class-operations)
  - [Of (Pointed)](#of-pointed)
  - [Map (Functor)](#map-functor)
  - [Ap (Applicative)](#ap-applicative)
  - [Chain (Monad)](#chain-monad)
- [Transformation Operations](#transformation-operations)
  - [Fold](#fold)
  - [Filter](#filter)
- [Composition Operations](#composition-operations)
  - [Sequence](#sequence)
  - [Traverse](#traverse)
- [Do-Notation Operations](#do-notation-operations)
  - [Do](#do)
  - [Bind](#bind)
  - [Let](#let)
  - [ApS](#aps)
- [Side Effect Operations](#side-effect-operations)
  - [ChainFirst](#chainfirst)
  - [Tap](#tap)
- [Structure Operations](#structure-operations)
  - [Flatten](#flatten)
- [Alternative Operations](#alternative-operations)
  - [Alt](#alt)
  - [OrElse](#orelse)
  - [GetOrElse](#getorelse)
- [Bifunctor Operations](#bifunctor-operations)
  - [BiMap](#bimap)
  - [MapLeft](#mapleft)
- [Monoid Operations](#monoid-operations)
  - [Overview](#overview)
  - [Core Monoid Functions](#core-monoid-functions)
  - [Specialized Monoids](#specialized-monoids)
  - [Monad-Specific Monoids](#monad-specific-monoids)
  - [Monoid Patterns](#monoid-patterns)
  - [Usage Guidelines](#usage-guidelines)

---

## Core Type Class Operations

### Of (Pointed)

**Concept**: Lifts a pure value into a monadic context. This is the fundamental operation for wrapping values in a computational context.

**Type Signature**: `func Of[A any](a A) M[A]`

**Purpose**: Creates the simplest possible monadic value containing the given value, with no effects or errors.

**Implementations**:

- [`option.Of`](../option/core.go:124) - Wraps value in Some
- [`either.Of`](../either/either.go:39) - Creates Right value
- [`result.Of`](../result/either.go:38) - Creates successful Result
- [`io.Of`](../io/io.go:44) - Creates pure IO computation
- [`ioeither.Of`](../ioeither/ioeither.go:60) - Creates successful IO effect
- [`ioresult.Of`](../ioresult/ioeither.go:54) - Creates successful IO Result
- [`iooption.Of`](../iooption/iooption.go:31) - Creates Some IO computation
- [`array.Of`](../array/array.go:581) - Creates single-element array
- [`lazy.Of`](../lazy/lazy.go:33) - Creates lazy computation
- [`reader.Of`](../reader/reader.go:226) - Creates constant Reader
- [`readereither.Of`](../readereither/reader.go:85) - Creates successful ReaderEither
- [`readeroption.Of`](../readeroption/reader.go:127) - Creates Some ReaderOption
- [`readerresult.Of`](../readerresult/reader.go:279) - Creates successful ReaderResult
- [`readerio.Of`](../readerio/reader.go:452) - Creates pure ReaderIO
- [`readerioeither.Of`](../readerioeither/reader.go:1126) - Creates successful ReaderIOEither
- [`readerioresult.Of`](../readerioresult/reader.go:616) - Creates successful ReaderIOResult
- [`readeriooption.Of`](../readeriooption/reader.go:138) - Creates Some ReaderIOOption
- [`readerreaderioeither.Of`](../readerreaderioeither/reader.go:638) - Creates successful nested Reader
- [`state.Of`](../state/state.go:102) - Creates stateful computation
- [`stateio.Of`](../stateio/state.go:34) - Creates stateful IO computation
- [`statereaderioeither.Of`](../statereaderioeither/state.go:56) - Creates successful stateful Reader
- [`effect.Of`](../effect/effect.go:139) - Creates pure Effect
- [`identity.Of`](../identity/identity.go:134) - Identity (no-op)
- [`pair.Of`](../pair/pair.go:28) - Creates Pair with same value twice
- [`tuple.Of`](../tuple/tuple.go:34) - Creates single-element Tuple
- [`iterator/stateless.Of`](../iterator/stateless/iterator.go:41) - Creates single-element iterator
- [`iterator/iter.Of`](../iterator/iter/iter.go:94) - Creates single-element sequence
- [`iterator/itereither.Of`](../iterator/itereither/ioeither.go:221) - Creates successful sequence
- [`iterator/iterresult.Of`](../iterator/iterresult/ioeither.go:206) - Creates successful result sequence
- [`optics/codec/decode.Of`](../optics/codec/decode/monad.go:17) - Creates successful decoder
- [`optics/codec/validation.Of`](../optics/codec/validation/monad.go:18) - Creates successful validation
- [`optics/codec/validate.Of`](../optics/codec/validate/validate.go:159) - Creates successful validator

**Example**:
```go
// Option
opt := option.Of(42) // Some(42)

// Either
e := either.Of[error](42) // Right(42)

// IO
computation := io.Of(42) // IO that returns 42
result := computation() // 42
```

---

### Map (Functor)

**Concept**: Transforms the value inside a context by applying a function, preserving the structure. This is the fundamental operation for value transformation.

**Type Signature**: `func Map[A, B any](f func(A) B) func(M[A]) M[B]`

**Purpose**: Applies a pure function to the wrapped value without changing the computational context (errors, effects, etc. are preserved).

**Implementations**:

- [`option.Map`](../option/option.go:158) - Transforms Some values, preserves None
- [`either.Map`](../either/either.go:160) - Transforms Right values, preserves Left
- [`result.Map`](../result/either.go:156) - Transforms success values, preserves errors
- [`io.Map`](../io/io.go:97) - Transforms IO computation result
- [`ioeither.Map`](../ioeither/ioeither.go:140) - Transforms successful IO effects
- [`ioresult.Map`](../ioresult/ioeither.go:142) - Transforms successful IO results
- [`iooption.Map`](../iooption/iooption.go:84) - Transforms Some IO values
- [`array.Map`](../array/array.go:158) - Transforms each array element
- [`lazy.Map`](../lazy/lazy.go:93) - Transforms lazy computation result
- [`reader.Map`](../reader/reader.go:137) - Transforms Reader result
- [`readereither.Map`](../readereither/reader.go:65) - Transforms successful ReaderEither
- [`readeroption.Map`](../readeroption/reader.go:85) - Transforms Some ReaderOption
- [`readerresult.Map`](../readerresult/reader.go:190) - Transforms successful ReaderResult
- [`readerio.Map`](../readerio/reader.go:229) - Transforms ReaderIO result
- [`readerioeither.Map`](../readerioeither/reader.go:118) - Transforms successful ReaderIOEither
- [`readerioresult.Map`](../readerioresult/reader.go:75) - Transforms successful ReaderIOResult
- [`readeriooption.Map`](../readeriooption/reader.go:87) - Transforms Some ReaderIOOption
- [`readerreaderioeither.Map`](../readerreaderioeither/reader.go:89) - Transforms nested Reader result
- [`state.Map`](../state/state.go:134) - Transforms stateful computation result
- [`stateio.Map`](../stateio/state.go:63) - Transforms stateful IO result
- [`statereaderioeither.Map`](../statereaderioeither/state.go:86) - Transforms stateful Reader result
- [`effect.Map`](../effect/effect.go:167) - Transforms Effect result
- [`identity.Map`](../identity/identity.go:89) - Identity transformation
- [`pair.Map`](../pair/pair.go:210) - Transforms second element
- [`record.Map`](../record/record.go:275) - Transforms each record value
- [`iterator/stateless.Map`](../iterator/stateless/iterator.go:66) - Transforms iterator elements
- [`iterator/iter.Map`](../iterator/iter/iter.go:220) - Transforms sequence elements
- [`iterator/itereither.Map`](../iterator/itereither/ioeither.go:321) - Transforms successful sequence elements
- [`iterator/iterresult.Map`](../iterator/iterresult/ioeither.go:281) - Transforms successful result sequence
- [`optics/codec/decode.Map`](../optics/codec/decode/monad.go:404) - Transforms decoded value
- [`optics/codec/validation.Map`](../optics/codec/validation/monad.go:146) - Transforms validated value
- [`optics/codec/validate.Map`](../optics/codec/validate/validate.go:348) - Transforms validator result

**Example**:
```go
// Option
double := option.Map(func(x int) int { return x * 2 })
result := double(option.Some(21)) // Some(42)

// Either
toString := either.Map[error](strconv.Itoa)
result := toString(either.Right[error](42)) // Right("42")

// Array
squares := array.Map(func(x int) int { return x * x })
result := squares([]int{1, 2, 3}) // [1, 4, 9]
```

---

### Ap (Applicative)

**Concept**: Applies a wrapped function to a wrapped value. This enables combining multiple independent computations.

**Type Signature**: `func Ap[B, A any](fa M[A]) func(M[func(A) B]) M[B]`

**Purpose**: Allows applying functions that are themselves wrapped in a context to values in a context. Essential for combining multiple independent effects or validations.

**Implementations**:

- [`option.Ap`](../option/option.go:123) - Applies Some function to Some value
- [`either.Ap`](../either/either.go:79) - Applies Right function to Right value (fail-fast)
- [`result.Ap`](../result/either.go:74) - Applies success function to success value
- [`io.Ap`](../io/io.go:178) - Applies IO function to IO value
- [`ioeither.Ap`](../ioeither/ioeither.go:191) - Applies successful IO function
- [`ioresult.Ap`](../ioresult/ioeither.go:205) - Applies successful IO result function
- [`iooption.Ap`](../iooption/iooption.go:103) - Applies Some IO function
- [`array.Ap`](../array/array.go:618) - Cartesian product of functions and values
- [`lazy.Ap`](../lazy/lazy.go:162) - Applies lazy function to lazy value
- [`reader.Ap`](../reader/reader.go:213) - Applies Reader function to Reader value
- [`readereither.Ap`](../readereither/reader.go:93) - Applies successful ReaderEither function
- [`readeroption.Ap`](../readeroption/reader.go:157) - Applies Some ReaderOption function
- [`readerresult.Ap`](../readerresult/reader.go:351) - Applies successful ReaderResult function
- [`readerio.Ap`](../readerio/reader.go:535) - Applies ReaderIO function
- [`readerioeither.Ap`](../readerioeither/reader.go:1054) - Applies successful ReaderIOEither function
- [`readerioresult.Ap`](../readerioresult/reader.go:561) - Applies successful ReaderIOResult function
- [`readeriooption.Ap`](../readeriooption/reader.go:173) - Applies Some ReaderIOOption function
- [`readerreaderioeither.Ap`](../readerreaderioeither/reader.go:581) - Applies nested Reader function
- [`state.Ap`](../state/state.go:216) - Applies stateful function
- [`stateio.Ap`](../stateio/state.go:126) - Applies stateful IO function
- [`statereaderioeither.Ap`](../statereaderioeither/state.go:150) - Applies stateful Reader function
- [`effect.Ap`](../effect/effect.go:435) - Applies Effect function
- [`identity.Ap`](../identity/identity.go:55) - Direct function application
- [`pair.Ap`](../pair/pair.go:479) - Applies function to second element
- [`record.Ap`](../record/record.go:1091) - Applies record of functions to record of values
- [`iterator/stateless.Ap`](../iterator/stateless/iterator.go:109) - Applies iterator of functions
- [`iterator/iter.Ap`](../iterator/iter/iter.go:679) - Applies sequence of functions
- [`iterator/itereither.Ap`](../iterator/itereither/ioeither.go:402) - Applies successful sequence functions
- [`iterator/iterresult.Ap`](../iterator/iterresult/ioeither.go:350) - Applies successful result sequence functions
- [`optics/codec/decode.Ap`](../optics/codec/decode/monad.go:444) - Applies decoder function
- [`optics/codec/validation.Ap`](../optics/codec/validation/monad.go:34) - Applies validation function (accumulates errors)
- [`optics/codec/validate.Ap`](../optics/codec/validate/validate.go:865) - Applies validator function

**Example**:
```go
// Option - combining independent values
add := func(x int) func(int) int {
    return func(y int) int { return x + y }
}
result := F.Pipe2(
    option.Some(add),
    option.Ap(option.Some(10)),
    option.Ap(option.Some(32)),
) // Some(42)

// Either - fail-fast on first error
result := F.Pipe1(
    either.Right[error](func(x int) int { return x * 2 }),
    either.Ap[int](either.Right[error](21)),
) // Right(42)
```

---

### Chain (Monad)

**Concept**: Sequences computations where each step depends on the previous result. Also known as flatMap or bind.

**Type Signature**: `func Chain[A, B any](f func(A) M[B]) func(M[A]) M[B]`

**Purpose**: Enables sequential composition of operations that return monadic values. The key operation that makes a Functor into a Monad.

**Implementations**:

- [`option.Chain`](../option/option.go:285) - Sequences Option-returning operations
- [`either.Chain`](../either/either.go:324) - Sequences Either-returning operations (fail-fast)
- [`result.Chain`](../result/either.go:248) - Sequences Result-returning operations
- [`io.Chain`](../io/io.go:134) - Sequences IO computations
- [`ioeither.Chain`](../ioeither/ioeither.go:160) - Sequences IO Either operations
- [`ioresult.Chain`](../ioresult/ioeither.go:163) - Sequences IO Result operations
- [`iooption.Chain`](../iooption/iooption.go:92) - Sequences IO Option operations
- [`array.Chain`](../array/array.go:602) - Flat-maps over arrays
- [`lazy.Chain`](../lazy/lazy.go:132) - Sequences lazy computations
- [`reader.Chain`](../reader/reader.go:334) - Sequences Reader computations
- [`readereither.Chain`](../readereither/reader.go:73) - Sequences ReaderEither operations
- [`readeroption.Chain`](../readeroption/reader.go:114) - Sequences ReaderOption operations
- [`readerresult.Chain`](../readerresult/reader.go:222) - Sequences ReaderResult operations
- [`readerio.Chain`](../readerio/reader.go:375) - Sequences ReaderIO operations
- [`readerioeither.Chain`](../readerioeither/reader.go:1062) - Sequences ReaderIOEither operations
- [`readerioresult.Chain`](../readerioresult/reader.go:570) - Sequences ReaderIOResult operations
- [`readeriooption.Chain`](../readeriooption/reader.go:121) - Sequences ReaderIOOption operations
- [`readerreaderioeither.Chain`](../readerreaderioeither/reader.go:595) - Sequences nested Reader operations
- [`state.Chain`](../state/state.go:178) - Sequences stateful computations
- [`stateio.Chain`](../stateio/state.go:99) - Sequences stateful IO operations
- [`statereaderioeither.Chain`](../statereaderioeither/state.go:122) - Sequences stateful Reader operations
- [`effect.Chain`](../effect/effect.go:198) - Sequences Effect operations
- [`identity.Chain`](../identity/identity.go:171) - Direct function application
- [`pair.Chain`](../pair/pair.go:374) - Chains operations on second element
- [`record.Chain`](../record/record.go:696) - Flat-maps over record values
- [`iterator/stateless.Chain`](../iterator/stateless/iterator.go:74) - Flat-maps over iterators
- [`iterator/iter.Chain`](../iterator/iter/iter.go:590) - Flat-maps over sequences
- [`iterator/itereither.Chain`](../iterator/itereither/ioeither.go:366) - Sequences Either sequence operations
- [`iterator/iterresult.Chain`](../iterator/iterresult/ioeither.go:326) - Sequences Result sequence operations
- [`optics/codec/decode.Chain`](../optics/codec/decode/monad.go:163) - Sequences decoder operations
- [`optics/codec/validation.Chain`](../optics/codec/validation/monad.go:203) - Sequences validation operations
- [`optics/codec/validate.Chain`](../optics/codec/validate/validate.go:409) - Sequences validator operations

**Example**:
```go
// Option - dependent operations
safeDivide := func(x int) option.Option[int] {
    if x == 0 { return option.None[int]() }
    return option.Some(100 / x)
}
result := F.Pipe1(
    option.Some(5),
    option.Chain(safeDivide),
) // Some(20)

// Either - error handling pipeline
parseAndValidate := F.Pipe1(
    either.Right[error]("42"),
    either.Chain(func(s string) either.Either[error, int] {
        return result.Eitherize1(strconv.Atoi)(s)
    }),
) // Right(42)
```

---

## Transformation Operations

### Fold

**Concept**: Collapses a monadic value into a plain value by providing handlers for all cases.

**Type Signature**: `func Fold[A, B any](onError func(E) B, onSuccess func(A) B) func(M[A]) B`

**Purpose**: Extracts the value from a monadic context by handling all possible cases (error/success, none/some, etc.).

**Implementations**:

- [`option.Fold`](../option/option.go:225) - Handles None and Some cases
- [`either.Fold`](../either/either.go:448) - Handles Left and Right cases
- [`result.Fold`](../result/either.go:377) - Handles error and success cases
- [`ioeither.Fold`](../ioeither/ioeither.go:283) - Converts IOEither to IO
- [`ioresult.Fold`](../ioresult/ioeither.go:278) - Converts IOResult to IO
- [`iooption.Fold`](../iooption/iooption.go:176) - Converts IOOption to IO
- [`array.Fold`](../array/array.go:1063) - Reduces array using Monoid
- [`readereither.Fold`](../readereither/reader.go:101) - Converts ReaderEither to Reader
- [`readeroption.Fold`](../readeroption/reader.go:189) - Converts ReaderOption to Reader
- [`readerresult.Fold`](../readerresult/reader.go:476) - Converts ReaderResult to Reader
- [`readerioeither.Fold`](../readerioeither/reader.go:1226) - Converts ReaderIOEither to ReaderIO
- [`readerioresult.Fold`](../readerioresult/reader.go:748) - Converts ReaderIOResult to ReaderIO
- [`readeriooption.Fold`](../readeriooption/reader.go:209) - Converts ReaderIOOption to ReaderIO
- [`record.Fold`](../record/record.go:799) - Reduces record values using Monoid
- [`iterator/stateless.Fold`](../iterator/stateless/iterator.go:140) - Reduces iterator using Monoid
- [`iterator/iter.Fold`](../iterator/iter/iter.go:883) - Reduces sequence using Monoid
- [`iterator/itereither.Fold`](../iterator/itereither/ioeither.go:458) - Converts SeqEither to Seq
- [`iterator/iterresult.Fold`](../iterator/iterresult/ioeither.go:396) - Converts SeqResult to Seq

**Example**:
```go
// Option
result := F.Pipe1(
    option.Some(42),
    option.Fold(
        func() string { return "no value" },
        func(x int) string { return fmt.Sprintf("value: %d", x) },
    ),
) // "value: 42"

// Either
result := F.Pipe1(
    either.Right[error](42),
    either.Fold(
        func(e error) string { return "error: " + e.Error() },
        func(x int) string { return fmt.Sprintf("success: %d", x) },
    ),
) // "success: 42"
```

---

### Filter

**Concept**: Keeps only values that satisfy a predicate, converting failures to the error case.

**Type Signature**: `func Filter[A any](pred func(A) bool) func(M[A]) M[A]`

**Purpose**: Conditionally preserves values based on a predicate, useful for validation and conditional logic.

**Implementations**:

- [`option.Filter`](../option/option.go:440) - Converts Some to None if predicate fails
- [`either.Filter`](../either/filterable.go:159) - Converts Right to Left if predicate fails
- [`result.Filter`](../result/filterable.go:142) - Converts success to error if predicate fails
- [`array.Filter`](../array/array.go:234) - Keeps only elements matching predicate
- [`record.Filter`](../record/record.go:568) - Keeps only record entries matching predicate
- [`iterator/iter.Filter`](../iterator/iter/iter.go:332) - Filters sequence elements
- [`iterator/itereither.Filter`](../iterator/itereither/filter.go:153) - Filters successful sequence elements
- [`iterator/iterresult.Filter`](../iterator/iterresult/filter.go:160) - Filters successful result sequence

**Example**:
```go
// Option
isPositive := option.Filter(func(x int) bool { return x > 0 })
result1 := isPositive(option.Some(42))  // Some(42)
result2 := isPositive(option.Some(-5))  // None

// Array
evens := array.Filter(func(x int) bool { return x%2 == 0 })
result := evens([]int{1, 2, 3, 4, 5}) // [2, 4]
```

---

## Composition Operations

### Sequence

**Concept**: Inverts the nesting of two monadic types, swapping their order.

**Type Signature**: `func Sequence[A any](mma M[N[A]]) N[M[A]]`

**Purpose**: Transforms a monad containing another monad into the reverse nesting. Essential for working with nested effects.

**Implementations**:

- [`option.Sequence`](../option/sequence.go:36) - Inverts Option[M[A]] to M[Option[A]]
- [`either.Sequence`](../either/traverse.go:63) - Inverts Either[E, M[A]] to M[Either[E, A]]
- [`result.Sequence`](../result/traverse.go:61) - Inverts Result[M[A]] to M[Result[A]]
- [`array.Sequence`](../array/sequence.go:66) - Inverts []M[A] to M[[]A]
- [`reader.Sequence`](../reader/flip.go:65) - Flips Reader nesting order
- [`readereither.Sequence`](../readereither/flip.go:81) - Flips ReaderEither nesting
- [`readeroption.Sequence`](../readeroption/flip.go:80) - Flips ReaderOption nesting
- [`readerresult.Sequence`](../readerresult/flip.go:84) - Flips ReaderResult nesting
- [`readerio.Sequence`](../readerio/flip.go:74) - Flips ReaderIO nesting
- [`readerioeither.Sequence`](../readerioeither/flip.go:45) - Flips ReaderIOEither nesting
- [`readerioresult.Sequence`](../readerioresult/flip.go:9) - Flips ReaderIOResult nesting
- [`readerreaderioeither.Sequence`](../readerreaderioeither/flip.go:76) - Flips nested Reader nesting

**Example**:
```go
// Option containing array to array of options
opts := option.Some([]int{1, 2, 3})
result := option.Sequence(
    array.Of[option.Option[int]],
    array.Map(option.Some[int]),
)(opts) // [Some(1), Some(2), Some(3)]

// Array of Options to Option of array
opts := []option.Option[int]{
    option.Some(1),
    option.Some(2),
    option.Some(3),
}
result := array.Sequence(
    option.Of[[]int],
    option.Map(array.Of[int]),
)(opts) // Some([1, 2, 3])
```

---

### Traverse

**Concept**: Maps each element with an effect-producing function, then sequences the results.

**Type Signature**: `func Traverse[A, B any](f func(A) M[B]) func([]A) M[[]B]`

**Purpose**: Combines mapping and sequencing in one operation. More efficient than Map followed by Sequence.

**Implementations**:

- [`option.Traverse`](../option/sequence.go:60) - Maps and sequences Option operations
- [`either.Traverse`](../either/traverse.go:36) - Maps and sequences Either operations
- [`result.Traverse`](../result/traverse.go:38) - Maps and sequences Result operations
- [`array.Traverse`](../array/traverse.go:61) - Maps and sequences array operations
- [`record.Traverse`](../record/traverse.go:91) - Maps and sequences record operations
- [`reader.Traverse`](../reader/flip.go:124) - Maps and sequences Reader operations
- [`readereither.Traverse`](../readereither/flip.go:210) - Maps and sequences ReaderEither operations
- [`readeroption.Traverse`](../readeroption/flip.go:145) - Maps and sequences ReaderOption operations
- [`readerresult.Traverse`](../readerresult/flip.go:148) - Maps and sequences ReaderResult operations
- [`readerio.Traverse`](../readerio/flip.go:147) - Maps and sequences ReaderIO operations
- [`readerioeither.Traverse`](../readerioeither/flip.go:158) - Maps and sequences ReaderIOEither operations
- [`readerioresult.Traverse`](../readerioresult/flip.go:34) - Maps and sequences ReaderIOResult operations
- [`readerreaderioeither.Traverse`](../readerreaderioeither/flip.go:247) - Maps and sequences nested Reader operations

**Example**:
```go
// Parse array of strings to Option of array of ints
parseInts := array.Traverse(
    option.Of[[]int],
    option.Map(array.Of[int]),
)(func(s string) option.Option[int] {
    n, err := strconv.Atoi(s)
    if err != nil { return option.None[int]() }
    return option.Some(n)
})

result := parseInts([]string{"1", "2", "3"}) // Some([1, 2, 3])
```

---

## Do-Notation Operations

### Do

**Concept**: Creates an empty context to start building up a result using do-notation style.

**Type Signature**: `func Do[S any](empty S) M[S]`

**Purpose**: Initializes a context structure that will be populated step-by-step using Bind and Let operations.

**Implementations**:

- [`option.Do`](../option/bind.go:36) - Creates initial Option context
- [`either.Do`](../either/bind.go:34) - Creates initial Either context

**Example**:
```go
type Result struct {
    x int
    y string
}

result := option.Do(Result{})
// Result: Some(Result{x: 0, y: ""})
```

---

### Bind

**Concept**: Attaches the result of a monadic computation to a context, building up a larger structure.

**Type Signature**: `func Bind[S1, S2, A any](setter func(A) func(S1) S2, f func(S1) M[A]) func(M[S1]) M[S2]`

**Purpose**: Enables sequential building of complex structures in do-notation style, where each step can depend on previous results.

**Implementations**:

- [`option.Bind`](../option/bind.go:54) - Binds Option computation result to context
- [`either.Bind`](../either/bind.go:59) - Binds Either computation result to context

**Example**:
```go
type State struct {
    x int
    y int
}

result := F.Pipe2(
    option.Do(State{}),
    option.Bind(
        func(x int) func(State) State {
            return func(s State) State { s.x = x; return s }
        },
        func(s State) option.Option[int] { return option.Some(42) },
    ),
    option.Bind(
        func(y int) func(State) State {
            return func(s State) State { s.y = y; return s }
        },
        func(s State) option.Option[int] { return option.Some(s.x * 2) },
    ),
) // Some(State{x: 42, y: 84})
```

---

### Let

**Concept**: Attaches the result of a pure computation to a context, similar to Bind but for non-monadic functions.

**Type Signature**: `func Let[S1, S2, B any](setter func(B) func(S1) S2, f func(S1) B) func(M[S1]) M[S2]`

**Purpose**: Adds computed values to the context without introducing monadic effects.

**Implementations**:

- [`option.Let`](../option/bind.go:78) - Adds pure computation result to Option context
- [`either.Let`](../either/bind.go:88) - Adds pure computation result to Either context

**Example**:
```go
type State struct {
    x int
    computed int
}

result := F.Pipe2(
    option.Do(State{x: 5}),
    option.Let(
        func(c int) func(State) State {
            return func(s State) State { s.computed = c; return s }
        },
        func(s State) int { return s.x * 2 },
    ),
) // Some(State{x: 5, computed: 10})
```

### ApS

**Concept**: Applies a wrapped value to a context using a setter function, combining applicative application with context building.

**Type Signature**: `func ApS[S1, S2, T any](setter func(T) func(S1) S2, fa M[T]) func(M[S1]) M[S2]`

**Purpose**: Enables building up complex structures by applying independent computations and setting their results into a context. This is the applicative version of Bind for do-notation.

**Implementations**:

- [`option.ApS`](../option/bind.go:142)
- [`either.ApS`](../either/bind.go:164)
- [`result.ApS`](../result/bind.go:145)
- [`io.ApS`](../io/bind.go:146)
- [`ioeither.ApS`](../ioeither/bind.go:158)
- [`ioresult.ApS`](../ioresult/bind.go:151)
- [`array.ApS`](../array/bind.go:163)
- [`lazy.ApS`](../lazy/bind.go:136)
- [`reader.ApS`](../reader/bind.go:202)
- [`readereither.ApS`](../readereither/bind.go:214)
- [`readerresult.ApS`](../readerresult/bind.go:213)
- [`state.ApS`](../state/bind.go:274)
- [`effect.ApS`](../effect/bind.go:247)
- [`identity.ApS`](../identity/bind.go:216)
- [`record.ApS`](../record/bind.go:186)
- And 15+ more implementations across Reader*, State*, Iterator*, and Optics modules

**Example**:
```go
type State struct { x, y int }

result := F.Pipe2(
    option.Do(State{}),
    option.ApS(
        func(x int) func(State) State {
            return func(s State) State { s.x = x; return s }
        },
        option.Some(42),
    ),
    option.ApS(
        func(y int) func(State) State {
            return func(s State) State { s.y = y; return s }
        },
        option.Some(100),
    ),
) // Some(State{x: 42, y: 100})
```

---

### ChainFirst

**Concept**: Executes a monadic computation for its side effects but returns the original value.

**Type Signature**: `func ChainFirst[A, B any](f func(A) M[B]) func(M[A]) M[A]`

**Purpose**: Allows performing side-effecting operations (logging, validation, etc.) while preserving the original value in the pipeline.

**Implementations**:

- [`option.ChainFirst`](../option/option.go:346)
- [`either.ChainFirst`](../either/either.go:329)
- [`result.ChainFirst`](../result/either.go:255)
- [`io.ChainFirst`](../io/io.go:225)
- [`ioeither.ChainFirst`](../ioeither/ioeither.go:333)
- [`ioresult.ChainFirst`](../ioresult/ioeither.go:323)
- [`lazy.ChainFirst`](../lazy/lazy.go:183)
- [`readerio.ChainFirst`](../readerio/reader.go:401)
- [`readerioeither.ChainFirst`](../readerioeither/reader.go:1073)
- [`state.ChainFirst`](../state/state.go:254)
- [`effect.ChainFirst`](../effect/effect.go:202)
- [`identity.ChainFirst`](../identity/identity.go:208)
- [`endomorphism.ChainFirst`](../endomorphism/endo.go:279)
- [`iterator/stateless.ChainFirst`](../iterator/stateless/iterator.go:148)
- [`iterator/itereither.ChainFirst`](../iterator/itereither/ioeither.go:508)

**Example**:
```go
logValue := option.ChainFirst(func(x int) option.Option[string] {
    fmt.Printf("Value: %d\n", x)
    return option.Some(fmt.Sprintf("logged: %d", x))
})

result := F.Pipe1(option.Some(42), logValue)
// Prints "Value: 42", returns Some(42)
```

---

### Tap

**Concept**: Alias for ChainFirst - executes a side effect while preserving the original value.

**Type Signature**: `func Tap[A, B any](f func(A) M[B]) func(M[A]) M[A]`

**Purpose**: More intuitive name for ChainFirst, commonly used for logging, debugging, or validation side effects.

**Implementations**:

- [`ioeither.Tap`](../ioeither/ioeither.go:342)
- [`ioresult.Tap`](../ioresult/ioeither.go:328)
- [`readerio.Tap`](../readerio/reader.go:431)
- [`readerioeither.Tap`](../readerioeither/reader.go:1100)
- [`readerreaderioeither.Tap`](../readerreaderioeither/reader.go:616)
- [`effect.Tap`](../effect/effect.go:497)
- [`iterator/itereither.Tap`](../iterator/itereither/ioeither.go:517)

**Example**:
```go
result := F.Pipe1(
    ioresult.Of(42),
    ioresult.Tap(func(x int) ioresult.IOResult[string] {
        return ioresult.Of(fmt.Sprintf("Processing: %d", x))
    }),
)() // Returns (42, nil) after side effect
```

---

### Flatten

**Concept**: Removes one level of nesting from a nested monadic structure.

**Type Signature**: `func Flatten[A any](mma M[M[A]]) M[A]`

**Purpose**: Collapses nested monads into a single layer. Equivalent to `Chain(identity)`.

**Implementations**:

- [`option.Flatten`](../option/option.go:362)
- [`either.Flatten`](../either/either.go:343)
- [`result.Flatten`](../result/either.go:267)
- [`io.Flatten`](../io/io.go:208)
- [`ioeither.Flatten`](../ioeither/ioeither.go:231)
- [`ioresult.Flatten`](../ioresult/ioeither.go:234)
- [`array.Flatten`](../array/array.go:768)
- [`lazy.Flatten`](../lazy/lazy.go:166)
- [`reader.Flatten`](../reader/reader.go:429)
- [`readereither.Flatten`](../readereither/reader.go:147)
- [`readerresult.Flatten`](../readerresult/reader.go:703)
- [`readerio.Flatten`](../readerio/reader.go:849)
- [`readerioeither.Flatten`](../readerioeither/reader.go:1132)
- [`state.Flatten`](../state/state.go:274)
- [`record.Flatten`](../record/record.go:713)
- [`endomorphism.Flatten`](../endomorphism/endo.go:347)
- [`iterator/iter.Flatten`](../iterator/iter/iter.go:636)

**Example**:
```go
nested := option.Some(option.Some(42))
result := option.Flatten(nested) // Some(42)

nested := [][]int{{1, 2}, {3, 4}}
result := array.Flatten(nested) // [1, 2, 3, 4]
```

---

### Alt

**Concept**: Provides an alternative computation if the first one fails or is empty.

**Type Signature**: `func Alt[A any](second Lazy[M[A]]) func(M[A]) M[A]`

**Purpose**: Implements the Alternative type class, allowing fallback to a second computation when the first fails.

**Implementations**:

- [`option.Alt`](../option/option.go:384)
- [`either.Alt`](../either/either.go:540)
- [`result.Alt`](../result/either.go:469)
- [`ioeither.Alt`](../ioeither/ioeither.go:456)
- [`ioresult.Alt`](../ioresult/ioeither.go:428)
- [`readereither.Alt`](../readereither/reader.go:485)
- [`readerresult.Alt`](../readerresult/reader.go:868)
- [`readerioeither.Alt`](../readerioeither/reader.go:1317)
- [`readerreaderioeither.Alt`](../readerreaderioeither/reader.go:757)
- [`iterator/itereither.Alt`](../iterator/itereither/ioeither.go:616)
- [`optics/codec/decode.Alt`](../optics/codec/decode/monad.go:602)
- [`optics/codec/validation.Alt`](../optics/codec/validation/monad.go:639)

**Example**:
```go
primary := option.None[int]()
fallback := option.Alt(func() option.Option[int] {
    return option.Some(42)
})
result := fallback(primary) // Some(42)
```

---

### OrElse

**Concept**: Recovers from an error by providing an alternative computation based on the error value.

**Type Signature**: `func OrElse[E, A any](onLeft func(E) M[A]) func(M[A]) M[A]`

**Purpose**: Allows error recovery by inspecting the error and providing a fallback computation.

**Implementations**:

- [`either.OrElse`](../either/either.go:567)
- [`result.OrElse`](../result/either.go:495)
- [`ioeither.OrElse`](../ioeither/ioeither.go:669)
- [`ioresult.OrElse`](../ioresult/ioeither.go:514)
- [`readereither.OrElse`](../readereither/reader.go:272)
- [`readerresult.OrElse`](../readerresult/reader.go:507)
- [`readerioeither.OrElse`](../readerioeither/reader.go:1686)
- [`iterator/itereither.OrElse`](../iterator/itereither/ioeither.go:811)
- [`optics/codec/decode.OrElse`](../optics/codec/decode/monad.go:372)
- [`optics/codec/validation.OrElse`](../optics/codec/validation/monad.go:474)

**Example**:
```go
recover := result.OrElse(func(err error) result.Result[int] {
    if err.Error() == "not found" {
        return result.Of(0)
    }
    return result.Left[int](err)
})
result := recover(result.Left[int](errors.New("not found"))) // Right(0)
```

---

### GetOrElse

**Concept**: Extracts the value from a monad or provides a default if it's empty/failed.

**Type Signature**: `func GetOrElse[A any](onEmpty func() A) func(M[A]) A`

**Purpose**: Safely extracts values from monadic contexts with a fallback for failure cases.

**Implementations**:

- [`option.GetOrElse`](../option/option.go:253)
- [`either.GetOrElse`](../either/either.go:504)
- [`result.GetOrElse`](../result/either.go:432)
- [`ioeither.GetOrElse`](../ioeither/ioeither.go:288)
- [`ioresult.GetOrElse`](../ioresult/ioeither.go:285)
- [`readereither.GetOrElse`](../readereither/reader.go:105)
- [`readeroption.GetOrElse`](../readeroption/reader.go:220)
- [`readerresult.GetOrElse`](../readerresult/reader.go:491)
- [`readerioeither.GetOrElse`](../readerioeither/reader.go:1239)
- [`iterator/itereither.GetOrElse`](../iterator/itereither/ioeither.go:463)

**Example**:
```go
result := F.Pipe1(
    option.None[int](),
    option.GetOrElse(func() int { return 42 }),
) // 42
```

---

### BiMap

**Concept**: Maps two functions over both type parameters of a bifunctor simultaneously.

**Type Signature**: `func BiMap[E1, E2, A, B any](f func(E1) E2, g func(A) B) func(M[E1, A]) M[E2, B]`

**Purpose**: Transforms both the error and success channels of a computation in one operation.

**Implementations**:

- [`either.BiMap`](../either/either.go:121)
- [`result.BiMap`](../result/either.go:114)
- [`ioeither.BiMap`](../ioeither/ioeither.go:278)
- [`ioresult.BiMap`](../ioresult/ioeither.go:271)
- [`readereither.BiMap`](../readereither/reader.go:156)
- [`readerresult.BiMap`](../readerresult/reader.go:749)
- [`readerioeither.BiMap`](../readerioeither/reader.go:1271)
- [`pair.BiMap`](../pair/pair.go:257)
- [`tuple.BiMap`](../tuple/tuple.go:102)
- [`iterator/itereither.BiMap`](../iterator/itereither/ioeither.go:453)

**Example**:
```go
transform := either.BiMap(
    func(e error) string { return "Error: " + e.Error() },
    func(x int) string { return fmt.Sprintf("Value: %d", x) },
)
result := transform(either.Right[error](42)) // Right("Value: 42")
```

---

### MapLeft

**Concept**: Transforms only the error/left channel of a bifunctor.

**Type Signature**: `func MapLeft[E1, E2, A any](f func(E1) E2) func(M[E1, A]) M[E2, A]`

**Purpose**: Allows transforming error values without affecting success values.

**Implementations**:

- [`either.MapLeft`](../either/either.go:166)
- [`result.MapLeft`](../result/either.go:164)
- [`ioeither.MapLeft`](../ioeither/ioeither.go:265)
- [`ioresult.MapLeft`](../ioresult/ioeither.go:259)
- [`readereither.MapLeft`](../readereither/reader.go:244)
- [`readerresult.MapLeft`](../readerresult/reader.go:826)
- [`readerioeither.MapLeft`](../readerioeither/reader.go:1362)
- [`readerreaderioeither.MapLeft`](../readerreaderioeither/reader.go:787)
- [`iterator/itereither.MapLeft`](../iterator/itereither/ioeither.go:432)

**Example**:
```go
enrichError := result.MapLeft(func(err error) error {
    return fmt.Errorf("operation failed: %w", err)
})
result := enrichError(result.Left[int](errors.New("timeout")))
// Left(error: "operation failed: timeout")
```

---

---

## Summary

This document covers the most common operations across monads in fp-go/v2:

- **Core Operations**: Of, Map, Ap, Chain - the fundamental type class operations
- **Transformation**: Fold, Filter - extracting and filtering values
- **Composition**: Sequence, Traverse - working with nested effects
- **Do-Notation**: Do, Bind, Let - building complex structures step-by-step

Each operation maintains the laws and properties of its type class, ensuring predictable and composable behavior across all implementations.

## Summary

This document covers the most common operations across monads in fp-go/v2:

- **Core Operations**: Of, Map, Ap, Chain - the fundamental type class operations
- **Transformation**: Fold, Filter - extracting and filtering values
- **Composition**: Sequence, Traverse - working with nested effects
- **Do-Notation**: Do, Bind, Let, ApS - building complex structures step-by-step
- **Side Effects**: ChainFirst, Tap - executing side effects while preserving values
- **Structure**: Flatten - collapsing nested monads
- **Alternatives**: Alt, OrElse, GetOrElse - providing fallbacks and recovery
- **Bifunctors**: BiMap, MapLeft - transforming error channels

Each operation maintains the laws and properties of its type class, ensuring predictable and composable behavior across all implementations.

---

## Monoid Operations

### Overview

A **Monoid** is an algebraic structure that extends Semigroup by adding an identity element. It consists of:
- A type `A`
- An associative binary operation `Concat: (A, A) → A`
- An identity element `Empty: () → A`

Monoids satisfy these laws:
1. **Associativity**: `Concat(Concat(x, y), z) = Concat(x, Concat(y, z))`
2. **Left Identity**: `Concat(Empty(), x) = x`
3. **Right Identity**: `Concat(x, Empty()) = x`

Monoids are fundamental for combining values in a consistent way, enabling operations like folding, reducing, and accumulating results across collections or computational contexts.

### Core Monoid Functions

#### MakeMonoid

**Concept**: Creates a monoid instance from a binary operation and an identity element.

**Type Signature**: `func MakeMonoid[A any](concat func(A, A) A, empty A) Monoid[A]`

**Purpose**: The primary constructor for creating custom monoid instances.

**Implementation**: [`monoid.MakeMonoid`](../monoid/monoid.go:150)

**Example**:
```go
// Integer addition monoid
addMonoid := monoid.MakeMonoid(
    func(a, b int) int { return a + b },
    0,  // identity element
)
result := addMonoid.Concat(5, 3)  // 8
empty := addMonoid.Empty()         // 0
```

#### ConcatAll / Fold

**Concept**: Combines all elements in a collection using the monoid's operation, starting with the identity element.

**Type Signature**: `func ConcatAll[A any](m Monoid[A]) func([]A) A`

**Purpose**: Reduces a collection to a single value using the monoid's combining operation.

**Implementations**:
- [`monoid.ConcatAll`](../monoid/array.go:85) - Standard version for slices
- [`monoid.Fold`](../monoid/array.go:111) - Alias for ConcatAll
- [`monoid.GenericConcatAll`](../monoid/array.go:50) - Generic version for custom slice types

**Example**:
```go
addMonoid := monoid.MakeMonoid(
    func(a, b int) int { return a + b },
    0,
)
sum := monoid.ConcatAll(addMonoid)([]int{1, 2, 3, 4, 5})  // 15
empty := monoid.ConcatAll(addMonoid)([]int{})              // 0
```

#### Reverse

**Concept**: Creates the dual of a monoid by swapping the order of arguments in the binary operation.

**Type Signature**: `func Reverse[A any](m Monoid[A]) Monoid[A]`

**Purpose**: Useful for non-commutative operations where order matters.

**Implementation**: [`monoid.Reverse`](../monoid/monoid.go:194)

**Example**:
```go
subMonoid := monoid.MakeMonoid(
    func(a, b int) int { return a - b },
    0,
)
reversedMonoid := monoid.Reverse(subMonoid)

result1 := subMonoid.Concat(10, 3)         // 10 - 3 = 7
result2 := reversedMonoid.Concat(10, 3)    // 3 - 10 = -7
```

### Specialized Monoids

#### Numeric Monoids

**Implementations**:
- [`number.MonoidSum`](../number/monoid.go:23) - Addition with 0 as identity
- [`number.MonoidProduct`](../number/monoid.go:32) - Multiplication with 1 as identity

**Example**:
```go
sumMonoid := number.MonoidSum[int]()
product := number.MonoidProduct[int]()

sum := monoid.ConcatAll(sumMonoid)([]int{1, 2, 3, 4})      // 10
prod := monoid.ConcatAll(product)([]int{2, 3, 4})          // 24
```

#### Array/Slice Monoids

**Implementations**:
- [`array.Monoid`](../array/monoid.go:35) - Concatenates slices with empty slice as identity
- [`array/generic.Monoid`](../array/generic/monoid.go:19) - Generic version for custom slice types

**Example**:
```go
sliceMonoid := array.Monoid[int]()
combined := sliceMonoid.Concat([]int{1, 2}, []int{3, 4})  // [1, 2, 3, 4]
empty := sliceMonoid.Empty()                               // []
```

#### Iterator Monoids

**Implementations**:
- [`iterator/iter.Monoid`](../iterator/iter/monoid.go:41) - Concatenates iterators
- [`iterator/iter.ConcatMonoid`](../iterator/iter/monoid.go:168) - Alias for Monoid
- [`iterator/iter.MergeMonoid`](../iterator/iter/mergeall.go:243) - Merges iterators with buffering
- [`iterator/stateless.Monoid`](../iterator/stateless/monoid.go:24) - For stateless iterators

**Example**:
```go
seqMonoid := iter.Monoid[int]()
seq1 := slices.Values([]int{1, 2, 3})
seq2 := slices.Values([]int{4, 5, 6})
combined := seqMonoid.Concat(seq1, seq2)  // Yields 1, 2, 3, 4, 5, 6
```

#### String and Record Monoids

**Implementations**:
- String concatenation: Use `monoid.MakeMonoid(func(a, b string) string { return a + b }, "")`
- [`record.UnionMonoid`](../record/monoid.go:26) - Merges maps with semigroup for values
- [`record.UnionLastMonoid`](../record/monoid.go:33) - Merges maps, keeping last value
- [`record.UnionFirstMonoid`](../record/monoid.go:40) - Merges maps, keeping first value
- [`record.MergeMonoid`](../record/monoid.go:47) - Alias for UnionLastMonoid

#### Tuple Monoids

**Implementations**: Monoids for tuples of 1-15 elements, combining component-wise
- [`tuple.Monoid1`](../tuple/gen.go:215) through [`tuple.Monoid15`](../tuple/gen.go:1895)

**Example**:
```go
// Monoid for pairs of integers (using addition for both)
pairMonoid := tuple.Monoid2(
    number.MonoidSum[int](),
    number.MonoidSum[int](),
)
result := pairMonoid.Concat(
    tuple.MakeTuple2(1, 2),
    tuple.MakeTuple2(3, 4),
)  // (4, 6)
```

#### Boolean Monoids

**Implementations**:
- Boolean AND: `monoid.MakeMonoid(func(a, b bool) bool { return a && b }, true)`
- Boolean OR: `monoid.MakeMonoid(func(a, b bool) bool { return a || b }, false)`

**Example**:
```go
andMonoid := monoid.MakeMonoid(
    func(a, b bool) bool { return a && b },
    true,
)
result := monoid.ConcatAll(andMonoid)([]bool{true, true, false})  // false
```

#### Predicate Monoids

**Implementations**:
- [`predicate.MonoidAny`](../predicate/monoid.go:111) - True if any predicate matches (OR semantics)
- [`predicate.MonoidAll`](../predicate/monoid.go:139) - True if all predicates match (AND semantics)

**Example**:
```go
anyMonoid := predicate.MonoidAny[int]()
isPositive := func(n int) bool { return n > 0 }
isEven := func(n int) bool { return n%2 == 0 }
combined := anyMonoid.Concat(isPositive, isEven)
combined(3)   // true (positive)
combined(-2)  // true (even)
combined(-3)  // false (neither)
```

#### Function Monoids

**Implementations**:
- [`monoid.FunctionMonoid`](../monoid/function.go:66) - Combines functions point-wise when codomain has a monoid
- [`endomorphism.Monoid`](../endomorphism/monoid.go:149) - Function composition with identity

**Example**:
```go
// Functions returning integers
intAddMonoid := number.MonoidSum[int]()
funcMonoid := monoid.FunctionMonoid[string, int](intAddMonoid)

f1 := func(s string) int { return len(s) }
f2 := func(s string) int { return len(s) * 2 }

combined := funcMonoid.Concat(f1, f2)
result := combined("hello")  // len("hello") + len("hello")*2 = 5 + 10 = 15
```

#### Pair Monoids

**Implementations**:
- [`pair.Monoid`](../pair/monoid.go:88) - Combines pairs component-wise
- [`pair.ApplicativeMonoid`](../pair/monoid.go:139) - Applicative-based combination
- [`pair.ApplicativeMonoidTail`](../pair/monoid.go:222) - Focuses on tail (right) component
- [`pair.ApplicativeMonoidHead`](../pair/monoid.go:309) - Focuses on head (left) component

#### Comparison Monoids

**Implementations**:
- [`ord.Monoid`](../ord/monoid.go:67) - Combines ordering comparisons
- [`eq.Monoid`](../eq/monoid.go:151) - Combines equality predicates

### Monad-Specific Monoids

Many monads provide specialized monoid instances that combine monadic values:

#### Option Monoids

**Implementations**:
- [`option.Monoid`](../option/monoid.go:71) - Requires semigroup for inner type
- [`option.ApplicativeMonoid`](../option/apply.go:47) - Combines Some values, None if any is None
- [`option.AlternativeMonoid`](../option/monoid.go:88) - First Some value, or None
- [`option.AltMonoid`](../option/monoid.go:110) - Alternative with lazy default
- [`option.FirstMonoid`](../option/monoid.go:149) - Takes first Some
- [`option.LastMonoid`](../option/monoid.go:183) - Takes last Some

**Example**:
```go
// ApplicativeMonoid - combines values inside Some
m := option.ApplicativeMonoid(number.MonoidSum[int]())
result := m.Concat(option.Some(5), option.Some(3))  // Some(8)
result2 := m.Concat(option.Some(5), option.None[int]())  // None

// FirstMonoid - takes first Some
firstM := option.FirstMonoid[int]()
result3 := firstM.Concat(option.Some(5), option.Some(3))  // Some(5)
```

#### Either/Result Monoids

**Implementations**:
- [`either.ApplicativeMonoid`](../either/apply.go:48) - Combines Right values
- [`either.AlternativeMonoid`](../either/monoid.go:33) - First Right, or accumulates Lefts
- [`either.AltMonoid`](../either/monoid.go:53) - Alternative with lazy default
- [`either.FirstMonoid`](../either/monoid.go:94) - Takes first Right
- [`either.LastMonoid`](../either/monoid.go:130) - Takes last Right
- [`result.ApplicativeMonoid`](../result/apply.go:49) - Result-specific version
- [`result.AlternativeMonoid`](../result/monoid.go:36) - Result-specific version
- [`result.AltMonoid`](../result/monoid.go:52) - Result-specific version
- [`result.FirstMonoid`](../result/monoid.go:82) - Result-specific version
- [`result.LastMonoid`](../result/monoid.go:110) - Result-specific version

#### IO Monoids

**Implementations**:
- [`io.ApplicativeMonoid`](../io/apply.go:45) - Combines IO effects
- [`lazy.ApplicativeMonoid`](../lazy/apply.go:92) - For lazy computations

#### Reader Monoids

**Implementations**:
- [`reader.ApplicativeMonoid`](../reader/semigroup.go:49) - Combines reader results
- [`reader/generic.ApplicativeMonoid`](../reader/generic/monoid.go:8) - Generic version
- [`readereither.ApplicativeMonoid`](../readereither/monoid.go:50) - For ReaderEither
- [`readereither.AlternativeMonoid`](../readereither/monoid.go:89) - Alternative semantics
- [`readereither.AltMonoid`](../readereither/monoid.go:130) - With lazy default
- [`readeroption.ApplicativeMonoid`](../readeroption/monoid.go:68) - For ReaderOption
- [`readeroption.AlternativeMonoid`](../readeroption/monoid.go:140) - Alternative semantics
- [`readeroption.AltMonoid`](../readeroption/monoid.go:209) - With lazy default
- [`readerresult.ApplicativeMonoid`](../readerresult/monoid.go:97) - For ReaderResult
- [`readerresult.AlternativeMonoid`](../readerresult/monoid.go:39) - Alternative semantics
- [`readerresult.AltMonoid`](../readerresult/monoid.go:68) - With lazy default

#### ReaderIO Monoids

**Implementations**:
- [`readerioeither.ApplicativeMonoid`](../readerioeither/monoid.go:37) - Sequential execution
- [`readerioeither.ApplicativeMonoidSeq`](../readerioeither/monoid.go:53) - Explicit sequential
- [`readerioeither.ApplicativeMonoidPar`](../readerioeither/monoid.go:69) - Parallel execution
- [`readerioeither.AlternativeMonoid`](../readerioeither/monoid.go:86) - Alternative semantics
- [`readerioeither.AltMonoid`](../readerioeither/monoid.go:104) - With lazy default
- [`readeriooption.ApplicativeMonoid`](../readeriooption/monoid.go:68) - For ReaderIOOption
- [`readeriooption.AlternativeMonoid`](../readeriooption/monoid.go:140) - Alternative semantics
- [`readeriooption.AltMonoid`](../readeriooption/monoid.go:211) - With lazy default
- [`readerioresult.ApplicativeMonoid`](../readerioresult/monoid.go:38) - For ReaderIOResult
- [`readerioresult.ApplicativeMonoidSeq`](../readerioresult/monoid.go:49) - Sequential
- [`readerioresult.ApplicativeMonoidPar`](../readerioresult/monoid.go:60) - Parallel
- [`readerioresult.AlternativeMonoid`](../readerioresult/monoid.go:72) - Alternative semantics
- [`readerioresult.AltMonoid`](../readerioresult/monoid.go:84) - With lazy default

#### ReaderReaderIO Monoids

**Implementations**:
- [`readerreaderioeither.ApplicativeMonoid`](../readerreaderioeither/monoid.go:26) - Default
- [`readerreaderioeither.ApplicativeMonoidSeq`](../readerreaderioeither/monoid.go:35) - Sequential
- [`readerreaderioeither.ApplicativeMonoidPar`](../readerreaderioeither/monoid.go:44) - Parallel
- [`readerreaderioeither.AlternativeMonoid`](../readerreaderioeither/monoid.go:53) - Alternative
- [`readerreaderioeither.AltMonoid`](../readerreaderioeither/monoid.go:63) - With lazy default

#### State Monoids

**Implementations**:
- [`state.ApplicativeMonoid`](../state/monoid.go:141) - Combines state computations
- [`stateio.ApplicativeMonoid`](../stateio/monoid.go:146) - For StateIO

#### Context-Aware Monoids

**Implementations**:
- [`context/readerioresult.ApplicativeMonoid`](../context/readerioresult/monoid.go:39) - Default
- [`context/readerioresult.ApplicativeMonoidSeq`](../context/readerioresult/monoid.go:50) - Sequential
- [`context/readerioresult.ApplicativeMonoidPar`](../context/readerioresult/monoid.go:61) - Parallel
- [`context/readerioresult.AlternativeMonoid`](../context/readerioresult/monoid.go:73) - Alternative
- [`context/readerioresult.AltMonoid`](../context/readerioresult/monoid.go:85) - With lazy default
- [`context/readerreaderioresult.ApplicativeMonoid`](../context/readerreaderioresult/monoid.go:49) - Default
- [`context/readerreaderioresult.ApplicativeMonoidSeq`](../context/readerreaderioresult/monoid.go:65) - Sequential
- [`context/readerreaderioresult.ApplicativeMonoidPar`](../context/readerreaderioresult/monoid.go:81) - Parallel
- [`context/readerreaderioresult.AlternativeMonoid`](../context/readerreaderioresult/monoid.go:110) - Alternative
- [`context/readerreaderioresult.AltMonoid`](../context/readerreaderioresult/monoid.go:143) - With lazy default
- [`idiomatic/context/readerresult.ApplicativeMonoid`](../idiomatic/context/readerresult/monoid.go:118) - Idiomatic version
- [`idiomatic/context/readerresult.AlternativeMonoid`](../idiomatic/context/readerresult/monoid.go:54) - Alternative
- [`idiomatic/context/readerresult.AltMonoid`](../idiomatic/context/readerresult/monoid.go:85) - With lazy default
- [`idiomatic/readerresult.ApplicativeMonoid`](../idiomatic/readerresult/monoid.go:97) - Idiomatic version
- [`idiomatic/readerresult.AlternativeMonoid`](../idiomatic/readerresult/monoid.go:39) - Alternative
- [`idiomatic/readerresult.AltMonoid`](../idiomatic/readerresult/monoid.go:68) - With lazy default

#### Effect Monoids

**Implementations**:
- [`effect.ApplicativeMonoid`](../effect/monoid.go:51) - Combines effects
- [`effect.AlternativeMonoid`](../effect/monoid.go:83) - Alternative semantics

#### Validation/Codec Monoids

**Implementations**:
- [`optics/codec.AltMonoid`](../optics/codec/alt.go:475) - For codec types
- [`optics/codec/decode.ApplicativeMonoid`](../optics/codec/decode/monoid.go:79) - Combines decoders
- [`optics/codec/decode.AlternativeMonoid`](../optics/codec/decode/monoid.go:214) - Alternative decoders
- [`optics/codec/decode.AltMonoid`](../optics/codec/decode/monoid.go:363) - With lazy default
- [`optics/codec/validate.ApplicativeMonoid`](../optics/codec/validate/monoid.go:117) - Combines validators
- [`optics/codec/validate.AlternativeMonoid`](../optics/codec/validate/monoid.go:244) - Alternative validators
- [`optics/codec/validate.AltMonoid`](../optics/codec/validate/monoid.go:384) - With lazy default
- [`optics/codec/validation.ApplicativeMonoid`](../optics/codec/validation/monoid.go:45) - Accumulates errors
- [`optics/codec/validation.AlternativeMonoid`](../optics/codec/validation/monoid.go:126) - First success
- [`optics/codec/validation.AltMonoid`](../optics/codec/validation/monoid.go:223) - With lazy default

#### Special Purpose Monoids

**Implementations**:
- [`monoid.VoidMonoid`](../monoid/void.go:60) - For unit/void type
- [`constant.Monoid`](../constant/monoid.go:64) - Constant monoid (always returns same value)

### Monoid Patterns

#### ApplicativeMonoid Pattern

Combines values inside applicative contexts. If all values are successful, combines them using the inner monoid. If any fail, the result is a failure.

**Common across**: Option, Either, Result, IO, Reader, and their combinations

**Example**:
```go
m := option.ApplicativeMonoid(number.MonoidSum[int]())
result := m.Concat(option.Some(5), option.Some(3))  // Some(8)
result2 := m.Concat(option.Some(5), option.None[int]())  // None
```

#### AlternativeMonoid Pattern

Provides fallback semantics. Takes the first successful value, or combines failures.

**Common across**: Option, Either, Result, and their Reader/IO variants

**Example**:
```go
m := result.AlternativeMonoid(number.MonoidSum[int]())
result := m.Concat(result.Of(5), result.Of(3))  // Right(5) - first success
result2 := m.Concat(result.Left(err1), result.Of(3))  // Right(3) - fallback
```

#### AltMonoid Pattern

Similar to AlternativeMonoid but with a lazy default value.

**Common across**: Option, Either, Result, and their Reader/IO variants

**Example**:
```go
defaultValue := lazy.Of(result.Of(42))
m := result.AltMonoid(defaultValue)
result := m.Concat(result.Left(err), result.Left(err2))  // Right(42) - uses default
```

### Usage Guidelines

1. **Choose the right monoid**: Select based on your combining semantics (addition, concatenation, first/last, etc.)
2. **Leverage existing monoids**: Use built-in monoids for common types before creating custom ones
3. **Compose monoids**: Build complex monoids from simpler ones (e.g., tuple monoids, function monoids)
4. **Consider execution strategy**: For IO-based monoids, choose between sequential and parallel execution
5. **Use appropriate pattern**: ApplicativeMonoid for combining successes, AlternativeMonoid for fallbacks

### Related Concepts

- **Semigroup**: A monoid without the identity element requirement
- **Magma**: A set with a binary operation (no laws required)
- **Group**: A monoid where every element has an inverse

For more details on monoid laws and theory, see the [`monoid` package documentation](../monoid/doc.go).