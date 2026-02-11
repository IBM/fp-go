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
	"github.com/IBM/fp-go/v2/readerioeither"
)

// ChainConsumer chains a consumer (side-effect function) into a ReaderIOResult computation,
// replacing the success value with Void (empty struct).
//
// This is useful for performing side effects (like logging, printing, or writing to a file)
// where you don't need to preserve the original value. The consumer is only executed if the
// computation succeeds; if it fails with an error, the consumer is skipped.
//
// Type parameters:
//   - R: The context/environment type
//   - A: The value type to consume
//
// Parameters:
//   - c: A consumer function that performs a side effect on the value
//
// Returns:
//
//	An Operator that executes the consumer and returns Void on success
//
// Example:
//
//	import (
//	    "context"
//	    "fmt"
//	    RIO "github.com/IBM/fp-go/v2/readerioresult"
//	)
//
//	// Log a value and discard it
//	logValue := RIO.ChainConsumer[context.Context](func(x int) {
//	    fmt.Printf("Value: %d\n", x)
//	})
//
//	computation := F.Pipe1(
//	    RIO.Of[context.Context](42),
//	    logValue,
//	)
//	// Prints "Value: 42" and returns result.Of(struct{}{})
//
//go:inline
func ChainConsumer[R, A any](c Consumer[A]) Operator[R, A, Void] {
	return readerioeither.ChainConsumer[R, error](c)
}

// ChainFirstConsumer chains a consumer into a ReaderIOResult computation while preserving
// the original value.
//
// This is useful for performing side effects (like logging, printing, or metrics collection)
// where you want to keep the original value for further processing. The consumer is only
// executed if the computation succeeds; if it fails with an error, the consumer is skipped
// and the error is propagated.
//
// Type parameters:
//   - R: The context/environment type
//   - A: The value type to consume and preserve
//
// Parameters:
//   - c: A consumer function that performs a side effect on the value
//
// Returns:
//
//	An Operator that executes the consumer and returns the original value on success
//
// Example:
//
//	import (
//	    "context"
//	    "fmt"
//	    F "github.com/IBM/fp-go/v2/function"
//	    N "github.com/IBM/fp-go/v2/number"
//	    RIO "github.com/IBM/fp-go/v2/readerioresult"
//	)
//
//	// Log a value but keep it for further processing
//	logValue := RIO.ChainFirstConsumer[context.Context](func(x int) {
//	    fmt.Printf("Processing: %d\n", x)
//	})
//
//	computation := F.Pipe2(
//	    RIO.Of[context.Context](10),
//	    logValue,
//	    RIO.Map[context.Context](N.Mul(2)),
//	)
//	// Prints "Processing: 10" and returns result.Of(20)
//
//go:inline
func ChainFirstConsumer[R, A any](c Consumer[A]) Operator[R, A, A] {
	return readerioeither.ChainFirstConsumer[R, error](c)
}
