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

package ioresult

import (
	"fmt"
	"testing"
	"time"

	R "github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

var expLogBackoff = R.ExponentialBackoff(10 * time.Millisecond)

// our retry policy with a 1s cap
var testLogPolicy = R.CapDelay(
	2*time.Second,
	R.Monoid.Concat(expLogBackoff, R.LimitRetries(20)),
)

func TestRetry(t *testing.T) {
	action := func(status R.RetryStatus) IOResult[string] {
		if status.IterNumber < 5 {
			return Left[string](fmt.Errorf("retrying %d", status.IterNumber))
		}
		return Of(fmt.Sprintf("Retrying %d", status.IterNumber))
	}
	check := func(val string, err error) bool {
		return err != nil
	}

	r := Retrying(testLogPolicy, action, check)

	result, err := r()
	assert.NoError(t, err)
	assert.Equal(t, "Retrying 5", result)
}
