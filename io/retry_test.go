package io

import (
	"fmt"
	"strings"
	"testing"
	"time"

	R "github.com/ibm/fp-go/retry"
	"github.com/stretchr/testify/assert"
)

var expLogBackoff = R.ExponentialBackoff(10)

// our retry policy with a 1s cap
var testLogPolicy = R.CapDelay(
	2*time.Second,
	R.Monoid.Concat(expLogBackoff, R.LimitRetries(20)),
)

func TestRetry(t *testing.T) {
	action := func(status R.RetryStatus) IO[string] {
		return Of(fmt.Sprintf("Retrying %d", status.IterNumber))
	}
	check := func(value string) bool {
		return !strings.Contains(value, "5")
	}

	r := Retrying(testLogPolicy, action, check)

	assert.Equal(t, "Retrying 5", r())
}
