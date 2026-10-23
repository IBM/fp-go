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

package readerio

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/lazy"
	"github.com/stretchr/testify/assert"
)

func TestFromLazy_Success(t *testing.T) {
	t.Run("evaluates lazily and ignores the context", func(t *testing.T) {
		calls := 0
		rio := FromLazy[ReaderTestConfig](func() int { calls++; return 42 })
		assert.Equal(t, 0, calls)
		assert.Equal(t, 42, rio(ReaderTestConfig{Value: 1})())
		assert.Equal(t, 42, rio(ReaderTestConfig{Value: 2})())
		assert.Equal(t, 2, calls)
	})
}

// ExampleFromLazy demonstrates lifting a lazy computation into a ReaderIO.
func ExampleFromLazy() {
	fmt.Println(FromLazy[ReaderTestConfig](lazy.Of(42))(ReaderTestConfig{})())
	// Output: 42
}
