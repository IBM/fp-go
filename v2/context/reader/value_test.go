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

package reader

import (
	"context"
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type valueTestKey string

const valueUserKey valueTestKey = "user"

func TestAskValue(t *testing.T) {
	t.Run("returns Some when the key holds a value of the expected type", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), valueUserKey, "Alice")
		assert.Equal(t, O.Some("Alice"), AskValue[string](valueUserKey)(ctx))
	})

	t.Run("returns None when the key is absent", func(t *testing.T) {
		assert.Equal(t, O.None[string](), AskValue[string](valueUserKey)(t.Context()))
	})

	t.Run("returns None when the stored value has a different type", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), valueUserKey, 42)
		assert.Equal(t, O.None[string](), AskValue[string](valueUserKey)(ctx))
	})

	t.Run("returns None for a stored nil interface", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), valueUserKey, nil)
		assert.Equal(t, O.None[error](), AskValue[error](valueUserKey)(ctx))
	})

	t.Run("round-trips with WithValue", func(t *testing.T) {
		ctx := WithValue[string](valueUserKey)("Alice")(t.Context())
		assert.Equal(t, O.Some("Alice"), AskValue[string](valueUserKey)(ctx))
	})
}
