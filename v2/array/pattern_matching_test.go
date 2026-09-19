// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package array_test

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/optics/prism"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	R "github.com/IBM/fp-go/v2/reader"
	RO "github.com/IBM/fp-go/v2/readeroption"
	RR "github.com/IBM/fp-go/v2/readerresult"
	"github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// The tests in this file pin down the semantics documented in
// docs/PATTERN_MATCHING.md. Each matcher ("case") is a function
// func(T) Option[B], which is the same type as RO.ReaderOption[T, B],
// O.Kleisli[T, B] and prism.Prism[T, B].GetOption.

// counting wraps a function and records how often it was invoked.
func counting[T, B any](calls *int, f func(T) B) func(T) B {
	return func(t T) B {
		*calls++
		return f(t)
	}
}

// matchFold is the lazy "first matching case" combinator used throughout the docs.
func matchFold[T, B any](cases ...RO.ReaderOption[T, B]) RO.ReaderOption[T, B] {
	return A.Fold(RO.AltMonoid[T, B]())(cases)
}

// matchFind is the equivalent formulation via FindFirstMap, point-free:
// t -> apply-to-t -> first case that matches -> run over cases.
func matchFind[T, B any](cases ...RO.ReaderOption[T, B]) RO.ReaderOption[T, B] {
	return F.Flow3(
		R.Read[O.Option[B], T],
		A.FindFirstMap[RO.ReaderOption[T, B], B],
		R.Read[O.Option[B]](cases),
	)
}

// matchAltAll runs every case (R.SequenceArray) and then picks the first
// Some. It is the point-free form of AltAllArray([]Option{c1(t), c2(t), ...})
// and, like it, evaluates eagerly.
func matchAltAll[T, B any](cases ...RO.ReaderOption[T, B]) RO.ReaderOption[T, B] {
	return F.Pipe2(
		cases,
		R.SequenceArray[T, O.Option[B]],
		R.Map[T](O.AltAllArray(O.None[B]())),
	)
}

// when builds a case from a guard and a constant result.
func when[T, B any](guard P.Predicate[T], b B) RO.ReaderOption[T, B] {
	return F.Flow2(O.FromPredicate(guard), O.Map(F.Constant1[T](b)))
}

func hasPrefix(prefix string) P.Predicate[string] {
	return F.Bind2nd(strings.HasPrefix, prefix)
}

func trimPrefix(prefix string) func(string) string {
	return F.Bind2nd(strings.TrimPrefix, prefix)
}

// parseBase is strconv.ParseInt with the base and a bit size of 64 bound.
func parseBase(base int) result.Kleisli[string, int64] {
	return F.Bind23of3(result.Eitherize3(strconv.ParseInt))(base, 64)
}

// withPrefix matches strings that start with prefix and parses the remainder.
func withPrefix(prefix string, base int) RO.ReaderOption[string, int64] {
	return F.Flow3(
		O.FromPredicate(hasPrefix(prefix)),
		O.Map(trimPrefix(prefix)),
		O.Chain(F.Flow2(parseBase(base), result.ToOption[int64])),
	)
}

// withPrefixR is withPrefix with a dedicated error for any failure.
func withPrefixR(prefix string, base int, err error) RR.ReaderResult[string, int64] {
	return F.Flow4(
		result.FromPredicate(hasPrefix(prefix), F.Constant1[string](err)),
		result.Map(trimPrefix(prefix)),
		result.Chain(parseBase(base)),
		result.MapLeft[int64](F.Constant1[error](err)),
	)
}

// decimalR parses base 10 with a dedicated error.
func decimalR(err error) RR.ReaderResult[string, int64] {
	return F.Flow2(parseBase(10), result.MapLeft[int64](F.Constant1[error](err)))
}

func sign() func(int) string {
	return F.Pipe2(
		A.From(
			when(N.MoreThan(0), "positive"),
			when(N.LessThan(0), "negative"),
		),
		A.Fold(RO.AltMonoid[int, string]()),
		RO.GetOrElse(F.Constant1[int]("zero")),
	)
}

func TestPatternMatching_AltAllArrayIsEagerButScanStopsEarly(t *testing.T) {
	var c1, c2, c3 int
	m := matchAltAll(
		counting(&c1, RO.None[int, string]()),
		counting(&c2, RO.Of[int]("second")),
		counting(&c3, RO.Of[int]("third")),
	)

	assert.Equal(t, O.Some("second"), m(1))
	assert.Equal(t, 1, c1)
	assert.Equal(t, 1, c2)
	assert.Equal(t, 1, c3, "every case runs before AltAllArray picks the first Some")

	// the same holds for the common slice-literal form
	var d1, d2 int
	first := counting(&d1, RO.Of[int]("first"))
	second := counting(&d2, RO.Of[int]("second"))
	assert.Equal(t, O.Some("first"), O.AltAllArray(O.None[string]())([]O.Option[string]{first(1), second(1)}))
	assert.Equal(t, 1, d2, "Go evaluates the slice literal before AltAllArray runs")
}

func TestPatternMatching_FoldAltMonoidIsLazy(t *testing.T) {
	var c1, c2, c3 int
	m := matchFold(
		counting(&c1, RO.None[int, string]()),
		counting(&c2, RO.Of[int]("second")),
		counting(&c3, RO.Of[int]("third")),
	)

	assert.Equal(t, 0, c1+c2+c3, "building the matcher must not run any case")

	assert.Equal(t, O.Some("second"), m(1))
	assert.Equal(t, 1, c1)
	assert.Equal(t, 1, c2)
	assert.Equal(t, 0, c3, "cases after the first match must not run")
}

func TestPatternMatching_FindFirstMapIsLazy(t *testing.T) {
	var c1, c2, c3 int
	m := matchFind(
		counting(&c1, RO.None[int, string]()),
		counting(&c2, RO.Of[int]("second")),
		counting(&c3, RO.Of[int]("third")),
	)

	assert.Equal(t, O.Some("second"), m(1))
	assert.Equal(t, 1, c1)
	assert.Equal(t, 1, c2)
	assert.Equal(t, 0, c3)
}

func TestPatternMatching_ReaderOptionAltIsLazy(t *testing.T) {
	var c1, c2 int
	m := F.Pipe1(
		counting(&c1, when(N.MoreThan(0), "positive")),
		RO.Alt(F.Constant(counting(&c2, RO.Of[int]("fallback")))),
	)

	assert.Equal(t, O.Some("positive"), m(1))
	assert.Equal(t, 0, c2, "fallback must not run after a match")

	assert.Equal(t, O.Some("fallback"), m(-1))
	assert.Equal(t, 1, c2)
	assert.Equal(t, 2, c1)
}

func TestPatternMatching_AllFormulationsAgree(t *testing.T) {
	cases := A.From(
		when(N.MoreThan(100), "huge"),
		when(N.MoreThan(10), "large"),
		when(N.MoreThan(0), "small"),
	)
	fold := matchFold(cases...)
	find := matchFind(cases...)
	altAll := matchAltAll(cases...)

	for _, n := range []int{-5, 0, 1, 10, 11, 100, 101} {
		assert.Equal(t, fold(n), find(n), "n=%d", n)
		assert.Equal(t, fold(n), altAll(n), "n=%d", n)
	}
	assert.Equal(t, O.Some("huge"), fold(101))
	assert.Equal(t, O.Some("large"), fold(100))
	assert.Equal(t, O.Some("small"), fold(1))
	assert.Equal(t, O.None[string](), fold(0))
}

func TestPatternMatching_OrderDefinesPriority(t *testing.T) {
	general := when(N.MoreThan(0), "positive")
	specific := when(N.MoreThan(100), "huge")

	assert.Equal(t, O.Some("huge"), matchFold(specific, general)(500))
	assert.Equal(t, O.Some("positive"), matchFold(general, specific)(500), "general case shadows specific one")
}

func TestPatternMatching_EmptyMatcherIsNone(t *testing.T) {
	assert.Equal(t, O.None[string](), matchFold[int, string]()(42))
	assert.Equal(t, O.None[string](), matchFind[int, string]()(42))
	assert.Equal(t, O.None[string](), matchAltAll[int, string]()(42))
}

func TestPatternMatching_DefaultMakesMatchTotal(t *testing.T) {
	sgn := sign()
	assert.Equal(t, "positive", sgn(3))
	assert.Equal(t, "negative", sgn(-3))
	assert.Equal(t, "zero", sgn(0))

	// RO.GetOrElse closes a ReaderOption with a default that sees the input
	total := F.Pipe1(
		matchFold(F.Flow2(O.FromPredicate(N.MoreThan(0)), O.Map(strconv.Itoa))),
		RO.GetOrElse(F.Flow2(strconv.Itoa, S.Prepend("non-positive: "))),
	)
	assert.Equal(t, "7", total(7))
	assert.Equal(t, "non-positive: -7", total(-7))

	// an always-Some case at the end has the same effect
	withCatchAll := matchFold(
		when(N.MoreThan(0), "positive"),
		RO.Of[int]("other"),
	)
	assert.Equal(t, O.Some("other"), withCatchAll(-1))
}

func TestPatternMatching_AltAllArrayStartWith(t *testing.T) {
	// Some as startWith wins over everything in the array
	assert.Equal(t, O.Some(1), O.AltAllArray(O.Some(1))(A.From(O.Some(2))))
	// None as startWith is the "no match" result
	assert.Equal(t, O.None[int](), O.AltAllArray(O.None[int]())(A.From(O.None[int]())))
	assert.Equal(t, O.None[int](), O.AltAllArray(O.None[int]())(nil))
}

func TestPatternMatching_AltChainIsLazy(t *testing.T) {
	var calls int
	expensive := counting(&calls, O.Of[int])

	assert.Equal(t, O.Some(1), F.Pipe1(O.Some(1), O.Alt(F.Nullary2(F.Constant(2), expensive))))
	assert.Equal(t, 0, calls)

	assert.Equal(t, O.Some(2), F.Pipe1(O.None[int](), O.Alt(F.Nullary2(F.Constant(2), expensive))))
	assert.Equal(t, 1, calls)
}

// Shape is a closed sum type used to demonstrate type-based matching.
type Shape interface{ isShape() }

type Circle struct{ Radius float64 }
type Rect struct{ W, H float64 }
type Square struct{ Side float64 }

func (Circle) isShape() {}
func (Rect) isShape()   {}
func (Square) isShape() {}

// caseOf turns a handler for one variant into a case for the whole sum type.
func caseOf[T Shape, B any](f func(T) B) RO.ReaderOption[Shape, B] {
	return F.Flow3(F.ToAny[Shape], O.InstanceOf[T], O.Map(f))
}

func TestPatternMatching_SumTypeCases(t *testing.T) {
	describe := matchFold(
		caseOf(S.Format[Circle]("circle %v")),
		caseOf(S.Format[Rect]("rect %v")),
	)

	assert.Equal(t, O.Some("circle {2}"), describe(Circle{Radius: 2}))
	assert.Equal(t, O.Some("rect {2 3}"), describe(Rect{W: 2, H: 3}))
	assert.Equal(t, O.None[string](), describe(Square{Side: 1}), "unhandled variant yields None")
}

func TestPatternMatching_PrismCases(t *testing.T) {
	// prism.GetOption is a ready-made case
	classify := matchFold(
		F.Flow2(prism.ParseInt().GetOption, O.Map(S.Format[int]("int %d"))),
		F.Flow2(prism.ParseFloat64().GetOption, O.Map(S.Format[float64]("float %g"))),
		F.Flow2(prism.ParseBool().GetOption, O.Map(S.Format[bool]("bool %t"))),
	)

	assert.Equal(t, O.Some("int 42"), classify("42"), "int wins because it comes first")
	assert.Equal(t, O.Some("float 3.5"), classify("3.5"))
	assert.Equal(t, O.Some("bool true"), classify("true"))
	assert.Equal(t, O.None[string](), classify("nope"))
}

func TestPatternMatching_FromValidationCases(t *testing.T) {
	// any func(A) (B, bool) is lifted into a case without a wrapper
	const key = "FP_GO_PATTERN_MATCHING_TEST"
	t.Setenv(key, "from env")

	lookup := matchFold(
		O.FromValidation(os.LookupEnv),
		RO.Of[string]("default"),
	)
	assert.Equal(t, O.Some("from env"), lookup(key))
	assert.Equal(t, O.Some("default"), lookup(key+"_UNSET"))
}

func TestPatternMatching_SliceMatching(t *testing.T) {
	parse := matchFold(withPrefix("0x", 16), prism.ParseInt64().GetOption)

	assert.Equal(t, O.Some[int64](255), parse("0xff"))
	assert.Equal(t, O.Some[int64](12), parse("12"))
	assert.Equal(t, O.None[int64](), parse("0xzz"))

	// FindFirstMap: first element of a slice that matches any case
	assert.Equal(t, O.Some[int64](16), A.FindFirstMap(parse)(A.From("a", "0x10", "7")))
	assert.Equal(t, O.None[int64](), A.FindFirstMap(parse)(A.From("a", "b")))
	// FindLastMap searches from the end
	assert.Equal(t, O.Some[int64](7), A.FindLastMap(parse)(A.From("a", "0x10", "7")))
	// FilterMap keeps all matches
	assert.Equal(t, A.From[int64](16, 7), A.FilterMap(parse)(A.From("a", "0x10", "7")))
}

func TestPatternMatching_ReaderResultAltKeepsLastError(t *testing.T) {
	errHex := errors.New("not hex")
	errDec := errors.New("not decimal")

	parse := F.Pipe1(
		withPrefixR("0x", 16, errHex),
		RR.Alt(F.Constant(decimalR(errDec))),
	)

	assert.Equal(t, result.Of[int64](255), parse("0xff"))
	assert.Equal(t, result.Of[int64](12), parse("12"))
	assert.Equal(t, result.Left[int64](errDec), parse("x"), "the last failing case determines the error")
}

func TestPatternMatching_ReaderResultFoldIsLazyAndKeepsLastError(t *testing.T) {
	errNoMatch := errors.New("no match")
	errHex := errors.New("not hex")
	errDec := errors.New("not decimal")
	altMonoid := RR.AltMonoid(F.Constant(RR.Left[string, int64](errNoMatch)))

	var hexCalls, decCalls int
	parse := A.Fold(altMonoid)(A.From(
		counting(&hexCalls, withPrefixR("0x", 16, errHex)),
		counting(&decCalls, decimalR(errDec)),
	))

	assert.Equal(t, result.Of[int64](255), parse("0xff"))
	assert.Equal(t, 1, hexCalls)
	assert.Equal(t, 0, decCalls, "dec must not run after hex matched")

	assert.Equal(t, result.Of[int64](12), parse("12"))
	assert.Equal(t, result.Left[int64](errDec), parse("x"), "the last failing case determines the error")

	// with no cases the zero value is the result
	assert.Equal(t, result.Left[int64](errNoMatch), A.Fold(altMonoid)(nil)("x"))
}
