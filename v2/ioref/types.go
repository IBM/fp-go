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

// Package ioref provides mutable references in the IO monad.
//
// IORef represents a mutable reference that can be read and written within IO computations.
// It provides thread-safe access to shared mutable state using read-write locks.
//
// This is inspired by Haskell's Data.IORef module and provides a functional approach
// to managing mutable state with explicit IO effects.
//
// Example usage:
//
//	// Create a new IORef
//	ref := ioref.MakeIORef(42)()
//
//	// Read the current value
//	value := ioref.Read(ref)()  // 42
//
//	// Write a new value
//	ioref.Write(100)(ref)()
//
//	// Modify the value
//	ioref.Modify(func(x int) int { return x * 2 })(ref)()
//
//	// Read the modified value
//	newValue := ioref.Read(ref)()  // 200
package ioref

import (
	"sync"

	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/pair"
)

type (
	// ioRef is the internal implementation of a mutable reference.
	// It uses a read-write mutex to ensure thread-safe access.
	ioRef[A any] struct {
		mu sync.RWMutex
		a  A
	}

	// IO represents a synchronous computation that may have side effects.
	// It's a function that takes no arguments and returns a value of type A.
	IO[A any] = io.IO[A]

	// IORef represents a mutable reference to a value of type A.
	// Operations on IORef are thread-safe and performed within the IO monad.
	//
	// IORef provides a way to work with mutable state in a functional style,
	// where mutations are explicit and contained within IO computations.
	IORef[A any] = *ioRef[A]

	// Endomorphism represents a function from A to A.
	// It's commonly used with Modify to transform the value in an IORef.
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	Pair[A, B any] = pair.Pair[A, B]
)
