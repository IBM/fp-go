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

package readerioresult

import (
	"context"
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/stretchr/testify/assert"
)

func TestEitherize0(t *testing.T) {
	f := func(ctx context.Context) (int, error) {
		return ctx.Value("key").(int), nil
	}

	result := Eitherize0(f)()
	ctx := context.WithValue(context.Background(), "key", 42)
	assert.Equal(t, E.Right[error](42), result(ctx)())
}

func TestEitherize1(t *testing.T) {
	f := func(ctx context.Context, x int) (int, error) {
		return x * 2, nil
	}

	result := Eitherize1(f)(5)
	ctx := context.Background()
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestEitherize2(t *testing.T) {
	f := func(ctx context.Context, x, y int) (int, error) {
		return x + y, nil
	}

	result := Eitherize2(f)(5, 3)
	ctx := context.Background()
	assert.Equal(t, E.Right[error](8), result(ctx)())
}

func TestEitherize3(t *testing.T) {
	f := func(ctx context.Context, x, y, z int) (int, error) {
		return x + y + z, nil
	}

	result := Eitherize3(f)(5, 3, 2)
	ctx := context.Background()
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestUneitherize0(t *testing.T) {
	f := func() ReaderIOResult[context.Context, int] {
		return Of[context.Context](42)
	}

	result := Uneitherize0(f)
	ctx := context.Background()
	val, err := result(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestUneitherize1(t *testing.T) {
	f := func(x int) ReaderIOResult[context.Context, int] {
		return Of[context.Context](x * 2)
	}

	result := Uneitherize1(f)
	ctx := context.Background()
	val, err := result(ctx, 5)
	assert.NoError(t, err)
	assert.Equal(t, 10, val)
}

func TestUneitherize2(t *testing.T) {
	f := func(x, y int) ReaderIOResult[context.Context, int] {
		return Of[context.Context](x + y)
	}

	result := Uneitherize2(f)
	ctx := context.Background()
	val, err := result(ctx, 5, 3)
	assert.NoError(t, err)
	assert.Equal(t, 8, val)
}

func TestFrom0(t *testing.T) {
	f := func(ctx context.Context) func() (int, error) {
		return func() (int, error) {
			return 42, nil
		}
	}

	result := From0(f)()
	ctx := context.Background()
	assert.Equal(t, E.Right[error](42), result(ctx)())
}

func TestFrom1(t *testing.T) {
	f := func(ctx context.Context, x int) func() (int, error) {
		return func() (int, error) {
			return x * 2, nil
		}
	}

	result := From1(f)(5)
	ctx := context.Background()
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestFrom2(t *testing.T) {
	f := func(ctx context.Context, x, y int) func() (int, error) {
		return func() (int, error) {
			return x + y, nil
		}
	}

	result := From2(f)(5, 3)
	ctx := context.Background()
	assert.Equal(t, E.Right[error](8), result(ctx)())
}

func TestEitherizeWithError(t *testing.T) {
	f := func(ctx context.Context, x int) (int, error) {
		if x < 0 {
			return 0, errors.New("negative value")
		}
		return x * 2, nil
	}

	result := Eitherize1(f)(-5)
	ctx := context.Background()
	assert.True(t, E.IsLeft(result(ctx)()))
}

func TestUneitherizeWithError(t *testing.T) {
	f := func(x int) ReaderIOResult[context.Context, int] {
		if x < 0 {
			return Left[context.Context, int](errors.New("negative value"))
		}
		return Of[context.Context](x * 2)
	}

	result := Uneitherize1(f)
	ctx := context.Background()
	_, err := result(ctx, -5)
	assert.Error(t, err)
}
