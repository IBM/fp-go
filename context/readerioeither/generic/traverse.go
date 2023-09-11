// Copyright (c) 2023 IBM Corp.
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

package generic

import (
	"context"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	RA "github.com/IBM/fp-go/internal/array"
	RR "github.com/IBM/fp-go/internal/record"
)

// MonadTraverseArray transforms an array
func MonadTraverseArray[
	AS ~[]A,
	GRBS ~func(context.Context) GIOBS,
	GRB ~func(context.Context) GIOB,
	GIOBS ~func() E.Either[error, BS],
	GIOB ~func() E.Either[error, B],
	BS ~[]B,
	A, B any](as AS, f func(A) GRB) GRBS {

	return RA.MonadTraverse[AS](
		Of[GRBS, GIOBS, BS],
		Map[GRBS, func(context.Context) func() E.Either[error, func(B) BS], GIOBS, func() E.Either[error, func(B) BS], BS, func(B) BS],
		Ap[GRBS, func(context.Context) func() E.Either[error, func(B) BS], GRB],

		as, f,
	)
}

// TraverseArray transforms an array
func TraverseArray[
	AS ~[]A,
	GRBS ~func(context.Context) GIOBS,
	GRB ~func(context.Context) GIOB,
	GIOBS ~func() E.Either[error, BS],
	GIOB ~func() E.Either[error, B],
	BS ~[]B,
	A, B any](f func(A) GRB) func(AS) GRBS {

	return RA.Traverse[AS](
		Of[GRBS, GIOBS, BS],
		Map[GRBS, func(context.Context) func() E.Either[error, func(B) BS], GIOBS, func() E.Either[error, func(B) BS], BS, func(B) BS],
		Ap[GRBS, func(context.Context) func() E.Either[error, func(B) BS], GRB],

		f,
	)
}

// TraverseArrayWithIndex transforms an array
func TraverseArrayWithIndex[
	AS ~[]A,
	GRBS ~func(context.Context) GIOBS,
	GRB ~func(context.Context) GIOB,
	GIOBS ~func() E.Either[error, BS],
	GIOB ~func() E.Either[error, B],
	BS ~[]B,
	A, B any](f func(int, A) GRB) func(AS) GRBS {

	return RA.TraverseWithIndex[AS](
		Of[GRBS, GIOBS, BS],
		Map[GRBS, func(context.Context) func() E.Either[error, func(B) BS], GIOBS, func() E.Either[error, func(B) BS], BS, func(B) BS],
		Ap[GRBS, func(context.Context) func() E.Either[error, func(B) BS], GRB],

		f,
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[
	AS ~[]A,
	GAS ~[]GRA,
	GRAS ~func(context.Context) GIOAS,
	GRA ~func(context.Context) GIOA,
	GIOAS ~func() E.Either[error, AS],
	GIOA ~func() E.Either[error, A],
	A any](ma GAS) GRAS {

	return MonadTraverseArray[GAS, GRAS](ma, F.Identity[GRA])
}

// MonadTraverseRecord transforms a record
func MonadTraverseRecord[K comparable,
	AS ~map[K]A,
	GRBS ~func(context.Context) GIOBS,
	GRB ~func(context.Context) GIOB,
	GIOBS ~func() E.Either[error, BS],
	GIOB ~func() E.Either[error, B],
	BS ~map[K]B,

	A, B any](ma AS, f func(A) GRB) GRBS {

	return RR.MonadTraverse[AS](
		Of[GRBS, GIOBS, BS],
		Map[GRBS, func(context.Context) func() E.Either[error, func(B) BS], GIOBS, func() E.Either[error, func(B) BS], BS, func(B) BS],
		Ap[GRBS, func(context.Context) func() E.Either[error, func(B) BS], GRB],

		ma, f,
	)
}

// TraverseRecord transforms a record
func TraverseRecord[K comparable,
	AS ~map[K]A,
	GRBS ~func(context.Context) GIOBS,
	GRB ~func(context.Context) GIOB,
	GIOBS ~func() E.Either[error, BS],
	GIOB ~func() E.Either[error, B],
	BS ~map[K]B,

	A, B any](f func(A) GRB) func(AS) GRBS {

	return RR.Traverse[AS](
		Of[GRBS, GIOBS, BS],
		Map[GRBS, func(context.Context) func() E.Either[error, func(B) BS], GIOBS, func() E.Either[error, func(B) BS], BS, func(B) BS],
		Ap[GRBS, func(context.Context) func() E.Either[error, func(B) BS], GRB],

		f,
	)
}

// TraverseRecordWithIndex transforms a record
func TraverseRecordWithIndex[K comparable,
	AS ~map[K]A,
	GRBS ~func(context.Context) GIOBS,
	GRB ~func(context.Context) GIOB,
	GIOBS ~func() E.Either[error, BS],
	GIOB ~func() E.Either[error, B],
	BS ~map[K]B,

	A, B any](f func(K, A) GRB) func(AS) GRBS {

	return RR.TraverseWithIndex[AS](
		Of[GRBS, GIOBS, BS],
		Map[GRBS, func(context.Context) func() E.Either[error, func(B) BS], GIOBS, func() E.Either[error, func(B) BS], BS, func(B) BS],
		Ap[GRBS, func(context.Context) func() E.Either[error, func(B) BS], GRB],

		f,
	)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable,
	AS ~map[K]A,
	GAS ~map[K]GRA,
	GRAS ~func(context.Context) GIOAS,
	GRA ~func(context.Context) GIOA,
	GIOAS ~func() E.Either[error, AS],
	GIOA ~func() E.Either[error, A],
	A any](ma GAS) GRAS {

	return MonadTraverseRecord[K, GAS, GRAS](ma, F.Identity[GRA])
}
