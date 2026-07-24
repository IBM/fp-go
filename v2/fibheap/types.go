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

package fibheap

import (
	"github.com/IBM/fp-go/v2/option"
)

type (
	// Option is an alias for option.Option, representing a value that may or may not be present.
	Option[A any] = option.Option[A]

	// Handle is an opaque reference to a node inside a Heap.
	// It is returned by Insert and must be used with DecreaseKey and Delete.
	// A Handle becomes invalid after the node has been deleted from the heap.
	Handle[A any] struct {
		n *node[A]
	}

	// Heap is a Fibonacci heap parameterised over element type A.
	// It is a value type; copy it before mutating if you need to preserve the original.
	// The zero value is a valid empty heap.
	Heap[A any] struct {
		min  *node[A]
		size int
	}
)
