package ioeither

import (
	"fmt"
	"testing"
	"time"

	E "github.com/ibm/fp-go/either"
	R "github.com/ibm/fp-go/retry"
	"github.com/stretchr/testify/assert"
)

var expLogBackoff = R.ExponentialBackoff(10 * time.Millisecond)

// our retry policy with a 1s cap
var testLogPolicy = R.CapDelay(
	2*time.Second,
	R.Monoid.Concat(expLogBackoff, R.LimitRetries(20)),
)

func TestRetry(t *testing.T) {
	action := func(status R.RetryStatus) IOEither[error, string] {
		if status.IterNumber < 5 {
			return Left[string](fmt.Errorf("retrying %d", status.IterNumber))
		}
		return Of[error](fmt.Sprintf("Retrying %d", status.IterNumber))
	}
	check := E.IsLeft[error, string]

	r := Retrying(testLogPolicy, action, check)

	assert.Equal(t, E.Of[error]("Retrying 5"), r())
}
