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

// nameOptional focuses on the name only when it is not empty
func nameOptional() Optional[person, string] {
	return common.MakeOptional(
		func(p person) option.Option[string] {
			if p.name != "" {
				return option.Some(p.name)
			}
			return option.None[string]()
		},
		func(p person, name string) person {
			p.name = name
			return p
		},
	)
}

func fromOptionalIdentity[S, A any]() func(Optional[S, A]) Traversal[S, A, Endomorphism[S], A] {
	return FromOptional[S, A](
		identity.Of[Endomorphism[S]],
		identity.Map[A, Endomorphism[S]],
	)
}

func fromOptionalOption[S, A any]() func(Optional[S, A]) Traversal[S, A, option.Option[Endomorphism[S]], option.Option[A]] {
	return FromOptional[S, A](
		option.Of[Endomorphism[S]],
		option.Map[A, Endomorphism[S]],
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

func alwaysFail(string) option.Option[string] {
	return option.None[string]()
}

func TestFromOptionalIdentity(t *testing.T) {
	trv := fromOptionalIdentity[person, string]()(nameOptional())

	modify := func(f func(string) string) func(person) person {
		return func(p person) person {
			return trv(f)(p)(p)
		}
	}

	t.Run("modifies the focus when the optional matches", func(t *testing.T) {
		res := modify(strings.ToUpper)(person{name: "alice", age: 30})
		assert.Equal(t, person{name: "ALICE", age: 30}, res)
	})

	t.Run("is a no-op when the optional does not match", func(t *testing.T) {
		src := person{name: "", age: 30}
		assert.Equal(t, src, modify(strings.ToUpper)(src))
	})

	t.Run("does not invoke the function when the optional does not match", func(t *testing.T) {
		calls := 0
		count := func(s string) string {
			calls++
			return s
		}
		modify(count)(person{name: "", age: 30})
		assert.Equal(t, 0, calls)

		modify(count)(person{name: "bob", age: 30})
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

	t.Run("the endomorphism is the identity when the source does not match", func(t *testing.T) {
		endo := trv(strings.ToUpper)(person{name: "", age: 30})
		assert.Equal(t, person{name: "bob", age: 99}, endo(person{name: "bob", age: 99}))
	})
}

func TestFromOptionalOption(t *testing.T) {
	trv := fromOptionalOption[person, string]()(nameOptional())

	t.Run("propagates a successful effect when the optional matches", func(t *testing.T) {
		src := person{name: "alice", age: 30}
		res := runEndo(src)(trv(nonEmptyUpper)(src))
		assert.Equal(t, option.Some(person{name: "ALICE", age: 30}), res)
	})

	t.Run("propagates a failing effect when the optional matches", func(t *testing.T) {
		src := person{name: "alice", age: 30}
		res := runEndo(src)(trv(alwaysFail)(src))
		assert.Equal(t, option.None[person](), res)
	})

	t.Run("lifts the identity endomorphism when the optional does not match", func(t *testing.T) {
		src := person{name: "", age: 30}
		res := runEndo(src)(trv(alwaysFail)(src))
		assert.Equal(t, option.Some(src), res)
	})

	t.Run("traversing with the pure effect is the pure source", func(t *testing.T) {
		for _, src := range []person{{name: "alice", age: 30}, {name: "", age: 40}} {
			assert.Equal(t, option.Some(src), runEndo(src)(trv(option.Of[string])(src)))
		}
	})
}
