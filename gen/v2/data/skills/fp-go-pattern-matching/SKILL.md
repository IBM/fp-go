---
name: fp-go-pattern-matching
description: >
  Guides writing, refactoring, and reviewing fp-go v2 code that replaces
  switch/case, if-else chains, or type switches with composable, point-free
  pattern matching. Apply whenever the user wants multi-branch dispatch as
  values: first-match-wins case lists, routers, classifiers, multi-format
  parsers, fallback chains, or matching over sum types and slices. Trigger on
  any mention of: pattern matching in Go, match/case, switch to functional,
  AltAllArray, AltAllSeq, AltMonoid, readeroption.AltMonoid, readeroption.Alt,
  FindFirstMap, FindLastMap, FilterMap, option.Alt, first Some, fallback chain,
  InstanceOf, prism GetOption as matcher, GetOrElse default branch, option.Fold,
  either.Fold, array.Match.
---

# fp-go Pattern Matching

All imports come from `github.com/IBM/fp-go/v2` and use the canonical aliases of the **fp-go** skill. The canonical long-form
reference is `v2/docs/PATTERN_MATCHING.md`. Tested examples are in
`v2/array/example_pattern_matching_test.go` and `v2/array/pattern_matching_test.go`.

```go
import (
    A  "github.com/IBM/fp-go/v2/array"
    F  "github.com/IBM/fp-go/v2/function"
    N  "github.com/IBM/fp-go/v2/number"
    O  "github.com/IBM/fp-go/v2/option"
    P  "github.com/IBM/fp-go/v2/predicate"
    R  "github.com/IBM/fp-go/v2/result"
    RD "github.com/IBM/fp-go/v2/reader"
    RO "github.com/IBM/fp-go/v2/readeroption"
    S  "github.com/IBM/fp-go/v2/string"
    "github.com/IBM/fp-go/v2/optics/prism"
    "github.com/IBM/fp-go/v2/readerresult" // generic ReaderResult[T, B]; not context/readerresult
)
```

---

## The model

- A **case** is `func(T) Option[B]`: `Some(b)` if it matches, `None` otherwise.
  `RO.ReaderOption[T, B]`, `O.Kleisli[T, B]` and `prism.Prism[T, B].GetOption`
  are all aliases of this type, so values from these packages mix freely in one slice.
- A **match** is an ordered slice of cases, and the first `Some` wins.
- A **default** turns the partial match (`Option[B]`) into a total function.
- **Everything is point-free.** Cases are composed from guards, getters and
  formatters with `Flow`/`Pipe`. Lambdas are allowed only for leaves (field
  accessors, multi-field formatters, deliberate side effects in tests).

## The canonical idiom (use this by default)

```go
func getMethod(r Request) string { return r.Method }   // leaf accessor (or lens.Get)
func getPath(r Request) string   { return r.Path }

func on(guard P.Predicate[Request], branch func(Request) string) RO.ReaderOption[Request, string] {
    return F.Flow2(O.FromPredicate(guard), O.Map(branch))
}

func isMethod(m string) P.Predicate[Request] {
    return F.Flow2(getMethod, P.IsStrictEqual[string]()(m))
}

func HandleRequest() func(Request) string {
    return F.Pipe2(
        A.From(                                           // ordered, specific -> general
            on(isMethod("GET"), F.Flow2(getPath, S.Format[string]("Fetching: %s"))),
            on(isMethod("DELETE"), F.Flow2(getPath, S.Format[string]("Deleting: %s"))),
        ),
        A.Fold(RO.AltMonoid[Request, string]()),          // lazy first-match
        RO.GetOrElse(F.Flow2(getMethod, S.Format[string]("Unsupported: %s"))), // default
    )
}
```

- **Lazy:** stops at the first match. Later cases are never invoked.
- Building the matcher runs no case. Wrap it in a constructor `func` (never a
  package-level `var`), call the constructor once during setup, and reuse the result.
- The folded match is itself a case, so matches nest.
- Zero cases give `Empty()`, which is `None` for all inputs.

Generic helper for the most common case shape:

```go
func when[T, B any](guard P.Predicate[T], b B) RO.ReaderOption[T, B] {
    return F.Flow2(O.FromPredicate(guard), O.Map(F.Constant1[T](b)))
}
```

Alternatives, all lazy and point-free:

```go
// append one or two fallback cases
F.Pipe2(c1, RO.Alt(F.Constant(c2)), RO.Alt(F.Constant(c3)))

// FindFirstMap over a slice of cases
F.Flow3(
    RD.Read[O.Option[B], T],                   // x -> apply-a-case-to-x
    A.FindFirstMap[RO.ReaderOption[T, B], B], // -> first case that matches
    RD.Read[O.Option[B]](cases),               // run over the cases
)
```

## Point-free building blocks

| Instead of | Write |
|---|---|
| `func(r Req) bool { return r.Method == m }` | `F.Flow2(getMethod, P.IsStrictEqual[string]()(m))` |
| `func(e Ev) bool { return e.Priority > 8 }` | `F.Flow2(getPriority, N.MoreThan(8))` |
| `p(x) && q(x)` | `F.Pipe1(p, P.And(q))`, also `P.Or`, `P.Not` |
| `strings.HasPrefix(s, p)` | `F.Bind2nd(strings.HasPrefix, p)` |
| `strings.TrimPrefix(s, p)` | `F.Bind2nd(strings.TrimPrefix, p)` |
| `fmt.Sprintf("x %v", v)` | `S.Format[T]("x %v")` |
| `"INFO: " + s` | `S.Prepend("INFO: ")`, `S.Append` |
| `func(T) B { return b }` | `F.Constant1[T](b)` |
| `strconv.ParseInt(s, base, 64)` | `F.Bind23of3(R.Eitherize3(strconv.ParseInt))(base, 64)` |
| `strconv.Atoi(s)` | `R.Eitherize1(strconv.Atoi)` |
| `Result` -> `Option` | `R.ToOption[A]` |
| `func() Option[B] { return c(x) }` | `F.Nullary2(F.Constant(x), c)` |
| `func(s Shape) any { return s }` | `F.ToAny[Shape]` |
| replace any error | `R.MapLeft[A](F.Constant1[error](err))` |

(`number` has **no** `Equal`/`LessThanEqual`. Use `P.IsStrictEqual`, `P.IsZero[T]()`, `N.LessThan(n+1)`.)

## Critical pitfall: AltAllArray over applied cases is eager

```go
// BAD: all three cases run, because Go evaluates the slice literal first
O.AltAllArray(O.None[B]())([]O.Option[B]{c1(x), c2(x), c3(x)})
// its point-free equivalent shows why: SequenceArray runs every case first
F.Pipe2(cases, RD.SequenceArray[T, O.Option[B]], RD.Map[T](O.AltAllArray(O.None[B]())))
```

`AltAllArray` stops *scanning* at the first `Some` but cannot prevent the calls
that built the slice. Use it only when the `Option` values already exist.
When reviewing code, rewrite this form to `A.Fold(RO.AltMonoid[T, B]())(A.From(c1, c2, c3))`.
The same applies to `either.AltAllArray` and `result.AltAllArray`.

## Building cases

| Need | Build with |
|---|---|
| guard -> constant | `when(guard, b)` = `F.Flow2(O.FromPredicate(guard), O.Map(F.Constant1[T](b)))` |
| guard -> computed | `F.Flow2(O.FromPredicate(guard), O.Map(f))` |
| guard -> transform -> partial parse | `F.Flow3(O.FromPredicate(g), O.Map(f), O.Chain(k))` |
| `func(A) (B, bool)` | `O.FromValidation(f)`, for example `O.FromValidation(os.LookupEnv)` |
| type switch on `any` | `F.Flow2(O.InstanceOf[Variant], O.Map(handler))` |
| type switch on named interface | `F.Flow3(F.ToAny[Shape], O.InstanceOf[Variant], O.Map(handler))` |
| parse / extract | `prism.ParseInt().GetOption`, `ParseInt64`, `ParseFloat64`, `ParseBool`, `ParseDate(layout)`, `RegexMatcher(re)` |
| pointer non-nil | `O.FromNillable2` |
| catch-all | `RO.Of[T](b)` as the **last** case |
| no combinator fits | small **named** `func(T) O.Option[B]` leaf |

Prefix-parse case, fully point-free:

```go
func withPrefix(prefix string, base int) RO.ReaderOption[string, int64] {
    return F.Flow3(
        O.FromPredicate(F.Bind2nd(strings.HasPrefix, prefix)),
        O.Map(F.Bind2nd(strings.TrimPrefix, prefix)),
        O.Chain(F.Flow2(F.Bind23of3(R.Eitherize3(strconv.ParseInt))(base, 64), R.ToOption[int64])),
    )
}
```

Sum-type cases via an adapter:

```go
func caseOf[T Shape, B any](f func(T) B) RO.ReaderOption[Shape, B] {
    return F.Flow3(F.ToAny[Shape], O.InstanceOf[T], O.Map(f))
}
describe := A.Fold(RO.AltMonoid[Shape, string]())(A.From(
    caseOf(S.Format[Circle]("circle %v")),
    caseOf(S.Format[Rect]("rect %v")),
))
```

Go has no exhaustiveness checking. A default guarantees a result, not
coverage. Add one test per variant.

## Closing the match

```go
F.Pipe1(match, RO.GetOrElse(F.Flow2(getX, S.Format[X]("no match: %v")))) // default sees input
F.Pipe1(match, RO.GetOrElse(F.Constant1[T](dflt)))                       // constant default
A.From(c1, c2, RO.Of[T](dflt))                                           // catch-all case, still Option
```

## Two-way matches: use the eliminator, not a case list

| Input | Use |
|---|---|
| `Option[A]` | `O.Fold(onNone func() B, onSome func(A) B)`, e.g. `O.Fold(F.Constant("none"), S.Format[int]("got %d"))` |
| `Either[E, A]` | `either.Fold(onLeft, onRight)` |
| `Result[A]` | `R.Fold(onErr func(error) B, onOk func(A) B)` |
| predicate | `P.Fold(onFalse, onTrue func(A) B)(pred)` |
| slice empty/non-empty | `A.Match(onEmpty func() B, onNonEmpty func([]A) B)`, `A.MatchLeft(onEmpty, func(head A, tail []A) B)` |

## Matching over slices

| Goal | Use |
|---|---|
| first element matching any case | `A.FindFirstMap(match)` (lazy over elements) |
| last element | `A.FindLastMap(match)` |
| with index | `A.FindFirstMapWithIndex(func(int, A) O.Option[B])` |
| all matches | `A.FilterMap(match)` -> `[]B` |

## Matching with errors

When a failing case should say why, use `Result` cases and the error-aware monoid:

```go
func withPrefixR(prefix string, base int, err error) readerresult.ReaderResult[string, int64] {
    return F.Flow4(
        R.FromPredicate(F.Bind2nd(strings.HasPrefix, prefix), F.Constant1[string](err)),
        R.Map(F.Bind2nd(strings.TrimPrefix, prefix)),
        R.Chain(F.Bind23of3(R.Eitherize3(strconv.ParseInt))(base, 64)),
        R.MapLeft[int64](F.Constant1[error](err)),
    )
}

parse := A.Fold(readerresult.AltMonoid(F.Constant(readerresult.Left[string, int64](errNoMatch))))(A.From(
    withPrefixR("0x", 16, errHex),
    withPrefixR("", 10, errDec),
))
// two alternatives: F.Pipe1(c1, readerresult.Alt(F.Constant(c2)))
```

- Lazy. If all cases fail, the **last** case's error is returned. With no cases, `zero` is returned.
- If the cases are declared as plain `func(string) R.Result[int64]`,
  `A.From` may need an explicit type argument: `A.From[readerresult.ReaderResult[string, int64]](...)`.
- `readereither.AltMonoid(zero)` is the equivalent for custom error types.

## Review checklist

1. Cases are ordered specific -> general (a general case first shadows later ones).
2. No `AltAllArray([]Option{c(x), ...})` over case applications. Rewrite to `A.Fold(RO.AltMonoid...)`.
3. Point-free: no lambdas that only compare a field, call one stdlib function,
   format a value or return a constant. Replace them using the building-blocks table.
4. The matcher is built by a constructor `func`, once, not inside a loop or per request.
5. Intended total functions end in `RO.GetOrElse`, so `Option` does not leak.
6. Cases are pure. If a branch needs IO, return the effect (`Option[IOResult[B]]`) instead of running it in the guard.
7. Each case has a unit test, plus a test of the combined match, including the default and priority.
8. Fixed, local, non-reused branching stays a plain `switch`. Do not force this pattern where it adds nothing.

## Before presenting code

Signatures are easy to misremember. Check anything not shown here with the
fp-go MCP `search_examples` tool (see **fp-go-mcp**), then run `go build ./...`,
`go vet ./...` and the tests.
