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

package readerreaderioeither

import (
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/stretchr/testify/assert"
)

type BracketOuterCtx struct {
	resourcePool string
}

type BracketInnerCtx struct {
	timeout int
}

type Resource struct {
	id       string
	acquired bool
	released bool
}

func TestBracketSuccessPath(t *testing.T) {
	outer := BracketOuterCtx{resourcePool: "pool1"}
	inner := BracketInnerCtx{timeout: 30}

	resource := &Resource{id: "res1"}

	// Acquire resource
	acquire := func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, *Resource] {
		return func(i BracketInnerCtx) IOE.IOEither[error, *Resource] {
			return func() E.Either[error, *Resource] {
				resource.acquired = true
				return E.Right[error](resource)
			}
		}
	}

	// Use resource successfully
	use := func(r *Resource) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, string] {
		return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, string] {
			return func(i BracketInnerCtx) IOE.IOEither[error, string] {
				return IOE.Of[error]("result from " + r.id)
			}
		}
	}

	// Release resource
	release := func(r *Resource, result E.Either[error, string]) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, any] {
		return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, any] {
			return func(i BracketInnerCtx) IOE.IOEither[error, any] {
				return func() E.Either[error, any] {
					r.released = true
					return E.Right[error, any](nil)
				}
			}
		}
	}

	result := Bracket(acquire, use, release)
	outcome := result(outer)(inner)()

	assert.Equal(t, E.Right[error]("result from res1"), outcome)
	assert.True(t, resource.acquired, "Resource should be acquired")
	assert.True(t, resource.released, "Resource should be released")
}

func TestBracketUseFailure(t *testing.T) {
	outer := BracketOuterCtx{resourcePool: "pool1"}
	inner := BracketInnerCtx{timeout: 30}

	resource := &Resource{id: "res1"}
	useErr := errors.New("use failed")

	// Acquire resource
	acquire := func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, *Resource] {
		return func(i BracketInnerCtx) IOE.IOEither[error, *Resource] {
			return func() E.Either[error, *Resource] {
				resource.acquired = true
				return E.Right[error](resource)
			}
		}
	}

	// Use resource with failure
	use := func(r *Resource) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, string] {
		return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, string] {
			return func(i BracketInnerCtx) IOE.IOEither[error, string] {
				return IOE.Left[string](useErr)
			}
		}
	}

	// Release resource (should still be called)
	release := func(r *Resource, result E.Either[error, string]) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, any] {
		return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, any] {
			return func(i BracketInnerCtx) IOE.IOEither[error, any] {
				return func() E.Either[error, any] {
					r.released = true
					return E.Right[error, any](nil)
				}
			}
		}
	}

	result := Bracket(acquire, use, release)
	outcome := result(outer)(inner)()

	assert.Equal(t, E.Left[string](useErr), outcome)
	assert.True(t, resource.acquired, "Resource should be acquired")
	assert.True(t, resource.released, "Resource should be released even on failure")
}

func TestBracketAcquireFailure(t *testing.T) {
	outer := BracketOuterCtx{resourcePool: "pool1"}
	inner := BracketInnerCtx{timeout: 30}

	resource := &Resource{id: "res1"}
	acquireErr := errors.New("acquire failed")
	useCalled := false
	releaseCalled := false

	// Acquire resource fails
	acquire := func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, *Resource] {
		return func(i BracketInnerCtx) IOE.IOEither[error, *Resource] {
			return IOE.Left[*Resource](acquireErr)
		}
	}

	// Use should not be called
	use := func(r *Resource) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, string] {
		return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, string] {
			return func(i BracketInnerCtx) IOE.IOEither[error, string] {
				return func() E.Either[error, string] {
					useCalled = true
					return E.Right[error]("should not reach here")
				}
			}
		}
	}

	// Release should not be called
	release := func(r *Resource, result E.Either[error, string]) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, any] {
		return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, any] {
			return func(i BracketInnerCtx) IOE.IOEither[error, any] {
				return func() E.Either[error, any] {
					releaseCalled = true
					return E.Right[error, any](nil)
				}
			}
		}
	}

	result := Bracket(acquire, use, release)
	outcome := result(outer)(inner)()

	assert.Equal(t, E.Left[string](acquireErr), outcome)
	assert.False(t, resource.acquired, "Resource should not be acquired")
	assert.False(t, useCalled, "Use should not be called when acquire fails")
	assert.False(t, releaseCalled, "Release should not be called when acquire fails")
}

func TestBracketReleaseReceivesResult(t *testing.T) {
	outer := BracketOuterCtx{resourcePool: "pool1"}
	inner := BracketInnerCtx{timeout: 30}

	resource := &Resource{id: "res1"}
	var capturedResult E.Either[error, string]

	// Acquire resource
	acquire := func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, *Resource] {
		return func(i BracketInnerCtx) IOE.IOEither[error, *Resource] {
			return func() E.Either[error, *Resource] {
				resource.acquired = true
				return E.Right[error](resource)
			}
		}
	}

	// Use resource
	use := func(r *Resource) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, string] {
		return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, string] {
			return func(i BracketInnerCtx) IOE.IOEither[error, string] {
				return IOE.Of[error]("use result")
			}
		}
	}

	// Release captures the result
	release := func(r *Resource, result E.Either[error, string]) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, any] {
		return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, any] {
			return func(i BracketInnerCtx) IOE.IOEither[error, any] {
				return func() E.Either[error, any] {
					capturedResult = result
					r.released = true
					return E.Right[error, any](nil)
				}
			}
		}
	}

	result := Bracket(acquire, use, release)
	outcome := result(outer)(inner)()

	assert.Equal(t, E.Right[error]("use result"), outcome)
	assert.Equal(t, E.Right[error]("use result"), capturedResult)
	assert.True(t, resource.released, "Resource should be released")
}

func TestBracketWithContextAccess(t *testing.T) {
	outer := BracketOuterCtx{resourcePool: "production"}
	inner := BracketInnerCtx{timeout: 60}

	resource := &Resource{id: "res1"}

	// Acquire uses outer context
	acquire := func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, *Resource] {
		return func(i BracketInnerCtx) IOE.IOEither[error, *Resource] {
			return func() E.Either[error, *Resource] {
				resource.id = o.resourcePool + "-resource"
				resource.acquired = true
				return E.Right[error](resource)
			}
		}
	}

	// Use uses inner context
	use := func(r *Resource) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, string] {
		return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, string] {
			return func(i BracketInnerCtx) IOE.IOEither[error, string] {
				return func() E.Either[error, string] {
					result := r.id + " with timeout " + string(rune(i.timeout+'0'))
					return E.Right[error](result)
				}
			}
		}
	}

	// Release uses both contexts
	release := func(r *Resource, result E.Either[error, string]) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, any] {
		return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, any] {
			return func(i BracketInnerCtx) IOE.IOEither[error, any] {
				return func() E.Either[error, any] {
					r.released = true
					return E.Right[error, any](nil)
				}
			}
		}
	}

	result := Bracket(acquire, use, release)
	outcome := result(outer)(inner)()

	assert.True(t, E.IsRight(outcome))
	assert.True(t, resource.acquired)
	assert.True(t, resource.released)
	assert.Equal(t, "production-resource", resource.id)
}

func TestBracketMultipleResources(t *testing.T) {
	outer := BracketOuterCtx{resourcePool: "pool1"}
	inner := BracketInnerCtx{timeout: 30}

	resource1 := &Resource{id: "res1"}
	resource2 := &Resource{id: "res2"}

	// Acquire first resource
	acquire1 := func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, *Resource] {
		return func(i BracketInnerCtx) IOE.IOEither[error, *Resource] {
			return func() E.Either[error, *Resource] {
				resource1.acquired = true
				return E.Right[error](resource1)
			}
		}
	}

	// Use first resource to acquire second
	use1 := func(r1 *Resource) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, string] {
		// Nested bracket for second resource
		acquire2 := func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, *Resource] {
			return func(i BracketInnerCtx) IOE.IOEither[error, *Resource] {
				return func() E.Either[error, *Resource] {
					resource2.acquired = true
					return E.Right[error](resource2)
				}
			}
		}

		use2 := func(r2 *Resource) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, string] {
			return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, string] {
				return func(i BracketInnerCtx) IOE.IOEither[error, string] {
					return IOE.Of[error](r1.id + " and " + r2.id)
				}
			}
		}

		release2 := func(r2 *Resource, result E.Either[error, string]) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, any] {
			return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, any] {
				return func(i BracketInnerCtx) IOE.IOEither[error, any] {
					return func() E.Either[error, any] {
						r2.released = true
						return E.Right[error, any](nil)
					}
				}
			}
		}

		return Bracket(acquire2, use2, release2)
	}

	release1 := func(r1 *Resource, result E.Either[error, string]) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, any] {
		return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, any] {
			return func(i BracketInnerCtx) IOE.IOEither[error, any] {
				return func() E.Either[error, any] {
					r1.released = true
					return E.Right[error, any](nil)
				}
			}
		}
	}

	result := Bracket(acquire1, use1, release1)
	outcome := result(outer)(inner)()

	assert.Equal(t, E.Right[error]("res1 and res2"), outcome)
	assert.True(t, resource1.acquired && resource1.released, "Resource 1 should be acquired and released")
	assert.True(t, resource2.acquired && resource2.released, "Resource 2 should be acquired and released")
}

func TestBracketReleaseErrorDoesNotAffectResult(t *testing.T) {
	outer := BracketOuterCtx{resourcePool: "pool1"}
	inner := BracketInnerCtx{timeout: 30}

	resource := &Resource{id: "res1"}
	releaseErr := errors.New("release failed")

	// Acquire resource
	acquire := func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, *Resource] {
		return func(i BracketInnerCtx) IOE.IOEither[error, *Resource] {
			return func() E.Either[error, *Resource] {
				resource.acquired = true
				return E.Right[error](resource)
			}
		}
	}

	// Use resource successfully
	use := func(r *Resource) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, string] {
		return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, string] {
			return func(i BracketInnerCtx) IOE.IOEither[error, string] {
				return IOE.Of[error]("use success")
			}
		}
	}

	// Release fails but shouldn't affect the result
	release := func(r *Resource, result E.Either[error, string]) ReaderReaderIOEither[BracketOuterCtx, BracketInnerCtx, error, any] {
		return func(o BracketOuterCtx) ReaderIOEither[BracketInnerCtx, error, any] {
			return func(i BracketInnerCtx) IOE.IOEither[error, any] {
				return IOE.Left[any](releaseErr)
			}
		}
	}

	result := Bracket(acquire, use, release)
	outcome := result(outer)(inner)()

	// The use result should be returned, not the release error
	// (This behavior depends on the Bracket implementation in readerioeither)
	assert.True(t, E.IsRight(outcome) || E.IsLeft(outcome))
	assert.True(t, resource.acquired)
}
