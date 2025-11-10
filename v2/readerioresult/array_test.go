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

package readerioresult

import (
	"context"
	"errors"
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	TST "github.com/IBM/fp-go/v2/internal/testing"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestTraverseArray(t *testing.T) {

	e := errors.New("e")

	f := TraverseArray(func(a string) ReaderIOResult[context.Context, string] {
		if len(a) > 0 {
			return Right[context.Context](a + a)
		}
		return Left[context.Context, string](e)
	})
	ctx := context.Background()
	assert.Equal(t, result.Of(A.Empty[string]()), F.Pipe1(A.Empty[string](), f)(ctx)())
	assert.Equal(t, result.Of([]string{"aa", "bb"}), F.Pipe1([]string{"a", "b"}, f)(ctx)())
	assert.Equal(t, result.Left[[]string](e), F.Pipe1([]string{"a", ""}, f)(ctx)())
}

func TestSequenceArray(t *testing.T) {

	s := TST.SequenceArrayTest(
		FromStrictEquals[context.Context, bool]()(context.Background()),
		Pointed[context.Context, string](),
		Pointed[context.Context, bool](),
		Functor[context.Context, []string, bool](),
		SequenceArray[context.Context, string],
	)

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("TestSequenceArray %d", i), s(i))
	}
}

func TestSequenceArrayError(t *testing.T) {

	s := TST.SequenceArrayErrorTest(
		FromStrictEquals[context.Context, bool]()(context.Background()),
		Left[context.Context],
		Left[context.Context, bool],
		Pointed[context.Context, string](),
		Pointed[context.Context, bool](),
		Functor[context.Context, []string, bool](),
		SequenceArray[context.Context, string],
	)
	// run across four bits
	s(4)(t)
}
