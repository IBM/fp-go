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
	"strings"
	"testing"

	AR "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/lens"
	AI "github.com/IBM/fp-go/v2/optics/traversal/array/identity"
	"github.com/stretchr/testify/assert"
)

type Team struct {
	Name    string
	Members []string
}

type Company struct {
	Name  string
	Teams []Team
}

func TestCompose_Success(t *testing.T) {
	t.Run("composes lens with array traversal to modify nested values", func(t *testing.T) {
		// Arrange
		membersLens := lens.MakeLens(
			func(team Team) []string { return team.Members },
			func(team Team, members []string) Team {
				team.Members = members
				return team
			},
		)
		arrayTraversal := AI.FromArray[string]()

		memberTraversal := F.Pipe1(
			membersLens,
			Compose[Team](arrayTraversal),
		)

		team := Team{
			Name:    "Engineering",
			Members: []string{"alice", "bob", "charlie"},
		}

		// Act - uppercase all member names
		result := memberTraversal(strings.ToUpper)(team)

		// Assert
		expected := Team{
			Name:    "Engineering",
			Members: []string{"ALICE", "BOB", "CHARLIE"},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("composes lens with array traversal on empty array", func(t *testing.T) {
		// Arrange
		membersLens := lens.MakeLens(
			func(team Team) []string { return team.Members },
			func(team Team, members []string) Team {
				team.Members = members
				return team
			},
		)
		arrayTraversal := AI.FromArray[string]()

		memberTraversal := F.Pipe1(
			membersLens,
			Compose[Team](arrayTraversal),
		)

		team := Team{
			Name:    "Engineering",
			Members: []string{},
		}

		// Act
		result := memberTraversal(strings.ToUpper)(team)

		// Assert
		assert.Equal(t, team, result)
	})

	t.Run("composes lens with array traversal to transform numbers", func(t *testing.T) {
		// Arrange
		type Stats struct {
			Name   string
			Scores []int
		}

		scoresLens := lens.MakeLens(
			func(s Stats) []int { return s.Scores },
			func(s Stats, scores []int) Stats {
				s.Scores = scores
				return s
			},
		)
		arrayTraversal := AI.FromArray[int]()

		scoreTraversal := F.Pipe1(
			scoresLens,
			Compose[Stats, []int, int](arrayTraversal),
		)

		stats := Stats{
			Name:   "Player1",
			Scores: []int{10, 20, 30},
		}

		// Act - double all scores
		result := scoreTraversal(func(n int) int { return n * 2 })(stats)

		// Assert
		expected := Stats{
			Name:   "Player1",
			Scores: []int{20, 40, 60},
		}
		assert.Equal(t, expected, result)
	})
}

func TestCompose_Integration(t *testing.T) {
	t.Run("composes multiple lenses and traversals", func(t *testing.T) {
		// Arrange - nested structure with Company -> Teams -> Members
		teamsLens := lens.MakeLens(
			func(c Company) []Team { return c.Teams },
			func(c Company, teams []Team) Company {
				c.Teams = teams
				return c
			},
		)

		// First compose: Company -> []Team -> Team
		teamArrayTraversal := AI.FromArray[Team]()
		companyToTeamTraversal := F.Pipe1(
			teamsLens,
			Compose[Company, []Team, Team](teamArrayTraversal),
		)

		// Second compose: Team -> []string -> string
		membersLens := lens.MakeLens(
			func(team Team) []string { return team.Members },
			func(team Team, members []string) Team {
				team.Members = members
				return team
			},
		)
		memberArrayTraversal := AI.FromArray[string]()
		teamToMemberTraversal := F.Pipe1(
			membersLens,
			Compose[Team](memberArrayTraversal),
		)

		company := Company{
			Name: "TechCorp",
			Teams: []Team{
				{Name: "Engineering", Members: []string{"alice", "bob"}},
				{Name: "Design", Members: []string{"charlie", "diana"}},
			},
		}

		// Act - uppercase all members in all teams
		// First traverse to teams, then for each team traverse to members
		result := companyToTeamTraversal(func(team Team) Team {
			return teamToMemberTraversal(strings.ToUpper)(team)
		})(company)

		// Assert
		expected := Company{
			Name: "TechCorp",
			Teams: []Team{
				{Name: "Engineering", Members: []string{"ALICE", "BOB"}},
				{Name: "Design", Members: []string{"CHARLIE", "DIANA"}},
			},
		}
		assert.Equal(t, expected, result)
	})
}

func TestCompose_EdgeCases(t *testing.T) {
	t.Run("preserves structure name when modifying members", func(t *testing.T) {
		// Arrange
		membersLens := lens.MakeLens(
			func(team Team) []string { return team.Members },
			func(team Team, members []string) Team {
				team.Members = members
				return team
			},
		)
		arrayTraversal := AI.FromArray[string]()

		memberTraversal := F.Pipe1(
			membersLens,
			Compose[Team](arrayTraversal),
		)

		team := Team{
			Name:    "Engineering",
			Members: []string{"alice"},
		}

		// Act
		result := memberTraversal(strings.ToUpper)(team)

		// Assert - Name should be unchanged
		assert.Equal(t, "Engineering", result.Name)
		assert.Equal(t, AR.From("ALICE"), result.Members)
	})

	t.Run("handles identity transformation", func(t *testing.T) {
		// Arrange
		membersLens := lens.MakeLens(
			func(team Team) []string { return team.Members },
			func(team Team, members []string) Team {
				team.Members = members
				return team
			},
		)
		arrayTraversal := AI.FromArray[string]()

		memberTraversal := F.Pipe1(
			membersLens,
			Compose[Team](arrayTraversal),
		)

		team := Team{
			Name:    "Engineering",
			Members: []string{"alice", "bob"},
		}

		// Act - apply identity function
		result := memberTraversal(F.Identity[string])(team)

		// Assert - should be unchanged
		assert.Equal(t, team, result)
	})
}
