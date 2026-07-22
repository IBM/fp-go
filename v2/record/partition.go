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

package record

import (
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/record/generic"
)

// MonadPartition splits a record into two records based on a predicate applied
// to each value. Entries for which the predicate returns true are placed in the
// right (tail) record; entries for which it returns false are placed in the left
// (head) record.
//
// Type Parameters:
//   - K: key type, must be comparable
//   - V: value type
//
// Parameters:
//   - kvs: the source record to partition
//   - pred: a function that tests each value
//
// Returns:
//   - Pair[Record[K, V], Record[K, V]]: a pair where Head holds the non-matching
//     entries and Tail holds the matching entries
//
// See Also:
//   - MonadPartitionWithIndex: index-aware variant
//   - Partition: curried variant
func MonadPartition[K comparable, V any](kvs Record[K, V], pred func(V) bool) pair.Pair[Record[K, V], Record[K, V]] {
	return generic.MonadPartition(kvs, pred)
}

// MonadPartitionWithIndex splits a record into two records based on a predicate
// applied to each key-value pair. Entries for which the predicate returns true
// are placed in the right (tail) record; entries for which it returns false are
// placed in the left (head) record.
//
// Type Parameters:
//   - K: key type, must be comparable
//   - V: value type
//
// Parameters:
//   - kvs: the source record to partition
//   - pred: a function that tests each key-value pair
//
// Returns:
//   - Pair[Record[K, V], Record[K, V]]: a pair where Head holds the non-matching
//     entries and Tail holds the matching entries
//
// See Also:
//   - MonadPartition: value-only variant
//   - PartitionWithIndex: curried variant
func MonadPartitionWithIndex[K comparable, V any](kvs Record[K, V], pred func(K, V) bool) pair.Pair[Record[K, V], Record[K, V]] {
	return generic.MonadPartitionWithIndex(kvs, pred)
}

// Partition returns a curried function that splits a record into two records
// based on a predicate applied to each value. The returned function accepts a
// record and returns a pair where Head holds the non-matching entries and Tail
// holds the matching entries.
//
// Type Parameters:
//   - K: key type, must be comparable
//   - V: value type
//
// Parameters:
//   - pred: a function that tests each value
//
// Returns:
//   - pair.Kleisli[Record[K, V], Record[K, V], Record[K, V]]: a function from
//     Record[K, V] to Pair[Record[K, V], Record[K, V]]
//
// See Also:
//   - PartitionWithIndex: index-aware variant
//   - MonadPartition: uncurried variant
func Partition[K comparable, V any](pred func(V) bool) pair.Kleisli[Record[K, V], Record[K, V], Record[K, V]] {
	return generic.Partition[Record[K, V]](pred)
}

// PartitionWithIndex returns a curried function that splits a record into two
// records based on a predicate applied to each key-value pair. The returned
// function accepts a record and returns a pair where Head holds the non-matching
// entries and Tail holds the matching entries.
//
// Type Parameters:
//   - K: key type, must be comparable
//   - V: value type
//
// Parameters:
//   - pred: a function that tests each key-value pair
//
// Returns:
//   - pair.Kleisli[Record[K, V], Record[K, V], Record[K, V]]: a function from
//     Record[K, V] to Pair[Record[K, V], Record[K, V]]
//
// See Also:
//   - Partition: value-only variant
//   - MonadPartitionWithIndex: uncurried variant
func PartitionWithIndex[K comparable, V any](pred func(K, V) bool) pair.Kleisli[Record[K, V], Record[K, V], Record[K, V]] {
	return generic.PartitionWithIndex[Record[K, V]](pred)
}
