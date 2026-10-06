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

package effect

import (
	"context"
	"errors"
	"fmt"

	"github.com/IBM/fp-go/v2/tailrec"
)

// exRecConfig is the context of the TailRec examples.
type exRecConfig struct {
	Step int
}

// exCountdown decrements n by the configured step until it reaches zero.
func exCountdown(n int) Effect[exRecConfig, Trampoline[int, string]] {
	if n <= 0 {
		return Of[exRecConfig](tailrec.Land[int]("liftoff"))
	}
	return Asks(func(c exRecConfig) Trampoline[int, string] {
		return tailrec.Bounce[string](n - c.Step)
	})
}

// ExampleTailRec demonstrates a loop with a million iterations that does not
// grow the stack. Each step reads the decrement from the context.
func ExampleTailRec() {
	countdown := TailRec(exCountdown)

	res, err := RunSync(Provide[string](exRecConfig{Step: 1})(countdown(1_000_000)))(context.Background())

	fmt.Println(res, err)
	// Output: liftoff <nil>
}

// ExampleTailRec_failure demonstrates that a failing step stops the recursion
// and its error becomes the result of the whole loop.
func ExampleTailRec_failure() {
	countdown := TailRec(func(n int) Effect[exRecConfig, Trampoline[int, string]] {
		if n == 3 {
			return Fail[exRecConfig, Trampoline[int, string]](errors.New("abort at 3"))
		}
		return exCountdown(n)
	})

	_, err := RunSync(Provide[string](exRecConfig{Step: 1})(countdown(10)))(context.Background())

	fmt.Println(err)
	// Output: abort at 3
}
