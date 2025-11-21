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

func TestFrom0(t *testing.T) {
	getConfig := func(ctx context.Context) (string, error) {
		return "config", nil
	}

	rr := From0(getConfig)()
	v, err := rr(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "config", v)

	// Test with error
	getConfigErr := func(ctx context.Context) (string, error) {
		return "", errors.New("config error")
	}

	rrErr := From0(getConfigErr)()
	_, err = rrErr(context.Background())
	assert.Error(t, err)
}

func TestFrom1(t *testing.T) {
	getUser := func(ctx context.Context, id int) (string, error) {
		if id == 42 {
			return "Alice", nil
		}
		return "", errors.New("user not found")
	}

	rr := From1(getUser)

	rr1 := rr(42)
	v, err := rr1(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "Alice", v)

	rr2 := rr(99)
	_, err = rr2(context.Background())
	assert.Error(t, err)
}

func TestFrom2(t *testing.T) {
	queryDB := func(ctx context.Context, table string, id int) (string, error) {
		if table == "users" && id == 42 {
			return "record", nil
		}
		return "", errors.New("not found")
	}

	rr := From2(queryDB)

	rr1 := rr("users", 42)
	v, err := rr1(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "record", v)

	rr2 := rr("posts", 1)
	_, err = rr2(context.Background())
	assert.Error(t, err)
}

func TestFrom3(t *testing.T) {
	updateRecord := func(ctx context.Context, table string, id int, data string) (string, error) {
		if table == "users" && id == 42 && data == "updated" {
			return "success", nil
		}
		return "", errors.New("update failed")
	}

	rr := From3(updateRecord)

	rr1 := rr("users", 42, "updated")
	v, err := rr1(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "success", v)

	rr2 := rr("posts", 1, "data")
	_, err = rr2(context.Background())
	assert.Error(t, err)
}
