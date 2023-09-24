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

package examples

import (
	"fmt"

	N "github.com/IBM/fp-go/number"
)

// addInts adds two integers
func addInts(left, right int) int {
	return left + right
}

// addNumbers adds two numbers
func addNumbers[T N.Number](left, right T) T {
	return left + right
}

func Example_generics() {
	// invoke the non generic version
	fmt.Println(addInts(1, 2))

	// invoke the generic version
	fmt.Println(addNumbers(1, 2))
	fmt.Println(addNumbers(1.0, 2.0))

	// Output:
	// 3
	// 3
	// 3
}
