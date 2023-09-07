// Copyright (c) 2023 IBM Corp.
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

package mostlyadequate

import (
	"fmt"
	"regexp"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	S "github.com/IBM/fp-go/string"
)

func findUserById(id int) IOE.IOEither[error, Chapter08User] {
	switch id {
	case 1:
		return IOE.Of[error](albert08)
	case 2:
		return IOE.Of[error](gary08)
	case 3:
		return IOE.Of[error](theresa08)
	default:
		return IOE.Left[Chapter08User](fmt.Errorf("user %d not found", id))
	}
}

func Example_solution11A() {
	// eitherToMaybe :: Either b a -> Maybe a
	eitherToMaybe := E.ToOption[error, string]

	fmt.Println(eitherToMaybe(E.Of[error]("one eyed willy")))
	fmt.Println(eitherToMaybe(E.Left[string](fmt.Errorf("some error"))))

	// Output:
	// Some[string](one eyed willy)
	// None[string]
}

func Example_solution11B() {
	findByNameId := F.Flow2(
		findUserById,
		IOE.Map[error](Chapter08User.getName),
	)

	fmt.Println(findByNameId(1)())
	fmt.Println(findByNameId(2)())
	fmt.Println(findByNameId(3)())
	fmt.Println(findByNameId(4)())

	// Output:
	// Right[<nil>, string](Albert)
	// Right[<nil>, string](Gary)
	// Right[<nil>, string](Theresa)
	// Left[*errors.errorString, string](user 4 not found)
}

func Example_solution11C() {
	// strToList :: String -> [Char
	strToList := Split(regexp.MustCompile(``))

	// listToStr :: [Char] -> String
	listToStr := A.Intercalate(S.Monoid)("")

	sortLetters := F.Flow3(
		strToList,
		A.Sort(S.Ord),
		listToStr,
	)

	fmt.Println(sortLetters("sortme"))

	// Output:
	// emorst
}
