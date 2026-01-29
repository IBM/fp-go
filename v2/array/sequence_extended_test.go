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
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestSequenceWithOption tests the generic Sequence function with Option monad
func TestSequenceWithOption(t *testing.T) {
	// Test with Option monad - all Some values
	opts := From(
		O.Some(1),
		O.Some(2),
		O.Some(3),
	)

	// Use the Sequence function with Option's applicative monoid
	monoid := O.ApplicativeMonoid(Monoid[int]())
	seq := Sequence(O.Map(Of[int]), monoid)
	result := seq(opts)

	assert.Equal(t, O.Of(From(1, 2, 3)), result)

	// Test with Option monad - contains None
	optsWithNone := From(
		O.Some(1),
		O.None[int](),
		O.Some(3),
	)

	result2 := seq(optsWithNone)
	assert.True(t, O.IsNone(result2))

	// Test with empty array
	empty := Empty[Option[int]]()
	result3 := seq(empty)
	assert.Equal(t, O.Some(Empty[int]()), result3)
}

// TestMonadSequence tests the MonadSequence function
func TestMonadSequence(t *testing.T) {
	// Test with Option monad
	opts := From(
		O.Some("hello"),
		O.Some("world"),
	)

	monoid := O.ApplicativeMonoid(Monoid[string]())
	result := MonadSequence(O.Map(Of[string]), monoid, opts)

	assert.Equal(t, O.Of(From("hello", "world")), result)

	// Test with None in the array
	optsWithNone := From(
		O.Some("hello"),
		O.None[string](),
	)

	result2 := MonadSequence(O.Map(Of[string]), monoid, optsWithNone)
	assert.Equal(t, O.None[[]string](), result2)
}
