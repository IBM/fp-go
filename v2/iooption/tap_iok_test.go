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

package iooption

import (
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestTapIOK(t *testing.T) {
	t.Run("runs the side effect for Some and keeps the value", func(t *testing.T) {
		var seen []int
		record := func(n int) io.IO[string] {
			return func() string { seen = append(seen, n); return "" }
		}
		assert.Equal(t, O.Some(1), TapIOK(record)(Some(1))())
		assert.Equal(t, O.Some(2), MonadTapIOK(Some(2), record)())
		assert.Equal(t, []int{1, 2}, seen)
	})

	t.Run("skips the side effect for None", func(t *testing.T) {
		calls := 0
		record := func(n int) io.IO[int] { return func() int { calls++; return n } }
		assert.Equal(t, O.None[int](), TapIOK(record)(None[int]())())
		assert.Equal(t, O.None[int](), MonadTapIOK(None[int](), record)())
		assert.Equal(t, 0, calls)
	})
}

// ExampleTapIOK runs a side effect for a Some value and keeps the value.
func ExampleTapIOK() {
	fmt.Println(F.Pipe1(Some(7), TapIOK(io.Printf[int]("log: %d\n")))())
	// Output:
	// log: 7
	// Some[int](7)
}

// ExampleMonadTapIOK is the non curried version of TapIOK.
func ExampleMonadTapIOK() {
	fmt.Println(MonadTapIOK(None[int](), io.Printf[int]("log: %d\n"))())
	// Output: None[int]
}
