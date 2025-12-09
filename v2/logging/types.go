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

package logging

import (
	"github.com/IBM/fp-go/v2/endomorphism"
)

type (
	// Endomorphism represents a function that takes a value of type A and returns a value of the same type A.
	// This is a type alias for endomorphism.Endomorphism[A], which is a fundamental concept in functional
	// programming representing transformations that preserve type.
	//
	// In the context of the logging package, this is primarily used for context transformations,
	// such as adding a logger to a context while maintaining the context.Context type.
	//
	// Type Parameters:
	//   - A: The type that the endomorphism operates on
	//
	// Example:
	//
	//	// An endomorphism that adds a logger to a context
	//	var addLogger Endomorphism[context.Context] = WithLogger(myLogger)
	//
	//	// Apply the transformation
	//	ctx := context.Background()
	//	newCtx := addLogger(ctx) // Both ctx and newCtx are context.Context
	Endomorphism[A any] = endomorphism.Endomorphism[A]
)
