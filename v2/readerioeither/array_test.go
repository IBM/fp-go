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
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	TST "github.com/IBM/fp-go/v2/internal/testing"
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

func TestSequenceArray(t *testing.T) {

	s := TST.SequenceArrayTest(
		FromStrictEquals[context.Context, error, bool]()(context.Background()),
		Pointed[context.Context, error, string](),
		Pointed[context.Context, error, bool](),
		Functor[context.Context, error, []string, bool](),
		SequenceArray[context.Context, error, string],
	)

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("TestSequenceArray %d", i), s(i))
	}
}

func TestSequenceArrayError(t *testing.T) {

	s := TST.SequenceArrayErrorTest(
		FromStrictEquals[context.Context, error, bool]()(context.Background()),
		Left[context.Context, string, error],
		Left[context.Context, bool, error],
		Pointed[context.Context, error, string](),
		Pointed[context.Context, error, bool](),
		Functor[context.Context, error, []string, bool](),
		SequenceArray[context.Context, error, string],
	)
	// run across four bits
	s(4)(t)
}
