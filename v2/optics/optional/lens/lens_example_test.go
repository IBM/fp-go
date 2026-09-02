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

package lens_test

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/internal/common"
	optlens "github.com/IBM/fp-go/v2/optics/optional/lens"
	O "github.com/IBM/fp-go/v2/option"
)

// exConfig is a configuration struct used by the examples.
type exConfig struct {
	Timeout int
	Retries int
}

// exWrapper holds an optional *exConfig, representing a configuration that may
// or may not have been initialised.
type exWrapper struct {
	cfg *exConfig
}

// makeExConfigOptional returns an Optional[exWrapper, *exConfig] that focuses
// on the cfg pointer, returning None when it is nil.
func makeExConfigOptional() C.Optional[exWrapper, *exConfig] {
	return C.MakeOptional(
		func(w exWrapper) O.Option[*exConfig] {
			if w.cfg != nil {
				return O.Some(w.cfg)
			}
			return O.None[*exConfig]()
		},
		func(w exWrapper, c *exConfig) exWrapper {
			w.cfg = c
			return w
		},
	)
}

// makeTimeoutLens returns a Lens[*exConfig, int] that focuses on the Timeout
// field, making a copy on every Set so immutability is preserved.
func makeTimeoutLens() C.Lens[*exConfig, int] {
	return C.MakeLens(
		func(c *exConfig) int { return c.Timeout },
		func(c *exConfig, t int) *exConfig { c2 := *c; c2.Timeout = t; return &c2 },
	)
}

// ExampleCompose demonstrates the basic usage of Compose: narrowing an
// Optional[S, A] with a Lens[A, B] to obtain an Optional[S, B].
//
// Here:
//   - outer optional  exWrapper → *exConfig  (None when cfg is nil)
//   - inner lens      *exConfig → int         (Timeout field, always present)
//
// The result is an Optional[exWrapper, int] that focuses on Timeout but
// returns None (and is a no-op for Set) whenever the wrapper holds no config.
func ExampleCompose() {
	sa := makeExConfigOptional()      // Optional[exWrapper, *exConfig]
	ab := makeTimeoutLens()           // Lens[*exConfig, int]
	sb := F.Pipe1(sa, optlens.Compose[exWrapper](ab)) // Optional[exWrapper, int]

	present := exWrapper{cfg: &exConfig{Timeout: 30, Retries: 3}}
	absent := exWrapper{cfg: nil}

	// GetOption – Some when cfg is present
	fmt.Println(sb.GetOption(present))

	// GetOption – None when cfg is absent
	fmt.Println(sb.GetOption(absent))

	// Set updates Timeout; other fields are preserved
	updated := sb.Set(60)(present)
	fmt.Println(updated.cfg.Timeout)
	fmt.Println(updated.cfg.Retries)

	// Output:
	// Some[int](30)
	// None[int]
	// 60
	// 3
}

// ExampleCompose_noOpOnNone demonstrates that Set is a no-op when the outer
// optional produces None — consistent with the Optional no-op law.
func ExampleCompose_noOpOnNone() {
	sa := makeExConfigOptional()
	ab := makeTimeoutLens()
	sb := F.Pipe1(sa, optlens.Compose[exWrapper](ab))

	absent := exWrapper{cfg: nil}
	result := sb.Set(999)(absent)

	fmt.Println(result.cfg == nil)
	// Output:
	// true
}

// ExampleCompose_immutability demonstrates that Set does not mutate the
// original *exConfig value; a fresh copy is returned instead.
func ExampleCompose_immutability() {
	sa := makeExConfigOptional()
	ab := makeTimeoutLens()
	sb := F.Pipe1(sa, optlens.Compose[exWrapper](ab))

	original := exWrapper{cfg: &exConfig{Timeout: 10, Retries: 2}}
	updated := sb.Set(99)(original)

	// original is unchanged
	fmt.Println(original.cfg.Timeout)
	// updated holds the new value
	fmt.Println(updated.cfg.Timeout)
	// Output:
	// 10
	// 99
}
