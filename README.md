# Functional programming library for golang

**🚧 Work in progress! 🚧** Despite major version 1 because of <https://github.com/semantic-release/semantic-release/issues/1507>. Trying to not make breaking changes, but devil is in the details.

![logo](resources/images/logo.png)

This library is strongly influenced by the awesome [fp-ts](https://github.com/gcanti/fp-ts).

## Getting started

```bash
go get github.com/IBM/fp-go
```

Refer to the [samples](./samples/).

## Design Goal

This library aims to provide a set of data types and functions that make it easy and fun to write maintainable and testable code in golang. It encourages the following patterns:

- write many small, testable and pure functions, i.e. functions that produce output only depending on their input and that do not execute side effects
- offer helpers to isolate side effects into lazily executed functions (IO)
- expose a consistent set of composition to create new functions from existing ones
  - for each data type there exists a small set of composition functions
  - these functions are called the same across all data types, so you only have to learn a small number of function names
  - the semantic of functions of the same name is consistent across all data types

### How does this play with the [🧘🏽 Zen Of Go](https://the-zen-of-go.netlify.app/)?

#### 🧘🏽 Each package fulfils a single purpose

✔️ Each of the top level packages (e.g. Option, Either, ReaderIOEither, ...) fulfils the purpose of defining the respective data type and implementing the set of common operations for this data type.

#### 🧘🏽 Handle errors explicitly

✔️ The library makes a clear distinction between that operations that cannot fail by design and operations that can fail. Failure is represented via the `Either` type and errors are handled explicitly by using `Either`'s monadic set of operations.

#### 🧘🏽 Return early rather than nesting deeply

✔️ We recommend to implement simple, small functions that implement one feature and that would typically not invoke other functions. Interaction with other functions is done by function composition and the composition makes sure to run one function after the other. In the error case the `Either` monad makes sure to skip the error path.

#### 🧘🏽 Leave concurrency to the caller

✔️ All pure are synchronous by default. The I/O operations are asynchronous per default.

#### 🧘🏽 Before you launch a goroutine, know when it will stop

🤷🏽 This is left to the user of the library since the library itself will not start goroutines on its own. The Task monad offers support for cancellation via the golang context, though.

#### 🧘🏽 Avoid package level state

✔️ No package level state anywhere, this would be a significant anti-pattern

#### 🧘🏽 Simplicity matters

✔️ The library is simple in the sense that it offers a small, consistent interface to a variety of data types. Users can concentrate on implementing business logic rather than dealing with low level data structures.

#### 🧘🏽 Write tests to lock in the behaviour of your package’s API

🟡 The programming pattern suggested by this library encourages writing test cases. The library itself also has a growing number of tests, but not enough, yet. TBD

#### 🧘🏽 If you think it’s slow, first prove it with a benchmark

✔️ Absolutely. If you think the function composition offered by this library is too slow, please provide a benchmark.

#### 🧘🏽 Moderation is a virtue

✔️ The library does not implement its own goroutines and also does not require any expensive synchronization primitives. Coordination of IO operations is implemented via atomic counters without additional primitives.

#### 🧘🏽 Maintainability counts

✔️ Code that consumes this library is easy to maintain because of the small and concise set of operations exposed. Also the suggested programming paradigm to decompose an application into small functions increases maintainability, because these functions are easy to understand and if they are pure, it's often sufficient to look at the type signature to understand the purpose.

The library itself also comprises many small functions, but it's admittedly harder to maintain than code that uses it. However this asymmetry is intended because it offloads complexity from users into a central component.

## Comparation to Idiomatic Go

In this section we discuss how the functional APIs differ from idiomatic go function signatures and how to convert back and forth.

### Pure functions

Pure functions are functions that take input parameters and that compute an output without changing any global state and without mutating the input parameters. They will always return the same output for the same input.

#### Without Errors

If your pure function does not return an error, the idiomatic signature is just fine and no changes are required.

#### With Errors

If your pure function can return an error, then it will have a `(T, error)` return value in idiomatic go. In functional style the return value is [Either[error, T]](https://pkg.go.dev/github.com/IBM/fp-go/either) because function composition is easier with such a return type. Use the `EitherizeXXX` methods in ["github.com/IBM/fp-go/either"](https://pkg.go.dev/github.com/IBM/fp-go/either) to convert from idiomatic to functional style and `UneitherizeXXX` to convert from functional to idiomatic style.

### Effectful functions

An effectful function (or function with a side effect) is one that changes data outside the scope of the function or that does not always produce the same output for the same input (because it depends on some external, mutable state). There is no special way in idiomatic go to identify such a function other than documentation. In functional style we represent them as functions that do not take an input but that produce an output. The base type for these functions is [IO[T]](https://pkg.go.dev/github.com/IBM/fp-go/io) because in many cases such functions represent `I/O` operations.

#### Without Errors

If your effectful function does not return an error, the functional signature is [IO[T]](https://pkg.go.dev/github.com/IBM/fp-go/io)

#### With Errors

If your effectful function can return an error, the functional signature is [IOEither[error, T]](https://pkg.go.dev/github.com/IBM/fp-go/ioeither). Use `EitherizeXXX` from ["github.com/IBM/fp-go/ioeither"](https://pkg.go.dev/github.com/IBM/fp-go/ioeither) to convert an idiomatic go function to functional style.

### Go Context

Functions that take a [context](https://pkg.go.dev/context) are per definition effectful because they depend on the context parameter that is designed to be mutable (it can e.g. be used to cancel a running operation). Furthermore in idiomatic go the parameter is typically passed as the first parameter to a function.

In functional style we isolate the [context](https://pkg.go.dev/context) and represent the nature of the effectful function as an [IOEither[error, T]](https://pkg.go.dev/github.com/IBM/fp-go/ioeither). The resulting type is [ReaderIOEither[T]](https://pkg.go.dev/github.com/IBM/fp-go/context/readerioeither), a function taking a [context](https://pkg.go.dev/context) that returns a function without parameters returning an [Either[error, T]](https://pkg.go.dev/github.com/IBM/fp-go/either). Use the `EitherizeXXX` methods from ["github.com/IBM/fp-go/context/readerioeither"](https://pkg.go.dev/github.com/IBM/fp-go/context/readerioeither) to convert an idiomatic go function with a [context](https://pkg.go.dev/context) to functional style.

## Implementation Notes

### Generics

All monadic operations are implemented via generics, i.e. they offer a type safe way to compose operations. This allows for convenient IDE support and also gives confidence about the correctness of the composition at compile time.

Downside is that this will result in different versions of each operation per type, these versions are generated by the golang compiler at build time (unlike type erasure in languages such as Java of TypeScript). This might lead to large binaries for codebases with many different types. If this is a concern, you can always implement type erasure on top, i.e. use the monadic operations with the `any` type as if generics were not supported. You loose type safety, but this might result in smaller binaries.

### Ordering of Generic Type Parameters

In go we need to specify all type parameters of a function on the global function definition, even if the function returns a higher order function and some of the type parameters are only applicable to the higher order function. So the following is not possible:

```go
func Map[A, B any](f func(A) B) [R, E any]func(fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B]
```

Note that the parameters `R` and `E` are not needed by the first level of `Map` but only by the resulting higher order function. Instead we need to specify the following:

```go
func Map[R, E, A, B any](f func(A) B) func(fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B]
```

which overspecifies `Map` on the global scope. As a result the go compiler will not be able to auto-detect these parameters, it can only auto detect `A` and `B` since they appear in the argument of `Map`. We need to explicitly pass values for these type parameters when `Map` is being used.

Because of this limitation the order of parameters on a function matters. We want to make sure that we define those parameters that cannot be auto-detected, first, and the parameters that can be auto-detected, last. This can lead to inconsistencies in parameter ordering, but we believe that the gain in convenience is worth it. The parameter order of `Ap` is e.g. different from that of `Map`:

```go
func Ap[B, R, E, A any](fa ReaderIOEither[R, E, A]) func(fab ReaderIOEither[R, E, func(A) B]) ReaderIOEither[R, E, B]
```

because `R`, `E` and `A` can be determined from the argument to `Ap` but `B` cannot.

### Use of the [~ Operator](https://go.googlesource.com/proposal/+/master/design/47781-parameterized-go-ast.md)

The FP library attempts to be easy to consume and one aspect of this is the definition of higher level type definitions instead of having to use their low level equivalent. It is e.g. more convenient and readable to use

```go
ReaderIOEither[R, E, A]
```

than

```go
func(R) func() Either.Either[E, A]
```

although both are logically equivalent. At the time of this writing the go type system does not support generic type aliases, only generic type definition, i.e. it is not possible to write:

```go
type ReaderIOEither[R, E, A any] = RD.Reader[R, IOE.IOEither[E, A]]
```

only

```go
type ReaderIOEither[R, E, A any] RD.Reader[R, IOE.IOEither[E, A]]
```

This makes a big difference, because in the second case the type `ReaderIOEither[R, E, A any]` is considered a completely new type, not compatible to its right hand side, so it's not just a shortcut but a fully new type.

From the implementation perspective however there is no reason to restrict the implementation to the new type, it can be generic for all compatible types. The way to express this in go is the [~](https://go.googlesource.com/proposal/+/master/design/47781-parameterized-go-ast.md) operator. This comes with some quite complicated type declarations in some cases, which undermines the goal of the library to be easy to use.

For that reason there exist sub-packages called `Generic` for all higher level types. These packages contain the fully generic implementation of the operations, preferring abstraction over usability. These packages are not meant to be used by end-users but are meant to be used by library extensions. The implementation for the convenient higher level types specializes the generic implementation for the particular higher level type, i.e. this layer does not contain any business logic but only *type magic*.

### Higher Kinded Types

Go does not support higher kinded types (HKT). Such types occur if a generic type itself is parametrized by another generic type. Example:

The `Map` operation for `ReaderIOEither` is defined as:

```go
func Map[R, E, A, B any](f func(A) B) func(fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B]
```

and in fact the equivalent operations for all other mondas follow the same pattern, we could try to introduce a new type for `ReaderIOEither` (without a parameter) as a HKT, e.g. like so (made-up syntax, does not work in go):

```go
func Map[HKT, R, E, A, B any](f func(A) B) func(HKT[R, E, A]) HKT[R, E, B]
```

this would be the completely generic method signature for all possible monads. In particular in many cases it is possible to compose functions independent of the concrete knowledge of the actual `HKT`. From the perspective of a library this is the ideal situation because then a particular algorithm only has to be implemented and tested once.

This FP library addresses this by introducing the HKTs as individual types, e.g. `HKT[A]` would be represented as a new generic type `HKTA`. This loses the correlation to the type `A` but allows to implement generic algorithms, at the price of readability.

For that reason these implementations are kept in the `internal` package. These are meant to be used by the library itself or by extensions, not by end users.
