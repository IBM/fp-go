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

package prism

import (
	"fmt"
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/tailrec"
	"github.com/stretchr/testify/assert"
)

type tramp = tailrec.Trampoline[int, string]

func TestFromLand_Success(t *testing.T) {
	p := FromLand[int, string]()

	t.Run("extracts landed result", func(t *testing.T) {
		assert.Equal(t, O.Some("done"), p.GetOption(tailrec.Land[int]("done")))
	})

	t.Run("reverse get creates landed trampoline", func(t *testing.T) {
		assert.Equal(t, tailrec.Land[int]("done"), p.ReverseGet("done"))
	})

	t.Run("set replaces landed result", func(t *testing.T) {
		assert.Equal(t, tailrec.Land[int]("new"), Set[tramp]("new")(p)(tailrec.Land[int]("old")))
	})
}

func TestFromLand_Failure(t *testing.T) {
	p := FromLand[int, string]()

	t.Run("bounce yields None", func(t *testing.T) {
		assert.Equal(t, O.None[string](), p.GetOption(tailrec.Bounce[string](42)))
	})

	t.Run("set leaves bounce unchanged", func(t *testing.T) {
		b := tailrec.Bounce[string](42)
		assert.Equal(t, b, Set[tramp]("new")(p)(b))
	})
}

func TestFromBounce_Success(t *testing.T) {
	p := FromBounce[int, string]()

	t.Run("extracts bounce state", func(t *testing.T) {
		assert.Equal(t, O.Some(42), p.GetOption(tailrec.Bounce[string](42)))
	})

	t.Run("reverse get creates bouncing trampoline", func(t *testing.T) {
		assert.Equal(t, tailrec.Bounce[string](42), p.ReverseGet(42))
	})

	t.Run("set replaces bounce state", func(t *testing.T) {
		assert.Equal(t, tailrec.Bounce[string](7), Set[tramp](7)(p)(tailrec.Bounce[string](42)))
	})
}

func TestFromBounce_Failure(t *testing.T) {
	p := FromBounce[int, string]()

	t.Run("land yields None", func(t *testing.T) {
		assert.Equal(t, O.None[int](), p.GetOption(tailrec.Land[int]("done")))
	})

	t.Run("set leaves land unchanged", func(t *testing.T) {
		l := tailrec.Land[int]("done")
		assert.Equal(t, l, Set[tramp](7)(p)(l))
	})
}

func TestTrampolinePrisms_EdgeCases(t *testing.T) {
	t.Run("landed zero value is still matched", func(t *testing.T) {
		assert.Equal(t, O.Some(0), FromLand[string, int]().GetOption(tailrec.Land[string](0)))
	})

	t.Run("zero value trampoline is a bounce", func(t *testing.T) {
		var zero tramp
		assert.Equal(t, O.Some(0), FromBounce[int, string]().GetOption(zero))
		assert.Equal(t, O.None[string](), FromLand[int, string]().GetOption(zero))
	})

	t.Run("prism laws", func(t *testing.T) {
		land := FromLand[int, string]()
		assert.Equal(t, O.Some("x"), land.GetOption(land.ReverseGet("x")))
		l := tailrec.Land[int]("x")
		assert.Equal(t, l, land.ReverseGet(O.GetOrElse(func() string { return "" })(land.GetOption(l))))

		bounce := FromBounce[int, string]()
		assert.Equal(t, O.Some(7), bounce.GetOption(bounce.ReverseGet(7)))
		b := tailrec.Bounce[string](7)
		assert.Equal(t, b, bounce.ReverseGet(O.GetOrElse(func() int { return 0 })(bounce.GetOption(b))))
	})
}

func TestTrampolinePrisms_Integration(t *testing.T) {
	t.Run("compose with FromOption", func(t *testing.T) {
		p := Compose[tailrec.Trampoline[int, Option[string]]](FromOption[string]())(FromLand[int, Option[string]]())
		assert.Equal(t, O.Some("x"), p.GetOption(tailrec.Land[int](O.Some("x"))))
		assert.Equal(t, O.None[string](), p.GetOption(tailrec.Land[int](O.None[string]())))
		assert.Equal(t, O.None[string](), p.GetOption(tailrec.Bounce[Option[string]](1)))
	})
}

// ExampleFromLand demonstrates focusing on the result of a landed trampoline.
func ExampleFromLand() {
	p := FromLand[int, string]()

	fmt.Println(p.GetOption(tailrec.Land[int]("done")))
	fmt.Println(p.GetOption(tailrec.Bounce[string](42)))
	fmt.Println(p.ReverseGet("done"))

	// Output:
	// Some[string](done)
	// None[string]
	// Land(done)
}

// ExampleFromBounce demonstrates focusing on the state of a bouncing trampoline.
func ExampleFromBounce() {
	p := FromBounce[int, string]()

	fmt.Println(p.GetOption(tailrec.Bounce[string](42)))
	fmt.Println(p.GetOption(tailrec.Land[int]("done")))
	fmt.Println(p.ReverseGet(42))

	// Output:
	// Some[int](42)
	// None[int]
	// Bounce(42)
}
