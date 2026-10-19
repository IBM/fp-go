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

package readerreaderioresult

import (
	"context"
	"testing"

	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// contextAwareRelease returns a release action that fails if its context is cancelled
// and otherwise records that it ran.
func contextAwareRelease(released *bool) ReaderReaderIOResult[int, Void] {
	return reader.Of[int](F.Pipe1(
		RIOE.FromIO(func() Void { *released = true; return F.VOID }),
		RIOE.WithContext[Void],
	))
}

func TestUninterruptible_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	var released bool
	assert.Equal(t, result.Of(F.VOID), Uninterruptible(contextAwareRelease(&released))(0)(ctx)())
	assert.True(t, released)
}

func TestBracket_CancelledDuringUse(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	var released bool
	use := func(int) ReaderReaderIOResult[int, int] {
		return reader.Of[int](RIOE.FromIO(func() int { cancel(); return 2 }))
	}
	release := func(int, Result[int]) ReaderReaderIOResult[int, Void] { return contextAwareRelease(&released) }

	assert.Equal(t, result.Of(2), Bracket(Of[int](1), use, release)(0)(ctx)())
	assert.True(t, released)
}
