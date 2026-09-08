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
	O "github.com/IBM/fp-go/v2/option"
	Ord "github.com/IBM/fp-go/v2/ord"
)

// Empty returns a new, empty Fibonacci heap.
//
// Complexity: O(1)
//
// Example:
//
//	h := fibheap.Empty[int]()
//
//go:inline
func Empty[A any]() Heap[A] {
	return Heap[A]{}
}

// IsEmpty returns true if the heap contains no elements.
//
// Complexity: O(1)
//
// Example:
//
//	h := fibheap.Empty[int]()
//	fibheap.IsEmpty(h) // true
//
//go:inline
func IsEmpty[A any](h Heap[A]) bool {
	return h.size == 0
}

// Size returns the number of elements in the heap.
//
// Complexity: O(1)
//
// Example:
//
//	h, _ := fibheap.Insert(ord.FromStrictCompare[int]())(1)(fibheap.Empty[int]())
//	fibheap.Size(h) // 1
//
//go:inline
func Size[A any](h Heap[A]) int {
	return h.size
}

// FindMin returns the minimum element wrapped in Some, or None if the heap is empty.
//
// Complexity: O(1)
//
// Example:
//
//	h, _ := fibheap.Insert(ord.FromStrictCompare[int]())(3)(fibheap.Empty[int]())
//	fibheap.FindMin(h) // Some(3)
//
//go:inline
func FindMin[A any](h Heap[A]) Option[A] {
	if h.min == nil {
		return O.None[A]()
	}
	return O.Some(h.min.key)
}

// Insert inserts a value into the heap and returns the updated heap together with
// a Handle that can later be used with DecreaseKey or Delete.
//
// Complexity: O(1) amortised
//
// The function is curried (data-last):
//
//	insert := fibheap.Insert(ord.FromStrictCompare[int]())
//	h, handle := insert(42)(h)
//
//go:inline
func Insert[A any](o Ord.Ord[A]) func(A) func(Heap[A]) (Heap[A], Handle[A]) {
	return func(key A) func(Heap[A]) (Heap[A], Handle[A]) {
		return func(h Heap[A]) (Heap[A], Handle[A]) {
			n := newNode(key)
			addToRootList(&h, n)
			if o.Compare(key, h.min.key) <= 0 {
				h.min = n
			}
			h.size++
			return h, Handle[A]{n: n}
		}
	}
}

// ExtractMin removes and returns the minimum element.
// Returns the updated heap and Some(min) if the heap was non-empty, or None otherwise.
//
// Complexity: O(log n) amortised
//
// The function is curried (data-last):
//
//	extract := fibheap.ExtractMin(ord.FromStrictCompare[int]())
//	h, minVal := extract(h)
//
//go:inline
func ExtractMin[A any](o Ord.Ord[A]) func(Heap[A]) (Heap[A], Option[A]) {
	return func(h Heap[A]) (Heap[A], Option[A]) {
		z := h.min
		if z == nil {
			return h, O.None[A]()
		}

		// add z's children to root list
		if z.child != nil {
			// collect children first to avoid pointer aliasing while iterating
			children := make([]*node[A], 0, z.degree)
			c := z.child
			start := c
			for {
				children = append(children, c)
				c = c.next
				if c == start {
					break
				}
			}
			for _, child := range children {
				addToRootList(&h, child)
				child.parent = nil
			}
		}

		// remove z from root list
		if z.next == z {
			// z was the only root
			h.min = nil
		} else {
			z.prev.next = z.next
			z.next.prev = z.prev
			h.min = z.next
			consolidate(&h, o)
		}
		h.size--

		return h, O.Some(z.key)
	}
}

// Merge combines two heaps into one. The original heaps must not be used afterwards.
//
// Complexity: O(1)
//
// The function is curried (data-last):
//
//	merge := fibheap.Merge(ord.FromStrictCompare[int]())
//	combined := merge(h1)(h2)
//
//go:inline
func Merge[A any](o Ord.Ord[A]) func(Heap[A]) func(Heap[A]) Heap[A] {
	return func(h1 Heap[A]) func(Heap[A]) Heap[A] {
		return func(h2 Heap[A]) Heap[A] {
			if h1.min == nil {
				return h2
			}
			if h2.min == nil {
				return h1
			}
			// concatenate two circular lists:
			// h1: ... <-> h1.min.prev <-> h1.min <-> ...
			// h2: ... <-> h2.min.prev <-> h2.min <-> ...
			last1 := h1.min.prev
			last2 := h2.min.prev

			last1.next = h2.min
			h2.min.prev = last1
			last2.next = h1.min
			h1.min.prev = last2

			if o.Compare(h2.min.key, h1.min.key) < 0 {
				h1.min = h2.min
			}
			h1.size += h2.size
			return h1
		}
	}
}

// DecreaseKey updates the key of the node referenced by handle to newKey.
// newKey must be less than or equal to the current key (panics otherwise).
//
// Complexity: O(1) amortised
//
// The function is curried (data-last):
//
//	decrease := fibheap.DecreaseKey(ord.FromStrictCompare[int]())
//	h = decrease(newKey)(handle)(h)
//
//go:inline
func DecreaseKey[A any](o Ord.Ord[A]) func(A) func(Handle[A]) func(Heap[A]) Heap[A] {
	return func(newKey A) func(Handle[A]) func(Heap[A]) Heap[A] {
		return func(handle Handle[A]) func(Heap[A]) Heap[A] {
			return func(h Heap[A]) Heap[A] {
				n := handle.n
				if o.Compare(newKey, n.key) > 0 {
					panic("fibheap.DecreaseKey: newKey is greater than current key")
				}
				n.key = newKey
				p := n.parent
				if p != nil && o.Compare(n.key, p.key) < 0 {
					cut(&h, n, p)
					cascadingCut(&h, p)
				}
				if o.Compare(n.key, h.min.key) < 0 {
					h.min = n
				}
				return h
			}
		}
	}
}

// Delete removes the node referenced by handle from the heap.
// The Handle must not be used after this call.
//
// Complexity: O(log n) amortised
//
// The function is curried (data-last):
//
//	delete := fibheap.Delete(ord.FromStrictCompare[int]())
//	h = delete(handle)(h)
//
//go:inline
func Delete[A any](o Ord.Ord[A]) func(Handle[A]) func(Heap[A]) Heap[A] {
	return func(handle Handle[A]) func(Heap[A]) Heap[A] {
		return func(h Heap[A]) Heap[A] {
			// promote handle to minimum by forcing it to be the min node
			handle.n.key = h.min.key // temporarily equal to min
			// cut it to root list if it has a parent
			if handle.n.parent != nil {
				p := handle.n.parent
				cut(&h, handle.n, p)
				cascadingCut(&h, p)
			}
			// force it to be the minimum
			h.min = handle.n
			h2, _ := ExtractMin[A](o)(h)
			return h2
		}
	}
}
