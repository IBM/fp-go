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

package identity

import (
	I "github.com/IBM/fp-go/v2/identity"
	G "github.com/IBM/fp-go/v2/optics/lens/traversal/generic"
)

// Compose composes a lens with a traversal to create a new traversal.
//
// This function allows you to focus deeper into a data structure by first using
// a lens to access a field, then using a traversal to access multiple values within
// that field. The result is a traversal that can operate on all the nested values.
//
// The composition follows the pattern: Lens[S, A] → Traversal[A, B] → Traversal[S, B]
// where the lens focuses on field A within structure S, and the traversal focuses on
// multiple B values within A.
//
// Type Parameters:
//   - S: The outer structure type
//   - A: The intermediate field type (target of the lens)
//   - B: The final focus type (targets of the traversal)
//
// Parameters:
//   - t: A traversal that focuses on B values within A
//
// Returns:
//   - A function that takes a Lens[S, A] and returns a Traversal[S, B]
//
// Example:
//
//	import (
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/optics/lens"
//	    LT "github.com/IBM/fp-go/v2/optics/lens/traversal"
//	    AI "github.com/IBM/fp-go/v2/optics/traversal/array/identity"
//	)
//
//	type Team struct {
//	    Name    string
//	    Members []string
//	}
//
//	// Lens to access the Members field
//	membersLens := lens.MakeLens(
//	    func(t Team) []string { return t.Members },
//	    func(t Team, m []string) Team { t.Members = m; return t },
//	)
//
//	// Traversal for array elements
//	arrayTraversal := AI.FromArray[string]()
//
//	// Compose lens with traversal to access all member names
//	memberTraversal := F.Pipe1(
//	    membersLens,
//	    LT.Compose[Team, []string, string](arrayTraversal),
//	)
//
//	team := Team{Name: "Engineering", Members: []string{"Alice", "Bob"}}
//	// Uppercase all member names
//	updated := memberTraversal(strings.ToUpper)(team)
//	// updated.Members: ["ALICE", "BOB"]
//
// See Also:
//   - Lens: A functional reference to a subpart of a data structure
//   - Traversal: A functional reference to multiple subparts
//   - traversal.Compose: Composes two traversals
func Compose[S, A, B any](t Traversal[A, B, A, B]) func(Lens[S, A]) Traversal[S, B, S, B] {
	return G.Compose[S, A, B, S, A, B](
		I.Map,
	)(t)
}
