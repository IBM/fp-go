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
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestFindFirstWithIndex(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	finder := FindFirstWithIndex(func(i, x int) bool {
		return i > 2 && x%2 == 0
	})
	result := finder(src)
	assert.Equal(t, O.Some(4), result)

	notFound := FindFirstWithIndex(func(i, x int) bool {
		return i > 10
	})
	assert.Equal(t, O.None[int](), notFound(src))
}

func TestFindFirstMap(t *testing.T) {
	src := []string{"a", "42", "b", "100"}
	finder := FindFirstMap(func(s string) O.Option[int] {
		if len(s) > 1 {
			return O.Some(len(s))
		}
		return O.None[int]()
	})
	result := finder(src)
	assert.Equal(t, O.Some(2), result)
}

func TestFindFirstMapWithIndex(t *testing.T) {
	src := []string{"a", "b", "c", "d"}
	finder := FindFirstMapWithIndex(func(i int, s string) O.Option[string] {
		if i > 1 {
			return O.Some(fmt.Sprintf("%d:%s", i, s))
		}
		return O.None[string]()
	})
	result := finder(src)
	assert.Equal(t, O.Some("2:c"), result)
}

func TestFindLast(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	finder := FindLast(func(x int) bool { return x%2 == 0 })
	result := finder(src)
	assert.Equal(t, O.Some(4), result)

	notFound := FindLast(func(x int) bool { return x > 10 })
	assert.Equal(t, O.None[int](), notFound(src))
}

func TestFindLastWithIndex(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	finder := FindLastWithIndex(func(i, x int) bool {
		return i < 3 && x%2 == 0
	})
	result := finder(src)
	assert.Equal(t, O.Some(2), result)
}

func TestFindLastMap(t *testing.T) {
	src := []string{"a", "42", "b", "100"}
	finder := FindLastMap(func(s string) O.Option[int] {
		if len(s) > 1 {
			return O.Some(len(s))
		}
		return O.None[int]()
	})
	result := finder(src)
	assert.Equal(t, O.Some(3), result)
}

func TestFindLastMapWithIndex(t *testing.T) {
	src := []string{"a", "b", "c", "d"}
	finder := FindLastMapWithIndex(func(i int, s string) O.Option[string] {
		if i < 3 {
			return O.Some(fmt.Sprintf("%d:%s", i, s))
		}
		return O.None[string]()
	})
	result := finder(src)
	assert.Equal(t, O.Some("2:c"), result)
}
