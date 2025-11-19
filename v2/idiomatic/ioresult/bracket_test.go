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

	"github.com/stretchr/testify/assert"
)

func TestBracket(t *testing.T) {
	t.Run("successful acquire, use, and release", func(t *testing.T) {
		acquired := false
		used := false
		released := false

		acquire := func() (string, error) {
			acquired = true
			return "resource", nil
		}

		use := func(r string) IOResult[int] {
			return func() (int, error) {
				used = true
				assert.Equal(t, "resource", r)
				return 42, nil
			}
		}

		release := func(b int, err error) func(string) IOResult[any] {
			return func(r string) IOResult[any] {
				return func() (any, error) {
					released = true
					assert.Equal(t, "resource", r)
					assert.Equal(t, 42, b)
					assert.NoError(t, err)
					return nil, nil
				}
			}
		}

		result := Bracket(acquire, use, release)
		val, err := result()

		assert.True(t, acquired)
		assert.True(t, used)
		assert.True(t, released)
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	})

	t.Run("acquire fails - use and release not called", func(t *testing.T) {
		used := false
		released := false

		acquire := func() (string, error) {
			return "", errors.New("acquire failed")
		}

		use := func(r string) IOResult[int] {
			return func() (int, error) {
				used = true
				return 42, nil
			}
		}

		release := func(b int, err error) func(string) IOResult[any] {
			return func(r string) IOResult[any] {
				return func() (any, error) {
					released = true
					return nil, nil
				}
			}
		}

		result := Bracket(acquire, use, release)
		_, err := result()

		assert.False(t, used)
		assert.False(t, released)
		assert.Error(t, err)
		assert.Equal(t, "acquire failed", err.Error())
	})

	t.Run("use fails - release is still called", func(t *testing.T) {
		acquired := false
		released := false

		acquire := func() (string, error) {
			acquired = true
			return "resource", nil
		}

		use := func(r string) IOResult[int] {
			return func() (int, error) {
				return 0, errors.New("use failed")
			}
		}

		release := func(b int, err error) func(string) IOResult[any] {
			return func(r string) IOResult[any] {
				return func() (any, error) {
					released = true
					assert.Equal(t, "resource", r)
					assert.Equal(t, 0, b)
					assert.Error(t, err)
					assert.Equal(t, "use failed", err.Error())
					return nil, nil
				}
			}
		}

		result := Bracket(acquire, use, release)
		_, err := result()

		assert.True(t, acquired)
		assert.True(t, released)
		assert.Error(t, err)
		assert.Equal(t, "use failed", err.Error())
	})

	t.Run("use succeeds but release fails", func(t *testing.T) {
		acquired := false
		used := false
		released := false

		acquire := func() (string, error) {
			acquired = true
			return "resource", nil
		}

		use := func(r string) IOResult[int] {
			return func() (int, error) {
				used = true
				return 42, nil
			}
		}

		release := func(b int, err error) func(string) IOResult[any] {
			return func(r string) IOResult[any] {
				return func() (any, error) {
					released = true
					return nil, errors.New("release failed")
				}
			}
		}

		result := Bracket(acquire, use, release)
		_, err := result()

		assert.True(t, acquired)
		assert.True(t, used)
		assert.True(t, released)
		assert.Error(t, err)
		assert.Equal(t, "release failed", err.Error())
	})

	t.Run("both use and release fail - use error is returned", func(t *testing.T) {
		acquire := func() (string, error) {
			return "resource", nil
		}

		use := func(r string) IOResult[int] {
			return func() (int, error) {
				return 0, errors.New("use failed")
			}
		}

		release := func(b int, err error) func(string) IOResult[any] {
			return func(r string) IOResult[any] {
				return func() (any, error) {
					assert.Error(t, err)
					assert.Equal(t, "use failed", err.Error())
					return nil, errors.New("release failed")
				}
			}
		}

		result := Bracket(acquire, use, release)
		_, err := result()

		assert.Error(t, err)
		// use error takes precedence
		assert.Equal(t, "use failed", err.Error())
	})

	t.Run("resource cleanup with file-like resource", func(t *testing.T) {
		type File struct {
			name   string
			closed bool
		}

		var file *File

		acquire := func() (*File, error) {
			file = &File{name: "test.txt", closed: false}
			return file, nil
		}

		use := func(f *File) IOResult[string] {
			return func() (string, error) {
				assert.False(t, f.closed)
				return "file content", nil
			}
		}

		release := func(content string, err error) func(*File) IOResult[any] {
			return func(f *File) IOResult[any] {
				return func() (any, error) {
					f.closed = true
					return nil, nil
				}
			}
		}

		result := Bracket(acquire, use, release)
		content, err := result()

		assert.NoError(t, err)
		assert.Equal(t, "file content", content)
		assert.True(t, file.closed)
	})

	t.Run("release receives both value and error from use", func(t *testing.T) {
		var receivedValue int
		var receivedError error

		acquire := func() (string, error) {
			return "resource", nil
		}

		use := func(r string) IOResult[int] {
			return func() (int, error) {
				return 100, errors.New("use error")
			}
		}

		release := func(b int, err error) func(string) IOResult[any] {
			return func(r string) IOResult[any] {
				return func() (any, error) {
					receivedValue = b
					receivedError = err
					return nil, nil
				}
			}
		}

		result := Bracket(acquire, use, release)
		_, _ = result()

		assert.Equal(t, 100, receivedValue)
		assert.Error(t, receivedError)
		assert.Equal(t, "use error", receivedError.Error())
	})

	t.Run("release receives zero value and nil when use succeeds", func(t *testing.T) {
		var receivedValue int
		var receivedError error

		acquire := func() (string, error) {
			return "resource", nil
		}

		use := func(r string) IOResult[int] {
			return func() (int, error) {
				return 42, nil
			}
		}

		release := func(b int, err error) func(string) IOResult[any] {
			return func(r string) IOResult[any] {
				return func() (any, error) {
					receivedValue = b
					receivedError = err
					return nil, nil
				}
			}
		}

		result := Bracket(acquire, use, release)
		val, err := result()

		assert.NoError(t, err)
		assert.Equal(t, 42, val)
		assert.Equal(t, 42, receivedValue)
		assert.NoError(t, receivedError)
	})
}
