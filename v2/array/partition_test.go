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

package array

import (
	"fmt"
	"strconv"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	N "github.com/IBM/fp-go/v2/number"
	P "github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

// isEven is a reusable predicate used across partition tests.
func isEven(v int) bool { return v%2 == 0 }

// TestMonadPartitionWithIndex_SplitsOnIndexAndValue verifies that matching
// elements go to Tail and non-matching elements go to Head, and that the
// predicate receives the correct index and value together.
func TestMonadPartitionWithIndex_SplitsOnIndexAndValue(t *testing.T) {
	src := []int{10, 21, 30, 43, 50}

	// Keep elements whose index is even (positions 0, 2, 4).
	result := MonadPartitionWithIndex(src, func(i int, _ int) bool {
		return i%2 == 0
	})

	assert.Equal(t, []int{21, 43}, P.Head(result))
	assert.Equal(t, []int{10, 30, 50}, P.Tail(result))
}

// TestMonadPartitionWithIndex_ValueAndIndexTogether verifies that the predicate
// can combine index and value in its decision.
func TestMonadPartitionWithIndex_ValueAndIndexTogether(t *testing.T) {
	src := []int{0, 1, 2, 3, 4}

	// Match when element equals its own index.
	result := MonadPartitionWithIndex(src, func(i, v int) bool {
		return i == v
	})

	// All elements equal their index: everything goes to Tail, Head is empty.
	assert.Nil(t, P.Head(result))
	assert.Equal(t, []int{0, 1, 2, 3, 4}, P.Tail(result))
}

// TestMonadPartitionWithIndex_AllMatch verifies that Head is nil when all
// elements satisfy the predicate.
func TestMonadPartitionWithIndex_AllMatch(t *testing.T) {
	src := []int{2, 4, 6}

	result := MonadPartitionWithIndex(src, func(_ int, v int) bool {
		return isEven(v)
	})

	assert.Nil(t, P.Head(result))
	assert.Equal(t, []int{2, 4, 6}, P.Tail(result))
}

// TestMonadPartitionWithIndex_NoneMatch verifies that Tail is nil when no
// element satisfies the predicate.
func TestMonadPartitionWithIndex_NoneMatch(t *testing.T) {
	src := []int{1, 3, 5}

	result := MonadPartitionWithIndex(src, func(_ int, v int) bool {
		return isEven(v)
	})

	assert.Equal(t, []int{1, 3, 5}, P.Head(result))
	assert.Nil(t, P.Tail(result))
}

// TestMonadPartitionWithIndex_EmptySlice verifies that an empty input produces
// nil on both sides.
func TestMonadPartitionWithIndex_EmptySlice(t *testing.T) {
	result := MonadPartitionWithIndex([]int{}, func(_ int, v int) bool {
		return isEven(v)
	})

	assert.Nil(t, P.Head(result))
	assert.Nil(t, P.Tail(result))
}

// TestMonadPartitionWithIndex_NilSlice verifies that a nil input is treated
// identically to an empty slice.
func TestMonadPartitionWithIndex_NilSlice(t *testing.T) {
	var src []int

	result := MonadPartitionWithIndex(src, func(_ int, v int) bool {
		return isEven(v)
	})

	assert.Nil(t, P.Head(result))
	assert.Nil(t, P.Tail(result))
}

// TestMonadPartitionWithIndex_PreservesOrder verifies that the relative order
// of elements is maintained in both output slices.
func TestMonadPartitionWithIndex_PreservesOrder(t *testing.T) {
	src := []string{"a", "bb", "ccc", "dddd", "eeeee"}

	// Route strings whose index is odd to Tail.
	result := MonadPartitionWithIndex(src, func(i int, _ string) bool {
		return i%2 != 0
	})

	assert.Equal(t, []string{"a", "ccc", "eeeee"}, P.Head(result))
	assert.Equal(t, []string{"bb", "dddd"}, P.Tail(result))
}

// TestPartitionWithIndex_CurriedIsReusable verifies that the curried form
// produces the same result as calling MonadPartitionWithIndex directly, and
// that the returned function can be applied to multiple inputs.
func TestPartitionWithIndex_CurriedIsReusable(t *testing.T) {
	// Keep elements whose value exceeds twice their index.
	splitFn := PartitionWithIndex(func(i, v int) bool {
		return v > 2*i
	})

	src1 := []int{5, 5, 5}
	src2 := []int{0, 1, 2}

	r1 := splitFn(src1)
	r2 := splitFn(src2)

	// src1: v=5 > 2*0=0 ✓, v=5 > 2*1=2 ✓, v=5 > 2*2=4 ✓ → all to Tail
	assert.Nil(t, P.Head(r1))
	assert.Equal(t, []int{5, 5, 5}, P.Tail(r1))

	// src2: v=0 > 0 ✗, v=1 > 2 ✗, v=2 > 4 ✗ → all to Head
	assert.Equal(t, []int{0, 1, 2}, P.Head(r2))
	assert.Nil(t, P.Tail(r2))
}

// TestPartitionWithIndex_NilSlice verifies that the curried form handles a nil
// slice without panicking.
func TestPartitionWithIndex_NilSlice(t *testing.T) {
	var src []int

	result := PartitionWithIndex(func(_ int, v int) bool { return isEven(v) })(src)

	assert.Nil(t, P.Head(result))
	assert.Nil(t, P.Tail(result))
}

// TestPartitionWithIndex_MatchesMonadPartitionWithIndex verifies that the
// curried PartitionWithIndex is equivalent to its uncurried counterpart.
func TestPartitionWithIndex_MatchesMonadPartitionWithIndex(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6}
	pred := func(i, v int) bool { return i%2 == 0 && isEven(v) }

	direct := MonadPartitionWithIndex(src, pred)
	curried := PartitionWithIndex(pred)(src)

	assert.Equal(t, P.Head(direct), P.Head(curried))
	assert.Equal(t, P.Tail(direct), P.Tail(curried))
}

// ExampleMonadPartitionWithIndex demonstrates splitting a log of HTTP status
// codes into errors (>= 400) and successes (< 400), using the request index
// to label each entry in a diagnostic message.
func ExampleMonadPartitionWithIndex() {
	statuses := []int{200, 404, 201, 500, 204}

	result := MonadPartitionWithIndex(statuses, func(_ int, code int) bool {
		return code >= 400
	})

	fmt.Println("successes:", P.Head(result))
	fmt.Println("errors:   ", P.Tail(result))

	// Output:
	// successes: [200 201 204]
	// errors:    [404 500]
}

// ExampleMonadPartitionWithIndex_usingIndex demonstrates using the element
// index to keep only the first occurrence of each parity class.
func ExampleMonadPartitionWithIndex_usingIndex() {
	words := []string{"alpha", "beta", "gamma", "delta", "epsilon"}

	// Route words at even positions to Tail (selected), odd positions to Head (rest).
	result := MonadPartitionWithIndex(words, func(i int, _ string) bool {
		return i%2 == 0
	})

	fmt.Println("odd-positioned: ", P.Head(result))
	fmt.Println("even-positioned:", P.Tail(result))

	// Output:
	// odd-positioned:  [beta delta]
	// even-positioned: [alpha gamma epsilon]
}

// ExamplePartitionWithIndex demonstrates building a reusable classifier that
// flags duplicate entries (all occurrences after the first) by routing them to
// the right (Tail) side.
func ExamplePartitionWithIndex() {
	// Keep only the first time a value is seen (route later duplicates to Tail).
	seen := map[int]bool{}
	firstOccurrence := PartitionWithIndex(func(_ int, v int) bool {
		if seen[v] {
			return true // duplicate → Tail
		}
		seen[v] = true
		return false // first time → Head
	})

	src := []int{1, 2, 3, 2, 4, 1, 5}
	result := firstOccurrence(src)

	fmt.Println("unique:     ", P.Head(result))
	fmt.Println("duplicates:", P.Tail(result))

	// Output:
	// unique:      [1 2 3 4 5]
	// duplicates: [2 1]
}

// ExamplePartitionWithIndex_splitByPosition demonstrates reusing the same
// partition function to separate elements that land at even and odd positions
// across two different slices.
func ExamplePartitionWithIndex_splitByPosition() {
	evenPositions := PartitionWithIndex(func(i int, _ int) bool {
		return i%2 == 0
	})

	a := []int{10, 20, 30, 40}
	b := []int{1, 2, 3}

	ra := evenPositions(a)
	rb := evenPositions(b)

	fmt.Println("a odd-pos: ", P.Head(ra))
	fmt.Println("a even-pos:", P.Tail(ra))
	fmt.Println("b odd-pos: ", P.Head(rb))
	fmt.Println("b even-pos:", P.Tail(rb))

	// Output:
	// a odd-pos:  [20 40]
	// a even-pos: [10 30]
	// b odd-pos:  [2]
	// b even-pos: [1 3]
}

// ExamplePartitionWithIndex_moreThanThreshold demonstrates applying the N
// helper package to the value side while ignoring the index.
func ExamplePartitionWithIndex_moreThanThreshold() {
	above50 := PartitionWithIndex(func(_ int, v int) bool {
		return N.MoreThan(50)(v)
	})

	scores := []int{40, 75, 55, 30, 90}
	result := above50(scores)

	fmt.Println("at-or-below:", P.Head(result))
	fmt.Println("above:      ", P.Tail(result))

	// Output:
	// at-or-below: [40 30]
	// above:       [75 55 90]
}

// ----------------------------------------------------------------------------
// MonadPartitionMap tests
// ----------------------------------------------------------------------------

// classifyInt routes positive integers to Right and non-positive integers
// to Left (stored as their string representation).
func classifyInt(v int) E.Either[string, int] {
	if v > 0 {
		return E.Right[string](v)
	}
	return E.Left[int](strconv.Itoa(v))
}

// TestMonadPartitionMap_MixedInput verifies that positive values land in Tail
// and non-positive values land in Head (as strings).
func TestMonadPartitionMap_MixedInput(t *testing.T) {
	src := []int{-1, 2, 0, 4, -3}

	result := MonadPartitionMap(src, classifyInt)

	assert.Equal(t, []string{"-1", "0", "-3"}, P.Head(result))
	assert.Equal(t, []int{2, 4}, P.Tail(result))
}

// TestMonadPartitionMap_AllRight verifies that Head is nil when every element
// maps to Right.
func TestMonadPartitionMap_AllRight(t *testing.T) {
	src := []int{1, 2, 3}

	result := MonadPartitionMap(src, classifyInt)

	assert.Nil(t, P.Head(result))
	assert.Equal(t, []int{1, 2, 3}, P.Tail(result))
}

// TestMonadPartitionMap_AllLeft verifies that Tail is nil when every element
// maps to Left.
func TestMonadPartitionMap_AllLeft(t *testing.T) {
	src := []int{-1, 0, -3}

	result := MonadPartitionMap(src, classifyInt)

	assert.Equal(t, []string{"-1", "0", "-3"}, P.Head(result))
	assert.Nil(t, P.Tail(result))
}

// TestMonadPartitionMap_EmptySlice verifies that an empty input produces nil
// on both sides.
func TestMonadPartitionMap_EmptySlice(t *testing.T) {
	result := MonadPartitionMap([]int{}, classifyInt)

	assert.Nil(t, P.Head(result))
	assert.Nil(t, P.Tail(result))
}

// TestMonadPartitionMap_NilSlice verifies that a nil input is treated
// identically to an empty slice.
func TestMonadPartitionMap_NilSlice(t *testing.T) {
	var src []int

	result := MonadPartitionMap(src, classifyInt)

	assert.Nil(t, P.Head(result))
	assert.Nil(t, P.Tail(result))
}

// TestMonadPartitionMap_PreservesOrder verifies that the relative order of
// elements is maintained in both output slices.
func TestMonadPartitionMap_PreservesOrder(t *testing.T) {
	src := []int{3, -1, 5, -2, 7}

	result := MonadPartitionMap(src, classifyInt)

	assert.Equal(t, []string{"-1", "-2"}, P.Head(result))
	assert.Equal(t, []int{3, 5, 7}, P.Tail(result))
}

// TestMonadPartitionMap_TypeTransformation verifies that Left and Right can
// carry values of different types than the input.
func TestMonadPartitionMap_TypeTransformation(t *testing.T) {
	// Classify strings: parse as int → Right(int), else Left(original string).
	classify := func(s string) E.Either[string, int] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return E.Left[int](s)
		}
		return E.Right[string](n)
	}

	src := []string{"42", "hello", "7", "world", "0"}
	result := MonadPartitionMap(src, classify)

	assert.Equal(t, []string{"hello", "world"}, P.Head(result))
	assert.Equal(t, []int{42, 7, 0}, P.Tail(result))
}

// ----------------------------------------------------------------------------
// PartitionMap tests
// ----------------------------------------------------------------------------

// TestPartitionMap_CurriedIsReusable verifies that the curried form produces
// the same result as calling MonadPartitionMap directly, and that the returned
// function can be applied to multiple inputs.
func TestPartitionMap_CurriedIsReusable(t *testing.T) {
	splitFn := PartitionMap(classifyInt)

	src1 := []int{1, -2, 3}
	src2 := []int{-4, -5, 6}

	r1 := splitFn(src1)
	r2 := splitFn(src2)

	assert.Equal(t, []string{"-2"}, P.Head(r1))
	assert.Equal(t, []int{1, 3}, P.Tail(r1))

	assert.Equal(t, []string{"-4", "-5"}, P.Head(r2))
	assert.Equal(t, []int{6}, P.Tail(r2))
}

// TestPartitionMap_NilSlice verifies that the curried form handles a nil slice
// without panicking.
func TestPartitionMap_NilSlice(t *testing.T) {
	var src []int

	result := PartitionMap(classifyInt)(src)

	assert.Nil(t, P.Head(result))
	assert.Nil(t, P.Tail(result))
}

// TestPartitionMap_MatchesMonadPartitionMap verifies that the curried
// PartitionMap is equivalent to its uncurried counterpart.
func TestPartitionMap_MatchesMonadPartitionMap(t *testing.T) {
	src := []int{-3, 1, -2, 4, 0, 5}

	direct := MonadPartitionMap(src, classifyInt)
	curried := PartitionMap(classifyInt)(src)

	assert.Equal(t, P.Head(direct), P.Head(curried))
	assert.Equal(t, P.Tail(direct), P.Tail(curried))
}

// ----------------------------------------------------------------------------
// MonadPartitionMap examples
// ----------------------------------------------------------------------------

// ExampleMonadPartitionMap demonstrates routing HTTP status codes: codes >= 400
// are mapped to an error string (Left/Head) while codes < 400 keep their
// integer value (Right/Tail).
func ExampleMonadPartitionMap() {
	classify := func(code int) E.Either[string, int] {
		if code >= 400 {
			return E.Left[int](fmt.Sprintf("error %d", code))
		}
		return E.Right[string](code)
	}

	statuses := []int{200, 404, 201, 500, 204}
	result := MonadPartitionMap(statuses, classify)

	fmt.Println("successes:", P.Tail(result))
	fmt.Println("errors:   ", P.Head(result))

	// Output:
	// successes: [200 201 204]
	// errors:    [error 404 error 500]
}

// ExampleMonadPartitionMap_parseStrings demonstrates splitting a mixed slice
// of strings into successfully-parsed integers (Right/Tail) and unparseable
// strings (Left/Head).
func ExampleMonadPartitionMap_parseStrings() {
	classify := func(s string) E.Either[string, int] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return E.Left[int](s)
		}
		return E.Right[string](n)
	}

	tokens := []string{"10", "foo", "42", "bar", "7"}
	result := MonadPartitionMap(tokens, classify)

	fmt.Println("valid:  ", P.Tail(result))
	fmt.Println("invalid:", P.Head(result))

	// Output:
	// valid:   [10 42 7]
	// invalid: [foo bar]
}

// ----------------------------------------------------------------------------
// PartitionMap examples
// ----------------------------------------------------------------------------

// ExamplePartitionMap demonstrates building a reusable classifier that
// separates positive numbers (Right/Tail) from non-positive numbers whose
// absolute value is stored as the Left/Head payload.
func ExamplePartitionMap() {
	splitPositive := PartitionMap(func(v int) E.Either[int, int] {
		if v > 0 {
			return E.Right[int](v)
		}
		return E.Left[int](-v) // store absolute value on the left
	})

	src := []int{3, -1, 0, 7, -4}
	result := splitPositive(src)

	fmt.Println("positive:         ", P.Tail(result))
	fmt.Println("non-positive (abs):", P.Head(result))

	// Output:
	// positive:          [3 7]
	// non-positive (abs): [1 0 4]
}

// ExamplePartitionMap_reuseAcrossSlices demonstrates that the curried function
// returned by PartitionMap can be applied to multiple independent slices.
func ExamplePartitionMap_reuseAcrossSlices() {
	aboveThreshold := PartitionMap(func(v int) E.Either[int, int] {
		if N.MoreThan(50)(v) {
			return E.Right[int](v)
		}
		return E.Left[int](v)
	})

	a := []int{40, 75, 55, 30}
	b := []int{60, 20, 90}

	ra := aboveThreshold(a)
	rb := aboveThreshold(b)

	fmt.Println("a above:", P.Tail(ra))
	fmt.Println("a below:", P.Head(ra))
	fmt.Println("b above:", P.Tail(rb))
	fmt.Println("b below:", P.Head(rb))

	// Output:
	// a above: [75 55]
	// a below: [40 30]
	// b above: [60 90]
	// b below: [20]
}
