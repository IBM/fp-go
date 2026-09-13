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

// Package fibheap provides a generic Fibonacci heap data structure.
//
// A Fibonacci heap is a priority-queue data structure with excellent amortised
// complexity for the operations used by graph algorithms such as Dijkstra and Prim:
//
//	Operation     Amortised complexity
//	──────────────────────────────────
//	Empty         O(1)
//	Insert        O(1)
//	FindMin       O(1)
//	Merge         O(1)
//	DecreaseKey   O(1)
//	ExtractMin    O(log n)
//	Delete        O(log n)
//	IsEmpty       O(1)
//	Size          O(1)
//
// # Ordering
//
// All operations that compare keys accept an [ord.Ord][A] instance, so the heap
// can be used with any type that has a total order — including custom structs.
//
//	intOrd := ord.FromStrictCompare[int]()
//	h := fibheap.Empty[int]()
//	h, h1 := fibheap.Insert(intOrd)(3)(h)
//	h, h2 := fibheap.Insert(intOrd)(1)(h)
//	h, h3 := fibheap.Insert(intOrd)(2)(h)
//	fibheap.FindMin(h) // Some(1)
//
// # Handles
//
// Insert returns a [Handle][A] alongside the updated heap. The handle is an
// opaque reference to the inserted node and must be used with [DecreaseKey]
// and [Delete]:
//
//	h, handle := fibheap.Insert(intOrd)(10)(h)
//	h = fibheap.DecreaseKey(intOrd)(5)(handle)(h)  // key 10 → 5
//	h = fibheap.Delete(intOrd)(handle)(h)           // remove it entirely
//
// A Handle is invalidated once [Delete] is called on it; do not reuse it.
//
// # API style
//
// All multi-parameter functions follow fp-go's data-last, curried convention so
// they compose cleanly with [function.Flow] and [function.Pipe]:
//
//	insert := fibheap.Insert(ord.FromStrictCompare[int]())
//	h, _ = insert(1)(h)
//	h, _ = insert(2)(h)
//
// # Merging heaps
//
//	combined := fibheap.Merge(intOrd)(h1)(h2)
//
// After merging, neither h1 nor h2 should be used independently.
//
// # Related packages
//
//   - ord: Total ordering type class used to compare keys
//   - option: Option type returned by FindMin and ExtractMin
package fibheap
