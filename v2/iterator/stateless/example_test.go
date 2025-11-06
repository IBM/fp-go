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

package stateless

import (
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

func Example_any() {
	// `Any` function that simply returns the boolean identity
	anyBool := Any(F.Identity[bool])

	fmt.Println(anyBool(FromArray(A.From(true, false, false))))
	fmt.Println(anyBool(FromArray(A.From(false, false, false))))
	fmt.Println(anyBool(Empty[bool]()))

	// Output:
	// true
	// false
	// false
}

func Example_next() {

	seq := MakeBy(F.Identity[int])

	first := seq()

	value := F.Pipe1(
		first,
		O.Map(Current[int]),
	)

	fmt.Println(value)

	// Output:
	// Some[int](0)
}
