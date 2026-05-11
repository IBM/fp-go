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

package itereither

import (
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/iterator/iter"
	"github.com/stretchr/testify/assert"
)

func TestBracket_Success(t *testing.T) {
	var released bool
	var useCalled bool

	acquire := iter.From(E.Right[error]("resource"))
	use := func(r string) SeqEither[error, int] {
		useCalled = true
		return iter.From(E.Right[error](len(r)))
	}
	release := func(r string, result Either[error, int]) SeqEither[error, F.Void] {
		released = true
		assert.Equal(t, "resource", r)
		assert.True(t, E.IsRight(result))
		return iter.From(E.Right[error](F.VOID))
	}

	result := collectEithers(Bracket(acquire, use, release))
	assert.Equal(t, []Either[error, int]{E.Right[error](8)}, result)
	assert.True(t, useCalled)
	assert.True(t, released)
}

func TestBracket_UseFailure(t *testing.T) {
	var released bool
	var releaseResult Either[error, int]

	acquire := iter.From(E.Right[error]("resource"))
	use := func(r string) SeqEither[error, int] {
		return iter.From(E.Left[int](errors.New("use failed")))
	}
	release := func(r string, result Either[error, int]) SeqEither[error, F.Void] {
		released = true
		releaseResult = result
		return iter.From(E.Right[error](F.VOID))
	}

	result := collectEithers(Bracket(acquire, use, release))
	assert.Len(t, result, 1)
	assert.True(t, E.IsLeft(result[0]))
	assert.True(t, released)
	assert.True(t, E.IsLeft(releaseResult))
}

func TestBracket_AcquireFailure(t *testing.T) {
	var useCalled bool
	var released bool

	acquire := iter.From(E.Left[string](errors.New("acquire failed")))
	use := func(r string) SeqEither[error, int] {
		useCalled = true
		return iter.From(E.Right[error](42))
	}
	release := func(r string, result Either[error, int]) SeqEither[error, F.Void] {
		released = true
		return iter.From(E.Right[error](F.VOID))
	}

	result := collectEithers(Bracket(acquire, use, release))
	assert.Len(t, result, 1)
	assert.True(t, E.IsLeft(result[0]))
	assert.False(t, useCalled)
	assert.False(t, released)
}

func TestBracket_ReleaseFailure(t *testing.T) {
	var released bool

	acquire := iter.From(E.Right[error]("resource"))
	use := func(r string) SeqEither[error, int] {
		return iter.From(E.Right[error](42))
	}
	release := func(r string, result Either[error, int]) SeqEither[error, F.Void] {
		released = true
		return iter.From(E.Left[F.Void](errors.New("release failed")))
	}

	result := collectEithers(Bracket(acquire, use, release))
	// Release failure propagates as an error
	assert.Len(t, result, 1)
	assert.True(t, E.IsLeft(result[0]))
	assert.True(t, released)
}

func TestBracket_MultipleResources(t *testing.T) {
	var released1, released2 bool

	acquire1 := iter.From(E.Right[error]("resource1"))
	acquire2 := iter.From(E.Right[error]("resource2"))

	use1 := func(r1 string) SeqEither[error, string] {
		return Bracket(
			acquire2,
			func(r2 string) SeqEither[error, string] {
				return iter.From(E.Right[error](r1 + "+" + r2))
			},
			func(r2 string, result Either[error, string]) SeqEither[error, F.Void] {
				released2 = true
				return iter.From(E.Right[error](F.VOID))
			},
		)
	}

	release1 := func(r1 string, result Either[error, string]) SeqEither[error, F.Void] {
		released1 = true
		return iter.From(E.Right[error](F.VOID))
	}

	result := collectEithers(Bracket(acquire1, use1, release1))
	assert.Equal(t, []Either[error, string]{E.Right[error]("resource1+resource2")}, result)
	assert.True(t, released1)
	assert.True(t, released2)
}

func TestBracket_ReleaseCalledOnError(t *testing.T) {
	var released bool
	var releaseError Either[error, int]

	acquire := iter.From(E.Right[error]("resource"))
	use := func(r string) SeqEither[error, int] {
		return iter.From(E.Left[int](errors.New("processing error")))
	}
	release := func(r string, result Either[error, int]) SeqEither[error, F.Void] {
		released = true
		releaseError = result
		return iter.From(E.Right[error](F.VOID))
	}

	result := collectEithers(Bracket(acquire, use, release))
	assert.Len(t, result, 1)
	assert.True(t, E.IsLeft(result[0]))
	assert.True(t, released)
	assert.True(t, E.IsLeft(releaseError))
}

func TestBracket_ResourcePassedToRelease(t *testing.T) {
	var releasedResource string

	acquire := iter.From(E.Right[error]("test-resource"))
	use := func(r string) SeqEither[error, int] {
		return iter.From(E.Right[error](42))
	}
	release := func(r string, result Either[error, int]) SeqEither[error, F.Void] {
		releasedResource = r
		return iter.From(E.Right[error](F.VOID))
	}

	collectEithers(Bracket(acquire, use, release))
	assert.Equal(t, "test-resource", releasedResource)
}

func TestBracket_EmptySequence(t *testing.T) {
	var released bool

	acquire := iter.From[Either[error, string]]()
	use := func(r string) SeqEither[error, int] {
		return iter.From(E.Right[error](42))
	}
	release := func(r string, result Either[error, int]) SeqEither[error, F.Void] {
		released = true
		return iter.From(E.Right[error](F.VOID))
	}

	result := collectEithers(Bracket(acquire, use, release))
	assert.Empty(t, result)
	assert.False(t, released)
}

func TestBracket_FileHandlePattern(t *testing.T) {
	type FileHandle struct {
		Name   string
		Closed bool
	}

	var handle *FileHandle

	acquire := iter.From(E.Right[error](&FileHandle{Name: "test.txt", Closed: false}))
	use := func(fh *FileHandle) SeqEither[error, string] {
		handle = fh
		return iter.From(E.Right[error]("file contents"))
	}
	release := func(fh *FileHandle, result Either[error, string]) SeqEither[error, F.Void] {
		fh.Closed = true
		return iter.From(E.Right[error](F.VOID))
	}

	result := collectEithers(Bracket(acquire, use, release))
	assert.Equal(t, []Either[error, string]{E.Right[error]("file contents")}, result)
	assert.NotNil(t, handle)
	assert.True(t, handle.Closed)
}

// Made with Bob
