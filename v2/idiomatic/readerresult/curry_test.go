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

package readerresult

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurry0(t *testing.T) {
	getConfig := func(ctx context.Context) (string, error) {
		return "config", nil
	}

	rr := Curry0(getConfig)
	v, err := rr(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "config", v)
}

func TestCurry1(t *testing.T) {
	getUser := func(ctx context.Context, id int) (string, error) {
		if id == 42 {
			return "Alice", nil
		}
		return "", errors.New("user not found")
	}

	curried := Curry1(getUser)

	rr1 := curried(42)
	v, err := rr1(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "Alice", v)

	rr2 := curried(99)
	_, err = rr2(context.Background())
	assert.Error(t, err)
}

func TestCurry2(t *testing.T) {
	queryDB := func(ctx context.Context, table string, id int) (string, error) {
		if table == "users" && id == 42 {
			return "record", nil
		}
		return "", errors.New("not found")
	}

	curried := Curry2(queryDB)

	rr1 := curried("users")(42)
	v, err := rr1(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "record", v)

	rr2 := curried("posts")(1)
	_, err = rr2(context.Background())
	assert.Error(t, err)
}

func TestCurry3(t *testing.T) {
	updateRecord := func(ctx context.Context, table string, id int, data string) (string, error) {
		if table == "users" && id == 42 && data == "updated" {
			return "success", nil
		}
		return "", errors.New("update failed")
	}

	curried := Curry3(updateRecord)

	rr1 := curried("users")(42)("updated")
	v, err := rr1(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "success", v)

	rr2 := curried("posts")(1)("data")
	_, err = rr2(context.Background())
	assert.Error(t, err)
}

func TestUncurry1(t *testing.T) {
	rrf := func(id int) ReaderResult[context.Context, string] {
		if id == 42 {
			return Of[context.Context]("Alice")
		}
		return Left[context.Context, string](errors.New("user not found"))
	}

	gofunc := Uncurry1(rrf)

	res1, err1 := gofunc(context.Background(), 42)
	assert.NoError(t, err1)
	assert.Equal(t, "Alice", res1)

	res2, err2 := gofunc(context.Background(), 99)
	assert.Error(t, err2)
	assert.Equal(t, "", res2)
}

func TestUncurry2(t *testing.T) {
	rrf := func(table string) func(int) ReaderResult[context.Context, string] {
		return func(id int) ReaderResult[context.Context, string] {
			if table == "users" && id == 42 {
				return Of[context.Context]("record")
			}
			return Left[context.Context, string](errors.New("not found"))
		}
	}

	gofunc := Uncurry2(rrf)

	res1, err1 := gofunc(context.Background(), "users", 42)
	assert.NoError(t, err1)
	assert.Equal(t, "record", res1)

	res2, err2 := gofunc(context.Background(), "posts", 1)
	assert.Error(t, err2)
	assert.Equal(t, "", res2)
}

func TestUncurry3(t *testing.T) {
	rrf := func(table string) func(int) func(string) ReaderResult[context.Context, string] {
		return func(id int) func(string) ReaderResult[context.Context, string] {
			return func(data string) ReaderResult[context.Context, string] {
				if table == "users" && id == 42 && data == "updated" {
					return Of[context.Context]("success")
				}
				return Left[context.Context, string](errors.New("update failed"))
			}
		}
	}

	gofunc := Uncurry3(rrf)

	res1, err1 := gofunc(context.Background(), "users", 42, "updated")
	assert.NoError(t, err1)
	assert.Equal(t, "success", res1)

	res2, err2 := gofunc(context.Background(), "posts", 1, "data")
	assert.Error(t, err2)
	assert.Equal(t, "", res2)
}

// Test round-trip conversions
func TestCurryUncurryRoundTrip(t *testing.T) {
	// Original Go function
	original := func(ctx context.Context, id int) (string, error) {
		if id == 42 {
			return "Alice", nil
		}
		return "", errors.New("not found")
	}

	// Curry then uncurry
	curried := Curry1(original)
	uncurried := Uncurry1(curried)

	// Should behave the same as original
	res1, err1 := original(context.Background(), 42)
	res2, err2 := uncurried(context.Background(), 42)
	assert.Equal(t, res1, res2)
	assert.Equal(t, err1, err2)

	res3, err3 := original(context.Background(), 99)
	res4, err4 := uncurried(context.Background(), 99)
	assert.Equal(t, res3, res4)
	assert.Equal(t, err3, err4)
}
