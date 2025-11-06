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
	I "github.com/IBM/fp-go/v2/number/integer"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/ord"
	S "github.com/IBM/fp-go/v2/string"
)

var (
	Exclaim   = S.Format[string]("%s!")
	Shout     = F.Flow2(ToUpper, Exclaim)
	Dasherize = F.Flow4(
		Replace(regexp.MustCompile(`\s{2,}`))(" "),
		Split(regexp.MustCompile(` `)),
		A.Map(ToLower),
		A.Intercalate(S.Monoid)("-"),
	)
)

func Example_shout() {
	fmt.Println(Shout("send in the clowns"))

	// Output: SEND IN THE CLOWNS!
}

func Example_dasherize() {
	fmt.Println(Dasherize("The world is a vampire"))

	// Output: the-world-is-a-vampire
}

func Example_pipe() {
	output := F.Pipe2(
		"send in the clowns",
		ToUpper,
		Exclaim,
	)

	fmt.Println(output)

	// Output: SEND IN THE CLOWNS!
}

func Example_solution05A() {
	IsLastInStock := F.Flow2(
		A.Last[Car],
		O.Map(Car.getInStock),
	)

	fmt.Println(IsLastInStock(Cars[0:3]))
	fmt.Println(IsLastInStock(Cars[3:]))

	// Output:
	// Some[bool](true)
	// Some[bool](false)
}

func Example_solution05B() {
	// averageDollarValue :: [Car] -> Int
	averageDollarValue := F.Flow2(
		A.Map(Car.getDollarValue),
		average,
	)

	fmt.Println(averageDollarValue(Cars))

	// Output:
	// 790700
}

func Example_solution05C() {
	// order by horsepower
	ordByHorsepower := ord.Contramap(Car.getHorsepower)(I.Ord)

	// fastestCar :: [Car] -> Option[String]
	fastestCar := F.Flow3(
		A.Sort(ordByHorsepower),
		A.Last[Car],
		O.Map(F.Flow2(
			Car.getName,
			S.Format[string]("%s is the fastest"),
		)),
	)

	fmt.Println(fastestCar(Cars))

	// Output:
	// Some[string](Aston Martin One-77 is the fastest)
}
