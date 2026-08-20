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
	"math"

	O "github.com/IBM/fp-go/v2/ord"
)

// node is an internal doubly-linked circular list node for the Fibonacci heap.
// siblings form a circular doubly-linked list; child points to one arbitrary child.
type node[A any] struct {
	key    A
	degree int
	marked bool
	parent *node[A]
	child  *node[A]
	prev   *node[A]
	next   *node[A]
}

// newNode allocates a new singleton node.
func newNode[A any](key A) *node[A] {
	n := &node[A]{key: key}
	n.prev = n
	n.next = n
	return n
}

// addToRootList inserts n into the root list of the heap (before the min pointer).
// n must already be a singleton (prev == next == n) or freshly unlinked.
func addToRootList[A any](h *Heap[A], n *node[A]) {
	n.parent = nil
	if h.min == nil {
		n.prev = n
		n.next = n
		h.min = n
		return
	}
	// splice n in between h.min.prev and h.min
	n.next = h.min
	n.prev = h.min.prev
	h.min.prev.next = n
	h.min.prev = n
}

// removeFromList unlinks n from its sibling ring without touching its children.
func removeFromList[A any](n *node[A]) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = n
	n.next = n
}

// link makes y a child of x (y.key >= x.key in min-heap sense).
func link[A any](x, y *node[A]) {
	removeFromList(y)
	y.parent = x
	if x.child == nil {
		x.child = y
		y.prev = y
		y.next = y
	} else {
		// insert y into x's child list
		y.next = x.child
		y.prev = x.child.prev
		x.child.prev.next = y
		x.child.prev = y
	}
	x.degree++
	y.marked = false
}

// consolidate restructures the root list so that no two roots have the same degree.
func consolidate[A any](h *Heap[A], o O.Ord[A]) {
	// upper bound on degree: log_phi(n) ≈ 1.44 * log2(n)
	maxDeg := int(math.Log2(float64(h.size))*1.5) + 2
	if maxDeg < 1 {
		maxDeg = 1
	}
	table := make([]*node[A], maxDeg+1)

	// collect all roots
	roots := make([]*node[A], 0, maxDeg+2)
	cur := h.min
	if cur != nil {
		start := cur
		for {
			roots = append(roots, cur)
			cur = cur.next
			if cur == start {
				break
			}
		}
	}

	for _, w := range roots {
		x := w
		d := x.degree
		for d < len(table) && table[d] != nil {
			y := table[d]
			if o.Compare(x.key, y.key) > 0 {
				x, y = y, x
			}
			link(x, y)
			table[d] = nil
			d++
		}
		// grow table if needed
		for d >= len(table) {
			table = append(table, nil)
		}
		table[d] = x
	}

	// rebuild root list and find new min
	h.min = nil
	for _, r := range table {
		if r == nil {
			continue
		}
		r.prev = r
		r.next = r
		r.parent = nil
		if h.min == nil {
			h.min = r
		} else {
			// insert r into root list
			r.next = h.min
			r.prev = h.min.prev
			h.min.prev.next = r
			h.min.prev = r
			if o.Compare(r.key, h.min.key) < 0 {
				h.min = r
			}
		}
	}
}

// cut removes n from its parent's child list and adds it to the root list.
func cut[A any](h *Heap[A], n, parent *node[A]) {
	// remove n from parent's children
	if n.next == n {
		parent.child = nil
	} else {
		if parent.child == n {
			parent.child = n.next
		}
		removeFromList(n)
	}
	parent.degree--
	n.marked = false
	addToRootList(h, n)
}

// cascadingCut propagates cuts up the tree.
func cascadingCut[A any](h *Heap[A], n *node[A]) {
	p := n.parent
	if p == nil {
		return
	}
	if !n.marked {
		n.marked = true
	} else {
		cut(h, n, p)
		cascadingCut(h, p)
	}
}
