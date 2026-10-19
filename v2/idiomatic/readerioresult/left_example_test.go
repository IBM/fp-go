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

package readerioresult

import (
	"errors"
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
)

// ExampleChainFirstLeft logs an error and keeps it.
func ExampleChainFirstLeft() {
	logged := F.Pipe1(
		Left[config, int](errors.New("connection refused")),
		ChainFirstLeft[int](func(err error) ReaderIOResult[config, F.Void] {
			fmt.Println("failed:", err)
			return Of[config](F.VOID)
		}),
	)

	fmt.Println(logged(config{})())
	// Output:
	// failed: connection refused
	// 0 connection refused
}

// ExampleTapLeft logs an error and keeps it.
func ExampleTapLeft() {
	logged := F.Pipe1(
		Left[config, int](errors.New("connection refused")),
		TapLeft[int](func(err error) ReaderIOResult[config, F.Void] {
			fmt.Println("failed:", err)
			return Of[config](F.VOID)
		}),
	)

	fmt.Println(logged(config{})())
	// Output:
	// failed: connection refused
	// 0 connection refused
}

// ExampleMonadChainFirstLeft logs an error and keeps it.
func ExampleMonadChainFirstLeft() {
	logged := MonadChainFirstLeft(Left[config, int](errors.New("connection refused")), func(err error) ReaderIOResult[config, F.Void] {
		fmt.Println("failed:", err)
		return Of[config](F.VOID)
	})

	fmt.Println(logged(config{})())
	// Output:
	// failed: connection refused
	// 0 connection refused
}

// ExampleMonadTapLeft logs an error and keeps it.
func ExampleMonadTapLeft() {
	logged := MonadTapLeft(Left[config, int](errors.New("connection refused")), func(err error) ReaderIOResult[config, F.Void] {
		fmt.Println("failed:", err)
		return Of[config](F.VOID)
	})

	fmt.Println(logged(config{})())
	// Output:
	// failed: connection refused
	// 0 connection refused
}

// ExampleChainLeft recovers from a failure with a default port.
func ExampleChainLeft() {
	port := F.Pipe1(
		Left[config, int](errors.New("connection refused")),
		ChainLeft(func(error) ReaderIOResult[config, int] { return Of[config](8080) }),
	)

	fmt.Println(port(config{})())
	// Output:
	// 8080 <nil>
}

// ExampleMonadChainLeft recovers from a failure with a default port.
func ExampleMonadChainLeft() {
	port := MonadChainLeft(Left[config, int](errors.New("connection refused")), func(error) ReaderIOResult[config, int] { return Of[config](8080) })

	fmt.Println(port(config{})())
	// Output:
	// 8080 <nil>
}

// ExampleOrElse recovers from a failure with a default port.
func ExampleOrElse() {
	port := F.Pipe1(
		Left[config, int](errors.New("connection refused")),
		OrElse(func(error) ReaderIOResult[config, int] { return Of[config](8080) }),
	)

	fmt.Println(port(config{})())
	// Output:
	// 8080 <nil>
}

// ExampleMapLeft adds context to an error.
func ExampleMapLeft() {
	port := F.Pipe1(
		Left[config, int](errors.New("connection refused")),
		MapLeft[config, int](func(err error) error { return fmt.Errorf("lookup: %w", err) }),
	)

	fmt.Println(port(config{})())
	// Output:
	// 0 lookup: connection refused
}

// ExampleMonadMapLeft adds context to an error.
func ExampleMonadMapLeft() {
	port := MonadMapLeft(Left[config, int](errors.New("connection refused")), func(err error) error {
		return fmt.Errorf("lookup: %w", err)
	})

	fmt.Println(port(config{})())
	// Output:
	// 0 lookup: connection refused
}
