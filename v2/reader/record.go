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

package reader

import (
	"github.com/IBM/fp-go/v2/function"
	RR "github.com/IBM/fp-go/v2/internal/record"
)

//go:inline
func MonadTraverseRecord[K comparable, R, A, B any](ma map[K]A, f Kleisli[R, A, B]) Reader[R, map[K]B] {
	return RR.MonadTraverse[map[K]A, map[K]B](
		Of,
		Map,
		Ap,
		ma,
		f,
	)
}

//go:inline
func TraverseRecord[K comparable, R, A, B any](f Kleisli[R, A, B]) func(map[K]A) Reader[R, map[K]B] {
	return RR.Traverse[map[K]A, map[K]B](
		Of,
		Map,
		Ap,
		f,
	)
}

//go:inline
func MonadTraverseRecordWithIndex[K comparable, R, A, B any](ma map[K]A, f func(K, A) Reader[R, B]) Reader[R, map[K]B] {
	return RR.MonadTraverseWithIndex[map[K]A, map[K]B](
		Of,
		Map,
		Ap,
		ma,
		f,
	)
}

//go:inline
func TraverseRecordWithIndex[K comparable, R, A, B any](f func(K, A) Reader[R, B]) func(map[K]A) Reader[R, map[K]B] {
	return RR.TraverseWithIndex[map[K]A, map[K]B](
		Of,
		Map,
		Ap,
		f,
	)
}

//go:inline
func SequenceRecord[K comparable, R, A any](ma map[K]Reader[R, A]) Reader[R, map[K]A] {
	return MonadTraverseRecord(ma, function.Identity[Reader[R, A]])
}
