# Pattern Matching in fp-go

This guide shows how to express multi-branch `switch`/`case` logic as composable
values in fp-go v2. The core idea: **a case is a function `func(T) Option[B]`**
that returns `Some` if it matches and `None` if it does not. A match is an
ordered list of cases, combined so that the first `Some` wins.

Cases and matches are written **point-free**: they are composed from guards,
getters and formatters with `Flow`/`Pipe`, rather than written as lambdas that
name their input.

All code in this guide is covered by runnable examples and tests:

- [`array/example_pattern_matching_test.go`](../array/example_pattern_matching_test.go): `Example_pattern_matching*` (shown on pkg.go.dev)
- [`array/pattern_matching_test.go`](../array/pattern_matching_test.go): unit tests that pin down the semantics described here (evaluation order, laziness, priority, defaults)

## Table of Contents

- [Quick Start](#quick-start)
- [The Case Type](#the-case-type)
- [Point-Free Building Blocks](#point-free-building-blocks)
- [Building Cases](#building-cases)
- [Combining Cases](#combining-cases)
- [Eager vs. Lazy Evaluation](#eager-vs-lazy-evaluation)
- [Closing the Match: Defaults](#closing-the-match-defaults)
- [Two-Way Matching with Fold](#two-way-matching-with-fold)
- [Matching over Slices](#matching-over-slices)
- [Matching with Errors](#matching-with-errors)
- [Recipes](#recipes)
- [Best Practices and Pitfalls](#best-practices-and-pitfalls)
- [Comparison with switch](#comparison-with-switch)
- [API Reference](#api-reference)

## Quick Start

```go
import (
    A  "github.com/IBM/fp-go/v2/array"
    F  "github.com/IBM/fp-go/v2/function"
    O  "github.com/IBM/fp-go/v2/option"
    P  "github.com/IBM/fp-go/v2/predicate"
    RO "github.com/IBM/fp-go/v2/readeroption"
    S  "github.com/IBM/fp-go/v2/string"
)

// leaf accessors (or lens.Get)
func getMethod(r Request) string { return r.Method }
func getPath(r Request) string   { return r.Path }

// guard -> branch result
func on(guard P.Predicate[Request], branch func(Request) string) RO.ReaderOption[Request, string] {
    return F.Flow2(O.FromPredicate(guard), O.Map(branch))
}

func isMethod(m string) P.Predicate[Request] {
    return F.Flow2(getMethod, P.IsStrictEqual[string]()(m))
}

func HandleRequest() func(Request) string {
    return F.Pipe2(
        A.From(                                                          // ordered cases
            on(isMethod("GET"), F.Flow2(getPath, S.Format[string]("Fetching: %s"))),
            on(isMethod("DELETE"), F.Flow2(getPath, S.Format[string]("Deleting: %s"))),
        ),
        A.Fold(RO.AltMonoid[Request, string]()),                         // first Some wins, lazily
        RO.GetOrElse(F.Flow2(getMethod, S.Format[string]("Unsupported method: %s"))), // default
    )
}

handle := HandleRequest()                          // build once, call many times
handle(Request{Method: "GET", Path: "/users"})     // "Fetching: /users"
handle(Request{Method: "PATCH"})                   // "Unsupported method: PATCH"
```

`handle` has type `func(Request) string`. It tries the cases in order,
**stops at the first match** (later cases are never called), and falls back to
the default.

## The Case Type

A case has type `func(T) Option[B]`. fp-go has several aliases for that shape,
so values from different packages can be used as cases directly:

| Alias | Package | Typical source |
|---|---|---|
| `readeroption.ReaderOption[T, B]` | `readeroption` | the type `RO.AltMonoid` combines |
| `option.Kleisli[T, B]` | `option` | `O.FromPredicate`, `O.FromValidation`, `O.InstanceOf` |
| `Prism[T, B].GetOption` | `optics/prism` | `prism.ParseInt()`, `prism.InstanceOf[T]()`, generated prisms |

Because they are all aliases of `func(T) Option[B]`, you can put them in one
slice without conversion.

## Point-Free Building Blocks

These combinators remove the need for lambdas in cases:

| Instead of | Write |
|---|---|
| `func(r Request) bool { return r.Method == m }` | `F.Flow2(getMethod, P.IsStrictEqual[string]()(m))` |
| `func(e Event) bool { return e.Priority > 8 }` | `F.Flow2(getPriority, N.MoreThan(8))` |
| `func(x T) bool { return p(x) && q(x) }` | `F.Pipe1(p, P.And(q))` (also `P.Or`, `P.Not`) |
| `func(s string) bool { return strings.HasPrefix(s, p) }` | `F.Bind2nd(strings.HasPrefix, p)` |
| `func(s string) string { return strings.TrimPrefix(s, p) }` | `F.Bind2nd(strings.TrimPrefix, p)` |
| `func(x T) string { return fmt.Sprintf("...%v", x) }` | `S.Format[T]("...%v")` |
| `func(s string) string { return "INFO: " + s }` | `S.Prepend("INFO: ")` |
| `func(T) B { return b }` | `F.Constant1[T](b)` |
| `func(s string) (int64, error) { return strconv.ParseInt(s, base, 64) }` | `F.Bind23of3(result.Eitherize3(strconv.ParseInt))(base, 64)` |
| `func() Option[B] { return c(x) }` | `F.Nullary2(F.Constant(x), c)` |

Leaf functions are the exception: field accessors (or `lens.Get`) and formatters
that combine several fields are ordinary named functions, which the pipeline
then references by name.

## Building Cases

### Guard + result

`O.FromPredicate` turns a guard into `Some(x)`/`None`, and `O.Map` computes the
branch result. Most cases have this shape, so give it a name:

```go
func when[T, B any](guard P.Predicate[T], b B) RO.ReaderOption[T, B] {
    return F.Flow2(O.FromPredicate(guard), O.Map(F.Constant1[T](b)))
}

func Classify() func(int) string {
    return F.Pipe2(
        A.From(
            when(P.IsZero[int](), "zero"),
            when(N.LessThan(0), "negative"),
            when(N.LessThan(11), "small positive"),
        ),
        A.Fold(RO.AltMonoid[int, string]()),
        RO.GetOrElse(F.Constant1[int]("large positive")),
    )
}
```

Compose guards with `P.And`, `P.Or` and `P.Not`:

```go
isUrgent := F.Flow2(getPriority, N.MoreThan(8))
isCritical := F.Pipe1(isType("error"), P.And(isUrgent))
```

### Guard + transformation

A case can also transform its input. For example, match a prefix, strip it,
and parse the rest:

```go
func parseBase(base int) result.Kleisli[string, int64] {
    return F.Bind23of3(result.Eitherize3(strconv.ParseInt))(base, 64)
}

func withPrefix(prefix string, base int) RO.ReaderOption[string, int64] {
    return F.Flow3(
        O.FromPredicate(F.Bind2nd(strings.HasPrefix, prefix)),
        O.Map(F.Bind2nd(strings.TrimPrefix, prefix)),
        O.Chain(F.Flow2(parseBase(base), result.ToOption[int64])),
    )
}
```

### From `(value, ok)` functions

`O.FromValidation` lifts any Go function of shape `func(A) (B, bool)` into a
case without a wrapper:

```go
lookup := A.Fold(RO.AltMonoid[string, string]())(A.From(
    O.FromValidation(os.LookupEnv),
    RO.Of[string]("default"),
))
```

### Type cases (the functional type switch)

`O.InstanceOf[T]` is a type assertion that returns an `Option`:

```go
area := F.Pipe2(
    A.From(
        F.Flow2(O.InstanceOf[Circle], O.Map(circleArea)),
        F.Flow2(O.InstanceOf[Rect], O.Map(rectArea)),
    ),
    A.Fold(RO.AltMonoid[any, float64]()),
    RO.GetOrElse(F.Constant1[any](-1.0)),
)
```

`O.InstanceOf` takes `any`. When the input is a named interface such as
`Shape`, prepend `F.ToAny` so the case has type `func(Shape) Option[B]`:

```go
func caseOf[T Shape, B any](f func(T) B) RO.ReaderOption[Shape, B] {
    return F.Flow3(F.ToAny[Shape], O.InstanceOf[T], O.Map(f))
}

describe := A.Fold(RO.AltMonoid[Shape, string]())(A.From(
    caseOf(S.Format[Circle]("circle %v")),
    caseOf(S.Format[Rect]("rect %v")),
))
```

### Prisms

A prism's `GetOption` is a case. The built-in prisms (`ParseInt`, `ParseInt64`,
`ParseFloat64`, `ParseBool`, `ParseDate`, `RegexMatcher`, `InstanceOf`, `Deref`,
`FromOption`, ...) and generated prisms for sum types can be used directly:

```go
classify := A.Fold(RO.AltMonoid[string, string]())(A.From(
    F.Flow2(prism.ParseInt().GetOption, O.Map(S.Format[int]("int %d"))),
    F.Flow2(prism.ParseFloat64().GetOption, O.Map(S.Format[float64]("float %g"))),
    F.Flow2(prism.ParseBool().GetOption, O.Map(S.Format[bool]("bool %t"))),
))
classify("42")   // Some("int 42"), because the int case comes first
classify("3.5")  // Some("float 3.5")
classify("nope") // None
```

### Catch-all case

`RO.Of[T](b)` always matches. Put it last to make the match total (see
[Defaults](#closing-the-match-defaults)).

### Hand-written cases

If a guard and its result depend on the same intermediate value and no
combinator expresses that cleanly, a plain `func(T) Option[B]` is fine. Keep it
small and named, and treat it as a leaf:

```go
func matchSession(r Request) O.Option[Session] {
    if s, ok := sessions[r.Cookie]; ok && s.Valid() {
        return O.Some(s)
    }
    return O.None[Session]()
}
```

## Combining Cases

Given ordered cases `c1, c2, ..., cn` of type `func(T) Option[B]`, there are
four ways to get "the first case that matches":

| Technique | Lazy? | When to use |
|---|---|---|
| `A.Fold(RO.AltMonoid[T, B]())(cases)` | yes | **Default choice.** Builds a reusable `func(T) Option[B]` once. |
| `F.Flow3(R.Read[O.Option[B], T], A.FindFirstMap[RO.ReaderOption[T, B], B], R.Read[O.Option[B]](cases))` | yes | Same result via `FindFirstMap`, useful when the cases are data. |
| `F.Pipe2(c1, RO.Alt(F.Constant(c2)), RO.Alt(F.Constant(c3)))` | yes | Two or three fixed alternatives. |
| `O.AltAllArray(O.None[B]())(options)` | **no** | The `Option` values already exist (for example, results computed earlier). |

### `A.Fold` with `RO.AltMonoid`

`RO.AltMonoid[T, B]()` is a monoid on `func(T) Option[B]`:

- `Empty()` is the case that never matches (`None` for every input).
- `Concat(c1, c2)` is the case "try `c1`, and only if it returns `None`, try `c2`".

Folding a slice of cases with it gives a single case that tries them in order.
The result is an ordinary `ReaderOption`, so it can itself be used as a case in
another match. Nesting works without extra glue.

```go
match := A.Fold(RO.AltMonoid[T, B]())(A.From(c1, c2, c3))
```

A reusable generic helper, if you need it often:

```go
func Match[T, B any](cases ...RO.ReaderOption[T, B]) RO.ReaderOption[T, B] {
    return A.Fold(RO.AltMonoid[T, B]())(cases)
}
```

### `A.FindFirstMap` with `reader.Read`

`FindFirstMap(sel)` applies `sel` to the elements of a slice and returns the
first `Some`. When the slice holds cases, `reader.Read[O.Option[B]](x)` (which
applies a case to `x`) is the selector. Composing the steps gives the matcher
point-free:

```go
route := F.Flow4(
    R.Read[O.Option[string], string],                        // path -> "apply a case to path"
    A.FindFirstMap[RO.ReaderOption[string, string], string], // -> first case that matches
    R.Read[O.Option[string]](cases),                         // run it over the cases
    O.GetOrElse(F.Constant("page")),                         // default
)
```

### `RO.Alt`

`RO.Alt` appends one fallback case. It takes the fallback as a thunk, and
`F.Constant(c2)` provides one:

```go
resolve := F.Pipe2(
    fromEnv,
    RO.Alt(F.Constant(fromFile)),
    RO.GetOrElse(S.Prepend("literal:")),
)
```

On plain `Option` values, `O.Alt` does the same:
`O.Alt(F.Nullary2(F.Constant(x), c2))` runs `c2(x)` only if needed.

### `O.AltAllArray`

`O.AltAllArray(startWith)(options)` returns `startWith` if it is `Some`,
otherwise the first `Some` in `options`, otherwise `startWith`. Use
`O.None[B]()` as `startWith` for "no match". `O.AltAllSeq` is the same for
`iter.Seq`.

Use it when the `Option` values already exist. **Do not** use it to combine
cases. See the next section for why.

## Eager vs. Lazy Evaluation

In Go, function arguments and slice literals are evaluated before the call. In
this code

```go
O.AltAllArray(O.None[string]())([]O.Option[string]{
    matchA(x), matchB(x), matchC(x),
})
```

**all three cases run**, even if `matchA` already matched. `AltAllArray` only
stops *scanning* early. It cannot stop the calls that produced the slice. That
matters when cases are expensive (parsing, lookups) or have side effects.

The point-free spelling makes the problem visible: to feed cases to
`AltAllArray` you must first run all of them, for example with
`R.SequenceArray`:

```go
// eager: SequenceArray runs every case, then AltAllArray picks the first Some
F.Pipe2(cases, R.SequenceArray[T, O.Option[B]], R.Map[T](O.AltAllArray(O.None[B]())))
```

The lazy techniques pass the *cases* (functions), not their *results*, so the
combinator decides what to call:

```go
match := A.Fold(RO.AltMonoid[int, string]())(A.From(
    trace("a", O.None[string]()),
    trace("b", O.Some("b")),
    trace("c", O.Some("c")),
))
match(0)
// try a
// try b          <- "c" is never invoked
// Some[string](b)
```

Guaranteed by tests: `TestPatternMatching_AltAllArrayIsEagerButScanStopsEarly`,
`TestPatternMatching_FoldAltMonoidIsLazy`, `TestPatternMatching_FindFirstMapIsLazy`,
`TestPatternMatching_ReaderOptionAltIsLazy`, `TestPatternMatching_AltChainIsLazy`.

Building the folded matcher does not run any case. Build it once (at startup
or when a component is constructed) and reuse it. Every call then costs at most
one invocation per case.

## Closing the Match: Defaults

A folded match is partial and returns `Option[B]`. There are three ways to
handle the "no match" branch:

```go
// 1. Default that sees the input -> func(T) B
F.Pipe1(match, RO.GetOrElse(F.Flow2(getMethod, S.Format[string]("Unsupported method: %s"))))

// 2. Constant default -> func(T) B
F.Pipe1(match, RO.GetOrElse(F.Constant1[T](defaultB)))

// 3. Catch-all as the last case -> still func(T) Option[B], but always Some
A.Fold(RO.AltMonoid[T, B]())(A.From(c1, c2, RO.Of[T](defaultB)))
```

Prefer (1) or (2) when the result should be total. They remove the `Option`
from the type, so callers cannot forget the default.

Go does not check exhaustiveness. Unlike languages with algebraic data types,
the compiler does not warn when a variant of a sum type has no case. A default
branch guarantees a *result*, not that every variant was considered. When
matching on sum types, add a test per variant (see
`TestPatternMatching_SumTypeCases`).

With no cases at all, the fold returns `Empty()`, which is `None` for every input.

## Two-Way Matching with Fold

For exactly two branches on a sum type you already have, use the type's
eliminator instead of a case list:

| Input | Eliminator | Branches |
|---|---|---|
| `Option[A]` | `O.Fold(onNone func() B, onSome func(A) B)` | none / some |
| `Either[E, A]` | `either.Fold(onLeft func(E) B, onRight func(A) B)` | left / right |
| `Result[A]` | `result.Fold(onErr func(error) B, onOk func(A) B)` | error / ok |
| `bool` via predicate | `predicate.Fold(onFalse, onTrue func(A) B)(p)` | false / true |
| `[]A` | `A.Match(onEmpty func() B, onNonEmpty func([]A) B)` | empty / non-empty |
| `[]A` | `A.MatchLeft(onEmpty func() B, onNonEmpty func(A, []A) B)` | empty / head + tail |

```go
describe := O.Fold(F.Constant("nothing"), S.Format[int]("got %d"))
describe(O.Some(3)) // "got 3"
```

## Matching over Slices

The same cases work on collections:

| Function | Result |
|---|---|
| `A.FindFirstMap(case)` | first element that matches, mapped |
| `A.FindLastMap(case)` | last element that matches, mapped |
| `A.FindFirstMapWithIndex(func(int, A) Option[B])` | as `FindFirstMap`, with index |
| `A.FilterMap(case)` | all matches, mapped (`[]B`) |

```go
parseNumber := A.Fold(RO.AltMonoid[string, int64]())(A.From(
    withPrefix("0x", 16),
    withPrefix("0b", 2),
    withPrefix("0o", 8),
    prism.ParseInt64().GetOption,
))

inputs := A.From("invalid", "0x2A", "also bad", "17")
A.FindFirstMap(parseNumber)(inputs)                            // Some(42)
A.FindLastMap(parseNumber)(inputs)                             // Some(17)
A.FilterMap(parseNumber)(A.From("0b101", "x", "0o17", "7"))    // [5 15 7]
```

`FindFirstMap` stops at the first match. The selector is not called for later
elements.

## Matching with Errors

If a failing case should explain why it failed, use `Result`/`Either` instead
of `Option`. The combinators are the same and differ only in the empty value:

```go
import RR "github.com/IBM/fp-go/v2/readerresult"

// withPrefixR fails with err if the prefix is missing or the rest does not parse
func withPrefixR(prefix string, base int, err error) RR.ReaderResult[string, int64] {
    return F.Flow4(
        result.FromPredicate(F.Bind2nd(strings.HasPrefix, prefix), F.Constant1[string](err)),
        result.Map(F.Bind2nd(strings.TrimPrefix, prefix)),
        result.Chain(parseBase(base)),
        result.MapLeft[int64](F.Constant1[error](err)),
    )
}

func decimalR(err error) RR.ReaderResult[string, int64] {
    return F.Flow2(parseBase(10), result.MapLeft[int64](F.Constant1[error](err)))
}

parse := A.Fold(RR.AltMonoid(F.Constant(RR.Left[string, int64](errNoMatch))))(A.From(
    withPrefixR("0x", 16, errHex),
    decimalR(errDec),
))

parse("0xff") // Ok(255), decimalR is not called
parse("x")    // Err(errDec): the error of the last failing case
```

- The fold is lazy, like `RO.AltMonoid`.
- If every case fails, the **last** case's error is returned. With no cases the
  `zero` value is returned.
- For two or three fixed alternatives: `F.Pipe1(c1, RR.Alt(F.Constant(c2)))`.
- `either.AltAllArray` / `result.AltAllArray` exist and are eager, like `O.AltAllArray`.
- `readereither.AltMonoid` is the equivalent for a custom error type.

## Recipes

### Route table

```go
func prefix(p string, h Handler) RO.ReaderOption[Request, Handler] {
    return F.Flow2(
        O.FromPredicate(F.Flow2(getPath, F.Bind2nd(strings.HasPrefix, p))),
        O.Map(F.Constant1[Request](h)),
    )
}

func Router() func(Request) Handler {
    return F.Pipe2(
        A.From(prefix("/api/", apiHandler), prefix("/static/", assetHandler)),
        A.Fold(RO.AltMonoid[Request, Handler]()),
        RO.GetOrElse(F.Constant1[Request](notFound)),
    )
}
```

### Nested matches

A folded match is itself a case, so matches nest:

```go
errorCases := Match(matchCritical, matchError)
match      := Match(errorCases, matchWarning, matchInfo)
```

### Extending a match

Cases are plain slices, so a specialized matcher can prepend cases to a shared base:

```go
verbose := F.Pipe2(
    A.From(matchError, matchWarning),
    A.Prepend(matchDebug),
    A.Fold(RO.AltMonoid[Event, string]()),
)
```

## Best Practices and Pitfalls

1. **Order from specific to general.** The first match wins, so a general case
   placed first shadows the specific ones after it
   (`TestPatternMatching_OrderDefinesPriority`).
2. **Pass cases, not results.** Combine `func(T) Option[B]` values with
   `A.Fold(RO.AltMonoid...)`, `A.FindFirstMap` or `RO.Alt`. Avoid
   `AltAllArray([]Option{c1(x), c2(x)})`, because it runs every case.
3. **Stay point-free.** Build guards from getters and predicate combinators,
   and branch results from `S.Format`, `S.Prepend`, `F.Constant1` or named
   functions. Keep lambdas for leaves only.
4. **Build the matcher once.** Wrap it in a constructor function
   (`func Router() func(Request) Handler`), call that once during setup, and
   reuse the result. Do not rebuild it inside a hot loop.
5. **Close with `RO.GetOrElse`** when a total function is intended, so the
   `Option` does not leak to callers.
6. **Keep cases pure.** Only a lazy combinator guarantees that side effects of
   later cases do not happen. Do not rely on that for correctness. If a case
   must perform IO, match on the input and return the effect (for example,
   `func(T) Option[IOResult[B]]`) instead of running it inside the guard.
7. **Test cases individually and the match as a whole.** Cases are ordinary
   functions: `assert.Equal(t, O.None[string](), matchGET(Request{Method: "POST"}))`.
8. **Name cases after what they match** (`isUrgent`, `withPrefix("0x", 16)`),
   not by position (`case1`).

## Comparison with switch

```go
// switch
func sign(n int) string {
    switch {
    case n > 0:
        return "positive"
    case n < 0:
        return "negative"
    default:
        return "zero"
    }
}

// fp-go
func Sign() func(int) string {
    return F.Pipe2(
        A.From(
            when(N.MoreThan(0), "positive"),
            when(N.LessThan(0), "negative"),
        ),
        A.Fold(RO.AltMonoid[int, string]()),
        RO.GetOrElse(F.Constant1[int]("zero")),
    )
}
```

Use `switch` when the branches are few, fixed, local to one function and not
reused. It is simpler and faster.

Use case lists when:

- cases are reused across several matchers or extended per context;
- the set of cases is data (configured, registered, or built at runtime);
- cases already exist as `Option`-returning functions (prisms, parsers, `FromPredicate` guards);
- the match itself must be a value to pass around, nest, or compose in a pipeline;
- each branch should be unit-tested in isolation.

Both approaches share one limitation: Go does not check exhaustiveness.

## API Reference

| Function | Signature (simplified) | Role |
|---|---|---|
| `readeroption.AltMonoid[T, B]()` | `Monoid[func(T) Option[B]]` | first-match combinator for cases (lazy) |
| `readeroption.Alt(thunk)` | `func(func(T) Option[B]) func(T) Option[B]` | appends one fallback case (lazy) |
| `array.Fold(m)` | `func([]A) A` | concatenates cases with a monoid |
| `readeroption.GetOrElse(onNone)` | `func(func(T) Option[B]) func(T) B` | default branch |
| `array.FindFirstMap(sel)` | `func([]A) Option[B]` | first element where `sel` is `Some` (lazy) |
| `array.FindLastMap(sel)` | `func([]A) Option[B]` | last element where `sel` is `Some` |
| `array.FilterMap(sel)` | `func([]A) []B` | all elements where `sel` is `Some` |
| `reader.Read[B](x)` | `func(func(T) B) B` | applies a case to `x` |
| `reader.SequenceArray` | `func([]func(T) A) func(T) []A` | runs all cases (eager) |
| `option.Alt(thunk)` | `func(Option[A]) Option[A]` | lazy binary fallback on values |
| `option.AltAllArray(start)` | `func([]Option[A]) Option[A]` | first `Some` among existing Options (eager input) |
| `option.AltAllSeq(start)` | `func(Seq[Option[A]]) Option[A]` | same for iterators |
| `option.AltMonoid[A]()` | `Monoid[Option[A]]` | first-`Some` monoid on values |
| `option.FromPredicate(p)` | `func(A) Option[A]` | guard, turns a predicate into a case |
| `option.FromValidation(f)` | `func(A) Option[B]` | lifts `func(A) (B, bool)` |
| `option.InstanceOf[T]` | `func(any) Option[T]` | type case |
| `function.ToAny[T]` | `func(T) any` | adapts a named interface for `InstanceOf` |
| `prism.Prism[S, A].GetOption` | `func(S) Option[A]` | prism as case |
| `predicate.IsStrictEqual`, `And`, `Or`, `Not`, `ContraMap` | | guard combinators |
| `function.Bind2nd`, `Bind23of3`, `Constant1`, `Nullary2` | | adapt Go functions point-free |
| `string.Format`, `Prepend`, `Append` | | branch results |
| `readerresult.AltMonoid(zero)` | `Monoid[func(T) Result[B]]` | first-match that keeps the error (lazy) |
| `option.Fold`, `either.Fold`, `array.Match` | | two-way eliminators |
