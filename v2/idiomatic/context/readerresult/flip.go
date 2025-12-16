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

package readerresult

import (
	"context"

	RR "github.com/IBM/fp-go/v2/idiomatic/readerresult"
	"github.com/IBM/fp-go/v2/reader"
)

// SequenceReader swaps the order of nested environment parameters when the inner type is a Reader.
//
// It transforms ReaderResult[Reader[R, A]] into a function that takes context.Context first,
// then R, and returns (A, error). This is useful when you have a ReaderResult computation
// that produces a Reader, and you want to sequence the environment dependencies.
//
// Type Parameters:
//   - R: The inner Reader's environment type
//   - A: The final result type
//
// Parameters:
//   - ma: A ReaderResult that produces a Reader[R, A]
//
// Returns:
//   - A Kleisli arrow that takes context.Context and R to produce (A, error)
//
// Example:
//
//	type Config struct {
//	    DatabaseURL string
//	}
//
//	// Returns a ReaderResult that produces a Reader
//	getDBReader := func(ctx context.Context) (reader.Reader[Config, string], error) {
//	    return func(cfg Config) string {
//	        return cfg.DatabaseURL
//	    }, nil
//	}
//
//	// Sequence the environments: context.Context -> Config -> string
//	sequenced := readerresult.SequenceReader[Config, string](getDBReader)
//	result, err := sequenced(ctx)(config)
//
//go:inline
func SequenceReader[R, A any](ma ReaderResult[Reader[R, A]]) Kleisli[R, A] {
	return WithContextK(RR.SequenceReader(ma))
}

// TraverseReader combines SequenceReader with a Kleisli arrow transformation.
//
// It takes a Reader Kleisli arrow (a function from A to Reader[R, B]) and returns
// a function that transforms ReaderResult[A] into a Kleisli arrow from context.Context
// and R to B. This is useful for transforming values within a ReaderResult while
// introducing an additional Reader dependency.
//
// Type Parameters:
//   - R: The Reader's environment type
//   - A: The input type
//   - B: The output type
//
// Parameters:
//   - f: A Kleisli arrow that transforms A into Reader[R, B]
//
// Returns:
//   - A function that transforms ReaderResult[A] into a Kleisli arrow from context.Context and R to B
//
// Example:
//
//	type Config struct {
//	    Multiplier int
//	}
//
//	// A Kleisli arrow that uses Config to transform int to string
//	formatWithConfig := func(n int) reader.Reader[Config, string] {
//	    return func(cfg Config) string {
//	        return fmt.Sprintf("Value: %d", n * cfg.Multiplier)
//	    }
//	}
//
//	// Create a ReaderResult[int]
//	getValue := readerresult.Of[int](42)
//
//	// Traverse: transform the int using the Reader Kleisli arrow
//	traversed := readerresult.TraverseReader[Config](formatWithConfig)(getValue)
//	result, err := traversed(ctx)(Config{Multiplier: 2})
//	// result == "Value: 84"
//
//go:inline
func TraverseReader[R, A, B any](
	f reader.Kleisli[R, A, B],
) func(ReaderResult[A]) Kleisli[R, B] {
	return RR.TraverseReader[context.Context](f)
}
