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

package lazy

import (
	"fmt"
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
)

func ExampleLazy_creation() {
	// lazy function of a constant value
	val := Of(42)
	// create another function to transform this
	valS := F.Pipe1(
		val,
		Map(strconv.Itoa),
	)

	fmt.Println(valS())

	// Output:
	// 42
}
