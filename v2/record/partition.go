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
	"github.com/IBM/fp-go/v2/either"
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

// MonadPartitionMap splits a record into two records by applying a function
// that returns an Either to each value. Entries for which the function returns
// a Right value are collected in the right (tail) record; entries for which it
// returns a Left value are collected in the left (head) record.
//
// This is the uncurried, monadic form of PartitionMap.
//
// Type Parameters:
//   - K: key type, must be comparable
//   - A: input value type
//   - L: type of Left (head) values produced by pred
//   - R: type of Right (tail) values produced by pred
//
// Parameters:
//   - kvs: the source record to partition
//   - pred: a Kleisli arrow from A to Either[L, R] that classifies each value
//
// Returns:
//   - Pair[Record[K, L], Record[K, R]]: a pair where Head holds keys mapped to
//     their Left results and Tail holds keys mapped to their Right results
//
// See Also:
//   - PartitionMap: curried variant
//   - MonadPartition: simpler boolean-predicate variant
func MonadPartitionMap[K comparable, A, L, R any](kvs Record[K, A], pred either.Kleisli[L, A, R]) pair.Pair[Record[K, L], Record[K, R]] {
	return generic.MonadPartitionMap[Record[K, A], Record[K, L], Record[K, R]](kvs, pred)
}

// PartitionMap returns a curried function that splits a record into two records
// by applying a function that returns an Either to each value. Entries for which
// the function returns a Right value are collected in the right (tail) record;
// entries for which it returns a Left value are collected in the left (head)
// record.
//
// This is the curried form of MonadPartitionMap.
//
// Type Parameters:
//   - K: key type, must be comparable
//   - A: input value type
//   - L: type of Left (head) values produced by pred
//   - R: type of Right (tail) values produced by pred
//
// Parameters:
//   - pred: a Kleisli arrow from A to Either[L, R] that classifies each value
//
// Returns:
//   - pair.Kleisli[Record[K, L], Record[K, A], Record[K, R]]: a reusable
//     function from Record[K, A] to Pair[Record[K, L], Record[K, R]]
//
// See Also:
//   - MonadPartitionMap: uncurried variant
//   - Partition: simpler boolean-predicate variant
func PartitionMap[K comparable, A, L, R any](pred either.Kleisli[L, A, R]) pair.Kleisli[Record[K, L], Record[K, A], Record[K, R]] {
	return generic.PartitionMap[Record[K, A], Record[K, L], Record[K, R]](pred)
}
