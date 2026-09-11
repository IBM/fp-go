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
	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/identity"
	AR "github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/internal/common"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type basket struct {
	owner string
	items []string
}

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

func identityArray[T any]() Traversable[T, T, []T, []T] {
	return AR.MakeTraversable[[]T](
		identity.Of[[]T],
		identity.Map[[]T, func(T) []T],
		identity.Ap[[]T, T],
	)
}

func optionArray[T any]() Traversable[T, option.Option[T], []T, option.Option[[]T]] {
	return AR.MakeTraversable[[]T](
		option.Of[[]T],
		option.Map[[]T, func(T) []T],
		option.Ap[[]T, T],
	)
}

// runEndo applies an optional endomorphism to its source value
func runEndo[S any](s S) func(option.Option[Endomorphism[S]]) option.Option[S] {
	return option.Map(func(e Endomorphism[S]) S {
		return e(s)
	})
}

func nonEmptyUpper(s string) option.Option[string] {
	if s == "" {
		return option.None[string]()
	}
	return option.Some(strings.ToUpper(s))
}

func TestFromTraversableLensIdentity(t *testing.T) {
	trv := FromTraversableLens[string, string, basket](
		identity.Map[[]string, Endomorphism[basket]],
	)(identityArray[string]())(itemsLens())

	modify := func(f func(string) string) func(basket) basket {
		return func(b basket) basket {
			return trv(f)(b)(b)
		}
	}

	t.Run("traverses all elements", func(t *testing.T) {
		res := modify(strings.ToUpper)(basket{owner: "alice", items: A.From("a", "b", "c")})
		assert.Equal(t, basket{owner: "alice", items: A.From("A", "B", "C")}, res)
	})

	t.Run("handles an empty traversable", func(t *testing.T) {
		res := modify(strings.ToUpper)(basket{owner: "alice", items: A.Empty[string]()})
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

	t.Run("traversing with identity is the identity", func(t *testing.T) {
		src := basket{owner: "alice", items: A.From("a", "b")}
		assert.Equal(t, src, modify(identity.Of[string])(src))
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

	t.Run("the endomorphism sets the elements computed from the source", func(t *testing.T) {
		endo := trv(strings.ToUpper)(basket{owner: "alice", items: A.From("a", "b")})
		assert.Equal(t, basket{owner: "bob", items: A.From("A", "B")}, endo(basket{owner: "bob", items: A.From("x")}))
	})
}

func TestFromTraversableLensOption(t *testing.T) {
	trv := FromTraversableLens[string, option.Option[string], basket](
		option.Map[[]string, Endomorphism[basket]],
	)(optionArray[string]())(itemsLens())

	t.Run("propagates success when every element succeeds", func(t *testing.T) {
		src := basket{owner: "alice", items: A.From("a", "b")}
		res := runEndo(src)(trv(nonEmptyUpper)(src))
		assert.Equal(t, option.Some(basket{owner: "alice", items: A.From("A", "B")}), res)
	})

	t.Run("propagates failure when any element fails", func(t *testing.T) {
		src := basket{owner: "alice", items: A.From("a", "", "c")}
		res := runEndo(src)(trv(nonEmptyUpper)(src))
		assert.Equal(t, option.None[basket](), res)
	})
}

type config struct {
	name     string
	database option.Option[string]
}

func databaseLens() Lens[config, option.Option[string]] {
	return common.MakeLens(
		func(c config) option.Option[string] {
			return c.database
		},
		func(c config, database option.Option[string]) config {
			c.database = database
			return c
		},
	)
}

func identityOption[T any]() Traversable[T, T, option.Option[T], option.Option[T]] {
	return option.MakeTraversable[T, T](
		identity.Of[option.Option[T]],
		identity.Map[T, option.Option[T]],
	)
}

func optionOption[T any]() Traversable[T, option.Option[T], option.Option[T], option.Option[option.Option[T]]] {
	return option.MakeTraversable[T, T](
		option.Of[option.Option[T]],
		option.Map[T, option.Option[T]],
	)
}

func toURL(db string) string {
	return "postgresql://" + db
}

func TestFromTraversableLensOptionFieldIdentity(t *testing.T) {
	trv := FromTraversableLens[string, string, config](
		identity.Map[option.Option[string], Endomorphism[config]],
	)(identityOption[string]())(databaseLens())

	modify := func(f func(string) string) func(config) config {
		return func(c config) config {
			return trv(f)(c)(c)
		}
	}

	t.Run("modifies a present value", func(t *testing.T) {
		res := modify(toURL)(config{name: "app", database: option.Some("localhost")})
		assert.Equal(t, config{name: "app", database: option.Some("postgresql://localhost")}, res)
	})

	t.Run("is a no-op for an absent value", func(t *testing.T) {
		calls := 0
		count := func(s string) string {
			calls++
			return s
		}
		src := config{name: "app", database: option.None[string]()}
		assert.Equal(t, src, modify(count)(src))
		assert.Equal(t, 0, calls)
	})
}

func TestFromTraversableLensOptionFieldOption(t *testing.T) {
	trv := FromTraversableLens[string, option.Option[string], config](
		option.Map[option.Option[string], Endomorphism[config]],
	)(optionOption[string]())(databaseLens())

	t.Run("propagates success for a present value", func(t *testing.T) {
		src := config{name: "app", database: option.Some("localhost")}
		res := runEndo(src)(trv(nonEmptyUpper)(src))
		assert.Equal(t, option.Some(config{name: "app", database: option.Some("LOCALHOST")}), res)
	})

	t.Run("propagates failure for a present value", func(t *testing.T) {
		src := config{name: "app", database: option.Some("")}
		res := runEndo(src)(trv(nonEmptyUpper)(src))
		assert.Equal(t, option.None[config](), res)
	})

	t.Run("succeeds unchanged for an absent value", func(t *testing.T) {
		src := config{name: "app", database: option.None[string]()}
		res := runEndo(src)(trv(nonEmptyUpper)(src))
		assert.Equal(t, option.Some(src), res)
	})
}

func TestFromTraversableLensThunk(t *testing.T) {
	trv := FromTraversableLens[string, thunk.ReaderIOResult[string], config](
		thunk.Map[option.Option[string], Endomorphism[config]],
	)(option.MakeTraversable[string, string](
		thunk.Of[option.Option[string]],
		thunk.Map[string, option.Option[string]],
	))(databaseLens())

	run := func(src config) result.Result[config] {
		return result.Map(endomorphism.Read(src))(trv(F.Flow2(toURL, thunk.Of[string]))(src)(t.Context())())
	}

	t.Run("modifies a present value", func(t *testing.T) {
		res := run(config{name: "app", database: option.Some("localhost")})
		assert.Equal(t, result.Of(config{name: "app", database: option.Some("postgresql://localhost")}), res)
	})

	t.Run("is a no-op for an absent value", func(t *testing.T) {
		src := config{name: "app", database: option.None[string]()}
		assert.Equal(t, result.Of(src), run(src))
	})
}
