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

package statereaderioresult

import (
	"context"
	"errors"
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
)

// ExampleChainFirstLeft logs an error and keeps it.
func ExampleChainFirstLeft() {
	logged := F.Pipe1(
		Left[counter, int](errors.New("connection refused")),
		ChainFirstLeft[int](func(err error) StateReaderIOResult[counter, F.Void] {
			fmt.Println("failed:", err)
			return Of[counter](F.VOID)
		}),
	)

	fmt.Println(result.Map(pair.Tail[counter, int])(logged(0)(context.Background())()))
	// Output:
	// failed: connection refused
	// Left[*errors.errorString](connection refused)
}

// ExampleTapLeft logs an error and keeps it.
func ExampleTapLeft() {
	logged := F.Pipe1(
		Left[counter, int](errors.New("connection refused")),
		TapLeft[int](func(err error) StateReaderIOResult[counter, F.Void] {
			fmt.Println("failed:", err)
			return Of[counter](F.VOID)
		}),
	)

	fmt.Println(result.Map(pair.Tail[counter, int])(logged(0)(context.Background())()))
	// Output:
	// failed: connection refused
	// Left[*errors.errorString](connection refused)
}

// ExampleMonadChainFirstLeft logs an error and keeps it.
func ExampleMonadChainFirstLeft() {
	logged := MonadChainFirstLeft(Left[counter, int](errors.New("connection refused")), func(err error) StateReaderIOResult[counter, F.Void] {
		fmt.Println("failed:", err)
		return Of[counter](F.VOID)
	})

	fmt.Println(result.Map(pair.Tail[counter, int])(logged(0)(context.Background())()))
	// Output:
	// failed: connection refused
	// Left[*errors.errorString](connection refused)
}

// ExampleMonadTapLeft logs an error and keeps it.
func ExampleMonadTapLeft() {
	logged := MonadTapLeft(Left[counter, int](errors.New("connection refused")), func(err error) StateReaderIOResult[counter, F.Void] {
		fmt.Println("failed:", err)
		return Of[counter](F.VOID)
	})

	fmt.Println(result.Map(pair.Tail[counter, int])(logged(0)(context.Background())()))
	// Output:
	// failed: connection refused
	// Left[*errors.errorString](connection refused)
}

// ExampleChainLeft recovers from a failure with a default port.
func ExampleChainLeft() {
	port := F.Pipe1(
		Left[counter, int](errors.New("connection refused")),
		ChainLeft(func(error) StateReaderIOResult[counter, int] { return Of[counter](8080) }),
	)

	fmt.Println(result.Map(pair.Tail[counter, int])(port(0)(context.Background())()))
	// Output:
	// Right[int](8080)
}

// ExampleMonadChainLeft recovers from a failure with a default port.
func ExampleMonadChainLeft() {
	port := MonadChainLeft(Left[counter, int](errors.New("connection refused")), func(error) StateReaderIOResult[counter, int] { return Of[counter](8080) })

	fmt.Println(result.Map(pair.Tail[counter, int])(port(0)(context.Background())()))
	// Output:
	// Right[int](8080)
}

// ExampleOrElse recovers from a failure with a default port.
func ExampleOrElse() {
	port := F.Pipe1(
		Left[counter, int](errors.New("connection refused")),
		OrElse(func(error) StateReaderIOResult[counter, int] { return Of[counter](8080) }),
	)

	fmt.Println(result.Map(pair.Tail[counter, int])(port(0)(context.Background())()))
	// Output:
	// Right[int](8080)
}

// ExampleMapLeft adds context to an error.
func ExampleMapLeft() {
	port := F.Pipe1(
		Left[counter, int](errors.New("connection refused")),
		MapLeft[counter, int](func(err error) error { return fmt.Errorf("lookup: %w", err) }),
	)

	fmt.Println(result.Map(pair.Tail[counter, int])(port(0)(context.Background())()))
	// Output:
	// Left[*fmt.wrapError](lookup: connection refused)
}

// ExampleMonadMapLeft adds context to an error.
func ExampleMonadMapLeft() {
	port := MonadMapLeft(Left[counter, int](errors.New("connection refused")), func(err error) error {
		return fmt.Errorf("lookup: %w", err)
	})

	fmt.Println(result.Map(pair.Tail[counter, int])(port(0)(context.Background())()))
	// Output:
	// Left[*fmt.wrapError](lookup: connection refused)
}
