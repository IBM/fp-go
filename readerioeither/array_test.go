// Copyright (c) 2023 IBM Corp.
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

package readerioeither

import (
	"context"
	"testing"

	A "github.com/IBM/fp-go/array"
	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"
)

func TestTraverseArray(t *testing.T) {
	f := TraverseArray(func(a string) ReaderIOEither[context.Context, string, string] {
		if len(a) > 0 {
			return Right[context.Context, string](a + a)
		}
		return Left[context.Context, string, string]("e")
	})
	ctx := context.Background()
	assert.Equal(t, ET.Right[string](A.Empty[string]()), F.Pipe1(A.Empty[string](), f)(ctx)())
	assert.Equal(t, ET.Right[string]([]string{"aa", "bb"}), F.Pipe1([]string{"a", "b"}, f)(ctx)())
	assert.Equal(t, ET.Left[[]string]("e"), F.Pipe1([]string{"a", ""}, f)(ctx)())
}
