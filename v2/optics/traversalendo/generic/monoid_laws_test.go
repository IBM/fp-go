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
	"strconv"
	"strings"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/identity"
	"github.com/IBM/fp-go/v2/internal/common"
	M "github.com/IBM/fp-go/v2/monoid"
	TLG "github.com/IBM/fp-go/v2/optics/traversalendo/lens/generic"
	"github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// logged is a minimal writer applicative used to observe the order of effects
type logged[T any] struct {
	log []string
	val T
}

func ofLogged[T any](t T) logged[T] {
	return logged[T]{val: t}
}

func mapLogged[T, U any](f func(T) U) func(logged[T]) logged[U] {
	return func(lt logged[T]) logged[U] {
		return logged[U]{log: lt.log, val: f(lt.val)}
	}
}

// apLogged runs the effects of the function first, then the effects of the argument
func apLogged[U, T any](fa logged[T]) func(logged[func(T) U]) logged[U] {
	return func(fab logged[func(T) U]) logged[U] {
		return logged[U]{log: A.ArrayConcatAll(fab.log, fa.log), val: fab.val(fa.val)}
	}
}

// visit records the focus it is applied to and upper cases it
func visit(s string) logged[string] {
	return logged[string]{log: A.Of(s), val: strings.ToUpper(s)}
}

func memberRoleLens() Lens[member, string] {
	return common.MakeLens(
		func(m member) string {
			return m.role
		},
		func(m member, role string) member {
			m.role = role
			return m
		},
	)
}

type memberTraversal = Traversal[member, string, Endomorphism[member], string]

func identityMonoid() M.Monoid[memberTraversal] {
	return MakeMonoid[string, string, member](
		identity.Of[Endomorphism[member]],
		identity.Map[Endomorphism[member], Endomorphism[Endomorphism[member]]],
		identity.Ap[Endomorphism[member], Endomorphism[member]],
	)
}

func identityField(l Lens[member, string]) memberTraversal {
	return TLG.FromLens[member, string](identity.Map[string, Endomorphism[member]])(l)
}

// runIdentity modifies a member through a traversal endomorphism using the identity functor
func runIdentity(trv memberTraversal, f func(string) string) func(member) member {
	return func(m member) member {
		return trv(f)(m)(m)
	}
}

func TestMonoidIdentity(t *testing.T) {
	m := identityMonoid()
	name := identityField(memberNameLens())
	role := identityField(memberRoleLens())
	src := member{"alice", "dev"}

	t.Run("concat focuses on the targets of both traversals", func(t *testing.T) {
		res := runIdentity(m.Concat(name, role), strings.ToUpper)(src)
		assert.Equal(t, member{"ALICE", "DEV"}, res)
	})

	t.Run("empty focuses on nothing", func(t *testing.T) {
		calls := 0
		count := func(s string) string {
			calls++
			return s
		}
		assert.Equal(t, src, runIdentity(m.Empty(), count)(src))
		assert.Equal(t, 0, calls)
	})

	t.Run("empty is a left and right identity", func(t *testing.T) {
		exp := runIdentity(name, strings.ToUpper)(src)
		assert.Equal(t, exp, runIdentity(m.Concat(m.Empty(), name), strings.ToUpper)(src))
		assert.Equal(t, exp, runIdentity(m.Concat(name, m.Empty()), strings.ToUpper)(src))
	})

	t.Run("concat is associative", func(t *testing.T) {
		left := m.Concat(m.Concat(name, role), name)
		right := m.Concat(name, m.Concat(role, name))
		assert.Equal(t, runIdentity(left, strings.ToUpper)(src), runIdentity(right, strings.ToUpper)(src))
	})

	t.Run("fold of no traversals is the empty traversal", func(t *testing.T) {
		assert.Equal(t, src, runIdentity(M.Fold(m)(A.Empty[memberTraversal]()), strings.ToUpper)(src))
	})
}

func TestMonoidOption(t *testing.T) {
	m := MakeMonoid[string, option.Option[string], member](
		option.Of[Endomorphism[member]],
		option.Map[Endomorphism[member], Endomorphism[Endomorphism[member]]],
		option.Ap[Endomorphism[member], Endomorphism[member]],
	)
	field := TLG.FromLens[member, string](option.Map[string, Endomorphism[member]])
	both := m.Concat(field(memberNameLens()), field(memberRoleLens()))

	t.Run("succeeds when every focus succeeds", func(t *testing.T) {
		src := member{"alice", "dev"}
		assert.Equal(t, option.Some(member{"ALICE", "DEV"}), runEndo(src)(both(nonEmptyUpper)(src)))
	})

	t.Run("fails when any focus fails", func(t *testing.T) {
		src := member{"alice", ""}
		assert.Equal(t, option.None[member](), runEndo(src)(both(nonEmptyUpper)(src)))
	})

	t.Run("empty lifts the identity endomorphism", func(t *testing.T) {
		src := member{"", ""}
		assert.Equal(t, option.Some(src), runEndo(src)(m.Empty()(nonEmptyUpper)(src)))
	})
}

func TestConcatEffectOrder(t *testing.T) {
	concat := Concat[string, logged[string], member](
		mapLogged[Endomorphism[member], Endomorphism[Endomorphism[member]]],
		apLogged[Endomorphism[member], Endomorphism[member]],
	)
	field := TLG.FromLens[member, string](mapLogged[string, Endomorphism[member]])
	src := member{"alice", "dev"}

	res := concat(field(memberNameLens()), field(memberRoleLens()))(visit)(src)

	assert.Equal(t, member{"ALICE", "DEV"}, res.val(src))
	assert.Equal(t, A.From("alice", "dev"), res.log, "effects run left to right")
}

func TestConcatEffectLaws(t *testing.T) {
	m := MakeMonoid[string, logged[string], member](
		ofLogged[Endomorphism[member]],
		mapLogged[Endomorphism[member], Endomorphism[Endomorphism[member]]],
		apLogged[Endomorphism[member], Endomorphism[member]],
	)
	field := TLG.FromLens[member, string](mapLogged[string, Endomorphism[member]])
	name := field(memberNameLens())
	role := field(memberRoleLens())
	src := member{"alice", "dev"}

	run := func(trv Traversal[member, string, logged[Endomorphism[member]], logged[string]]) logged[member] {
		return mapLogged(endomorphism.Read(src))(trv(visit)(src))
	}

	t.Run("concat is associative including effects", func(t *testing.T) {
		left := run(m.Concat(m.Concat(name, role), name))
		right := run(m.Concat(name, m.Concat(role, name)))
		assert.Equal(t, left, right)
		assert.Equal(t, A.From("alice", "dev", "alice"), left.log)
	})

	t.Run("empty is a left and right identity including effects", func(t *testing.T) {
		assert.Equal(t, run(name), run(m.Concat(m.Empty(), name)))
		assert.Equal(t, run(name), run(m.Concat(name, m.Empty())))
	})

	t.Run("empty produces no effects", func(t *testing.T) {
		res := run(m.Empty())
		assert.Equal(t, src, res.val)
		assert.Empty(t, res.log)
	})
}

func TestConcatOverlappingFocus(t *testing.T) {
	name := identityField(memberNameLens())

	// numbered replaces each focus by the ordinal number of its visit
	calls := 0
	numbered := func(string) string {
		calls++
		return strconv.Itoa(calls)
	}

	// effects run left to right, but the resulting endomorphisms compose as l ∘ r,
	// so the update produced by the left traversal is applied last
	res := runIdentity(identityMonoid().Concat(name, name), numbered)(member{"alice", "dev"})

	assert.Equal(t, 2, calls)
	assert.Equal(t, member{"1", "dev"}, res, "the update of the left traversal wins")
}
