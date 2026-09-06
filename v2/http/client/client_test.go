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

package client

import (
	"net/http"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/eq"
	LT "github.com/IBM/fp-go/v2/optics/lens/testing"
	"github.com/stretchr/testify/assert"
)

// eqDuration is an Eq instance for time.Duration using strict equality.
var eqDuration = eq.FromStrictEquals[time.Duration]()

// eqClientPtr is an Eq instance for *http.Client that compares pointer identity.
// The lens laws require structural equality of the source type. Since
// MakeLensRefWithName copies the struct on every Set, two independently
// produced *http.Client values with identical fields are different pointers;
// we therefore compare only Timeout and Transport to check structural equality.
var eqClientPtr = eq.FromEquals(func(a, b *http.Client) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Timeout == b.Timeout && a.Transport == b.Transport
})

// eqRoundTripper is an Eq instance for http.RoundTripper using pointer identity.
var eqRoundTripper = eq.FromEquals(func(a, b http.RoundTripper) bool {
	return a == b
})

// ---------------------------------------------------------------------------
// TimeoutLens
// ---------------------------------------------------------------------------

func TestTimeoutLens_Get(t *testing.T) {
	t.Run("returns zero duration for default client", func(t *testing.T) {
		lens := TimeoutLens()
		c := &http.Client{}
		assert.Equal(t, time.Duration(0), lens.Get(c))
	})

	t.Run("returns configured timeout", func(t *testing.T) {
		lens := TimeoutLens()
		c := &http.Client{Timeout: 15 * time.Second}
		assert.Equal(t, 15*time.Second, lens.Get(c))
	})
}

func TestTimeoutLens_Set(t *testing.T) {
	t.Run("sets timeout on a new copy", func(t *testing.T) {
		lens := TimeoutLens()
		original := &http.Client{}
		updated := lens.Set(30 * time.Second)(original)

		// The updated client has the new timeout.
		assert.Equal(t, 30*time.Second, updated.Timeout)
		// The original is not modified.
		assert.Equal(t, time.Duration(0), original.Timeout)
	})

	t.Run("overwrites existing timeout", func(t *testing.T) {
		lens := TimeoutLens()
		c := &http.Client{Timeout: 5 * time.Second}
		updated := lens.Set(60 * time.Second)(c)

		assert.Equal(t, 60*time.Second, updated.Timeout)
		assert.Equal(t, 5*time.Second, c.Timeout)
	})

	t.Run("immutability: original pointer is not reused", func(t *testing.T) {
		lens := TimeoutLens()
		original := &http.Client{}
		updated := lens.Set(10 * time.Second)(original)
		assert.NotSame(t, original, updated)
	})
}

func TestTimeoutLens_Laws(t *testing.T) {
	timeoutLaws := LT.AssertLaws(t, eqDuration, eqClientPtr)(TimeoutLens())

	c0 := &http.Client{}
	c1 := &http.Client{Timeout: 5 * time.Second}
	c2 := &http.Client{Timeout: 1 * time.Minute}

	for _, s := range []*http.Client{c0, c1, c2} {
		for _, a := range []time.Duration{0, time.Second, 30 * time.Second} {
			assert.True(t, timeoutLaws(s, a))
		}
	}
}

// ---------------------------------------------------------------------------
// TransportLens
// ---------------------------------------------------------------------------

func TestTransportLens_Get(t *testing.T) {
	t.Run("returns nil for default client", func(t *testing.T) {
		lens := TransportLens()
		c := &http.Client{}
		assert.Nil(t, lens.Get(c))
	})

	t.Run("returns configured transport", func(t *testing.T) {
		lens := TransportLens()
		transport := http.DefaultTransport
		c := &http.Client{Transport: transport}
		assert.Same(t, transport.(*http.Transport), lens.Get(c).(*http.Transport))
	})
}

func TestTransportLens_Set(t *testing.T) {
	t.Run("sets transport on a new copy", func(t *testing.T) {
		lens := TransportLens()
		transport := http.DefaultTransport
		original := &http.Client{}
		updated := lens.Set(transport)(original)

		assert.Equal(t, transport, updated.Transport)
		assert.Nil(t, original.Transport)
	})

	t.Run("overwrites existing transport", func(t *testing.T) {
		lens := TransportLens()
		t1 := http.DefaultTransport
		t2 := &http.Transport{MaxIdleConns: 5}

		c := &http.Client{Transport: t1}
		updated := lens.Set(t2)(c)

		assert.Equal(t, t2, updated.Transport)
		assert.Equal(t, t1, c.Transport)
	})

	t.Run("immutability: original pointer is not reused", func(t *testing.T) {
		lens := TransportLens()
		original := &http.Client{}
		updated := lens.Set(http.DefaultTransport)(original)
		assert.NotSame(t, original, updated)
	})
}

func TestTransportLens_Laws(t *testing.T) {
	transportLaws := LT.AssertLaws(t, eqRoundTripper, eqClientPtr)(TransportLens())

	c0 := &http.Client{}
	c1 := &http.Client{Transport: http.DefaultTransport}

	t1 := http.DefaultTransport
	var t2 http.RoundTripper // nil

	for _, s := range []*http.Client{c0, c1} {
		for _, a := range []http.RoundTripper{t1, t2} {
			assert.True(t, transportLaws(s, a))
		}
	}
}

// ---------------------------------------------------------------------------
// Combined lens composition
// ---------------------------------------------------------------------------

func TestComposedLenses(t *testing.T) {
	t.Run("timeout and transport can be set independently on the same client", func(t *testing.T) {
		tl := TimeoutLens()
		trl := TransportLens()

		base := &http.Client{}
		withTimeout := tl.Set(30 * time.Second)(base)
		withBoth := trl.Set(http.DefaultTransport)(withTimeout)

		assert.Equal(t, 30*time.Second, withBoth.Timeout)
		assert.Equal(t, http.DefaultTransport, withBoth.Transport)

		// originals untouched
		assert.Equal(t, time.Duration(0), base.Timeout)
		assert.Nil(t, base.Transport)
	})
}
