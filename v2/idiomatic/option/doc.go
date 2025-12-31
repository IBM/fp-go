// Copyright (c) 2025 IBM Corp.
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

// Package option implements the Option monad using idiomatic Go tuple signatures.
//
// Unlike the standard option package which uses wrapper structs, this package represents
// Options as tuples (value, bool) where the boolean indicates presence (true) or absence (false).
// This approach is more idiomatic in Go and has better performance characteristics.
//
// # Type Signatures
//
// The core types used in this package are:
//
//	Operator[A, B any] = func(A, bool) (B, bool)  // Transforms an Option[A] to Option[B]
//	Kleisli[A, B any]  = func(A) (B, bool)        // Monadic function from A to Option[B]
//
// # Basic Usage
//
// Create an Option with Some or None:
//
//	some := Some(42)           // (42, true)
//	none := None[int]()        // (0, false)
//	opt := Of(42)              // Alternative to Some: (42, true)
//
// Check if an Option contains a value:
//
//	value, ok := Some(42)
//	if ok {
//	    // value == 42
//	}
//
//	if IsSome(Some(42)) {
//	    // Option contains a value
//	}
//	if IsNone(None[int]()) {
//	    // Option is empty
//	}
//
// Extract values:
//
//	value, ok := Some(42)      // Direct tuple unpacking: value == 42, ok == true
//	value := GetOrElse(func() int { return 0 })(Some(42))  // Returns 42
//	value := GetOrElse(func() int { return 0 })(None[int]())  // Returns 0
//
// # Transformations
//
// Map transforms the contained value:
//
//	double := Map(N.Mul(2))
//	result := double(Some(21))  // (42, true)
//	result := double(None[int]())  // (0, false)
//
// Chain sequences operations that may fail:
//
//	validate := Chain(func(x int) (int, bool) {
//	    if x > 0 { return x * 2, true }
//	    return 0, false
//	})
//	result := validate(Some(5))  // (10, true)
//	result := validate(Some(-1))  // (0, false)
//
// Filter keeps values that satisfy a predicate:
//
//	isPositive := Filter(N.MoreThan(0))
//	result := isPositive(Some(5))  // (5, true)
//	result := isPositive(Some(-1))  // (0, false)
//
// # Working with Collections
//
// Transform arrays using TraverseArray:
//
//	doublePositive := func(x int) (int, bool) {
//	    if x > 0 { return x * 2, true }
//	    return 0, false
//	}
//	result := TraverseArray(doublePositive)([]int{1, 2, 3})  // ([2, 4, 6], true)
//	result := TraverseArray(doublePositive)([]int{1, -2, 3})  // ([], false)
//
// Transform with indexes:
//
//	f := func(i int, x int) (int, bool) {
//	    if x > i { return x, true }
//	    return 0, false
//	}
//	result := TraverseArrayWithIndex(f)([]int{1, 2, 3})  // ([1, 2, 3], true)
//
// Transform records (maps):
//
//	double := func(x int) (int, bool) { return x * 2, true }
//	result := TraverseRecord(double)(map[string]int{"a": 1, "b": 2})
//	// (map[string]int{"a": 2, "b": 4}, true)
//
// # Algebraic Operations
//
// Option supports various algebraic structures:
//
//   - Functor: Map operations for transforming values
//   - Applicative: Ap operations for applying wrapped functions
//   - Monad: Chain operations for sequencing computations
//   - Alternative: Alt operations for providing fallbacks
//
// Applicative example:
//
//	fab := Some(N.Mul(2))
//	fa := Some(21)
//	result := Ap[int](fa)(fab)  // (42, true)
//
// Alternative example:
//
//	withDefault := Alt(func() (int, bool) { return 100, true })
//	result := withDefault(Some(42))  // (42, true)
//	result := withDefault(None[int]())  // (100, true)
//
// # Conversion Functions
//
// Convert predicates to Options:
//
//	isPositive := FromPredicate(N.MoreThan(0))
//	result := isPositive(5)  // (5, true)
//	result := isPositive(-1)  // (0, false)
//
// Convert nullable pointers to Options:
//
//	var ptr *int = nil
//	result := FromNillable(ptr)  // (nil, false)
//	val := 42
//	result := FromNillable(&val)  // (&val, true)
//
// Convert zero/non-zero values to Options:
//
//	result := FromZero[int]()(0)  // (0, true)
//	result := FromZero[int]()(5)  // (0, false)
//	result := FromNonZero[int]()(5)  // (5, true)
//	result := FromNonZero[int]()(0)  // (0, false)
//
// Use equality-based conversion:
//
//	import "github.com/IBM/fp-go/v2/eq"
//	equals42 := FromEq(eq.FromStrictEquals[int]())(42)
//	result := equals42(42)  // (42, true)
//	result := equals42(10)  // (0, false)
//
// # Do-Notation Style
//
// Build complex computations using do-notation:
//
//	type Result struct {
//	    x int
//	    y int
//	    sum int
//	}
//
//	result := F.Pipe3(
//	    Do(Result{}),
//	    Bind(func(x int) func(Result) Result {
//	        return func(r Result) Result { r.x = x; return r }
//	    }, func(r Result) (int, bool) { return Some(10) }),
//	    Bind(func(y int) func(Result) Result {
//	        return func(r Result) Result { r.y = y; return r }
//	    }, func(r Result) (int, bool) { return Some(20) }),
//	    Let(func(sum int) func(Result) Result {
//	        return func(r Result) Result { r.sum = sum; return r }
//	    }, func(r Result) int { return r.x + r.y }),
//	)  // (Result{x: 10, y: 20, sum: 30}, true)
//
// # Lens-Based Operations
//
// Use lenses for cleaner field updates:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	ageLens := lens.MakeLens(
//	    func(p Person) int { return p.Age },
//	    func(p Person, age int) Person { p.Age = age; return p },
//	)
//
//	// Update using a lens
//	incrementAge := BindL(ageLens, func(age int) (int, bool) {
//	    if age < 120 { return age + 1, true }
//	    return 0, false
//	})
//	result := incrementAge(Some(Person{Name: "Alice", Age: 30}))
//	// (Person{Name: "Alice", Age: 31}, true)
//
//	// Set using a lens
//	setAge := LetToL(ageLens, 25)
//	result := setAge(Some(Person{Name: "Bob", Age: 30}))
//	// (Person{Name: "Bob", Age: 25}, true)
//
// # Folding and Reducing
//
// Fold provides a way to handle both Some and None cases:
//
//	handler := Fold(
//	    func() string { return "no value" },
//	    func(x int) string { return fmt.Sprintf("value: %d", x) },
//	)
//	result := handler(Some(42))  // "value: 42"
//	result := handler(None[int]())  // "no value"
//
// Reduce folds an Option into a single value:
//
//	sum := Reduce(func(acc, val int) int { return acc + val }, 0)
//	result := sum(Some(5))  // 5
//	result := sum(None[int]())  // 0
//
// # Debugging
//
// Convert Options to strings for debugging:
//
//	str := ToString(Some(42))  // "Some[int](42)"
//	str := ToString(None[int]())  // "None[int]"
//
// # Subpackages
//
//   - option/number: Number conversion utilities for working with Options
package option

//go:generate go run .. option --count 10 --filename gen.go
