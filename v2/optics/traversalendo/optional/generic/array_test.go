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

package generic

import (
	"strings"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/identity"
	"github.com/IBM/fp-go/v2/internal/common"
	"github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type tagged struct {
	enabled bool
	tags    []string
}

// tagsOptional focuses on the tags only when the structure is enabled
func tagsOptional() Optional[tagged, []string] {
	return common.MakeOptional(
		func(t tagged) option.Option[[]string] {
			if t.enabled {
				return option.Some(t.tags)
			}
			return option.None[[]string]()
		},
		func(t tagged, tags []string) tagged {
			t.tags = tags
			return t
		},
	)
}

func fromArrayOptional[S, T any]() func(Optional[S, []T]) Traversal[S, T, Endomorphism[S], T] {
	return FromArrayOptional[[]T, S, T](
		identity.Of[[]T],
		identity.Map[[]T, func(T) []T],
		identity.Of[Endomorphism[S]],
		identity.Map[[]T, Endomorphism[S]],
		identity.Ap[[]T, T],
	)
}

func fromArrayOptionalOption[S, T any]() func(Optional[S, []T]) Traversal[S, T, option.Option[Endomorphism[S]], option.Option[T]] {
	return FromArrayOptional[[]T, S, T](
		option.Of[[]T],
		option.Map[[]T, func(T) []T],
		option.Of[Endomorphism[S]],
		option.Map[[]T, Endomorphism[S]],
		option.Ap[[]T, T],
	)
}

func TestFromArrayOptional(t *testing.T) {
	trv := fromArrayOptional[tagged, string]()(tagsOptional())

	modify := func(f func(string) string) func(tagged) tagged {
		return func(s tagged) tagged {
			return trv(f)(s)(s)
		}
	}

	t.Run("traverses all elements when the optional matches", func(t *testing.T) {
		res := modify(strings.ToUpper)(tagged{enabled: true, tags: A.From("a", "b", "c")})
		assert.Equal(t, tagged{enabled: true, tags: A.From("A", "B", "C")}, res)
	})

	t.Run("is a no-op when the optional does not match", func(t *testing.T) {
		src := tagged{enabled: false, tags: A.From("a", "b")}
		assert.Equal(t, src, modify(strings.ToUpper)(src))
	})

	t.Run("handles an empty array when the optional matches", func(t *testing.T) {
		res := modify(strings.ToUpper)(tagged{enabled: true, tags: A.Empty[string]()})
		assert.Empty(t, res.tags)
		assert.True(t, res.enabled)
	})

	t.Run("invokes the function once per element only when the optional matches", func(t *testing.T) {
		calls := 0
		count := func(s string) string {
			calls++
			return s
		}
		modify(count)(tagged{enabled: false, tags: A.From("a", "b", "c")})
		assert.Equal(t, 0, calls)

		modify(count)(tagged{enabled: true, tags: A.From("a", "b", "c")})
		assert.Equal(t, 3, calls)
	})

	t.Run("visits elements from left to right", func(t *testing.T) {
		var visits []string
		record := func(s string) string {
			visits = append(visits, s)
			return s
		}
		modify(record)(tagged{enabled: true, tags: A.From("a", "b", "c")})
		assert.Equal(t, A.From("a", "b", "c"), visits)
	})

	t.Run("traversing with identity is the identity", func(t *testing.T) {
		for _, src := range []tagged{{enabled: true, tags: A.From("a", "b")}, {enabled: false, tags: A.From("c")}} {
			assert.Equal(t, src, modify(identity.Of[string])(src))
		}
	})

	t.Run("does not mutate the source array", func(t *testing.T) {
		src := tagged{enabled: true, tags: A.From("a", "b")}
		modify(strings.ToUpper)(src)
		assert.Equal(t, A.From("a", "b"), src.tags)
	})
}

func TestFromArrayOptionalOption(t *testing.T) {
	trv := fromArrayOptionalOption[tagged, string]()(tagsOptional())

	t.Run("propagates success when the optional matches and every element succeeds", func(t *testing.T) {
		src := tagged{enabled: true, tags: A.From("a", "b")}
		res := runEndo(src)(trv(nonEmptyUpper)(src))
		assert.Equal(t, option.Some(tagged{enabled: true, tags: A.From("A", "B")}), res)
	})

	t.Run("propagates failure when the optional matches and any element fails", func(t *testing.T) {
		src := tagged{enabled: true, tags: A.From("a", "", "c")}
		res := runEndo(src)(trv(nonEmptyUpper)(src))
		assert.Equal(t, option.None[tagged](), res)
	})

	t.Run("lifts the identity endomorphism when the optional does not match", func(t *testing.T) {
		src := tagged{enabled: false, tags: A.From("a", "b")}
		res := runEndo(src)(trv(alwaysFail)(src))
		assert.Equal(t, option.Some(src), res)
	})

	t.Run("succeeds for an empty array when the optional matches", func(t *testing.T) {
		src := tagged{enabled: true, tags: A.Empty[string]()}
		res := runEndo(src)(trv(alwaysFail)(src))
		assert.True(t, option.IsSome(res))
	})
}
