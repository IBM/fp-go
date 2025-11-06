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

package function

import (
	"fmt"
)

func addSthg(value int) int {
	return value + 1
}

func doSthgElse(value int) int {
	return value * 2
}

func doFinalSthg(value int) string {
	return fmt.Sprintf("final value: %d", value)
}

func Example() {
	// start point
	value := 1
	// imperative style
	value1 := addSthg(value)                    // 2
	value2 := doSthgElse(value1)                // 4
	finalValueImperative := doFinalSthg(value2) // "final value: 4"

	// the same but inline
	finalValueInline := doFinalSthg(doSthgElse(addSthg(value)))

	// with pipe
	finalValuePipe := Pipe3(value, addSthg, doSthgElse, doFinalSthg)

	// with flow
	transform := Flow3(addSthg, doSthgElse, doFinalSthg)
	finalValueFlow := transform(value)

	fmt.Println(finalValueImperative)
	fmt.Println(finalValueInline)
	fmt.Println(finalValuePipe)
	fmt.Println(finalValueFlow)

	// Output:
	// final value: 4
	// final value: 4
	// final value: 4
	// final value: 4
}
