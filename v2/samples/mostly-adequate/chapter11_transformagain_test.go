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

package mostlyadequate

import (
	"fmt"
	"regexp"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioresult"
	R "github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
)

func findUserByID(id int) IOResult[Chapter08User] {
	switch id {
	case 1:
		return ioresult.Of(albert08)
	case 2:
		return ioresult.Of(gary08)
	case 3:
		return ioresult.Of(theresa08)
	default:
		return ioresult.Left[Chapter08User](fmt.Errorf("user %d not found", id))
	}
}

func Example_solution11A() {
	// eitherToMaybe :: Either b a -> Maybe a
	eitherToMaybe := R.ToOption[string]

	fmt.Println(eitherToMaybe(R.Of("one eyed willy")))
	fmt.Println(eitherToMaybe(R.Left[string](fmt.Errorf("some error"))))

	// Output:
	// Some[string](one eyed willy)
	// None[string]
}

func Example_solution11B() {
	findByNameID := F.Flow2(
		findUserByID,
		ioresult.Map(Chapter08User.getName),
	)

	fmt.Println(findByNameID(1)())
	fmt.Println(findByNameID(2)())
	fmt.Println(findByNameID(3)())
	fmt.Println(findByNameID(4)())

	// Output:
	// Right[string](Albert)
	// Right[string](Gary)
	// Right[string](Theresa)
	// Left[*errors.errorString](user 4 not found)
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
