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
	"fmt"
	"strconv"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/optics/prism"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	R "github.com/IBM/fp-go/v2/reader"
	RO "github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
)

// Example_pattern_matching demonstrates multi-branch pattern matching on a
// single value, the functional counterpart of a switch statement.
//
// Each case is a func(Request) Option[string] (a ReaderOption), composed
// point-free from a guard (O.FromPredicate) and a branch result (O.Map).
// Folding the cases with RO.AltMonoid yields one ReaderOption that tries the
// cases in order and stops at the first Some. RO.GetOrElse supplies the
// default branch and turns the partial match into a total function.
func Example_pattern_matching() {
	type Request struct {
		Method string
		Path   string
		Body   string
	}

	getMethod := func(r Request) string { return r.Method }
	getPath := func(r Request) string { return r.Path }
	describeCreate := func(r Request) string {
		return fmt.Sprintf("Creating: %s with body: %s", r.Path, r.Body)
	}

	// case combinator: guard -> branch result
	on := func(guard P.Predicate[Request], branch func(Request) string) RO.ReaderOption[Request, string] {
		return F.Flow2(O.FromPredicate(guard), O.Map(branch))
	}
	isMethod := func(m string) P.Predicate[Request] {
		return F.Flow2(getMethod, P.IsStrictEqual[string]()(m))
	}

	handleRequest := F.Pipe2(
		A.From(
			on(isMethod("GET"), F.Flow2(getPath, S.Format[string]("Fetching: %s"))),
			on(isMethod("POST"), describeCreate),
			on(isMethod("DELETE"), F.Flow2(getPath, S.Format[string]("Deleting: %s"))),
		),
		A.Fold(RO.AltMonoid[Request, string]()),
		RO.GetOrElse(F.Flow2(getMethod, S.Format[string]("Unsupported method: %s"))),
	)

	requests := []Request{
		{Method: "GET", Path: "/users"},
		{Method: "POST", Path: "/users", Body: `{"name":"Alice"}`},
		{Method: "DELETE", Path: "/users/123"},
		{Method: "PATCH", Path: "/users/123"},
	}

	for _, req := range requests {
		fmt.Println(handleRequest(req))
	}

	// Output:
	// Fetching: /users
	// Creating: /users with body: {"name":"Alice"}
	// Deleting: /users/123
	// Unsupported method: PATCH
}

// Example_pattern_matching_lazy shows that folding cases with RO.AltMonoid is
// lazy: cases after the first match are never invoked. Collecting the results
// into a []Option and calling option.AltAllArray instead would run every case
// up front, because Go evaluates the slice literal eagerly.
func Example_pattern_matching_lazy() {
	// trace is deliberately side-effecting so the evaluation order is visible
	trace := func(name string, res O.Option[string]) RO.ReaderOption[int, string] {
		return func(int) O.Option[string] {
			fmt.Println("try", name)
			return res
		}
	}

	match := A.Fold(RO.AltMonoid[int, string]())(A.From(
		trace("a", O.None[string]()),
		trace("b", O.Some("b")),
		trace("c", O.Some("c")),
	))

	fmt.Println(match(0))

	// Output:
	// try a
	// try b
	// Some[string](b)
}

// Example_pattern_matching_numeric classifies numbers using guard cases
// built point-free from predicates: FromPredicate turns the guard into an
// Option and Map produces the branch result.
func Example_pattern_matching_numeric() {
	when := func(guard P.Predicate[int], label string) RO.ReaderOption[int, string] {
		return F.Flow2(O.FromPredicate(guard), O.Map(F.Constant1[int](label)))
	}

	classify := F.Pipe2(
		A.From(
			when(P.IsZero[int](), "zero"),
			when(N.LessThan(0), "negative"),
			when(N.LessThan(11), "small positive"),
		),
		A.Fold(RO.AltMonoid[int, string]()),
		RO.GetOrElse(F.Constant1[int]("large positive")),
	)

	for _, n := range []int{0, -5, 3, 15, -100, 10} {
		fmt.Printf("%d: %s\n", n, classify(n))
	}

	// Output:
	// 0: zero
	// -5: negative
	// 3: small positive
	// 15: large positive
	// -100: negative
	// 10: small positive
}

// Example_pattern_matching_with_guards demonstrates guards combined with
// predicate.And and cases ordered by specificity: the critical-error case
// must come before the general error case, otherwise it would be shadowed.
func Example_pattern_matching_with_guards() {
	type Event struct {
		Type     string
		Priority int
		Message  string
	}

	getType := func(e Event) string { return e.Type }
	getPriority := func(e Event) int { return e.Priority }
	getMessage := func(e Event) string { return e.Message }

	isType := func(t string) P.Predicate[Event] {
		return F.Flow2(getType, P.IsStrictEqual[string]()(t))
	}
	isUrgent := F.Flow2(getPriority, N.MoreThan(8))

	on := func(guard P.Predicate[Event], prefix string) RO.ReaderOption[Event, string] {
		return F.Flow2(O.FromPredicate(guard), O.Map(F.Flow2(getMessage, S.Prepend(prefix))))
	}

	formatEvent := F.Pipe2(
		A.From(
			on(F.Pipe1(isType("error"), P.And(isUrgent)), "CRITICAL: "),
			on(isType("error"), "ERROR: "),
			on(isType("warning"), "WARNING: "),
		),
		A.Fold(RO.AltMonoid[Event, string]()),
		RO.GetOrElse(F.Flow2(getMessage, S.Prepend("INFO: "))),
	)

	events := []Event{
		{Type: "error", Priority: 10, Message: "System failure"},
		{Type: "error", Priority: 5, Message: "Connection lost"},
		{Type: "warning", Priority: 3, Message: "High memory usage"},
		{Type: "info", Priority: 1, Message: "User logged in"},
	}

	for _, event := range events {
		fmt.Println(formatEvent(event))
	}

	// Output:
	// CRITICAL: System failure
	// ERROR: Connection lost
	// WARNING: High memory usage
	// INFO: User logged in
}

// Example_pattern_matching_array combines both levels of matching: a
// multi-format parser (first case that matches a string) and FindFirstMap
// (first element of a slice that the parser accepts).
func Example_pattern_matching_array() {
	// parseBase(b) is strconv.ParseInt with base b and bit size 64 bound
	parseBase := func(base int) O.Kleisli[string, int64] {
		return F.Flow2(
			F.Bind23of3(result.Eitherize3(strconv.ParseInt))(base, 64),
			result.ToOption[int64],
		)
	}
	// withPrefix matches strings with the prefix and parses the rest
	withPrefix := func(prefix string, base int) RO.ReaderOption[string, int64] {
		return F.Flow3(
			O.FromPredicate(F.Bind2nd(strings.HasPrefix, prefix)),
			O.Map(F.Bind2nd(strings.TrimPrefix, prefix)),
			O.Chain(parseBase(base)),
		)
	}

	parseNumber := A.Fold(RO.AltMonoid[string, int64]())(A.From(
		withPrefix("0x", 16),
		withPrefix("0b", 2),
		withPrefix("0o", 8),
		prism.ParseInt64().GetOption,
	))

	inputs := []string{"invalid", "0x2A", "also bad", "17"}

	fmt.Println(A.FindFirstMap(parseNumber)(inputs))
	fmt.Println(A.FindLastMap(parseNumber)(inputs))
	fmt.Println(A.FilterMap(parseNumber)([]string{"0b101", "x", "0o17", "7"}))

	// Output:
	// Some[int64](42)
	// Some[int64](17)
	// [5 15 7]
}

// Example_pattern_matching_sum_type is the functional counterpart of a type
// switch: O.InstanceOf performs the type assertion, so each case handles
// exactly one variant and the fold picks the one matching the dynamic type.
func Example_pattern_matching_sum_type() {
	type Circle struct{ R float64 }
	type Rect struct{ W, H float64 }

	circleArea := func(c Circle) float64 { return 3 * c.R * c.R }
	rectArea := func(r Rect) float64 { return r.W * r.H }

	area := F.Pipe2(
		A.From(
			F.Flow2(O.InstanceOf[Circle], O.Map(circleArea)),
			F.Flow2(O.InstanceOf[Rect], O.Map(rectArea)),
		),
		A.Fold(RO.AltMonoid[any, float64]()),
		RO.GetOrElse(F.Constant1[any](-1.0)),
	)

	fmt.Println(area(Circle{R: 2}))
	fmt.Println(area(Rect{W: 2, H: 3}))
	fmt.Println(area("not a shape"))

	// Output:
	// 12
	// 6
	// -1
}

// Example_pattern_matching_find_first_map shows the alternative formulation
// with FindFirstMap: the slice holds the cases, and reader.Read(path) applies
// each case to the input until one returns Some. Composing reader.Read,
// FindFirstMap and reader.Read(cases) yields the matcher point-free.
func Example_pattern_matching_find_first_map() {
	prefix := func(p, target string) RO.ReaderOption[string, string] {
		return F.Flow2(O.FromPredicate(F.Bind2nd(strings.HasPrefix, p)), O.Map(F.Constant1[string](target)))
	}

	cases := A.From(
		prefix("/api/", "api"),
		prefix("/static/", "assets"),
	)

	route := F.Flow4(
		R.Read[O.Option[string], string],                        // path -> "apply a case to path"
		A.FindFirstMap[RO.ReaderOption[string, string], string], // -> first case that matches
		R.Read[O.Option[string]](cases),                         // run it over the cases
		O.GetOrElse(F.Constant("page")),                         // default
	)

	fmt.Println(route("/api/users"))
	fmt.Println(route("/static/app.js"))
	fmt.Println(route("/about"))

	// Output:
	// api
	// assets
	// page
}

// Example_pattern_matching_alt shows the binary form: RO.Alt appends a
// fallback case. F.Constant makes the fallback lazy, so it only runs when the
// preceding cases return None.
func Example_pattern_matching_alt() {
	fromEnv := F.Flow2(O.FromPredicate(F.Bind2nd(strings.HasPrefix, "$")), O.Map(S.Prepend("env:")))
	fromFile := F.Flow2(O.FromPredicate(F.Bind2nd(strings.HasPrefix, "@")), O.Map(S.Prepend("file:")))

	resolve := F.Pipe2(
		fromEnv,
		RO.Alt(F.Constant(fromFile)),
		RO.GetOrElse(S.Prepend("literal:")),
	)

	fmt.Println(resolve("$HOME"))
	fmt.Println(resolve("@config.json"))
	fmt.Println(resolve("value"))

	// Output:
	// env:$HOME
	// file:@config.json
	// literal:value
}
