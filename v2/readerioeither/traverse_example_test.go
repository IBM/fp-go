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

package readerioeither_test

import (
	"fmt"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	RIE "github.com/IBM/fp-go/v2/readerioeither"
	S "github.com/IBM/fp-go/v2/string"
)

type Config struct {
	APIKey  string
	BaseURL string
}

type User struct {
	ID   int
	Name string
}

// Simulate fetching a user by ID
func fetchUser(id int) RIE.ReaderIOEither[Config, error, User] {
	return func(cfg Config) func() E.Either[error, User] {
		return func() E.Either[error, User] {
			if S.IsEmpty(cfg.APIKey) {
				return E.Left[User](fmt.Errorf("missing API key"))
			}
			if id <= 0 {
				return E.Left[User](fmt.Errorf("invalid user ID: %d", id))
			}
			// Simulate successful fetch
			return E.Right[error](User{ID: id, Name: fmt.Sprintf("User%d", id)})
		}
	}
}

// Example of TraverseArray - fetch multiple users
func ExampleTraverseArray() {
	cfg := Config{APIKey: "secret", BaseURL: "https://api.example.com"}

	// Fetch users with IDs 1, 2, 3
	userIDs := []int{1, 2, 3}

	fetchUsers := RIE.TraverseArray(fetchUser)
	result := fetchUsers(userIDs)(cfg)()

	E.Fold(
		func(err error) string {
			fmt.Println("Failed to fetch users")
			return ""
		},
		func(users []User) string {
			fmt.Println("Successfully fetched all users")
			return ""
		},
	)(result)
	// Output: Successfully fetched all users
}

// Example of TraverseArrayWithIndex - process items with their positions
func ExampleTraverseArrayWithIndex() {
	cfg := Config{APIKey: "secret"}

	items := []string{"apple", "banana", "cherry"}

	processWithIndex := RIE.TraverseArrayWithIndex(func(i int, item string) RIE.ReaderIOEither[Config, error, string] {
		return RIE.Of[Config, error](fmt.Sprintf("%d: %s", i+1, item))
	})

	result := processWithIndex(items)(cfg)()

	E.Fold(
		func(err error) int {
			fmt.Printf("Error: %v\n", err)
			return 0
		},
		func(processed []string) int {
			for _, item := range processed {
				fmt.Println(item)
			}
			return len(processed)
		},
	)(result)
	// Output:
	// 1: apple
	// 2: banana
	// 3: cherry
}

// Example of SequenceArray - execute multiple independent computations
func ExampleSequenceArray() {
	cfg := Config{APIKey: "secret"}

	computations := []RIE.ReaderIOEither[Config, error, int]{
		RIE.Of[Config, error](10),
		RIE.Of[Config, error](20),
		RIE.Of[Config, error](30),
	}

	result := RIE.SequenceArray(computations)(cfg)()

	E.Fold(
		func(err error) string {
			fmt.Printf("Error: %v\n", err)
			return ""
		},
		func(values []int) string {
			sum := 0
			for _, v := range values {
				sum += v
			}
			fmt.Printf("Sum: %d\n", sum)
			return ""
		},
	)(result)
	// Output: Sum: 60
}

// Example showing error handling with TraverseArray
func ExampleTraverseArray_errorHandling() {
	cfg := Config{APIKey: "secret"}

	// One invalid ID will cause the entire operation to fail
	userIDs := []int{1, -1, 3} // -1 is invalid

	fetchUsers := RIE.TraverseArray(fetchUser)
	result := fetchUsers(userIDs)(cfg)()

	E.Fold(
		func(err error) string {
			fmt.Printf("Failed: %v\n", err)
			return ""
		},
		func(users []User) string {
			fmt.Printf("Success: fetched %d users\n", len(users))
			return ""
		},
	)(result)
	// Output: Failed: invalid user ID: -1
}

// Example of combining Traverse with other operations
func ExampleTraverseArray_pipeline() {
	cfg := Config{APIKey: "secret"}

	// Pipeline: fetch users, then extract their names
	userIDs := []int{1, 2}

	pipeline := F.Pipe2(
		userIDs,
		RIE.TraverseArray(fetchUser),
		RIE.Map[Config, error](func(users []User) []string {
			names := make([]string, len(users))
			for i, user := range users {
				names[i] = user.Name
			}
			return names
		}),
	)

	result := pipeline(cfg)()

	E.Fold(
		func(err error) string {
			fmt.Printf("Error: %v\n", err)
			return ""
		},
		func(names []string) string {
			fmt.Printf("User names: %v\n", names)
			return ""
		},
	)(result)
	// Output: User names: [User1 User2]
}
