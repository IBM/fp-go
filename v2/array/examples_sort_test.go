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

package array

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/number/integer"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/ord"
	S "github.com/IBM/fp-go/v2/string"
)

type user struct {
	name string
	age  O.Option[int]
}

func (user user) GetName() string {
	return user.name
}

func (user user) GetAge() O.Option[int] {
	return user.age
}

// Example_sort adapts examples from [https://github.com/inato/fp-ts-cheatsheet#sort-elements-with-ord]
func Example_sort() {

	strings := From("zyx", "abc", "klm")

	sortedStrings := F.Pipe1(
		strings,
		Sort(S.Ord),
	) // => ['abc', 'klm', 'zyx']

	// reverse sort
	reverseSortedStrings := F.Pipe1(
		strings,
		Sort(ord.Reverse(S.Ord)),
	) // => ['zyx', 'klm', 'abc']

	// sort Option
	optionalNumbers := From(O.Some(1337), O.None[int](), O.Some(42))

	sortedNums := F.Pipe1(
		optionalNumbers,
		Sort(O.Ord(I.Ord)),
	)

	// complex object with different rules
	byName := F.Pipe1(
		S.Ord,
		ord.Contramap(user.GetName),
	) // ord.Ord[user]

	byAge := F.Pipe1(
		O.Ord(I.Ord),
		ord.Contramap(user.GetAge),
	) // ord.Ord[user]

	sortedUsers := F.Pipe1(
		From(user{name: "a", age: O.Of(30)}, user{name: "d", age: O.Of(10)}, user{name: "c"}, user{name: "b", age: O.Of(10)}),
		SortBy(From(byAge, byName)),
	)

	fmt.Println(sortedStrings)
	fmt.Println(reverseSortedStrings)
	fmt.Println(sortedNums)
	fmt.Println(sortedUsers)

	// Output:
	// [abc klm zyx]
	// [zyx klm abc]
	// [None[int] Some[int](42) Some[int](1337)]
	// [{c {0 false}} {b {10 true}} {d {10 true}} {a {30 true}}]

}
