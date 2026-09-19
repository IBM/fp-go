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

package readerioresult

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	ER "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/logging"
)

// ---------------------------------------------------------------------------
// Shared fixtures for the logging examples
// ---------------------------------------------------------------------------

type exUser struct {
	ID   int
	Name string
}

var errExNotFound = errors.New("not found")

// exLogger writes plain text records to stdout so that the examples can assert
// on them. Time and duration are dropped and the correlation ID is masked
// because they differ between runs.
func exLogger() *slog.Logger {
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

// exRunLogged scopes exLogger to the computation and runs it. Every context-aware
// logging operator inside the computation picks the logger up from the context.
func exRunLogged[A any](rio ReaderIOResult[A]) (A, error) {
	return exRun(context.Background(), F.Pipe1(
		rio,
		Local[A](logging.WithLogger(exLogger())),
	))
}

// exFetchUser is an effectful lookup that fails for unknown IDs.
func exFetchUser(id int) ReaderIOResult[exUser] {
	if id == 42 {
		return Of(exUser{ID: 42, Name: "Alice"})
	}
	return Left[exUser](errExNotFound)
}

// exGetName is a named accessor, used to keep the pipelines point-free.
func exGetName(u exUser) string {
	return u.Name
}

// ---------------------------------------------------------------------------
// Simple log statements
// ---------------------------------------------------------------------------

// ExampleTapSLog adds a single structured log statement to a pipeline. TapSLog
// logs the value on success and the error on failure and passes the
// computation through unchanged.
func ExampleTapSLog() {
	pipeline := func(id int) ReaderIOResult[string] {
		return F.Pipe2(
			exFetchUser(id),
			TapSLog[exUser]("fetched user"),
			Map(exGetName),
		)
	}

	fmt.Println(exRunLogged(pipeline(42)))
	fmt.Println(exRunLogged(pipeline(7)))

	// Output:
	// level=INFO msg="fetched user" value="{ID:42 Name:Alice}"
	// Alice <nil>
	// level=INFO msg="fetched user" error="not found"
	//  not found
}

// ExampleTapIOK_logging adds a printf-style log statement. Any io.Kleisli can be
// used with TapIOK; io.Logf writes to the standard logger, io.Printf (used here
// so that the output is testable) writes to stdout.
func ExampleTapIOK_logging() {
	pipeline := F.Pipe2(
		exFetchUser(42),
		TapIOK(io.Printf[exUser]("fetched user %+v\n")),
		Map(exGetName),
	)

	fmt.Println(exRunLogged(pipeline))

	// Output:
	// fetched user {ID:42 Name:Alice}
	// Alice <nil>
}

// exLogUser is a custom log statement with explicit attributes. It reads the
// logger from the context, so it honors the scoped logger (and the correlation
// ID of an enclosing LogEntryExit).
func exLogUser(u exUser) ReaderIO[Void] {
	return func(ctx context.Context) IO[Void] {
		return func() Void {
			logging.GetLoggerFromContext(ctx).InfoContext(ctx, "user loaded", "id", u.ID, "name", u.Name)
			return F.VOID
		}
	}
}

// ExampleTapReaderIOK_logging uses a hand written ReaderIO to control exactly
// which attributes are logged.
func ExampleTapReaderIOK_logging() {
	pipeline := F.Pipe2(
		exFetchUser(42),
		TapReaderIOK(exLogUser),
		Map(exGetName),
	)

	fmt.Println(exRunLogged(pipeline))

	// Output:
	// level=INFO msg="user loaded" id=42 name=Alice
	// Alice <nil>
}

// ExampleSLogRight logs only the success channel.
func ExampleSLogRight() {
	pipeline := F.Pipe1(
		exFetchUser(42),
		Tap(SLogRight[exUser]("fetched user")),
	)

	fmt.Println(exRunLogged(pipeline))

	// Output:
	// level=INFO msg="fetched user" value="{ID:42 Name:Alice}"
	// {42 Alice} <nil>
}

// ExampleSLogLeft logs only the error channel, at ERROR level. The original
// error is propagated.
func ExampleSLogLeft() {
	pipeline := F.Pipe1(
		exFetchUser(7),
		TapLeft[exUser](SLogLeft("fetching user failed")),
	)

	_, err := exRunLogged(pipeline)
	fmt.Println(errors.Is(err, errExNotFound))

	// Output:
	// level=ERROR msg="fetching user failed" error="not found"
	// true
}

// ---------------------------------------------------------------------------
// Entry and exit logs
// ---------------------------------------------------------------------------

// ExampleLogEntryExit wraps a computation with entry and exit logs. The exit
// log is "[exiting ]" on success and "[throwing]" (with the error) on failure.
func ExampleLogEntryExit() {
	loadUser := func(id int) ReaderIOResult[exUser] {
		return F.Pipe1(
			exFetchUser(id),
			LogEntryExit[exUser]("loadUser"),
		)
	}

	fmt.Println(exRunLogged(loadUser(42)))
	fmt.Println(exRunLogged(loadUser(7)))

	// Output:
	// level=INFO msg=[entering] ID=<n> name=loadUser
	// level=INFO msg="[exiting ]" ID=<n> name=loadUser
	// {42 Alice} <nil>
	// level=INFO msg=[entering] ID=<n> name=loadUser
	// level=INFO msg=[throwing] ID=<n> name=loadUser error="not found"
	// {0 } not found
}

// ExampleLogEntryExit_nested shows that LogEntryExit wraps everything above it
// in the pipe. Log statements inside the wrapped computation run between the
// entry and the exit log and carry the same correlation ID.
func ExampleLogEntryExit_nested() {
	loadName := F.Pipe3(
		exFetchUser(42),
		TapSLog[exUser]("fetched user"), // inside: logged with ID
		Map(exGetName),
		LogEntryExit[string]("loadName"), // wraps the three steps above
	)

	fmt.Println(exRunLogged(loadName))

	// Output:
	// level=INFO msg=[entering] ID=<n> name=loadName
	// level=INFO msg="fetched user" ID=<n> value="{ID:42 Name:Alice}"
	// level=INFO msg="[exiting ]" ID=<n> name=loadName
	// Alice <nil>
}

// ExampleLogEntryExitWithCallback logs entry and exit at DEBUG level.
func ExampleLogEntryExitWithCallback() {
	pipeline := F.Pipe1(
		exFetchUser(42),
		LogEntryExitWithCallback[exUser](slog.LevelDebug, logging.GetLoggerFromContext, "loadUser"),
	)

	fmt.Println(exRunLogged(pipeline))

	// Output:
	// level=DEBUG msg=[entering] ID=<n> name=loadUser
	// level=DEBUG msg="[exiting ]" ID=<n> name=loadUser
	// {42 Alice} <nil>
}

// ---------------------------------------------------------------------------
// Adding context to the error channel
// ---------------------------------------------------------------------------

// ExampleMapLeft adds context to an error. errors.OnError wraps the original
// error with %w, so errors.Is and errors.As keep working downstream.
func ExampleMapLeft() {
	loadUser := func(id int) ReaderIOResult[exUser] {
		return F.Pipe1(
			exFetchUser(id),
			MapLeft[exUser](ER.OnError("loading user %d", id)),
		)
	}

	_, err := exRunLogged(loadUser(7))
	fmt.Println(err)
	fmt.Println(errors.Is(err, errExNotFound))

	// Output:
	// loading user 7, Caused By: not found
	// true
}

// exLoadGreeting is a step that enriches its own errors with the input that
// caused them, so that the caller does not need to know about it.
func exLoadGreeting(id int) ReaderIOResult[string] {
	return F.Pipe2(
		exFetchUser(id),
		Map(exGetName),
		MapLeft[string](ER.OnError("loading greeting for user %d", id)),
	)
}

// ExampleMapLeft_logging enriches the error close to where it occurs and logs
// it once, at the boundary of the pipeline. The single log record contains the
// full chain of context.
func ExampleMapLeft_logging() {
	pipeline := F.Pipe2(
		exLoadGreeting(7),
		MapLeft[string](ER.OnError("handling request")),
		TapLeft[string](SLogLeft("request failed")),
	)

	_, err := exRunLogged(pipeline)
	fmt.Println(errors.Is(err, errExNotFound))

	// Output:
	// level=ERROR msg="request failed" error="handling request, Caused By: loading greeting for user 7, Caused By: not found"
	// true
}
