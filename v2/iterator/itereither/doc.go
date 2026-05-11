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

// Package itereither provides the SeqEither type, which represents an iterator
// that yields Either values, combining iteration with error handling.
//
// SeqEither[E, A] is defined as iter.Seq[Either[E, A]], representing a sequence
// of values that can each be either an error of type E or a success value of type A.
//
// # Fantasy Land Specification
//
// This package implements a monad transformer combining iteration with Either for error handling.
// It follows the Fantasy Land specification for functional programming patterns.
//
// Implemented algebras:
//   - Functor: Map operations over successful values
//   - Bifunctor: Map operations over both error and success values
//   - Apply: Apply wrapped functions to wrapped values
//   - Applicative: Lift pure values into the context
//   - Chain: Sequence dependent computations
//   - Monad: Full monadic operations
//   - Alt: Alternative computations for error recovery
//
// # Core Concepts
//
// SeqEither combines two powerful abstractions:
//   - Iteration: Processing sequences of values lazily
//   - Either: Representing computations that can fail with typed errors
//
// This allows for elegant error handling in iterator pipelines, where each
// element can independently succeed or fail, and operations can be chained
// while preserving error information.
//
// Common Use Cases
//
//   - Processing streams of data where individual items may fail
//   - Parsing sequences where each element requires validation
//   - Transforming collections with operations that can error
//   - Building pipelines with graceful error propagation
package itereither

//go:generate go run .. ioeither --count 10 --filename gen.go
