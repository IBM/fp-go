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
	"errors"
	"strings"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/identity"
	AR "github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/internal/common"
	TTG "github.com/IBM/fp-go/v2/optics/traversalendo/traversable/generic"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type basket struct {
	owner string
	items []string
}

var errEmpty = errors.New("empty item")

func itemsLens() Lens[basket, []string] {
	return common.MakeLens(
		func(b basket) []string {
			return b.items
		},
		func(b basket, items []string) basket {
			b.items = items
			return b
		},
	)
}

func fromArrayLensIdentity[S, T any]() func(Lens[S, []T]) Traversal[S, T, Endomorphism[S], T] {
	return FromArrayLens[[]T, S, T](
		identity.Of[[]T],
		identity.Map[[]T, func(T) []T],
		identity.Map[[]T, Endomorphism[S]],
		identity.Ap[[]T, T],
	)
}

func fromArrayLensOption[S, T any]() func(Lens[S, []T]) Traversal[S, T, option.Option[Endomorphism[S]], option.Option[T]] {
	return FromArrayLens[[]T, S, T](
		option.Of[[]T],
		option.Map[[]T, func(T) []T],
		option.Map[[]T, Endomorphism[S]],
		option.Ap[[]T, T],
	)
}

func nonEmptyUpper(s string) option.Option[string] {
	if s == "" {
		return option.None[string]()
	}
	return option.Some(strings.ToUpper(s))
}

func TestFromArrayLensIdentity(t *testing.T) {
	trv := fromArrayLensIdentity[basket, string]()(itemsLens())

	modify := func(f func(string) string) func(basket) basket {
		return func(b basket) basket {
			return trv(f)(b)(b)
		}
	}

	t.Run("traverses all elements", func(t *testing.T) {
		res := modify(strings.ToUpper)(basket{owner: "alice", items: A.From("a", "b", "c")})
		assert.Equal(t, basket{owner: "alice", items: A.From("A", "B", "C")}, res)
	})

	t.Run("handles an empty array", func(t *testing.T) {
		res := modify(strings.ToUpper)(basket{owner: "alice", items: A.Empty[string]()})
		assert.Empty(t, res.items)
		assert.Equal(t, "alice", res.owner)
	})

	t.Run("handles a nil array", func(t *testing.T) {
		res := modify(strings.ToUpper)(basket{owner: "alice"})
		assert.Empty(t, res.items)
		assert.Equal(t, "alice", res.owner)
	})

	t.Run("invokes the function once per element", func(t *testing.T) {
		calls := 0
		count := func(s string) string {
			calls++
			return s
		}
		modify(count)(basket{owner: "alice", items: A.From("a", "b", "c")})
		assert.Equal(t, 3, calls)
	})

	t.Run("visits elements from left to right", func(t *testing.T) {
		var visits []string
		record := func(s string) string {
			visits = append(visits, s)
			return s
		}
		modify(record)(basket{owner: "alice", items: A.From("a", "b", "c")})
		assert.Equal(t, A.From("a", "b", "c"), visits)
	})

	t.Run("traversing with identity is the identity", func(t *testing.T) {
		src := basket{owner: "alice", items: A.From("a", "b")}
		assert.Equal(t, src, modify(identity.Of[string])(src))
	})

	t.Run("the endomorphism sets the elements computed from the source", func(t *testing.T) {
		endo := trv(strings.ToUpper)(basket{owner: "alice", items: A.From("a", "b")})
		assert.Equal(t, basket{owner: "bob", items: A.From("A", "B")}, endo(basket{owner: "bob", items: A.From("x")}))
	})

	t.Run("does not mutate the source array", func(t *testing.T) {
		src := basket{owner: "alice", items: A.From("a", "b")}
		modify(strings.ToUpper)(src)
		assert.Equal(t, A.From("a", "b"), src.items)
	})
}

func TestFromArrayLensOption(t *testing.T) {
	trv := fromArrayLensOption[basket, string]()(itemsLens())

	run := func(src basket) option.Option[basket] {
		return option.Map(endomorphism.Read(src))(trv(nonEmptyUpper)(src))
	}

	t.Run("propagates success when every element succeeds", func(t *testing.T) {
		res := run(basket{owner: "alice", items: A.From("a", "b")})
		assert.Equal(t, option.Some(basket{owner: "alice", items: A.From("A", "B")}), res)
	})

	t.Run("propagates failure when any element fails", func(t *testing.T) {
		res := run(basket{owner: "alice", items: A.From("a", "", "c")})
		assert.Equal(t, option.None[basket](), res)
	})

	t.Run("succeeds for an empty array", func(t *testing.T) {
		res := run(basket{owner: "alice", items: A.Empty[string]()})
		assert.True(t, option.IsSome(res))
	})
}

func TestFromArrayLensMatchesTraversableLens(t *testing.T) {
	arrayTrv := fromArrayLensIdentity[basket, string]()(itemsLens())

	traversableTrv := TTG.FromTraversableLens[string, string, basket](
		identity.Map[[]string, Endomorphism[basket]],
	)(AR.MakeTraversable[[]string](
		identity.Of[[]string],
		identity.Map[[]string, func(string) []string],
		identity.Ap[[]string, string],
	))(itemsLens())

	for _, src := range []basket{
		{owner: "alice", items: A.From("a", "b")},
		{owner: "bob", items: A.Empty[string]()},
	} {
		assert.Equal(t, traversableTrv(strings.ToUpper)(src)(src), arrayTrv(strings.ToUpper)(src)(src))
	}
}

func TestFromArrayLensThunk(t *testing.T) {
	trv := FromArrayLens[[]string, basket, string](
		thunk.Of[[]string],
		thunk.Map[[]string, func(string) []string],
		thunk.Map[[]string, Endomorphism[basket]],
		thunk.Ap[[]string, string],
	)(itemsLens())

	upper := func(s string) thunk.ReaderIOResult[string] {
		if s == "" {
			return thunk.Left[string](errEmpty)
		}
		return thunk.Of(strings.ToUpper(s))
	}

	run := func(src basket) result.Result[basket] {
		return result.Map(endomorphism.Read(src))(trv(upper)(src)(t.Context())())
	}

	t.Run("propagates success when every element succeeds", func(t *testing.T) {
		res := run(basket{owner: "alice", items: A.From("a", "b")})
		assert.Equal(t, result.Of(basket{owner: "alice", items: A.From("A", "B")}), res)
	})

	t.Run("propagates the error when any element fails", func(t *testing.T) {
		res := run(basket{owner: "alice", items: A.From("a", "")})
		assert.Equal(t, result.Left[basket](errEmpty), res)
	})
}
