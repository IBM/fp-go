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

package cli

import (
	"testing"
	"time"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
	C "github.com/urfave/cli/v3"
)

func TestStringFlagPrism_Success(t *testing.T) {
	t.Run("extracts StringFlag from Flag", func(t *testing.T) {
		// Arrange
		prism := StringFlagPrism()
		var flag C.Flag = &C.StringFlag{Name: "input", Value: "test"}

		// Act
		result := prism.GetOption(flag)

		// Assert
		assert.True(t, O.IsSome(result))
		extracted := O.MonadFold(result, func() *C.StringFlag { return nil }, func(f *C.StringFlag) *C.StringFlag { return f })
		assert.NotNil(t, extracted)
		assert.Equal(t, "input", extracted.Name)
		assert.Equal(t, "test", extracted.Value)
	})
}

func TestStringFlagPrism_Failure(t *testing.T) {
	t.Run("returns None for non-StringFlag", func(t *testing.T) {
		// Arrange
		prism := StringFlagPrism()
		var flag C.Flag = &C.IntFlag{Name: "count"}

		// Act
		result := prism.GetOption(flag)

		// Assert
		assert.True(t, O.IsNone(result))
	})
}

func TestStringFlagPrism_ReverseGet(t *testing.T) {
	t.Run("converts StringFlag back to Flag", func(t *testing.T) {
		// Arrange
		prism := StringFlagPrism()
		strFlag := &C.StringFlag{Name: "output", Value: "result"}

		// Act
		flag := prism.ReverseGet(strFlag)

		// Assert
		assert.NotNil(t, flag)
		assert.IsType(t, &C.StringFlag{}, flag)
	})
}

func TestIntFlagPrism_Success(t *testing.T) {
	t.Run("extracts IntFlag from Flag", func(t *testing.T) {
		// Arrange
		prism := IntFlagPrism()
		var flag C.Flag = &C.IntFlag{Name: "count", Value: 42}

		// Act
		result := prism.GetOption(flag)

		// Assert
		assert.True(t, O.IsSome(result))
		extracted := O.MonadFold(result, func() *C.IntFlag { return nil }, func(f *C.IntFlag) *C.IntFlag { return f })
		assert.NotNil(t, extracted)
		assert.Equal(t, "count", extracted.Name)
		assert.Equal(t, 42, extracted.Value)
	})
}

func TestBoolFlagPrism_Success(t *testing.T) {
	t.Run("extracts BoolFlag from Flag", func(t *testing.T) {
		// Arrange
		prism := BoolFlagPrism()
		var flag C.Flag = &C.BoolFlag{Name: "verbose", Value: true}

		// Act
		result := prism.GetOption(flag)

		// Assert
		assert.True(t, O.IsSome(result))
		extracted := O.MonadFold(result, func() *C.BoolFlag { return nil }, func(f *C.BoolFlag) *C.BoolFlag { return f })
		assert.NotNil(t, extracted)
		assert.Equal(t, "verbose", extracted.Name)
		assert.Equal(t, true, extracted.Value)
	})
}

func TestFloat64FlagPrism_Success(t *testing.T) {
	t.Run("extracts Float64Flag from Flag", func(t *testing.T) {
		// Arrange
		prism := Float64FlagPrism()
		var flag C.Flag = &C.Float64Flag{Name: "ratio", Value: 0.5}

		// Act
		result := prism.GetOption(flag)

		// Assert
		assert.True(t, O.IsSome(result))
		extracted := O.MonadFold(result, func() *C.Float64Flag { return nil }, func(f *C.Float64Flag) *C.Float64Flag { return f })
		assert.NotNil(t, extracted)
		assert.Equal(t, "ratio", extracted.Name)
		assert.Equal(t, 0.5, extracted.Value)
	})
}

func TestDurationFlagPrism_Success(t *testing.T) {
	t.Run("extracts DurationFlag from Flag", func(t *testing.T) {
		// Arrange
		prism := DurationFlagPrism()
		duration := 30 * time.Second
		var flag C.Flag = &C.DurationFlag{Name: "timeout", Value: duration}

		// Act
		result := prism.GetOption(flag)

		// Assert
		assert.True(t, O.IsSome(result))
		extracted := O.MonadFold(result, func() *C.DurationFlag { return nil }, func(f *C.DurationFlag) *C.DurationFlag { return f })
		assert.NotNil(t, extracted)
		assert.Equal(t, "timeout", extracted.Name)
		assert.Equal(t, duration, extracted.Value)
	})
}

func TestTimestampFlagPrism_Success(t *testing.T) {
	t.Run("extracts TimestampFlag from Flag", func(t *testing.T) {
		// Arrange
		prism := TimestampFlagPrism()
		var flag C.Flag = &C.TimestampFlag{Name: "created"}

		// Act
		result := prism.GetOption(flag)

		// Assert
		assert.True(t, O.IsSome(result))
		extracted := O.MonadFold(result, func() *C.TimestampFlag { return nil }, func(f *C.TimestampFlag) *C.TimestampFlag { return f })
		assert.NotNil(t, extracted)
		assert.Equal(t, "created", extracted.Name)
	})
}

func TestStringSliceFlagPrism_Success(t *testing.T) {
	t.Run("extracts StringSliceFlag from Flag", func(t *testing.T) {
		// Arrange
		prism := StringSliceFlagPrism()
		var flag C.Flag = &C.StringSliceFlag{Name: "tags"}

		// Act
		result := prism.GetOption(flag)

		// Assert
		assert.True(t, O.IsSome(result))
		extracted := O.MonadFold(result, func() *C.StringSliceFlag { return nil }, func(f *C.StringSliceFlag) *C.StringSliceFlag { return f })
		assert.NotNil(t, extracted)
		assert.Equal(t, "tags", extracted.Name)
	})
}

func TestIntSliceFlagPrism_Success(t *testing.T) {
	t.Run("extracts IntSliceFlag from Flag", func(t *testing.T) {
		// Arrange
		prism := IntSliceFlagPrism()
		var flag C.Flag = &C.IntSliceFlag{Name: "ports"}

		// Act
		result := prism.GetOption(flag)

		// Assert
		assert.True(t, O.IsSome(result))
		extracted := O.MonadFold(result, func() *C.IntSliceFlag { return nil }, func(f *C.IntSliceFlag) *C.IntSliceFlag { return f })
		assert.NotNil(t, extracted)
		assert.Equal(t, "ports", extracted.Name)
	})
}

func TestFloat64SliceFlagPrism_Success(t *testing.T) {
	t.Run("extracts Float64SliceFlag from Flag", func(t *testing.T) {
		// Arrange
		prism := Float64SliceFlagPrism()
		var flag C.Flag = &C.Float64SliceFlag{Name: "ratios"}

		// Act
		result := prism.GetOption(flag)

		// Assert
		assert.True(t, O.IsSome(result))
		extracted := O.MonadFold(result, func() *C.Float64SliceFlag { return nil }, func(f *C.Float64SliceFlag) *C.Float64SliceFlag { return f })
		assert.NotNil(t, extracted)
		assert.Equal(t, "ratios", extracted.Name)
	})
}

func TestUintFlagPrism_Success(t *testing.T) {
	t.Run("extracts UintFlag from Flag", func(t *testing.T) {
		// Arrange
		prism := UintFlagPrism()
		var flag C.Flag = &C.UintFlag{Name: "workers", Value: 4}

		// Act
		result := prism.GetOption(flag)

		// Assert
		assert.True(t, O.IsSome(result))
		extracted := O.MonadFold(result, func() *C.UintFlag { return nil }, func(f *C.UintFlag) *C.UintFlag { return f })
		assert.NotNil(t, extracted)
		assert.Equal(t, "workers", extracted.Name)
		assert.Equal(t, uint(4), extracted.Value)
	})
}

func TestUint64FlagPrism_Success(t *testing.T) {
	t.Run("extracts Uint64Flag from Flag", func(t *testing.T) {
		// Arrange
		prism := Uint64FlagPrism()
		var flag C.Flag = &C.Uint64Flag{Name: "size", Value: 1024}

		// Act
		result := prism.GetOption(flag)

		// Assert
		assert.True(t, O.IsSome(result))
		extracted := O.MonadFold(result, func() *C.Uint64Flag { return nil }, func(f *C.Uint64Flag) *C.Uint64Flag { return f })
		assert.NotNil(t, extracted)
		assert.Equal(t, "size", extracted.Name)
		assert.Equal(t, uint64(1024), extracted.Value)
	})
}

func TestInt64FlagPrism_Success(t *testing.T) {
	t.Run("extracts Int64Flag from Flag", func(t *testing.T) {
		// Arrange
		prism := Int64FlagPrism()
		var flag C.Flag = &C.Int64Flag{Name: "offset", Value: -100}

		// Act
		result := prism.GetOption(flag)

		// Assert
		assert.True(t, O.IsSome(result))
		extracted := O.MonadFold(result, func() *C.Int64Flag { return nil }, func(f *C.Int64Flag) *C.Int64Flag { return f })
		assert.NotNil(t, extracted)
		assert.Equal(t, "offset", extracted.Name)
		assert.Equal(t, int64(-100), extracted.Value)
	})
}

func TestPrisms_EdgeCases(t *testing.T) {
	t.Run("all prisms return None for wrong type", func(t *testing.T) {
		// Arrange
		var flag C.Flag = &C.StringFlag{Name: "test"}

		// Act & Assert
		assert.True(t, O.IsNone(IntFlagPrism().GetOption(flag)))
		assert.True(t, O.IsNone(BoolFlagPrism().GetOption(flag)))
		assert.True(t, O.IsNone(Float64FlagPrism().GetOption(flag)))
		assert.True(t, O.IsNone(DurationFlagPrism().GetOption(flag)))
		assert.True(t, O.IsNone(TimestampFlagPrism().GetOption(flag)))
		assert.True(t, O.IsNone(StringSliceFlagPrism().GetOption(flag)))
		assert.True(t, O.IsNone(IntSliceFlagPrism().GetOption(flag)))
		assert.True(t, O.IsNone(Float64SliceFlagPrism().GetOption(flag)))
		assert.True(t, O.IsNone(UintFlagPrism().GetOption(flag)))
		assert.True(t, O.IsNone(Uint64FlagPrism().GetOption(flag)))
		assert.True(t, O.IsNone(Int64FlagPrism().GetOption(flag)))
	})
}
