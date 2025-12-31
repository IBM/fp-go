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

package string

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemigroup(t *testing.T) {
	// Test basic concatenation
	result := Semigroup.Concat("hello", " world")
	assert.Equal(t, "hello world", result)

	// Test associativity: (a • b) • c = a • (b • c)
	a, b, c := "foo", "bar", "baz"
	left := Semigroup.Concat(Semigroup.Concat(a, b), c)
	right := Semigroup.Concat(a, Semigroup.Concat(b, c))
	assert.Equal(t, left, right)
	assert.Equal(t, "foobarbaz", left)

	// Test with empty strings
	assert.Equal(t, "hello", Semigroup.Concat("", "hello"))
	assert.Equal(t, "hello", Semigroup.Concat("hello", ""))
}

func TestIntersperseSemigroup(t *testing.T) {
	// Test with comma separator
	commaSemigroup := IntersperseSemigroup(", ")
	result := commaSemigroup.Concat("a", "b")
	assert.Equal(t, "a, b", result)

	// Test associativity
	a, b, c := "x", "y", "z"
	left := commaSemigroup.Concat(commaSemigroup.Concat(a, b), c)
	right := commaSemigroup.Concat(a, commaSemigroup.Concat(b, c))
	assert.Equal(t, left, right)
	assert.Equal(t, "x, y, z", left)

	// Test with dash separator
	dashSemigroup := IntersperseSemigroup("-")
	assert.Equal(t, "foo-bar", dashSemigroup.Concat("foo", "bar"))

	// Test with empty strings - should not add separator (monoid identity)
	assert.Equal(t, "b", commaSemigroup.Concat("", "b"))
	assert.Equal(t, "a", commaSemigroup.Concat("a", ""))
}

func TestConcat(t *testing.T) {
	// Test the internal concat function
	result := concat("hello", " world")
	assert.Equal(t, "hello world", result)

	result = concat("foo", "bar")
	assert.Equal(t, "foobar", result)

	result = concat("", "test")
	assert.Equal(t, "test", result)
}
