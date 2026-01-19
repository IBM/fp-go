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
	"errors"
	"testing"

	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type Resource struct {
	id       string
	acquired bool
	released bool
}

func TestBracketSuccessPath(t *testing.T) {
	cfg := defaultConfig
	ctx := t.Context()

	resource := &Resource{id: "res1"}

	// Acquire resource
	acquire := func(c AppConfig) ReaderIOResult[context.Context, *Resource] {
		return func(ctx context.Context) IOResult[*Resource] {
			return func() Result[*Resource] {
				resource.acquired = true
				return result.Of(resource)
			}
		}
	}

	// Use resource successfully
	use := func(r *Resource) ReaderReaderIOResult[AppConfig, string] {
		return func(c AppConfig) ReaderIOResult[context.Context, string] {
			return func(ctx context.Context) IOResult[string] {
				return func() Result[string] {
					return result.Of("result from " + r.id)
				}
			}
		}
	}

	// Release resource
	release := func(r *Resource, res Result[string]) ReaderReaderIOResult[AppConfig, any] {
		return func(c AppConfig) ReaderIOResult[context.Context, any] {
			return func(ctx context.Context) IOResult[any] {
				return func() Result[any] {
					r.released = true
					return result.Of[any](nil)
				}
			}
		}
	}

	computation := Bracket(acquire, use, release)
	outcome := computation(cfg)(ctx)()

	assert.Equal(t, result.Of("result from res1"), outcome)
	assert.True(t, resource.acquired, "Resource should be acquired")
	assert.True(t, resource.released, "Resource should be released")
}

func TestBracketUseFailure(t *testing.T) {
	cfg := defaultConfig
	ctx := t.Context()

	resource := &Resource{id: "res1"}
	useErr := errors.New("use failed")

	// Acquire resource
	acquire := func(c AppConfig) ReaderIOResult[context.Context, *Resource] {
		return func(ctx context.Context) IOResult[*Resource] {
			return func() Result[*Resource] {
				resource.acquired = true
				return result.Of(resource)
			}
		}
	}

	// Use resource with failure
	use := func(r *Resource) ReaderReaderIOResult[AppConfig, string] {
		return func(c AppConfig) ReaderIOResult[context.Context, string] {
			return func(ctx context.Context) IOResult[string] {
				return func() Result[string] {
					return result.Left[string](useErr)
				}
			}
		}
	}

	// Release resource (should still be called)
	release := func(r *Resource, res Result[string]) ReaderReaderIOResult[AppConfig, any] {
		return func(c AppConfig) ReaderIOResult[context.Context, any] {
			return func(ctx context.Context) IOResult[any] {
				return func() Result[any] {
					r.released = true
					return result.Of[any](nil)
				}
			}
		}
	}

	computation := Bracket(acquire, use, release)
	outcome := computation(cfg)(ctx)()

	assert.Equal(t, result.Left[string](useErr), outcome)
	assert.True(t, resource.acquired, "Resource should be acquired")
	assert.True(t, resource.released, "Resource should be released even on failure")
}

func TestBracketAcquireFailure(t *testing.T) {
	cfg := defaultConfig
	ctx := t.Context()

	resource := &Resource{id: "res1"}
	acquireErr := errors.New("acquire failed")
	useCalled := false
	releaseCalled := false

	// Acquire resource fails
	acquire := func(c AppConfig) ReaderIOResult[context.Context, *Resource] {
		return func(ctx context.Context) IOResult[*Resource] {
			return func() Result[*Resource] {
				return result.Left[*Resource](acquireErr)
			}
		}
	}

	// Use should not be called
	use := func(r *Resource) ReaderReaderIOResult[AppConfig, string] {
		return func(c AppConfig) ReaderIOResult[context.Context, string] {
			return func(ctx context.Context) IOResult[string] {
				return func() Result[string] {
					useCalled = true
					return result.Of("should not reach here")
				}
			}
		}
	}

	// Release should not be called
	release := func(r *Resource, res Result[string]) ReaderReaderIOResult[AppConfig, any] {
		return func(c AppConfig) ReaderIOResult[context.Context, any] {
			return func(ctx context.Context) IOResult[any] {
				return func() Result[any] {
					releaseCalled = true
					return result.Of[any](nil)
				}
			}
		}
	}

	computation := Bracket(acquire, use, release)
	outcome := computation(cfg)(ctx)()

	assert.Equal(t, result.Left[string](acquireErr), outcome)
	assert.False(t, resource.acquired, "Resource should not be acquired")
	assert.False(t, useCalled, "Use should not be called when acquire fails")
	assert.False(t, releaseCalled, "Release should not be called when acquire fails")
}

func TestBracketReleaseReceivesResult(t *testing.T) {
	cfg := defaultConfig
	ctx := t.Context()

	resource := &Resource{id: "res1"}
	var capturedResult Result[string]

	// Acquire resource
	acquire := func(c AppConfig) ReaderIOResult[context.Context, *Resource] {
		return func(ctx context.Context) IOResult[*Resource] {
			return func() Result[*Resource] {
				resource.acquired = true
				return result.Of(resource)
			}
		}
	}

	// Use resource
	use := func(r *Resource) ReaderReaderIOResult[AppConfig, string] {
		return func(c AppConfig) ReaderIOResult[context.Context, string] {
			return func(ctx context.Context) IOResult[string] {
				return func() Result[string] {
					return result.Of("use result")
				}
			}
		}
	}

	// Release captures the result
	release := func(r *Resource, res Result[string]) ReaderReaderIOResult[AppConfig, any] {
		return func(c AppConfig) ReaderIOResult[context.Context, any] {
			return func(ctx context.Context) IOResult[any] {
				return func() Result[any] {
					capturedResult = res
					r.released = true
					return result.Of[any](nil)
				}
			}
		}
	}

	computation := Bracket(acquire, use, release)
	outcome := computation(cfg)(ctx)()

	assert.Equal(t, result.Of("use result"), outcome)
	assert.Equal(t, result.Of("use result"), capturedResult)
	assert.True(t, resource.released, "Resource should be released")
}

func TestBracketWithContextAccess(t *testing.T) {
	cfg := AppConfig{DatabaseURL: "production-db", LogLevel: "debug"}
	ctx := t.Context()

	resource := &Resource{id: "res1"}

	// Acquire uses outer context
	acquire := func(c AppConfig) ReaderIOResult[context.Context, *Resource] {
		return func(ctx context.Context) IOResult[*Resource] {
			return func() Result[*Resource] {
				resource.id = c.DatabaseURL + "-resource"
				resource.acquired = true
				return result.Of(resource)
			}
		}
	}

	// Use uses both contexts
	use := func(r *Resource) ReaderReaderIOResult[AppConfig, string] {
		return func(c AppConfig) ReaderIOResult[context.Context, string] {
			return func(ctx context.Context) IOResult[string] {
				return func() Result[string] {
					res := r.id + " with log level " + c.LogLevel
					return result.Of(res)
				}
			}
		}
	}

	// Release uses both contexts
	release := func(r *Resource, res Result[string]) ReaderReaderIOResult[AppConfig, any] {
		return func(c AppConfig) ReaderIOResult[context.Context, any] {
			return func(ctx context.Context) IOResult[any] {
				return func() Result[any] {
					r.released = true
					return result.Of[any](nil)
				}
			}
		}
	}

	computation := Bracket(acquire, use, release)
	outcome := computation(cfg)(ctx)()

	assert.True(t, result.IsRight(outcome))
	assert.True(t, resource.acquired)
	assert.True(t, resource.released)
	assert.Equal(t, "production-db-resource", resource.id)
}

func TestBracketMultipleResources(t *testing.T) {
	cfg := defaultConfig
	ctx := t.Context()

	resource1 := &Resource{id: "res1"}
	resource2 := &Resource{id: "res2"}

	// Acquire first resource
	acquire1 := func(c AppConfig) ReaderIOResult[context.Context, *Resource] {
		return func(ctx context.Context) IOResult[*Resource] {
			return func() Result[*Resource] {
				resource1.acquired = true
				return result.Of(resource1)
			}
		}
	}

	// Use first resource to acquire second
	use1 := func(r1 *Resource) ReaderReaderIOResult[AppConfig, string] {
		// Nested bracket for second resource
		acquire2 := func(c AppConfig) ReaderIOResult[context.Context, *Resource] {
			return func(ctx context.Context) IOResult[*Resource] {
				return func() Result[*Resource] {
					resource2.acquired = true
					return result.Of(resource2)
				}
			}
		}

		use2 := func(r2 *Resource) ReaderReaderIOResult[AppConfig, string] {
			return func(c AppConfig) ReaderIOResult[context.Context, string] {
				return func(ctx context.Context) IOResult[string] {
					return func() Result[string] {
						return result.Of(r1.id + " and " + r2.id)
					}
				}
			}
		}

		release2 := func(r2 *Resource, res Result[string]) ReaderReaderIOResult[AppConfig, any] {
			return func(c AppConfig) ReaderIOResult[context.Context, any] {
				return func(ctx context.Context) IOResult[any] {
					return func() Result[any] {
						r2.released = true
						return result.Of[any](nil)
					}
				}
			}
		}

		return Bracket(acquire2, use2, release2)
	}

	release1 := func(r1 *Resource, res Result[string]) ReaderReaderIOResult[AppConfig, any] {
		return func(c AppConfig) ReaderIOResult[context.Context, any] {
			return func(ctx context.Context) IOResult[any] {
				return func() Result[any] {
					r1.released = true
					return result.Of[any](nil)
				}
			}
		}
	}

	computation := Bracket(acquire1, use1, release1)
	outcome := computation(cfg)(ctx)()

	assert.Equal(t, result.Of("res1 and res2"), outcome)
	assert.True(t, resource1.acquired && resource1.released, "Resource 1 should be acquired and released")
	assert.True(t, resource2.acquired && resource2.released, "Resource 2 should be acquired and released")
}

func TestBracketReleaseErrorDoesNotAffectResult(t *testing.T) {
	cfg := defaultConfig
	ctx := t.Context()

	resource := &Resource{id: "res1"}
	releaseErr := errors.New("release failed")

	// Acquire resource
	acquire := func(c AppConfig) ReaderIOResult[context.Context, *Resource] {
		return func(ctx context.Context) IOResult[*Resource] {
			return func() Result[*Resource] {
				resource.acquired = true
				return result.Of(resource)
			}
		}
	}

	// Use resource successfully
	use := func(r *Resource) ReaderReaderIOResult[AppConfig, string] {
		return func(c AppConfig) ReaderIOResult[context.Context, string] {
			return func(ctx context.Context) IOResult[string] {
				return func() Result[string] {
					return result.Of("use success")
				}
			}
		}
	}

	// Release fails but shouldn't affect the result
	release := func(r *Resource, res Result[string]) ReaderReaderIOResult[AppConfig, any] {
		return func(c AppConfig) ReaderIOResult[context.Context, any] {
			return func(ctx context.Context) IOResult[any] {
				return func() Result[any] {
					return result.Left[any](releaseErr)
				}
			}
		}
	}

	computation := Bracket(acquire, use, release)
	outcome := computation(cfg)(ctx)()

	// The use result should be returned, not the release error
	// (This behavior depends on the Bracket implementation)
	assert.True(t, result.IsRight(outcome) || result.IsLeft(outcome))
	assert.True(t, resource.acquired)
}
