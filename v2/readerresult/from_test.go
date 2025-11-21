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

	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestFrom0(t *testing.T) {
	getConfig := func(ctx context.Context) (string, error) {
		return "config", nil
	}

	rr := From0(getConfig)()
	res := rr(context.Background())
	assert.Equal(t, result.Of("config"), res)

	// Test with error
	getConfigErr := func(ctx context.Context) (string, error) {
		return "", errors.New("config error")
	}

	rrErr := From0(getConfigErr)()
	resErr := rrErr(context.Background())
	assert.True(t, result.IsLeft(resErr))
}

func TestFrom1(t *testing.T) {
	getUser := func(ctx context.Context, id int) (string, error) {
		if id == 42 {
			return "Alice", nil
		}
		return "", errors.New("user not found")
	}

	rr := From1(getUser)

	res1 := rr(42)(context.Background())
	assert.Equal(t, result.Of("Alice"), res1)

	res2 := rr(99)(context.Background())
	assert.True(t, result.IsLeft(res2))
}

func TestFrom2(t *testing.T) {
	queryDB := func(ctx context.Context, table string, id int) (string, error) {
		if table == "users" && id == 42 {
			return "record", nil
		}
		return "", errors.New("not found")
	}

	rr := From2(queryDB)

	res1 := rr("users", 42)(context.Background())
	assert.Equal(t, result.Of("record"), res1)

	res2 := rr("posts", 1)(context.Background())
	assert.True(t, result.IsLeft(res2))
}

func TestFrom3(t *testing.T) {
	updateRecord := func(ctx context.Context, table string, id int, data string) (string, error) {
		if table == "users" && id == 42 && data == "updated" {
			return "success", nil
		}
		return "", errors.New("update failed")
	}

	rr := From3(updateRecord)

	res1 := rr("users", 42, "updated")(context.Background())
	assert.Equal(t, result.Of("success"), res1)

	res2 := rr("posts", 1, "data")(context.Background())
	assert.True(t, result.IsLeft(res2))
}
