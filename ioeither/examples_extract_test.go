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

package ioeither

import (
	"fmt"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IO "github.com/IBM/fp-go/io"
)

func ExampleIOEither_extraction() {
	// IOEither
	someIOEither := Right[error](42)
	eitherValue := someIOEither()                            // E.Right(42)
	value := E.GetOrElse(F.Constant1[error](0))(eitherValue) // 42

	// Or more directly
	infaillibleIO := GetOrElse(F.Constant1[error](IO.Of(0)))(someIOEither) // => IO.Right(42)
	valueFromIO := infaillibleIO()                                         // => 42

	fmt.Println(eitherValue)
	fmt.Println(value)
	fmt.Println(valueFromIO)

	// Output:
	// Right[int](42)
	// 42
	// 42

}
