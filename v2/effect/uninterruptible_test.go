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

package effect

import (
	"context"
	"testing"

	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// cancellingUse returns a use function that cancels the context while it runs.
func cancellingUse(cancel context.CancelFunc) Kleisli[int, int, int] {
	return func(int) Effect[int, int] {
		return FromThunk[int](thunk.FromIO(func() int { cancel(); return 2 }))
	}
}

// contextAwareRelease returns a release effect that fails if its context is cancelled
// and otherwise records that it ran.
func contextAwareRelease(released *bool) Effect[int, Void] {
	return FromThunk[int](F.Pipe1(
		thunk.FromIO(func() Void { *released = true; return F.VOID }),
		thunk.WithContext[Void],
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
	release := func(int, Result[int]) Effect[int, Void] { return contextAwareRelease(&released) }

	assert.Equal(t, result.Of(2), Bracket(Of[int](1), cancellingUse(cancel), release)(0)(ctx)())
	assert.True(t, released)
}

func TestWithResource_CancelledDuringUse(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	var released bool
	withRes := WithResource[int](Of[int](1), func(int) Effect[int, Void] { return contextAwareRelease(&released) })

	assert.Equal(t, result.Of(2), withRes(cancellingUse(cancel))(0)(ctx)())
	assert.True(t, released)
}

func TestWithLock_CancelledDuringUse(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	var unlocked bool
	lock := Of[int](context.CancelFunc(func() { unlocked = true }))

	assert.Equal(t, result.Of(2), WithLock[int, int](lock)(cancellingUse(cancel)(1))(0)(ctx)())
	assert.True(t, unlocked)
}
