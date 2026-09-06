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

package effect

import (
	"fmt"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	R "github.com/IBM/fp-go/v2/result"
)

// ---------------------------------------------------------------------------
// Shared fixtures
//
// Do-notation assembles a record field by field.  exOrder is that record and
// the setters below are the "one field at a time" updates each combinator needs.
// ---------------------------------------------------------------------------

// exOrder is the state accumulated by the do-notation examples.
type exOrder struct {
	Customer string
	Quantity int
	Total    int
}

// Curried setters — the shape `func(T) func(S) S` that Bind/Let/ApS expect.
func exSetCustomer(v string) func(exOrder) exOrder {
	return func(o exOrder) exOrder { o.Customer = v; return o }
}

func exSetQuantity(v int) func(exOrder) exOrder {
	return func(o exOrder) exOrder { o.Quantity = v; return o }
}

func exSetTotal(v int) func(exOrder) exOrder {
	return func(o exOrder) exOrder { o.Total = v; return o }
}

// Lenses — the *L variants take a lens instead of a setter, which removes the
// hand-written closure entirely.
var (
	exCustomerLens = L.MakeLens(
		func(o exOrder) string { return o.Customer },
		func(o exOrder, v string) exOrder { o.Customer = v; return o },
	)
	exQuantityLens = L.MakeLens(
		func(o exOrder) int { return o.Quantity },
		func(o exOrder, v int) exOrder { o.Quantity = v; return o },
	)
)

// exOrderQuantity is a named getter, so Let can stay point-free.
func exOrderQuantity(o exOrder) int { return o.Quantity }

// ---------------------------------------------------------------------------
// Do / BindTo — starting a do block
// ---------------------------------------------------------------------------

// ExampleDo demonstrates starting a do-notation block from an empty state.
func ExampleDo() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApS(exSetCustomer, Of[exCoreConfig]("acme")),
		ApS(exSetQuantity, Of[exCoreConfig](3)),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// {acme 3 0} <nil>
}

// ExampleBindTo demonstrates the other way to start a block: take an existing
// effect and lift its value into the first field of the state.
func ExampleBindTo() {
	eff := F.Pipe2(
		Of[exCoreConfig]("acme"),
		BindTo[exCoreConfig](func(c string) exOrder { return exOrder{Customer: c} }),
		ApS(exSetQuantity, Of[exCoreConfig](3)),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// {acme 3 0} <nil>
}

// ---------------------------------------------------------------------------
// Bind / ApS / Let / LetTo
// ---------------------------------------------------------------------------

// ExampleBind demonstrates adding a field whose effect *depends* on the state
// accumulated so far — the sequential case.
func ExampleBind() {
	// The total needs the quantity that an earlier step put into the state.
	priceIt := func(o exOrder) Effect[exCoreConfig, int] {
		return F.Pipe1(Asks(exCoreMultiplier), Map[exCoreConfig](N.Mul(o.Quantity)))
	}

	eff := F.Pipe3(
		Do[exCoreConfig](exOrder{}),
		ApS(exSetCustomer, Of[exCoreConfig]("acme")),
		ApS(exSetQuantity, Of[exCoreConfig](3)),
		Bind(exSetTotal, priceIt),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 10}, eff))
	// Output:
	// {acme 3 30} <nil>
}

// ExampleBind_shortCircuits demonstrates that a failing step abandons the whole
// block.
func ExampleBind_shortCircuits() {
	reject := func(exOrder) Effect[exCoreConfig, int] {
		return Fail[exCoreConfig, int](fmt.Errorf("pricing unavailable"))
	}

	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApS(exSetQuantity, Of[exCoreConfig](3)),
		Bind(exSetTotal, reject),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// { 0 0} pricing unavailable
}

// ExampleApS demonstrates adding a field whose effect is *independent* of the
// accumulated state — use ApS for these and Bind only when there is a real
// dependency.
func ExampleApS() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApS(exSetCustomer, Of[exCoreConfig]("acme")),
		ApS(exSetQuantity, Asks(exCoreMultiplier)),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 7}, eff))
	// Output:
	// {acme 7 0} <nil>
}

// ExampleLet demonstrates deriving a field from the state with a pure function
// — no effect involved, so no failure is possible.
func ExampleLet() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApS(exSetQuantity, Of[exCoreConfig](4)),
		// Point-free: read the quantity, then multiply by the unit price.
		Let[exCoreConfig](exSetTotal, F.Flow2(exOrderQuantity, N.Mul(25))),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// { 4 100} <nil>
}

// ExampleLetTo demonstrates attaching a constant to the state.
func ExampleLetTo() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		LetTo[exCoreConfig](exSetCustomer, "walk-in"),
		LetTo[exCoreConfig](exSetQuantity, 1),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// {walk-in 1 0} <nil>
}

// ---------------------------------------------------------------------------
// Lens-based variants
// ---------------------------------------------------------------------------

// ExampleApSL demonstrates the lens form of ApS: the lens replaces the
// hand-written setter closure.
func ExampleApSL() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApSL(exCustomerLens, Of[exCoreConfig]("acme")),
		ApSL(exQuantityLens, Asks(exCoreMultiplier)),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 5}, eff))
	// Output:
	// {acme 5 0} <nil>
}

// ExampleBindL demonstrates the lens form of Bind.  The Kleisli arrow receives
// the *current field value* rather than the whole state, so it can refine a
// field in place.
func ExampleBindL() {
	// Scale the quantity already in the state by the configured multiplier.
	scale := func(q int) Effect[exCoreConfig, int] {
		return F.Pipe1(Asks(exCoreMultiplier), Map[exCoreConfig](N.Mul(q)))
	}

	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApSL(exQuantityLens, Of[exCoreConfig](3)),
		BindL(exQuantityLens, scale),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 4}, eff))
	// Output:
	// { 12 0} <nil>
}

// ExampleLetL demonstrates the lens form of Let: a pure, point-free update of
// one field.
func ExampleLetL() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApSL(exQuantityLens, Of[exCoreConfig](6)),
		LetL[exCoreConfig](exQuantityLens, N.Mul(2)),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// { 12 0} <nil>
}

// ExampleLetToL demonstrates the lens form of LetTo.
func ExampleLetToL() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		LetToL[exCoreConfig](exCustomerLens, "walk-in"),
		LetToL[exCoreConfig](exQuantityLens, 1),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// {walk-in 1 0} <nil>
}

// ---------------------------------------------------------------------------
// Bind*K — binding foreign Kleisli arrows
// ---------------------------------------------------------------------------

// ExampleBindEitherK demonstrates binding a pure computation that may fail.
func ExampleBindEitherK() {
	// exOrder -> Result[int]: reject an empty order.
	price := func(o exOrder) Result[int] {
		if o.Quantity <= 0 {
			return R.Left[int](fmt.Errorf("quantity must be positive"))
		}
		return R.Of(o.Quantity * 25)
	}

	ok := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApS(exSetQuantity, Of[exCoreConfig](4)),
		BindEitherK[exCoreConfig](exSetTotal, price),
	)
	fmt.Println(exCoreRun(exCoreConfig{}, ok))

	bad := F.Pipe1(Do[exCoreConfig](exOrder{}), BindEitherK[exCoreConfig](exSetTotal, price))
	fmt.Println(exCoreRun(exCoreConfig{}, bad))
	// Output:
	// { 4 100} <nil>
	// { 0 0} quantity must be positive
}

// ExampleBindIOK demonstrates binding a side effect that cannot fail.
func ExampleBindIOK() {
	// Point-free: read the quantity, multiply, lift into IO.
	price := F.Flow3(exOrderQuantity, N.Mul(25), io.Of[int])

	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApS(exSetQuantity, Of[exCoreConfig](4)),
		BindIOK[exCoreConfig](exSetTotal, price),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// { 4 100} <nil>
}

// ExampleBindIOResultK demonstrates binding a side effect that may fail.
func ExampleBindIOResultK() {
	price := func(o exOrder) ioresult.IOResult[int] {
		if o.Quantity <= 0 {
			return ioresult.Left[int](fmt.Errorf("quantity must be positive"))
		}
		return ioresult.Of(o.Quantity * 25)
	}

	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApS(exSetQuantity, Of[exCoreConfig](2)),
		BindIOResultK[exCoreConfig](exSetTotal, price),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// { 2 50} <nil>
}

// ExampleBindIOEitherK demonstrates the IOEither spelling of BindIOResultK, for
// code that works with an explicit error type parameter.
func ExampleBindIOEitherK() {
	price := F.Flow3(exOrderQuantity, N.Mul(25), ioeither.Of[error, int])

	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApS(exSetQuantity, Of[exCoreConfig](3)),
		BindIOEitherK[exCoreConfig](exSetTotal, price),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// { 3 75} <nil>
}

// ExampleBindReaderK demonstrates binding a pure computation that reads the
// dependency.
func ExampleBindReaderK() {
	// exOrder -> Reader[exCoreConfig, int]
	price := func(o exOrder) reader.Reader[exCoreConfig, int] {
		return F.Flow2(exCoreMultiplier, N.Mul(o.Quantity))
	}

	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApS(exSetQuantity, Of[exCoreConfig](3)),
		BindReaderK(exSetTotal, price),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 10}, eff))
	// Output:
	// { 3 30} <nil>
}

// ExampleBindReaderIOK demonstrates binding a computation that reads the
// dependency and performs IO.
func ExampleBindReaderIOK() {
	price := func(o exOrder) readerio.ReaderIO[exCoreConfig, int] {
		return func(cfg exCoreConfig) io.IO[int] {
			return io.Of(o.Quantity * cfg.Multiplier)
		}
	}

	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApS(exSetQuantity, Of[exCoreConfig](5)),
		BindReaderIOK(exSetTotal, price),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 10}, eff))
	// Output:
	// { 5 50} <nil>
}

// ---------------------------------------------------------------------------
// Bind*KL — lens forms of the Bind*K family
// ---------------------------------------------------------------------------

// ExampleBindIOKL demonstrates refining one field in place with an IO action.
func ExampleBindIOKL() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApSL(exQuantityLens, Of[exCoreConfig](6)),
		BindIOKL[exCoreConfig](exQuantityLens, F.Flow2(N.Mul(2), io.Of[int])),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// { 12 0} <nil>
}

// ExampleBindIOEitherKL demonstrates refining one field in place with a
// failure-capable IO action.
func ExampleBindIOEitherKL() {
	double := F.Flow2(N.Mul(2), ioeither.Of[error, int])

	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApSL(exQuantityLens, Of[exCoreConfig](6)),
		BindIOEitherKL[exCoreConfig](exQuantityLens, double),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// { 12 0} <nil>
}

// ExampleBindReaderKL demonstrates refining one field in place with a
// dependency-reading computation.
func ExampleBindReaderKL() {
	scale := func(q int) reader.Reader[exCoreConfig, int] {
		return F.Flow2(exCoreMultiplier, N.Mul(q))
	}

	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApSL(exQuantityLens, Of[exCoreConfig](3)),
		BindReaderKL(exQuantityLens, scale),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 4}, eff))
	// Output:
	// { 12 0} <nil>
}

// ExampleBindReaderIOKL demonstrates refining one field in place with a
// dependency-reading IO computation.
func ExampleBindReaderIOKL() {
	scale := func(q int) readerio.ReaderIO[exCoreConfig, int] {
		return func(cfg exCoreConfig) io.IO[int] { return io.Of(q * cfg.Multiplier) }
	}

	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApSL(exQuantityLens, Of[exCoreConfig](3)),
		BindReaderIOKL(exQuantityLens, scale),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 5}, eff))
	// Output:
	// { 15 0} <nil>
}

// ---------------------------------------------------------------------------
// Ap*S — binding foreign *values* (independent of the state)
// ---------------------------------------------------------------------------

// ExampleApEitherS demonstrates binding an already-computed Result.
func ExampleApEitherS() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApEitherS[exCoreConfig](exSetCustomer, E.Of[error]("acme")),
		ApEitherS[exCoreConfig](exSetQuantity, R.Of(3)),
	)
	fmt.Println(exCoreRun(exCoreConfig{}, eff))

	failing := F.Pipe1(
		Do[exCoreConfig](exOrder{}),
		ApEitherS[exCoreConfig](exSetQuantity, R.Left[int](fmt.Errorf("no stock"))),
	)
	fmt.Println(exCoreRun(exCoreConfig{}, failing))
	// Output:
	// {acme 3 0} <nil>
	// { 0 0} no stock
}

// ExampleApIOS demonstrates binding an IO value that cannot fail.
func ExampleApIOS() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApIOS[exCoreConfig](exSetCustomer, io.Of("acme")),
		ApIOS[exCoreConfig](exSetQuantity, io.Of(3)),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// {acme 3 0} <nil>
}

// ExampleApIOEitherS demonstrates binding an IO value that may fail.
func ExampleApIOEitherS() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApIOEitherS[exCoreConfig](exSetCustomer, ioeither.Of[error]("acme")),
		ApIOEitherS[exCoreConfig](exSetQuantity, ioeither.Of[error](3)),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// {acme 3 0} <nil>
}

// ExampleApReaderS demonstrates binding a value read from the dependency.
func ExampleApReaderS() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApReaderS(exSetCustomer, reader.Reader[exCoreConfig, string](exCorePrefix)),
		ApReaderS(exSetQuantity, reader.Reader[exCoreConfig, int](exCoreMultiplier)),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 3, Prefix: "acme"}, eff))
	// Output:
	// {acme 3 0} <nil>
}

// ExampleApReaderIOS demonstrates binding a value read from the dependency by
// an IO action.
func ExampleApReaderIOS() {
	quantity := func(cfg exCoreConfig) io.IO[int] { return io.Of(cfg.Multiplier) }

	eff := F.Pipe1(
		Do[exCoreConfig](exOrder{}),
		ApReaderIOS(exSetQuantity, readerio.ReaderIO[exCoreConfig, int](quantity)),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 8}, eff))
	// Output:
	// { 8 0} <nil>
}

// ---------------------------------------------------------------------------
// Ap*SL — lens forms of the Ap*S family
// ---------------------------------------------------------------------------

// ExampleApEitherSL demonstrates the lens form of ApEitherS.
func ExampleApEitherSL() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApEitherSL[exCoreConfig](exCustomerLens, R.Of("acme")),
		ApEitherSL[exCoreConfig](exQuantityLens, R.Of(3)),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// {acme 3 0} <nil>
}

// ExampleApIOSL demonstrates the lens form of ApIOS.
func ExampleApIOSL() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApIOSL[exCoreConfig](exCustomerLens, io.Of("acme")),
		ApIOSL[exCoreConfig](exQuantityLens, io.Of(3)),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// {acme 3 0} <nil>
}

// ExampleApIOEitherSL demonstrates the lens form of ApIOEitherS.
func ExampleApIOEitherSL() {
	eff := F.Pipe1(
		Do[exCoreConfig](exOrder{}),
		ApIOEitherSL[exCoreConfig](exQuantityLens, ioeither.Of[error](3)),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// { 3 0} <nil>
}

// ExampleApReaderSL demonstrates the lens form of ApReaderS.
func ExampleApReaderSL() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApReaderSL(exCustomerLens, reader.Reader[exCoreConfig, string](exCorePrefix)),
		ApReaderSL(exQuantityLens, reader.Reader[exCoreConfig, int](exCoreMultiplier)),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 3, Prefix: "acme"}, eff))
	// Output:
	// {acme 3 0} <nil>
}

// ExampleApReaderIOSL demonstrates the lens form of ApReaderIOS.
func ExampleApReaderIOSL() {
	quantity := func(cfg exCoreConfig) io.IO[int] { return io.Of(cfg.Multiplier) }

	eff := F.Pipe1(
		Do[exCoreConfig](exOrder{}),
		ApReaderIOSL(exQuantityLens, readerio.ReaderIO[exCoreConfig, int](quantity)),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 8}, eff))
	// Output:
	// { 8 0} <nil>
}

// ---------------------------------------------------------------------------
// Putting it together
// ---------------------------------------------------------------------------

// ExampleDo_fullPipeline demonstrates the rule of thumb for choosing a
// combinator: ApS for independent fields, Bind when a field needs what came
// before, Let for pure derivations.
func ExampleDo_fullPipeline() {
	priceIt := func(o exOrder) Effect[exCoreConfig, int] {
		return F.Pipe1(Asks(exCoreMultiplier), Map[exCoreConfig](N.Mul(o.Quantity)))
	}

	eff := F.Pipe4(
		Do[exCoreConfig](exOrder{}),
		ApSL(exCustomerLens, Asks(exCorePrefix)),  // independent
		ApSL(exQuantityLens, Of[exCoreConfig](3)), // independent
		Bind(exSetTotal, priceIt),                 // depends on Quantity
		Map[exCoreConfig](func(o exOrder) string {
			return fmt.Sprintf("%s x%d = %d", o.Customer, o.Quantity, o.Total)
		}),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 25, Prefix: "acme"}, eff))
	// Output:
	// acme x3 = 75 <nil>
}
