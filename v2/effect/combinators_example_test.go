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
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/logging"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	R "github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
)

// exCombConfig is the dependency type of the combinator examples.
type exCombConfig struct {
	Factor int
}

func exCombFactor(c exCombConfig) int { return c.Factor }

// exCombScale multiplies a number by the configured factor.
func exCombScale(n int) Effect[exCombConfig, int] {
	return F.Pipe1(
		Asks(exCombFactor),
		Map[exCombConfig](N.Mul(n)),
	)
}

// exCombParse parses a number, failing for invalid input.
var exCombParse = F.Flow2(R.Eitherize1(strconv.Atoi), FromResult[exCombConfig, int])

// exCombRun provides the dependencies and runs an effect.
func exCombRun[A any](eff Effect[exCombConfig, A]) (A, error) {
	return RunSync(Provide[A](exCombConfig{Factor: 10})(eff))(context.Background())
}

// exCombLogger writes plain text records without time and correlation IDs to stdout.
func exCombLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey, "duration":
				return slog.Attr{}
			case "ID":
				return slog.String("ID", "<n>")
			}
			return a
		},
	}))
}

// exCombRunLogged runs an effect with exCombLogger as the context logger.
func exCombRunLogged[A any](eff Effect[exCombConfig, A]) (A, error) {
	ctx := pair.Second(logging.WithLogger(exCombLogger())(context.Background()))
	return RunSync(Provide[A](exCombConfig{Factor: 10})(eff))(ctx)
}

func ExampleBracket() {
	acquire := FromIO[exCombConfig](io.FromImpure(func() { fmt.Println("open") }))
	use := func(Void) Effect[exCombConfig, int] { return exCombScale(4) }
	release := func(_ Void, res Result[int]) Effect[exCombConfig, Void] {
		return FromIO[exCombConfig](io.FromImpure(func() { fmt.Println("close", res) }))
	}

	fmt.Println(exCombRun(Bracket(acquire, use, release)))
	// Output:
	// open
	// close Right[int](40)
	// 40 <nil>
}

func ExampleWithResource() {
	open := FromIO[exCombConfig](io.FromImpure(func() { fmt.Println("open") }))
	closeIt := func(Void) Effect[exCombConfig, Void] {
		return FromIO[exCombConfig](io.FromImpure(func() { fmt.Println("close") }))
	}
	withConnection := WithResource[int](open, closeIt)

	fmt.Println(exCombRun(withConnection(func(Void) Effect[exCombConfig, int] {
		return exCombParse("x")
	})))
	// Output:
	// open
	// close
	// 0 strconv.Atoi: parsing "x": invalid syntax
}

func ExampleWithLock() {
	var mu sync.Mutex
	lock := FromIO[exCombConfig](func() context.CancelFunc {
		mu.Lock()
		return mu.Unlock
	})

	fmt.Println(exCombRun(F.Pipe1(
		[]int{1, 2, 3},
		TraverseArrayPar(F.Flow2(exCombScale, WithLock[exCombConfig, int](lock))),
	)))
	// Output:
	// [10 20 30] <nil>
}

func ExampleWithSemaphore() {
	sem := make(chan Void, 2)

	fmt.Println(exCombRun(F.Pipe1(
		[]int{1, 2, 3, 4},
		TraverseArrayPar(F.Flow2(exCombScale, WithSemaphore[exCombConfig, int](sem))),
	)))
	// Output:
	// [10 20 30 40] <nil>
}

func ExampleTraverseArraySeq() {
	fmt.Println(exCombRun(F.Pipe1(
		[]string{"1", "2", "3"},
		TraverseArraySeq(exCombParse),
	)))
	// Output:
	// [1 2 3] <nil>
}

func ExampleTraverseArrayPar() {
	fmt.Println(exCombRun(F.Pipe1(
		[]int{1, 2, 3},
		TraverseArrayPar(exCombScale),
	)))
	// Output:
	// [10 20 30] <nil>
}

func ExampleTraverseArrayParN() {
	fmt.Println(exCombRun(F.Pipe1(
		[]int{1, 2, 3, 4, 5},
		TraverseArrayParN(2, exCombScale),
	)))
	// Output:
	// [10 20 30 40 50] <nil>
}

func ExampleTraverseArrayParN_failure() {
	fmt.Println(exCombRun(F.Pipe1(
		[]string{"1", "x", "3"},
		TraverseArrayParN(2, exCombParse),
	)))
	// Output:
	// [] strconv.Atoi: parsing "x": invalid syntax
}

func ExampleApSeq() {
	fmt.Println(exCombRun(F.Pipe1(
		Of[exCombConfig](N.Add(1)),
		ApSeq[int](exCombScale(4)),
	)))
	// Output:
	// 41 <nil>
}

func ExampleApPar() {
	fmt.Println(exCombRun(F.Pipe1(
		Of[exCombConfig](S.Format[int]("value=%d")),
		ApPar[string](exCombScale(4)),
	)))
	// Output:
	// value=40 <nil>
}

func ExampleDelay() {
	fmt.Println(exCombRun(F.Pipe1(
		exCombScale(2),
		Delay[exCombConfig, int](time.Millisecond),
	)))
	// Output:
	// 20 <nil>
}

func ExampleTimer() {
	start := time.Now()
	t, err := exCombRun(Timer[exCombConfig](5 * time.Millisecond))

	fmt.Println(!t.Before(start.Add(5*time.Millisecond)), err)
	// Output:
	// true <nil>
}

func ExampleWithTimeout() {
	slow := F.Pipe1(
		exCombScale(1),
		Delay[exCombConfig, int](time.Second),
	)

	fmt.Println(exCombRun(F.Pipe1(
		slow,
		WithTimeout[exCombConfig, int](10*time.Millisecond),
	)))
	// Output:
	// 0 context deadline exceeded
}

func ExampleTapSLog() {
	fmt.Println(exCombRunLogged(F.Pipe1(
		exCombScale(3),
		TapSLog[exCombConfig, int]("scaled"),
	)))
	// Output:
	// level=INFO msg=scaled value=30
	// 30 <nil>
}

func ExampleTapSLogDebug() {
	fmt.Println(exCombRunLogged(F.Pipe1(
		exCombParse("x"),
		TapSLogDebug[exCombConfig, int]("parsed"),
	)))
	// Output:
	// level=DEBUG msg=parsed error="strconv.Atoi: parsing \"x\": invalid syntax"
	// 0 strconv.Atoi: parsing "x": invalid syntax
}

func ExampleLogEntryExit() {
	fmt.Println(exCombRunLogged(F.Pipe1(
		exCombScale(5),
		LogEntryExit[exCombConfig, int]("scale"),
	)))
	// Output:
	// level=INFO msg=[entering] ID=<n> name=scale
	// level=INFO msg="[exiting ]" ID=<n> name=scale
	// 50 <nil>
}

func ExampleAnnotateLogs() {
	fmt.Println(exCombRunLogged(F.Pipe2(
		exCombScale(3),
		TapSLog[exCombConfig, int]("scaled"),
		AnnotateLogs[exCombConfig, int]("requestID", "r-1"),
	)))
	// Output:
	// level=INFO msg=scaled requestID=r-1 value=30
	// 30 <nil>
}
