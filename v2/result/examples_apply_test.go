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

package result_test

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	R "github.com/IBM/fp-go/v2/result"
)

func ExampleApplySemigroup() {
	intAdd := N.MonoidSum[int]()
	eitherSemi := R.ApplySemigroup(intAdd)
	result := eitherSemi.Concat(R.Of(2), R.Of(3))
	fmt.Println(R.GetOrElse(F.Constant1[error](0))(result))
	// Output:
	// 5
}

func ExampleApplicativeMonoid() {
	intAddMonoid := N.MonoidSum[int]()
	eitherMon := R.ApplicativeMonoid(intAddMonoid)
	empty := eitherMon.Empty()
	fmt.Println(R.GetOrElse(F.Constant1[error](0))(empty))
	// Output:
	// 0
}
