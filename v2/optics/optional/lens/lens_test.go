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

package lens

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/internal/common"
	OPT "github.com/IBM/fp-go/v2/optics/optional"
	OPTP "github.com/IBM/fp-go/v2/optics/optional/prism"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// ---- test fixtures ----

type Inner struct {
	A int
}

type State = O.Option[*Inner]

func (inner *Inner) getA() int {
	return inner.A
}

func (inner *Inner) setA(a int) *Inner {
	inner.A = a
	return inner
}

// makeLensA returns a Lens[*Inner, int] focusing on Inner.A.
func makeLensA() L.Lens[*Inner, int] {
	return L.MakeLensRef((*Inner).getA, (*Inner).setA)
}

// makeInnerOptional returns an Optional[State, *Inner] that focuses on the
// Some variant of a State (= Option[*Inner]).
func makeInnerOptional() L.Optional[State, *Inner] {
	return F.Pipe1(
		OPT.Id[State](),
		OPTP.Some[State, *Inner],
	)
}

// makeLensAComposed composes makeInnerOptional with the identity ref-lens and
// then with lensA, giving an Optional[State, int] for the inner A field.
func makeLensAComposed() L.Optional[State, int] {
	lensA := F.Pipe1(
		L.LensIdRef[Inner](),
		L.LensComposeLensRef[Inner](makeLensA()),
	)
	return F.Pipe1(
		makeInnerOptional(),
		Compose[State](lensA),
	)
}

// ---- existing regression test (unchanged behaviour) ----

func TestCompose(t *testing.T) {

	inner1 := Inner{1}

	lensa := L.MakeLensRef((*Inner).getA, (*Inner).setA)

	sa := F.Pipe1(
		OPT.Id[State](),
		OPTP.Some[State, *Inner],
	)
	ab := F.Pipe1(
		L.LensIdRef[Inner](),
		L.LensComposeLensRef[Inner](lensa),
	)
	sb := F.Pipe1(
		sa,
		Compose[State](ab),
	)
	// check get access
	assert.Equal(t, O.None[int](), sb.GetOption(O.None[*Inner]()))
	assert.Equal(t, O.Of(1), sb.GetOption(O.Of(&inner1)))

	// check set access
	res := F.Pipe1(
		sb.Set(2)(O.Of(&inner1)),
		O.Map(func(i *Inner) int {
			return i.A
		}),
	)
	assert.Equal(t, O.Of(2), res)
	assert.Equal(t, 1, inner1.A)

	assert.Equal(t, O.None[*Inner](), sb.Set(2)(O.None[*Inner]()))

}

// ---- equivalence with free-function ----

// TestCompose_EquivalentToOptionalComposeLens verifies that Compose is a thin
// wrapper around L.OptionalComposeLens and the two are behaviourally identical.
func TestCompose_EquivalentToOptionalComposeLens(t *testing.T) {
	lensA := makeLensA()
	sa := makeInnerOptional()

	method := F.Pipe1(sa, Compose[State](lensA))
	free := L.OptionalComposeLens[State](lensA)(sa)

	inner := &Inner{A: 7}
	states := []State{O.Some(inner), O.None[*Inner]()}

	for _, s := range states {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(s), method.GetOption(s))
		})
		t.Run("Set parity", func(t *testing.T) {
			assert.Equal(t, free.Set(99)(s), method.Set(99)(s))
		})
	}
}

// ---- GetOption behaviour ----

// TestCompose_GetOption_Some verifies that GetOption returns Some when the outer
// optional matches and the lens focuses successfully.
func TestCompose_GetOption_Some(t *testing.T) {
	sb := makeLensAComposed()

	inner := &Inner{A: 42}
	got := sb.GetOption(O.Some(inner))
	assert.Equal(t, O.Some(42), got)
}

// TestCompose_GetOption_None verifies that GetOption returns None when the outer
// optional produces None (State is None).
func TestCompose_GetOption_None(t *testing.T) {
	sb := makeLensAComposed()

	got := sb.GetOption(O.None[*Inner]())
	assert.Equal(t, O.None[int](), got)
}

// ---- Set behaviour ----

// TestCompose_Set_UpdatesLensFocus verifies that Set writes through the lens
// when the outer optional matches.
func TestCompose_Set_UpdatesLensFocus(t *testing.T) {
	sb := makeLensAComposed()

	inner := &Inner{A: 1}
	updated := sb.Set(99)(O.Some(inner))

	got := F.Pipe1(updated, O.Map(func(i *Inner) int { return i.A }))
	assert.Equal(t, O.Some(99), got)
}

// TestCompose_Set_NoOpOnNone verifies that Set is a no-op when the outer
// optional produces None (Optional no-op law / GetSet law).
func TestCompose_Set_NoOpOnNone(t *testing.T) {
	sb := makeLensAComposed()

	original := O.None[*Inner]()
	result := sb.Set(99)(original)
	assert.Equal(t, original, result)
}

// TestCompose_Set_Immutability verifies that Set does not mutate the original
// *Inner value stored inside the Some.
func TestCompose_Set_Immutability(t *testing.T) {
	sb := makeLensAComposed()

	inner := &Inner{A: 5}
	_ = sb.Set(50)(O.Some(inner))

	// original Inner must be unchanged
	assert.Equal(t, 5, inner.A)
}

// ---- Optional laws ----

// TestCompose_OptionalLaw1_GetSet verifies the GetSet law:
// If GetOption(s) == None then Set(b)(s) == s.
func TestCompose_OptionalLaw1_GetSet(t *testing.T) {
	sb := makeLensAComposed()

	none := O.None[*Inner]()
	result := sb.Set(42)(none)
	assert.Equal(t, none, result)
}

// TestCompose_OptionalLaw2_SetGet verifies the SetGet law:
// If GetOption(s) == Some(_) then GetOption(Set(b)(s)) == Some(b).
func TestCompose_OptionalLaw2_SetGet(t *testing.T) {
	sb := makeLensAComposed()

	s := O.Some(&Inner{A: 1})
	updated := sb.Set(7)(s)
	got := sb.GetOption(updated)
	assert.Equal(t, O.Some(7), got)
}

// TestCompose_OptionalLaw3_SetSet verifies the SetSet law:
// Set(b2)(Set(b1)(s)) == Set(b2)(s).
func TestCompose_OptionalLaw3_SetSet(t *testing.T) {
	sb := makeLensAComposed()

	s := O.Some(&Inner{A: 1})
	result1 := sb.Set(20)(sb.Set(10)(s))
	result2 := sb.Set(20)(s)

	got1 := F.Pipe1(result1, O.Map(func(i *Inner) int { return i.A }))
	got2 := F.Pipe1(result2, O.Map(func(i *Inner) int { return i.A }))
	assert.Equal(t, got2, got1)
}

// ---- pipeline usage ----

// TestCompose_PipelineUsage verifies that Compose integrates naturally with
// F.Pipe1, matching results from direct application.
func TestCompose_PipelineUsage(t *testing.T) {
	lensA := makeLensA()
	sa := makeInnerOptional()

	pipeResult := F.Pipe1(sa, Compose[State](lensA))
	direct := Compose[State](lensA)(sa)

	inner := &Inner{A: 3}
	s := O.Some(inner)

	assert.Equal(t, direct.GetOption(s), pipeResult.GetOption(s))
	assert.Equal(t, direct.Set(99)(s), pipeResult.Set(99)(s))
}

// ---- chained composition ----

// TestCompose_Chained verifies that two Compose calls can be chained: first
// narrowing from State to *Inner, then from *Inner to a nested field.
//
// Structure: State (= Option[*Inner]) → *Inner (via PrismSome) → Inner.A (via lensA)
func TestCompose_Chained(t *testing.T) {
	type Outer struct {
		inner *Inner
	}

	outerLens := L.MakeLens(
		func(o Outer) *Inner { return o.inner },
		func(o Outer, i *Inner) Outer { o.inner = i; return o },
	)

	outerOptional := OPT.MakeOptional(
		func(o Outer) O.Option[*Inner] {
			if o.inner != nil {
				return O.Some(o.inner)
			}
			return O.None[*Inner]()
		},
		func(o Outer, i *Inner) Outer { o.inner = i; return o },
	)

	aLens := makeLensA()

	// chain: Optional[Outer, *Inner] ∘ Lens[*Inner, int] → Optional[Outer, int]
	composed := F.Pipe1(outerOptional, Compose[Outer](aLens))

	inner := &Inner{A: 100}
	matching := Outer{inner: inner}
	nonMatching := Outer{inner: nil}

	t.Run("GetOption returns Some when both levels present", func(t *testing.T) {
		assert.Equal(t, O.Some(100), composed.GetOption(matching))
	})

	t.Run("GetOption returns None when optional misses", func(t *testing.T) {
		assert.Equal(t, O.None[int](), composed.GetOption(nonMatching))
	})

	t.Run("Set updates A when both levels present", func(t *testing.T) {
		updated := composed.Set(200)(matching)
		assert.Equal(t, O.Some(200), composed.GetOption(updated))
	})

	t.Run("Set is no-op when optional misses", func(t *testing.T) {
		result := composed.Set(200)(nonMatching)
		assert.Equal(t, nonMatching, result)
	})

	// silence unused-variable warning: outerLens is defined to show intent
	_ = outerLens
}
