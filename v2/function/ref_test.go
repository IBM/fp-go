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

package function

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRef_Extended covers cases not present in function_test.go.
func TestRef_Extended(t *testing.T) {
	t.Run("each call returns a distinct pointer", func(t *testing.T) {
		p1 := Ref(7)
		p2 := Ref(7)
		assert.NotSame(t, p1, p2)
	})

	t.Run("works with zero value", func(t *testing.T) {
		p := Ref(0)
		assert.Equal(t, 0, *p)
	})

	t.Run("works with bool", func(t *testing.T) {
		assert.True(t, *Ref(true))
		assert.False(t, *Ref(false))
	})
}

// TestDeref_Extended covers cases not present in function_test.go.
func TestDeref_Extended(t *testing.T) {
	t.Run("panics on nil pointer", func(t *testing.T) {
		assert.Panics(t, func() {
			var p *int
			Deref(p)
		})
	})
}

// TestDerefSafe tests the DerefSafe function.
func TestDerefSafe(t *testing.T) {
	t.Run("returns pointed-to int when non-nil", func(t *testing.T) {
		v := 42
		assert.Equal(t, 42, DerefSafe(Constant(0))(&v))
	})

	t.Run("returns zero fallback for nil int pointer", func(t *testing.T) {
		assert.Equal(t, 0, DerefSafe(Constant(0))(nil))
	})

	t.Run("returns pointed-to string when non-nil", func(t *testing.T) {
		s := "hello"
		assert.Equal(t, "hello", DerefSafe(Constant("default"))(&s))
	})

	t.Run("returns string fallback for nil pointer", func(t *testing.T) {
		assert.Equal(t, "default", DerefSafe(Constant("default"))(nil))
	})

	t.Run("fallback is independent of pointed-to value", func(t *testing.T) {
		v := 99
		safe := DerefSafe(Constant(-1))
		assert.Equal(t, 99, safe(&v))
		assert.Equal(t, -1, safe(nil))
	})

	t.Run("zero thunk is not called when pointer is non-nil", func(t *testing.T) {
		called := false
		zero := func() int { called = true; return -1 }
		v := 42
		DerefSafe(zero)(&v)
		assert.False(t, called)
	})
}

// ExampleRef demonstrates wrapping a value in a pointer.
func ExampleRef() {
	p := Ref(42)
	fmt.Println(*p)
	// Output:
	// 42
}

// ExampleDeref demonstrates dereferencing a pointer to recover its value.
func ExampleDeref() {
	v := 42
	fmt.Println(Deref(&v))
	// Output:
	// 42
}

// ExampleDeref_roundTrip demonstrates that Deref(Ref(v)) == v.
func ExampleDeref_roundTrip() {
	fmt.Println(Deref(Ref("hello")))
	// Output:
	// hello
}

// ExampleDerefSafe demonstrates safe dereferencing with a fallback thunk for nil.
func ExampleDerefSafe() {
	safe := DerefSafe(Constant(0))

	v := 42
	fmt.Println(safe(&v))
	fmt.Println(safe(nil))
	// Output:
	// 42
	// 0
}

// ExampleDerefSafe_string demonstrates DerefSafe with a string fallback thunk.
func ExampleDerefSafe_string() {
	safe := DerefSafe(Constant("(none)"))

	s := "hello"
	fmt.Println(safe(&s))
	fmt.Println(safe(nil))
	// Output:
	// hello
	// (none)
}
