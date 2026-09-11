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

	"github.com/IBM/fp-go/v2/identity"
	"github.com/IBM/fp-go/v2/internal/common"
	"github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type person struct {
	name string
	age  int
}

func nameLens() Lens[person, string] {
	return common.MakeLens(
		func(p person) string {
			return p.name
		},
		func(p person, name string) person {
			p.name = name
			return p
		},
	)
}

func fromLensIdentity[S, A any]() func(Lens[S, A]) Traversal[S, A, Endomorphism[S], A] {
	return FromLens[S, A](identity.Map[A, Endomorphism[S]])
}

func fromLensOption[S, A any]() func(Lens[S, A]) Traversal[S, A, option.Option[Endomorphism[S]], option.Option[A]] {
	return FromLens[S, A](option.Map[A, Endomorphism[S]])
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

func TestFromLensIdentity(t *testing.T) {
	trv := fromLensIdentity[person, string]()(nameLens())

	modify := func(f func(string) string) func(person) person {
		return func(p person) person {
			return trv(f)(p)(p)
		}
	}

	t.Run("modifies the focus", func(t *testing.T) {
		res := modify(strings.ToUpper)(person{name: "alice", age: 30})
		assert.Equal(t, person{name: "ALICE", age: 30}, res)
	})

	t.Run("invokes the function exactly once", func(t *testing.T) {
		calls := 0
		count := func(s string) string {
			calls++
			return s
		}
		modify(count)(person{name: "", age: 30})
		assert.Equal(t, 1, calls)
	})

	t.Run("traversing with identity is the identity", func(t *testing.T) {
		for _, src := range []person{{name: "alice", age: 30}, {name: "", age: 40}} {
			assert.Equal(t, src, modify(identity.Of[string])(src))
		}
	})

	t.Run("the endomorphism sets the focus computed from the source", func(t *testing.T) {
		endo := trv(strings.ToUpper)(person{name: "alice", age: 30})
		assert.Equal(t, person{name: "ALICE", age: 99}, endo(person{name: "bob", age: 99}))
	})
}

func TestFromLensOption(t *testing.T) {
	trv := fromLensOption[person, string]()(nameLens())

	t.Run("propagates a successful effect", func(t *testing.T) {
		src := person{name: "alice", age: 30}
		res := runEndo(src)(trv(nonEmptyUpper)(src))
		assert.Equal(t, option.Some(person{name: "ALICE", age: 30}), res)
	})

	t.Run("propagates a failing effect", func(t *testing.T) {
		src := person{name: "", age: 30}
		res := runEndo(src)(trv(nonEmptyUpper)(src))
		assert.Equal(t, option.None[person](), res)
	})

	t.Run("traversing with the pure effect is the pure source", func(t *testing.T) {
		for _, src := range []person{{name: "alice", age: 30}, {name: "", age: 40}} {
			assert.Equal(t, option.Some(src), runEndo(src)(trv(option.Of[string])(src)))
		}
	})
}
