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
	AR "github.com/IBM/fp-go/v2/internal/array"
	"github.com/IBM/fp-go/v2/internal/common"
	TLG "github.com/IBM/fp-go/v2/optics/traversalendo/lens/generic"
	TTG "github.com/IBM/fp-go/v2/optics/traversalendo/traversable/generic"
	"github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type (
	member struct {
		name string
		role string
	}

	team struct {
		name    string
		members []member
	}
)

func memberNameLens() Lens[member, string] {
	return common.MakeLens(
		func(m member) string {
			return m.name
		},
		func(m member, name string) member {
			m.name = name
			return m
		},
	)
}

func teamMembersLens() Lens[team, []member] {
	return common.MakeLens(
		func(t team) []member {
			return t.members
		},
		func(t team, members []member) team {
			t.members = members
			return t
		},
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

// memberNamesIdentity focuses on the names of all members of a team using the identity functor
func memberNamesIdentity() Traversal[team, string, Endomorphism[team], string] {
	nameTrav := TLG.FromLens[member, string](identity.Map[string, Endomorphism[member]])(memberNameLens())

	membersTrav := TTG.FromTraversableLens[member, member, team](
		identity.Map[[]member, Endomorphism[team]],
	)(AR.MakeTraversable[[]member](
		identity.Of[[]member],
		identity.Map[[]member, func(member) []member],
		identity.Ap[[]member, member],
	))(teamMembersLens())

	return Compose[team, string, Endomorphism[team], string](
		identity.Map[Endomorphism[member], member],
	)(nameTrav)(membersTrav)
}

// memberNamesOption focuses on the names of all members of a team using the option functor
func memberNamesOption() Traversal[team, string, option.Option[Endomorphism[team]], option.Option[string]] {
	nameTrav := TLG.FromLens[member, string](option.Map[string, Endomorphism[member]])(memberNameLens())

	membersTrav := TTG.FromTraversableLens[member, option.Option[member], team](
		option.Map[[]member, Endomorphism[team]],
	)(AR.MakeTraversable[[]member](
		option.Of[[]member],
		option.Map[[]member, func(member) []member],
		option.Ap[[]member, member],
	))(teamMembersLens())

	return Compose[team, string, option.Option[Endomorphism[team]], option.Option[string]](
		option.Map[Endomorphism[member], member],
	)(nameTrav)(membersTrav)
}

func TestComposeIdentity(t *testing.T) {
	trv := memberNamesIdentity()

	modify := func(f func(string) string) func(team) team {
		return func(t team) team {
			return trv(f)(t)(t)
		}
	}

	t.Run("modifies the nested focus of every element", func(t *testing.T) {
		src := team{name: "core", members: A.From(member{"alice", "dev"}, member{"bob", "ops"})}
		exp := team{name: "core", members: A.From(member{"ALICE", "dev"}, member{"BOB", "ops"})}
		assert.Equal(t, exp, modify(strings.ToUpper)(src))
	})

	t.Run("handles an empty outer structure", func(t *testing.T) {
		res := modify(strings.ToUpper)(team{name: "core", members: A.Empty[member]()})
		assert.Empty(t, res.members)
		assert.Equal(t, "core", res.name)
	})

	t.Run("invokes the function once per nested focus", func(t *testing.T) {
		calls := 0
		count := func(s string) string {
			calls++
			return s
		}
		modify(count)(team{name: "core", members: A.From(member{"a", "x"}, member{"b", "y"}, member{"c", "z"})})
		assert.Equal(t, 3, calls)
	})

	t.Run("traversing with identity is the identity", func(t *testing.T) {
		src := team{name: "core", members: A.From(member{"alice", "dev"}, member{"bob", "ops"})}
		assert.Equal(t, src, modify(identity.Of[string])(src))
	})
}

func TestComposeOption(t *testing.T) {
	trv := memberNamesOption()

	t.Run("propagates success when every nested focus succeeds", func(t *testing.T) {
		src := team{name: "core", members: A.From(member{"alice", "dev"}, member{"bob", "ops"})}
		exp := team{name: "core", members: A.From(member{"ALICE", "dev"}, member{"BOB", "ops"})}
		assert.Equal(t, option.Some(exp), runEndo(src)(trv(nonEmptyUpper)(src)))
	})

	t.Run("propagates failure when any nested focus fails", func(t *testing.T) {
		src := team{name: "core", members: A.From(member{"alice", "dev"}, member{"", "ops"})}
		assert.Equal(t, option.None[team](), runEndo(src)(trv(nonEmptyUpper)(src)))
	})
}

type division struct {
	name  string
	teams []team
}

func divisionTeamsLens() Lens[division, []team] {
	return common.MakeLens(
		func(d division) []team {
			return d.teams
		},
		func(d division, teams []team) division {
			d.teams = teams
			return d
		},
	)
}

func teamNameLens() Lens[team, string] {
	return common.MakeLens(
		func(t team) string {
			return t.name
		},
		func(t team, name string) team {
			t.name = name
			return t
		},
	)
}

func identityArray[T any]() TTG.Traversable[T, T, []T, []T] {
	return AR.MakeTraversable[[]T](
		identity.Of[[]T],
		identity.Map[[]T, func(T) []T],
		identity.Ap[[]T, T],
	)
}

// teamMembersIdentity focuses on all members of a team using the identity functor
func teamMembersIdentity() Traversal[team, member, Endomorphism[team], member] {
	return TTG.FromTraversableLens[member, member, team](
		identity.Map[[]member, Endomorphism[team]],
	)(identityArray[member]())(teamMembersLens())
}

// divisionTeamsIdentity focuses on all teams of a division using the identity functor
func divisionTeamsIdentity() Traversal[division, team, Endomorphism[division], team] {
	return TTG.FromTraversableLens[team, team, division](
		identity.Map[[]team, Endomorphism[division]],
	)(identityArray[team]())(divisionTeamsLens())
}

// modifyWith modifies a source through a traversal endomorphism using the identity functor
func modifyWith[S, A any](trv Traversal[S, A, Endomorphism[S], A], f func(A) A) func(S) S {
	return func(s S) S {
		return trv(f)(s)(s)
	}
}

// recordUpper upper cases its input and records every visited focus
func recordUpper(visits *[]string) func(string) string {
	return func(s string) string {
		*visits = append(*visits, s)
		return strings.ToUpper(s)
	}
}

func TestComposeAssociativity(t *testing.T) {
	sa := divisionTeamsIdentity()
	ab := teamMembersIdentity()
	bc := identityField(memberNameLens())

	// Compose(bc)(Compose(ab)(sa))
	left := Compose[division, string, Endomorphism[division], string](
		identity.Map[Endomorphism[member], member],
	)(bc)(Compose[division, member, Endomorphism[division], member](
		identity.Map[Endomorphism[team], team],
	)(ab)(sa))

	// Compose(Compose(bc)(ab))(sa)
	right := Compose[division, string, Endomorphism[division], string](
		identity.Map[Endomorphism[team], team],
	)(Compose[team, string, Endomorphism[team], string](
		identity.Map[Endomorphism[member], member],
	)(bc)(ab))(sa)

	src := division{name: "eng", teams: A.From(
		team{name: "core", members: A.From(member{"alice", "dev"}, member{"bob", "ops"})},
		team{name: "web", members: A.From(member{"carol", "dev"})},
	)}
	exp := division{name: "eng", teams: A.From(
		team{name: "core", members: A.From(member{"ALICE", "dev"}, member{"BOB", "ops"})},
		team{name: "web", members: A.From(member{"CAROL", "dev"})},
	)}

	var leftVisits, rightVisits []string
	assert.Equal(t, exp, modifyWith(left, recordUpper(&leftVisits))(src))
	assert.Equal(t, exp, modifyWith(right, recordUpper(&rightVisits))(src))

	assert.Equal(t, A.From("alice", "bob", "carol"), leftVisits, "visits nested foci in order")
	assert.Equal(t, leftVisits, rightVisits)

	assert.Equal(t, "alice", src.teams[0].members[0].name, "does not mutate the source")
}

func TestComposeWithMonoid(t *testing.T) {
	src := team{name: "core", members: A.From(member{"alice", "dev"}, member{"bob", "ops"})}
	composeMembers := Compose[team, string, Endomorphism[team], string](
		identity.Map[Endomorphism[member], member],
	)

	t.Run("concat combines a field with a nested traversal", func(t *testing.T) {
		m := MakeMonoid[string, string, team](
			identity.Of[Endomorphism[team]],
			identity.Map[Endomorphism[team], Endomorphism[Endomorphism[team]]],
			identity.Ap[Endomorphism[team], Endomorphism[team]],
		)
		teamName := TLG.FromLens[team, string](identity.Map[string, Endomorphism[team]])(teamNameLens())

		var visits []string
		res := modifyWith(m.Concat(teamName, memberNamesIdentity()), recordUpper(&visits))(src)

		assert.Equal(t, team{name: "CORE", members: A.From(member{"ALICE", "dev"}, member{"BOB", "ops"})}, res)
		assert.Equal(t, A.From("core", "alice", "bob"), visits)
	})

	t.Run("compose accepts a concatenated inner traversal", func(t *testing.T) {
		m := identityMonoid()
		inner := m.Concat(identityField(memberNameLens()), identityField(memberRoleLens()))

		var visits []string
		res := modifyWith(composeMembers(inner)(teamMembersIdentity()), recordUpper(&visits))(src)

		assert.Equal(t, team{name: "core", members: A.From(member{"ALICE", "DEV"}, member{"BOB", "OPS"})}, res)
		assert.Equal(t, A.From("alice", "dev", "bob", "ops"), visits)
	})

	t.Run("an empty inner traversal focuses on nothing", func(t *testing.T) {
		var visits []string
		res := modifyWith(composeMembers(identityMonoid().Empty())(teamMembersIdentity()), recordUpper(&visits))(src)

		assert.Equal(t, src, res)
		assert.Empty(t, visits)
	})

	t.Run("an empty outer traversal focuses on nothing", func(t *testing.T) {
		outer := Empty[member, member, team](identity.Of[Endomorphism[team]])

		var visits []string
		res := modifyWith(composeMembers(identityField(memberNameLens()))(outer), recordUpper(&visits))(src)

		assert.Equal(t, src, res)
		assert.Empty(t, visits)
	})
}
