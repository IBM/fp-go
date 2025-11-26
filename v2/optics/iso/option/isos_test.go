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

package option

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/optics/iso"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestFromZeroInt tests the FromZero isomorphism with integer type
func TestFromZeroInt(t *testing.T) {
	isoInt := FromZero[int]()

	t.Run("Get converts zero to None", func(t *testing.T) {
		result := isoInt.Get(0)
		assert.True(t, O.IsNone(result))
	})

	t.Run("Get converts non-zero to Some", func(t *testing.T) {
		result := isoInt.Get(42)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 42, O.MonadGetOrElse(result, func() int { return 0 }))
	})

	t.Run("Get converts negative to Some", func(t *testing.T) {
		result := isoInt.Get(-5)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, -5, O.MonadGetOrElse(result, func() int { return 0 }))
	})

	t.Run("ReverseGet converts None to zero", func(t *testing.T) {
		result := isoInt.ReverseGet(O.None[int]())
		assert.Equal(t, 0, result)
	})

	t.Run("ReverseGet converts Some to value", func(t *testing.T) {
		result := isoInt.ReverseGet(O.Some(42))
		assert.Equal(t, 42, result)
	})
}

// TestFromZeroString tests the FromZero isomorphism with string type
func TestFromZeroString(t *testing.T) {
	isoStr := FromZero[string]()

	t.Run("Get converts empty string to None", func(t *testing.T) {
		result := isoStr.Get("")
		assert.True(t, O.IsNone(result))
	})

	t.Run("Get converts non-empty string to Some", func(t *testing.T) {
		result := isoStr.Get("hello")
		assert.True(t, O.IsSome(result))
		assert.Equal(t, "hello", O.MonadGetOrElse(result, func() string { return "" }))
	})

	t.Run("ReverseGet converts None to empty string", func(t *testing.T) {
		result := isoStr.ReverseGet(O.None[string]())
		assert.Equal(t, "", result)
	})

	t.Run("ReverseGet converts Some to value", func(t *testing.T) {
		result := isoStr.ReverseGet(O.Some("world"))
		assert.Equal(t, "world", result)
	})
}

// TestFromZeroFloat tests the FromZero isomorphism with float64 type
func TestFromZeroFloat(t *testing.T) {
	isoFloat := FromZero[float64]()

	t.Run("Get converts 0.0 to None", func(t *testing.T) {
		result := isoFloat.Get(0.0)
		assert.True(t, O.IsNone(result))
	})

	t.Run("Get converts non-zero float to Some", func(t *testing.T) {
		result := isoFloat.Get(3.14)
		assert.True(t, O.IsSome(result))
		assert.InDelta(t, 3.14, O.MonadGetOrElse(result, func() float64 { return 0.0 }), 0.001)
	})

	t.Run("ReverseGet converts None to 0.0", func(t *testing.T) {
		result := isoFloat.ReverseGet(O.None[float64]())
		assert.Equal(t, 0.0, result)
	})

	t.Run("ReverseGet converts Some to value", func(t *testing.T) {
		result := isoFloat.ReverseGet(O.Some(2.718))
		assert.InDelta(t, 2.718, result, 0.001)
	})
}

// TestFromZeroPointer tests the FromZero isomorphism with pointer type
func TestFromZeroPointer(t *testing.T) {
	isoPtr := FromZero[*int]()

	t.Run("Get converts nil to None", func(t *testing.T) {
		result := isoPtr.Get(nil)
		assert.True(t, O.IsNone(result))
	})

	t.Run("Get converts non-nil pointer to Some", func(t *testing.T) {
		num := 42
		result := isoPtr.Get(&num)
		assert.True(t, O.IsSome(result))
		ptr := O.MonadGetOrElse(result, func() *int { return nil })
		assert.NotNil(t, ptr)
		assert.Equal(t, 42, *ptr)
	})

	t.Run("ReverseGet converts None to nil", func(t *testing.T) {
		result := isoPtr.ReverseGet(O.None[*int]())
		assert.Nil(t, result)
	})

	t.Run("ReverseGet converts Some to pointer", func(t *testing.T) {
		num := 99
		result := isoPtr.ReverseGet(O.Some(&num))
		assert.NotNil(t, result)
		assert.Equal(t, 99, *result)
	})
}

// TestFromZeroBool tests the FromZero isomorphism with bool type
func TestFromZeroBool(t *testing.T) {
	isoBool := FromZero[bool]()

	t.Run("Get converts false to None", func(t *testing.T) {
		result := isoBool.Get(false)
		assert.True(t, O.IsNone(result))
	})

	t.Run("Get converts true to Some", func(t *testing.T) {
		result := isoBool.Get(true)
		assert.True(t, O.IsSome(result))
		assert.True(t, O.MonadGetOrElse(result, func() bool { return false }))
	})

	t.Run("ReverseGet converts None to false", func(t *testing.T) {
		result := isoBool.ReverseGet(O.None[bool]())
		assert.False(t, result)
	})

	t.Run("ReverseGet converts Some to true", func(t *testing.T) {
		result := isoBool.ReverseGet(O.Some(true))
		assert.True(t, result)
	})
}

// TestFromZeroRoundTripLaws verifies the isomorphism laws
func TestFromZeroRoundTripLaws(t *testing.T) {
	t.Run("Law 1: ReverseGet(Get(t)) == t for integers", func(t *testing.T) {
		isoInt := FromZero[int]()

		// Test with zero value
		assert.Equal(t, 0, isoInt.ReverseGet(isoInt.Get(0)))

		// Test with non-zero values
		assert.Equal(t, 42, isoInt.ReverseGet(isoInt.Get(42)))
		assert.Equal(t, -10, isoInt.ReverseGet(isoInt.Get(-10)))
	})

	t.Run("Law 1: ReverseGet(Get(t)) == t for strings", func(t *testing.T) {
		isoStr := FromZero[string]()

		// Test with zero value
		assert.Equal(t, "", isoStr.ReverseGet(isoStr.Get("")))

		// Test with non-zero values
		assert.Equal(t, "hello", isoStr.ReverseGet(isoStr.Get("hello")))
	})

	t.Run("Law 2: Get(ReverseGet(opt)) == opt for None", func(t *testing.T) {
		isoInt := FromZero[int]()

		none := O.None[int]()
		result := isoInt.Get(isoInt.ReverseGet(none))
		assert.Equal(t, none, result)
	})

	t.Run("Law 2: Get(ReverseGet(opt)) == opt for Some", func(t *testing.T) {
		isoInt := FromZero[int]()

		some := O.Some(42)
		result := isoInt.Get(isoInt.ReverseGet(some))
		assert.Equal(t, some, result)
	})
}

// TestFromZeroWithModify tests using FromZero with iso.Modify
func TestFromZeroWithModify(t *testing.T) {
	isoInt := FromZero[int]()

	t.Run("Modify applies transformation to non-zero value", func(t *testing.T) {
		double := func(opt O.Option[int]) O.Option[int] {
			return O.MonadMap(opt, N.Mul(2))
		}

		result := iso.Modify[int](double)(isoInt)(5)
		assert.Equal(t, 10, result)
	})

	t.Run("Modify preserves zero value", func(t *testing.T) {
		double := func(opt O.Option[int]) O.Option[int] {
			return O.MonadMap(opt, N.Mul(2))
		}

		result := iso.Modify[int](double)(isoInt)(0)
		assert.Equal(t, 0, result)
	})
}

// TestFromZeroWithCompose tests composing FromZero with other isomorphisms
func TestFromZeroWithCompose(t *testing.T) {
	isoInt := FromZero[int]()

	// Create an isomorphism that doubles/halves values
	doubleIso := iso.MakeIso(
		func(opt O.Option[int]) O.Option[int] {
			return O.MonadMap(opt, N.Mul(2))
		},
		func(opt O.Option[int]) O.Option[int] {
			return O.MonadMap(opt, N.Div(2))
		},
	)

	composed := F.Pipe1(isoInt, iso.Compose[int](doubleIso))

	t.Run("Composed isomorphism works with non-zero", func(t *testing.T) {
		result := composed.Get(5)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 10, O.MonadGetOrElse(result, func() int { return 0 }))
	})

	t.Run("Composed isomorphism works with zero", func(t *testing.T) {
		result := composed.Get(0)
		assert.True(t, O.IsNone(result))
	})

	t.Run("Composed isomorphism reverse works", func(t *testing.T) {
		result := composed.ReverseGet(O.Some(20))
		assert.Equal(t, 10, result)
	})
}

// TestFromZeroWithUnwrapWrap tests using Unwrap and Wrap helpers
func TestFromZeroWithUnwrapWrap(t *testing.T) {
	isoInt := FromZero[int]()

	t.Run("Unwrap extracts Option from value", func(t *testing.T) {
		result := iso.Unwrap[O.Option[int]](42)(isoInt)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 42, O.MonadGetOrElse(result, func() int { return 0 }))
	})

	t.Run("Wrap creates value from Option", func(t *testing.T) {
		result := iso.Wrap[int](O.Some(99))(isoInt)
		assert.Equal(t, 99, result)
	})

	t.Run("To is alias for Unwrap", func(t *testing.T) {
		result := iso.To[O.Option[int]](42)(isoInt)
		assert.True(t, O.IsSome(result))
	})

	t.Run("From is alias for Wrap", func(t *testing.T) {
		result := iso.From[int](O.Some(99))(isoInt)
		assert.Equal(t, 99, result)
	})
}

// TestFromZeroWithReverse tests reversing the isomorphism
func TestFromZeroWithReverse(t *testing.T) {
	isoInt := FromZero[int]()
	reversed := iso.Reverse(isoInt)

	t.Run("Reversed Get is original ReverseGet", func(t *testing.T) {
		result := reversed.Get(O.Some(42))
		assert.Equal(t, 42, result)
	})

	t.Run("Reversed ReverseGet is original Get", func(t *testing.T) {
		result := reversed.ReverseGet(42)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 42, O.MonadGetOrElse(result, func() int { return 0 }))
	})
}

// TestFromZeroCustomType tests FromZero with a custom comparable type
func TestFromZeroCustomType(t *testing.T) {
	type UserID int

	isoUserID := FromZero[UserID]()

	t.Run("Get converts zero UserID to None", func(t *testing.T) {
		result := isoUserID.Get(UserID(0))
		assert.True(t, O.IsNone(result))
	})

	t.Run("Get converts non-zero UserID to Some", func(t *testing.T) {
		result := isoUserID.Get(UserID(123))
		assert.True(t, O.IsSome(result))
		assert.Equal(t, UserID(123), O.MonadGetOrElse(result, func() UserID { return 0 }))
	})

	t.Run("ReverseGet converts None to zero UserID", func(t *testing.T) {
		result := isoUserID.ReverseGet(O.None[UserID]())
		assert.Equal(t, UserID(0), result)
	})

	t.Run("ReverseGet converts Some to UserID", func(t *testing.T) {
		result := isoUserID.ReverseGet(O.Some(UserID(456)))
		assert.Equal(t, UserID(456), result)
	})
}

// TestFromZeroEdgeCases tests edge cases and boundary conditions
func TestFromZeroEdgeCases(t *testing.T) {
	t.Run("Works with maximum int value", func(t *testing.T) {
		isoInt := FromZero[int]()
		maxInt := int(^uint(0) >> 1)

		result := isoInt.Get(maxInt)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, maxInt, isoInt.ReverseGet(result))
	})

	t.Run("Works with minimum int value", func(t *testing.T) {
		isoInt := FromZero[int]()
		minInt := -int(^uint(0)>>1) - 1

		result := isoInt.Get(minInt)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, minInt, isoInt.ReverseGet(result))
	})

	t.Run("Works with very long strings", func(t *testing.T) {
		isoStr := FromZero[string]()
		longStr := string(make([]byte, 10000))
		for i := range longStr {
			longStr = longStr[:i] + "a" + longStr[i+1:]
		}

		result := isoStr.Get(longStr)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, longStr, isoStr.ReverseGet(result))
	})
}
