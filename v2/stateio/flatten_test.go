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

package stateio

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

// incCounter increments the counter and returns the new value.
func incCounter(s TestState) IO[Pair[TestState, int]] {
	return func() Pair[TestState, int] {
		s.Counter++
		return pair.MakePair(s, s.Counter)
	}
}

func TestFlatten_Success(t *testing.T) {
	t.Run("runs the inner computation", func(t *testing.T) {
		res := Flatten(Of[TestState](Of[TestState](42)))(TestState{Counter: 1})()
		assert.Equal(t, pair.MakePair(TestState{Counter: 1}, 42), res)
	})

	t.Run("threads the state from the outer to the inner computation", func(t *testing.T) {
		outer := MonadMap(StateIO[TestState, int](incCounter), func(int) StateIO[TestState, int] { return incCounter })
		res := Flatten(outer)(TestState{Counter: 0})()
		assert.Equal(t, pair.MakePair(TestState{Counter: 2}, 2), res)
	})
}

// ExampleFlatten demonstrates removing one level of nesting from a StateIO.
func ExampleFlatten() {
	fmt.Println(Flatten(Of[TestState](Of[TestState]("done")))(TestState{Counter: 3})())
	// Output: Pair[stateio.TestState, string]({3 }, done)
}
