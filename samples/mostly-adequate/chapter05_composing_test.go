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
	F "github.com/IBM/fp-go/function"
	S "github.com/IBM/fp-go/string"
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
