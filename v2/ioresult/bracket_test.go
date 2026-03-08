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

package ioresult

import (
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestBracket_Success(t *testing.T) {
	acquired := false
	used := false
	released := false

	acquire := func() IOResult[int] {
		return func() Result[int] {
			acquired = true
			return result.Of(42)
		}
	}()

	use := func(n int) IOResult[string] {
		return func() Result[string] {
			used = true
			return result.Of("success")
		}
	}

	release := func(n int, res Result[string]) IOResult[F.Void] {
		return func() Result[F.Void] {
			released = true
			return result.Of(F.VOID)
		}
	}

	res := Bracket(acquire, use, release)()

	assert.True(t, acquired, "Resource should be acquired")
	assert.True(t, used, "Resource should be used")
	assert.True(t, released, "Resource should be released")
	assert.Equal(t, result.Of("success"), res)
}

func TestBracket_UseFailure(t *testing.T) {
	acquired := false
	released := false
	releaseResult := result.Result[string]{}

	acquire := func() IOResult[int] {
		return func() Result[int] {
			acquired = true
			return result.Of(42)
		}
	}()

	useErr := errors.New("use error")
	use := func(n int) IOResult[string] {
		return func() Result[string] {
			return result.Left[string](useErr)
		}
	}

	release := func(n int, res Result[string]) IOResult[F.Void] {
		return func() Result[F.Void] {
			released = true
			releaseResult = res
			return result.Of(F.VOID)
		}
	}

	res := Bracket(acquire, use, release)()

	assert.True(t, acquired, "Resource should be acquired")
	assert.True(t, released, "Resource should be released even on use failure")
	assert.Equal(t, result.Left[string](useErr), res)
	assert.Equal(t, result.Left[string](useErr), releaseResult)
}

func TestBracket_AcquireFailure(t *testing.T) {
	used := false
	released := false

	acquireErr := errors.New("acquire error")
	acquire := func() IOResult[int] {
		return func() Result[int] {
			return result.Left[int](acquireErr)
		}
	}()

	use := func(n int) IOResult[string] {
		return func() Result[string] {
			used = true
			return result.Of("success")
		}
	}

	release := func(n int, res Result[string]) IOResult[F.Void] {
		return func() Result[F.Void] {
			released = true
			return result.Of(F.VOID)
		}
	}

	res := Bracket(acquire, use, release)()

	assert.False(t, used, "Use should not be called if acquire fails")
	assert.False(t, released, "Release should not be called if acquire fails")
	assert.Equal(t, result.Left[string](acquireErr), res)
}

func TestBracket_ReleaseFailure(t *testing.T) {
	acquired := false
	used := false
	released := false

	acquire := func() IOResult[int] {
		return func() Result[int] {
			acquired = true
			return result.Of(42)
		}
	}()

	use := func(n int) IOResult[string] {
		return func() Result[string] {
			used = true
			return result.Of("success")
		}
	}

	releaseErr := errors.New("release error")
	release := func(n int, res Result[string]) IOResult[F.Void] {
		return func() Result[F.Void] {
			released = true
			return result.Left[F.Void](releaseErr)
		}
	}

	res := Bracket(acquire, use, release)()

	assert.True(t, acquired, "Resource should be acquired")
	assert.True(t, used, "Resource should be used")
	assert.True(t, released, "Release should be attempted")
	// When release fails, the release error is returned
	assert.Equal(t, result.Left[string](releaseErr), res)
}

func TestBracket_BothUseAndReleaseFail(t *testing.T) {
	acquired := false
	released := false

	acquire := func() IOResult[int] {
		return func() Result[int] {
			acquired = true
			return result.Of(42)
		}
	}()

	useErr := errors.New("use error")
	use := func(n int) IOResult[string] {
		return func() Result[string] {
			return result.Left[string](useErr)
		}
	}

	releaseErr := errors.New("release error")
	release := func(n int, res Result[string]) IOResult[F.Void] {
		return func() Result[F.Void] {
			released = true
			return result.Left[F.Void](releaseErr)
		}
	}

	res := Bracket(acquire, use, release)()

	assert.True(t, acquired, "Resource should be acquired")
	assert.True(t, released, "Release should be attempted")
	// When both fail, the release error is returned
	assert.Equal(t, result.Left[string](releaseErr), res)
}

func TestBracket_ResourceValue(t *testing.T) {
	// Test that the acquired resource value is passed correctly
	var usedValue int
	var releasedValue int

	acquire := Of(100)

	use := func(n int) IOResult[string] {
		usedValue = n
		return Of("result")
	}

	release := func(n int, res Result[string]) IOResult[F.Void] {
		releasedValue = n
		return Of(F.VOID)
	}

	Bracket(acquire, use, release)()

	assert.Equal(t, 100, usedValue, "Use should receive acquired value")
	assert.Equal(t, 100, releasedValue, "Release should receive acquired value")
}

func TestBracket_ResultValue(t *testing.T) {
	// Test that the use result is passed to release
	var releaseReceivedResult Result[string]

	acquire := Of(42)

	use := func(n int) IOResult[string] {
		return Of("test result")
	}

	release := func(n int, res Result[string]) IOResult[F.Void] {
		releaseReceivedResult = res
		return Of(F.VOID)
	}

	Bracket(acquire, use, release)()

	assert.Equal(t, result.Of("test result"), releaseReceivedResult)
}


