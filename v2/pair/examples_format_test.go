// Copyright (c) 2024 - 2025 IBM Corp.
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

package pair_test

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	P "github.com/IBM/fp-go/v2/pair"
)

// ExamplePair_String demonstrates the fmt.Stringer interface implementation.
func ExamplePair_String() {
	p1 := P.MakePair("username", 42)
	p2 := P.MakePair(100, "active")

	fmt.Println(p1.String())
	fmt.Println(p2.String())

	// Output:
	// Pair[string, int](username, 42)
	// Pair[int, string](100, active)
}

// ExamplePair_GoString demonstrates the fmt.GoStringer interface implementation.
func ExamplePair_GoString() {
	p1 := P.MakePair("key", 42)
	p2 := P.MakePair(errors.New("error"), "value")

	fmt.Printf("%#v\n", p1)
	fmt.Printf("%#v\n", p2)

	// Output:
	// pair.MakePair[string, int]("key", 42)
	// pair.MakePair[error, string](&errors.errorString{s:"error"}, "value")
}

// ExamplePair_Format demonstrates the fmt.Formatter interface implementation.
func ExamplePair_Format() {
	p := P.MakePair("config", 8080)

	// Different format verbs
	fmt.Printf("%%s: %s\n", p)
	fmt.Printf("%%v: %v\n", p)
	fmt.Printf("%%+v: %+v\n", p)
	fmt.Printf("%%#v: %#v\n", p)

	// Output:
	// %s: Pair[string, int](config, 8080)
	// %v: Pair[string, int](config, 8080)
	// %+v: Pair[string, int](config, 8080)
	// %#v: pair.MakePair[string, int]("config", 8080)
}

// ExamplePair_LogValue demonstrates the slog.LogValuer interface implementation.
func ExamplePair_LogValue() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time for consistent output
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	// Pair with string and int
	p1 := P.MakePair("username", 42)
	logger.Info("user data", "data", p1)

	// Pair with error and string
	p2 := P.MakePair(errors.New("connection failed"), "retry")
	logger.Error("operation failed", "status", p2)

	// Output:
	// level=INFO msg="user data" data.head=username data.tail=42
	// level=ERROR msg="operation failed" status.head="connection failed" status.tail=retry
}

// ExamplePair_formatting_comparison demonstrates different formatting options.
func ExamplePair_formatting_comparison() {
	type Config struct {
		Host string
		Port int
	}

	config := Config{Host: "localhost", Port: 8080}
	p := P.MakePair(config, []string{"api", "web"})

	fmt.Printf("String():   %s\n", p.String())
	fmt.Printf("GoString(): %s\n", p.GoString())
	fmt.Printf("%%v:         %v\n", p)
	fmt.Printf("%%#v:        %#v\n", p)

	// Output:
	// String():   Pair[pair_test.Config, []string]({localhost 8080}, [api web])
	// GoString(): pair.MakePair[pair_test.Config, []string](pair_test.Config{Host:"localhost", Port:8080}, []string{"api", "web"})
	// %v:         Pair[pair_test.Config, []string]({localhost 8080}, [api web])
	// %#v:        pair.MakePair[pair_test.Config, []string](pair_test.Config{Host:"localhost", Port:8080}, []string{"api", "web"})
}

// ExamplePair_LogValue_structured demonstrates structured logging with Pair.
func ExamplePair_LogValue_structured() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	// Simulate a key-value store operation
	operation := func(key string, value int) P.Pair[string, int] {
		return P.MakePair(key, value)
	}

	// Log successful operation
	result1 := operation("counter", 42)
	logger.Info("store operation", "key", "counter", "result", result1)

	// Log another operation
	result2 := operation("timeout", 30)
	logger.Info("store operation", "key", "timeout", "result", result2)

	// Output:
	// level=INFO msg="store operation" key=counter result.head=counter result.tail=42
	// level=INFO msg="store operation" key=timeout result.head=timeout result.tail=30
}

// ExamplePair_formatting_with_maps demonstrates formatting pairs containing maps.
func ExamplePair_formatting_with_maps() {
	metadata := map[string]string{
		"version": "1.0",
		"author":  "Alice",
	}
	p := P.MakePair("config", metadata)

	fmt.Printf("%%v: %v\n", p)
	fmt.Printf("%%s: %s\n", p)

	// Output:
	// %v: Pair[string, map[string]string](config, map[author:Alice version:1.0])
	// %s: Pair[string, map[string]string](config, map[author:Alice version:1.0])
}

// ExamplePair_formatting_nested demonstrates formatting nested pairs.
func ExamplePair_formatting_nested() {
	inner := P.MakePair("inner", 10)
	outer := P.MakePair(inner, "outer")

	fmt.Printf("%%v: %v\n", outer)
	fmt.Printf("%%#v: %#v\n", outer)

	// Output:
	// %v: Pair[pair.Pair[string,int], string](Pair[string, int](inner, 10), outer)
	// %#v: pair.MakePair[pair.Pair[string,int], string](pair.MakePair[string, int]("inner", 10), "outer")
}
