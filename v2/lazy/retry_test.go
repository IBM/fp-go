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

package lazy

import (
	"fmt"
	"strings"
	"testing"
	"time"

	R "github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

var expLogBackoff = R.ExponentialBackoff(10)

// our retry policy with a 1s cap
var testLogPolicy = R.CapDelay(
	2*time.Second,
	R.Monoid.Concat(expLogBackoff, R.LimitRetries(20)),
)

func TestRetry(t *testing.T) {
	action := func(status R.RetryStatus) Lazy[string] {
		return Of(fmt.Sprintf("Retrying %d", status.IterNumber))
	}
	check := func(value string) bool {
		return !strings.Contains(value, "5")
	}

	r := Retrying(testLogPolicy, action, check)

	assert.Equal(t, "Retrying 5", r())
}
