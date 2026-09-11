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
	TLG "github.com/IBM/fp-go/v2/optics/traversalendo/lens/generic"
	"github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestToTraversalIdentity(t *testing.T) {
	nameTrav := ToTraversal[string, string, member](
		identity.Map[Endomorphism[member], member],
	)(TLG.FromLens[member, string](identity.Map[string, Endomorphism[member]])(memberNameLens()))

	t.Run("modifies the focus of the source", func(t *testing.T) {
		assert.Equal(t, member{"ALICE", "dev"}, nameTrav(strings.ToUpper)(member{"alice", "dev"}))
	})

	t.Run("traversing with identity is the identity", func(t *testing.T) {
		src := member{"alice", "dev"}
		assert.Equal(t, src, nameTrav(identity.Of[string])(src))
	})

	t.Run("converts a composed traversal endomorphism", func(t *testing.T) {
		teamTrav := ToTraversal[string, string, team](
			identity.Map[Endomorphism[team], team],
		)(memberNamesIdentity())

		src := team{name: "core", members: A.From(member{"alice", "dev"}, member{"bob", "ops"})}
		exp := team{name: "core", members: A.From(member{"ALICE", "dev"}, member{"BOB", "ops"})}
		assert.Equal(t, exp, teamTrav(strings.ToUpper)(src))
	})
}

func TestToTraversalOption(t *testing.T) {
	nameTrav := ToTraversal[string, option.Option[string], member](
		option.Map[Endomorphism[member], member],
	)(TLG.FromLens[member, string](option.Map[string, Endomorphism[member]])(memberNameLens()))

	t.Run("propagates a successful effect", func(t *testing.T) {
		assert.Equal(t, option.Some(member{"ALICE", "dev"}), nameTrav(nonEmptyUpper)(member{"alice", "dev"}))
	})

	t.Run("propagates a failing effect", func(t *testing.T) {
		assert.Equal(t, option.None[member](), nameTrav(nonEmptyUpper)(member{"", "dev"}))
	})

	t.Run("converts a composed traversal endomorphism", func(t *testing.T) {
		teamTrav := ToTraversal[string, option.Option[string], team](
			option.Map[Endomorphism[team], team],
		)(memberNamesOption())

		src := team{name: "core", members: A.From(member{"alice", "dev"}, member{"bob", "ops"})}
		exp := team{name: "core", members: A.From(member{"ALICE", "dev"}, member{"BOB", "ops"})}
		assert.Equal(t, option.Some(exp), teamTrav(nonEmptyUpper)(src))

		failing := team{name: "core", members: A.From(member{"alice", "dev"}, member{"", "ops"})}
		assert.Equal(t, option.None[team](), teamTrav(nonEmptyUpper)(failing))
	})
}

func TestToTraversalEmpty(t *testing.T) {
	t.Run("returns the source unchanged without invoking the function", func(t *testing.T) {
		trv := ToTraversal[string, string, member](
			identity.Map[Endomorphism[member], member],
		)(Empty[string, string, member](identity.Of[Endomorphism[member]]))

		calls := 0
		count := func(s string) string {
			calls++
			return s
		}
		src := member{"alice", "dev"}
		assert.Equal(t, src, trv(count)(src))
		assert.Equal(t, 0, calls)
	})

	t.Run("lifts the source into the effect", func(t *testing.T) {
		trv := ToTraversal[string, option.Option[string], member](
			option.Map[Endomorphism[member], member],
		)(Empty[string, option.Option[string], member](option.Of[Endomorphism[member]]))

		// nonEmptyUpper would fail on both fields, but empty focuses on neither
		src := member{"", ""}
		assert.Equal(t, option.Some(src), trv(nonEmptyUpper)(src))
	})
}

func TestToTraversalPreservesEffectOrder(t *testing.T) {
	field := TLG.FromLens[member, string](mapLogged[string, Endomorphism[member]])
	both := Concat[string, logged[string], member](
		mapLogged[Endomorphism[member], Endomorphism[Endomorphism[member]]],
		apLogged[Endomorphism[member], Endomorphism[member]],
	)(field(memberNameLens()), field(memberRoleLens()))

	trv := ToTraversal[string, logged[string], member](mapLogged[Endomorphism[member], member])(both)

	res := trv(visit)(member{"alice", "dev"})
	assert.Equal(t, member{"ALICE", "DEV"}, res.val)
	assert.Equal(t, A.From("alice", "dev"), res.log)
}
