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

package result

import (
	"errors"
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
)

// ExampleChainFirstLeft logs an error and keeps it.
func ExampleChainFirstLeft() {
	logged := F.Pipe1(
		Left[int](errors.New("connection refused")),
		ChainFirstLeft[int](func(err error) Result[F.Void] {
			fmt.Println("failed:", err)
			return Of(F.VOID)
		}),
	)

	fmt.Println(logged)
	// Output:
	// failed: connection refused
	// Left[*errors.errorString](connection refused)
}

// ExampleTapLeft logs an error and keeps it.
func ExampleTapLeft() {
	logged := F.Pipe1(
		Left[int](errors.New("connection refused")),
		TapLeft[int](func(err error) Result[F.Void] {
			fmt.Println("failed:", err)
			return Of(F.VOID)
		}),
	)

	fmt.Println(logged)
	// Output:
	// failed: connection refused
	// Left[*errors.errorString](connection refused)
}

// ExampleMonadChainFirstLeft logs an error and keeps it.
func ExampleMonadChainFirstLeft() {
	logged := MonadChainFirstLeft(Left[int](errors.New("connection refused")), func(err error) Result[F.Void] {
		fmt.Println("failed:", err)
		return Of(F.VOID)
	})

	fmt.Println(logged)
	// Output:
	// failed: connection refused
	// Left[*errors.errorString](connection refused)
}

// ExampleMonadTapLeft logs an error and keeps it.
func ExampleMonadTapLeft() {
	logged := MonadTapLeft(Left[int](errors.New("connection refused")), func(err error) Result[F.Void] {
		fmt.Println("failed:", err)
		return Of(F.VOID)
	})

	fmt.Println(logged)
	// Output:
	// failed: connection refused
	// Left[*errors.errorString](connection refused)
}
