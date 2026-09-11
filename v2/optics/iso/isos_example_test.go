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

package iso

import (
	"fmt"
	"strings"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/array/nonempty"
	"github.com/IBM/fp-go/v2/boolean"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
)

// ---------------------------------------------------------------------------
// UTF8String
// ---------------------------------------------------------------------------

// ExampleUTF8String demonstrates converting between a UTF-8 byte slice and a string.
func ExampleUTF8String() {
	iso := UTF8String()

	// []byte → string
	fmt.Println(iso.Get([]byte{104, 101, 108, 108, 111}))

	// string → []byte; the accented rune expands to two bytes
	fmt.Println(iso.ReverseGet("héllo"))
	// Output:
	// hello
	// [104 195 169 108 108 111]
}

// ExampleUTF8String_roundTrip demonstrates the round-trip isomorphism laws.
func ExampleUTF8String_roundTrip() {
	iso := UTF8String()

	// ReverseGet(Get(bs)) == bs
	fmt.Println(string(iso.ReverseGet(iso.Get([]byte("round trip")))))

	// Get(ReverseGet(s)) == s
	fmt.Println(iso.Get(iso.ReverseGet("hello")))
	// Output:
	// round trip
	// hello
}

// ExampleUTF8String_modify demonstrates what an Iso buys you: Modify lifts a
// plain string function into []byte space, so no conversion is written by hand.
func ExampleUTF8String_modify() {
	shout := F.Pipe1(UTF8String(), Modify[[]byte](strings.ToUpper))

	fmt.Println(string(shout([]byte("hello"))))
	// Output:
	// HELLO
}

// ---------------------------------------------------------------------------
// Lines
// ---------------------------------------------------------------------------

// ExampleLines demonstrates joining a slice into newline-separated text and
// splitting that text back into lines.
func ExampleLines() {
	iso := Lines()

	fmt.Printf("%q\n", iso.Get([]string{"alpha", "beta", "gamma"}))
	fmt.Printf("%q\n", iso.ReverseGet("one\ntwo"))
	// Output:
	// "alpha\nbeta\ngamma"
	// ["one" "two"]
}

// ExampleLines_roundTrip demonstrates the round-trip laws, including the
// trailing-separator edge case called out in the documentation.
func ExampleLines_roundTrip() {
	iso := Lines()

	// ReverseGet(Get(ls)) == ls, empty elements included
	fmt.Printf("%q\n", iso.ReverseGet(iso.Get([]string{"a", "", "b"})))

	// A trailing newline yields an empty final element
	fmt.Printf("%q\n", iso.ReverseGet("a\nb\n"))
	// Output:
	// ["a" "" "b"]
	// ["a" "b" ""]
}

// ExampleLines_modify demonstrates lifting a whole-text transformation into the
// []string space with Modify.
func ExampleLines_modify() {
	shout := F.Pipe1(Lines(), Modify[[]string](strings.ToUpper))

	fmt.Printf("%q\n", shout([]string{"hello", "world"}))
	// Output:
	// ["HELLO" "WORLD"]
}

// ---------------------------------------------------------------------------
// UnixMilli
// ---------------------------------------------------------------------------

// ExampleUnixMilli demonstrates converting between Unix millisecond timestamps
// and time.Time values.
func ExampleUnixMilli() {
	iso := UnixMilli()

	// int64 milliseconds → time.Time
	fmt.Println(iso.Get(1609459200000).UTC().Format(time.RFC3339))

	// time.Time → int64 milliseconds
	fmt.Println(iso.ReverseGet(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
	// Output:
	// 2021-01-01T00:00:00Z
	// 1609459200000
}

// ExampleUnixMilli_roundTrip demonstrates the round-trip laws and the
// millisecond precision limit.
func ExampleUnixMilli_roundTrip() {
	iso := UnixMilli()

	// ReverseGet(Get(ms)) == ms always holds
	fmt.Println(iso.ReverseGet(iso.Get(1234567890000)))

	// Get(ReverseGet(t)) == t holds when t carries no sub-millisecond part
	exact := time.Date(2021, 1, 1, 12, 30, 45, 123000000, time.UTC)
	fmt.Println(iso.Get(iso.ReverseGet(exact)).Equal(exact))

	// Sub-millisecond precision is truncated by the round trip
	precise := time.Date(2021, 1, 1, 12, 30, 45, 123456789, time.UTC)
	fmt.Println(iso.Get(iso.ReverseGet(precise)).Nanosecond())
	// Output:
	// 1234567890000
	// true
	// 123000000
}

// ---------------------------------------------------------------------------
// Add / Sub
// ---------------------------------------------------------------------------

// ExampleAdd demonstrates shifting a number by a constant offset.
func ExampleAdd() {
	translateX := Add(100)

	fmt.Println(translateX.Get(50))
	fmt.Println(translateX.ReverseGet(150))
	// Output:
	// 150
	// 50
}

// ExampleAdd_float demonstrates that Add works for any numeric type.
func ExampleAdd_float() {
	iso := Add(0.5)

	fmt.Println(iso.Get(1.0))
	fmt.Println(iso.ReverseGet(1.5))
	// Output:
	// 1.5
	// 1
}

// ExampleAdd_roundTrip demonstrates the round-trip laws: adding then
// subtracting the offset is the identity in both directions.
func ExampleAdd_roundTrip() {
	iso := Add(7)

	fmt.Println(iso.ReverseGet(iso.Get(35)))
	fmt.Println(iso.Get(iso.ReverseGet(35)))
	// Output:
	// 35
	// 35
}

// ExampleAdd_modify demonstrates computing in a shifted coordinate space:
// Modify moves into 1-based indices, doubles, and shifts back to 0-based.
func ExampleAdd_modify() {
	shiftedDouble := F.Pipe1(Add(1), Modify[int](N.Mul(2)))

	fmt.Println(shiftedDouble(0)) // 0 → 1 → 2 → 1
	fmt.Println(shiftedDouble(4)) // 4 → 5 → 10 → 9
	// Output:
	// 1
	// 9
}

// ExampleSub demonstrates shifting a number by a constant offset in the
// opposite direction of Add.
func ExampleSub() {
	iso := Sub(5)

	fmt.Println(iso.Get(10))
	fmt.Println(iso.ReverseGet(5))
	// Output:
	// 5
	// 10
}

// ExampleSub_roundTrip demonstrates the round-trip laws for Sub.
func ExampleSub_roundTrip() {
	iso := Sub(3)

	fmt.Println(iso.ReverseGet(iso.Get(20)))
	fmt.Println(iso.Get(iso.ReverseGet(20)))
	// Output:
	// 20
	// 20
}

// ExampleSub_equivalentToAdd demonstrates that Sub(n) and Add(-n) describe the
// same isomorphism.
func ExampleSub_equivalentToAdd() {
	sub5 := Sub(5)
	addNeg5 := Add(-5)

	fmt.Println(sub5.Get(10) == addNeg5.Get(10))
	fmt.Println(sub5.ReverseGet(5) == addNeg5.ReverseGet(5))
	// Output:
	// true
	// true
}

// ---------------------------------------------------------------------------
// SwapPair
// ---------------------------------------------------------------------------

// ExampleSwapPair demonstrates exchanging the two elements of a Pair.
func ExampleSwapPair() {
	iso := SwapPair[string, int]()

	fmt.Println(iso.Get(pair.MakePair("hello", 42)))
	fmt.Println(iso.ReverseGet(pair.MakePair(7, "world")))
	// Output:
	// Pair[int, string](42, hello)
	// Pair[string, int](world, 7)
}

// ExampleSwapPair_roundTrip demonstrates that SwapPair is self-inverse:
// swapping twice returns the original pair.
func ExampleSwapPair_roundTrip() {
	iso := SwapPair[string, int]()

	fmt.Println(iso.ReverseGet(iso.Get(pair.MakePair("hello", 42))))
	fmt.Println(iso.Get(iso.ReverseGet(pair.MakePair(1, "one"))))
	// Output:
	// Pair[string, int](hello, 42)
	// Pair[int, string](1, one)
}

// ExampleSwapPair_modify demonstrates editing the *first* element of a pair
// through a function that only knows how to edit the second one.
func ExampleSwapPair_modify() {
	// Modify works on Pair[int, string]; SwapPair lifts it to Pair[string, int].
	upperHead := F.Pipe1(
		SwapPair[string, int](),
		Modify[Pair[string, int]](pair.Map[int](strings.ToUpper)),
	)

	fmt.Println(upperHead(pair.MakePair("hello", 42)))
	// Output:
	// Pair[string, int](HELLO, 42)
}

// ---------------------------------------------------------------------------
// SwapEither
// ---------------------------------------------------------------------------

// ExampleSwapEither demonstrates exchanging the Left and Right positions of an
// Either.  A value on the Right moves to the Left and vice versa.
func ExampleSwapEither() {
	iso := SwapEither[string, int]()

	// Right(42) becomes Left(42)
	fmt.Println(iso.Get(either.Right[string](42)))

	// Left("boom") becomes Right("boom")
	fmt.Println(iso.Get(either.Left[int]("boom")))

	// ReverseGet swaps back
	fmt.Println(iso.ReverseGet(either.Right[int]("recovered")))
	// Output:
	// Left[int](42)
	// Right[string](boom)
	// Left[string](recovered)
}

// ExampleSwapEither_roundTrip demonstrates that SwapEither is self-inverse.
func ExampleSwapEither_roundTrip() {
	iso := SwapEither[string, int]()

	fmt.Println(iso.ReverseGet(iso.Get(either.Right[string](42))))
	fmt.Println(iso.ReverseGet(iso.Get(either.Left[int]("boom"))))
	// Output:
	// Right[int](42)
	// Left[string](boom)
}

// ---------------------------------------------------------------------------
// ReverseArray
// ---------------------------------------------------------------------------

// ExampleReverseArray demonstrates reversing the order of a slice in both
// directions.
func ExampleReverseArray() {
	iso := ReverseArray[int]()

	fmt.Println(iso.Get([]int{1, 2, 3, 4}))
	fmt.Println(iso.ReverseGet([]int{4, 3, 2, 1}))
	// Output:
	// [4 3 2 1]
	// [1 2 3 4]
}

// ExampleReverseArray_roundTrip demonstrates that ReverseArray is self-inverse:
// Get and ReverseGet are the same function, so applying either twice is the
// identity.
func ExampleReverseArray_roundTrip() {
	iso := ReverseArray[string]()

	fmt.Println(iso.ReverseGet(iso.Get([]string{"a", "b", "c"})))
	fmt.Println(iso.Get(iso.Get([]string{"a", "b", "c"})))
	// Output:
	// [a b c]
	// [a b c]
}

// ExampleReverseArray_modify demonstrates appending to a slice by prepending in
// reversed space — a head-only operation reused as a tail operation.
func ExampleReverseArray_modify() {
	appendZero := F.Pipe1(ReverseArray[int](), Modify[[]int](A.Prepend(0)))

	fmt.Println(appendZero([]int{1, 2, 3}))
	// Output:
	// [1 2 3 0]
}

// ---------------------------------------------------------------------------
// Head
// ---------------------------------------------------------------------------

// ExampleHead demonstrates wrapping a value into a singleton non-empty array
// and reading the head back out.
func ExampleHead() {
	iso := Head[string]()

	fmt.Println(iso.Get("solo"))
	fmt.Println(iso.ReverseGet(nonempty.From("first", "second")))
	// Output:
	// [solo]
	// first
}

// ExampleHead_roundTrip demonstrates the round-trip laws.  ReverseGet(Get(v))
// always returns v, while Get(ReverseGet(as)) only reconstructs as when as is a
// singleton — every element after the head is discarded.
func ExampleHead_roundTrip() {
	iso := Head[int]()

	fmt.Println(iso.ReverseGet(iso.Get(42)))
	fmt.Println(iso.Get(iso.ReverseGet(nonempty.Of(7))))
	fmt.Println(iso.Get(iso.ReverseGet(nonempty.From(1, 2, 3))))
	// Output:
	// 42
	// [7]
	// [1]
}

// ---------------------------------------------------------------------------
// FromStrictEquals / FromEquals
// ---------------------------------------------------------------------------

// ExampleFromStrictEquals demonstrates basic usage: mapping "no"/"yes" to false/true.
func ExampleFromStrictEquals() {
	iso := FromStrictEquals(lazy.Of("no"), lazy.Of("yes"))

	fmt.Println(iso.Get("yes"))
	fmt.Println(iso.Get("no"))
	fmt.Println(iso.ReverseGet(true))
	fmt.Println(iso.ReverseGet(false))
	// Output:
	// true
	// false
	// yes
	// no
}

// ExampleFromStrictEquals_integer demonstrates FromStrictEquals with integer sentinels.
func ExampleFromStrictEquals_integer() {
	iso := FromStrictEquals(lazy.Of(0), lazy.Of(1))

	fmt.Println(iso.Get(1))
	fmt.Println(iso.Get(0))
	fmt.Println(iso.ReverseGet(true))
	fmt.Println(iso.ReverseGet(false))
	// Output:
	// true
	// false
	// 1
	// 0
}

// ExampleFromStrictEquals_roundTrip demonstrates the round-trip isomorphism laws.
func ExampleFromStrictEquals_roundTrip() {
	iso := FromStrictEquals(lazy.Of(0), lazy.Of(1))

	// Law 1: ReverseGet(Get(sentinel)) == sentinel
	fmt.Println(iso.ReverseGet(iso.Get(1)))
	fmt.Println(iso.ReverseGet(iso.Get(0)))

	// Law 2: Get(ReverseGet(bool)) == bool
	fmt.Println(iso.Get(iso.ReverseGet(true)))
	fmt.Println(iso.Get(iso.ReverseGet(false)))
	// Output:
	// 1
	// 0
	// true
	// false
}

// ExampleFromStrictEquals_modify demonstrates using FromStrictEquals with Modify
// to toggle a sentinel value through bool space.
func ExampleFromStrictEquals_modify() {
	toggle := F.Pipe1(
		FromStrictEquals(lazy.Of("inactive"), lazy.Of("active")),
		Modify[string](boolean.Not),
	)

	fmt.Println(toggle("active"))
	fmt.Println(toggle("inactive"))
	// Output:
	// inactive
	// active
}

// ExampleFromEquals demonstrates FromEquals with a custom case-insensitive Eq.
func ExampleFromEquals() {
	caseInsensitiveEq := eq.FromEquals(strings.EqualFold)

	iso := FromEquals(lazy.Of("NO"), lazy.Of("YES"))(caseInsensitiveEq)

	// Any casing of "YES" maps to true
	fmt.Println(iso.Get("yes"))
	fmt.Println(iso.Get("YES"))
	fmt.Println(iso.Get("Yes"))

	// Any casing of "NO" maps to false
	fmt.Println(iso.Get("no"))

	// ReverseGet always returns the canonical sentinel
	fmt.Println(iso.ReverseGet(true))
	fmt.Println(iso.ReverseGet(false))
	// Output:
	// true
	// true
	// true
	// false
	// YES
	// NO
}

// ExampleFromEquals_roundTrip demonstrates round-trip laws for FromEquals.
func ExampleFromEquals_roundTrip() {
	iso := FromEquals(lazy.Of("off"), lazy.Of("on"))(eq.FromStrictEquals[string]())

	// Get(ReverseGet(b)) == b
	fmt.Println(iso.Get(iso.ReverseGet(true)))
	fmt.Println(iso.Get(iso.ReverseGet(false)))
	// Output:
	// true
	// false
}

// ExampleNegate demonstrates basic Negate usage: NOT in both directions.
func ExampleNegate() {
	iso := Negate()

	fmt.Println(iso.Get(true))
	fmt.Println(iso.Get(false))
	fmt.Println(iso.ReverseGet(true))
	fmt.Println(iso.ReverseGet(false))
	// Output:
	// false
	// true
	// false
	// true
}

// ExampleNegate_roundTrip demonstrates the involution law: double negation is identity.
func ExampleNegate_roundTrip() {
	iso := Negate()

	fmt.Println(iso.Get(iso.Get(true)))
	fmt.Println(iso.Get(iso.Get(false)))
	// Output:
	// true
	// false
}

// ExampleNegate_modify demonstrates using Negate with Modify to double-negate a bool
// that was lifted from a string sentinel.
func ExampleNegate_modify() {
	// Map "disabled"/"enabled" → bool, then negate to flip the meaning
	toggle := F.Pipe1(
		FromStrictEquals(lazy.Of("disabled"), lazy.Of("enabled")),
		Modify[string](Negate().Get),
	)

	fmt.Println(toggle("enabled"))
	fmt.Println(toggle("disabled"))
	// Output:
	// disabled
	// enabled
}

// ExampleRef demonstrates basic Ref usage: wrapping a value into a pointer,
// and returning the lazy fallback when ReverseGet receives nil.
func ExampleRef() {
	iso := Ref(lazy.Of(0))

	p := iso.Get(42)
	fmt.Println(*p)
	fmt.Println(iso.ReverseGet(p))
	fmt.Println(iso.ReverseGet(nil))
	// Output:
	// 42
	// 42
	// 0
}

// ExampleRef_roundTrip demonstrates the round-trip isomorphism laws for Ref.
func ExampleRef_roundTrip() {
	iso := Ref(lazy.Of(""))

	// ReverseGet(Get(v)) == v
	fmt.Println(iso.ReverseGet(iso.Get("hello")))

	// *Get(ReverseGet(p)) == *p
	original := "world"
	fmt.Println(*iso.Get(iso.ReverseGet(&original)))
	// Output:
	// hello
	// world
}

// ExampleDeref demonstrates basic Deref usage: dereferencing a pointer to its value,
// and returning the lazy fallback when Get receives nil.
func ExampleDeref() {
	iso := Deref(lazy.Of(0))

	v := 42
	fmt.Println(iso.Get(&v))
	fmt.Println(iso.Get(nil))

	p := iso.ReverseGet(99)
	fmt.Println(*p)
	// Output:
	// 42
	// 0
	// 99
}

// ExampleDeref_roundTrip demonstrates the round-trip isomorphism laws for Deref.
func ExampleDeref_roundTrip() {
	iso := Deref(lazy.Of(""))

	// Get(ReverseGet(v)) == v
	fmt.Println(iso.Get(iso.ReverseGet("hello")))

	// *ReverseGet(Get(p)) == *p
	original := "world"
	fmt.Println(*iso.ReverseGet(iso.Get(&original)))
	// Output:
	// hello
	// world
}

// ExampleToNillable demonstrates basic ToNillable usage: converting between
// Option[A] and nullable *A in both directions.
func ExampleToNillable() {
	iso := ToNillable[int]()

	// Some(v) → non-nil pointer
	p := iso.Get(O.Some(42))
	fmt.Println(*p)

	// None → nil
	fmt.Println(iso.Get(O.None[int]()) == nil)

	// non-nil pointer → Some
	v := 99
	fmt.Println(iso.ReverseGet(&v))

	// nil → None
	fmt.Println(iso.ReverseGet(nil))
	// Output:
	// 42
	// true
	// Some[int](99)
	// None[int]
}

// ExampleToNillable_roundTrip demonstrates the round-trip isomorphism laws for ToNillable.
func ExampleToNillable_roundTrip() {
	iso := ToNillable[int]()

	// ReverseGet(Get(Some(v))) == Some(v)
	fmt.Println(iso.ReverseGet(iso.Get(O.Some(7))))

	// ReverseGet(Get(None)) == None
	fmt.Println(iso.ReverseGet(iso.Get(O.None[int]())))
	// Output:
	// Some[int](7)
	// None[int]
}

// ExampleFromNillable demonstrates basic FromNillable usage: converting between
// nullable *A and Option[A] in both directions.
func ExampleFromNillable() {
	iso := FromNillable[int]()

	// non-nil pointer → Some
	v := 42
	fmt.Println(iso.Get(&v))

	// nil → None
	fmt.Println(iso.Get(nil))

	// Some(v) → non-nil pointer
	p := iso.ReverseGet(O.Some(99))
	fmt.Println(*p)

	// None → nil
	fmt.Println(iso.ReverseGet(O.None[int]()) == nil)
	// Output:
	// Some[int](42)
	// None[int]
	// 99
	// true
}

// ExampleFromNillable_roundTrip demonstrates the round-trip isomorphism laws for FromNillable.
func ExampleFromNillable_roundTrip() {
	iso := FromNillable[int]()

	// Get(ReverseGet(Some(v))) == Some(v)
	fmt.Println(iso.Get(iso.ReverseGet(O.Some(7))))

	// Get(ReverseGet(None)) == None
	fmt.Println(iso.Get(iso.ReverseGet(O.None[int]())))
	// Output:
	// Some[int](7)
	// None[int]
}

// ExampleNanoDuration demonstrates basic NanoDuration usage: converting between
// raw nanosecond counts and time.Duration in both directions.
func ExampleNanoDuration() {
	iso := NanoDuration()

	// int64 nanoseconds → time.Duration
	fmt.Println(iso.Get(int64(2 * time.Second)))
	fmt.Println(iso.Get(int64(500 * time.Millisecond)))

	// time.Duration → int64 nanoseconds
	fmt.Println(iso.ReverseGet(time.Minute))
	// Output:
	// 2s
	// 500ms
	// 60000000000
}

// ExampleNanoDuration_roundTrip demonstrates the lossless round-trip laws for NanoDuration.
func ExampleNanoDuration_roundTrip() {
	iso := NanoDuration()

	// ReverseGet(Get(n)) == n
	n := int64(3 * time.Second)
	fmt.Println(iso.ReverseGet(iso.Get(n)) == n)

	// Get(ReverseGet(d)) == d
	d := 750 * time.Millisecond
	fmt.Println(iso.Get(iso.ReverseGet(d)) == d)
	// Output:
	// true
	// true
}

// ExampleSecondsDuration demonstrates basic SecondsDuration usage: converting between
// whole-second counts and time.Duration in both directions.
func ExampleSecondsDuration() {
	iso := SecondsDuration()

	// int seconds → time.Duration
	fmt.Println(iso.Get(5))
	fmt.Println(iso.Get(0))
	fmt.Println(iso.Get(-2))

	// time.Duration → int (whole seconds, sub-second part truncated)
	fmt.Println(iso.ReverseGet(90 * time.Second))
	fmt.Println(iso.ReverseGet(1500 * time.Millisecond)) // 1.5 s → 1
	// Output:
	// 5s
	// 0s
	// -2s
	// 90
	// 1
}

// ExampleSecondsDuration_roundTrip demonstrates the round-trip laws for SecondsDuration.
func ExampleSecondsDuration_roundTrip() {
	iso := SecondsDuration()

	// ReverseGet(Get(n)) == n always holds
	fmt.Println(iso.ReverseGet(iso.Get(42)))

	// Get(ReverseGet(d)) == d holds only when d is an exact multiple of time.Second
	exact := 10 * time.Second
	fmt.Println(iso.Get(iso.ReverseGet(exact)) == exact)
	// Output:
	// 42
	// true
}

// ---------------------------------------------------------------------------
// EmptyToNil
// ---------------------------------------------------------------------------

// ExampleEmptyToNil demonstrates the EmptyToNil isomorphism: both Get and
// ReverseGet normalise an empty (or nil) slice to nil and leave non-empty
// slices unchanged.
func ExampleEmptyToNil() {
	iso := EmptyToNil[int]()

	// nil → nil (Get)
	fmt.Println(iso.Get(nil) == nil)

	// empty non-nil → nil (Get)
	fmt.Println(iso.Get([]int{}) == nil)

	// non-empty → unchanged (Get)
	fmt.Println(iso.Get([]int{1, 2, 3}))

	// nil → nil (ReverseGet — same normalisation)
	fmt.Println(iso.ReverseGet(nil) == nil)

	// empty non-nil → nil (ReverseGet)
	fmt.Println(iso.ReverseGet([]int{}) == nil)

	// non-empty → unchanged (ReverseGet)
	fmt.Println(iso.ReverseGet([]int{4, 5}))
	// Output:
	// true
	// true
	// [1 2 3]
	// true
	// true
	// [4 5]
}

// ExampleEmptyToNil_roundTrip demonstrates the round-trip laws for EmptyToNil.
// Note that []int{} and nil are identified by this iso.
func ExampleEmptyToNil_roundTrip() {
	iso := EmptyToNil[int]()

	// Get(ReverseGet(s)) == s for non-empty s
	fmt.Println(iso.Get(iso.ReverseGet([]int{1, 2, 3})))

	// ReverseGet(Get(s)) == s for non-empty s
	fmt.Println(iso.ReverseGet(iso.Get([]int{7, 8})))

	// []int{} and nil are identified: both map to nil
	fmt.Println(iso.Get([]int{}) == nil)
	fmt.Println(iso.ReverseGet(iso.Get([]int{})) == nil)
	// Output:
	// [1 2 3]
	// [7 8]
	// true
	// true
}
